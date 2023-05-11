package selectors

import (
	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	flowcontrolpointsv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
)

// ControlPointID is the struct that represents a ControlPoint.
//
// Agent group is implied.
// Type is ignored.
type ControlPointID struct {
	ControlPoint string
	Service      string
}

// TypedControlPointID is the struct that represents a FlowControlPoint.
//
// Agent group is implied.
//
// Note: We need to mirror flowcontrolpointsv1.FlowControlPoint, because
// protobuf-generated struct cannot be used as map keys.
type TypedControlPointID struct {
	ControlPointID
	Type string
}

// TypedGlobalControlPointID is ControlPointID with explicit agent group.
//
// Note: We need to mirror cmdv1.GlobalFlowControlPoint, because
// protobuf-generated struct cannot be used as map keys.
type TypedGlobalControlPointID struct {
	TypedControlPointID
	AgentGroup string
}

// GlobalControlPointID is just like TypedGlobalControlPointID but embedding the ControlPointID instead of TypedControlPointID.
//
// Useful for defining a control point to find, without having to specify the source.
type GlobalControlPointID struct {
	ControlPointID
	AgentGroup string
}

// NewControlPointID returns a controlPointID.
func NewControlPointID(controlPoint, service string) ControlPointID {
	return ControlPointID{
		ControlPoint: controlPoint,
		Service:      service,
	}
}

// WithType returns the controlpoint as TypedControlPointID.
func (cp ControlPointID) WithType(controlPointType string) TypedControlPointID {
	return TypedControlPointID{
		ControlPointID: cp,
		Type:           controlPointType,
	}
}

// String returns the string representation of the control point.
func (cp ControlPointID) String() string {
	return cp.Service + "/" + cp.ControlPoint
}

// NewTypedControlPointID returns a typedControlPointID.
func NewTypedControlPointID(controlPoint, controlPointType, service string) TypedControlPointID {
	return TypedControlPointID{
		ControlPointID: NewControlPointID(controlPoint, service),
		Type:           controlPointType,
	}
}

// ToProto returns protobuf representation of control point.
func (cp *TypedControlPointID) ToProto() *flowcontrolpointsv1.FlowControlPoint {
	return &flowcontrolpointsv1.FlowControlPoint{
		ControlPoint: cp.ControlPoint,
		Type:         cp.Type,
		Service:      cp.Service,
	}
}

// InAgentGroup returns the controlpoint as TypedGlobalControlPointID with given agent group.
func (cp TypedControlPointID) InAgentGroup(agentGroup string) TypedGlobalControlPointID {
	return TypedGlobalControlPointID{
		TypedControlPointID: cp,
		AgentGroup:          agentGroup,
	}
}

// TypedControlPointIDFromProto creates TypedControlPointID from protobuf representation.
func TypedControlPointIDFromProto(protoCP *flowcontrolpointsv1.FlowControlPoint) TypedControlPointID {
	return TypedControlPointID{
		ControlPointID: ControlPointID{
			Service:      protoCP.GetService(),
			ControlPoint: protoCP.GetControlPoint(),
		},
		Type: protoCP.GetType(),
	}
}

// ToProto returns protobuf representation of control point.
func (cp *TypedGlobalControlPointID) ToProto() *cmdv1.GlobalFlowControlPoint {
	return &cmdv1.GlobalFlowControlPoint{
		FlowControlPoint: &flowcontrolpointsv1.FlowControlPoint{
			ControlPoint: cp.ControlPoint,
			Type:         cp.Type,
			Service:      cp.Service,
		},
		AgentGroup: cp.AgentGroup,
	}
}

// TypedGlobalControlPointIDFromProto creates ControlPointID from protobuf representation.
func TypedGlobalControlPointIDFromProto(protoCP *cmdv1.GlobalFlowControlPoint) TypedGlobalControlPointID {
	return TypedGlobalControlPointID{
		TypedControlPointID: TypedControlPointID{
			ControlPointID: ControlPointID{
				ControlPoint: protoCP.FlowControlPoint.GetControlPoint(),
				Service:      protoCP.FlowControlPoint.GetService(),
			},
			Type: protoCP.FlowControlPoint.GetType(),
		},
		AgentGroup: protoCP.GetAgentGroup(),
	}
}
