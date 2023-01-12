package actuators

import (
	"github.com/fluxninja/aperture/pkg/policies/infra/actuators/podautoscaler"
	"go.uber.org/fx"
)

// Module returns the fx options for infra actuator integrations.
func Module() fx.Option {
	return fx.Options(
		podautoscaler.Module(),
	)
}
