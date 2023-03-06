package cmd

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	"github.com/fluxninja/aperture/pkg/agentfunctions/agents"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

type Handler struct {
	cmdv1.UnimplementedControllerServer
	agents agents.Agents
}

func NewHandler(agents agents.Agents) *Handler {
	return &Handler{agents: agents}
}

func (h *Handler) ListAgents(
	ctx context.Context,
	_ *emptypb.Empty,
) (*cmdv1.ListAgentsResponse, error) {
	return &cmdv1.ListAgentsResponse{
		Agents: h.agents.List(),
	}, nil
}

func (h *Handler) ListControlPoints(
	ctx context.Context,
	_ *cmdv1.ListFlowControlPointsRequest,
) (*cmdv1.ListFlowControlPointsControllerResponse, error) {
	agentsControlPoints, err := h.agents.ListControlPoints()
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

		for _, protoCp := range resp.Success.FlowControlPoints.FlowControlPoints {
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
	// FIXME we can add argument to ListControlPoints for agents to filter
	// non-matching control points.
	agentsControlPoints, err := h.agents.ListControlPoints()
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
