package flowcontrol

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/exp/maps"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	concurrencyLimiters []iface.ConcurrencyLimiter
	fluxMeters          []iface.FluxMeter
	rateLimiters        []iface.RateLimiter
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[string, multiMatchResult]

// PopulateFromMultiMatcher populates result object with results from MultiMatcher.
func (result *multiMatchResult) populateFromMultiMatcher(mm *multimatcher.MultiMatcher[string, multiMatchResult], labels map[string]string) {
	resultCollection := mm.Match(multimatcher.Labels(labels))
	result.concurrencyLimiters = append(result.concurrencyLimiters, resultCollection.concurrencyLimiters...)
	result.fluxMeters = append(result.fluxMeters, resultCollection.fluxMeters...)
	result.rateLimiters = append(result.rateLimiters, resultCollection.rateLimiters...)
}

// NewEngine Main fx app.
func NewEngine() iface.Engine {
	e := &Engine{
		multiMatchers:  make(map[selectors.ControlPointID]*multiMatcher),
		fluxMetersMap:  make(map[iface.FluxMeterID]iface.FluxMeter),
		conLimiterMap:  make(map[iface.LimiterID]iface.ConcurrencyLimiter),
		rateLimiterMap: make(map[iface.LimiterID]iface.RateLimiter),
	}
	return e
}

// Engine APIs to
// (1) Get schedulers given a service, control point and set of labels.
// (2) Get flux meter histogram given a metric id.
type Engine struct {
	fluxMeterMapMutex   sync.RWMutex
	fluxMetersMap       map[iface.FluxMeterID]iface.FluxMeter
	conLimiterMapMutex  sync.RWMutex
	conLimiterMap       map[iface.LimiterID]iface.ConcurrencyLimiter
	rateLimiterMapMutex sync.RWMutex
	rateLimiterMap      map[iface.LimiterID]iface.RateLimiter
	multiMatchersMutex  sync.RWMutex
	multiMatchers       map[selectors.ControlPointID]*multiMatcher
}

// ProcessRequest .
func (e *Engine) ProcessRequest(
	ctx context.Context,
	controlPoint string,
	serviceIDs []string,
	labels map[string]string,
) (response *flowcontrolv1.CheckResponse) {
	response = &flowcontrolv1.CheckResponse{
		DecisionType:  flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys: maps.Keys(labels),
		Services:      serviceIDs,
		ControlPoint:  controlPoint,
	}

	mmr := e.getMatches(controlPoint, serviceIDs, labels)
	if mmr == nil {
		return
	}

	fluxMeters := mmr.fluxMeters
	fluxMeterProtos := make([]*flowcontrolv1.FluxMeterInfo, len(fluxMeters))
	for i, fluxMeter := range fluxMeters {
		fluxMeterProtos[i] = &flowcontrolv1.FluxMeterInfo{
			FluxMeterName: fluxMeter.GetFluxMeterName(),
		}
	}
	response.FluxMeterInfos = fluxMeterProtos

	// execute rate limiters first
	rateLimiters := make([]iface.Limiter, len(mmr.rateLimiters))
	for i, rl := range mmr.rateLimiters {
		rateLimiters[i] = rl
	}

	rateLimiterDecisions, rateLimitersDecisionType := runLimiters(ctx, rateLimiters, labels)
	response.LimiterDecisions = rateLimiterDecisions

	defer func() {
		if response.DecisionType == flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED {
			returnExtraTokens(mmr.rateLimiters, rateLimiterDecisions, labels)
		}
	}()

	// If any rate limiter dropped, then mark this as a decision reason and return.
	// Do not execute concurrency limiters.
	if rateLimitersDecisionType == flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED {
		response.DecisionType = rateLimitersDecisionType
		response.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED
		return
	}

	// execute concurrency limiters
	concurrencyLimiters := make([]iface.Limiter, len(mmr.concurrencyLimiters))
	for i, cl := range mmr.concurrencyLimiters {
		concurrencyLimiters[i] = cl
	}

	concurrencyLimiterDecisions, concurrencyLimitersDecisionType := runLimiters(ctx, concurrencyLimiters, labels)
	response.LimiterDecisions = append(response.LimiterDecisions, concurrencyLimiterDecisions...)

	if concurrencyLimitersDecisionType == flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED {
		response.DecisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		response.RejectReason = flowcontrolv1.CheckResponse_REJECT_REASON_CONCURRENCY_LIMITED
		return
	}

	return
}

