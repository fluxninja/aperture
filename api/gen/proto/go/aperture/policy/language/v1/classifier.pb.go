// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: aperture/policy/language/v1/classifier.proto

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

// Set of classification rules sharing a common selector
//
// :::info
//
// See also [Classifier overview](/concepts/integrations/flow-control/resources/classifier.md).
//
// :::
// Example
// ```yaml
// flow_selector:
//   service_selector:
//      agent_group: demoapp
//      service: service1-demo-app.demoapp.svc.cluster.local
//   flow_matcher:
//      control_point: ingress
//      label_matcher:
//         match_labels:
//           user_tier: gold
//         match_expressions:
//           - key: user_type
//             operator: In
// rules:
//   user:
//    extractor:
//      from: request.http.headers.user-agent
//   telemetry: false
// ```
type Classifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Defines where to apply the flow classification rule.
	FlowSelector *FlowSelector `protobuf:"bytes,1,opt,name=flow_selector,json=flowSelector,proto3" json:"flow_selector,omitempty" validate:"required"` // @gotags: validate:"required"
	// A map of {key, value} pairs mapping from
	// [flow label](/concepts/integrations/flow-control/flow-label.md) keys to rules that define
	// how to extract and propagate flow labels with that key.
	Rules map[string]*Rule `protobuf:"bytes,2,rep,name=rules,proto3" json:"rules,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" validate:"required,gt=0,dive,keys,required,endkeys,required"` // @gotags: validate:"required,gt=0,dive,keys,required,endkeys,required"
}

func (x *Classifier) Reset() {
	*x = Classifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Classifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Classifier) ProtoMessage() {}

func (x *Classifier) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[0]
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
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{0}
}

func (x *Classifier) GetFlowSelector() *FlowSelector {
	if x != nil {
		return x.FlowSelector
	}
	return nil
}

func (x *Classifier) GetRules() map[string]*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

// Rule describes a single classification Rule
//
// Classification rule extracts a value from request metadata.
// More specifically, from `input`, which has the same spec as [Envoy's External Authorization Attribute Context][attribute-context].
// See https://play.openpolicyagent.org/p/gU7vcLkc70 for an example input.
// There are two ways to define a flow classification rule:
// * Using a declarative extractor – suitable from simple cases, such as directly reading a value from header or a field from json body.
// * Rego expression.
//
// Performance note: It's recommended to use declarative extractors where possible, as they may be slightly performant than Rego expressions.
//
// Example of Declarative JSON extractor:
// ```yaml
// extractor:
//   json:
//     from: request.http.body
//     pointer: /user/name
// ```
//
// Example of Rego module which also disables telemetry visibility of label:
// ```yaml
// rego:
//   query: data.user_from_cookie.user
//   source: |
//     package user_from_cookie
//     cookies := split(input.attributes.request.http.headers.cookie, "; ")
//     user := user {
//         cookie := cookies[_]
//         startswith(cookie, "session=")
//         session := substring(cookie, count("session="), -1)
//         parts := split(session, ".")
//         object := json.unmarshal(base64url.decode(parts[0]))
//         user := object.user
//     }
// telemetry: false
// ```
// [attribute-context]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto
type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Source:
	//	*Rule_Extractor
	//	*Rule_Rego_
	Source isRule_Source `protobuf_oneof:"source"`
	// Decides if the created flow label should be available as an attribute in OLAP telemetry and
	// propagated in [baggage](/concepts/integrations/flow-control/flow-label.md#baggage)
	//
	// :::note
	//
	// The flow label is always accessible in Aperture Policies regardless of this setting.
	//
	// :::
	//
	// :::caution
	//
	// When using [FluxNinja ARC plugin](arc/plugin.md), telemetry enabled
	// labels are sent to FluxNinja ARC for observability. Telemetry should be disabled for
	// sensitive labels.
	//
	// :::
	Telemetry bool `protobuf:"varint,3,opt,name=telemetry,proto3" json:"telemetry,omitempty" default:"true"` // @gotags: default:"true"
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[1]
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
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{1}
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

func (x *Rule) GetTelemetry() bool {
	if x != nil {
		return x.Telemetry
	}
	return false
}

type isRule_Source interface {
	isRule_Source()
}

type Rule_Extractor struct {
	// High-level declarative extractor.
	Extractor *Extractor `protobuf:"bytes,1,opt,name=extractor,proto3,oneof"`
}

