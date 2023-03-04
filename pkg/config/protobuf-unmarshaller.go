package config

import (
	"errors"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
)

// ProtobufUnmarshaller is an unmarshaller that can unmarshalls
// protobuf-encoded data into messages.
//
// No defaults handling nor validation is performed.
//
// Note: Unmarshalling keys or unmarshalling into non-protobuf structs in not
// supported by this unmarshaller.
type ProtobufUnmarshaller struct {
	bytes atomic.Value // holding []byte
}

// Make sure ProtobufUnmarshaller implements Unmarshaller.
var _ Unmarshaller = &ProtobufUnmarshaller{}

// UnmarshalKey is not supported by ProtobufUnmarshaler.
func (u *ProtobufUnmarshaller) UnmarshalKey(key string, output interface{}) error {
	panic("unimplemented")
}

// IsSet is not supported by ProtobufUnmarshaler.
func (u *ProtobufUnmarshaller) IsSet(key string) bool { panic("unimplemented") }

// Get is not supported by ProtobufUnmarshaler.
func (u *ProtobufUnmarshaller) Get(key string) interface{} { panic("unimplemented") }

// Reload sets the protobuf-encoded bytes.
//
// Previous state is completely forgotten.
func (u *ProtobufUnmarshaller) Reload(bytes []byte) error {
	if bytes == nil {
		return errors.New("attempt to reload with nil bytes")
	}
	u.bytes.Store(bytes)
	return nil
}

// Marshal the config into bytes.
func (u *ProtobufUnmarshaller) Marshal() ([]byte, error) {
	bytes, err := u.bytes.Load().([]byte)
	if !err {
		return nil, errors.New("attempt to marshal nil bytes")
	}
	return bytes, nil
}

// Unmarshal unmarshals previously set protobuf-encoded bytes into output.
//
// Output should be a proto.Message.
func (u *ProtobufUnmarshaller) Unmarshal(output interface{}) error {
	msg, ok := output.(proto.Message)
	if !ok {
		return errors.New("attempt to unmarshal into non-proto.Message")
	}

	bytes, _ := u.bytes.Load().([]byte)

	return proto.Unmarshal(bytes, msg)
}

// NewProtobufUnmarshaller crates a new ProtobufUnmarshaller.
func NewProtobufUnmarshaller(bytes []byte) (Unmarshaller, error) {
	pu := &ProtobufUnmarshaller{}
	if bytes != nil {
		ok := pu.Reload(bytes)
		if ok != nil {
			return nil, ok
		}
	}
	return pu, nil
}
