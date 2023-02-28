package main

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/plugins"
)

// ServicePlugin returns the plugin.
func ServicePlugin() plugins.ServicePluginIface {
	return &Plugin{}
}

// Plugin implements the plugin interface.
type Plugin struct{}

// Module returns the plugin module.
func (p *Plugin) Module() fx.Option {
	log.Info().Msg("Loading Kubernetes Attributes Processor")
	return fx.Options(
		module(),
	)
}

func module() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideProcessor,
				fx.ResultTags(config.GroupTag(otelcollector.ProcessorFactoriesFxTag)),
			),
		),
	)
}

func provideProcessor() ([]processor.Factory, error) {
	return []processor.Factory{
		k8sattributesprocessor.NewFactory(),
	}, nil
}
