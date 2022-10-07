package controlplane

import (
	"context"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PolicyService is the implementation of policylangv1.PolicyService interface.
type PolicyService struct {
	policylangv1.UnimplementedPolicyServiceServer
	policyFactory *PolicyFactory
}

// RegisterPolicyService registers a service for policy.
func RegisterPolicyService(server *grpc.Server, policyFactory *PolicyFactory) {
	svc := &PolicyService{
		policyFactory: policyFactory,
	}
	policylangv1.RegisterPolicyServiceServer(server, svc)
}

// GetPolicies returns all the policies running in the system.
func (s *PolicyService) GetPolicies(ctx context.Context, _ *emptypb.Empty) (*policylangv1.GetPoliciesResponse, error) {
	policies := s.policyFactory.GetPolicies()
	return &policylangv1.GetPoliciesResponse{
		Policies: policies,
	}, nil
}
