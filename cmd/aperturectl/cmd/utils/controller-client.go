package utils

import (
	"context"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	v1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	v11 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// IntrospectionClient is a subset of cmdv1.ControllerClient that covers APIs
// that need controller to grab information from agents via reverse rpc.
//
// These are currently not supported for the cloud controller.
//
// FIXME: Perhaps it'd be better to split the service on proto level (keep backcompat in mind).
type IntrospectionClient interface {
	ListAgents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*cmdv1.ListAgentsResponse, error)
	// Seems to be unimplemented on the controller at all?!
	ListServices(ctx context.Context, in *cmdv1.ListServicesRequest, opts ...grpc.CallOption) (*cmdv1.ListServicesControllerResponse, error)
	ListFlowControlPoints(ctx context.Context, in *cmdv1.ListFlowControlPointsRequest, opts ...grpc.CallOption) (*cmdv1.ListFlowControlPointsControllerResponse, error)
	ListAutoScaleControlPoints(ctx context.Context, in *cmdv1.ListAutoScaleControlPointsRequest, opts ...grpc.CallOption) (*cmdv1.ListAutoScaleControlPointsControllerResponse, error)
	ListDiscoveryEntities(ctx context.Context, in *cmdv1.ListDiscoveryEntitiesRequest, opts ...grpc.CallOption) (*cmdv1.ListDiscoveryEntitiesControllerResponse, error)
	ListDiscoveryEntity(ctx context.Context, in *cmdv1.ListDiscoveryEntityRequest, opts ...grpc.CallOption) (*cmdv1.ListDiscoveryEntityAgentResponse, error)
	PreviewFlowLabels(ctx context.Context, in *cmdv1.PreviewFlowLabelsRequest, opts ...grpc.CallOption) (*cmdv1.PreviewFlowLabelsControllerResponse, error)
	PreviewHTTPRequests(ctx context.Context, in *cmdv1.PreviewHTTPRequestsRequest, opts ...grpc.CallOption) (*cmdv1.PreviewHTTPRequestsControllerResponse, error)
}

// PolicyClient is a subset of cmdv1.ControllerClient that covers APIs related to policies.
//
// FIXME: Perhaps it'd be better to split the service on proto level (keep backcompat in mind).
type PolicyClient interface {
	ListPolicies(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.GetPoliciesResponse, error)
	UpsertPolicy(ctx context.Context, in *v1.UpsertPolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PostDynamicConfig(ctx context.Context, in *v1.PostDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeletePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetDecisions(ctx context.Context, in *v1.GetDecisionsRequest, opts ...grpc.CallOption) (*v1.GetDecisionsResponse, error)
}

// StatusClient is a subset of cmdv1.ControllerClient that covers APIs related to status.
//
// FIXME: Perhaps it'd be better to split the service on proto level (keep backcompat in mind).
type StatusClient interface {
	GetStatus(ctx context.Context, in *v11.GroupStatusRequest, opts ...grpc.CallOption) (*v11.GroupStatus, error)
}

// CloudPolicyClient is a subset of cloudv1.CloudControllerClient that covers APIs related to policies.
type CloudPolicyClient interface {
	UpsertPolicy(ctx context.Context, in *v1.UpsertPolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeletePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
