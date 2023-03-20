package agent

import (
	"crypto/tls"
	"fmt"
	"strings"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"k8s.io/client-go/rest"

	agentconfig "github.com/fluxninja/aperture/cmd/aperture-agent/config"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/tracestologsprocessor"
)

func provideAgent(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
) (*otelconfig.OTELConfig, error) {
	var agentCfg agentconfig.AgentOTELConfig
	if err := unmarshaller.UnmarshalKey("otel", &agentCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.NewOTELConfig()
	otelCfg.SetDebugPort(&agentCfg.CommonOTELConfig)
	otelCfg.AddDebugExtensions(&agentCfg.CommonOTELConfig)

	addLogsPipeline(otelCfg, &agentCfg)
	addTracesPipeline(otelCfg, lis)
	addMetricsPipeline(otelCfg, &agentCfg, tlsConfig, lis, promClient)
	addCustomMetricsPipelines(otelCfg, &agentCfg)
	otelconfig.AddAlertsPipeline(otelCfg, agentCfg.CommonOTELConfig, otelconsts.ProcessorAgentResourceLabels)
	return otelCfg, nil
}

func addLogsPipeline(
	config *otelconfig.OTELConfig,
	userConfig *agentconfig.AgentOTELConfig,
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
		Receivers:  []string{otelconsts.ReceiverOTLP},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterLogging},
	})
}

func addTracesPipeline(config *otelconfig.OTELConfig, lis *listener.Listener) {
	config.AddExporter(otelconsts.ExporterOTLPLoopback, map[string]any{
		"endpoint": lis.GetAddr(),
		"tls": map[string]any{
			"insecure": true,
		},
	})
	config.AddProcessor(otelconsts.ProcessorTracesToLogs, tracestologsprocessor.Config{
		LogsExporter: otelconsts.ExporterOTLPLoopback,
	})

	config.Service.AddPipeline("traces", otelconfig.Pipeline{
		Receivers:  []string{otelconsts.ReceiverOTLP},
		Processors: []string{otelconsts.ProcessorTracesToLogs},
		// We need some exporter configured to make this pipeline correct. Actual
		// Log exporting is done inside the processor.
		Exporters: []string{otelconsts.ExporterLogging},
	})

	// TODO This receiver should be replaced with some receiver which really does nothing.
	config.AddReceiver("filelog", map[string]any{
		"include":       []string{"/var/log/myservice/*.json"},
		"poll_interval": "1000h",
	})
	// We need a fake log pipeline which will initialize the ExporterOTLPLoopback
	// for logs type.
	config.Service.AddPipeline("logs/fake", otelconfig.Pipeline{
		Receivers: []string{"filelog"},
		Exporters: []string{otelconsts.ExporterOTLPLoopback},
	})
}

func addMetricsPipeline(
	config *otelconfig.OTELConfig,
	agentConfig *agentconfig.AgentOTELConfig,
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
	config *otelconfig.OTELConfig,
	agentConfig *agentconfig.AgentOTELConfig,
) {
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
		pipelineName = strings.TrimPrefix(pipelineName, "metrics/")
		for receiverName, receiverConfig := range metricConfig.Receivers {
			config.AddReceiver(normalizeComponentName(pipelineName, receiverName), receiverConfig)
		}
		for processorName, processorConfig := range metricConfig.Processors {
			config.AddProcessor(normalizeComponentName(pipelineName, processorName), processorConfig)
		}
		config.Service.AddPipeline(normalizePipelineName(pipelineName), otelconfig.Pipeline{
			Receivers: normalizeComponentNames(pipelineName, metricConfig.Pipeline.Receivers),
			Processors: append(
				normalizeComponentNames(pipelineName, metricConfig.Pipeline.Processors),
				otelconsts.ProcessorCustomMetrics,
				otelconsts.ProcessorAgentResourceLabels,
			),
			Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
		})
	}
}

// normalizePipelineName normalizes user defined pipeline name by adding
// `metrics/user-defined-` prefix.
// This ensures no builtin metrics pipeline is overwritten.
func normalizePipelineName(pipelineName string) string {
	return fmt.Sprintf("metrics/user-defined-%s", pipelineName)
}

// normalizeComponentNames calls `normalizeComponentName` for each element of the
// slice. Returns new slice with modified elements.
func normalizeComponentNames(pipelineName string, components []string) []string {
	renamed := make([]string, len(components))
	for i, c := range components {
		renamed[i] = normalizeComponentName(pipelineName, c)
	}
	return renamed
}

// normalizeComponentName normalizes user defines component name by adding
// `user-defined-<pipeline_name>` suffix.
// This ensures no builtin components are overwritten.
func normalizeComponentName(pipelineName, componentName string) string {
	return fmt.Sprintf("%s/user-defined-%s", componentName, pipelineName)
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
	config *otelconfig.OTELConfig,
	agentConfig *agentconfig.AgentOTELConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
) {
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-self", tlsConfig, lis),
		otelconfig.BuildOTELScrapeConfig("aperture-otel", agentConfig.CommonOTELConfig),
	}

	_, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		log.Debug().Msg("K8s environment not detected. Skipping K8s scrape configurations.")
	} else if err != nil {
		log.Warn().Err(err).Msg("Error when discovering k8s environment")
	} else {
		log.Debug().Msg("K8s environment detected. Adding K8s scrape configurations.")
		scrapeConfigs = append(scrapeConfigs, buildKubernetesPodsScrapeConfig())
	}

	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
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
