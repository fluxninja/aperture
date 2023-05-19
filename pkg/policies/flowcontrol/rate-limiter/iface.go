package ratelimiter

// RateLimiter is a generic limiter interface.
type RateLimiter interface {
	TakeIfAvailable(label string, count float64) (ok bool, remaining float64, current float64)
}
