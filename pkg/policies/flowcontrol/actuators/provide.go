package actuators

import (
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/loadscheduler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/rate"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/regulator"
	"go.uber.org/fx"
)

// Module returns the fx options for flowcontrol side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		loadscheduler.Module(),
		rate.Module(),
		regulator.Module(),
	)
}
