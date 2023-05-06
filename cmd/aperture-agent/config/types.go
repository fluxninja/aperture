// +kubebuilder:validation:Optional
package config

import (
	"github.com/fluxninja/aperture/pkg/config"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
	otelcustom "github.com/fluxninja/aperture/pkg/otelcollector/custom"
)

// swagger:operation POST /otel agent-configuration OTel
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/AgentOTelConfig"

// AgentOTelConfig is the configuration for Agent's OTel collector.
//
// Example configuration:
//
// ```yaml
//
//	otel:
//	  batch_alerts:
//	    send_batch_max_size: 100
//	    send_batch_size: 100
//	    timeout: 1s
//	  batch_prerollup:
//	    send_batch_max_size: 10000
//	    send_batch_size: 10000
//	    timeout: 10s
//	  batch_postrollup:
//	    send_batch_max_size: 100
//	    send_batch_size: 100
//	    timeout: 1s
//	  custom_metrics:
//	    rabbitmq:
//	      processors:
//	        batch:
//	          send_batch_size: 10
//	          timeout: 10s
//	      receivers:
//	        rabbitmq:
//	          collection_interval: 10s
//	          endpoint: http://<rabbitmq-svc-fqdn>:15672
//	          password: secretpassword
//	          username: admin
//	      per_agent_group: true
//
// ```
//
// +kubebuilder:object:generate=true
//
//swagger:model
type AgentOTelConfig struct {
	// CustomMetrics configures custom metrics OTel pipelines, which will send data to
	// the controller Prometheus.
	// Key in this map refers to OTel pipeline name. Prefixing pipeline name with `metrics/`
	// is optional, as all the components and pipeline names would be normalized.
	// By default `kubeletstats` custom metrics is added, which can be overwritten.
	//
	// Below is example to overwrite `kubeletstats` custom metrics:
	//
	// ```yaml
	//  otel:
	//    custom_metrics:
	//      kubeletstats: {}
	// ```
	//
	CustomMetrics map[string]otelcustom.CustomMetricsConfig `json:"custom_metrics,omitempty"`
	// BatchPrerollup configures the OTel batch pre-processor.
	BatchPrerollup BatchPrerollupConfig `json:"batch_prerollup"`
	// BatchPostrollup configures the OTel batch post-processor.
	BatchPostrollup             BatchPostrollupConfig `json:"batch_postrollup"`
	otelconfig.CommonOTelConfig `json:",inline"`
	// DisableKubernetesScraper disables metrics collection for Kubernetes resources.
	DisableKubernetesScraper bool `json:"disable_kubernetes_scraper" default:"false"`
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
