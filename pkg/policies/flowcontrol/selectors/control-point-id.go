package selectors

import (
	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// ControlPointID is the struct that represents a ControlPoint.
//
// Note: We need to mirror cmdv1.ServiceControlPoint, because
// protobuf-generated struct cannot be used as map keys.
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

// ToProto returns protobuf representation of control point.
func (cp *ControlPointID) ToProto() *cmdv1.ServiceControlPoint {
	return &cmdv1.ServiceControlPoint{
		ServiceName: cp.Service,
		Name:        cp.ControlPoint,
	}
}

// ControlPointIDFromProto creates ControlPointID from protobuf representation.
func ControlPointIDFromProto(protoCP *cmdv1.ServiceControlPoint) ControlPointID {
	return ControlPointID{
		Service:      protoCP.GetServiceName(),
		ControlPoint: protoCP.GetName(),
	}
}

func controlPointIDFromSelectorProto(flowSelectorMsg *policylangv1.FlowSelector) (ControlPointID, error) {
	ctrlPt := flowSelectorMsg.FlowMatcher.GetControlPoint()
	return ControlPointID{
		Service:      flowSelectorMsg.ServiceSelector.GetService(),
		ControlPoint: ctrlPt,
	}, nil
}
