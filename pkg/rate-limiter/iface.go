package ratelimiter

import "time"

// RateLimiter is a generic limiter interface.
type RateLimiter interface {
	Name() string
	TakeIfAvailable(label string, count float64) (ok bool, remaining float64, current float64)
	Take(label string, count float64) (ok bool, waitTime time.Duration, remaining float64, current float64)
	Return(label string, count float64) (remaining float64, current float64)
	SetPassThrough(passthrough bool)
	GetPassThrough() bool
	Close() error
}
