package actuators

import (
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/loadregulator"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/loadscheduler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/rate"
	"go.uber.org/fx"
)

// Module returns the fx options for flowcontrol side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		loadscheduler.Module(),
		rate.Module(),
		loadregulator.Module(),
	)
}
