package agent

import (
	"crypto/tls"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"k8s.io/client-go/rest"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/tracestologsprocessor"
)

// swagger:operation POST /otel agent-configuration OTEL
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/AgentOTELConfig"

// AgentOTELConfig is the configuration for Agent's OTEL collector.
// swagger:model
// +kubebuilder:object:generate=true
type AgentOTELConfig struct {
	otelconfig.CommonOTELConfig `json:",inline"`
	// BatchPrerollup configures batch prerollup processor.
	BatchPrerollup BatchPrerollupConfig `json:"batch_prerollup"`
	// BatchPostrollup configures batch postrollup processor.
	BatchPostrollup BatchPostrollupConfig `json:"batch_postrollup"`
}

// BatchPrerollupConfig defines configuration for OTEL batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchPrerollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"10s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"10000"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"10000"`
}

// BatchPostrollupConfig defines configuration for OTEL batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchPostrollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"100"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"100"`
}

func provideAgent(
	unmarshaller config.Unmarshaller,
	lis *listener.Listener,
	promClient promapi.Client,
	tlsConfig *tls.Config,
) (*otelconfig.OTELConfig, error) {
	var agentCfg AgentOTELConfig
	if err := unmarshaller.UnmarshalKey("otel", &agentCfg); err != nil {
		return nil, err
	}

	otelCfg := otelconfig.NewOTELConfig()
	otelCfg.SetDebugPort(&agentCfg.CommonOTELConfig)
	otelCfg.AddDebugExtensions(&agentCfg.CommonOTELConfig)

	addLogsPipeline(otelCfg, &agentCfg)
	addTracesPipeline(otelCfg, lis)
	addMetricsPipeline(otelCfg, &agentCfg, tlsConfig, lis, promClient)
	otelconfig.AddAlertsPipeline(otelCfg, agentCfg.CommonOTELConfig, otelconsts.ProcessorAgentResourceLabels)
	return otelCfg, nil
}

func addLogsPipeline(config *otelconfig.OTELConfig, userConfig *AgentOTELConfig) {
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
	agentConfig *AgentOTELConfig,
	tlsConfig *tls.Config,
	lis *listener.Listener,
	promClient promapi.Client,
) {
	addPrometheusReceiver(config, agentConfig, tlsConfig, lis)
	config.AddProcessor(otelconsts.ProcessorEnrichment, nil)
	otelconfig.AddPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/fast", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverPrometheus},
		Processors: []string{
			otelconsts.ProcessorEnrichment,
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addPrometheusReceiver(
	config *otelconfig.OTELConfig,
	agentConfig *AgentOTELConfig,
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
		scrapeConfigs = append(scrapeConfigs, buildKubernetesNodesScrapeConfig(), buildKubernetesPodsScrapeConfig())
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

func buildKubernetesNodesScrapeConfig() map[string]any {
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
