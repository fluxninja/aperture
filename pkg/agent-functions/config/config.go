// +kubebuilder:validation:Optional
package config

import (
	grpcclient "github.com/fluxninja/aperture/v2/pkg/net/grpc"
)

// swagger:operation POST /agent_functions common-configuration Agent Functions
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     $ref: "#/definitions/AgentFunctionsConfig"

// Key is the key for agentfunctions configuration.
const Key = "agent_functions"

// AgentFunctionsConfig is configuration for agent functions.
// swagger:model AgentFunctionsConfig
// +kubebuilder:object:generate=true
type AgentFunctionsConfig struct {
	// RPC servers to connect to (which will be able to call agent functions)
	Endpoints []string `json:"endpoints,omitempty" validate:"omitempty,dive,omitempty"`

	// Network client configuration
	ClientConfig ClientConfig `json:"client"`
}

// ClientConfig is configuration for network clients used by agent-functions.
// swagger:model
// +kubebuilder:object:generate=true
type ClientConfig struct {
	// gRPC client settings.
	GRPCClient grpcclient.GRPCClientConfig `json:"grpc"`
}
