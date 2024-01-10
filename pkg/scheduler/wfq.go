package scheduler

import (
	"container/heap"
	"container/list"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/fluxninja/aperture/v2/pkg/log"

	"github.com/fluxninja/aperture/v2/pkg/metrics"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
)

// Internal structure for tracking the request in the scheduler queue.
type queuedRequest struct {
	fInfo                      *flowInfo
	request                    *Request
	workloadState              *WorkloadState
	ready                      chan struct{} // Ready signal -- true = schedule, false = cancel/timeout
	flowID                     string        // Flow ID
	vft                        float64       // Virtual finish time
	cost                       float64       // Cost of the request (invPriority * tokens)
	onHeap                     bool          // Whether the request is on the heap or not
	workloadPreemptionCounters preemptionCounters
	fairnessPreemptionCounters preemptionCounters
}

type preemptionCounters struct {
	tokensInQueue float64 // tokens in queue when request was added to queue
	tokensAllowed float64 // tokens allowed counter when request was added to queue
}

////////

// Memory pool for heapRequest(s).
var (
	qRequestPool      sync.Pool
	NumFairnessQueues = 1 << 8
)

func newHeapRequest() interface{} {
	qRequest := new(queuedRequest)
	qRequest.ready = make(chan struct{}, 1)
	return qRequest
}

func getHeapRequest(request *Request, workloadState *WorkloadState) *queuedRequest {
	qRequest := qRequestPool.Get().(*queuedRequest)
	qRequest.request = request
	qRequest.workloadState = workloadState
	return qRequest
}

func putHeapRequest(qRequest *queuedRequest) {
	qRequestPool.Put(qRequest)
}

////////

type requestHeap []*queuedRequest

// make sure we implement the heap interface.
var _ heap.Interface = &requestHeap{}

// Len returns the number of heap requests in the scheduler queue.
func (h *requestHeap) Len() int {
	return len(*h)
}

// Less compares heap requests by their virtual finish time.
// It's a min-heap -- i.e. requests with smallest vft are popped first.
// That's why we need to invert priority values - lower priority requests get larger vft values.
func (h *requestHeap) Less(i, j int) bool {
	return (*h)[i].vft < (*h)[j].vft
}

