package heartbeats

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/extensions/fluxninja/extconfig"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/peers"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

// PeersWatcherModule is a fx module that watches all agent peers.
func PeersWatcherModule() fx.Option {
	return fx.Options(
		// Etcd peers watcher
		fx.Provide(
			setupPeersWatcher,
		),
	)
}

// PeersOut is a return struct provided to fx.
type PeersOut struct {
	fx.Out
	PeerWatcher *peers.PeerDiscovery `name:"fluxninja-peers-watcher" optional:"true"`
}

func setupPeersWatcher(
	extensionConfig *extconfig.FluxNinjaExtensionConfig,
	etcdClient *etcdclient.Client,
	lc fx.Lifecycle,
) (PeersOut, error) {
	apiKey := extensionConfig.APIKey
	if apiKey == "" {
		//nolint:staticcheck // SA1019 read AgentAPIKey config for backward compatibility
		apiKey = extensionConfig.AgentAPIKey
	}
	if apiKey == "" {
		return PeersOut{}, nil
	}

	if info.Service != utils.ApertureController {
		return PeersOut{}, nil
	}
	pd, err := peers.NewPeerDiscovery("aperture-agent", etcdClient, nil)
	if err != nil {
		return PeersOut{}, err
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := pd.Start()
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			return pd.Stop()
		},
	})
	return PeersOut{PeerWatcher: pd}, nil
}
