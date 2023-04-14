// +kubebuilder:validation:Optional
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
	defaultServerKey   = "server.http"
	defaultHandlerName = "default"
)

type monitoringContext struct {
	context.Context
	handlerName string
}

// ServerModule is an fx module that provides annotated HTTP Server using the default listener and registers its metrics with the prometheus registry.
func ServerModule() fx.Option {
	return fx.Options(
		ServerConstructor{}.Annotate(),
	)
}

// HTTPServerConfig holds configuration for HTTP Server.
// swagger:model
// +kubebuilder:object:generate=true
type HTTPServerConfig struct {
	// Idle timeout
	IdleTimeout config.Duration `json:"idle_timeout" validate:"gte=0s" default:"30s"`
	// Read header timeout
	ReadHeaderTimeout config.Duration `json:"read_header_timeout" validate:"gte=0s" default:"10s"`
	// Read timeout
	ReadTimeout config.Duration `json:"read_timeout" validate:"gte=0s" default:"10s"`
	// Write timeout
	WriteTimeout config.Duration `json:"write_timeout" validate:"gte=0s" default:"45s"`
	// Max header size in bytes
	MaxHeaderBytes int `json:"max_header_bytes" validate:"gte=0" default:"1048576"`
	// Buckets specification in latency histogram
	LatencyBucketsMS []float64 `json:"latency_buckets_ms" validate:"gte=0" default:"[10.0,25.0,100.0,250.0,1000.0]"`
	// Disable HTTP Keepalive
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
	ConfigKey string
	// Default Server Config
	DefaultConfig HTTPServerConfig
}

// Annotate creates an annotated instance of HTTP Server.
func (constructor ServerConstructor) Annotate() fx.Option {
	if constructor.ConfigKey == "" {
		constructor.ConfigKey = defaultServerKey
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
	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize httpserver configuration!")
		return nil, nil, nil, err
	}

	// Register metrics
	defaultLabels := []string{metrics.MethodLabel, metrics.StatusCodeLabel, metrics.HandlerName}
	errorCounters := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.HTTPErrorMetricName,
		Help: "The total number of errors that occurred",
	}, defaultLabels)
	requestCounters := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.HTTPRequestMetricName,
		Help: "The total number of requests that occurred",
	}, defaultLabels)
	// We record latency milliseconds
	latencyHistograms := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metrics.HTTPRequestLatencyMetricName,
		Help:    "Latency of the requests processed by the server",
		Buckets: config.LatencyBucketsMS,
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
		IdleTimeout:       config.IdleTimeout.AsDuration(),
		ReadHeaderTimeout: config.ReadHeaderTimeout.AsDuration(),
		ReadTimeout:       config.ReadTimeout.AsDuration(),
		WriteTimeout:      config.WriteTimeout.AsDuration(),
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
		OnStart: func(_ context.Context) error {
			panichandler.Go(func() {
				// request shutdown if this server exits
				defer func() { _ = shutdowner.Shutdown() }()

				listener := listener.GetListener()

				if tlsConfig != nil {
					listener = tls.NewListener(listener, tlsConfig)
				}

				log.Info().Str("constructor", constructor.ConfigKey).Str("addr", listener.Addr().String()).Msg("Starting HTTP server")
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
			log.Info().Str("constructor", constructor.ConfigKey).Str("addr", listener.Addr().String()).Msg("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})

	return router, server, httpServer, nil
}

func (s *Server) monitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &monitoringContext{}
		ctx.Context = r.Context()
		ctx.handlerName = defaultHandlerName
		startTime := time.Now()
		rec := newStatusRecorder(w)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rec, r.WithContext(ctx))

		duration := time.Since(startTime)

		labels := map[string]string{
			metrics.MethodLabel:     r.Method,
			metrics.StatusCodeLabel: fmt.Sprintf("%d", rec.statusCode),
			metrics.HandlerName:     ctx.handlerName,
		}

		requestCounter, err := s.RequestCounters.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract request counter metric from registry: %v", err)
		} else {
			requestCounter.Inc()
		}

		errorCounter, err := s.ErrorCounters.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract error counter metric from registry: %v", err)
		} else if rec.statusCode >= http.StatusBadRequest {
			errorCounter.Inc()
		}

		latencyHistogram, err := s.Latencies.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract latency histogram metric from registry: %v", err)
		} else {
			latencyHistogram.Observe(float64(duration.Milliseconds()))
		}
	})
}

// HandlerNameMiddleware sets handler name in monitoring context.
func HandlerNameMiddleware(handlerName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if mCtx, ok := ctx.(*monitoringContext); ok {
				mCtx.handlerName = handlerName
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{ResponseWriter: w}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader records statusCode and calls wrapped WriteHeader method.
func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.statusCode = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}
