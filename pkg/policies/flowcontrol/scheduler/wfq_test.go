package scheduler

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
)

const (
	_testTolerance    = 0.15
	_testSlotDuration = time.Millisecond * 100
	_testSlotCount    = 10
)

var (
	prometheusRegistry              *prometheus.Registry
	wfqFlowsGauge                   prometheus.Gauge
	wfqHeapRequestsGauge            prometheus.Gauge
	tokenBucketLMGauge              prometheus.Gauge
	tokenBucketFillRateGauge        prometheus.Gauge
	tokenBucketBucketCapacityGauge  prometheus.Gauge
	tokenBucketAvailableTokensGauge prometheus.Gauge
)

func getMetrics() *TokenBucketLoadMultiplierMetrics {
	prometheusRegistry = prometheus.NewRegistry()

	constLabels := make(prometheus.Labels)

	wfqFlowsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        metrics.WFQFlowsMetricName,
		Help:        "A gauge that tracks the number of flows in the WFQScheduler",
		ConstLabels: constLabels,
	})
	_ = prometheusRegistry.Register(wfqFlowsGauge)
	wfqHeapRequestsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        metrics.WFQRequestsMetricName,
		Help:        "A gauge that tracks the number of queued requests in the WFQScheduler",
		ConstLabels: constLabels,
	})
	_ = prometheusRegistry.Register(wfqHeapRequestsGauge)
	tokenBucketLMGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metrics.TokenBucketLMMetricName,
		Help: "A gauge that tracks the load multiplier",
	})
	_ = prometheusRegistry.Register(tokenBucketLMGauge)
	tokenBucketFillRateGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metrics.TokenBucketFillRateMetricName,
		Help: "A gauge that tracks the fill rate of token bucket",
	})
	_ = prometheusRegistry.Register(tokenBucketFillRateGauge)
	tokenBucketBucketCapacityGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metrics.TokenBucketCapacityMetricName,
		Help: "A gauge that tracks the capacity of token bucket",
	})
	_ = prometheusRegistry.Register(tokenBucketBucketCapacityGauge)
	tokenBucketAvailableTokensGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metrics.TokenBucketAvailableMetricName,
		Help: "A gauge that tracks the number of tokens available in token bucket",
	})
	_ = prometheusRegistry.Register(tokenBucketAvailableTokensGauge)
	tbbMetrics := &TokenBucketMetrics{
		FillRateGauge:        tokenBucketFillRateGauge,
		BucketCapacityGauge:  tokenBucketBucketCapacityGauge,
		AvailableTokensGauge: tokenBucketAvailableTokensGauge,
	}
	metrics := &TokenBucketLoadMultiplierMetrics{
		LMGauge:            tokenBucketLMGauge,
		TokenBucketMetrics: tbbMetrics,
	}
	return metrics
}

type flowTracker struct {
	fairnessLabel string // what label it needs
	tokens        uint64 // how many tokens it needs
	priority      uint8
	timeout       time.Duration
	requestRate   uint64
	// counters
	totalRequests    uint64 // total requests made
	acceptedRequests uint64 // total requests accepted
}

func (flow *flowTracker) String() string {
	return fmt.Sprintf("FlowTracker<"+
		"Flow<FairnessLabel: %s "+
		"| RequestTokens: %d "+
		"| Priority: %d "+
		"| Timeout: %v "+
		"| RequestRate: %d "+
		"| TotalRequests: %d "+
		"| AcceptedRequests: %d "+
		"| PercentageSuccess: %0.2f "+
		">",
		flow.fairnessLabel,
		flow.tokens,
		flow.priority,
		flow.timeout,
		flow.requestRate,
		flow.totalRequests,
		flow.acceptedRequests,
		float64(flow.acceptedRequests)/float64(flow.totalRequests)*100,
	)
}

type flowTrackers []*flowTracker

// // // Ensures clock is updated periodically
func updateClock(t *testing.T, clk clockwork.FakeClock, flows flowTrackers) {
	// request delay
	minRequestDelay := 1 * time.Second
	// loop through flows
	for _, flow := range flows {
		requestDelay := time.Duration(1e9 / flow.requestRate)
		if requestDelay < minRequestDelay {
			minRequestDelay = requestDelay
		}
		if flow.timeout != 0 && flow.timeout < minRequestDelay {
			minRequestDelay = flow.timeout
		}
	}
	t.Logf("FakeClock ticks - minRequestDelay: %v\n", minRequestDelay)
	for {
		time.Sleep(1 * time.Millisecond)
		clk.Advance(minRequestDelay)
	}
}

