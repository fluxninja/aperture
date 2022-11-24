package tracestologsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const typeStr = "tracestologs" // The value of "type" key in configuration.

// NewFactory returns a new factory for the tracestologsprocessor.
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesProcessor(createTracesProcessor, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.ProcessorConfig {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
	}
}

func createTracesProcessor(
	_ context.Context,
	params component.ProcessorCreateSettings,
	cfg component.ProcessorConfig,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	return newProcessor(cfg, nextConsumer)
}
