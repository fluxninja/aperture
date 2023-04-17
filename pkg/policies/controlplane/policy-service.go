package controlplane

import (
	"context"
	"fmt"
	"path"
	"strings"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	etcdwriter "github.com/fluxninja/aperture/pkg/etcd/writer"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/paths"
)

// PolicyService is the implementation of policylangv1.PolicyService interface.
type PolicyService struct {
	policylangv1.UnimplementedPolicyServiceServer
	policyFactory               *PolicyFactory
	policyTrackers              notifiers.Trackers
	policyDynamicConfigTrackers notifiers.Trackers
	etcdWriter                  *etcdwriter.Writer
}

// RegisterPolicyService registers a service for policy.
func RegisterPolicyService(
	policyTrackers, policyDynamicConfigTrackers notifiers.Trackers,
	server *grpc.Server,
	policyFactory *PolicyFactory,
	etcdClient *etcdclient.Client,
	lifecycle fx.Lifecycle,
) *PolicyService {
	svc := &PolicyService{
		policyFactory:               policyFactory,
		policyTrackers:              policyTrackers,
		policyDynamicConfigTrackers: policyDynamicConfigTrackers,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			svc.etcdWriter = etcdwriter.NewWriter(etcdClient, false)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return svc.etcdWriter.Close()
		},
	})

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
func (s *PolicyService) PostPolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest) (*policylangv1.PostResponse, error) {
	return s.updatePolicies(ctx, policies, true)
}

// PatchPolicies patches policies to the system.
func (s *PolicyService) PatchPolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest) (*policylangv1.PostResponse, error) {
	return s.updatePolicies(ctx, policies, false)
}

// PostDynamicConfigs updates dynamic configs to the system.
func (s *PolicyService) PostDynamicConfigs(ctx context.Context, dynamicConfigs *policylangv1.PostDynamicConfigsRequest) (*policylangv1.PostResponse, error) {
	var errs []string
	for idx, request := range dynamicConfigs.DynamicConfigs {
		if request.Name == "" {
			errs = append(errs, fmt.Sprintf("policy name is empty at index '%d'", idx))
			continue
		}

		_, err := s.GetPolicy(ctx, &policylangv1.GetPolicyRequest{Name: request.Name})
		if err != nil {
			errs = append(errs, fmt.Sprintf("policy '%s' not found", request.Name))
			continue
		}

		if request.DynamicConfig == nil {
			errs = append(errs, fmt.Sprintf("dynamic config is missing for policy name '%s'", request.Name))
			continue
		}

		jsonDynamicConfig, err := request.DynamicConfig.MarshalJSON()
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to marshal dynamic config '%s': '%s'", request.Name, err))
			continue
		}
		s.etcdWriter.Write(path.Join(paths.PoliciesAPIDynamicConfigPath, request.Name), jsonDynamicConfig)
	}

	response := &policylangv1.PostResponse{}
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
	s.etcdWriter.Delete(path.Join(paths.PoliciesAPIConfigPath, policy.Name))
	s.etcdWriter.Delete(path.Join(paths.PoliciesAPIDynamicConfigPath, policy.Name))
	return &emptypb.Empty{}, nil
}

// updatePolicies manages Post/Patch operations on policies.
func (s *PolicyService) updatePolicies(ctx context.Context, policies *policylangv1.PostPoliciesRequest, isPatch bool) (*policylangv1.PostResponse, error) {
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
		s.etcdWriter.Write(path.Join(paths.PoliciesAPIConfigPath, request.Name), jsonPolicy)
	}

	response := &policylangv1.PostResponse{}
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
