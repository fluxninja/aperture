// +kubebuilder:validation:Optional
package extconfig

import (
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/http"
)

const (
	// ExtensionConfigKey is the key for the extension configuration.
	ExtensionConfigKey = "fluxninja"
)

// FluxNinjaExtensionConfig is the configuration for FluxNinja ARC integration.
// swagger:model
// +kubebuilder:object:generate=true
type FluxNinjaExtensionConfig struct {
	// Interval between each heartbeat.
	HeartbeatInterval config.Duration `json:"heartbeat_interval" validate:"gte=0s" default:"5s"`
	// Address to gRPC or HTTP(s) server listening in agent service. To use HTTP protocol, the address must start with `http(s)://`.
	Endpoint string `json:"endpoint" validate:"omitempty,hostname_port|url|fqdn"`
	// API Key for this agent. If this key is not set, the extension won't be enabled.
	APIKey string `json:"api_key"`
	// Client configuration.
	ClientConfig ClientConfig `json:"client"`
	// Installation mode describes on which underlying platform the Agent or the Controller is being run.
	InstallationMode string `json:"installation_mode" validate:"oneof=KUBERNETES_SIDECAR KUBERNETES_DAEMONSET LINUX_BARE_METAL" default:"LINUX_BARE_METAL"`
}

// ClientConfig is the client configuration.
// swagger:model
// +kubebuilder:object:generate=true
type ClientConfig struct {
	// HTTP client settings.
	HTTPClient http.HTTPClientConfig `json:"http"`
	// gRPC client settings.
	GRPCClient grpc.GRPCClientConfig `json:"grpc"`
}

// Module provides the FluxNinja extension configuration.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(provideConfig),
	)
}

// provideConfig provides the extension configuration.
func provideConfig(unmarshaller config.Unmarshaller) (*FluxNinjaExtensionConfig, error) {
	var extensionConfig FluxNinjaExtensionConfig
	if err := unmarshaller.UnmarshalKey(ExtensionConfigKey, &extensionConfig); err != nil {
		return nil, err
	}
	return &extensionConfig, nil
}
