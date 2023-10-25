// +kubebuilder:validation:Optional
package agentinfo

import (
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/utils"
	"go.uber.org/fx"
)

const (
	configKey = "agent_info"
)

// InstallationModeConfig can be provided by an extension to provide mode of installation.
type InstallationModeConfig struct {
	InstallationMode string
}

// swagger:operation POST /agent_info common-configuration AgentInfo
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/AgentInfoConfig"

// AgentInfoConfig is the configuration for the agent group and other agent attributes.
// swagger:model
// +kubebuilder:object:generate=true
type AgentInfoConfig struct {
	// All agents within an agent group receive the same data-plane configuration (for example, Flux Meters, Rate Limiters and so on).
	//
	// [Read more about agent groups here](/concepts/selector.md#agent-group).
	AgentGroup string `json:"agent_group" default:"default"`
}

// AgentInfo is the agent info.
type AgentInfo struct {
	agentGroup string
}

// AgentInfoIn holds parameters for ProvideAgentInfo.
type AgentInfoIn struct {
	fx.In

	Unmarshaller           config.Unmarshaller
	InstallationModeConfig *InstallationModeConfig `optional:"true"`
}

// ProvideAgentInfo provides the agent info via Fx.
func ProvideAgentInfo(in AgentInfoIn) (*AgentInfo, error) {
	var config AgentInfoConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &config); err != nil {
		return nil, err
	}

	if in.InstallationModeConfig != nil && in.InstallationModeConfig.InstallationMode != utils.InstallationModeCloudAgent && config.AgentGroup == utils.ApertureCloudAgentGroup {
		return nil, fmt.Errorf("'%s' is a reserved group name for FluxNinja Cloud Agents. Please use a different agent group name", utils.ApertureCloudAgentGroup)
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
