package heartbeats

import (
	"context"
	"net/http"

	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/pkg/agentinfo"
	"github.com/FluxNinja/aperture/pkg/config"
	"github.com/FluxNinja/aperture/pkg/entitycache"
	"github.com/FluxNinja/aperture/pkg/jobs"
	"github.com/FluxNinja/aperture/pkg/log"
	grpcclient "github.com/FluxNinja/aperture/pkg/net/grpc"
	httpclient "github.com/FluxNinja/aperture/pkg/net/http"
	"github.com/FluxNinja/aperture/pkg/peers"
	"github.com/FluxNinja/aperture/pkg/status"
	"github.com/FluxNinja/aperture/plugins/service/aperture-plugin-fluxninja/pluginconfig"
)

// Module returns the module for heartbeats.
func Module() fx.Option {
	log.Info().Msg("Loading Heartbeats plugin")
	return fx.Options(
		grpcclient.ClientConstructor{Name: "heartbeats-grpc-client", Key: pluginconfig.PluginConfigKey + ".client_grpc"}.Annotate(),
		httpclient.ClientConstructor{Name: "heartbeats-http-client", Key: pluginconfig.PluginConfigKey + ".client_http"}.Annotate(),
		fx.Provide(Provide),
		PeersWatcherModule(),
		jobs.JobGroupConstructor{Group: heartbeatsGroup}.Annotate(),
		fx.Invoke(Invoke),
	)
}

// ConstructorIn injects dependencies into the Heartbeats constructor.
type ConstructorIn struct {
	fx.In

	Lifecycle                  fx.Lifecycle
	Unmarshaller               config.Unmarshaller
	JobGroup                   *jobs.JobGroup                     `name:"heartbeats-job-group"`
	GRPClientConnectionBuilder grpcclient.ClientConnectionBuilder `name:"heartbeats-grpc-client"`
	HTTPClient                 *http.Client                       `name:"heartbeats-http-client"`
	StatusRegistry             *status.Registry
	EntityCache                *entitycache.EntityCache `optional:"true"`
	AgentInfo                  *agentinfo.AgentInfo     `optional:"true"`
	PeersWatcher               *peers.PeerDiscovery     `name:"fluxninja-peers-watcher" optional:"true"`
}

// Provide provides a new instance of Heartbeats.
func Provide(in ConstructorIn) (*heartbeats, error) {
	var config pluginconfig.FluxNinjaPluginConfig
	if err := in.Unmarshaller.UnmarshalKey(pluginconfig.PluginConfigKey, &config); err != nil {
		return nil, err
	}

	heartbeats := newHeartbeats(in.JobGroup, config, in.StatusRegistry, in.EntityCache, in.AgentInfo, in.PeersWatcher)
	runCtx, cancel := context.WithCancel(context.Background())

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := heartbeats.start(runCtx, &in)
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

// Invoke enables heartbeats in FX.
func Invoke(*heartbeats) {}
