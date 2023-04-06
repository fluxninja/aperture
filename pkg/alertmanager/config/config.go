package alertmgrconfig

import (
	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
)

// AlertManagerConfig main level config for alertmanager.
// swagger:model
// +kubebuilder:object:generate=true
type AlertManagerConfig struct {
	Clients []AlertManagerClientConfig `json:"clients,omitempty"`
}

// AlertManagerClientConfig config for single alertmanager client.
// swagger:model AlertManagerClientConfig
// +kubebuilder:object:generate=true
type AlertManagerClientConfig struct {
	Name       string                      `json:"name"`
	Address    string                      `json:"address" validate:"hostname_port|url|fqdn"`
	BasePath   string                      `json:"base_path" default:"/"`
	HTTPConfig commonhttp.HTTPClientConfig `json:"http_client"`
}
