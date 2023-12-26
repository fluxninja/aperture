package concurrencylimiter

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	guuid "github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	tokencounterv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/tokencounter/v1"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	deadlinemargin "github.com/fluxninja/aperture/v2/pkg/dmap-funcs/deadline-margin"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

const (
	// TakeNFunction is the name of the function used to take N tokens from the global counter.
	TakeNFunction = "TakeN"
	// ReturnFunction is the name of the function used to return N tokens to the global counter.
	ReturnFunction = "Return"
	// CancelQueuedFunction is the name of the function used to cancel a queued request.
	CancelQueuedFunction = "CancelQueued"
	// CancelInflightFunction is the name of the function used to cancel an inflight request.
	CancelInflightFunction = "CancelInflight"
	// MinimumWaitTime is the minimum wait time for a request in the queue. This is used if the token rate is 0.
	MinimumWaitTime = 10 * time.Millisecond
	// MaximumWaitTime is the maximum wait time for a request in the queue.
	MaximumWaitTime = 5 * time.Second
	// WaitExpiryMargin is the margin for expiry in addition to wait time for a request in the queue.
	WaitExpiryMargin = 100 * time.Millisecond
	// TokenRateWindowSize is the number of requests after which token rate is updated.
	TokenRateWindowSize = 10
)

// GlobalTokenCounter is a global token counter that uses a distributed map to store the state.
type GlobalTokenCounter struct {
	dMap                olric.DMap
	dc                  *distcache.DistCache
	name                string
	capacity            float64
	maxInflightDuration time.Duration
	mu                  sync.RWMutex
	passThrough         bool
}

// GlobalCounter implements the ConcurrencyLimiter interface.
var _ ConcurrencyLimiter = (*GlobalTokenCounter)(nil)

// NewGlobalTokenCounter creates a new GlobalTokenCounter.
func NewGlobalTokenCounter(dc *distcache.DistCache,
	name string,
	maxIdleDuration time.Duration,
	maxInflightDuration time.Duration,
) (*GlobalTokenCounter, error) {
	gtc := &GlobalTokenCounter{
		dc:                  dc,
		name:                name,
		maxInflightDuration: maxInflightDuration,
	}

	dmapConfig := config.DMap{
		MaxIdleDuration: maxIdleDuration,
		Functions: map[string]config.Function{
			TakeNFunction:          gtc.takeN,
			ReturnFunction:         gtc.returnTokens,
			CancelQueuedFunction:   gtc.cancelQueued,
			CancelInflightFunction: gtc.cancelInflight,
		},
	}

	dmap, err := dc.NewDMap(name, dmapConfig, false)
	if err != nil {
		return nil, err
	}

	gtc.dMap = dmap

	return gtc, nil
}

// Name returns the name of the global token counter.
func (gtc *GlobalTokenCounter) Name() string {
	return gtc.name
}

// SetCapacity sets the capacity of the global token counter.
func (gtc *GlobalTokenCounter) SetCapacity(capacity float64) {
	gtc.mu.Lock()
	defer gtc.mu.Unlock()
	gtc.capacity = capacity
}

// GetCapacity returns the capacity of the global token counter.
func (gtc *GlobalTokenCounter) GetCapacity() float64 {
	gtc.mu.RLock()
	defer gtc.mu.RUnlock()
	return gtc.capacity
}

// Close closes the global token counter.
func (gtc *GlobalTokenCounter) Close() error {
	gtc.mu.Lock()
	defer gtc.mu.Unlock()
	err := gtc.dc.DeleteDMap(gtc.name)
	if err != nil {
		return err
	}
	return nil
}

