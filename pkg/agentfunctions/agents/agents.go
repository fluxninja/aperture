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

// ListFlowControlPoints lists control points of all agents.
//
// Handled by agentfunctions.ControlPointsHandler.
func (a Agents) ListFlowControlPoints() ([]rpc.Result[*cmdv1.ListFlowControlPointsAgentResponse], error) {
	var req cmdv1.ListFlowControlPointsRequest
	return rpc.CallAll[cmdv1.ListFlowControlPointsAgentResponse](a.Clients, &req)
}

// ListAutoScaleControlPoints lists auto-scale control points of all agents.
func (a Agents) ListAutoScaleControlPoints() ([]rpc.Result[*cmdv1.ListAutoScaleControlPointsAgentResponse], error) {
	var req cmdv1.ListAutoScaleControlPointsRequest
	return rpc.CallAll[cmdv1.ListAutoScaleControlPointsAgentResponse](a.Clients, &req)
}

// ListDiscoveryEntities lists discovery entities.
func (a Agents) ListDiscoveryEntities() ([]rpc.Result[*cmdv1.ListDiscoveryEntitiesAgentResponse], error) {
	var req cmdv1.ListDiscoveryEntitiesRequest
	return rpc.CallAll[cmdv1.ListDiscoveryEntitiesAgentResponse](a.Clients, &req)
}

// ListDiscoveryEntity lists discovery entity by ip address or name.
func (a Agents) ListDiscoveryEntity(req *cmdv1.ListDiscoveryEntityRequest) (*cmdv1.ListDiscoveryEntityAgentResponse, error) {
	return rpc.Call[cmdv1.ListDiscoveryEntityAgentResponse](a.Clients, a.List()[0], req)
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
