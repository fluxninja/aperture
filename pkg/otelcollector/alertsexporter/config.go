package alertsexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/fluxninja/aperture/v2/pkg/alertmanager"
)

// Config for alerts exporter.
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	alertMgr                       *alertmanager.AlertManager
}

var _ component.Config = (*Config)(nil)
