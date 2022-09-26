package iface

import (
	"time"

	"github.com/fluxninja/aperture/pkg/status"
)

const (
	// PoliciesRoot - path in config and status registry for policies results.
	PoliciesRoot = "policies"
)

// PolicyBase is for read only access to base policy info.
type PolicyBase interface {
	GetPolicyName() string
	GetPolicyHash() string
}

// Policy is for read only access to full policy state.
type Policy interface {
	PolicyBase
	GetEvaluationInterval() time.Duration
	GetStatusRegistry() status.Registry
}
