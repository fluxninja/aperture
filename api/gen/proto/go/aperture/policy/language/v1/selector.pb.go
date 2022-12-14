// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/policy/language/v1/selector.proto

package languagev1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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

// Describes which flow in which service a [flow control
// component](/concepts/flow-control/flow-control.md#components) should apply
// to
//
// :::info
// See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).
// :::
type FlowSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceSelector *ServiceSelector `protobuf:"bytes,1,opt,name=service_selector,json=serviceSelector,proto3" json:"service_selector,omitempty" validate:"required"` // @gotags: validate:"required"
	FlowMatcher     *FlowMatcher     `protobuf:"bytes,2,opt,name=flow_matcher,json=flowMatcher,proto3" json:"flow_matcher,omitempty" validate:"required"`             // @gotags: validate:"required"
}

func (x *FlowSelector) Reset() {
	*x = FlowSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowSelector) ProtoMessage() {}

func (x *FlowSelector) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowSelector.ProtoReflect.Descriptor instead.
func (*FlowSelector) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_selector_proto_rawDescGZIP(), []int{0}
}

func (x *FlowSelector) GetServiceSelector() *ServiceSelector {
	if x != nil {
		return x.ServiceSelector
	}
	return nil
}

func (x *FlowSelector) GetFlowMatcher() *FlowMatcher {
	if x != nil {
		return x.FlowMatcher
	}
	return nil
}

// Describes which service a [flow control or observability
// component](/concepts/flow-control/flow-control.md#components) should apply
// to
//
// :::info
// See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).
// :::
type ServiceSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Which [agent-group](/concepts/service.md#agent-group) this
	// selector applies to.
	AgentGroup string `protobuf:"bytes,1,opt,name=agent_group,json=agentGroup,proto3" json:"agent_group,omitempty" default:"default"` // @gotags: default:"default"
	// The Fully Qualified Domain Name of the
	// [service](/concepts/service.md) to select.
	//
	// In kubernetes, this is the FQDN of the Service object.
	//
	// Empty string means all services within an agent group (catch-all).
	//
	// :::note
	// One entity may belong to multiple services.
	// :::
	Service string `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
}

func (x *ServiceSelector) Reset() {
	*x = ServiceSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceSelector) ProtoMessage() {}

func (x *ServiceSelector) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceSelector.ProtoReflect.Descriptor instead.
func (*ServiceSelector) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_selector_proto_rawDescGZIP(), []int{1}
}

func (x *ServiceSelector) GetAgentGroup() string {
	if x != nil {
		return x.AgentGroup
	}
	return ""
}

func (x *ServiceSelector) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

// Describes which flows a [flow control
// component](/concepts/flow-control/flow-control.md#components) should apply
// to
//
// :::info
// See also [FlowSelector overview](/concepts/flow-control/flow-selector.md).
// :::
//
// Example:
// ```yaml
// control_point: ingress
// label_matcher:
//   match_labels:
//     user_tier: gold
//   match_expressions:
//     - key: query
//       operator: In
//       values:
//         - insert
//         - delete
//     - label: user_agent
//       regex: ^(?!.*Chrome).*Safari
// ```
type FlowMatcher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// [Control Point](/concepts/flow-control/flow-control.md#control-point)
	// identifies the location of a Flow within a Service. For an SDK based insertion, a Control Point can represent a particular feature or execution
	// block within a Service. In case of Service Mesh or Middleware insertion, a Control Point can identify ingress vs egress calls or distinct listeners
	// or filter chains.
	ControlPoint string `protobuf:"bytes,1,opt,name=control_point,json=controlPoint,proto3" json:"control_point,omitempty" validate:"required"` // @gotags: validate:"required"
	// Label matcher allows to add _additional_ condition on
	// [flow labels](/concepts/flow-control/flow-label.md)
	// must also be satisfied (in addition to service+control point matching)
	//
	// :::info
	// See also [Label Matcher overview](/concepts/flow-control/flow-selector.md#label-matcher).
	// :::
	//
	// :::note
	// [Classifiers](#v1-classifier) _can_ use flow labels created by some other
	// classifier, but only if they were created at some previous control point
	// (and propagated in baggage).
	//
	// This limitation doesn't apply to selectors of other entities, like
	// Flux Meters or Actuators. It's valid to create a flow label on a control
	// point using classifier, and immediately use it for matching on the same
	// control point.
	// :::
	LabelMatcher *LabelMatcher `protobuf:"bytes,2,opt,name=label_matcher,json=labelMatcher,proto3" json:"label_matcher,omitempty"`
}

