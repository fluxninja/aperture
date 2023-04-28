package adapterconnector

import (
	"context"
	"sync"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	typeStr = "adapter"
)

// NewFactory returns a connector.Factory.
func NewFactory() connector.Factory {
	return connector.NewFactory(
		typeStr,
		createDefaultConfig,
		connector.WithTracesToLogs(createTracesToLogs, component.StabilityLevelDevelopment),
	)
}

// createDefaultConfig creates the default configuration.
func createDefaultConfig() component.Config {
	return &struct{}{}
}

// createTracesToLogs creates a trace receiver based on provided config.
func createTracesToLogs(
	_ context.Context,
	_ connector.CreateSettings,
	_ component.Config,
	nextConsumer consumer.Logs,
) (connector.Traces, error) {
	return &adapter{Logs: nextConsumer}, nil
}

// adapter is used to pass signals of one type to a pipeline with another type.
type adapter struct {
	lock sync.Mutex
	consumer.Logs
}

// ConsumeTraces implements the consumer.Traces interface.
func (a *adapter) ConsumeTraces(ctx context.Context, traces ptrace.Traces) error {
	return a.tracesToLogs(ctx, traces)
}

func (a *adapter) tracesToLogs(ctx context.Context, traces ptrace.Traces) error {
	a.lock.Lock()

	l := a.buildLogs(traces)

	// This component no longer needs to read the log once built, so it is safe to unlock.
	a.lock.Unlock()

	if err := a.Logs.ConsumeLogs(ctx, l); err != nil {
		return err
	}

	return nil
}

// buildLogs collects the computed raw log data, builds the log object and
// writes the raw log data into the log object.
func (a *adapter) buildLogs(traces ptrace.Traces) plog.Logs {
	l := plog.NewLogs()

	otelcollector.IterateSpans(traces, func(s ptrace.Span) {
		spanAttributes := s.Attributes()

		ill := l.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()
		ill.Scope().SetName("adapterconnector")

		lr := ill.LogRecords().AppendEmpty()

		spanAttributes.CopyTo(lr.Attributes())
	})

	return l
}

// Start implements the component.Component interface.
func (a *adapter) Start(_ context.Context, _ component.Host) error {
	log.Info().Msg("Starting adapterconnector")
	return nil
}

// Shutdown implements the component.Component interface.
func (a *adapter) Shutdown(_ context.Context) error {
	log.Info().Msg("Shutting down adapterconnector")
	return nil
}

// Capabilities returns the capabilities of the adapter.
func (a *adapter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
