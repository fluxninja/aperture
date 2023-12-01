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

	policysyncv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
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
	wfqAcceptedTokensCounter        prometheus.Counter
	wfqRejectedTokensCounter        prometheus.Counter
	wfqIncomingTokensCounter        prometheus.Counter
	requestInQueueDurationSummary   *prometheus.SummaryVec
	tokenBucketLMGauge              prometheus.Gauge
	tokenBucketFillRateGauge        prometheus.Gauge
	tokenBucketBucketCapacityGauge  prometheus.Gauge
	tokenBucketAvailableTokensGauge prometheus.Gauge

	labels = prometheus.Labels{
		metrics.PolicyNameLabel:  "test-policy",
		metrics.PolicyHashLabel:  "test-hash",
		metrics.ComponentIDLabel: "test-component",
	}
)

func getMetrics() (prometheus.Gauge, *TokenBucketMetrics) {
	prometheusRegistry = prometheus.NewRegistry()

	constLabels := labels

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
	wfqIncomingTokensCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        metrics.IncomingTokensMetricName,
		Help:        "A counter measuring work incoming into Scheduler",
		ConstLabels: constLabels,
	})
	_ = prometheusRegistry.Register(wfqIncomingTokensCounter)
	wfqAcceptedTokensCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        metrics.AcceptedTokensMetricName,
		Help:        "A counter measuring work admitted by Scheduler",
		ConstLabels: constLabels,
	})
	_ = prometheusRegistry.Register(wfqAcceptedTokensCounter)
	wfqRejectedTokensCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        metrics.RejectedTokensMetricName,
		Help:        "A counter measuring work rejected by Scheduler",
		ConstLabels: constLabels,
	})
	_ = prometheusRegistry.Register(wfqRejectedTokensCounter)
	requestInQueueDurationSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.RequestInQueueDurationMetricName,
		Help: "A summary that tracks the duration of requests in queue",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIDLabel,
		metrics.WorkloadIndexLabel,
	},
	)
	_ = prometheusRegistry.Register(requestInQueueDurationSummary)
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

	return tokenBucketLMGauge, tbbMetrics
}

type flowTracker struct {
	fairnessLabel string  // what label it needs
	tokens        float64 // how many tokens it needs
	priority      float64
	timeout       time.Duration
	requestRate   uint64
	// counters
	totalRequests    uint64 // total requests made
	acceptedRequests uint64 // total requests accepted
}

func (flow *flowTracker) String() string {
	return fmt.Sprintf("FlowTracker<"+
		"Flow<FairnessLabel: %s "+
		"| RequestTokens: %0.2f "+
		"| Priority: %f "+
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
	ok, _, _ := sched.Schedule(ctx, flow.makeRequest())
	if ok {
		atomic.AddUint64(&flow.acceptedRequests, 1)
	}
}

func (flow *flowTracker) makeRequest() *Request {
	return NewRequest(flow.fairnessLabel, flow.tokens, flow.priority)
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
		{fairnessLabel: "workload1", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload2", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload3", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload4", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload5", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
	}
	c := clockwork.NewRealClock()
	_, metrics := getMetrics()

	manager := NewBasicTokenBucket(c, 0, metrics)

	schedMetrics := &WFQMetrics{
		FlowsGauge:                    wfqFlowsGauge,
		HeapRequestsGauge:             wfqHeapRequestsGauge,
		IncomingTokensCounter:         wfqIncomingTokensCounter,
		AcceptedTokensCounter:         wfqAcceptedTokensCounter,
		RejectedTokensCounter:         wfqRejectedTokensCounter,
		RequestInQueueDurationSummary: requestInQueueDurationSummary,
	}
	sched := NewWFQScheduler(c, manager, schedMetrics, labels)

	b.Logf("iterations: %d", b.N)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), flows[i%len(flows)].timeout)
			_, _, _ = sched.Schedule(ctx, flows[i%len(flows)].makeRequest())
			cancel()
		}
	})
}

