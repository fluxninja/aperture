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
	Stage                     Stage                    `mapstructure:"stage"`

	promRegistry *prometheus.Registry
}

// Stage defines whether this is initial or aggregating rollup.
type Stage string

const (
	// InitialStage is rollup made directly from individual logs.
	InitialStage Stage = "initial"

	// Aggregation is rollup of rollup results.
	Aggregation Stage = "aggregation"
)

var _ component.ProcessorConfig = (*Config)(nil)
