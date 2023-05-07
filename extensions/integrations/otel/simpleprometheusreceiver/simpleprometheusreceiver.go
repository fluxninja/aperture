package simpleprometheusreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading simpleprometheusreceiver Receiver")
	if info.Service != utils.ApertureAgent {
		return nil
	}
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideReceiver,
				fx.ResultTags(config.GroupTag(otelconsts.ReceiverFactoriesFxTag)),
			),
		),
	)
}

func provideReceiver() receiver.Factory {
	return simpleprometheusreceiver.NewFactory()
}
