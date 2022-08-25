package enrichmentprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

const (
	// The value of "type" key in configuration.
	typeStr = "enrichment"
)

// NewFactory returns a new factory for the enrichment processor.
func NewFactory(cache *entitycache.EntityCache) component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig(cache),
		component.WithLogsProcessor(createLogsProcessor, component.StabilityLevelInDevelopment),
		component.WithTracesProcessor(createTracesProcessor, component.StabilityLevelInDevelopment),
		component.WithMetricsProcessor(createMetricsProcessor, component.StabilityLevelInDevelopment),
	)
}

func createDefaultConfig(cache *entitycache.EntityCache) component.ProcessorCreateDefaultConfigFunc {
	return func() config.Processor {
		return &Config{
			ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
			entityCache:       cache,
		}
	}
}

func createLogsProcessor(
	ctx context.Context,
	params component.ProcessorCreateSettings,
	cfg config.Processor,
	nextLogsConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	cfgTyped := cfg.(*Config)
	proc := newProcessor(cfgTyped.entityCache)
	return processorhelper.NewLogsProcessorWithCreateSettings(
		ctx,
		params,
		cfg,
		nextLogsConsumer,
		proc.ConsumeLogs,
		processorhelper.WithCapabilities(proc.Capabilities()),
		processorhelper.WithStart(proc.Start),
		processorhelper.WithShutdown(proc.Shutdown),
	)
}

func createTracesProcessor(
	ctx context.Context,
	params component.ProcessorCreateSettings,
	cfg config.Processor,
	nextTracesConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	cfgTyped := cfg.(*Config)
	proc := newProcessor(cfgTyped.entityCache)
	return processorhelper.NewTracesProcessorWithCreateSettings(
		ctx,
		params,
		cfg,
		nextTracesConsumer,
		proc.ConsumeTraces,
		processorhelper.WithCapabilities(proc.Capabilities()),
		processorhelper.WithStart(proc.Start),
		processorhelper.WithShutdown(proc.Shutdown),
	)
}

func createMetricsProcessor(
	ctx context.Context,
	params component.ProcessorCreateSettings,
	cfg config.Processor,
	nextMetricsConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	cfgTyped := cfg.(*Config)
	proc := newProcessor(cfgTyped.entityCache)
	return processorhelper.NewMetricsProcessorWithCreateSettings(
		ctx,
		params,
		cfg,
		nextMetricsConsumer,
		proc.ConsumeMetrics,
		processorhelper.WithCapabilities(proc.Capabilities()),
		processorhelper.WithStart(proc.Start),
		processorhelper.WithShutdown(proc.Shutdown),
	)
}
