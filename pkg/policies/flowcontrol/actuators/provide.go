package actuators

import (
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/loadscheduler"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/regulator"
	"go.uber.org/fx"
)

// Module returns the fx options for flowcontrol side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		loadscheduler.Module(),
		ratelimiter.Module(),
		regulator.Module(),
	)
}
