package tokenbucket

import (
	"math"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"

	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	"github.com/fluxninja/aperture/v2/pkg/log"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

const (
	// TakeNFunction is the name of the function used to take N tokens from the bucket.
	TakeNFunction = "TakeN"
)

// TokenBucketRateLimiter implements Limiter.
type TokenBucketRateLimiter struct {
	mu             sync.RWMutex
	dMap           *olric.DMap
	name           string
	bucketCapacity float64
	fillAmount     float64
	interval       time.Duration
	continuousFill bool
	passThrough    bool
}

// NewTokenBucket creates a new instance of DistCacheRateTracker.
func NewTokenBucket(dc *distcache.DistCache,
	name string,
	interval time.Duration,
	maxIdleDuration time.Duration,
	continuousFill bool,
) (*TokenBucketRateLimiter, error) {
	dc.Mutex.Lock()
	defer dc.Mutex.Unlock()

	lbrl := &TokenBucketRateLimiter{
		name:           name,
		interval:       interval,
		passThrough:    true,
		continuousFill: continuousFill,
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

// SetBucketCapacity sets the rate limit for the rate limiter.
func (lbrl *TokenBucketRateLimiter) SetBucketCapacity(bucketCapacity float64) {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	lbrl.bucketCapacity = bucketCapacity
}

// GetBucketCapacity returns the rate limit for the rate limiter.
func (lbrl *TokenBucketRateLimiter) GetBucketCapacity() float64 {
	lbrl.mu.RLock()
	defer lbrl.mu.RUnlock()
	return lbrl.bucketCapacity
}

// SetFillAmount sets the default fill amount for the rate limiter.
func (lbrl *TokenBucketRateLimiter) SetFillAmount(fillAmount float64) {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	lbrl.fillAmount = fillAmount
}

// Name returns the name of the DistCacheRateTracker.
func (lbrl *TokenBucketRateLimiter) Name() string {
	return lbrl.name
}

// Close cleans up DMap held within the DistCacheRateTracker.
func (lbrl *TokenBucketRateLimiter) Close() error {
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
func (lbrl *TokenBucketRateLimiter) TakeIfAvailable(label string, n float64) (bool, float64, float64) {
	if lbrl.GetPassThrough() {
		return true, 0, 0
	}

	req := takeNRequest{
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
		log.Autosample().Errorf("error taking from token bucket: %v", err)
		return true, 0, 0
	}

	var resp takeNResponse

	err = utils.UnmarshalGob(resultBytes, &resp)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return true, 0, 0
	}

	return resp.Ok, resp.Remaining, resp.Current
}

// Take increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// It also returns the wait time at which the tokens will be available.
func (lbrl *TokenBucketRateLimiter) Take(label string, n float64) (bool, time.Duration, float64, float64) {
	if lbrl.GetPassThrough() {
		return true, 0, 0, 0
	}

	req := takeNRequest{
		Want:    n,
		CanWait: true,
	}
	// encode request
	reqBytes, err := utils.MarshalGob(req)
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return true, 0, 0, 0
	}

	resultBytes, err := lbrl.dMap.Function(label, TakeNFunction, reqBytes)
	if err != nil {
		log.Autosample().Errorf("error taking from token bucket: %v", err)
		return true, 0, 0, 0
	}

	var resp takeNResponse

	err = utils.UnmarshalGob(resultBytes, &resp)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return true, 0, 0, 0
	}
	var waitTime time.Duration
	if !resp.AvailableAt.IsZero() {
		waitTime = time.Until(resp.AvailableAt)
		if waitTime < 0 {
			waitTime = 0
		}
	}

	return resp.Ok, waitTime, resp.Remaining, resp.Current
}

// Return returns n tokens to the bucket.
func (lbrl *TokenBucketRateLimiter) Return(label string, n float64) (float64, float64) {
	if lbrl.GetPassThrough() {
		return 0, 0
	}
	_, remaining, current := lbrl.TakeIfAvailable(label, -n)
	return remaining, current
}

// Per-label tracking in distributed cache.
type tokenBucketState struct {
	LastFill  time.Time
	Available float64
}

type takeNRequest struct {
	Want    float64
	CanWait bool
}

type takeNResponse struct {
	AvailableAt time.Time
	Current     float64
	Remaining   float64
	Ok          bool
}

// takeN takes a number of tokens from the bucket.
func (lbrl *TokenBucketRateLimiter) takeN(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	lbrl.mu.RLock()
	defer lbrl.mu.RUnlock()

	// Decode currentState from gob encoded currentStateBytes
	now := time.Now()

	state, err := lbrl.fastForwardState(now, stateBytes)
	if err != nil {
		return nil, nil, err
	}

	// Decode arg from gob encoded argBytes
	var arg takeNRequest
	if argBytes != nil {
		err = utils.UnmarshalGob(argBytes, &arg)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	result := takeNResponse{
		Ok:          true,
		AvailableAt: now,
	}

	state.Available -= arg.Want
	if math.Signbit(state.Available) != math.Signbit(lbrl.bucketCapacity) {
		if lbrl.fillAmount != 0 {
			waitTime := time.Duration(math.Abs(state.Available) / math.Abs(lbrl.fillAmount) * float64(lbrl.interval))
			availableAt := now.Add(waitTime)
			result.AvailableAt = availableAt
		}
		result.Ok = arg.CanWait
		// return the tokens to the bucket if the request is not ok
		if !result.Ok {
			state.Available += arg.Want
		}
	} else if math.Abs(state.Available) > math.Abs(lbrl.bucketCapacity) {
		state.Available = lbrl.bucketCapacity
	}

	result.Remaining = state.Available
	result.Current = lbrl.bucketCapacity - state.Available

	// Encode result to gob encoded resultBytes
	resultBytes, err := utils.MarshalGob(result)
	if err != nil {
		log.Autosample().Errorf("error encoding result: %v", err)
		return nil, nil, err
	}

	// Encode currentState to gob encoded newStateBytes
	newStateBytes, err := utils.MarshalGob(state)
	if err != nil {
		log.Autosample().Errorf("error encoding new state: %v", err)
		return nil, nil, err
	}

	return newStateBytes, resultBytes, nil
}

func (lbrl *TokenBucketRateLimiter) fastForwardState(now time.Time, stateBytes []byte) (*tokenBucketState, error) {
	var state tokenBucketState

	if stateBytes != nil {
		err := utils.UnmarshalGob(stateBytes, &state)
		if err != nil {
			log.Autosample().Errorf("error decoding current state: %v", err)
			return nil, err
		}
	} else {
		state.LastFill = now
		state.Available = lbrl.bucketCapacity
	}

	if lbrl.fillAmount != 0 {
		// Calculate the time passed since the last fill
		sinceLastFill := now.Sub(state.LastFill)
		fillAmount := 0.0
		if lbrl.continuousFill {
			fillAmount = lbrl.fillAmount * float64(sinceLastFill) / float64(lbrl.interval)
			state.LastFill = now
		} else if sinceLastFill >= lbrl.interval {
			fills := int(sinceLastFill / lbrl.interval)
			if fills > 0 {
				fillAmount = lbrl.fillAmount * float64(fills)
				state.LastFill = state.LastFill.Add(time.Duration(fills) * lbrl.interval)
			}
		}
		// Fill the calculated amount
		state.Available += fillAmount

		if math.Signbit(state.Available) == math.Signbit(lbrl.bucketCapacity) &&
			math.Abs(state.Available) > math.Abs(lbrl.bucketCapacity) {
			state.Available = lbrl.bucketCapacity
		}
	}
	return &state, nil
}

// SetPassThrough sets the pass through flag.
func (lbrl *TokenBucketRateLimiter) SetPassThrough(passThrough bool) {
	lbrl.mu.Lock()
	defer lbrl.mu.Unlock()
	lbrl.passThrough = passThrough
}

// GetPassThrough returns the pass through flag.
func (lbrl *TokenBucketRateLimiter) GetPassThrough() bool {
	lbrl.mu.RLock()
	defer lbrl.mu.RUnlock()
	return lbrl.passThrough
}

// Make sure TokenBucketRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*TokenBucketRateLimiter)(nil)
