// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: aperture/policy/private/v1/flowcontrol.proto

package privatev1

import (
	v1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
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
	WorkloadLatencyBasedTokens bool              `protobuf:"varint,4,opt,name=workload_latency_based_tokens,json=workloadLatencyBasedTokens,proto3" json:"workload_latency_based_tokens,omitempty"`
	Scheduler                  *v1.Scheduler     `protobuf:"bytes,5,opt,name=scheduler,proto3" json:"scheduler,omitempty"`
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

func (x *LoadActuator) GetWorkloadLatencyBasedTokens() bool {
	if x != nil {
		return x.WorkloadLatencyBasedTokens
	}
	return false
}

func (x *LoadActuator) GetScheduler() *v1.Scheduler {
	if x != nil {
		return x.Scheduler
	}
	return nil
}

type RateLimiter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InPorts           *v1.RateLimiter_Ins `protobuf:"bytes,1,opt,name=in_ports,json=inPorts,proto3" json:"in_ports,omitempty" validate:"required"` // @gotags: validate:"required"
	Selectors         []*v1.Selector      `protobuf:"bytes,2,rep,name=selectors,proto3" json:"selectors,omitempty" validate:"required,gt=0,dive"`            // @gotags: validate:"required,gt=0,dive"
	ParentComponentId string              `protobuf:"bytes,3,opt,name=parent_component_id,json=parentComponentId,proto3" json:"parent_component_id,omitempty"`
}

func (x *RateLimiter) Reset() {
	*x = RateLimiter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimiter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimiter) ProtoMessage() {}

func (x *RateLimiter) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use RateLimiter.ProtoReflect.Descriptor instead.
func (*RateLimiter) Descriptor() ([]byte, []int) {
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP(), []int{1}
}

func (x *RateLimiter) GetInPorts() *v1.RateLimiter_Ins {
	if x != nil {
		return x.InPorts
	}
	return nil
}

func (x *RateLimiter) GetSelectors() []*v1.Selector {
	if x != nil {
		return x.Selectors
	}
	return nil
}

func (x *RateLimiter) GetParentComponentId() string {
	if x != nil {
		return x.ParentComponentId
	}
	return ""
}

type QuotaScheduler struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InPorts           *v1.RateLimiter_Ins `protobuf:"bytes,1,opt,name=in_ports,json=inPorts,proto3" json:"in_ports,omitempty" validate:"required"` // @gotags: validate:"required"
	Selectors         []*v1.Selector      `protobuf:"bytes,2,rep,name=selectors,proto3" json:"selectors,omitempty" validate:"required,gt=0,dive"`            // @gotags: validate:"required,gt=0,dive"
	ParentComponentId string              `protobuf:"bytes,3,opt,name=parent_component_id,json=parentComponentId,proto3" json:"parent_component_id,omitempty"`
	Scheduler         *v1.Scheduler       `protobuf:"bytes,4,opt,name=scheduler,proto3" json:"scheduler,omitempty"`
}

func (x *QuotaScheduler) Reset() {
	*x = QuotaScheduler{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuotaScheduler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuotaScheduler) ProtoMessage() {}

func (x *QuotaScheduler) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuotaScheduler.ProtoReflect.Descriptor instead.
func (*QuotaScheduler) Descriptor() ([]byte, []int) {
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP(), []int{2}
}

func (x *QuotaScheduler) GetInPorts() *v1.RateLimiter_Ins {
	if x != nil {
		return x.InPorts
	}
	return nil
}

func (x *QuotaScheduler) GetSelectors() []*v1.Selector {
	if x != nil {
		return x.Selectors
	}
	return nil
}

func (x *QuotaScheduler) GetParentComponentId() string {
	if x != nil {
		return x.ParentComponentId
	}
	return ""
}

func (x *QuotaScheduler) GetScheduler() *v1.Scheduler {
	if x != nil {
		return x.Scheduler
	}
	return nil
}

type ConcurrencyLimiter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InPorts           *v1.ConcurrencyLimiter_Ins `protobuf:"bytes,1,opt,name=in_ports,json=inPorts,proto3" json:"in_ports,omitempty" validate:"required"` // @gotags: validate:"required"
	Selectors         []*v1.Selector             `protobuf:"bytes,2,rep,name=selectors,proto3" json:"selectors,omitempty" validate:"required,gt=0,dive"`            // @gotags: validate:"required,gt=0,dive"
	ParentComponentId string                     `protobuf:"bytes,3,opt,name=parent_component_id,json=parentComponentId,proto3" json:"parent_component_id,omitempty"`
}

