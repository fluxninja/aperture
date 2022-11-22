package otelcollector

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"google.golang.org/protobuf/proto"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/utils"
)

// IterAction describes whether to keep or discard an item processed by an iteration callback.
type IterAction bool

const (
	// Keep means keep this item and continue.
	Keep IterAction = true
	// Discard means remove this item and continue.
	Discard = false
)

// IterateLogRecords calls given function for each logRecord.
//
// The callback should return whether a log record should be kept or removed.
func IterateLogRecords(ld plog.Logs, fn func(plog.LogRecord) IterAction) {
	ld.ResourceLogs().RemoveIf(func(resourceLogs plog.ResourceLogs) bool {
		resourceLogs.ScopeLogs().RemoveIf(func(scopeLogs plog.ScopeLogs) bool {
			scopeLogs.LogRecords().RemoveIf(func(logRecord plog.LogRecord) bool {
				return fn(logRecord) == Discard //nolint:gosimple
			})
			return scopeLogs.LogRecords().Len() == 0
		})
		return resourceLogs.ScopeLogs().Len() == 0
	})
}

// IterateSpans calls given function for each span.
func IterateSpans(td ptrace.Traces, fn func(ptrace.Span)) {
	resourceSpansSlice := td.ResourceSpans()
	for resourceSpansIt := 0; resourceSpansIt < resourceSpansSlice.Len(); resourceSpansIt++ {
		resourceSpans := resourceSpansSlice.At(resourceSpansIt)
		scopeSpansSlice := resourceSpans.ScopeSpans()
		for scopeSpansIt := 0; scopeSpansIt < scopeSpansSlice.Len(); scopeSpansIt++ {
			scopeSpans := scopeSpansSlice.At(scopeSpansIt)
			spansSlice := scopeSpans.Spans()
			for spansIt := 0; spansIt < spansSlice.Len(); spansIt++ {
				span := spansSlice.At(spansIt)
				fn(span)
			}
		}
	}
}

// IterateMetrics calls given function for each metric.
func IterateMetrics(md pmetric.Metrics, fn func(pmetric.Metric)) {
	resourceMetricsSlice := md.ResourceMetrics()
	for resourceMetricsIt := 0; resourceMetricsIt < resourceMetricsSlice.Len(); resourceMetricsIt++ {
		resourceMetrics := resourceMetricsSlice.At(resourceMetricsIt)
		scopeMetricsSlice := resourceMetrics.ScopeMetrics()
		for scopeMetricsIt := 0; scopeMetricsIt < scopeMetricsSlice.Len(); scopeMetricsIt++ {
			scopeMetrics := scopeMetricsSlice.At(scopeMetricsIt)
			metricsSlice := scopeMetrics.Metrics()
			for metricsIt := 0; metricsIt < metricsSlice.Len(); metricsIt++ {
				span := metricsSlice.At(metricsIt)
				fn(span)
			}
		}
	}
}

// IterateDataPoints calls given function for each metric data point.
func IterateDataPoints(metric pmetric.Metric, fn func(pcommon.Map)) {
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		dataPoints := metric.Gauge().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			fn(dataPoints.At(i).Attributes())
		}
	case pmetric.MetricTypeSum:
		dataPoints := metric.Sum().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			fn(dataPoints.At(i).Attributes())
		}
	case pmetric.MetricTypeSummary:
		dataPoints := metric.Summary().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			fn(dataPoints.At(i).Attributes())
		}
	case pmetric.MetricTypeHistogram:
		dataPoints := metric.Histogram().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			fn(dataPoints.At(i).Attributes())
		}
	case pmetric.MetricTypeExponentialHistogram:
		dataPoints := metric.ExponentialHistogram().DataPoints()
		for i := 0; i < dataPoints.Len(); i++ {
			fn(dataPoints.At(i).Attributes())
		}
	}
}

// GetStruct is a helper for decoding complex structs encoded into an attribute
// as string.
//
// The attribute can be encoded as either:
// * JSON.
// * base64'd protobuf wire format, if `output` is `proto.Message`.
//
// Takes:
// attributes to read from
// label key to read in attributes
// output interface that is filled via json/proto unmarshal
// treatAsMissing is a list of values that are treated as attribute missing from source
//
// Returns true is label was decoded successfully, false otherwise.
func GetStruct(attributes pcommon.Map, label string, output interface{}, treatAsMissing []string) bool {
	value, ok := attributes.Get(label)
	if !ok {
		log.Sample(noAttrSampler).
			Warn().Str("label", label).Msg("Label does not exist in attributes map")
		return false
	}
	if value.Type() != pcommon.ValueTypeStr {
		log.Sample(notStringSampler).Warn().Str("label", label).Msg("Label is not a string")
		return false
	}

	stringVal := value.Str()

	for _, markerForMissing := range treatAsMissing {
		if stringVal == markerForMissing {
			log.Sample(missingAttrSampler).
				Info().Str("label", label).Msg("Missing attribute from source")
			return false
		}
	}

	if msg, isProto := output.(proto.Message); isProto && !strings.HasPrefix(stringVal, "{") {
		wireMsg, err := base64.StdEncoding.DecodeString(stringVal)
		if err != nil {
			log.Sample(failedBase64Sampler).
				Warn().Err(err).Str("label", label).Msg("Failed to unmarshal as base64")
		}
		err = proto.Unmarshal(wireMsg, msg)
		if err != nil {
			log.Sample(failedProtoSampler).
				Warn().Err(err).Str("label", label).Msg("Failed to unmarshal as protobuf")
		}
		return true
	}

	err := json.Unmarshal([]byte(stringVal), output)
	if err != nil {
		log.Sample(failedJSONSampler).
			Warn().Err(err).Str("label", label).Msg("Failed to unmarshal as json")
	}

	return true
}

// GetFloat64 returns float64 value from given attribute map at given key.
func GetFloat64(attributes pcommon.Map, key string, treatAsMissing []string) (float64, bool) {
	rawNewValue, exists := attributes.Get(key)
	if !exists {
		log.Trace().Str("key", key).Msg("Key not found")
		return 0, false
	}
	if rawNewValue.Type() == pcommon.ValueTypeDouble {
		return rawNewValue.Double(), true
	} else if rawNewValue.Type() == pcommon.ValueTypeInt {
		return float64(rawNewValue.Int()), true
	} else if rawNewValue.Type() == pcommon.ValueTypeStr {
		newValue, err := strconv.ParseFloat(rawNewValue.Str(), 64)
		if err != nil {
			for _, treatAsMissingValue := range treatAsMissing {
				if rawNewValue.Str() == treatAsMissingValue {
					return 0, false
				}
			}
			log.Sample(failedFloatSampler).
				Warn().Str("key", key).Str("value", rawNewValue.AsString()).
				Msg("Failed parsing value as float")
			return 0, false
		}
		return newValue, true
	}
	log.Sample(unsupportedFloatTypeSampler).
		Warn().Str("key", key).Str("value", rawNewValue.AsString()).Msg("Unsupported value type")
	return 0, false
}

var (
	noAttrSampler               = log.NewRatelimitingSampler()
	notStringSampler            = log.NewRatelimitingSampler()
	missingAttrSampler          = log.NewRatelimitingSampler()
	failedBase64Sampler         = log.NewRatelimitingSampler()
	failedProtoSampler          = log.NewRatelimitingSampler()
	failedJSONSampler           = log.NewRatelimitingSampler()
	failedFloatSampler          = log.NewRatelimitingSampler()
	unsupportedFloatTypeSampler = log.NewRatelimitingSampler()
)

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
