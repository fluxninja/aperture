package main

import (
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/fx"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver"

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
	log.Info().Msg("Loading SQL Query Receiver")
	return fx.Options(
		module(),
	)
}

func module() fx.Option {
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
		sqlqueryreceiver.NewFactory(),
	}, nil
}
