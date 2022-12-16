package tracestologsprocessor

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/multierr"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

type tracesToLogsProcessor struct {
	lock         sync.Mutex
	config       Config
	logsExporter exporter.Logs
	nextConsumer consumer.Traces
}

func newProcessor(config component.Config, nextConsumer consumer.Traces) (*tracesToLogsProcessor, error) {
	log.Info().Msg("Building tracestologsprocessor")

	pConfig := config.(*Config)

	return &tracesToLogsProcessor{
		config:       *pConfig,
		nextConsumer: nextConsumer,
	}, nil
}

// Start implements the component.Component interface.
func (p *tracesToLogsProcessor) Start(ctx context.Context, host component.Host) error {
	log.Info().Msg("Starting tracestologsprocessor")

	exporters := host.GetExporters()

	var availableLogsExporters []string

	// The available list of exporters come from any configured metrics pipelines' exporters.
	for k, exp := range exporters[component.DataTypeLogs] {
		logsExp, ok := exp.(exporter.Logs)
		if !ok {
			return fmt.Errorf("the exporter %q isn't a metrics exporter", k.String())
		}

		availableLogsExporters = append(availableLogsExporters, k.String())

		log.Debug().
			Str("configured-exporter", p.config.LogsExporter).
			Strs("available-exporters", availableLogsExporters).
			Msg("Looking for configured exporter from available exporters")
		if k.String() == p.config.LogsExporter {
			p.logsExporter = logsExp
			log.Info().Str("configured-exporter", p.config.LogsExporter).Msg("Found exporter")
			break
		}
	}
	if p.logsExporter == nil {
		return fmt.Errorf("failed to find logs exporter: '%s'. Available exporters: %+v",
			p.config.LogsExporter, availableLogsExporters)
	}
	log.Info().Msg("Started tracestologsprocessor")
	return nil
}

// Shutdown implements the component.Component interface.
func (p *tracesToLogsProcessor) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down tracestologsprocessor")
	return nil
}

// Capabilities implements the consumer interface.
func (p *tracesToLogsProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeTraces implements the consumer.Traces interface.
// It generates logs, forwarding these logs to the discovered logs exporter.
// The original input trace data will be forwarded to the next consumer, unmodified.
func (p *tracesToLogsProcessor) ConsumeTraces(ctx context.Context, traces ptrace.Traces) error {
	// Forward trace data unmodified and propagate both logs and trace pipeline errors, if any.
	return multierr.Combine(p.tracesToLogs(ctx, traces), p.nextConsumer.ConsumeTraces(ctx, traces))
}

func (p *tracesToLogsProcessor) tracesToLogs(ctx context.Context, traces ptrace.Traces) error {
	p.lock.Lock()

	l := p.buildLogs(traces)

	// This component no longer needs to read the log once built, so it is safe to unlock.
	p.lock.Unlock()

	if err := p.logsExporter.ConsumeLogs(ctx, l); err != nil {
		return err
	}

	return nil
}

// buildLogs collects the computed raw log data, builds the log object and
// writes the raw log data into the log object.
func (p *tracesToLogsProcessor) buildLogs(traces ptrace.Traces) plog.Logs {
	l := plog.NewLogs()

	otelcollector.IterateSpans(traces, func(s ptrace.Span) {
		spanAttributes := s.Attributes()

		ill := l.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()
		ill.Scope().SetName("tracestologsprocessor")

		lr := ill.LogRecords().AppendEmpty()

		spanAttributes.CopyTo(lr.Attributes())
	})

	return l
}
