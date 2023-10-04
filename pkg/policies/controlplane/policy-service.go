package controlplane

import (
	"context"
	"fmt"
	"path"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/paths"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// PolicyService is the implementation of policylangv1.PolicyService interface.
type PolicyService struct {
	policylangv1.UnimplementedPolicyServiceServer
	policyFactory *PolicyFactory
	etcdClient    *etcdclient.Client
}

// RegisterPolicyServiceIn bundles and annotates parameters.
type RegisterPolicyServiceIn struct {
	fx.In
	Server        *grpc.Server `name:"default"`
	PolicyFactory *PolicyFactory
	ETCDClient    *etcdclient.Client
	Lifecycle     fx.Lifecycle
}

// RegisterPolicyService registers a service for policy.
func RegisterPolicyService(in RegisterPolicyServiceIn) *PolicyService {
	svc := &PolicyService{
		policyFactory: in.PolicyFactory,
		etcdClient:    in.ETCDClient,
	}

	policylangv1.RegisterPolicyServiceServer(in.Server, svc)
	return svc
}

// GetPolicies returns all the policies running (or supposed to be running) in the system.
func (s *PolicyService) GetPolicies(ctx context.Context, _ *emptypb.Empty) (*policylangv1.GetPoliciesResponse, error) {
	localPolicies := s.policyFactory.GetPolicyWrappers()
	remotePolicies, err := s.etcdClient.Client.KV.Get(ctx, paths.PoliciesAPIConfigPath, clientv3.WithPrefix())
	if err != nil {
		log.Warn().Err(err).Msg("GetPolicies: failed to fetch policies from etcd")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch policies from etcd: %s", err))
	}

	resp := make(map[string]*policylangv1.GetPolicyResponse, len(remotePolicies.Kvs))
	for _, kv := range remotePolicies.Kvs {
		name, ok := strings.CutPrefix(string(kv.Key), fmt.Sprintf("%s/", paths.PoliciesAPIConfigPath))
		if !ok {
			log.Bug().Str("key", string(kv.Key)).Msg("failed to parse policy name from key")
			continue
		}
		resp[name] = getPolicyResponse(kv.Value, localPolicies[name])
	}

	for name, policy := range localPolicies {
		if _, ok := resp[name]; ok {
			continue
		}
		resp[name] = getLocalOnlyPolicyResponse(policy)
	}

	return &policylangv1.GetPoliciesResponse{
		Policies: &policylangv1.Policies{
			Policies: resp,
		},
	}, nil
}

// GetPolicy returns the policy which matches the given name.
//
// Returns error if policy cannot be found in *neither* etcd nor locally.
func (s *PolicyService) GetPolicy(ctx context.Context, request *policylangv1.GetPolicyRequest) (*policylangv1.GetPolicyResponse, error) {
	localPolicy := s.policyFactory.GetPolicyWrapper(request.Name)
	remotePolicies, err := s.etcdClient.Client.KV.Get(ctx, path.Join(paths.PoliciesAPIConfigPath, request.Name))
	if err != nil {
		log.Warn().Err(err).Str("policy", request.Name).Msg("GetPolicy: failed to fetch policy from etcd")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch policy from etcd: %s", err))

	}

	if len(remotePolicies.Kvs) == 0 {
		if localPolicy == nil {
			return nil, status.Error(codes.NotFound, "policy not found")
		}
		return getLocalOnlyPolicyResponse(localPolicy), nil
	}

	return getPolicyResponse(remotePolicies.Kvs[0].Value, localPolicy), nil
}

// getPolicyResponse combines info about the policy from etcd and local
// policyFactory to create a response for GetPolicy/GetPolicies
//
// localPolicy can be nil.
func getPolicyResponse(remoteBytes []byte, localPolicy *policysyncv1.PolicyWrapper) *policylangv1.GetPolicyResponse {
	remotePolicy := &policylangv1.Policy{}
	err := proto.Unmarshal(remoteBytes, remotePolicy)
	if err != nil {
		if localPolicy == nil {
			return &policylangv1.GetPolicyResponse{
				Status: policylangv1.GetPolicyResponse_INVALID,
				Reason: fmt.Sprintf("parse error: %s", err),
			}
		}
		return &policylangv1.GetPolicyResponse{
			Policy: localPolicy.Policy,
			Status: policylangv1.GetPolicyResponse_OUTDATED,
			Reason: fmt.Sprintf("parse error: %s", err),
		}
	}

	if localPolicy == nil {
		return &policylangv1.GetPolicyResponse{
			Policy: remotePolicy,
			Status: policylangv1.GetPolicyResponse_NOT_LOADED,
			Reason: "not loaded into controller",
		}
	}

	remoteHash, err := HashPolicy(remotePolicy)
	if err != nil {
		return &policylangv1.GetPolicyResponse{
			Policy: localPolicy.Policy,
			Status: policylangv1.GetPolicyResponse_OUTDATED,
			Reason: fmt.Sprintf("failed to compute remote policy hash: %s", err),
		}
	}

	localHash := localPolicy.GetCommonAttributes().GetPolicyHash()
	if remoteHash != localHash {
		return &policylangv1.GetPolicyResponse{
			Policy: localPolicy.Policy,
			Status: policylangv1.GetPolicyResponse_OUTDATED,
			Reason: fmt.Sprintf(
				"policy mismatch, remote hash: %s, local hash %s",
				remoteHash,
				localHash,
			),
		}
	}

	return &policylangv1.GetPolicyResponse{
		Policy: localPolicy.Policy,
		Status: policylangv1.GetPolicyResponse_VALID,
	}
}

func getLocalOnlyPolicyResponse(localPolicy *policysyncv1.PolicyWrapper) *policylangv1.GetPolicyResponse {
	if localPolicy.Source == policysyncv1.PolicyWrapper_ETCD {
		return &policylangv1.GetPolicyResponse{
			Policy: localPolicy.Policy,
			Status: policylangv1.GetPolicyResponse_STALE,
			Reason: "deleted, but still running",
		}
	}

	return &policylangv1.GetPolicyResponse{
		// It's expected for policies from K8S CRs to not be present in /config/api/policies
		Policy: localPolicy.Policy,
		Status: policylangv1.GetPolicyResponse_VALID,
	}
}

// UpsertPolicy creates/updates policy to the system.
func (s *PolicyService) UpsertPolicy(ctx context.Context, req *policylangv1.UpsertPolicyRequest) (*emptypb.Empty, error) {
	if wrapper := s.policyFactory.GetPolicyWrapper(req.PolicyName); wrapper.GetSource() == policysyncv1.PolicyWrapper_K8S {
		return nil, status.Error(codes.AlreadyExists, "policy already exists and is managed by K8S object")
	}

	etcdPolicy, err := s.etcdClient.Client.KV.Get(ctx, path.Join(paths.PoliciesAPIConfigPath, req.PolicyName))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	oldPolicy := &policylangv1.Policy{}
	if len(etcdPolicy.Kvs) > 0 {
		uerr := config.UnmarshalJSON(etcdPolicy.Kvs[0].Value, oldPolicy)
		if uerr != nil {
			// Deprecated: v3.0.0. Older way of string policy on etcd.
			// Remove this code in v3.0.0.
			uerr = proto.Unmarshal(etcdPolicy.Kvs[0].Value, oldPolicy)
			if uerr != nil {
				return nil, status.Error(codes.FailedPrecondition, "cannot patch, existing policy is invalid")
			}
		}
	}

	newPolicy := &policylangv1.Policy{}
	if req.PolicyString != "" {
		err = config.UnmarshalYAML([]byte(req.PolicyString), newPolicy)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal policy: %s", err)
		}
	} else if req.Policy != nil { // Deprecated: v2.20.0. Should stop accepting policy as proto message
		newPolicy = proto.Clone(req.Policy).(*policylangv1.Policy)
	} else {
		return nil, status.Error(codes.InvalidArgument, "policy is empty")
	}

	if len(req.GetUpdateMask().GetPaths()) > 0 {
		if len(etcdPolicy.Kvs) == 0 {
			return nil, status.Error(codes.NotFound, "cannot patch, no such policy")
		}

		if len(req.UpdateMask.Paths) == 1 && req.UpdateMask.Paths[0] == "all" {
			// We use update_mask: { paths: ["all"] } to convey that we want to
			// update the whole Policy.  This is non-standard usage of update
			// masks, so we cannot call ApplyFieldMask here.
			// FIXME Do we need field masks at all? Looks like they're only
			// used to pass a "force" flag from aperturectl.
		} else {
			utils.ApplyFieldMask(oldPolicy, newPolicy, req.UpdateMask)
			newPolicy = oldPolicy
		}
	} else if len(etcdPolicy.Kvs) > 0 {
		// We want to prevent accidentally overwriting valid policy, that's why
		// it's only possible when PATCH (update mask) is specified.
		_, _, err = ValidateAndCompileProto(ctx, "dummy-name", oldPolicy)
		if err == nil {
			// Note: Older aperturectl versions and Aperture Cloud rely on
			// exact string of error message!
			return nil, status.Error(codes.AlreadyExists, "policy already exists. Use UpsertPolicy with PATCH call to update it.")
		}
		if ctx.Err() != nil {
			return nil, err
		}
		// Otherwise, we know policy is invalid and thus we allow overwriting it.
	}

	_, _, err = ValidateAndCompileProto(ctx, req.PolicyName, newPolicy)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to compile policy: %s", err)
	}

	newPolicyString, err := newPolicy.MarshalJSON()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal policy: %s", err)
	}

	// FIXME compare original mod revision to make sure the policy we're patching hasn't changed meanwhile
	_, err = s.etcdClient.KV.Put(ctx, path.Join(paths.PoliciesAPIConfigPath, req.PolicyName), string(newPolicyString))
	if err != nil {
		return nil, err
	}

	if len(etcdPolicy.Kvs) > 0 {
		log.Info().Str("policy", req.PolicyName).Msg("Policy updated in etcd")
	} else {
		log.Info().Str("policy", req.PolicyName).Msg("Policy created in etcd")
	}

	return new(emptypb.Empty), nil
}

