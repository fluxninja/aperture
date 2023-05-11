package iface

import (
	"context"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/v2/pkg/agentinfo"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

// RequestContext provides the request parameters for the Check method.
type RequestContext struct {
	FlowLabels   map[string]string
	ControlPoint string
	Services     []string
}

// Engine is an interface for registering fluxmeters and schedulers.
type Engine interface {
	ProcessRequest(
		ctx context.Context,
		requestContext RequestContext,
	) *flowcontrolv1.CheckResponse

	GetAgentInfo() *agentinfo.AgentInfo

	RegisterLoadScheduler(ls LoadScheduler) error
	UnregisterLoadScheduler(ls LoadScheduler) error
	GetLoadScheduler(limiterID LimiterID) LoadScheduler

	RegisterFluxMeter(fm FluxMeter) error
	UnregisterFluxMeter(fm FluxMeter) error
	GetFluxMeter(fluxMeterName string) FluxMeter

	RegisterRateLimiter(l RateLimiter) error
	UnregisterRateLimiter(l RateLimiter) error
	GetRateLimiter(limiterID LimiterID) RateLimiter

	RegisterRegulator(l Limiter) error
	UnregisterRegulator(l Limiter) error
	GetRegulator(limiterID LimiterID) Limiter

	RegisterLabelPreview(l LabelPreview) error
	UnregisterLabelPreview(l LabelPreview) error
}
