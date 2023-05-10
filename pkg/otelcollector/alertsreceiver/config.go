package alertsreceiver

import (
	"go.opentelemetry.io/collector/component"

	"github.com/fluxninja/aperture/v2/pkg/alerts"
)

// Config for alerts receiver.
type Config struct {
	alerter alerts.Alerter
}

var _ component.Config = (*Config)(nil)
