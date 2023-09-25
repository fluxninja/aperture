package scheduler

import (
	"container/heap"
	"container/list"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/log"

	"github.com/fluxninja/aperture/v2/pkg/metrics"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
)

// Internal structure for tracking the request in the scheduler queue.
type queuedRequest struct {
	fInfo  *flowInfo
	ready  chan struct{} // Ready signal -- true = schedule, false = cancel/timeout
	flowID string        // Flow ID
	vft    float64       // Virtual finish time
	cost   float64       // Cost of the request (invPriority * tokens)
	onHeap bool          // Whether the request is on the heap or not
}

////////

// Memory pool for heapRequest(s).
var qRequestPool sync.Pool

func newHeapRequest() interface{} {
	qRequest := new(queuedRequest)
	qRequest.ready = make(chan struct{}, 1)
	return qRequest
}

func getHeapRequest() *queuedRequest {
	return qRequestPool.Get().(*queuedRequest)
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

// WFQScheduler : Weighted Fair Queue Scheduler.
type WFQScheduler struct {
	clk            clockwork.Clock
	lastAccessTime time.Time
	manager        TokenManager
	// metrics
	metrics *WFQMetrics
	// metrics labels
	metricsLabels prometheus.Labels
	// flows
	flows    map[string]*flowInfo
	requests requestHeap
	vt       float64 // virtual time
	// generation helps close the queue in face of concurrent requests leaving the queue while new requests also arrive.
	generation uint64
	lock       sync.Mutex
	queueOpen  bool // This tracks overload state
}

// NewWFQScheduler creates a new weighted fair queue scheduler.
func NewWFQScheduler(clk clockwork.Clock, tokenManger TokenManager, metrics *WFQMetrics, metricsLabels prometheus.Labels) Scheduler {
	sched := new(WFQScheduler)
	sched.queueOpen = false
	sched.generation = 0
	sched.clk = clk
	sched.lastAccessTime = sched.clk.Now()
	sched.vt = 0
	sched.flows = make(map[string]*flowInfo)
	sched.manager = tokenManger
	sched.metricsLabels = metricsLabels

	if metrics != nil {
		sched.metrics = metrics
	}

	return sched
}

func (sched *WFQScheduler) updateRequestInQueueMetrics(accepted bool, request *Request, startTime time.Time) {
	metricsLabels := make(prometheus.Labels, len(sched.metricsLabels)+1)
	metricsLabels[metrics.WorkloadIndexLabel] = strings.TrimPrefix(request.FairnessLabel, "workload:")
	for k, v := range sched.metricsLabels {
		metricsLabels[k] = v
	}
	requestInQueueMetricsObserver, err := sched.metrics.RequestInQueueDurationSummary.GetMetricWith(metricsLabels)
	if err == nil {
		requestInQueueMetricsObserver.Observe(float64(time.Since(startTime).Nanoseconds() / 1e6))
	} else {
		log.Error().Err(err).Msg("Failed to get request in queue duration summary")
	}
}

func (sched *WFQScheduler) updateMetricsAndReturn(accepted bool, remaining float64, current float64, request *Request, startTime time.Time) (bool, float64, float64) {
	if accepted {
		sched.metrics.AcceptedTokensCounter.Add(request.Tokens)
	}
	sched.metrics.IncomingTokensCounter.Add(request.Tokens)
	sched.updateRequestInQueueMetrics(accepted, request, startTime)
	return accepted, remaining, current
}

// Schedule blocks until the request is scheduled or until timeout.
// Return value - true: Accept, false: Reject.
func (sched *WFQScheduler) Schedule(ctx context.Context, request *Request) (bool, float64, float64) {
	startTime := time.Now()
	if request.Tokens == 0 {
		return sched.updateMetricsAndReturn(true, 0, 0, request, startTime)
	}

	sched.lock.Lock()
	queueOpen := sched.queueOpen
	sched.lastAccessTime = sched.clk.Now()
	sched.lock.Unlock()

	if sched.manager.PreprocessRequest(ctx, request) {
		return sched.updateMetricsAndReturn(true, 0, 0, request, startTime)
	}

	// try to schedule right now
	if !queueOpen {
		accepted, _, remaining, current := sched.manager.TakeIfAvailable(ctx, request.Tokens)
		if accepted {
			// we got the tokens, no need to queue
			return sched.updateMetricsAndReturn(true, remaining, current, request, startTime)
		}
	}

	// Unable to schedule right now, so queue the request
	qRequest := sched.queueRequest(ctx, request)

	// scheduler is in overload situation and we have to wait for ready signal and tokens
	select {
	case <-qRequest.ready:
		accepted, remaining, current := sched.scheduleRequest(ctx, request, qRequest)
		return sched.updateMetricsAndReturn(accepted, remaining, current, request, startTime)
	case <-ctx.Done():
		sched.cancelRequest(qRequest)
		return sched.updateMetricsAndReturn(false, 0, 0, request, startTime)
	}
}

// Construct FlowID by appending RequestLabel and Priority.
func (sched *WFQScheduler) flowID(fairnessLabel string, priority float64, generation uint64) string {
	return fmt.Sprintf("%s_%s_%d", fairnessLabel, strconv.FormatFloat(priority, 'f', -1, 64), generation)
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

	qRequest = getHeapRequest()

	flowID := sched.flowID(request.FairnessLabel, request.InvPriority, sched.generation)

	qRequest.flowID = flowID

	cost := float64(request.Tokens) * request.InvPriority

	// Get FlowInfo
	fInfo, ok := sched.flows[flowID]
	if !ok {
		fInfo = getFlowInfo()
		sched.setFlowsGauge(float64(len(sched.flows)))
		fInfo.vt = sched.vt
		sched.flows[flowID] = fInfo
	}
	// Increment reference counter
	fInfo.refCnt++

	// Store flowInfo pointer in the request
	qRequest.fInfo = fInfo

	// Store the cost of the request
	qRequest.cost = cost

	if !firstRequest {
		if !fInfo.requestOnHeap {
			qRequest.vft = fInfo.vt + cost
			qRequest.onHeap = true
			heap.Push(&sched.requests, qRequest)
			sched.setRequestsGauge(float64(sched.requests.Len()))
			fInfo.requestOnHeap = true
		} else {
			// push to flow queue
			fInfo.queue.PushBack(qRequest)
		}
	} else {
		// This is the only request in queue at this time, wake it up
		qRequest.ready <- struct{}{}
	}

	return qRequest
}

// adjust queue counters. Note: qRequest pointer should not be used after calling this function as it will get recycled via Pool.
func (sched *WFQScheduler) scheduleRequest(ctx context.Context, request *Request, qRequest *queuedRequest) (bool, float64, float64) {
	// This request has been selected to be executed next
	allowed, waitTime, remaining, current := sched.manager.Take(ctx, float64(request.Tokens))
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
					sched.manager.Return(returnCtx, request.Tokens)
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
					sched.manager.Return(returnCtx, request.Tokens)
				}(cancel)
			case <-timer.Chan():
			}
		}
	}

	sched.lock.Lock()
	defer sched.lock.Unlock()

	if allowed {
		// move the flow's VT forward
		qRequest.fInfo.vt += qRequest.cost
		// set new virtual time of scheduler
		sched.vt = qRequest.fInfo.vt
	}

	// This request is responsible for waking up the next request
	sched.wakeNextRequest(qRequest.fInfo)

	sched.cleanup(qRequest)

	return allowed, remaining, current
}

