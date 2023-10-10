// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: aperture/cloud/v1/blueprints.proto

package cloudv1

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

type Blueprint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Blueprint) Reset() {
	*x = Blueprint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Blueprint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Blueprint) ProtoMessage() {}

func (x *Blueprint) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Blueprint.ProtoReflect.Descriptor instead.
func (*Blueprint) Descriptor() ([]byte, []int) {
	return file_aperture_cloud_v1_blueprints_proto_rawDescGZIP(), []int{0}
}

func (x *Blueprint) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blueprints []*Blueprint `protobuf:"bytes,1,rep,name=blueprints,proto3" json:"blueprints,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_aperture_cloud_v1_blueprints_proto_rawDescGZIP(), []int{1}
}

func (x *ListResponse) GetBlueprints() []*Blueprint {
	if x != nil {
		return x.Blueprints
	}
	return nil
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_aperture_cloud_v1_blueprints_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blueprint *Blueprint `protobuf:"bytes,1,opt,name=blueprint,proto3" json:"blueprint,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_aperture_cloud_v1_blueprints_proto_rawDescGZIP(), []int{3}
}

func (x *GetResponse) GetBlueprint() *Blueprint {
	if x != nil {
		return x.Blueprint
	}
	return nil
}

type ApplyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Blueprint *Blueprint `protobuf:"bytes,2,opt,name=blueprint,proto3" json:"blueprint,omitempty"`
}

func (x *ApplyRequest) Reset() {
	*x = ApplyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplyRequest) ProtoMessage() {}