func runFlows(sched Scheduler, wg *sync.WaitGroup, flows flowTrackers, duration time.Duration, clk clockwork.Clock) {
	for _, flow := range flows {
		// reset counters
		flow.totalRequests = 0
		flow.acceptedRequests = 0
		wg.Add(1)
		go runFlow(sched, wg, flow, duration, clk)
	}
}

func runFlow(sched Scheduler, wg *sync.WaitGroup, flow *flowTracker, duration time.Duration, clk clockwork.Clock) {
	defer wg.Done()
	// run each configured flow
	stopTime := clk.Now().Add(duration)
	requestDelay := time.Duration(1e9 / flow.requestRate)

	for {
		atomic.AddUint64(&flow.totalRequests, 1)
		wg.Add(1)
		go runRequest(sched, wg, flow)
		requestDelay := requestDelay
		nextReqTime := clk.Now().Add(requestDelay)
		if nextReqTime.Before(stopTime) {
			clk.Sleep(requestDelay)
		} else {
			return
		}
	}
}

func runRequest(sched Scheduler, wg *sync.WaitGroup, flow *flowTracker) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), flow.timeout)
	defer cancel()
	ok := sched.Schedule(ctx, flow.makeRequest())
	if ok {
		atomic.AddUint64(&flow.acceptedRequests, 1)
	}
}

func (flow *flowTracker) makeRequest() Request {
	return Request{
		FairnessLabel: flow.fairnessLabel,
		Tokens:        flow.tokens,
		Priority:      flow.priority,
	}
}

func printPrettyFlowTracker(t *testing.T, flows flowTrackers) {
	var totalRequests uint64
	var totalAccepted uint64

	for _, flow := range flows {
		totalRequests += flow.totalRequests
		totalAccepted += flow.acceptedRequests
	}

	t.Logf("\n\n\nSummary Statistics:\n")
	t.Logf("totalRequests: %d | totalAccepted: %d | successRate: %.02f\n", totalRequests, totalAccepted, float64(totalAccepted)/float64(totalRequests)*100)
}

// calculate tokenRate for flowTrackers
func tokenRate(flows flowTrackers) uint64 {
	var totalTokenRate uint64
	for _, flow := range flows {
		totalTokenRate += flow.requestRate * (flow.tokens)
	}
	return totalTokenRate
}

// Wait for some time
func Time(duration string) {
	sleep, err := time.ParseDuration(duration)
	if err != nil {
		log.Panic().Err(err).Msgf("Failed at parsing duration: %v", err)
	}
	time.Sleep(sleep)
}

// ------------------------- Benchmark Testing -------------------------
func BenchmarkBasicTokenBucket(b *testing.B) {
	flows := flowTrackers{
		{fairnessLabel: "workload1", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload2", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload3", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload4", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload5", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
	}
	c := clockwork.NewRealClock()
	startTime := c.Now()
	manager := NewBasicTokenBucket(startTime, 0, getMetrics().TokenBucketMetrics)

	schedMetrics := &WFQMetrics{
		FlowsGauge:        wfqFlowsGauge,
		HeapRequestsGauge: wfqHeapRequestsGauge,
	}
	sched := NewWFQScheduler(manager, c, schedMetrics)

	b.Logf("iterations: %d", b.N)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), flows[i%len(flows)].timeout)
			_ = sched.Schedule(ctx, flows[i%len(flows)].makeRequest())
			cancel()
		}
	})
}

func BenchmarkTokenBucketLoadMultiplier(b *testing.B) {
	flows := flowTrackers{
		{fairnessLabel: "workload1", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload2", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload3", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload4", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload5", tokens: 1, priority: 0, timeout: 5 * time.Millisecond},
	}
	c := clockwork.NewRealClock()
	startTime := c.Now()
	manager := NewTokenBucketLoadMultiplier(startTime, _testSlotCount, _testSlotDuration, getMetrics())
	manager.SetContinuousTracking(true)
	manager.SetLoadMultiplier(startTime, 1.0)

	schedMetrics := &WFQMetrics{
		FlowsGauge:        wfqFlowsGauge,
		HeapRequestsGauge: wfqHeapRequestsGauge,
	}
	sched := NewWFQScheduler(manager, c, schedMetrics)

	// bootstrap bucket
	bootstrapTime := time.Second * 1
	for c.Now().Before(startTime.Add(bootstrapTime)) {
		for _, flow := range flows {
			ctx, cancel := context.WithTimeout(context.Background(), flow.timeout)
			_ = sched.Schedule(ctx, flow.makeRequest())
			cancel()
		}
	}

	b.Logf("iterations: %d", b.N)

	// make sure tokens a large negative number
	manager.tbb.addTokens(-math.MaxFloat64)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), flows[i%len(flows)].timeout)
			_ = sched.Schedule(ctx, flows[i%len(flows)].makeRequest())
			cancel()
		}
	})
}

