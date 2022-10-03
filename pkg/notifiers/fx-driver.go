package notifiers

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/status"
)

// FxOptionsFunc is a function that returns fx.Option.
type FxOptionsFunc func(Key, config.Unmarshaller, status.Registry) (fx.Option, error)

type fxRunner struct {
	statusRegistry         status.Registry
	fxRunnerStatusRegistry status.Registry
	app                    *fx.App
	prometheusRegistry     *prometheus.Registry
	UnmarshalKeyNotifier
	fxOptionsFuncs []FxOptionsFunc
}

// Make sure fxRunner implements KeyNotifier.
var _ KeyNotifier = (*fxRunner)(nil)

// Notify is the main function that notifies the application of the key change.
func (fr *fxRunner) Notify(event Event) {
	fr.UnmarshalKeyNotifier.Notify(event)
	fr.processEvent(event)
}

func (fr *fxRunner) processEvent(event Event) {
	logger := fr.fxRunnerStatusRegistry.GetLogger()
	switch event.Type {
	case Write:
		logger.Info().Str("event", event.String()).Msg("key update")
		if fr.app != nil {
			// stop existing app
			err := fr.deinitApp()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to stop existing app")
			}
		}
		// instantiate and start a new app
		err := fr.initApp(event.Key)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to instantiate and start a new app")
		}
	case Remove:
		logger.Info().Str("event", event.String()).Msg("key removed")
		// deinit the app
		err := fr.deinitApp()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to deinit app")
		}
		fr.statusRegistry.Detach()
	}
}

func (fr *fxRunner) initApp(key Key) error {
	fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner initializing"), nil))
	logger := fr.fxRunnerStatusRegistry.GetLogger()

	if fr.app == nil && fr.Unmarshaller != nil {
		var options []fx.Option
		for _, fxOptionsFunc := range fr.fxOptionsFuncs {
			o, e := fxOptionsFunc(key, fr.Unmarshaller, fr.statusRegistry)
			if e != nil {
				logger.Error().Err(e).Msg("fxOptionsFunc failed")
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
			// Supply status registry for the key
			fx.Supply(fr.statusRegistry),
			// Supply prometheus registry
			fx.Supply(fr.prometheusRegistry),
			option,
		)

		var err error
		if err = fr.app.Err(); err != nil {
			visualize, _ := fx.VisualizeError(err)
			logger.Error().Err(err).Str("visualize", visualize).Msg("fx.New failed")
			fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(nil, err))
			_ = fr.deinitApp()
			return err
		}

		fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner starting"), nil))

		ctx, cancel := context.WithTimeout(context.Background(), fr.app.StartTimeout())
		defer cancel()
		if err = fr.app.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("Could not start application")
			fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(nil, err))
			return err
		}
		fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner started"), nil))
	} else {
		fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(nil, fr.err))
	}
	return nil
}

func (fr *fxRunner) deinitApp() error {
	fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner stopping"), nil))
	logger := fr.fxRunnerStatusRegistry.GetLogger()
	if fr.app != nil {
		ctx, cancel := context.WithTimeout(context.Background(), fr.app.StopTimeout())
		defer func() { fr.app = nil }()
		defer cancel()
		if err := fr.app.Stop(ctx); err != nil {
			logger.Error().Err(err).Msg("Could not stop application")
			return err
		}
	}
	fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner stopped"), nil))
	return nil
}

// FxDriver tracks prefix and allows spawning "mini FX-based apps" per key in the prefix.
type FxDriver struct {
	StatusRegistry     status.Registry
	PrometheusRegistry *prometheus.Registry
	// Options for new unmarshaller instances
	UnmarshalPrefixNotifier

	// function to provide fx.Options.
	//
	// Resulting fx.Options will be used to create a "mini FX-based apps" per key.
	// The lifecycle of the app will be tied to the existence of the key.
	// Note that when key's contents change the previous App will be stopped
	// and a fresh one will be created.
	FxOptionsFuncs []FxOptionsFunc
}

// Make sure FxDriver implements PrefixNotifier.
var _ PrefixNotifier = (*FxDriver)(nil)

// GetKeyNotifier returns a KeyNotifier that will notify the driver of key changes.
func (fxDriver *FxDriver) GetKeyNotifier(key Key) KeyNotifier {
	statusRegistry := fxDriver.StatusRegistry.Child(key.String())
	fr := &fxRunner{
		UnmarshalKeyNotifier:   fxDriver.getUnmarshalKeyNotifier(key),
		fxOptionsFuncs:         fxDriver.FxOptionsFuncs,
		statusRegistry:         statusRegistry,
		fxRunnerStatusRegistry: statusRegistry.Child("fx_runner"),
		prometheusRegistry:     fxDriver.PrometheusRegistry,
	}

	return fr
}
