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

	Lifecycle      fx.Lifecycle
	StatusRegistry status.Registry
}

func platformReadinessStatus(in platformReadinessStatusIn) error {
	readinessStatusRegistry := status.NewRegistry(in.StatusRegistry, readinessStatusPath)
	platformStatusRegistry := status.NewRegistry(readinessStatusRegistry, platformStatusPath)

	s := status.NewStatus(nil, nil)
	err := platformStatusRegistry.Push(s)
	if err != nil {
		return err
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s := status.NewStatus(nil, errors.New("platform starting"))
			err := platformStatusRegistry.Push(s)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			s := status.NewStatus(nil, errors.New("platform stopping"))
			err := platformStatusRegistry.Push(s)
			if err != nil {
				return err
			}
			return nil
		},
	})

	return nil
}
