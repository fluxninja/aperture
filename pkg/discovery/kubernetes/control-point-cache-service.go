package kubernetes

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	controlpointcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/controlpointcache/v1"
)

// Handler is the gRPC server handler.
type Handler struct {
	controlpointcachev1.UnimplementedControlPointCacheServer
	ControlPointCache ControlPointCache
}

// NewHandler returns a new Handler.
func NewHandler(cpc ControlPointCache) *Handler {
	return &Handler{
		ControlPointCache: cpc,
	}
}

// GetControlPoints returns a ControlPoint from the cache.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*controlpointcachev1.KubernetesControlPoints, error) {
	keys := h.ControlPointCache.Keys()
	controlPoints := make([]*controlpointcachev1.KubernetesControlPoint, len(keys))
	for _, key := range keys {
		cp := key.ToProto()
		controlPoints = append(controlPoints, cp)
	}
	return &controlpointcachev1.KubernetesControlPoints{
		ControlPoints: controlPoints,
	}, nil
}
