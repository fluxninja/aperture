package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/otelcollector"
)

/*
 * IncludeList: This IncludeList is applied to logs and spans at during the enrichment process, after check response based labels are attached and metrics have been parsed.
 */
var (
	_includeAttributesCommon = []string{
		otelcollector.ApertureSourceLabel,
		otelcollector.WorkloadDurationLabel,
		otelcollector.FlowDurationLabel,
		otelcollector.ApertureProcessingDurationLabel,
		otelcollector.ApertureDecisionTypeLabel,
		otelcollector.ApertureErrorLabel,
		otelcollector.ApertureRejectReasonLabel,
		otelcollector.ApertureRateLimitersLabel,
		otelcollector.ApertureDroppingRateLimitersLabel,
		otelcollector.ApertureConcurrencyLimitersLabel,
		otelcollector.ApertureDroppingConcurrencyLimitersLabel,
		otelcollector.ApertureWorkloadsLabel,
		otelcollector.ApertureDroppingWorkloadsLabel,
		otelcollector.ApertureFluxMetersLabel,
		otelcollector.ApertureFlowLabelKeysLabel,
		otelcollector.ApertureClassifiersLabel,
		otelcollector.ApertureClassifierErrorsLabel,
		otelcollector.ApertureServicesLabel,
		otelcollector.ApertureControlPointLabel,
		otelcollector.ApertureResponseStatusLabel,
		otelcollector.ResponseReceivedLabel,
	}

	_includeAttributesHTTP = []string{
		otelcollector.HTTPStatusCodeLabel,
		otelcollector.HTTPRequestContentLength,
		otelcollector.HTTPResponseContentLength,
	}

	_includeAttributesSDK = []string{
		otelcollector.ApertureFeatureStatusLabel,
	}

	includeListHTTP = otelcollector.FormIncludeList(append(_includeAttributesCommon, _includeAttributesHTTP...))
	includeListSDK  = otelcollector.FormIncludeList(append(_includeAttributesCommon, _includeAttributesSDK...))
)

// EnforceIncludeListHTTP filters attributes for HTTP telemetry.
func EnforceIncludeListHTTP(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListHTTP)
}

// EnforceIncludeListSDK filters attributes for SDK telemetry.
func EnforceIncludeListSDK(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListSDK)
}
