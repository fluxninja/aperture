package actuators

import (
	"go.uber.org/fx"

	concurrencylimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/concurrency-limiter"
	concurrencyscheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/concurrency-scheduler"
	loadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/load-scheduler"
	quotascheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/quota-scheduler"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/sampler"
	workloadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/actuators/workload-scheduler"
)

// Module returns the fx options for flowcontrol side pieces of actuator.
func Module() fx.Option {
	return fx.Options(
		workloadscheduler.Module(),
		loadscheduler.Module(),
		ratelimiter.Module(),
		sampler.Module(),
		quotascheduler.Module(),
		concurrencylimiter.Module(),
		concurrencyscheduler.Module(),
	)
}
