package watcher

import (
	"context"
	"io/fs"
	"os"
	"path"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/config"
	etcdclient "github.com/fluxninja/aperture/v2/pkg/etcd/client"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
)

// Constructor holds fields to create an annotated instance of etcd Watcher.
type Constructor struct {
	Name     string
	EtcdPath string
}

// Annotate creates an annotated instance of etcd Watcher.
func (c Constructor) Annotate() fx.Option {
	name := ``
	if c.Name != "" {
		name = config.NameTag(c.Name)
	}

	return fx.Options(
		fx.Provide(
			fx.Annotate(
				c.provideWatcher,
				fx.ResultTags(name),
			),
		),
	)
}

func (c Constructor) provideWatcher(client *etcdclient.Client, unmarshaller config.Unmarshaller, lifecycle fx.Lifecycle) (notifiers.Watcher, error) {
	watcher, err := NewWatcher(client, c.EtcdPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create etcd watcher")
		return nil, err
	}

	// make test notifier directory and create a new PrefixToFS notifier on that directory
	testNotifierPath := path.Join(config.DefaultTempDirectory, c.EtcdPath)
	err = os.MkdirAll(testNotifierPath, fs.ModePerm)
	if err != nil {
		log.Error().Err(err).Str("dir", testNotifierPath).Msg("Unable to create test notifier directory")
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := watcher.Start()
			if err != nil {
				log.Error().Err(err).Msg("Failed to start etcd watcher")
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			return watcher.Stop()
		},
	})

	return watcher, nil
}
