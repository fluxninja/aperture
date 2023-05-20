package leakybucket

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"

	"github.com/fluxninja/aperture/v2/pkg/distcache"
	"github.com/fluxninja/aperture/v2/pkg/log"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

const (
	// TakeNFunction is the name of the function used to take N tokens from the bucket.
	TakeNFunction = "TakeN"
)

// LeakyBucketRateLimiter implements Limiter.
type LeakyBucketRateLimiter struct {
	mu             sync.RWMutex
	dMap           *olric.DMap
	name           string
	bucketCapacity float64
	leakAmount     float64
	interval       time.Duration
}

// NewLeakyBucket creates a new instance of DistCacheRateTracker.
func NewLeakyBucket(dc *distcache.DistCache,
	name string,
	interval time.Duration,
	maxIdleDuration time.Duration,
) (*LeakyBucketRateLimiter, error) {
	dc.Mutex.Lock()
	defer dc.Mutex.Unlock()

	lbrl := &LeakyBucketRateLimiter{
		name:           name,
		interval:       interval,
		bucketCapacity: -1,
	}

	dmapConfig := config.DMap{
		MaxIdleDuration: maxIdleDuration,
		Functions: map[string]config.Function{
			TakeNFunction: lbrl.takeN,
		},
	}

	dc.AddDMapCustomConfig(name, dmapConfig)
	dMap, err := dc.Olric.NewDMap(name)
	if err != nil {
		return nil, err
	}
	dc.RemoveDMapCustomConfig(name)

	lbrl.dMap = dMap

	return lbrl, nil
}

// SetRateLimit sets the rate limit for the rate limiter.
func (lbrl *LeakyBucketRateLimiter) SetRateLimit(rateLimit float64) {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	lbrl.bucketCapacity = rateLimit
}

// GetRateLimit returns the rate limit for the rate limiter.
func (lbrl *LeakyBucketRateLimiter) GetRateLimit() float64 {
	lbrl.mu.RLock()
	defer lbrl.mu.RUnlock()
	return lbrl.bucketCapacity
}

// SetLeakAmount sets the default leak amount for the rate limiter.
func (lbrl *LeakyBucketRateLimiter) SetLeakAmount(leakAmount float64) {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	lbrl.leakAmount = leakAmount
}

// Name returns the name of the DistCacheRateTracker.
func (lbrl *LeakyBucketRateLimiter) Name() string {
	return lbrl.name
}

// Close cleans up DMap held within the DistCacheRateTracker.
func (lbrl *LeakyBucketRateLimiter) Close() error {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	err := lbrl.dMap.Destroy()
	if err != nil {
		return err
	}
	return nil
}

// TakeIfAvailable increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// If an error occurred it returns true, 0 and 0 (fail open).
func (lbrl *LeakyBucketRateLimiter) TakeIfAvailable(label string, n float64) (bool, float64, float64) {
	req := request{
		Want:    n,
		CanWait: false,
	}
	// encode request
	reqBytes, err := utils.MarshalGob(req)
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return true, 0, 0
	}

	resultBytes, err := lbrl.dMap.Function(label, TakeNFunction, reqBytes)
	if err != nil {
		log.Autosample().Errorf("error taking from leaky bucket: %v", err)
		return true, 0, 0
	}

	var resp response

	err = utils.UnmarshalGob(resultBytes, &resp)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return true, 0, 0
	}

	return resp.Ok, resp.Remaining, resp.Current
}

// Per-label tracking in distributed cache.
type leakyBucketState struct {
	LastLeak time.Time
	Current  float64
}

type request struct {
	Want    float64
	CanWait bool
}

type response struct {
	AvailableAt time.Time
	Current     float64
	Remaining   float64
	Ok          bool
}

// takeN takes a number of tokens from the bucket.
func (lbrl *LeakyBucketRateLimiter) takeN(key string, currentStateBytes, argBytes []byte) ([]byte, []byte, error) {
	lbrl.mu.RLock()
	defer lbrl.mu.RUnlock()

	// Decode currentState from gob encoded currentStateBytes
	now := time.Now()
	currentState := leakyBucketState{
		LastLeak: now,
	}

	if currentStateBytes != nil {
		buf := bytes.NewBuffer(currentStateBytes)
		err := gob.NewDecoder(buf).Decode(&currentState)
		if err != nil {
			log.Autosample().Errorf("error decoding current state: %v", err)
			return nil, nil, err
		}
	}

	// Decode arg from gob encoded argBytes
	var arg request
	if argBytes != nil {
		buf := bytes.NewBuffer(argBytes)
		err := gob.NewDecoder(buf).Decode(&arg)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	result := response{
		Ok:          true,
		AvailableAt: now,
	}

	// Calculate the time passed since the last leak
	timeSinceLastLeak := now.Sub(currentState.LastLeak)

	// Calculate the amount to leak based on the time passed and leak rate
	leakAmount := lbrl.leakAmount * float64(timeSinceLastLeak) / float64(lbrl.interval)

	// Leak the calculated amount
	currentState.Current -= leakAmount
	if currentState.Current < 0 {
		currentState.Current = 0
	}

	// Update lastLeak
	currentState.LastLeak = now

	currentState.Current += arg.Want
	if currentState.Current > lbrl.bucketCapacity {
		if arg.CanWait {
			waitTime := time.Duration((currentState.Current - lbrl.bucketCapacity) / lbrl.leakAmount * float64(lbrl.interval))
			availableAt := now.Add(waitTime)
			result = response{
				Ok:          true,
				AvailableAt: availableAt,
			}
		} else {
			result.Ok = false
		}
	}

	// return the tokens to the bucket if the request is not ok
	if !result.Ok {
		currentState.Current -= arg.Want
	}

	result.Remaining = lbrl.bucketCapacity - currentState.Current
	result.Current = currentState.Current

	// Encode result to gob encoded resultBytes
	resultBuf := bytes.Buffer{}
	err := gob.NewEncoder(&resultBuf).Encode(result)
	if err != nil {
		log.Autosample().Errorf("error encoding result: %v", err)
		return nil, nil, err
	}

	// Encode currentState to gob encoded newStateBytes
	newStateBuf := bytes.Buffer{}
	err = gob.NewEncoder(&newStateBuf).Encode(currentState)
	if err != nil {
		log.Autosample().Errorf("error encoding new state: %v", err)
		return nil, nil, err
	}

	return newStateBuf.Bytes(), resultBuf.Bytes(), nil
}

// Make sure DistCacheRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*LeakyBucketRateLimiter)(nil)