type Rule_Rego_ struct {
	// Rego module to extract a value from.
	Rego *Rule_Rego `protobuf:"bytes,2,opt,name=rego,proto3,oneof"`
}

func (*Rule_Extractor) isRule_Source() {}

func (*Rule_Rego_) isRule_Source() {}

// Defines a high-level way to specify how to extract a flow label value given http request metadata, without a need to write rego code
//
// There are multiple variants of extractor, specify exactly one.
type Extractor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Variant:
	//	*Extractor_From
	//	*Extractor_Json
	//	*Extractor_Address
	//	*Extractor_Jwt
	//	*Extractor_PathTemplates
	Variant isExtractor_Variant `protobuf_oneof:"variant"`
}

func (x *Extractor) Reset() {
	*x = Extractor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Extractor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Extractor) ProtoMessage() {}

func (x *Extractor) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Extractor.ProtoReflect.Descriptor instead.
func (*Extractor) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{2}
}

func (m *Extractor) GetVariant() isExtractor_Variant {
	if m != nil {
		return m.Variant
	}
	return nil
}

func (x *Extractor) GetFrom() string {
	if x, ok := x.GetVariant().(*Extractor_From); ok {
		return x.From
	}
	return ""
}

func (x *Extractor) GetJson() *JSONExtractor {
	if x, ok := x.GetVariant().(*Extractor_Json); ok {
		return x.Json
	}
	return nil
}

func (x *Extractor) GetAddress() *AddressExtractor {
	if x, ok := x.GetVariant().(*Extractor_Address); ok {
		return x.Address
	}
	return nil
}

func (x *Extractor) GetJwt() *JWTExtractor {
	if x, ok := x.GetVariant().(*Extractor_Jwt); ok {
		return x.Jwt
	}
	return nil
}

func (x *Extractor) GetPathTemplates() *PathTemplateMatcher {
	if x, ok := x.GetVariant().(*Extractor_PathTemplates); ok {
		return x.PathTemplates
	}
	return nil
}

type isExtractor_Variant interface {
	isExtractor_Variant()
}

type Extractor_From struct {
	// Use an attribute with no conversion
	//
	// Attribute path is a dot-separated path to attribute.
	//
	// Should be either:
	// * one of the fields of [Attribute Context][attribute-context], or
	// * a special "request.http.bearer" pseudo-attribute.
	// Eg. "request.http.method" or "request.http.header.user-agent"
	//
	// Note: The same attribute path syntax is shared by other extractor variants,
	// wherever attribute path is needed in their "from" syntax.
	//
	// Example:
	// ```yaml
	// from: request.http.headers.user-agent
	// ```
	// [attribute-context]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto
	From string `protobuf:"bytes,1,opt,name=from,proto3,oneof"`
}

type Extractor_Json struct {
	// Deserialize a json, and extract one of the fields.
	Json *JSONExtractor `protobuf:"bytes,2,opt,name=json,proto3,oneof"`
}

type Extractor_Address struct {
	// Display an address as a single string - `<ip>:<port>`.
	Address *AddressExtractor `protobuf:"bytes,3,opt,name=address,proto3,oneof"`
}

type Extractor_Jwt struct {
	// Parse the attribute as JWT and read the payload.
	Jwt *JWTExtractor `protobuf:"bytes,4,opt,name=jwt,proto3,oneof"`
}

type Extractor_PathTemplates struct {
	// Match HTTP Path to given path templates.
	PathTemplates *PathTemplateMatcher `protobuf:"bytes,5,opt,name=path_templates,json=pathTemplates,proto3,oneof"`
}

func (*Extractor_From) isExtractor_Variant() {}

func (*Extractor_Json) isExtractor_Variant() {}

func (*Extractor_Address) isExtractor_Variant() {}

func (*Extractor_Jwt) isExtractor_Variant() {}

func (*Extractor_PathTemplates) isExtractor_Variant() {}

// Deserialize a json, and extract one of the fields
//
// Example:
// ```yaml
// from: request.http.body
// pointer: /user/name
// ```
type JSONExtractor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Attribute path pointing to some strings - eg. "request.http.body".
	From string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty" validate:"required"` //@gotags: validate:"required"
	// Json pointer represents a parsed json pointer which allows to select a specified field from the json payload.
	//
	// Note: Uses [json pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
	// eg. `/foo/bar`. If the pointer points into an object, it'd be stringified.
	Pointer string `protobuf:"bytes,2,opt,name=pointer,proto3" json:"pointer,omitempty"`
}

