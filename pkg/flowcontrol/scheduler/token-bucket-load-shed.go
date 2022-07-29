package scheduler

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/fluxninja/aperture/pkg/log"
)

// TokenBucketLoadShedMetrics holds metrics related to internals of TokenBucketLoadShed.
type TokenBucketLoadShedMetrics struct {
	LSFGauge           prometheus.Gauge
	TokenBucketMetrics *TokenBucketMetrics
}

// TokenBucketLoadShed is a token bucket with load shedding.
type TokenBucketLoadShed struct {
	lock    sync.Mutex
	tbb     *tokenBucketBase
	counter *WindowedCounter
	lsf     float64 // load shed factor between 0 and 1
	// metrics
	lsfGauge prometheus.Gauge
}

// NewTokenBucketLoadShed creates a new TokenBucketLoadShed.
func NewTokenBucketLoadShed(now time.Time, metrics *TokenBucketLoadShedMetrics) *TokenBucketLoadShed {
	tbls := &TokenBucketLoadShed{
		tbb: &tokenBucketBase{},
	}

	if metrics != nil {
		tbls.lsfGauge = metrics.LSFGauge
		tbls.tbb.metrics = metrics.TokenBucketMetrics
	}

	slotDuration, _ := time.ParseDuration("100ms")
	// maintain a 1 sec window with 100ms slots for counters
	tbls.counter = NewWindowedCounter(now, 10, slotDuration)

	tbls.lsf = 0
	tbls.setLSFGauge(float64(tbls.lsf))

	return tbls
}

// SetLoadShedFactor sets the load shed factor number between [0,1] --> 0 = no load shedding, 1 = load shed 100%.
func (tbls *TokenBucketLoadShed) SetLoadShedFactor(now time.Time, lsf float64) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()

	if lsf >= 0 && lsf <= 1 {
		tbls.lsf = lsf
		tbls.setLSFGauge(float64(tbls.lsf))
		if !tbls.counter.IsBootstrapping() {
			// set fillRate based on latest sched.tokenRate
			tbls.tbb.setFillRate(now, tbls.counter.CalculateTokenRate()*(1-tbls.lsf))
			log.Debug().
				Float64("loadShedFactor", lsf).
				Float64("calculated fillRate", tbls.tbb.getFillRate()).
				Float64("calculated token rate", tbls.counter.CalculateTokenRate()).
				Msg("Controller update - Setting fill rate")
		}
	} else {
		log.Panic().Msgf("Load shed factor must be between 0 and 1, got %f", lsf)
	}
}

// LoadShedFactor returns the current load shed factor.
func (tbls *TokenBucketLoadShed) LoadShedFactor() float64 {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.lsf
}

// PreprocessRequest preprocesses a request and makes decision whether to accept or reject the request.
func (tbls *TokenBucketLoadShed) PreprocessRequest(now time.Time, rContext RequestContext) bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()

	wasBootstrapping := tbls.counter.IsBootstrapping()

	// Shift counter slot if needed
	ready := tbls.counter.AddTokens(now, rContext.Tokens)

	// recalculate token rate
	if ready {
		// adjust fillRate based on the new tokenRate
		tbls.tbb.setFillRate(now, tbls.counter.CalculateTokenRate()*(1-tbls.lsf))
		log.Trace().
			Float64("loadShedFactor", tbls.lsf).
			Float64("calculated fillRate", tbls.tbb.getFillRate()).
			Float64("calculated token rate", tbls.counter.CalculateTokenRate()).
			Msg("Sliding window update - Setting fill rate")

		if wasBootstrapping {
			// This is the first time we are learning the tokenRate, initialize the token bucket
			tbls.tbb.addTokens(tbls.tbb.getFillRate())
		}
	}

	// Accept this request if we are still learning the tokenRate
	// or if we don't intend to load shed at all
	if tbls.counter.IsBootstrapping() || tbls.lsf == 0 {
		tbls.tbb.adjustTokens(now)
		return true
	}

	return false
}

// TakeIfAvailable takes tokens from the token bucket if available, otherwise return false.
func (tbls *TokenBucketLoadShed) TakeIfAvailable(now time.Time, tokens float64) bool {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.tbb.takeIfAvailable(now, tokens)
}

// Take takes tokens from the token bucket even if available tokens are less than asked.
// If tokens are not available at the moment, it will return amount of wait time and checks
// whether the operation was successful or not.
func (tbls *TokenBucketLoadShed) Take(now time.Time, timeout time.Duration, tokens float64) (time.Duration, bool) {
	tbls.lock.Lock()
	defer tbls.lock.Unlock()
	return tbls.tbb.take(now, timeout, tokens)
}

func (tbls *TokenBucketLoadShed) setLSFGauge(v float64) {
	if tbls.lsfGauge != nil {
		tbls.lsfGauge.Set(v)
	}
}
