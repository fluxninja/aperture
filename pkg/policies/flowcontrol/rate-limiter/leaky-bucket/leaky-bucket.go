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
)

const (
	// TakeNFunction is the name of the function used to take N tokens from the bucket.
	TakeNFunction = "TakeN"
)

// LeakyBucketRateLimiter implements Limiter.
type LeakyBucketRateLimiter struct {
	mu         sync.RWMutex
	dMap       *olric.DMap
	name       string
	parameters Parameters
}

// NewLeakyBucket creates a new instance of DistCacheRateTracker.
func NewLeakyBucket(dc *distcache.DistCache,
	name string,
	maxIdleDuration time.Duration,
) (*LeakyBucketRateLimiter, error) {
	dc.Mutex.Lock()
	defer dc.Mutex.Unlock()

	lbrl := &LeakyBucketRateLimiter{
		name: name,
		parameters: Parameters{
			BucketCapacity: -1,
		},
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

// SetParameters sets the default parameters for the rate limiter.
func (lbrl *LeakyBucketRateLimiter) SetParameters(params Parameters) {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	lbrl.parameters = params
}

// GetParameters returns the default parameters for the rate limiter.
func (lbrl *LeakyBucketRateLimiter) GetParameters() Parameters {
	lbrl.mu.RLock()
	defer lbrl.mu.RUnlock()
	return lbrl.parameters
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
	var reqBuf bytes.Buffer
	enc := gob.NewEncoder(&reqBuf)
	err := enc.Encode(req)
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return true, 0, 0
	}

	resultBytes, err := lbrl.dMap.Function(label, TakeNFunction, reqBuf.Bytes())
	if err != nil {
		log.Autosample().Errorf("error taking from leaky bucket: %v", err)
		return true, 0, 0
	}

	var resp response
	resBuf := bytes.NewBuffer(resultBytes)
	err = gob.NewDecoder(resBuf).Decode(&resp)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return true, 0, 0
	}

	return resp.Ok, resp.Remaining, resp.Current
}

// Parameters are the parameters for the leaky bucket rate limiter.
type Parameters struct {
	BucketCapacity float64
	LeakAmount     float64
	LeakInterval   time.Duration
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
		Current:  lbrl.parameters.BucketCapacity,
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

	if lbrl.parameters.BucketCapacity >= 0 {
		// First calculate current size of bucket based on time passed since last leak
		// and leak rate
		timeSinceLastLeak := now.Sub(currentState.LastLeak)

		if timeSinceLastLeak > lbrl.parameters.LeakInterval {
			// Calculate number of leaks since last leak
			leaks := int(timeSinceLastLeak / lbrl.parameters.LeakInterval)
			// Calculate amount to leak
			leakAmount := lbrl.parameters.LeakAmount * float64(leaks)
			// Leak
			currentState.Current -= leakAmount
			if currentState.Current < 0 {
				currentState.Current = 0
			}
			// Update lastLeak
			currentState.LastLeak = currentState.LastLeak.Add(time.Duration(leaks) * lbrl.parameters.LeakInterval)
		}

		currentState.Current += arg.Want

		if currentState.Current > lbrl.parameters.BucketCapacity {
			if arg.CanWait {
				waitTime := time.Duration((currentState.Current - lbrl.parameters.BucketCapacity) / lbrl.parameters.LeakAmount * float64(lbrl.parameters.LeakInterval))
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
	}

	result.Remaining = lbrl.parameters.BucketCapacity - currentState.Current
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
