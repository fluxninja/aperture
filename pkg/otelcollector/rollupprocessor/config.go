package rollupprocessor

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for rollup processor.
type Config struct {
	config.ProcessorSettings  `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	AttributeCardinalityLimit int                      `mapstructure:"attribute_cardinality_limit"`
	RollupBuckets             []float64                `mapstructure:"rollup_buckets"`

	promRegistry *prometheus.Registry
}

var _ component.Config = (*Config)(nil)
