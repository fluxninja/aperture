package otel

import (
	"fmt"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/rollupprocessor"
)

const (
	// ReceiverOTLP collects traces and logs from libraries and SDKs.
	ReceiverOTLP = "otlp"
	// ReceiverPrometheus collects metrics from environment and services.
	ReceiverPrometheus = "prometheus"

	// ProcessorEnrichment enriches traces, logs and metrics with discovery data.
	ProcessorEnrichment = "enrichment"
	// ProcessorMetrics generates metrics based on traces and logs and exposes them
	// on application prometheus metrics endpoint.
	ProcessorMetrics = "metrics"
	// ProcessorBatchPrerollup batches incoming data before rolling up. This is
	// required, as rollup processor can only roll up data inside a single batch.
	ProcessorBatchPrerollup = "batch/prerollup"
	// ProcessorBatchPostrollup batches data after rolling up, as roll up process
	// shrinks number of data points significantly.
	ProcessorBatchPostrollup = "batch/postrollup"
	// ProcessorBatchMetricsFast batches metrics in small and fast packages. Used
	// in flow control policy.
	ProcessorBatchMetricsFast = "batch/metrics-fast"
	// ProcessorRollup rolls up data to decrease cardinality.
	ProcessorRollup = "rollup"

	// ExporterLogging exports telemetry using Aperture logger.
	ExporterLogging = "aperturelogging"
	// ExporterPrometheusRemoteWrite exports metrics to local prometheus instance.
	ExporterPrometheusRemoteWrite = "prometheusremotewrite"
)

type otelConfig struct {
	// NodeName is name of the node from which OTEL should collect metrics.
	NodeName string `json:"node_name"`
	// Addr is an address on which this app is serving metrics.
	// TODO this should be inherited from the listener.Listener config, but it's
	// not initialized at the provide state of app.
	Addr string `json:"addr" validate:"hostname_port" default:":8080"`
}

// ProvideAnnotatedAgentConfig provides annotated OTEL config for agent.
func ProvideAnnotatedAgentConfig() fx.Option {
	return fx.Option(
		fx.Provide(
			fx.Annotate(
				newAgentOTELConfig,
				fx.ResultTags(`name:"base"`),
			),
		),
	)
}

// ProvideAnnotatedControllerConfig provides annotated OTEL config for controller.
func ProvideAnnotatedControllerConfig() fx.Option {
	return fx.Option(
		fx.Provide(
			fx.Annotate(
				newControllerOTELConfig,
				fx.ResultTags(`name:"base"`),
			),
		),
	)
}

func newAgentOTELConfig(unmarshaller config.Unmarshaller, promClient promapi.Client) (*otelcollector.OTELConfig, error) {
	var cfg otelConfig
	if err := unmarshaller.UnmarshalKey("otel", &cfg); err != nil {
		return nil, err
	}
	config := otelcollector.NewOTELConfig()
	config.AddDebugExtensions()
	addLogsAndTracesPipelines(config)
	addMetricsPipeline(config, promClient, cfg)
	return config, nil
}

func newControllerOTELConfig(unmarshaller config.Unmarshaller, promClient promapi.Client) (*otelcollector.OTELConfig, error) {
	var cfg otelConfig
	if err := unmarshaller.UnmarshalKey("otel", &cfg); err != nil {
		return nil, err
	}
	config := otelcollector.NewOTELConfig()
	config.AddDebugExtensions()
	addControllerMetricsPipeline(config, promClient, cfg)
	return config, nil
}

func addLogsAndTracesPipelines(config *otelcollector.OTELConfig) {
	// Common dependencies for pipelines
	addOTLPReceiver(config)
	config.AddProcessor(ProcessorEnrichment, nil)
	addMetricsProcessor(config)
	config.AddBatchProcessor(ProcessorBatchPrerollup, 1*time.Second, 10000)
	addRollupProcessor(config)
	config.AddBatchProcessor(ProcessorBatchPostrollup, 1*time.Second, 10000)
	config.AddExporter(ExporterLogging, nil)

	processors := []string{
		ProcessorEnrichment,
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

	config.Service.AddPipeline("traces", otelcollector.Pipeline{
		Receivers:  []string{ReceiverOTLP},
		Processors: processors,
		Exporters:  []string{ExporterLogging},
	})
}

func addMetricsPipeline(config *otelcollector.OTELConfig, promClient promapi.Client, cfg otelConfig) {
	addPrometheusReceiver(config, cfg)
	config.AddProcessor(ProcessorEnrichment, nil)
	config.AddBatchProcessor(ProcessorBatchMetricsFast, 1*time.Second, 1000)
	addPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/fast", otelcollector.Pipeline{
		Receivers: []string{ReceiverPrometheus},
		Processors: []string{
			ProcessorEnrichment,
			ProcessorBatchMetricsFast,
		},
		Exporters: []string{ExporterPrometheusRemoteWrite},
	})
}

func addControllerMetricsPipeline(config *otelcollector.OTELConfig, promClient promapi.Client, cfg otelConfig) {
	addControllerPrometheusReceiver(config, cfg)
	config.AddBatchProcessor(ProcessorBatchMetricsFast, 1*time.Second, 1000)
	addPrometheusRemoteWriteExporter(config, promClient)
	config.Service.AddPipeline("metrics/controller-fast", otelcollector.Pipeline{
		Receivers:  []string{ReceiverPrometheus},
		Processors: []string{ProcessorBatchMetricsFast},
		Exporters:  []string{ExporterPrometheusRemoteWrite},
	})
}

