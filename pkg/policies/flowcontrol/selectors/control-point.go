package selectors

import (
	"fmt"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// ControlPoint is the interface for controlPoint.
type ControlPoint interface {
	Type() flowcontrolv1.ControlPointInfo_Type
	Feature() string
	ToControlPointInfoProto() *flowcontrolv1.ControlPointInfo
}

type controlPoint struct {
	feature string
	type_   flowcontrolv1.ControlPointInfo_Type
}

// controlPoint implements the ControlPoint interface.
var _ ControlPoint = (*controlPoint)(nil)

// NewControlPoint returns a controlPoint.
func NewControlPoint(type_ flowcontrolv1.ControlPointInfo_Type, feature string) ControlPoint {
	return controlPoint{
		type_:   type_,
		feature: feature,
	}
}

func controlPointFromSelectorControlPointProto(controlPointMsg *policylangv1.ControlPoint) (ControlPoint, error) {
	if controlPointMsg != nil && controlPointMsg.Controlpoint != nil {
		switch cp := controlPointMsg.Controlpoint.(type) {
		case *policylangv1.ControlPoint_Feature:
			return NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_FEATURE, cp.Feature), nil
		case *policylangv1.ControlPoint_Traffic:
			switch cp.Traffic {
			case "ingress":
				return NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_INGRESS, ""), nil
			case "egress":
				return NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_EGRESS, ""), nil
			default:
				return NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_UNKNOWN, ""), fmt.Errorf("invalid traffic direction")
			}
		}
	}
	return NewControlPoint(flowcontrolv1.ControlPointInfo_TYPE_UNKNOWN, ""), fmt.Errorf("unknown/missing control point")
}

// Type returns the control point type.
func (p controlPoint) Type() flowcontrolv1.ControlPointInfo_Type {
	return p.type_
}

// Feature returns the control point feature.
func (p controlPoint) Feature() string {
	return p.feature
}

// ToControlPointInfoProto returns a flow control control point proto.
func (p controlPoint) ToControlPointInfoProto() *flowcontrolv1.ControlPointInfo {
	return &flowcontrolv1.ControlPointInfo{
		Type:    p.type_,
		Feature: p.feature,
	}
}
