package alertsreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"

	"github.com/fluxninja/aperture/pkg/alerts"
)

const (
	typeStr   = "alerts"
	stability = component.StabilityLevelDevelopment
)

// NewFactory creates a factory for alerts receiver.
func NewFactory(alerter alerts.Alerter) component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig(alerter),
		component.WithLogsReceiver(createLogsReceiver, stability))
}

func createDefaultConfig(alerter alerts.Alerter) func() component.ReceiverConfig {
	return func() component.ReceiverConfig {
		return &Config{
			ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
			alerter:          alerter,
		}
	}
}

func createLogsReceiver(
	_ context.Context,
	_ component.ReceiverCreateSettings,
	rConf component.ReceiverConfig,
	consumer consumer.Logs,
) (component.LogsReceiver, error) {
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
