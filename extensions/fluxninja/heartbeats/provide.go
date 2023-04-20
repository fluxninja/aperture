package heartbeats

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	heartbeatv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/fluxninja/v1"
	"github.com/fluxninja/aperture/extensions/fluxninja/extconfig"
	"github.com/fluxninja/aperture/pkg/agentinfo"
	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/discovery/entities"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/etcd/election"
	"github.com/fluxninja/aperture/pkg/jobs"
	"github.com/fluxninja/aperture/pkg/log"
	grpcclient "github.com/fluxninja/aperture/pkg/net/grpc"
	"github.com/fluxninja/aperture/pkg/net/grpcgateway"
	httpclient "github.com/fluxninja/aperture/pkg/net/http"
	"github.com/fluxninja/aperture/pkg/peers"
	autoscalediscovery "github.com/fluxninja/aperture/pkg/policies/autoscale/kubernetes/discovery"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
	"github.com/fluxninja/aperture/pkg/status"
)

// Module returns the module for heartbeats.
func Module() fx.Option {
	log.Info().Msg("Loading Heartbeats extension")
	return fx.Options(
		grpcgateway.RegisterHandler{Handler: heartbeatv1.RegisterControllerInfoServiceHandlerFromEndpoint}.Annotate(),
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", ConfigKey: extconfig.ExtensionConfigKey + ".client.http"}.Annotate(),
		PeersWatcherModule(),
		jobs.JobGroupConstructor{Name: heartbeatsGroup}.Annotate(),
		fx.Provide(provide),
	)
}

// ConstructorIn injects dependencies into the Heartbeats constructor.
type ConstructorIn struct {
	fx.In

	Lifecycle                        fx.Lifecycle
	ExtensionConfig                  *extconfig.FluxNinjaExtensionConfig
	GRPCServer                       *grpc.Server
	JobGroup                         *jobs.JobGroup                     `name:"heartbeats-job-group"`
	GRPClientConnectionBuilder       grpcclient.ClientConnectionBuilder `name:"heartbeats-grpc-client"`
	HTTPClient                       *http.Client                       `name:"heartbeats-http-client"`
	StatusRegistry                   status.Registry
	Entities                         *entities.Entities   `optional:"true"`
	AgentInfo                        *agentinfo.AgentInfo `optional:"true"`
	PeersWatcher                     *peers.PeerDiscovery `name:"fluxninja-peers-watcher" optional:"true"`
	EtcdClient                       *etcdclient.Client
	Election                         *election.Election                          `optional:"true"`
	PolicyFactory                    *controlplane.PolicyFactory                 `optional:"true"`
	FlowControlPoints                *cache.Cache[selectors.TypedControlPointID] `optional:"true"`
	AutoscaleKubernetesControlPoints autoscalediscovery.AutoScaleControlPoints   `optional:"true"`
}

// provide provides a new instance of Heartbeats.
func provide(in ConstructorIn) (*Heartbeats, error) {
	if in.ExtensionConfig.APIKey == "" {
		log.Info().Msg("Heartbeats API key not set, skipping")
		return nil, nil
	}

	heartbeats := newHeartbeats(
		in.JobGroup,
		in.ExtensionConfig,
		in.StatusRegistry,
		in.Entities,
		in.AgentInfo,
		in.PeersWatcher,
		in.PolicyFactory,
		in.Election,
		in.FlowControlPoints,
		in.AutoscaleKubernetesControlPoints,
	)

	heartbeatv1.RegisterControllerInfoServiceServer(in.GRPCServer, heartbeats)

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
