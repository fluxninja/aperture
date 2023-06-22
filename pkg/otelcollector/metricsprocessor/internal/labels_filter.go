package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/v2/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/utils"
)

/*
 * IncludeList: This IncludeList is applied to logs and spans at during the enrichment process, after check response based labels are attached and metrics have been parsed.
 */
var (
	_includeAttributesCommon = []string{
		otelconsts.ApertureSourceLabel,
		otelconsts.WorkloadDurationLabel,
		otelconsts.FlowDurationLabel,
		otelconsts.ApertureProcessingDurationLabel,
		otelconsts.ApertureDecisionTypeLabel,
		otelconsts.ApertureRejectReasonLabel,
		otelconsts.ApertureRateLimitersLabel,
		otelconsts.ApertureDroppingRateLimitersLabel,
		otelconsts.ApertureLoadSchedulersLabel,
		otelconsts.ApertureDroppingLoadSchedulersLabel,
		otelconsts.ApertureWorkloadsLabel,
		otelconsts.ApertureDroppingWorkloadsLabel,
		otelconsts.ApertureFluxMetersLabel,
		otelconsts.ApertureFlowLabelKeysLabel,
		otelconsts.ApertureClassifiersLabel,
		otelconsts.ApertureClassifierErrorsLabel,
		otelconsts.ApertureServicesLabel,
		otelconsts.ApertureControlPointLabel,
		otelconsts.ApertureControlPointTypeLabel,
		otelconsts.ApertureFlowStatusLabel,
		otelconsts.ResponseReceivedLabel,
		otelconsts.ApertureSourceServiceLabel,
		otelconsts.ApertureDestinationServiceLabel,
	}

	_includeAttributesHTTP = []string{
		otelconsts.HTTPStatusCodeLabel,
		otelconsts.HTTPRequestContentLength,
		otelconsts.HTTPResponseContentLength,
	}

	_includeAttributesSDK = []string{
		otelconsts.ApertureFlowStatusLabel,
	}

	includeListHTTP = utils.SliceToSet(append(_includeAttributesCommon, _includeAttributesHTTP...))
	includeListSDK  = utils.SliceToSet(append(_includeAttributesCommon, _includeAttributesSDK...))
)

// EnforceIncludeListHTTP filters attributes for HTTP telemetry.
func EnforceIncludeListHTTP(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListHTTP)
}

// EnforceIncludeListSDK filters attributes for SDK telemetry.
func EnforceIncludeListSDK(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListSDK)
}
