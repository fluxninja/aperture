package controller

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	"github.com/fluxninja/aperture/pkg/alertmanager"
	"github.com/fluxninja/aperture/pkg/alerts"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsexporter"
	"github.com/fluxninja/aperture/pkg/otelcollector/alertsreceiver"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
)

// ModuleForControllerOTEL provides fx options for ControllerOTELComponents.
func ModuleForControllerOTEL() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideController,
				fx.ResultTags(otelconfig.BaseFxTag),
			),
			fx.Annotate(
				ControllerOTELComponents,
				fx.ParamTags(alerts.AlertsFxTag),
			),
		),
	)
}

// ControllerOTELComponents constructs OTEL Collector Factories for Controller.
func ControllerOTELComponents(
	alerter alerts.Alerter,
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

	receivers, err := receiver.MakeFactoryMap(
		prometheusreceiver.NewFactory(),
		alertsreceiver.NewFactory(alerter),
	)
	errs = multierr.Append(errs, err)

	exporters, err := exporter.MakeFactoryMap(
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		loggingexporter.NewFactory(),
		alertsexporter.NewFactory(alertMgr),
	)
	errs = multierr.Append(errs, err)

	processors, err := processor.MakeFactoryMap(
		batchprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		transformprocessor.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}

	return factories, errs
}
