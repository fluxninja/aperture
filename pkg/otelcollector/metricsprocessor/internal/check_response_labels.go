package internal

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
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
// * otelconsts.ApertureDroppingLoadSchedulersLabel
// * otelconsts.ApertureWorkloadsLabel
// * otelconsts.ApertureDroppingWorkloadsLabel
// * otelconsts.ApertureFluxMetersLabel
// * otelconsts.ApertureFlowLabelKeysLabel
// * otelconsts.ApertureClassifiersLabel
// * otelconsts.ApertureClassifierErrorsLabel
// * otelconsts.ApertureDecisionTypeLabel
// * otelconsts.ApertureRejectReasonLabel
// * dynamic flow labels.
func AddCheckResponseBasedLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse, sourceStr string) {
	// Aperture Processing Duration
	startTime := checkResponse.GetStart().AsTime()
	endTime := checkResponse.GetEnd().AsTime()
	if !startTime.IsZero() && !endTime.IsZero() {
		attributes.PutDouble(otelconsts.ApertureProcessingDurationLabel, float64(endTime.Sub(startTime).Milliseconds()))
	} else {
		log.Sample(noDurationSampler).
			Warn().Msgf("Aperture processing duration not found in %s access logs", sourceStr)
	}
	// Services
	servicesValue := pcommon.NewValueSlice()
	for _, service := range checkResponse.Services {
		servicesValue.Slice().AppendEmpty().SetStr(service)
	}
	servicesValue.CopyTo(attributes.PutEmpty(otelconsts.ApertureServicesLabel))

	// Control Point
	attributes.PutStr(otelconsts.ApertureControlPointLabel, checkResponse.GetControlPoint())

	labels := map[string]pcommon.Value{
		otelconsts.ApertureRateLimitersLabel:           pcommon.NewValueSlice(),
		otelconsts.ApertureDroppingRateLimitersLabel:   pcommon.NewValueSlice(),
		otelconsts.ApertureLoadSchedulersLabel:         pcommon.NewValueSlice(),
		otelconsts.ApertureDroppingLoadSchedulersLabel: pcommon.NewValueSlice(),
		otelconsts.ApertureWorkloadsLabel:              pcommon.NewValueSlice(),
		otelconsts.ApertureDroppingWorkloadsLabel:      pcommon.NewValueSlice(),
		otelconsts.ApertureFluxMetersLabel:             pcommon.NewValueSlice(),
		otelconsts.ApertureFlowLabelKeysLabel:          pcommon.NewValueSlice(),
		otelconsts.ApertureClassifiersLabel:            pcommon.NewValueSlice(),
		otelconsts.ApertureClassifierErrorsLabel:       pcommon.NewValueSlice(),
		otelconsts.ApertureDecisionTypeLabel:           pcommon.NewValueStr(checkResponse.DecisionType.String()),
		otelconsts.ApertureRejectReasonLabel:           pcommon.NewValueStr(checkResponse.GetRejectReason().String()),
	}
	for _, decision := range checkResponse.LimiterDecisions {
		if decision.GetRateLimiterInfo() != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIDLabel, decision.GetComponentId()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelconsts.ApertureRateLimitersLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelconsts.ApertureDroppingRateLimitersLabel].Slice().AppendEmpty().SetStr(value)
			}
		}
		if cl := decision.GetLoadSchedulerInfo(); cl != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIDLabel, decision.GetComponentId()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelconsts.ApertureLoadSchedulersLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelconsts.ApertureDroppingLoadSchedulersLabel].Slice().AppendEmpty().SetStr(value)
			}

			workloadsRawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIDLabel, decision.GetComponentId()),
				fmt.Sprintf("%s:%v", metrics.WorkloadIndexLabel, cl.GetWorkloadIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value = strings.Join(workloadsRawValue, ",")
			labels[otelconsts.ApertureWorkloadsLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelconsts.ApertureDroppingWorkloadsLabel].Slice().AppendEmpty().SetStr(value)
			}
		}
		if decision.GetSamplerInfo() != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIDLabel, decision.GetComponentId()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelconsts.ApertureRateLimitersLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelconsts.ApertureDroppingRateLimitersLabel].Slice().AppendEmpty().SetStr(value)
			}
		}
	}
	for _, fluxMeter := range checkResponse.FluxMeterInfos {
		value := fluxMeter.GetFluxMeterName()
		labels[otelconsts.ApertureFluxMetersLabel].Slice().AppendEmpty().SetStr(value)
	}

	for _, flowLabelKey := range checkResponse.GetFlowLabelKeys() {
		labels[otelconsts.ApertureFlowLabelKeysLabel].Slice().AppendEmpty().SetStr(flowLabelKey)
	}

	for _, classifier := range checkResponse.ClassifierInfos {
		rawValue := []string{
			fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, classifier.PolicyName),
			fmt.Sprintf("%s:%v", metrics.ClassifierIndexLabel, classifier.ClassifierIndex),
		}
		value := strings.Join(rawValue, ",")
		labels[otelconsts.ApertureClassifiersLabel].Slice().AppendEmpty().SetStr(value)

		// add errors as attributes as well
		if classifier.Error != flowcontrolv1.ClassifierInfo_ERROR_NONE {
			errorsValue := []string{
				classifier.Error.String(),
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, classifier.PolicyName),
				fmt.Sprintf("%s:%v", metrics.ClassifierIndexLabel, classifier.ClassifierIndex),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, classifier.PolicyHash),
			}
			joinedValue := strings.Join(errorsValue, ",")
			labels[otelconsts.ApertureClassifierErrorsLabel].Slice().AppendEmpty().SetStr(joinedValue)
		}
	}

	for key, value := range labels {
		value.CopyTo(attributes.PutEmpty(key))
	}
}

var noDurationSampler = log.NewRatelimitingSampler()

// AddFlowLabels adds dynamic from labels.
func AddFlowLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse) {
	for key, value := range checkResponse.GetTelemetryFlowLabels() {
		pcommon.NewValueStr(value).CopyTo(attributes.PutEmpty(key))
	}
}
