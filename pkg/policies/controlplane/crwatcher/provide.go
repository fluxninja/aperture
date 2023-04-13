package crwatcher

import (
	"context"
	"os"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"go.uber.org/fx"
)

const (
	// ConfigKey is the key to the Kubernetes watcher config.
	ConfigKey = "policies.cr_watcher"
)

// CRWatcherConfig holds fields to configure the Kubernetes watcher for Aperture Policy custom resource.
// swagger:model
// +kubebuilder:object:generate=true
type CRWatcherConfig struct {
	// Enabled indicates whether the Kubernetes watcher is enabled.
	Enabled bool `json:"enabled" default:"false"`
}

// Constructor holds fields to create an annotated instance of Kubernetes Watcher.
type Constructor struct {
	// Name of tracker instance.
	Name string
	// Name of dynamic config tracker instance.
	DynamicConfigName string
}

// Annotate creates an annotated instance of Kubernetes Watcher.
func (constructor Constructor) Annotate() fx.Option {
	if constructor.Name == "" || constructor.DynamicConfigName == "" {
		log.Panic().Msg("Kubernetes watcher name is required")
	}
	name := config.NameTag(constructor.Name)
	dynamicConfigName := config.NameTag(constructor.DynamicConfigName)
	return fx.Options(fx.Invoke(
		fx.Annotate(
			constructor.provideWatcher,
			fx.ParamTags(name, dynamicConfigName),
		),
	))
}

// provideWatcher creates a Kubernetes watcher to watch the Policy Custom Resource.
func (constructor Constructor) provideWatcher(trackers, dynamicConfigTrackers notifiers.Trackers, unmarshaller config.Unmarshaller, lifecycle fx.Lifecycle) error {
	var config CRWatcherConfig
	err := unmarshaller.UnmarshalKey(ConfigKey, &config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal Kubernetes watcher config")
		return err
	}

	if !config.Enabled {
		log.Info().Msg("Kubernetes watcher is disabled")
		return nil
	}

	if os.Getenv("APERTURE_CONTROLLER_NAMESPACE") == "" {
		os.Setenv("APERTURE_CONTROLLER_NAMESPACE", "default")
	}

	watcher, err := NewWatcher(trackers, dynamicConfigTrackers)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Policy Kubernetes watcher")
		return err
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

	return nil
}
