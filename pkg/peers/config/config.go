// +kubebuilder:validation:Optional
package config

// PeerDiscoveryConfig holds configuration for Agent Peer Discovery.
// swagger:model
// +kubebuilder:object:generate=true
type PeerDiscoveryConfig struct {
	// Network address of aperture server to advertise to peers - this address
	// should be reachable from other agents. Used for nat traversal when
	// provided.
	AdvertisementAddr string `json:"advertisement_addr" validate:"omitempty,hostname_port"`
}
