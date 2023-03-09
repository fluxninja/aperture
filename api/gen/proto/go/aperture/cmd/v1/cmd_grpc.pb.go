// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cmdv1

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

// ControllerClient is the client API for Controller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControllerClient interface {
	ListAgents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListAgentsResponse, error)
	ListServices(ctx context.Context, in *ListServicesRequest, opts ...grpc.CallOption) (*ListServicesControllerResponse, error)
	ListFlowControlPoints(ctx context.Context, in *ListFlowControlPointsRequest, opts ...grpc.CallOption) (*ListFlowControlPointsControllerResponse, error)
	ListAutoScaleControlPoints(ctx context.Context, in *ListAutoScaleControlPointsRequest, opts ...grpc.CallOption) (*ListAutoScaleControlPointsControllerResponse, error)
	ListDiscoveryEntities(ctx context.Context, in *ListDiscoveryEntitiesRequest, opts ...grpc.CallOption) (*ListDiscoveryEntitiesControllerResponse, error)
	// duplicating a bit preview.v1.FlowPreviewService to keep controller APIs in one place.
	PreviewFlowLabels(ctx context.Context, in *PreviewFlowLabelsRequest, opts ...grpc.CallOption) (*PreviewFlowLabelsControllerResponse, error)
	PreviewHTTPRequests(ctx context.Context, in *PreviewHTTPRequestsRequest, opts ...grpc.CallOption) (*PreviewHTTPRequestsControllerResponse, error)
}

type controllerClient struct {
	cc grpc.ClientConnInterface
}

func NewControllerClient(cc grpc.ClientConnInterface) ControllerClient {
	return &controllerClient{cc}
}

