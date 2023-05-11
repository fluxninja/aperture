package servicegraphprocessor

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/servicegraphprocessor"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading servicegraphprocessor Processor")
	if info.Service != utils.ApertureAgent {
		return nil
	}
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideProcessor,
				fx.ResultTags(otelconsts.ProcessorFactoriesFxTag),
			),
		),
	)
}

func provideProcessor() processor.Factory {
	return servicegraphprocessor.NewFactory()
}
