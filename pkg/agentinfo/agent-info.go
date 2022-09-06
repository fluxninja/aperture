package agentinfo

import "github.com/fluxninja/aperture/pkg/config"

const (
	configKey = "agent_info"
)

// swagger:operation POST /agent_info common-configuration AgentInfo
// ---
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/AgentInfoConfig"

// AgentInfoConfig is the configuration for the agent group etc.
// swagger:model
// +kubebuilder:object:generate=true
type AgentInfoConfig struct {
	// All agents within an agent_group receive the same data-plane configuration (e.g. schedulers, FluxMeters, rate limiter).
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="default"
	AgentGroup string `json:"agent_group" default:"default"`
}

// AgentInfo is the agent info.
type AgentInfo struct {
	agentGroup string
}

// ProvideAgentInfo provides the agent info via Fx.
func ProvideAgentInfo(unmarshaller config.Unmarshaller) (*AgentInfo, error) {
	var config AgentInfoConfig
	if err := unmarshaller.UnmarshalKey(configKey, &config); err != nil {
		return nil, err
	}
	return NewAgentInfo(config.AgentGroup), nil
}

// NewAgentInfo creates a new agent info.
func NewAgentInfo(agentGroup string) *AgentInfo {
	return &AgentInfo{
		agentGroup: agentGroup,
	}
}

// GetAgentGroup returns the agent group.
func (a *AgentInfo) GetAgentGroup() string {
	return a.agentGroup
}
