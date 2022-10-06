// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/policy/wrappers/v1/policy.proto

package wrappersv1

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

type PolicyWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Policy
	Policy *v1.Policy `protobuf:"bytes,2,opt,name=policy,proto3" json:"policy,omitempty"`
}

func (x *PolicyWrapper) Reset() {
	*x = PolicyWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PolicyWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PolicyWrapper) ProtoMessage() {}

func (x *PolicyWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PolicyWrapper.ProtoReflect.Descriptor instead.
func (*PolicyWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_wrappers_v1_policy_proto_rawDescGZIP(), []int{0}
}

func (x *PolicyWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *PolicyWrapper) GetPolicy() *v1.Policy {
	if x != nil {
		return x.Policy
	}
	return nil
}

type FluxMeterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Flux Meter
	FluxMeter *v1.FluxMeter `protobuf:"bytes,1,opt,name=flux_meter,json=fluxMeter,proto3" json:"flux_meter,omitempty"`
	// Name of Flux Meter metric.
	FluxMeterName string `protobuf:"bytes,4,opt,name=flux_meter_name,json=fluxMeterName,proto3" json:"flux_meter_name,omitempty"`
}

func (x *FluxMeterWrapper) Reset() {
	*x = FluxMeterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FluxMeterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FluxMeterWrapper) ProtoMessage() {}

func (x *FluxMeterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FluxMeterWrapper.ProtoReflect.Descriptor instead.
func (*FluxMeterWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_wrappers_v1_policy_proto_rawDescGZIP(), []int{1}
}

func (x *FluxMeterWrapper) GetFluxMeter() *v1.FluxMeter {
	if x != nil {
		return x.FluxMeter
	}
	return nil
}

func (x *FluxMeterWrapper) GetFluxMeterName() string {
	if x != nil {
		return x.FluxMeterName
	}
	return ""
}

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
		mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassifierWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassifierWrapper) ProtoMessage() {}

func (x *ClassifierWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[2]
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
	return file_aperture_policy_wrappers_v1_policy_proto_rawDescGZIP(), []int{2}
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

type ConcurrencyLimiterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Concurrency Limiter
	ConcurrencyLimiter *v1.ConcurrencyLimiter `protobuf:"bytes,2,opt,name=concurrency_limiter,json=concurrencyLimiter,proto3" json:"concurrency_limiter,omitempty"`
}

func (x *ConcurrencyLimiterWrapper) Reset() {
	*x = ConcurrencyLimiterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcurrencyLimiterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcurrencyLimiterWrapper) ProtoMessage() {}

func (x *ConcurrencyLimiterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConcurrencyLimiterWrapper.ProtoReflect.Descriptor instead.
func (*ConcurrencyLimiterWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_wrappers_v1_policy_proto_rawDescGZIP(), []int{3}
}

func (x *ConcurrencyLimiterWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *ConcurrencyLimiterWrapper) GetConcurrencyLimiter() *v1.ConcurrencyLimiter {
	if x != nil {
		return x.ConcurrencyLimiter
	}
	return nil
}

type RateLimiterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Rate Limiter
	RateLimiter *v1.RateLimiter `protobuf:"bytes,2,opt,name=rate_limiter,json=rateLimiter,proto3" json:"rate_limiter,omitempty"`
}

func (x *RateLimiterWrapper) Reset() {
	*x = RateLimiterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterWrapper) ProtoMessage() {}

func (x *RateLimiterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_wrappers_v1_policy_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimiterWrapper.ProtoReflect.Descriptor instead.
func (*RateLimiterWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_wrappers_v1_policy_proto_rawDescGZIP(), []int{4}
}

func (x *RateLimiterWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *RateLimiterWrapper) GetRateLimiter() *v1.RateLimiter {
	if x != nil {
		return x.RateLimiter
	}
	return nil
}

var File_aperture_policy_wrappers_v1_policy_proto protoreflect.FileDescriptor

var file_aperture_policy_wrappers_v1_policy_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x77, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f,
	0x76, 0x31, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x28, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa8, 0x01, 0x0a, 0x0d, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x5a, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x22, 0x81, 0x01, 0x0a, 0x10, 0x46, 0x6c, 0x75, 0x78, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x57,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x0a, 0x66, 0x6c, 0x75, 0x78, 0x5f, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x75, 0x78, 0x4d, 0x65, 0x74,
	0x65, 0x72, 0x52, 0x09, 0x66, 0x6c, 0x75, 0x78, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x26, 0x0a,
	0x0f, 0x66, 0x6c, 0x75, 0x78, 0x5f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66, 0x6c, 0x75, 0x78, 0x4d, 0x65, 0x74, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xb8, 0x01, 0x0a, 0x11, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x5a, 0x0a, 0x11, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x47, 0x0a, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x22, 0xd9, 0x01, 0x0a, 0x19, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79,
	0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x5a,
	0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x77, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x60, 0x0a, 0x13, 0x63, 0x6f,
	0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63,
	0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x22, 0xbd, 0x01, 0x0a,
	0x12, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x12, 0x5a, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d,
	0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2e, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12,
	0x4b, 0x0a, 0x0c, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x52,
	0x0b, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x42, 0x94, 0x02, 0x0a,
	0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2e, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31,
	0x42, 0x0b, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x55, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f,
	0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x57, 0xaa, 0x02, 0x1b, 0x41,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x57,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1b, 0x41, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x57, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x57, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_wrappers_v1_policy_proto_rawDescOnce sync.Once
	file_aperture_policy_wrappers_v1_policy_proto_rawDescData = file_aperture_policy_wrappers_v1_policy_proto_rawDesc
)

