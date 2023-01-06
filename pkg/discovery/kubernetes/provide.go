// +kubebuilder:validation:Optional
package kubernetes

import (
	"context"

	"go.uber.org/fx"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/common"
	"github.com/fluxninja/aperture/pkg/entitycache"
	"github.com/fluxninja/aperture/pkg/etcd/election"
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
	AutoscaleEnabled bool   `json:"autoscale_enabled" default:"true"`
}

// Module returns an fx.Option that provides the Kubernetes discovery module.
func Module() fx.Option {
	return fx.Options(
		notifiers.TrackersConstructor{Name: "kubernetes_control_points"}.Annotate(),
		fx.Invoke(
			ProvideAutoscaler,
			InvokeServiceDiscovery,
		),
	)
}

// FxInAutoScaler is the input for the ProvideKuberetesControlPointsCache function.
type FxInAutoScaler struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   status.Registry
	KubernetesClient k8s.K8sClient
	Election         *election.Election
	Trackers         notifiers.Trackers `name:"kubernetes_control_points"`
}

// ProvideAutoscaler provides Kubernetes AutoScaler and starts Kubernetes control point discovery if enabled.
func ProvideAutoscaler(in FxInAutoScaler) (AutoScaler, error) {
	var cfg KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(configKey, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return nil, err
	}

	discoveryClient := in.KubernetesClient.GetClientSet().DiscoveryClient
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	cachedDiscoveryClient := cacheddiscovery.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(discoveryClient)
	scaleClient, err := scale.NewForConfig(config, mapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create scale client")
		return nil, err
	}

	autoScaler := newAutoScaler(scaleClient, in.Trackers)

	if cfg.DiscoveryEnabled {
		cpd, err := newControlPointDiscovery(in.Election, in.KubernetesClient, discoveryClient, dynClient, autoScaler)
		if err != nil {
			log.Info().Err(err).Msg("Failed to create Kubernetes control point discovery")
			return nil, err
		}

		in.Lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				cpd.start()
				return nil
			},
			OnStop: func(_ context.Context) error {
				cpd.stop()
				return nil
			},
		})
	} else {
		log.Info().Msg("Skipping Kubernetes discovery service creation")
	}

	return autoScaler, nil
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
