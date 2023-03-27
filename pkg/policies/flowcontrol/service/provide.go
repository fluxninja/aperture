package service

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/check"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/checkhttp"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/controlpoints"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/envoy"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/preview"
)

// Module is a set of default providers for flowcontrol components.
func Module() fx.Option {
	return fx.Options(
		check.Module(),
		checkhttp.Module(),
		envoy.Module(),
		preview.Module(),
		controlpoints.Module(),
	)
}
