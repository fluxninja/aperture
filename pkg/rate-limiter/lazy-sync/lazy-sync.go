package lazysync

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
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
	interval     time.Duration
	syncInterval time.Duration
}

// NewLazySyncRateLimiter creates a new LazySyncLimiter.
func NewLazySyncRateLimiter(limiter ratelimiter.RateLimiter,
	interval time.Duration,
	numSync uint32,
	jobGroup *jobs.JobGroup,
) (*LazySyncRateLimiter, error) {
	lazySyncInterval := time.Duration(int64(interval) / int64(numSync))
	lsl := &LazySyncRateLimiter{
		limiter:      limiter,
		jobGroup:     jobGroup,
		name:         limiter.Name() + "-lazy-sync",
		interval:     interval,
		syncInterval: lazySyncInterval,
	}

	job := jobs.NewBasicJob(lsl.name, lsl.audit)
	// register job with job group
	err := lsl.jobGroup.RegisterJob(job, jobs.JobConfig{
		ExecutionPeriod: config.MakeDuration(interval),
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

// SetPassThrough sets whether the limiter should pass through requests when it is not initialized.
func (lsl *LazySyncRateLimiter) SetPassThrough(passThrough bool) {
	lsl.limiter.SetPassThrough(passThrough)
}

// GetPassThrough returns whether the limiter should pass through requests when it is not initialized.
func (lsl *LazySyncRateLimiter) GetPassThrough() bool {
	return lsl.limiter.GetPassThrough()
}

func (lsl *LazySyncRateLimiter) audit(ctx context.Context) (proto.Message, error) {
	now := time.Now()
	// range through the map and sync the counters
	lsl.counters.Range(func(label, value interface{}) bool {
		c := value.(*counter)
		c.lock.Lock()
		defer c.lock.Unlock()

		// if this counter has not synced in a while, then remove it from the map
		if now.After(c.nextSync.Add(lsl.interval)) {
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

func (lsl *LazySyncRateLimiter) takeN(label string, n float64, canWait bool) (bool, time.Duration, float64, float64) {
	if lsl.GetPassThrough() {
		return true, 0, 0, 0
	}

	now := time.Now()
	syncRemote := func(c *counter, n float64) (bool, time.Duration, float64, float64) {
		toCheck := c.credit
		if canWait {
			toCheck += n
		}
		ok, waitTime, remaining, current := lsl.limiter.Take(label, toCheck)
		c.credit = 0
		if waitTime > 0 {
			c.nextSync = now.Add(waitTime)
		} else {
			c.nextSync = now.Add(lsl.syncInterval)
		}

		if ok && !canWait {
			// try to take n from the remaining
			if n > remaining {
				ok = false
			} else {
				c.credit = n
				remaining -= n
				current += n
			}
		}
		c.available = remaining
		return ok, waitTime, remaining, current
	}

	checkLimit := func(c *counter) (bool, time.Duration, float64, float64) {
		c.lock.Lock()
		defer c.lock.Unlock()
		// check if we need to sync
		if now.After(c.nextSync) {
			return syncRemote(c, n)
		}

		waitTime := time.Duration(0)

		if c.credit+n > c.available {
			if canWait {
				// need accurate wait time
				return syncRemote(c, n)
			}
			waitTime = c.nextSync.Sub(now)
			return false, waitTime, math.Max(0, c.available-c.credit), c.credit
		}
		c.credit += n
		return true, waitTime, math.Max(0, c.available-c.credit), c.credit
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

// TakeIfAvailable takes n tokens from the limiter if they are available.
func (lsl *LazySyncRateLimiter) TakeIfAvailable(label string, n float64) (bool, float64, float64) {
	ok, _, remaining, current := lsl.takeN(label, n, false)
	return ok, remaining, current
}

// Take takes n tokens from the limiter.
func (lsl *LazySyncRateLimiter) Take(label string, n float64) (bool, time.Duration, float64, float64) {
	return lsl.takeN(label, n, true)
}

// Return returns n tokens to the limiter.
func (lsl *LazySyncRateLimiter) Return(label string, n float64) (float64, float64) {
	if lsl.GetPassThrough() {
		return 0, 0
	}
	_, remaining, current := lsl.TakeIfAvailable(label, -n)
	return remaining, current
}

// Make sure TokenBucketRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*LazySyncRateLimiter)(nil)
