// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: aperture/policy/language/v1/label_matcher.proto

package languagev1

import (
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

// Allows to define rules whether a map of
// [labels](/concepts/flow-control/flow-label.md)
// should be considered a match or not
//
// It provides three ways to define requirements:
// - match labels
// - match expressions
// - arbitrary expression
//
// If multiple requirements are set, they're all combined using the logical AND operator.
// An empty label matcher always matches.
type LabelMatcher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A map of {key,value} pairs representing labels to be matched.
	// A single {key,value} in the `match_labels` requires that the label `key` is present and equal to `value`.
	//
	// Note: The requirements are combined using the logical AND operator.
	MatchLabels map[string]string `protobuf:"bytes,1,rep,name=match_labels,json=matchLabels,proto3" json:"match_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// List of Kubernetes-style label matcher requirements.
	//
	// Note: The requirements are combined using the logical AND operator.
	MatchExpressions []*K8SLabelMatcherRequirement `protobuf:"bytes,2,rep,name=match_expressions,json=matchExpressions,proto3" json:"match_expressions,omitempty"`
	// An arbitrary expression to be evaluated on the labels.
	Expression *MatchExpression `protobuf:"bytes,3,opt,name=expression,proto3" json:"expression,omitempty"`
}

func (x *LabelMatcher) Reset() {
	*x = LabelMatcher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelMatcher) ProtoMessage() {}

func (x *LabelMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelMatcher.ProtoReflect.Descriptor instead.
func (*LabelMatcher) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP(), []int{0}
}

func (x *LabelMatcher) GetMatchLabels() map[string]string {
	if x != nil {
		return x.MatchLabels
	}
	return nil
}

func (x *LabelMatcher) GetMatchExpressions() []*K8SLabelMatcherRequirement {
	if x != nil {
		return x.MatchExpressions
	}
	return nil
}

func (x *LabelMatcher) GetExpression() *MatchExpression {
	if x != nil {
		return x.Expression
	}
	return nil
}

// Label selector requirement which is a selector that contains values, a key, and an operator that relates the key and values.
type K8SLabelMatcherRequirement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Label key that the selector applies to.
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty" validate:"required"` // @gotags: validate:"required"
	// Logical operator which represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator string `protobuf:"bytes,2,opt,name=operator,proto3" json:"operator,omitempty" validate:"oneof=In NotIn Exists DoesNotExists"` // @gotags: validate:"oneof=In NotIn Exists DoesNotExists"
	// An array of string values that relates to the key by an operator.
	// If the operator is In or NotIn, the values array must be non-empty.
	// If the operator is Exists or DoesNotExist, the values array must be empty.
	Values []string `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *K8SLabelMatcherRequirement) Reset() {
	*x = K8SLabelMatcherRequirement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *K8SLabelMatcherRequirement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SLabelMatcherRequirement) ProtoMessage() {}

func (x *K8SLabelMatcherRequirement) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SLabelMatcherRequirement.ProtoReflect.Descriptor instead.
func (*K8SLabelMatcherRequirement) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP(), []int{1}
}

func (x *K8SLabelMatcherRequirement) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *K8SLabelMatcherRequirement) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *K8SLabelMatcherRequirement) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

// Defines a `[map<string, string> → bool]` expression to be evaluated on labels
//
// MatchExpression has multiple variants, exactly one should be set.
//
// Example:
// ```yaml
// all:
//
//	of:
//	  - label_exists: foo
//	  - label_equals: { label = app, value = frobnicator }
//
// ```
type MatchExpression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gotags: validate:"required"
	//
	// Types that are assignable to Variant:
	//
	//	*MatchExpression_Not
	//	*MatchExpression_All
	//	*MatchExpression_Any
	//	*MatchExpression_LabelExists
	//	*MatchExpression_LabelEquals
	//	*MatchExpression_LabelMatches
	Variant isMatchExpression_Variant `protobuf_oneof:"variant" validate:"required"`
}

func (x *MatchExpression) Reset() {
	*x = MatchExpression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MatchExpression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchExpression) ProtoMessage() {}

