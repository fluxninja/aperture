package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
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
	passThrough     bool
	clk             clockwork.Clock
}

func newTokenBucket(clk clockwork.Clock, metrics *TokenBucketMetrics) *tokenBucketBase {
	return &tokenBucketBase{
		clk:        clk,
		metrics:    metrics,
		latestTime: clk.Now(),
	}
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

func (tbb *tokenBucketBase) adjustTokens() {
	now := tbb.clk.Now()
	toAdd := float64(now.Sub(tbb.latestTime)) * float64(tbb.fillRate) / 1e9
	tbb.latestTime = now
	tbb.addTokens(toAdd)
}

// setFillRate : Adjust tokens fill rate (tokens/sec) at runtime.
func (tbb *tokenBucketBase) setFillRate(fillRate float64) {
	// adjust availableTokens
	tbb.adjustTokens()
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

// return tokens.
func (tbb *tokenBucketBase) returnTokens(toReturn float64) {
	if !tbb.passThrough {
		tbb.addTokens(toReturn)
	}
}

func (tbb *tokenBucketBase) take(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64) {
	// if tokens aren't coming do not provide any more tokens
	if tbb.fillRate == 0 {
		return false, time.Duration(0), tbb.availableTokens, tbb.bucketCapacity - tbb.availableTokens
	}

	tbb.adjustTokens()
	tbb.addTokens(-tokens)

	if tbb.availableTokens >= 0 {
		return true, time.Duration(0), tbb.availableTokens, tbb.bucketCapacity - tbb.availableTokens
	}
	// figure out when the tokens will be available
	waitTime := time.Duration(-tbb.availableTokens * float64(1e9) / float64(tbb.fillRate))

	// get deadline from context
	deadline, ok := ctx.Deadline()
	if ok && deadline.Sub(tbb.latestTime) < waitTime {
		// return tokens
		tbb.addTokens(tokens)
		return false, time.Duration(0), tbb.availableTokens, tbb.bucketCapacity - tbb.availableTokens
	}

	return true, waitTime, tbb.availableTokens, tbb.bucketCapacity - tbb.availableTokens
}

func (tbb *tokenBucketBase) takeIfAvailable(tokens float64) (bool, time.Duration, float64, float64) {
	tbb.adjustTokens()
	if tbb.availableTokens > tokens {
		tbb.addTokens(-tokens)
		return true, time.Duration(0), tbb.availableTokens, tbb.bucketCapacity - tbb.availableTokens
	}
	var waitTime time.Duration
	if tbb.fillRate > 0 {
		waitTime = time.Duration((tokens - tbb.availableTokens) * float64(1e9) / float64(tbb.fillRate))
	}
	return false, waitTime, tbb.availableTokens, tbb.bucketCapacity - tbb.availableTokens
}

func (tbb *tokenBucketBase) getFillRate() float64 {
	return tbb.fillRate
}

func (tbb *tokenBucketBase) setPassThrough(passThrough bool) {
	tbb.passThrough = passThrough
}

func (tbb *tokenBucketBase) getPassThrough() bool {
	return tbb.passThrough
}

// Check BasicTokenBucket implements TokenManager interface.
var _ TokenManager = &BasicTokenBucket{}

// BasicTokenBucket is a basic token bucket implementation.
type BasicTokenBucket struct {
	tbb  *tokenBucketBase
	lock sync.Mutex
}

// NewBasicTokenBucket creates a new BasicTokenBucket with adjusted fill rate.
func NewBasicTokenBucket(clk clockwork.Clock, fillRate float64, metrics *TokenBucketMetrics) *BasicTokenBucket {
	btb := &BasicTokenBucket{}
	btb.tbb = newTokenBucket(clk, metrics)
	btb.tbb.setFillRate(fillRate)
	return btb
}

// TakeIfAvailable takes tokens from the basic token bucket if available, otherwise return false.
func (btb *BasicTokenBucket) TakeIfAvailable(_ context.Context, tokens float64) (bool, time.Duration, float64, float64, string) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	ok, waitTime, remaining, current := btb.tbb.takeIfAvailable(tokens)
	return ok, waitTime, remaining, current, ""
}

// Take takes tokens from the basic token bucket even if available tokens are less than asked.
// If tokens are not available at the moment, it will return amount of wait time and checks
// whether the operation was successful or not.
func (btb *BasicTokenBucket) Take(ctx context.Context, tokens float64) (bool, time.Duration, float64, float64, string) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	ok, waitTime, remaining, current := btb.tbb.take(ctx, tokens)
	return ok, waitTime, remaining, current, ""
}

// Return returns tokens to the basic token bucket.
func (btb *BasicTokenBucket) Return(_ context.Context, tokens float64, _ string) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	btb.tbb.returnTokens(tokens)
}

// PreprocessRequest decides whether to proactively accept a request.
func (btb *BasicTokenBucket) PreprocessRequest(_ context.Context, request *Request) bool {
	return btb.tbb.getPassThrough()
}

// SetFillRate adjusts the fill rate of the BasicTokenBucket.
func (btb *BasicTokenBucket) SetFillRate(fillRate float64) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	btb.tbb.setFillRate(fillRate)
}

// GetFillRate returns the fill rate of the BasicTokenBucket.
func (btb *BasicTokenBucket) GetFillRate() float64 {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	return btb.tbb.getFillRate()
}

// SetPassThrough sets the passThrough flag of the BasicTokenBucket.
func (btb *BasicTokenBucket) SetPassThrough(passThrough bool) {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	btb.tbb.setPassThrough(passThrough)
}

// GetPassThrough returns the passThrough flag of the BasicTokenBucket.
func (btb *BasicTokenBucket) GetPassThrough() bool {
	btb.lock.Lock()
	defer btb.lock.Unlock()
	return btb.tbb.getPassThrough()
}
