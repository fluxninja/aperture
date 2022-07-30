//go:generate swagger generate spec --scan-models --include="github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/*" --include-tag=plugin-configuration -o ../../../docs/gen/config/aperture-plugin-fluxninja/plugin-swagger.yaml

// FluxNinja Cloud Plugin
//   BasePath: /aperture-controller
// swagger:meta

package main

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/plugins"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/heartbeats"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/otel"
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
		heartbeats.Module(),
		otel.ProvideAnnotatedPluginConfig(),
	)
}
