package controlplane

import (
	"context"
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PolicyService is the implementation of policylangv1.PolicyService interface.
type PolicyService struct {
	policylangv1.UnimplementedPolicyServiceServer
	policyFactory        *PolicyFactory
	etcdTracker          notifiers.Trackers
	dynamicConfigTracker notifiers.Trackers
}

// RegisterPolicyService registers a service for policy.
func RegisterPolicyService(etcdTracker, dynamicConfigYTracker notifiers.Trackers, server *grpc.Server, policyFactory *PolicyFactory) {
	svc := &PolicyService{
		policyFactory:        policyFactory,
		etcdTracker:          etcdTracker,
		dynamicConfigTracker: dynamicConfigYTracker,
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

// PostPolicy posts policies to the system.
func (s *PolicyService) PostPolicy(ctx context.Context, policies *policylangv1.PostPolicyRequest) (*policylangv1.PostPoliciesResponse, error) {
	var errs []string
	for _, request := range policies.Policies {
		_, err := CompilePolicy(request.Policy, s.policyFactory.registry)
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to compile policy %s: %s", request.Name, err))
			continue
		}

		jsonPolicy, err := request.Policy.MarshalJSON()
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to marshal policy %s: %s", request.Name, err))
			continue
		}

		jsonDynamicConfig, err := request.DynamicConfig.MarshalJSON()
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to marshal dynamic config %s: %s", request.Name, err))
			continue
		}
		s.etcdTracker.WriteEvent(notifiers.Key(request.Name), jsonPolicy)
		s.dynamicConfigTracker.WriteEvent(notifiers.Key(request.Name), jsonDynamicConfig)
	}

	response := &policylangv1.PostPoliciesResponse{}
	if len(errs) > 0 {
		response.Errors = errs
		response.Message = "failed to post policies"
	} else {
		response.Message = "successfully posted policies"
	}
	return response, nil
}

// DeletePolicy deletes a policy from the system.
func (s *PolicyService) DeletePolicy(ctx context.Context, policy *policylangv1.DeletePolicyRequest) (*emptypb.Empty, error) {
	s.etcdTracker.RemoveEvent(notifiers.Key(policy.Name))
	return &emptypb.Empty{}, nil
}
