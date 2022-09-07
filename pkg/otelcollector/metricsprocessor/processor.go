package metricsprocessor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/otelcollector"
)

type metricsProcessor struct {
	cfg                    *Config
	workloadLatencySummary *prometheus.SummaryVec
}

func newProcessor(cfg *Config) (*metricsProcessor, error) {
	p := &metricsProcessor{
		cfg: cfg,
	}
	err := p.registerRequestLatencyHistogram()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *metricsProcessor) registerRequestLatencyHistogram() error {
	p.workloadLatencySummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.WorkloadLatencyMetricName,
		Help: "Latency summary of workload",
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIndexLabel,
		metrics.DecisionTypeLabel,
		metrics.WorkloadIndexLabel,
	})
	err := p.cfg.promRegistry.Register(p.workloadLatencySummary)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			// We're registering this histogram vec from multiple processors
			// (logs processor and traces processor), so if both processors are
			// enabled, it's expected that whichever processor is created
			// second, it will see that the histogram vec was already
			// registered. Use the existing histogram vec from now on.
			p.workloadLatencySummary = are.ExistingCollector.(*prometheus.SummaryVec)
			return nil
		}
	}
	return err
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
		checkResponse := p.addCheckResponseBasedLabels(logRecord.Attributes(), []string{otelcollector.MissingAttributeSourceValue})

		p.addAuthzResponseBasedLabels(logRecord.Attributes(), []string{otelcollector.MissingAttributeSourceValue})

		return p.updateMetrics(logRecord.Attributes(), checkResponse)
		// TODO tgill: Pass attributes through white list to ensure we drop any extra fields before rollup
		// p.whitelistLogAttributes(logRecord.Attributes())
	})
	return ld, err
}

// ConsumeTraces receives ptrace.Traces for consumption then returns updated traces with policy labels and metrics.
func (p *metricsProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	err := otelcollector.IterateSpans(td, func(span ptrace.Span) error {
		checkResponse := p.addCheckResponseBasedLabels(span.Attributes(), []string{})

		// TODO tgill: move span latency addition to its own function
		endTimestamp := span.EndTimestamp()
		startTimeStamp := span.StartTimestamp()
		latency := float64(endTimestamp.AsTime().Sub(startTimeStamp.AsTime())) / float64(time.Millisecond)
		span.Attributes().UpsertDouble(otelcollector.DurationLabel, latency)

		return p.updateMetrics(span.Attributes(), checkResponse)
		// TODO tgill: Pass attributes through white list to ensure we drop any extra fields before rollup
		// p.whitelistSpanAttributes(span.Attributes())
	})
	return td, err
}

func (p *metricsProcessor) addAuthzResponseBasedLabels(attributes pcommon.Map, treatAsMissing []string) {
	var authzResponse flowcontrolv1.AuthzResponse
	otelcollector.GetStruct(attributes, otelcollector.MarshalledAuthzResponseLabel, &authzResponse, treatAsMissing)
	labels := map[string]pcommon.Value{
		otelcollector.AuthzStatusLabel: pcommon.NewValueString(authzResponse.GetStatus().String()),
	}
	for key, value := range labels {
		attributes.Upsert(key, value)
	}
	attributes.Remove(otelcollector.MarshalledAuthzResponseLabel)
}