func (x *JSONExtractor) Reset() {
	*x = JSONExtractor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JSONExtractor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JSONExtractor) ProtoMessage() {}

func (x *JSONExtractor) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JSONExtractor.ProtoReflect.Descriptor instead.
func (*JSONExtractor) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{3}
}

func (x *JSONExtractor) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *JSONExtractor) GetPointer() string {
	if x != nil {
		return x.Pointer
	}
	return ""
}

// Display an [Address][ext-authz-address] as a single string, eg. `<ip>:<port>`
//
// IP addresses in attribute context are defined as objects with separate ip and port fields.
// This is a helper to display an address as a single string.
//
// Note: Use with care, as it might accidentally introduce a high-cardinality flow label values.
//
// [ext-authz-address]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/address.proto#config-core-v3-address
//
// Example:
// ```yaml
// from: "source.address # or destination.address"
// ```
type AddressExtractor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Attribute path pointing to some string - eg. "source.address".
	From string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty" validate:"required"` //@gotags: validate:"required"
}

func (x *AddressExtractor) Reset() {
	*x = AddressExtractor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressExtractor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressExtractor) ProtoMessage() {}

func (x *AddressExtractor) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressExtractor.ProtoReflect.Descriptor instead.
func (*AddressExtractor) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{4}
}

func (x *AddressExtractor) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

// Parse the attribute as JWT and read the payload
//
// Specify a field to be extracted from payload using "json_pointer".
//
// Note: The signature is not verified against the secret (we're assuming there's some
// other parts of the system that handles such verification).
//
// Example:
// ```yaml
// from: request.http.bearer
// json_pointer: /user/email
// ```
type JWTExtractor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Jwt token can be pulled from any input attribute, but most likely you'd want to use "request.http.bearer".
	From string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty" validate:"required"` //@gotags: validate:"required"
	// Json pointer allowing to select a specified field from the json payload.
	//
	// Note: Uses [json pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
	// eg. `/foo/bar`. If the pointer points into an object, it'd be stringified.
	JsonPointer string `protobuf:"bytes,2,opt,name=json_pointer,json=jsonPointer,proto3" json:"json_pointer,omitempty"`
}

func (x *JWTExtractor) Reset() {
	*x = JWTExtractor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JWTExtractor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JWTExtractor) ProtoMessage() {}

func (x *JWTExtractor) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JWTExtractor.ProtoReflect.Descriptor instead.
func (*JWTExtractor) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{5}
}

func (x *JWTExtractor) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *JWTExtractor) GetJsonPointer() string {
	if x != nil {
		return x.JsonPointer
	}
	return ""
}

// Matches HTTP Path to given path templates
//
// HTTP path will be matched against given path templates.
// If a match occurs, the value associated with the path template will be treated as a result.
// In case of multiple path templates matching, the most specific one will be chosen.
type PathTemplateMatcher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Template value keys are OpenAPI-inspired path templates.
	//
	// * Static path segment `/foo` matches a path segment exactly
	// * `/{param}` matches arbitrary path segment.
	//   (The param name is ignored and can be omitted (`{}`))
	// * The parameter must cover whole segment.
	// * Additionally, path template can end with `/*` wildcard to match
	//   arbitrary number of trailing segments (0 or more).
	// * Multiple consecutive `/` are ignored, as well as trailing `/`.
	// * Parametrized path segments must come after static segments.
	// * `*`, if present, must come last.
	// * Most specific template "wins" (`/foo` over `/{}` and `/{}` over `/*`).
	//
	// See also <https://swagger.io/specification/#path-templating-matching>
	//
	// Example:
	// ```yaml
	// /register: register
	// "/user/{userId}": user
	// /static/*: other
	// ```
	TemplateValues map[string]string `protobuf:"bytes,1,rep,name=template_values,json=templateValues,proto3" json:"template_values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" validate:"gt=0,dive,keys,required,endkeys,required"` // @gotags: validate:"gt=0,dive,keys,required,endkeys,required"
}

func (x *PathTemplateMatcher) Reset() {
	*x = PathTemplateMatcher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PathTemplateMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PathTemplateMatcher) ProtoMessage() {}

