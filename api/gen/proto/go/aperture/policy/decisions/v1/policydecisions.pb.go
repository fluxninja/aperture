// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/policy/decisions/v1/policydecisions.proto

package decisionsv1

import (
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

type LoadShedDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoadShedFactor float64 `protobuf:"fixed64,1,opt,name=load_shed_factor,json=loadShedFactor,proto3" json:"load_shed_factor,omitempty"`
}

func (x *LoadShedDecision) Reset() {
	*x = LoadShedDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadShedDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadShedDecision) ProtoMessage() {}

func (x *LoadShedDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadShedDecision.ProtoReflect.Descriptor instead.
func (*LoadShedDecision) Descriptor() ([]byte, []int) {
	return file_aperture_policy_decisions_v1_policydecisions_proto_rawDescGZIP(), []int{0}
}

func (x *LoadShedDecision) GetLoadShedFactor() float64 {
	if x != nil {
		return x.LoadShedFactor
	}
	return 0
}

type TokensDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefaultWorkloadTokens uint64 `protobuf:"varint,1,opt,name=default_workload_tokens,json=defaultWorkloadTokens,proto3" json:"default_workload_tokens,omitempty" default:"1"` // @gotags: default:"1"
	// Key in map is a string representation of WorkloadDesc message.
	TokensByWorkload map[string]uint64 `protobuf:"bytes,2,rep,name=tokens_by_workload,json=tokensByWorkload,proto3" json:"tokens_by_workload,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *TokensDecision) Reset() {
	*x = TokensDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokensDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokensDecision) ProtoMessage() {}

func (x *TokensDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokensDecision.ProtoReflect.Descriptor instead.
func (*TokensDecision) Descriptor() ([]byte, []int) {
	return file_aperture_policy_decisions_v1_policydecisions_proto_rawDescGZIP(), []int{1}
}

func (x *TokensDecision) GetDefaultWorkloadTokens() uint64 {
	if x != nil {
		return x.DefaultWorkloadTokens
	}
	return 0
}

func (x *TokensDecision) GetTokensByWorkload() map[string]uint64 {
	if x != nil {
		return x.TokensByWorkload
	}
	return nil
}

type WorkloadDesc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkloadKey   string `protobuf:"bytes,1,opt,name=workload_key,json=workloadKey,proto3" json:"workload_key,omitempty" default:"default_workload_key"`       // @gotags: default:"default_workload_key"
	WorkloadValue string `protobuf:"bytes,2,opt,name=workload_value,json=workloadValue,proto3" json:"workload_value,omitempty" default:"default_workload_value"` // @gotags: default:"default_workload_value"
}

func (x *WorkloadDesc) Reset() {
	*x = WorkloadDesc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkloadDesc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkloadDesc) ProtoMessage() {}

func (x *WorkloadDesc) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkloadDesc.ProtoReflect.Descriptor instead.
func (*WorkloadDesc) Descriptor() ([]byte, []int) {
	return file_aperture_policy_decisions_v1_policydecisions_proto_rawDescGZIP(), []int{2}
}

func (x *WorkloadDesc) GetWorkloadKey() string {
	if x != nil {
		return x.WorkloadKey
	}
	return ""
}

func (x *WorkloadDesc) GetWorkloadValue() string {
	if x != nil {
		return x.WorkloadValue
	}
	return ""
}

type RateLimiterDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit float64 `protobuf:"fixed64,1,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *RateLimiterDecision) Reset() {
	*x = RateLimiterDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterDecision) ProtoMessage() {}