func (x *FlowMatcher) Reset() {
	*x = FlowMatcher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowMatcher) ProtoMessage() {}

func (x *FlowMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowMatcher.ProtoReflect.Descriptor instead.
func (*FlowMatcher) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_selector_proto_rawDescGZIP(), []int{2}
}

func (x *FlowMatcher) GetControlPoint() string {
	if x != nil {
		return x.ControlPoint
	}
	return ""
}

func (x *FlowMatcher) GetLabelMatcher() *LabelMatcher {
	if x != nil {
		return x.LabelMatcher
	}
	return nil
}

// Describes which pods a control or observability
// component should apply to.
type KubernetesSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Which [agent-group](/concepts/service.md#agent-group) this
	// selector applies to.
	AgentGroup string `protobuf:"bytes,1,opt,name=agent_group,json=agentGroup,proto3" json:"agent_group,omitempty" default:"default"` // @gotags: default:"default"
	// Kubernetes namespace that the resource belongs to.
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty" validate:"required"` // @gotags: validate:"required"
	// Kubernetes resource type.
	Kind string `protobuf:"bytes,3,opt,name=kind,proto3" json:"kind,omitempty" validate:"required"` // @gotags: validate:"required"
	// Kubernetes resource name.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty" validate:"required"` // @gotags: validate:"required"
}

func (x *KubernetesSelector) Reset() {
	*x = KubernetesSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KubernetesSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KubernetesSelector) ProtoMessage() {}

func (x *KubernetesSelector) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_selector_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KubernetesSelector.ProtoReflect.Descriptor instead.
func (*KubernetesSelector) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_selector_proto_rawDescGZIP(), []int{3}
}

func (x *KubernetesSelector) GetAgentGroup() string {
	if x != nil {
		return x.AgentGroup
	}
	return ""
}

func (x *KubernetesSelector) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *KubernetesSelector) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *KubernetesSelector) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_aperture_policy_language_v1_selector_proto protoreflect.FileDescriptor

var file_aperture_policy_language_v1_selector_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x5f, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfa, 0x01, 0x0a, 0x0c, 0x46,
	0x6c, 0x6f, 0x77, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x7a, 0x0a, 0x10, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x65, 0x6c, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x42, 0x21, 0x92, 0x41, 0x1e, 0x82, 0x03, 0x1b, 0x0a, 0x0d, 0x78, 0x2d, 0x67,
	0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x1a, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x6e, 0x0a, 0x0c, 0x66, 0x6c, 0x6f, 0x77, 0x5f,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x42, 0x21, 0x92, 0x41, 0x1e, 0x82, 0x03, 0x1b, 0x0a,
	0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0a,
	0x1a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x0b, 0x66, 0x6c, 0x6f, 0x77,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x22, 0x6d, 0x0a, 0x0f, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x40, 0x0a, 0x0b, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x1f, 0x92, 0x41, 0x1c, 0x82, 0x03, 0x19, 0x0a, 0x0c, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x09, 0x1a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x52, 0x0a, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0xa5, 0x01, 0x0a, 0x0b, 0x46, 0x6c, 0x6f, 0x77, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x46, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0x92,
	0x41, 0x1e, 0x82, 0x03, 0x1b, 0x0a, 0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x1a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x4e,
	0x0a, 0x0d, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72,
	0x52, 0x0c, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x22, 0x85,
	0x02, 0x0a, 0x12, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x53, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x40, 0x0a, 0x0b, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1f, 0x92, 0x41, 0x1c, 0x82,
	0x03, 0x19, 0x0a, 0x0c, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x12, 0x09, 0x1a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x52, 0x0a, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x3f, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0x92, 0x41, 0x1e, 0x82,
	0x03, 0x1b, 0x0a, 0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x0a, 0x1a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0x92, 0x41, 0x1e, 0x82, 0x03, 0x1b, 0x0a, 0x0d,
	0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x1a,
	0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12,
	0x35, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x21, 0x92,
	0x41, 0x1e, 0x82, 0x03, 0x1b, 0x0a, 0x0d, 0x78, 0x2d, 0x67, 0x6f, 0x2d, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x0a, 0x1a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0xaa, 0x02, 0x0a, 0x33, 0x63, 0x6f, 0x6d, 0x2e, 0x66,
	0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0d,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x55, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78,
	0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x4c, 0xaa, 0x02, 0x1b, 0x41,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x4c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1b, 0x41, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x4c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x4c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_language_v1_selector_proto_rawDescOnce sync.Once
	file_aperture_policy_language_v1_selector_proto_rawDescData = file_aperture_policy_language_v1_selector_proto_rawDesc
)

