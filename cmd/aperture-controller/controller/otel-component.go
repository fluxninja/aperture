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

	"github.com/fluxninja/aperture/v2/pkg/alertmanager"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/alertsexporter"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector/alertsreceiver"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

// ModuleForControllerOTel provides fx options for ControllerOTelComponents.
func ModuleForControllerOTel() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideController,
				fx.ResultTags(otelconsts.BaseFxTag),
			),
			fx.Annotate(
				ControllerOTelComponents,
				fx.ParamTags(
					alerts.AlertsFxTag,
					otelconsts.ReceiverFactoriesFxTag,
					otelconsts.ProcessorFactoriesFxTag,
				),
			),
		),
	)
}

// ControllerOTelComponents constructs OTel Collector Factories for Controller.
func ControllerOTelComponents(
	alerter alerts.Alerter,
	receiverFactories []receiver.Factory,
	processorFactories []processor.Factory,
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

	rf := []receiver.Factory{
		prometheusreceiver.NewFactory(),
		alertsreceiver.NewFactory(alerter),
	}
	// receiversFactory = append(receiversFactory, otelContribReceivers()...)
	rf = append(rf, receiverFactories...)
	receivers, err := receiver.MakeFactoryMap(rf...)
	errs = multierr.Append(errs, err)

	ef := []exporter.Factory{
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		loggingexporter.NewFactory(),
		alertsexporter.NewFactory(alertMgr),
	}
	exporters, err := exporter.MakeFactoryMap(ef...)
	errs = multierr.Append(errs, err)

	pf := []processor.Factory{
		batchprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		transformprocessor.NewFactory(),
	}
	// processorsFactory = append(processorsFactory, otelContribProcessors()...)
	pf = append(pf, processorFactories...)
	processors, err := processor.MakeFactoryMap(pf...)
	errs = multierr.Append(errs, err)

	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}

	return factories, errs
}
