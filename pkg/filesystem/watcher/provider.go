package watcher

import (
	"context"
	"errors"

	"github.com/spf13/cast"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
)

// Module is a fx module that provides filesystem watcher.
func Module() fx.Option {
	return fx.Options(
		Constructor{ConfigKey: config.ConfigPathFlag, Path: config.DefaultConfigDirectory}.Annotate(),
	)
}

// Constructor holds fields to create an annotated instance of Watcher.
type Constructor struct {
	// Name of watcher instance where "" means main watcher.
	Name string
	// Config key from which to read directory path from.
	// If the pathkey is empty or is not set then fallback to Path.
	ConfigKey string
	// Directory path to fall back to.
	Path string
	// File types to watch. The extension needs to include leading dot.
	FileExt string
}

// Annotate creates an annotated instance of filesystem Watcher.
func (constructor Constructor) Annotate() fx.Option {
	var name string
	if constructor.Name == "" {
		name = ``
	} else {
		name = config.NameTag(constructor.Name)
	}

	return fx.Options(fx.Provide(
		fx.Annotate(
			constructor.provideWatcher,
			fx.ResultTags(name),
		),
	))
}

func (constructor Constructor) provideWatcher(unmarshaller config.Unmarshaller, lifecycle fx.Lifecycle) (notifiers.Watcher, error) {
	var directory string

	if constructor.ConfigKey != "" {
		// Get directory from config
		directory = cast.ToString(unmarshaller.Get(constructor.ConfigKey))
	}

	if directory == "" {
		if constructor.Path != "" {
			directory = constructor.Path
		} else {
			err := errors.New("no directory provided to watcher")
			log.Error().Err(err).Msg("Unable to create watcher!")
			return nil, err
		}
	}

	fileExt := constructor.FileExt
	if fileExt == "" {
		fileExt = config.DefaultConfigFileExt
	}

	watcher, err := NewWatcher(directory, fileExt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create fs watcher")
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info().Str("directory", watcher.directory).Msg("Starting watcher for directory")
			err := watcher.Start()
			if err != nil {
				log.Error().Err(err).Msg("Failed to start watcher")
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			log.Info().Str("directory", watcher.directory).Msg("Stopping watcher for directory")
			return watcher.Stop()
		},
	})

	return watcher, nil
}