// Swap swaps two heap requests in the scheduler queue.
func (h *requestHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// Push appends a heap request to the scheduler queue.
func (h *requestHeap) Push(x interface{}) {
	request := x.(*queuedRequest)
	*h = append(*h, request)
}

// Pop removes the next heap request from the scheduler queue and returns the request.
func (h *requestHeap) Pop() interface{} {
	old := *h
	n := len(old)
	request := old[n-1]
	*h = old[0 : n-1]
	old[n-1] = nil
	return request
}

////////

type flowInfo struct {
	queue         *list.List
	workloadID    string
	vt            float64
	refCnt        int
	requestOnHeap bool
}

var fInfoPool sync.Pool

func newFlowInfo() interface{} {
	fInfo := new(flowInfo)
	fInfo.queue = list.New()
	return fInfo
}

func getFlowInfo() *flowInfo {
	return fInfoPool.Get().(*flowInfo)
}

func putFlowInfo(fInfo *flowInfo) {
	fInfoPool.Put(fInfo)
}

////////

func init() {
	qRequestPool.New = newHeapRequest
	fInfoPool.New = newFlowInfo
}

////////

// WorkloadState holds the state of a workload.
type WorkloadState struct {
	preemptionMetrics *PreemptionMetrics
	numFairnessQueues int
}

// WFQScheduler : Weighted Fair Queue Scheduler.
type WFQScheduler struct {
	clk                    clockwork.Clock
	lastAccessTime         time.Time
	fairnessSeedUpdateTime time.Time
	manager                TokenManager
	// metrics
	metrics *WFQMetrics
	// preemption metrics
	preemptionMetrics *PreemptionMetrics
	// metrics labels
	metricsLabels prometheus.Labels
	// flows
	flows     map[string]*flowInfo
	workloads map[string]*WorkloadState // map from workloadID to WorkloadState
	requests  requestHeap
	vt        float64 // virtual time
	// generation helps close the queue in face of concurrent requests leaving the queue while new requests also arrive.
	generation   uint64
	fairnessSeed uint32
	lock         sync.Mutex
	queueOpen    bool // This tracks overload state
}

// NewWFQScheduler creates a new weighted fair queue scheduler.
func NewWFQScheduler(clk clockwork.Clock, tokenManger TokenManager, metrics *WFQMetrics, metricsLabels prometheus.Labels) Scheduler {
	sched := new(WFQScheduler)
	sched.queueOpen = false
	sched.generation = 0
	sched.clk = clk
	now := sched.clk.Now()
	sched.lastAccessTime = now
	sched.fairnessSeedUpdateTime = now
	sched.fairnessSeed = 0
	sched.vt = 0
	sched.flows = make(map[string]*flowInfo)
	sched.workloads = make(map[string]*WorkloadState)
	sched.manager = tokenManger
	sched.metricsLabels = metricsLabels

	if metrics != nil {
		sched.metrics = metrics
		sched.preemptionMetrics = NewPreemptionMetrics(metricsLabels, metrics.WorkloadPreemptedTokensSummary, metrics.WorkloadDelayedTokensSummary, metrics.WorkloadOnTimeCounter)
	} else {
		sched.preemptionMetrics = NewPreemptionMetrics(metricsLabels, nil, nil, nil)
	}

	return sched
}

func (sched *WFQScheduler) updateRequestInQueueMetrics(accepted bool, request *Request, startTime time.Time) {
	metricsLabels := appendWorkloadLabel(sched.metricsLabels, request.WorkloadLabel)
	requestInQueueMetricsObserver, err := sched.metrics.RequestInQueueDurationSummary.GetMetricWith(metricsLabels)
	if err == nil {
		requestInQueueMetricsObserver.Observe(float64(time.Since(startTime).Nanoseconds() / 1e6))
	} else {
		log.Error().Err(err).Msg("Failed to get request in queue duration summary")
	}
}

func (sched *WFQScheduler) updateOutgoingTokenMetrics(accepted bool, tokens float64) {
	if tokens <= 0 {
		return
	}
	if accepted {
		sched.metrics.AcceptedTokensCounter.Add(tokens)
	} else {
		sched.metrics.RejectedTokensCounter.Add(tokens)
	}
}

func (sched *WFQScheduler) updateMetricsAndReturn(accepted bool, remaining float64, current float64, request *Request, startTime time.Time, reqID string) (bool, float64, float64, string) {
	sched.updateOutgoingTokenMetrics(accepted, request.Tokens)
	sched.updateRequestInQueueMetrics(accepted, request, startTime)
	return accepted, remaining, current, reqID
}

func (sched *WFQScheduler) updateIncomingTokensMetric(tokens float64) {
	if tokens > 0 {
		sched.metrics.IncomingTokensCounter.Add(tokens)
	}
}

// Schedule blocks until the request is scheduled or until timeout.
// Return value - true: Accept, false: Reject.
func (sched *WFQScheduler) Schedule(ctx context.Context, request *Request) (bool, float64, float64, string) {
	startTime := time.Now()
	sched.updateIncomingTokensMetric(request.Tokens)
	if request.Tokens == 0 {
		return sched.updateMetricsAndReturn(true, 0, 0, request, startTime, "")
	}

	sched.lock.Lock()
	queueOpen := sched.queueOpen
	sched.lastAccessTime = sched.clk.Now()
	sched.lock.Unlock()

	if sched.manager.PreprocessRequest(ctx, request) {
		return sched.updateMetricsAndReturn(true, 0, 0, request, startTime, "")
	}

	// try to schedule right now
	if !queueOpen {
		accepted, _, remaining, current, reqID := sched.manager.TakeIfAvailable(ctx, request.Tokens)
		if accepted {
			// we got the tokens, no need to queue
			return sched.updateMetricsAndReturn(true, remaining, current, request, startTime, reqID)
		}
	}

	// Unable to schedule right now, so queue the request
	// Update count for tokens in the queue here
	qRequest := sched.queueRequest(ctx, request)

	// scheduler is in overload situation and we have to wait for ready signal and tokens
	select {
	case <-qRequest.ready:
		accepted, remaining, current, reqID := sched.scheduleRequest(ctx, request, qRequest)
		// Update count for tokens in the queue here
		return sched.updateMetricsAndReturn(accepted, remaining, current, request, startTime, reqID)
	case <-ctx.Done():
		sched.cancelRequest(qRequest)
		// Update count for tokens in the queue here
		return sched.updateMetricsAndReturn(false, 0, 0, request, startTime, "")
	}
}

// Identifiers computes fairnessQueueID by hashing fairnessLabel and doing a bit-wise AND with number of fairness queues. Constructs workloadID by appending workloadLabel, Priority and Generation. Constructs flowID by appending workloadID and fairnessQueueID.
func (sched *WFQScheduler) Identifiers(workloadLabel, fairnessLabel string, priority float64, generation uint64) (string, string) {
	// workloadID is the workloadLabel + priority + generation
	workloadID := fmt.Sprintf("%s-%f-%d", workloadLabel, priority, generation)

	var flowID string
	if fairnessLabel != "" {
		if sched.lastAccessTime.Sub(sched.fairnessSeedUpdateTime) > 300*time.Second {
			// randomize fairnessSeed every 5 minutes
			//nolint:gosec // G404: Use of math/rand is acceptable here due to non-security context
			sched.fairnessSeed = rand.Uint32()
			sched.fairnessSeedUpdateTime = sched.lastAccessTime
		}
		// Compute the hash using xxHash
		hash := xxhash.Sum64String(fairnessLabel + fmt.Sprintf("%d", sched.fairnessSeed))

		// Assuming NumFairnessQueues is a power of 2, use bitwise AND for modulo
		fairnessQueueID := hash & uint64(NumFairnessQueues-1)
		// flowID is the workloadID + fairnessQueueID
		flowID = fmt.Sprintf("%s-%d", workloadID, fairnessQueueID)
	} else {
		// flowID is the workloadID
		flowID = workloadID
	}

	return flowID, workloadID
}

// Attempt to queue this request.
//
// Returns whether request was admitted right away without queueing.
// If admitted == false, might return a valid heapRequest
// If admitted == false and qRequest == nil, request was neither admitted nor
// queued (rejected right away).
func (sched *WFQScheduler) queueRequest(ctx context.Context, request *Request) (qRequest *queuedRequest) {
	sched.lock.Lock()
	defer sched.lock.Unlock()

	firstRequest := false

	// check if this is the first request entering this queue
	if !sched.queueOpen {
		firstRequest = true
		sched.queueOpen = true
		// reset sched virtual time
		sched.vt = 0
	}

	// Proceed to queueing

	flowID, workloadID := sched.Identifiers(request.WorkloadLabel, request.FairnessLabel, request.InvPriority, sched.generation)

	workloadState, workloadStateFound := sched.workloads[workloadID]
	if !workloadStateFound {
		workloadState = &WorkloadState{
			preemptionMetrics: NewPreemptionMetrics(sched.metricsLabels, sched.metrics.FairnessPreemptedTokensSummary, sched.metrics.FairnessDelayedTokensSummary, sched.metrics.FairnessOnTimeCounter),
			numFairnessQueues: 0,
		}
		sched.workloads[workloadID] = workloadState
	}

	// Get FlowInfo
	fInfo, ok := sched.flows[flowID]
	if !ok {
		fInfo = getFlowInfo()
		fInfo.workloadID = workloadID
		fInfo.vt = sched.vt
		sched.flows[flowID] = fInfo
		workloadState.numFairnessQueues++
	}

	qRequest = getHeapRequest(request, workloadState)

	qRequest.flowID = flowID
	// Store flowInfo pointer in the request
	qRequest.fInfo = fInfo
	// Increment reference counter
	fInfo.refCnt++

	sched.preemptionMetrics.onQueueEntry(&qRequest.workloadPreemptionCounters, request)

	fairnessAdjustment := workloadState.numFairnessQueues

	cost := float64(request.Tokens) * request.InvPriority * float64(fairnessAdjustment)

	// Store the cost of the request
	qRequest.cost = cost

	if !firstRequest {
		if !fInfo.requestOnHeap {
			qRequest.vft = fInfo.vt + cost
			qRequest.onHeap = true
			heap.Push(&sched.requests, qRequest)
			fInfo.requestOnHeap = true
		} else {
			// push to flow queue
			fInfo.queue.PushBack(qRequest)
			workloadState.preemptionMetrics.onQueueEntry(&qRequest.fairnessPreemptionCounters, request)
		}
	} else {
		// This is the only request in queue at this time, wake it up
		qRequest.ready <- struct{}{}
	}

	return qRequest
}

// adjust queue counters. Note: qRequest pointer should not be used after calling this function as it will get recycled via Pool.
func (sched *WFQScheduler) scheduleRequest(ctx context.Context, request *Request, qRequest *queuedRequest) (bool, float64, float64, string) {
	// This request has been selected to be executed next
	allowed, waitTime, remaining, current, reqID := sched.manager.Take(ctx, float64(request.Tokens))
	// check if we need to wait
	if allowed && waitTime > 0 {
		// check whether ctx has deadline
		// and if deadline is less than waitTime
		// return tokens immediately
		if dl, o := ctx.Deadline(); o {
			if dl.Sub(sched.clk.Now()) < waitTime {
				allowed = false
				returnCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				go func(cancel context.CancelFunc) {
					defer cancel()
					sched.manager.Return(returnCtx, request.Tokens, reqID)
				}(cancel)
			}
		}

		if allowed {
			timer := sched.clk.NewTimer(waitTime)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				allowed = false
				returnCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				// return the tokens
				go func(cancel context.CancelFunc) {
					defer cancel()
					sched.manager.Return(returnCtx, request.Tokens, reqID)
				}(cancel)
			case <-timer.Chan():
			}
		}
	}

	sched.lock.Lock()
	defer sched.lock.Unlock()

	// Update metrics for preemption and delay
	sched.preemptionMetrics.onQueueExit(&qRequest.workloadPreemptionCounters, request, allowed)

	if allowed {
		// move the flow's VT forward
		qRequest.fInfo.vt += qRequest.cost
		// set new virtual time of scheduler
		sched.vt = qRequest.fInfo.vt
	}

	// This request is responsible for waking up the next request
	sched.wakeNextRequest(qRequest.fInfo)

	sched.cleanup(qRequest)

	return allowed, remaining, current, reqID
}

