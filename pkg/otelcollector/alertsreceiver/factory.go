package alertsreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/fluxninja/aperture/pkg/alerts"
)

const (
	typeStr   = "alerts"
	stability = component.StabilityLevelDevelopment
)

// NewFactory creates a factory for alerts receiver.
func NewFactory(alerter alerts.Alerter) receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig(alerter),
		receiver.WithLogs(createLogsReceiver, stability))
}

func createDefaultConfig(alerter alerts.Alerter) func() component.Config {
	return func() component.Config {
		return &Config{
			alerter: alerter,
		}
	}
}

func createLogsReceiver(
	_ context.Context,
	_ receiver.CreateSettings,
	rConf component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	cfg := rConf.(*Config)
	p, err := newReceiver(cfg)
	if err != nil {
		return nil, err
	}
	err = p.registerLogsConsumer(consumer)
	if err != nil {
		return nil, err
	}
	return p, err
}
