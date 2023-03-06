// Server-side for handling agent functions
package agents

import (
	"go.uber.org/fx"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// Module is fx module for controlling Agents on controller side.
var Module = fx.Provide(NewAgents)

// Agents wraps rpc.Clients where clients are agents and provides wrapper
//
// Agents wraps functions registered in agentfunctions, types should match.
type Agents struct{ *rpc.Clients }

// NewAgents wraps Clients with Agent-specific function wrappers.
func NewAgents(clients *rpc.Clients) Agents { return Agents{Clients: clients} }

// ListControlPoints lists control points of all agents.
//
// Handled by agentfunctions.ControlPointsHandler.
func (a Agents) ListControlPoints() ([]rpc.Result[*cmdv1.ListFlowControlPointsAgentResponse], error) {
	var req cmdv1.ListFlowControlPointsRequest
	return rpc.CallAll[cmdv1.ListFlowControlPointsAgentResponse](a.Clients, &req)
}

// PreviewFlowLabels previews flow labels on a given agent.
//
// Handled by agentfunctions.PreviewHandler.
func (a Agents) PreviewFlowLabels(
	agent string,
	req *previewv1.PreviewRequest,
) (*previewv1.PreviewFlowLabelsResponse, error) {
	return rpc.Call[previewv1.PreviewFlowLabelsResponse](
		a.Clients,
		agent,
		&cmdv1.PreviewFlowLabelsRequest{Request: req},
	)
}

// PreviewHTTPRequests previews flow labels on a given agent.
//
// Handled by agentfunctions.PreviewHandler.
func (a Agents) PreviewHTTPRequests(
	agent string,
	req *previewv1.PreviewRequest,
) (*previewv1.PreviewHTTPRequestsResponse, error) {
	return rpc.Call[previewv1.PreviewHTTPRequestsResponse](
		a.Clients,
		agent,
		&cmdv1.PreviewHTTPRequestsRequest{Request: req},
	)
}
