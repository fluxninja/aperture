package lazysync

import (
	"context"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
)

type counter struct {
	nextSync  time.Time
	credit    float64
	current   float64
	available float64
	lock      sync.Mutex
	waiting   bool
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

func (lsl *LazySyncRateLimiter) takeN(ctx context.Context, label string, n float64, canWait bool) (bool, time.Duration, float64, float64) {
	if lsl.GetPassThrough() {
		return true, 0, 0, 0
	}

	now := time.Now()
	syncRemote := func(c *counter, n float64) (bool, time.Duration, float64, float64) {
		tokens := c.credit
		if canWait {
			tokens += n
		}
		ok, waitTime, remaining, current := lsl.limiter.Take(ctx, label, tokens)
		c.credit = 0
		if waitTime > 0 {
			c.waiting = true
			c.nextSync = now.Add(waitTime)
		} else {
			c.waiting = false
			c.nextSync = now.Add(lsl.syncInterval)
		}

		if ok && !canWait {
			if n <= remaining {
				c.credit = n
				remaining -= n
				current += n
			} else {
				ok = false
			}
		}
		c.available = remaining
		c.current = current
		return ok, waitTime, remaining, current
	}

	checkLimit := func(c *counter) (bool, time.Duration, float64, float64) {
		c.lock.Lock()
		defer c.lock.Unlock()

		ret := func(ok bool, waitTime time.Duration) (bool, time.Duration, float64, float64) {
			return ok, waitTime, c.available - c.credit, c.current + c.credit
		}

		if n <= 0 {
			c.credit += n
			return ret(true, 0)
		}

		if now.After(c.nextSync) {
			return syncRemote(c, n)
		}

		waitTime := time.Duration(0)

		if c.waiting && !canWait {
			return ret(false, 0)
		}

		if c.waiting {
			waitTime = c.nextSync.Sub(now)
		}

		if c.credit+n <= c.available {
			c.credit += n
			return ret(true, waitTime)
		}

		if canWait && !c.waiting {
			return syncRemote(c, n)
		}

		return ret(false, waitTime)
	}

	c := &counter{}
	existing, loaded := lsl.counters.LoadOrStore(label, c)
	if loaded {
		c = existing.(*counter)
	}
	return checkLimit(c)
}

// TakeIfAvailable takes n tokens from the limiter if they are available.
func (lsl *LazySyncRateLimiter) TakeIfAvailable(ctx context.Context, label string, n float64) (bool, time.Duration, float64, float64) {
	return lsl.takeN(ctx, label, n, false)
}

// Take takes n tokens from the limiter.
func (lsl *LazySyncRateLimiter) Take(ctx context.Context, label string, n float64) (bool, time.Duration, float64, float64) {
	return lsl.takeN(ctx, label, n, true)
}

// Return returns n tokens to the limiter.
func (lsl *LazySyncRateLimiter) Return(ctx context.Context, label string, n float64) (float64, float64) {
	_, _, remaining, current := lsl.takeN(ctx, label, -n, false)
	return remaining, current
}

// Make sure TokenBucketRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*LazySyncRateLimiter)(nil)
