package actuators

import (
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/concurrency"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/flowregulator"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/actuators/rate"
	"go.uber.org/fx"
)

// Module returns the fx options for flowcontrol side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		concurrency.Module(),
		rate.Module(),
		flowregulator.Module(),
	)
}
