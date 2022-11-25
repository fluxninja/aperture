package agent

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/grpc"

	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/controlpointcache"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsreceiver"
	"github.com/fluxninja/aperture/pkg/otelcollector/enrichmentprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/metricsprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/rollupprocessor"
	"github.com/fluxninja/aperture/pkg/otelcollector/tracestologsprocessor"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
)

// ModuleForAgentOTEL provides fx options for AgentOTELComponent.
func ModuleForAgentOTEL() fx.Option {
	return fx.Options(
		fx.Provide(
			otelcollector.NewOtelConfig,
			fx.Annotate(
				provideAgent,
				fx.ResultTags(otelcollector.BaseFxTag),
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
	controlPointCache *controlpointcache.ControlPointCache,
) (component.Factories, error) {
	var errs error

	extensions, err := component.MakeExtensionFactoryMap(
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

	receivers, err := component.MakeReceiverFactoryMap(
		otlpreceiver.NewFactory(tsw, msw, lsw),
		prometheusreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		alertsreceiver.NewFactory(alerter),
	)
	errs = multierr.Append(errs, err)

	exporters, err := component.MakeExporterFactoryMap(
		fileexporter.NewFactory(),
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	processors, err := component.MakeProcessorFactoryMap(
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
		enrichmentprocessor.NewFactory(cache),
		rollupprocessor.NewFactory(),
		metricsprocessor.NewFactory(promRegistry, engine, clasEng, controlPointCache),
		attributesprocessor.NewFactory(),
		tracestologsprocessor.NewFactory(),
		transformprocessor.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	factories := component.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}

	return factories, errs
}

func provideAgent(cfg *otelcollector.OtelParams) *otelcollector.OTELConfig {
	addLogsPipeline(cfg)
	addTracesPipeline(cfg)
	otelcollector.AddMetricsPipeline(cfg)
	otelcollector.AddAlertsPipeline(cfg, otelcollector.ProcessorAgentResourceLabels)
	return cfg.Config
}

func addLogsPipeline(cfg *otelcollector.OtelParams) {
	config := cfg.Config
	// Common dependencies for pipelines
	config.AddReceiver(otelcollector.ReceiverOTLP, otlpreceiver.Config{})
	// Note: Passing map[string]interface{}{} instead of real config, so that
	// processors' configs' default work.
	config.AddProcessor(otelcollector.ProcessorMetrics, map[string]interface{}{})
	config.AddBatchProcessor(
		otelcollector.ProcessorBatchPrerollup,
		cfg.BatchPrerollup.Timeout.AsDuration(),
		cfg.BatchPrerollup.SendBatchSize,
		cfg.BatchPrerollup.SendBatchMaxSize,
	)
	config.AddProcessor(otelcollector.ProcessorRollup, map[string]interface{}{})
	config.AddBatchProcessor(
		otelcollector.ProcessorBatchPostrollup,
		cfg.BatchPostrollup.Timeout.AsDuration(),
		cfg.BatchPostrollup.SendBatchSize,
		cfg.BatchPostrollup.SendBatchMaxSize,
	)
	config.AddExporter(otelcollector.ExporterLogging, nil)

	processors := []string{
		otelcollector.ProcessorMetrics,
		otelcollector.ProcessorBatchPrerollup,
		otelcollector.ProcessorRollup,
		otelcollector.ProcessorBatchPostrollup,
		otelcollector.ProcessorAgentGroup,
	}

	config.Service.AddPipeline("logs", otelcollector.Pipeline{
		Receivers:  []string{otelcollector.ReceiverOTLP},
		Processors: processors,
		Exporters:  []string{otelcollector.ExporterLogging},
	})
}

func addTracesPipeline(cfg *otelcollector.OtelParams) {
	config := cfg.Config
	config.AddExporter(otelcollector.ExporterOTLPLoopback, map[string]any{
		"endpoint": cfg.Listener.GetAddr(),
		"tls": map[string]any{
			"insecure": true,
		},
	})
	config.AddProcessor(otelcollector.ProcessorTracesToLogs, tracestologsprocessor.Config{
		LogsExporter: otelcollector.ExporterOTLPLoopback,
	})

	config.Service.AddPipeline("traces", otelcollector.Pipeline{
		Receivers:  []string{otelcollector.ReceiverOTLP},
		Processors: []string{otelcollector.ProcessorTracesToLogs},
		// We need some exporter configured to make this pipeline correct. Actual
		// Log exporting is done inside the processor.
		Exporters: []string{otelcollector.ExporterLogging},
	})

	// TODO This receiver should be replaced with some receiver which really does nothing.
	config.AddReceiver("filelog", map[string]any{
		"include":       []string{"/var/log/myservice/*.json"},
		"poll_interval": "1000h",
	})
	// We need a fake log pipeline which will initialize the ExporterOTLPLoopback
	// for logs type.
	config.Service.AddPipeline("logs/fake", otelcollector.Pipeline{
		Receivers: []string{"filelog"},
		Exporters: []string{otelcollector.ExporterOTLPLoopback},
	})
}
