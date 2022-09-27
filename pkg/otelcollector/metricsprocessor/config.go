package metricsprocessor

import (
	"github.com/fluxninja/aperture/pkg/policies/dataplane/iface"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/config"
)

// Config holds configuration for the metrics processor.
type Config struct {
	promRegistry             *prometheus.Registry
	metricsAPI               iface.ResponseMetricsAPI
	engine                   iface.Engine
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}
