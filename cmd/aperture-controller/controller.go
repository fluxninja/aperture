//go:generate swagger generate spec --scan-models --include="github.com/FluxNinja*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-controller/config-swagger.yaml

//go:generate swagger generate spec --include="github.com/FluxNinja*" --include-tag=policies-configuration -o ../../docs/gen/policies/config-swagger.yaml

//go:generate swagger generate spec --include="github.com/FluxNinja*" --include-tag=classification-configuration -o ../../docs/gen/classification/config-swagger.yaml

// Aperture Controller
//   BasePath: /aperture-controller
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/cmd/aperture-controller/controller"
	"github.com/FluxNinja/aperture/pkg/classification"
	"github.com/FluxNinja/aperture/pkg/flowcontrol"
	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/otel"
	"github.com/FluxNinja/aperture/pkg/otelcollector"
	"github.com/FluxNinja/aperture/pkg/platform"
	"github.com/FluxNinja/aperture/pkg/webhooks"
	"github.com/FluxNinja/aperture/pkg/webhooks/validation"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		otel.ProvideAnnotatedControllerConfig(),
		fx.Provide(
			clockwork.NewRealClock,
			classification.ProvideCMFileValidator,
			flowcontrol.ProvideCMFileValidator,
			validation.ProvideCMValidator,
			otel.ControllerOTELComponents,
		),
		otelcollector.Module(),
		controller.Module(),
		webhooks.Module(),
		fx.Invoke(
			classification.RegisterCMFileValidator,
			flowcontrol.RegisterCMFileValidator,
			validation.RegisterCMValidator,
		),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Fatal().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-controller app created")
	platform.Run(app)
}
