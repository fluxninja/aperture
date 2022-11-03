// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/policy/sync/v1/classifier.proto

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

type ClassifierWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Classifier
	Classifier *v1.Classifier `protobuf:"bytes,2,opt,name=classifier,proto3" json:"classifier,omitempty"`
}

func (x *ClassifierWrapper) Reset() {
	*x = ClassifierWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_classifier_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassifierWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassifierWrapper) ProtoMessage() {}

func (x *ClassifierWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_classifier_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassifierWrapper.ProtoReflect.Descriptor instead.
func (*ClassifierWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_classifier_proto_rawDescGZIP(), []int{0}
}

func (x *ClassifierWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *ClassifierWrapper) GetClassifier() *v1.Classifier {
	if x != nil {
		return x.Classifier
	}
	return nil
}

var File_aperture_policy_sync_v1_classifier_proto protoreflect.FileDescriptor

var file_aperture_policy_sync_v1_classifier_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63,
	0x2e, 0x76, 0x31, 0x1a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb4, 0x01, 0x0a, 0x11, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x12, 0x47, 0x0a, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x42, 0x90, 0x02, 0x0a, 0x2f, 0x63, 0x6f,
	0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75,
	0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f,
	0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x79, 0x6e, 0x63, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x41, 0x50, 0x53, 0xaa, 0x02, 0x17, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x17, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x5c, 0x53, 0x79, 0x6e, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x23, 0x41, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x53, 0x79, 0x6e, 0x63,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x1a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x3a, 0x3a, 0x53, 0x79, 0x6e, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_sync_v1_classifier_proto_rawDescOnce sync.Once
	file_aperture_policy_sync_v1_classifier_proto_rawDescData = file_aperture_policy_sync_v1_classifier_proto_rawDesc
)

func file_aperture_policy_sync_v1_classifier_proto_rawDescGZIP() []byte {
	file_aperture_policy_sync_v1_classifier_proto_rawDescOnce.Do(func() {
		file_aperture_policy_sync_v1_classifier_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_sync_v1_classifier_proto_rawDescData)
	})
	return file_aperture_policy_sync_v1_classifier_proto_rawDescData
}

var file_aperture_policy_sync_v1_classifier_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_aperture_policy_sync_v1_classifier_proto_goTypes = []interface{}{
	(*ClassifierWrapper)(nil), // 0: aperture.policy.sync.v1.ClassifierWrapper
	(*CommonAttributes)(nil),  // 1: aperture.policy.sync.v1.CommonAttributes
	(*v1.Classifier)(nil),     // 2: aperture.policy.language.v1.Classifier
}
var file_aperture_policy_sync_v1_classifier_proto_depIdxs = []int32{
	1, // 0: aperture.policy.sync.v1.ClassifierWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	2, // 1: aperture.policy.sync.v1.ClassifierWrapper.classifier:type_name -> aperture.policy.language.v1.Classifier
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_aperture_policy_sync_v1_classifier_proto_init() }
func file_aperture_policy_sync_v1_classifier_proto_init() {
	if File_aperture_policy_sync_v1_classifier_proto != nil {
		return
	}
	file_aperture_policy_sync_v1_common_attributes_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_sync_v1_classifier_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassifierWrapper); i {
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
			RawDescriptor: file_aperture_policy_sync_v1_classifier_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_sync_v1_classifier_proto_goTypes,
		DependencyIndexes: file_aperture_policy_sync_v1_classifier_proto_depIdxs,
		MessageInfos:      file_aperture_policy_sync_v1_classifier_proto_msgTypes,
	}.Build()
	File_aperture_policy_sync_v1_classifier_proto = out.File
	file_aperture_policy_sync_v1_classifier_proto_rawDesc = nil
	file_aperture_policy_sync_v1_classifier_proto_goTypes = nil
	file_aperture_policy_sync_v1_classifier_proto_depIdxs = nil
}
