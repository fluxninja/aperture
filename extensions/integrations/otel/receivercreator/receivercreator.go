package receivercreator

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading receivercreator Receiver")
	if info.Service != utils.ApertureAgent {
		return nil
	}
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideReceiver,
				fx.ResultTags(otelconsts.ReceiverFactoriesFxTag),
			),
		),
	)
}

func provideReceiver() receiver.Factory {
	return receivercreator.NewFactory()
}
