package main

import (
	"github.com/fluxninja/aperture/extensions/fluxninja"
	"github.com/fluxninja/aperture/extensions/sentry"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fluxninja.Module(),
		sentry.Module(),
	)
}

func GetExtensions() []string {
	return []string{
		"github.com/fluxninja/aperture/extensions/fluxninja",
		"github.com/fluxninja/aperture/extensions/sentry",
	}
}
