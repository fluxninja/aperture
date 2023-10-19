package cmd

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	autoscalecontrolpointsv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	entitiesv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/discovery/entities/v1"
	previewv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/agent-functions/agents"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/consts"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/selectors"
	apertureStatus "github.com/fluxninja/aperture/v2/pkg/status"
)

// Handler is a gRPC server for the controller service.
type Handler struct {
	cmdv1.UnimplementedControllerServer
	agents        agents.Agents
	policyService *controlplane.PolicyService
	statusService *apertureStatus.StatusService
}

// NewHandler creates a new Handler.
func NewHandler(
	agents agents.Agents,
	policyService *controlplane.PolicyService,
	statusService *apertureStatus.StatusService,
) *Handler {
	return &Handler{
		agents:        agents,
		policyService: policyService,
		statusService: statusService,
	}
}

// ListAgents lists all agents.
func (h *Handler) ListAgents(
	ctx context.Context,
	req *cmdv1.ListAgentsRequest,
) (*cmdv1.ListAgentsResponse, error) {
	agents, err := h.agents.GetAgentsForGroup(req.AgentGroup)
	if err != nil {
		return nil, err
	}
	return &cmdv1.ListAgentsResponse{
		Agents: agents,
	}, nil
}

// ListFlowControlPoints lists all FlowControlPoints.
func (h *Handler) ListFlowControlPoints(
	ctx context.Context,
	_ *cmdv1.ListFlowControlPointsRequest,
) (*cmdv1.ListFlowControlPointsControllerResponse, error) {
	agentsControlPoints, err := h.agents.ListFlowControlPoints()
	if err != nil {
		return nil, err
	}

	numErrors := uint32(0)
	allControlPoints := map[selectors.TypedGlobalControlPointID]struct{}{}
	for _, resp := range agentsControlPoints {
		if resp.Err != nil {
			numErrors += 1
			continue
		}

		for _, protoCp := range resp.Success.FlowControlPoints.FlowControlPoints {
			gcp := selectors.TypedControlPointIDFromProto(protoCp).InAgentGroup(resp.Success.AgentGroup)
			allControlPoints[gcp] = struct{}{}
		}
	}

	protoControlPoints := make([]*cmdv1.GlobalFlowControlPoint, 0, len(allControlPoints))
	for cp := range allControlPoints {
		protoControlPoints = append(protoControlPoints, cp.ToProto())
	}

	return &cmdv1.ListFlowControlPointsControllerResponse{
		GlobalFlowControlPoints: protoControlPoints,
		ErrorsCount:             numErrors,
	}, nil
}

// AutoScaleControlPointID is a ControlPointID without an agent group.
type AutoScaleControlPointID struct {
	APIVersion string
	Kind       string
	Namespace  string
	Name       string
}

// GlobalAutoScaleControlPointID is a ControlPointID with an agent group.
type GlobalAutoScaleControlPointID struct {
	AutoScaleControlPointID
	AgentGroup string
}

// ToProto converts ControlPointID to protobuf representation.
func (gcp GlobalAutoScaleControlPointID) ToProto() *cmdv1.GlobalAutoScaleControlPoint {
	return &cmdv1.GlobalAutoScaleControlPoint{
		AgentGroup: gcp.AgentGroup,
		AutoScaleControlPoint: &autoscalecontrolpointsv1.AutoScaleKubernetesControlPoint{
			ApiVersion: gcp.APIVersion,
			Kind:       gcp.Kind,
			Namespace:  gcp.Namespace,
			Name:       gcp.Name,
		},
	}
}

// ListAutoScaleControlPoints lists all AutoScaleControlPoints.
func (h *Handler) ListAutoScaleControlPoints(
	ctx context.Context,
	_ *cmdv1.ListAutoScaleControlPointsRequest,
) (*cmdv1.ListAutoScaleControlPointsControllerResponse, error) {
	agentsControlPoints, err := h.agents.ListAutoScaleControlPoints()
	if err != nil {
		return nil, err
	}

	numErrors := uint32(0)
	allControlPoints := map[GlobalAutoScaleControlPointID]struct{}{}
	for _, resp := range agentsControlPoints {
		if resp.Err != nil {
			numErrors += 1
			continue
		}

		for _, protoCp := range resp.Success.AutoScaleControlPoints.AutoScaleKubernetesControlPoints {
			gcp := GlobalAutoScaleControlPointID{
				AutoScaleControlPointID: AutoScaleControlPointIDFromProto(protoCp),
				AgentGroup:              resp.Success.AgentGroup,
			}
			allControlPoints[gcp] = struct{}{}
		}
	}

	protoControlPoints := make([]*cmdv1.GlobalAutoScaleControlPoint, 0, len(allControlPoints))
	for cp := range allControlPoints {
		protoControlPoints = append(protoControlPoints, cp.ToProto())
	}

	return &cmdv1.ListAutoScaleControlPointsControllerResponse{
		GlobalAutoScaleControlPoints: protoControlPoints,
		ErrorsCount:                  numErrors,
	}, nil
}

