// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: aperture/policy/private/v1/flowcontrol.proto

package privatev1

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

type LoadActuator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InPorts                    *LoadActuator_Ins `protobuf:"bytes,1,opt,name=in_ports,json=inPorts,proto3" json:"in_ports,omitempty"`
	LoadSchedulerComponentId   string            `protobuf:"bytes,2,opt,name=load_scheduler_component_id,json=loadSchedulerComponentId,proto3" json:"load_scheduler_component_id,omitempty"`
	Selectors                  []*v1.Selector    `protobuf:"bytes,3,rep,name=selectors,proto3" json:"selectors,omitempty"`
	DryRun                     bool              `protobuf:"varint,4,opt,name=dry_run,json=dryRun,proto3" json:"dry_run,omitempty"`
	DryRunConfigKey            string            `protobuf:"bytes,5,opt,name=dry_run_config_key,json=dryRunConfigKey,proto3" json:"dry_run_config_key,omitempty"`
	WorkloadLatencyBasedTokens bool              `protobuf:"varint,6,opt,name=workload_latency_based_tokens,json=workloadLatencyBasedTokens,proto3" json:"workload_latency_based_tokens,omitempty"`
}

func (x *LoadActuator) Reset() {
	*x = LoadActuator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadActuator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadActuator) ProtoMessage() {}

func (x *LoadActuator) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadActuator.ProtoReflect.Descriptor instead.
func (*LoadActuator) Descriptor() ([]byte, []int) {
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP(), []int{0}
}

func (x *LoadActuator) GetInPorts() *LoadActuator_Ins {
	if x != nil {
		return x.InPorts
	}
	return nil
}

func (x *LoadActuator) GetLoadSchedulerComponentId() string {
	if x != nil {
		return x.LoadSchedulerComponentId
	}
	return ""
}

func (x *LoadActuator) GetSelectors() []*v1.Selector {
	if x != nil {
		return x.Selectors
	}
	return nil
}

func (x *LoadActuator) GetDryRun() bool {
	if x != nil {
		return x.DryRun
	}
	return false
}

func (x *LoadActuator) GetDryRunConfigKey() string {
	if x != nil {
		return x.DryRunConfigKey
	}
	return ""
}

func (x *LoadActuator) GetWorkloadLatencyBasedTokens() bool {
	if x != nil {
		return x.WorkloadLatencyBasedTokens
	}
	return false
}

type LoadActuator_Ins struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Load multiplier is proportion of incoming
	// token rate that needs to be accepted.
	LoadMultiplier *v1.InPort `protobuf:"bytes,1,opt,name=load_multiplier,json=loadMultiplier,proto3" json:"load_multiplier,omitempty"`
	PassThrough    *v1.InPort `protobuf:"bytes,2,opt,name=pass_through,json=passThrough,proto3" json:"pass_through,omitempty"`
}

func (x *LoadActuator_Ins) Reset() {
	*x = LoadActuator_Ins{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadActuator_Ins) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadActuator_Ins) ProtoMessage() {}

func (x *LoadActuator_Ins) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadActuator_Ins.ProtoReflect.Descriptor instead.
func (*LoadActuator_Ins) Descriptor() ([]byte, []int) {
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP(), []int{0, 0}
}

func (x *LoadActuator_Ins) GetLoadMultiplier() *v1.InPort {
	if x != nil {
		return x.LoadMultiplier
	}
	return nil
}

func (x *LoadActuator_Ins) GetPassThrough() *v1.InPort {
	if x != nil {
		return x.PassThrough
	}
	return nil
}

var File_aperture_policy_private_v1_flowcontrol_proto protoreflect.FileDescriptor

