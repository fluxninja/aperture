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
		component.WithMetricsProcessor(createMetricsProcessor, component.StabilityLevelInDevelopment),
	)
}

func createDefaultConfig(cache *entitycache.EntityCache) component.ProcessorCreateDefaultConfigFunc {
	return func() component.ProcessorConfig {
		return &Config{
			ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
			entityCache:       cache,
		}
	}
}

func createMetricsProcessor(
	ctx context.Context,
	params component.ProcessorCreateSettings,
	cfg component.ProcessorConfig,
	nextMetricsConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	cfgTyped := cfg.(*Config)
	proc := newProcessor(cfgTyped.entityCache)
	return processorhelper.NewMetricsProcessor(
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
