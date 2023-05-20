package scheduler

import (
	"time"

	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
)

var _ TokenManager = &RateLimiterTokenBucket{}

// RateLimiterTokenBucket is a distributed rate-limiter token bucket implementation.
type RateLimiterTokenBucket struct {
	limiter ratelimiter.RateLimiter
	key     string
}

// NewRateLimiterTokenBucket creates a new instance of RateLimiterTokenBucket.
func NewRateLimiterTokenBucket(key string, limiter ratelimiter.RateLimiter) *RateLimiterTokenBucket {
	return &RateLimiterTokenBucket{
		key:     key,
		limiter: limiter,
	}
}

// SetPassThrough sets the passthrough value.
func (rltb *RateLimiterTokenBucket) SetPassThrough(passthrough bool) {
	rltb.limiter.SetPassThrough(passthrough)
}

// GetPassThrough returns the passthrough value.
func (rltb *RateLimiterTokenBucket) GetPassThrough() bool {
	return rltb.limiter.GetPassThrough()
}

// PreprocessRequest is a no-op.
func (rltb *RateLimiterTokenBucket) PreprocessRequest(now time.Time, request Request) bool {
	return rltb.GetPassThrough()
}

// TakeIfAvailable takes tokens if available.
func (rltb *RateLimiterTokenBucket) TakeIfAvailable(now time.Time, tokens float64) bool {
	ok, _, _ := rltb.limiter.TakeIfAvailable(rltb.key, tokens)
	return ok
}

// Take takes tokens.
func (rltb *RateLimiterTokenBucket) Take(now time.Time, tokens float64) (time.Duration, bool) {
	ok, waitTime, _, _ := rltb.limiter.Take(rltb.key, tokens)
	return waitTime, ok
}

// Return returns tokens.
func (rltb *RateLimiterTokenBucket) Return(tokens float64) {
	rltb.limiter.Return(rltb.key, tokens)
}
