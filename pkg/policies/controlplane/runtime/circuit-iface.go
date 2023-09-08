package runtime

import "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"

// TickEndCallback is a function that is called when a tick ends.
type TickEndCallback func(circuit CircuitAPI) error

// TickStartCallback is a function that is called when a tick starts.
type TickStartCallback func(circuit CircuitAPI) error

// CircuitAPI is for read only access to policy and scheduling of background jobs.
type CircuitAPI interface {
	iface.Policy
	ScheduleConditionalBackgroundJob(backgroundJob BackgroundJob, ticksPerExecution int)
	GetTickInfo() TickInfo
}

// CircuitSuperAPI is for read only access to policy, scheduling of background jobs and also provides methods for acquiring & releasing circuit execution lock.
type CircuitSuperAPI interface {
	CircuitAPI
	RegisterTickEndCallback(ec TickEndCallback)
	RegisterTickStartCallback(sc TickStartCallback)
	LockExecution()
	UnlockExecution()
}