// TakeIfAvailable takes N tokens from the global token counter if they are available.
func (gtc *GlobalTokenCounter) TakeIfAvailable(ctx context.Context, label string, n float64) (ok bool, waitTime time.Duration, remaining float64, current float64, reqID string) {
	ok = true
	if gtc.GetPassThrough() {
		return
	}

	if deadlinemargin.IsMarginExceeded(ctx) {
		ok = false
		return
	}

	req := tokencounterv1.TakeNRequest{
		Deadline:  timestamppb.New(time.Time{}),
		Tokens:    n,
		CanWait:   false,
		RequestId: guuid.NewString(),
	}

	reqBytes, err := req.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return
	}

	resultBytes, err := gtc.dMap.Function(ctx, label, TakeNFunction, reqBytes)
	if err != nil {
		log.Error().Err(err).Str("dmapName", gtc.dMap.Name()).Float64("tokens", n).Msg("error taking from token counter")
		return
	}

	var resp tokencounterv1.TakeNResponse
	err = resp.UnmarshalVT(resultBytes)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return
	}

	availableAt := resp.GetCheckBackAt().AsTime()
	if !availableAt.IsZero() {
		waitTime = time.Until(availableAt)
		if waitTime < 0 {
			waitTime = 0
		}
	}

	return resp.AvailableNow, waitTime, resp.Remaining, resp.Current, req.RequestId
}

// Take takes N tokens from the global token counter.
func (gtc *GlobalTokenCounter) Take(ctx context.Context, label string, count float64) (ok bool, _ time.Duration, remaining float64, current float64, reqID string) {
	ok = true
	if gtc.GetPassThrough() {
		return
	}

	deadline := time.Time{}
	d, isSet := ctx.Deadline()
	if isSet {
		deadline = d
	}

	req := tokencounterv1.TakeNRequest{
		Deadline:  timestamppb.New(deadline),
		Tokens:    count,
		CanWait:   true,
		RequestId: guuid.NewString(),
	}

	reqBytes, err := req.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		ok = true
		return
	}

	// Helper function to cancel the queued request if we are exiting early
	var isQueued bool
	cancelQueued := func() {
		if !isQueued {
			return
		}
		cancelQueuedReq := tokencounterv1.CancelQueuedRequest{
			RequestId: req.RequestId,
		}
		cancelQueuedReqBytes, err := cancelQueuedReq.MarshalVT()
		if err != nil {
			log.Autosample().Errorf("error encoding request: %v", err)
			return
		}
		_, err = gtc.dMap.Function(ctx, label, CancelQueuedFunction, cancelQueuedReqBytes)
		if err != nil {
			log.Autosample().Error().Err(err).Str("dmapName", gtc.dMap.Name()).Float64("tokens", count).Msg("error canceling queued request")
			return
		}
	}
	// TakeNRequest in a loop until we get a response that is ok or we exceed the deadline
	for {
		if deadlinemargin.IsMarginExceeded(ctx) {
			cancelQueued()
			ok = false
			return
		}

		resultBytes, err := gtc.dMap.Function(ctx, label, TakeNFunction, reqBytes)
		if err != nil {
			log.Autosample().Error().Err(err).Str("dmapName", gtc.dMap.Name()).Float64("tokens", count).Msg("error taking from token counter")
			cancelQueued()
			ok = true
			return
		}
		isQueued = true

		var resp tokencounterv1.TakeNResponse
		err = resp.UnmarshalVT(resultBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding response: %v", err)
			cancelQueued()
			ok = true
			return
		}

		if !resp.GetOk() {
			current = resp.GetCurrent()
			remaining = resp.GetRemaining()
			ok = false
			return
		}

		if resp.GetAvailableNow() {
			current = resp.GetCurrent()
			remaining = resp.GetRemaining()
			reqID = req.RequestId
			ok = true
			return
		}

		checkBackAt := resp.GetCheckBackAt().AsTime()
		var waitTime time.Duration
		if !checkBackAt.IsZero() {
			waitTime = time.Until(checkBackAt)
			if waitTime < 0 {
				waitTime = 0
			}
			// Try more aggressively than the wait time to increase utilization
			waitTime /= 2
		}

		if waitTime > MaximumWaitTime {
			waitTime = MaximumWaitTime
		}

		if waitTime < MinimumWaitTime {
			waitTime = MinimumWaitTime
		}

		if deadlinemargin.IsWaitMarginExceeded(ctx, waitTime) {
			cancelQueued()
			ok = false
			return
		}

		// Pause until waitTime has elapsed or context is canceled
		select {
		case <-time.After(waitTime):
		case <-ctx.Done():
			ok = false
			cancelQueued()
			return
		}
	}
}

