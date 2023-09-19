package agentelection

import (
	"go.uber.org/fx"

	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/etcd/election"
)

// Module is a fx module that provides etcd based leader election per agent group.
func Module() fx.Option {
	return fx.Provide(ProvideAgentElection)
}

// ProvideAgentElection provides a wrapper around etcd based leader election for agents.
func ProvideAgentElection(
	in election.ElectionIn,
	agentInfo *agentinfo.AgentInfo,
) *election.Election {
	return election.ProvideElection("/election/"+agentInfo.GetAgentGroup(), in)
}
