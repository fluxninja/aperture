package metricsprocessor

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/fluxninja/aperture/pkg/controlpointcache"
	"github.com/fluxninja/aperture/pkg/policies/flowcontrol/iface"
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
	controlPointCache *controlpointcache.ControlPointCache,
) component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig(promRegistry, engine, clasEng, controlPointCache),
		component.WithLogsProcessor(createLogsProcessor, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig(
	promRegistry *prometheus.Registry,
	engine iface.Engine,
	clasEng iface.ClassificationEngine,
	controlPointCache *controlpointcache.ControlPointCache,
) component.CreateDefaultConfigFunc {
	return func() component.Config {
		return &Config{
			ProcessorSettings:    config.NewProcessorSettings(component.NewID(typeStr)),
			promRegistry:         promRegistry,
			engine:               engine,
			classificationEngine: clasEng,
			controlPointCache:    controlPointCache,
		}
	}
}

func createLogsProcessor(
	ctx context.Context,
	params component.ProcessorCreateSettings,
	cfg component.Config,
	nextLogsConsumer consumer.Logs,
) (component.LogsProcessor, error) {
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
