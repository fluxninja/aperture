package nginxreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading nginxreceiver Receiver")
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideReceiver,
				fx.ResultTags(config.GroupTag(otelcollector.ReceiverFactoriesFxTag)),
			),
		),
	)
}

func provideReceiver() receiver.Factory {
	return nginxreceiver.NewFactory()
}
