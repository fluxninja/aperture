package scheduler

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// TokenBucketMetrics holds metrics related to internals of TokenBucket.
type TokenBucketMetrics struct {
	FillRateGauge        prometheus.Gauge
	BucketCapacityGauge  prometheus.Gauge
	AvailableTokensGauge prometheus.Gauge
}

// Base for all token bucket. Implementations must take locks for concurrent access.
type tokenBucketBase struct {
	// Token Bucket
	latestTime      time.Time // Latest time for which we know the tokens in the bucket
	metrics         *TokenBucketMetrics
	fillRate        float64 // Tokens added per second
	bucketCapacity  float64 // Overall capacity of the bucket, currently the same as fillRate
	availableTokens float64 // Available Tokens as of latestTime (can be a negative or a fractional number)
}

func (tbb *tokenBucketBase) setFillRateGauge(v float64) {
	if tbb.metrics != nil && tbb.metrics.FillRateGauge != nil {
		tbb.metrics.FillRateGauge.Set(v)
	}
}

func (tbb *tokenBucketBase) setBucketCapacityGauge(v float64) {
	if tbb.metrics != nil && tbb.metrics.BucketCapacityGauge != nil {
		tbb.metrics.BucketCapacityGauge.Set(v)
	}
}

func (tbb *tokenBucketBase) setAvailableTokensGauge(v float64) {
	if tbb.metrics != nil && tbb.metrics.AvailableTokensGauge != nil {
		tbb.metrics.AvailableTokensGauge.Set(v)
	}
}

func (tbb *tokenBucketBase) adjustTokens(now time.Time) {
	if tbb.latestTime.IsZero() {
		tbb.latestTime = now
	}

	toAdd := float64(now.Sub(tbb.latestTime)) * float64(tbb.fillRate) / 1e9
	tbb.latestTime = now
	tbb.addTokens(toAdd)
}

// setFillRate : Adjust tokens fill rate (tokens/sec) at runtime.
func (tbb *tokenBucketBase) setFillRate(now time.Time, fillRate float64) {
	// adjust availableTokens
	tbb.adjustTokens(now)
	// subsequent requests get scheduled based on the new rate
	tbb.fillRate = fillRate
	tbb.setFillRateGauge(tbb.fillRate)
	// reset bucket
	tbb.bucketCapacity = fillRate
	tbb.setBucketCapacityGauge(tbb.bucketCapacity)
}

// add tokens.
func (tbb *tokenBucketBase) addTokens(toAdd float64) {
	tbb.availableTokens += toAdd
	if tbb.availableTokens > float64(tbb.bucketCapacity) {
		tbb.availableTokens = float64(tbb.bucketCapacity)
	}
	tbb.setAvailableTokensGauge(tbb.availableTokens)
}

func (tbb *tokenBucketBase) take(now time.Time, tokens float64) (time.Duration, bool) {
	// if tokens aren't coming don't provide any more tokens
	if tbb.fillRate == 0 {
		return time.Duration(0), false
	}

	tbb.adjustTokens(now)
	tbb.addTokens(-tokens)

	if tbb.availableTokens >= 0 {
		return time.Duration(0), true
	}
	// figure out when the tokens will be available
	waitTime := time.Duration(-tbb.availableTokens * float64(1e9) / float64(tbb.fillRate))

	return waitTime, true
}

func (tbb *tokenBucketBase) takeIfAvailable(now time.Time, tokens float64) bool {
	tbb.adjustTokens(now)
	if tbb.availableTokens > tokens {
		tbb.addTokens(-tokens)
		return true
	}
	return false
}

func (tbb *tokenBucketBase) getFillRate() float64 {
	return tbb.fillRate
}

// Check BasicTokenBucket implements TokenManager interface.
var _ TokenManager = &BasicTokenBucket{}

// BasicTokenBucket is a basic token bucket implementation.
type BasicTokenBucket struct {
	lock        sync.Mutex
	tbb         *tokenBucketBase
	passThrough bool
}

// NewBasicTokenBucket creates a new BasicTokenBucket with adjusted fill rate.
func NewBasicTokenBucket(now time.Time, fillRate float64, metrics *TokenBucketMetrics) *BasicTokenBucket {
	btb := &BasicTokenBucket{}
	btb.tbb = &tokenBucketBase{}

	if metrics != nil {
		btb.tbb.metrics = metrics
	}

	btb.tbb.setFillRate(now, fillRate)
	return btb
}

// TakeIfAvailable takes tokens from the basic token bucket if available, otherwise return false.
func (btb *BasicTokenBucket) TakeIfAvailable(now time.Time, tokens float64) bool {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	return btb.tbb.takeIfAvailable(now, tokens)
}

// Take takes tokens from the basic token bucket even if available tokens are less than asked.
// If tokens are not available at the moment, it will return amount of wait time and checks
// whether the operation was successful or not.
func (btb *BasicTokenBucket) Take(now time.Time, tokens float64) (time.Duration, bool) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	return btb.tbb.take(now, tokens)
}

// Return returns tokens to the basic token bucket.
func (btb *BasicTokenBucket) Return(tokens float64) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	btb.tbb.addTokens(tokens)
}

// PreprocessRequest is a no-op for BasicTokenBucket and by default, it rejects the request.
func (btb *BasicTokenBucket) PreprocessRequest(now time.Time, rContext Request) bool {
	return btb.passThrough
}

// SetFillRate adjusts the fill rate of the BasicTokenBucket.
func (btb *BasicTokenBucket) SetFillRate(now time.Time, fillRate float64) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	btb.tbb.setFillRate(now, fillRate)
}

// GetFillRate returns the fill rate of the BasicTokenBucket.
func (btb *BasicTokenBucket) GetFillRate() float64 {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	return btb.tbb.getFillRate()
}

// SetPassThrough sets the passThrough flag of the BasicTokenBucket.
func (btb *BasicTokenBucket) SetPassThrough(passThrough bool) {
	btb.passThrough = passThrough
}

// PassThrough returns the passThrough flag of the BasicTokenBucket.
func (btb *BasicTokenBucket) PassThrough() bool {
	return btb.passThrough
}
