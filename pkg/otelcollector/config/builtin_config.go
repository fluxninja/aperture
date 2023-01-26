package config

import (
	"fmt"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	promapi "github.com/prometheus/client_golang/api"
	"k8s.io/client-go/rest"
)

// AddMetricsPipeline adds metrics to pipeline for agent OTEL collector.
func AddMetricsPipeline(cfg *OTELParams) {
	config := cfg.Config
	addPrometheusReceiver(cfg)
	config.AddProcessor(otelconsts.ProcessorEnrichment, nil)
	addPrometheusRemoteWriteExporter(config, cfg.promClient)
	config.Service.AddPipeline("metrics/fast", Pipeline{
		Receivers: []string{otelconsts.ReceiverPrometheus},
		Processors: []string{
			otelconsts.ProcessorEnrichment,
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

// AddControllerMetricsPipeline adds metrics to pipeline for controller OTEL collector.
func AddControllerMetricsPipeline(cfg *OTELParams) {
	config := cfg.Config
	addControllerPrometheusReceiver(config, cfg)
	addPrometheusRemoteWriteExporter(config, cfg.promClient)
	config.Service.AddPipeline("metrics/controller-fast", Pipeline{
		Receivers:  []string{otelconsts.ReceiverPrometheus},
		Processors: []string{},
		Exporters:  []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

// AddAlertsPipeline adds reusable alerts pipeline.
func AddAlertsPipeline(cfg *OTELParams, extraProcessors ...string) {
	config := cfg.Config
	config.AddReceiver(otelconsts.ReceiverAlerts, map[string]any{})
	config.AddProcessor(otelconsts.ProcessorAlertsNamespace, map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    otelconsts.AlertNamespaceLabel,
				"action": "insert",
				"value":  info.Hostname,
			},
		},
	})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchAlerts,
		cfg.BatchAlerts.Timeout.AsDuration(),
		cfg.BatchAlerts.SendBatchSize,
		cfg.BatchAlerts.SendBatchMaxSize,
	)
	config.AddExporter(otelconsts.ExporterAlerts, nil)

	processors := []string{
		otelconsts.ProcessorBatchAlerts,
		otelconsts.ProcessorAlertsNamespace,
	}
	processors = append(processors, extraProcessors...)
	config.Service.AddPipeline("logs/alerts", Pipeline{
		Receivers:  []string{otelconsts.ReceiverAlerts},
		Processors: processors,
		Exporters:  []string{otelconsts.ExporterLogging, otelconsts.ExporterAlerts},
	})
}

func addPrometheusReceiver(cfg *OTELParams) {
	config := cfg.Config
	scrapeConfigs := []map[string]any{
		buildApertureSelfScrapeConfig("aperture-self", cfg),
		buildOTELScrapeConfig("aperture-otel", cfg),
	}

	_, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		log.Debug().Msg("K8s environment not detected. Skipping K8s scrape configurations.")
	} else if err != nil {
		log.Warn().Err(err).Msg("Error when discovering k8s environment")
	} else {
		log.Debug().Msg("K8s environment detected. Adding K8s scrape configurations.")
		scrapeConfigs = append(scrapeConfigs, buildKubernetesNodesScrapeConfig(cfg), buildKubernetesPodsScrapeConfig(cfg))
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

func addControllerPrometheusReceiver(config *OTELConfig, cfg *OTELParams) {
	scrapeConfigs := []map[string]any{
		buildApertureSelfScrapeConfig("aperture-controller-self", cfg),
		buildOTELScrapeConfig("aperture-controller-otel", cfg),
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

func addPrometheusRemoteWriteExporter(config *OTELConfig, promClient promapi.Client) {
	endpoint := promClient.URL("api/v1/write", nil)
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddExporter(otelconsts.ExporterPrometheusRemoteWrite, map[string]any{
		"endpoint": endpoint.String(),
	})
}

func buildApertureSelfScrapeConfig(name string, cfg *OTELParams) map[string]any {
	scheme := "http"
	if cfg.tlsConfig != nil {
		scheme = "https"
	}
	return map[string]any{
		"job_name": name,
		"scheme":   scheme,
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"metrics_path": "/metrics",
		"static_configs": []map[string]any{
			{
				"targets": []string{cfg.Listener.GetAddr()},
				"labels": map[string]any{
					metrics.InstanceLabel:    info.Hostname,
					metrics.ProcessUUIDLabel: info.UUID,
				},
			},
		},
	}
}

func buildOTELScrapeConfig(name string, cfg *OTELParams) map[string]any {
	otelDebugTarget := fmt.Sprintf(":%d", cfg.Ports.DebugPort)
	return map[string]any{
		"job_name": name,
		"scheme":   "http",
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"metrics_path": "/metrics",
		"static_configs": []map[string]any{
			{
				"targets": []string{otelDebugTarget},
				"labels": map[string]any{
					metrics.InstanceLabel:    info.Hostname,
					metrics.ProcessUUIDLabel: info.UUID,
				},
			},
		},
	}
}

func buildKubernetesNodesScrapeConfig(cfg *OTELParams) map[string]any {
	return map[string]any{
		"job_name":     "kubernetes-nodes",
		"scheme":       "https",
		"metrics_path": "/metrics/cadvisor",
		"authorization": map[string]any{
			"credentials_file": "/var/run/secrets/kubernetes.io/serviceaccount/token",
		},
		"tls_config": map[string]any{
			"insecure_skip_verify": true,
		},
		"kubernetes_sd_configs": []map[string]any{
			{"role": "node"},
		},
		"relabel_configs": []map[string]any{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_node_name"},
				"action":        "keep",
				"regex":         info.Hostname,
			},
		},
		"metric_relabel_configs": []map[string]any{
			{
				"source_labels": []string{"__name__"},
				"action":        "keep",
				"regex":         "container_memory_working_set_bytes|container_spec_memory_limit_bytes|container_spec_cpu_(?:quota|period)|container_cpu_usage_seconds_total",
			},
			{
				"source_labels": []string{"pod"},
				"action":        "replace",
				"target_label":  "entity_name",
			},
		},
	}
}

func buildKubernetesPodsScrapeConfig(cfg *OTELParams) map[string]any {
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
