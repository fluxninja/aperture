package globaltokenbucket

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	tokenbucketv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/tokenbucket/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	deadlinemargin "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/deadline-margin"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

const (
	// TakeNFunction is the name of the function used to take N tokens from the bucket.
	TakeNFunction = "TakeN"
)

// counter is used to track the state of a label.
type counter struct {
	availableAt         time.Time
	remaining           float64
	current             float64
	lastRequestedTokens float64
}

// GlobalTokenBucket implements Limiter.
type GlobalTokenBucket struct {
	dMap             olric.DMap
	dc               *distcache.DistCache
	jobGroup         *jobs.JobGroup
	counters         sync.Map
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
func NewGlobalTokenBucket(
	dc *distcache.DistCache,
	name string,
	interval time.Duration,
	maxIdleDuration time.Duration,
	continuousFill bool,
	delayInitialFill bool,
	jobGroup *jobs.JobGroup,
) (*GlobalTokenBucket, error) {
	gtb := &GlobalTokenBucket{
		name:             name,
		interval:         interval,
		passThrough:      true,
		continuousFill:   continuousFill,
		delayInitialFill: delayInitialFill,
		dc:               dc,
		jobGroup:         jobGroup,
	}

	job := jobs.NewBasicJob(gtb.name, gtb.audit)
	// register job with job group
	err := gtb.jobGroup.RegisterJob(job, jobs.JobConfig{
		ExecutionPeriod: config.MakeDuration(interval),
	})
	if err != nil {
		return nil, err
	}

	dmapConfig := olricconfig.DMap{
		MaxIdleDuration: maxIdleDuration,
		Functions: map[string]olricconfig.Function{
			TakeNFunction: gtb.takeN,
		},
	}

	dMap, err := dc.NewDMap(name, dmapConfig, false)
	if err != nil {
		return nil, err
	}

	gtb.dMap = dMap

	return gtb, nil
}

func (gtb *GlobalTokenBucket) audit(ctx context.Context) (proto.Message, error) {
	now := time.Now()
	// range through the map and sync the counters
	gtb.counters.Range(func(label, value interface{}) bool {
		c := value.(*counter)

		// if the availableAt time is in the past, remove the counter
		if now.After(c.availableAt) {
			gtb.counters.Delete(label)
			return true
		}
		return true
	})
	return nil, nil
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
	err = gtb.jobGroup.DeregisterJob(gtb.name)
	if err != nil {
		return err
	}
	return nil
}

func (gtb *GlobalTokenBucket) executeTakeRequest(ctx context.Context, label string, n float64, canWait bool, deadline time.Time) (bool, time.Duration, float64, float64) {
	getWaitTime := func(availableAt time.Time) time.Duration {
		var waitTime time.Duration
		if !availableAt.IsZero() {
			waitTime = time.Until(availableAt)
			if waitTime < 0 {
				waitTime = 0
			}
		}
		return waitTime
	}

	if gtb.GetPassThrough() {
		return true, 0, 0, 0
	}

	c := &counter{}

	if !canWait {
		existing, ok := gtb.counters.Load(label)
		if ok {
			c = existing.(*counter)
			waitTime := getWaitTime(c.availableAt)
			if c.lastRequestedTokens >= n && waitTime > 0 {
				return false, waitTime, c.remaining, c.current
			}
		}
	}

	if deadlinemargin.IsMarginExceeded(ctx) {
		return false, 0, 0, 0
	}

	req := tokenbucketv1.TakeNRequest{
		Deadline: timestamppb.New(deadline),
		Want:     n,
		CanWait:  canWait,
	}
	reqBytes, err := req.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return true, 0, 0, 0
	}

	resultBytes, err := gtb.dMap.Function(ctx, label, TakeNFunction, reqBytes)
	if err != nil {
		log.Autosample().Error().Err(err).Str("dmapName", gtb.dMap.Name()).Float64("tokens", n).Msg("error taking from token bucket")
		return true, 0, 0, 0
	}

	var resp tokenbucketv1.TakeNResponse
	err = resp.UnmarshalVT(resultBytes)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return true, 0, 0, 0
	}

	if !canWait {
		if n < 0 {
			gtb.counters.Delete(label)
		} else if !resp.GetOk() {
			c = &counter{
				availableAt:         resp.GetAvailableAt().AsTime(),
				remaining:           resp.GetRemaining(),
				current:             resp.GetCurrent(),
				lastRequestedTokens: n,
			}
			gtb.counters.Store(label, c)
		}
	}

	availableAt := resp.GetAvailableAt().AsTime()
	waitTime := getWaitTime(availableAt)

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

