package listener

import (
	"context"
	"net"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

const (
	defaultKey = "server.listener"
)

// Listener wraps net.Listener, that can be potentially not-yet-started.
type Listener struct {
	// This will be available in the Start stage.
	lis net.Listener
	// This will be available in the Provide stage.
	addr string
}

// GetListener returns wrapped Listener
//
// This function is supposed to only be called in the Start stage. Otherwise,
// this function will panic.
func (l *Listener) GetListener() net.Listener {
	if l.lis == nil {
		log.Panic().Msg("Listener is not yet started")
	}
	return l.lis
}

// GetAddr returns the address of the listener.
func (l *Listener) GetAddr() string {
	return l.addr
}

// Module is an fx module that provides annotated Listener.
func Module() fx.Option {
	return fx.Options(fx.Provide(Constructor{}.ProvideAnnotated()))
}

// Constructor holds fields to create an annotated Listener.
type Constructor struct {
	ConfigKey     string
	Name          string
	DefaultConfig ListenerConfig
}

// ProvideAnnotated provides an annotated instance of Listener.
func (constructor Constructor) ProvideAnnotated() fx.Annotated {
	if constructor.ConfigKey == "" {
		constructor.ConfigKey = defaultKey
	}
	return fx.Annotated{
		Name:   constructor.Name,
		Target: constructor.provideListener,
	}
}

// ListenerIn is the input to Listener constructor.
type ListenerIn struct {
	fx.In

	Unmarshaller config.Unmarshaller
	Lifecycle    fx.Lifecycle
}

func (constructor Constructor) provideListener(in ListenerIn) (*Listener, error) {
	config := constructor.DefaultConfig

	if err := in.Unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize listener configuration!")
		return nil, err
	}

	return provideListenerFromConfig(in.Lifecycle, config), nil
}

// provideListenerFromConfig provides listener using an already-parsed config.
func provideListenerFromConfig(lc fx.Lifecycle, config ListenerConfig) *Listener {
	listener := Listener{}
	listener.addr = config.Addr

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var err error
			listener.lis, err = newListener(config)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create listener")
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			_ = listener.lis.Close()
			return nil
		},
	})

	return &listener
}

// ProvideTestListener provides a listener which will choose the port to listen
// to on its own.
// Use fx.Populate an `port = lis.GetListener().Addr().(*net.TCPAddr).Port` to
// obtain actual port number
//
// Note: This is a separate function, because Addr: ":0" fails hostname_port
// validation thus cannot go through Unmarshaller-based constructor.
func ProvideTestListener(lc fx.Lifecycle) *Listener {
	return provideListenerFromConfig(lc, ListenerConfig{
		Addr:    "127.0.0.1:0",
		Network: "tcp",
	})
}