// addCheckResponseBasedLabels adds the following labels:
// * `decision_type`
// * `decision_reason`
// * `rate_limiters`
// * `dropping_rate_limiters`
// * `concurrency_limiters`
// * `dropping_concurrency_limiters`
// * `flux_meters`.
func (p *metricsProcessor) addCheckResponseBasedLabels(attributes pcommon.Map, treatAsMissing []string) *flowcontrolv1.CheckResponse {
	var checkResponse flowcontrolv1.CheckResponse
	success := otelcollector.GetStruct(attributes, otelcollector.MarshalledCheckResponseLabel, &checkResponse, treatAsMissing)
	if !success {
		return nil
	}
	labels := map[string]pcommon.Value{
		otelcollector.RateLimitersLabel:                pcommon.NewValueSlice(),
		otelcollector.DroppingRateLimitersLabel:        pcommon.NewValueSlice(),
		otelcollector.ConcurrencyLimitersLabel:         pcommon.NewValueSlice(),
		otelcollector.DroppingConcurrencyLimitersLabel: pcommon.NewValueSlice(),
		otelcollector.FluxMetersLabel:                  pcommon.NewValueSlice(),
		otelcollector.FlowLabelKeysLabel:               pcommon.NewValueSlice(),
		otelcollector.DecisionTypeLabel:                pcommon.NewValueString(checkResponse.DecisionType.String()),
		otelcollector.DecisionRejectReasonLabel:        pcommon.NewValueString(""),
		otelcollector.DecisionErrorReasonLabel:         pcommon.NewValueString(""),
	}
	if checkResponse.DecisionReason != nil {
		labels[otelcollector.DecisionErrorReasonLabel] = pcommon.NewValueString(checkResponse.DecisionReason.GetErrorReason().String())
		labels[otelcollector.DecisionRejectReasonLabel] = pcommon.NewValueString(checkResponse.DecisionReason.GetRejectReason().String())
	}
	for _, decision := range checkResponse.LimiterDecisions {
		if decision.GetRateLimiter() != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.RateLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			if decision.Dropped {
				labels[otelcollector.DroppingRateLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			}
		}
		if cl := decision.GetConcurrencyLimiter(); cl != nil {
			rawValue := []string{
				fmt.Sprintf("%s:%v", metrics.PolicyNameLabel, decision.GetPolicyName()),
				fmt.Sprintf("%s:%v", metrics.ComponentIndexLabel, decision.GetComponentIndex()),
				fmt.Sprintf("%s:%v", metrics.WorkloadIndexLabel, cl.GetWorkloadIndex()),
				fmt.Sprintf("%s:%v", metrics.PolicyHashLabel, decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.ConcurrencyLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			if decision.Dropped {
				labels[otelcollector.DroppingConcurrencyLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			}
		}
	}
	for _, fluxMeter := range checkResponse.FluxMeters {
		rawValue := []string{
			fmt.Sprintf("%s:%v", metrics.FluxMeterNameLabel, fluxMeter.GetFluxMeterName()),
		}
		value := strings.Join(rawValue, ",")
		labels[otelcollector.FluxMetersLabel].SliceVal().AppendEmpty().SetStringVal(value)
	}

	for _, flowLabelKey := range checkResponse.GetFlowLabelKeys() {
		labels[otelcollector.FlowLabelKeysLabel].SliceVal().AppendEmpty().SetStringVal(flowLabelKey)
	}

	for key, value := range labels {
		attributes.Upsert(key, value)
	}
	attributes.Remove(otelcollector.MarshalledCheckResponseLabel)
	return &checkResponse
}

func (p *metricsProcessor) updateMetrics(
	attributes pcommon.Map,
	checkResponse *flowcontrolv1.CheckResponse,
) error {
	if checkResponse == nil {
		return nil
	}
	if len(checkResponse.LimiterDecisions) > 0 {
		// Update workload metrics
		latency, _ := otelcollector.GetFloat64(attributes, otelcollector.DurationLabel, []string{})
		for _, decision := range checkResponse.LimiterDecisions {
			if cl := decision.GetConcurrencyLimiter(); cl != nil {
				labels := map[string]string{
					metrics.PolicyNameLabel:     decision.PolicyName,
					metrics.PolicyHashLabel:     decision.PolicyHash,
					metrics.ComponentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
					metrics.DecisionTypeLabel:   checkResponse.DecisionType.String(),
					metrics.WorkloadIndexLabel:  cl.GetWorkloadIndex(),
				}

				err := p.updateMetricsForWorkload(labels, latency)
				if err != nil {
					return err
				}
			} // TODO: add rate limiter metrics
		}
	}

	if len(checkResponse.FluxMeters) > 0 {
		// Update flux meter metrics
		statusCodeStr := ""
		statusCode, exists := attributes.Get(otelcollector.StatusCodeLabel)
		if exists {
			statusCodeStr = statusCode.StringVal()
		}
		featureStatusStr := ""
		featureStatus, exists := attributes.Get(otelcollector.FeatureStatusLabel)
		if exists {
			featureStatusStr = featureStatus.StringVal()
		}
		for _, fluxMeter := range checkResponse.FluxMeters {
			p.updateMetricsForFluxMeters(fluxMeter, checkResponse.DecisionType, statusCodeStr, featureStatusStr, attributes)
		}
	}

	return nil
}

func (p *metricsProcessor) updateMetricsForWorkload(labels map[string]string, latency float64) error {
	latencyHistogram, err := p.workloadLatencySummary.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return err
	}
	latencyHistogram.Observe(latency)

	return nil
}

func (p *metricsProcessor) updateMetricsForFluxMeters(
	fluxMeterMessage *flowcontrolv1.FluxMeter,
	decisionType flowcontrolv1.DecisionType,
	statusCode string,
	featureStatus string,
	attributes pcommon.Map,
) {
	// TODO tgill: Remove additional metric keys from attributes map. User can refer to additional metrics using attribute_key in flux meter. The extra metrics need to be discarded to prevent cardinality blow up in rollup processor
	// defer

	fluxMeter := p.cfg.engine.GetFluxMeter(fluxMeterMessage.FluxMeterName)
	if fluxMeter == nil {
		log.Debug().Str(metrics.FluxMeterNameLabel, fluxMeterMessage.GetFluxMeterName()).
			Str(metrics.DecisionTypeLabel, decisionType.String()).
			Str(metrics.StatusCodeLabel, statusCode).
			Str(metrics.FeatureStatusLabel, featureStatus).
			Msg("FluxMeter not found")
		return
	}

	// metricValue is the value at fluxMeter's AttributeKey
	metricValue, _ := otelcollector.GetFloat64(attributes, fluxMeter.GetAttributeKey(), []string{})

	fluxMeterHistogram := fluxMeter.GetHistogram(decisionType, statusCode, featureStatus)
	if fluxMeterHistogram != nil {
		fluxMeterHistogram.Observe(metricValue)
	}
}

// TODO tgill
/*func (p *metricsProcessor) whitelistLogAttributes(attributes pcommon.Map) {
	p._whitelistCommonAttributes(attributes)
	p._whitelistLogAttributes(attributes)
}

func (p *metricsProcessor) _whitelistCommonAttributes(attributes pcommon.Map) {
	for key := range attributes {
		if !p.cfg.whitelistedLogAttributes[key] {
			attributes.Delete(key)
		}
	}
}*/
