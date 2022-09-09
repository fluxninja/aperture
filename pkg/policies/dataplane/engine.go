package dataplane

import (
	"fmt"
	"sync"

	"golang.org/x/exp/maps"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

// multiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type multiMatchResult struct {
	concurrencyLimiters []iface.Limiter
	fluxMeters          []iface.FluxMeter
	rateLimiters        []iface.RateLimiter
	classifiers         []iface.Classifier
}

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[string, multiMatchResult]

// PopulateFromMultiMatcher populates result object with results from MultiMatcher.
func (result *multiMatchResult) populateFromMultiMatcher(mm *multimatcher.MultiMatcher[string, multiMatchResult], labels selectors.Labels) {
	resultCollection := mm.Match(multimatcher.Labels(labels.ToPlainMap()))
	result.concurrencyLimiters = append(result.concurrencyLimiters, resultCollection.concurrencyLimiters...)
	result.fluxMeters = append(result.fluxMeters, resultCollection.fluxMeters...)
	result.rateLimiters = append(result.rateLimiters, resultCollection.rateLimiters...)
	result.classifiers = append(result.classifiers, resultCollection.classifiers...)
}

// ProvideEngineAPI Main fx app.
func ProvideEngineAPI() iface.Engine {
	e := &Engine{
		multiMatchers: make(map[selectors.ControlPointID]*multiMatcher),
		fluxMetersMap: make(map[iface.FluxMeterID]iface.FluxMeter),
	}
	return e
}

// Engine APIs to
// (1) Get schedulers given a service, control point and set of labels
// (2) Get flux meter histogram given a metric id.
type Engine struct {
	fluxMeterMapMutex  sync.RWMutex
	fluxMetersMap      map[iface.FluxMeterID]iface.FluxMeter
	multiMatchersMutex sync.RWMutex
	multiMatchers      map[selectors.ControlPointID]*multiMatcher
}

// ProcessRequest .
func (e *Engine) ProcessRequest(controlPoint selectors.ControlPoint, serviceIDs []services.ServiceID, labels selectors.Labels) (response *flowcontrolv1.CheckResponse) {
	response = &flowcontrolv1.CheckResponse{
		DecisionType:  flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
		FlowLabelKeys: maps.Keys(labels),
	}

	mmr := e.getMatches(controlPoint, serviceIDs, labels)
	if mmr == nil {
		return
	}

	fluxMeters := mmr.fluxMeters
	fluxMeterProtos := make([]*flowcontrolv1.FluxMeter, len(fluxMeters))
	for i, fluxMeter := range fluxMeters {
		fluxMeterProtos[i] = &flowcontrolv1.FluxMeter{
			FluxMeterName: fluxMeter.GetFluxMeterName(),
		}
	}
	response.FluxMeters = fluxMeterProtos

	classifiers := mmr.classifiers
	classifierProtos := make([]*flowcontrolv1.Classifier, len(classifiers))
	for i, classifier := range classifiers {
		classifierProtos[i] = &flowcontrolv1.Classifier{
			PolicyName:      classifier.GetClassifierID().PolicyName,
			PolicyHash:      classifier.GetClassifierID().PolicyHash,
			ClassifierIndex: classifier.GetClassifierID().ClassifierIndex,
		}
	}
	response.Classifiers = classifierProtos

	// execute rate limiters first
	rateLimiters := make([]iface.Limiter, len(mmr.rateLimiters))
	for i, rl := range mmr.rateLimiters {
		rateLimiters[i] = rl
	}
	rateLimiterDecisions, rateLimitersDecisionType := runLimiters(rateLimiters, labels)
	response.LimiterDecisions = rateLimiterDecisions

	defer func() {
		if response.DecisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
			returnExtraTokens(mmr.rateLimiters, rateLimiterDecisions, labels)
		}
	}()

	// If any rate limiter dropped, then mark this as a decision reason and return.
	// Do not execute concurrency limiters.
	if rateLimitersDecisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
		response.DecisionType = rateLimitersDecisionType
		response.DecisionReason = &flowcontrolv1.DecisionReason{
			RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_RATE_LIMITED,
		}
		return
	}

	// execute concurrency limiters
	concurrencyLimiters := make([]iface.Limiter, len(mmr.concurrencyLimiters))
	copy(concurrencyLimiters, mmr.concurrencyLimiters)

	concurrencyLimiterDecisions, concurrencyLimitersDecisionType := runLimiters(concurrencyLimiters, labels)
	response.LimiterDecisions = append(response.LimiterDecisions, concurrencyLimiterDecisions...)

	if concurrencyLimitersDecisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
		response.DecisionType = flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED
		response.DecisionReason = &flowcontrolv1.DecisionReason{
			RejectReason: flowcontrolv1.DecisionReason_REJECT_REASON_CONCURRENCY_LIMITED,
		}
		return
	}

	return
}

