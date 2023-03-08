package agentfunctions

import (
	"context"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	flowcontrolControlPoints "github.com/fluxninja/aperture/pkg/policies/flowcontrol/service/controlpoints"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// FlowControlControlPointsHandler is a handler for ListFlowControlPoints function
//
// Note: There's no requirement every handler needs to be in a separate struct.
// More methods can be added to this one.
type FlowControlControlPointsHandler struct {
	cache      *cache.Cache[selectors.ControlPointID]
	agentGroup string
}

// NewFlowControlControlPointsHandler returns a new FlowControlControlPointsHandler.
func NewFlowControlControlPointsHandler(
	cache *cache.Cache[selectors.ControlPointID],
	agentInfo *agentinfo.AgentInfo,
) FlowControlControlPointsHandler {
	return FlowControlControlPointsHandler{
		cache:      cache,
		agentGroup: agentInfo.GetAgentGroup(),
	}
}

// ListFlowControlPoints lists currently discovered control points.
func (h *FlowControlControlPointsHandler) ListFlowControlPoints(
	ctx context.Context,
	_ *cmdv1.ListFlowControlPointsRequest,
) (*cmdv1.ListFlowControlPointsAgentResponse, error) {
	return &cmdv1.ListFlowControlPointsAgentResponse{
		FlowControlPoints: flowcontrolControlPoints.ToProto(h.cache),
		AgentGroup:        h.agentGroup,
	}, nil
}

// RegisterControlPointsHandler registers ControlPointsHandler in handler registry.
func RegisterControlPointsHandler(handler FlowControlControlPointsHandler, registry *rpc.HandlerRegistry) error {
	return rpc.RegisterFunction(registry, handler.ListFlowControlPoints)
}
