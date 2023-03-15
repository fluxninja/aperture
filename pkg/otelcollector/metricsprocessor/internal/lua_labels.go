package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

// AddLuaSpecificLabels adds labels specific to data source.
func AddLuaSpecificLabels(attributes pcommon.Map) {
	treatAsMissing := []string{""}
	// Retrieve request length
	requestLength, requestLengthFound := otelcollector.GetFloat64(attributes, otelconsts.BytesSentLabel, treatAsMissing)
	if requestLengthFound {
		attributes.PutDouble(otelconsts.HTTPRequestContentLength, requestLength)
	}
	// Retrieve response lengths
	responseLength, responseLengthFound := otelcollector.GetFloat64(attributes, otelconsts.BytesReceivedLabel, treatAsMissing)
	if responseLengthFound {
		attributes.PutDouble(otelconsts.HTTPResponseContentLength, responseLength)
	}

	// Compute durations
	responseDuration, responseDurationExists := otelcollector.GetFloat64(attributes, otelconsts.ResponseDurationLabel, treatAsMissing)
	checkHTTPDuration, checkHTTPDurationExists := otelcollector.GetFloat64(attributes, otelconsts.CheckHTTPDurationLabel, treatAsMissing)

	// Add ResponseReceivedLabel based on whether responseDuration is present and non-zero
	if responseDurationExists && responseDuration > 0 {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
	} else {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedFalse)
	}

	if responseDurationExists {
		attributes.PutDouble(otelconsts.FlowDurationLabel, responseDuration)
	}

	if responseDurationExists && checkHTTPDurationExists {
		workloadDuration := responseDuration - checkHTTPDuration
		// discard negative values which can happen in case of connection resets
		if workloadDuration > 0 {
			attributes.PutDouble(otelconsts.WorkloadDurationLabel, workloadDuration)
		}
	}

	log.Error().Msgf("Updated attributes '%+v': ", attributes)
}