func runLimiters(limiters []iface.Limiter, labels selectors.Labels) ([]*flowcontrolv1.LimiterDecision, flowcontrolv1.DecisionType) {
	var wg sync.WaitGroup
	var once sync.Once
	decisions := make([]*flowcontrolv1.LimiterDecision, len(limiters))

	decisionType := flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED

	setDecisionRejected := func() {
		decisionType = flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED
	}

	execLimiter := func(limiter iface.Limiter, i int) func() {
		return func() {
			defer wg.Done()
			decisions[i] = limiter.RunLimiter(labels)
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
	labels selectors.Labels,
) {
	for i, l := range rateLimiterDecisions {
		if !l.Dropped && l.Reason == flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED {
			go rateLimiters[i].TakeN(labels, -1)
		}
	}
}

// RegisterConcurrencyLimiter adds concurrency limiter to multimatcher.
func (e *Engine) RegisterConcurrencyLimiter(cl iface.Limiter) error {
	concurrencyLimiterMatchedCB := func(mmr multiMatchResult) multiMatchResult {
		mmr.concurrencyLimiters = append(mmr.concurrencyLimiters, cl)
		return mmr
	}
	return e.register("ConcurrencyLimiter:"+cl.GetLimiterID().String(), cl.GetSelector(), concurrencyLimiterMatchedCB)
}

// UnregisterConcurrencyLimiter removes concurrency limiter from multimatcher.
func (e *Engine) UnregisterConcurrencyLimiter(cl iface.Limiter) error {
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

// RegisterRateLimiter adds limiter actuator to multimatcher.
func (e *Engine) RegisterRateLimiter(rl iface.RateLimiter) error {
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
	selectorProto := rl.GetSelector()
	return e.unregister("RateLimiter:"+rl.GetLimiterID().String(), selectorProto)
}

// RegisterClassifier adds classifier to multimatcher.
func (e *Engine) RegisterClassifier(c iface.Classifier) error {
	classifierMatchedCB := func(mmr multiMatchResult) multiMatchResult {
		mmr.classifiers = append(mmr.classifiers, c)
		return mmr
	}
	return e.register("Classifier:"+c.GetClassifierID().String(), c.GetSelector(), classifierMatchedCB)
}

// UnregisterClassifier removes classifier from multimatcher.
func (e *Engine) UnregisterClassifier(c iface.Classifier) error {
	selectorProto := c.GetSelector()
	return e.unregister("Classifier:"+c.GetClassifierID().String(), selectorProto)
}

// getMatches returns schedulers and fluxmeters for given labels.
func (e *Engine) getMatches(controlPoint selectors.ControlPoint, serviceIDs []services.ServiceID, labels selectors.Labels) *multiMatchResult {
	e.multiMatchersMutex.RLock()
	defer e.multiMatchersMutex.RUnlock()

	mmResult := &multiMatchResult{}

	// Lookup catchall multi matchers for controlPoint
	controlPointID := selectors.ControlPointID{
		ControlPoint: controlPoint,
		ServiceID: services.ServiceID{
			Service: "",
		},
	}
	camm, ok := e.multiMatchers[controlPointID]
	if ok {
		mmResult.populateFromMultiMatcher(camm, labels)
	}

	for _, serviceID := range serviceIDs {
		controlPointID := selectors.ControlPointID{
			ControlPoint: controlPoint,
			ServiceID:    serviceID,
		}
		// Lookup multi matcher for controlPointID
		mm, ok := e.multiMatchers[controlPointID]
		if ok {
			mmResult.populateFromMultiMatcher(mm, labels)
		}
	}

	return mmResult
}

func (e *Engine) register(key string, selectorProto *selectorv1.Selector, matchedCB multimatcher.MatchCallback[multiMatchResult]) error {
	e.multiMatchersMutex.Lock()
	defer e.multiMatchersMutex.Unlock()

	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	mm, ok := e.multiMatchers[selector.ControlPointID]
	if !ok {
		mm = multimatcher.New[string, multiMatchResult]()
		e.multiMatchers[selector.ControlPointID] = mm
	}
	err = mm.AddEntry(key, selector.LabelMatcher, matchedCB)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) unregister(key string, selectorProto *selectorv1.Selector) error {
	e.multiMatchersMutex.Lock()
	defer e.multiMatchersMutex.Unlock()

	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	// check if multi matcher exists for this control point id
	mm, ok := e.multiMatchers[selector.ControlPointID]
	if !ok {
		return fmt.Errorf("unable to unregister, multi matcher not found for control point id: %v", selector.ControlPointID)
	}
	err = mm.RemoveEntry(key)
	if err != nil {
		return err
	}
	// remove this multi matcher if this was the last entry
	if mm.Length() == 0 {
		delete(e.multiMatchers, selector.ControlPointID)
	}

	return nil
}
