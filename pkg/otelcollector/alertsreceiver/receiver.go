package alertsreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"github.com/fluxninja/aperture/pkg/log"
)

type alertsReceiver struct {
	cfg          *Config
	logsConsumer consumer.Logs

	// shutdown kills the long running context. Should be set in Start()
	shutdown func()
}

func newReceiver(cfg *Config) (*alertsReceiver, error) {
	p := &alertsReceiver{
		cfg: cfg,
	}

	return p, nil
}

// Start creates long running context, saves shutdown function and starts run() goroutine.
func (p *alertsReceiver) Start(_ context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.shutdown = cancel
	go p.run(ctx)
	return nil
}

// Shutdown calls shutdown function which kills run() goroutine.
func (p *alertsReceiver) Shutdown(_ context.Context) error {
	p.shutdown()
	return nil
}

func (p *alertsReceiver) registerLogsConsumer(lc consumer.Logs) error {
	if lc == nil {
		return component.ErrNilNextConsumer
	}
	p.logsConsumer = lc
	return nil
}

func (p *alertsReceiver) run(ctx context.Context) {
	for {
		select {
		case alert := <-p.cfg.alerter.AlertsChan():
			err := p.logsConsumer.ConsumeLogs(ctx, alert.AsLogs())
			// We do not care much about those errors. Alerts can be dropped sometimes,
			// they are sent all the time anyway.
			log.Autosample().Debug().Err(err).Msg("ConsumeLogs failed")
		case <-ctx.Done():
			return
		}
	}
}
