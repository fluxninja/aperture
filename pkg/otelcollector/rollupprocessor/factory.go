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
		component.WithLogsProcessor(CreateLogsProcessor, component.StabilityLevelInDevelopment),
		component.WithTracesProcessor(CreateTracesProcessor, component.StabilityLevelInDevelopment))
}

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
	}
}

// CreateLogsProcessor returns rollupProcessor handling logs.
func CreateLogsProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	return newRollupLogsProcessor(set, nextConsumer, cfg.(*Config))
}

// CreateTracesProcessor returns rollupProcessor handling traces.
func CreateTracesProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	return newRollupTracesProcessor(set, nextConsumer, cfg.(*Config))
}
