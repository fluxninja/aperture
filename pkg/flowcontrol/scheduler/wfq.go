package scheduler

import (
	"container/heap"
	"container/list"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
)

// Internal structure for tracking the request in the scheduler queue.
type queuedRequest struct {
	enqueueTime time.Time // Time when this request was enqueued
	fInfo       *flowInfo
	ready       chan bool     // Ready signal -- true = schedule, false = cancel/timeout
	flowID      string        // Flow ID
	timeout     time.Duration // Timeout for this request
	vft         uint64        // Virtual finish time
	cost        uint64        // Cost of the request (invPriority * tokens)
}

////////

// Memory pool for heapRequest(s).
var qRequestPool sync.Pool

func newHeapRequest() interface{} {
	qRequest := new(queuedRequest)
	qRequest.ready = make(chan bool, 1)
	return qRequest
}

func getHeapRequest() *queuedRequest {
	return qRequestPool.Get().(*queuedRequest)
}

func putHeapRequest(qRequest *queuedRequest) {
	qRequest.fInfo = nil
	qRequestPool.Put(qRequest)
}

////////

type requestHeap []*queuedRequest

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
	vt            uint64
	refCnt        int
	requestOnHeap bool
	auditTime     time.Time
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

// WFQScheduler : WFQ + CoDel.
type WFQScheduler struct {
	auditTime time.Time // Time when audit last ran
	lock      sync.Mutex
	manager   TokenManager
	clk       clockwork.Clock
	// metrics
	metrics *WFQMetrics
	// flows
	flows         map[string]*flowInfo
	requests      requestHeap
	minTimeout    time.Duration
	auditDuration time.Duration
	vt            uint64 // virtual time
	// generation helps close the queue in face of concurrent requests leaving the queue while new requests also arrive.
	generation     uint64
	defaultTimeout time.Duration
	queueOpen      bool // This tracks overload state
}

// WFQMetrics holds metrics related to internal workings of WFQScheduler.
type WFQMetrics struct {
	FlowsGauge        prometheus.Gauge
	HeapRequestsGauge prometheus.Gauge
}

// GetPendingFlows returns the number of flows in the scheduler.
func (sched *WFQScheduler) GetPendingFlows() int {
	return len(sched.flows)
}

// GetPendingRequests returns the number of requests in the scheduler.
func (sched *WFQScheduler) GetPendingRequests() int {
	return len(sched.requests)
}

// NewWFQScheduler creates a new weighted fair queue scheduler,
// timeout -- timeout for requests.
func NewWFQScheduler(timeout time.Duration, tokenManger TokenManager, clk clockwork.Clock, metrics *WFQMetrics) Scheduler {
	sched := new(WFQScheduler)
	sched.queueOpen = false
	sched.generation = 0
	sched.vt = 0
	sched.flows = make(map[string]*flowInfo)
	sched.minTimeout = timeout
	sched.auditDuration = sched.minTimeout / 2
	sched.defaultTimeout = timeout
	sched.manager = tokenManger
	sched.clk = clk

	now := sched.clk.Now()
	sched.auditTime = now

	if metrics != nil {
		sched.metrics = metrics
	}

	return sched
}

// Schedule blocks until the request is scheduled or until timeout (CoDel).
// Return value - true: Accept, false: Reject.
func (sched *WFQScheduler) Schedule(rContext RequestContext) bool {
	if rContext.Tokens == 0 {
		return true
	}

	// Assign default timeout to this request if it is not set
	if rContext.Timeout == 0 {
		rContext.Timeout = sched.defaultTimeout
	}

	// Unable to schedule right now, so queue the request
	qRequest := sched.enter(rContext)

	if qRequest == nil {
		// scheduler is not in overload situation and the request was able to get tokens
		return true
	}

	// scheduler is in overload situation and we have to wait for ready signal and tokens
	ready := <-qRequest.ready
	// leave the scheduler
	return sched.leave(rContext, qRequest, ready)
}

