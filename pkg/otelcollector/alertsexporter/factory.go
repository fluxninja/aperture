package alertsexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/fluxninja/aperture/v2/pkg/alertmanager"
)

const (
	typeStr   = "alerts"
	stability = component.StabilityLevelDevelopment
)

// NewFactory creates a factory for alerts exporter.
func NewFactory(alertMgr *alertmanager.AlertManager) exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig(alertMgr),
		exporter.WithLogs(createLogsExporter, stability))
}

func createDefaultConfig(alertMgr *alertmanager.AlertManager) func() component.Config {
	return func() component.Config {
		return &Config{
			TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
			alertMgr:        alertMgr,
		}
	}
}

func createLogsExporter(
	_ context.Context,
	_ exporter.CreateSettings,
	eConf component.Config,
) (exporter.Logs, error) {
	cfg := eConf.(*Config)
	ex, err := newExporter(cfg)
	if err != nil {
		return nil, err
	}

	return ex, err
}
