package grpc

import (
	"container/list"
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
)

// ClientModule is an fx module that provides annotated grpc ClientConnectionBuilder.
func ClientModule() fx.Option {
	return ClientConstructor{}.Annotate()
}

// ClientConstructor holds fields to create an annotated instance of ClientConnectionBuilder.
type ClientConstructor struct {
	Name          string
	Key           string
	DefaultConfig GRPCClientConfig
}

// GRPCClientConfig holds configuration for GRPC Client.
// swagger:model
type GRPCClientConfig struct {
	// Minimum connection timeout
	MinConnectionTimeout config.Duration `json:"min_connection_timeout" validate:"gte=0" default:"20s"`
	// Client TLS configuration
	ClientTLSConfig tlsconfig.ClientTLSConfig `json:"tls"`
	// Backoff config
	Backoff BackoffConfig `json:"backoff"`
	// Disable ClientTLS
	Insecure bool `json:"insecure" default:"false"`
	// Use HTTP CONNECT Proxy
	UseProxy bool `json:"use_proxy" default:"false"`
}

// BackoffConfig holds configuration for GRPC Client Backoff.
// swagger:model
type BackoffConfig struct {
	// Base Delay
	BaseDelay config.Duration `json:"base_delay" validate:"gte=0" default:"1s"`
	// Max Delay
	MaxDelay config.Duration `json:"max_delay" validate:"gte=0" default:"120s"`
	// Backoff multiplier
	Multiplier float64 `json:"multiplier" validate:"gte=0" default:"1.6"`
	// Jitter
	Jitter float64 `json:"jitter" validate:"gte=0" default:"0.2"`
}

// Annotate creates an annotated instance of GRPC ClientConnectionBuilder.
func (c ClientConstructor) Annotate() fx.Option {
	if c.Key == "" {
		log.Panic().Msg("config key not provided")
	}

	name := config.NameTag(c.Name)
	cfgName := config.NameTag(c.Name + "-config")
	return fx.Provide(
		fx.Annotate(
			c.provideClientConnectionBuilder,
			fx.ResultTags(name, cfgName),
		),
	)
}

func (c ClientConstructor) provideClientConnectionBuilder(unmarshaller config.Unmarshaller) (ClientConnectionBuilder, *GRPCClientConfig, error) {
	config := c.DefaultConfig
	err := unmarshaller.UnmarshalKey(c.Key, &config)
	if err != nil {
		return nil, nil, err
	}

	dialOptions, err := config.ClientTLSConfig.GetGRPCDialOptions(config.Insecure)
	if err != nil {
		return nil, nil, err
	}

	builder := newClientConnectionBuilder()

	dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(
		grpc_prometheus.UnaryClientInterceptor,
		otelgrpc.UnaryClientInterceptor(),
	))
	dialOptions = append(dialOptions, grpc.WithChainStreamInterceptor(
		grpc_prometheus.StreamClientInterceptor,
		otelgrpc.StreamClientInterceptor(),
	))
	dialOptions = append(dialOptions, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  config.Backoff.BaseDelay.Duration.AsDuration(),
			Multiplier: config.Backoff.Multiplier,
			Jitter:     config.Backoff.Jitter,
			MaxDelay:   config.Backoff.MaxDelay.Duration.AsDuration(),
		},
		MinConnectTimeout: config.MinConnectionTimeout.Duration.AsDuration(),
	}))

	if !config.UseProxy {
		dialOptions = append(dialOptions, grpc.WithNoProxy())
	}

	return builder.AddOptions(dialOptions...), &config, nil
}

// ClientConnectionBuilder is a convenience builder to gather []grpc.DialOption.
type ClientConnectionBuilder interface {
	AddOptions(opts ...grpc.DialOption) ClientConnectionBuilder
	Build() ClientConnectionWrapper
}

// ClientConnectionWrapper is a convenience wrapper to support predefined dial Options provided by ClientConnectionBuilder.
type ClientConnectionWrapper interface {
	// Context can be nil
	Dial(ctx context.Context, target string, extraOptions ...grpc.DialOption) (*grpc.ClientConn, error)
}

type clientConnectionBuilder struct {
	ll *list.List
}

// make sure clientConnectionBuilder implements ClientConnectionBuilder.
var _ ClientConnectionBuilder = (*clientConnectionBuilder)(nil)

func newClientConnectionBuilder() ClientConnectionBuilder {
	return &clientConnectionBuilder{
		ll: list.New(),
	}
}

type clientConnectionOptions struct {
	options []grpc.DialOption
}

// AddOptions adds grpc.DialOptions to the ClientConnectionBuilder.
func (b *clientConnectionBuilder) AddOptions(opts ...grpc.DialOption) ClientConnectionBuilder {
	b.ll.PushBack(func(cco *clientConnectionOptions) {
		cco.options = append(cco.options, opts...)
	})
	return b
}

// Build iterates through and collects grpc.DialOptions, builds a ClientConnectionWrapper and returns it.
func (b *clientConnectionBuilder) Build() ClientConnectionWrapper {
	cco := new(clientConnectionOptions)
	for e := b.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(connOptions *clientConnectionOptions))
		f(cco)
	}
	return &clientConnectionWrapper{
		options: cco,
	}
}

type clientConnectionWrapper struct {
	options *clientConnectionOptions
}

// make sure clientConnectionWrapper implements ClientConnectionWrapper.
var _ ClientConnectionWrapper = (*clientConnectionWrapper)(nil)

// Dial creates a client connection to the given target with the given options and returns the connection interface.
func (w *clientConnectionWrapper) Dial(ctx context.Context, target string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOptions := w.options.options
	dialOptions = append(dialOptions, options...)
	if ctx == nil {
		ctx = context.Background()
	}
	return grpc.DialContext(ctx, target, dialOptions...)
}
