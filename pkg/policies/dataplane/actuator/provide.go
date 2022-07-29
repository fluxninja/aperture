package actuator

import (
	"github.com/FluxNinja/aperture/pkg/policies/dataplane/actuator/concurrency"
	"github.com/FluxNinja/aperture/pkg/policies/dataplane/actuator/rate"
	"go.uber.org/fx"
)

// Module returns the fx options for dataplane side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		concurrency.Module(),
		rate.Module(),
	)
}