func (x *MatchExpression) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MatchExpression.ProtoReflect.Descriptor instead.
func (*MatchExpression) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP(), []int{2}
}

func (m *MatchExpression) GetVariant() isMatchExpression_Variant {
	if m != nil {
		return m.Variant
	}
	return nil
}

func (x *MatchExpression) GetNot() *MatchExpression {
	if x, ok := x.GetVariant().(*MatchExpression_Not); ok {
		return x.Not
	}
	return nil
}

func (x *MatchExpression) GetAll() *MatchExpression_List {
	if x, ok := x.GetVariant().(*MatchExpression_All); ok {
		return x.All
	}
	return nil
}

func (x *MatchExpression) GetAny() *MatchExpression_List {
	if x, ok := x.GetVariant().(*MatchExpression_Any); ok {
		return x.Any
	}
	return nil
}

func (x *MatchExpression) GetLabelExists() string {
	if x, ok := x.GetVariant().(*MatchExpression_LabelExists); ok {
		return x.LabelExists
	}
	return ""
}

func (x *MatchExpression) GetLabelEquals() *EqualsMatchExpression {
	if x, ok := x.GetVariant().(*MatchExpression_LabelEquals); ok {
		return x.LabelEquals
	}
	return nil
}

func (x *MatchExpression) GetLabelMatches() *MatchesMatchExpression {
	if x, ok := x.GetVariant().(*MatchExpression_LabelMatches); ok {
		return x.LabelMatches
	}
	return nil
}

type isMatchExpression_Variant interface {
	isMatchExpression_Variant()
}

type MatchExpression_Not struct {
	// The expression negates the result of sub expression.
	Not *MatchExpression `protobuf:"bytes,1,opt,name=not,proto3,oneof"`
}

type MatchExpression_All struct {
	// The expression is true when all sub expressions are true.
	All *MatchExpression_List `protobuf:"bytes,2,opt,name=all,proto3,oneof"`
}

type MatchExpression_Any struct {
	// The expression is true when any sub expression is true.
	Any *MatchExpression_List `protobuf:"bytes,3,opt,name=any,proto3,oneof"`
}

type MatchExpression_LabelExists struct {
	// The expression is true when label with given name exists.
	LabelExists string `protobuf:"bytes,4,opt,name=label_exists,json=labelExists,proto3,oneof" validate:"required"` // @gotags: validate:"required"
}

type MatchExpression_LabelEquals struct {
	// The expression is true when label value equals given value.
	LabelEquals *EqualsMatchExpression `protobuf:"bytes,5,opt,name=label_equals,json=labelEquals,proto3,oneof"`
}

type MatchExpression_LabelMatches struct {
	// The expression is true when label matches given regular expression.
	LabelMatches *MatchesMatchExpression `protobuf:"bytes,6,opt,name=label_matches,json=labelMatches,proto3,oneof"`
}

func (*MatchExpression_Not) isMatchExpression_Variant() {}

func (*MatchExpression_All) isMatchExpression_Variant() {}

func (*MatchExpression_Any) isMatchExpression_Variant() {}

func (*MatchExpression_LabelExists) isMatchExpression_Variant() {}

func (*MatchExpression_LabelEquals) isMatchExpression_Variant() {}

func (*MatchExpression_LabelMatches) isMatchExpression_Variant() {}

// Label selector expression of the equal form `label == value`.
type EqualsMatchExpression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the label to equal match the value.
	Label string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty" validate:"required"` // @gotags: validate:"required"
	// Exact value that the label should be equal to.
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *EqualsMatchExpression) Reset() {
	*x = EqualsMatchExpression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EqualsMatchExpression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EqualsMatchExpression) ProtoMessage() {}

func (x *EqualsMatchExpression) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EqualsMatchExpression.ProtoReflect.Descriptor instead.
func (*EqualsMatchExpression) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP(), []int{3}
}

func (x *EqualsMatchExpression) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *EqualsMatchExpression) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Label selector expression of the form `label matches regex`.
type MatchesMatchExpression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the label to match the regular expression.
	Label string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty" validate:"required"` // @gotags: validate:"required"
	// Regular expression that should match the label value.
	// It uses [Go's regular expression syntax](https://github.com/google/re2/wiki/Syntax).
	Regex string `protobuf:"bytes,2,opt,name=regex,proto3" json:"regex,omitempty" validate:"required"` // @gotags: validate:"required"
}

