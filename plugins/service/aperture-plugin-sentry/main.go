//go:generate swagger generate spec --scan-models --include="github.com/fluxninja/aperture/plugins/service/aperture-plugin-sentry/*" --include-tag=plugin-configuration -o ../../../docs/gen/config/aperture-plugin-sentry/plugin-swagger.yaml

// Sentry Plugin
//   BasePath: /aperture-controller
// swagger:meta

package main

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/plugins"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-sentry/sentry"
)

const (
	Plugin        = "aperture-plugin-sentry"
	BuildHost     = "unknown"
	BuildOS       = "unknown"
	BuildTime     = "unknown"
	GitBranch     = "unknown"
	GitCommitHash = "unknown"
)

func ServicePlugin() plugins.ServicePluginIface {
	return &SentryPlugin{}
}

type SentryPlugin struct{}

func (sp *SentryPlugin) Module() fx.Option {
	log.Info().Msg("Loading SentryPlugin")
	return fx.Options(
		sentry.SentryWriterConstructor{ConfigKey: Plugin}.Annotate(),
	)
}

func main() {
	sp := &SentryPlugin{}
	app := platform.New(
		platform.Config{}.Module(),
		sp.Module(),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("sentry-plugin app created")
	platform.Run(app)
}
