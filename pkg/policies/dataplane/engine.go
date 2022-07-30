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
	e := &Engine{}
	e.multiMatchers = make(map[selectors.ControlPointID]*multiMatcher)
	e.fluxMeterHists = make(map[string]prometheus.Histogram)
	return e
}

// Engine APIs to
// (1) Get schedulers given a service, control point and set of labels
// (2) Get flux meter histogram given a metric id
// TODO: Will implement 3 APIs described in policy-config-api.go.
type Engine struct {
	fluxMeterHistMutex sync.RWMutex
	multiMatchMutex    sync.RWMutex
	fluxMeterHists     map[string]prometheus.Histogram
	multiMatchers      map[selectors.ControlPointID]*multiMatcher
}

// ProcessRequest .
func (e *Engine) ProcessRequest(controlPoint selectors.ControlPoint, serviceIDs []services.ServiceID, labels selectors.Labels) *flowcontrolv1.CheckResponse {
	resp := &flowcontrolv1.CheckResponse{}

	decisionType := flowcontrolv1.DecisionType_DECISION_TYPE_ACCEPTED
	var decisionReason *flowcontrolv1.Reason

	setDecisionRejected := func() {
		decisionType = flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED
	}

	// Run applicable service protection policies (schedulers) in parallel and reject request if any policy decides to reject
	multiMatchResult := e.getMatches(controlPoint, serviceIDs, labels)

	// append fluxmeter IDs
	fluxMeters := multiMatchResult.FluxMeters
	// TODO: change the return signature of GetMatches to return only FluxMeterIDs.
	fluxMeterIDs := make([]string, 0, len(fluxMeters))
	for _, fluxMeter := range fluxMeters {
		fluxMeterIDs = append(fluxMeterIDs, fluxMeter.GetMetricID())
	}

	limiterDecisions := []*flowcontrolv1.LimiterDecision{}

	var wg sync.WaitGroup
	var once sync.Once

	execLimiter := func(limiter iface.Limiter, decision *flowcontrolv1.LimiterDecision) func() {
		return func() {
			defer wg.Done()
			decision = limiter.RunLimiter(labels)
			if decision.Dropped {
				once.Do(setDecisionRejected)
			}
		}
	}

	execLimiters := func(limiters []iface.Limiter) []*flowcontrolv1.LimiterDecision {
		decisions := make([]*flowcontrolv1.LimiterDecision, len(limiters))
		// execute limiters
		for i, limiter := range limiters {
			wg.Add(1)
			decision := &flowcontrolv1.LimiterDecision{}
			decisions[i] = decision
			if i == len(limiters)-1 {
				execLimiter(limiter, decision)()
			} else {
				panichandler.Go(execLimiter(limiter, decision))
			}
		}
		wg.Wait()
		return decisions
	}

	returnResponse := func() *flowcontrolv1.CheckResponse {
		resp = &flowcontrolv1.CheckResponse{
			DecisionType:     decisionType,
			LimiterDecisions: limiterDecisions,
			FluxMeterIds:     fluxMeterIDs,
			Reason:           decisionReason,
		}
		return resp
	}

	///////////// Execute rate limiters and concurrency limiters //////////////

	// execute rate limiters first
	rateLimiters := multiMatchResult.RateLimiters
	rateLimiterDecisions := make([]*flowcontrolv1.LimiterDecision, 0, len(rateLimiters))
	// return extra tokens back to rate limiters in case some rate limiters or concurrent limiters decides to reject request
	defer func() {
		if resp.DecisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
			for i, l := range rateLimiterDecisions {
				if !l.Dropped && l.Reason == flowcontrolv1.LimiterDecision_LIMITER_REASON_UNSPECIFIED {
					go rateLimiters[i].TakeN(labels, -1)
				}
			}
		}
	}()

	rLimiters := []iface.Limiter{}
	for _, limiter := range rateLimiters {
		rLimiters = append(rLimiters, limiter)
	}

	limiterDecisions = append(limiterDecisions, execLimiters(rLimiters)...)
	// save rate limiter decisions separately in case we need to return back the tokens
	rateLimiterDecisions = limiterDecisions

	if decisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
		decisionReason = &flowcontrolv1.Reason{
			Reason: &flowcontrolv1.Reason_RejectReason_{
				RejectReason: flowcontrolv1.Reason_REJECT_REASON_RATE_LIMITED,
			},
		}
		return returnResponse()
	}

	// execute concurrency limiters
	concurrencyLimiters := multiMatchResult.ConcurrencyLimiters

	limiterDecisions = append(limiterDecisions, execLimiters(concurrencyLimiters)...)

	if decisionType == flowcontrolv1.DecisionType_DECISION_TYPE_REJECTED {
		decisionReason = &flowcontrolv1.Reason{
			Reason: &flowcontrolv1.Reason_RejectReason_{
				RejectReason: flowcontrolv1.Reason_REJECT_REASON_CONCURRENCY_LIMITED,
			},
		}
		return returnResponse()
	}

	return returnResponse()
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

	return e.register(
		"ConcurrencyLimiter:"+cl.GetPolicyName(),
		cl.GetSelector(),
		concurrencyLimiterMatchedCB,
	)
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

	return e.register(
		"RateLimiter:"+rl.GetPolicyName(),
		rl.GetSelector(),
		limiterActuatorMatchedCB,
	)
}

