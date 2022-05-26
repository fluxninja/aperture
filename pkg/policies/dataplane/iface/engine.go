package iface

import (
	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "aperture.tech/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"aperture.tech/aperture/pkg/selectors"
	"aperture.tech/aperture/pkg/services"
)

// EngineAPI is an interface for registering fluxmeters and schedulers.
type EngineAPI interface {
	ProcessRequest(controlPoint selectors.ControlPoint, serviceIDs []services.ServiceID, labels selectors.Labels) *flowcontrolv1.CheckResponse

	RegisterConcurrencyLimiter(sa Limiter) error
	UnregisterConcurrencyLimiter(sa Limiter) error

	RegisterFluxMeter(fm FluxMeter) error
	UnregisterFluxMeter(fm FluxMeter) error
	GetFluxMeterHist(metricID string) prometheus.Histogram

	RegisterRateLimiter(l RateLimiter) error
	UnregisterRateLimiter(l RateLimiter) error
}

// MultiMatchResult is used as return value of PolicyConfigAPI.GetMatches.
type MultiMatchResult struct {
	ConcurrencyLimiters []Limiter
	// TODO: Can be FluxMeterIDs
	FluxMeters   []FluxMeter
	RateLimiters []RateLimiter
}
