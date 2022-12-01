package alertsexporter

import (
	"context"

	"github.com/go-openapi/strfmt"
	prommodels "github.com/prometheus/alertmanager/api/v2/models"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
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

// ConsumeLogs sends alert from logs to alert manager clients.
func (ex *alertsExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	alerts := prommodels.PostableAlerts{}

	otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) otelcollector.IterAction {
		attributes := logRecord.Attributes()

		genURLStr := ""
		if genURL, ok := attributes.Get(otelcollector.AlertGeneratorURLLabel); ok {
			genURLStr = genURL.Str()
		}
		singleAlert := prommodels.Alert{
			GeneratorURL: strfmt.URI(genURLStr),
		}
		postableAlert := &prommodels.PostableAlert{
			Alert: singleAlert,
		}
		alerts = append(alerts, postableAlert)

		return otelcollector.Keep
	})

	for _, amClient := range ex.cfg.alertMgr.Clients {
		err := amClient.SendAlerts(ctx, alerts)
		if err != nil {
			log.Warn().Err(err).Msgf("could not send alerts for client: %+v", amClient.GetName())
		}
	}

	return nil
}
