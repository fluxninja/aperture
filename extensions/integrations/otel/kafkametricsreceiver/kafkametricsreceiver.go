package kafkametricsreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading kafkametricsreceiver Receiver")
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
	return kafkametricsreceiver.NewFactory()
}
