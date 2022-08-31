package platform

import (
	"context"
	"errors"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/status"
)

const (
	readinessStatusPath = "readiness"
	platformStatusPath  = "platform"
)

func platformStatusModule() fx.Option {
	return fx.Options(
		fx.Invoke(platformReadinessStatus),
	)
}

type platformReadinessStatusIn struct {
	fx.In
	Lifecycle fx.Lifecycle
}

func platformReadinessStatus(in platformReadinessStatusIn) error {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			platform.statusRegistry.SetStatus(status.NewStatus(nil, errors.New("platform starting")))
			return nil
		},
		OnStop: func(context.Context) error {
			platform.statusRegistry.SetStatus(status.NewStatus(nil, errors.New("platform stopped")))
			return nil
		},
	})

	return nil
}
