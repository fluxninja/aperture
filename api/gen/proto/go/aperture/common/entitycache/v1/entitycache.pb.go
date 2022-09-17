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

type GetEntityByIPAddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IpAddress string `protobuf:"bytes,1,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
}

func (x *GetEntityByIPAddressRequest) Reset() {
	*x = GetEntityByIPAddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEntityByIPAddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEntityByIPAddressRequest) ProtoMessage() {}

func (x *GetEntityByIPAddressRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetEntityByIPAddressRequest.ProtoReflect.Descriptor instead.
func (*GetEntityByIPAddressRequest) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{0}
}

func (x *GetEntityByIPAddressRequest) GetIpAddress() string {
	if x != nil {
		return x.IpAddress
	}
	return ""
}

type GetEntityByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetEntityByNameRequest) Reset() {
	*x = GetEntityByNameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEntityByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEntityByNameRequest) ProtoMessage() {}

func (x *GetEntityByNameRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetEntityByNameRequest.ProtoReflect.Descriptor instead.
func (*GetEntityByNameRequest) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{1}
}

func (x *GetEntityByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

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
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServicesList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServicesList) ProtoMessage() {}

func (x *ServicesList) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ServicesList.ProtoReflect.Descriptor instead.
func (*ServicesList) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{2}
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
// particular agent.
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
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[3]
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
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{3}
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

// OverlappingService contains info about a service that overlaps with another one.
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
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverlappingService) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverlappingService) ProtoMessage() {}

func (x *OverlappingService) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[4]
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
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{4}
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

// EntityCache contains both mappings of ip address to entity and entity name to entity.
type EntityCache struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntitiesByIpAddress  *EntityCache_Entities `protobuf:"bytes,1,opt,name=entities_by_ip_address,json=entitiesByIpAddress,proto3" json:"entities_by_ip_address,omitempty"`
	EntitiesByEntityName *EntityCache_Entities `protobuf:"bytes,2,opt,name=entities_by_entity_name,json=entitiesByEntityName,proto3" json:"entities_by_entity_name,omitempty"`
}

func (x *EntityCache) Reset() {
	*x = EntityCache{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityCache) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityCache) ProtoMessage() {}

func (x *EntityCache) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityCache.ProtoReflect.Descriptor instead.
func (*EntityCache) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{5}
}

func (x *EntityCache) GetEntitiesByIpAddress() *EntityCache_Entities {
	if x != nil {
		return x.EntitiesByIpAddress
	}
	return nil
}

func (x *EntityCache) GetEntitiesByEntityName() *EntityCache_Entities {
	if x != nil {
		return x.EntitiesByEntityName
	}
	return nil
}

// Entity represents a pod, vm, etc.
type Entity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityId   *Entity_EntityID `protobuf:"bytes,1,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
	IpAddress  string           `protobuf:"bytes,2,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	Services   []string         `protobuf:"bytes,3,rep,name=services,proto3" json:"services,omitempty"`
	EntityName string           `protobuf:"bytes,4,opt,name=entity_name,json=entityName,proto3" json:"entity_name,omitempty"`
}

func (x *Entity) Reset() {
	*x = Entity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{6}
}

func (x *Entity) GetEntityId() *Entity_EntityID {
	if x != nil {
		return x.EntityId
	}
	return nil
}

func (x *Entity) GetIpAddress() string {
	if x != nil {
		return x.IpAddress
	}
	return ""
}

func (x *Entity) GetServices() []string {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *Entity) GetEntityName() string {
	if x != nil {
		return x.EntityName
	}
	return ""
}

// Entities defines mapping of entities.
type EntityCache_Entities struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entities map[string]*Entity `protobuf:"bytes,1,rep,name=entities,proto3" json:"entities,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EntityCache_Entities) Reset() {
	*x = EntityCache_Entities{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityCache_Entities) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityCache_Entities) ProtoMessage() {}

func (x *EntityCache_Entities) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityCache_Entities.ProtoReflect.Descriptor instead.
func (*EntityCache_Entities) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{5, 0}
}

func (x *EntityCache_Entities) GetEntities() map[string]*Entity {
	if x != nil {
		return x.Entities
	}
	return nil
}

