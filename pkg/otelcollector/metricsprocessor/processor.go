package metricsprocessor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/rs/zerolog"
)

type metricsProcessor struct {
	cfg *Config
}

func newProcessor(cfg *Config) (*metricsProcessor, error) {
	p := &metricsProcessor{
		cfg: cfg,
	}

	return p, nil
}

// Start indicates and logs the start of the metrics processor.
func (p *metricsProcessor) Start(_ context.Context, _ component.Host) error {
	log.Debug().Msg("metrics processor start")
	return nil
}

// Shutdown indicates and logs the shutdown of the metrics processor.
func (p *metricsProcessor) Shutdown(context.Context) error {
	log.Debug().Msg("metrics processor shutdown")
	return nil
}

// Capabilities returns the capabilities of the processor with MutatesData set to true.
func (p *metricsProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

// ConsumeLogs receives plog.Logs for consumption then returns updated logs with policy labels and metrics.
func (p *metricsProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	err := otelcollector.IterateLogRecords(ld, func(logRecord plog.LogRecord) error {
		retErr := func(errMsg string) error {
			log.Sample(zerolog.Sometimes).Warn().Msg(errMsg)
			return fmt.Errorf(errMsg)
		}
		// Attributes
		attributes := logRecord.Attributes()

		// CheckResponse
		checkResponse := &flowcontrolv1.CheckResponse{}

		// Source specific processing
		source, exists := attributes.Get(otelcollector.ApertureSourceLabel)
		if !exists {
			return retErr("aperture source label not found")
		}
		sourceStr := source.StringVal()
		if sourceStr == otelcollector.ApertureSourceSDK {
			success := otelcollector.GetStruct(attributes, otelcollector.ApertureCheckResponseLabel, checkResponse, []string{})
			if !success {
				return retErr("aperture check response label not found in Envoy access logs")
			}

			addSDKSpecificLabels(attributes)
		} else if sourceStr == otelcollector.ApertureSourceEnvoy {
			success := otelcollector.GetStruct(attributes, otelcollector.ApertureCheckResponseLabel, checkResponse, []string{otelcollector.EnvoyMissingAttributeValue})
			if !success {
				return retErr("aperture check response label not found in SDK access logs")
			}

			addEnvoySpecificLabels(attributes)
		} else {
			return retErr("aperture source label not recognized")
		}

		addCheckResponseBasedLabels(attributes, checkResponse, sourceStr)

		// Update metrics and enforce include list to eliminate any excess attributes
		if sourceStr == otelcollector.ApertureSourceSDK {
			p.updateMetrics(attributes, checkResponse, []string{})
			enforceIncludeListSDK(attributes)
		} else if sourceStr == otelcollector.ApertureSourceEnvoy {
			p.updateMetrics(attributes, checkResponse, []string{otelcollector.EnvoyMissingAttributeValue})
			enforceIncludeListHTTP(attributes)
		}

		// Add dynamic Flow labels
		addFlowLabels(attributes, checkResponse)

		return nil
	})
	return ld, err
}

func addSDKSpecificLabels(attributes pcommon.Map) {
	// Compute durations
	flowStart, flowStartExists := getSDKLabelTimestampValue(attributes, otelcollector.ApertureFlowStartTimestampLabel)
	workloadStart, workloadStartExists := getSDKLabelTimestampValue(attributes, otelcollector.ApertureWorkloadStartTimestampLabel)
	flowEnd, flowEndExists := getSDKLabelTimestampValue(attributes, otelcollector.ApertureFlowEndTimestampLabel)

	if flowStartExists && flowEndExists {
		flowDuration := flowEnd.Sub(flowStart)
		attributes.PutDouble(otelcollector.FlowDurationLabel, float64(flowDuration.Milliseconds()))
	}
	if workloadStartExists && flowEndExists {
		workloadDuration := flowEnd.Sub(workloadStart)
		attributes.PutDouble(otelcollector.WorkloadDurationLabel, float64(workloadDuration.Milliseconds()))
	}
}

func addEnvoySpecificLabels(attributes pcommon.Map) {
	treatAsZero := []string{otelcollector.EnvoyMissingAttributeValue}
	// Retrieve request length
	requestLength, _ := otelcollector.GetFloat64(attributes, otelcollector.EnvoyBytesSentLabel, treatAsZero)
	attributes.PutDouble(otelcollector.HTTPRequestContentLength, requestLength)
	// Retrieve response lengths
	responseLength, _ := otelcollector.GetFloat64(attributes, otelcollector.EnvoyBytesReceivedLabel, treatAsZero)
	attributes.PutDouble(otelcollector.HTTPResponseContentLength, responseLength)

	// Compute durations
	responseDuration, responseDurationExists := otelcollector.GetFloat64(attributes, otelcollector.EnvoyResponseDurationLabel, treatAsZero)
	authzDuration, authzDurationExists := otelcollector.GetFloat64(attributes, otelcollector.EnvoyAuthzDurationLabel, []string{})

	if responseDurationExists {
		attributes.PutDouble(otelcollector.FlowDurationLabel, responseDuration)
	}

	if responseDurationExists && authzDurationExists {
		attributes.PutDouble(otelcollector.WorkloadDurationLabel, responseDuration-authzDuration)
	}
}

func getSDKLabelTimestampValue(attributes pcommon.Map, labelKey string) (time.Time, bool) {
	return getLabelTimestampValue(attributes, labelKey, otelcollector.ApertureSourceSDK)
}

func getLabelTimestampValue(attributes pcommon.Map, labelKey, source string) (time.Time, bool) {
	value, exists := getLabelValue(attributes, labelKey, source)
	if !exists {
		return time.Time{}, false
	}
	return _getLabelTimestampValue(value, labelKey, source)
}

func _getLabelTimestampValue(value pcommon.Value, labelKey, source string) (time.Time, bool) {
	var valueInt int64
	if value.Type() == pcommon.ValueTypeInt {
		valueInt = value.IntVal()
	} else {
		log.Sample(zerolog.Sometimes).Warn().Str("source", source).Str("key", labelKey).Msg("Failed to parse a timestamp field")
		return time.Time{}, false
	}

	return time.Unix(0, valueInt), true
}

func getLabelValue(attributes pcommon.Map, labelKey, source string) (pcommon.Value, bool) {
	value, exists := attributes.Get(labelKey)
	if !exists {
		log.Sample(zerolog.Sometimes).Warn().Str("source", source).Str("key", labelKey).Msg("Label not found")
		return pcommon.Value{}, false
	}
	return value, exists
}

// addCheckResponseBasedLabels adds the following labels:
// * `decision_type`
// * `decision_reason`
// * `rate_limiters`
// * `dropping_rate_limiters`
// * `concurrency_limiters`
// * `dropping_concurrency_limiters`
// * `flux_meters`.
// * `flow_label_keys`.
// * `classifiers`.
func addCheckResponseBasedLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse, sourceStr string) {
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
		servicesValue.SliceVal().AppendEmpty().SetStringVal(service)
	}
	servicesValue.CopyTo(attributes.PutEmpty(otelcollector.ApertureServicesLabel))
	// Control Point
	attributes.PutString(otelcollector.ApertureControlPointLabel, checkResponse.GetControlPointInfo().String())

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
		otelcollector.ApertureDecisionTypeLabel:                pcommon.NewValueString(checkResponse.DecisionType.String()),
		otelcollector.ApertureRejectReasonLabel:                pcommon.NewValueString(checkResponse.GetRejectReason().String()),
		otelcollector.ApertureErrorLabel:                       pcommon.NewValueString(checkResponse.GetError().String()),
	}
	for _, decision := range checkResponse.LimiterDecisions {
		if decision.GetRateLimiterInfo() != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.ApertureRateLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			if decision.Dropped {
				labels[otelcollector.ApertureDroppingRateLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			}
		}
		if cl := decision.GetConcurrencyLimiterInfo(); cl != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.ApertureConcurrencyLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			if decision.Dropped {
				labels[otelcollector.ApertureDroppingConcurrencyLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			}

			workloadsRawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.WorkloadIndexLabel, cl.GetWorkloadIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value = strings.Join(workloadsRawValue, ",")
			labels[otelcollector.ApertureWorkloadsLabel].SliceVal().AppendEmpty().SetStringVal(value)
			if decision.Dropped {
				labels[otelcollector.ApertureDroppingWorkloadsLabel].SliceVal().AppendEmpty().SetStringVal(value)
			}
		}
	}
	for _, fluxMeter := range checkResponse.FluxMeterInfos {
		value := fluxMeter.GetFluxMeterName()
		labels[otelcollector.ApertureFluxMetersLabel].SliceVal().AppendEmpty().SetStringVal(value)
	}

	for _, flowLabelKey := range checkResponse.GetFlowLabelKeys() {
		labels[otelcollector.ApertureFlowLabelKeysLabel].SliceVal().AppendEmpty().SetStringVal(flowLabelKey)
	}

	for _, classifier := range checkResponse.ClassifierInfos {
		rawValue := []string{
			fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, classifier.PolicyName),
			fmt.Sprintf("%s:%v", metrics.ClassifierIndexLabel, classifier.ClassifierIndex),
		}
		value := strings.Join(rawValue, ",")
		labels[otelcollector.ApertureClassifiersLabel].SliceVal().AppendEmpty().SetStringVal(value)

		// add errors as attributes as well
		if classifier.Error != flowcontrolv1.ClassifierInfo_ERROR_NONE {
			errorsValue := []string{
				classifier.Error.String(),
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, classifier.PolicyName),
				fmt.Sprintf("%s:%v", metrics.ClassifierIndexLabel, classifier.ClassifierIndex),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, classifier.PolicyHash),
			}
			joinedValue := strings.Join(errorsValue, ",")
			labels[otelcollector.ApertureClassifierErrorsLabel].SliceVal().AppendEmpty().SetStringVal(joinedValue)
		}
	}

	for key, value := range labels {
		value.CopyTo(attributes.PutEmpty(key))
	}
}

func addFlowLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse) {
	for key, value := range checkResponse.TelemetryFlowLabels {
		pcommon.NewValueString(value).CopyTo(attributes.PutEmpty(key))
	}
}

func (p *metricsProcessor) updateMetrics(
	attributes pcommon.Map,
	checkResponse *flowcontrolv1.CheckResponse,
	treatAsZero []string,
) {
	if checkResponse == nil {
		return
	}
	if len(checkResponse.LimiterDecisions) > 0 {
		// Update workload metrics
		latency, _ := otelcollector.GetFloat64(attributes, otelcollector.WorkloadDurationLabel, []string{})
		for _, decision := range checkResponse.LimiterDecisions {
			if cl := decision.GetConcurrencyLimiterInfo(); cl != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:     decision.PolicyName,
					metrics.PolicyHashLabel:     decision.PolicyHash,
					metrics.ComponentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
					metrics.DecisionTypeLabel:   checkResponse.DecisionType.String(),
					metrics.WorkloadIndexLabel:  cl.GetWorkloadIndex(),
				}

				p.updateMetricsForWorkload(labels, latency)
			} // TODO: add rate limiter metrics
		}
	}

	if len(checkResponse.FluxMeterInfos) > 0 {
		// Update flux meter metrics
		statusCodeStr := ""
		statusCode, exists := attributes.Get(otelcollector.HTTPStatusCodeLabel)
		if exists {
			statusCodeStr = statusCode.StringVal()
		}
		featureStatusStr := ""
		featureStatus, exists := attributes.Get(otelcollector.ApertureFeatureStatusLabel)
		if exists {
			featureStatusStr = featureStatus.StringVal()
		}
		for _, fluxMeter := range checkResponse.FluxMeterInfos {
			p.updateMetricsForFluxMeters(fluxMeter, checkResponse.DecisionType, statusCodeStr, featureStatusStr, attributes, treatAsZero)
		}
	}
}

