// +kubebuilder:validation:Optional
// +kubebuilder:object:generate=true
package config

// Entity represents a pod, VM, and so on.
type Entity struct {
	// Unique identifier of the entity.
	UID string `json:"uid,omitempty" validate:"required"` // @gotags: validate:"required"
	// IP address of the entity.
	IPAddress string `json:"ip_address,omitempty" validate:"required,ip"` // @gotags: validate:"required,ip"
	// Name of the entity. For example, pod name.
	Name string `json:"name,omitempty"`
	// Namespace of the entity. For example, pod namespace.
	Namespace string `json:"namespace,omitempty"`
	// Node name of the entity. For example, hostname.
	NodeName string ` json:"node_name,omitempty"`
	// Services of the entity.
	Services []string `json:"services,omitempty" validate:"gt=0"` // @gotags: validate:"gt=0"
}

// StaticDiscoveryConfig for pre-determined list of services.
// swagger:model
type StaticDiscoveryConfig struct {
	Entities []Entity `json:"entities,omitempty"`
}
