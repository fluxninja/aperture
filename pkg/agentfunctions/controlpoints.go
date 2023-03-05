package agentfunctions

import (
	"context"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	flowcontrolcontrolpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// ControlPointsHandler is a handler for ListControlPoints function
//
// Note: There's no requirement every handler needs to be in a separate struct.
// More methods can be added to this one.
type ControlPointsHandler struct {
	cache      *cache.Cache[selectors.ControlPointID]
	agentGroup string
}

// NewControlPointsHandler returns a new ControlPointsHandler.
func NewControlPointsHandler(
	cache *cache.Cache[selectors.ControlPointID],
	agentInfo *agentinfo.AgentInfo,
) ControlPointsHandler {
	return ControlPointsHandler{
		cache:      cache,
		agentGroup: agentInfo.GetAgentGroup(),
	}
}

// ListControlPoints lists currently discovered control points.
func (h *ControlPointsHandler) ListControlPoints(
	ctx context.Context,
	_ *cmdv1.ListFlowControlControlPointsRequest,
) (*cmdv1.ListFlowControlControlPointsAgentResponse, error) {
	controlPoints := h.cache.GetAll()

	controlPointsProto := &flowcontrolcontrolpointsv1.FlowControlControlPoints{
		FlowControlControlPoints: make([]*flowcontrolcontrolpointsv1.FlowControlControlPoint, 0, len(controlPoints)),
	}
	for _, controlPoint := range controlPoints {
		controlPointsProto.FlowControlControlPoints = append(controlPointsProto.FlowControlControlPoints,
			&flowcontrolcontrolpointsv1.FlowControlControlPoint{
				Service:      controlPoint.Service,
				ControlPoint: controlPoint.ControlPoint,
			})
	}

	return &cmdv1.ListFlowControlControlPointsAgentResponse{
		FlowControlControlPoints: controlPointsProto,
		AgentGroup:               h.agentGroup,
	}, nil
}

// RegisterControlPointsHandler registers ControlPointsHandler in handler registry.
func RegisterControlPointsHandler(handler ControlPointsHandler, registry *rpc.HandlerRegistry) error {
	return rpc.RegisterFunction(registry, handler.ListControlPoints)
}
