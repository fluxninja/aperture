package iface

import (
	"context"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

// Engine is an interface for registering fluxmeters and schedulers.
type Engine interface {
	ProcessRequest(
		ctx context.Context,
		controlPoint selectors.ControlPoint,
		serviceIDs []string,
		labels map[string]string,
	) *flowcontrolv1.CheckResponse

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