// takeN takes a number of tokens from the bucket.
func (gtb *GlobalTokenBucket) takeN(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	gtb.mu.RLock()
	defer gtb.mu.RUnlock()

	// Decode currentState from proto encoded currentStateBytes
	now := time.Now()

	state, err := gtb.fastForwardState(now, stateBytes, key)
	if err != nil {
		return nil, nil, err
	}

	// Decode arg from proto encoded argBytes
	var arg tokenbucketv1.TakeNRequest
	if argBytes != nil {
		err = arg.UnmarshalVT(argBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	result := tokenbucketv1.TakeNResponse{
		Ok:          true,
		AvailableAt: timestamppb.New(now),
	}

	// if we are first time drawing from the bucket, set the start fill time
	if gtb.delayInitialFill && state.Available == gtb.bucketCapacity {
		state.StartFillAt = timestamppb.New(now.Add(gtb.timeToFill(gtb.bucketCapacity)))
		if gtb.continuousFill {
			state.LastFillAt = state.StartFillAt
		} else {
			startFilleAtTime := state.StartFillAt.AsTime()
			state.LastFillAt = timestamppb.New(startFilleAtTime.Add(-gtb.interval))
		}
	}

	state.Available -= arg.Want

	if arg.Want > 0 {
		if state.Available < 0 {
			result.Ok = arg.CanWait && gtb.fillAmount != 0
			if gtb.fillAmount != 0 {
				result.AvailableAt = timestamppb.New(gtb.getAvailableAt(now, state))
				deadlineTime := arg.Deadline.AsTime()
				if arg.CanWait && !deadlineTime.IsZero() && result.AvailableAt.AsTime().After(arg.Deadline.AsTime()) {
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

	// Encode result to proto encoded resultBytes
	resultBytes, err := result.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding result: %v", err)
		return nil, nil, err
	}

	// Encode currentState to proto encoded newStateBytes
	newStateBytes, err := state.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding new state: %v", err)
		return nil, nil, err
	}

	return newStateBytes, resultBytes, nil
}

func (gtb *GlobalTokenBucket) fastForwardState(now time.Time, stateBytes []byte, key string) (*tokenbucketv1.State, error) {
	var state tokenbucketv1.State

	if stateBytes != nil {
		err := state.UnmarshalVT(stateBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding current state: %v", err)
			return nil, err
		}
	} else {
		log.Trace().Str("key", key).Str("dMap", gtb.dMap.Name()).Msg("Creating new token bucket state")
		state.LastFillAt = timestamppb.New(now)
		state.Available = gtb.bucketCapacity
	}

	startFillAtTime := state.StartFillAt.AsTime()
	lastFillAtTime := state.LastFillAt.AsTime()

	// do not fill the bucket until the start fill time
	if startFillAtTime.IsZero() || now.After(startFillAtTime) {
		// Calculate the time passed since the last fill
		sinceLastFill := now.Sub(lastFillAtTime)
		fillAmount := 0.0
		if gtb.continuousFill {
			fillAmount = gtb.fillAmount * float64(sinceLastFill) / float64(gtb.interval)
			state.LastFillAt = timestamppb.New(now)
		} else if sinceLastFill >= gtb.interval {
			fills := int(sinceLastFill / gtb.interval)
			if fills > 0 {
				fillAmount = gtb.fillAmount * float64(fills)
				state.LastFillAt = timestamppb.New(lastFillAtTime.Add(time.Duration(fills) * gtb.interval))
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
func (gtb *GlobalTokenBucket) getAvailableAt(now time.Time, state *tokenbucketv1.State) time.Time {
	if state.Available >= 0 {
		return now
	}
	timeToFill := gtb.timeToFill(-state.Available)
	startFillAtTime := state.StartFillAt.AsTime()
	if now.Before(startFillAtTime) {
		return startFillAtTime.Add(timeToFill)
	} else {
		// this code assumes that other parts of the code are correct, such as
		// LastFill is not in the future if now is after StartFillAt
		// And timeSinceLastFill is not greater than interval
		lastFillAtTime := state.LastFillAt.AsTime()
		timeSinceLastFill := now.Sub(lastFillAtTime)
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
