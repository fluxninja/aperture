package internal

import (
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// AddSDKSpecificLabels adds labels specific to SDK data source.
func AddSDKSpecificLabels(attributes pcommon.Map) {
	// Compute durations
	flowStart, flowStartExists := getSDKLabelTimestampValue(attributes, otelcollector.ApertureFlowStartTimestampLabel)
	workloadStart, workloadStartExists := getSDKLabelTimestampValue(attributes, otelcollector.ApertureWorkloadStartTimestampLabel)
	flowEnd, flowEndExists := getSDKLabelTimestampValue(attributes, otelcollector.ApertureFlowEndTimestampLabel)

	if flowStartExists && flowEndExists {
		flowDuration := flowEnd.Sub(flowStart)
		attributes.PutDouble(otelcollector.FlowDurationLabel, float64(flowDuration.Milliseconds()))
	}
	if workloadStartExists && flowEndExists {
		workloadDuration := flowEnd.Sub(workloadStart)
		attributes.PutDouble(otelcollector.WorkloadDurationLabel, float64(workloadDuration.Milliseconds()))
	}
}

func getSDKLabelTimestampValue(attributes pcommon.Map, labelKey string) (time.Time, bool) {
	return getLabelTimestampValue(attributes, labelKey, otelcollector.ApertureSourceSDK)
}

func getLabelTimestampValue(attributes pcommon.Map, labelKey, source string) (time.Time, bool) {
	value, exists := getLabelValue(attributes, labelKey, source)
	if !exists {
		return time.Time{}, false
	}
	return _getLabelTimestampValue(value, labelKey, source)
}

func getLabelValue(attributes pcommon.Map, labelKey, source string) (pcommon.Value, bool) {
	value, exists := attributes.Get(labelKey)
	if !exists {
		log.Sample(zerolog.Sometimes).Warn().Str("source", source).Str("key", labelKey).Msg("Label not found")
		return pcommon.Value{}, false
	}
	return value, exists
}

func _getLabelTimestampValue(value pcommon.Value, labelKey, source string) (time.Time, bool) {
	var valueInt int64
	if value.Type() == pcommon.ValueTypeInt {
		valueInt = value.Int()
	} else {
		log.Sample(zerolog.Sometimes).Warn().Str("source", source).Str("key", labelKey).Msg("Failed to parse a timestamp field")
		return time.Time{}, false
	}

	return time.Unix(0, valueInt), true
}
