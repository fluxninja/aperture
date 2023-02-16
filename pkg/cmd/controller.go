package cmd

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
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
	_ *cmdv1.ListControlPointsRequest,
) (*cmdv1.ListControlPointsResponse, error) {
	agentsControlPoints, err := h.agents.ListControlPoints()
	if err != nil {
		return nil, err
	}

	numErrors := uint32(0)
	allControlPoints := map[selectors.ControlPointID]struct{}{}
	for _, resp := range agentsControlPoints {
		if resp.Err != nil {
			numErrors += 1
			continue
		}

		for _, cp := range resp.Success.ControlPoints {
			allControlPoints[selectors.ControlPointIDFromProto(cp)] = struct{}{}
		}
	}

	protoControlPoints := make([]*cmdv1.ServiceControlPoint, 0, len(allControlPoints))
	for cp := range allControlPoints {
		protoControlPoints = append(protoControlPoints, cp.ToProto())
	}

	return &cmdv1.ListControlPointsResponse{
		ControlPoints: protoControlPoints,
		ErrorsCount:   numErrors,
	}, nil
}
