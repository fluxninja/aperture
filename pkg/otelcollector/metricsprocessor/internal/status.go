package internal

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/pkg/metrics"
	otelconsts "github.com/fluxninja/aperture/pkg/otelcollector/consts"
)

// StatusesFromAttributes gets HTTP status code and Flow status from attributes.
func StatusesFromAttributes(attributes pcommon.Map) (statusCode string, flowStatus string) {
	rawStatusCode, exists := attributes.Get(otelconsts.HTTPStatusCodeLabel)
	if exists {
		statusCode = rawStatusCode.Str()
	}
	rawFlowStatus, exists := attributes.Get(otelconsts.ApertureFlowStatusLabel)
	if exists {
		flowStatus = rawFlowStatus.Str()
	}
	return
}

// StatusLabelsForMetrics returns labels maps used which describe Histogram.
func StatusLabelsForMetrics(
	decisionType flowcontrolv1.CheckResponse_DecisionType,
	statusCode string,
	flowStatus string,
) map[string]string {
	return map[string]string{
		metrics.DecisionTypeLabel: decisionType.String(),
		metrics.StatusCodeLabel:   statusCode,
		metrics.FlowStatusLabel:   flowStatusForMetrics(statusCode, flowStatus),
	}
}

func flowStatusForMetrics(statusCode, flowStatusStr string) string {
	return flowStatus(
		statusCode,
		flowStatusStr,
		metrics.FlowStatusOK,
		metrics.FlowStatusError)
}

// FlowStatusForTelemetry returns protocol independent Flow status for telemetry based on
// HTTP status code and Flow status.
func FlowStatusForTelemetry(statusCode, flowStatusStr string) string {
	return flowStatus(
		statusCode,
		flowStatusStr,
		otelconsts.ApertureFlowStatusOK,
		otelconsts.ApertureFlowStatusError)
}

func flowStatus(statusCode, flowStatusStr, okStatus, errorStatus string) string {
	// Checking status code this way instead of parsing int properly handles empty
	// string as well.
	if strings.HasPrefix(statusCode, "1") ||
		strings.HasPrefix(statusCode, "2") ||
		strings.HasPrefix(statusCode, "3") {
		return okStatus
	}
	if flowStatusStr != "" {
		return flowStatusStr
	}
	return errorStatus
}
