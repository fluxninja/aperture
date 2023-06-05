// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: aperture/distcache/v1/stats.proto

package distcachev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DistCacheService_GetStats_FullMethodName = "/aperture.distcache.v1.DistCacheService/GetStats"
)

// DistCacheServiceClient is the client API for DistCacheService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DistCacheServiceClient interface {
	GetStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*structpb.Struct, error)
}

type distCacheServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDistCacheServiceClient(cc grpc.ClientConnInterface) DistCacheServiceClient {
	return &distCacheServiceClient{cc}
}

func (c *distCacheServiceClient) GetStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*structpb.Struct, error) {
	out := new(structpb.Struct)
	err := c.cc.Invoke(ctx, DistCacheService_GetStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DistCacheServiceServer is the server API for DistCacheService service.
// All implementations should embed UnimplementedDistCacheServiceServer
// for forward compatibility
type DistCacheServiceServer interface {
	GetStats(context.Context, *emptypb.Empty) (*structpb.Struct, error)
}

// UnimplementedDistCacheServiceServer should be embedded to have forward compatible implementations.
type UnimplementedDistCacheServiceServer struct {
}

func (UnimplementedDistCacheServiceServer) GetStats(context.Context, *emptypb.Empty) (*structpb.Struct, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}

// UnsafeDistCacheServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DistCacheServiceServer will
// result in compilation errors.
type UnsafeDistCacheServiceServer interface {
	mustEmbedUnimplementedDistCacheServiceServer()
}

func RegisterDistCacheServiceServer(s grpc.ServiceRegistrar, srv DistCacheServiceServer) {
	s.RegisterService(&DistCacheService_ServiceDesc, srv)
}

func _DistCacheService_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistCacheServiceServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DistCacheService_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistCacheServiceServer).GetStats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// DistCacheService_ServiceDesc is the grpc.ServiceDesc for DistCacheService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DistCacheService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aperture.distcache.v1.DistCacheService",
	HandlerType: (*DistCacheServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStats",
			Handler:    _DistCacheService_GetStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aperture/distcache/v1/stats.proto",
}