func totalSentTokens(flows flowTrackers) []uint64 {
	var total uint64
	totalTokens := make([]uint64, len(flows))
	for i, flow := range flows {
		totalTokens[i] = flow.totalRequests * flow.tokens
		total += totalTokens[i]
	}
	return totalTokens
}

func calculateFillRate(flows flowTrackers, lm float64) float64 {
	fillRate := float64(0)
	for _, flow := range flows {
		fillRate += float64(flow.tokens * flow.requestRate)
	}
	return fillRate * lm
}

func baseOfBasicBucketTest(t *testing.T, flows flowTrackers, fillRate float64, noOfRuns int) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := clockwork.NewFakeClock()
	go updateClock(t, c, flows)

	startTime := c.Now()

	basicBucket := NewBasicTokenBucket(startTime, fillRate, getMetrics().TokenBucketMetrics)
	metrics := &WFQMetrics{
		FlowsGauge:        wfqFlowsGauge,
		HeapRequestsGauge: wfqHeapRequestsGauge,
	}
	sched := NewWFQScheduler(basicBucket, c, metrics)
	var wg sync.WaitGroup
	var acceptedTokenRatio float64

	// ------------------------- Run Experiment -------------------------

	flowRunTime := time.Second * 10

	sumPriority := float64(0)
	minInvPrio := uint8(math.MaxUint8)
	minPrio := uint8(math.MaxUint8)
	for _, flow := range flows {
		if minInvPrio > (math.MaxUint8 - flow.priority) {
			minInvPrio = math.MaxUint8 - flow.priority
		}
		if minPrio > flow.priority {
			minPrio = flow.priority
		}
	}
	t.Logf("minPrio: %d, minInvPrio: %d\n", minPrio, minInvPrio)
	adjustedPrio := make([]uint16, len(flows))
	for i, flow := range flows {
		adjustedPrio[i] = uint16(flow.priority-minPrio+minInvPrio) + 1
		sumPriority += float64(adjustedPrio[i])
	}

	// Estimate the tokens -- It's a rough approach but seems to work so far for variety of loads and priorities
	totalEstimatedtokens := uint64(0)
	estimatedTokens := make([]uint64, len(flows))

	for i := range flows {
		estimatedTokens[i] = uint64(fillRate*float64(adjustedPrio[i])/sumPriority) * uint64(flowRunTime.Seconds())
		totalEstimatedtokens += estimatedTokens[i]
	}
	t.Logf("Flows: %v", flows)
	t.Logf("Adjusted Prios: %v\n", adjustedPrio)
	t.Logf("Estimated minimum tokens per flow before run are: %v\n", estimatedTokens)
	t.Logf("Total estimated allocated tokens: %d\n", totalEstimatedtokens)

	for i := 0; i < noOfRuns; i++ {
		runFlows(sched, &wg, flows, flowRunTime, c)
		wg.Wait()

		totalTokens := make([]uint64, len(flows))
		acceptedTokenSum := uint64(0)
		acceptedTokens := make([]uint64, len(flows))
		for i, flow := range flows {
			acceptedTokens[i] = flow.acceptedRequests * flow.tokens
			acceptedTokenSum += acceptedTokens[i]
			totalTokens[i] = flow.totalRequests * flow.tokens
		}
		t.Logf("Tokens sent per flow: %v\n", totalTokens)
		t.Logf("Total accepted tokens per flow after run are: %v\n", acceptedTokens)
		t.Logf("Total accepted tokens after run are: %v\n", acceptedTokenSum)

		for i := range acceptedTokens {
			// Will not get an accurate reading if traffic rate is very low
			if totalTokens[i] > 20 {
				ratio := float64(acceptedTokens[i]) / float64(totalTokens[i])
				if (ratio < 1) && math.Abs(1-ratio) > _testTolerance {
					ratio := float64(acceptedTokens[i]) / float64(estimatedTokens[i])
					if (ratio < 1) && math.Abs(1-ratio) > _testTolerance {
						t.Logf("Test failed because of more than %f percent unfairness %f: acceptedTokens: %d, allocatedTokens: %d\n", _testTolerance*100, math.Abs(1-ratio), acceptedTokens[i], estimatedTokens[i])
						t.Fail()
					}
				}
			}
		}
		acceptedTokenRatio = float64(acceptedTokenSum) / float64(fillRate*float64(flowRunTime.Seconds()))
		t.Logf("Accepted token ratio: %f\n", acceptedTokenRatio)
		if math.Abs(1-acceptedTokenRatio) > _testTolerance {
			t.Logf("Test failed because of more than %f percent unfairness %f\n", _testTolerance*100, math.Abs(1-acceptedTokenRatio))
			t.Fail()
		}
		if sched.(*WFQScheduler).queueOpen {
			t.Logf("Test failed because queue has some elements stuck in it\n")
			t.Fail()
		}
		if sched.(*WFQScheduler).GetPendingFlows() > 0 {
			t.Logf("Test failed because scheduler has some flow stuck in it\n")
			t.Fail()
		}
		if sched.(*WFQScheduler).GetPendingRequests() > 0 {
			t.Logf("Test failed because some request is still stuck in the heap\n")
			t.Fail()
		}
		if sched.(*WFQScheduler).generation < uint64(i+1) {
			t.Logf("There are less generations than expected\n")
			t.Fail()
		}
	}
	stopTime := c.Now()
	basicBucket.SetFillRate(stopTime, 0)
	if basicBucket.GetFillRate() != 0 {
		t.Logf("Fill rate is not 0 after stop\n")
		t.Fail()
	}
	t.Logf("Testing end\nDuration: %s", stopTime.Sub(startTime).String())
}

