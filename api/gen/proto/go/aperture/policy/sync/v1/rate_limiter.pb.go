// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: aperture/policy/sync/v1/rate_limiter.proto

package syncv1

import (
	v1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

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
		mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterWrapper) ProtoMessage() {}

func (x *RateLimiterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[0]
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
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP(), []int{0}
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

type RateLimiterDecisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Rate Limiter Decision
	RateLimiterDecision *RateLimiterDecision `protobuf:"bytes,2,opt,name=rate_limiter_decision,json=rateLimiterDecision,proto3" json:"rate_limiter_decision,omitempty"`
	// RateLimiterDynamicConfig is the dynamic configuration for the rate limiter.
	RateLimiterDynamicConfig *v1.RateLimiter_DynamicConfig `protobuf:"bytes,3,opt,name=rate_limiter_dynamic_config,json=rateLimiterDynamicConfig,proto3" json:"rate_limiter_dynamic_config,omitempty"`
}

func (x *RateLimiterDecisionWrapper) Reset() {
	*x = RateLimiterDecisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterDecisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterDecisionWrapper) ProtoMessage() {}

func (x *RateLimiterDecisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimiterDecisionWrapper.ProtoReflect.Descriptor instead.
func (*RateLimiterDecisionWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP(), []int{1}
}

func (x *RateLimiterDecisionWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *RateLimiterDecisionWrapper) GetRateLimiterDecision() *RateLimiterDecision {
	if x != nil {
		return x.RateLimiterDecision
	}
	return nil
}

func (x *RateLimiterDecisionWrapper) GetRateLimiterDynamicConfig() *v1.RateLimiter_DynamicConfig {
	if x != nil {
		return x.RateLimiterDynamicConfig
	}
	return nil
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
		mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterDecision) ProtoMessage() {}

func (x *RateLimiterDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[2]
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
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP(), []int{2}
}

func (x *RateLimiterDecision) GetLimit() float64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type LeakyBucketRateLimiterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Rate Limiter
	RateLimiter *v1.LeakyBucketRateLimiter `protobuf:"bytes,2,opt,name=rate_limiter,json=rateLimiter,proto3" json:"rate_limiter,omitempty"`
}

func (x *LeakyBucketRateLimiterWrapper) Reset() {
	*x = LeakyBucketRateLimiterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeakyBucketRateLimiterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeakyBucketRateLimiterWrapper) ProtoMessage() {}

func (x *LeakyBucketRateLimiterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeakyBucketRateLimiterWrapper.ProtoReflect.Descriptor instead.
func (*LeakyBucketRateLimiterWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP(), []int{3}
}

func (x *LeakyBucketRateLimiterWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *LeakyBucketRateLimiterWrapper) GetRateLimiter() *v1.LeakyBucketRateLimiter {
	if x != nil {
		return x.RateLimiter
	}
	return nil
}

type LeakyBucketRateLimiterDecisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CommonAttributes
	CommonAttributes *CommonAttributes `protobuf:"bytes,1,opt,name=common_attributes,json=commonAttributes,proto3" json:"common_attributes,omitempty"`
	// Rate Limiter Decision
	RateLimiterDecision *LeakyBucketRateLimiterDecision `protobuf:"bytes,2,opt,name=rate_limiter_decision,json=rateLimiterDecision,proto3" json:"rate_limiter_decision,omitempty"`
}

func (x *LeakyBucketRateLimiterDecisionWrapper) Reset() {
	*x = LeakyBucketRateLimiterDecisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeakyBucketRateLimiterDecisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeakyBucketRateLimiterDecisionWrapper) ProtoMessage() {}

func (x *LeakyBucketRateLimiterDecisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeakyBucketRateLimiterDecisionWrapper.ProtoReflect.Descriptor instead.
func (*LeakyBucketRateLimiterDecisionWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP(), []int{4}
}

func (x *LeakyBucketRateLimiterDecisionWrapper) GetCommonAttributes() *CommonAttributes {
	if x != nil {
		return x.CommonAttributes
	}
	return nil
}

func (x *LeakyBucketRateLimiterDecisionWrapper) GetRateLimiterDecision() *LeakyBucketRateLimiterDecision {
	if x != nil {
		return x.RateLimiterDecision
	}
	return nil
}

type LeakyBucketRateLimiterDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BucketCapacity float64              `protobuf:"fixed64,1,opt,name=bucket_capacity,json=bucketCapacity,proto3" json:"bucket_capacity,omitempty"`
	LeakInterval   *durationpb.Duration `protobuf:"bytes,2,opt,name=leak_interval,json=leakInterval,proto3" json:"leak_interval,omitempty"`
	LeakAmount     float64              `protobuf:"fixed64,3,opt,name=leak_amount,json=leakAmount,proto3" json:"leak_amount,omitempty"`
}

