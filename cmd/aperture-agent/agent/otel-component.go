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
			AgentOTelComponents,
		),
	)
}

// AgentOTelComponentsIn bundles and annotates parameters.
type AgentOTelComponentsIn struct {
	fx.In
	Alerter            alerts.Alerter      `name:"AlertsFx"`
	ReceiverFactories  []receiver.Factory  `group:"otel-collector-receiver-factories"`
	ProcessorFactories []processor.Factory `group:"otel-collector-processor-factories"`
	PromRegistry       *prometheus.Registry
	Engine             iface.Engine
	ClasEng            iface.ClassificationEngine
	ServerGRPC         *grpc.Server `name:"default"`
	ControlPointCache  *cache.Cache[selectors.TypedControlPointID]
	AlertMgr           *alertmanager.AlertManager
}

// AgentOTelComponents constructs OTel Collector Factories for Agent.
func AgentOTelComponents(in AgentOTelComponentsIn) (otelcol.Factories, error) {
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
	ptraceotlp.RegisterGRPCServer(in.ServerGRPC, tsw)
	pmetricotlp.RegisterGRPCServer(in.ServerGRPC, msw)
	plogotlp.RegisterGRPCServer(in.ServerGRPC, lsw)

	rf := []receiver.Factory{
		otlpreceiver.NewFactory(tsw, msw, lsw),
		alertsreceiver.NewFactory(in.Alerter),
	}
	// receiversFactory = append(receiversFactory, otelContribReceivers()...)
	rf = append(rf, in.ReceiverFactories...)
	receivers, err := receiver.MakeFactoryMap(rf...)
	errs = multierr.Append(errs, err)

	ef := []exporter.Factory{
		fileexporter.NewFactory(),
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		alertsexporter.NewFactory(in.AlertMgr),
	}
	exporters, err := exporter.MakeFactoryMap(ef...)
	errs = multierr.Append(errs, err)

	pf := []processor.Factory{
		rollupprocessor.NewFactory(in.PromRegistry),
		metricsprocessor.NewFactory(in.PromRegistry, in.Engine, in.ClasEng, in.ControlPointCache),
	}
	// processorsFactory = append(processorsFactory, otelContribProcessors()...)
	pf = append(pf, in.ProcessorFactories...)
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
