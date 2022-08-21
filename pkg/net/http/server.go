package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/net/listener"
	"github.com/fluxninja/aperture/pkg/panichandler"
)

const (
	defaultServerKey = "server.http"
)

// ServerModule is an fx module that provides annotated HTTP Server using the default listener and registers its metrics with the prometheus registry.
func ServerModule() fx.Option {
	return fx.Options(
		ServerConstructor{}.Annotate(),
	)
}

// HTTPServerConfig holds configuration for HTTP Server.
// swagger:model
type HTTPServerConfig struct {
	// Idle timeout
	IdleTimeout config.Duration `json:"idle_timeout" validate:"gte=0s" default:"30s"`
	// Read header timeout
	ReadHeaderTimeout config.Duration `json:"read_header_timeout" validate:"gte=0s" default:"10s"`
	// Read timeout
	ReadTimeout config.Duration `json:"read_timeout" validate:"gte=0s" default:"10s"`
	// Write timeout
	WriteTimeout config.Duration `json:"write_timeout" validate:"gte=0s" default:"10s"`
	// The lowest bucket in latency histogram
	LatencyBucketStartMS float64 `json:"latency_bucket_start_ms" validate:"gte=0" default:"20"`
	// Max header size in bytes
	MaxHeaderBytes int `json:"max_header_bytes" validate:"gte=0" default:"1048576"`
	// The bucket width in latency histogram
	LatencyBucketWidthMS float64 `json:"latency_bucket_width_ms" validate:"gte=0" default:"20"`
	// The number of buckets in latency histogram
	LatencyBucketCount int `json:"latency_bucket_count" validate:"gte=0" default:"100"`
	// Disable HTTP Keep Alives
	DisableHTTPKeepAlives bool `json:"disable_http_keep_alives" default:"false"`
}

// ServerConstructor holds fields to create an annotated HTTP Server.
type ServerConstructor struct {
	// Name of http server instance -- empty for main server
	Name string
	// Name of listener instance
	ListenerName string
	// Name of tls config instance
	TLSConfigName string
	// Viper config key/server name
	Key string
	// Default Server Config
	DefaultConfig HTTPServerConfig
}

// Annotate creates an annotated instance of HTTP Server.
func (constructor ServerConstructor) Annotate() fx.Option {
	if constructor.Key == "" {
		constructor.Key = defaultServerKey
	}
	tlsName := config.NameTag(constructor.TLSConfigName) + ` optional:"true"`
	name := config.NameTag(constructor.Name)
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor.provideServer,
				fx.ParamTags(config.NameTag(constructor.ListenerName), tlsName),
				fx.ResultTags(name, name, name),
			),
		),
	)
}

// Server holds fields for custom HTTP server.
type Server struct {
	Server    *http.Server
	Mux       *mux.Router
	Listener  *listener.Listener
	TLSConfig *tls.Config
	// As we are using Gorilla Mux, root handler is registered the last as a catch all
	RootHandler     http.Handler
	RequestCounters *prometheus.CounterVec
	ErrorCounters   *prometheus.CounterVec
	Latencies       *prometheus.HistogramVec
}

func (constructor ServerConstructor) provideServer(
	listener *listener.Listener,
	tlsConfig *tls.Config,
	unmarshaller config.Unmarshaller,
	lifecycle fx.Lifecycle,
	shutdowner fx.Shutdowner,
	pr *prometheus.Registry,
) (*mux.Router, *http.Server, *Server, error) {
	config := constructor.DefaultConfig
	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize httpserver configuration!")
		return nil, nil, nil, err
	}

	// Register metrics
	defaultLabels := []string{metrics.MethodLabel, metrics.ResponseStatusCodeLabel}
	errorCounters := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.ErrorCountMetricName,
		Help: "The total number of errors that occurred",
	}, defaultLabels)
	requestCounters := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.RequestCounterMetricName,
		Help: "The total number of requests that occurred",
	}, defaultLabels)
	// We record latency milliseconds
	latencyHistograms := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metrics.LatencyHistogramMetricName,
		Help:    "Latency of the requests processed by the server",
		Buckets: prometheus.LinearBuckets(config.LatencyBucketStartMS, config.LatencyBucketWidthMS, config.LatencyBucketCount),
	}, defaultLabels)
	for _, metric := range []prometheus.Collector{errorCounters, requestCounters, latencyHistograms} {
		err := pr.Register(metric)
		if err != nil {
			// Ignore already registered error, as this is not harmful. Metrics may
			// be registered by other running server.
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, nil, nil, fmt.Errorf("couldn't register prometheus metrics: %w", err)
			}
		}
	}

	router := mux.NewRouter()

	server := &http.Server{
		Handler:           router,
		MaxHeaderBytes:    config.MaxHeaderBytes,
		IdleTimeout:       config.IdleTimeout.Duration.AsDuration(),
		ReadHeaderTimeout: config.ReadHeaderTimeout.Duration.AsDuration(),
		ReadTimeout:       config.ReadTimeout.Duration.AsDuration(),
		WriteTimeout:      config.WriteTimeout.Duration.AsDuration(),
		TLSConfig:         tlsConfig,
	}

	httpServer := &Server{
		Server:          server,
		Mux:             router,
		Listener:        listener,
		TLSConfig:       tlsConfig,
		RequestCounters: requestCounters,
		ErrorCounters:   errorCounters,
		Latencies:       latencyHistograms,
	}
	router.Use(httpServer.monitoringMiddleware)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			panichandler.Go(func() {
				// request shutdown if this server exits
				defer func() { _ = shutdowner.Shutdown() }()

				listener := listener.GetListener()

				if tlsConfig != nil {
					listener = tls.NewListener(listener, tlsConfig)
				}

				log.Info().Str("constructor", constructor.Key).Str("addr", listener.Addr().String()).Msg("Starting HTTP server")
				// check if RootHandler is set
				if httpServer.RootHandler != nil {
					log.Info().Msg("Registering RootHandlerFunc!")
					router.PathPrefix("/").Handler(httpServer.RootHandler)
				}
				if err := server.Serve(listener); err != http.ErrServerClosed {
					log.Error().Err(err).Msg("Unable to start HTTP server!")
				}
			})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			listener := listener.GetListener()
			log.Info().Str("constructor", constructor.Key).Str("addr", listener.Addr().String()).Msg("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})

	return router, server, httpServer, nil
}

func (s *Server) monitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := statusRecorder{w, 200}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(&rec, r)
		duration := time.Since(startTime)

		labels := map[string]string{
			metrics.MethodLabel:             r.Method,
			metrics.ResponseStatusCodeLabel: fmt.Sprintf("%d", rec.status),
		}

		requestCounter, err := s.RequestCounters.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract request counter metric from registry: %v", err)
		} else {
			requestCounter.Inc()
		}

		latencyHistogram, err := s.Latencies.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract latency histogram metric from registry: %v", err)
		} else {
			latencyHistogram.Observe(float64(duration.Milliseconds()))
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}
