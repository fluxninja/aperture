package cmd

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	autoscalecontrolpointsv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/autoscale/kubernetes/controlpoints/v1"
	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	"github.com/fluxninja/aperture/pkg/agentfunctions/agents"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

// Handler is a gRPC server for the controller service.
type Handler struct {
	cmdv1.UnimplementedControllerServer
	agents agents.Agents
}

// NewHandler creates a new Handler.
func NewHandler(agents agents.Agents) *Handler {
	return &Handler{agents: agents}
}

// ListAgents lists all agents.
func (h *Handler) ListAgents(
	ctx context.Context,
	_ *emptypb.Empty,
) (*cmdv1.ListAgentsResponse, error) {
	return &cmdv1.ListAgentsResponse{
		Agents: h.agents.List(),
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
	allControlPoints := map[selectors.GlobalControlPointID]struct{}{}
	for _, resp := range agentsControlPoints {
		if resp.Err != nil {
			numErrors += 1
			continue
		}

		for _, protoCp := range resp.Success.FlowControlPoints {
			gcp := selectors.ControlPointIDFromProto(protoCp).InAgentGroup(resp.Success.AgentGroup)
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

		for _, protoCp := range resp.Success.AutoScaleControlPoints {
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

		for _, cpProto := range agent.Success.FlowControlPoints {
			cp := selectors.ControlPointIDFromProto(cpProto)

			if cp.ControlPoint != needle.ControlPoint {
				continue
			}

			// FIXME This "all" thing shouldn't be hardcoded.
			if needle.Service == "all" || cp.Service == needle.Service {
				agents = append(agents, agent.Client)
				continue agentsLoop // avoid duplicating agent
			}
		}
	}

	return agents, nil
}