func file_aperture_policy_language_v1_selector_proto_rawDescGZIP() []byte {
	file_aperture_policy_language_v1_selector_proto_rawDescOnce.Do(func() {
		file_aperture_policy_language_v1_selector_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_language_v1_selector_proto_rawDescData)
	})
	return file_aperture_policy_language_v1_selector_proto_rawDescData
}

var file_aperture_policy_language_v1_selector_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_aperture_policy_language_v1_selector_proto_goTypes = []interface{}{
	(*FlowSelector)(nil),       // 0: aperture.policy.language.v1.FlowSelector
	(*ServiceSelector)(nil),    // 1: aperture.policy.language.v1.ServiceSelector
	(*FlowMatcher)(nil),        // 2: aperture.policy.language.v1.FlowMatcher
	(*KubernetesSelector)(nil), // 3: aperture.policy.language.v1.KubernetesSelector
	(*LabelMatcher)(nil),       // 4: aperture.policy.language.v1.LabelMatcher
}
var file_aperture_policy_language_v1_selector_proto_depIdxs = []int32{
	1, // 0: aperture.policy.language.v1.FlowSelector.service_selector:type_name -> aperture.policy.language.v1.ServiceSelector
	2, // 1: aperture.policy.language.v1.FlowSelector.flow_matcher:type_name -> aperture.policy.language.v1.FlowMatcher
	4, // 2: aperture.policy.language.v1.FlowMatcher.label_matcher:type_name -> aperture.policy.language.v1.LabelMatcher
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_aperture_policy_language_v1_selector_proto_init() }
func file_aperture_policy_language_v1_selector_proto_init() {
	if File_aperture_policy_language_v1_selector_proto != nil {
		return
	}
	file_aperture_policy_language_v1_label_matcher_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_language_v1_selector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowSelector); i {
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
		file_aperture_policy_language_v1_selector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceSelector); i {
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
		file_aperture_policy_language_v1_selector_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowMatcher); i {
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
		file_aperture_policy_language_v1_selector_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KubernetesSelector); i {
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
			RawDescriptor: file_aperture_policy_language_v1_selector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_language_v1_selector_proto_goTypes,
		DependencyIndexes: file_aperture_policy_language_v1_selector_proto_depIdxs,
		MessageInfos:      file_aperture_policy_language_v1_selector_proto_msgTypes,
	}.Build()
	File_aperture_policy_language_v1_selector_proto = out.File
	file_aperture_policy_language_v1_selector_proto_rawDesc = nil
	file_aperture_policy_language_v1_selector_proto_goTypes = nil
	file_aperture_policy_language_v1_selector_proto_depIdxs = nil
}
