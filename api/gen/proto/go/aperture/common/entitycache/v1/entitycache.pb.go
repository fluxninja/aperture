// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/common/entitycache/v1/entitycache.proto

package entitycachev1

import (
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

type ServicesList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Services            []*Service            `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	OverlappingServices []*OverlappingService `protobuf:"bytes,2,rep,name=overlapping_services,json=overlappingServices,proto3" json:"overlapping_services,omitempty"`
}

func (x *ServicesList) Reset() {
	*x = ServicesList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServicesList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServicesList) ProtoMessage() {}

func (x *ServicesList) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServicesList.ProtoReflect.Descriptor instead.
func (*ServicesList) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{0}
}

func (x *ServicesList) GetServices() []*Service {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *ServicesList) GetOverlappingServices() []*OverlappingService {
	if x != nil {
		return x.OverlappingServices
	}
	return nil
}

// Service contains information about single service discovered in agent group by a
// particular agent
type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	EntitiesCount int32  `protobuf:"varint,2,opt,name=entities_count,json=entitiesCount,proto3" json:"entities_count,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{1}
}

func (x *Service) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Service) GetEntitiesCount() int32 {
	if x != nil {
		return x.EntitiesCount
	}
	return 0
}

// OverlappingService contains info about a service that overlaps with another one
type OverlappingService struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service1      string `protobuf:"bytes,1,opt,name=service1,proto3" json:"service1,omitempty"`
	Service2      string `protobuf:"bytes,2,opt,name=service2,proto3" json:"service2,omitempty"`
	EntitiesCount int32  `protobuf:"varint,3,opt,name=entities_count,json=entitiesCount,proto3" json:"entities_count,omitempty"`
}

func (x *OverlappingService) Reset() {
	*x = OverlappingService{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverlappingService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverlappingService) ProtoMessage() {}

func (x *OverlappingService) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverlappingService.ProtoReflect.Descriptor instead.
func (*OverlappingService) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{2}
}

func (x *OverlappingService) GetService1() string {
	if x != nil {
		return x.Service1
	}
	return ""
}

func (x *OverlappingService) GetService2() string {
	if x != nil {
		return x.Service2
	}
	return ""
}

func (x *OverlappingService) GetEntitiesCount() int32 {
	if x != nil {
		return x.EntitiesCount
	}
	return 0
}

var File_aperture_common_entitycache_v1_entitycache_proto protoreflect.FileDescriptor

var file_aperture_common_entitycache_v1_entitycache_proto_rawDesc = []byte{
	0x0a, 0x30, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xba, 0x01,
	0x0a, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x43,
	0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x12, 0x65, 0x0a, 0x14, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69,
	0x6e, 0x67, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x32, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x13, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69,
	0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x44, 0x0a, 0x07, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0x73, 0x0a, 0x12, 0x4f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x31, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x12, 0x25,
	0x0a, 0x0e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x85, 0x01, 0x0a, 0x0f, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x72, 0x0a, 0x0f, 0x47, 0x65, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x4c, 0x69,
	0x73, 0x74, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x76, 0x31, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2d, 0x6c, 0x69, 0x73, 0x74, 0x42, 0xae, 0x02,
	0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63,
	0x61, 0x63, 0x68, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61,
	0x63, 0x68, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x43, 0x45, 0xaa, 0x02, 0x1e, 0x41, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1e, 0x41,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x2a,
	0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x21, 0x41, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_common_entitycache_v1_entitycache_proto_rawDescOnce sync.Once
	file_aperture_common_entitycache_v1_entitycache_proto_rawDescData = file_aperture_common_entitycache_v1_entitycache_proto_rawDesc
)

func file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP() []byte {
	file_aperture_common_entitycache_v1_entitycache_proto_rawDescOnce.Do(func() {
		file_aperture_common_entitycache_v1_entitycache_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_common_entitycache_v1_entitycache_proto_rawDescData)
	})
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescData
}

var file_aperture_common_entitycache_v1_entitycache_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_aperture_common_entitycache_v1_entitycache_proto_goTypes = []interface{}{
	(*ServicesList)(nil),       // 0: aperture.common.entitycache.v1.ServicesList
	(*Service)(nil),            // 1: aperture.common.entitycache.v1.Service
	(*OverlappingService)(nil), // 2: aperture.common.entitycache.v1.OverlappingService
	(*emptypb.Empty)(nil),      // 3: google.protobuf.Empty
}
var file_aperture_common_entitycache_v1_entitycache_proto_depIdxs = []int32{
	1, // 0: aperture.common.entitycache.v1.ServicesList.services:type_name -> aperture.common.entitycache.v1.Service
	2, // 1: aperture.common.entitycache.v1.ServicesList.overlapping_services:type_name -> aperture.common.entitycache.v1.OverlappingService
	3, // 2: aperture.common.entitycache.v1.EntitiesService.GetServicesList:input_type -> google.protobuf.Empty
	0, // 3: aperture.common.entitycache.v1.EntitiesService.GetServicesList:output_type -> aperture.common.entitycache.v1.ServicesList
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_aperture_common_entitycache_v1_entitycache_proto_init() }
func file_aperture_common_entitycache_v1_entitycache_proto_init() {
	if File_aperture_common_entitycache_v1_entitycache_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServicesList); i {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service); i {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OverlappingService); i {
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
			RawDescriptor: file_aperture_common_entitycache_v1_entitycache_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_aperture_common_entitycache_v1_entitycache_proto_goTypes,
		DependencyIndexes: file_aperture_common_entitycache_v1_entitycache_proto_depIdxs,
		MessageInfos:      file_aperture_common_entitycache_v1_entitycache_proto_msgTypes,
	}.Build()
	File_aperture_common_entitycache_v1_entitycache_proto = out.File
	file_aperture_common_entitycache_v1_entitycache_proto_rawDesc = nil
	file_aperture_common_entitycache_v1_entitycache_proto_goTypes = nil
	file_aperture_common_entitycache_v1_entitycache_proto_depIdxs = nil
}
