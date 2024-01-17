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
	// DisableKubeletScraper disables the default metrics collection for Kubelet.
	// Deprecated: Kubelet scraper is removed entirely, so this flag makes no difference.
	DisableKubeletScraper bool `json:"disable_kubelet_scraper" default:"false"`
	// EnableHighCardinalityPlatformMetrics filters out high cardinality Aperture platform metrics from being
	// published to Prometheus. Filtered out metrics are:
	//   * "grpc_server_handled_total.*"
	//   * "grpc_server_handling_seconds.*"
	//   * "grpc_server_handling_seconds_bucket.*"
	//   * "grpc_server_handling_seconds_count.*"
	//   * "grpc_server_handling_seconds_sum.*"
	//   * "grpc_server_msg_received_total.*"
	//   * "grpc_server_msg_sent_total.*"
	//   * "grpc_server_started_total.*"
	EnableHighCardinalityPlatformMetrics bool `json:"enable_high_cardinality_platform_metrics" default:"false"`
}

// BatchPrerollupConfig defines pre-rollup configuration for OTel batch processor.
// +kubebuilder:object:generate=true
//
//swagger:model
type BatchPrerollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"10s"`

	// SendBatchSize is the number of metrics to send in a batch.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"10000"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"10000"`
}

// BatchPostrollupConfig defines post-rollup configuration for OTel batch processor.
// +kubebuilder:object:generate=true
//
//swagger:model
type BatchPostrollupConfig struct {
	// Timeout sets the time after which a batch will be sent regardless of size.
	Timeout config.Duration `json:"timeout" validate:"gt=0" default:"1s"`

	// SendBatchSize is the number of metrics to send in a batch.
	SendBatchSize uint32 `json:"send_batch_size" validate:"gt=0" default:"100"`

	// SendBatchMaxSize is the upper limit of the batch size. Bigger batches will be split
	// into smaller units.
	SendBatchMaxSize uint32 `json:"send_batch_max_size" validate:"gte=0" default:"100"`
}
