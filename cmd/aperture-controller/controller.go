//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration -o ../../docs/gen/config/controller/config-swagger.yaml

// Package main Controller
//
// Aperture Controller
//
//	BasePath: /aperture-controller
//
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/cmd/aperture-controller/controller"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/webhooks"
	"github.com/fluxninja/aperture/pkg/webhooks/policyvalidator"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		controller.ModuleForControllerOTEL(),
		fx.Provide(
			clockwork.NewRealClock,
		),
		otelcollector.Module(),
		controlplane.Module(),
		webhooks.Module(),
		policyvalidator.Module(),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-controller app created")
	platform.Run(app)
}
