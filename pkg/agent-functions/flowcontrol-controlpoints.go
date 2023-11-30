package agentfunctions

import (
	"context"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/cache"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	flowcontrolControlPoints "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service/controlpoints"
)

// FlowControlControlPointsHandler is a handler for ListFlowControlPoints function
//
// Note: There's no requirement every handler needs to be in a separate struct.
// More methods can be added to this one.
type FlowControlControlPointsHandler struct {
	cache      *cache.Cache[selectors.TypedControlPointID]
	agentGroup string
}

// NewFlowControlControlPointsHandler returns a new FlowControlControlPointsHandler.
func NewFlowControlControlPointsHandler(
	cache *cache.Cache[selectors.TypedControlPointID],
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
func RegisterControlPointsHandler(handler FlowControlControlPointsHandler, t *transport.EtcdTransportServer) error {
	return transport.RegisterFunction(t, handler.ListFlowControlPoints)
}
