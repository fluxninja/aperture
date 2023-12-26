package iface

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/labels"
)

//go:generate mockgen -source=limiter.go -destination=../../mocks/mock_limiter.go -package=mocks

// LimiterID is the ID of the Limiter.
type LimiterID struct {
	PolicyName  string
	PolicyHash  string
	ComponentID string
}

// String function returns the LimiterID as a string.
func (limiterID LimiterID) String() string {
	return "policy_name-" + limiterID.PolicyName + "-component_id-" + limiterID.ComponentID + "-policy_hash-" + limiterID.PolicyHash
}

// Limiter interface.
// Lifetime of this interface is per policy/component.
type Limiter interface {
	GetPolicyName() string
	GetSelectors() []*policylangv1.Selector
	Decide(context.Context, labels.Labels) *flowcontrolv1.LimiterDecision
	Revert(context.Context, labels.Labels, *flowcontrolv1.LimiterDecision)
	GetLimiterID() LimiterID
	GetRequestCounter(labels map[string]string) prometheus.Counter
	GetRampMode() bool
}

// Scheduler interface.
type Scheduler interface {
	Limiter
	GetLatencyObserver(labels map[string]string) prometheus.Observer
}

// FlowEnder interface.
type FlowEnder interface {
	Return(ctx context.Context, label string, tokens float64, requestID string) (bool, error)
}

// ConcurrencyLimiter interface.
type ConcurrencyLimiter interface {
	Limiter
	FlowEnder
}

// ConcurrencyScheduler interface.
type ConcurrencyScheduler interface {
	Scheduler
	FlowEnder
}
