package discovery

import (
	"fmt"

	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/k8s"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	autoscalek8sconfig "github.com/fluxninja/aperture/v2/pkg/policies/autoscale/kubernetes/config"
	"github.com/fluxninja/aperture/v2/pkg/status"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

var (
	// FxTagBase is the tag's base used to identify the Kubernetes Control Points Tracker.
	FxTagBase = "kubernetes_control_points"
	// FxTag is the tag used to identify the Kubernetes Control Points Tracker.
	FxTag = config.NameTag(FxTagBase)
)

// Module returns the fx options for Kubernetes Control Point Discovery.
func Module() fx.Option {
	return fx.Options(
		notifiers.TrackersConstructor{Name: FxTagBase}.Annotate(),
		fx.Provide(provideAutoScaleControlPoints),
	)
}

// FxIn is the input for the ProvideKuberetesControlPointsCache function.
type FxIn struct {
	fx.In
	Unmarshaller       config.Unmarshaller
	StatusRegistry     status.Registry
	KubernetesClient   k8s.K8sClient      `optional:"true"`
	Trackers           notifiers.Trackers `name:"kubernetes_control_points"`
	EtcdClient         *etcdclient.Client
	Config             autoscalek8sconfig.AutoScaleKubernetesConfig
	PrometheusRegistry *prometheus.Registry
	AgentInfo          *agentinfo.AgentInfo
}

// provideAutoScaleControlPoints provides Kubernetes AutoScaler and starts Kubernetes control point discovery if enabled.
func provideAutoScaleControlPoints(in FxIn) (AutoScaleControlPoints, error) {
	if in.KubernetesClient == nil {
		log.Error().Msg("Kubernetes client is not available, skipping Kubernetes AutoScaler creation and control point discovery")
		return nil, nil
	}
	controlPointCache, err := newAutoScaleControlPoints(in.Trackers, in.EtcdClient, in.KubernetesClient)
	if err != nil {
		return nil, fmt.Errorf("could not create auto scale control points: %w", err)
	}

	if !in.Config.Enabled {
		log.Info().Msg("Skipping Kubernetes Control Point Discovery since AutoScale is disabled")
		return controlPointCache, nil
	}
	_, err = newControlPointDiscovery(in.EtcdClient, in.KubernetesClient, controlPointCache)
	if err != nil {
		log.Info().Err(err).Msg("Failed to create Kubernetes Control Point Discovery")
		return nil, err
	}

	return controlPointCache, nil
}
