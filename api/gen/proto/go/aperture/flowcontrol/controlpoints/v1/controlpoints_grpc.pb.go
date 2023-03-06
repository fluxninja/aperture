// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package controlpointsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FlowControlPointsServiceClient is the client API for FlowControlPointsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FlowControlPointsServiceClient interface {
	GetControlPoints(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FlowControlPoints, error)
}

type flowControlPointsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFlowControlPointsServiceClient(cc grpc.ClientConnInterface) FlowControlPointsServiceClient {
	return &flowControlPointsServiceClient{cc}
}

func (c *flowControlPointsServiceClient) GetControlPoints(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FlowControlPoints, error) {
	out := new(FlowControlPoints)
	err := c.cc.Invoke(ctx, "/aperture.flowcontrol.controlpoints.v1.FlowControlPointsService/GetControlPoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FlowControlPointsServiceServer is the server API for FlowControlPointsService service.
// All implementations should embed UnimplementedFlowControlPointsServiceServer
// for forward compatibility
type FlowControlPointsServiceServer interface {
	GetControlPoints(context.Context, *emptypb.Empty) (*FlowControlPoints, error)
}

// UnimplementedFlowControlPointsServiceServer should be embedded to have forward compatible implementations.
type UnimplementedFlowControlPointsServiceServer struct {
}

func (UnimplementedFlowControlPointsServiceServer) GetControlPoints(context.Context, *emptypb.Empty) (*FlowControlPoints, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetControlPoints not implemented")
}

// UnsafeFlowControlPointsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FlowControlPointsServiceServer will
// result in compilation errors.
type UnsafeFlowControlPointsServiceServer interface {
	mustEmbedUnimplementedFlowControlPointsServiceServer()
}

func RegisterFlowControlPointsServiceServer(s grpc.ServiceRegistrar, srv FlowControlPointsServiceServer) {
	s.RegisterService(&FlowControlPointsService_ServiceDesc, srv)
}

func _FlowControlPointsService_GetControlPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlowControlPointsServiceServer).GetControlPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.flowcontrol.controlpoints.v1.FlowControlPointsService/GetControlPoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlowControlPointsServiceServer).GetControlPoints(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// FlowControlPointsService_ServiceDesc is the grpc.ServiceDesc for FlowControlPointsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FlowControlPointsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aperture.flowcontrol.controlpoints.v1.FlowControlPointsService",
	HandlerType: (*FlowControlPointsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetControlPoints",
			Handler:    _FlowControlPointsService_GetControlPoints_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aperture/flowcontrol/controlpoints/v1/controlpoints.proto",
}