// UnregisterRateLimiter removes limiter actuator from multimatcher.
func (e *Engine) UnregisterRateLimiter(rl iface.RateLimiter) error {
	selectorProto := rl.GetSelector()
	return e.unregister("RateLimiter:"+rl.GetPolicyName(), selectorProto)
}

// GetMatches returns schedulers and fluxmeters for given labels.
func (e *Engine) getMatches(controlPoint selectors.ControlPoint, svcs []services.ServiceID, labels selectors.Labels) iface.MultiMatchResult {
	e.multiMatchMutex.RLock()
	defer e.multiMatchMutex.RUnlock()

	retMMRslt := iface.MultiMatchResult{}
	for _, service := range svcs {
		controlPointID := selectors.ControlPointID{
			ControlPoint: controlPoint,
			Service:      service,
		}

		// Lookup multi matcher for this control point id
		mm, ok := e.multiMatchers[controlPointID]
		if ok {
			// Run match
			resultCollection := mm.Match(multimatcher.Labels(labels.ToPlainMap()))
			// Append the matching Schedulers
			retMMRslt.ConcurrencyLimiters = append(retMMRslt.ConcurrencyLimiters, resultCollection.ConcurrencyLimiters...)
			// Append the matching FluxMeters
			retMMRslt.FluxMeters = append(retMMRslt.FluxMeters, resultCollection.FluxMeters...)
			// Append the matching Limiters
			retMMRslt.RateLimiters = append(retMMRslt.RateLimiters, resultCollection.RateLimiters...)
		}
	}
	return retMMRslt
}

func (e *Engine) register(
	key string,
	selectorProto *policylangv1.Selector,
	matchedCB multimatcher.MatchCallback[iface.MultiMatchResult],
) error {
	e.multiMatchMutex.Lock()
	defer e.multiMatchMutex.Unlock()
	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse selector")
		return fmt.Errorf("failed to parse selector: %v", err)
	}
	// check if multi matcher exists for this control point id
	mm, ok := e.multiMatchers[selector.ControlPointID]
	if !ok {
		mm = multimatcher.New[string, iface.MultiMatchResult]()
		e.multiMatchers[selector.ControlPointID] = mm
	}

	return mm.AddEntry(key, selector.LabelMatcher, matchedCB)
}

func (e *Engine) unregister(key string, selectorProto *policylangv1.Selector) error {
	e.multiMatchMutex.Lock()
	defer e.multiMatchMutex.Unlock()
	selector, err := selectors.FromProto(selectorProto)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse selector")
		return fmt.Errorf("failed to parse selector: %v", err)
	}
	// check if multi matcher exists for this control point id
	mm, ok := e.multiMatchers[selector.ControlPointID]
	if !ok {
		log.Warn().Msgf("Unable to unregister, multi matcher not found for control point id")
		return nil
	}
	retval := mm.RemoveEntry(key)
	// remove this multi matcher if this was the last entry
	if mm.Length() == 0 {
		delete(e.multiMatchers, selector.ControlPointID)
	}
	return retval
}
