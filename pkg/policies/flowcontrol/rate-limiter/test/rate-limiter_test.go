package ratelimiter

import (
	"context"
	"fmt"
	stdlog "log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	"github.com/hashicorp/memberlist"

	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/distcache"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/log"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter"
	lazysync "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter/lazy-sync"
	tokenbucket "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter/token-bucket"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

func newTestLimiter(t *testing.T, distCache *distcache.DistCache, limit float64, interval time.Duration) (ratelimiter.RateLimiter, error) {
	limiter, err := tokenbucket.NewTokenBucket(distCache, "Limiter", interval, time.Hour, true)
	if err != nil {
		t.Logf("Failed to create DistCacheLimiter: %v", err)
		return nil, err
	}
	limiter.SetRateLimit(limit)
	limiter.SetFillAmount(limit)

	t.Logf("Successfully created new Limiter")
	return limiter, nil
}

func newTestDistCacheWithConfig(t *testing.T, c *olricconfig.Config) (*distcache.DistCache, error) {
	distCache := &distcache.DistCache{}

	ctx, cancel := context.WithCancel(context.Background())

	if c != nil {
		distCache.Config = c
	} else {
		oc := olricconfig.New("local")

		err := oc.DMaps.Sanitize()
		if err != nil {
			t.Errorf("Failed to sanitize olric config: %v", err)
		}
		err = oc.DMaps.Validate()
		if err != nil {
			t.Errorf("Failed to validate olric config: %v", err)
		}
		distCache.Config = oc
	}

	distCache.Config.Started = func() {
		t.Log("Started olric server")
		cancel()
	}

	o, err := olric.New(distCache.Config)
	if err != nil {
		return nil, err
	}

	distCache.Olric = o

	go func() {
		t.Log("Starting DistCacheLimiter")
		err = distCache.Olric.Start()
		if err != nil {
			t.Errorf("Failed to start olric: %v", err)
		}
	}()

	select {
	case <-time.After(time.Second):
		t.Fatal("DistCache cannot be started in one second")
	case <-ctx.Done():
		// everything is fine
	}

	return distCache, nil
}

type testDistCacheCluster struct {
	mu      sync.Mutex
	members map[string]*distcache.DistCache
}

func newTestOlricConfig() *olricconfig.Config {
	c := olricconfig.New("local")
	mc := memberlist.DefaultLocalConfig()
	mc.BindAddr = "127.0.0.1"
	mc.BindPort = 0
	c.MemberlistConfig = mc
	c.Logger = stdlog.New(&distcache.OlricLogWriter{Logger: log.GetGlobalLogger()}, "", 0)

	port, err := getFreePort()
	if err != nil {
		panic(fmt.Sprintf("GetFreePort returned an error: %v", err))
	}
	c.BindAddr = "127.0.0.1"
	c.BindPort = port
	c.MemberlistConfig.Name = net.JoinHostPort(c.BindAddr, strconv.Itoa(c.BindPort))
	if err := c.Sanitize(); err != nil {
		panic(fmt.Sprintf("failed to sanitize default config: %v", err))
	}
	return c
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func (cl *testDistCacheCluster) addDistCacheWithConfig(t *testing.T, c *olricconfig.Config) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if c == nil {
		c = newTestOlricConfig()
	}

	for _, member := range cl.members {
		c.Peers = append(c.Peers, fmt.Sprintf("%s:%d", member.Config.MemberlistConfig.BindAddr, member.Config.MemberlistConfig.BindPort))
	}

	dc, err := newTestDistCacheWithConfig(t, c)
	if err != nil {
		return err
	}
	cl.members[dc.Config.MemberlistConfig.Name] = dc
	return nil
}

func newTestDistCacheCluster(t *testing.T, n int) *testDistCacheCluster {
	cl := &testDistCacheCluster{
		members: make(map[string]*distcache.DistCache, n),
	}
	for i := 0; i < n; i++ {
		t.Log("Adding cluster member")
		_ = cl.addDistCacheWithConfig(t, nil)
	}
	return cl
}

type flow struct {
	requestlabel string // what label it needs
	requestRate  int32
	limit        int32
	// counters
	totalRequests    int32 // total requests made
	acceptedRequests int32 // total requests accepted
}

type flowRunner struct {
	wg       sync.WaitGroup
	limiters []ratelimiter.RateLimiter
	flows    []*flow
	duration time.Duration
}

// runFlows runs the flows for the given duration.
func (fr *flowRunner) runFlows(t *testing.T) {
	for _, f := range fr.flows {
		f.limit = int32(fr.limiters[0].GetRateLimit())

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
				accepted, _, _ := limiter.TakeIfAvailable(f.requestlabel, 1)
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
func createOlricLimiters(t *testing.T, cl *testDistCacheCluster, limit float64, interval time.Duration) []ratelimiter.RateLimiter {
	var limiters []ratelimiter.RateLimiter
	for _, distCache := range cl.members {
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
func createLazySyncLimiters(t *testing.T, limiters []ratelimiter.RateLimiter, syncDuration time.Duration) []ratelimiter.RateLimiter {
	var lazySyncLimiters []ratelimiter.RateLimiter
	for _, limiter := range limiters {
		jobGroup := createJobGroup(limiter)
		lazySyncLimiter, err := lazysync.NewLazySyncRateLimiter(limiter, syncDuration, jobGroup)
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
			actualRequestsExpected += f.limit
		}
		t.Logf("flow (%s) @ %d requests/sec: \n fills=%f, totalRequests=%d, limit=%d, acceptedRequests=%d, acceptedRequestsExpected=%d",
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

// closeTestDistCacheCluster closes the test dist cache cluster.
func closeTestDistCacheCluster(t *testing.T, cl *testDistCacheCluster) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	for _, member := range cl.members {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := member.Olric.Shutdown(ctx)
		if err != nil {
			t.Log("Failed to shutdown olric")
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
	syncDuration          time.Duration
	enableLazySyncLimiter bool
}

// baseOfLimiterTest is the base test for all limiter tests.
func baseOfLimiterTest(config testConfig) {
	var fr *flowRunner
	var lazySyncLimiters []ratelimiter.RateLimiter
	t := config.t
	cl := newTestDistCacheCluster(t, config.numOlrics)

	if len(cl.members) != config.numOlrics {
		t.Logf("Expected %d members, got %d", config.numOlrics, len(cl.members))
		t.FailNow()
	}

	limiters := createOlricLimiters(t, cl, config.limit, config.interval)

	if config.enableLazySyncLimiter {
		limiters = createLazySyncLimiters(t, limiters, config.syncDuration)
	}

	t.Log("Starting flows...")

	fr = &flowRunner{
		wg:       sync.WaitGroup{},
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
	closeTestDistCacheCluster(t, cl)
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
		syncDuration:          time.Millisecond * 100,
		tolerance:             0.2,
	})
}
