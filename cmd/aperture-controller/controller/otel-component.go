package controller

import (
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

// ModuleForControllerOTEL provides fx options for ControllerOTELComponents.
func ModuleForControllerOTEL() fx.Option {
	return fx.Options(
		fx.Provide(
			otelcollector.NewOtelConfig,
			fx.Annotate(
				provideController,
				fx.ResultTags(otelcollector.BaseFxTag),
			),
			ControllerOTELComponents,
		),
	)
}

// ControllerOTELComponents constructs OTEL Collector Factories for Controller.
func ControllerOTELComponents() (component.Factories, error) {
	var errs error

	extensions, err := component.MakeExtensionFactoryMap(
		zpagesextension.NewFactory(),
		ballastextension.NewFactory(),
		healthcheckextension.NewFactory(),
		pprofextension.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	receivers, err := component.MakeReceiverFactoryMap(
		prometheusreceiver.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	exporters, err := component.MakeExporterFactoryMap(
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	processors, err := component.MakeProcessorFactoryMap(
		batchprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
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

func provideController(cfg *otelcollector.OtelParams) *otelcollector.OTELConfig {
	otelcollector.AddControllerMetricsPipeline(cfg)
	return cfg.Config
}