func (x *ConcurrencyLimiter) Reset() {
	*x = ConcurrencyLimiter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcurrencyLimiter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcurrencyLimiter) ProtoMessage() {}

func (x *ConcurrencyLimiter) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConcurrencyLimiter.ProtoReflect.Descriptor instead.
func (*ConcurrencyLimiter) Descriptor() ([]byte, []int) {
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP(), []int{3}
}

func (x *ConcurrencyLimiter) GetInPorts() *v1.ConcurrencyLimiter_Ins {
	if x != nil {
		return x.InPorts
	}
	return nil
}

func (x *ConcurrencyLimiter) GetSelectors() []*v1.Selector {
	if x != nil {
		return x.Selectors
	}
	return nil
}

func (x *ConcurrencyLimiter) GetParentComponentId() string {
	if x != nil {
		return x.ParentComponentId
	}
	return ""
}

type ConcurrencyScheduler struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InPorts           *v1.ConcurrencyLimiter_Ins `protobuf:"bytes,1,opt,name=in_ports,json=inPorts,proto3" json:"in_ports,omitempty" validate:"required"` // @gotags: validate:"required"
	Selectors         []*v1.Selector             `protobuf:"bytes,2,rep,name=selectors,proto3" json:"selectors,omitempty" validate:"required,gt=0,dive"`            // @gotags: validate:"required,gt=0,dive"
	ParentComponentId string                     `protobuf:"bytes,3,opt,name=parent_component_id,json=parentComponentId,proto3" json:"parent_component_id,omitempty"`
	Scheduler         *v1.Scheduler              `protobuf:"bytes,4,opt,name=scheduler,proto3" json:"scheduler,omitempty"`
}

func (x *ConcurrencyScheduler) Reset() {
	*x = ConcurrencyScheduler{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcurrencyScheduler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcurrencyScheduler) ProtoMessage() {}

func (x *ConcurrencyScheduler) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConcurrencyScheduler.ProtoReflect.Descriptor instead.
func (*ConcurrencyScheduler) Descriptor() ([]byte, []int) {
	return file_aperture_policy_private_v1_flowcontrol_proto_rawDescGZIP(), []int{4}
}

func (x *ConcurrencyScheduler) GetInPorts() *v1.ConcurrencyLimiter_Ins {
	if x != nil {
		return x.InPorts
	}
	return nil
}

func (x *ConcurrencyScheduler) GetSelectors() []*v1.Selector {
	if x != nil {
		return x.Selectors
	}
	return nil
}

func (x *ConcurrencyScheduler) GetParentComponentId() string {
	if x != nil {
		return x.ParentComponentId
	}
	return ""
}

func (x *ConcurrencyScheduler) GetScheduler() *v1.Scheduler {
	if x != nil {
		return x.Scheduler
	}
	return nil
}

type LoadActuator_Ins struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoadMultiplier *v1.InPort `protobuf:"bytes,1,opt,name=load_multiplier,json=loadMultiplier,proto3" json:"load_multiplier,omitempty"`
}

func (x *LoadActuator_Ins) Reset() {
	*x = LoadActuator_Ins{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadActuator_Ins) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadActuator_Ins) ProtoMessage() {}

