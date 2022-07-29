package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

// Module is a fx module that provides Prometheus registry and invokes registering metrics handler.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(ProvidePrometheusRegistry),
		fx.Invoke(RegisterMetricsHandler),
	)
}

const (
	// MetricsConfigKey is the key for the Prometheus configuration.
	metricsConfigKey = "metrics"
	// MetricsEndpoint is default endpoint which is used to register metrics handler.
	metricsEndpoint = "/metrics"
)

// swagger:operation POST /metrics common-configuration Metrics
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/MetricsConfig"

// MetricsConfig holds configuration for service metrics.
// swagger:model
type MetricsConfig struct {
	// Pedantic controls whether a pedantic Registerer is used as the prometheus backend. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewPedanticRegistry>
	Pedantic bool `json:"pedantic" default:"false"`

	// EnableGoCollector controls whether the go collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewGoCollector>
	EnableGoCollector bool `json:"enable_go_metrics" default:"false"`

	// EnableProcessCollector controls whether the process collector is registered on startup. See <https://godoc.org/github.com/prometheus/client_golang/prometheus#NewProcessCollector>
	EnableProcessCollector bool `json:"enable_process_collector" default:"false"`
}

// ProvidePrometheusRegistry creates a new Prometheus Registry and provides it via Fx.
// Metrics from the Registry are served on /metrics via the default http server of Fx application.
func ProvidePrometheusRegistry(unmarshaller config.Unmarshaller) (*prometheus.Registry, error) {
	var pr *prometheus.Registry
	var config MetricsConfig

	if err := unmarshaller.UnmarshalKey(metricsConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("Unable to deserialize prometheus configuration!")
		return nil, err
	}

	if config.Pedantic {
		pr = prometheus.NewRegistry()
	} else {
		pr = prometheus.NewPedanticRegistry()
	}

	if config.EnableGoCollector {
		if err := pr.Register(collectors.NewGoCollector()); err == nil {
			log.Warn().Err(err).Msg("Unable to register prometheus go collector!")
		} else {
			log.Info().Msg("Registering Go Metrics collector")
		}
	}

	if config.EnableProcessCollector {
		if err := pr.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
			log.Warn().Err(err).Msg("Unable to register prometheus process collector!")
		} else {
			log.Info().Msg("Registering Process Metrics collector")
		}
	}

	return pr, nil
}

// RegisterMetricsHandler registers the metrics handler on the promhttp server.
func RegisterMetricsHandler(router *mux.Router, pr *prometheus.Registry) {
	log.Info().Msg("Registering Prometheus metrics endpoint")
	logger := log.Component("PROM_HTTP")
	router.Handle(metricsEndpoint, promhttp.HandlerFor(pr, promhttp.HandlerOpts{ErrorLog: &logger}))
}
