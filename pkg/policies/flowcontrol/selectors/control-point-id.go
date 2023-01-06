package selectors

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// ControlPointID is the struct that represents a ControlPoint.
type ControlPointID struct {
	Service      string
	ControlPoint string
}

// NewControlPointID returns a controlPointID.
func NewControlPointID(service string, controlPoint string) ControlPointID {
	return ControlPointID{
		Service:      service,
		ControlPoint: controlPoint,
	}
}

func controlPointIDFromSelectorProto(flowSelectorMsg *policylangv1.FlowSelector) (ControlPointID, error) {
	ctrlPt := flowSelectorMsg.FlowMatcher.GetControlPoint()
	return ControlPointID{
		Service:      flowSelectorMsg.ServiceSelector.GetService(),
		ControlPoint: ctrlPt,
	}, nil
}
