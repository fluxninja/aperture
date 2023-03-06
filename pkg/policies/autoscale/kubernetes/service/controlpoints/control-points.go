package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	controlpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	"github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/discovery"
)

// Handler is the gRPC server handler.
type Handler struct {
	controlpointsv1.UnimplementedAutoScaleKubernetesControlPointsServiceServer
	AutoScaleControlPoints discovery.AutoScaleControlPoints
}

// NewHandler returns a new Handler.
func NewHandler(cpc discovery.AutoScaleControlPoints) *Handler {
	return &Handler{
		AutoScaleControlPoints: cpc,
	}
}

// GetControlPoints returns a ControlPoint from the cache.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*controlpointsv1.AutoScaleKubernetesControlPoints, error) {
	return h.AutoScaleControlPoints.ToProto(), nil
}
