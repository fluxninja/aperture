// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/classification/v1/ruleset.proto

package classificationv1

import (
	v1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AllRulesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AllRules *AllRules `protobuf:"bytes,1,opt,name=all_rules,json=allRules,proto3" json:"all_rules,omitempty"`
}

func (x *AllRulesResponse) Reset() {
	*x = AllRulesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllRulesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllRulesResponse) ProtoMessage() {}

func (x *AllRulesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllRulesResponse.ProtoReflect.Descriptor instead.
func (*AllRulesResponse) Descriptor() ([]byte, []int) {
	return file_aperture_classification_v1_ruleset_proto_rawDescGZIP(), []int{0}
}

func (x *AllRulesResponse) GetAllRules() *AllRules {
	if x != nil {
		return x.AllRules
	}
	return nil
}

// All the ruleset name to ruleset mapping.
type AllRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A map of {key, value} pairs mapping from the name of the classifier to the values of the classifier.
	AllRules map[string]*Classifier `protobuf:"bytes,1,rep,name=all_rules,json=allRules,proto3" json:"all_rules,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AllRules) Reset() {
	*x = AllRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllRules) ProtoMessage() {}

func (x *AllRules) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllRules.ProtoReflect.Descriptor instead.
func (*AllRules) Descriptor() ([]byte, []int) {
	return file_aperture_classification_v1_ruleset_proto_rawDescGZIP(), []int{1}
}

func (x *AllRules) GetAllRules() map[string]*Classifier {
	if x != nil {
		return x.AllRules
	}
	return nil
}

// Set of classification rules sharing a common selector.
//
// Example:
// ```yaml
// selector:
//   service: service1.default.svc.cluster.local
//   control_point:
//     traffic: ingress
// rules:
//   user:
//     extractor:
//       from: request.http.headers.user
// ```
type Classifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Defines where to apply the flow classification rule.
	Selector *v1.Selector `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.
	Rules map[string]*Rule `protobuf:"bytes,2,rep,name=rules,proto3" json:"rules,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Classifier) Reset() {
	*x = Classifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Classifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Classifier) ProtoMessage() {}

func (x *Classifier) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Classifier.ProtoReflect.Descriptor instead.
func (*Classifier) Descriptor() ([]byte, []int) {
	return file_aperture_classification_v1_ruleset_proto_rawDescGZIP(), []int{2}
}

func (x *Classifier) GetSelector() *v1.Selector {
	if x != nil {
		return x.Selector
	}
	return nil
}

func (x *Classifier) GetRules() map[string]*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

