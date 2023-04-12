package controlplane

import (
	"context"
	"fmt"
	"strings"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func RegisterPolicyService(etcdTracker, dynamicConfigYTracker notifiers.Trackers, server *grpc.Server, policyFactory *PolicyFactory) *PolicyService {
	svc := &PolicyService{
		policyFactory:        policyFactory,
		etcdTracker:          etcdTracker,
		dynamicConfigTracker: dynamicConfigYTracker,
	}
	policylangv1.RegisterPolicyServiceServer(server, svc)
	return svc
}

// GetPolicies returns all the policies running in the system.
func (s *PolicyService) GetPolicies(ctx context.Context, _ *emptypb.Empty) (*policylangv1.GetPoliciesResponse, error) {
	policies := s.policyFactory.GetPolicies()
	return &policylangv1.GetPoliciesResponse{
		Policies: policies,
	}, nil
}

// GetPolicy returns the policy running in the system which matches the given name.
func (s *PolicyService) GetPolicy(ctx context.Context, request *policylangv1.GetPolicyRequest) (*policylangv1.GetPolicyResponse, error) {
	policy := s.policyFactory.GetPolicy(request.Name)
	if policy == nil {
		return nil, status.Error(codes.NotFound, "policy not found")
	}
	return &policylangv1.GetPolicyResponse{
		Policy: policy,
	}, nil
}

// PostPolicies posts policies to the system.
func (s *PolicyService) PostPolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest) (*policylangv1.PostPoliciesResponse, error) {
	return s.managePolicies(ctx, policies, true)
}

// PatchPolicies patches policies to the system.
func (s *PolicyService) PatchPolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest) (*policylangv1.PostPoliciesResponse, error) {
	return s.managePolicies(ctx, policies, false)
}

// DeletePolicy deletes a policy from the system.
func (s *PolicyService) DeletePolicy(ctx context.Context, policy *policylangv1.DeletePolicyRequest) (*emptypb.Empty, error) {
	s.etcdTracker.RemoveEvent(notifiers.Key(policy.Name))
	s.dynamicConfigTracker.RemoveEvent(notifiers.Key(policy.Name))
	return &emptypb.Empty{}, nil
}

// managePolicies manages Post/Patch operations on policies.
func (s *PolicyService) managePolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest, checkRequired bool) (*policylangv1.PostPoliciesResponse, error) {
	var errs []string
	for idx, request := range policies.Policies {
		if request.Name == "" {
			errs = append(errs, fmt.Sprintf("policy name is empty at index '%d'", idx))
			continue
		}

		if checkRequired {
			if request.Policy == nil {
				errs = append(errs, fmt.Sprintf("policy is missing for policy name '%s'", request.Name))
				continue
			}
		} else {
			_, err := s.GetPolicy(ctx, &policylangv1.GetPolicyRequest{Name: request.Name})
			if err != nil {
				errs = append(errs, fmt.Sprintf("policy '%s' not found", request.Name))
				continue
			}
		}

		jsonPolicy, err := s.getPolicyBytes(request.Name, request.Policy)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		s.etcdTracker.WriteEvent(notifiers.Key(request.Name), jsonPolicy)

		if request.DynamicConfig == nil {
			continue
		}

		jsonDynamicConfig, err := request.DynamicConfig.MarshalJSON()
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to marshal dynamic config '%s': '%s'", request.Name, err))
			continue
		}
		s.dynamicConfigTracker.WriteEvent(notifiers.Key(request.Name), jsonDynamicConfig)
	}

	response := &policylangv1.PostPoliciesResponse{}
	var err error
	if len(errs) > 0 {
		response.Errors = errs
		response.Message = "failed to upload policies"
		err = status.Error(codes.InvalidArgument, strings.Join(errs, ","))
	} else {
		response.Message = "successfully uploaded policies"
	}
	return response, err
}

// getPolicyBytes returns the policy bytes after checking validity of the policy.
func (s *PolicyService) getPolicyBytes(name string, policy *policylangv1.Policy) ([]byte, error) {
	jsonPolicy, err := policy.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy '%s': '%s'", name, err)
	}

	_, _, err = ValidateAndCompile(context.Background(), name, jsonPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to compile policy '%s': '%s'", name, err)
	}

	return jsonPolicy, nil
}
