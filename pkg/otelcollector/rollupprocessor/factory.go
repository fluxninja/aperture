package rollupprocessor

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// The value of "type" key in configuration.
	typeStr = "rollup"
)

var defaultRollupBuckets = []float64{10, 25, 100, 250, 1000, 2500, 10000}

// NewFactory returns a new factory for the Rollup processor.
func NewFactory(promRegistry *prometheus.Registry) component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig(promRegistry),
		component.WithLogsProcessor(CreateLogsProcessor, component.StabilityLevelDevelopment))
}

func createDefaultConfig(promRegistry *prometheus.Registry) func() component.ProcessorConfig {
	return func() component.ProcessorConfig {
		return &Config{
			ProcessorSettings:         config.NewProcessorSettings(component.NewID(typeStr)),
			AttributeCardinalityLimit: defaultAttributeCardinalityLimit,
			RollupBuckets:             defaultRollupBuckets,
			promRegistry:              promRegistry,
		}
	}
}

// CreateLogsProcessor returns rollupProcessor handling logs.
func CreateLogsProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg component.ProcessorConfig,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	return newRollupLogsProcessor(set, nextConsumer, cfg.(*Config))
}
