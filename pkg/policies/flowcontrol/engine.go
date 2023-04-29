package flowcontrol

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/exp/maps"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/consts"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	loadSchedulers map[iface.Limiter]struct{}
	fluxMeters     map[iface.FluxMeter]struct{}
	rateLimiters   map[iface.Limiter]struct{}
	regulators     map[iface.Limiter]struct{}
	labelPreviews  map[iface.LabelPreview]struct{}
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[string, *multiMatchResult]

// NewEngine Main fx app.
func NewEngine() iface.Engine {
	e := &Engine{
		multiMatchers:   make(map[selectors.ControlPointID]*multiMatcher),
		fluxMetersMap:   make(map[iface.FluxMeterID]iface.FluxMeter),
		conLimiterMap:   make(map[iface.LimiterID]iface.LoadScheduler),
		rateLimiterMap:  make(map[iface.LimiterID]iface.RateLimiter),
		regulatorMap:    make(map[iface.LimiterID]iface.Limiter),
		labelPreviewMap: make(map[iface.PreviewID]iface.LabelPreview),
	}
	return e
}

// Engine APIs to
// (1) Get schedulers given a service, control point and set of labels.
// (2) Get flux meter histogram given a metric id.
type Engine struct {
	mutex           sync.RWMutex
	fluxMetersMap   map[iface.FluxMeterID]iface.FluxMeter
	conLimiterMap   map[iface.LimiterID]iface.LoadScheduler
	rateLimiterMap  map[iface.LimiterID]iface.RateLimiter
	regulatorMap    map[iface.LimiterID]iface.Limiter
	labelPreviewMap map[iface.PreviewID]iface.LabelPreview
	multiMatchers   map[selectors.ControlPointID]*multiMatcher
}

// ProcessRequest .
func (e *Engine) ProcessRequest(
	ctx context.Context,
	requestContext iface.RequestContext,
) (response *flowcontrolv1.CheckResponse) {
	controlPoint := requestContext.ControlPoint
	services := requestContext.Services
	flowLabels := requestContext.FlowLabels
	labelKeys := maps.Keys(flowLabels)
	sort.Strings(labelKeys)
	response = &flowcontrolv1.CheckResponse{
		DecisionType:  flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys: labelKeys,
		Services:      services,
		ControlPoint:  controlPoint,
	}

	mmr := e.getMatches(controlPoint, services, flowLabels)
	if mmr == nil {
		return
	}

	labelPreviews := mmr.labelPreviews
	for labelPreview := range labelPreviews {
		labelPreview.AddLabelPreview(flowLabels)
	}

	fluxMeters := mmr.fluxMeters
	fluxMeterProtos := make([]*flowcontrolv1.FluxMeterInfo, len(fluxMeters))
	i := 0
	for fluxMeter := range fluxMeters {
		fluxMeterProtos[i] = &flowcontrolv1.FluxMeterInfo{
			FluxMeterName: fluxMeter.GetFluxMeterName(),
		}
		i++
	}
	response.FluxMeterInfos = fluxMeterProtos

	limiterTypes := []struct {
		limiters     map[iface.Limiter]struct{}
		rejectReason flowcontrolv1.CheckResponse_RejectReason
	}{
		{mmr.regulators, flowcontrolv1.CheckResponse_REJECT_REASON_REGULATED},
		{mmr.rateLimiters, flowcontrolv1.CheckResponse_REJECT_REASON_RATE_LIMITED},
		{mmr.loadSchedulers, flowcontrolv1.CheckResponse_REJECT_REASON_NO_TOKENS},
	}

	for _, limiterType := range limiterTypes {
		limiterDecisions, decisionType := runLimiters(ctx, limiterType.limiters, flowLabels)
		for _, limiterDecision := range limiterDecisions {
			response.LimiterDecisions = append(response.LimiterDecisions, limiterDecision)
		}

		defer func() {
			if response.DecisionType == flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED {
				revertRemaining(flowLabels, limiterDecisions)
			}
		}()

		if decisionType == flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED {
			response.DecisionType = decisionType
			response.RejectReason = limiterType.rejectReason
			return
		}
	}

	return
}

func runLimiters(ctx context.Context, limiters map[iface.Limiter]struct{}, labels map[string]string) (
	map[iface.Limiter]*flowcontrolv1.LimiterDecision,
	flowcontrolv1.CheckResponse_DecisionType,
) {
	var wg sync.WaitGroup
	var once sync.Once

	// make a child context with a cancel function
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	decisions := make(map[iface.Limiter]*flowcontrolv1.LimiterDecision, len(limiters))

	decisionType := flowcontrolv1.CheckResponse_DECISION_TYPE_ACCEPTED

	setDecisionDropped := func() {
		decisionType = flowcontrolv1.CheckResponse_DECISION_TYPE_REJECTED
		// cancel all the limiter calls
		cancel()
	}

	lock := &sync.Mutex{}

	execLimiter := func(limiter iface.Limiter) func() {
		return func() {
			defer wg.Done()
			decision := limiter.Decide(ctx, labels)
			if decision.Dropped {
				once.Do(setDecisionDropped)
			}
			lock.Lock()
			defer lock.Unlock()
			decisions[limiter] = decision
		}
	}

	// execute limiters
	i := 0
	for limiter := range limiters {
		wg.Add(1)
		if i == len(limiters)-1 {
			execLimiter(limiter)()
		} else {
			panichandler.Go(execLimiter(limiter))
		}
		i++
	}
	wg.Wait()

	return decisions, decisionType
}

