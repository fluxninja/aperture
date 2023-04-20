// +kubebuilder:validation:Optional
package config

// PrometheusConfig holds configuration for Prometheus Server.
// swagger:model
// +kubebuilder:object:generate=true
type PrometheusConfig struct {
	// Address of the Prometheus server
	Address string `json:"address" validate:"required,hostname_port|url|fqdn"`
}
