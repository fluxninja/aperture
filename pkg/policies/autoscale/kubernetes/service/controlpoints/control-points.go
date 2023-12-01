package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	controlpointsv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	cmdv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/cmd/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/etcd/transport"
	"github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/discovery"
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
	return h.AutoScaleControlPoints.ToProto(), nil
}

// RegisterControlPointsHandler registers ControlPointsHandler in RPC handler registry.
func RegisterControlPointsHandler(handler *Handler, t *transport.EtcdTransportServer) error {
	return transport.RegisterFunction(t, handler.ListAutoScaleControlPoints)
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
