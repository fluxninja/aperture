package runtime

import "aperture.tech/aperture/pkg/policies/apis/policyapi"

// TickEndCallback is a function that is called when a tick ends.
type TickEndCallback func(tickInfo TickInfo) error

// CircuitAPI is for read only access to policy and also provides methods for acquiring & releasing circuit execution lock.
type CircuitAPI interface {
	policyapi.PolicyReadAPI
	RegisterTickEndCallback(cb TickEndCallback)
	LockExecution()
	UnlockExecution()
}
