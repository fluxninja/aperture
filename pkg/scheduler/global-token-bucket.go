package scheduler

import (
	"context"
	"time"

	ratelimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/rate-limiter"
)

var _ TokenManager = &GlobalTokenBucket{}

// GlobalTokenBucket is a distributed rate-limiter token bucket implementation.
type GlobalTokenBucket struct {
	limiter ratelimiter.RateLimiter
	key     string
}

// NewGlobalTokenBucket creates a new instance of GlobalTokenBucket.
func NewGlobalTokenBucket(key string, limiter ratelimiter.RateLimiter) *GlobalTokenBucket {
	return &GlobalTokenBucket{
		key:     key,
		limiter: limiter,
	}
}

// SetPassThrough sets the passthrough value.
func (gtb *GlobalTokenBucket) SetPassThrough(passthrough bool) {
	gtb.limiter.SetPassThrough(passthrough)
}

// GetPassThrough returns the passthrough value.
func (gtb *GlobalTokenBucket) GetPassThrough() bool {
	return gtb.limiter.GetPassThrough()
}

// PreprocessRequest is a no-op.
func (gtb *GlobalTokenBucket) PreprocessRequest(_ context.Context, request *Request) bool {
	return gtb.GetPassThrough()
}

// TakeIfAvailable takes tokens if available.
func (gtb *GlobalTokenBucket) TakeIfAvailable(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64, string) {
	ok, waitTime, remaining, current := gtb.limiter.TakeIfAvailable(ctx, gtb.key, tokens)
	return ok, waitTime, remaining, current, ""
}

// Take takes tokens.
func (gtb *GlobalTokenBucket) Take(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64, string) {
	ok, waitTime, remaining, current := gtb.limiter.Take(ctx, gtb.key, tokens)
	return ok, waitTime, remaining, current, ""
}

// Return returns tokens.
func (gtb *GlobalTokenBucket) Return(ctx context.Context, tokens float64, _ string) {
	gtb.limiter.Return(ctx, gtb.key, tokens)
}