// Return returns N tokens to the global token counter.
func (gtc *GlobalTokenCounter) Return(ctx context.Context, label string, count float64, reqID string) (ok bool, err error) {
	if reqID == "" {
		err = errors.New("request id cannot be empty")
		return
	}

	if gtc.GetPassThrough() {
		return
	}

	if deadlinemargin.IsMarginExceeded(ctx) {
		err = errors.New("deadline exceeded")
		return
	}

	req := tokencounterv1.ReturnNRequest{
		Tokens:    count,
		RequestId: reqID,
	}

	reqBytes, err := req.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding request: %v", err)
		return
	}

	resultBytes, err := gtc.dMap.Function(ctx, label, ReturnFunction, reqBytes)
	if err != nil {
		log.Autosample().Error().Err(err).Str("dmapName", gtc.dMap.Name()).Float64("tokens", count).Msg("error returning to token counter")
		return
	}

	var resp tokencounterv1.ReturnNResponse
	err = resp.UnmarshalVT(resultBytes)
	if err != nil {
		log.Autosample().Errorf("error decoding response: %v", err)
		return
	}
	return resp.GetOk(), err
}

// SetPassThrough sets the value of passThrough.
func (gtc *GlobalTokenCounter) SetPassThrough(passthrough bool) {
	gtc.mu.Lock()
	defer gtc.mu.Unlock()
	gtc.passThrough = passthrough
}

// GetPassThrough returns the current value of passThrough.
func (gtc *GlobalTokenCounter) GetPassThrough() bool {
	gtc.mu.RLock()
	defer gtc.mu.RUnlock()
	return gtc.passThrough
}

func (gtc *GlobalTokenCounter) takeN(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	gtc.mu.RLock()
	defer gtc.mu.RUnlock()

	state, err := gtc.decodeState(stateBytes, key)
	if err != nil {
		return nil, nil, err
	}

	// Decode takeNReq from proto encoded argBytes
	var takeNReq tokencounterv1.TakeNRequest
	if argBytes != nil {
		err = takeNReq.UnmarshalVT(argBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	now := time.Now()
	takeNResp := tokencounterv1.TakeNResponse{
		Ok:          true,
		CheckBackAt: timestamppb.New(now),
	}

	// Audit and find requestID in the queued requests
	var requestsQueued []*tokencounterv1.Request
	var found, headOfQueue bool
	var requestQueued *tokencounterv1.Request
	var reqIndex int
	var tokensAhead float64

	for i, r := range state.RequestsQueued {
		if isExpired(r, now) {
			continue
		}
		if !found && r.RequestId == takeNReq.RequestId && takeNReq.CanWait {
			found = true
			requestQueued = r
			reqIndex = i
			if i == 0 {
				headOfQueue = true
			}
		}
		if !found {
			tokensAhead += r.Tokens
		}
		requestsQueued = append(requestsQueued, r)
	}

	if !found && takeNReq.CanWait {
		// Add the request to the queue
		requestQueued = &tokencounterv1.Request{
			RequestId: takeNReq.RequestId,
			Tokens:    takeNReq.Tokens,
		}
		// Queue the request
		requestsQueued = append(requestsQueued, requestQueued)
		reqIndex = len(requestsQueued) - 1
	}

	var waitTime time.Duration
	if state.TokenRate > 0 {
		waitTime = time.Second * time.Duration(tokensAhead/state.TokenRate)
	} else {
		waitTime = MinimumWaitTime
	}

	if requestQueued != nil {
		requestQueued.NumRetries += 1
		if requestQueued.NumRetries%5 == 0 {
			// Increase wait time exponentially after every 5 re-tries
			newWaitTime := requestQueued.WaitFor.AsDuration() * 2
			if newWaitTime > waitTime {
				waitTime = newWaitTime
			}
		}
	}

	var requestsInflight []*tokencounterv1.Request
	var tokensInflight float64
	for _, r := range state.RequestsInflight {
		if isExpired(r, now) {
			continue
		}
		requestsInflight = append(requestsInflight, r)
		tokensInflight += r.Tokens
	}
	haveSpareCapacity := tokensInflight < gtc.capacity

	// Criteria for adding the request to the inflight requests
	if (len(requestsQueued) <= 1 || headOfQueue) && haveSpareCapacity {
		// Add the request to the inflight requests and mark it as available now
		takeNResp.AvailableNow = true
		requestsInflight = append(requestsInflight, &tokencounterv1.Request{
			RequestId: takeNReq.RequestId,
			Tokens:    takeNReq.Tokens,
			ExpiresAt: timestamppb.New(now.Add(gtc.maxInflightDuration)),
			WaitFor:   durationpb.New(gtc.maxInflightDuration),
		})
	}

	// Check if the request is queued
	if requestQueued != nil {
		effectiveWaitTime := waitTime + WaitExpiryMargin
		// update the expiry time of the request
		requestQueued.ExpiresAt = timestamppb.New(now.Add(effectiveWaitTime))
		requestQueued.WaitFor = durationpb.New(effectiveWaitTime)
		// Check if request has exceeded its deadline
		exceededDeadline := now.Add(waitTime).After(takeNReq.Deadline.AsTime())
		if exceededDeadline || takeNResp.AvailableNow {
			// Remove request from the queue if it has exceeded its deadline or if it is available now
			requestsQueued = append(requestsQueued[:reqIndex], requestsQueued[reqIndex+1:]...)
		}
		if exceededDeadline {
			// If the request has exceeded its deadline, mark it as not ok
			takeNResp.Ok = false
		} else {
			// Else, set the check back time
			takeNResp.CheckBackAt = timestamppb.New(now.Add(waitTime))
		}
	}

	state.RequestsQueued = requestsQueued
	state.RequestsInflight = requestsInflight

	stateBytes, err = state.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding state: %v", err)
		return nil, nil, err
	}

	takeNResp.Current = tokensInflight
	takeNResp.Remaining = gtc.capacity - tokensInflight

	takeNRespBytes, err := takeNResp.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding response: %v", err)
		return nil, nil, err
	}

	return stateBytes, takeNRespBytes, nil
}