// PostDynamicConfig updates dynamic config to the system.
func (s *PolicyService) PostDynamicConfig(ctx context.Context, req *policylangv1.PostDynamicConfigRequest) (*emptypb.Empty, error) {
	etcdPolicy, err := s.etcdClient.Client.KV.Get(
		ctx,
		path.Join(paths.PoliciesAPIConfigPath, req.PolicyName),
		clientv3.WithKeysOnly(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(etcdPolicy.Kvs) == 0 && s.policyFactory.GetPolicyWrapper(req.PolicyName) == nil {
		return nil, status.Error(codes.NotFound, "no such policy")
	}

	jsonDynamicConfig, err := req.DynamicConfig.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dynamic config '%s': '%s'", req.PolicyName, err)
	}

	// FIXME use a transaction so that it's not possible to post dynamic config
	// to non-existing policy.
	// Note: Right now we allow setting dynamic config for k8s-managed policy.
	// Do we want to continue supporting that?
	_, err = s.etcdClient.KV.Put(ctx, path.Join(paths.PoliciesAPIDynamicConfigPath, req.PolicyName), string(jsonDynamicConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to write dynamic config '%s' to etcd: '%s'", req.PolicyName, err)
	}

	return new(emptypb.Empty), nil
}

// DeletePolicy deletes a policy from the system.
func (s *PolicyService) DeletePolicy(ctx context.Context, policy *policylangv1.DeletePolicyRequest) (*emptypb.Empty, error) {
	resp, err := s.etcdClient.KV.Delete(ctx, path.Join(paths.PoliciesAPIConfigPath, policy.Name))
	if err != nil {
		return nil, err
	}

	if resp.Deleted > 0 {
		log.Info().Str("policy", policy.Name).Msg("Policy removed from etcd")
	}
	// FIXME: If Deleted==0 we should return NotFound, but first we need to ensure
	// that the delete activity in cloud handles such status correctly.

	_, err = s.etcdClient.KV.Delete(ctx, path.Join(paths.PoliciesAPIDynamicConfigPath, policy.Name))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetDecisions returns the decisions.
func (s *PolicyService) GetDecisions(ctx context.Context, req *policylangv1.GetDecisionsRequest) (*policylangv1.GetDecisionsResponse, error) {
	decisionsPathPrefix := paths.DecisionsPrefix + "/"
	decisionType := ""
	all := true
	if req != nil {
		if req.DecisionType != "" {
			all = false
			decisionType = req.DecisionType
			decisionsPathPrefix += decisionType + "/"
		}
	}

	resp, err := s.etcdClient.Client.KV.Get(ctx, decisionsPathPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	decisions := map[string]string{}
	for _, kv := range resp.Kvs {
		decisionName, ok := strings.CutPrefix(string(kv.Key), decisionsPathPrefix)
		if !ok {
			continue
		}

		if all {
			decisionType = strings.Split(decisionName, "/")[0]
		}

		var m protoreflect.ProtoMessage
		switch decisionType {
		case "load_scheduler":
			m = &policysyncv1.LoadDecisionWrapper{}
		case "rate_limiter", "quota_scheduler":
			m = &policysyncv1.RateLimiterDecisionWrapper{}
		case "pod_scaler":
			m = &policysyncv1.ScaleDecisionWrapper{}
		case "sampler":
			m = &policysyncv1.SamplerDecisionWrapper{}
		}

		err := proto.Unmarshal(kv.Value, m)
		if err != nil {
			return nil, err
		}
		mjson := protojson.Format(m)
		decisions[decisionName] = mjson
	}

	return &policylangv1.GetDecisionsResponse{
		Decisions: decisions,
	}, nil
}