func runLimiters(ctx context.Context, limiters []iface.Limiter, labels map[string]string) ([]*flowcontrolv1.LimiterDecision,
	flowcontrolv1.CheckResponse_DecisionType,
) {
	var wg sync.WaitGroup
	var once sync.Once
	decisions := make([]*flowcontrolv1.LimiterDecision, len(limiters))

	decisionType := flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED

	setDecisionRejected := func() {
		decisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
	}

	execLimiter := func(limiter iface.Limiter, i int) func() {
		return func() {
			defer wg.Done()
			decisions[i] = limiter.RunLimiter(ctx, labels)
			if decisions[i].Dropped {
				once.Do(setDecisionRejected)
			}
		}
	}

	// execute limiters
	for i, limiter := range limiters {
		wg.Add(1)
		if i == len(limiters)-1 {
			execLimiter(limiter, i)()
		} else {
			panichandler.Go(execLimiter(limiter, i))
		}
	}
	wg.Wait()

	return decisions, decisionType
}

func returnExtraTokens(
	rateLimiters []iface.RateLimiter,
	rateLimiterDecisions []*flowcontrolv1.LimiterDecision,
	labels map[string]string,
) {
	for i, l := range rateLimiterDecisions {
		if !l.Dropped && l.Reason == flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED {
			go rateLimiters[i].TakeN(labels, -1)
		}
	}
}

// RegisterConcurrencyLimiter adds concurrency limiter to multimatcher.
func (e *Engine) RegisterConcurrencyLimiter(cl iface.ConcurrencyLimiter) error {
	e.conLimiterMapMutex.Lock()
	defer e.conLimiterMapMutex.Unlock()
	if _, ok := e.conLimiterMap[cl.GetLimiterID()]; !ok {
		e.conLimiterMap[cl.GetLimiterID()] = cl
	} else {
		return fmt.Errorf("metric id already registered")
	}

	concurrencyLimiterMatchedCB := func(mmr multiMatchResult) multiMatchResult {
		mmr.concurrencyLimiters = append(mmr.concurrencyLimiters, cl)
		return mmr
	}
	return e.register("ConcurrencyLimiter:"+cl.GetLimiterID().String(), cl.GetSelector(), concurrencyLimiterMatchedCB)
}

// UnregisterConcurrencyLimiter removes concurrency limiter from multimatcher.
func (e *Engine) UnregisterConcurrencyLimiter(cl iface.ConcurrencyLimiter) error {
	e.conLimiterMapMutex.Lock()
	defer e.conLimiterMapMutex.Unlock()
	delete(e.conLimiterMap, cl.GetLimiterID())

	selectorProto := cl.GetSelector()
	return e.unregister("ConcurrencyLimiter:"+cl.GetLimiterID().String(), selectorProto)
}

// RegisterFluxMeter adds fluxmeter to histogram map and multimatcher.
func (e *Engine) RegisterFluxMeter(fm iface.FluxMeter) error {
	// Save the histogram in fluxMeterHists indexed by metric id
	e.fluxMeterMapMutex.Lock()
	defer e.fluxMeterMapMutex.Unlock()
	if _, ok := e.fluxMetersMap[fm.GetFluxMeterID()]; !ok {
		e.fluxMetersMap[fm.GetFluxMeterID()] = fm
	} else {
		return fmt.Errorf("metric id already registered")
	}

	// Save the fluxMeterAPI in multiMatchers
	fluxMeterMatchedCB := func(mmr multiMatchResult) multiMatchResult {
		mmr.fluxMeters = append(mmr.fluxMeters, fm)
		return mmr
	}

	selectorProto := fm.GetSelector()
	return e.register("FluxMeter:"+fm.GetFluxMeterID().String(), selectorProto, fluxMeterMatchedCB)
}

// UnregisterFluxMeter removes fluxmeter from histogram map and multimatcher.
func (e *Engine) UnregisterFluxMeter(fm iface.FluxMeter) error {
	// Remove the histogram from fluxMeterHists indexed by metric id
	e.fluxMeterMapMutex.Lock()
	defer e.fluxMeterMapMutex.Unlock()
	delete(e.fluxMetersMap, fm.GetFluxMeterID())

	// Remove the fluxMeterAPI from multiMatchers
	selectorProto := fm.GetSelector()
	return e.unregister("FluxMeter:"+fm.GetFluxMeterID().String(), selectorProto)
}