// Rule describes a single Flow Classification Rule.
// Flow classification rule extracts a value from request metadata.
// More specifically, from `input`, which has the same spec as [Envoy's External Authorization Attribute Context][attribute-context].
// See <https://play.openpolicyagent.org/p/gU7vcLkc70> for an example input.
// There are two ways to define a flow classification rule:
// * Using a declarative extractor – suitable from simple cases, such as directly reading a value from header or a field from json body.
// * Rego expression.
//
// Performance note: It's recommended to use declarative extractors where possible, as they may be slightly performant than Rego expressions.
// [attribute-context](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto)
//
// Example:
// ```yaml
// Example of Declarative JSON extractor:
//   yaml:
//     extractor:
//       json:
//         from: request.http.body
//         pointer: /user/name
//     propagate: true
//     hidden: false
// Example of Rego module:
//   yaml:
//     rego:
//       query: data.user_from_cookie.user
//       source:
//         package: user_from_cookie
//         cookies: "split(input.attributes.request.http.headers.cookie, ';')"
//         cookie: "cookies[_]"
//         cookie.startswith: "('session=')"
//         session: "substring(cookie, count('session='), -1)"
//         parts: "split(session, '.')"
//         object: "json.unmarshal(base64url.decode(parts[0]))"
//         user: object.user
//     propagate: false
//     hidden: true
// ```
type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Source:
	//	*Rule_Extractor
	//	*Rule_Rego_
	Source isRule_Source `protobuf_oneof:"source"`
	// Decides if the created label should be applied to the whole flow (propagated in baggage) (default=true).
	Propagate *wrapperspb.BoolValue `protobuf:"bytes,3,opt,name=propagate,proto3" json:"propagate,omitempty"`
	// Decides if the created flow label should be hidden from the telemetry.
	Hidden bool `protobuf:"varint,4,opt,name=hidden,proto3" json:"hidden,omitempty"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_aperture_classification_v1_ruleset_proto_rawDescGZIP(), []int{3}
}

func (m *Rule) GetSource() isRule_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (x *Rule) GetExtractor() *Extractor {
	if x, ok := x.GetSource().(*Rule_Extractor); ok {
		return x.Extractor
	}
	return nil
}

func (x *Rule) GetRego() *Rule_Rego {
	if x, ok := x.GetSource().(*Rule_Rego_); ok {
		return x.Rego
	}
	return nil
}

func (x *Rule) GetPropagate() *wrapperspb.BoolValue {
	if x != nil {
		return x.Propagate
	}
	return nil
}

func (x *Rule) GetHidden() bool {
	if x != nil {
		return x.Hidden
	}
	return false
}

type isRule_Source interface {
	isRule_Source()
}

type Rule_Extractor struct {
	// High-level flow label declarative extractor.
	// Rego extractor extracts a value from the rego module.
	Extractor *Extractor `protobuf:"bytes,1,opt,name=extractor,proto3,oneof"`
}

type Rule_Rego_ struct {
	// Rego module to extract a value from the rego module.
	Rego *Rule_Rego `protobuf:"bytes,2,opt,name=rego,proto3,oneof"`
}

func (*Rule_Extractor) isRule_Source() {}

func (*Rule_Rego_) isRule_Source() {}

// Raw rego rules are compiled 1:1 to rego queries.
// High-level extractor-based rules are compiled into a single rego query.
type Rule_Rego struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Source code of the rego module.
	//
	// Note: Must include a "package" declaration.
	Source string `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	// Query string to extract a value (eg. `data.<mymodulename>.<variablename>`).
	//
	// Note: The module name must match the package name from the "source".
	Query string `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *Rule_Rego) Reset() {
	*x = Rule_Rego{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule_Rego) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule_Rego) ProtoMessage() {}

func (x *Rule_Rego) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_classification_v1_ruleset_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule_Rego.ProtoReflect.Descriptor instead.
func (*Rule_Rego) Descriptor() ([]byte, []int) {
	return file_aperture_classification_v1_ruleset_proto_rawDescGZIP(), []int{3, 0}
}

func (x *Rule_Rego) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Rule_Rego) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

var File_aperture_classification_v1_ruleset_proto protoreflect.FileDescriptor

var file_aperture_classification_v1_ruleset_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c,
	0x65, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x2a, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2a, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x55, 0x0a, 0x10, 0x41, 0x6c, 0x6c,
	0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x09, 0x61, 0x6c, 0x6c, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x24, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73,
	0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c,
	0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x08, 0x61, 0x6c, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73,
	0x22, 0xc0, 0x01, 0x0a, 0x08, 0x41, 0x6c, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x4f, 0x0a,
	0x09, 0x61, 0x6c, 0x6c, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x32, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73,
	0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c,
	0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x2e, 0x41, 0x6c, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x61, 0x6c, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x1a, 0x63,
	0x0a, 0x0d, 0x41, 0x6c, 0x6c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x3c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73,
	0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c,
	0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0xf4, 0x01, 0x0a, 0x0a, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69,
	0x65, 0x72, 0x12, 0x41, 0x0a, 0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x47, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x52, 0x75, 0x6c,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x1a, 0x5a,
	0x0a, 0x0a, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x36,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x9c, 0x02, 0x0a, 0x04, 0x52,
	0x75, 0x6c, 0x65, 0x12, 0x45, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x48, 0x00, 0x52,
	0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x3b, 0x0a, 0x04, 0x72, 0x65,
	0x67, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x6f, 0x48,
	0x00, 0x52, 0x04, 0x72, 0x65, 0x67, 0x6f, 0x12, 0x38, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x70, 0x61,
	0x67, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f,
	0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x70, 0x61, 0x67, 0x61, 0x74,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x1a, 0x34, 0x0a, 0x04, 0x52, 0x65, 0x67,
	0x6f, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x42,
	0x08, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x32, 0x73, 0x0a, 0x0c, 0x52, 0x75, 0x6c,
	0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x63, 0x0a, 0x08, 0x41, 0x6c, 0x6c,
	0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c, 0x6c, 0x52, 0x75,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x42, 0x94,
	0x02, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x42, 0x0c, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x5a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c,
	0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6c, 0x61,
	0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x41, 0x43, 0x58, 0xaa, 0x02, 0x1a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x1a, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x26,
	0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1c, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x3a, 0x3a, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_classification_v1_ruleset_proto_rawDescOnce sync.Once
	file_aperture_classification_v1_ruleset_proto_rawDescData = file_aperture_classification_v1_ruleset_proto_rawDesc
)

func file_aperture_classification_v1_ruleset_proto_rawDescGZIP() []byte {
	file_aperture_classification_v1_ruleset_proto_rawDescOnce.Do(func() {
		file_aperture_classification_v1_ruleset_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_classification_v1_ruleset_proto_rawDescData)
	})
	return file_aperture_classification_v1_ruleset_proto_rawDescData
}

var file_aperture_classification_v1_ruleset_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_aperture_classification_v1_ruleset_proto_goTypes = []interface{}{
	(*AllRulesResponse)(nil),     // 0: aperture.classification.v1.AllRulesResponse
	(*AllRules)(nil),             // 1: aperture.classification.v1.AllRules
	(*Classifier)(nil),           // 2: aperture.classification.v1.Classifier
	(*Rule)(nil),                 // 3: aperture.classification.v1.Rule
	nil,                          // 4: aperture.classification.v1.AllRules.AllRulesEntry
	nil,                          // 5: aperture.classification.v1.Classifier.RulesEntry
	(*Rule_Rego)(nil),            // 6: aperture.classification.v1.Rule.Rego
	(*v1.Selector)(nil),          // 7: aperture.policy.language.v1.Selector
	(*Extractor)(nil),            // 8: aperture.classification.v1.Extractor
	(*wrapperspb.BoolValue)(nil), // 9: google.protobuf.BoolValue
	(*emptypb.Empty)(nil),        // 10: google.protobuf.Empty
}
var file_aperture_classification_v1_ruleset_proto_depIdxs = []int32{
	1,  // 0: aperture.classification.v1.AllRulesResponse.all_rules:type_name -> aperture.classification.v1.AllRules
	4,  // 1: aperture.classification.v1.AllRules.all_rules:type_name -> aperture.classification.v1.AllRules.AllRulesEntry
	7,  // 2: aperture.classification.v1.Classifier.selector:type_name -> aperture.policy.language.v1.Selector
	5,  // 3: aperture.classification.v1.Classifier.rules:type_name -> aperture.classification.v1.Classifier.RulesEntry
	8,  // 4: aperture.classification.v1.Rule.extractor:type_name -> aperture.classification.v1.Extractor
	6,  // 5: aperture.classification.v1.Rule.rego:type_name -> aperture.classification.v1.Rule.Rego
	9,  // 6: aperture.classification.v1.Rule.propagate:type_name -> google.protobuf.BoolValue
	2,  // 7: aperture.classification.v1.AllRules.AllRulesEntry.value:type_name -> aperture.classification.v1.Classifier
	3,  // 8: aperture.classification.v1.Classifier.RulesEntry.value:type_name -> aperture.classification.v1.Rule
	10, // 9: aperture.classification.v1.RulesService.AllRules:input_type -> google.protobuf.Empty
	0,  // 10: aperture.classification.v1.RulesService.AllRules:output_type -> aperture.classification.v1.AllRulesResponse
	10, // [10:11] is the sub-list for method output_type
	9,  // [9:10] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_aperture_classification_v1_ruleset_proto_init() }
func file_aperture_classification_v1_ruleset_proto_init() {
	if File_aperture_classification_v1_ruleset_proto != nil {
		return
	}
	file_aperture_classification_v1_extractor_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_classification_v1_ruleset_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllRulesResponse); i {
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
		file_aperture_classification_v1_ruleset_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllRules); i {
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
		file_aperture_classification_v1_ruleset_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Classifier); i {
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
		file_aperture_classification_v1_ruleset_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule); i {
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
		file_aperture_classification_v1_ruleset_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule_Rego); i {
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
	file_aperture_classification_v1_ruleset_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Rule_Extractor)(nil),
		(*Rule_Rego_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_classification_v1_ruleset_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_aperture_classification_v1_ruleset_proto_goTypes,
		DependencyIndexes: file_aperture_classification_v1_ruleset_proto_depIdxs,
		MessageInfos:      file_aperture_classification_v1_ruleset_proto_msgTypes,
	}.Build()
	File_aperture_classification_v1_ruleset_proto = out.File
	file_aperture_classification_v1_ruleset_proto_rawDesc = nil
	file_aperture_classification_v1_ruleset_proto_goTypes = nil
	file_aperture_classification_v1_ruleset_proto_depIdxs = nil
}