func TestHighRpsFlows(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 5, priority: 0, requestRate: 100},
		{fairnessLabel: "workload1", tokens: 5, priority: 0, requestRate: 100},
		{fairnessLabel: "workload2", tokens: 5, priority: 0, requestRate: 100},
		{fairnessLabel: "workload3", tokens: 5, priority: 50, requestRate: 100},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestLowRpsFlows(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 20, priority: 2, requestRate: 50},
		{fairnessLabel: "workload1", tokens: 20, priority: 3, requestRate: 50},
		{fairnessLabel: "workload2", tokens: 20, priority: 3, requestRate: 50},
		{fairnessLabel: "workload3", tokens: 20, priority: 2, requestRate: 50},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestMixRpsFlows(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 10, priority: 5, requestRate: 100},
		{fairnessLabel: "workload1", tokens: 15, priority: 4, requestRate: 200},
		{fairnessLabel: "workload2", tokens: 15, priority: 3, requestRate: 80},
		{fairnessLabel: "workload3", tokens: 10, priority: 2, requestRate: 120},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestSingleHighRequest(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 15, priority: 20, requestRate: 200},
		{fairnessLabel: "workload1", tokens: 5, priority: 1, requestRate: 50},
		{fairnessLabel: "workload2", tokens: 5, priority: 1, requestRate: 50},
		{fairnessLabel: "workload3", tokens: 5, priority: 1, requestRate: 50},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestSingleLowRequest(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 1, priority: 50, requestRate: 20},
		{fairnessLabel: "workload1", tokens: 8, priority: 75, requestRate: 100},
		{fairnessLabel: "workload2", tokens: 8, priority: 100, requestRate: 100},
		{fairnessLabel: "workload3", tokens: 8, priority: 125, requestRate: 100},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestSingleLowRequestLowTimeout(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 1, priority: 0, requestRate: 1, timeout: time.Millisecond},
		{fairnessLabel: "workload1", tokens: 8, priority: 75, requestRate: 100, timeout: 50 * time.Millisecond},
		{fairnessLabel: "workload2", tokens: 8, priority: 100, requestRate: 100, timeout: 50 * time.Millisecond},
		{fairnessLabel: "workload3", tokens: 8, priority: 125, requestRate: 100, timeout: 50 * time.Millisecond},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestIncreasingPriority(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 5, priority: 0, requestRate: 50},
		{fairnessLabel: "workload1", tokens: 5, priority: 50, requestRate: 50},
		{fairnessLabel: "workload2", tokens: 5, priority: 100, requestRate: 50},
		{fairnessLabel: "workload3", tokens: 5, priority: 150, requestRate: 50},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 3)
}