// GetFluxMeter Lookup function for getting flux meter.
func (e *Engine) GetFluxMeter(fluxMeterName string) iface.FluxMeter {
	e.fluxMeterMapMutex.RLock()
	defer e.fluxMeterMapMutex.RUnlock()
	fmID := iface.FluxMeterID{
		FluxMeterName: fluxMeterName,
	}
	return e.fluxMetersMap[fmID]
}

// GetConcurrencyLimiter Lookup function for getting concurrency limiter.
func (e *Engine) GetConcurrencyLimiter(limiterID iface.LimiterID) iface.ConcurrencyLimiter {
	e.conLimiterMapMutex.RLock()
	defer e.conLimiterMapMutex.RUnlock()
	return e.conLimiterMap[limiterID]
}

// RegisterRateLimiter adds limiter actuator to multimatcher.
func (e *Engine) RegisterRateLimiter(rl iface.RateLimiter) error {
	e.rateLimiterMapMutex.Lock()
	defer e.rateLimiterMapMutex.Unlock()
	if _, ok := e.rateLimiterMap[rl.GetLimiterID()]; !ok {
		e.rateLimiterMap[rl.GetLimiterID()] = rl
	} else {
		return fmt.Errorf("metric id already registered")
	}

	limiterActuatorMatchedCB := func(mmr multiMatchResult) multiMatchResult {
		mmr.rateLimiters = append(
			mmr.rateLimiters,
			rl,
		)
		return mmr
	}

	return e.register("RateLimiter:"+rl.GetLimiterID().String(), rl.GetSelector(), limiterActuatorMatchedCB)
}

// UnregisterRateLimiter removes limiter actuator from multimatcher.
func (e *Engine) UnregisterRateLimiter(rl iface.RateLimiter) error {
	e.rateLimiterMapMutex.Lock()
	defer e.rateLimiterMapMutex.Unlock()
	delete(e.rateLimiterMap, rl.GetLimiterID())

	selectorProto := rl.GetSelector()
	return e.unregister("RateLimiter:"+rl.GetLimiterID().String(), selectorProto)
}

// GetRateLimiter Lookup function for getting rate limiter.
func (e *Engine) GetRateLimiter(limiterID iface.LimiterID) iface.RateLimiter {
	e.rateLimiterMapMutex.RLock()
	defer e.rateLimiterMapMutex.RUnlock()
	return e.rateLimiterMap[limiterID]
}

// getMatches returns schedulers and fluxmeters for given labels.
func (e *Engine) getMatches(controlPoint string, serviceIDs []string, labels map[string]string) *multiMatchResult {
	e.multiMatchersMutex.RLock()
	defer e.multiMatchersMutex.RUnlock()

	mmResult := &multiMatchResult{}

	// Lookup catchall multi matchers for controlPoint
	controlPointID := selectors.NewControlPointID("", controlPoint)
	camm, ok := e.multiMatchers[controlPointID]
	if ok {
		mmResult.populateFromMultiMatcher(camm, labels)
	}

	for _, serviceID := range serviceIDs {
		controlPointID := selectors.NewControlPointID(serviceID, controlPoint)
		// Lookup multi matcher for controlPointID
		mm, ok := e.multiMatchers[controlPointID]
		if ok {
			mmResult.populateFromMultiMatcher(mm, labels)
		}
	}

	return mmResult
}

func (e *Engine) register(key string, selectorProto *policylangv1.Selector, matchedCB multimatcher.MatchCallback[multiMatchResult]) error {
	e.multiMatchersMutex.Lock()
	defer e.multiMatchersMutex.Unlock()

	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	mm, ok := e.multiMatchers[selector.ControlPointID()]
	if !ok {
		mm = multimatcher.New[string, multiMatchResult]()
		e.multiMatchers[selector.ControlPointID()] = mm
	}
	err = mm.AddEntry(key, selector.LabelMatcher(), matchedCB)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) unregister(key string, selectorProto *policylangv1.Selector) error {
	e.multiMatchersMutex.Lock()
	defer e.multiMatchersMutex.Unlock()

	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	// check if multi matcher exists for this control point id
	mm, ok := e.multiMatchers[selector.ControlPointID()]
	if !ok {
		return fmt.Errorf("unable to unregister, multi matcher not found for control point id")
	}
	err = mm.RemoveEntry(key)
	if err != nil {
		return err
	}
	// remove this multi matcher if this was the last entry
	if mm.Length() == 0 {
		delete(e.multiMatchers, selector.ControlPointID())
	}

	return nil
}
