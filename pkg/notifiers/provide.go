package notifiers

import (
	"context"

	"go.uber.org/fx"

	"github.com/FluxNinja/aperture/pkg/config"
)

// TrackersConstructor is a Fx constructor for Trackers.
type TrackersConstructor struct {
	Name string
}

// Annotate provides Fx annotated Tracker.
func (t TrackersConstructor) Annotate() fx.Option {
	var name string
	if t.Name == "" {
		name = ``
	} else {
		name = config.NameTag(t.Name)
	}

	return fx.Provide(
		fx.Annotate(
			t.provideTrackers,
			fx.ResultTags(name),
		),
	)
}

func (t TrackersConstructor) provideTrackers(lifecycle fx.Lifecycle) (Trackers, error) {
	trackers := NewDefaultTrackers()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return trackers.Start()
		},
		OnStop: func(context.Context) error {
			return trackers.Stop()
		},
	})

	return trackers, nil
}
