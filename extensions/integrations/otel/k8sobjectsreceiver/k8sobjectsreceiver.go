package k8sobjectsreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/pkg/utils"
)

// Module returns the extension module.
func Module() fx.Option {
	log.Info().Msg("Loading k8sobjectsreceiver Receiver")
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
	return k8sobjectsreceiver.NewFactory()
}
