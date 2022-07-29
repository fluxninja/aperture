package agentinfo

import "github.com/FluxNinja/aperture/pkg/config"

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
type AgentInfoConfig struct {
	// All agents within an agent_group receive the same data-plane configuration (e.g. schedulers, FluxMeters, rate limiter).
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
	return &AgentInfo{
		agentGroup: config.AgentGroup,
	}, nil
}

// GetAgentGroup returns the agent group.
func (a *AgentInfo) GetAgentGroup() string {
	return a.agentGroup
}
