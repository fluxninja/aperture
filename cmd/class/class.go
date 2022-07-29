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

	"github.com/fluxninja/aperture/pkg/authz"
	"github.com/fluxninja/aperture/pkg/classification"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/flowcontrol"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/platform"
	"github.com/fluxninja/aperture/pkg/policies/dataplane"
	"github.com/fluxninja/aperture/pkg/webhooks"
	"github.com/fluxninja/aperture/pkg/webhooks/validation"
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
