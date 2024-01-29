package platform

import (
	"context"
	"errors"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/status"
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
			platform.statusRegistry.Child("system", readinessStatusPath).
				Child("component", platformStatusPath).
				SetStatus(status.NewStatus(nil, errors.New("platform starting")), nil)
			return nil
		},
		OnStop: func(context.Context) error {
			platform.statusRegistry.Child("system", readinessStatusPath).
				Child("component", platformStatusPath).
				SetStatus(status.NewStatus(nil, errors.New("platform stopped")), nil)
			return nil
		},
	})

	return nil
}
