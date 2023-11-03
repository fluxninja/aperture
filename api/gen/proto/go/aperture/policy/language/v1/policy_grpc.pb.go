// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: aperture/policy/language/v1/policy.proto

package languagev1

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

const (
	PolicyService_GetPolicy_FullMethodName           = "/aperture.policy.language.v1.PolicyService/GetPolicy"
	PolicyService_GetPolicies_FullMethodName         = "/aperture.policy.language.v1.PolicyService/GetPolicies"
	PolicyService_UpsertPolicy_FullMethodName        = "/aperture.policy.language.v1.PolicyService/UpsertPolicy"
	PolicyService_PostDynamicConfig_FullMethodName   = "/aperture.policy.language.v1.PolicyService/PostDynamicConfig"
	PolicyService_GetDynamicConfig_FullMethodName    = "/aperture.policy.language.v1.PolicyService/GetDynamicConfig"
	PolicyService_DeleteDynamicConfig_FullMethodName = "/aperture.policy.language.v1.PolicyService/DeleteDynamicConfig"
	PolicyService_DeletePolicy_FullMethodName        = "/aperture.policy.language.v1.PolicyService/DeletePolicy"
	PolicyService_GetDecisions_FullMethodName        = "/aperture.policy.language.v1.PolicyService/GetDecisions"
)

