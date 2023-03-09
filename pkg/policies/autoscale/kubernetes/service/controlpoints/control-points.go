package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	controlpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/discovery"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// Handler is the gRPC server handler.
type Handler struct {
	controlpointsv1.UnimplementedAutoScaleKubernetesControlPointsServiceServer
	AutoScaleControlPoints discovery.AutoScaleControlPoints
	AgentGroup             string
}

// NewHandler returns a new Handler.
func NewHandler(cpc discovery.AutoScaleControlPoints, agentInfo *agentinfo.AgentInfo) *Handler {
	return &Handler{
		AutoScaleControlPoints: cpc,
		AgentGroup:             agentInfo.GetAgentGroup(),
	}
}

// GetControlPoints returns a ControlPoint from the cache.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*controlpointsv1.AutoScaleKubernetesControlPoints, error) {
	return &controlpointsv1.AutoScaleKubernetesControlPoints{
		AutoScaleKubernetesControlPoints: h.AutoScaleControlPoints.ToProto(),
	}, nil
}

// RegisterControlPointsHandler registers ControlPointsHandler in RPC handler registry.
func RegisterControlPointsHandler(handler *Handler, registry *rpc.HandlerRegistry) error {
	return rpc.RegisterFunction(registry, handler.ListAutoScaleControlPoints)
}

// ListAutoScaleControlPoints lists currently discovered control points.
func (h *Handler) ListAutoScaleControlPoints(
	ctx context.Context,
	_ *cmdv1.ListAutoScaleControlPointsRequest,
) (*cmdv1.ListAutoScaleControlPointsAgentResponse, error) {
	return &cmdv1.ListAutoScaleControlPointsAgentResponse{
		AutoScaleControlPoints: h.AutoScaleControlPoints.ToProto(),
		AgentGroup:             h.AgentGroup,
	}, nil
}
