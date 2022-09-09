// +kubebuilder:validation:Optional
package kubernetes

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/status"
)

var configKey = common.DiscoveryConfigKey + ".kubernetes"

// KubernetesDiscoveryConfig for Kubernetes service discovery.
// swagger:model
// +kubebuilder:object:generate=true
type KubernetesDiscoveryConfig struct {
	// NodeName is the name of the k8s node the agent should be monitoring
	NodeName         string `json:"node_name"`
	PodName          string `json:"pod_name"`
	DiscoveryEnabled bool   `json:"discovery_enabled" default:"true"`
}

// FxIn describes parameters passed to k8s discovery constructor.
type FxIn struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   status.Registry
	KubernetesClient k8s.K8sClient
	EntityTrackers   notifiers.Trackers `name:"entity_trackers"`
}

// InvokeKubernetesServiceDiscovery creates a Kubernetes service discovery.
func InvokeKubernetesServiceDiscovery(in FxIn) error {
	var cfg KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return err
	}

	if cfg.DiscoveryEnabled {
		kd, err := newKubernetesServiceDiscovery(in.EntityTrackers, cfg.NodeName, in.KubernetesClient)
		if err != nil {
			log.Info().Err(err).Msg("Failed to create Kubernetes discovery service")
			return nil
		}

		in.Lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				kd.start()
				return nil
			},
			OnStop: func(_ context.Context) error {
				kd.stop()
				return nil
			},
		})
	} else {
		log.Info().Msg("Skipping Kubernetes discovery service creation")
	}

	return nil
}
