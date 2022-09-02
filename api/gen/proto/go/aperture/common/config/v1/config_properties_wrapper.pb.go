// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/common/config/v1/config_properties_wrapper.proto

package configv1

import (
	v11 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
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

	// Policy
	Policy *v1.Policy `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,2,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,3,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *PolicyWrapper) Reset() {
	*x = PolicyWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PolicyWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PolicyWrapper) ProtoMessage() {}

func (x *PolicyWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[0]
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
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{0}
}

func (x *PolicyWrapper) GetPolicy() *v1.Policy {
	if x != nil {
		return x.Policy
	}
	return nil
}

func (x *PolicyWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *PolicyWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

type FluxMeterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Flux Meter
	FluxMeter *v1.FluxMeter `protobuf:"bytes,1,opt,name=flux_meter,json=fluxMeter,proto3" json:"flux_meter,omitempty"`
	// Name of fluxmeter metric.
	FluxmeterName string `protobuf:"bytes,4,opt,name=fluxmeter_name,json=fluxmeterName,proto3" json:"fluxmeter_name,omitempty"`
}

func (x *FluxMeterWrapper) Reset() {
	*x = FluxMeterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FluxMeterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FluxMeterWrapper) ProtoMessage() {}

func (x *FluxMeterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[1]
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
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{1}
}

func (x *FluxMeterWrapper) GetFluxMeter() *v1.FluxMeter {
	if x != nil {
		return x.FluxMeter
	}
	return nil
}

func (x *FluxMeterWrapper) GetFluxmeterName() string {
	if x != nil {
		return x.FluxmeterName
	}
	return ""
}

type ClassifierWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Classifier
	Classifier *v1.Classifier `protobuf:"bytes,1,opt,name=classifier,proto3" json:"classifier,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,2,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,3,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *ClassifierWrapper) Reset() {
	*x = ClassifierWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassifierWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassifierWrapper) ProtoMessage() {}

func (x *ClassifierWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[2]
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
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{2}
}

func (x *ClassifierWrapper) GetClassifier() *v1.Classifier {
	if x != nil {
		return x.Classifier
	}
	return nil
}

func (x *ClassifierWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *ClassifierWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

type ConcurrencyLimiterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Concurrency Limiter
	ConcurrencyLimiter *v1.ConcurrencyLimiter `protobuf:"bytes,1,opt,name=concurrency_limiter,json=concurrencyLimiter,proto3" json:"concurrency_limiter,omitempty"`
	// The index of Component in the Circuit.
	ComponentIndex int64 `protobuf:"varint,2,opt,name=component_index,json=componentIndex,proto3" json:"component_index,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,3,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,4,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *ConcurrencyLimiterWrapper) Reset() {
	*x = ConcurrencyLimiterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcurrencyLimiterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcurrencyLimiterWrapper) ProtoMessage() {}

func (x *ConcurrencyLimiterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[3]
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
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{3}
}

func (x *ConcurrencyLimiterWrapper) GetConcurrencyLimiter() *v1.ConcurrencyLimiter {
	if x != nil {
		return x.ConcurrencyLimiter
	}
	return nil
}

func (x *ConcurrencyLimiterWrapper) GetComponentIndex() int64 {
	if x != nil {
		return x.ComponentIndex
	}
	return 0
}

func (x *ConcurrencyLimiterWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *ConcurrencyLimiterWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

type RateLimiterWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rate Limiter
	RateLimiter *v1.RateLimiter `protobuf:"bytes,1,opt,name=rate_limiter,json=rateLimiter,proto3" json:"rate_limiter,omitempty"`
	// The index of Component in the Circuit.
	ComponentIndex int64 `protobuf:"varint,2,opt,name=component_index,json=componentIndex,proto3" json:"component_index,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,3,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,4,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *RateLimiterWrapper) Reset() {
	*x = RateLimiterWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterWrapper) ProtoMessage() {}

func (x *RateLimiterWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[4]
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
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{4}
}

func (x *RateLimiterWrapper) GetRateLimiter() *v1.RateLimiter {
	if x != nil {
		return x.RateLimiter
	}
	return nil
}

func (x *RateLimiterWrapper) GetComponentIndex() int64 {
	if x != nil {
		return x.ComponentIndex
	}
	return 0
}

