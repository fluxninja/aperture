// +kubebuilder:validation:Optional
package config

import (
	otelconfig "github.com/fluxninja/aperture/pkg/otelcollector/config"
)

//swagger:operation POST /otel controller-configuration OTEL
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/ControllerOTELConfig"

// ControllerOTELConfig is the configuration for Agent's OTEL collector.
// +kubebuilder:object:generate=true
//
//swagger:model
type ControllerOTELConfig struct {
	otelconfig.CommonOTELConfig `json:",inline"`
}