func (x *LeakyBucketRateLimiterDecision) Reset() {
	*x = LeakyBucketRateLimiterDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeakyBucketRateLimiterDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeakyBucketRateLimiterDecision) ProtoMessage() {}

func (x *LeakyBucketRateLimiterDecision) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeakyBucketRateLimiterDecision.ProtoReflect.Descriptor instead.
func (*LeakyBucketRateLimiterDecision) Descriptor() ([]byte, []int) {
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP(), []int{5}
}

func (x *LeakyBucketRateLimiterDecision) GetBucketCapacity() float64 {
	if x != nil {
		return x.BucketCapacity
	}
	return 0
}

func (x *LeakyBucketRateLimiterDecision) GetLeakInterval() *durationpb.Duration {
	if x != nil {
		return x.LeakInterval
	}
	return nil
}

func (x *LeakyBucketRateLimiterDecision) GetLeakAmount() float64 {
	if x != nil {
		return x.LeakAmount
	}
	return 0
}

var File_aperture_policy_sync_v1_rate_limiter_proto protoreflect.FileDescriptor

var file_aperture_policy_sync_v1_rate_limiter_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x2d, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f,
	0x76, 0x31, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb9, 0x01, 0x0a, 0x12, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x11,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x12, 0x4b, 0x0a, 0x0c, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x65, 0x72, 0x52, 0x0b, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65,
	0x72, 0x22, 0xcd, 0x02, 0x0a, 0x1a, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65,
	0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x60, 0x0a, 0x15, 0x72, 0x61, 0x74, 0x65,
	0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x13, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x75, 0x0a, 0x1b, 0x72, 0x61,
	0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x64, 0x79, 0x6e, 0x61, 0x6d,
	0x69, 0x63, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x36, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69,
	0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x18, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x65, 0x72, 0x44, 0x79, 0x6e, 0x61, 0x6d, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x22, 0x2b, 0x0a, 0x13, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72,
	0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0xcf,
	0x01, 0x0a, 0x1d, 0x4c, 0x65, 0x61, 0x6b, 0x79, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x56, 0x0a, 0x0c, 0x72, 0x61, 0x74, 0x65,
	0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33,
	0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x65, 0x61,
	0x6b, 0x79, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x65, 0x72, 0x52, 0x0b, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72,
	0x22, 0xec, 0x01, 0x0a, 0x25, 0x4c, 0x65, 0x61, 0x6b, 0x79, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x56, 0x0a, 0x11, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x52, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x12, 0x6b, 0x0a, 0x15, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x65, 0x72, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x37, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x65, 0x61, 0x6b,
	0x79, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x13, 0x72, 0x61, 0x74, 0x65,
	0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x22,
	0xaa, 0x01, 0x0a, 0x1e, 0x4c, 0x65, 0x61, 0x6b, 0x79, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52,
	0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x63, 0x61, 0x70,
	0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x3e, 0x0a, 0x0d, 0x6c,
	0x65, 0x61, 0x6b, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6c,
	0x65, 0x61, 0x6b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x6c,
	0x65, 0x61, 0x6b, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x0a, 0x6c, 0x65, 0x61, 0x6b, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x94, 0x02, 0x0a,
	0x2f, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x76, 0x31,
	0x42, 0x10, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x50, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70,
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
	file_aperture_policy_sync_v1_rate_limiter_proto_rawDescOnce sync.Once
	file_aperture_policy_sync_v1_rate_limiter_proto_rawDescData = file_aperture_policy_sync_v1_rate_limiter_proto_rawDesc
)

func file_aperture_policy_sync_v1_rate_limiter_proto_rawDescGZIP() []byte {
	file_aperture_policy_sync_v1_rate_limiter_proto_rawDescOnce.Do(func() {
		file_aperture_policy_sync_v1_rate_limiter_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_sync_v1_rate_limiter_proto_rawDescData)
	})
	return file_aperture_policy_sync_v1_rate_limiter_proto_rawDescData
}