func (x *LoadActuator_Ins) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[5]
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
	0x74, 0x6f, 0x22, 0xb9, 0x03, 0x0a, 0x0c, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x63, 0x74, 0x75, 0x61,
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
	0x12, 0x41, 0x0a, 0x1d, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6c, 0x61, 0x74,
	0x65, 0x6e, 0x63, 0x79, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x64, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x1a, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x42, 0x61, 0x73, 0x65, 0x64, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x12, 0x44, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x52, 0x09,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x1a, 0x53, 0x0a, 0x03, 0x49, 0x6e, 0x73,
	0x12, 0x4c, 0x0a, 0x0f, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c,
	0x69, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x0e,
	0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x22, 0xcb,
	0x01, 0x0a, 0x0b, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x47,
	0x0a, 0x08, 0x69, 0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x73, 0x52, 0x07,
	0x69, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x2e, 0x0a, 0x13,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x94, 0x02, 0x0a,
	0x0e, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x12,
	0x47, 0x0a, 0x08, 0x69, 0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x73, 0x52,
	0x07, 0x69, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x2e, 0x0a,
	0x13, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x44, 0x0a,
	0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x72, 0x22, 0xd9, 0x01, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x63, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x12, 0x4e, 0x0a, 0x08, 0x69, 0x6e,
	0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x49, 0x6e,
	0x73, 0x52, 0x07, 0x69, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12,
	0x2e, 0x0a, 0x13, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22,
	0xa1, 0x02, 0x0a, 0x14, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x12, 0x4e, 0x0a, 0x08, 0x69, 0x6e, 0x5f, 0x70,
	0x6f, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x73, 0x52,
	0x07, 0x69, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x09, 0x73, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x2e, 0x0a,
	0x13, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x44, 0x0a,
	0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x72, 0x42, 0xab, 0x02, 0x0a, 0x32, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x46, 0x6c, 0x6f, 0x77,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x56,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e,
	0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x32, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x50, 0xaa, 0x02, 0x1a, 0x41,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x50,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1b, 0x41, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x50, 0x72, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x5f, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x5f, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x1d, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_aperture_policy_private_v1_flowcontrol_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_aperture_policy_private_v1_flowcontrol_proto_goTypes = []interface{}{
	(*LoadActuator)(nil),              // 0: aperture.policy.private.v1.LoadActuator
	(*RateLimiter)(nil),               // 1: aperture.policy.private.v1.RateLimiter
	(*QuotaScheduler)(nil),            // 2: aperture.policy.private.v1.QuotaScheduler
	(*ConcurrencyLimiter)(nil),        // 3: aperture.policy.private.v1.ConcurrencyLimiter
	(*ConcurrencyScheduler)(nil),      // 4: aperture.policy.private.v1.ConcurrencyScheduler
	(*LoadActuator_Ins)(nil),          // 5: aperture.policy.private.v1.LoadActuator.Ins
	(*v1.Selector)(nil),               // 6: aperture.policy.language.v1.Selector
	(*v1.Scheduler)(nil),              // 7: aperture.policy.language.v1.Scheduler
	(*v1.RateLimiter_Ins)(nil),        // 8: aperture.policy.language.v1.RateLimiter.Ins
	(*v1.ConcurrencyLimiter_Ins)(nil), // 9: aperture.policy.language.v1.ConcurrencyLimiter.Ins
	(*v1.InPort)(nil),                 // 10: aperture.policy.language.v1.InPort
}
var file_aperture_policy_private_v1_flowcontrol_proto_depIdxs = []int32{
	5,  // 0: aperture.policy.private.v1.LoadActuator.in_ports:type_name -> aperture.policy.private.v1.LoadActuator.Ins
	6,  // 1: aperture.policy.private.v1.LoadActuator.selectors:type_name -> aperture.policy.language.v1.Selector
	7,  // 2: aperture.policy.private.v1.LoadActuator.scheduler:type_name -> aperture.policy.language.v1.Scheduler
	8,  // 3: aperture.policy.private.v1.RateLimiter.in_ports:type_name -> aperture.policy.language.v1.RateLimiter.Ins
	6,  // 4: aperture.policy.private.v1.RateLimiter.selectors:type_name -> aperture.policy.language.v1.Selector
	8,  // 5: aperture.policy.private.v1.QuotaScheduler.in_ports:type_name -> aperture.policy.language.v1.RateLimiter.Ins
	6,  // 6: aperture.policy.private.v1.QuotaScheduler.selectors:type_name -> aperture.policy.language.v1.Selector
	7,  // 7: aperture.policy.private.v1.QuotaScheduler.scheduler:type_name -> aperture.policy.language.v1.Scheduler
	9,  // 8: aperture.policy.private.v1.ConcurrencyLimiter.in_ports:type_name -> aperture.policy.language.v1.ConcurrencyLimiter.Ins
	6,  // 9: aperture.policy.private.v1.ConcurrencyLimiter.selectors:type_name -> aperture.policy.language.v1.Selector
	9,  // 10: aperture.policy.private.v1.ConcurrencyScheduler.in_ports:type_name -> aperture.policy.language.v1.ConcurrencyLimiter.Ins
	6,  // 11: aperture.policy.private.v1.ConcurrencyScheduler.selectors:type_name -> aperture.policy.language.v1.Selector
	7,  // 12: aperture.policy.private.v1.ConcurrencyScheduler.scheduler:type_name -> aperture.policy.language.v1.Scheduler
	10, // 13: aperture.policy.private.v1.LoadActuator.Ins.load_multiplier:type_name -> aperture.policy.language.v1.InPort
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
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
			switch v := v.(*RateLimiter); i {
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
		file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuotaScheduler); i {
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
		file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConcurrencyLimiter); i {
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
		file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConcurrencyScheduler); i {
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
		file_aperture_policy_private_v1_flowcontrol_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
			NumMessages:   6,
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
