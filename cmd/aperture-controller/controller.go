//go:generate swagger generate spec --scan-models --include="github.com/fluxninja*" --include-tag=common-configuration -o ../../docs/gen/config/aperture-controller/config-swagger.yaml

//go:generate swagger generate spec --include="github.com/fluxninja*" --include-tag=policies-configuration -o ../../docs/gen/policies/config-swagger.yaml

//go:generate swagger generate spec --include="github.com/fluxninja*" --include-tag=classification-configuration -o ../../docs/gen/classification/config-swagger.yaml

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
	"github.com/fluxninja/aperture/pkg/classification"
	"github.com/fluxninja/aperture/pkg/flowcontrol"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otel"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/webhooks"
	"github.com/fluxninja/aperture/pkg/webhooks/validation"
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
		log.Panic().Err(err).Msg("fx.New failed: " + visualize)
	}

	log.Info().Msg("aperture-controller app created")
	platform.Run(app)
}
