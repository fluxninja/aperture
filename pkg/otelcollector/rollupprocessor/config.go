package rollupprocessor

import (
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for rollup processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}

var _ config.Processor = (*Config)(nil)
