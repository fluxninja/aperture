// Code generated by protoc-gen-go-json. DO NOT EDIT.
// source: aperture/policy/wrappers/v1/decisions.proto

package wrappersv1

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// MarshalJSON implements json.Marshaler
func (msg *LoadShedDecsisionWrapper) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *LoadShedDecsisionWrapper) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *TokensDecisionWrapper) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *TokensDecisionWrapper) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *RateLimiterDecisionWrapper) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *RateLimiterDecisionWrapper) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}
