package enrichmentprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/fluxninja/aperture/pkg/entitycache"
)

const (
	// The value of "type" key in configuration.
	typeStr = "enrichment"
)

// NewFactory returns a new factory for the enrichment processor.
func NewFactory(cache *entitycache.EntityCache) processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig(cache),
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig(cache *entitycache.EntityCache) component.CreateDefaultConfigFunc {
	return func() component.Config {
		return &Config{
			entityCache: cache,
		}
	}
}

func createMetricsProcessor(
	ctx context.Context,
	params processor.CreateSettings,
	cfg component.Config,
	nextMetricsConsumer consumer.Metrics,
) (processor.Metrics, error) {
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
