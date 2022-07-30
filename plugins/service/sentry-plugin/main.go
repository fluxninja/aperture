//go:generate swagger generate spec --scan-models --include="github.com/fluxninja/aperture/plugins/service/sentry-plugin/*" --include-tag=plugin-configuration -o ../../../docs/gen/config/sentry-plugin/plugin-swagger.yaml

// Sentry Plugin
//   BasePath: /aperture-controller
// swagger:meta

package main

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/plugins"
	"github.com/fluxninja/aperture/plugins/service/sentry-plugin/sentry"
)

const (
	Plugin        = "sentry-plugin"
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
		sentry.SentryWriterConstructor{Key: Plugin}.Annotate(),
		fx.Invoke(sentry.RegisterSentryPanicHandler),
	)
}
