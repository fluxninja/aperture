package utils

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
)

// IntrospectionClient is a subset of cmdv1.ControllerClient that covers APIs
// that need controller to grab information from agents via reverse rpc.
//
// These are currently not supported for the cloud controller.
//
// FIXME: Perhaps it'd be better to split the service on proto level (keep backcompat in mind).
type IntrospectionClient interface {
	ListAgents(ctx context.Context, in *cmdv1.ListAgentsRequest, opts ...grpc.CallOption) (*cmdv1.ListAgentsResponse, error)
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
type PolicyClient interface {
	UpsertPolicy(ctx context.Context, in *policylangv1.UpsertPolicyRequest, opts ...grpc.CallOption) (*policylangv1.UpsertPolicyResponse, error)
	DeletePolicy(ctx context.Context, in *policylangv1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// SelfHostedPolicyClient is a subset of cmdv1.ControllerClient that covers APIs related to policies.
//
// FIXME: Perhaps it'd be better to split the service on proto level (keep backcompat in mind).
type SelfHostedPolicyClient interface {
	ListPolicies(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*policylangv1.GetPoliciesResponse, error)
	UpsertPolicy(ctx context.Context, in *policylangv1.UpsertPolicyRequest, opts ...grpc.CallOption) (*policylangv1.UpsertPolicyResponse, error)
	PostDynamicConfig(ctx context.Context, in *policylangv1.PostDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetDynamicConfig(ctx context.Context, in *policylangv1.GetDynamicConfigRequest, opts ...grpc.CallOption) (*policylangv1.GetDynamicConfigResponse, error)
	DeleteDynamicConfig(ctx context.Context, in *policylangv1.DeleteDynamicConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeletePolicy(ctx context.Context, in *policylangv1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetDecisions(ctx context.Context, in *policylangv1.GetDecisionsRequest, opts ...grpc.CallOption) (*policylangv1.GetDecisionsResponse, error)
	GetPolicy(ctx context.Context, in *policylangv1.GetPolicyRequest, opts ...grpc.CallOption) (*policylangv1.GetPolicyResponse, error)
}

// StatusClient is a subset of cmdv1.ControllerClient that covers APIs related to status.
//
// FIXME: Perhaps it'd be better to split the service on proto level (keep backcompat in mind).
type StatusClient interface {
	GetStatus(ctx context.Context, in *statusv1.GroupStatusRequest, opts ...grpc.CallOption) (*statusv1.GroupStatus, error)
}

// CloudPolicyClient is a subset of cloudv1.CloudControllerClient that covers APIs related to policies.
type CloudPolicyClient interface {
	UpsertPolicy(ctx context.Context, in *policylangv1.UpsertPolicyRequest, opts ...grpc.CallOption) (*policylangv1.UpsertPolicyResponse, error)
	DeletePolicy(ctx context.Context, in *policylangv1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ArchivePolicy(ctx context.Context, in *policylangv1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

var _ CloudPolicyClient = cloudv1.NewPolicyServiceClient(nil)

// CloudBlueprintsClient is a subset of cloudv1.CloudControllerClient that covers APIs related to cloud blueprints.
type CloudBlueprintsClient interface {
	List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*cloudv1.ListResponse, error)
	Get(ctx context.Context, in *cloudv1.GetRequest, opts ...grpc.CallOption) (*cloudv1.GetResponse, error)
	Apply(ctx context.Context, in *cloudv1.ApplyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *cloudv1.DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Archive(ctx context.Context, in *cloudv1.DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

var _ CloudBlueprintsClient = cloudv1.NewBlueprintsServiceClient(nil)
