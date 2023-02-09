// +kubebuilder:validation:Optional
package config

import (
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
)

// swagger:operation POST /otel controller-configuration OTEL
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ControllerOTELConfig"

// ControllerOTELConfig is the configuration for Agent's OTEL collector.
// swagger:model
// +kubebuilder:object:generate=true
type ControllerOTELConfig struct {
	otelconfig.CommonOTELConfig `json:",inline"`
}