// ListDiscoveryEntities lists all Discovery entities.
func (h *Handler) ListDiscoveryEntities(ctx context.Context, req *cmdv1.ListDiscoveryEntitiesRequest) (*cmdv1.ListDiscoveryEntitiesControllerResponse, error) {
	discoveryEntities, err := h.agents.ListDiscoveryEntities(req.AgentGroup)
	if err != nil {
		return nil, err
	}

	numErrors := uint32(0)
	entities := make(map[string]*entitiesv1.Entity)
	for _, resp := range discoveryEntities {
		if resp.Err != nil {
			numErrors += 1
			continue
		}

		for k, v := range resp.Success.Entities {
			entities[k] = v
		}
	}

	return &cmdv1.ListDiscoveryEntitiesControllerResponse{
		Entities: &cmdv1.ListDiscoveryEntitiesAgentResponse{
			Entities: entities,
		},
		ErrorsCount: numErrors,
	}, nil
}

// ListDiscoveryEntity lists all Discovery entity.
func (h *Handler) ListDiscoveryEntity(ctx context.Context, req *cmdv1.ListDiscoveryEntityRequest) (*cmdv1.ListDiscoveryEntityAgentResponse, error) {
	discoveryEntity, err := h.agents.ListDiscoveryEntity(req)
	if err != nil {
		return nil, err
	}

	var entity *entitiesv1.Entity
	if discoveryEntity.Entity != nil {
		services := make([]string, 0, len(discoveryEntity.Entity.Services))
		services = append(services, discoveryEntity.Entity.Services...)
		entity = &entitiesv1.Entity{
			Uid:       discoveryEntity.Entity.Uid,
			IpAddress: discoveryEntity.Entity.IpAddress,
			Name:      discoveryEntity.Entity.Name,
			Namespace: discoveryEntity.Entity.Namespace,
			NodeName:  discoveryEntity.Entity.NodeName,
			Services:  services,
		}
	}

	return &cmdv1.ListDiscoveryEntityAgentResponse{
		Entity: entity,
	}, nil
}

// AutoScaleControlPointIDFromProto creates ControlPointID from protobuf representation.
func AutoScaleControlPointIDFromProto(protoCP *autoscalecontrolpointsv1.AutoScaleKubernetesControlPoint) AutoScaleControlPointID {
	return AutoScaleControlPointID{
		APIVersion: protoCP.ApiVersion,
		Kind:       protoCP.Kind,
		Namespace:  protoCP.Namespace,
		Name:       protoCP.Name,
	}
}

// PreviewFlowLabels previews flow labels.
func (h *Handler) PreviewFlowLabels(
	ctx context.Context,
	req *cmdv1.PreviewFlowLabelsRequest,
) (*cmdv1.PreviewFlowLabelsControllerResponse, error) {
	cp := mkGlobalControlPoint(req.AgentGroup, req.Request)
	return doPreview(
		ctx,
		h,
		req.Request,
		cp,
		h.agents.PreviewFlowLabels,
		func(
			resp *previewv1.PreviewFlowLabelsResponse,
		) *cmdv1.PreviewFlowLabelsControllerResponse {
			return &cmdv1.PreviewFlowLabelsControllerResponse{Response: resp}
		},
	)
}

// PreviewHTTPRequests previews HTTP requests.
func (h *Handler) PreviewHTTPRequests(
	ctx context.Context,
	req *cmdv1.PreviewHTTPRequestsRequest,
) (*cmdv1.PreviewHTTPRequestsControllerResponse, error) {
	cp := mkGlobalControlPoint(req.AgentGroup, req.Request)
	return doPreview(
		ctx,
		h,
		req.Request,
		cp,
		h.agents.PreviewHTTPRequests,
		func(
			resp *previewv1.PreviewHTTPRequestsResponse,
		) *cmdv1.PreviewHTTPRequestsControllerResponse {
			return &cmdv1.PreviewHTTPRequestsControllerResponse{Response: resp}
		},
	)
}

