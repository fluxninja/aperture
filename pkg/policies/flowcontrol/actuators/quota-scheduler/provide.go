package quotascheduler

import "go.uber.org/fx"

// Module returns the fx options for dataplane side pieces of rate limiter.
func Module() fx.Option {
	return fx.Options(
		quotaSchedulerModule(),
	)
}
