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
	LeaseTTL config.Duration `json:"lease_ttl" validate:"gte=1s" default:"10s"`
	// Authentication
	Username string `json:"username"`
	Password string `json:"password"`
	// Client TLS configuration
	ClientTLSConfig tlsconfig.ClientTLSConfig `json:"tls"`
	// List of etcd server endpoints
	Endpoints []string `json:"endpoints,omitempty" validate:"omitempty,dive,hostname_port|url|fqdn,omitempty"`
	// LogLevel of logs coming from inside the etcd client
	LogLevel string `json:"log_level,omitempty" validate:"omitempty,oneof=debug DEBUG info INFO warn WARN error ERROR dpanic DPANIC panic PANIC fatal FATAL,omitempty" default:"error"`
}
