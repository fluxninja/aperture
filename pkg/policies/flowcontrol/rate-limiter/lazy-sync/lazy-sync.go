package lazysync

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter"
	"google.golang.org/protobuf/proto"
)

type counter struct {
	lock      sync.Mutex
	nextSync  time.Time
	available float64
	credit    float64
}

// LazySyncRateLimiter is a limiter that syncs its state lazily with another limiter.
type LazySyncRateLimiter struct {
	counters     sync.Map
	limiter      ratelimiter.RateLimiter
	jobGroup     *jobs.JobGroup
	name         string
	syncDuration time.Duration
}

// NewLazySyncRateLimiter creates a new LazySyncLimiter.
func NewLazySyncRateLimiter(limiter ratelimiter.RateLimiter,
	syncDuration time.Duration,
	jobGroup *jobs.JobGroup,
) (*LazySyncRateLimiter, error) {
	lsl := &LazySyncRateLimiter{
		limiter:      limiter,
		jobGroup:     jobGroup,
		name:         limiter.Name() + "-lazy-sync",
		syncDuration: syncDuration,
	}

	job := jobs.NewBasicJob(lsl.name, lsl.audit)
	// register job with job group
	err := lsl.jobGroup.RegisterJob(job, jobs.JobConfig{
		ExecutionPeriod:  config.MakeDuration(syncDuration * 4),
		ExecutionTimeout: config.MakeDuration(syncDuration * 4),
	})
	if err != nil {
		return nil, err
	}

	return lsl, nil
}

// Close closes the limiter.
func (lsl *LazySyncRateLimiter) Close() error {
	err := lsl.jobGroup.DeregisterJob(lsl.name)
	if err != nil {
		return err
	}
	lsl.limiter.Close()
	return nil
}

// GetRateLimit returns the rate limit.
func (lsl *LazySyncRateLimiter) GetRateLimit() float64 {
	return lsl.limiter.GetRateLimit()
}

// SetRateLimit sets the rate limit.
func (lsl *LazySyncRateLimiter) SetRateLimit(limit float64) {
	lsl.limiter.SetRateLimit(limit)
}

func (lsl *LazySyncRateLimiter) audit(ctx context.Context) (proto.Message, error) {
	now := time.Now()
	// range through the map and sync the counters
	lsl.counters.Range(func(label, value interface{}) bool {
		c := value.(*counter)
		c.lock.Lock()
		defer c.lock.Unlock()

		// if this counter has not synced in a while, then remove it from the map
		if now.After(c.nextSync.Add(lsl.syncDuration * 4)) {
			lsl.counters.Delete(label)
			return true
		}
		return true
	})
	return nil, nil
}

// Name returns the name of the limiter.
func (lsl *LazySyncRateLimiter) Name() string {
	return lsl.name
}

// TakeIfAvailable takes n tokens from the limiter.
func (lsl *LazySyncRateLimiter) TakeIfAvailable(label string, n float64) (bool, float64, float64) {
	if lsl.GetRateLimit() < 0 {
		return true, -1, -1
	}

	now := time.Now()
	syncRemote := func(c *counter, n float64) (bool, float64, float64) {
		ok, waitTime, remaining, current := lsl.limiter.Take(label, c.credit+n)
		c.credit = 0
		c.available = remaining
		if waitTime > 0 {
			c.nextSync = now.Add(waitTime)
		} else {
			c.nextSync = now.Add(lsl.syncDuration)
		}
		return ok, remaining, current
	}

	checkLimit := func(c *counter) (bool, float64, float64) {
		c.lock.Lock()
		defer c.lock.Unlock()
		// check if we need to sync
		if now.After(c.nextSync) {
			return syncRemote(c, n)
		}

		if c.credit+n > c.available {
			return false, math.Max(0, c.available-c.credit), c.credit
		}
		c.credit += n
		return true, math.Max(0, c.available-c.credit), c.credit
	}

	// check if the counter exists in the map
	if v, ok := lsl.counters.Load(label); ok {
		c := v.(*counter)
		return checkLimit(c)
	}
	// fallback to using the underlying limiter
	c := &counter{}
	c.lock.Lock()
	defer c.lock.Unlock()
	existing, loaded := lsl.counters.LoadOrStore(label, c)
	if loaded {
		ce := existing.(*counter)
		return checkLimit(ce)
	}
	return syncRemote(c, n)
}

// Take takes n tokens from the limiter.
func (lsl *LazySyncRateLimiter) Take(label string, n float64) (bool, time.Duration, float64, float64) {
	waitTime := time.Duration(0)
	ok, remaining, current := lsl.TakeIfAvailable(label, n)
	return ok, waitTime, remaining, current
}

// Make sure TokenBucketRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*LazySyncRateLimiter)(nil)
