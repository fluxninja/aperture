package config

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"go.uber.org/fx"
)

// Module returns the fx options for autoscale config.
func Module() fx.Option {
	return fx.Provide(provideConfig)
}

// ConfigKey is config path where FlowPreviewConfig is located.
const ConfigKey = "auto_scale.kubernetes"

// AutoScaleKubernetesConfig is the configuration for the flow preview service.
// swagger:model
// +kubebuilder:object:generate=true
type AutoScaleKubernetesConfig struct {
	// Enables the Kubernetes auto scale capability.
	Enabled bool `json:"enabled" default:"true"`
}

func provideConfig(unmarshaller config.Unmarshaller) (AutoScaleKubernetesConfig, error) {
	var cfg AutoScaleKubernetesConfig
	if err := unmarshaller.UnmarshalKey(ConfigKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return cfg, err
	}
	return cfg, nil
}
