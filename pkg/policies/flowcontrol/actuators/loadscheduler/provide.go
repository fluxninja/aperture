package loadscheduler

import (
	"go.uber.org/fx"
)

// Module returns the fx options for dataplane side pieces of concurrency control.
func Module() fx.Option {
	return fx.Options(
		loadSchedulerModule(),
	)
}
