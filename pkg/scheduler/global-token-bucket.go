package scheduler

import (
	"context"
	"time"

	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
)

var _ TokenManager = &GlobalTokenBucket{}

// GlobalTokenBucket is a distributed rate-limiter token bucket implementation.
type GlobalTokenBucket struct {
	limiter ratelimiter.RateLimiter
	key     string
}

// NewGlobalTokenBucket creates a new instance of RateLimiterTokenBucket.
func NewGlobalTokenBucket(key string, limiter ratelimiter.RateLimiter) *GlobalTokenBucket {
	return &GlobalTokenBucket{
		key:     key,
		limiter: limiter,
	}
}

// SetPassThrough sets the passthrough value.
func (rltb *GlobalTokenBucket) SetPassThrough(passthrough bool) {
	rltb.limiter.SetPassThrough(passthrough)
}

// GetPassThrough returns the passthrough value.
func (rltb *GlobalTokenBucket) GetPassThrough() bool {
	return rltb.limiter.GetPassThrough()
}

// PreprocessRequest is a no-op.
func (rltb *GlobalTokenBucket) PreprocessRequest(_ context.Context, request *Request) bool {
	return rltb.GetPassThrough()
}

// TakeIfAvailable takes tokens if available.
func (rltb *GlobalTokenBucket) TakeIfAvailable(ctx context.Context, tokens float64) bool {
	ok, _, _ := rltb.limiter.TakeIfAvailable(ctx, rltb.key, tokens)
	return ok
}

// Take takes tokens.
func (rltb *GlobalTokenBucket) Take(ctx context.Context, tokens float64) (time.Duration, bool) {
	ok, waitTime, _, _ := rltb.limiter.Take(ctx, rltb.key, tokens)
	return waitTime, ok
}

// Return returns tokens.
func (rltb *GlobalTokenBucket) Return(ctx context.Context, tokens float64) {
	rltb.limiter.Return(ctx, rltb.key, tokens)
}