func BenchmarkTokenBucketLoadMultiplier(b *testing.B) {
	flows := flowTrackers{
		{fairnessLabel: "workload1", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload2", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload3", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload4", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
		{fairnessLabel: "workload5", tokens: 1, priority: 1, timeout: 5 * time.Millisecond},
	}
	c := clockwork.NewRealClock()
	startTime := c.Now()

	lmGauge, metrics := getMetrics()
	manager := NewLoadMultiplierTokenBucket(c, _testSlotCount, _testSlotDuration, lmGauge, metrics)
	manager.SetContinuousTracking(true)
	manager.SetLoadDecisionValues(&policysyncv1.LoadDecision{
		LoadMultiplier: 1.0,
	})

	schedMetrics := &WFQMetrics{
		FlowsGauge:                    wfqFlowsGauge,
		HeapRequestsGauge:             wfqHeapRequestsGauge,
		IncomingTokensCounter:         wfqIncomingTokensCounter,
		AcceptedTokensCounter:         wfqAcceptedTokensCounter,
		RejectedTokensCounter:         wfqRejectedTokensCounter,
		RequestInQueueDurationSummary: requestInQueueDurationSummary,
	}
	sched := NewWFQScheduler(c, manager, schedMetrics, labels)

	// bootstrap bucket
	bootstrapTime := time.Second * 1
	for c.Now().Before(startTime.Add(bootstrapTime)) {
		for _, flow := range flows {
			ctx, cancel := context.WithTimeout(context.Background(), flow.timeout)
			_, _, _ = sched.Schedule(ctx, flow.makeRequest())
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
			_, _, _ = sched.Schedule(ctx, flows[i%len(flows)].makeRequest())
			cancel()
		}
	})
}

func calculateFillRate(flows flowTrackers, lm float64) float64 {
	fillRate := float64(0)
	for _, flow := range flows {
		fillRate += flow.tokens * float64(flow.requestRate)
	}
	return fillRate * lm
}

func baseOfBasicBucketTest(t *testing.T, flows flowTrackers, fillRate float64, noOfRuns int) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := clockwork.NewFakeClock()
	go updateClock(t, c, flows)

	startTime := c.Now()

	_, tbMetrics := getMetrics()
	basicBucket := NewBasicTokenBucket(c, fillRate, tbMetrics)
	metrics := &WFQMetrics{
		FlowsGauge:                    wfqFlowsGauge,
		HeapRequestsGauge:             wfqHeapRequestsGauge,
		IncomingTokensCounter:         wfqIncomingTokensCounter,
		AcceptedTokensCounter:         wfqAcceptedTokensCounter,
		RejectedTokensCounter:         wfqRejectedTokensCounter,
		RequestInQueueDurationSummary: requestInQueueDurationSummary,
	}
	sched := NewWFQScheduler(c, basicBucket, metrics, labels)
	var wg sync.WaitGroup
	var acceptedTokenRatio float64

	// ------------------------- Run Experiment -------------------------

	flowRunTime := time.Second * 10

	sumPriority := float64(0)
	adjustedPriority := make([]float64, len(flows))
	for i, flow := range flows {
		adjustedPriority[i] = 1 / flow.priority
		sumPriority += float64(adjustedPriority[i])
	}

	// Estimate the tokens -- It's a rough approach but seems to work so far for variety of loads and priorities
	totalEstimatedtokens := uint64(0)
	estimatedTokens := make([]uint64, len(flows))

	for i := range flows {
		estimatedTokens[i] = uint64(fillRate*float64(adjustedPriority[i])/sumPriority) * uint64(flowRunTime.Seconds())
		totalEstimatedtokens += estimatedTokens[i]
	}
	t.Logf("Flows: %v", flows)
	t.Logf("Adjusted Prios: %v\n", adjustedPriority)
	t.Logf("Estimated minimum tokens per flow before run are: %v\n", estimatedTokens)
	t.Logf("Total estimated allocated tokens: %d\n", totalEstimatedtokens)

	for i := 0; i < noOfRuns; i++ {
		runFlows(sched, &wg, flows, flowRunTime, c)
		wg.Wait()

		totalTokens := make([]float64, len(flows))
		acceptedTokenSum := float64(0)
		acceptedTokens := make([]float64, len(flows))
		for i, flow := range flows {
			acceptedTokens[i] = float64(flow.acceptedRequests) * flow.tokens
			acceptedTokenSum += acceptedTokens[i]
			totalTokens[i] = float64(flow.totalRequests) * flow.tokens
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
						t.Logf("Test failed because of more than %f percent unfairness %f: acceptedTokens: %0.2f, allocatedTokens: %d\n", _testTolerance*100, math.Abs(1-ratio), acceptedTokens[i], estimatedTokens[i])
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
	basicBucket.SetFillRate(0)
	if basicBucket.GetFillRate() != 0 {
		t.Logf("Fill rate is not 0 after stop\n")
		t.Fail()
	}
	t.Logf("Testing end\nDuration: %s", stopTime.Sub(startTime).String())
}

func TestHighRpsFlows(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 5, priority: 1, requestRate: 100},
		{fairnessLabel: "workload1", tokens: 5, priority: 1, requestRate: 100},
		{fairnessLabel: "workload2", tokens: 5, priority: 1, requestRate: 100},
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
		{fairnessLabel: "workload0", tokens: 1, priority: 1, requestRate: 1, timeout: time.Millisecond},
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
		{fairnessLabel: "workload0", tokens: 5, priority: 1, requestRate: 50},
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
		{fairnessLabel: "workload0", tokens: 0, priority: 1, requestRate: 50},
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
		tokensA := float64(flows[i].acceptedRequests) * flows[i].tokens
		tokensB := float64(flows[i+1].acceptedRequests) * flows[i+1].tokens
		// check fairness within priority levels
		if math.Abs(1-float64(tokensA)/float64(tokensB)) > _testTolerance {
			t.Logf("Fairness within priority levels is not within %f percent. Accepted tokens: %0.2f, %0.2f", _testTolerance*100, tokensA, tokensB)
			t.Fail()
		}
	}
}

