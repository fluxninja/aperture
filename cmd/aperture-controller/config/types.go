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
}
