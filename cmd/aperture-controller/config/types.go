// +kubebuilder:validation:Optional
package config

import (
	otelconfig "github.com/fluxninja/aperture/v2/pkg/otelcollector/config"
)

// swagger:operation POST /otel controller-configuration OTel
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ControllerOTelConfig"

// ControllerOTelConfig is the configuration for Controller's OTel collector.
// swagger:model
// +kubebuilder:object:generate=true
type ControllerOTelConfig struct {
	otelconfig.CommonOTelConfig `json:",inline"`
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
