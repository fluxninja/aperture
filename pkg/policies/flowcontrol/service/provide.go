package service

import (
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy"
	"go.uber.org/fx"
)

// Module is a set of default providers for flowcontrol components.
func Module() fx.Option {
	return fx.Options(
		check.Module(),
		envoy.Module(),
	)
}
