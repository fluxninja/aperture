// Code generated by protoc-gen-go-json. DO NOT EDIT.
// source: aperture/policy/sync/v1/infra_meter.proto

package syncv1

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// MarshalJSON implements json.Marshaler
func (msg *InfraMeterWrapper) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *InfraMeterWrapper) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}.Unmarshal(b, msg)
}
