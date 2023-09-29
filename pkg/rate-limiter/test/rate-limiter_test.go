package ratelimiter

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/alerts"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/rate-limiter"
	globaltokenbucket "github.com/fluxninja/aperture/v2/pkg/rate-limiter/global-token-bucket"
	lazysync "github.com/fluxninja/aperture/v2/pkg/rate-limiter/lazy-sync"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

func newTestLimiter(t *testing.T, distCache *distcache.DistCache, limit float64, interval time.Duration) (ratelimiter.RateLimiter, error) {
	limiter, err := globaltokenbucket.NewGlobalTokenBucket(distCache, "Limiter", interval, time.Hour, true, false)
	if err != nil {
		t.Logf("Failed to create DistCacheLimiter: %v", err)
		return nil, err
	}
	limiter.SetBucketCapacity(limit)
	limiter.SetFillAmount(limit)
	limiter.SetPassThrough(false)

	t.Logf("Successfully created new Limiter")
	return limiter, nil
}

type flow struct {
	requestlabel string // what label it needs
	requestRate  int32
	limit        float64
	// counters
	totalRequests    int32 // total requests made
	acceptedRequests int32 // total requests accepted
}

type flowRunner struct {
	limiters []ratelimiter.RateLimiter
	flows    []*flow
	wg       sync.WaitGroup
	duration time.Duration
	limit    float64
}

// runFlows runs the flows for the given duration.
func (fr *flowRunner) runFlows(t *testing.T) {
	for _, f := range fr.flows {
		f.limit = fr.limit

		fr.wg.Add(1)
		go func(f *flow) {
			defer fr.wg.Done()

			stopTime := time.Now().Add(fr.duration)
			// delay between requests (in nanoseconds) = 1,000,000,000 / requestRate
			requestDelay := time.Duration(1e9 / f.requestRate)

			for {
				// nolint
				randomLimiterIndex := rand.Intn(len(fr.limiters))
				limiter := fr.limiters[randomLimiterIndex]
				atomic.AddInt32(&f.totalRequests, 1)
				accepted, _, _, _ := limiter.TakeIfAvailable(context.TODO(), f.requestlabel, 1)
				if accepted {
					atomic.AddInt32(&f.acceptedRequests, 1)
				}

				nextReqTime := time.Now().Add(requestDelay)
				if nextReqTime.Before(stopTime) {
					time.Sleep(requestDelay)
				} else {
					break
				}
			}
		}(f)
	}
	fr.wg.Wait()
}

// createJobGroup creates a job group for the given limiter..
func createJobGroup(_ ratelimiter.RateLimiter) *jobs.JobGroup {
	var gws jobs.GroupWatchers

	alerter := alerts.NewSimpleAlerter(100)
	reg := status.NewRegistry(log.GetGlobalLogger(), alerter).Child("test", "jobs")

	group, err := jobs.NewJobGroup(reg, jobs.JobGroupConfig{}, gws)
	if err != nil {
		panic(fmt.Sprintf("Failed to create job group: %v", err))
	}
	_ = group.Start()
	return group
}

// createOlricLimiters creates a set of Olric limiters.
func createOlricLimiters(t *testing.T, cl *distcache.TestDistCacheCluster, limit float64, interval time.Duration) []ratelimiter.RateLimiter {
	cl.Lock.Lock()
	defer cl.Lock.Unlock()
	var limiters []ratelimiter.RateLimiter
	for _, distCache := range cl.Members {
		limiter, err := newTestLimiter(t, distCache, limit, interval)
		if err != nil {
			t.Logf("Error creating limiter: %v", err)
			t.FailNow()
		}
		limiters = append(limiters, limiter)
	}
	return limiters
}

// createLazySyncLimiters creates a set of lazy-sync-limiters.
func createLazySyncLimiters(t *testing.T, limiters []ratelimiter.RateLimiter, interval time.Duration, numSyncs uint32) []ratelimiter.RateLimiter {
	var lazySyncLimiters []ratelimiter.RateLimiter
	for _, limiter := range limiters {
		jobGroup := createJobGroup(limiter)
		lazySyncLimiter, err := lazysync.NewLazySyncRateLimiter(limiter, interval, numSyncs, jobGroup)
		if err != nil {
			t.Logf("Error creating lazy sync limiter: %v", err)
			t.FailNow()
		}
		lazySyncLimiters = append(lazySyncLimiters, lazySyncLimiter)
	}
	return lazySyncLimiters
}