func (x *ApplyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplyRequest.ProtoReflect.Descriptor instead.
func (*ApplyRequest) Descriptor() ([]byte, []int) {
	return file_aperture_cloud_v1_blueprints_proto_rawDescGZIP(), []int{4}
}

func (x *ApplyRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ApplyRequest) GetBlueprint() *Blueprint {
	if x != nil {
		return x.Blueprint
	}
	return nil
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_cloud_v1_blueprints_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_aperture_cloud_v1_blueprints_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_aperture_cloud_v1_blueprints_proto protoreflect.FileDescriptor

var file_aperture_cloud_v1_blueprints_proto_rawDesc = []byte{
	0x0a, 0x22, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x2f, 0x76, 0x31, 0x2f, 0x62, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x25, 0x0a, 0x09, 0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x4c, 0x0a, 0x0c, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0a, 0x62, 0x6c, 0x75,
	0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76,
	0x31, 0x2e, 0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x52, 0x0a, 0x62, 0x6c, 0x75,
	0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x20, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x49, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x09, 0x62, 0x6c, 0x75, 0x65,
	0x70, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x2e,
	0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x52, 0x09, 0x62, 0x6c, 0x75, 0x65, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x22, 0x5e, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x09, 0x62, 0x6c, 0x75, 0x65,
	0x70, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x2e,
	0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x52, 0x09, 0x62, 0x6c, 0x75, 0x65, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x22, 0x23, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0xbd, 0x03, 0x0a, 0x11, 0x42, 0x6c,
	0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x61, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x1f, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x12, 0x18, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e,
	0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e,
	0x74, 0x73, 0x12, 0x6d, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1d, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21,
	0x12, 0x1f, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x2f,
	0x62, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65,
	0x7d, 0x12, 0x69, 0x0a, 0x05, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x12, 0x1f, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x70, 0x70, 0x6c, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x22, 0x1f, 0x2f, 0x66, 0x6c,
	0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x6c, 0x75, 0x65, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12, 0x6b, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x20, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x2a, 0x1f, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e,
	0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69, 0x6e,
	0x74, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x42, 0xef, 0x01, 0x0a, 0x29, 0x63, 0x6f,
	0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x76, 0x31, 0x42, 0x0f, 0x42, 0x6c, 0x75, 0x65, 0x70, 0x72, 0x69,
	0x6e, 0x74, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61,
	0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x76, 0x31, 0x3b,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x43, 0x58, 0xaa, 0x02, 0x11,
	0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x11, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x5c, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x3a, 0x3a, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_aperture_cloud_v1_blueprints_proto_rawDescOnce sync.Once
	file_aperture_cloud_v1_blueprints_proto_rawDescData = file_aperture_cloud_v1_blueprints_proto_rawDesc
)

func file_aperture_cloud_v1_blueprints_proto_rawDescGZIP() []byte {
	file_aperture_cloud_v1_blueprints_proto_rawDescOnce.Do(func() {
		file_aperture_cloud_v1_blueprints_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_cloud_v1_blueprints_proto_rawDescData)
	})
	return file_aperture_cloud_v1_blueprints_proto_rawDescData
}

var file_aperture_cloud_v1_blueprints_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_aperture_cloud_v1_blueprints_proto_goTypes = []interface{}{
	(*Blueprint)(nil),     // 0: aperture.cloud.v1.Blueprint
	(*ListResponse)(nil),  // 1: aperture.cloud.v1.ListResponse
	(*GetRequest)(nil),    // 2: aperture.cloud.v1.GetRequest
	(*GetResponse)(nil),   // 3: aperture.cloud.v1.GetResponse
	(*ApplyRequest)(nil),  // 4: aperture.cloud.v1.ApplyRequest
	(*DeleteRequest)(nil), // 5: aperture.cloud.v1.DeleteRequest
	(*emptypb.Empty)(nil), // 6: google.protobuf.Empty
}
var file_aperture_cloud_v1_blueprints_proto_depIdxs = []int32{
	0, // 0: aperture.cloud.v1.ListResponse.blueprints:type_name -> aperture.cloud.v1.Blueprint
	0, // 1: aperture.cloud.v1.GetResponse.blueprint:type_name -> aperture.cloud.v1.Blueprint
	0, // 2: aperture.cloud.v1.ApplyRequest.blueprint:type_name -> aperture.cloud.v1.Blueprint
	6, // 3: aperture.cloud.v1.BlueprintsService.List:input_type -> google.protobuf.Empty
	2, // 4: aperture.cloud.v1.BlueprintsService.Get:input_type -> aperture.cloud.v1.GetRequest
	4, // 5: aperture.cloud.v1.BlueprintsService.Apply:input_type -> aperture.cloud.v1.ApplyRequest
	5, // 6: aperture.cloud.v1.BlueprintsService.Delete:input_type -> aperture.cloud.v1.DeleteRequest
	1, // 7: aperture.cloud.v1.BlueprintsService.List:output_type -> aperture.cloud.v1.ListResponse
	3, // 8: aperture.cloud.v1.BlueprintsService.Get:output_type -> aperture.cloud.v1.GetResponse
	6, // 9: aperture.cloud.v1.BlueprintsService.Apply:output_type -> google.protobuf.Empty
	6, // 10: aperture.cloud.v1.BlueprintsService.Delete:output_type -> google.protobuf.Empty
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_aperture_cloud_v1_blueprints_proto_init() }
func file_aperture_cloud_v1_blueprints_proto_init() {
	if File_aperture_cloud_v1_blueprints_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_cloud_v1_blueprints_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Blueprint); i {
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
		file_aperture_cloud_v1_blueprints_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResponse); i {
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
		file_aperture_cloud_v1_blueprints_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_aperture_cloud_v1_blueprints_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
		file_aperture_cloud_v1_blueprints_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplyRequest); i {
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
		file_aperture_cloud_v1_blueprints_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
			RawDescriptor: file_aperture_cloud_v1_blueprints_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_aperture_cloud_v1_blueprints_proto_goTypes,
		DependencyIndexes: file_aperture_cloud_v1_blueprints_proto_depIdxs,
		MessageInfos:      file_aperture_cloud_v1_blueprints_proto_msgTypes,
	}.Build()
	File_aperture_cloud_v1_blueprints_proto = out.File
	file_aperture_cloud_v1_blueprints_proto_rawDesc = nil
	file_aperture_cloud_v1_blueprints_proto_goTypes = nil
	file_aperture_cloud_v1_blueprints_proto_depIdxs = nil
}
