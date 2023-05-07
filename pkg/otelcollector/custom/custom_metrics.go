package config

import (
	"fmt"
	"sort"
	"strings"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/leaderonlyreceiver"
	"go.opentelemetry.io/collector/component"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/structpb"
)

// AddCustomMetricsPipelines adds custom metrics pipelines to the given OTelConfig.
func AddCustomMetricsPipelines(
	config *otelconfig.OTelConfig,
	customMetrics map[string]*policylangv1.InfraMeter,
) error {
	config.AddProcessor(otelconsts.ProcessorCustomMetrics, map[string]any{
		"attributes": []map[string]interface{}{
			{
				"key":    "service.name",
				"action": "upsert",
				"value":  "aperture-custom-metrics",
			},
		},
	})
	if customMetrics == nil {
		customMetrics = map[string]*policylangv1.InfraMeter{}
	}
	for pipelineName, metricConfig := range customMetrics {
		if err := addCustomMetricsPipeline(config, pipelineName, metricConfig); err != nil {
			return fmt.Errorf("failed to add custom metric pipeline %s: %w", pipelineName, err)
		}
	}
	return nil
}

func addCustomMetricsPipeline(
	config *otelconfig.OTelConfig,
	pipelineName string,
	metricConfig *policylangv1.InfraMeter,
) error {
	pipelineName = strings.TrimPrefix(pipelineName, "metrics/")

	receiverIDs := map[string]string{}
	processorIDs := map[string]string{}

	for origName, receiverConfig := range metricConfig.Receivers {
		var id component.ID
		if err := id.UnmarshalText([]byte(origName)); err != nil {
			return fmt.Errorf("invalid id %q: %w", origName, err)
		}
		id = component.NewIDWithName(id.Type(), normalizeComponentName(pipelineName, id.Name()))
		var cfg any
		cfg = receiverConfig.AsMap()
		id, cfg = leaderonlyreceiver.WrapConfigIf(metricConfig.PerAgentGroup, id, cfg)
		receiverIDs[origName] = id.String()
		config.AddReceiver(id.String(), cfg)
	}

	for origName, processorConfig := range metricConfig.Processors {
		id := normalizeComponentName(pipelineName, origName)
		processorIDs[origName] = id
		config.AddProcessor(id, processorConfig)
	}

	if len(metricConfig.Pipeline.Receivers) == 0 && len(metricConfig.Pipeline.Processors) == 0 {
		if len(metricConfig.Processors) >= 1 {
			return fmt.Errorf("empty pipeline, inferring pipeline is supported only with 0 or 1 processors")
		}

		// Skip adding pipeline if there are no receivers and processors.
		if len(metricConfig.Receivers) == 0 && len(metricConfig.Processors) == 0 {
			return nil
		}

		// When pipeline not set explicitly, create pipeline with all defined receivers and processors.
		if len(metricConfig.Receivers) > 0 {
			metricConfig.Pipeline.Receivers = maps.Keys(metricConfig.Receivers)
			sort.Strings(metricConfig.Pipeline.Receivers)
		}
		if len(metricConfig.Processors) > 0 {
			metricConfig.Pipeline.Processors = maps.Keys(metricConfig.Processors)
		}
	}

	config.Service.AddPipeline(normalizePipelineName(pipelineName), otelconfig.Pipeline{
		Receivers: mapSlice(receiverIDs, metricConfig.Pipeline.Receivers),
		Processors: append(
			mapSlice(processorIDs, metricConfig.Pipeline.Processors),
			otelconsts.ProcessorCustomMetrics,
			otelconsts.ProcessorAgentResourceLabels,
		),
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})

	return nil
}

// normalizePipelineName normalizes user defined pipeline name by adding
// `metrics/user-defined-` prefix.
// This ensures no builtin metrics pipeline is overwritten.
func normalizePipelineName(pipelineName string) string {
	return fmt.Sprintf("metrics/user-defined-%s", pipelineName)
}

func mapSlice(mapping map[string]string, xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		y, ok := mapping[x]
		if !ok {
			y = x
		}
		ys = append(ys, y)
	}
	return ys
}

// normalizeComponentName normalizes user defines component name by adding
// `user-defined-<pipeline_name>` suffix.
// This ensures no builtin components are overwritten.
func normalizeComponentName(pipelineName, componentName string) string {
	suffix := fmt.Sprintf("user-defined-%s", pipelineName)
	if componentName == "" {
		return suffix
	}
	return fmt.Sprintf("%s/%s", componentName, suffix)
}

// InfraMeterForKubeletStats returns an InfraMeter for kubelet stats.
func InfraMeterForKubeletStats() *policylangv1.InfraMeter {
	kubeletStatsReceiver, err := structpb.NewStruct(map[string]any{
		"collection_interval":  "10s",
		"auth_type":            "serviceAccount",
		"endpoint":             "https://${NODE_NAME}:10250",
		"insecure_skip_verify": true,
		"metric_groups": []any{
			"pod",
		},
	})
	if err != nil {
		log.Panic().Err(err).Msg("failed to create kubelet stats config")
	}

	receivers := map[string]*structpb.Struct{
		otelconsts.ReceiverKubeletStats: kubeletStatsReceiver,
	}

	kubeletStatsProcessor, err := structpb.NewStruct(map[string]any{
		"metrics": map[string]any{
			"include": map[string]any{
				"match_type": "strict",
				"metric_names": []any{
					"k8s.pod.cpu.utilization",
					"k8s.pod.memory.available",
					"k8s.pod.memory.usage",
					"k8s.pod.memory.working_set",
				},
			},
		},
	})
	if err != nil {
		log.Panic().Err(err).Msg("failed to create kubelet stats processor config")
	}

	k8sAttributesProcessor, err := structpb.NewStruct(map[string]any{
		"auth_type":   "serviceAccount",
		"passthrough": false,
		"filter": map[string]any{
			"node_from_env_var": "NODE_NAME",
		},
		"extract": map[string]any{
			"metadata": []any{
				"k8s.daemonset.name",
				"k8s.cronjob.name",
				"k8s.deployment.name",
				"k8s.job.name",
				"k8s.namespace.name",
				"k8s.node.name",
				"k8s.pod.name",
				"k8s.pod.uid",
				"k8s.replicaset.name",
				"k8s.statefulset.name",
			},
			"labels": []any{
				map[string]any{
					"key_regex": "^app.kubernetes.io/.*",
				},
			},
		},
		"pod_association": []any{
			map[string]any{
				"sources": map[string]any{
					"from": "resource_attribute",
					"name": "k8s.pod.uid",
				},
			},
		},
	})
	if err != nil {
		log.Panic().Err(err).Msg("failed to create k8s attributes processor config")
	}

	processors := map[string]*structpb.Struct{
		otelconsts.ProcessorFilterKubeletStats: kubeletStatsProcessor,
		otelconsts.ProcessorK8sAttributes:      k8sAttributesProcessor,
	}

	return &policylangv1.InfraMeter{
		Receivers:  receivers,
		Processors: processors,
		Pipeline: &policylangv1.InfraMeter_MetricsPipeline{
			Receivers: []string{
				otelconsts.ReceiverKubeletStats,
			},
			Processors: []string{
				otelconsts.ProcessorFilterKubeletStats,
				otelconsts.ProcessorK8sAttributes,
			},
		},
	}
}
