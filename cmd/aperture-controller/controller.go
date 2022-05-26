//go:generate swagger generate spec --scan-models --include="aperture.tech*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-controller/config-swagger.yaml

//go:generate swagger generate spec --include="aperture.tech*" --include-tag=policies-configuration -o ../../docs/gen/policies/config-swagger.yaml

//go:generate swagger generate spec --include="aperture.tech*" --include-tag=classification-configuration -o ../../docs/gen/classification/config-swagger.yaml

// Aperture Controller
//   BasePath: /aperture-controller
// swagger:meta
package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"

	"aperture.tech/aperture/cmd/aperture-controller/controller"
	"aperture.tech/aperture/pkg/classification"
	"aperture.tech/aperture/pkg/flowcontrol"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/otel"
	"aperture.tech/aperture/pkg/otelcollector"
	"aperture.tech/aperture/pkg/platform"
	"aperture.tech/aperture/pkg/webhooks"
	"aperture.tech/aperture/pkg/webhooks/validation"
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