func TestTimeouts(t *testing.T) {
	flows := flowTrackers{
		{fairnessLabel: "workload0", tokens: 5, priority: 1, requestRate: 50, timeout: time.Millisecond * 5},
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
		FlowsGauge:                    wfqFlowsGauge,
		HeapRequestsGauge:             wfqHeapRequestsGauge,
		IncomingTokensCounter:         wfqIncomingTokensCounter,
		AcceptedTokensCounter:         wfqAcceptedTokensCounter,
		RejectedTokensCounter:         wfqRejectedTokensCounter,
		RequestInQueueDurationSummary: requestInQueueDurationSummary,
	}

	c := clockwork.NewFakeClock()
	go updateClock(t, c, flows)

	lmGauge, tbMetrics := getMetrics()
	loadMultiplierBucket := NewLoadMultiplierTokenBucket(c, _testSlotCount, _testSlotDuration, lmGauge, tbMetrics)
	loadMultiplierBucket.SetContinuousTracking(true)
	sched := NewWFQScheduler(c, loadMultiplierBucket, schedMetrics, labels)

	trainAndDeplete := func() {
		// Running Train and deplete the bucket
		depleteRunTime := time.Second * 2
		loadMultiplierBucket.SetLoadDecisionValues(&policysyncv1.LoadDecision{
			LoadMultiplier: 0.0,
		})
		loadMultiplierBucket.SetPassThrough(false)

		runFlows(sched, &wg, flows, depleteRunTime, c)
		wg.Wait()
	}

	runExperiment := func() {
		// Running Actual Experiment
		loadMultiplierBucket.SetLoadDecisionValues(&policysyncv1.LoadDecision{
			LoadMultiplier: lm,
		})
		flowRunTime := time.Second * 10
		runFlows(sched, &wg, flows, flowRunTime, c)
		wg.Wait()

		totalAcceptedTokens := float64(0)
		totalSentToken := float64(0)
		for _, flow := range flows {
			totalAcceptedTokens += float64(flow.acceptedRequests) * flow.tokens
			totalSentToken += float64(flow.totalRequests) * flow.tokens
		}
		t.Logf("Total tokens sent: %0.2f", totalSentToken)
		t.Logf("Total tokens accepted: %0.2f", totalAcceptedTokens)
		totalSentToken = float64(totalSentToken) * lm
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
	lmGauge, tbMetrics := getMetrics()
	manager := NewLoadMultiplierTokenBucket(c, _testSlotCount, _testSlotDuration, lmGauge, tbMetrics)
	manager.SetContinuousTracking(true)
	manager.SetLoadDecisionValues(&policysyncv1.LoadDecision{
		LoadMultiplier: 0.5,
	})
	if manager.LoadMultiplier() != 0.5 {
		t.Logf("LoadMultiplier is not 0.5\n")
	}
	manager.SetLoadDecisionValues(&policysyncv1.LoadDecision{
		LoadMultiplier: -1.5,
	})

	// If the panic is not thrown, the test will fail.
	t.Errorf("Expected panic has not been caught")
}