func (sched *WFQScheduler) wakeNextRequest(fInfo *flowInfo) {
	// load next request for this flow onto the heap
	if !fInfo.requestOnHeap {
		elm := fInfo.queue.Front()
		if elm != nil {
			fInfo.queue.Remove(elm)
			nextReq := elm.Value.(*queuedRequest)
			nextReq.workloadState.preemptionMetrics.onQueueExit(&nextReq.fairnessPreemptionCounters, nextReq.request, true)
			nextReq.vft = fInfo.vt + nextReq.cost
			heap.Push(&sched.requests, nextReq)
			nextReq.onHeap = true
			fInfo.requestOnHeap = true
		}
	}

	// no more requests?
	if sched.requests.Len() == 0 {
		// close the queue
		sched.generation++
		sched.queueOpen = false
		return
	}
	// Pop from queue and wake next request
	qRequest := heap.Pop(&sched.requests).(*queuedRequest)
	qRequest.onHeap = false
	qRequest.fInfo.requestOnHeap = false
	// wake up this request
	qRequest.ready <- struct{}{}
}

func (sched *WFQScheduler) cancelRequest(qRequest *queuedRequest) {
	sched.lock.Lock()
	defer sched.lock.Unlock()

	sched.preemptionMetrics.onQueueExit(&qRequest.workloadPreemptionCounters, qRequest.request, false)

	select {
	case <-qRequest.ready:
		// This request is responsible for waking up the next request
		sched.wakeNextRequest(qRequest.fInfo)
	default:
		// remove from heap
		if qRequest.onHeap {
			qRequest.onHeap = false
			for i := 0; i < sched.requests.Len(); i++ {
				if sched.requests[i] == qRequest {
					// replace with the next request in the flow
					elm := qRequest.fInfo.queue.Front()
					if elm != nil {
						qRequest.fInfo.queue.Remove(elm)
						nextReq := elm.Value.(*queuedRequest)
						nextReq.workloadState.preemptionMetrics.onQueueExit(&nextReq.fairnessPreemptionCounters, nextReq.request, true)
						nextReq.vft = qRequest.fInfo.vt + nextReq.cost
						sched.requests[i] = nextReq
						nextReq.onHeap = true
						qRequest.fInfo.requestOnHeap = true
					} else {
						// swap with the last element in the heap
						sched.requests.Swap(i, sched.requests.Len()-1)
						// trim the last element from the slice
						sched.requests = sched.requests[:sched.requests.Len()-1]
						qRequest.fInfo.requestOnHeap = false
					}
					// Fix the heap
					if i < sched.requests.Len() {
						heap.Fix(&sched.requests, i)
					}
					break
				}
			}
		} else {
			// search within the flow queue
			var next *list.Element
			for elm := qRequest.fInfo.queue.Front(); elm != nil; elm = next {
				request := elm.Value.(*queuedRequest)
				next = elm.Next()
				if request == qRequest {
					qRequest.fInfo.queue.Remove(elm)
					qRequest.workloadState.preemptionMetrics.onQueueExit(&qRequest.fairnessPreemptionCounters, qRequest.request, false)
					break
				}
			}
		}
	}

	sched.cleanup(qRequest)
}

