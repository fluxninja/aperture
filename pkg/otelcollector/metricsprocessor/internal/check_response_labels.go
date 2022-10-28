package internal

import (
	"fmt"
	"strings"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

// AddCheckResponseBasedLabels adds the following labels:
// * otelcollector.ApertureProcessingDurationLabel
// * otelcollector.ApertureServicesLabel
// * otelcollector.ApertureControlPointLabel
// * otelcollector.ApertureRateLimitersLabel
// * otelcollector.ApertureDroppingRateLimitersLabel
// * otelcollector.ApertureConcurrencyLimitersLabel
// * otelcollector.ApertureDroppingConcurrencyLimitersLabel
// * otelcollector.ApertureWorkloadsLabel
// * otelcollector.ApertureDroppingWorkloadsLabel
// * otelcollector.ApertureFluxMetersLabel
// * otelcollector.ApertureFlowLabelKeysLabel
// * otelcollector.ApertureClassifiersLabel
// * otelcollector.ApertureClassifierErrorsLabel
// * otelcollector.ApertureDecisionTypeLabel
// * otelcollector.ApertureRejectReasonLabel
// * otelcollector.ApertureErrorLabel,
// * dynamic flow labels.
func AddCheckResponseBasedLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse, sourceStr string) {
	// Aperture Processing Duration
	startTime := checkResponse.GetStart().AsTime()
	endTime := checkResponse.GetEnd().AsTime()
	if !startTime.IsZero() && !endTime.IsZero() {
		attributes.PutDouble(otelcollector.ApertureProcessingDurationLabel, float64(endTime.Sub(startTime).Milliseconds()))
	} else {
		log.Sample(zerolog.Sometimes).Warn().Msgf("Aperture processing duration not found in %s access logs", sourceStr)
	}
	// Services
	servicesValue := pcommon.NewValueSlice()
	for _, service := range checkResponse.Services {
		servicesValue.Slice().AppendEmpty().SetStr(service)
	}
	servicesValue.CopyTo(attributes.PutEmpty(otelcollector.ApertureServicesLabel))

	// Control Point
	attributes.PutStr(otelcollector.ApertureControlPointLabel, checkResponse.GetControlPointInfo().String())

	labels := map[string]pcommon.Value{
		otelcollector.ApertureRateLimitersLabel:                pcommon.NewValueSlice(),
		otelcollector.ApertureDroppingRateLimitersLabel:        pcommon.NewValueSlice(),
		otelcollector.ApertureConcurrencyLimitersLabel:         pcommon.NewValueSlice(),
		otelcollector.ApertureDroppingConcurrencyLimitersLabel: pcommon.NewValueSlice(),
		otelcollector.ApertureWorkloadsLabel:                   pcommon.NewValueSlice(),
		otelcollector.ApertureDroppingWorkloadsLabel:           pcommon.NewValueSlice(),
		otelcollector.ApertureFluxMetersLabel:                  pcommon.NewValueSlice(),
		otelcollector.ApertureFlowLabelKeysLabel:               pcommon.NewValueSlice(),
		otelcollector.ApertureClassifiersLabel:                 pcommon.NewValueSlice(),
		otelcollector.ApertureClassifierErrorsLabel:            pcommon.NewValueSlice(),
		otelcollector.ApertureDecisionTypeLabel:                pcommon.NewValueStr(checkResponse.DecisionType.String()),
		otelcollector.ApertureRejectReasonLabel:                pcommon.NewValueStr(checkResponse.GetRejectReason().String()),
		otelcollector.ApertureErrorLabel:                       pcommon.NewValueStr(checkResponse.GetError().String()),
	}
	for _, decision := range checkResponse.LimiterDecisions {
		if decision.GetRateLimiterInfo() != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.ApertureRateLimitersLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelcollector.ApertureDroppingRateLimitersLabel].Slice().AppendEmpty().SetStr(value)
			}
		}
		if cl := decision.GetConcurrencyLimiterInfo(); cl != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.ApertureConcurrencyLimitersLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelcollector.ApertureDroppingConcurrencyLimitersLabel].Slice().AppendEmpty().SetStr(value)
			}

			workloadsRawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.WorkloadIndexLabel, cl.GetWorkloadIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value = strings.Join(workloadsRawValue, ",")
			labels[otelcollector.ApertureWorkloadsLabel].Slice().AppendEmpty().SetStr(value)
			if decision.Dropped {
				labels[otelcollector.ApertureDroppingWorkloadsLabel].Slice().AppendEmpty().SetStr(value)
			}
		}
	}
	for _, fluxMeter := range checkResponse.FluxMeterInfos {
		value := fluxMeter.GetFluxMeterName()
		labels[otelcollector.ApertureFluxMetersLabel].Slice().AppendEmpty().SetStr(value)
	}

	for _, flowLabelKey := range checkResponse.GetFlowLabelKeys() {
		labels[otelcollector.ApertureFlowLabelKeysLabel].Slice().AppendEmpty().SetStr(flowLabelKey)
	}

	for key, value := range checkResponse.GetTelemetryFlowLabels() {
		pcommon.NewValueStr(value).CopyTo(attributes.PutEmpty(key))
	}

	for _, classifier := range checkResponse.ClassifierInfos {
		rawValue := []string{
			fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, classifier.PolicyName),
			fmt.Sprintf("%s:%v", metrics.ClassifierIndexLabel, classifier.ClassifierIndex),
		}
		value := strings.Join(rawValue, ",")
		labels[otelcollector.ApertureClassifiersLabel].Slice().AppendEmpty().SetStr(value)

		// add errors as attributes as well
		if classifier.Error != flowcontrolv1.ClassifierInfo_ERROR_NONE {
			errorsValue := []string{
				classifier.Error.String(),
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, classifier.PolicyName),
				fmt.Sprintf("%s:%v", metrics.ClassifierIndexLabel, classifier.ClassifierIndex),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, classifier.PolicyHash),
			}
			joinedValue := strings.Join(errorsValue, ",")
			labels[otelcollector.ApertureClassifierErrorsLabel].Slice().AppendEmpty().SetStr(joinedValue)
		}
	}

	for key, value := range labels {
		value.CopyTo(attributes.PutEmpty(key))
	}
}