func (x *RateLimiterWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *RateLimiterWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

type LoadShedDecsisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Load Shed Decision
	LoadShedDecision *v11.LoadShedDecision `protobuf:"bytes,1,opt,name=load_shed_decision,json=loadShedDecision,proto3" json:"load_shed_decision,omitempty"`
	// The index of Component in the Circuit.
	ComponentIndex int64 `protobuf:"varint,2,opt,name=component_index,json=componentIndex,proto3" json:"component_index,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,3,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,4,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *LoadShedDecsisionWrapper) Reset() {
	*x = LoadShedDecsisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadShedDecsisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadShedDecsisionWrapper) ProtoMessage() {}

func (x *LoadShedDecsisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadShedDecsisionWrapper.ProtoReflect.Descriptor instead.
func (*LoadShedDecsisionWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{5}
}

func (x *LoadShedDecsisionWrapper) GetLoadShedDecision() *v11.LoadShedDecision {
	if x != nil {
		return x.LoadShedDecision
	}
	return nil
}

func (x *LoadShedDecsisionWrapper) GetComponentIndex() int64 {
	if x != nil {
		return x.ComponentIndex
	}
	return 0
}

func (x *LoadShedDecsisionWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *LoadShedDecsisionWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

type TokensDecisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Tokens Decision
	TokensDecision *v11.TokensDecision `protobuf:"bytes,1,opt,name=tokens_decision,json=tokensDecision,proto3" json:"tokens_decision,omitempty"`
	// The index of Component in the Circuit.
	ComponentIndex int64 `protobuf:"varint,2,opt,name=component_index,json=componentIndex,proto3" json:"component_index,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,3,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,4,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *TokensDecisionWrapper) Reset() {
	*x = TokensDecisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokensDecisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokensDecisionWrapper) ProtoMessage() {}

