package ratelimiter

// RateLimiter is a generic limiter interface.
type RateLimiter interface {
	Name() string
	TakeIfAvailable(label string, count float64) (ok bool, remaining float64, current float64)
	GetRateLimit() float64
	SetRateLimit(limit float64)
	Close() error
}
