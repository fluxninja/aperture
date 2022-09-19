// +kubebuilder:validation:Optional
package otel

import (
	"crypto/tls"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/rollupprocessor"
)

const (
	// ReceiverOTLP collects logs from libraries and SDKs.
	ReceiverOTLP = "otlp"
	// ReceiverPrometheus collects metrics from environment and services.
	ReceiverPrometheus = "prometheus"

	// ProcessorEnrichment enriches metrics with discovery data.
	ProcessorEnrichment = "enrichment"
	// ProcessorMetrics generates metrics based on logs and exposes them
	// on application prometheus metrics endpoint.
	ProcessorMetrics = "metrics"
	// ProcessorBatchPrerollup batches incoming data before rolling up. This is
	// required, as rollup processor can only roll up data inside a single batch.
	ProcessorBatchPrerollup = "batch/prerollup"
	// ProcessorBatchPostrollup batches data after rolling up, as roll up process
	// shrinks number of data points significantly.
	ProcessorBatchPostrollup = "batch/postrollup"
	// ProcessorRollup rolls up data to decrease cardinality.
	ProcessorRollup = "rollup"
	// ProcessorAgentGroup adds `agent_group` attribute.
	ProcessorAgentGroup = "attributes/agent_group"

	// ExporterLogging exports telemetry using Aperture logger.
	ExporterLogging = "aperturelogging"
	// ExporterPrometheusRemoteWrite exports metrics to local prometheus instance.
	ExporterPrometheusRemoteWrite = "prometheusremotewrite"
)

var baseFxTag = config.NameTag("base")

type otelParams struct {
	promClient      promapi.Client
	config          *otelcollector.OTELConfig
	listener        *listener.Listener
	tlsConfig       *tls.Config
	serverTLSConfig tlsconfig.ServerTLSConfig
	OtelConfig
}

// swagger:operation POST /otel common-configuration Otel
// ---
// x-fn-config-env: true
// parameters:
// - name: proxy
//   in: body
//   schema:
//     "$ref": "#/definitions/OtelConfig"

// OtelConfig is the configuration for the OTEL collector.
// swagger:model
// +kubebuilder:object:generate=true
type OtelConfig struct {
	// GRPC listener addr for OTEL Collector.
	GRPCAddr string `json:"grpc_addr" validate:"hostname_port" default:":4317"`
	// HTTP listener addr for OTEL Collector.
	HTTPAddr string `json:"http_addr" validate:"hostname_port" default:":4318"`
	// BatchPrerollup configures batch prerollup processor.
	BatchPrerollup BatchConfig `json:"batch_prerollup"`
	// BatchPostrollup configures batch postrollup processor.
	BatchPostrollup BatchConfig `json:"batch_postrollup"`
}

// BatchConfig defines configuration for OTEL batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"10000"`
}

// Type decides which configuration to use.
type Type int

const (
	// AgentType instantiates agent pipeline.
	AgentType Type = iota
	// ControllerType instantiates controller pipeline.
	ControllerType
)

// OTELConfigConstructor is the constructor for the OTEL collector configuration.
type OTELConfigConstructor struct {
	Type Type
}

// Annotate provides fx options.
func (c OTELConfigConstructor) Annotate() fx.Option {
	options := fx.Provide(newOtelConfig)
	switch c.Type {
	case AgentType:
		options = fx.Options(options, fx.Provide(fx.Annotate(provideAgent, fx.ResultTags(baseFxTag))))
	case ControllerType:
		options = fx.Options(options, fx.Provide(fx.Annotate(provideController, fx.ResultTags(baseFxTag))))
	}
	return options
}

// FxIn consumes parameters via Fx.
type FxIn struct {
	fx.In
	Unmarshaller    config.Unmarshaller
	Listener        *listener.Listener
	PromClient      promapi.Client
	TLSConfig       *tls.Config
	ServerTLSConfig tlsconfig.ServerTLSConfig
}

func newOtelConfig(in FxIn) (*otelParams, error) {
	config := otelcollector.NewOTELConfig()
	config.AddDebugExtensions()

	var userCfg OtelConfig
	if err := in.Unmarshaller.UnmarshalKey("otel", &userCfg); err != nil {
		return nil, err
	}
	cfg := &otelParams{
		OtelConfig:      userCfg,
		listener:        in.Listener,
		promClient:      in.PromClient,
		tlsConfig:       in.TLSConfig,
		config:          config,
		serverTLSConfig: in.ServerTLSConfig,
	}
	return cfg, nil
}

func provideAgent(cfg *otelParams) *otelcollector.OTELConfig {
	addLogsPipeline(cfg)
	addMetricsPipeline(cfg)
	return cfg.config
}

func provideController(cfg *otelParams) *otelcollector.OTELConfig {
	addControllerMetricsPipeline(cfg)
	return cfg.config
}

