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
	stability = component.StabilityLevelInDevelopment
)

// NewFactory creates a factory for alerts receiver.
func NewFactory(alerter alerts.Alerter) component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig(alerter),
		component.WithLogsReceiver(createLogsReceiver, stability))
}

func createDefaultConfig(alerter alerts.Alerter) func() config.Receiver {
	return func() config.Receiver {
		return &Config{
			ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
			alerter:          alerter,
		}
	}
}

func createLogsReceiver(
	_ context.Context,
	_ component.ReceiverCreateSettings,
	rConf config.Receiver,
	consumer consumer.Logs,
) (component.LogsReceiver, error) {
	cfg := rConf.(*Config)
	p, err := newProcessor(cfg)
	if err != nil {
		return nil, err
	}
	err = p.registerLogsConsumer(consumer)
	if err != nil {
		return nil, err
	}
	return p, err
}
