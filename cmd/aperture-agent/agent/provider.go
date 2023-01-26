package agent

import (
	"fmt"
	"path"

	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/peers"
	"go.uber.org/fx"
)

// ProvidePeersPrefix provides the peers prefix.
func ProvidePeersPrefix(agentInfo *agentinfo.AgentInfo) (peers.PeerDiscoveryPrefix, error) {
	agentGroup := agentInfo.GetAgentGroup()
	prefix := path.Join(info.Service, agentGroup)
	return peers.PeerDiscoveryPrefix(prefix), nil
}

type FxIn struct {
	fx.In
	BaseConfig *otelcollector.OTELConfig `name:"base"`
	AgentInfo  *agentinfo.AgentInfo
}

func AddAgentInfoAttribute(in FxIn) {
	in.BaseConfig.AddProcessor(otelcollector.ProcessorAgentGroup, map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    otelcollector.AgentGroupLabel,
				"action": "insert",
				"value":  in.AgentInfo.GetAgentGroup(),
			},
		},
	})
	in.BaseConfig.AddProcessor(otelcollector.ProcessorAgentResourceLabels, map[string]interface{}{
		"log_statements": []map[string]interface{}{
			{
				"context": "resource",
				"statements": []string{
					fmt.Sprintf(`set(attributes["%v"], "%v")`,
						otelcollector.AgentGroupLabel, in.AgentInfo.GetAgentGroup()),
					fmt.Sprintf(`set(attributes["%v"], "%v")`,
						otelcollector.InstanceLabel, info.Hostname),
				},
			},
		},
	})
}
