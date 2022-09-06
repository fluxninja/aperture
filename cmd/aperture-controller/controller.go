//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-controller/config-swagger.yaml

// Aperture Controller
//
//	BasePath: /aperture-controller
//
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otel"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/webhooks"
	"github.com/fluxninja/aperture/pkg/webhooks/validation"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		otel.OTELConfigConstructor{Type: otel.ControllerType}.Annotate(),
		fx.Provide(
			clockwork.NewRealClock,
			otel.ControllerOTELComponents,
		),
		otelcollector.Module(),
		controlplane.Module(),
		webhooks.Module(),
		validation.Module(),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-controller app created")
	platform.Run(app)
}
