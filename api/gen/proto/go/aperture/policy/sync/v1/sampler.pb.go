// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: aperture/policy/sync/v1/sampler.proto

package syncv1

import (
	v1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
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

type SamplerWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Sampler
	Sampler *v1.Sampler `protobuf:"bytes,2,opt,name=sampler,proto3" json:"sampler,omitempty"`
}

func (x *SamplerWrapper) Reset() {
	*x = SamplerWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_sampler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SamplerWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SamplerWrapper) ProtoMessage() {}

func (x *SamplerWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_sampler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SamplerWrapper.ProtoReflect.Descriptor instead.
func (*SamplerWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_sampler_proto_rawDescGZIP(), []int{0}
}

func (x *SamplerWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *SamplerWrapper) GetSampler() *v1.Sampler {
	if x != nil {
		return x.Sampler
	}
	return nil
}

type SamplerDecisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Sampler Decision
	SamplerDecision *SamplerDecision `protobuf:"bytes,2,opt,name=sampler_decision,json=samplerDecision,proto3" json:"sampler_decision,omitempty"`
}

func (x *SamplerDecisionWrapper) Reset() {
	*x = SamplerDecisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_sampler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SamplerDecisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SamplerDecisionWrapper) ProtoMessage() {}

func (x *SamplerDecisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_sampler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SamplerDecisionWrapper.ProtoReflect.Descriptor instead.
func (*SamplerDecisionWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_sampler_proto_rawDescGZIP(), []int{1}
}

func (x *SamplerDecisionWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *SamplerDecisionWrapper) GetSamplerDecision() *SamplerDecision {
	if x != nil {
		return x.SamplerDecision
	}
	return nil
}

type SamplerDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptPercentage float64 `protobuf:"fixed64,1,opt,name=accept_percentage,json=acceptPercentage,proto3" json:"accept_percentage,omitempty"`
	// PassThroughLabelValues dynamic config.
	PassThroughLabelValues []string `protobuf:"bytes,2,rep,name=pass_through_label_values,json=passThroughLabelValues,proto3" json:"pass_through_label_values,omitempty"`
}

func (x *SamplerDecision) Reset() {
	*x = SamplerDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_sampler_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SamplerDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SamplerDecision) ProtoMessage() {}

func (x *SamplerDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_sampler_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SamplerDecision.ProtoReflect.Descriptor instead.
func (*SamplerDecision) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_sampler_proto_rawDescGZIP(), []int{2}
}

func (x *SamplerDecision) GetAcceptPercentage() float64 {
	if x != nil {
		return x.AcceptPercentage
	}
	return 0
}

func (x *SamplerDecision) GetPassThroughLabelValues() []string {
	if x != nil {
		return x.PassThroughLabelValues
	}
	return nil
}

var File_aperture_policy_sync_v1_sampler_proto protoreflect.FileDescriptor

var file_aperture_policy_sync_v1_sampler_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31,
	0x1a, 0x2d, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6c,
	0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f,
	0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xa8, 0x01, 0x0a, 0x0e, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29,
	0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x3e, 0x0a, 0x07, 0x73,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x72, 0x52, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x22, 0xc5, 0x01, 0x0a, 0x16,
	0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x57,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x53,
	0x0a, 0x10, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x0f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0x79, 0x0a, 0x0f, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x44, 0x65,
	0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2b, 0x0a, 0x11, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x10, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x61, 0x67, 0x65, 0x12, 0x39, 0x0a, 0x19, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x74, 0x68, 0x72, 0x6f,
	0x75, 0x67, 0x68, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x16, 0x70, 0x61, 0x73, 0x73, 0x54, 0x68, 0x72, 0x6f,
	0x75, 0x67, 0x68, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x90,
	0x02, 0x0a, 0x2f, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61,
	0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e,
	0x76, 0x31, 0x42, 0x0c, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x50, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x79,
	0x6e, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x53, 0xaa, 0x02, 0x17, 0x41, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x53, 0x79, 0x6e,
	0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x17, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x53, 0x79, 0x6e, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x23, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x5c, 0x53, 0x79, 0x6e, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a,
	0x3a, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x53, 0x79, 0x6e, 0x63, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_sync_v1_sampler_proto_rawDescOnce sync.Once
	file_aperture_policy_sync_v1_sampler_proto_rawDescData = file_aperture_policy_sync_v1_sampler_proto_rawDesc
)

func file_aperture_policy_sync_v1_sampler_proto_rawDescGZIP() []byte {
	file_aperture_policy_sync_v1_sampler_proto_rawDescOnce.Do(func() {
		file_aperture_policy_sync_v1_sampler_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_sync_v1_sampler_proto_rawDescData)
	})
	return file_aperture_policy_sync_v1_sampler_proto_rawDescData
}

var file_aperture_policy_sync_v1_sampler_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_aperture_policy_sync_v1_sampler_proto_goTypes = []interface{}{
	(*SamplerWrapper)(nil),         // 0: aperture.policy.sync.v1.SamplerWrapper
	(*SamplerDecisionWrapper)(nil), // 1: aperture.policy.sync.v1.SamplerDecisionWrapper
	(*SamplerDecision)(nil),        // 2: aperture.policy.sync.v1.SamplerDecision
	(*CommonAttributes)(nil),       // 3: aperture.policy.sync.v1.CommonAttributes
	(*v1.Sampler)(nil),             // 4: aperture.policy.language.v1.Sampler
}
var file_aperture_policy_sync_v1_sampler_proto_depIdxs = []int32{
	3, // 0: aperture.policy.sync.v1.SamplerWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	4, // 1: aperture.policy.sync.v1.SamplerWrapper.sampler:type_name -> aperture.policy.language.v1.Sampler
	3, // 2: aperture.policy.sync.v1.SamplerDecisionWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	2, // 3: aperture.policy.sync.v1.SamplerDecisionWrapper.sampler_decision:type_name -> aperture.policy.sync.v1.SamplerDecision
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_aperture_policy_sync_v1_sampler_proto_init() }
func file_aperture_policy_sync_v1_sampler_proto_init() {
	if File_aperture_policy_sync_v1_sampler_proto != nil {
		return
	}
	file_aperture_policy_sync_v1_common_attributes_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_sync_v1_sampler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SamplerWrapper); i {
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
		file_aperture_policy_sync_v1_sampler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SamplerDecisionWrapper); i {
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
		file_aperture_policy_sync_v1_sampler_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SamplerDecision); i {
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
			RawDescriptor: file_aperture_policy_sync_v1_sampler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_sync_v1_sampler_proto_goTypes,
		DependencyIndexes: file_aperture_policy_sync_v1_sampler_proto_depIdxs,
		MessageInfos:      file_aperture_policy_sync_v1_sampler_proto_msgTypes,
	}.Build()
	File_aperture_policy_sync_v1_sampler_proto = out.File
	file_aperture_policy_sync_v1_sampler_proto_rawDesc = nil
	file_aperture_policy_sync_v1_sampler_proto_goTypes = nil
	file_aperture_policy_sync_v1_sampler_proto_depIdxs = nil
}
