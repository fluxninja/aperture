package agent

import (
	"path"

	"aperture.tech/aperture/pkg/agentinfo"
	"aperture.tech/aperture/pkg/info"
	"aperture.tech/aperture/pkg/peers"
)

// ProvidePeersPrefix provides the peers prefix.
func ProvidePeersPrefix(agentInfo *agentinfo.AgentInfo) (peers.PeerDiscoveryPrefix, error) {
	agentGroup := agentInfo.GetAgentGroup()
	prefix := path.Join(info.Service, agentGroup)
	return peers.PeerDiscoveryPrefix(prefix), nil
}
