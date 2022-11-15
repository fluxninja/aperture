package internal

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// StatusesFromAttributes gets HTTP status code and Feature status from attributes.
func StatusesFromAttributes(attributes pcommon.Map) (statusCode string, featureStatus string) {
	rawStatusCode, exists := attributes.Get(otelcollector.HTTPStatusCodeLabel)
	if exists {
		statusCode = rawStatusCode.Str()
	}
	rawFeatureStatus, exists := attributes.Get(otelcollector.ApertureFeatureStatusLabel)
	if exists {
		featureStatus = rawFeatureStatus.Str()
	}
	return
}

// StatusLabelsForMetrics returns labels maps used which describe Histogram.
func StatusLabelsForMetrics(
	decisionType flowcontrolv1.CheckResponse_DecisionType,
	statusCode string,
	featureStatus string,
) map[string]string {
	return map[string]string{
		metrics.ResponseStatusLabel: responseStatusForMetrics(statusCode, featureStatus),
		metrics.DecisionTypeLabel:   decisionType.String(),
		metrics.StatusCodeLabel:     statusCode,
		metrics.FeatureStatusLabel:  featureStatus,
	}
}

func responseStatusForMetrics(statusCode, featureStatus string) string {
	return responseStatus(
		statusCode,
		featureStatus,
		metrics.ResponseStatusOK,
		metrics.ResponseStatusError)
}

// ResponseStatusForTelemetry returns response status for telemetry based on
// HTTP status code and Feature status.
func ResponseStatusForTelemetry(statusCode, featureStatus string) string {
	return responseStatus(
		statusCode,
		featureStatus,
		otelcollector.ApertureResponseStatusOK,
		otelcollector.ApertureResponseStatusError)
}

func responseStatus(statusCode, featureStatus, okStatus, errorStatus string) string {
	// Checking status code this way instead of parsing int properly handles empty
	// string as well.
	if strings.HasPrefix(statusCode, "1") ||
		strings.HasPrefix(statusCode, "2") ||
		strings.HasPrefix(statusCode, "3") {
		return okStatus
	}
	if featureStatus != "" {
		return featureStatus
	}
	return errorStatus
}
