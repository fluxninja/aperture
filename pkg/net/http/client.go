package http

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/net/tlsconfig"
)

// ClientModule is an fx module that provides annotated HTTP client.
func ClientModule() fx.Option {
	return fx.Options(
		ClientConstructor{}.Annotate(),
	)
}

// ClientConstructor holds fields to create an annotated instance of HTTP client.
type ClientConstructor struct {
	Name          string
	ConfigKey     string
	DefaultConfig HTTPClientConfig
}

// HTTPClientConfig holds configuration for HTTP Client.
// swagger:model
// +kubebuilder:object:generate=true
type HTTPClientConfig struct {
	// Network level keep-alive duration
	NetworkKeepAlive config.Duration `json:"network_keep_alive" validate:"gte=0s" default:"30s"`
	// Timeout for making network connection
	NetworkTimeout config.Duration `json:"network_timeout" validate:"gte=0s" default:"30s"`
	// HTTP client timeout - Timeouts includes connection time, redirects, reading the response etc. 0 = no timeout.
	Timeout config.Duration `json:"timeout" validate:"gte=0s" default:"60s"`
	// Proxy Connect Header - map[string][]string
	ProxyConnectHeader http.Header `json:"proxy_connect_header"`
	// TLS Handshake Timeout. 0 = no timeout
	TLSHandshakeTimeout config.Duration `json:"tls_handshake_timeout" validate:"gte=0s" default:"10s"`
	// Expect Continue Timeout. 0 = no timeout.
	ExpectContinueTimeout config.Duration `json:"expect_continue_timeout" validate:"gte=0s" default:"1s"`
	// Response Header Timeout. 0 = no timeout.
	ResponseHeaderTimeout config.Duration `json:"response_header_timeout" validate:"gte=0s" default:"0s"`
	// Idle Connection Timeout. 0 = no timeout.
	IdleConnTimeout config.Duration `json:"idle_connection_timeout" validate:"gte=0s" default:"90s"`
	// SSL key log file (useful for debugging with wireshark)
	KeyLogWriter string `json:"key_log_file" validate:"omitempty,file"`
	// Client TLS configuration
	ClientTLSConfig tlsconfig.ClientTLSConfig `json:"tls"`
	// Max Idle Connections per host. 0 = no limit.
	MaxIdleConnsPerHost int `json:"max_idle_connections_per_host" validate:"gte=0" default:"5"`
	// Max Idle Connections. 0 = no limit.
	MaxIdleConns int `json:"max_idle_connections" validate:"gte=0" default:"100"`
	// Max Connections Per Host. 0 = no limit.
	MaxConnsPerHost int `json:"max_conns_per_host" validate:"gte=0" default:"0"`
	// Max Response Header Bytes. 0 = no limit.
	MaxResponseHeaderBytes int64 `json:"max_response_header_bytes" validate:"gte=0" default:"0"`
	// Write Buffer Size. 0 = 4KB.
	WriteBufferSize int `json:"write_buffer_size" validate:"gte=0" default:"0"`
	// Read Buffer Size. 0 = 4KB
	ReadBufferSize int `json:"read_buffer_size" validate:"gte=0" default:"0"`
	// Disable Compression
	DisableCompression bool `json:"disable_compression" default:"false"`
	// Use Proxy
	UseProxy bool `json:"use_proxy" default:"false"`
	// Disable HTTP Keep Alives
	DisableKeepAlives bool `json:"disable_keep_alives" default:"false"`
}

// Annotate creates an annotated instance of HTTP Client.
func (constructor ClientConstructor) Annotate() fx.Option {
	if constructor.ConfigKey == "" {
		log.Panic().Msg("config key not provided")
	}

	name := config.NameTag(constructor.Name)
	cfgName := config.NameTag(constructor.Name + "-config")

	return fx.Provide(
		fx.Annotate(
			constructor.provideHTTPClient,
			fx.ResultTags(name, name, cfgName),
		),
	)
}

func (constructor ClientConstructor) provideHTTPClient(unmarshaller config.Unmarshaller, lifecycle fx.Lifecycle) (*http.Client, *MiddlewareChain, *HTTPClientConfig, error) {
	var err error

	config := constructor.DefaultConfig
	if err = unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize httpclient configuration!")
		return nil, nil, nil, err
	}

	tlsConfig, err := config.ClientTLSConfig.GetTLSConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: (&net.Dialer{
			Timeout:   config.NetworkTimeout.AsDuration(),
			KeepAlive: config.NetworkKeepAlive.AsDuration(),
		}).DialContext,
		TLSHandshakeTimeout:    config.TLSHandshakeTimeout.AsDuration(),
		DisableKeepAlives:      config.DisableKeepAlives,
		DisableCompression:     config.DisableCompression,
		MaxIdleConns:           config.MaxIdleConns,
		MaxIdleConnsPerHost:    config.MaxIdleConnsPerHost,
		MaxConnsPerHost:        config.MaxConnsPerHost,
		IdleConnTimeout:        config.IdleConnTimeout.AsDuration(),
		ResponseHeaderTimeout:  config.ResponseHeaderTimeout.AsDuration(),
		ExpectContinueTimeout:  config.ExpectContinueTimeout.AsDuration(),
		ProxyConnectHeader:     config.ProxyConnectHeader,
		MaxResponseHeaderBytes: config.MaxResponseHeaderBytes,
		WriteBufferSize:        config.WriteBufferSize,
		ReadBufferSize:         config.ReadBufferSize,
	}

	if config.UseProxy {
		transport.Proxy = http.ProxyFromEnvironment
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout.AsDuration(),
	}

	// return a middleware chain -- call invokes on this object to chain middleware functions
	mwc := &MiddlewareChain{
		client:      client,
		middlewares: []Middleware{},
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// build middleware chain
			mwc.buildChain()
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})

	return client, mwc, &config, nil
}

// inspired by https://github.com/improbable-eng/go-httpwares/blob/master/tripperware.go

// RoundTripperFunc wraps a func to make it into a http.RoundTripper. Similar to http.HandleFunc.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements http.RoundTripper.
func (rtf RoundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return rtf(request)
}

// Middleware is signature of all http middleware.
type Middleware func(next http.RoundTripper) http.RoundTripper

// MiddlewareChain holds a chain of middleware.
type MiddlewareChain struct {
	client      *http.Client
	middlewares []Middleware
}

// Chain appends provided middleware to the MiddlewareChain.
// Middleware will be chained based on the order of Invokes.
func (mwc *MiddlewareChain) Chain(middlewares ...Middleware) {
	mwc.middlewares = append(mwc.middlewares, middlewares...)
}

func (mwc *MiddlewareChain) buildChain() {
	if len(mwc.middlewares) == 0 {
		return
	}

	transport := mwc.client.Transport

	for i := len(mwc.middlewares) - 1; i >= 0; i-- {
		transport = mwc.middlewares[i](transport)
	}
	mwc.client.Transport = transport
}
