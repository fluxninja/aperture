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

func newTestLimiter(t *testing.T, distCache *distcache.DistCache, config testConfig) (ratelimiter.RateLimiter, error) {
	limiter, err := globaltokenbucket.NewGlobalTokenBucket(distCache, "Limiter", config.interval, time.Hour, config.continuousFill, config.delayInitialFill)
	if err != nil {
		t.Logf("Failed to create DistCacheLimiter: %v", err)
		return nil, err
	}
	limiter.SetBucketCapacity(config.bucketCapacity)
	limiter.SetFillAmount(config.fillAmount)
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
func createOlricLimiters(t *testing.T, cl *distcache.TestDistCacheCluster, config testConfig) []ratelimiter.RateLimiter {
	cl.Lock.Lock()
	defer cl.Lock.Unlock()
	var limiters []ratelimiter.RateLimiter
	for _, distCache := range cl.Members {
		limiter, err := newTestLimiter(t, distCache, config)
		if err != nil {
			t.Logf("Error creating limiter: %v", err)
			t.FailNow()
		}
		limiters = append(limiters, limiter)
	}
	return limiters
}

// createLazySyncLimiters creates a set of lazy-sync-limiters.
func createLazySyncLimiters(t *testing.T, limiters []ratelimiter.RateLimiter, config testConfig) []ratelimiter.RateLimiter {
	var lazySyncLimiters []ratelimiter.RateLimiter
	for _, limiter := range limiters {
		jobGroup := createJobGroup(limiter)
		lazySyncLimiter, err := lazysync.NewLazySyncRateLimiter(limiter, config.interval, config.numSyncs, jobGroup)
		if err != nil {
			t.Logf("Error creating lazy sync limiter: %v", err)
			t.FailNow()
		}
		lazySyncLimiters = append(lazySyncLimiters, lazySyncLimiter)
	}
	return lazySyncLimiters
}

// checkResults checks if a certain number of requests were accepted under a given tolerance.
func checkResults(t *testing.T, fr *flowRunner, duration time.Duration, config testConfig) {
	t.Logf("duration: %v", duration)
	// if delayedFilling is enabled, then subtract the time it takes to fill the bucket
	// from the duration
	if config.delayInitialFill {
		if config.continuousFill {
			timeToFillBucket := time.Duration(config.bucketCapacity/config.fillAmount) * config.interval
			duration -= timeToFillBucket
		} else {
			// find fills needed to fill the bucket
			fills := math.Ceil(config.bucketCapacity/config.fillAmount) - 1
			timeToFillBucket := time.Duration(fills) * config.interval
			duration -= timeToFillBucket
		}
	}

	availableAmount := config.bucketCapacity
	if config.continuousFill {
		// assume continuous filling of limit amount per interval
		// take into account end and start times
		availableAmount += float64(duration) / float64(config.interval) * config.fillAmount
	} else {
		// assume filling at the end of each interval
		// find out exact number of fills (integer) that happened
		// take into account end and start times
		fills := math.Floor(float64(duration) / float64(config.interval))
		availableAmount += float64(fills) * config.fillAmount
	}

	for _, f := range fr.flows {
		acceptedRequestsExpected := int32(math.Min(availableAmount, float64(f.totalRequests)))

		t.Logf("flow (%s) @ %d requests/sec: \n fillAmount=%f, totalRequests=%d, limit=%f, acceptedRequests=%d, acceptedRequestsExpected=%d",
			f.requestlabel,
			f.requestRate,
			availableAmount,
			f.totalRequests,
			f.limit,
			f.acceptedRequests,
			acceptedRequestsExpected)
		acceptedReqRatio := float64(f.acceptedRequests) / float64(acceptedRequestsExpected)
		if math.Abs(1-acceptedReqRatio) > config.tolerance {
			t.Logf("Accepted request ratio is %f, which is outside the tolerance of %f", acceptedReqRatio, config.tolerance)
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
	fillAmount            float64
	bucketCapacity        float64
	interval              time.Duration
	tolerance             float64
	duration              time.Duration
	numSyncs              uint32
	enableLazySyncLimiter bool
	continuousFill        bool
	delayInitialFill      bool
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

	limiters := createOlricLimiters(t, cl, config)

	if config.enableLazySyncLimiter {
		limiters = createLazySyncLimiters(t, limiters, config)
	}

	t.Log("Starting flows...")

	fr = &flowRunner{
		wg:       sync.WaitGroup{},
		limit:    config.fillAmount,
		limiters: limiters,
		flows:    config.flows,
		duration: config.duration,
	}

	start := time.Now()
	fr.runFlows(t)
	end := time.Now()
	duration := end.Sub(start)

	checkResults(t, fr, duration, config)

	if config.enableLazySyncLimiter {
		closeLimiters(t, lazySyncLimiters)
	}

	closeLimiters(t, limiters)
	distcache.CloseTestDistCacheCluster(t, cl)
}

type combination struct {
	continuousFill   bool
	delayInitialFill bool
}

// vary combinations of continuousFill and delayInitialFill
var combinations = []combination{
	{true, true},
	{true, false},
	{false, true},
	{false, false},
}

// TestOlricLimiterWithBasicLimit tests the basic limit functionality of the limiter and if it accepts the limit of requests sent within interval.
func TestOlricLimiterWithBasicLimit(t *testing.T) {
	for _, c := range combinations {
		c := c // capture range variable
		t.Run(fmt.Sprintf("continuousFill=%v,delayInitialFill=%v", c.continuousFill, c.delayInitialFill), func(t *testing.T) {
			t.Parallel() // run subtests in parallel
			flows := []*flow{
				{requestlabel: "user-0", requestRate: 50},
			}
			baseOfLimiterTest(testConfig{
				t:                t,
				numOlrics:        1,
				fillAmount:       10,
				bucketCapacity:   30,
				interval:         time.Second * 1,
				flows:            flows,
				duration:         time.Second*10 - time.Millisecond*100,
				tolerance:        0.02,
				continuousFill:   c.continuousFill,
				delayInitialFill: c.delayInitialFill,
			})
		})
	}
}

// TestOlricClusterMultiLimiter tests the behavior of a cluster of OlricLimiter and if it accepts the limit of requests sent within a given interval.
func TestOlricClusterMultiLimiter(t *testing.T) {
	for _, c := range combinations {
		c := c // capture range variable
		t.Run(fmt.Sprintf("continuousFill=%v,delayInitialFill=%v", c.continuousFill, c.delayInitialFill), func(t *testing.T) {
			t.Parallel() // marks each subtest to run in parallel
			flows := []*flow{
				{requestlabel: "user-0", requestRate: 200},
				{requestlabel: "user-1", requestRate: 30},
				{requestlabel: "user-2", requestRate: 50},
				{requestlabel: "user-3", requestRate: 90},
			}
			baseOfLimiterTest(testConfig{
				t:                t,
				numOlrics:        6,
				fillAmount:       10,
				bucketCapacity:   30,
				interval:         time.Second * 1,
				flows:            flows,
				duration:         time.Second*10 - time.Millisecond*100,
				tolerance:        0.02,
				continuousFill:   c.continuousFill,
				delayInitialFill: c.delayInitialFill,
			})
		})
	}
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
		fillAmount:            10,
		bucketCapacity:        30,
		interval:              time.Second * 1,
		flows:                 flows,
		duration:              time.Second*10 - time.Millisecond*100,
		enableLazySyncLimiter: true,
		numSyncs:              10,
		tolerance:             0.1,
		continuousFill:        true,
		delayInitialFill:      true,
	})
}
