package selectors

import (
	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	flowcontrolpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/controlpoints/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/consts"
)

// ControlPointID is the struct that represents a ControlPoint.
//
// Agent group is implied.
// Type is ignored.
type ControlPointID struct {
	Service      string
	ControlPoint string
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
func NewControlPointID(service string, controlPoint string) ControlPointID {
	return ControlPointID{
		Service:      service,
		ControlPoint: controlPoint,
	}
}

// WithType returns the controlpoint as TypedControlPointID.
func (cp ControlPointID) WithType(controlPointType string) TypedControlPointID {
	return TypedControlPointID{
		ControlPointID: cp,
		Type:           controlPointType,
	}
}

// NewTypedControlPointID returns a typedControlPointID.
func NewTypedControlPointID(service string, controlPoint string, controlPointType string) TypedControlPointID {
	return TypedControlPointID{
		ControlPointID: NewControlPointID(service, controlPoint),
		Type:           controlPointType,
	}
}

// ToProto returns protobuf representation of control point.
func (cp *TypedControlPointID) ToProto() *flowcontrolpointsv1.FlowControlPoint {
	return &flowcontrolpointsv1.FlowControlPoint{
		Service:      cp.Service,
		ControlPoint: cp.ControlPoint,
		Type:         cp.Type,
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
			Service:      cp.Service,
			ControlPoint: cp.ControlPoint,
			Type:         cp.Type,
		},
		AgentGroup: cp.AgentGroup,
	}
}

// TypedGlobalControlPointIDFromProto creates ControlPointID from protobuf representation.
func TypedGlobalControlPointIDFromProto(protoCP *cmdv1.GlobalFlowControlPoint) TypedGlobalControlPointID {
	return TypedGlobalControlPointID{
		TypedControlPointID: TypedControlPointID{
			ControlPointID: ControlPointID{
				Service:      protoCP.FlowControlPoint.GetService(),
				ControlPoint: protoCP.FlowControlPoint.GetControlPoint(),
			},
			Type: protoCP.FlowControlPoint.GetType(),
		},
		AgentGroup: protoCP.GetAgentGroup(),
	}
}

func controlPointIDFromSelectorProto(flowSelectorMsg *policylangv1.FlowSelector) (ControlPointID, error) {
	ctrlPt := flowSelectorMsg.FlowMatcher.GetControlPoint()
	service := flowSelectorMsg.ServiceSelector.GetService()
	// map all to catch-all service for backward compatibility
	// Deprecated: v1.5.0
	if service == "all" {
		service = consts.AnyService
	}
	return ControlPointID{
		Service:      service,
		ControlPoint: ctrlPt,
	}, nil
}
