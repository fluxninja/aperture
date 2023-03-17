package metricsgenerationprocessor

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading metricsgenerationprocessor Processor")
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideProcessor,
				fx.ResultTags(config.GroupTag(otelcollector.ProcessorFactoriesFxTag)),
			),
		),
	)
}

func provideProcessor() processor.Factory {
	return metricsgenerationprocessor.NewFactory()
}
