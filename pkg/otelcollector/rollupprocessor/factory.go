package rollupprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// The value of "type" key in configuration.
	typeStr = "rollup"
)

// NewFactory returns a new factory for the Rollup processor.
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithLogsProcessor(CreateLogsProcessor, component.StabilityLevelDevelopment))
}

func createDefaultConfig() component.ProcessorConfig {
	return &Config{
		ProcessorSettings:         config.NewProcessorSettings(component.NewID(typeStr)),
		AttributeCardinalityLimit: defaultAttributeCardinalityLimit,
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
