// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: aperture/policy/sync/v1/flow_regulator.proto

package syncv1

import (
	v1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
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

type FlowRegulatorWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Flow Regulator
	FlowRegulator *v1.FlowRegulator `protobuf:"bytes,2,opt,name=flow_regulator,json=flowRegulator,proto3" json:"flow_regulator,omitempty"`
}

func (x *FlowRegulatorWrapper) Reset() {
	*x = FlowRegulatorWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowRegulatorWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowRegulatorWrapper) ProtoMessage() {}

func (x *FlowRegulatorWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowRegulatorWrapper.ProtoReflect.Descriptor instead.
func (*FlowRegulatorWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_flow_regulator_proto_rawDescGZIP(), []int{0}
}

func (x *FlowRegulatorWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *FlowRegulatorWrapper) GetFlowRegulator() *v1.FlowRegulator {
	if x != nil {
		return x.FlowRegulator
	}
	return nil
}

type FlowRegulatorDynamicConfigWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// FlowRegulatorDynamicConfig is the dynamic configuration for the Flow Regulator.
	FlowRegulatorDynamicConfig *v1.FlowRegulator_DynamicConfig `protobuf:"bytes,2,opt,name=flow_regulator_dynamic_config,json=flowRegulatorDynamicConfig,proto3" json:"flow_regulator_dynamic_config,omitempty"`
}

func (x *FlowRegulatorDynamicConfigWrapper) Reset() {
	*x = FlowRegulatorDynamicConfigWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowRegulatorDynamicConfigWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowRegulatorDynamicConfigWrapper) ProtoMessage() {}

func (x *FlowRegulatorDynamicConfigWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowRegulatorDynamicConfigWrapper.ProtoReflect.Descriptor instead.
func (*FlowRegulatorDynamicConfigWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_flow_regulator_proto_rawDescGZIP(), []int{1}
}

func (x *FlowRegulatorDynamicConfigWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *FlowRegulatorDynamicConfigWrapper) GetFlowRegulatorDynamicConfig() *v1.FlowRegulator_DynamicConfig {
	if x != nil {
		return x.FlowRegulatorDynamicConfig
	}
	return nil
}

type FlowRegulatorDecisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Flow Regulator Decision
	FlowRegulatorDecision *FlowRegulatorDecision `protobuf:"bytes,2,opt,name=flow_regulator_decision,json=flowRegulatorDecision,proto3" json:"flow_regulator_decision,omitempty"`
}

func (x *FlowRegulatorDecisionWrapper) Reset() {
	*x = FlowRegulatorDecisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowRegulatorDecisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowRegulatorDecisionWrapper) ProtoMessage() {}

func (x *FlowRegulatorDecisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowRegulatorDecisionWrapper.ProtoReflect.Descriptor instead.
func (*FlowRegulatorDecisionWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_flow_regulator_proto_rawDescGZIP(), []int{2}
}

func (x *FlowRegulatorDecisionWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *FlowRegulatorDecisionWrapper) GetFlowRegulatorDecision() *FlowRegulatorDecision {
	if x != nil {
		return x.FlowRegulatorDecision
	}
	return nil
}

type FlowRegulatorDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptPercentage float64 `protobuf:"fixed64,1,opt,name=accept_percentage,json=acceptPercentage,proto3" json:"accept_percentage,omitempty"`
}

func (x *FlowRegulatorDecision) Reset() {
	*x = FlowRegulatorDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowRegulatorDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowRegulatorDecision) ProtoMessage() {}

func (x *FlowRegulatorDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowRegulatorDecision.ProtoReflect.Descriptor instead.
func (*FlowRegulatorDecision) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_flow_regulator_proto_rawDescGZIP(), []int{3}
}

func (x *FlowRegulatorDecision) GetAcceptPercentage() float64 {
	if x != nil {
		return x.AcceptPercentage
	}
	return 0
}

var File_aperture_policy_sync_v1_flow_regulator_proto protoreflect.FileDescriptor

