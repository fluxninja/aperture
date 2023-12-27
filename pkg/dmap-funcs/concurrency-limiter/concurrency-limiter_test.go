package concurrencylimiter_test

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	concurrencylimiter "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/concurrency-limiter"
)

type flow struct {
	requestlabel string
	requestRate  int32

	// config info
	capacity float64

	// counters
	totalRequests    int32 // total requests made
	acceptedRequests int32 // total requests accepted
}

type testConfig struct {
	// test config
	t         *testing.T
	numOlrics int

	// limiter config
	capacity            float64
	maxInflightDuration time.Duration

	// flow config
	flows     []*flow
	tolerance float64
	duration  time.Duration
	numSyncs  uint32
}

func newTestLimiter(t *testing.T, distCache *distcache.DistCache, config testConfig) (concurrencylimiter.ConcurrencyLimiter, error) {
	limiter, err := concurrencylimiter.NewGlobalTokenCounter(distCache, "Limiter", time.Hour, config.maxInflightDuration)
	if err != nil {
		t.Logf("Failed to create GlobalTokenCounter: %v", err)
		return nil, err
	}
	limiter.SetCapacity(config.capacity)
	limiter.SetPassThrough(false)
	t.Logf("Successfully created new ConcurrencyLimiter")
	return limiter, nil
}

type flowRunner struct {
	wg sync.WaitGroup

	limiters []concurrencylimiter.ConcurrencyLimiter
	flows    []*flow
	duration time.Duration
	capacity float64
}

// createLimiters creates a set of Olric limiters.
func createLimiters(t *testing.T, cl *distcache.TestDistCacheCluster, config testConfig) []concurrencylimiter.ConcurrencyLimiter {
	cl.Lock.Lock()
	defer cl.Lock.Unlock()
	var limiters []concurrencylimiter.ConcurrencyLimiter
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

// runTest is the base test for all limiter tests.
func runTest(config testConfig) {
	var fr *flowRunner
	t := config.t
	cl := distcache.NewTestDistCacheCluster(t, config.numOlrics)

	if len(cl.Members) != config.numOlrics {
		t.Logf("Expected %d members, got %d", config.numOlrics, len(cl.Members))
		t.FailNow()
	}

	limiters := createLimiters(t, cl, config)

	t.Log("Starting flows")

	fr = &flowRunner{
		wg:       sync.WaitGroup{},
		capacity: config.capacity,
		limiters: limiters,
		flows:    config.flows,
		duration: config.duration,
	}

	start := time.Now()
	fr.runFlows(t)
	end := time.Now()
	duration := end.Sub(start)

	checkResults(t, fr, duration, config)

	closeLimiters(t, limiters)
	distcache.CloseTestDistCacheCluster(t, cl)
}

// runFlows runs the flows for the given duration.
func (fr *flowRunner) runFlows(t *testing.T) {
	for _, f := range fr.flows {
		f.capacity = fr.capacity

		fr.wg.Add(1)
		go func(f *flow) {
			defer fr.wg.Done()

			stopTime := time.Now().Add(fr.duration)
			// delay between requests (in nanoseconds) = 1,000,000,000 / requestRate
			requestDelay := time.Duration(1e9 / f.requestRate)

			for {
				randomLimiterIndex := rand.Intn(len(fr.limiters))
				limiter := fr.limiters[randomLimiterIndex]
				atomic.AddInt32(&f.totalRequests, 1)
				accepted, waitTime, remaining, current, reqID := limiter.TakeIfAvailable(context.TODO(), f.requestlabel, 1)
				if accepted {
					atomic.AddInt32(&f.acceptedRequests, 1)
				}
				t.Logf("accepted=%v, waitTime=%v, remaining=%v, current=%v, reqID=%v", accepted, waitTime, remaining, current, reqID)

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

// checkResults checks if a certain number of requests were accepted under a given tolerance.
func checkResults(t *testing.T, fr *flowRunner, duration time.Duration, config testConfig) {
	// availableAmount := config.capacity

	for _, f := range fr.flows {
		acceptedRequestsExpected := int32(math.Min(config.capacity, float64(f.totalRequests)))

		t.Logf("flow (%s) @ %d requests/sec for %v secs: \n totalRequests=%d, capacity=%f, acceptedRequests=%d, acceptedRequestsExpected=%d",
			f.requestlabel,
			f.requestRate,
			duration,
			f.totalRequests,
			f.capacity,
			f.acceptedRequests,
			acceptedRequestsExpected,
		)
		acceptedReqRatio := float64(f.acceptedRequests) / float64(acceptedRequestsExpected)
		if math.Abs(1-acceptedReqRatio) > config.tolerance {
			t.Logf("Accepted request ratio is %f, which is outside the tolerance of %f", acceptedReqRatio, config.tolerance)
			t.Fail()
		}
	}
}

// closeLimiters closes all the limiters.
func closeLimiters(t *testing.T, limiters []concurrencylimiter.ConcurrencyLimiter) {
	for _, limiter := range limiters {
		err := limiter.Close()
		if err != nil {
			t.Logf("Failed to close Limiter: %v", err)
		}
	}
}

// TestLimiterWithBasicLimit tests the limiter with a basic limit.
func TestLimiterWithBasicLimit(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		// run subtests in parallel
		t.Parallel()
		flows := []*flow{
			{
				requestlabel: "user-0",
				requestRate:  2,
			},
		}
		runTest(testConfig{
			t:                   t,
			numOlrics:           1,
			flows:               flows,
			duration:            (time.Second * 10) - (time.Millisecond * 100),
			tolerance:           0.02,
			capacity:            10,
			maxInflightDuration: 7200 * time.Second,
		})
	})
}
