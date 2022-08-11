package platform

import (
	"context"
	"errors"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/status"
)

const platformReadinessStatusName = "readiness.platform"

func platformStatusModule() fx.Option {
	return fx.Options(
		fx.Invoke(providePlatformReadinessStatus),
	)
}

type platformReadinessStatusIn struct {
	fx.In

	Lifecycle      fx.Lifecycle
	StatusRegistry *status.Registry
}

func providePlatformReadinessStatus(in platformReadinessStatusIn) error {
	s := status.NewStatus(nil, nil)
	err := in.StatusRegistry.Push(platformReadinessStatusName, s)
	if err != nil {
		return err
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s := status.NewStatus(nil, errors.New("platform starting"))
			err := in.StatusRegistry.Push(platformReadinessStatusName, s)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			s := status.NewStatus(nil, errors.New("platform stopping"))
			err := in.StatusRegistry.Push(platformReadinessStatusName, s)
			if err != nil {
				return err
			}
			return nil
		},
	})

	return nil
}
