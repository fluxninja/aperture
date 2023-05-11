package notifiers

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Watcher is a generic interface for watchers/trackers.
type Watcher interface {
	AddPrefixNotifier(PrefixNotifier) error
	RemovePrefixNotifier(PrefixNotifier) error
	AddKeyNotifier(KeyNotifier) error
	RemoveKeyNotifier(KeyNotifier) error
	Start() error
	Stop() error
}

// WatcherLifecycle starts/stops watcher and adds/removes prefix notifier(s) to etcd watcher.
func WatcherLifecycle(lc fx.Lifecycle, watcher Watcher, notifiers []PrefixNotifier) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := watcher.Start()
			if err != nil {
				return err
			}

			for _, notifier := range notifiers {
				err := watcher.AddPrefixNotifier(notifier)
				if err != nil {
					log.Error().Err(err).Msg("Failed to add notifier")
					return err
				}
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			var errMulti error
			for _, notifier := range notifiers {
				err := watcher.RemovePrefixNotifier(notifier)
				if err != nil {
					errMulti = multierror.Append(errMulti, err)
					log.Error().Err(err).Msg("Failed to remove notifier")
				}
			}

			err := watcher.Stop()
			if err != nil {
				errMulti = multierror.Append(errMulti, err)
			}
			return errMulti
		},
	})
}

// NotifierLifecycle adds/removed prefix notifier to etcd watcher.
func NotifierLifecycle(lc fx.Lifecycle, watcher Watcher, notifier PrefixNotifier) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := watcher.AddPrefixNotifier(notifier)
			if err != nil {
				log.Error().Err(err).Msg("Failed to add notifier")
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			err := watcher.RemovePrefixNotifier(notifier)
			if err != nil {
				log.Error().Err(err).Msg("Failed to remove notifier")
				return err
			}
			return nil
		},
	})
}
