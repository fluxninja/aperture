package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter is a generic limiter interface.
type RateLimiter interface {
	Name() string
	TakeIfAvailable(ctx context.Context, label string, count float64) (ok bool, remaining float64, current float64)
	Take(ctx context.Context, label string, count float64) (ok bool, waitTime time.Duration, remaining float64, current float64)
	Return(ctx context.Context, label string, count float64) (remaining float64, current float64)
	SetPassThrough(passthrough bool)
	GetPassThrough() bool
	Close() error
}
