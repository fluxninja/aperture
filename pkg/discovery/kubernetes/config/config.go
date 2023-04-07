// +kubebuilder:validation:Optional
package config

import "github.com/fluxninja/aperture/pkg/discovery/common"

// Key is the key for the Kubernetes discovery configuration.
var Key = common.DiscoveryConfigKey + ".kubernetes"

// KubernetesDiscoveryConfig for Kubernetes service discovery.
// swagger:model
// +kubebuilder:object:generate=true
type KubernetesDiscoveryConfig struct {
	Enabled bool `json:"enabled" default:"true"`
}
