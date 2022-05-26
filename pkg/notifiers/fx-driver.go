package notifiers

import (
	"context"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"aperture.tech/aperture/pkg/config"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/status"
)

// FxOptionsFunc is a function that returns fx.Option.
type FxOptionsFunc func(Key, config.Unmarshaller, *status.Registry) (fx.Option, error)

type fxRunner struct {
	UnmarshalKeyNotifier
	app                *fx.App
	fxOptionsFuncs     []FxOptionsFunc
	statusRegistry     *status.Registry
	prometheusRegistry *prometheus.Registry
	registryPath       string
}

// Make sure fxRunner implements KeyNotifier.
var _ KeyNotifier = (*fxRunner)(nil)

// Notify is the main function that notifies the application of the key change.
func (fr *fxRunner) Notify(event Event) {
	fr.UnmarshalKeyNotifier.Notify(event)
	fr.processEvent(event)
}

func (fr *fxRunner) processEvent(event Event) {
	switch event.Type {
	case Write:
		log.Info().Str("event", event.String()).Msg("key update")
		if fr.app != nil {
			// stop existing app
			err := fr.deinitApp()
			if err != nil {
				log.Error().Err(err).Msg("Failed to stop existing app")
			}
		}
		// instantiate and start a new app
		err := fr.initApp(event.Key)
		if err != nil {
			log.Error().Err(err).Msg("Failed to instanticate and start a new app")
		}
	case Remove:
		log.Info().Str("event", event.String()).Msg("key removed")
		// deinit the app
		err := fr.deinitApp()
		if err != nil {
			log.Error().Err(err).Msg("Failed to deinit app")
		}
		fr.statusRegistry.Delete(fr.registryPath)
	}
}

func (fr *fxRunner) initApp(key Key) error {
	s := status.NewStatus(nil, errors.New("policy runner initializing"))
	err := fr.statusRegistry.Push(fr.registryPath, s)
	if err != nil {
		log.Error().Err(err).Str("path", fr.registryPath).Msg("Failed to push status to registry")
	}

	if fr.app == nil && fr.Unmarshaller != nil {
		var options []fx.Option
		for _, fxOptionsFunc := range fr.fxOptionsFuncs {
			o, e := fxOptionsFunc(key, fr.Unmarshaller, fr.statusRegistry)
			if e != nil {
				log.Error().Err(e).Msg("fxOptionsFunc failed")
				return e
			}
			options = append(options, o)
		}
		option := fx.Options(options...)

		fr.app = fx.New(
			// Note: Supplying fr.Unmarshaller directly results in supplying
			// concrete type instead of interface, thus supplying via Provide.
			fx.Provide(func() config.Unmarshaller { return fr.Unmarshaller }),
			// Supply keyinfo
			fx.Supply(key),
			// Supply status registry
			fx.Supply(fr.statusRegistry),
			// Supply prometheus registry
			fx.Supply(fr.prometheusRegistry),
			option,
		)

		if err = fr.app.Err(); err != nil {
			visualize, _ := fx.VisualizeError(err)
			log.Error().Err(err).Str("visualize", visualize).Msg("fx.New failed")
			s := status.NewStatus(nil, err)
			_ = fr.statusRegistry.Push(fr.registryPath, s)
			_ = fr.deinitApp()
			return err
		}

		s := status.NewStatus(nil, errors.New("policy runner starting"))
		err = fr.statusRegistry.Push(fr.registryPath, s)
		if err != nil {
			log.Error().Err(err).Str("path", fr.registryPath).Msg("Failed to push status to registry")
		}

		ctx, cancel := context.WithTimeout(context.Background(), fr.app.StartTimeout())
		defer cancel()
		if err = fr.app.Start(ctx); err != nil {
			log.Error().Str("key", string(key)).Err(err).Msg("Could not start application")
			s = status.NewStatus(nil, err)
			_ = fr.statusRegistry.Push(fr.registryPath, s)
			return err
		}
		s = status.NewStatus(wrapperspb.String("policy runner started"), nil)
		err = fr.statusRegistry.Push(fr.registryPath, s)
		if err != nil {
			log.Error().Err(err).Msg("Failed to push status to registry")
		}
	} else {
		s := status.NewStatus(nil, fr.err)
		err = fr.statusRegistry.Push(fr.registryPath, s)
		if err != nil {
			log.Error().Err(err).Msg("Failed to push status to registry")
		}
	}
	return nil
}

func (fr *fxRunner) deinitApp() error {
	if fr.app != nil {
		ctx, cancel := context.WithTimeout(context.Background(), fr.app.StopTimeout())
		defer func() { fr.app = nil }()
		defer cancel()
		if err := fr.app.Stop(ctx); err != nil {
			log.Error().Err(err).Msg("Could not stop application")
			return err
		}
	}
	return nil
}

// FxDriver tracks prefix and allows spawning "mini FX-based apps" per key in the prefix.
type FxDriver struct {
	// Options for new unmarshaller instances
	UnmarshalPrefixNotifier

	// function to provide fx.Options.
	//
	// Resulting fx.Options will be used to create a "mini FX-based apps" per key.
	// The lifecycle of the app will be tied to the existence of the key.
	// Note that when key's contents change the previous App will be stopped
	// and a fresh one will be created.
	FxOptionsFuncs     []FxOptionsFunc
	StatusRegistry     *status.Registry
	PrometheusRegistry *prometheus.Registry
	// Registry path prefix to push status updates to
	StatusPath string
}

// Make sure FxDriver implements PrefixNotifier.
var _ PrefixNotifier = (*FxDriver)(nil)

// GetKeyNotifier returns a KeyNotifier that will notify the driver of key changes.
func (fxDriver *FxDriver) GetKeyNotifier(key Key) KeyNotifier {
	log.Info().Str("key", key.String()).Msg("GetKeyNotifier")

	statusPath := fmt.Sprintf("%s.%s.driver", fxDriver.StatusPath, key)

	fr := &fxRunner{
		UnmarshalKeyNotifier: fxDriver.getUnmarshalKeyNotifier(key),
		fxOptionsFuncs:       fxDriver.FxOptionsFuncs,
		registryPath:         statusPath,
		statusRegistry:       fxDriver.StatusRegistry,
		prometheusRegistry:   fxDriver.PrometheusRegistry,
	}

	return fr
}
