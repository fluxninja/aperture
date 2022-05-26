package main

// Minimal App for experimenting with Classifier,
//
// Run:
// go run cmd/class/class.go --config-path .
//
// class.yaml:
// classification:
//   rules_dir: "/tmp/foo"
//
// webhooks:
//   certs_dir: certs
//   addr: ":8089"
//
// log:
//   pretty_console: true

import (
	"fmt"

	"go.uber.org/fx"

	"aperture.tech/aperture/pkg/authz"
	"aperture.tech/aperture/pkg/classification"
	"aperture.tech/aperture/pkg/entitycache"
	"aperture.tech/aperture/pkg/flowcontrol"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/platform"
	"aperture.tech/aperture/pkg/policies/dataplane"
	"aperture.tech/aperture/pkg/webhooks"
	"aperture.tech/aperture/pkg/webhooks/validation"
)

func main() {
	app := platform.New(
		platform.Config{}.Module(),
		fx.Supply(dataplane.Engine{}),
		fx.Provide(
			flowcontrol.ProvideNopMetrics,
			flowcontrol.ProvideDummyHandler, // stub – empty
			entitycache.NewEntityCache,      // stub – empty
		),
		webhooks.Module(),
		classification.Module,
		authz.Module,
		fx.Invoke(
			authz.Register,
			classification.RegisterCMFileValidator,
			flowcontrol.RegisterCMFileValidator,
			validation.RegisterCMValidator,
		),
	)

	if err := app.Err(); err != nil {
		visualize, _ := fx.VisualizeError(err)
		log.Fatal().Err(err).Msg("fx.New failed: " + visualize)
		fmt.Println(err)
	}

	log.Info().Msg("App Created")
	log.Info().Msg("Starting app")
	platform.Run(app)
}