func Test0FillRate(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 0, priority: 0, requestRate: 50},
		{fairnessLabel: "workload1", tokens: 5, priority: 50, requestRate: 50},
		{fairnessLabel: "workload2", tokens: 5, priority: 100, requestRate: 50},
		{fairnessLabel: "workload3", tokens: 5, priority: 200, requestRate: 50},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, 0, 1)
}

func TestFairnessWithinPriority(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 16, priority: 50, requestRate: 50},
		{fairnessLabel: "workload1", tokens: 16, priority: 50, requestRate: 50},
		{fairnessLabel: "workload2", tokens: 16, priority: 100, requestRate: 50},
		{fairnessLabel: "workload3", tokens: 16, priority: 100, requestRate: 50},
	}
	for _, flow := range flows {
		flow.timeout = 50 * time.Millisecond
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)

	for i := 0; i < len(flows); i += 2 {
		tokensA := flows[i].acceptedRequests * flows[i].tokens
		tokensB := flows[i+1].acceptedRequests * flows[i+1].tokens
		// check fairness within priority levels
		if math.Abs(1-float64(tokensA)/float64(tokensB)) > _testTolerance {
			t.Logf("Fairness within priority levels is not within %f percent. Accepted tokens: %d, %d", _testTolerance*100, tokensA, tokensB)
			t.Fail()
		}
	}
}

func TestTimeouts(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 5, priority: 0, requestRate: 50, timeout: time.Millisecond * 5},
	}
	baseOfBasicBucketTest(t, flows, calculateFillRate(flows, 0.5), 1)
}

func TestLoadMultiplierBucket(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	lm := 0.7
	flows := flowTrackers{
		{
			fairnessLabel: "workload0",
			tokens:        5,
			priority:      0,
			requestRate:   500,
			timeout:       30 * time.Millisecond,
		},
	}
	schedMetrics := &WFQMetrics{
		FlowsGauge:        wfqFlowsGauge,
		HeapRequestsGauge: wfqHeapRequestsGauge,
	}

	c := clockwork.NewFakeClock()
	go updateClock(t, c, flows)

	loadMultiplierBucket := NewTokenBucketLoadMultiplier(c.Now(), _testSlotCount, _testSlotDuration, getMetrics())
	loadMultiplierBucket.SetContinuousTracking(true)
	sched := NewWFQScheduler(loadMultiplierBucket, c, schedMetrics)

	trainAndDeplete := func() {
		// Running Train and deplete the bucket
		depleteRunTime := time.Second * 2
		loadMultiplierBucket.SetLoadMultiplier(c.Now(), 0.0)

		runFlows(sched, &wg, flows, depleteRunTime, c)
		wg.Wait()
	}

	runExperiment := func() {
		// Running Actual Experiment
		loadMultiplierBucket.SetLoadMultiplier(c.Now(), lm)
		flowRunTime := time.Second * 10
		runFlows(sched, &wg, flows, flowRunTime, c)
		wg.Wait()

		totalAcceptedTokens := uint64(0)
		totalSentToken := uint64(0)
		for _, flow := range flows {
			totalAcceptedTokens += (flow.acceptedRequests * flow.tokens)
			totalSentToken += (flow.totalRequests * flow.tokens)
		}
		t.Logf("Total tokens sent: %d", totalSentToken)
		t.Logf("Total tokens accepted: %d", totalAcceptedTokens)
		totalSentToken = uint64(float64(totalSentToken) * lm)
		ratio := float64(totalAcceptedTokens) / float64(totalSentToken)

		if math.Abs(ratio-1) > _testTolerance {
			t.Errorf("Load Multiplier Bucket Test Failed, unfairness detected: %f", ratio)
		}
	}

	trainAndDeplete()
	runExperiment()

	// Experiment to send the traffic after some time so that the token bucket adds back to bootstrap
	c.Sleep(time.Second * 5)

	trainAndDeplete()
	runExperiment()
}

func TestPanic(t *testing.T) {
	// No need to check whether `recover()` is nil. Just turn off the panic.
	defer func() { _ = recover() }()

	c := clockwork.NewRealClock()
	startTime := c.Now()
	manager := NewTokenBucketLoadMultiplier(startTime, _testSlotCount, _testSlotDuration, getMetrics())
	manager.SetContinuousTracking(true)
	manager.SetLoadMultiplier(startTime, 0.5)
	if manager.LoadMultiplier() != 0.5 {
		t.Logf("LoadMultiplier is not 0.5\n")
	}
	manager.SetLoadMultiplier(startTime, -1.5)

	// If the panic is not thrown, the test will fail.
	t.Errorf("Expected panic has not been caught")
}
