package controlpoints

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	flowcontrolcontrolpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

// Handler implements FlowControlControlPointsService.
type Handler struct {
	flowcontrolcontrolpointsv1.UnimplementedFlowControlControlPointsServiceServer
	serviceControlPointCache *cache.Cache[selectors.ControlPointID]
}

// NewHandler returns a new Handler.
func NewHandler(serviceControlPointCache *cache.Cache[selectors.ControlPointID]) *Handler {
	return &Handler{
		serviceControlPointCache: serviceControlPointCache,
	}
}

// GetControlPoints returns all control points.
func (h *Handler) GetControlPoints(ctx context.Context, _ *emptypb.Empty) (*flowcontrolcontrolpointsv1.FlowControlControlPoints, error) {
	serviceControlPointObjects := make(map[selectors.ControlPointID]struct{})
	if h.serviceControlPointCache != nil {
		serviceControlPointObjects = h.serviceControlPointCache.GetAll()
	}
	controlpoints := make([]*flowcontrolcontrolpointsv1.FlowControlControlPoint, 0, len(serviceControlPointObjects))
	for controlPointID := range serviceControlPointObjects {
		cp := &flowcontrolcontrolpointsv1.FlowControlControlPoint{
			Service:      controlPointID.Service,
			ControlPoint: controlPointID.ControlPoint,
		}
		controlpoints = append(controlpoints, cp)
	}
	return &flowcontrolcontrolpointsv1.FlowControlControlPoints{ControlPoints: controlpoints}, nil
}
