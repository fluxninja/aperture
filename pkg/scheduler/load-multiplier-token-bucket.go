package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"

	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Check if TokenBucketLoadMultiplier implements TokenManager interface.
var _ TokenManager = &LoadMultiplierTokenBucket{}

// LoadMultiplierTokenBucket is a token bucket with load multiplier.
type LoadMultiplierTokenBucket struct {
	lmGauge            prometheus.Gauge // metrics
	tbb                *tokenBucketBase
	counter            *WindowedCounter
	lm                 float64 // load multiplier >=0
	lock               sync.Mutex
	continuousTracking bool
	validUntil         time.Time
}

// NewLoadMultiplierTokenBucket creates a new TokenBucketLoadMultiplier.
func NewLoadMultiplierTokenBucket(
	clk clockwork.Clock,
	slotCount uint8,
	slotDuration time.Duration,
	lmGauge prometheus.Gauge,
	tbMetrics *TokenBucketMetrics,
) *LoadMultiplierTokenBucket {
	tbb := newTokenBucket(clk, tbMetrics)

	tbls := &LoadMultiplierTokenBucket{
		tbb:                tbb,
		counter:            NewWindowedCounter(clk, slotCount, slotDuration),
		lm:                 0,
		continuousTracking: false,
		validUntil:         time.Time{},
	}
	tbls.tbb.setPassThrough(true)

	tbls.lmGauge = lmGauge
	tbls.setLMGauge(float64(tbls.lm))

	return tbls
}

// SetContinuousTracking sets whether to continuously track the token rate and adjust the fill rate based on load multiplier.
func (tbls *LoadMultiplierTokenBucket) SetContinuousTracking(continuousTracking bool) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	tbls.continuousTracking = continuousTracking
}

// SetLoadDecisionValues sets the load multiplier number --> 0 = no load accepted, 1 = accept up to 100% of current load, 2 = accept up to 200% of current load.
func (tbls *LoadMultiplierTokenBucket) SetLoadDecisionValues(loadDecision *policysyncv1.LoadDecision) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()

	lm := loadDecision.GetLoadMultiplier()
	if lm >= 0 {
		// set valid till timestamp if it is set else set it to 0
		loadDecisionValidUntil := loadDecision.GetValidUntil()
		if loadDecisionValidUntil != nil {
			tbls.validUntil = loadDecisionValidUntil.AsTime()
		} else {
			tbls.validUntil = time.Time{}
		}
		tbls.lm = lm
		tbls.setLMGauge(float64(tbls.lm))
		if !tbls.counter.IsBootstrapping() {
			// set fillRate based on latest sched.tokenRate
			tbls.setFillRate()
			log.Trace().
				Float64("loadMultiplier", tbls.lm).
				Float64("calculated fillRate", tbls.tbb.getFillRate()).
				Float64("calculated token rate", tbls.counter.CalculateTokenRate()).
				Msg("Controller update - Setting fill rate")
		}
	} else {
		log.Panic().Msgf("Load multiplier must be greater than 0, got %f", lm)
	}
}

// setFillRate - unsage, must be called with lock held.
func (tbls *LoadMultiplierTokenBucket) setFillRate() {
	fillRate := tbls.counter.CalculateTokenRate() * tbls.lm
	tbls.tbb.setFillRate(fillRate)
}

// LoadMultiplier returns the current load multiplier.
func (tbls *LoadMultiplierTokenBucket) LoadMultiplier() float64 {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.lm
}

// PreprocessRequest preprocesses a request and makes decision whether to pro-actively accept a request.
func (tbls *LoadMultiplierTokenBucket) PreprocessRequest(_ context.Context, request *Request) bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()

	wasBootstrapping := tbls.counter.IsBootstrapping()

	// Shift counter slot if needed
	ready := tbls.counter.AddTokens(request)

	// recalculate token rate
	if ready {
		if wasBootstrapping || tbls.continuousTracking {
			// adjust fillRate based on the new tokenRate
			tbls.setFillRate()
			log.Trace().
				Float64("loadMultiplier", tbls.lm).
				Float64("calculated fillRate", tbls.tbb.getFillRate()).
				Float64("calculated token rate", tbls.counter.CalculateTokenRate()).
				Msg("Sliding window update - Setting fill rate")
			if wasBootstrapping {
				// This is the first time we are learning the tokenRate, initialize the token bucket
				tbls.tbb.addTokens(tbls.tbb.getFillRate())
			}
		}
	}

	// Accept this request if we are still learning the tokenRate
	if tbls.counter.IsBootstrapping() {
		tbls.tbb.adjustTokens()
		return true
	}

	return tbls.getPassThrough()
}

// TakeIfAvailable takes tokens from the token bucket if available, otherwise return false.
func (tbls *LoadMultiplierTokenBucket) TakeIfAvailable(_ context.Context, tokens float64) (bool, time.Duration, float64, float64) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.tbb.takeIfAvailable(tokens)
}

// Take takes tokens from the token bucket even if available tokens are less than asked.
// If tokens are not available at the moment, it will return amount of wait time and checks
// whether the operation was successful or not.
func (tbls *LoadMultiplierTokenBucket) Take(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.tbb.take(ctx, tokens)
}

// Return returns tokens to the token bucket.
func (tbls *LoadMultiplierTokenBucket) Return(_ context.Context, tokens float64) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	tbls.tbb.returnTokens(tokens)
}

func (tbls *LoadMultiplierTokenBucket) setLMGauge(v float64) {
	if tbls.lmGauge != nil {
		tbls.lmGauge.Set(v)
	}
}

// GetPassThrough gets value of passThrough flag.
func (tbls *LoadMultiplierTokenBucket) GetPassThrough() bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.getPassThrough()
}

// getPassThrough gets value of passThrough flag. Unsafe, must be called with lock held.
func (tbls *LoadMultiplierTokenBucket) getPassThrough() bool {
	if !tbls.validUntil.IsZero() && tbls.validUntil.Before(time.Now()) {
		return true
	}
	return tbls.tbb.getPassThrough()
}

// SetPassThrough sets PassThrough flag which decides whether to pass through requests.
func (tbls *LoadMultiplierTokenBucket) SetPassThrough(passThrough bool) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	log.Trace().Msgf("Setting passThrough to %v", passThrough)
	tbls.tbb.setPassThrough(passThrough)
}
