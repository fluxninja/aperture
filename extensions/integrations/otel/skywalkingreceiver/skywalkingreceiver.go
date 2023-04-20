package skywalkingreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading skywalkingreceiver Receiver")
	if info.Service != utils.ApertureAgent {
		return nil
	}
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
	return skywalkingreceiver.NewFactory()
}
