package otelcollector

import (
	"encoding/json"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// IterateLogRecords calls given function for each logRecord. If the function
// returns error further logRecords will not be processed and the error will be returned.
func IterateLogRecords(ld plog.Logs, fn func(plog.LogRecord) error) error {
	resourceLogsSlice := ld.ResourceLogs()
	for resourceLogsIt := 0; resourceLogsIt < resourceLogsSlice.Len(); resourceLogsIt++ {
		resourceLogs := resourceLogsSlice.At(resourceLogsIt)
		scopeLogsSlice := resourceLogs.ScopeLogs()
		for scopeLogsIt := 0; scopeLogsIt < scopeLogsSlice.Len(); scopeLogsIt++ {
			scopeLogs := scopeLogsSlice.At(scopeLogsIt)
			logsSlice := scopeLogs.LogRecords()
			for logsIt := 0; logsIt < logsSlice.Len(); logsIt++ {
				logRecord := logsSlice.At(logsIt)
				err := fn(logRecord)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// IterateSpans calls given function for each span. If the function returns error
// further span will not be processed and the error will be returned.
func IterateSpans(td ptrace.Traces, fn func(ptrace.Span) error) error {
	resourceSpansSlice := td.ResourceSpans()
	for resourceSpansIt := 0; resourceSpansIt < resourceSpansSlice.Len(); resourceSpansIt++ {
		resourceSpans := resourceSpansSlice.At(resourceSpansIt)
		scopeSpansSlice := resourceSpans.ScopeSpans()
		for scopeSpansIt := 0; scopeSpansIt < scopeSpansSlice.Len(); scopeSpansIt++ {
			scopeSpans := scopeSpansSlice.At(scopeSpansIt)
			spansSlice := scopeSpans.Spans()
			for spansIt := 0; spansIt < spansSlice.Len(); spansIt++ {
				span := spansSlice.At(spansIt)
				err := fn(span)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// IterateMetrics calls given function for each metric. If the function returns error
// further metric will not be processed and the error will be returned.
func IterateMetrics(md pmetric.Metrics, fn func(pmetric.Metric) error) error {
	resourceMetricsSlice := md.ResourceMetrics()
	for resourceMetricsIt := 0; resourceMetricsIt < resourceMetricsSlice.Len(); resourceMetricsIt++ {
		resourceMetrics := resourceMetricsSlice.At(resourceMetricsIt)
		scopeMetricsSlice := resourceMetrics.ScopeMetrics()
		for scopeMetricsIt := 0; scopeMetricsIt < scopeMetricsSlice.Len(); scopeMetricsIt++ {
			scopeMetrics := scopeMetricsSlice.At(scopeMetricsIt)
			metricsSlice := scopeMetrics.Metrics()
			for metricsIt := 0; metricsIt < metricsSlice.Len(); metricsIt++ {
				span := metricsSlice.At(metricsIt)
				err := fn(span)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// IterateDataPoints calls given function for each metric data point. If the function returns error
// further data point will not be processed and the error will be returned.
func IterateDataPoints(metric pmetric.Metric, fn func(pcommon.Map) error) error {
	switch metric.DataType() {
	case pmetric.MetricDataTypeGauge:
		dataPoints := metric.Gauge().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			err := fn(dataPoints.At(i).Attributes())
			if err != nil {
				return err
			}
		}
	case pmetric.MetricDataTypeSum:
		dataPoints := metric.Sum().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			err := fn(dataPoints.At(i).Attributes())
			if err != nil {
				return err
			}
		}
	case pmetric.MetricDataTypeSummary:
		dataPoints := metric.Summary().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			err := fn(dataPoints.At(i).Attributes())
			if err != nil {
				return err
			}
		}
	case pmetric.MetricDataTypeHistogram:
		dataPoints := metric.Histogram().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			err := fn(dataPoints.At(i).Attributes())
			if err != nil {
				return err
			}
		}
	case pmetric.MetricDataTypeExponentialHistogram:
		dataPoints := metric.ExponentialHistogram().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			err := fn(dataPoints.At(i).Attributes())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO (hasit): The following unmarshaling function
// 1. should be renamed and moved to a common package along with the corresponding marshal method

// UnmarshalStringVal is a helper for cases we're sending more complex
// structure json-encoded in a string label
//
// Returns whether label was actually a string and unmarshaling was attempted.
func UnmarshalStringVal(value pcommon.Value, labelName string, output interface{}) bool {
	if value.Type() != pcommon.ValueTypeString {
		return false
	}

	stringVal := value.StringVal()

	if stringVal == MissingAttributeSourceValue {
		log.Trace().Str("label", labelName).Msg("Missing attribute source")
		return true
	}

	err := json.Unmarshal([]byte(stringVal), output)
	if err != nil {
		log.Debug().Err(err).Str("label", labelName).Msg("Unmarshalling")
		// This is almost impossible to happen (eg. broken sdk), so ignoring error is ok
	}

	return true
}

// GetCheckResponse unmarshalls limiter decisions from string label.
func GetCheckResponse(attributes pcommon.Map) (checkResponse *flowcontrolv1.CheckResponse) {
	rawLimiterDecisions, exists := attributes.Get(MarshalledCheckResponseLabel)
	if !exists {
		log.Debug().Str("label", MarshalledCheckResponseLabel).Msg("Label does not exist")
		return
	}
	if !UnmarshalStringVal(rawLimiterDecisions, MarshalledCheckResponseLabel, &checkResponse) {
		log.Debug().Str("label", MarshalledCheckResponseLabel).Msg("Label is not a string")
	}
	return
}
