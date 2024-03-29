package grpc

import (
	"fmt"

	gogoproto "github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"google.golang.org/protobuf/proto"
)

// Codec that supports both VTProto-enabled messages (that have helpers
// generated by vtprotobuf plugin), but also regular messages (including
// support for gogoproto messages from etcd api and collector/pdata).
//
// Adapted from https://github.com/vitessio/vitess/blob/main/go/vt/servenv/grpc_codec.go
// See also https://github.com/planetscale/vtprotobuf#vtprotobuf-with-grpc
type vtprotoCodec struct{}

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

// Marshal implements grpc/encoding.Codec.
func (vtprotoCodec) Marshal(v any) ([]byte, error) {
	switch v := v.(type) {
	// FIXME: This is commented out because marshaling CheckHTTPResponse is
	// interestingly slower with VTMarshal (with most time consumed by proto.Size).
	// This is caused by https://github.com/planetscale/vtprotobuf/issues/54.
	// case vtprotoMessage:
	// 	return v.MarshalVT()
	case proto.Message:
		return proto.Marshal(v)
	// Required for etcd client & collector/pdata.
	case gogoproto.Message:
		return gogoproto.Marshal(v)
	default:
		return nil, fmt.Errorf("failed to marshal, message is %T, want vtprotoMessage, proto.Message or gogoproto.Message", v)
	}
}

// Unmarshal implements grpc/encoding.Codec.
func (vtprotoCodec) Unmarshal(data []byte, v any) error {
	switch v := v.(type) {
	case vtprotoMessage:
		return v.UnmarshalVT(data)
	case proto.Message:
		return proto.Unmarshal(data, v)
	case gogoproto.Message:
		return gogoproto.Unmarshal(data, v)
	default:
		return fmt.Errorf("failed to unmarshal, message is %T, want vtprotoMessage, proto.Message or gogoproto.Message", v)
	}
}

// Name implements grpc/encoding.Codec.
func (vtprotoCodec) Name() string { return "proto" }

func init() {
	// Use the optimized vtproto codecs by default in the grpc server.
	encoding.RegisterCodec(vtprotoCodec{})
}
