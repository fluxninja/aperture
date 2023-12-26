package scheduler

import (
	"context"
	"time"
)

// Scheduler : Interface for schedulers.
type Scheduler interface {
	// Schedule sends RequestContext to the underlying scheduler and returns a boolean value,
	// where true means accept and false means reject.
	Schedule(ctx context.Context, request *Request) (accept bool, remaining float64, current float64, reqID string)
	// Info returns the last access time and number of requests that are currently in the queue.
	Info() (time.Time, int)
}

// TokenManager : Interface for token managers.
type TokenManager interface {
	// Take tokens if available, otherwise return false
	TakeIfAvailable(ctx context.Context, tokens float64) (allowed bool, waitTime time.Duration, remaining float64, current float64, requestID string)
	// Take tokens even if available tokens are less than asked - returns wait time if tokens are not available immediately. The other return value conveys whether the operation was successful or not.
	Take(ctx context.Context, tokens float64) (allowed bool, waitTime time.Duration, remaining float64, current float64, requestID string)
	// Return tokens, useful when requests choose to drop themselves on timeout or cancellation
	Return(ctx context.Context, tokens float64, requestID string)
	// Provides TokenManager the request that the scheduler processing -- some TokenManager implementations use this level of visibility for their algorithms. Return value decides whether the request has to be accepted right away in case TokenManger is not yet ready or configured to accept all traffic (short circuit).
	PreprocessRequest(ctx context.Context, request *Request) (accept bool)
}
