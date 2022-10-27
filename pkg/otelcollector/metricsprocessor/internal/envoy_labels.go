package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// AddEnvoySpecificLabels adds labels specific to Envoy data source.
func AddEnvoySpecificLabels(attributes pcommon.Map) {
	treatAsZero := []string{otelcollector.EnvoyMissingAttributeValue}
	// Retrieve request length
	requestLength, _ := otelcollector.GetFloat64(attributes, otelcollector.EnvoyBytesSentLabel, treatAsZero)
	attributes.PutDouble(otelcollector.HTTPRequestContentLength, requestLength)
	// Retrieve response lengths
	responseLength, _ := otelcollector.GetFloat64(attributes, otelcollector.EnvoyBytesReceivedLabel, treatAsZero)
	attributes.PutDouble(otelcollector.HTTPResponseContentLength, responseLength)

	// Compute durations
	responseDuration, responseDurationExists := otelcollector.GetFloat64(attributes, otelcollector.EnvoyResponseDurationLabel, treatAsZero)
	authzDuration, authzDurationExists := otelcollector.GetFloat64(attributes, otelcollector.EnvoyAuthzDurationLabel, treatAsZero)

	if responseDurationExists {
		attributes.PutDouble(otelcollector.FlowDurationLabel, responseDuration)
	}

	if responseDurationExists && authzDurationExists {
		attributes.PutDouble(otelcollector.WorkloadDurationLabel, responseDuration-authzDuration)
	}
}
