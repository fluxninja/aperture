package notifiers

import (
	"context"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// FxOptionsFunc is a function that returns fx.Option.
type FxOptionsFunc func(Key, config.Unmarshaller, status.Registry) (fx.Option, error)

type fxRunner struct {
	statusRegistry         status.Registry
	fxRunnerStatusRegistry status.Registry
	app                    *fx.App
	prometheusRegistry     *prometheus.Registry
	*unmarshalKeyNotifier
	fxOptionsFuncs []FxOptionsFunc
}

// Make sure fxRunner implements KeyNotifier.
var _ KeyNotifier = (*fxRunner)(nil)

func (fr *fxRunner) processEvent(event Event, unmarshaller config.Unmarshaller) {
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
		err := fr.initApp(event.Key, unmarshaller)
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

func (fr *fxRunner) initApp(key Key, unmarshaller config.Unmarshaller) error {
	fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner initializing"), nil))
	logger := fr.fxRunnerStatusRegistry.GetLogger()

	if fr.app == nil && unmarshaller != nil {
		var options []fx.Option
		for _, fxOptionsFunc := range fr.fxOptionsFuncs {
			o, e := fxOptionsFunc(key, unmarshaller, fr.statusRegistry)
			if e != nil {
				logger.Error().Err(e).Msg("fxOptionsFunc failed")
				return e
			}
			options = append(options, o)
		}
		option := fx.Options(options...)

		fr.app = fx.New(
			// Note: Supplying unmarshaller directly results in supplying
			// concrete type instead of interface, thus supplying via Provide.
			fx.Provide(func() config.Unmarshaller { return unmarshaller }),
			// Supply keyinfo
			fx.Supply(key),
			// Supply status registry for the key
			fx.Supply(fr.statusRegistry),
			// Supply prometheus registry
			fx.Supply(fr.prometheusRegistry),
			option,
			fx.WithLogger(func() fxevent.Logger {
				logger := zap.New(
					log.NewZapAdapter(log.GetGlobalLogger(), fmt.Sprintf("fxdriver-%s", key.String())),
				)
				return &fxevent.ZapLogger{Logger: logger}
			}),
		)

		var err error
		if err = fr.app.Err(); err != nil {
			visualize, _ := fx.VisualizeError(err)
			logger.Error().Err(err).Str("visualize", visualize).Msg("fx.New failed")
			fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(nil, err))
			if deinitErr := fr.deinitApp(); deinitErr != nil {
				logger.Error().Err(deinitErr).Msg("Failed to deinitialize application after start failure")
			}
			return err
		}

		fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner starting"), nil))

		ctx, cancel := context.WithTimeout(context.Background(), fr.app.StartTimeout())
		defer cancel()
		if err = fr.app.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("Could not start application")
			fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(nil, err))
			if deinitErr := fr.deinitApp(); deinitErr != nil {
				logger.Error().Err(deinitErr).Msg("Failed to deinitialize application after start failure")
			}
			return err
		}
		fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner started"), nil))
	} else {
		fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(nil, errors.New("fxRunner is not initialized")))
	}
	return nil
}

func (fr *fxRunner) deinitApp() error {
	fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner stopping"), nil))
	logger := fr.fxRunnerStatusRegistry.GetLogger()
	if fr.app != nil {
		ctx, cancel := context.WithTimeout(context.Background(), fr.app.StopTimeout())
		defer func() {
			cancel()
			fr.app = nil
		}()
		if err := fr.app.Stop(ctx); err != nil {
			logger.Error().Err(err).Msg("Could not stop application")
			return err
		}
	}
	fr.fxRunnerStatusRegistry.SetStatus(status.NewStatus(wrapperspb.String("policy runner stopped"), nil))
	return nil
}

// fxDriver tracks prefix and allows spawning "mini FX-based apps" per key in the prefix.
type fxDriver struct {
	statusRegistry     status.Registry
	prometheusRegistry *prometheus.Registry
	// Options for new unmarshaller instances
	*unmarshalPrefixNotifier

	// function to provide fx.Options.
	//
	// Resulting fx.Options will be used to create a "mini FX-based apps" per key.
	// The lifecycle of the app will be tied to the existence of the key.
	// Note that when key's contents change the previous App will be stopped
	// and a fresh one will be created.
	fxOptionsFuncs []FxOptionsFunc
}

// Make sure FxDriver implements PrefixNotifier.
var _ PrefixNotifier = (*fxDriver)(nil)

// NewFxDriver creates a new FxDriver.
func NewFxDriver(
	statusRegistry status.Registry,
	prometheusRegistry *prometheus.Registry,
	getUnmarshallerFunc GetUnmarshallerFunc,
	fxOptionsFuncs []FxOptionsFunc,
) (*fxDriver, error) {
	// Subscribe to all prefixes and additional notifier callback is nil
	upn, err := NewUnmarshalPrefixNotifier("", nil, getUnmarshallerFunc)
	if err != nil {
		return nil, err
	}
	return &fxDriver{
		statusRegistry:          statusRegistry,
		prometheusRegistry:      prometheusRegistry,
		unmarshalPrefixNotifier: upn,
		fxOptionsFuncs:          fxOptionsFuncs,
	}, nil
}

// GetKeyNotifier returns a KeyNotifier that will notify the driver of key changes.
func (fxd *fxDriver) GetKeyNotifier(key Key) (KeyNotifier, error) {
	unmarshaller, err := fxd.getUnmarshallerFunc(nil)
	if err != nil {
		return nil, err
	}

	statusRegistry := fxd.statusRegistry.Child("key", key.String())
	fr := &fxRunner{
		fxOptionsFuncs:         fxd.fxOptionsFuncs,
		statusRegistry:         statusRegistry,
		fxRunnerStatusRegistry: statusRegistry.Child("system", "fx_runner"),
		prometheusRegistry:     fxd.prometheusRegistry,
	}

	ukn, err := NewUnmarshalKeyNotifier(key, unmarshaller, fr.processEvent)
	if err != nil {
		return nil, err
	}
	fr.unmarshalKeyNotifier = ukn

	return fr, nil
}
