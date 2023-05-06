package config

import (
	"fmt"
	"sort"
	"strings"

	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/leaderonlyreceiver"
	"go.opentelemetry.io/collector/component"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/runtime"
)

// CustomMetricsConfig defines receivers, processors, and single metrics pipeline which will be exported to the controller Prometheus.
// Environment variables can be used in the configuration using format `${ENV_VAR_NAME}`.
// +kubebuilder:object:generate=true
//
// :::info
//
// See also [Get Started / Setup Integrations / Metrics](/get-started/integrations/metrics/metrics.md).
//
// :::
//
//swagger:model
type CustomMetricsConfig struct {
	// Receivers define receivers to be used in custom metrics pipelines. This should
	// be in OTel format - https://opentelemetry.io/docs/collector/configuration/#receivers.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Receivers Components `json:"receivers"`
	// Processors define processors to be used in custom metrics pipelines. This should
	// be in OTel format - https://opentelemetry.io/docs/collector/configuration/#processors.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Processors Components `json:"processors,omitempty"`
	// Pipeline is an OTel metrics pipeline definition, which **only** uses receivers
	// and processors defined above. Exporter would be added automatically.
	//
	// If there are no processors defined or only one processor is defined, the
	// pipeline definition can be omitted. In such cases, the pipeline will
	// automatically use all given receivers and the defined processor (if
	// any).  However, if there are more than one processor, the pipeline must
	// be defined explicitly.
	Pipeline CustomMetricsPipelineConfig `json:"pipeline"`
	// PerAgentGroup marks the pipeline to be instantiated only once per agent
	// group. This is helpful for receivers that scrape for example, some cluster-wide
	// metrics. When not set, pipeline will be instantiated on every Agent.
	PerAgentGroup bool `json:"per_agent_group"`
}

// CustomMetricsPipelineConfig defines a custom metrics pipeline.
// +kubebuilder:object:generate=true
//
//swagger:model
type CustomMetricsPipelineConfig struct {
	Receivers  []string `json:"receivers,omitempty"`
	Processors []string `json:"processors,omitempty"`
}

// Components is an alias type for map[string]any. This needs to be used
// because of the CRD requirements for the operator.
// https://github.com/kubernetes-sigs/controller-tools/issues/636
// https://github.com/kubernetes-sigs/kubebuilder/issues/528
// +kubebuilder:object:generate=false
// +kubebuilder:pruning:PreserveUnknownFields
// +kubebuilder:validation:Schemaless
type Components map[string]any

// DeepCopyInto is an deepcopy function, copying the receiver, writing into out.
// In must be non-nil.
// We need to specify this manually, as the generator does not support `any`.
func (in *Components) DeepCopyInto(out *Components) {
	if in == nil {
		*out = nil
	} else {
		*out = runtime.DeepCopyJSON(*in)
	}
}

// DeepCopy is an deepcopy function, copying the receiver, creating a new
// Components.
// We need to specify this manually, as the generator does not support `any`.
func (in *Components) DeepCopy() *Components {
	if in == nil {
		return nil
	}
	out := new(Components)
	in.DeepCopyInto(out)
	return out
}

// AddCustomMetricsPipelines adds custom metrics pipelines to the given OTelConfig.
func AddCustomMetricsPipelines(
	config *otelconfig.OTelConfig,
	customMetrics map[string]CustomMetricsConfig,
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
		customMetrics = map[string]CustomMetricsConfig{}
	}
	if _, ok := customMetrics[otelconsts.ReceiverKubeletStats]; !ok {
		customMetrics[otelconsts.ReceiverKubeletStats] = makeCustomMetricsConfigForKubeletStats()
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
	metricConfig CustomMetricsConfig,
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
		id, receiverConfig = leaderonlyreceiver.WrapConfigIf(metricConfig.PerAgentGroup, id, receiverConfig)
		receiverIDs[origName] = id.String()
		config.AddReceiver(id.String(), receiverConfig)
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

func makeCustomMetricsConfigForKubeletStats() CustomMetricsConfig {
	receivers := map[string]any{
		otelconsts.ReceiverKubeletStats: map[string]any{
			"collection_interval":  "10s",
			"auth_type":            "serviceAccount",
			"endpoint":             "https://${NODE_NAME}:10250",
			"insecure_skip_verify": true,
			"metric_groups": []any{
				"pod",
			},
		},
	}
	processors := map[string]any{
		otelconsts.ProcessorFilterKubeletStats: map[string]any{
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
		},
		otelconsts.ProcessorK8sAttributes: map[string]any{
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
		},
	}
	return CustomMetricsConfig{
		Receivers:  receivers,
		Processors: processors,
		Pipeline: CustomMetricsPipelineConfig{
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
