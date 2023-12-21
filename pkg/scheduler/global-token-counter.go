package scheduler

import (
	"context"
	"time"

	concurrencylimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/concurrency-limiter"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var _ TokenManager = &GlobalTokenCounter{}

// GlobalTokenCounter is a distributed rate-limiter token bucket implementation.
type GlobalTokenCounter struct {
	limiter concurrencylimiter.ConcurrencyLimiter
	key     string
}

// NewGlobalTokenCounter creates a new instance of GlobalTokenCounter.
func NewGlobalTokenCounter(key string, limiter concurrencylimiter.ConcurrencyLimiter) *GlobalTokenCounter {
	return &GlobalTokenCounter{
		key:     key,
		limiter: limiter,
	}
}

// SetPassThrough sets the passthrough value.
func (gtc *GlobalTokenCounter) SetPassThrough(passthrough bool) {
	gtc.limiter.SetPassThrough(passthrough)
}

// GetPassThrough returns the passthrough value.
func (gtc *GlobalTokenCounter) GetPassThrough() bool {
	return gtc.limiter.GetPassThrough()
}

// PreprocessRequest is a no-op.
func (gtc *GlobalTokenCounter) PreprocessRequest(_ context.Context, request *Request) bool {
	return gtc.GetPassThrough()
}

// TakeIfAvailable takes tokens if available.
func (gtc *GlobalTokenCounter) TakeIfAvailable(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64, string) {
	return gtc.limiter.TakeIfAvailable(ctx, gtc.key, tokens)
}

// Take takes tokens.
func (gtc *GlobalTokenCounter) Take(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64, string) {
	return gtc.limiter.Take(ctx, gtc.key, tokens)
}

// Return returns tokens.
func (gtc *GlobalTokenCounter) Return(ctx context.Context, tokens float64, reqID string) {
	_, err := gtc.limiter.Return(ctx, gtc.key, tokens, reqID)
	if err != nil {
		log.Autosample().Error().Err(err).Msg("Error returning tokens")
	}
}