func revertRemaining(
	labels map[string]string,
	limiterDecisions map[iface.Limiter]*flowcontrolv1.LimiterDecision,
) {
	labelsCopy := make(map[string]string, len(labels))
	for k, v := range labels {
		labelsCopy[k] = v
	}
	for l, d := range limiterDecisions {
		if !d.Dropped && d.Reason == flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED {
			go l.Revert(labelsCopy, d)
		}
	}
}

// RegisterLoadScheduler adds load scheduler to multimatcher.
func (e *Engine) RegisterLoadScheduler(cl iface.LoadScheduler) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.conLimiterMap[cl.GetLimiterID()]; !ok {
		e.conLimiterMap[cl.GetLimiterID()] = cl
	} else {
		return fmt.Errorf("load scheduler already registered")
	}

	loadSchedulerMatchedCB := func(mmr *multiMatchResult) *multiMatchResult {
		if mmr.loadSchedulers == nil {
			mmr.loadSchedulers = make(map[iface.Limiter]struct{})
		}
		mmr.loadSchedulers[cl] = struct{}{}
		return mmr
	}
	return e.register("LoadScheduler:"+cl.GetLimiterID().String(), cl.GetFlowSelector(), loadSchedulerMatchedCB)
}

// UnregisterLoadScheduler removes load scheduler from multimatcher.
func (e *Engine) UnregisterLoadScheduler(cl iface.LoadScheduler) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.conLimiterMap, cl.GetLimiterID())

	return e.unregister("LoadScheduler:"+cl.GetLimiterID().String(), cl.GetFlowSelector())
}

// RegisterFluxMeter adds fluxmeter to histogram map and multimatcher.
func (e *Engine) RegisterFluxMeter(fm iface.FluxMeter) error {
	// Save the histogram in fluxMeterHists indexed by metric id
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.fluxMetersMap[fm.GetFluxMeterID()]; !ok {
		e.fluxMetersMap[fm.GetFluxMeterID()] = fm
	} else {
		return fmt.Errorf("fluxmeter already registered")
	}

	// Save the fluxMeterAPI in multiMatchers
	fluxMeterMatchedCB := func(mmr *multiMatchResult) *multiMatchResult {
		if mmr.fluxMeters == nil {
			mmr.fluxMeters = make(map[iface.FluxMeter]struct{})
		}
		mmr.fluxMeters[fm] = struct{}{}
		return mmr
	}

	flowSelectorProto := fm.GetFlowSelector()
	return e.register("FluxMeter:"+fm.GetFluxMeterID().String(), flowSelectorProto, fluxMeterMatchedCB)
}

// UnregisterFluxMeter removes fluxmeter from histogram map and multimatcher.
func (e *Engine) UnregisterFluxMeter(fm iface.FluxMeter) error {
	// Remove the histogram from fluxMeterHists indexed by metric id
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.fluxMetersMap, fm.GetFluxMeterID())

	// Remove the fluxMeterAPI from multiMatchers
	return e.unregister("FluxMeter:"+fm.GetFluxMeterID().String(), fm.GetFlowSelector())
}

// GetFluxMeter Lookup function for getting flux meter.
func (e *Engine) GetFluxMeter(fluxMeterName string) iface.FluxMeter {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	fmID := iface.FluxMeterID{
		FluxMeterName: fluxMeterName,
	}
	return e.fluxMetersMap[fmID]
}

// GetLoadScheduler Lookup function for getting load scheduler.
func (e *Engine) GetLoadScheduler(limiterID iface.LimiterID) iface.LoadScheduler {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.conLimiterMap[limiterID]
}

// RegisterRateLimiter adds limiter actuator to multimatcher.
func (e *Engine) RegisterRateLimiter(rl iface.RateLimiter) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.rateLimiterMap[rl.GetLimiterID()]; !ok {
		e.rateLimiterMap[rl.GetLimiterID()] = rl
	} else {
		return fmt.Errorf("rate limiter already registered")
	}

	limiterActuatorMatchedCB := func(mmr *multiMatchResult) *multiMatchResult {
		if mmr.rateLimiters == nil {
			mmr.rateLimiters = make(map[iface.Limiter]struct{})
		}
		mmr.rateLimiters[rl] = struct{}{}
		return mmr
	}

	return e.register("RateLimiter:"+rl.GetLimiterID().String(), rl.GetFlowSelector(), limiterActuatorMatchedCB)
}

