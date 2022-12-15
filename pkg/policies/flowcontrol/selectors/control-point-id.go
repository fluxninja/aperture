package selectors

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// ControlPointID is the interface for controlPointID.
type ControlPointID interface {
	Service() string
	ControlPoint() string
}

type controlPointID struct {
	service      string
	controlPoint string
}

// controlPointID implements the ControlPointID interface.
var _ ControlPointID = (*controlPointID)(nil)

// NewControlPointID returns a controlPointID.
func NewControlPointID(service string, controlPoint string) ControlPointID {
	return controlPointID{
		service:      service,
		controlPoint: controlPoint,
	}
}

func controlPointIDFromSelectorProto(flowSelectorMsg *policylangv1.FlowSelector) (ControlPointID, error) {
	ctrlPt := flowSelectorMsg.FlowMatcher.GetControlPoint()
	return controlPointID{
		service:      flowSelectorMsg.ServiceSelector.GetService(),
		controlPoint: ctrlPt,
	}, nil
}

// Service returns the service name.
func (p controlPointID) Service() string {
	return p.service
}

// ControlPoint returns the control point.
func (p controlPointID) ControlPoint() string {
	return p.controlPoint
}
