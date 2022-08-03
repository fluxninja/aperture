package iface

import (
	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/selectors"
	"github.com/fluxninja/aperture/pkg/services"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

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
	FluxMeters          []FluxMeter
	RateLimiters        []RateLimiter
}