func (x *TokensDecisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokensDecisionWrapper.ProtoReflect.Descriptor instead.
func (*TokensDecisionWrapper) Descriptor() ([]byte, []int) {
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{6}
}

func (x *TokensDecisionWrapper) GetTokensDecision() *v11.TokensDecision {
	if x != nil {
		return x.TokensDecision
	}
	return nil
}

func (x *TokensDecisionWrapper) GetComponentIndex() int64 {
	if x != nil {
		return x.ComponentIndex
	}
	return 0
}

func (x *TokensDecisionWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *TokensDecisionWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

type RateLimiterDecisionWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rate Limiter Decision
	RateLimiterDecision *v11.RateLimiterDecision `protobuf:"bytes,1,opt,name=rate_limiter_decision,json=rateLimiterDecision,proto3" json:"rate_limiter_decision,omitempty"`
	// The index of Component in the Circuit.
	ComponentIndex int64 `protobuf:"varint,2,opt,name=component_index,json=componentIndex,proto3" json:"component_index,omitempty"`
	// Name of the Policy.
	PolicyName string `protobuf:"bytes,3,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	// Hash of the entire Policy spec.
	PolicyHash string `protobuf:"bytes,4,opt,name=policy_hash,json=policyHash,proto3" json:"policy_hash,omitempty"`
}

func (x *RateLimiterDecisionWrapper) Reset() {
	*x = RateLimiterDecisionWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiterDecisionWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiterDecisionWrapper) ProtoMessage() {}

func (x *RateLimiterDecisionWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[7]
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
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP(), []int{7}
}

func (x *RateLimiterDecisionWrapper) GetRateLimiterDecision() *v11.RateLimiterDecision {
	if x != nil {
		return x.RateLimiterDecision
	}
	return nil
}

func (x *RateLimiterDecisionWrapper) GetComponentIndex() int64 {
	if x != nil {
		return x.ComponentIndex
	}
	return 0
}

func (x *RateLimiterDecisionWrapper) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

func (x *RateLimiterDecisionWrapper) GetPolicyHash() string {
	if x != nil {
		return x.PolicyHash
	}
	return ""
}

var File_aperture_common_config_v1_config_properties_wrapper_proto protoreflect.FileDescriptor

var file_aperture_common_config_v1_config_properties_wrapper_proto_rawDesc = []byte{
	0x0a, 0x39, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x5f, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2b, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x66, 0x6c, 0x75, 0x78, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x28, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x01, 0x0a, 0x0d, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x3b, 0x0a, 0x06, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0x80, 0x01, 0x0a, 0x10, 0x46,
	0x6c, 0x75, 0x78, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12,
	0x45, 0x0a, 0x0a, 0x66, 0x6c, 0x75, 0x78, 0x5f, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x46, 0x6c, 0x75, 0x78, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x09, 0x66, 0x6c, 0x75,
	0x78, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x6c, 0x75, 0x78, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x66, 0x6c, 0x75, 0x78, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x9e, 0x01,
	0x0a, 0x11, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x12, 0x47, 0x0a, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x52, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0xe8,
	0x01, 0x0a, 0x19, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x60, 0x0a, 0x13,
	0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x63, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x27,
	0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
	0x6e, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0xcc, 0x01, 0x0a, 0x12, 0x52, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x12, 0x4b, 0x0a, 0x0c, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72,
	0x52, 0x0b, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x27, 0x0a,
	0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e,
	0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0xe3, 0x01, 0x0a, 0x18, 0x4c, 0x6f, 0x61,
	0x64, 0x53, 0x68, 0x65, 0x64, 0x44, 0x65, 0x63, 0x73, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x57, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x5c, 0x0a, 0x12, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x73, 0x68,
	0x65, 0x64, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2e, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x53, 0x68, 0x65, 0x64, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x10, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x68, 0x65, 0x64, 0x44, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0xd9,
	0x01, 0x0a, 0x15, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x55, 0x0a, 0x0f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x73, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x0e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0xee, 0x01, 0x0a, 0x1a, 0x52,
	0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x65, 0x0a, 0x15, 0x72, 0x61, 0x74,
	0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x64, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x13, 0x72, 0x61, 0x74,
	0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x48, 0x61, 0x73, 0x68, 0x42, 0x97, 0x02, 0x0a, 0x1d,
	0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x1c, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x57,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x51, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69,
	0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x41, 0x43, 0x43, 0xaa, 0x02, 0x19, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x56, 0x31, 0xca, 0x02, 0x19, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x25, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescOnce sync.Once
	file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescData = file_aperture_common_config_v1_config_properties_wrapper_proto_rawDesc
)

func file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescGZIP() []byte {
	file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescOnce.Do(func() {
		file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescData)
	})
	return file_aperture_common_config_v1_config_properties_wrapper_proto_rawDescData
}

var file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_aperture_common_config_v1_config_properties_wrapper_proto_goTypes = []interface{}{
	(*PolicyWrapper)(nil),              // 0: aperture.common.config.v1.PolicyWrapper
	(*FluxMeterWrapper)(nil),           // 1: aperture.common.config.v1.FluxMeterWrapper
	(*ClassifierWrapper)(nil),          // 2: aperture.common.config.v1.ClassifierWrapper
	(*ConcurrencyLimiterWrapper)(nil),  // 3: aperture.common.config.v1.ConcurrencyLimiterWrapper
	(*RateLimiterWrapper)(nil),         // 4: aperture.common.config.v1.RateLimiterWrapper
	(*LoadShedDecsisionWrapper)(nil),   // 5: aperture.common.config.v1.LoadShedDecsisionWrapper
	(*TokensDecisionWrapper)(nil),      // 6: aperture.common.config.v1.TokensDecisionWrapper
	(*RateLimiterDecisionWrapper)(nil), // 7: aperture.common.config.v1.RateLimiterDecisionWrapper
	(*v1.Policy)(nil),                  // 8: aperture.policy.language.v1.Policy
	(*v1.FluxMeter)(nil),               // 9: aperture.policy.language.v1.FluxMeter
	(*v1.Classifier)(nil),              // 10: aperture.policy.language.v1.Classifier
	(*v1.ConcurrencyLimiter)(nil),      // 11: aperture.policy.language.v1.ConcurrencyLimiter
	(*v1.RateLimiter)(nil),             // 12: aperture.policy.language.v1.RateLimiter
	(*v11.LoadShedDecision)(nil),       // 13: aperture.policy.decisions.v1.LoadShedDecision
	(*v11.TokensDecision)(nil),         // 14: aperture.policy.decisions.v1.TokensDecision
	(*v11.RateLimiterDecision)(nil),    // 15: aperture.policy.decisions.v1.RateLimiterDecision
}
var file_aperture_common_config_v1_config_properties_wrapper_proto_depIdxs = []int32{
	8,  // 0: aperture.common.config.v1.PolicyWrapper.policy:type_name -> aperture.policy.language.v1.Policy
	9,  // 1: aperture.common.config.v1.FluxMeterWrapper.flux_meter:type_name -> aperture.policy.language.v1.FluxMeter
	10, // 2: aperture.common.config.v1.ClassifierWrapper.classifier:type_name -> aperture.policy.language.v1.Classifier
	11, // 3: aperture.common.config.v1.ConcurrencyLimiterWrapper.concurrency_limiter:type_name -> aperture.policy.language.v1.ConcurrencyLimiter
	12, // 4: aperture.common.config.v1.RateLimiterWrapper.rate_limiter:type_name -> aperture.policy.language.v1.RateLimiter
	13, // 5: aperture.common.config.v1.LoadShedDecsisionWrapper.load_shed_decision:type_name -> aperture.policy.decisions.v1.LoadShedDecision
	14, // 6: aperture.common.config.v1.TokensDecisionWrapper.tokens_decision:type_name -> aperture.policy.decisions.v1.TokensDecision
	15, // 7: aperture.common.config.v1.RateLimiterDecisionWrapper.rate_limiter_decision:type_name -> aperture.policy.decisions.v1.RateLimiterDecision
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_aperture_common_config_v1_config_properties_wrapper_proto_init() }
func file_aperture_common_config_v1_config_properties_wrapper_proto_init() {
	if File_aperture_common_config_v1_config_properties_wrapper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadShedDecsisionWrapper); i {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokensDecisionWrapper); i {
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
		file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_common_config_v1_config_properties_wrapper_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_common_config_v1_config_properties_wrapper_proto_goTypes,
		DependencyIndexes: file_aperture_common_config_v1_config_properties_wrapper_proto_depIdxs,
		MessageInfos:      file_aperture_common_config_v1_config_properties_wrapper_proto_msgTypes,
	}.Build()
	File_aperture_common_config_v1_config_properties_wrapper_proto = out.File
	file_aperture_common_config_v1_config_properties_wrapper_proto_rawDesc = nil
	file_aperture_common_config_v1_config_properties_wrapper_proto_goTypes = nil
	file_aperture_common_config_v1_config_properties_wrapper_proto_depIdxs = nil
}
