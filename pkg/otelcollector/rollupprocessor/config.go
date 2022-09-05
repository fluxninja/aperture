package rollupprocessor

import (
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for rollup processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	Rollups []*otelcollector.Rollup `mapstructure:"rollups"`
}

var _ config.Processor = (*Config)(nil)
