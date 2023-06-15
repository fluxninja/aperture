package agent

import (
	"fmt"
	"path"

	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/peers"
)

// ProvidePeersPrefix provides the peers prefix.
func ProvidePeersPrefix(agentInfo *agentinfo.AgentInfo) (peers.PeerDiscoveryPrefix, error) {
	agentGroup := agentInfo.GetAgentGroup()
	prefix := path.Join(info.Service, agentGroup)
	return peers.PeerDiscoveryPrefix(prefix), nil
}

// AddAgentInfoAttribute adds the agent group and instance labels to OTel config.
func AddAgentInfoAttribute(
	agentInfo *agentinfo.AgentInfo,
	otelConfigProvider *otelconfig.Provider,
) {
	// Note: This hook is added before starting the collector, so there's no
	// risk of collector seeing config with missing processors.
	otelConfigProvider.AddMutatingHook(func(cfg *otelconfig.Config) {
		cfg.AddProcessor(otelconsts.ProcessorAgentGroup, map[string]interface{}{
			"actions": []map[string]interface{}{
				{
					"key":    otelconsts.AgentGroupLabel,
					"action": "insert",
					"value":  agentInfo.GetAgentGroup(),
				},
			},
		})
		transformStatements := []map[string]interface{}{
			{
				"context": "resource",
				"statements": []string{
					fmt.Sprintf(`set(attributes["%v"], "%v")`,
						otelconsts.AgentGroupLabel, agentInfo.GetAgentGroup()),
					fmt.Sprintf(`set(attributes["%v"], "%v")`,
						otelconsts.InstanceLabel, info.Hostname),
					fmt.Sprintf(`set(attributes["%v"], "%v")`,
						metrics.ProcessUUIDLabel, info.UUID),
				},
			},
		}
		cfg.AddProcessor(otelconsts.ProcessorAgentResourceLabels, map[string]interface{}{
			"log_statements":    transformStatements,
			"metric_statements": transformStatements,
		})
	})
}
