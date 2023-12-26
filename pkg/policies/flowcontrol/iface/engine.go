package iface

import (
	"context"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/labels"
)

//go:generate mockgen -source=engine.go -destination=../../mocks/mock_engine.go -package=mocks

// RequestContext provides the request parameters for the Check method.
type RequestContext struct {
	FlowLabels         labels.Labels
	ControlPoint       string
	CacheLookupRequest *flowcontrolv1.CacheLookupRequest
	Services           []string
	RampMode           bool
	ExpectEnd          bool
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

	RegisterRateLimiter(l Limiter) error
	UnregisterRateLimiter(l Limiter) error
	GetRateLimiter(limiterID LimiterID) Limiter

	RegisterSampler(l Limiter) error
	UnregisterSampler(l Limiter) error
	GetSampler(limiterID LimiterID) Limiter

	RegisterLabelPreview(l LabelPreview) error
	UnregisterLabelPreview(l LabelPreview) error

	RegisterCache(c Cache)

	RegisterConcurrencyLimiter(l ConcurrencyLimiter) error
	UnregisterConcurrencyLimiter(l ConcurrencyLimiter) error
	// Note: Use GetRateLimiter and GetFlowEnder for retrieving the limiter and flowender respectively.

	RegisterConcurrencyScheduler(l ConcurrencyScheduler) error
	UnregisterConcurrencyScheduler(l ConcurrencyScheduler) error
	// Note: Use GetScheduler and GetFlowEnder for retrieving the scheduler and flowender respectively.

	GetFlowEnder(limiterID LimiterID) FlowEnder

	FlowEnd(
		ctx context.Context,
		request *flowcontrolv1.FlowEndRequest,
	) *flowcontrolv1.FlowEndResponse
}