func file_aperture_policy_wrappers_v1_policy_proto_rawDescGZIP() []byte {
	file_aperture_policy_wrappers_v1_policy_proto_rawDescOnce.Do(func() {
		file_aperture_policy_wrappers_v1_policy_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_wrappers_v1_policy_proto_rawDescData)
	})
	return file_aperture_policy_wrappers_v1_policy_proto_rawDescData
}

var file_aperture_policy_wrappers_v1_policy_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_aperture_policy_wrappers_v1_policy_proto_goTypes = []interface{}{
	(*PolicyWrapper)(nil),             // 0: aperture.policy.wrappers.v1.PolicyWrapper
	(*FluxMeterWrapper)(nil),          // 1: aperture.policy.wrappers.v1.FluxMeterWrapper
	(*ClassifierWrapper)(nil),         // 2: aperture.policy.wrappers.v1.ClassifierWrapper
	(*ConcurrencyLimiterWrapper)(nil), // 3: aperture.policy.wrappers.v1.ConcurrencyLimiterWrapper
	(*RateLimiterWrapper)(nil),        // 4: aperture.policy.wrappers.v1.RateLimiterWrapper
	(*CommonAttributes)(nil),          // 5: aperture.policy.wrappers.v1.CommonAttributes
	(*v1.Policy)(nil),                 // 6: aperture.policy.language.v1.Policy
	(*v1.FluxMeter)(nil),              // 7: aperture.policy.language.v1.FluxMeter
	(*v1.Classifier)(nil),             // 8: aperture.policy.language.v1.Classifier
	(*v1.ConcurrencyLimiter)(nil),     // 9: aperture.policy.language.v1.ConcurrencyLimiter
	(*v1.RateLimiter)(nil),            // 10: aperture.policy.language.v1.RateLimiter
}
var file_aperture_policy_wrappers_v1_policy_proto_depIdxs = []int32{
	5,  // 0: aperture.policy.wrappers.v1.PolicyWrapper.common_attributes:type_name -> aperture.policy.wrappers.v1.CommonAttributes
	6,  // 1: aperture.policy.wrappers.v1.PolicyWrapper.policy:type_name -> aperture.policy.language.v1.Policy
	7,  // 2: aperture.policy.wrappers.v1.FluxMeterWrapper.flux_meter:type_name -> aperture.policy.language.v1.FluxMeter
	5,  // 3: aperture.policy.wrappers.v1.ClassifierWrapper.common_attributes:type_name -> aperture.policy.wrappers.v1.CommonAttributes
	8,  // 4: aperture.policy.wrappers.v1.ClassifierWrapper.classifier:type_name -> aperture.policy.language.v1.Classifier
	5,  // 5: aperture.policy.wrappers.v1.ConcurrencyLimiterWrapper.common_attributes:type_name -> aperture.policy.wrappers.v1.CommonAttributes
	9,  // 6: aperture.policy.wrappers.v1.ConcurrencyLimiterWrapper.concurrency_limiter:type_name -> aperture.policy.language.v1.ConcurrencyLimiter
	5,  // 7: aperture.policy.wrappers.v1.RateLimiterWrapper.common_attributes:type_name -> aperture.policy.wrappers.v1.CommonAttributes
	10, // 8: aperture.policy.wrappers.v1.RateLimiterWrapper.rate_limiter:type_name -> aperture.policy.language.v1.RateLimiter
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_aperture_policy_wrappers_v1_policy_proto_init() }
func file_aperture_policy_wrappers_v1_policy_proto_init() {
	if File_aperture_policy_wrappers_v1_policy_proto != nil {
		return
	}
	file_aperture_policy_wrappers_v1_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_wrappers_v1_policy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PolicyWrapper); i {
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
		file_aperture_policy_wrappers_v1_policy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FluxMeterWrapper); i {
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
		file_aperture_policy_wrappers_v1_policy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_policy_wrappers_v1_policy_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConcurrencyLimiterWrapper); i {
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
		file_aperture_policy_wrappers_v1_policy_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimiterWrapper); i {
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
			RawDescriptor: file_aperture_policy_wrappers_v1_policy_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_wrappers_v1_policy_proto_goTypes,
		DependencyIndexes: file_aperture_policy_wrappers_v1_policy_proto_depIdxs,
		MessageInfos:      file_aperture_policy_wrappers_v1_policy_proto_msgTypes,
	}.Build()
	File_aperture_policy_wrappers_v1_policy_proto = out.File
	file_aperture_policy_wrappers_v1_policy_proto_rawDesc = nil
	file_aperture_policy_wrappers_v1_policy_proto_goTypes = nil
	file_aperture_policy_wrappers_v1_policy_proto_depIdxs = nil
}
