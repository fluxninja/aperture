package batchprocessor

import (
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading batchprocessor Processor")
	if info.Service != utils.ApertureAgent {
		return nil
	}
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
	return batchprocessor.NewFactory()
}