func (sched *WFQScheduler) wakeNextRequest(fInfo *flowInfo) {
	// load next request for this flow onto the heap
	if !fInfo.requestOnHeap {
		elm := fInfo.queue.Front()
		if elm != nil {
			fInfo.queue.Remove(elm)
			nextReq := elm.Value.(*queuedRequest)
			nextReq.vft = fInfo.vt + nextReq.cost
			heap.Push(&sched.requests, nextReq)
			nextReq.onHeap = true
			fInfo.requestOnHeap = true
			sched.setRequestsGauge(float64(sched.requests.Len()))
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
	sched.setRequestsGauge(float64(sched.requests.Len()))
	qRequest.onHeap = false
	qRequest.fInfo.requestOnHeap = false
	// wake up this request
	qRequest.ready <- struct{}{}
}

func (sched *WFQScheduler) cancelRequest(qRequest *queuedRequest) {
	sched.lock.Lock()
	defer sched.lock.Unlock()

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
		sched.setFlowsGauge(float64(len(sched.flows)))
	}
	putHeapRequest(qRequest)
}

// Revert returns tokens to the token bucket.
func (sched *WFQScheduler) Revert(ctx context.Context, tokens float64) {
	sched.manager.Return(ctx, tokens)
}

func (sched *WFQScheduler) setFlowsGauge(v float64) {
	if sched.metrics != nil && sched.metrics.FlowsGauge != nil {
		sched.metrics.FlowsGauge.Set(v)
	}
}

func (sched *WFQScheduler) setRequestsGauge(v float64) {
	if sched.metrics != nil && sched.metrics.HeapRequestsGauge != nil {
		sched.metrics.HeapRequestsGauge.Set(v)
	}
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

// WFQMetrics holds metrics related to internal workings of WFQScheduler.
type WFQMetrics struct {
	FlowsGauge                    prometheus.Gauge
	HeapRequestsGauge             prometheus.Gauge
	IncomingTokensCounter         prometheus.Counter
	AcceptedTokensCounter         prometheus.Counter
	RequestInQueueDurationSummary *prometheus.SummaryVec
}
