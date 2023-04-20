package service

import (
	"github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/service/controlpoints"
	"go.uber.org/fx"
)

// Module is a set of default providers for flowcontrol components.
func Module() fx.Option {
	return fx.Options(
		controlpoints.Module(),
	)
}