func (x *MatchesMatchExpression) Reset() {
	*x = MatchesMatchExpression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MatchesMatchExpression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchesMatchExpression) ProtoMessage() {}

func (x *MatchesMatchExpression) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MatchesMatchExpression.ProtoReflect.Descriptor instead.
func (*MatchesMatchExpression) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP(), []int{4}
}

func (x *MatchesMatchExpression) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *MatchesMatchExpression) GetRegex() string {
	if x != nil {
		return x.Regex
	}
	return ""
}

// List of MatchExpressions that's used for all or any matching
//
// for example, `{any: {of: [expr1, expr2]}}`.
type MatchExpression_List struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of sub expressions of the match expression.
	Of []*MatchExpression `protobuf:"bytes,1,rep,name=of,proto3" json:"of,omitempty"`
}

func (x *MatchExpression_List) Reset() {
	*x = MatchExpression_List{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MatchExpression_List) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchExpression_List) ProtoMessage() {}

func (x *MatchExpression_List) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_label_matcher_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MatchExpression_List.ProtoReflect.Descriptor instead.
func (*MatchExpression_List) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP(), []int{2, 0}
}

func (x *MatchExpression_List) GetOf() []*MatchExpression {
	if x != nil {
		return x.Of
	}
	return nil
}

var File_aperture_policy_language_v1_label_matcher_proto protoreflect.FileDescriptor

var file_aperture_policy_language_v1_label_matcher_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1b, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x22, 0xe1,
	0x02, 0x0a, 0x0c, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12,
	0x5d, 0x0a, 0x0c, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3a, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72,
	0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0b, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x64,
	0x0a, 0x11, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x38, 0x73, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x10, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x4c, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x1a, 0x3e, 0x0a, 0x10, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x62, 0x0a, 0x1a, 0x4b, 0x38, 0x73, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x8c, 0x04, 0x0a, 0x0f, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x03, 0x6e, 0x6f,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x03, 0x6e, 0x6f, 0x74, 0x12, 0x45, 0x0a, 0x03,
	0x61, 0x6c, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x00, 0x52, 0x03,
	0x61, 0x6c, 0x6c, 0x12, 0x45, 0x0a, 0x03, 0x61, 0x6e, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x31, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x48, 0x00, 0x52, 0x03, 0x61, 0x6e, 0x79, 0x12, 0x23, 0x0a, 0x0c, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x0b, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12,
	0x57, 0x0a, 0x0c, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x5f, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45,
	0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0b, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x12, 0x5a, 0x0a, 0x0d, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x33, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x73, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0c, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x73, 0x1a, 0x44, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3c, 0x0a, 0x02,
	0x6f, 0x66, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x02, 0x6f, 0x66, 0x42, 0x09, 0x0a, 0x07, 0x76, 0x61,
	0x72, 0x69, 0x61, 0x6e, 0x74, 0x22, 0x43, 0x0a, 0x15, 0x45, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x44, 0x0a, 0x16, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x73, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65,
	0x67, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x67, 0x65, 0x78,
	0x42, 0xae, 0x02, 0x0a, 0x33, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e,
	0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x55, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69,
	0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70,
	0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61,
	0x67, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x4c, 0xaa, 0x02, 0x1b, 0x41, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x4c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1b, 0x41, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x4c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x3a, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x3a, 0x3a, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_language_v1_label_matcher_proto_rawDescOnce sync.Once
	file_aperture_policy_language_v1_label_matcher_proto_rawDescData = file_aperture_policy_language_v1_label_matcher_proto_rawDesc
)

func file_aperture_policy_language_v1_label_matcher_proto_rawDescGZIP() []byte {
	file_aperture_policy_language_v1_label_matcher_proto_rawDescOnce.Do(func() {
		file_aperture_policy_language_v1_label_matcher_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_language_v1_label_matcher_proto_rawDescData)
	})
	return file_aperture_policy_language_v1_label_matcher_proto_rawDescData
}

