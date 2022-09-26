package otelcollector

import (
	"encoding/json"
	"strconv"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/utils"
	"github.com/rs/zerolog"
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

// GetStruct is a helper for decoding complex structs encoded into an attribute
// as a json-encoded string.
// Takes:
// attributes to read from
// label key to read in attributes
// output interface that is filled via json unmarshal
// treatAsMissing is a list of values that are treated as attribute missing from source
//
// Returns true is label was decoded successfully, false otherwise.
func GetStruct(attributes pcommon.Map, label string, output interface{}, treatAsMissing []string) bool {
	value, ok := attributes.Get(label)
	if !ok {
		log.Sample(zerolog.Sometimes).Warn().Str("label", label).Msg("Label does not exist in attributes map")
		return false
	}
	if value.Type() != pcommon.ValueTypeString {
		log.Sample(zerolog.Sometimes).Warn().Str("label", label).Msg("Label is not a string")
		return false
	}

	stringVal := value.StringVal()

	for _, markerForMissing := range treatAsMissing {
		if stringVal == markerForMissing {
			log.Sample(zerolog.Sometimes).Info().Str("label", label).Msg("Missing attribute from source")
			return false
		}
	}

	err := json.Unmarshal([]byte(stringVal), output)
	if err != nil {
		log.Sample(zerolog.Sometimes).Error().Err(err).Str("label", label).Msg("Failed to unmarshal")
	}

	return true
}

// GetFloat64 returns float64 value from given attribute map at given key.
func GetFloat64(attributes pcommon.Map, key string, treatAsZero []string) (float64, bool) {
	rawNewValue, exists := attributes.Get(key)
	if !exists {
		log.Sample(zerolog.Sometimes).Trace().Str("key", key).Msg("Key not found")
		return 0, false
	}
	if rawNewValue.Type() == pcommon.ValueTypeDouble {
		return rawNewValue.DoubleVal(), true
	} else if rawNewValue.Type() == pcommon.ValueTypeInt {
		return float64(rawNewValue.IntVal()), true
	} else if rawNewValue.Type() == pcommon.ValueTypeString {
		newValue, err := strconv.ParseFloat(rawNewValue.StringVal(), 64)
		if err != nil {
			for _, treatAsZeroValue := range treatAsZero {
				if rawNewValue.StringVal() == treatAsZeroValue {
					return 0, true
				}
			}
			log.Sample(zerolog.Sometimes).Warn().Str("key", key).Str("value", rawNewValue.AsString()).Msg("Failed parsing value as float")
			return 0, false
		}
		return newValue, true
	}
	log.Sample(zerolog.Sometimes).Warn().Str("key", key).Str("value", rawNewValue.AsString()).Msg("Unsupported value type")
	return 0, false
}

// Max returns the maximum value of the given values.
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum value of the given values.
func Min(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

// FormIncludeList returns a map of all the keys in the given list with a value of true.
func FormIncludeList(attributes []string) map[string]bool {
	return utils.SliceToMap(attributes)
}

// FormExcludeList returns a map of all the keys in the given list with a value of false.
func FormExcludeList(attributes []string) map[string]bool {
	return utils.SliceToMap(attributes)
}

type enforceCriteria uint8

const (
	include enforceCriteria = iota
	exclude
)

// EnforceIncludeList enforces the given include list on the given attributes.
func EnforceIncludeList(attributes pcommon.Map, includeList map[string]bool) {
	enforceList(attributes, includeList, include)
}

// EnforceExcludeList enforces the given exclude list on the given attributes.
func EnforceExcludeList(attributes pcommon.Map, excludeList map[string]bool) {
	enforceList(attributes, excludeList, exclude)
}

func enforceList(attributes pcommon.Map, list map[string]bool, enforceCriteria enforceCriteria) {
	keysToRemove := make([]string, 0)
	attributes.Range(func(key string, _ pcommon.Value) bool {
		if enforceCriteria == include && !list[key] {
			keysToRemove = append(keysToRemove, key)
		} else if enforceCriteria == exclude && list[key] {
			keysToRemove = append(keysToRemove, key)
		}
		return true
	})
	for _, key := range keysToRemove {
		attributes.Remove(key)
	}
}
