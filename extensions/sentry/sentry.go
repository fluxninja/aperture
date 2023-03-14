//go:generate swagger generate spec --scan-models --include="github.com/fluxninja/aperture/extensions/sentry/*" --include-tag=extension-configuration -o ../../docs/gen/config/extensions/sentry/extension-swagger.yaml
//go:generate go run ../../docs/tools/swagger/process-go-tags.go ../../docs/gen/config/extensions/sentry/extension-swagger.yaml

// Sentry Extension
//   BasePath: /aperture-controller
// swagger:meta

package sentry

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/log"
)

// Module returns the Sentry extension module for the platform.
func Module() fx.Option {
	log.Info().Msg("Loading Sentry Extension")
	constructor := &sentryWriterConstructor{ConfigKey: "sentry"}
	return fx.Options(
		constructor.annotate(),
	)
}
