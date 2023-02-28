package agent

import (
	"fmt"

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

	"github.com/fluxninja/aperture/pkg/alertmanager"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsexporter"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsreceiver"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
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
			fx.Annotate(
				provideAgent,
				fx.ResultTags(otelconfig.BaseFxTag),
			),
			fx.Annotate(
				AgentOTELComponents,
				fx.ParamTags(
					alerts.AlertsFxTag,
					fmt.Sprintf("%s,optional:\"true\"", config.NameTag(otelcollector.ReceiverFactoriesFxTag)),
					fmt.Sprintf("%s,optional:\"true\"", config.NameTag(otelcollector.ProcessorFactoriesFxTag)),
				),
			),
		),
	)
}

// AgentOTELComponents constructs OTEL Collector Factories for Agent.
func AgentOTELComponents(
	alerter alerts.Alerter,
	receiverFactories []receiver.Factory,
	processorFactories []processor.Factory,
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

	rf := []receiver.Factory{
		otlpreceiver.NewFactory(tsw, msw, lsw),
		alertsreceiver.NewFactory(alerter),
	}
	// receiversFactory = append(receiversFactory, otelContribReceivers()...)
	rf = append(rf, receiverFactories...)
	receivers, err := receiver.MakeFactoryMap(rf...)
	errs = multierr.Append(errs, err)

	ef := []exporter.Factory{
		fileexporter.NewFactory(),
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		alertsexporter.NewFactory(alertMgr),
	}
	exporters, err := exporter.MakeFactoryMap(ef...)
	errs = multierr.Append(errs, err)

	pf := []processor.Factory{
		rollupprocessor.NewFactory(promRegistry),
		metricsprocessor.NewFactory(promRegistry, engine, clasEng, controlPointCache),
		tracestologsprocessor.NewFactory(),
	}
	// processorsFactory = append(processorsFactory, otelContribProcessors()...)
	pf = append(pf, processorFactories...)
	processors, err := processor.MakeFactoryMap(pf...)
	errs = multierr.Append(errs, err)

	return otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}, errs
}