// checkResults checks if a certain number of requests were accepted under a given tolerance.
func checkResults(t *testing.T, fr *flowRunner, fills float64, tolerance float64) {
	for _, f := range fr.flows {
		// calculate expected requests taking into account the burst capacity in the limiter
		actualRequestsExpected := int32(math.Min(float64(f.limit)*(fills), float64(f.totalRequests)))
		if actualRequestsExpected != f.totalRequests {
			// add burst capacity
			actualRequestsExpected += int32(f.limit)
		}
		t.Logf("flow (%s) @ %d requests/sec: \n fills=%f, totalRequests=%d, limit=%f, acceptedRequests=%d, acceptedRequestsExpected=%d",
			f.requestlabel,
			f.requestRate,
			fills,
			f.totalRequests,
			f.limit,
			f.acceptedRequests,
			actualRequestsExpected)
		acceptedReqRatio := float64(f.acceptedRequests) / float64(actualRequestsExpected)
		if math.Abs(1-acceptedReqRatio) > tolerance {
			t.Logf("Accepted request ratio is %f, which is outside the tolerance of %f", acceptedReqRatio, tolerance)
			t.Fail()
		}
	}
}

// closeLimiters closes all the limiters.
func closeLimiters(t *testing.T, limiters []ratelimiter.RateLimiter) {
	for _, limiter := range limiters {
		err := limiter.Close()
		if err != nil {
			t.Logf("Failed to close Limiter: %v", err)
		}
	}
}

type testConfig struct {
	t                     *testing.T
	flows                 []*flow
	numOlrics             int
	limit                 float64
	interval              time.Duration
	tolerance             float64
	duration              time.Duration
	numSyncs              uint32
	enableLazySyncLimiter bool
}

// baseOfLimiterTest is the base test for all limiter tests.
func baseOfLimiterTest(config testConfig) {
	var fr *flowRunner
	var lazySyncLimiters []ratelimiter.RateLimiter
	t := config.t
	cl := distcache.NewTestDistCacheCluster(t, config.numOlrics)

	if len(cl.Members) != config.numOlrics {
		t.Logf("Expected %d members, got %d", config.numOlrics, len(cl.Members))
		t.FailNow()
	}

	limiters := createOlricLimiters(t, cl, config.limit, config.interval)

	if config.enableLazySyncLimiter {
		limiters = createLazySyncLimiters(t, limiters, config.interval, config.numSyncs)
	}

	t.Log("Starting flows...")

	fr = &flowRunner{
		wg:       sync.WaitGroup{},
		limit:    config.limit,
		limiters: limiters,
		flows:    config.flows,
		duration: config.duration,
	}

	start := time.Now()
	fr.runFlows(t)
	end := time.Now()

	fills := float64(end.Sub(start)) / float64(config.interval)

	checkResults(t, fr, fills, config.tolerance)

	if config.enableLazySyncLimiter {
		closeLimiters(t, lazySyncLimiters)
	}

	closeLimiters(t, limiters)
	distcache.CloseTestDistCacheCluster(t, cl)
}

// TestOlricLimiterWithBasicLimit tests the basic limit functionality of the limiter and if it accepts the limit of requests sent within interval.
func TestOlricLimiterWithBasicLimit(t *testing.T) {
	flows := []*flow{
		{requestlabel: "user-0", requestRate: 50},
	}
	baseOfLimiterTest(testConfig{
		t:         t,
		numOlrics: 1,
		limit:     10,
		interval:  time.Second * 1,
		flows:     flows,
		duration:  time.Second * 10,
		tolerance: 0.1,
	})
}

// TestOlricClusterMultiLimiter tests the behavior of a cluster of OlricLimiter and if it accepts the limit of requests sent within a given interval.
func TestOlricClusterMultiLimiter(t *testing.T) {
	flows := []*flow{
		{requestlabel: "user-0", requestRate: 200},
		{requestlabel: "user-1", requestRate: 30},
		{requestlabel: "user-2", requestRate: 50},
		{requestlabel: "user-3", requestRate: 90},
	}
	baseOfLimiterTest(testConfig{
		t:         t,
		numOlrics: 6,
		limit:     10,
		interval:  time.Second * 1,
		flows:     flows,
		duration:  time.Second * 10,
		tolerance: 0.1,
	})
}

// TestLazySyncClusterLimiter tests the lazy sync limiter which has a non-deterministic behavior and results may vary for each run.
// In order to pass the test, a 5% tolerance is allowed.
func TestLazySyncClusterLimiter(t *testing.T) {
	flows := []*flow{
		{requestlabel: "user-0", requestRate: 50},
		{requestlabel: "user-1", requestRate: 20},
	}

	baseOfLimiterTest(testConfig{
		t:                     t,
		numOlrics:             3,
		limit:                 10,
		interval:              time.Second * 1,
		flows:                 flows,
		duration:              time.Second * 10,
		enableLazySyncLimiter: true,
		numSyncs:              10,
		tolerance:             0.2,
	})
}