// queueRequest is going to be recycled and must not be used
// after calling this function.
func (sched *WFQScheduler) cleanup(qRequest *queuedRequest) {
	// decrement reference counter
	qRequest.fInfo.refCnt--
	// check if the flow is empty
	if qRequest.fInfo.refCnt == 0 {
		// delete the flow
		delete(sched.flows, qRequest.flowID)
		// send flowInfo back to the Pool
		putFlowInfo(qRequest.fInfo)
		workloadState, ok := sched.workloads[qRequest.fInfo.workloadID]
		if ok {
			workloadState.numFairnessQueues--
			if workloadState.numFairnessQueues == 0 {
				delete(sched.workloads, qRequest.fInfo.workloadID)
			}
		}
	}
	putHeapRequest(qRequest)
}

// Info returns the last access time and number of requests that are currently in the queue.
func (sched *WFQScheduler) Info() (time.Time, int) {
	sched.lock.Lock()
	defer sched.lock.Unlock()
	return sched.lastAccessTime, sched.requests.Len()
}

// GetPendingFlows returns the number of flows in the scheduler.
func (sched *WFQScheduler) GetPendingFlows() int {
	return len(sched.flows)
}

// GetPendingRequests returns the number of requests in the scheduler.
func (sched *WFQScheduler) GetPendingRequests() int {
	return len(sched.requests)
}

