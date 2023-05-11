//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration --include-tag=controller-configuration -o ../../docs/gen/config/controller/config-swagger.yaml
//go:generate go run ../../docs/tools/swagger/process-go-tags.go ../../docs/gen/config/controller/config-swagger.yaml

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

	"github.com/fluxninja/aperture/v2/cmd/aperture-controller/controller"
	"github.com/fluxninja/aperture/v2/pkg/agentfunctions/agents"
	"github.com/fluxninja/aperture/v2/pkg/cmd"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/otelcollector"
	"github.com/fluxninja/aperture/v2/pkg/platform"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/v2/pkg/rpc"
	"github.com/fluxninja/aperture/v2/pkg/webhooks"
	"github.com/fluxninja/aperture/v2/pkg/webhooks/policyvalidator"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		controller.ModuleForControllerOTel(),
		fx.Provide(
			clockwork.NewRealClock,
		),
		otelcollector.Module(),
		controlplane.Module(),
		webhooks.Module(),
		policyvalidator.Module(),
		rpc.ServerModule,
		agents.Module,
		cmd.Module,
		Module(),
	)

	if err := app.Err(); err != nil {
		visualize, viserr := fx.VisualizeError(err)
		if viserr != nil {
			log.Panic().Err(viserr).Msgf("Failed to visualize fx error: %s", viserr)
		}
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-controller app created")

	platform.Run(app)
}
