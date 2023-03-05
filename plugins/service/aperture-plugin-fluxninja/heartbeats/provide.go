package heartbeats

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/discovery/entities"
	"github.com/fluxninja/aperture/pkg/discovery/kubernetes"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/peers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/fluxninja/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

// Module returns the module for heartbeats.
func Module() fx.Option {
	log.Info().Msg("Loading Heartbeats plugin")
	return fx.Options(
		RegisterControllerInfoServiceHTTP(),
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", ConfigKey: pluginconfig.PluginConfigKey + ".client.grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", ConfigKey: pluginconfig.PluginConfigKey + ".client.http"}.Annotate(),
		fx.Provide(Provide),
		PeersWatcherModule(),
		jobs.JobGroupConstructor{Name: heartbeatsGroup}.Annotate(),
		fx.Invoke(
			Invoke,
			RegisterControllerInfoService,
		),
	)
}

// ConstructorIn injects dependencies into the Heartbeats constructor.
type ConstructorIn struct {
	fx.In

	Lifecycle                        fx.Lifecycle
	Unmarshaller                     config.Unmarshaller
	JobGroup                         *jobs.JobGroup                     `name:"heartbeats-job-group"`
	GRPClientConnectionBuilder       grpcclient.ClientConnectionBuilder `name:"heartbeats-grpc-client"`
	HTTPClient                       *http.Client                       `name:"heartbeats-http-client"`
	StatusRegistry                   status.Registry
	Entities                         *entities.Entities   `optional:"true"`
	AgentInfo                        *agentinfo.AgentInfo `optional:"true"`
	PeersWatcher                     *peers.PeerDiscovery `name:"fluxninja-peers-watcher" optional:"true"`
	EtcdClient                       *etcdclient.Client
	PolicyFactory                    *controlplane.PolicyFactory            `optional:"true"`
	FlowControlControlPoints         *cache.Cache[selectors.ControlPointID] `optional:"true"`
	AutoscaleKubernetesControlPoints kubernetes.AutoscaleControlPoints      `optional:"true"`
}

// Provide provides a new instance of Heartbeats.
func Provide(in ConstructorIn) (*Heartbeats, error) {
	var config pluginconfig.FluxNinjaPluginConfig
	if err := in.Unmarshaller.UnmarshalKey(pluginconfig.PluginConfigKey, &config); err != nil {
		return nil, err
	}

	var discoveryConfig kubernetes.KubernetesDiscoveryConfig
	if err := in.Unmarshaller.UnmarshalKey(kubernetes.ConfigKey, &discoveryConfig); err != nil {
		return nil, err
	}

	installationMode := getInstallationMode()

	heartbeats := newHeartbeats(
		in.JobGroup,
		config,
		in.StatusRegistry,
		in.Entities,
		in.AgentInfo,
		in.PeersWatcher,
		in.PolicyFactory,
		in.FlowControlControlPoints,
		in.AutoscaleKubernetesControlPoints,
		installationMode,
	)

	runCtx, cancel := context.WithCancel(context.Background())

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := heartbeats.setupControllerInfo(runCtx, in.EtcdClient)
			if err != nil {
				log.Error().Err(err).Msg("Could not read/create controller id in heartbeats")
				return err
			}

			err = heartbeats.start(runCtx, &in)
			if err != nil {
				log.Error().Err(err).Msg("Heartbeats start had an error")
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			heartbeats.stop()
			return nil
		},
	})

	return heartbeats, nil
}

func getInstallationMode() string {
	_, kubernetesServiceHostExists := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	_, kubernetesServicePortExists := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	if kubernetesServiceHostExists && kubernetesServicePortExists {
		// Check if the code is running in a sidecar container
		if os.Getenv("SIDECAR") != "" {
			return "KUBERNETES_SIDECAR"
		} else {
			return "KUBERNETES_DAEMONSET"
		}
	}
	// Check if running on bare metal Linux
	if _, err := os.Stat("/proc/cpuinfo"); err == nil {
		return "BAREMETAL_LINUX"
	}
	return "UNKNOWN"
}

// Invoke enables heartbeats in FX.
func Invoke(*Heartbeats) {}