func addLogsPipeline(cfg *otelParams) {
	config := cfg.config
	// Common dependencies for pipelines
	addOTLPReceiver(cfg)
	addMetricsProcessor(config)
	config.AddBatchProcessor(ProcessorBatchPrerollup, cfg.BatchPrerollup.Timeout.AsDuration(), cfg.BatchPrerollup.SendBatchSize)
	addRollupProcessor(config)
	config.AddBatchProcessor(ProcessorBatchPostrollup, cfg.BatchPostrollup.Timeout.AsDuration(), cfg.BatchPostrollup.SendBatchSize)
	config.AddExporter(ExporterLogging, nil)

	processors := []string{
		ProcessorAgentGroup,
		ProcessorMetrics,
		ProcessorBatchPrerollup,
		ProcessorRollup,
		ProcessorBatchPostrollup,
	}

	config.Service.AddPipeline("logs", otelcollector.Pipeline{
		Receivers:  []string{ReceiverOTLP},
		Processors: processors,
		Exporters:  []string{ExporterLogging},
	})
}

func addMetricsPipeline(cfg *otelParams) {
	config := cfg.config
	addPrometheusReceiver(cfg)
	config.AddProcessor(ProcessorEnrichment, nil)
	addPrometheusRemoteWriteExporter(config, cfg.promClient)
	config.Service.AddPipeline("metrics/fast", otelcollector.Pipeline{
		Receivers: []string{ReceiverPrometheus},
		Processors: []string{
			ProcessorEnrichment,
			ProcessorAgentGroup,
		},
		Exporters: []string{ExporterPrometheusRemoteWrite},
	})
}

func addControllerMetricsPipeline(cfg *otelParams) {
	config := cfg.config
	addControllerPrometheusReceiver(config, cfg)
	addPrometheusRemoteWriteExporter(config, cfg.promClient)
	config.Service.AddPipeline("metrics/controller-fast", otelcollector.Pipeline{
		Receivers:  []string{ReceiverPrometheus},
		Processors: []string{},
		Exporters:  []string{ExporterPrometheusRemoteWrite},
	})
}

func addOTLPReceiver(cfg *otelParams) {
	config := cfg.config
	config.AddReceiver(ReceiverOTLP, otlpreceiver.Config{
		Protocols: otlpreceiver.Protocols{
			GRPC: &configgrpc.GRPCServerSettings{
				NetAddr: confignet.NetAddr{
					Endpoint:  cfg.GRPCAddr,
					Transport: "tcp",
				},
			},
			HTTP: &confighttp.HTTPServerSettings{
				Endpoint: cfg.HTTPAddr,
			},
		},
	})
}

func addMetricsProcessor(config *otelcollector.OTELConfig) {
	config.AddProcessor(ProcessorMetrics, metricsprocessor.Config{})
}

func addRollupProcessor(config *otelcollector.OTELConfig) {
	config.AddProcessor(ProcessorRollup, rollupprocessor.Config{})
}

func addPrometheusReceiver(cfg *otelParams) {
	config := cfg.config
	scrapeConfigs := []map[string]any{
		buildApertureSelfScrapeConfig("aperture-self", cfg),
		buildKubernetesNodesScrapeConfig(cfg),
		buildKubernetesPodsScrapeConfig(cfg),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddReceiver(ReceiverPrometheus, map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval":     "1m",
				"scrape_timeout":      "10s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}

func addControllerPrometheusReceiver(config *otelcollector.OTELConfig, cfg *otelParams) {
	scrapeConfigs := []map[string]any{
		buildApertureSelfScrapeConfig("aperture-controller-self", cfg),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddReceiver(ReceiverPrometheus, map[string]any{
		"config": map[string]any{
			"global": map[string]any{
				"scrape_interval":     "1m",
				"scrape_timeout":      "10s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}

func addPrometheusRemoteWriteExporter(config *otelcollector.OTELConfig, promClient promapi.Client) {
	endpoint := promClient.URL("api/v1/write", nil)
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddExporter(ExporterPrometheusRemoteWrite, map[string]any{
		"endpoint": endpoint.String(),
	})
}

func buildApertureSelfScrapeConfig(name string, cfg *otelParams) map[string]any {
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
		"scrape_interval": "1s",
		"scrape_timeout":  "900ms",
		"metrics_path":    "/metrics",
		"static_configs": []map[string]any{
			{
				"targets": []string{cfg.listener.GetAddr()},
				"labels": map[string]any{
					"instance":     info.Hostname,
					"process_uuid": info.UUID,
				},
			},
		},
	}
}

func buildKubernetesNodesScrapeConfig(cfg *otelParams) map[string]any {
	return map[string]any{
		"job_name":        "kubernetes-nodes",
		"scheme":          "https",
		"scrape_interval": "2s",
		"scrape_timeout":  "1500ms",
		"metrics_path":    "/metrics/cadvisor",
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

func buildKubernetesPodsScrapeConfig(cfg *otelParams) map[string]any {
	return map[string]any{
		"job_name":        "kubernetes-pods",
		"scheme":          "http",
		"scrape_interval": "2s",
		"scrape_timeout":  "1500ms",
		"metrics_path":    "/metrics",
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