var file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_aperture_policy_sync_v1_rate_limiter_proto_goTypes = []interface{}{
	(*RateLimiterWrapper)(nil),                    // 0: aperture.policy.sync.v1.RateLimiterWrapper
	(*RateLimiterDecisionWrapper)(nil),            // 1: aperture.policy.sync.v1.RateLimiterDecisionWrapper
	(*RateLimiterDecision)(nil),                   // 2: aperture.policy.sync.v1.RateLimiterDecision
	(*LeakyBucketRateLimiterWrapper)(nil),         // 3: aperture.policy.sync.v1.LeakyBucketRateLimiterWrapper
	(*LeakyBucketRateLimiterDecisionWrapper)(nil), // 4: aperture.policy.sync.v1.LeakyBucketRateLimiterDecisionWrapper
	(*LeakyBucketRateLimiterDecision)(nil),        // 5: aperture.policy.sync.v1.LeakyBucketRateLimiterDecision
	(*CommonAttributes)(nil),                      // 6: aperture.policy.sync.v1.CommonAttributes
	(*v1.RateLimiter)(nil),                        // 7: aperture.policy.language.v1.RateLimiter
	(*v1.RateLimiter_DynamicConfig)(nil),          // 8: aperture.policy.language.v1.RateLimiter.DynamicConfig
	(*v1.LeakyBucketRateLimiter)(nil),             // 9: aperture.policy.language.v1.LeakyBucketRateLimiter
	(*durationpb.Duration)(nil),                   // 10: google.protobuf.Duration
}
var file_aperture_policy_sync_v1_rate_limiter_proto_depIdxs = []int32{
	6,  // 0: aperture.policy.sync.v1.RateLimiterWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	7,  // 1: aperture.policy.sync.v1.RateLimiterWrapper.rate_limiter:type_name -> aperture.policy.language.v1.RateLimiter
	6,  // 2: aperture.policy.sync.v1.RateLimiterDecisionWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	2,  // 3: aperture.policy.sync.v1.RateLimiterDecisionWrapper.rate_limiter_decision:type_name -> aperture.policy.sync.v1.RateLimiterDecision
	8,  // 4: aperture.policy.sync.v1.RateLimiterDecisionWrapper.rate_limiter_dynamic_config:type_name -> aperture.policy.language.v1.RateLimiter.DynamicConfig
	6,  // 5: aperture.policy.sync.v1.LeakyBucketRateLimiterWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	9,  // 6: aperture.policy.sync.v1.LeakyBucketRateLimiterWrapper.rate_limiter:type_name -> aperture.policy.language.v1.LeakyBucketRateLimiter
	6,  // 7: aperture.policy.sync.v1.LeakyBucketRateLimiterDecisionWrapper.common_attributes:type_name -> aperture.policy.sync.v1.CommonAttributes
	5,  // 8: aperture.policy.sync.v1.LeakyBucketRateLimiterDecisionWrapper.rate_limiter_decision:type_name -> aperture.policy.sync.v1.LeakyBucketRateLimiterDecision
	10, // 9: aperture.policy.sync.v1.LeakyBucketRateLimiterDecision.leak_interval:type_name -> google.protobuf.Duration
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_aperture_policy_sync_v1_rate_limiter_proto_init() }
func file_aperture_policy_sync_v1_rate_limiter_proto_init() {
	if File_aperture_policy_sync_v1_rate_limiter_proto != nil {
		return
	}
	file_aperture_policy_sync_v1_common_attributes_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimiterDecisionWrapper); i {
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
		file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeakyBucketRateLimiterWrapper); i {
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
		file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeakyBucketRateLimiterDecisionWrapper); i {
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
		file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeakyBucketRateLimiterDecision); i {
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
			RawDescriptor: file_aperture_policy_sync_v1_rate_limiter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_sync_v1_rate_limiter_proto_goTypes,
		DependencyIndexes: file_aperture_policy_sync_v1_rate_limiter_proto_depIdxs,
		MessageInfos:      file_aperture_policy_sync_v1_rate_limiter_proto_msgTypes,
	}.Build()
	File_aperture_policy_sync_v1_rate_limiter_proto = out.File
	file_aperture_policy_sync_v1_rate_limiter_proto_rawDesc = nil
	file_aperture_policy_sync_v1_rate_limiter_proto_goTypes = nil
	file_aperture_policy_sync_v1_rate_limiter_proto_depIdxs = nil
}
