package actuators

import (
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/actuators/podscaler"
	"go.uber.org/fx"
)

// Module returns the fx options for infra actuator integrations.
func Module() fx.Option {
	return fx.Options(
		podscaler.Module(),
	)
}
