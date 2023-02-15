package rollupprocessor

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for rollup processor.
type Config struct {
	AttributeCardinalityLimit int       `mapstructure:"attribute_cardinality_limit"`
	RollupBuckets             []float64 `mapstructure:"rollup_buckets"`

	promRegistry *prometheus.Registry
}

var _ component.Config = (*Config)(nil)
