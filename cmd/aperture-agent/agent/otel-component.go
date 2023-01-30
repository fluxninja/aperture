package agent

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"k8s.io/client-go/rest"

	"github.com/fluxninja/aperture/pkg/alertmanager"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsexporter"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsreceiver"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/otelcollector/enrichmentprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/rollupprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/tracestologsprocessor"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

// ModuleForAgentOTEL provides fx options for AgentOTELComponent.
func ModuleForAgentOTEL() fx.Option {
	return fx.Options(
		fx.Provide(
			cache.Provide[selectors.ControlPointID],
			otelconfig.NewOTELParams,
			fx.Annotate(
				provideAgent,
				fx.ResultTags(otelconfig.BaseFxTag),
			),
			fx.Annotate(
				AgentOTELComponents,
				fx.ParamTags(alerts.AlertsFxTag),
			),
		),
	)
}

// AgentOTELComponents constructs OTEL Collector Factories for Agent.
func AgentOTELComponents(
	alerter alerts.Alerter,
	cache *entitycache.EntityCache,
	promRegistry *prometheus.Registry,
	engine iface.Engine,
	clasEng iface.ClassificationEngine,
	serverGRPC *grpc.Server,
	controlPointCache *cache.Cache[selectors.ControlPointID],
	alertMgr *alertmanager.AlertManager,
) (otelcol.Factories, error) {
	var errs error

	extensions, err := extension.MakeFactoryMap(
		zpagesextension.NewFactory(),
		ballastextension.NewFactory(),
		healthcheckextension.NewFactory(),
		pprofextension.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	// We need to create and register empty server wrappers in GRPC server, as OTEL
	// receivers are created after our GRPC server is started.
	// Inside the otlpreceiver the wrappers are filled with proper servers.
	tsw := &otlpreceiver.TraceServerWrapper{}
	msw := &otlpreceiver.MetricServerWrapper{}
	lsw := &otlpreceiver.LogServerWrapper{}
	ptraceotlp.RegisterGRPCServer(serverGRPC, tsw)
	pmetricotlp.RegisterGRPCServer(serverGRPC, msw)
	plogotlp.RegisterGRPCServer(serverGRPC, lsw)

	receiversFactory := []receiver.Factory{
		otlpreceiver.NewFactory(tsw, msw, lsw),
		alertsreceiver.NewFactory(alerter),
	}

	receiversFactory = append(receiversFactory, otelContribReceivers()...)

	receivers, err := receiver.MakeFactoryMap(receiversFactory...)
	errs = multierr.Append(errs, err)

	exporters, err := exporter.MakeFactoryMap(
		fileexporter.NewFactory(),
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		alertsexporter.NewFactory(alertMgr),
	)
	errs = multierr.Append(errs, err)

	processorsFactory := []processor.Factory{
		enrichmentprocessor.NewFactory(cache),
		rollupprocessor.NewFactory(promRegistry),
		metricsprocessor.NewFactory(promRegistry, engine, clasEng, controlPointCache),
		tracestologsprocessor.NewFactory(),
	}
	processorsFactory = append(processorsFactory, otelContribProcessors()...)
	processors, err := processor.MakeFactoryMap(processorsFactory...)
	errs = multierr.Append(errs, err)

	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}

	return factories, errs
}

func provideAgent(cfg *otelconfig.OTELParams) *otelconfig.OTELConfig {
	addLogsPipeline(cfg)
	addTracesPipeline(cfg)
	addMetricsPipeline(cfg)
	otelconfig.AddAlertsPipeline(cfg, otelconsts.ProcessorAgentResourceLabels)
	return cfg.Config
}

func addLogsPipeline(cfg *otelconfig.OTELParams) {
	config := cfg.Config
	// Common dependencies for pipelines
	config.AddReceiver(otelconsts.ReceiverOTLP, otlpreceiver.Config{})
	// Note: Passing map[string]interface{}{} instead of real config, so that
	// processors' configs' default work.
	config.AddProcessor(otelconsts.ProcessorMetrics, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPrerollup,
		cfg.BatchPrerollup.Timeout.AsDuration(),
		cfg.BatchPrerollup.SendBatchSize,
		cfg.BatchPrerollup.SendBatchMaxSize,
	)
	config.AddProcessor(otelconsts.ProcessorRollup, map[string]interface{}{})
	config.AddBatchProcessor(
		otelconsts.ProcessorBatchPostrollup,
		cfg.BatchPostrollup.Timeout.AsDuration(),
		cfg.BatchPostrollup.SendBatchSize,
		cfg.BatchPostrollup.SendBatchMaxSize,
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

func addTracesPipeline(cfg *otelconfig.OTELParams) {
	config := cfg.Config
	config.AddExporter(otelconsts.ExporterOTLPLoopback, map[string]any{
		"endpoint": cfg.Listener.GetAddr(),
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

func addMetricsPipeline(cfg *otelconfig.OTELParams) {
	config := cfg.Config
	addPrometheusReceiver(cfg)
	config.AddProcessor(otelconsts.ProcessorEnrichment, nil)
	otelconfig.AddPrometheusRemoteWriteExporter(config, cfg.PromClient)
	config.Service.AddPipeline("metrics/fast", otelconfig.Pipeline{
		Receivers: []string{otelconsts.ReceiverPrometheus},
		Processors: []string{
			otelconsts.ProcessorEnrichment,
			otelconsts.ProcessorAgentGroup,
		},
		Exporters: []string{otelconsts.ExporterPrometheusRemoteWrite},
	})
}

func addPrometheusReceiver(cfg *otelconfig.OTELParams) {
	config := cfg.Config
	scrapeConfigs := []map[string]any{
		otelconfig.BuildApertureSelfScrapeConfig("aperture-self", cfg),
		otelconfig.BuildOTELScrapeConfig("aperture-otel", cfg),
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

func buildKubernetesNodesScrapeConfig(cfg *otelconfig.OTELParams) map[string]any {
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

func buildKubernetesPodsScrapeConfig(cfg *otelconfig.OTELParams) map[string]any {
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
