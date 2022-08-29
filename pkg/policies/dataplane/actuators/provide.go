package actuators

import (
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/concurrency"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/actuators/rate"
	"go.uber.org/fx"
)

// Module returns the fx options for dataplane side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		concurrency.Module(),
		rate.Module(),
	)
}