func (x *PathTemplateMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PathTemplateMatcher.ProtoReflect.Descriptor instead.
func (*PathTemplateMatcher) Descriptor() ([]byte, []int) {
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{6}
}

func (x *PathTemplateMatcher) GetTemplateValues() map[string]string {
	if x != nil {
		return x.TemplateValues
	}
	return nil
}

// Raw rego rules are compiled 1:1 to rego queries
//
// High-level extractor-based rules are compiled into a single rego query.
type Rule_Rego struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Source code of the rego module.
	//
	// Note: Must include a "package" declaration.
	Source string `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty" validate:"required"` // @gotags: validate:"required"
	// Query string to extract a value (eg. `data.<mymodulename>.<variablename>`).
	//
	// Note: The module name must match the package name from the "source".
	Query string `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty" validate:"required"` // @gotags: validate:"required"
}

func (x *Rule_Rego) Reset() {
	*x = Rule_Rego{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule_Rego) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule_Rego) ProtoMessage() {}

func (x *Rule_Rego) ProtoReflect() protoreflect.Message {
	mi := &file_aperture_policy_language_v1_classifier_proto_msgTypes[8]
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
	return file_aperture_policy_language_v1_classifier_proto_rawDescGZIP(), []int{1, 0}
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

var File_aperture_policy_language_v1_classifier_proto protoreflect.FileDescriptor

var file_aperture_policy_language_v1_classifier_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c,
	0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x2a, 0x61, 0x70, 0x65,
	0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x83, 0x02, 0x0a, 0x0a, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x4e, 0x0a, 0x0d, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x73,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0c, 0x66, 0x6c, 0x6f, 0x77, 0x53, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x48, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65,
	0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73,
	0x1a, 0x5b, 0x0a, 0x0a, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x37, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75,
	0x6c, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xea, 0x01,
	0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x46, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x65, 0x72,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x48, 0x00, 0x52, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x3c,
	0x0a, 0x04, 0x72, 0x65, 0x67, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e,
	0x52, 0x65, 0x67, 0x6f, 0x48, 0x00, 0x52, 0x04, 0x72, 0x65, 0x67, 0x6f, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x09, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x1a, 0x34, 0x0a, 0x04, 0x52, 0x65,
	0x67, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x42, 0x08, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0xd3, 0x02, 0x0a, 0x09, 0x45,
	0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x40,
	0x0a, 0x04, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x61,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x53, 0x4f, 0x4e, 0x45,
	0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x04, 0x6a, 0x73, 0x6f, 0x6e,
	0x12, 0x49, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x48, 0x00, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3d, 0x0a, 0x03, 0x6a,
	0x77, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x57, 0x54, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x12, 0x59, 0x0a, 0x0e, 0x70, 0x61,
	0x74, 0x68, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x30, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x61, 0x74, 0x68, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x72, 0x48, 0x00, 0x52, 0x0d, 0x70, 0x61, 0x74, 0x68, 0x54, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74,
	0x22, 0x3d, 0x0a, 0x0d, 0x4a, 0x53, 0x4f, 0x4e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x22,
	0x26, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x22, 0x45, 0x0a, 0x0c, 0x4a, 0x57, 0x54, 0x45, 0x78,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x6a,
	0x73, 0x6f, 0x6e, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6a, 0x73, 0x6f, 0x6e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x22, 0xc7,
	0x01, 0x0a, 0x13, 0x50, 0x61, 0x74, 0x68, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x6d, 0x0a, 0x0f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x44, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61,
	0x74, 0x68, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x1a, 0x41, 0x0a, 0x13, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0xac, 0x02, 0x0a, 0x33, 0x63, 0x6f, 0x6d,
	0x2e, 0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x64, 0x2e, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x2e, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x42, 0x0f, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x55, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x6c, 0x75, 0x78, 0x6e, 0x69, 0x6e, 0x6a, 0x61, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75,
	0x72, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x3b,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x50, 0x4c,
	0xaa, 0x02, 0x1b, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x1b, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x5c, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x27, 0x41,
	0x70, 0x65, 0x72, 0x74, 0x75, 0x72, 0x65, 0x5c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5c, 0x4c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1e, 0x41, 0x70, 0x65, 0x72, 0x74, 0x75, 0x72,
	0x65, 0x3a, 0x3a, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3a, 0x3a, 0x4c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aperture_policy_language_v1_classifier_proto_rawDescOnce sync.Once
	file_aperture_policy_language_v1_classifier_proto_rawDescData = file_aperture_policy_language_v1_classifier_proto_rawDesc
)

