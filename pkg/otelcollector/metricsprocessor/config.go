package metricsprocessor

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/config"

	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
)

// Config holds configuration for the metrics processor.
type Config struct {
	promRegistry             *prometheus.Registry
	engine                   iface.Engine
	classificationEngine     iface.ClassificationEngine
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}
