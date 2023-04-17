package controlplane

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// PolicyService is the implementation of policylangv1.PolicyService interface.
type PolicyService struct {
	policylangv1.UnimplementedPolicyServiceServer
	policyFactory               *PolicyFactory
	policyTrackers              notifiers.Trackers
	policyDynamicConfigTrackers notifiers.Trackers
}

// RegisterPolicyService registers a service for policy.
func RegisterPolicyService(
	policyTrackers, policyDynamicConfigTrackers notifiers.Trackers,
	server *grpc.Server,
	policyFactory *PolicyFactory,
) *PolicyService {
	svc := &PolicyService{
		policyFactory:               policyFactory,
		policyTrackers:              policyTrackers,
		policyDynamicConfigTrackers: policyDynamicConfigTrackers,
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
	return s.updatePolicies(ctx, policies, true)
}

// PatchPolicies patches policies to the system.
func (s *PolicyService) PatchPolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest) (*policylangv1.PostPoliciesResponse, error) {
	return s.updatePolicies(ctx, policies, false)
}

// PostDynamicConfigs updtaes dynamic configs to the system.
func (s *PolicyService) PostDynamicConfigs(ctx context.Context, dynamicConfigs *policylangv1.PostDynamicConfigsRequest) (*policylangv1.PostPoliciesResponse, error) {
	var errs []string
	for idx, request := range dynamicConfigs.DynamicConfigs {
		if request.PolicyName == "" {
			errs = append(errs, fmt.Sprintf("policy name is empty at index '%d'", idx))
			continue
		}

		_, err := s.GetPolicy(ctx, &policylangv1.GetPolicyRequest{Name: request.PolicyName})
		if err != nil {
			errs = append(errs, fmt.Sprintf("policy '%s' not found", request.PolicyName))
			continue
		}

		if request.DynamicConfig == nil {
			errs = append(errs, fmt.Sprintf("dynamic config is missing for policy name '%s'", request.PolicyName))
			continue
		}

		jsonDynamicConfig, err := request.DynamicConfig.MarshalJSON()
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to marshal dynamic config '%s': '%s'", request.PolicyName, err))
			continue
		}
		s.policyDynamicConfigTrackers.WriteEvent(notifiers.Key(request.PolicyName), jsonDynamicConfig)
	}

	response := &policylangv1.PostPoliciesResponse{}
	var err error
	if len(errs) > 0 {
		response.Errors = errs
		response.Message = "failed to upload dynamic configs"
		err = status.Error(codes.InvalidArgument, strings.Join(errs, ","))
	} else {
		response.Message = "successfully uploaded dynamic configs"
	}
	return response, err
}

// DeletePolicy deletes a policy from the system.
func (s *PolicyService) DeletePolicy(ctx context.Context, policy *policylangv1.DeletePolicyRequest) (*emptypb.Empty, error) {
	s.policyTrackers.RemoveEvent(notifiers.Key(policy.Name))
	s.policyDynamicConfigTrackers.RemoveEvent(notifiers.Key(policy.Name))
	return &emptypb.Empty{}, nil
}

// updatePolicies manages Post/Patch operations on policies.
func (s *PolicyService) updatePolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest, isPatch bool) (*policylangv1.PostPoliciesResponse, error) {
	var errs []string
	for idx, request := range policies.Policies {
		if request.Name == "" {
			errs = append(errs, fmt.Sprintf("policy name is not provided at index '%d'", idx))
			continue
		}

		if isPatch {
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
		s.policyTrackers.WriteEvent(notifiers.Key(request.Name), jsonPolicy)
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
