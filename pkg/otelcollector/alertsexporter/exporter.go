package alertsexporter

import (
	"context"

	"github.com/fluxninja/aperture/pkg/log"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
)

type alertsExporter struct {
	cfg          *Config
	logsConsumer consumer.Logs

	// shutdown kills the long running context. Should be set in Start()
	shutdown func()
}

func newExporter(cfg *Config) (*alertsExporter, error) {
	ex := &alertsExporter{
		cfg: cfg,
	}

	return ex, nil
}

// Start is TODO.
func (ex *alertsExporter) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Shutdown is TODO.
func (ex *alertsExporter) Shutdown(_ context.Context) error {
	return nil
}

// Capabilities returns the capabilities of the exporter.
func (ex *alertsExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: false,
	}
}

// ConsumeLogs is TODO.
func (ex *alertsExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	log.Error().Msgf("DARIA LOG CONSUME LOGS IN EXPORTER")
	return nil
}