var file_aperture_policy_private_v1_flowcontrol_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6c, 0x6f,
	0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x2d, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x82, 0x04, 0x0a, 0x0c, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x63, 0x74, 0x75, 0x61,
	0x74, 0x6f, 0x72, 0x12, 0x47, 0x0a, 0x08, 0x69, 0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x63, 0x74, 0x75, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x49, 0x6e, 0x73, 0x52, 0x07, 0x69, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x3d, 0x0a, 0x1b,
	0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x5f, 0x63,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x18, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72,
	0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x43, 0x0a, 0x09, 0x73,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25,
	0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73,
	0x12, 0x17, 0x0a, 0x07, 0x64, 0x72, 0x79, 0x5f, 0x72, 0x75, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x64, 0x72, 0x79, 0x52, 0x75, 0x6e, 0x12, 0x2b, 0x0a, 0x12, 0x64, 0x72, 0x79,
	0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x72, 0x79, 0x52, 0x75, 0x6e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x4b, 0x65, 0x79, 0x12, 0x41, 0x0a, 0x1d, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x64,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x1a, 0x77,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x42, 0x61,
	0x73, 0x65, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x1a, 0x9b, 0x01, 0x0a, 0x03, 0x49, 0x6e,
	0x73, 0x12, 0x4c, 0x0a, 0x0f, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70,
	0x6c, 0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x52,
	0x0e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x12,
	0x46, 0x0a, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x0b, 0x70, 0x61, 0x73, 0x73,
	0x54, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x42, 0xab, 0x02, 0x0a, 0x32, 0x63, 0x6f, 0x6d, 0x2e,
	0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x10,
	0x46, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x56, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31,
	0x3b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x50,
	0xaa, 0x02, 0x1a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1b,
	0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c,
	0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x50, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1d, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x3a, 0x3a, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_private_v1_flowcontrol_proto_rawDescOnce sync.Once
	file_aperture_policy_private_v1_flowcontrol_proto_rawDescData = file_aperture_policy_private_v1_flowcontrol_proto_rawDesc
)

func file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP() []byte {
	file_aperture_policy_private_v1_flowcontrol_proto_rawDescOnce.Do(func() {
		file_aperture_policy_private_v1_flowcontrol_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_private_v1_flowcontrol_proto_rawDescData)
	})
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescData
}

var file_aperture_policy_private_v1_flowcontrol_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_aperture_policy_private_v1_flowcontrol_proto_goTypes = []interface{}{
	(*LoadActuator)(nil),     // 0: aperture.policy.private.v1.LoadActuator
	(*LoadActuator_Ins)(nil), // 1: aperture.policy.private.v1.LoadActuator.Ins
	(*v1.Selector)(nil),      // 2: aperture.policy.language.v1.Selector
	(*v1.InPort)(nil),        // 3: aperture.policy.language.v1.InPort
}
var file_aperture_policy_private_v1_flowcontrol_proto_depIdxs = []int32{
	1, // 0: aperture.policy.private.v1.LoadActuator.in_ports:type_name -> aperture.policy.private.v1.LoadActuator.Ins
	2, // 1: aperture.policy.private.v1.LoadActuator.selectors:type_name -> aperture.policy.language.v1.Selector
	3, // 2: aperture.policy.private.v1.LoadActuator.Ins.load_multiplier:type_name -> aperture.policy.language.v1.InPort
	3, // 3: aperture.policy.private.v1.LoadActuator.Ins.pass_through:type_name -> aperture.policy.language.v1.InPort
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_aperture_policy_private_v1_flowcontrol_proto_init() }
func file_aperture_policy_private_v1_flowcontrol_proto_init() {
	if File_aperture_policy_private_v1_flowcontrol_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadActuator); i {
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
		file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadActuator_Ins); i {
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
			RawDescriptor: file_aperture_policy_private_v1_flowcontrol_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_private_v1_flowcontrol_proto_goTypes,
		DependencyIndexes: file_aperture_policy_private_v1_flowcontrol_proto_depIdxs,
		MessageInfos:      file_aperture_policy_private_v1_flowcontrol_proto_msgTypes,
	}.Build()
	File_aperture_policy_private_v1_flowcontrol_proto = out.File
	file_aperture_policy_private_v1_flowcontrol_proto_rawDesc = nil
	file_aperture_policy_private_v1_flowcontrol_proto_goTypes = nil
	file_aperture_policy_private_v1_flowcontrol_proto_depIdxs = nil
}
