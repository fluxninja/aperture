package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

// AddEnvoySpecificLabels adds labels specific to Envoy data source.
func AddEnvoySpecificLabels(attributes pcommon.Map) {
	treatAsMissing := []string{otelconsts.EnvoyMissingAttributeValue}
	// Retrieve request length
	requestLength, requestLengthFound := otelcollector.GetFloat64(attributes, otelconsts.EnvoyBytesSentLabel, treatAsMissing)
	if requestLengthFound {
		attributes.PutDouble(otelconsts.HTTPRequestContentLength, requestLength)
	}
	// Retrieve response lengths
	responseLength, responseLengthFound := otelcollector.GetFloat64(attributes, otelconsts.EnvoyBytesReceivedLabel, treatAsMissing)
	if responseLengthFound {
		attributes.PutDouble(otelconsts.HTTPResponseContentLength, responseLength)
	}

	// Compute durations
	responseDuration, responseDurationExists := otelcollector.GetFloat64(attributes, otelconsts.EnvoyResponseDurationLabel, treatAsMissing)
	authzDuration, authzDurationExists := otelcollector.GetFloat64(attributes, otelconsts.EnvoyAuthzDurationLabel, treatAsMissing)

	// Add ResponseReceivedLabel based on whether responseDuration is present and non-zero
	if responseDurationExists && responseDuration > 0 {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
	} else {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedFalse)
	}

	if responseDurationExists {
		attributes.PutDouble(otelconsts.FlowDurationLabel, responseDuration)
	}

	if responseDurationExists && authzDurationExists {
		workloadDuration := responseDuration - authzDuration
		// discard negative values which can happen in case of connection resets
		if workloadDuration > 0 {
			attributes.PutDouble(otelconsts.WorkloadDurationLabel, workloadDuration)
		}
	}
}
