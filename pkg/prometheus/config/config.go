// +kubebuilder:validation:Optional
package config

// PrometheusLabel holds Name->Value mapping for the label that will
// be attached to every PromQL query executed by the controller.
type PrometheusLabel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PrometheusConfig holds configuration for Prometheus Server.
// swagger:model
// +kubebuilder:object:generate=true
type PrometheusConfig struct {
	// A list of labels to be attached to every query
	Labels []PrometheusLabel `json:"labels,omitempty"`
	// Address of the Prometheus server
	Address string `json:"address" validate:"omitempty,hostname_port|url|fqdn"`
}
