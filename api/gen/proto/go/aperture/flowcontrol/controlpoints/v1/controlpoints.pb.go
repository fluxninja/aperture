// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: aperture/flowcontrol/controlpoints/v1/controlpoints.proto

package controlpointsv1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FlowControlPoints struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FlowControlPoints []*FlowControlPoint `protobuf:"bytes,1,rep,name=flow_control_points,json=flowControlPoints,proto3" json:"flow_control_points,omitempty"`
}

func (x *FlowControlPoints) Reset() {
	*x = FlowControlPoints{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowControlPoints) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowControlPoints) ProtoMessage() {}

func (x *FlowControlPoints) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowControlPoints.ProtoReflect.Descriptor instead.
func (*FlowControlPoints) Descriptor() ([]byte, []int) {
	return file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescGZIP(), []int{0}
}

func (x *FlowControlPoints) GetFlowControlPoints() []*FlowControlPoint {
	if x != nil {
		return x.FlowControlPoints
	}
	return nil
}

type FlowControlPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ControlPoint string `protobuf:"bytes,2,opt,name=control_point,json=controlPoint,proto3" json:"control_point,omitempty"`
	Type         string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Service      string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
}

func (x *FlowControlPoint) Reset() {
	*x = FlowControlPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowControlPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowControlPoint) ProtoMessage() {}

func (x *FlowControlPoint) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowControlPoint.ProtoReflect.Descriptor instead.
func (*FlowControlPoint) Descriptor() ([]byte, []int) {
	return file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescGZIP(), []int{1}
}

func (x *FlowControlPoint) GetControlPoint() string {
	if x != nil {
		return x.ControlPoint
	}
	return ""
}

func (x *FlowControlPoint) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *FlowControlPoint) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

var File_aperture_flowcontrol_controlpoints_v1_controlpoints_proto protoreflect.FileDescriptor

var file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDesc = []byte{
	0x0a, 0x39, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x25, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70,
	0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7c, 0x0a,
	0x11, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x12, 0x67, 0x0a, 0x13, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x37, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x11, 0x66, 0x6c, 0x6f, 0x77, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x65, 0x0a, 0x10, 0x46,
	0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x32, 0xbb, 0x01, 0x0a, 0x18, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x9e, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x38, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x38, 0x92, 0x41, 0x10, 0x0a, 0x0e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1f, 0x12, 0x1d, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x42, 0xf3, 0x02, 0x0a, 0x3d, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e,
	0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e,
	0x76, 0x31, 0x42, 0x12, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x67, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x76,
	0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x41, 0x46, 0x43, 0xaa, 0x02, 0x25, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x25, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x46, 0x6c, 0x6f, 0x77, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x31, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x5c, 0x46, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5c, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x28, 0x41, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x46, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x3a, 0x3a, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescOnce sync.Once
	file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescData = file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDesc
)

func file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescGZIP() []byte {
	file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescOnce.Do(func() {
		file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescData)
	})
	return file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDescData
}

var file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_goTypes = []interface{}{
	(*FlowControlPoints)(nil), // 0: aperture.flowcontrol.controlpoints.v1.FlowControlPoints
	(*FlowControlPoint)(nil),  // 1: aperture.flowcontrol.controlpoints.v1.FlowControlPoint
	(*emptypb.Empty)(nil),     // 2: google.protobuf.Empty
}
var file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_depIdxs = []int32{
	1, // 0: aperture.flowcontrol.controlpoints.v1.FlowControlPoints.flow_control_points:type_name -> aperture.flowcontrol.controlpoints.v1.FlowControlPoint
	2, // 1: aperture.flowcontrol.controlpoints.v1.FlowControlPointsService.GetControlPoints:input_type -> google.protobuf.Empty
	0, // 2: aperture.flowcontrol.controlpoints.v1.FlowControlPointsService.GetControlPoints:output_type -> aperture.flowcontrol.controlpoints.v1.FlowControlPoints
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_init() }
func file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_init() {
	if File_aperture_flowcontrol_controlpoints_v1_controlpoints_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowControlPoints); i {
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
		file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowControlPoint); i {
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
			RawDescriptor: file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_goTypes,
		DependencyIndexes: file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_depIdxs,
		MessageInfos:      file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_msgTypes,
	}.Build()
	File_aperture_flowcontrol_controlpoints_v1_controlpoints_proto = out.File
	file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_rawDesc = nil
	file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_goTypes = nil
	file_aperture_flowcontrol_controlpoints_v1_controlpoints_proto_depIdxs = nil
}