func (c *controllerClient) ListAgents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListAgentsResponse, error) {
	out := new(ListAgentsResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/ListAgents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerClient) ListServices(ctx context.Context, in *ListServicesRequest, opts ...grpc.CallOption) (*ListServicesControllerResponse, error) {
	out := new(ListServicesControllerResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/ListServices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerClient) ListFlowControlPoints(ctx context.Context, in *ListFlowControlPointsRequest, opts ...grpc.CallOption) (*ListFlowControlPointsControllerResponse, error) {
	out := new(ListFlowControlPointsControllerResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/ListFlowControlPoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerClient) ListAutoScaleControlPoints(ctx context.Context, in *ListAutoScaleControlPointsRequest, opts ...grpc.CallOption) (*ListAutoScaleControlPointsControllerResponse, error) {
	out := new(ListAutoScaleControlPointsControllerResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/ListAutoScaleControlPoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerClient) ListDiscoveryEntities(ctx context.Context, in *ListDiscoveryEntitiesRequest, opts ...grpc.CallOption) (*ListDiscoveryEntitiesControllerResponse, error) {
	out := new(ListDiscoveryEntitiesControllerResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/ListDiscoveryEntities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerClient) PreviewFlowLabels(ctx context.Context, in *PreviewFlowLabelsRequest, opts ...grpc.CallOption) (*PreviewFlowLabelsControllerResponse, error) {
	out := new(PreviewFlowLabelsControllerResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/PreviewFlowLabels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerClient) PreviewHTTPRequests(ctx context.Context, in *PreviewHTTPRequestsRequest, opts ...grpc.CallOption) (*PreviewHTTPRequestsControllerResponse, error) {
	out := new(PreviewHTTPRequestsControllerResponse)
	err := c.cc.Invoke(ctx, "/aperture.cmd.v1.Controller/PreviewHTTPRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControllerServer is the server API for Controller service.
// All implementations should embed UnimplementedControllerServer
// for forward compatibility
type ControllerServer interface {
	ListAgents(context.Context, *emptypb.Empty) (*ListAgentsResponse, error)
	ListServices(context.Context, *ListServicesRequest) (*ListServicesControllerResponse, error)
	ListFlowControlPoints(context.Context, *ListFlowControlPointsRequest) (*ListFlowControlPointsControllerResponse, error)
	ListAutoScaleControlPoints(context.Context, *ListAutoScaleControlPointsRequest) (*ListAutoScaleControlPointsControllerResponse, error)
	ListDiscoveryEntities(context.Context, *ListDiscoveryEntitiesRequest) (*ListDiscoveryEntitiesControllerResponse, error)
	// duplicating a bit preview.v1.FlowPreviewService to keep controller APIs in one place.
	PreviewFlowLabels(context.Context, *PreviewFlowLabelsRequest) (*PreviewFlowLabelsControllerResponse, error)
	PreviewHTTPRequests(context.Context, *PreviewHTTPRequestsRequest) (*PreviewHTTPRequestsControllerResponse, error)
}

// UnimplementedControllerServer should be embedded to have forward compatible implementations.
type UnimplementedControllerServer struct {
}

func (UnimplementedControllerServer) ListAgents(context.Context, *emptypb.Empty) (*ListAgentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAgents not implemented")
}
func (UnimplementedControllerServer) ListServices(context.Context, *ListServicesRequest) (*ListServicesControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListServices not implemented")
}
func (UnimplementedControllerServer) ListFlowControlPoints(context.Context, *ListFlowControlPointsRequest) (*ListFlowControlPointsControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFlowControlPoints not implemented")
}
func (UnimplementedControllerServer) ListAutoScaleControlPoints(context.Context, *ListAutoScaleControlPointsRequest) (*ListAutoScaleControlPointsControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAutoScaleControlPoints not implemented")
}
func (UnimplementedControllerServer) ListDiscoveryEntities(context.Context, *ListDiscoveryEntitiesRequest) (*ListDiscoveryEntitiesControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDiscoveryEntities not implemented")
}
func (UnimplementedControllerServer) PreviewFlowLabels(context.Context, *PreviewFlowLabelsRequest) (*PreviewFlowLabelsControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreviewFlowLabels not implemented")
}
func (UnimplementedControllerServer) PreviewHTTPRequests(context.Context, *PreviewHTTPRequestsRequest) (*PreviewHTTPRequestsControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreviewHTTPRequests not implemented")
}

// UnsafeControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControllerServer will
// result in compilation errors.
type UnsafeControllerServer interface {
	mustEmbedUnimplementedControllerServer()
}

func RegisterControllerServer(s grpc.ServiceRegistrar, srv ControllerServer) {
	s.RegisterService(&Controller_ServiceDesc, srv)
}

func _Controller_ListAgents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).ListAgents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/ListAgents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).ListAgents(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Controller_ListServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).ListServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/ListServices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).ListServices(ctx, req.(*ListServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Controller_ListFlowControlPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFlowControlPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).ListFlowControlPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/ListFlowControlPoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).ListFlowControlPoints(ctx, req.(*ListFlowControlPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Controller_ListAutoScaleControlPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAutoScaleControlPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).ListAutoScaleControlPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/ListAutoScaleControlPoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).ListAutoScaleControlPoints(ctx, req.(*ListAutoScaleControlPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Controller_ListDiscoveryEntities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDiscoveryEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).ListDiscoveryEntities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/ListDiscoveryEntities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).ListDiscoveryEntities(ctx, req.(*ListDiscoveryEntitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Controller_PreviewFlowLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PreviewFlowLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).PreviewFlowLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/PreviewFlowLabels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).PreviewFlowLabels(ctx, req.(*PreviewFlowLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Controller_PreviewHTTPRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PreviewHTTPRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).PreviewHTTPRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aperture.cmd.v1.Controller/PreviewHTTPRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).PreviewHTTPRequests(ctx, req.(*PreviewHTTPRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Controller_ServiceDesc is the grpc.ServiceDesc for Controller service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Controller_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aperture.cmd.v1.Controller",
	HandlerType: (*ControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAgents",
			Handler:    _Controller_ListAgents_Handler,
		},
		{
			MethodName: "ListServices",
			Handler:    _Controller_ListServices_Handler,
		},
		{
			MethodName: "ListFlowControlPoints",
			Handler:    _Controller_ListFlowControlPoints_Handler,
		},
		{
			MethodName: "ListAutoScaleControlPoints",
			Handler:    _Controller_ListAutoScaleControlPoints_Handler,
		},
		{
			MethodName: "ListDiscoveryEntities",
			Handler:    _Controller_ListDiscoveryEntities_Handler,
		},
		{
			MethodName: "PreviewFlowLabels",
			Handler:    _Controller_PreviewFlowLabels_Handler,
		},
		{
			MethodName: "PreviewHTTPRequests",
			Handler:    _Controller_PreviewHTTPRequests_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aperture/cmd/v1/cmd.proto",
}
