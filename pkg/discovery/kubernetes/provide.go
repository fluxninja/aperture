package kubernetes

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/discovery/entities"
	kubernetesconfig "github.com/fluxninja/aperture/v2/pkg/discovery/kubernetes/config"
	"github.com/fluxninja/aperture/v2/pkg/k8s"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// Module returns an fx.Option that provides the Kubernetes discovery module.
func Module() fx.Option {
	return fx.Options(
		fx.Invoke(InvokeServiceDiscovery),
	)
}

// FxInSvc describes parameters passed to k8s discovery constructor.
type FxInSvc struct {
	fx.In
	Unmarshaller     config.Unmarshaller
	Lifecycle        fx.Lifecycle
	StatusRegistry   status.Registry
	KubernetesClient k8s.K8sClient
	EntityTrackers   *entities.EntityTrackers
}

// InvokeServiceDiscovery creates a Kubernetes service discovery.
func InvokeServiceDiscovery(in FxInSvc) error {
	var cfg kubernetesconfig.KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(kubernetesconfig.Key, &cfg); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize K8S discovery configuration!")
		return err
	}

	if !cfg.Enabled {
		log.Info().Msg("Skipping Kubernetes discovery since it is disabled")
		return nil
	}
	if in.KubernetesClient == nil {
		// No error, but Genuinely nil, example not in Kubernetes cluster
		log.Info().Msg("Kubernetes client is nil, skipping Kubernetes discovery")
		return nil
	}
	entityEvents := in.EntityTrackers.RegisterServiceDiscovery(podTrackerPrefix)
	ksd, err := newServiceDiscovery(entityEvents, in.KubernetesClient)
	if err != nil {
		log.Info().Err(err).Msg("Failed to create Kubernetes service discovery")
		return err
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ksd.start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return ksd.stop(ctx)
		},
	})

	return nil
}
