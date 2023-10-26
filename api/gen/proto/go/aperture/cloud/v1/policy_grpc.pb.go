// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: aperture/cloud/v1/policy.proto

// Messages for aperturectl → cloud controller communication.

package cloudv1

import (
	context "context"
	v1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PolicyService_UpsertPolicy_FullMethodName  = "/aperture.cloud.v1.PolicyService/UpsertPolicy"
	PolicyService_DeletePolicy_FullMethodName  = "/aperture.cloud.v1.PolicyService/DeletePolicy"
	PolicyService_ArchivePolicy_FullMethodName = "/aperture.cloud.v1.PolicyService/ArchivePolicy"
)

// PolicyServiceClient is the client API for PolicyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PolicyServiceClient interface {
	UpsertPolicy(ctx context.Context, in *v1.UpsertPolicyRequest, opts ...grpc.CallOption) (*v1.UpsertPolicyResponse, error)
	DeletePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ArchivePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type policyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPolicyServiceClient(cc grpc.ClientConnInterface) PolicyServiceClient {
	return &policyServiceClient{cc}
}

func (c *policyServiceClient) UpsertPolicy(ctx context.Context, in *v1.UpsertPolicyRequest, opts ...grpc.CallOption) (*v1.UpsertPolicyResponse, error) {
	out := new(v1.UpsertPolicyResponse)
	err := c.cc.Invoke(ctx, PolicyService_UpsertPolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) DeletePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PolicyService_DeletePolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) ArchivePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PolicyService_ArchivePolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PolicyServiceServer is the server API for PolicyService service.
// All implementations should embed UnimplementedPolicyServiceServer
// for forward compatibility
type PolicyServiceServer interface {
	UpsertPolicy(context.Context, *v1.UpsertPolicyRequest) (*v1.UpsertPolicyResponse, error)
	DeletePolicy(context.Context, *v1.DeletePolicyRequest) (*emptypb.Empty, error)
	ArchivePolicy(context.Context, *v1.DeletePolicyRequest) (*emptypb.Empty, error)
}

// UnimplementedPolicyServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPolicyServiceServer struct {
}

func (UnimplementedPolicyServiceServer) UpsertPolicy(context.Context, *v1.UpsertPolicyRequest) (*v1.UpsertPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertPolicy not implemented")
}
func (UnimplementedPolicyServiceServer) DeletePolicy(context.Context, *v1.DeletePolicyRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePolicy not implemented")
}
func (UnimplementedPolicyServiceServer) ArchivePolicy(context.Context, *v1.DeletePolicyRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArchivePolicy not implemented")
}

// UnsafePolicyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PolicyServiceServer will
// result in compilation errors.
type UnsafePolicyServiceServer interface {
	mustEmbedUnimplementedPolicyServiceServer()
}

func RegisterPolicyServiceServer(s grpc.ServiceRegistrar, srv PolicyServiceServer) {
	s.RegisterService(&PolicyService_ServiceDesc, srv)
}

func _PolicyService_UpsertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.UpsertPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).UpsertPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_UpsertPolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).UpsertPolicy(ctx, req.(*v1.UpsertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_DeletePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DeletePolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).DeletePolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_DeletePolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).DeletePolicy(ctx, req.(*v1.DeletePolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_ArchivePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DeletePolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).ArchivePolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_ArchivePolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).ArchivePolicy(ctx, req.(*v1.DeletePolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PolicyService_ServiceDesc is the grpc.ServiceDesc for PolicyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PolicyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aperture.cloud.v1.PolicyService",
	HandlerType: (*PolicyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertPolicy",
			Handler:    _PolicyService_UpsertPolicy_Handler,
		},
		{
			MethodName: "DeletePolicy",
			Handler:    _PolicyService_DeletePolicy_Handler,
		},
		{
			MethodName: "ArchivePolicy",
			Handler:    _PolicyService_ArchivePolicy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aperture/cloud/v1/policy.proto",
}