// Construct FlowID by appending RequestLabel and Priority.
func (sched *WFQScheduler) flowID(fairnessLabel string, priority uint8, generation uint64) string {
	return fmt.Sprintf("%s_%d_%d", fairnessLabel, priority, generation)
}

func (sched *WFQScheduler) auditFlow(now time.Time, fInfo *flowInfo) {
	// Reset auditTime

	fInfo.auditTime = now

	if sched.queueOpen {
		var next *list.Element
		// range through flow queue and timeout the requests
		for e := fInfo.queue.Front(); e != nil; e = next {
			request := e.Value.(*queuedRequest)
			next = e.Next()
			if now.Sub(request.enqueueTime) > request.timeout {
				// remove this request from the queue
				fInfo.queue.Remove(e)
				request.ready <- false
			}
		}
	}
}

// Run an audit of the queue every timeout duration
// parent function must hold the lock.
func (sched *WFQScheduler) auditHeap(now time.Time) {
	// Reset auditTime

	sched.auditTime = now

	if sched.queueOpen {
		heapify := false
		// wake up and cancel stale requests
		for i := 0; i < len(sched.requests); i++ {
			request := sched.requests[i]
			if now.Sub(request.enqueueTime) > request.timeout {
				heapify = true
				// request is stale, cancel it
				request.ready <- false
				// replace with the next request in the flow
				elm := request.fInfo.queue.Front()
				if elm != nil {
					request.fInfo.queue.Remove(elm)
					nextReq := elm.Value.(*queuedRequest)
					nextReq.vft = request.fInfo.vt + nextReq.cost
					sched.requests[i] = nextReq
					request.fInfo.requestOnHeap = true
				} else {
					// swap with the last element in the heap
					sched.requests.Swap(i, len(sched.requests)-1)
					// trim last element from the slice
					sched.requests = sched.requests[:len(sched.requests)-1]
					request.fInfo.requestOnHeap = false
				}
				i--
			}
		}
		if heapify {
			// reinitialize the heap
			heap.Init(&sched.requests)
		}
	}
	sched.setFlowsGauge(float64(len(sched.flows)))
	sched.setRequestsGauge(float64(sched.requests.Len()))
}

// Attempt to queue this request.
// If queued successfully, return a valid heapRequest.
// Otherwise, request was handled right away without queuing - return nil.
func (sched *WFQScheduler) enter(rContext RequestContext) (qRequest *queuedRequest) {
	sched.lock.Lock()
	defer sched.lock.Unlock()

	now := sched.clk.Now()

	// update timeout value to be global minimum
	if rContext.Timeout != 0 && rContext.Timeout < sched.minTimeout {
		sched.minTimeout = rContext.Timeout
		sched.auditDuration = sched.minTimeout / 2
	}

	if sched.manager.PreprocessRequest(now, rContext) {
		return nil
	}

	// try to schedule right now
	if !sched.queueOpen {
		ok := sched.manager.TakeIfAvailable(now, float64(rContext.Tokens))
		if ok {
			// we got the tokens, no need to queue
			return nil
		}
	}

	// we are in overload situation, attempt to queue

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

	flowID := sched.flowID(rContext.FairnessLabel, rContext.Priority, sched.generation)

	qRequest.flowID = flowID

	// invPriority range [1, 256]
	invPriority := uint64(math.MaxUint8-rContext.Priority) + 1
	cost := rContext.Tokens * invPriority

	// Get FlowInfo
	fInfo, ok := sched.flows[flowID]
	if !ok {
		fInfo = getFlowInfo()
		fInfo.vt = sched.vt
		fInfo.auditTime = now
		sched.flows[flowID] = fInfo
	}
	// Increment reference counter
	fInfo.refCnt++

	// Store flowInfo pointer in the request
	qRequest.fInfo = fInfo

	// Store current timestamp
	qRequest.enqueueTime = now
	// Store the timeout duration
	qRequest.timeout = rContext.Timeout

	// Store the cost of the request
	qRequest.cost = cost

	// Run audits if needed (audit runs at half the minTimeout, i.e. nyquist frequency)
	if now.Sub(fInfo.auditTime) > sched.auditDuration {
		sched.auditFlow(now, fInfo)
	}

	// Run heap audit at half the frequency of flow audits
	if now.Sub(sched.auditTime) > sched.minTimeout {
		sched.auditHeap(now)
	}

	if !firstRequest {
		if !fInfo.requestOnHeap {
			qRequest.vft = fInfo.vt + cost
			heap.Push(&sched.requests, qRequest)
			fInfo.requestOnHeap = true
		} else {
			// push to flow queue
			fInfo.queue.PushBack(qRequest)
		}
	} else {
		// This is the only request in queue at this time, wake it up
		qRequest.ready <- true
	}
	return qRequest
}