var file_aperture_policy_sync_v1_flow_regulator_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x72,
	0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x2d, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x01, 0x0a, 0x14, 0x46, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x51, 0x0a, 0x0e, 0x66, 0x6c, 0x6f, 0x77,
	0x5f, 0x72, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2a, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x0d, 0x66, 0x6c,
	0x6f, 0x77, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x22, 0xf8, 0x01, 0x0a, 0x21,
	0x46, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x79, 0x6e,
	0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73,
	0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x7b, 0x0a, 0x1d, 0x66, 0x6c, 0x6f,
	0x77, 0x5f, 0x72, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x64, 0x79, 0x6e, 0x61,
	0x6d, 0x69, 0x63, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x38, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x44, 0x79, 0x6e,
	0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x1a, 0x66, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xde, 0x01, 0x0a, 0x1c, 0x46, 0x6c, 0x6f, 0x77, 0x52,
	0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12,
	0x66, 0x0a, 0x17, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x72, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f,
	0x72, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2e, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x52,
	0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x15, 0x66, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x44,
	0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x44, 0x0a, 0x15, 0x46, 0x6c, 0x6f, 0x77, 0x52,
	0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x2b, 0x0a, 0x11, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x61, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x42, 0x93, 0x02,
	0x0a, 0x2f, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76,
	0x31, 0x42, 0x12, 0x46, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x3b,
	0x73, 0x79, 0x6e, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x53, 0xaa, 0x02, 0x17, 0x41,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x53,
	0x79, 0x6e, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x17, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x53, 0x79, 0x6e, 0x63, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x23, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x5c, 0x53, 0x79, 0x6e, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x3a, 0x3a, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x53, 0x79, 0x6e, 0x63, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_sync_v1_flow_regulator_proto_rawDescOnce sync.Once
	file_aperture_policy_sync_v1_flow_regulator_proto_rawDescData = file_aperture_policy_sync_v1_flow_regulator_proto_rawDesc
)

func file_aperture_policy_sync_v1_flow_regulator_proto_rawDescGZIP() []byte {
	file_aperture_policy_sync_v1_flow_regulator_proto_rawDescOnce.Do(func() {
		file_aperture_policy_sync_v1_flow_regulator_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_sync_v1_flow_regulator_proto_rawDescData)
	})
	return file_aperture_policy_sync_v1_flow_regulator_proto_rawDescData
}

var file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_aperture_policy_sync_v1_flow_regulator_proto_goTypes = []interface{}{
	(*FlowRegulatorWrapper)(nil),              // 0: aperture.policy.sync.v1.FlowRegulatorWrapper
	(*FlowRegulatorDynamicConfigWrapper)(nil), // 1: aperture.policy.sync.v1.FlowRegulatorDynamicConfigWrapper
	(*FlowRegulatorDecisionWrapper)(nil),      // 2: aperture.policy.sync.v1.FlowRegulatorDecisionWrapper
	(*FlowRegulatorDecision)(nil),             // 3: aperture.policy.sync.v1.FlowRegulatorDecision
	(*CommonAttributes)(nil),                  // 4: aperture.policy.sync.v1.CommonAttributes
	(*v1.FlowRegulator)(nil),                  // 5: aperture.policy.language.v1.FlowRegulator
	(*v1.FlowRegulator_DynamicConfig)(nil),    // 6: aperture.policy.language.v1.FlowRegulator.DynamicConfig
}
var file_aperture_policy_sync_v1_flow_regulator_proto_depIdxs = []int32{
	4, // 0: aperture.policy.sync.v1.FlowRegulatorWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	5, // 1: aperture.policy.sync.v1.FlowRegulatorWrapper.flow_regulator:type_name -> aperture.policy.language.v1.FlowRegulator
	4, // 2: aperture.policy.sync.v1.FlowRegulatorDynamicConfigWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	6, // 3: aperture.policy.sync.v1.FlowRegulatorDynamicConfigWrapper.flow_regulator_dynamic_config:type_name -> aperture.policy.language.v1.FlowRegulator.DynamicConfig
	4, // 4: aperture.policy.sync.v1.FlowRegulatorDecisionWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	3, // 5: aperture.policy.sync.v1.FlowRegulatorDecisionWrapper.flow_regulator_decision:type_name -> aperture.policy.sync.v1.FlowRegulatorDecision
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_aperture_policy_sync_v1_flow_regulator_proto_init() }
func file_aperture_policy_sync_v1_flow_regulator_proto_init() {
	if File_aperture_policy_sync_v1_flow_regulator_proto != nil {
		return
	}
	file_aperture_policy_sync_v1_common_attributes_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowRegulatorWrapper); i {
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
		file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowRegulatorDynamicConfigWrapper); i {
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
		file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowRegulatorDecisionWrapper); i {
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
		file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowRegulatorDecision); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_policy_sync_v1_flow_regulator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_sync_v1_flow_regulator_proto_goTypes,
		DependencyIndexes: file_aperture_policy_sync_v1_flow_regulator_proto_depIdxs,
		MessageInfos:      file_aperture_policy_sync_v1_flow_regulator_proto_msgTypes,
	}.Build()
	File_aperture_policy_sync_v1_flow_regulator_proto = out.File
	file_aperture_policy_sync_v1_flow_regulator_proto_rawDesc = nil
	file_aperture_policy_sync_v1_flow_regulator_proto_goTypes = nil
	file_aperture_policy_sync_v1_flow_regulator_proto_depIdxs = nil
}
