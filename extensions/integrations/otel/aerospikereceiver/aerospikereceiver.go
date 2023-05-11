package aerospikereceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading aerospikereceiver Receiver")
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
	return aerospikereceiver.NewFactory()
}
