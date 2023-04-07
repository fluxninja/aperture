// +kubebuilder:validation:Optional
package config

import (
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
)

// Key is the key for agentfunctions configuration.
const Key = "agent_functions"

// Config is configuration for agent functions.
// swagger:model
// +kubebuilder:object:generate=true
type Config struct {
	// RPC servers to connect to (which will be able to call agent functions)
	Endpoints []string `json:"endpoints,omitempty" validate:"omitempty,dive,omitempty"`

	// Network client configuration
	ClientConfig ClientConfig `json:"client"`
}

// ClientConfig is configuration for network clients used by agent-functions.
// swagger:model
// +kubebuilder:object:generate=true
type ClientConfig struct {
	// GRPC client settings.
	GRPCClient grpcclient.GRPCClientConfig `json:"grpc"`
}
