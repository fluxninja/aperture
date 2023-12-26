package concurrencylimiter

import (
	"context"
	"time"
)

// ConcurrencyLimiter is a generic limiter interface.
type ConcurrencyLimiter interface {
	Name() string
	TakeIfAvailable(ctx context.Context, label string, count float64) (ok bool, waitTime time.Duration, remaining float64, current float64, reqID string)
	Take(ctx context.Context, label string, count float64) (ok bool, waitTime time.Duration, remaining float64, current float64, reqID string)
	Return(ctx context.Context, label string, count float64, reqID string) (ok bool, err error)
	SetPassThrough(passthrough bool)
	GetPassThrough() bool
	Close() error
}
