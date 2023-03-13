package heartbeats

import (
	"context"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/extensions/fluxninja/extconfig"
	etcdclient "github.com/fluxninja/aperture/pkg/etcd/client"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/peers"
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
	client *etcdclient.Client,
	lc fx.Lifecycle,
) (PeersOut, error) {
	if extensionConfig.APIKey == "" {
		return PeersOut{}, nil
	}

	if info.Service != "aperture-controller" {
		return PeersOut{}, nil
	}
	pd, err := peers.NewPeerDiscovery("aperture-agent", client, nil)
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
