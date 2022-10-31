package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// AddEnvoySpecificLabels adds labels specific to Envoy data source.
func AddEnvoySpecificLabels(attributes pcommon.Map) {
	treatAsMissing := []string{otelcollector.EnvoyMissingAttributeValue}
	// Retrieve request length
	requestLength, requestLengthFound := otelcollector.GetFloat64(attributes, otelcollector.EnvoyBytesSentLabel, treatAsMissing)
	if requestLengthFound {
		attributes.PutDouble(otelcollector.HTTPRequestContentLength, requestLength)
	}
	// Retrieve response lengths
	responseLength, responseLengthFound := otelcollector.GetFloat64(attributes, otelcollector.EnvoyBytesReceivedLabel, treatAsMissing)
	if responseLengthFound {
		attributes.PutDouble(otelcollector.HTTPResponseContentLength, responseLength)
	}

	// Compute durations
	responseDuration, responseDurationExists := otelcollector.GetFloat64(attributes, otelcollector.EnvoyResponseDurationLabel, treatAsMissing)
	authzDuration, authzDurationExists := otelcollector.GetFloat64(attributes, otelcollector.EnvoyAuthzDurationLabel, treatAsMissing)

	// Add ResponseReceivedLabel based on whether responseDuration is present and non-zero
	if responseDurationExists && responseDuration > 0 {
		attributes.PutStr(otelcollector.ResponseReceivedLabel, otelcollector.ResponseReceivedTrue)
	} else {
		attributes.PutStr(otelcollector.ResponseReceivedLabel, otelcollector.ResponseReceivedFalse)
	}

	if responseDurationExists {
		attributes.PutDouble(otelcollector.FlowDurationLabel, responseDuration)
	}

	if responseDurationExists && authzDurationExists {
		workloadDuration := responseDuration - authzDuration
		// discard negative values which can happen in case of connection resets
		if workloadDuration > 0 {
			attributes.PutDouble(otelcollector.WorkloadDurationLabel, workloadDuration)
		}
	}
}
