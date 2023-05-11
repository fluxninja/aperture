package internal

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/v2/pkg/otelcollector"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

// AddLuaSpecificLabels adds labels specific to data source.
func AddLuaSpecificLabels(attributes pcommon.Map) {
	treatAsMissing := []string{otelconsts.LuaMissingAttributeValue}
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

	flowStart, flowStartExists := otelcollector.GetFloat64(attributes, otelconsts.ApertureFlowStartTimestampLabel, treatAsMissing)
	workloadStart, workloadStartExists := otelcollector.GetFloat64(attributes, otelconsts.ApertureWorkloadStartTimestampLabel, treatAsMissing)
	flowEnd, flowEndExists := otelcollector.GetFloat64(attributes, otelconsts.ApertureFlowEndTimestampLabel, treatAsMissing)

	// Add ResponseReceivedLabel based on whether flowEnd is present
	if flowEndExists {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedTrue)
	} else {
		attributes.PutStr(otelconsts.ResponseReceivedLabel, otelconsts.ResponseReceivedFalse)
	}

	if flowStartExists && flowEndExists {
		attributes.PutDouble(otelconsts.FlowDurationLabel, flowEnd-flowStart)
	}
	if workloadStartExists && flowEndExists {
		attributes.PutDouble(otelconsts.WorkloadDurationLabel, flowEnd-workloadStart)
	}
}
