package alertsexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"

	"github.com/fluxninja/aperture/pkg/alertmanager"
)

const (
	typeStr   = "alerts"
	stability = component.StabilityLevelInDevelopment
)

// NewFactory creates a factory for alerts exporter.
func NewFactory(alertMgr *alertmanager.AlertManager) component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig(alertMgr),
		component.WithLogsExporter(createLogsExporter, stability))
}

func createDefaultConfig(alertMgr *alertmanager.AlertManager) func() config.Exporter {
	return func() config.Exporter {
		return &Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
			alertMgr:         alertMgr,
		}
	}
}

func createLogsExporter(
	_ context.Context,
	_ component.ExporterCreateSettings,
	eConf config.Exporter,
) (component.LogsExporter, error) {
	cfg := eConf.(*Config)
	ex, err := newExporter(cfg)
	if err != nil {
		return nil, err
	}

	return ex, err
}
