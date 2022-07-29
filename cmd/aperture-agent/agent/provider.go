package agent

import (
	"path"

	"github.com/FluxNinja/aperture/pkg/agentinfo"
	"github.com/FluxNinja/aperture/pkg/info"
	"github.com/FluxNinja/aperture/pkg/peers"
)

// ProvidePeersPrefix provides the peers prefix.
func ProvidePeersPrefix(agentInfo *agentinfo.AgentInfo) (peers.PeerDiscoveryPrefix, error) {
	agentGroup := agentInfo.GetAgentGroup()
	prefix := path.Join(info.Service, agentGroup)
	return peers.PeerDiscoveryPrefix(prefix), nil
}
