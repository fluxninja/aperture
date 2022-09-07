// +kubebuilder:validation:Optional
package grpc

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/panichandler"
)

const (
	defaultServerKey = "server.grpc"
	// Name of gmux based listener.
	defaultGMuxListener = "grpc-gmux-listener"
)

// ServerModule is an fx module that provides annotated GRPC Server using the default listener and registers its metrics with the prometheus registry.
func ServerModule() fx.Option {
	return fx.Options(
		ServerConstructor{}.Annotate(),
		fx.Invoke(RegisterGRPCServerMetrics),
	)
}

// GMuxServerModule is an fx module that provides annotated grpc Server using gmux provided listener and registers its metrics with the prometheus registry.
func GMuxServerModule() fx.Option {
	return fx.Options(
		listener.GMuxConstructor{ListenerName: defaultGMuxListener}.Annotate(),
		ServerConstructor{ListenerName: defaultGMuxListener}.Annotate(),
		fx.Invoke(RegisterGRPCServerMetrics),
	)
}

// GRPCServerConfig holds configuration for GRPC Server.
// swagger:model
// +kubebuilder:object:generate=true
type GRPCServerConfig struct {
	// Connection timeout
	ConnectionTimeout config.Duration `json:"connection_timeout" validate:"gte=0s" default:"120s"`
	// Enable Reflection
	EnableReflection bool `json:"enable_reflection" default:"false"`
}

// ServerConstructor holds fields to create an annotated GRPC Server.
type ServerConstructor struct {
	// Name of grpc server instance -- empty for main server
	Name string
	// Name of listener instance
	ListenerName string
	// Viper config key/server name
	Key string
	// Default Server Config
	DefaultConfig GRPCServerConfig
	// Additional server Options
	ServerOptions []grpc.ServerOption
}

// Annotate creates an annotated instance of GRPC Server.
func (constructor ServerConstructor) Annotate() fx.Option {
	if constructor.Key == "" {
		constructor.Key = defaultServerKey
	}

	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor.provideServer,
				fx.ParamTags(config.NameTag(constructor.ListenerName)),
				fx.ResultTags(config.NameTag(constructor.Name), config.NameTag(constructor.Name)),
			),
		),
	)
}

func (constructor ServerConstructor) provideServer(
	listener *listener.Listener,
	unmarshaller config.Unmarshaller,
	lifecycle fx.Lifecycle,
	shutdowner fx.Shutdowner,
) (*grpc.Server, *grpc_prometheus.ServerMetrics, error) {
	config := constructor.DefaultConfig
	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize grpcserver configuration!")
		return nil, nil, err
	}

	grpcServerMetrics := grpc_prometheus.NewServerMetrics()

	// Connection timeout from config
	constructor.ServerOptions = append(constructor.ServerOptions, grpc.ConnectionTimeout(config.ConnectionTimeout.AsDuration()))

	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		grpcServerMetrics.UnaryServerInterceptor(),
		otelgrpc.UnaryServerInterceptor(),
	}
	constructor.ServerOptions = append(constructor.ServerOptions, grpc.ChainUnaryInterceptor(unaryServerInterceptors...))

	streamServerInterceptors := []grpc.StreamServerInterceptor{
		grpcServerMetrics.StreamServerInterceptor(),
		otelgrpc.StreamServerInterceptor(),
	}
	constructor.ServerOptions = append(constructor.ServerOptions, grpc.ChainStreamInterceptor(streamServerInterceptors...))

	server := grpc.NewServer(constructor.ServerOptions...)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			panichandler.Go(func() {
				// request shutdown if this server exits
				defer func() { _ = shutdowner.Shutdown() }()
				listener := listener.GetListener()
				log.Info().Str("constructor", constructor.Key).Str("addr", listener.Addr().String()).Msg("Starting GRPC server")

				grpcServerMetrics.InitializeMetrics(server)

				if config.EnableReflection {
					reflection.Register(server)
				}

				if err := server.Serve(listener); err != nil {
					log.Error().Err(err).Msg("Unable to start GRPC server!")
				}
			})
			return nil
		},
		OnStop: func(context.Context) error {
			listener := listener.GetListener()
			log.Info().Str("constructor", constructor.Key).Str("addr", listener.Addr().String()).Msg("Stopping GRPC server")
			server.GracefulStop()
			return nil
		},
	})

	return server, grpcServerMetrics, nil
}

// RegisterGRPCServerMetrics registers a collection of metrics provided by grpc_prometheus.ServerMetrics with a prometheus registry.
func RegisterGRPCServerMetrics(metrics *grpc_prometheus.ServerMetrics, pr *prometheus.Registry) error {
	err := pr.Register(metrics)
	if err != nil {
		return err
	}
	return nil
}
