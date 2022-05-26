//go:generate swagger generate spec --scan-models --include="aperture.tech/aperture/plugins/*" --include-tag=plugin-configuration -o ../../../docs/gen/config/aperture-plugin-fluxninja/plugin-swagger.yaml

// FluxNinja Cloud Plugin
//   BasePath: /aperture-controller
// swagger:meta

package main

import (
	"go.uber.org/fx"

	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/plugins"
	"aperture.tech/aperture/plugins/service/aperture-plugin-fluxninja/heartbeats"
	"aperture.tech/aperture/plugins/service/aperture-plugin-fluxninja/otel"
	"aperture.tech/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
	"aperture.tech/aperture/plugins/service/aperture-plugin-fluxninja/sentry"
)

// Set via ldflags.
var (
	Plugin        = "aperture-plugin-fluxninja"
	BuildHost     = "unknown"
	BuildOS       = "unknown"
	BuildTime     = "unknown"
	GitBranch     = "unknown"
	GitCommitHash = "unknown"
)

func ServicePlugin() plugins.ServicePluginIface {
	return &FluxNinjaPlugin{}
}

type FluxNinjaPlugin struct{}

func (fn *FluxNinjaPlugin) Module() fx.Option {
	log.Info().Msg("Loading FluxNinjaPlugin")
	return fx.Options(
		sentry.SentryWriterConstructor{Key: pluginconfig.PluginConfigKey + "." + sentry.SentryConfigKey}.Annotate(),
		heartbeats.Module(),
		otel.ProvideAnnotatedPluginConfig(),
	)
}
