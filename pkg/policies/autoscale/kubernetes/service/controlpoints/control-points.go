package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	controlpointcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	"github.com/fluxninja/aperture/pkg/discovery/kubernetes"
)

// Handler is the gRPC server handler.
type Handler struct {
	controlpointcachev1.UnimplementedControlPointsServiceServer
	AutoscaleControlPoints kubernetes.AutoscaleControlPoints
}

// NewHandler returns a new Handler.
func NewHandler(cpc kubernetes.AutoscaleControlPoints) *Handler {
	return &Handler{
		AutoscaleControlPoints: cpc,
	}
}

// GetControlPoints returns a ControlPoint from the cache.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*controlpointcachev1.AutoscaleKubernetesControlPoints, error) {
	keys := h.AutoscaleControlPoints.Keys()
	controlPoints := make([]*controlpointcachev1.AutoscaleKubernetesControlPoint, len(keys))
	for _, key := range keys {
		cp := key.ToProto()
		controlPoints = append(controlPoints, cp)
	}
	return &controlpointcachev1.AutoscaleKubernetesControlPoints{
		AutoscaleKubernetesControlPoints: controlPoints,
	}, nil
}
