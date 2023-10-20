// +kubebuilder:validation:Optional
package extconfig

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/net/grpc"
	"github.com/fluxninja/aperture/v2/pkg/net/http"
	"github.com/fluxninja/aperture/v2/pkg/prometheus"
)

const (
	// ExtensionConfigKey is the key for the extension configuration.
	ExtensionConfigKey = "fluxninja"
)

// FluxNinjaExtensionConfig is the configuration for [FluxNinja integration](/reference/fluxninja.md).
// swagger:model
// +kubebuilder:object:generate=true
type FluxNinjaExtensionConfig struct {
	// Address to gRPC or HTTP(s) server listening in agent service.
	// For connecting to Aperture Cloud Controller, the `endpoint` should be a `grpc/http2` address.
	// For self-hosted controller, the HTTP protocol address can start with `http(s)://`.
	Endpoint string `json:"endpoint" validate:"omitempty,hostname_port|url|fqdn"`
	// Deprecated: v3.0.0. Use AgentAPIKey instead
	APIKey string `json:"api_key" validate:"deprecated"`
	// API Key for this agent. If this key is not set, the extension won't be enabled.
	AgentAPIKey string `json:"agent_api_key"`
	// Installation mode describes on which underlying platform the Agent or the Controller is being run.
	InstallationMode string `json:"installation_mode" validate:"oneof=KUBERNETES_SIDECAR KUBERNETES_DAEMONSET LINUX_BARE_METAL" default:"LINUX_BARE_METAL"`
	// Whether to connect to [Aperture Cloud Controller](/reference/fluxninja.md).
	//
	// Enabling this flag configures various agent components to point to the
	// Aperture Cloud Controller, for example configures remote etcd endpoint and disables
	// local Prometheus OTel pipelines.
	//
	// Disable this flag only if using [Self-Hosted](/get-started/self-hosting/self-hosting.md) Aperture Controller.
	EnableCloudController bool `json:"enable_cloud_controller" default:"false"`
	// Interval between each heartbeat.
	HeartbeatInterval config.Duration `json:"heartbeat_interval" validate:"gte=0s" default:"5s"`
	// Client configuration.
	ClientConfig ClientConfig `json:"client"`
	// Disables local Prometheus OTel pipelines for metrics. Implied by EnableCloudController.
	DisableLocalOTelPipeline bool `json:"disable_local_otel_pipeline" default:"false"`
	// Overrides Controller ID for Aperture Controller. If not set, random id will be generated and persisted in etcd.
	ControllerID string `json:"controller_id,omitempty"`
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
		fx.Provide(provideEtcdConfigOverride),
		fx.Provide(providePrometheusConfigOverride),
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

func provideEtcdConfigOverride(extensionConfig *FluxNinjaExtensionConfig) *etcdclient.ConfigOverride {
	if extensionConfig.EnableCloudController {
		apiKey := extensionConfig.AgentAPIKey
		if apiKey == "" {
			//nolint:staticcheck // SA1019 read APIKey config for backward compatibility
			apiKey = extensionConfig.APIKey
		}
		return &etcdclient.ConfigOverride{
			Namespace: "",
			Endpoints: []string{extensionConfig.Endpoint},
			PerRPCCredentials: perRPCHeaders{
				headers: map[string]string{
					"apiKey": apiKey,
				},
			},
			OverriderName: "fluxninja extension",
		}
	} else {
		return nil
	}
}

func providePrometheusConfigOverride(extensionConfig *FluxNinjaExtensionConfig) *prometheus.ConfigOverride {
	if extensionConfig.EnableCloudController {
		return &prometheus.ConfigOverride{
			SkipClientCreation: true,
		}
	} else {
		return nil
	}
}

type perRPCHeaders struct {
	headers map[string]string
}

// GetRequestMetadata returns the request headers to be used with the RPC.
func (p perRPCHeaders) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return p.headers, nil
}

// RequireTransportSecurity always returns true for this implementation.
func (p perRPCHeaders) RequireTransportSecurity() bool {
	return true
}
