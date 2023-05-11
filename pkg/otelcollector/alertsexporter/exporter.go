package alertsexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

type alertsExporter struct {
	cfg *Config
}

func newExporter(cfg *Config) (*alertsExporter, error) {
	ex := &alertsExporter{
		cfg: cfg,
	}

	return ex, nil
}

// Start implements the Component interface.
func (ex *alertsExporter) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Shutdown implements the Component interface.
func (ex *alertsExporter) Shutdown(_ context.Context) error {
	return nil
}

// Capabilities returns the capabilities of the exporter.
func (ex *alertsExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: false,
	}
}

// ConsumeLogs sends alert from logs to alert manager clients.
func (ex *alertsExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	alerts := alerts.AlertsFromLogs(ld)

	for _, amClient := range ex.cfg.alertMgr.Clients {
		log.Trace().Int("alerts", ld.LogRecordCount()).Str("client", amClient.GetName()).Msg("Exporting alerts")
		err := amClient.SendAlerts(ctx, alerts)
		if err != nil {
			log.Warn().Err(err).Msgf("could not send alerts for client: %+v", amClient.GetName())
		}
	}

	return nil
}
