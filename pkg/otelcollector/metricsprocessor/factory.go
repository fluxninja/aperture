package metricsprocessor

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/fluxninja/aperture/pkg/cache"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/selectors"
)

const (
	// The value of "type" key in configuration.
	typeStr = "metrics"
)

// NewFactory returns a new factory for the metrics processor.
func NewFactory(
	promRegistry *prometheus.Registry,
	engine iface.Engine,
	clasEng iface.ClassificationEngine,
	controlPointCache *cache.Cache[selectors.TypedControlPointID],
) processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig(promRegistry, engine, clasEng, controlPointCache),
		processor.WithLogs(createLogsProcessor, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig(
	promRegistry *prometheus.Registry,
	engine iface.Engine,
	clasEng iface.ClassificationEngine,
	controlPointCache *cache.Cache[selectors.TypedControlPointID],
) component.CreateDefaultConfigFunc {
	return func() component.Config {
		return &Config{
			promRegistry:         promRegistry,
			engine:               engine,
			classificationEngine: clasEng,
			controlPointCache:    controlPointCache,
		}
	}
}

func createLogsProcessor(
	ctx context.Context,
	params processor.CreateSettings,
	cfg component.Config,
	nextLogsConsumer consumer.Logs,
) (processor.Logs, error) {
	cfgTyped := cfg.(*Config)
	proc, err := newProcessor(cfgTyped)
	if err != nil {
		return nil, err
	}
	return processorhelper.NewLogsProcessor(
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
