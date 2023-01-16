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
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/status"
)

var (
	configKey = common.DiscoveryConfigKey + ".kubernetes"
	// FxTagBase is the tag's base used to identify the Kubernetes Control Points Tracker.
	FxTagBase = "kubernetes_control_points"
	// FxTag is the tag used to identify the Kubernetes Control Points Tracker.
	FxTag = config.NameTag(FxTagBase)
)

// KubernetesDiscoveryConfig for Kubernetes service discovery.
// swagger:model
// +kubebuilder:object:generate=true
type KubernetesDiscoveryConfig struct {
	// NodeName is the name of the k8s node the agent should be monitoring
	NodeName         string `json:"node_name"`
	PodName          string `json:"pod_name"`
	DiscoveryEnabled bool   `json:"discovery_enabled" default:"true"`
	AutoscaleEnabled bool   `json:"autoscale_enabled" default:"true"`
}

// Module returns an fx.Option that provides the Kubernetes discovery module.
func Module() fx.Option {
	return fx.Options(
		notifiers.TrackersConstructor{Name: FxTagBase}.Annotate(),
		fx.Provide(
			ProvideControlPointCache,
		),
		fx.Invoke(
			InvokeServiceDiscovery,
		),
	)
}

// FxInK8sScale is the input for the ProvideKuberetesControlPointsCache function.
type FxInK8sScale struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   status.Registry
	KubernetesClient k8s.K8sClient `optional:"true"`
	Election         *election.Election
	Trackers         notifiers.Trackers `name:"kubernetes_control_points"`
}

// ProvideControlPointCache provides Kubernetes AutoScaler and starts Kubernetes control point discovery if enabled.
func ProvideControlPointCache(in FxInK8sScale) (ControlPointCache, error) {
	var cfg KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return nil, err
	}

	controlPointCache := newControlPointCache(in.Trackers, in.KubernetesClient)

	if !cfg.AutoscaleEnabled {
		log.Info().Msg("Skipping Kubernetes Control Point Discovery since Autoscale is disabled")
		return controlPointCache, nil
	}
	if in.KubernetesClient == nil {
		log.Error().Msg("Kubernetes client is not available, skipping Kubernetes Control Point Discovery")
		return controlPointCache, nil
	}
	cpd, err := newControlPointDiscovery(in.Election, in.KubernetesClient, controlPointCache)
	if err != nil {
		log.Info().Err(err).Msg("Failed to create Kubernetes Control Point Discovery")
		return nil, err
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			controlPointCache.start()
			cpd.start()
			return nil
		},
		OnStop: func(_ context.Context) error {
			cpd.stop()
			controlPointCache.stop()
			return nil
		},
	})

	return controlPointCache, nil
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

// InvokeServiceDiscovery creates a Kubernetes service discovery.
func InvokeServiceDiscovery(in FxInSvc) error {
	var cfg KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return err
	}

	if !cfg.DiscoveryEnabled {
		log.Info().Msg("Skipping Kubernetes discovery since it is disabled")
		return nil
	}
	if in.KubernetesClient == nil {
		// No error, but Genuinely nil, example not in Kubernetes cluster
		log.Info().Msg("Kubernetes client is nil, skipping Kubernetes discovery")
		return nil
	}
	entityEvents := in.EntityTrackers.RegisterServiceDiscovery(podTrackerPrefix)
	ksd, err := newServiceDiscovery(entityEvents, cfg.NodeName, in.KubernetesClient)
	if err != nil {
		log.Info().Err(err).Msg("Failed to create Kubernetes service discovery")
		return err
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

	return nil
}
