package alertsexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/fluxninja/aperture/pkg/alertmanager"
)

const (
	typeStr   = "alerts"
	stability = component.StabilityLevelDevelopment
)

// NewFactory creates a factory for alerts exporter.
func NewFactory(alertMgr *alertmanager.AlertManager) component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig(alertMgr),
		component.WithLogsExporter(createLogsExporter, stability))
}

func createDefaultConfig(alertMgr *alertmanager.AlertManager) func() component.ExporterConfig {
	return func() component.ExporterConfig {
		return &Config{
			ExporterSettings: config.NewExporterSettings(component.NewID(typeStr)),
			TimeoutSettings:  exporterhelper.NewDefaultTimeoutSettings(),
			alertMgr:         alertMgr,
		}
	}
}

func createLogsExporter(
	_ context.Context,
	_ component.ExporterCreateSettings,
	eConf component.ExporterConfig,
) (component.LogsExporter, error) {
	cfg := eConf.(*Config)
	ex, err := newExporter(cfg)
	if err != nil {
		return nil, err
	}

	return ex, err
}