var file_aperture_policy_language_v1_label_matcher_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_aperture_policy_language_v1_label_matcher_proto_goTypes = []interface{}{
	(*LabelMatcher)(nil),               // 0: aperture.policy.language.v1.LabelMatcher
	(*K8SLabelMatcherRequirement)(nil), // 1: aperture.policy.language.v1.K8sLabelMatcherRequirement
	(*MatchExpression)(nil),            // 2: aperture.policy.language.v1.MatchExpression
	(*EqualsMatchExpression)(nil),      // 3: aperture.policy.language.v1.EqualsMatchExpression
	(*MatchesMatchExpression)(nil),     // 4: aperture.policy.language.v1.MatchesMatchExpression
	nil,                                // 5: aperture.policy.language.v1.LabelMatcher.MatchLabelsEntry
	(*MatchExpression_List)(nil),       // 6: aperture.policy.language.v1.MatchExpression.List
}
var file_aperture_policy_language_v1_label_matcher_proto_depIdxs = []int32{
	5, // 0: aperture.policy.language.v1.LabelMatcher.match_labels:type_name -> aperture.policy.language.v1.LabelMatcher.MatchLabelsEntry
	1, // 1: aperture.policy.language.v1.LabelMatcher.match_expressions:type_name -> aperture.policy.language.v1.K8sLabelMatcherRequirement
	2, // 2: aperture.policy.language.v1.LabelMatcher.expression:type_name -> aperture.policy.language.v1.MatchExpression
	2, // 3: aperture.policy.language.v1.MatchExpression.not:type_name -> aperture.policy.language.v1.MatchExpression
	6, // 4: aperture.policy.language.v1.MatchExpression.all:type_name -> aperture.policy.language.v1.MatchExpression.List
	6, // 5: aperture.policy.language.v1.MatchExpression.any:type_name -> aperture.policy.language.v1.MatchExpression.List
	3, // 6: aperture.policy.language.v1.MatchExpression.label_equals:type_name -> aperture.policy.language.v1.EqualsMatchExpression
	4, // 7: aperture.policy.language.v1.MatchExpression.label_matches:type_name -> aperture.policy.language.v1.MatchesMatchExpression
	2, // 8: aperture.policy.language.v1.MatchExpression.List.of:type_name -> aperture.policy.language.v1.MatchExpression
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_aperture_policy_language_v1_label_matcher_proto_init() }
func file_aperture_policy_language_v1_label_matcher_proto_init() {
	if File_aperture_policy_language_v1_label_matcher_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_language_v1_label_matcher_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelMatcher); i {
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
		file_aperture_policy_language_v1_label_matcher_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*K8SLabelMatcherRequirement); i {
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
		file_aperture_policy_language_v1_label_matcher_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MatchExpression); i {
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
		file_aperture_policy_language_v1_label_matcher_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EqualsMatchExpression); i {
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
		file_aperture_policy_language_v1_label_matcher_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MatchesMatchExpression); i {
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
		file_aperture_policy_language_v1_label_matcher_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MatchExpression_List); i {
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
	file_aperture_policy_language_v1_label_matcher_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*MatchExpression_Not)(nil),
		(*MatchExpression_All)(nil),
		(*MatchExpression_Any)(nil),
		(*MatchExpression_LabelExists)(nil),
		(*MatchExpression_LabelEquals)(nil),
		(*MatchExpression_LabelMatches)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_policy_language_v1_label_matcher_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_language_v1_label_matcher_proto_goTypes,
		DependencyIndexes: file_aperture_policy_language_v1_label_matcher_proto_depIdxs,
		MessageInfos:      file_aperture_policy_language_v1_label_matcher_proto_msgTypes,
	}.Build()
	File_aperture_policy_language_v1_label_matcher_proto = out.File
	file_aperture_policy_language_v1_label_matcher_proto_rawDesc = nil
	file_aperture_policy_language_v1_label_matcher_proto_goTypes = nil
	file_aperture_policy_language_v1_label_matcher_proto_depIdxs = nil
}
