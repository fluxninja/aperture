// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loggingexporter // import "go.opentelemetry.io/collector/exporter/loggingexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/FluxNinja/aperture/pkg/log"
	"github.com/FluxNinja/aperture/pkg/otelcollector/loggingexporter/otlptext"
)

type loggingExporter struct {
	logsMarshaler    plog.Marshaler
	metricsMarshaler pmetric.Marshaler
	tracesMarshaler  ptrace.Marshaler
}

func (s *loggingExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	log.Trace().Int("#spans", td.SpanCount()).Msg("TracesExporter")

	buf, err := s.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	log.Trace().Msg(string(buf))
	return nil
}

func (s *loggingExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	log.Trace().Int("#metrics", md.MetricCount()).Msg("MetricsExporter")

	buf, err := s.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	log.Trace().Msg(string(buf))
	return nil
}

func (s *loggingExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	log.Trace().Int("#logs", ld.LogRecordCount()).Msg("LogsExporter")

	buf, err := s.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	log.Trace().Msg(string(buf))
	return nil
}

func newLoggingExporter() *loggingExporter {
	return &loggingExporter{
		logsMarshaler:    otlptext.NewTextLogsMarshaler(),
		metricsMarshaler: otlptext.NewTextMetricsMarshaler(),
		tracesMarshaler:  otlptext.NewTextTracesMarshaler(),
	}
}

// newTracesExporter creates an exporter.TracesExporter that just drops the
// received data and logs debugging messages.
func newTracesExporter(config config.Exporter, set component.ExporterCreateSettings) (component.TracesExporter, error) {
	s := newLoggingExporter()
	return exporterhelper.NewTracesExporter(
		config,
		set,
		s.pushTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}

// newMetricsExporter creates an exporter.MetricsExporter that just drops the
// received data and logs debugging messages.
func newMetricsExporter(config config.Exporter, set component.ExporterCreateSettings) (component.MetricsExporter, error) {
	s := newLoggingExporter()
	return exporterhelper.NewMetricsExporter(
		config,
		set,
		s.pushMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}

// newLogsExporter creates an exporter.LogsExporter that just drops the
// received data and logs debugging messages.
func newLogsExporter(config config.Exporter, set component.ExporterCreateSettings) (component.LogsExporter, error) {
	s := newLoggingExporter()
	return exporterhelper.NewLogsExporter(
		config,
		set,
		s.pushLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// Disable Timeout/RetryOnFailure and SendingQueue
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}
