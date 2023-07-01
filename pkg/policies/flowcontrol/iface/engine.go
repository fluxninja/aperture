package iface

import (
	"context"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/labels"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

// RequestContext provides the request parameters for the Check method.
type RequestContext struct {
	FlowLabels   labels.Labels
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

	RegisterScheduler(ls Scheduler) error
	UnregisterScheduler(ls Scheduler) error
	GetScheduler(limiterID LimiterID) Scheduler

	RegisterFluxMeter(fm FluxMeter) error
	UnregisterFluxMeter(fm FluxMeter) error
	GetFluxMeter(fluxMeterName string) FluxMeter

	RegisterRateLimiter(l RateLimiter) error
	UnregisterRateLimiter(l RateLimiter) error
	GetRateLimiter(limiterID LimiterID) RateLimiter

	RegisterSampler(l Limiter) error
	UnregisterSampler(l Limiter) error
	GetSampler(limiterID LimiterID) Limiter

	RegisterLabelPreview(l LabelPreview) error
	UnregisterLabelPreview(l LabelPreview) error
}
