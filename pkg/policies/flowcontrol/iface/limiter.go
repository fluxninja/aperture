package iface

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
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
}

// RateLimiter interface.
type RateLimiter interface {
	Limiter
	TakeIfAvailable(ctx context.Context, labels labels.Labels, count float64) (label string, ok bool, waitTime time.Duration, remaining float64, current float64)
}

// Scheduler interface.
type Scheduler interface {
	Limiter
	GetLatencyObserver(labels map[string]string) prometheus.Observer
}
