//go:generate swagger generate spec --scan-models --include="github.com/fluxninja/aperture/v2/extensions/fluxninja/*" --include-tag=extension-configuration -o ../../docs/gen/config/extensions/fluxninja/extension-swagger.yaml
//go:generate go run ../../docs/tools/swagger/process-go-tags.go ../../docs/gen/config/extensions/fluxninja/extension-swagger.yaml

// FluxNinja ARC Extension
//   BasePath: /aperture-controller
// swagger:meta

package fluxninja

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/heartbeats"
	"github.com/fluxninja/aperture/v2/extensions/fluxninja/otel"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Module returns the FluxNinja extension module for the platform.
func Module() fx.Option {
	log.Info().Msg("Loading FluxNinjaExtension")
	return fx.Options(
		heartbeats.Module(),
		otel.Module(),
		extconfig.Module(),
	)
}