func (p *metricsProcessor) updateMetricsForWorkload(labels map[string]string, latency float64) {
	latencyHistogram, err := p.cfg.metricsAPI.GetTokenLatencyHistogram(labels)
	if err != nil {
		log.Sample(zerolog.Sometimes).Warn().Err(err).Msg("Getting latency histogram")
		return
	}
	latencyHistogram.Observe(latency)
}

func (p *metricsProcessor) updateMetricsForFluxMeters(
	fluxMeterMessage *flowcontrolv1.FluxMeterInfo,
	decisionType flowcontrolv1.CheckResponse_DecisionType,
	statusCode string,
	featureStatus string,
	attributes pcommon.Map,
	treatAsZero []string,
) {
	fluxMeter := p.cfg.engine.GetFluxMeter(fluxMeterMessage.FluxMeterName)
	if fluxMeter == nil {
		log.Sample(zerolog.Sometimes).Warn().Str(metrics.FluxMeterNameLabel, fluxMeterMessage.GetFluxMeterName()).
			Str(metrics.DecisionTypeLabel, decisionType.String()).
			Str(metrics.StatusCodeLabel, statusCode).
			Str(metrics.FeatureStatusLabel, featureStatus).
			Msg("FluxMeter not found")
		return
	}

	// metricValue is the value at fluxMeter's AttributeKey
	metricValue, _ := otelcollector.GetFloat64(attributes, fluxMeter.GetAttributeKey(), treatAsZero)

	fluxMeterHistogram := fluxMeter.GetHistogram(decisionType, statusCode, featureStatus)
	if fluxMeterHistogram != nil {
		fluxMeterHistogram.Observe(metricValue)
	}
}

/*
 * IncludeList: This IncludeList is applied to logs and spans at the beginning of enrichment process.
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

func enforceIncludeListHTTP(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListHTTP)
}

func enforceIncludeListSDK(attributes pcommon.Map) {
	otelcollector.EnforceIncludeList(attributes, includeListSDK)
}
