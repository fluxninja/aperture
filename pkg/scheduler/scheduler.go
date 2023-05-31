package scheduler

import (
	"context"
	"time"
)

// Request is metadata for request in a flow that is to be allowed or dropped based on controlled delay and queue limits.
type Request struct {
	FairnessLabel string // for enforcing fairness
	Tokens        uint64 // tokens (e.g. expected latency or complexity) for this request
	Priority      uint8  // larger values represent higher priority
}

// Scheduler : Interface for schedulers.
type Scheduler interface {
	// Schedule sends RequestContext to the underlying scheduler and returns a boolean value,
	// where true means accept and false means reject.
	Schedule(ctx context.Context, request Request) bool
	// Revert "unschedules" a request.
	// Useful in case the request was rejected by any
	// other scheduler and the tokens are returned
	// back to the scheduler.
	Revert(ctx context.Context, tokens uint64)
	// Info returns the last access time and number of requests that are currently in the queue.
	Info() (time.Time, int)
}

// TokenManager : Interface for token managers.
type TokenManager interface {
	// Take tokens if available, otherwise return false
	TakeIfAvailable(ctx context.Context, tokens float64) bool
	// Take tokens even if available tokens are less than asked - returns wait time if tokens are not available immediately. The other return value conveys whether the operation was successful or not.
	Take(ctx context.Context, tokens float64) (time.Duration, bool)
	// Return tokens, useful when requests choose to drop themselves on timeout or cancellation
	Return(ctx context.Context, tokens float64)
	// Provides TokenManager the request that the scheduler processing -- some TokenManager implementations use this level of visibility for their algorithms. Return value decides whether the request has to be accepted right away in case TokenManger is not yet ready or configured to accept all traffic (short circuit).
	PreprocessRequest(ctx context.Context, request Request) (accept bool)
	// SetPassThrough sets the pass through flag for the token manager.
	SetPassThrough(passThrough bool)
	// GetPassThrough returns the pass through flag for the token manager.
	GetPassThrough() bool
}
