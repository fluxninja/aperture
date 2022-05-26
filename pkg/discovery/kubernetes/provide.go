package kubernetes

import (
	"context"

	"go.uber.org/fx"

	"aperture.tech/aperture/pkg/config"
	"aperture.tech/aperture/pkg/discovery/common"
	"aperture.tech/aperture/pkg/k8s"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/notifiers"
	"aperture.tech/aperture/pkg/status"
)

var configKey = common.DiscoveryConfigKey + ".kubernetes"

// KubernetesDiscoveryConfig for Kubernetes service discovery.
// swagger:model
type KubernetesDiscoveryConfig struct {
	// NodeName is the name of the k8s node the agent should be monitoring
	NodeName    string `json:"node_name"`
	PodName     string `json:"pod_name"`
	SidecarMode bool   `json:"sidecar_mode"`
}

// FxIn describes parameters passed to k8s discovery constructor.
type FxIn struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   *status.Registry
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

	if !cfg.SidecarMode {
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
		log.Info().Msg("Operating in sidecar mode. Skipping Kubernetes discovery service creation")
	}

	return nil
}
