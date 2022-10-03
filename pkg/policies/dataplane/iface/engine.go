package iface

import (
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/multimatcher"
	"github.com/fluxninja/aperture/pkg/policies/dataplane/selectors"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

// Engine is an interface for registering fluxmeters and schedulers.
type Engine interface {
	ProcessRequest(controlPoint selectors.ControlPoint, serviceIDs []string, labels map[string]string) *flowcontrolv1.CheckResponse

	RegisterConcurrencyLimiter(sa ConcurrencyLimiter) error
	UnregisterConcurrencyLimiter(sa ConcurrencyLimiter) error
	GetConcurrencyLimiter(limiterID LimiterID) ConcurrencyLimiter

	RegisterFluxMeter(fm FluxMeter) error
	UnregisterFluxMeter(fm FluxMeter) error
	GetFluxMeter(fluxMeterName string) FluxMeter

	RegisterRateLimiter(l RateLimiter) error
	UnregisterRateLimiter(l RateLimiter) error
	GetRateLimiter(limiterID LimiterID) RateLimiter
}

// MultiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type MultiMatchResult struct {
	ConcurrencyLimiters []ConcurrencyLimiter
	FluxMeters          []FluxMeter
	RateLimiters        []RateLimiter
	Classifiers         []Classifier
}

// PopulateFromMultiMatcher populates result object with results from MultiMatcher.
func (result *MultiMatchResult) PopulateFromMultiMatcher(mm *multimatcher.MultiMatcher[string, MultiMatchResult], labels map[string]string) {
	resultCollection := mm.Match(multimatcher.Labels(labels))
	result.ConcurrencyLimiters = append(result.ConcurrencyLimiters, resultCollection.ConcurrencyLimiters...)
	result.FluxMeters = append(result.FluxMeters, resultCollection.FluxMeters...)
	result.RateLimiters = append(result.RateLimiters, resultCollection.RateLimiters...)
	result.Classifiers = append(result.Classifiers, resultCollection.Classifiers...)
}
