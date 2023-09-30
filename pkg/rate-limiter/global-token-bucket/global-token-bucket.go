package globaltokenbucket

import (
	"context"
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
	lookupMargin  = 10 * time.Millisecond
)

// GlobalTokenBucket implements Limiter.
type GlobalTokenBucket struct {
	dMap             olric.DMap
	dc               *distcache.DistCache
	name             string
	bucketCapacity   float64
	fillAmount       float64
	interval         time.Duration
	mu               sync.RWMutex
	continuousFill   bool
	delayInitialFill bool
	passThrough      bool
}

// NewGlobalTokenBucket creates a new instance of DistCacheRateTracker.
func NewGlobalTokenBucket(dc *distcache.DistCache,
	name string,
	interval time.Duration,
	maxIdleDuration time.Duration,
	continuousFill bool,
	delayInitialFill bool,
) (*GlobalTokenBucket, error) {
	gtb := &GlobalTokenBucket{
		name:             name,
		interval:         interval,
		passThrough:      true,
		continuousFill:   continuousFill,
		delayInitialFill: delayInitialFill,
		dc:               dc,
	}

	dmapConfig := config.DMap{
		MaxIdleDuration: maxIdleDuration,
		Functions: map[string]config.Function{
			TakeNFunction: gtb.takeN,
		},
	}

	dMap, err := dc.NewDMap(name, dmapConfig)
	if err != nil {
		return nil, err
	}

	gtb.dMap = dMap

	return gtb, nil
}

// SetBucketCapacity sets the rate limit for the rate limiter.
func (gtb *GlobalTokenBucket) SetBucketCapacity(bucketCapacity float64) {
	gtb.mu.Lock()
	defer gtb.mu.Unlock()
	gtb.bucketCapacity = bucketCapacity
}

// GetBucketCapacity returns the rate limit for the rate limiter.
func (gtb *GlobalTokenBucket) GetBucketCapacity() float64 {
	gtb.mu.RLock()
	defer gtb.mu.RUnlock()
	return gtb.bucketCapacity
}

func isMarginExceeded(ctx context.Context) bool {
	deadline, deadlineOK := ctx.Deadline()
	if deadlineOK {
		// check if deadline will be passed in the next 10ms
		deadline = deadline.Add(-lookupMargin)
		return time.Now().After(deadline)
	}
	return false
}

// SetFillAmount sets the default fill amount for the rate limiter.
func (gtb *GlobalTokenBucket) SetFillAmount(fillAmount float64) {
	gtb.mu.Lock()
	defer gtb.mu.Unlock()
	gtb.fillAmount = fillAmount
}

// Name returns the name of the DistCacheRateTracker.
func (gtb *GlobalTokenBucket) Name() string {
	return gtb.name
}

// Close cleans up DMap held within the DistCacheRateTracker.
func (gtb *GlobalTokenBucket) Close() error {
	gtb.mu.Lock()
	defer gtb.mu.Unlock()
	err := gtb.dc.DeleteDMap(gtb.name)
	if err != nil {
		return err
	}
	return nil
}

