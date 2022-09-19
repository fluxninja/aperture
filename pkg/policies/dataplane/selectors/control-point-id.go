package selectors

import (
	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
)

// ControlPointID is the interface for controlPointID.
type ControlPointID interface {
	Service() string
	ControlPoint() ControlPoint
}

type controlPointID struct {
	service      string
	controlPoint ControlPoint
}

// controlPointID implements the ControlPointID interface.
var _ ControlPointID = (*controlPointID)(nil)

// NewControlPointID returns a controlPointID.
func NewControlPointID(service string, controlPoint ControlPoint) ControlPointID {
	return controlPointID{
		service:      service,
		controlPoint: controlPoint,
	}
}

func controlPointIDFromSelectorProto(selectorMsg *selectorv1.Selector) (ControlPointID, error) {
	ctrlPt, err := controlPointFromSelectorControlPointProto(selectorMsg.ControlPoint)
	return controlPointID{
		service:      selectorMsg.Service,
		controlPoint: ctrlPt,
	}, err
}

// Service returns the service name.
func (p controlPointID) Service() string {
	return p.service
}

// ControlPoint returns the control point.
func (p controlPointID) ControlPoint() ControlPoint {
	return p.controlPoint
}
