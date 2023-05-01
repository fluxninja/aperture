package agent

import (
	"crypto/tls"
	"fmt"
	"sort"
	"strings"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"golang.org/x/exp/maps"
	"k8s.io/client-go/rest"

	agentconfig "github.com/fluxninja/aperture/cmd/aperture-agent/config"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/leaderonlyreceiver"
)

func provideAgent(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
) (*otelconfig.OTelConfig, error) {
	var agentCfg agentconfig.AgentOTelConfig
	if err := unmarshaller.UnmarshalKey("otel", &agentCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.NewOTelConfig()
	otelCfg.SetDebugPort(&agentCfg.CommonOTelConfig)
	otelCfg.AddDebugExtensions(&agentCfg.CommonOTelConfig)

	addLogsPipeline(otelCfg, &agentCfg)
	addTracesPipeline(otelCfg, lis)
	addMetricsPipeline(otelCfg, &agentCfg, tlsConfig, lis, promClient)
	if err := addCustomMetricsPipelines(otelCfg, &agentCfg); err != nil {
		return nil, err
	}
	otelconfig.AddAlertsPipeline(otelCfg, agentCfg.CommonOTelConfig, otelconsts.ProcessorAgentResourceLabels)
	return otelCfg, nil
}

func addLogsPipeline(
	config *otelconfig.OTelConfig,
	userConfig *agentconfig.AgentOTelConfig,
) {
	// Common dependencies for pipelines
	config.AddReceiver(otelconsts.ReceiverOTLP, otlpreceiver.Config{})
	// Note: Passing map[string]interface{}{} instead of real config, so that
	// processors' configs' default work.
	config.AddProcessor(otelconsts.ProcessorMetrics, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPrerollup,
		userConfig.BatchPrerollup.Timeout.AsDuration(),
		userConfig.BatchPrerollup.SendBatchSize,
		userConfig.BatchPrerollup.SendBatchMaxSize,
	)
	config.AddProcessor(otelconsts.ProcessorRollup, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPostrollup,
		userConfig.BatchPostrollup.Timeout.AsDuration(),
		userConfig.BatchPostrollup.SendBatchSize,
		userConfig.BatchPostrollup.SendBatchMaxSize,
	)
	config.AddExporter(otelconsts.ExporterLogging, nil)

	processors := []string{
		otelconsts.ProcessorMetrics,
		otelconsts.ProcessorBatchPrerollup,
		otelconsts.ProcessorRollup,
		otelconsts.ProcessorBatchPostrollup,
		otelconsts.ProcessorAgentGroup,
	}

	config.Service.AddPipeline("logs", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverOTLP, otelconsts.ConnectorAdapter},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterLogging},
	})
}

func addTracesPipeline(config *otelconfig.OTelConfig, lis *listener.Listener) {
	config.AddConnector(otelconsts.ConnectorAdapter, map[string]any{})
	config.Service.AddPipeline("traces", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverOTLP},
		Exporters: []string{otelconsts.ConnectorAdapter},
	})
}

func addMetricsPipeline(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
	promClient promapi.Client,
) {
	addPrometheusReceiver(config, agentConfig, tlsConfig, lis)
	otelconfig.AddPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/fast", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverPrometheus},
		Processors: []string{
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addCustomMetricsPipelines(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
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
	if agentConfig.CustomMetrics == nil {
		agentConfig.CustomMetrics = map[string]agentconfig.CustomMetricsConfig{}
	}
	if _, ok := agentConfig.CustomMetrics[otelconsts.ReceiverKubeletStats]; !ok {
		agentConfig.CustomMetrics[otelconsts.ReceiverKubeletStats] = makeCustomMetricsConfigForKubeletStats()
	}
	for pipelineName, metricConfig := range agentConfig.CustomMetrics {
		if err := addCustomMetricsPipeline(config, pipelineName, metricConfig); err != nil {
			return fmt.Errorf("failed to add custom metric pipeline %s: %w", pipelineName, err)
		}
	}
	return nil
}

func addCustomMetricsPipeline(
	config *otelconfig.OTelConfig,
	pipelineName string,
	metricConfig agentconfig.CustomMetricsConfig,
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

func makeCustomMetricsConfigForKubeletStats() agentconfig.CustomMetricsConfig {
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
	return agentconfig.CustomMetricsConfig{
		Receivers:  receivers,
		Processors: processors,
		Pipeline: agentconfig.CustomMetricsPipelineConfig{
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

func addPrometheusReceiver(
	config *otelconfig.OTelConfig,
	agentConfig *agentconfig.AgentOTelConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-self", tlsConfig, lis),
		otelconfig.BuildOTelScrapeConfig("aperture-otel", agentConfig.CommonOTelConfig),
	}

	if !agentConfig.DisableKubernetesScraper {
		_, err := rest.InClusterConfig()
		if err == rest.ErrNotInCluster {
			log.Debug().Msg("K8s environment not detected. Skipping K8s scrape configurations.")
		} else if err != nil {
			log.Warn().Err(err).Msg("Error when discovering k8s environment")
		} else {
			log.Debug().Msg("K8s environment detected. Adding K8s scrape configurations.")
			scrapeConfigs = append(scrapeConfigs, buildKubernetesPodsScrapeConfig())
		}
	}

	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTel. Need to use bare maps instead.
	config.AddReceiver(otelconsts.ReceiverPrometheus, map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval":     "1s",
				"scrape_timeout":      "1s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}

func buildKubernetesPodsScrapeConfig() map[string]any {
	return map[string]any{
		"job_name":     "kubernetes-pods",
		"scheme":       "http",
		"metrics_path": "/metrics",
		"kubernetes_sd_configs": []map[string]any{
			{"role": "pod"},
		},
		"relabel_configs": []map[string]any{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_pod_node_name"},
				"action":        "keep",
				"regex":         info.Hostname,
			},
			// Scrape only pods which have github.com/fluxninja/scrape=true annotation.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_aperture_tech_scrape"},
				"action":        "keep",
				"regex":         "true",
			},
			// Allow rewrite of scheme, path and port where prometheus metrics are served.
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
				"action":        "replace",
				"regex":         "(https?)",
				"target_label":  "__scheme__",
			},
			{
				"source_labels": []string{"__meta_kubernetes_pod_annotation_prometheus_io_path"},
				"action":        "replace",
				"target_label":  "__metrics_path__",
				"regex":         "(.+)",
			},
			{
				"source_labels": []string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
				"action":        "replace",
				"regex":         `([^:]+)(?::\d+)?;(\d+)`,
				"replacement":   "$$1:$$2",
				"target_label":  "__address__",
			},
		},
		"metric_relabel_configs": []map[string]any{
			// For now, dropping everything. In future, we'll want to filter in some
			// metrics based on policies. See #4632.
			{
				"source_labels": []string{},
				"action":        "drop",
			},
		},
	}
}