func (x *RateLimiterDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimiterDecision.ProtoReflect.Descriptor instead.
func (*RateLimiterDecision) Descriptor() ([]byte, []int) {
	return file_aperture_policy_decisions_v1_policydecisions_proto_rawDescGZIP(), []int{3}
}

func (x *RateLimiterDecision) GetLimit() float64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

var File_aperture_policy_decisions_v1_policydecisions_proto protoreflect.FileDescriptor

var file_aperture_policy_decisions_v1_policydecisions_proto_rawDesc = []byte{
	0x0a, 0x32, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x76, 0x31, 0x22, 0x3c, 0x0a, 0x10, 0x4c, 0x6f, 0x61, 0x64, 0x53, 0x68, 0x65, 0x64, 0x44, 0x65,
	0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x10, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x73,
	0x68, 0x65, 0x64, 0x5f, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x0e, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x68, 0x65, 0x64, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x22, 0xff, 0x01, 0x0a, 0x0e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x17, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x77,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x15, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x57, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x12, 0x70, 0x0a, 0x12, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x5f, 0x62, 0x79, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x63,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x42, 0x79, 0x57, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x42, 0x79, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x43, 0x0a,
	0x15, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x42, 0x79, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x58, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x44, 0x65,
	0x73, 0x63, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x77,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x2b, 0x0a, 0x13,
	0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x42, 0xa4, 0x02, 0x0a, 0x20, 0x63, 0x6f,
	0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x14,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x57, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x46, 0x6c, 0x75, 0x78, 0x4e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x76, 0x31, 0x3b, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x41, 0x50, 0x44, 0xaa, 0x02, 0x1c, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1c, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x28, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x1f, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x3a, 0x3a, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_decisions_v1_policydecisions_proto_rawDescOnce sync.Once
	file_aperture_policy_decisions_v1_policydecisions_proto_rawDescData = file_aperture_policy_decisions_v1_policydecisions_proto_rawDesc
)

func file_aperture_policy_decisions_v1_policydecisions_proto_rawDescGZIP() []byte {
	file_aperture_policy_decisions_v1_policydecisions_proto_rawDescOnce.Do(func() {
		file_aperture_policy_decisions_v1_policydecisions_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_decisions_v1_policydecisions_proto_rawDescData)
	})
	return file_aperture_policy_decisions_v1_policydecisions_proto_rawDescData
}

var file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_aperture_policy_decisions_v1_policydecisions_proto_goTypes = []interface{}{
	(*LoadShedDecision)(nil),    // 0: aperture.policy.decisions.v1.LoadShedDecision
	(*TokensDecision)(nil),      // 1: aperture.policy.decisions.v1.TokensDecision
	(*WorkloadDesc)(nil),        // 2: aperture.policy.decisions.v1.WorkloadDesc
	(*RateLimiterDecision)(nil), // 3: aperture.policy.decisions.v1.RateLimiterDecision
	nil,                         // 4: aperture.policy.decisions.v1.TokensDecision.TokensByWorkloadEntry
}
var file_aperture_policy_decisions_v1_policydecisions_proto_depIdxs = []int32{
	4, // 0: aperture.policy.decisions.v1.TokensDecision.tokens_by_workload:type_name -> aperture.policy.decisions.v1.TokensDecision.TokensByWorkloadEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_aperture_policy_decisions_v1_policydecisions_proto_init() }
func file_aperture_policy_decisions_v1_policydecisions_proto_init() {
	if File_aperture_policy_decisions_v1_policydecisions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadShedDecision); i {
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
		file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokensDecision); i {
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
		file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkloadDesc); i {
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
		file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimiterDecision); i {
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
			RawDescriptor: file_aperture_policy_decisions_v1_policydecisions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_decisions_v1_policydecisions_proto_goTypes,
		DependencyIndexes: file_aperture_policy_decisions_v1_policydecisions_proto_depIdxs,
		MessageInfos:      file_aperture_policy_decisions_v1_policydecisions_proto_msgTypes,
	}.Build()
	File_aperture_policy_decisions_v1_policydecisions_proto = out.File
	file_aperture_policy_decisions_v1_policydecisions_proto_rawDesc = nil
	file_aperture_policy_decisions_v1_policydecisions_proto_goTypes = nil
	file_aperture_policy_decisions_v1_policydecisions_proto_depIdxs = nil
}
