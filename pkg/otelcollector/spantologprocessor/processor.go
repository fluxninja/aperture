package spantologprocessor

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/multierr"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

type processorImp struct {
	lock         sync.Mutex
	config       Config
	logsExporter component.LogsExporter
	nextConsumer consumer.Traces
}

func newProcessor(config config.Processor, nextConsumer consumer.Traces) (*processorImp, error) {
	log.Info().Msg("Building spantologprocessor")

	pConfig := config.(*Config)

	return &processorImp{
		config:       *pConfig,
		nextConsumer: nextConsumer,
	}, nil
}

// Start implements the component.Component interface.
func (p *processorImp) Start(ctx context.Context, host component.Host) error {
	log.Info().Msg("Starting spantologprocessor")

	exporters := host.GetExporters()

	var availableLogsExporters []string

	// The available list of exporters come from any configured metrics pipelines' exporters.
	for k, exp := range exporters[config.LogsDataType] {
		logsExp, ok := exp.(component.LogsExporter)
		if !ok {
			return fmt.Errorf("the exporter %q isn't a metrics exporter", k.String())
		}

		availableLogsExporters = append(availableLogsExporters, k.String())

		log.Debug().Str("spantolog-exporter", p.config.LogsExporter).Strs("available-exporters", availableLogsExporters).Msg("Looking for spanmetrics exporter from available exporters")
		if k.String() == p.config.LogsExporter {
			p.logsExporter = logsExp
			log.Info().Str("spantolog-exporter", p.config.LogsExporter).Msg("Found exporter")
			break
		}
	}
	if p.logsExporter == nil {
		return fmt.Errorf("failed to find metrics exporter: '%s'; please configure metrics_exporter from one of: %+v", p.config.LogsExporter, availableLogsExporters)
	}
	log.Info().Msg("Started spantologprocessor")
	return nil
}

// Shutdown implements the component.Component interface.
func (p *processorImp) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down spantologprocessor")
	return nil
}

// Capabilities implements the consumer interface.
func (p *processorImp) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeTraces implements the consumer.Traces interface.
// It generates logs, forwarding these logs to the discovered logs exporter.
// The original input trace data will be forwarded to the next consumer, unmodified.
func (p *processorImp) ConsumeTraces(ctx context.Context, traces ptrace.Traces) error {
	// Forward trace data unmodified and propagate both logs and trace pipeline errors, if any.
	return multierr.Combine(p.tracesToLogs(ctx, traces), p.nextConsumer.ConsumeTraces(ctx, traces))
}

func (p *processorImp) tracesToLogs(ctx context.Context, traces ptrace.Traces) error {
	p.lock.Lock()

	l, err := p.buildLogs(traces)

	// This component no longer needs to read the log once built, so it is safe to unlock.
	p.lock.Unlock()

	if err != nil {
		return err
	}

	if err = p.logsExporter.ConsumeLogs(ctx, l); err != nil {
		return err
	}

	return nil
}

// buildLogs collects the computed raw log data, builds the log object and
// writes the raw log data into the log object.
func (p *processorImp) buildLogs(traces ptrace.Traces) (plog.Logs, error) {
	l := plog.NewLogs()

	err := otelcollector.IterateSpans(traces, func(s ptrace.Span) error {
		spanAttributes := s.Attributes()

		ill := l.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()
		ill.Scope().SetName("spantologsprocessor")

		lr := ill.LogRecords().AppendEmpty()

		spanAttributes.CopyTo(lr.Attributes())

		return nil
	})
	if err != nil {
		return plog.Logs{}, err
	}

	return l, nil
}
