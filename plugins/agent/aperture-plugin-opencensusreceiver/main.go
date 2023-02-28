package main

import (
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/plugins"
)

// ServicePlugin returns the plugin.
func ServicePlugin() plugins.ServicePluginIface {
	return &Plugin{}
}

// Plugin implements the plugin interface.
type Plugin struct{}

// Module returns the plugin module.
func (p *Plugin) Module() fx.Option {
	log.Info().Msg("Loading OpenCensus Receiver")
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideReceiver,
				fx.ResultTags(config.GroupTag(otelcollector.ReceiverFactoriesFxTag)),
			),
		),
	)
}

func provideReceiver() ([]receiver.Factory, error) {
	return []receiver.Factory{
		opencensusreceiver.NewFactory(),
	}, nil
}
