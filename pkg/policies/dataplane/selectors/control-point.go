package selectors

import (
	"fmt"

	selectorv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/selector/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
)

// ControlPoint is the interface for controlPoint.
type ControlPoint interface {
	Type() flowcontrolv1.ControlPoint_Type
	Feature() string
	ToFlowControlPointProto() *flowcontrolv1.ControlPoint
}

type controlPoint struct {
	feature string
	type_   flowcontrolv1.ControlPoint_Type
}

// controlPoint implements the ControlPoint interface.
var _ ControlPoint = (*controlPoint)(nil)

// NewControlPoint returns a controlPoint.
func NewControlPoint(type_ flowcontrolv1.ControlPoint_Type, feature string) ControlPoint {
	return controlPoint{
		type_:   type_,
		feature: feature,
	}
}

func controlPointFromSelectorControlPointProto(controlPointMsg *selectorv1.ControlPoint) (ControlPoint, error) {
	if controlPointMsg != nil && controlPointMsg.Controlpoint != nil {
		switch cp := controlPointMsg.Controlpoint.(type) {
		case *selectorv1.ControlPoint_Feature:
			return NewControlPoint(flowcontrolv1.ControlPoint_TYPE_FEATURE, cp.Feature), nil
		case *selectorv1.ControlPoint_Traffic:
			switch cp.Traffic {
			case "ingress":
				return NewControlPoint(flowcontrolv1.ControlPoint_TYPE_INGRESS, ""), nil
			case "egress":
				return NewControlPoint(flowcontrolv1.ControlPoint_TYPE_EGRESS, ""), nil
			default:
				return NewControlPoint(flowcontrolv1.ControlPoint_TYPE_UNKNOWN, ""), fmt.Errorf("invalid traffic direction")
			}
		}
	}
	return NewControlPoint(flowcontrolv1.ControlPoint_TYPE_UNKNOWN, ""), fmt.Errorf("unknown/missing control point")
}

// Type returns the control point type.
func (p controlPoint) Type() flowcontrolv1.ControlPoint_Type {
	return p.type_
}

// Feature returns the control point feature.
func (p controlPoint) Feature() string {
	return p.feature
}

// ToFlowControlPointProto returns a flow control control point proto.
func (p controlPoint) ToFlowControlPointProto() *flowcontrolv1.ControlPoint {
	return &flowcontrolv1.ControlPoint{
		Type:    p.type_,
		Feature: p.feature,
	}
}
