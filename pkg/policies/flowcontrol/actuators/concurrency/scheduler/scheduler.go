package scheduler

import (
	"time"
)

// RequestContext is metadata for request in a flow that is to be allowed or dropped based on controlled delay and queue limits.
type RequestContext struct {
	FairnessLabel string // for enforcing fairness
	Tokens        uint64 // tokens (e.g. expected latency or complexity) for this request
	Priority      uint8  // larger values represent higher priority
	Timeout       time.Duration
}

// Scheduler : Interface for schedulers.
type Scheduler interface {
	// Schedule sends RequestContext to the underlying scheduler and returns a boolean value,
	// where true means accept and false means reject.
	Schedule(rContext RequestContext) bool
}

// TokenManager : Interface for token managers.
type TokenManager interface {
	// Take tokens if available, otherwise return false
	TakeIfAvailable(now time.Time, tokens float64) bool
	// Take tokens even if available tokens are less than asked - returns wait time if tokens are not available immediately. The other return value conveys whether the operation was successful or not.
	Take(now time.Time, timeout time.Duration, tokens float64) (time.Duration, bool)
	// Provides TokenManager the request that the scheduler processing -- some TokenManager implementations use this level of visibility for their algorithms. Return value decides whether the request has to be accepted right away in case TokenManger is not yet ready or configured to accept all traffic (short circuit).
	PreprocessRequest(now time.Time, rContext RequestContext) bool
}
