package utils

import (
	"context"

	v1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CloudPolicyClient is a subset of cloudv1.CloudControllerClient that covers APIs related to policies.
type CloudPolicyClient interface {
	UpsertPolicy(ctx context.Context, in *v1.UpsertPolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeletePolicy(ctx context.Context, in *v1.DeletePolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