func addOTLPReceiver(config *otelcollector.OTELConfig) {
	config.AddReceiver(ReceiverOTLP, otlpreceiver.Config{
		Protocols: otlpreceiver.Protocols{
			GRPC: &configgrpc.GRPCServerSettings{
				NetAddr: confignet.NetAddr{
					Endpoint:  "0.0.0.0:4317",
					Transport: "tcp",
				},
			},
			HTTP: &confighttp.HTTPServerSettings{
				Endpoint: "0.0.0.0:4318",
			},
		},
	})
}

func addMetricsProcessor(config *otelcollector.OTELConfig) {
	config.AddProcessor(ProcessorMetrics, metricsprocessor.Config{
		LatencyBucketStartMS: 20,
		LatencyBucketWidthMS: 20,
		LatencyBucketCount:   100,
	})
}

func addRollupProcessor(config *otelcollector.OTELConfig) {
	rollupFields := []string{
		otelcollector.HTTPDurationLabel,
		otelcollector.HTTPRequestContentLength,
		otelcollector.HTTPResponseContentLength,
	}
	rollups := []rollupprocessor.Rollup{}
	for _, field := range rollupFields {
		for _, t := range rollupprocessor.RollupTypes {
			rollups = append(rollups, rollupprocessor.Rollup{
				FromField: field,
				ToField:   fmt.Sprintf("%s_%s", field, t),
				Type:      t,
			})
		}
	}
	config.AddProcessor(ProcessorRollup, rollupprocessor.Config{
		Rollups: rollups,
	})
}

func addPrometheusReceiver(config *otelcollector.OTELConfig, cfg otelConfig) {
	scrapeConfigs := []map[string]interface{}{
		buildApertureSelfScrapeConfig("aperture-self", cfg),
		buildKubernetesNodesScrapeConfig(cfg),
		buildKubernetesPodsScrapeConfig(cfg),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddReceiver(ReceiverPrometheus, map[string]interface{}{
		"config": map[string]interface{}{
			"global": map[string]interface{}{
				"scrape_interval":     "1m",
				"scrape_timeout":      "10s",
				"evaluation_interval": "1m",
			},
			"scrape_configs": scrapeConfigs,
		},
	})
}

func addControllerPrometheusReceiver(config *otelcollector.OTELConfig, cfg otelConfig) {
	scrapeConfigs := []map[string]interface{}{
		buildApertureSelfScrapeConfig("aperture-controller-self", cfg),
	}
	// Unfortunately prometheus config structs do not have proper `mapstructure`
	// tags, so they are not properly read by OTEL. Need to use bare maps instead.
	config.AddReceiver(ReceiverPrometheus, map[string]interface{}{
		"config": map[string]interface{}{
			"global": map[string]interface{}{
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
	config.AddExporter(ExporterPrometheusRemoteWrite, map[string]interface{}{
		"endpoint": endpoint.String(),
	})
}

func buildApertureSelfScrapeConfig(name string, cfg otelConfig) map[string]interface{} {
	return map[string]interface{}{
		"job_name":        name,
		"scheme":          "http",
		"scrape_interval": "1s",
		"scrape_timeout":  "900ms",
		"metrics_path":    "/metrics",
		"static_configs": []map[string]interface{}{
			{
				"targets": []string{cfg.Addr},
				"labels": map[string]interface{}{
					"instance": cfg.NodeName,
				},
			},
		},
	}
}

func buildKubernetesNodesScrapeConfig(cfg otelConfig) map[string]interface{} {
	return map[string]interface{}{
		"job_name":        "kubernetes-nodes",
		"scheme":          "https",
		"scrape_interval": "2s",
		"scrape_timeout":  "1500ms",
		"metrics_path":    "/metrics/cadvisor",
		"authorization": map[string]interface{}{
			"credentials_file": "/var/run/secrets/kubernetes.io/serviceaccount/token",
		},
		"tls_config": map[string]interface{}{
			"insecure_skip_verify": true,
		},
		"kubernetes_sd_configs": []map[string]interface{}{
			{"role": "node"},
		},
		"relabel_configs": []map[string]interface{}{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_node_name"},
				"action":        "keep",
				"regex":         cfg.NodeName,
			},
		},
		"metric_relabel_configs": []map[string]interface{}{
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

func buildKubernetesPodsScrapeConfig(cfg otelConfig) map[string]interface{} {
	return map[string]interface{}{
		"job_name":        "kubernetes-pods",
		"scheme":          "http",
		"scrape_interval": "2s",
		"scrape_timeout":  "1500ms",
		"metrics_path":    "/metrics",
		"kubernetes_sd_configs": []map[string]interface{}{
			{"role": "pod"},
		},
		"relabel_configs": []map[string]interface{}{
			// Scrape only the node on which this agent is running.
			{
				"source_labels": []string{"__meta_kubernetes_pod_node_name"},
				"action":        "keep",
				"regex":         cfg.NodeName,
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
		"metric_relabel_configs": []map[string]interface{}{
			// For now, dropping everything. In future, we'll want to filter in some
			// metrics based on policies. See #4632.
			{
				"source_labels": []string{},
				"action":        "drop",
			},
		},
	}
}
