package dataplane

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

// multiMatcher is MultiMatcher instantiation used in this package.
type multiMatcher = multimatcher.MultiMatcher[string, iface.MultiMatchResult]

// ProvideEngineAPI Main fx app.
func ProvideEngineAPI() iface.EngineAPI {
	e := &Engine{
		multiMatchers:  make(map[selectors.ControlPointID]*multiMatcher),
		fluxMeterHists: make(map[string]prometheus.Histogram),
	}
	return e
}

// Engine APIs to
// (1) Get schedulers given a service, control point and set of labels
// (2) Get flux meter histogram given a metric id.
type Engine struct {
	fluxMeterHistMutex sync.RWMutex
	fluxMeterHists     map[string]prometheus.Histogram
	multiMatchersMutex sync.RWMutex
	multiMatchers      map[selectors.ControlPointID]*multiMatcher
}

// ProcessRequest .
func (e *Engine) ProcessRequest(controlPoint selectors.ControlPoint, serviceIDs []services.ServiceID, labels selectors.Labels) (response *flowcontrolv1.CheckResponse) {
	response = &flowcontrolv1.CheckResponse{
		DecisionType: flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED,
	}

	multiMatchResult := e.getMatches(controlPoint, serviceIDs, labels)
	if multiMatchResult == nil {
		return
	}

	rawFluxMeters := multiMatchResult.FluxMeters
	fluxMeters := make([]*flowcontrolv1.FluxMeter, len(rawFluxMeters))
	for i, rawFluxMeter := range rawFluxMeters {
		fluxMeters[i] = &flowcontrolv1.FluxMeter{
			AgentGroup:    rawFluxMeter.GetAgentGroup(),
			PolicyName:    rawFluxMeter.GetPolicyName(),
			PolicyHash:    rawFluxMeter.GetPolicyHash(),
			FluxMeterName: rawFluxMeter.GetMetricName(),
		}
	}
	response.FluxMeters = fluxMeters

	// execute rate limiters first
	rateLimiters := make([]iface.Limiter, len(multiMatchResult.RateLimiters))
	for i, rl := range multiMatchResult.RateLimiters {
		rateLimiters[i] = rl
	}
	rateLimiterDecisions, rateLimitersDecisionType := runLimiters(rateLimiters, labels)
	response.LimiterDecisions = rateLimiterDecisions

	defer func() {
		if response.DecisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
			returnExtraTokens(multiMatchResult.RateLimiters, rateLimiterDecisions, labels)
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

	// execute rate limiters first
	concurrencyLimiters := make([]iface.Limiter, len(multiMatchResult.ConcurrencyLimiters))
	copy(concurrencyLimiters, multiMatchResult.ConcurrencyLimiters)

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
	decisionType := flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED
	var wg sync.WaitGroup
	limiterDecisions := make([]*flowcontrolv1.LimiterDecision, len(limiters))
	for i, limiter := range limiters {
		wg.Add(1)
		panichandler.Go(func() {
			defer wg.Done()
			decision := limiter.RunLimiter(labels)
			if decision.Dropped {
				decisionType = flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED
			}
			limiterDecisions[i] = decision
		})
	}
	wg.Wait()
	return limiterDecisions, decisionType
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
	concurrencyLimiterMatchedCB := func(mmr iface.MultiMatchResult) iface.MultiMatchResult {
		mmr.ConcurrencyLimiters = append(
			mmr.ConcurrencyLimiters,
			cl,
		)
		return mmr
	}
	return e.register("ConcurrencyLimiter:"+cl.GetPolicyName(), cl.GetSelector(), concurrencyLimiterMatchedCB)
}

// UnregisterConcurrencyLimiter removes concurrency limiter from multimatcher.
func (e *Engine) UnregisterConcurrencyLimiter(cl iface.Limiter) error {
	selectorProto := cl.GetSelector()
	return e.unregister("ConcurrencyLimiter:"+cl.GetPolicyName(), selectorProto)
}

// RegisterFluxMeter adds fluxmeter to histogram map.
func (e *Engine) RegisterFluxMeter(fm iface.FluxMeter) error {
	// Save the histogram in fluxMeterHists indexed by metric id
	e.fluxMeterHistMutex.Lock()
	defer e.fluxMeterHistMutex.Unlock()
	if _, ok := e.fluxMeterHists[fm.GetMetricID()]; !ok {
		e.fluxMeterHists[fm.GetMetricID()] = fm.GetHistogram()
	} else {
		return fmt.Errorf("metric id already registered")
	}

	// Save the fluxMeterAPI in multiMatchers
	fluxMeterMatchedCB := func(mmr iface.MultiMatchResult) iface.MultiMatchResult {
		mmr.FluxMeters = append(
			mmr.FluxMeters,
			fm,
		)
		return mmr
	}

	selectorProto := fm.GetSelector()
	return e.register("FluxMeter:"+fm.GetMetricID(), selectorProto, fluxMeterMatchedCB)
}

// UnregisterFluxMeter removes fluxmeter from histogram map.
func (e *Engine) UnregisterFluxMeter(fm iface.FluxMeter) error {
	// Remove the histogram from fluxMeterHists indexed by metric id
	e.fluxMeterHistMutex.Lock()
	defer e.fluxMeterHistMutex.Unlock()
	delete(e.fluxMeterHists, fm.GetMetricID())

	// Remove the fluxMeterAPI from multiMatchers
	selectorProto := fm.GetSelector()
	return e.unregister("FluxMeter:"+fm.GetMetricID(), selectorProto)
}

// GetFluxMeterHist Lookup function for getting histogram by metric id.
// TODO: this method should move under policies/dataplane/fluxmeter.
func (e *Engine) GetFluxMeterHist(metricID string) prometheus.Histogram {
	e.fluxMeterHistMutex.RLock()
	defer e.fluxMeterHistMutex.RUnlock()
	return e.fluxMeterHists[metricID]
}

// RegisterRateLimiter adds limiter actuator to multimatcher.
func (e *Engine) RegisterRateLimiter(rl iface.RateLimiter) error {
	limiterActuatorMatchedCB := func(mmr iface.MultiMatchResult) iface.MultiMatchResult {
		mmr.RateLimiters = append(
			mmr.RateLimiters,
			rl,
		)
		return mmr
	}

	return e.register("RateLimiter:"+rl.GetPolicyName(), rl.GetSelector(), limiterActuatorMatchedCB)
}

// UnregisterRateLimiter removes limiter actuator from multimatcher.
func (e *Engine) UnregisterRateLimiter(rl iface.RateLimiter) error {
	selectorProto := rl.GetSelector()
	return e.unregister("RateLimiter:"+rl.GetPolicyName(), selectorProto)
}

// getMatches returns schedulers and fluxmeters for given labels.
func (e *Engine) getMatches(controlPoint selectors.ControlPoint, serviceIDs []services.ServiceID, labels selectors.Labels) *iface.MultiMatchResult {
	e.multiMatchersMutex.RLock()
	defer e.multiMatchersMutex.RUnlock()

	mmResult := &iface.MultiMatchResult{}

	// Lookup catchall multi matchers for controlPoint
	controlPointID := selectors.ControlPointID{
		ControlPoint: controlPoint,
		ServiceID: services.ServiceID{
			Service: "",
		},
	}
	camm, ok := e.multiMatchers[controlPointID]
	if ok {
		mmResult.PopulateFromMultiMatcher(camm, labels)
	}

	for _, serviceID := range serviceIDs {
		controlPointID := selectors.ControlPointID{
			ControlPoint: controlPoint,
			ServiceID:    serviceID,
		}
		// Lookup multi matcher for controlPointID
		mm, ok := e.multiMatchers[controlPointID]
		if ok {
			mmResult.PopulateFromMultiMatcher(mm, labels)
		}
	}

	return mmResult
}

func (e *Engine) register(key string, selectorProto *policylangv1.Selector, matchedCB multimatcher.MatchCallback[iface.MultiMatchResult]) error {
	e.multiMatchersMutex.Lock()
	defer e.multiMatchersMutex.Unlock()

	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse selector")
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	mm, ok := e.multiMatchers[selector.ControlPointID]
	if !ok {
		mm = multimatcher.New[string, iface.MultiMatchResult]()
		e.multiMatchers[selector.ControlPointID] = mm
	}
	err = mm.AddEntry(key, selector.LabelMatcher, matchedCB)
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
		log.Warn().Err(err).Msg("Failed to parse selector")
		return fmt.Errorf("failed to parse selector: %v", err)
	}

	// check if multi matcher exists for this control point id
	mm, ok := e.multiMatchers[selector.ControlPointID]
	if !ok {
		log.Warn().Msg("Unable to unregister, multi matcher not found for control point id")
		return nil
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
