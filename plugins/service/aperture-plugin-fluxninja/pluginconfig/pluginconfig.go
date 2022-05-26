package pluginconfig

import (
	"aperture.tech/aperture/pkg/config"
)

const (
	// PluginConfigKey is the key for the plugin configuration.
	PluginConfigKey = "fluxninja_plugin"
)

// FluxNinjaPluginConfig is the configuration for FluxNinja cloud integration plugin.
// swagger:model
type FluxNinjaPluginConfig struct {
	// Interval between each heartbeat.
	HeartbeatInterval config.Duration `json:"heartbeat_interval" validate:"gte=0s" default:"5s"`
	// Address to grpc or http(s) server listening in agent service. To use http protocol, the address must start with http(s)://.
	FluxNinjaEndpoint string `json:"fluxninja_endpoint" validate:"omitempty,hostname_port|url|fqdn"`
	// API Key for this agent.
	APIKey string `json:"api_key"`
}
