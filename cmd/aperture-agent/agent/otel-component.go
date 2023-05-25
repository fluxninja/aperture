package agent

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/connector"
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

	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/attributesprocessor"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/batchprocessor"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/filelogreceiver"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/filterprocessor"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/k8sattributesprocessor"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/kubeletstatsreceiver"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/memorylimiterprocessor"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/prometheusreceiver"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/resourceprocessor"
	"github.com/fluxninja/aperture/v2/cmd/aperture-agent/agent/otel/transformprocessor"
	alertmanager "github.com/fluxninja/aperture/v2/pkg/alert-manager"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/cache"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/adapterconnector"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/alertsexporter"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/alertsreceiver"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/leaderonlyreceiver"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/metricsprocessor"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/rollupprocessor"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
)

// ModuleForAgentOTel provides fx options for AgentOTelComponent.
func ModuleForAgentOTel() fx.Option {
	return fx.Options(
		kubeletstatsreceiver.Module(),
		k8sattributesprocessor.Module(),
		prometheusreceiver.Module(),
		filelogreceiver.Module(),
		leaderonlyreceiver.Module(),
		batchprocessor.Module(),
		memorylimiterprocessor.Module(),
		attributesprocessor.Module(),
		transformprocessor.Module(),
		resourceprocessor.Module(),
		filterprocessor.Module(),
		fx.Provide(
			cache.Provide[selectors.TypedControlPointID],
			provideAgent,
			fx.Annotate(
				AgentOTelComponents,
				fx.ParamTags(
					alerts.AlertsFxTag,
					otelconsts.ReceiverFactoriesFxTag,
					otelconsts.ProcessorFactoriesFxTag,
				),
			),
		),
	)
}

// AgentOTelComponents constructs OTel Collector Factories for Agent.
func AgentOTelComponents(
	alerter alerts.Alerter,
	receiverFactories []receiver.Factory,
	processorFactories []processor.Factory,
	promRegistry *prometheus.Registry,
	engine iface.Engine,
	clasEng iface.ClassificationEngine,
	serverGRPC *grpc.Server,
	controlPointCache *cache.Cache[selectors.TypedControlPointID],
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

	// We need to create and register empty server wrappers in gRPC server, as OTel
	// receivers are created after our gRPC server is started.
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
	}
	// processorsFactory = append(processorsFactory, otelContribProcessors()...)
	pf = append(pf, processorFactories...)
	processors, err := processor.MakeFactoryMap(pf...)
	errs = multierr.Append(errs, err)

	cf := []connector.Factory{
		adapterconnector.NewFactory(),
	}
	connectors, err := connector.MakeFactoryMap(cf...)
	errs = multierr.Append(errs, err)

	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
		Connectors: connectors,
	}

	return factories, errs
}
