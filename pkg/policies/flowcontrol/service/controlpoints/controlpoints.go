package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	flowcontrolpointsv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	"github.com/fluxninja/aperture/v2/pkg/cache"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
)

// Handler implements FlowControlPointsService.
type Handler struct {
	flowcontrolpointsv1.UnimplementedFlowControlPointsServiceServer
	serviceControlPointCache *cache.Cache[selectors.TypedControlPointID]
}

// NewHandler returns a new Handler.
func NewHandler(serviceControlPointCache *cache.Cache[selectors.TypedControlPointID]) *Handler {
	return &Handler{
		serviceControlPointCache: serviceControlPointCache,
	}
}

// GetControlPoints returns all control points.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*flowcontrolpointsv1.FlowControlPoints, error) {
	return ToProto(h.serviceControlPointCache), nil
}

// ToProto converts cache to proto message.
func ToProto(cache *cache.Cache[selectors.TypedControlPointID]) *flowcontrolpointsv1.FlowControlPoints {
	cpObjects := cache.GetAll()
	fcp := &flowcontrolpointsv1.FlowControlPoints{
		FlowControlPoints: make([]*flowcontrolpointsv1.FlowControlPoint, 0, len(cpObjects)),
	}
	for _, controlPointID := range cpObjects {
		fcp.FlowControlPoints = append(fcp.FlowControlPoints, controlPointID.ToProto())
	}
	return fcp
}
