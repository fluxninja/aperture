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

// FluxNinjaPluginConfig is the configuration for FluxNinja cloud integration plugin.
// swagger:model
// +kubebuilder:object:generate=true
type FluxNinjaPluginConfig struct {
	// Interval between each heartbeat.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="5s"
	HeartbeatInterval config.Duration `json:"heartbeat_interval" validate:"gte=0s" default:"5s"`
	// Address to grpc or http(s) server listening in agent service. To use http protocol, the address must start with http(s)://.
	//+kubebuilder:validation:Optional
	FluxNinjaEndpoint string `json:"fluxninja_endpoint,omitempty" validate:"omitempty,hostname_port|url|fqdn"`
	// API Key for this agent.
	//+kubebuilder:validation:Optional
	APIKey string `json:"api_key,omitempty"`
	// Client configuration.
	//+kubebuilder:validation:Optional
	ClientConfig ClientConfig `json:"client"`
}

// ClientConfig is the client configuration.
// swagger:model
// +kubebuilder:object:generate=true
type ClientConfig struct {
	// HTTP client settings.
	//+kubebuilder:validation:Optional
	HTTPClient http.HTTPClientConfig `json:"http"`
	// GRPC client settings.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={backoff:{base_delay:"1s",multiplier:1.6}}
	GRPCClient grpc.GRPCClientConfig `json:"grpc"`
}
