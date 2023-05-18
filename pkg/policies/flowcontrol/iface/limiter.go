package iface

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
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
	Decide(ctx context.Context, labels map[string]string) *flowcontrolv1.LimiterDecision
	Revert(labels map[string]string, decision *flowcontrolv1.LimiterDecision)
	GetLimiterID() LimiterID
	GetRequestCounter(labels map[string]string) prometheus.Counter
}

// RateLimiter interface.
type RateLimiter interface {
	Limiter
	TakeN(labels map[string]string, count float64) (label string, ok bool, remaining float64, current float64)
}

// LoadScheduler interface.
type LoadScheduler interface {
	Limiter
	GetLatencyObserver(labels map[string]string) prometheus.Observer
}
