package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	controlpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	"github.com/fluxninja/aperture/pkg/discovery/kubernetes"
)

// Handler is the gRPC server handler.
type Handler struct {
	controlpointsv1.UnimplementedAutoscaleKubernetesControlPointsServiceServer
	AutoscaleControlPoints kubernetes.AutoscaleControlPoints
}

// NewHandler returns a new Handler.
func NewHandler(cpc kubernetes.AutoscaleControlPoints) *Handler {
	return &Handler{
		AutoscaleControlPoints: cpc,
	}
}

// GetControlPoints returns a ControlPoint from the cache.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*controlpointsv1.AutoscaleKubernetesControlPoints, error) {
	keys := h.AutoscaleControlPoints.Keys()
	controlPoints := make([]*controlpointsv1.AutoscaleKubernetesControlPoint, len(keys))
	for _, key := range keys {
		cp := key.ToProto()
		controlPoints = append(controlPoints, cp)
	}
	return &controlpointsv1.AutoscaleKubernetesControlPoints{
		AutoscaleKubernetesControlPoints: controlPoints,
	}, nil
}
