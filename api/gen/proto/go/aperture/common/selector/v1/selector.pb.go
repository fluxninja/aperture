// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/common/selector/v1/selector.proto

package selectorv1

import (
	v1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/labelmatcher/v1"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Describes which flows a [dataplane
// component](/concepts/flow-control/flow-control.md#components) should apply
// to
//
// :::info
// See also [Selector overview](/concepts/flow-control/selector/selector.md).
// :::
//
// Example:
// ```yaml
// service: service1.default.svc.cluster.local
// control_point:
//   traffic: ingress # Allowed values are `ingress` and `egress`.
// label_matcher:
//   match_labels:
//     user_tier: gold
//   match_expressions:
//     - key: query
//       operator: In
//       values:
//         - insert
//         - delete
//     - label: user_agent
//       regex: ^(?!.*Chrome).*Safari
// ```
type Selector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Which [agent-group](/concepts/flow-control/selector/service.md#agent-group) this
	// selector applies to.
	AgentGroup string `protobuf:"bytes,1,opt,name=agent_group,json=agentGroup,proto3" json:"agent_group,omitempty" default:"default"` // @gotags: default:"default"
	// The Fully Qualified Domain Name of the
	// [service](/concepts/flow-control/selector/service.md) to select.
	//
	// In kubernetes, this is the FQDN of the Service object.
	//
	// Empty string means all services within an agent group (catch-all).
	//
	// :::note
	// One entity may belong to multiple services.
	// :::
	Service string `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	// Describes
	// [control point](/concepts/flow-control/flow-control.md#control-point)
	// within the entity where the policy should apply to.
	ControlPoint *ControlPoint `protobuf:"bytes,3,opt,name=control_point,json=controlPoint,proto3" json:"control_point,omitempty" validate:"required"` // @gotags: validate:"required"
	// Label matcher allows to add _additional_ condition on
	// [flow labels](/concepts/flow-control/selector/flow-label.md)
	// must also be satisfied (in addition to service+control point matching)
	//
	// :::info
	// See also [Label Matcher overview](/concepts/flow-control/selector/selector.md#label-matcher).
	// :::
	//
	// :::note
	// [Classifiers](#v1-classifier) _can_ use flow labels created by some other
	// classifier, but only if they were created at some previous control point
	// (and propagated in baggage).
	//
	// This limitation doesn't apply to selectors of other entities, like
	// FluxMeters or actuators. It's valid to create a flow label on a control
	// point using classifier, and immediately use it for matching on the same
	// control point.
	// :::
	LabelMatcher *v1.LabelMatcher `protobuf:"bytes,4,opt,name=label_matcher,json=labelMatcher,proto3" json:"label_matcher,omitempty"`
}

func (x *Selector) Reset() {
	*x = Selector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_selector_v1_selector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Selector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Selector) ProtoMessage() {}

func (x *Selector) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_selector_v1_selector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Selector.ProtoReflect.Descriptor instead.
func (*Selector) Descriptor() ([]byte, []int) {
	return file_aperture_common_selector_v1_selector_proto_rawDescGZIP(), []int{0}
}

func (x *Selector) GetAgentGroup() string {
	if x != nil {
		return x.AgentGroup
	}
	return ""
}

func (x *Selector) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *Selector) GetControlPoint() *ControlPoint {
	if x != nil {
		return x.ControlPoint
	}
	return nil
}

func (x *Selector) GetLabelMatcher() *v1.LabelMatcher {
	if x != nil {
		return x.LabelMatcher
	}
	return nil
}

// Identifies control point within a service that the rule or policy should apply to.
// Controlpoint is either a library feature name or one of ingress/egress traffic control point.
type ControlPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gotags: validate:"required"
	//
	// Types that are assignable to Controlpoint:
	//	*ControlPoint_Feature
	//	*ControlPoint_Traffic
	Controlpoint isControlPoint_Controlpoint `protobuf_oneof:"controlpoint" validate:"required"`
}

func (x *ControlPoint) Reset() {
	*x = ControlPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_selector_v1_selector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ControlPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ControlPoint) ProtoMessage() {}

func (x *ControlPoint) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_selector_v1_selector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ControlPoint.ProtoReflect.Descriptor instead.
func (*ControlPoint) Descriptor() ([]byte, []int) {
	return file_aperture_common_selector_v1_selector_proto_rawDescGZIP(), []int{1}
}

func (m *ControlPoint) GetControlpoint() isControlPoint_Controlpoint {
	if m != nil {
		return m.Controlpoint
	}
	return nil
}

func (x *ControlPoint) GetFeature() string {
	if x, ok := x.GetControlpoint().(*ControlPoint_Feature); ok {
		return x.Feature
	}
	return ""
}

func (x *ControlPoint) GetTraffic() string {
	if x, ok := x.GetControlpoint().(*ControlPoint_Traffic); ok {
		return x.Traffic
	}
	return ""
}

type isControlPoint_Controlpoint interface {
	isControlPoint_Controlpoint()
}

type ControlPoint_Feature struct {
	// Name of Aperture SDK's feature.
	// Feature corresponds to a block of code that can be "switched off" which usually is a "named opentelemetry's Span".
	//
	// Note: Flowcontrol only.
	Feature string `protobuf:"bytes,1,opt,name=feature,proto3,oneof" validate:"required"` //@gotags: validate:"required"
}

type ControlPoint_Traffic struct {
	// Type of traffic service, either "ingress" or "egress".
	// Apply the policy to the whole incoming/outgoing traffic of a service.
	// Usually powered by integration with a proxy (like envoy) or a web framework.
	//
	// * Flowcontrol: Blockable atom here is a single HTTP-transaction.
	// * Classification: Apply the classification rules to every incoming/outgoing request and attach the resulting flow labels to baggage and telemetry.
	Traffic string `protobuf:"bytes,2,opt,name=traffic,proto3,oneof" validate:"required,oneof=ingress egress"` // @gotags: validate:"required,oneof=ingress egress"
}

func (*ControlPoint_Feature) isControlPoint_Controlpoint() {}

func (*ControlPoint_Traffic) isControlPoint_Controlpoint() {}

var File_aperture_common_selector_v1_selector_proto protoreflect.FileDescriptor

var file_aperture_common_selector_v1_selector_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x32, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70,
	0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad, 0x02,
	0x0a, 0x08, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x40, 0x0a, 0x0b, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x1f, 0x92, 0x41, 0x1c, 0x82, 0x03, 0x19, 0x0a, 0x0c, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x09, 0x1a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x52, 0x0a, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x71, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x21, 0x92, 0x41, 0x1e, 0x82, 0x03, 0x1b,
	0x0a, 0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x0a, 0x1a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x0c, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x52, 0x0a, 0x0d, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52,
	0x0c, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x22, 0xb1, 0x01,
	0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x3d,
	0x0a, 0x07, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x21, 0x92, 0x41, 0x1e, 0x82, 0x03, 0x1b, 0x0a, 0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x1a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72,
	0x65, 0x64, 0x48, 0x00, 0x52, 0x07, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x52, 0x0a,
	0x07, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36,
	0x92, 0x41, 0x33, 0x82, 0x03, 0x30, 0x0a, 0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x1a, 0x1d, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x2c, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x3d, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x73, 0x73, 0x20,
	0x65, 0x67, 0x72, 0x65, 0x73, 0x73, 0x48, 0x00, 0x52, 0x07, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69,
	0x63, 0x42, 0x0e, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x42, 0x96, 0x02, 0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x55, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f,
	0x76, 0x31, 0x3b, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x41, 0x43, 0x53, 0xaa, 0x02, 0x1b, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x1b, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x5c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5c, 0x56, 0x31, 0xe2,
	0x02, 0x27, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x53, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_aperture_common_selector_v1_selector_proto_rawDescOnce sync.Once
	file_aperture_common_selector_v1_selector_proto_rawDescData = file_aperture_common_selector_v1_selector_proto_rawDesc
)

func file_aperture_common_selector_v1_selector_proto_rawDescGZIP() []byte {
	file_aperture_common_selector_v1_selector_proto_rawDescOnce.Do(func() {
		file_aperture_common_selector_v1_selector_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_common_selector_v1_selector_proto_rawDescData)
	})
	return file_aperture_common_selector_v1_selector_proto_rawDescData
}

var file_aperture_common_selector_v1_selector_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_aperture_common_selector_v1_selector_proto_goTypes = []interface{}{
	(*Selector)(nil),        // 0: aperture.common.selector.v1.Selector
	(*ControlPoint)(nil),    // 1: aperture.common.selector.v1.ControlPoint
	(*v1.LabelMatcher)(nil), // 2: aperture.common.labelmatcher.v1.LabelMatcher
}
var file_aperture_common_selector_v1_selector_proto_depIdxs = []int32{
	1, // 0: aperture.common.selector.v1.Selector.control_point:type_name -> aperture.common.selector.v1.ControlPoint
	2, // 1: aperture.common.selector.v1.Selector.label_matcher:type_name -> aperture.common.labelmatcher.v1.LabelMatcher
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_aperture_common_selector_v1_selector_proto_init() }
func file_aperture_common_selector_v1_selector_proto_init() {
	if File_aperture_common_selector_v1_selector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_common_selector_v1_selector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Selector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_aperture_common_selector_v1_selector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ControlPoint); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_aperture_common_selector_v1_selector_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ControlPoint_Feature)(nil),
		(*ControlPoint_Traffic)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_common_selector_v1_selector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_common_selector_v1_selector_proto_goTypes,
		DependencyIndexes: file_aperture_common_selector_v1_selector_proto_depIdxs,
		MessageInfos:      file_aperture_common_selector_v1_selector_proto_msgTypes,
	}.Build()
	File_aperture_common_selector_v1_selector_proto = out.File
	file_aperture_common_selector_v1_selector_proto_rawDesc = nil
	file_aperture_common_selector_v1_selector_proto_goTypes = nil
	file_aperture_common_selector_v1_selector_proto_depIdxs = nil
}
