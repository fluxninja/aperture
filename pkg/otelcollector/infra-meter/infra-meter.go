package inframeter

import (
	"fmt"
	"sort"
	"strings"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/leaderonlyreceiver"
	"go.opentelemetry.io/collector/component"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/structpb"
)

// AddInfraMeters adds infra metrics pipelines to the given OTelConfig.
func AddInfraMeters(
	config *otelconfig.Config,
	infraMeters map[string]*policylangv1.InfraMeter,
) error {
	config.AddProcessor(otelconsts.ProcessorInfraMeter, map[string]any{
		"attributes": []map[string]interface{}{
			{
				"key":    "service.name",
				"action": "upsert",
				"value":  "aperture-infra-meter",
			},
		},
	})
	if infraMeters == nil {
		infraMeters = map[string]*policylangv1.InfraMeter{}
	}
	for pipelineName, metricConfig := range infraMeters {
		if err := addInfraMeter(config, pipelineName, metricConfig); err != nil {
			return fmt.Errorf("failed to add infra metric pipeline %s: %w", pipelineName, err)
		}
	}
	return nil
}

func addInfraMeter(
	config *otelconfig.Config,
	pipelineName string,
	infraMeter *policylangv1.InfraMeter,
) error {
	pipelineName = strings.TrimPrefix(pipelineName, "metrics/")

	receiverIDs := map[string]string{}
	processorIDs := map[string]string{}

	for origName, receiverConfig := range infraMeter.Receivers {
		var id component.ID
		if err := id.UnmarshalText([]byte(origName)); err != nil {
			return fmt.Errorf("invalid id %q: %w", origName, err)
		}
		id = component.NewIDWithName(id.Type(), normalizeComponentName(pipelineName, id.Name()))
		var cfg any
		cfg = receiverConfig.AsMap()
		id, cfg = leaderonlyreceiver.WrapConfigIf(infraMeter.PerAgentGroup, id, cfg)
		receiverIDs[origName] = id.String()
		config.AddReceiver(id.String(), cfg)
	}

	for origName, processorConfig := range infraMeter.Processors {
		id := normalizeComponentName(pipelineName, origName)
		processorIDs[origName] = id
		var cfg any = processorConfig.AsMap()
		config.AddProcessor(id, cfg)
	}

	if infraMeter.Pipeline == nil {
		// We treat empty pipeline the same way as not-set pipeline, normalize.
		// This also allows to avoid nil checks below.
		infraMeter.Pipeline = &policylangv1.InfraMeter_MetricsPipeline{}
	}

	if len(infraMeter.Pipeline.Receivers) == 0 && len(infraMeter.Pipeline.Processors) == 0 {
		if len(infraMeter.Processors) >= 1 {
			return fmt.Errorf("empty pipeline, inferring pipeline is supported only with 0 or 1 processors")
		}

		// Skip adding pipeline if there are no receivers and processors.
		if len(infraMeter.Receivers) == 0 && len(infraMeter.Processors) == 0 {
			return nil
		}

		// When pipeline not set explicitly, create pipeline with all defined receivers and processors.
		if len(infraMeter.Receivers) > 0 {
			infraMeter.Pipeline.Receivers = maps.Keys(infraMeter.Receivers)
			sort.Strings(infraMeter.Pipeline.Receivers)
		}
		if len(infraMeter.Processors) > 0 {
			infraMeter.Pipeline.Processors = maps.Keys(infraMeter.Processors)
		}
	}

	config.Service.AddPipeline(normalizePipelineName(pipelineName), otelconfig.Pipeline{
		Receivers: mapSlice(receiverIDs, infraMeter.Pipeline.Receivers),
		Processors: append(
			mapSlice(processorIDs, infraMeter.Pipeline.Processors),
			otelconsts.ProcessorInfraMeter,
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
