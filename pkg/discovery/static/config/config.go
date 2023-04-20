// +kubebuilder:validation:Optional
package config

import (
	entitiesv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/discovery/entities/v1"
)

// StaticDiscoveryConfig for pre-determined list of services.
// swagger:model
// +kubebuilder:object:generate=true
type StaticDiscoveryConfig struct {
	Entities []entitiesv1.Entity `json:"entities,omitempty"`
}