func file_aperture_policy_language_v1_classifier_proto_rawDescGZIP() []byte {
	file_aperture_policy_language_v1_classifier_proto_rawDescOnce.Do(func() {
		file_aperture_policy_language_v1_classifier_proto_rawDescData = protoimpl.X.CompressGZIP(file_aperture_policy_language_v1_classifier_proto_rawDescData)
	})
	return file_aperture_policy_language_v1_classifier_proto_rawDescData
}

var file_aperture_policy_language_v1_classifier_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_aperture_policy_language_v1_classifier_proto_goTypes = []interface{}{
	(*Classifier)(nil),          // 0: aperture.policy.language.v1.Classifier
	(*Rule)(nil),                // 1: aperture.policy.language.v1.Rule
	(*Extractor)(nil),           // 2: aperture.policy.language.v1.Extractor
	(*JSONExtractor)(nil),       // 3: aperture.policy.language.v1.JSONExtractor
	(*AddressExtractor)(nil),    // 4: aperture.policy.language.v1.AddressExtractor
	(*JWTExtractor)(nil),        // 5: aperture.policy.language.v1.JWTExtractor
	(*PathTemplateMatcher)(nil), // 6: aperture.policy.language.v1.PathTemplateMatcher
	nil,                         // 7: aperture.policy.language.v1.Classifier.RulesEntry
	(*Rule_Rego)(nil),           // 8: aperture.policy.language.v1.Rule.Rego
	nil,                         // 9: aperture.policy.language.v1.PathTemplateMatcher.TemplateValuesEntry
	(*FlowSelector)(nil),        // 10: aperture.policy.language.v1.FlowSelector
}
var file_aperture_policy_language_v1_classifier_proto_depIdxs = []int32{
	10, // 0: aperture.policy.language.v1.Classifier.flow_selector:type_name -> aperture.policy.language.v1.FlowSelector
	7,  // 1: aperture.policy.language.v1.Classifier.rules:type_name -> aperture.policy.language.v1.Classifier.RulesEntry
	2,  // 2: aperture.policy.language.v1.Rule.extractor:type_name -> aperture.policy.language.v1.Extractor
	8,  // 3: aperture.policy.language.v1.Rule.rego:type_name -> aperture.policy.language.v1.Rule.Rego
	3,  // 4: aperture.policy.language.v1.Extractor.json:type_name -> aperture.policy.language.v1.JSONExtractor
	4,  // 5: aperture.policy.language.v1.Extractor.address:type_name -> aperture.policy.language.v1.AddressExtractor
	5,  // 6: aperture.policy.language.v1.Extractor.jwt:type_name -> aperture.policy.language.v1.JWTExtractor
	6,  // 7: aperture.policy.language.v1.Extractor.path_templates:type_name -> aperture.policy.language.v1.PathTemplateMatcher
	9,  // 8: aperture.policy.language.v1.PathTemplateMatcher.template_values:type_name -> aperture.policy.language.v1.PathTemplateMatcher.TemplateValuesEntry
	1,  // 9: aperture.policy.language.v1.Classifier.RulesEntry.value:type_name -> aperture.policy.language.v1.Rule
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_aperture_policy_language_v1_classifier_proto_init() }
func file_aperture_policy_language_v1_classifier_proto_init() {
	if File_aperture_policy_language_v1_classifier_proto != nil {
		return
	}
	file_aperture_policy_language_v1_selector_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_aperture_policy_language_v1_classifier_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Extractor); i {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JSONExtractor); i {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddressExtractor); i {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JWTExtractor); i {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PathTemplateMatcher); i {
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
		file_aperture_policy_language_v1_classifier_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
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
	file_aperture_policy_language_v1_classifier_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Rule_Extractor)(nil),
		(*Rule_Rego_)(nil),
	}
	file_aperture_policy_language_v1_classifier_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Extractor_From)(nil),
		(*Extractor_Json)(nil),
		(*Extractor_Address)(nil),
		(*Extractor_Jwt)(nil),
		(*Extractor_PathTemplates)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aperture_policy_language_v1_classifier_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aperture_policy_language_v1_classifier_proto_goTypes,
		DependencyIndexes: file_aperture_policy_language_v1_classifier_proto_depIdxs,
		MessageInfos:      file_aperture_policy_language_v1_classifier_proto_msgTypes,
	}.Build()
	File_aperture_policy_language_v1_classifier_proto = out.File
	file_aperture_policy_language_v1_classifier_proto_rawDesc = nil
	file_aperture_policy_language_v1_classifier_proto_goTypes = nil
	file_aperture_policy_language_v1_classifier_proto_depIdxs = nil
}
