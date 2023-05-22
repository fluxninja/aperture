package distcache

import (
	"errors"
	"log"

	"github.com/fluxninja/aperture/v2/pkg/peers"
)

// ServiceDiscovery holds fields needed to implement Olric's Service Discovery interface.
type ServiceDiscovery struct {
	discovery *peers.PeerDiscovery
}

// Initialize initializes the plugin: registers some internal data structures, clients etc.
// This method is not implemented.
func (s *ServiceDiscovery) Initialize() error {
	return nil
}

// SetLogger sets an appropriate logger.
// This method is not implemented.
func (s *ServiceDiscovery) SetLogger(l *log.Logger) {
}

// SetConfig registers plugin configuration.
// This method is not implemented.
func (s *ServiceDiscovery) SetConfig(cfg map[string]interface{}) error {
	return nil
}

// DiscoverPeers returns a list of known Olric nodes.
func (s *ServiceDiscovery) DiscoverPeers() ([]string, error) {
	peers := []string{}

	peerInfos := s.discovery.GetPeers()
	for _, peerInfo := range peerInfos.Peers {
		if olricMemberlistAddr, ok := peerInfo.Services[olricMemberlistServiceName]; ok {
			peers = append(peers, olricMemberlistAddr)
		}
	}

	if len(peers) == 0 {
		return nil, errors.New("no peers found")
	}

	return peers, nil
}

// Register registers this node to a service discovery directory.
// This method is not implemented.
func (s *ServiceDiscovery) Register() error { return nil }

// Deregister removes this node from a service discovery directory.
// This method is not implemented.
func (s *ServiceDiscovery) Deregister() error { return nil }

// Close stops underlying goroutines, if there is any. It should be a blocking call.
// This method is not implemented.
func (s *ServiceDiscovery) Close() error { return nil }
