package alertsreceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"

	"github.com/fluxninja/aperture/pkg/alerts"
)

// Config for alerts receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`

	alerter alerts.Alerter
}

var _ component.Config = (*Config)(nil)
