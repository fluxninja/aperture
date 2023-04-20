// +kubebuilder:validation:Optional
package config

// Key is config path where FlowPreviewConfig is located.
const Key = "flow_control.preview_service"

// FlowPreviewConfig is the configuration for the flow control preview service.
// swagger:model
// +kubebuilder:object:generate=true
type FlowPreviewConfig struct {
	// Enables the flow preview service.
	Enabled bool `json:"enabled" default:"true"`
}
