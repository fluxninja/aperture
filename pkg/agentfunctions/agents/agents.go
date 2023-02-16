// Server-side for handling agent functions
package agents

import (
	"go.uber.org/fx"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/pkg/rpc"
)

// Module is fx module for controlling Agents on controller side
var Module = fx.Provide(NewAgents)

// Agents wraps rpc.Clients where clients are agents
type Agents struct{ *rpc.Clients }

// NewAgents wraps Clients with Agent-specific function wrappers.
func NewAgents(clients *rpc.Clients) Agents { return Agents{Clients: clients} }

// ListControlPoints lists control points of all agents
func (a Agents) ListControlPoints() ([]rpc.Result[*cmdv1.ListControlPointsResponse], error) {
	var req cmdv1.ListControlPointsRequest
	return rpc.CallAll[cmdv1.ListControlPointsResponse](a.Clients, &req)
}
