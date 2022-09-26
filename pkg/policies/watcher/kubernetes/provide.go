package kubernetes

import (
	"context"
	"os"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"go.uber.org/fx"
)

// Constructor holds fields to create an annotated instance of Kubernetes Watcher.
type Constructor struct {
	// Name of watcher instance where "" means main watcher.
	Name string
}

// Module is a fx module that provides Kubernetes watcher.
func Module() fx.Option {
	return fx.Options(
		Constructor{}.Annotate(),
	)
}

// Annotate creates an annotated instance of Kubernetes Watcher.
func (constructor Constructor) Annotate() fx.Option {
	name := ``
	if constructor.Name != "" {
		name = config.NameTag(constructor.Name)
	}
	return fx.Options(fx.Provide(
		fx.Annotate(
			constructor.provideWatcher,
			fx.ResultTags(name),
		),
	))
}

// provideWatcher creates a Kubernetes watcher to watch the Policy Custom Resource.
func (constructor Constructor) provideWatcher(unmarshaller config.Unmarshaller, lifecycle fx.Lifecycle) (notifiers.Watcher, error) {
	if os.Getenv("APERTURE_CONTROLLER_NAMESPACE") == "" {
		os.Setenv("APERTURE_CONTROLLER_NAMESPACE", "default")
	}

	watcher, err := NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Policy Kubernetes watcher")
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := watcher.Start()
			if err != nil {
				log.Error().Err(err).Msg("Failed to start watcher")
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			return watcher.Stop()
		},
	})

	return watcher, nil
}
