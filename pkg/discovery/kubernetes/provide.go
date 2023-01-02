// +kubebuilder:validation:Optional
package kubernetes

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/etcd/election"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
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

// Module returns an fx.Option that provides the Kubernetes discovery module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideKubernetesControlPointsCache,
		),
		fx.Invoke(
			InvokeKubernetesServiceDiscovery,
			InvokeKubernetesControlPointDiscovery,
		),
	)
}

// ProvideKubernetesControlPointsCache creates a Kubernetes control points cache.
func ProvideKubernetesControlPointsCache() (*ControlPointCache, error) {
	return newControlPointCache(), nil
}

// FxInCtrlPt is the input for the ProvideKuberetesControlPointsCache function.
type FxInCtrlPt struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   status.Registry
	KubernetesClient k8s.K8sClient
	Election         *election.Election
}

// InvokeKubernetesControlPointDiscovery creates a Kubernetes control point discovery.
func InvokeKubernetesControlPointDiscovery(in FxInCtrlPt) error {
	var cfg KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return err
	}

	if cfg.DiscoveryEnabled {
		kcpd, err := newControlPointDiscovery(in.Election, in.KubernetesClient)
		if err != nil {
			log.Info().Err(err).Msg("Failed to create Kubernetes control point discovery")
			return nil
		}

		in.Lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				kcpd.start()
				return nil
			},
			OnStop: func(_ context.Context) error {
				kcpd.stop()
				return nil
			},
		})
	} else {
		log.Info().Msg("Skipping Kubernetes discovery service creation")
	}

	return nil
}

// FxInSvc describes parameters passed to k8s discovery constructor.
type FxInSvc struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   status.Registry
	KubernetesClient k8s.K8sClient
	EntityTrackers   *entitycache.EntityTrackers
}

// InvokeKubernetesServiceDiscovery creates a Kubernetes service discovery.
func InvokeKubernetesServiceDiscovery(in FxInSvc) error {
	var cfg KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return err
	}

	if cfg.DiscoveryEnabled {
		entityEvents := in.EntityTrackers.RegisterServiceDiscovery(podTrackerPrefix)
		ksd, err := newServiceDiscovery(entityEvents, cfg.NodeName, in.KubernetesClient)
		if err != nil {
			log.Info().Err(err).Msg("Failed to create Kubernetes service discovery")
			return nil
		}

		in.Lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				ksd.start()
				return nil
			},
			OnStop: func(_ context.Context) error {
				ksd.stop()
				return nil
			},
		})
	} else {
		log.Info().Msg("Skipping Kubernetes discovery service creation")
	}

	return nil
}
