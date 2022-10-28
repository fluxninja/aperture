package scheduler

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/fluxninja/aperture/pkg/log"
)

// TokenBucketLoadMultiplierMetrics holds metrics related to internals of TokenBucketLoadMultiplier.
type TokenBucketLoadMultiplierMetrics struct {
	LMGauge            prometheus.Gauge
	TokenBucketMetrics *TokenBucketMetrics
}

// TokenBucketLoadMultiplier is a token bucket with load multiplier.
type TokenBucketLoadMultiplier struct {
	lock               sync.Mutex
	lmGauge            prometheus.Gauge // metrics
	tbb                *tokenBucketBase
	counter            *WindowedCounter
	lm                 float64 // load multiplier >=0
	continuousTracking bool
	passThrough        bool
}

// NewTokenBucketLoadMultiplier creates a new TokenBucketLoadMultiplier.
func NewTokenBucketLoadMultiplier(now time.Time,
	slotCount uint8,
	slotDuration time.Duration,
	metrics *TokenBucketLoadMultiplierMetrics,
) *TokenBucketLoadMultiplier {
	tbls := &TokenBucketLoadMultiplier{
		tbb:                &tokenBucketBase{},
		counter:            NewWindowedCounter(now, slotCount, slotDuration),
		lm:                 0,
		continuousTracking: false,
	}

	if metrics != nil {
		tbls.lmGauge = metrics.LMGauge
		tbls.tbb.metrics = metrics.TokenBucketMetrics
	}

	tbls.setLMGauge(float64(tbls.lm))

	return tbls
}

// SetContinuousTracking sets whether to continuously track the token rate and adjust the fill rate based on load multiplier.
func (tbls *TokenBucketLoadMultiplier) SetContinuousTracking(continuousTracking bool) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	tbls.continuousTracking = continuousTracking
}

// SetLoadMultiplier sets the load multiplier number --> 0 = no load accepted, 1 = accept up to 100% of current load, 2 = accept up to 200% of current load.
func (tbls *TokenBucketLoadMultiplier) SetLoadMultiplier(now time.Time, lm float64) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()

	if lm >= 0 {
		tbls.lm = lm
		tbls.setLMGauge(float64(tbls.lm))
		if !tbls.counter.IsBootstrapping() {
			// set fillRate based on latest sched.tokenRate
			tbls.tbb.setFillRate(now, tbls.counter.CalculateTokenRate()*tbls.lm)
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

// LoadMultiplier returns the current load multiplier.
func (tbls *TokenBucketLoadMultiplier) LoadMultiplier() float64 {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.lm
}

// PreprocessRequest preprocesses a request and makes decision whether to accept or reject the request.
func (tbls *TokenBucketLoadMultiplier) PreprocessRequest(now time.Time, rContext RequestContext) bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()

	wasBootstrapping := tbls.counter.IsBootstrapping()

	// Shift counter slot if needed
	ready := tbls.counter.AddTokens(now, rContext.Tokens)

	// recalculate token rate
	if ready {
		if wasBootstrapping || tbls.continuousTracking {
			// adjust fillRate based on the new tokenRate
			tbls.tbb.setFillRate(now, tbls.counter.CalculateTokenRate()*(tbls.lm))
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
		tbls.tbb.adjustTokens(now)
		return true
	}

	if tbls.passThrough {
		return true
	}

	return false
}

// TakeIfAvailable takes tokens from the token bucket if available, otherwise return false.
func (tbls *TokenBucketLoadMultiplier) TakeIfAvailable(now time.Time, tokens float64) bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.tbb.takeIfAvailable(now, tokens)
}

// Take takes tokens from the token bucket even if available tokens are less than asked.
// If tokens are not available at the moment, it will return amount of wait time and checks
// whether the operation was successful or not.
func (tbls *TokenBucketLoadMultiplier) Take(now time.Time, timeout time.Duration, tokens float64) (time.Duration, bool) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.tbb.take(now, timeout, tokens)
}

func (tbls *TokenBucketLoadMultiplier) setLMGauge(v float64) {
	if tbls.lmGauge != nil {
		tbls.lmGauge.Set(v)
	}
}

// PassThrough gets value of PassThrough flag.
func (tbls *TokenBucketLoadMultiplier) PassThrough() bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.passThrough
}

// SetPassThrough sets PassThrough flag which decides whether to pass through requests.
func (tbls *TokenBucketLoadMultiplier) SetPassThrough(passThrough bool) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	log.Trace().Msgf("Setting passThrough to %v", passThrough)
	tbls.passThrough = passThrough
}