func (gtb *GlobalTokenBucket) executeTakeRequest(ctx context.Context, label string, n float64, canWait bool, deadline time.Time) (bool, time.Duration, float64, float64) {
	if gtb.GetPassThrough() {
		return true, 0, 0, 0
	}

	if isMarginExceeded(ctx) {
		return false, 0, 0, 0
	}

	req := takeNRequest{
		Want:     n,
		CanWait:  canWait,
		Deadline: deadline,
	}
	reqBytes, err := utils.MarshalGob(req)
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return true, 0, 0, 0
	}

	resultBytes, err := gtb.dMap.Function(ctx, label, TakeNFunction, reqBytes)
	if err != nil {
		log.Autosample().Error().Err(err).Str("dmapName", gtb.dMap.Name()).Float64("tokens", n).Msg("error taking from token bucket")
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

// TakeIfAvailable increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// If an error occurred it returns true, 0, 0 and 0 (fail open).
// It also may return the wait time at which the tokens will be available.
func (gtb *GlobalTokenBucket) TakeIfAvailable(ctx context.Context, label string, n float64) (bool, time.Duration, float64, float64) {
	return gtb.executeTakeRequest(ctx, label, n, false, time.Time{})
}

// Take increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// It also returns the wait time at which the tokens will be available.
func (gtb *GlobalTokenBucket) Take(ctx context.Context, label string, n float64) (bool, time.Duration, float64, float64) {
	deadline := time.Time{}
	d, ok := ctx.Deadline()
	if ok {
		deadline = d
	}

	return gtb.executeTakeRequest(ctx, label, n, true, deadline)
}

// Return returns n tokens to the bucket.
func (gtb *GlobalTokenBucket) Return(ctx context.Context, label string, n float64) (float64, float64) {
	_, _, remaining, current := gtb.TakeIfAvailable(ctx, label, -n)
	return remaining, current
}

// Per-label tracking in distributed cache.
type tokenBucketState struct {
	StartFillAt time.Time // To prevent more tokens than fillRate in a time window while using burst capacity
	LastFill    time.Time
	Available   float64
}

type takeNRequest struct {
	Deadline time.Time
	Want     float64
	CanWait  bool
}

type takeNResponse struct {
	AvailableAt time.Time
	Current     float64
	Remaining   float64
	Ok          bool
}

// takeN takes a number of tokens from the bucket.
func (gtb *GlobalTokenBucket) takeN(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	gtb.mu.RLock()
	defer gtb.mu.RUnlock()

	// Decode currentState from gob encoded currentStateBytes
	now := time.Now()

	state, err := gtb.fastForwardState(now, stateBytes)
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

	// if we are first time drawing from the bucket, set the start fill time
	if gtb.delayInitialFill && state.Available == gtb.bucketCapacity {
		state.StartFillAt = now.Add(gtb.timeToFill(gtb.bucketCapacity))
		if gtb.continuousFill {
			state.LastFill = state.StartFillAt
		} else {
			state.LastFill = state.StartFillAt.Add(-gtb.interval)
		}
	}

	state.Available -= arg.Want

	if arg.Want > 0 {
		if state.Available < 0 {
			result.Ok = arg.CanWait && gtb.fillAmount != 0
			if gtb.fillAmount != 0 {
				result.AvailableAt = gtb.getAvailableAt(now, state)
				if arg.CanWait && !arg.Deadline.IsZero() && result.AvailableAt.After(arg.Deadline) {
					result.Ok = false
				}
			}
			// return the tokens to the bucket if the request is not ok
			if !result.Ok {
				state.Available += arg.Want
			}
		}
	}

	if state.Available > gtb.bucketCapacity {
		state.Available = gtb.bucketCapacity
	}

	result.Remaining = state.Available
	result.Current = gtb.bucketCapacity - state.Available

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

func (gtb *GlobalTokenBucket) fastForwardState(now time.Time, stateBytes []byte) (*tokenBucketState, error) {
	var state tokenBucketState

	if stateBytes != nil {
		err := utils.UnmarshalGob(stateBytes, &state)
		if err != nil {
			log.Autosample().Errorf("error decoding current state: %v", err)
			return nil, err
		}
	} else {
		state.LastFill = now
		state.Available = gtb.bucketCapacity
	}

	// do not fill the bucket until the start fill time
	if state.StartFillAt.IsZero() || now.After(state.StartFillAt) {
		// Calculate the time passed since the last fill
		sinceLastFill := now.Sub(state.LastFill)
		fillAmount := 0.0
		if gtb.continuousFill {
			fillAmount = gtb.fillAmount * float64(sinceLastFill) / float64(gtb.interval)
			state.LastFill = now
		} else if sinceLastFill >= gtb.interval {
			fills := int(sinceLastFill / gtb.interval)
			if fills > 0 {
				fillAmount = gtb.fillAmount * float64(fills)
				state.LastFill = state.LastFill.Add(time.Duration(fills) * gtb.interval)
			}
		}
		// Fill the calculated amount
		state.Available += fillAmount

		if state.Available > gtb.bucketCapacity {
			state.Available = gtb.bucketCapacity
		}
	}
	return &state, nil
}

// timeToFill calculates the wait time for the given number of tokens based on the fill rate.
func (gtb *GlobalTokenBucket) timeToFill(tokens float64) time.Duration {
	if gtb.fillAmount != 0 {
		if gtb.continuousFill {
			return time.Duration(tokens / gtb.fillAmount * float64(gtb.interval))
		} else {
			// calculate how many fills we need
			fills := math.Ceil(tokens / gtb.fillAmount)
			return time.Duration(fills) * gtb.interval
		}
	}
	return 0
}

// getAvailableAt calculates the time at which the given number of tokens will be available.
func (gtb *GlobalTokenBucket) getAvailableAt(now time.Time, state *tokenBucketState) time.Time {
	if state.Available >= 0 {
		return now
	}
	timeToFill := gtb.timeToFill(-state.Available)
	if now.Before(state.StartFillAt) {
		return state.StartFillAt.Add(timeToFill)
	} else {
		// this code assumes that other parts of the code are correct, such as
		// LastFill is not in the future if now is after StartFillAt
		// And timeSinceLastFill is not greater than interval
		timeSinceLastFill := now.Sub(state.LastFill)
		if timeSinceLastFill > gtb.interval {
			log.Autosample().Errorf("time since last fill is greater than interval: %v", timeSinceLastFill)
			timeSinceLastFill = time.Duration(0)
		}
		return now.Add(timeToFill - timeSinceLastFill)
	}
}

// SetPassThrough sets the pass through flag.
func (gtb *GlobalTokenBucket) SetPassThrough(passThrough bool) {
	gtb.mu.Lock()
	defer gtb.mu.Unlock()
	gtb.passThrough = passThrough
}

// GetPassThrough returns the pass through flag.
func (gtb *GlobalTokenBucket) GetPassThrough() bool {
	gtb.mu.RLock()
	defer gtb.mu.RUnlock()
	return gtb.passThrough
}

// Make sure TokenBucketRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*GlobalTokenBucket)(nil)