// PreemptionMetrics holds metrics related to preemption and delay for a queuing system.
type PreemptionMetrics struct {
	workloadPreemptedTokensSummary *prometheus.SummaryVec
	workloadDelayedTokensSummary   *prometheus.SummaryVec
	workloadOnTimeCounter          *prometheus.CounterVec
	metricsLabels                  prometheus.Labels
	tokensInQueue                  float64
	tokensAllowed                  float64
}

// NewPreemptionMetrics creates a new PreemptionMetrics object.
func NewPreemptionMetrics(
	metricsLabels prometheus.Labels,
	workloadPreemptedTokensSummary *prometheus.SummaryVec,
	workloadDelayedTokensSummary *prometheus.SummaryVec,
	workloadOnTimeCounter *prometheus.CounterVec,
) *PreemptionMetrics {
	return &PreemptionMetrics{
		metricsLabels:                  metricsLabels,
		workloadPreemptedTokensSummary: workloadPreemptedTokensSummary,
		workloadDelayedTokensSummary:   workloadDelayedTokensSummary,
		workloadOnTimeCounter:          workloadOnTimeCounter,
	}
}

// Maintain token counters used for calculating preemption and delay metrics.
// WARNING: Unsafe and should be called with scheduler lock.
func (pMetrics *PreemptionMetrics) onQueueEntry(pCounter *preemptionCounters, request *Request) {
	pCounter.tokensInQueue = pMetrics.tokensInQueue
	pCounter.tokensAllowed = pMetrics.tokensAllowed
	pMetrics.tokensInQueue += request.Tokens
}