func (gtc *GlobalTokenCounter) returnTokens(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	gtc.mu.RLock()
	defer gtc.mu.RUnlock()

	state, err := gtc.decodeState(stateBytes, key)
	if err != nil {
		return nil, nil, err
	}

	// Decode returnNReq from proto encoded argBytes
	var returnNReq tokencounterv1.ReturnNRequest
	if argBytes != nil {
		err = returnNReq.UnmarshalVT(argBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	now := time.Now()
	returnNResp := tokencounterv1.ReturnNResponse{}

	// Audit and find requestID in the inflight requests
	var requestsInflight []*tokencounterv1.Request
	var tokens float64

	for _, r := range state.RequestsInflight {
		if isExpired(r, now) {
			continue
		}
		if r.RequestId == returnNReq.RequestId {
			returnNResp.Ok = true
			tokens = r.Tokens
			continue
		}
		requestsInflight = append(requestsInflight, r)
	}
	// Fall back to the request tokens if we can't find the request in the inflight requests
	if tokens == 0 {
		tokens = returnNReq.Tokens
	}
	state.RequestsInflight = requestsInflight

	// Token rate calculation
	tokenWindow := state.TokenWindow
	// Check if we need to invalidate the token rate
	if tokenWindow.End != nil && now.Sub(tokenWindow.End.AsTime()) > 2*gtc.maxInflightDuration {
		state.TokenRate = 0
		tokenWindow.Start = nil
		tokenWindow.End = nil
		tokenWindow.Sum = 0
		tokenWindow.Count = 0
	}

	if tokenWindow.Start == nil {
		tokenWindow.Start = timestamppb.New(now)
	}
	tokenWindow.End = timestamppb.New(now)

	tokenWindow.Sum += tokens
	tokenWindow.Count += 1

	if tokenWindow.Count >= TokenRateWindowSize {
		if state.TokenRate == 0 {
			state.TokenRate = tokenWindow.Sum / tokenWindow.End.AsTime().Sub(tokenWindow.Start.AsTime()).Seconds()
		}
		tokenWindow.Start = tokenWindow.End
		tokenWindow.End = nil
		tokenWindow.Sum = 0
		tokenWindow.Count = 0
	}

	stateBytes, err = state.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding state: %v", err)
		return nil, nil, err
	}

	returnNRespBytes, err := returnNResp.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding response: %v", err)
		return nil, nil, err
	}

	return stateBytes, returnNRespBytes, nil
}