// EntityID is a unique Entity identifier.
type Entity_EntityID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prefix string `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Uid    string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *Entity_EntityID) Reset() {
	*x = Entity_EntityID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity_EntityID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity_EntityID) ProtoMessage() {}

func (x *Entity_EntityID) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity_EntityID.ProtoReflect.Descriptor instead.
func (*Entity_EntityID) Descriptor() ([]byte, []int) {
	return file_aperture_common_entitycache_v1_entitycache_proto_rawDescGZIP(), []int{6, 0}
}

func (x *Entity_EntityID) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *Entity_EntityID) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
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
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a,
	0x1b, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x42, 0x79, 0x49, 0x50, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x2c, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xba, 0x01, 0x0a, 0x0c, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x43, 0x0a, 0x08, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12,
	0x65, 0x0a, 0x14, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4f,
	0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x13, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x44, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x73, 0x0a, 0x12,
	0x4f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x31, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x31, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0d, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0xb7, 0x03, 0x0a, 0x0b, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x61, 0x63, 0x68,
	0x65, 0x12, 0x69, 0x0a, 0x16, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x62, 0x79,
	0x5f, 0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x34, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x13, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x42, 0x79, 0x49, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x6b, 0x0a, 0x17,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x62, 0x79, 0x5f, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x52, 0x14, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x42, 0x79, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0xcf, 0x01, 0x0a, 0x08, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x5e, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x43, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x1a, 0x63, 0x0a, 0x0d, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3c, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xe8, 0x01, 0x0a, 0x06,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x4c, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x44, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x1a, 0x34, 0x0a, 0x08, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72,
	0x65, 0x66, 0x69, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x32, 0xbb, 0x04, 0x0a, 0x12, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x43, 0x61, 0x63, 0x68, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x76, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15,
	0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x6b, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x43, 0x61, 0x63, 0x68, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x2b, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x43, 0x61, 0x63, 0x68, 0x65, 0x22, 0x14, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x12, 0xa9, 0x01, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x42, 0x79, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3b, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x42, 0x79, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x22, 0x2c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26, 0x12, 0x24, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x69, 0x70, 0x2d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x2f, 0x7b, 0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x7d, 0x12, 0x93,
	0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x42, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x36, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x42, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x12, 0x18, 0x2f, 0x76, 0x31, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x6e, 0x61, 0x6d, 0x65, 0x2f, 0x7b, 0x6e,
	0x61, 0x6d, 0x65, 0x7d, 0x42, 0xae, 0x02, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41,
	0x43, 0x45, 0xaa, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x2a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x21, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_aperture_common_entitycache_v1_entitycache_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_aperture_common_entitycache_v1_entitycache_proto_goTypes = []interface{}{
	(*GetEntityByIPAddressRequest)(nil), // 0: aperture.common.entitycache.v1.GetEntityByIPAddressRequest
	(*GetEntityByNameRequest)(nil),      // 1: aperture.common.entitycache.v1.GetEntityByNameRequest
	(*ServicesList)(nil),                // 2: aperture.common.entitycache.v1.ServicesList
	(*Service)(nil),                     // 3: aperture.common.entitycache.v1.Service
	(*OverlappingService)(nil),          // 4: aperture.common.entitycache.v1.OverlappingService
	(*EntityCache)(nil),                 // 5: aperture.common.entitycache.v1.EntityCache
	(*Entity)(nil),                      // 6: aperture.common.entitycache.v1.Entity
	(*EntityCache_Entities)(nil),        // 7: aperture.common.entitycache.v1.EntityCache.Entities
	nil,                                 // 8: aperture.common.entitycache.v1.EntityCache.Entities.EntitiesEntry
	(*Entity_EntityID)(nil),             // 9: aperture.common.entitycache.v1.Entity.EntityID
	(*emptypb.Empty)(nil),               // 10: google.protobuf.Empty
}
var file_aperture_common_entitycache_v1_entitycache_proto_depIdxs = []int32{
	3,  // 0: aperture.common.entitycache.v1.ServicesList.services:type_name -> aperture.common.entitycache.v1.Service
	4,  // 1: aperture.common.entitycache.v1.ServicesList.overlapping_services:type_name -> aperture.common.entitycache.v1.OverlappingService
	7,  // 2: aperture.common.entitycache.v1.EntityCache.entities_by_ip_address:type_name -> aperture.common.entitycache.v1.EntityCache.Entities
	7,  // 3: aperture.common.entitycache.v1.EntityCache.entities_by_entity_name:type_name -> aperture.common.entitycache.v1.EntityCache.Entities
	9,  // 4: aperture.common.entitycache.v1.Entity.entity_id:type_name -> aperture.common.entitycache.v1.Entity.EntityID
	8,  // 5: aperture.common.entitycache.v1.EntityCache.Entities.entities:type_name -> aperture.common.entitycache.v1.EntityCache.Entities.EntitiesEntry
	6,  // 6: aperture.common.entitycache.v1.EntityCache.Entities.EntitiesEntry.value:type_name -> aperture.common.entitycache.v1.Entity
	10, // 7: aperture.common.entitycache.v1.EntityCacheService.GetServicesList:input_type -> google.protobuf.Empty
	10, // 8: aperture.common.entitycache.v1.EntityCacheService.GetEntityCache:input_type -> google.protobuf.Empty
	0,  // 9: aperture.common.entitycache.v1.EntityCacheService.GetEntityByIPAddress:input_type -> aperture.common.entitycache.v1.GetEntityByIPAddressRequest
	1,  // 10: aperture.common.entitycache.v1.EntityCacheService.GetEntityByName:input_type -> aperture.common.entitycache.v1.GetEntityByNameRequest
	2,  // 11: aperture.common.entitycache.v1.EntityCacheService.GetServicesList:output_type -> aperture.common.entitycache.v1.ServicesList
	5,  // 12: aperture.common.entitycache.v1.EntityCacheService.GetEntityCache:output_type -> aperture.common.entitycache.v1.EntityCache
	6,  // 13: aperture.common.entitycache.v1.EntityCacheService.GetEntityByIPAddress:output_type -> aperture.common.entitycache.v1.Entity
	6,  // 14: aperture.common.entitycache.v1.EntityCacheService.GetEntityByName:output_type -> aperture.common.entitycache.v1.Entity
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_aperture_common_entitycache_v1_entitycache_proto_init() }
func file_aperture_common_entitycache_v1_entitycache_proto_init() {
	if File_aperture_common_entitycache_v1_entitycache_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEntityByIPAddressRequest); i {
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
			switch v := v.(*GetEntityByNameRequest); i {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntityCache); i {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity); i {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntityCache_Entities); i {
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
		file_aperture_common_entitycache_v1_entitycache_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity_EntityID); i {
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
			NumMessages:   10,
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
