// +kubebuilder:validation:Optional
package config

// Config stores configuration for the Google Token Source.
// swagger:model
// +kubebuilder:object:generate=true
type Config struct {
	Enabled bool     `json:"enabled" default:"false"`
	Scopes  []string `json:"scopes,omitempty"`
}