func (gtc *GlobalTokenCounter) cancelQueued(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	gtc.mu.RLock()
	defer gtc.mu.RUnlock()

	state, err := gtc.decodeState(stateBytes, key)
	if err != nil {
		return nil, nil, err
	}

	// Decode cancelQueuedReq from proto encoded argBytes
	var cancelQueuedReq tokencounterv1.CancelQueuedRequest
	if argBytes != nil {
		err = cancelQueuedReq.UnmarshalVT(argBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	now := time.Now()
	cancelQueuedResp := tokencounterv1.CancelQueuedResponse{}

	// Audit and find requestID in the queued requests
	var requestsQueued []*tokencounterv1.Request

	for _, r := range state.RequestsQueued {
		if isExpired(r, now) {
			continue
		}
		if r.RequestId == cancelQueuedReq.RequestId {
			cancelQueuedResp.Ok = true
			continue
		}
		requestsQueued = append(requestsQueued, r)
	}

	state.RequestsQueued = requestsQueued

	stateBytes, err = state.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding state: %v", err)
		return nil, nil, err
	}

	cancelQueuedRespBytes, err := cancelQueuedResp.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding response: %v", err)
		return nil, nil, err
	}

	return stateBytes, cancelQueuedRespBytes, nil
}

func (gtc *GlobalTokenCounter) cancelInflight(key string, stateBytes, argBytes []byte) ([]byte, []byte, error) {
	gtc.mu.RLock()
	defer gtc.mu.RUnlock()

	state, err := gtc.decodeState(stateBytes, key)
	if err != nil {
		return nil, nil, err
	}

	// Decode cancelInflightReq from proto encoded argBytes
	var cancelInflightReq tokencounterv1.CancelInflightRequest
	if argBytes != nil {
		err = cancelInflightReq.UnmarshalVT(argBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding arg: %v", err)
			return nil, nil, err
		}
	}

	now := time.Now()
	cancelInflightResp := tokencounterv1.CancelInflightResponse{}

	// Audit and find requestID in the inflight requests
	var requestsInflight []*tokencounterv1.Request

	for _, r := range state.RequestsInflight {
		if isExpired(r, now) {
			continue
		}
		if r.RequestId == cancelInflightReq.RequestId {
			cancelInflightResp.Ok = true
			continue
		}
		requestsInflight = append(requestsInflight, r)
	}

	state.RequestsInflight = requestsInflight

	stateBytes, err = state.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding state: %v", err)
		return nil, nil, err
	}

	cancelInflightRespBytes, err := cancelInflightResp.MarshalVT()
	if err != nil {
		log.Autosample().Errorf("error encoding response: %v", err)
		return nil, nil, err
	}

	return stateBytes, cancelInflightRespBytes, nil
}

func (gtc *GlobalTokenCounter) decodeState(stateBytes []byte, key string) (*tokencounterv1.State, error) {
	var state tokencounterv1.State

	if stateBytes != nil {
		err := state.UnmarshalVT(stateBytes)
		if err != nil {
			log.Autosample().Errorf("error decoding current state: %v", err)
			return nil, err
		}
	} else {
		log.Info().Msgf("Creating new token counter state for key %s in dmap %s", key, gtc.dMap.Name())
		state.TokenWindow = &tokencounterv1.State_TokenWindow{}
	}

	return &state, nil
}

func isExpired(r *tokencounterv1.Request, now time.Time) bool {
	return r.ExpiresAt.AsTime().Before(now)
}
