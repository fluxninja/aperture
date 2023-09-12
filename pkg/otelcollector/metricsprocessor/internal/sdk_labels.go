package internal

import (
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/v2/pkg/log"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

// AddSDKSpecificLabels adds labels specific to SDK data source.
func AddSDKSpecificLabels(attributes pcommon.Map) {
	// Compute durations
	flowStart, flowStartExists := getSDKLabelTimestampValue(attributes, otelconsts.ApertureFlowStartTimestampLabelMs)
	if !flowStartExists {
		flowStart, flowStartExists = getSDKLabelTimestampValue(attributes, otelconsts.ApertureFlowStartTimestampLabel)
	}

	workloadStart, workloadStartExists := getSDKLabelTimestampValue(attributes, otelconsts.ApertureWorkloadStartTimestampLabelMs)
	if !workloadStartExists {
		workloadStart, workloadStartExists = getSDKLabelTimestampValue(attributes, otelconsts.ApertureWorkloadStartTimestampLabel)
	}

	flowEnd, flowEndExists := getSDKLabelTimestampValue(attributes, otelconsts.ApertureFlowEndTimestampLabelMs)
	if !flowEndExists {
		flowEnd, flowEndExists = getSDKLabelTimestampValue(attributes, otelconsts.ApertureFlowEndTimestampLabel)
	}

	// Add ResponseReceivedLabel based on whether flowEnd is present
	if flowEndExists {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
	} else {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedFalse)
	}

	if flowStartExists && flowEndExists {
		flowDuration := flowEnd.Sub(flowStart)
		attributes.PutDouble(otelconsts.FlowDurationLabel, float64(flowDuration.Milliseconds()))
	}
	if workloadStartExists && flowEndExists {
		workloadDuration := flowEnd.Sub(workloadStart)
		attributes.PutDouble(otelconsts.WorkloadDurationLabel, float64(workloadDuration.Milliseconds()))
	}
}

func getSDKLabelTimestampValue(attributes pcommon.Map, labelKey string) (time.Time, bool) {
	return getLabelTimestampValue(attributes, labelKey, otelconsts.ApertureSourceSDK)
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
		log.Sample(noLabelSampler).
			Warn().Str("source", source).Str("key", labelKey).Msg("Label not found")
		return pcommon.Value{}, false
	}
	return value, exists
}

func _getLabelTimestampValue(value pcommon.Value, labelKey, source string) (time.Time, bool) {
	var valueInt int64
	if value.Type() == pcommon.ValueTypeInt {
		valueInt = value.Int()
	} else {
		log.Sample(badTimestampSampler).
			Warn().Str("source", source).Str("key", labelKey).
			Msg("Failed to parse a timestamp field")
		return time.Time{}, false
	}

	if strings.HasSuffix(labelKey, "_ms") {
		return time.UnixMilli(valueInt), true
	} else {
		return time.Unix(0, valueInt), true
	}
}

var (
	noLabelSampler      = log.NewRatelimitingSampler()
	badTimestampSampler = log.NewRatelimitingSampler()
)
