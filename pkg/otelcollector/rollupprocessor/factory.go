package rollupprocessor

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

const (
	// The value of "type" key in configuration.
	typeStr = "rollup"
)

var defaultRollupBuckets = []float64{10, 25, 100, 250, 1000, 2500, 10000}

// NewFactory returns a new factory for the Rollup processor.
func NewFactory(promRegistry *prometheus.Registry) processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig(promRegistry),
		processor.WithLogs(CreateLogsProcessor, component.StabilityLevelDevelopment))
}

func createDefaultConfig(promRegistry *prometheus.Registry) func() component.Config {
	return func() component.Config {
		return &Config{
			AttributeCardinalityLimit: defaultAttributeCardinalityLimit,
			RollupBuckets:             defaultRollupBuckets,
			promRegistry:              promRegistry,
		}
	}
}

// CreateLogsProcessor returns rollupProcessor handling logs.
func CreateLogsProcessor(
	_ context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (processor.Logs, error) {
	return newRollupLogsProcessor(set, nextConsumer, cfg.(*Config))
}
