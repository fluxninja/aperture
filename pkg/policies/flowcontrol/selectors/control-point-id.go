package selectors

import (
	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	flowcontrolcontrolpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
)

// ControlPointID is the struct that represents a ControlPoint.
//
// Agent group is implied.
//
// Note: We need to mirror cmdv1.ServiceControlPoint, because
// protobuf-generated struct cannot be used as map keys.
type ControlPointID struct {
	Service      string
	ControlPoint string
}

// GlobalControlPointID is ControlPointID with explicit agent group.
//
// Note: We need to mirror cmdv1.ServiceControlPoint, because
// protobuf-generated struct cannot be used as map keys.
type GlobalControlPointID struct {
	ControlPointID
	AgentGroup string
}

// NewControlPointID returns a controlPointID.
func NewControlPointID(service string, controlPoint string) ControlPointID {
	return ControlPointID{
		Service:      service,
		ControlPoint: controlPoint,
	}
}

// ToProto returns protobuf representation of control point.
func (cp *ControlPointID) ToProto() *flowcontrolcontrolpointsv1.FlowControlControlPoint {
	return &flowcontrolcontrolpointsv1.FlowControlControlPoint{
		Service:      cp.Service,
		ControlPoint: cp.ControlPoint,
	}
}

// InAgentGroup returns the controlpoint as GlobalControlPointID with given agent group.
func (cp ControlPointID) InAgentGroup(agentGroup string) GlobalControlPointID {
	return GlobalControlPointID{
		ControlPointID: cp,
		AgentGroup:     agentGroup,
	}
}

// ControlPointIDFromProto creates ControlPointID from protobuf representation.
func ControlPointIDFromProto(protoCP *flowcontrolcontrolpointsv1.FlowControlControlPoint) ControlPointID {
	return ControlPointID{
		Service:      protoCP.GetService(),
		ControlPoint: protoCP.GetControlPoint(),
	}
}

// ToProto returns protobuf representation of control point.
func (cp *GlobalControlPointID) ToProto() *cmdv1.GlobalFlowControlControlPoint {
	return &cmdv1.GlobalFlowControlControlPoint{
		FlowControlControlPoint: &flowcontrolcontrolpointsv1.FlowControlControlPoint{
			Service:      cp.Service,
			ControlPoint: cp.ControlPoint,
		},
		AgentGroup: cp.AgentGroup,
	}
}

// GlobalControlPointIDFromProto creates ControlPointID from protobuf representation.
func GlobalControlPointIDFromProto(protoCP *cmdv1.GlobalFlowControlControlPoint) GlobalControlPointID {
	return GlobalControlPointID{
		ControlPointID: ControlPointID{
			Service:      protoCP.FlowControlControlPoint.GetService(),
			ControlPoint: protoCP.FlowControlControlPoint.GetControlPoint(),
		},
		AgentGroup: protoCP.GetAgentGroup(),
	}
}

func controlPointIDFromSelectorProto(flowSelectorMsg *policylangv1.FlowSelector) (ControlPointID, error) {
	ctrlPt := flowSelectorMsg.FlowMatcher.GetControlPoint()
	return ControlPointID{
		Service:      flowSelectorMsg.ServiceSelector.GetService(),
		ControlPoint: ctrlPt,
	}, nil
}
