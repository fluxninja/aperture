// +kubebuilder:validation:Optional
package config

import (
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/fluxninja/aperture/pkg/config"
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
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
//		batch_alerts:
//			send_batch_max_size: 100
//			send_batch_size: 100
//			timeout: 1s
//		batch_prerollup:
//			send_batch_max_size: 10000
//			send_batch_size: 10000
//			timeout: 10s
//		batch_postrollup:
//			send_batch_max_size: 100
//			send_batch_size: 100
//			timeout: 1s
//		custom_metrics:
//			rabbitmq:
//				processors:
//					batch:
//						send_batch_size: 10
//		 				timeout: 10s
//				receivers:
//		 			rabbitmq:
//		 				collection_interval: 10s
//						endpoint: http://<rabbitmq-svc-fqdn>:15672
//						password: secretpassword
//						username: admin
//				per_agent_group: true
//
// ```
//
// +kubebuilder:object:generate=true
//
//swagger:model
type AgentOTelConfig struct {
	otelconfig.CommonOTelConfig `json:",inline"`
	// DisableKubernetesScraper disables metrics collection for Kubernetes resources.
	DisableKubernetesScraper bool `json:"disable_kubernetes_scraper" default:"false"`
	// BatchPrerollup configures batch prerollup processor.
	BatchPrerollup BatchPrerollupConfig `json:"batch_prerollup"`
	// BatchPostrollup configures batch postrollup processor.
	BatchPostrollup BatchPostrollupConfig `json:"batch_postrollup"`
	// CustomMetrics configures custom metrics OTel pipelines, which will send data to
	// the controller Prometheus.
	// Key in this map refers to OTel pipeline name. Prefixing pipeline name with `metrics/`
	// is optional, as all the components and pipeline names would be normalized.
	// By default `kubeletstats` custom metrics is added, which can be overwritten.
	//
	// Below is example to overwrite `kubeletstats` custom metrics:
	//
	//	otel:
	//		custom_metrics:
	//			kubeletstats: {}
	//
	CustomMetrics map[string]CustomMetricsConfig `json:"custom_metrics,omitempty"`
}

// BatchPrerol[.*?]upConfig defines configuration for OTel batch processor.
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

// CustomMetricsConfig defines receivers, processors, and single metrics pipeline which will be exported to the controller Prometheus.
// Environment variables can be used in the configuration using format `${ENV_VAR_NAME}`.
// +kubebuilder:object:generate=true
//
// :::info
// See also [Get Started / Setup Integrations / Metrics](/get-started/integrations/metrics/metrics.md).
// :::
//
//swagger:model
type CustomMetricsConfig struct {
	// Receivers define receivers to be used in custom metrics pipelines. This should
	// be in OTel format - https://opentelemetry.io/docs/collector/configuration/#receivers.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Receivers Components `json:"receivers"`
	// Processors define processors to be used in custom metrics pipelines. This should
	// be in OTel format - https://opentelemetry.io/docs/collector/configuration/#processors.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Processors Components `json:"processors,omitempty"`
	// Pipeline is an OTel metrics pipeline definition, which **only** uses receivers
	// and processors defined above. Exporter would be added automatically.
	//
	// If there are no processors defined or only one processor is defined, the
	// pipeline definition can be omitted. In such cases, the pipeline will
	// automatically use all given receivers and the defined processor (if
	// any).  However, if there are more than one processor, the pipeline must
	// be defined explicitly.
	Pipeline CustomMetricsPipelineConfig `json:"pipeline"`
	// PerAgentGroup marks the pipeline to be instantiated only once per agent
	// group. This is helpful for receivers that scrape for example, some cluster-wide
	// metrics. When not set, pipeline will be instantiated on every Agent.
	PerAgentGroup bool `json:"per_agent_group"`
}

// Components is an alias type for map[string]any. This needs to be used
// because of the CRD requirements for the operator.
// https://github.com/kubernetes-sigs/controller-tools/issues/636
// https://github.com/kubernetes-sigs/kubebuilder/issues/528
// +kubebuilder:object:generate=false
// +kubebuilder:pruning:PreserveUnknownFields
// +kubebuilder:validation:Schemaless
type Components map[string]any

// DeepCopyInto is an deepcopy function, copying the receiver, writing into out.
// In must be non-nil.
// We need to specify this manyually, as the generator does not support `any`.
func (in *Components) DeepCopyInto(out *Components) {
	if in == nil {
		*out = nil
	} else {
		*out = runtime.DeepCopyJSON(*in)
	}
}

// DeepCopy is an deepcopy function, copying the receiver, creating a new
// Components.
// We need to specify this manyually, as the generator does not support `any`.
func (in *Components) DeepCopy() *Components {
	if in == nil {
		return nil
	}
	out := new(Components)
	in.DeepCopyInto(out)
	return out
}

// CustomMetricsPipelineConfig defines a custom metrics pipeline.
// +kubebuilder:object:generate=true
//
//swagger:model
type CustomMetricsPipelineConfig struct {
	Receivers  []string `json:"receivers,omitempty"`
	Processors []string `json:"processors,omitempty"`
}