// adjust queue counters. Note: qRequest pointer should not be used after calling this function as it will get recycled via Pool.
func (sched *WFQScheduler) leave(rContext RequestContext, qRequest *queuedRequest, ready bool) bool {
	defer func() {
		if ready {
			// This request is now responsible for waking up the next request (after it gets the tokens)
			sched.wakeNextRequest()
		}
	}()

	// incase we need to sleep
	waitTime := time.Duration(0)

	// wait for tokens if needed
	defer func() {
		// we got the tokens -- check if we need to wait
		if waitTime > 0 {
			sched.clk.Sleep(waitTime)
		}
	}()

	sched.lock.Lock()
	// Unlock
	defer sched.lock.Unlock()

	defer func() {
		// decrement reference counter
		qRequest.fInfo.refCnt--
		// check if the flow is empty
		if qRequest.fInfo.refCnt == 0 {
			// delete the flow
			delete(sched.flows, qRequest.flowID)
			// send flowInfo back to the Pool
			putFlowInfo(qRequest.fInfo)
		} else { // check whether new requests arrived in between that need loading
			sched.loadNextFlowReq(qRequest.fInfo)
		}
		putHeapRequest(qRequest)
		qRequest = nil
	}()

	if !ready {
		// Request has been canceled/timedout
		return false
	}

	// This request has been selected to be executed next
	now := sched.clk.Now()

	remainingTimeout := rContext.Timeout - now.Sub(qRequest.enqueueTime)
	waitTime, ok := sched.manager.Take(now, remainingTimeout, float64(rContext.Tokens))

	if ok {
		// move the flow's VT forward and add the next request from flow queue to the heap
		qRequest.fInfo.vt += qRequest.cost
		// set new virtual time of scheduler
		sched.vt = qRequest.fInfo.vt
	}

	// deferred stack -
	// 1. remove flow and put back heapRequest
	// 2. release lock
	// 3. sleep i.e. wait for tokens
	// 4. wake the next task
	return ok
}

func (sched *WFQScheduler) wakeNextRequest() {
	sched.lock.Lock()
	defer sched.lock.Unlock()

	now := sched.clk.Now()
	for {
		if sched.requests.Len() == 0 {
			// close the queue
			sched.generation++
			sched.queueOpen = false
			return
		}
		// Pop from queue and wake a valid request
		qRequest := heap.Pop(&sched.requests).(*queuedRequest)
		qRequest.fInfo.requestOnHeap = false
		if now.Sub(qRequest.enqueueTime) < qRequest.timeout {
			// wake up this request
			qRequest.ready <- true
			return
		}
		sched.loadNextFlowReq(qRequest.fInfo)
		// Timeout this request
		qRequest.ready <- false
	}
}

func (sched *WFQScheduler) loadNextFlowReq(fInfo *flowInfo) {
	if !fInfo.requestOnHeap {
		elm := fInfo.queue.Front()
		if elm != nil {
			fInfo.queue.Remove(elm)
			nextReq := elm.Value.(*queuedRequest)
			nextReq.vft = fInfo.vt + nextReq.cost
			heap.Push(&sched.requests, nextReq)
			fInfo.requestOnHeap = true
		}
	}
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
