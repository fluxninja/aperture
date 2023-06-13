// +kubebuilder:validation:Optional
package config

import (
	"github.com/fluxninja/aperture/v2/pkg/config"
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
)

// swagger:operation POST /otel agent-configuration OTel
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/AgentOTelConfig"

// AgentOTelConfig is the configuration for Agent's OTel collector.
// +kubebuilder:object:generate=true
//
//swagger:model
type AgentOTelConfig struct {
	// BatchPrerollup configures the OTel batch pre-processor.
	BatchPrerollup BatchPrerollupConfig `json:"batch_prerollup"`
	// BatchPostrollup configures the OTel batch post-processor.
	BatchPostrollup             BatchPostrollupConfig `json:"batch_postrollup"`
	otelconfig.CommonOTelConfig `json:",inline"`
	// DisableKubernetesScraper disables the default metrics collection for Kubernetes resources.
	DisableKubernetesScraper bool `json:"disable_kubernetes_scraper" default:"false"`
	// DisableKubeletScraper disables the default metrics collection for kubelet.
	// Deprecated: kubelet scraper is removed entirely, so this flag makes no difference.
	DisableKubeletScraper bool `json:"disable_kubelet_scraper" default:"false"`
}

// BatchPrerollupConfig defines configuration for OTel batch processor.
// +kubebuilder:object:generate=true
//
//swagger:model
type BatchPrerollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"10s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"10000"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"10000"`
}

// BatchPostrollupConfig defines configuration for OTel batch processor.
// +kubebuilder:object:generate=true
//
//swagger:model
type BatchPostrollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the size of a batch which after hit, will trigger it to be sent.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"100"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"100"`
}
