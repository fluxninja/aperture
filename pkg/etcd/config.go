// +kubebuilder:validation:Optional
package etcd

import (
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/net/tlsconfig"
)

// EtcdConfig holds configuration for etcd client.
// swagger:model
// +kubebuilder:object:generate=true
type EtcdConfig struct {
	// etcd namespace
	Namespace string `json:"namespace" default:"aperture"`
	// Lease time-to-live
	LeaseTTL config.Duration `json:"lease_ttl" validate:"gte=1s" default:"60s"`
	// Authentication
	Username string `json:"username"`
	Password string `json:"password"`
	// Client TLS configuration
	ClientTLSConfig tlsconfig.ClientTLSConfig `json:"tls"`
	// List of etcd server endpoints
	Endpoints []string `json:"endpoints,omitempty" validate:"omitempty,dive,hostname_port|url|fqdn,omitempty"`
}
