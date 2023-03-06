package kubernetes

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/actuators"
	"github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/service"
)

// Module returns the fx options for infra integrations of policy.
func Module() fx.Option {
	return fx.Options(
		actuators.Module(),
		service.Module(),
	)
}