// Helper to handle both flow labels and HTTP requests preview.
func doPreview[AgentResp, ControllerResp any](
	ctx context.Context,
	h *Handler,
	req *previewv1.PreviewRequest,
	cp selectors.GlobalControlPointID,
	preview func(string, *previewv1.PreviewRequest) (AgentResp, error),
	wrap func(AgentResp) ControllerResp,
) (ControllerResp, error) {
	var nilResp ControllerResp

	candidates, err := h.agentsWithControlPoint(ctx, cp)
	if err != nil {
		return nilResp, err
	}

	if len(candidates) == 0 {
		return nilResp, status.Error(
			codes.FailedPrecondition,
			"no agent with such control point",
		)
	}

	disabledCount := 0

	for _, agent := range candidates {
		select {
		// FIXME Propagate ctx as upper bound rpc timeout instead of checking manually.
		case <-ctx.Done():
			return nilResp, status.Error(codes.DeadlineExceeded, "timeout")
		default:
		}

		resp, err := preview(agent, req)
		if err != nil {
			if status.Code(err) == codes.FailedPrecondition {
				// This error is only returned by agent in case of disabled preview.
				disabledCount += 1
			}
			continue
		}

		return wrap(resp), nil
	}

	if len(candidates) == disabledCount {
		return nilResp, status.Error(codes.FailedPrecondition, "preview disabled")
	}

	return nilResp, status.Error(codes.Unavailable, "no agent with enough samples")
}

func mkGlobalControlPoint(
	agentGroup string,
	req *previewv1.PreviewRequest,
) selectors.GlobalControlPointID {
	return selectors.GlobalControlPointID{
		ControlPointID: selectors.ControlPointID{
			ControlPoint: req.GetControlPoint(),
			Service:      req.GetService(),
		},
		AgentGroup: agentGroup,
	}
}

func (h *Handler) agentsWithControlPoint(
	ctx context.Context,
	needle selectors.GlobalControlPointID,
) ([]string, error) {
	// FIXME We could narrow down list of agents to ask for control points if
	// we'd cache agent groups.
	// FIXME we can add argument to ListFlowControlPoints for agents to filter
	// non-matching control points.
	agentsControlPoints, err := h.agents.ListFlowControlPoints()
	if err != nil {
		return nil, err
	}

	var agents []string
agentsLoop:
	for _, agent := range agentsControlPoints {
		if agent.Err != nil {
			// FIXME Perhaps gather errors in some stats?
			continue
		}

		if agent.Success.AgentGroup != needle.AgentGroup {
			continue
		}

		for _, cpProto := range agent.Success.FlowControlPoints.FlowControlPoints {
			cp := selectors.TypedControlPointIDFromProto(cpProto)

			if cp.ControlPoint != needle.ControlPoint {
				continue
			}

			if needle.Service == consts.AnyService || cp.Service == needle.Service {
				agents = append(agents, agent.Client)
				continue agentsLoop // avoid duplicating agent
			}
		}
	}

	return agents, nil
}

// UpsertPolicy creates/updates policies in the system.
func (h *Handler) UpsertPolicy(ctx context.Context, req *policylangv1.UpsertPolicyRequest) (*policylangv1.UpsertPolicyResponse, error) {
	return h.policyService.UpsertPolicy(ctx, req)
}

// PostDynamicConfig updates dynamic-config in the system.
func (h *Handler) PostDynamicConfig(ctx context.Context, req *policylangv1.PostDynamicConfigRequest) (*emptypb.Empty, error) {
	return h.policyService.PostDynamicConfig(ctx, req)
}

// GetDynamicConfig gets dynamic-config of a policy.
func (h *Handler) GetDynamicConfig(ctx context.Context, req *policylangv1.GetDynamicConfigRequest) (*policylangv1.GetDynamicConfigResponse, error) {
	return h.policyService.GetDynamicConfig(ctx, req)
}

// DeleteDynamicConfig deletes dynamic-config of a policy.
func (h *Handler) DeleteDynamicConfig(ctx context.Context, req *policylangv1.DeleteDynamicConfigRequest) (*emptypb.Empty, error) {
	return h.policyService.DeleteDynamicConfig(ctx, req)
}

// DeletePolicy deletes policies from the system.
func (h *Handler) DeletePolicy(ctx context.Context, req *policylangv1.DeletePolicyRequest) (*emptypb.Empty, error) {
	return h.policyService.DeletePolicy(ctx, req)
}

// GetDecisions returns decisions.
func (h *Handler) GetDecisions(ctx context.Context, req *policylangv1.GetDecisionsRequest) (*policylangv1.GetDecisionsResponse, error) {
	return h.policyService.GetDecisions(ctx, req)
}

// ListPolicies returns all applied policies.
func (h *Handler) ListPolicies(ctx context.Context, _ *emptypb.Empty) (*policylangv1.GetPoliciesResponse, error) {
	return h.policyService.GetPolicies(ctx, nil)
}

func (h *Handler) GetPolicy(ctx context.Context, req *policylangv1.GetPolicyRequest) (*policylangv1.GetPolicyResponse, error) {
	return h.policyService.GetPolicy(ctx, req)
}

// GetStatus returns status of jobs in the system.
func (h *Handler) GetStatus(ctx context.Context, req *statusv1.GroupStatusRequest) (*statusv1.GroupStatus, error) {
	return h.statusService.GetGroupStatus(ctx, req)
}
