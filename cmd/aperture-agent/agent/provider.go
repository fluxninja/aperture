package agent

import (
	"fmt"
	"path"

	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/info"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
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
	BaseConfig *otelconfig.OTELConfig `name:"base"`
	AgentInfo  *agentinfo.AgentInfo
}

func AddAgentInfoAttribute(in FxIn) {
	in.BaseConfig.AddProcessor(otelconsts.ProcessorAgentGroup, map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"key":    otelconsts.AgentGroupLabel,
				"action": "insert",
				"value":  in.AgentInfo.GetAgentGroup(),
			},
		},
	})
	transformStatements := []map[string]interface{}{
		{
			"context": "resource",
			"statements": []string{
				fmt.Sprintf(`set(attributes["%v"], "%v")`,
					otelconsts.AgentGroupLabel, in.AgentInfo.GetAgentGroup()),
				fmt.Sprintf(`set(attributes["%v"], "%v")`,
					otelconsts.InstanceLabel, info.Hostname),
			},
		},
	}
	in.BaseConfig.AddProcessor(otelconsts.ProcessorAgentResourceLabels, map[string]interface{}{
		"log_statements":    transformStatements,
		"metric_statements": transformStatements,
	})
}
