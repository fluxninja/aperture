package ratelimiter

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/jobs"
)

type counter struct {
	local  int32
	global int64
}

// LazySyncRateLimiter is a limiter that syncs its state lazily with another limiter.
type LazySyncRateLimiter struct {
	counters      sync.Map
	limiter       RateLimiter
	jobGroup      *jobs.JobGroup
	name          string
	syncDuration  time.Duration
	totalCounters int64
}

// NewLazySyncRateLimiter creates a new LazySyncLimiter.
func NewLazySyncRateLimiter(limiter RateLimiter,
	syncDuration time.Duration,
	jobGroup *jobs.JobGroup,
) (RateLimiter, error) {
	lsl := &LazySyncRateLimiter{
		limiter:      limiter,
		jobGroup:     jobGroup,
		name:         limiter.Name() + "-lazy-sync",
		syncDuration: syncDuration,
	}

	job := &jobs.BasicJob{
		JobFunc: lsl.sync,
	}
	job.JobName = lsl.name
	// register job with job group
	err := lsl.jobGroup.RegisterJob(job, jobs.JobConfig{
		ExecutionPeriod: config.Duration{
			Duration: durationpb.New(syncDuration),
		},
		ExecutionTimeout: config.Duration{
			Duration: durationpb.New(syncDuration),
		},
		InitialDelay: config.Duration{
			Duration: durationpb.New(-1),
		},
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

func (lsl *LazySyncRateLimiter) sync(ctx context.Context) (proto.Message, error) {
	requestDelay := time.Duration(0)

	totalCount := atomic.LoadInt64(&lsl.totalCounters)
	if totalCount == 0 {
		return nil, nil
	}

	// get deadline
	deadline, ok := ctx.Deadline()
	if ok {
		// spread out requests over the deadline
		duration := time.Until(deadline)
		// reduce duration by 5ms to account for any processing overheads and the last few sync's
		duration -= 5 * time.Millisecond
		// if duration is less than 0, set it to 0
		if duration < 0 {
			duration = 0
		}
		requestDelay = time.Duration(int64(duration) / totalCount)
	}

	var i int64
	// range through the map and sync the counters
	lsl.counters.Range(func(label, value interface{}) bool {
		c := value.(*counter)

		// reset the local counter to 0 and asynchronously update the global counter if needed
		local := atomic.SwapInt32(&c.local, 0)
		// if local counter is 0, then remove the label from the map
		if local == 0 {
			// decrement total counters
			atomic.AddInt64(&lsl.totalCounters, -1)
			lsl.counters.Delete(label)
		} else {
			go func(i int64) {
				dur := time.Duration(i * int64(requestDelay))
				time.Sleep(dur)
				_, _, global := lsl.limiter.TakeN(label.(string), int(local))
				atomic.StoreInt64(&c.global, int64(global))
			}(i)
			i++
		}
		return true
	})
	return nil, nil
}

// Name returns the name of the limiter.
func (lsl *LazySyncRateLimiter) Name() string {
	return lsl.name
}

// TakeN takes n tokens from the limiter.
func (lsl *LazySyncRateLimiter) TakeN(label string, n int) (bool, int, int) {
	checkLimit := func(c *counter) (bool, int, int) {
		// atomic increment local counter
		local := atomic.AddInt32(&c.local, int32(n))
		total := int(local) + int(atomic.LoadInt64(&c.global))
		// check limit
		ok, remaining := lsl.limiter.GetRateLimitCheck().CheckRateLimit(label, total)
		return ok, remaining, total
	}

	// check if the counter exists in the map
	if v, ok := lsl.counters.Load(label); ok {
		c := v.(*counter)
		return checkLimit(c)
	}
	// fallback to using the underlying limiter
	ok, remaining, current := lsl.limiter.TakeN(label, n)
	c := &counter{
		local:  int32(0), // we have already taken these tokens from the underlying limiter
		global: int64(current),
	}
	existing, loaded := lsl.counters.LoadOrStore(label, c)
	if loaded {
		c := existing.(*counter)
		return checkLimit(c)
	} else {
		// increment total counters
		atomic.AddInt64(&lsl.totalCounters, 1)
	}
	return ok, remaining, current
}

// Take is a wrapper for TakeN(label, 1).
func (lsl *LazySyncRateLimiter) Take(label string) (bool, int, int) {
	return lsl.TakeN(label, 1)
}

// GetRateLimitCheck returns the limit checker of the limiter.
func (lsl *LazySyncRateLimiter) GetRateLimitCheck() RateLimitCheck {
	return lsl.limiter.GetRateLimitCheck()
}
