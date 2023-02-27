// +kubebuilder:validation:Optional
package pluginconfig

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/http"
)

const (
	// PluginConfigKey is the key for the plugin configuration.
	PluginConfigKey = "fluxninja_plugin"
)

// FluxNinjaPluginConfig is the configuration for FluxNinja ARC integration plugin.
// swagger:model
// +kubebuilder:object:generate=true
type FluxNinjaPluginConfig struct {
	// Interval between each heartbeat.
	HeartbeatInterval config.Duration `json:"heartbeat_interval" validate:"gte=0s" default:"5s"`
	// Address to grpc or http(s) server listening in agent service. To use http protocol, the address must start with http(s)://.
	FluxNinjaEndpoint string `json:"fluxninja_endpoint" validate:"omitempty,hostname_port|url|fqdn"`
	// API Key for this agent.
	APIKey string `json:"api_key"`
	// Client configuration.
	ClientConfig ClientConfig `json:"client"`
}

// ClientConfig is the client configuration.
// swagger:model
// +kubebuilder:object:generate=true
type ClientConfig struct {
	// HTTP client settings.
	HTTPClient http.HTTPClientConfig `json:"http"`
	// GRPC client settings.
	GRPCClient grpc.GRPCClientConfig `json:"grpc"`
}
