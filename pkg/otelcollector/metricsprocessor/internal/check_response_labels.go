package internal

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
)

// AddCheckResponseBasedLabels adds the following labels:
// * otelconsts.ApertureProcessingDurationLabel
// * otelconsts.ApertureServicesLabel
// * otelconsts.ApertureControlPointLabel
// * otelconsts.ApertureRateLimitersLabel
// * otelconsts.ApertureDroppingRateLimitersLabel
// * otelconsts.ApertureLoadSchedulersLabel
// * otelconsts.ApertureQuotaSchedulersLabel
// * otelconsts.ApertureDroppingQuotaSchedulersLabel
// * otelconsts.ApertureDroppingLoadSchedulersLabel
// * otelconsts.ApertureWorkloadsLabel
// * otelconsts.ApertureDroppingWorkloadsLabel
// * otelconsts.ApertureFluxMetersLabel
// * otelconsts.ApertureFlowLabelKeysLabel
// * otelconsts.ApertureClassifiersLabel
// * otelconsts.ApertureClassifierErrorsLabel
// * otelconsts.ApertureDecisionTypeLabel
// * otelconsts.ApertureRejectReasonLabel
// * otelconsts.ApertureResultCacheLookupStatusLabel
// * otelconsts.ApertureResultCacheOperationStatusLabel
// * otelconsts.ApertureGlobalCacheLookupStatusesLabel
// * otelconsts.ApertureGlobalCacheOperationStatusesLabel
// * dynamic flow labels.
func AddCheckResponseBasedLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse, sourceStr string) {
	// Aperture Processing Duration
	startTime := checkResponse.Start.AsTime()
	endTime := checkResponse.End.AsTime()
	if !startTime.IsZero() && !endTime.IsZero() {
		attributes.PutDouble(otelconsts.ApertureProcessingDurationLabel, float64(endTime.Sub(startTime).Milliseconds()))
	} else {
		log.Sample(noDurationSampler).
			Warn().Msgf("Aperture processing duration not found in %s access logs", sourceStr)
	}
	// Services
	putStrSlice(attributes, otelconsts.ApertureServicesLabel, checkResponse.Services)

	// Control Point
	attributes.PutStr(otelconsts.ApertureControlPointLabel, checkResponse.ControlPoint)

	// Decision type
	attributes.PutStr(
		otelconsts.ApertureDecisionTypeLabel,
		flowcontrolv1.CheckResponse_DecisionType_name[int32(checkResponse.DecisionType)],
	)

	// Reason
	attributes.PutStr(
		otelconsts.ApertureRejectReasonLabel,
		flowcontrolv1.CheckResponse_RejectReason_name[int32(checkResponse.RejectReason)],
	)

	// Cache
	if checkResponse.CacheLookupResponse != nil {
		if checkResponse.CacheLookupResponse.ResultCacheResponse != nil {
			resultCacheResponse := checkResponse.CacheLookupResponse.ResultCacheResponse
			attributes.PutStr(otelconsts.ApertureResultCacheLookupStatusLabel, resultCacheResponse.LookupStatus.String())
			attributes.PutStr(otelconsts.ApertureResultCacheOperationStatusLabel, resultCacheResponse.OperationStatus.String())
		}
		globalCacheLookupStatuses := attributes.PutEmptySlice(otelconsts.ApertureGlobalCacheLookupStatusesLabel)
		globalCacheOperationStatuses := attributes.PutEmptySlice(otelconsts.ApertureGlobalCacheOperationStatusesLabel)
		for _, globalCacheResponse := range checkResponse.CacheLookupResponse.GlobalCacheResponses {
			if globalCacheResponse == nil {
				continue
			}
			globalCacheLookupStatuses.AppendEmpty().SetStr(globalCacheResponse.LookupStatus.String())
			globalCacheOperationStatuses.AppendEmpty().SetStr(globalCacheResponse.OperationStatus.String())
		}
	}

	// Note: Sorted alphabetically to help sorting attributes in rollupprocessor.key at least a bit.
	droppingLoadSchedulersSlice := attributes.PutEmptySlice(otelconsts.ApertureDroppingLoadSchedulersLabel)
	droppingQuotaSchedulersSlice := attributes.PutEmptySlice(otelconsts.ApertureDroppingQuotaSchedulersLabel)
	droppingRateLimitersSlice := attributes.PutEmptySlice(otelconsts.ApertureDroppingRateLimitersLabel)
	droppingSamplersSlice := attributes.PutEmptySlice(otelconsts.ApertureDroppingSamplersLabel)
	droppingWorkloadsSlice := attributes.PutEmptySlice(otelconsts.ApertureDroppingWorkloadsLabel)
	loadSchedulersSlice := attributes.PutEmptySlice(otelconsts.ApertureLoadSchedulersLabel)
	quotaSchedulersSlice := attributes.PutEmptySlice(otelconsts.ApertureQuotaSchedulersLabel)
	rateLimitersSlice := attributes.PutEmptySlice(otelconsts.ApertureRateLimitersLabel)
	samplersSlice := attributes.PutEmptySlice(otelconsts.ApertureSamplersLabel)
	workloadsSlice := attributes.PutEmptySlice(otelconsts.ApertureWorkloadsLabel)

	for _, decision := range checkResponse.LimiterDecisions {
		if decision.GetRateLimiterInfo() != nil {
			value := fmt.Sprintf(
				"%s:%v,%s:%v,%s:%v",
				metrics.PolicyNameLabel, decision.GetPolicyName(),
				metrics.ComponentIDLabel, decision.GetComponentId(),
				metrics.PolicyHashLabel, decision.GetPolicyHash(),
			)
			rateLimitersSlice.AppendEmpty().SetStr(value)
			if decision.Dropped {
				droppingRateLimitersSlice.AppendEmpty().SetStr(value)
			}
		}
		if cl := decision.GetLoadSchedulerInfo(); cl != nil {
			value := fmt.Sprintf(
				"%s:%v,%s:%v,%s:%v",
				metrics.PolicyNameLabel, decision.GetPolicyName(),
				metrics.ComponentIDLabel, decision.GetComponentId(),
				metrics.PolicyHashLabel, decision.GetPolicyHash(),
			)
			loadSchedulersSlice.AppendEmpty().SetStr(value)
			if decision.Dropped {
				droppingLoadSchedulersSlice.AppendEmpty().SetStr(value)
			}

			workloadsValue := fmt.Sprintf(
				"%s:%v,%s:%v,%s:%v,%s:%v",
				metrics.PolicyNameLabel, decision.GetPolicyName(),
				metrics.ComponentIDLabel, decision.GetComponentId(),
				metrics.WorkloadIndexLabel, cl.GetWorkloadIndex(),
				metrics.PolicyHashLabel, decision.GetPolicyHash(),
			)
			workloadsSlice.AppendEmpty().SetStr(workloadsValue)
			if decision.Dropped {
				droppingWorkloadsSlice.AppendEmpty().SetStr(workloadsValue)
			}
		}
		if decision.GetSamplerInfo() != nil {
			value := fmt.Sprintf(
				"%s:%v,%s:%v,%s:%v",
				metrics.PolicyNameLabel, decision.GetPolicyName(),
				metrics.ComponentIDLabel, decision.GetComponentId(),
				metrics.PolicyHashLabel, decision.GetPolicyHash(),
			)
			samplersSlice.AppendEmpty().SetStr(value)
			if decision.Dropped {
				droppingSamplersSlice.AppendEmpty().SetStr(value)
			}
		}
		if cl := decision.GetQuotaSchedulerInfo(); cl != nil {
			value := fmt.Sprintf(
				"%s:%v,%s:%v,%s:%v",
				metrics.PolicyNameLabel, decision.GetPolicyName(),
				metrics.ComponentIDLabel, decision.GetComponentId(),
				metrics.PolicyHashLabel, decision.GetPolicyHash(),
			)
			quotaSchedulersSlice.AppendEmpty().SetStr(value)
			if decision.Dropped {
				droppingQuotaSchedulersSlice.AppendEmpty().SetStr(value)
			}

			workloadsValue := fmt.Sprintf(
				"%s:%v,%s:%v,%s:%v,%s:%v",
				metrics.PolicyNameLabel, decision.GetPolicyName(),
				metrics.ComponentIDLabel, decision.GetComponentId(),
				metrics.WorkloadIndexLabel, cl.GetWorkloadIndex(),
				metrics.PolicyHashLabel, decision.GetPolicyHash(),
			)
			workloadsSlice.AppendEmpty().SetStr(workloadsValue)
			if decision.Dropped {
				droppingWorkloadsSlice.AppendEmpty().SetStr(workloadsValue)
			}
		}
	}

	fluxMetersSlice := attributes.PutEmptySlice(otelconsts.ApertureFluxMetersLabel)
	fluxMetersSlice.EnsureCapacity(len(checkResponse.FluxMeterInfos))
	for _, fluxMeter := range checkResponse.FluxMeterInfos {
		fluxMetersSlice.AppendEmpty().SetStr(fluxMeter.GetFluxMeterName())
	}

	putStrSlice(attributes, otelconsts.ApertureFlowLabelKeysLabel, checkResponse.FlowLabelKeys)

	classifiersSlice := attributes.PutEmptySlice(otelconsts.ApertureClassifiersLabel)
	classifiersSlice.EnsureCapacity(len(checkResponse.ClassifierInfos))
	classifierErrorsSlice := attributes.PutEmptySlice(otelconsts.ApertureClassifierErrorsLabel)
	for _, classifier := range checkResponse.ClassifierInfos {
		value := fmt.Sprintf(
			"%s:%v,%s:%v",
			metrics.PolicyNameLabel, classifier.PolicyName,
			metrics.ClassifierIndexLabel, classifier.ClassifierIndex,
		)
		classifiersSlice.AppendEmpty().SetStr(value)

		// add errors as attributes as well
		if classifier.Error != flowcontrolv1.ClassifierInfo_ERROR_NONE {
			errorsValue := fmt.Sprintf(
				"%s,%s:%v,%s:%v,%s:%v",
				classifier.Error.String(),
				metrics.PolicyNameLabel, classifier.PolicyName,
				metrics.ClassifierIndexLabel, classifier.ClassifierIndex,
				metrics.PolicyHashLabel, classifier.PolicyHash,
			)
			classifierErrorsSlice.AppendEmpty().SetStr(errorsValue)
		}
	}
}

var noDurationSampler = log.NewRatelimitingSampler()

// AddFlowLabels adds dynamic from labels.
func AddFlowLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse) {
	for key, value := range checkResponse.GetTelemetryFlowLabels() {
		pcommon.NewValueStr(value).CopyTo(attributes.PutEmpty(key))
	}
}

func putStrSlice(m pcommon.Map, key string, strs []string) {
	slice := m.PutEmptySlice(key)
	slice.EnsureCapacity(len(strs))
	for _, str := range strs {
		slice.AppendEmpty().SetStr(str)
	}
}
