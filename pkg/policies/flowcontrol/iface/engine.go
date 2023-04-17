package iface

import (
	"context"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

// RequestContext provides the request parameters for the Check method.
type RequestContext struct {
	FlowLabels   map[string]string
	ControlPoint string
	Services     []string
	Tokens       uint64
}

// Engine is an interface for registering fluxmeters and schedulers.
type Engine interface {
	ProcessRequest(
		ctx context.Context,
		requestContext RequestContext,
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

	RegisterFlowRegulator(l Limiter) error
	UnregisterFlowRegulator(l Limiter) error
	GetFlowRegulator(limiterID LimiterID) Limiter

	RegisterLabelPreview(l LabelPreview) error
	UnregisterLabelPreview(l LabelPreview) error
}
