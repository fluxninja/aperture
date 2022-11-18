package alertsexporter

import (
	"go.opentelemetry.io/collector/config"

	"github.com/fluxninja/aperture/pkg/alertmanager"
)

// Config for alerts exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	alertMgr                *alertmanager.AlertManager
}

var _ config.Exporter = (*Config)(nil)