// PolicyServiceClient is the client API for PolicyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PolicyServiceClient interface {
	// GetPolicy returns a policy with the specified name.
	GetPolicy(ctx context.Context, in *GetPolicyRequest, opts ...grpc.CallOption) (*GetPolicyResponse, error)
	// GetPolicies returns all policies.
	GetPolicies(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetPoliciesResponse, error)
	// UpsertPolicy creates/updates policy based on the provided request.
	UpsertPolicy(ctx context.Context, in *UpsertPolicyRequest, opts ...grpc.CallOption) (*UpsertPolicyResponse, error)
	// PostDynamicConfig creates/updates dynamic configuration based on the provided request.
	PostDynamicConfig(ctx context.Context, in *PostDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// GetDynamicConfig lists dynamic configuration for a policy.
	GetDynamicConfig(ctx context.Context, in *GetDynamicConfigRequest, opts ...grpc.CallOption) (*GetDynamicConfigResponse, error)
	// DeleteDynamicConfig deletes dynamic configuration for a policy.
	DeleteDynamicConfig(ctx context.Context, in *DeleteDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeletePolicy removes a policy with the specified name.
	DeletePolicy(ctx context.Context, in *DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetDecisions(ctx context.Context, in *GetDecisionsRequest, opts ...grpc.CallOption) (*GetDecisionsResponse, error)
}

type policyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPolicyServiceClient(cc grpc.ClientConnInterface) PolicyServiceClient {
	return &policyServiceClient{cc}
}

func (c *policyServiceClient) GetPolicy(ctx context.Context, in *GetPolicyRequest, opts ...grpc.CallOption) (*GetPolicyResponse, error) {
	out := new(GetPolicyResponse)
	err := c.cc.Invoke(ctx, PolicyService_GetPolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) GetPolicies(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetPoliciesResponse, error) {
	out := new(GetPoliciesResponse)
	err := c.cc.Invoke(ctx, PolicyService_GetPolicies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) UpsertPolicy(ctx context.Context, in *UpsertPolicyRequest, opts ...grpc.CallOption) (*UpsertPolicyResponse, error) {
	out := new(UpsertPolicyResponse)
	err := c.cc.Invoke(ctx, PolicyService_UpsertPolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) PostDynamicConfig(ctx context.Context, in *PostDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PolicyService_PostDynamicConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) GetDynamicConfig(ctx context.Context, in *GetDynamicConfigRequest, opts ...grpc.CallOption) (*GetDynamicConfigResponse, error) {
	out := new(GetDynamicConfigResponse)
	err := c.cc.Invoke(ctx, PolicyService_GetDynamicConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) DeleteDynamicConfig(ctx context.Context, in *DeleteDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PolicyService_DeleteDynamicConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) DeletePolicy(ctx context.Context, in *DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PolicyService_DeletePolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) GetDecisions(ctx context.Context, in *GetDecisionsRequest, opts ...grpc.CallOption) (*GetDecisionsResponse, error) {
	out := new(GetDecisionsResponse)
	err := c.cc.Invoke(ctx, PolicyService_GetDecisions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PolicyServiceServer is the server API for PolicyService service.
// All implementations should embed UnimplementedPolicyServiceServer
// for forward compatibility
type PolicyServiceServer interface {
	// GetPolicy returns a policy with the specified name.
	GetPolicy(context.Context, *GetPolicyRequest) (*GetPolicyResponse, error)
	// GetPolicies returns all policies.
	GetPolicies(context.Context, *emptypb.Empty) (*GetPoliciesResponse, error)
	// UpsertPolicy creates/updates policy based on the provided request.
	UpsertPolicy(context.Context, *UpsertPolicyRequest) (*UpsertPolicyResponse, error)
	// PostDynamicConfig creates/updates dynamic configuration based on the provided request.
	PostDynamicConfig(context.Context, *PostDynamicConfigRequest) (*emptypb.Empty, error)
	// GetDynamicConfig lists dynamic configuration for a policy.
	GetDynamicConfig(context.Context, *GetDynamicConfigRequest) (*GetDynamicConfigResponse, error)
	// DeleteDynamicConfig deletes dynamic configuration for a policy.
	DeleteDynamicConfig(context.Context, *DeleteDynamicConfigRequest) (*emptypb.Empty, error)
	// DeletePolicy removes a policy with the specified name.
	DeletePolicy(context.Context, *DeletePolicyRequest) (*emptypb.Empty, error)
	GetDecisions(context.Context, *GetDecisionsRequest) (*GetDecisionsResponse, error)
}

// UnimplementedPolicyServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPolicyServiceServer struct {
}

func (UnimplementedPolicyServiceServer) GetPolicy(context.Context, *GetPolicyRequest) (*GetPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPolicy not implemented")
}
func (UnimplementedPolicyServiceServer) GetPolicies(context.Context, *emptypb.Empty) (*GetPoliciesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPolicies not implemented")
}
func (UnimplementedPolicyServiceServer) UpsertPolicy(context.Context, *UpsertPolicyRequest) (*UpsertPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertPolicy not implemented")
}
func (UnimplementedPolicyServiceServer) PostDynamicConfig(context.Context, *PostDynamicConfigRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostDynamicConfig not implemented")
}
func (UnimplementedPolicyServiceServer) GetDynamicConfig(context.Context, *GetDynamicConfigRequest) (*GetDynamicConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDynamicConfig not implemented")
}
func (UnimplementedPolicyServiceServer) DeleteDynamicConfig(context.Context, *DeleteDynamicConfigRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDynamicConfig not implemented")
}
func (UnimplementedPolicyServiceServer) DeletePolicy(context.Context, *DeletePolicyRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePolicy not implemented")
}
func (UnimplementedPolicyServiceServer) GetDecisions(context.Context, *GetDecisionsRequest) (*GetDecisionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDecisions not implemented")
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

func _PolicyService_GetPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).GetPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_GetPolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).GetPolicy(ctx, req.(*GetPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_GetPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).GetPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_GetPolicies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).GetPolicies(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_UpsertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertPolicyRequest)
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
		return srv.(PolicyServiceServer).UpsertPolicy(ctx, req.(*UpsertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_PostDynamicConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostDynamicConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).PostDynamicConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_PostDynamicConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).PostDynamicConfig(ctx, req.(*PostDynamicConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_GetDynamicConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDynamicConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).GetDynamicConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_GetDynamicConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).GetDynamicConfig(ctx, req.(*GetDynamicConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_DeleteDynamicConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDynamicConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).DeleteDynamicConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_DeleteDynamicConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).DeleteDynamicConfig(ctx, req.(*DeleteDynamicConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_DeletePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePolicyRequest)
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
		return srv.(PolicyServiceServer).DeletePolicy(ctx, req.(*DeletePolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_GetDecisions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDecisionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).GetDecisions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_GetDecisions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).GetDecisions(ctx, req.(*GetDecisionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PolicyService_ServiceDesc is the grpc.ServiceDesc for PolicyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PolicyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aperture.policy.language.v1.PolicyService",
	HandlerType: (*PolicyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPolicy",
			Handler:    _PolicyService_GetPolicy_Handler,
		},
		{
			MethodName: "GetPolicies",
			Handler:    _PolicyService_GetPolicies_Handler,
		},
		{
			MethodName: "UpsertPolicy",
			Handler:    _PolicyService_UpsertPolicy_Handler,
		},
		{
			MethodName: "PostDynamicConfig",
			Handler:    _PolicyService_PostDynamicConfig_Handler,
		},
		{
			MethodName: "GetDynamicConfig",
			Handler:    _PolicyService_GetDynamicConfig_Handler,
		},
		{
			MethodName: "DeleteDynamicConfig",
			Handler:    _PolicyService_DeleteDynamicConfig_Handler,
		},
		{
			MethodName: "DeletePolicy",
			Handler:    _PolicyService_DeletePolicy_Handler,
		},
		{
			MethodName: "GetDecisions",
			Handler:    _PolicyService_GetDecisions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aperture/policy/language/v1/policy.proto",
}
