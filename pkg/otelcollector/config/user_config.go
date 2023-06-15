// +kubebuilder:validation:Optional
package otelconfig

import "github.com/fluxninja/aperture/v2/pkg/config"

// NewDefaultCommonOTelConfig creates CommonOTelConfig with all the default values set.
func NewDefaultCommonOTelConfig() *CommonOTelConfig {
	return &CommonOTelConfig{
		Ports: PortsConfig{
			DebugPort:       8888,
			HealthCheckPort: 13133,
			PprofPort:       1777,
			ZpagesPort:      55679,
		},
	}
}

// CommonOTelConfig is the configuration for the OTel collector.
// swagger:model
// +kubebuilder:object:generate=true
type CommonOTelConfig struct {
	// BatchAlerts configures batch alerts processor.
	BatchAlerts BatchAlertsConfig `json:"batch_alerts"`
	// Ports configures debug, health and extension ports values.
	Ports PortsConfig `json:"ports"`
}

// PortsConfig defines configuration for OTel debug and extension ports.
// swagger:model
// +kubebuilder:object:generate=true
type PortsConfig struct {
	// Port on which OTel collector exposes Prometheus metrics on /metrics path.
	DebugPort uint32 `json:"debug_port" validate:"gte=0" default:"8888"`
	// Port on which health check extension in exposed.
	HealthCheckPort uint32 `json:"health_check_port" validate:"gte=0" default:"13133"`
	// Port on which `pprof` extension in exposed.
	PprofPort uint32 `json:"pprof_port" validate:"gte=0" default:"1777"`
	// Port on which `zpages` extension in exposed.
	ZpagesPort uint32 `json:"zpages_port" validate:"gte=0" default:"55679"`
}

// BatchAlertsConfig defines configuration for OTel batch processor.
// swagger:model
// +kubebuilder:object:generate=true
type BatchAlertsConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"100"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"100"`
}