// UnregisterRateLimiter removes limiter actuator from multimatcher.
func (e *Engine) UnregisterRateLimiter(rl iface.RateLimiter) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.rateLimiterMap, rl.GetLimiterID())

	return e.unregister("RateLimiter:"+rl.GetLimiterID().String(), rl.GetFlowSelector())
}

// GetRateLimiter Lookup function for getting rate limiter.
func (e *Engine) GetRateLimiter(limiterID iface.LimiterID) iface.RateLimiter {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.rateLimiterMap[limiterID]
}

// RegisterRegulator adds limiter actuator to multimatcher.
func (e *Engine) RegisterRegulator(l iface.Limiter) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.regulatorMap[l.GetLimiterID()]; !ok {
		e.regulatorMap[l.GetLimiterID()] = l
	} else {
		return fmt.Errorf("load regulator already registered")
	}

	regulatorMatchedCB := func(mmr *multiMatchResult) *multiMatchResult {
		if mmr.regulators == nil {
			mmr.regulators = make(map[iface.Limiter]struct{})
		}
		mmr.regulators[l] = struct{}{}
		return mmr
	}

	return e.register("Regulator:"+l.GetLimiterID().String(), l.GetFlowSelector(), regulatorMatchedCB)
}

// UnregisterRegulator removes limiter actuator from multimatcher.
func (e *Engine) UnregisterRegulator(rl iface.Limiter) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.regulatorMap, rl.GetLimiterID())

	return e.unregister("Regulator:"+rl.GetLimiterID().String(), rl.GetFlowSelector())
}

// GetRegulator Lookup function for getting load regulator.
func (e *Engine) GetRegulator(limiterID iface.LimiterID) iface.Limiter {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.regulatorMap[limiterID]
}

// RegisterLabelPreview adds label preview to multimatcher.
func (e *Engine) RegisterLabelPreview(lp iface.LabelPreview) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, ok := e.labelPreviewMap[lp.GetPreviewID()]; !ok {
		e.labelPreviewMap[lp.GetPreviewID()] = lp
	} else {
		return fmt.Errorf("label preview already registered")
	}

	labelPreviewMatchedCB := func(mmr *multiMatchResult) *multiMatchResult {
		if mmr.labelPreviews == nil {
			mmr.labelPreviews = make(map[iface.LabelPreview]struct{})
		}
		mmr.labelPreviews[lp] = struct{}{}
		return mmr
	}
	return e.register("LabelPreview:"+lp.GetPreviewID().String(), lp.GetFlowSelector(), labelPreviewMatchedCB)
}

// UnregisterLabelPreview removes label preview from multimatcher.
func (e *Engine) UnregisterLabelPreview(lp iface.LabelPreview) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.labelPreviewMap, lp.GetPreviewID())

	return e.unregister("LabelPreview:"+lp.GetPreviewID().String(), lp.GetFlowSelector())
}

// getMatches returns schedulers and fluxmeters for given labels.
func (e *Engine) getMatches(controlPoint string, serviceIDs []string, labels map[string]string) *multiMatchResult {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	mmResult := &multiMatchResult{}

	// Lookup any service multi matchers for controlPoint
	controlPointID := selectors.NewControlPointID(controlPoint, consts.AnyService)
	camm, ok := e.multiMatchers[controlPointID]
	if ok {
		camm.MatchWithResultCollection(multimatcher.Labels(labels), mmResult)
	}

	for _, serviceID := range serviceIDs {
		controlPointID := selectors.NewControlPointID(controlPoint, serviceID)
		// Lookup multi matcher for controlPointID
		mm, ok := e.multiMatchers[controlPointID]
		if ok {
			mm.MatchWithResultCollection(multimatcher.Labels(labels), mmResult)
		}
	}

	return mmResult
}

// Lock must be held by caller.
func (e *Engine) register(key string, flowSelectorProto *policylangv1.FlowSelector,
	matchedCB multimatcher.MatchCallback[*multiMatchResult],
) error {
	selector, err := selectors.FromProto(flowSelectorProto)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	mm, ok := e.multiMatchers[selector.ControlPointID()]
	if !ok {
		mm = multimatcher.New[string, *multiMatchResult]()
		e.multiMatchers[selector.ControlPointID()] = mm
	}
	err = mm.AddEntry(key, selector.LabelMatcher(), matchedCB)
	if err != nil {
		return err
	}

	return nil
}

// Lock must be held by caller.
func (e *Engine) unregister(key string, flowSelectorProto *policylangv1.FlowSelector) error {
	selector, err := selectors.FromProto(flowSelectorProto)
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