// Update metrics for preemption and delay
// WARNING: Unsafe and should be called with scheduler lock.
func (pMetrics *PreemptionMetrics) onQueueExit(pCounter *preemptionCounters, request *Request, allowed bool) {
	initMetrics := func(labels prometheus.Labels) error {
		var err error
		_, err = pMetrics.workloadPreemptedTokensSummary.GetMetricWith(labels)
		if err != nil {
			return fmt.Errorf("%w: failed to get workload_preempted_tokens summary", err)
		}
		_, err = pMetrics.workloadDelayedTokensSummary.GetMetricWith(labels)
		if err != nil {
			return fmt.Errorf("%w: failed to get workload_delayed_tokens summary", err)
		}
		_, err = pMetrics.workloadOnTimeCounter.GetMetricWith(labels)
		if err != nil {
			return fmt.Errorf("%w: failed to get workload_on_time_total counter", err)
		}
		return nil
	}

	publishSummary := func(summary *prometheus.SummaryVec, value float64) {
		if summary == nil {
			return
		}
		metricsLabels := appendWorkloadLabel(pMetrics.metricsLabels, request.WorkloadLabel)
		err := initMetrics(metricsLabels)
		if err != nil {
			log.Error().Err(err).Msg("Failed to initialize metrics")
		}
		observer, err := summary.GetMetricWith(metricsLabels)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get workload preempted tokens summary")
			return
		}
		observer.Observe(value)
	}

	publishCounter := func(counterVec *prometheus.CounterVec, value float64) {
		if counterVec == nil {
			return
		}
		metricsLabels := appendWorkloadLabel(pMetrics.metricsLabels, request.WorkloadLabel)
		err := initMetrics(metricsLabels)
		if err != nil {
			log.Error().Err(err).Msg("Failed to initialize metrics")
		}
		counter, err := counterVec.GetMetricWith(metricsLabels)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get workload on time counter")
		}
		counter.Add(value)
	}

	if allowed {
		// Update metrics
		realIncrement := pMetrics.tokensAllowed - pCounter.tokensAllowed
		expectedIncrement := pCounter.tokensInQueue
		bump := expectedIncrement - realIncrement
		if bump == 0 && pMetrics.workloadOnTimeCounter != nil {
			publishCounter(pMetrics.workloadOnTimeCounter, 1)
		} else if bump > 0 && pMetrics.workloadPreemptedTokensSummary != nil {
			publishSummary(pMetrics.workloadPreemptedTokensSummary, bump)
		} else if bump < 0 && pMetrics.workloadDelayedTokensSummary != nil {
			publishSummary(pMetrics.workloadDelayedTokensSummary, -bump)
		}

		// Update count for tokens accepted
		pMetrics.tokensAllowed += request.Tokens
	}

	// Update tokens in the queue for calculating preemption and delay metrics
	pMetrics.tokensInQueue -= request.Tokens
}

func appendWorkloadLabel(baseMetricsLabels prometheus.Labels, workloadLabel string) prometheus.Labels {
	metricsLabels := make(prometheus.Labels, len(baseMetricsLabels)+1)
	metricsLabels[metrics.WorkloadIndexLabel] = workloadLabel
	for k, v := range baseMetricsLabels {
		metricsLabels[k] = v
	}

	return metricsLabels
}

// WFQMetrics holds metrics related to internal workings of WFQScheduler.
type WFQMetrics struct {
	IncomingTokensCounter          prometheus.Counter
	AcceptedTokensCounter          prometheus.Counter
	RejectedTokensCounter          prometheus.Counter
	RequestInQueueDurationSummary  *prometheus.SummaryVec
	WorkloadPreemptedTokensSummary *prometheus.SummaryVec
	WorkloadDelayedTokensSummary   *prometheus.SummaryVec
	WorkloadOnTimeCounter          *prometheus.CounterVec
	FairnessPreemptedTokensSummary *prometheus.SummaryVec
	FairnessDelayedTokensSummary   *prometheus.SummaryVec
	FairnessOnTimeCounter          *prometheus.CounterVec
}
