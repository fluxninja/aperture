package rollupprocessor

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for rollup processor.
type Config struct {
	config.ProcessorSettings  `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	AttributeCardinalityLimit int                      `mapstructure:"attribute_cardinality_limit"`
}

var _ component.ProcessorConfig = (*Config)(nil)
