package metricsprocessor

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	cfg                      *Config
	workloadLatencyHistogram *prometheus.HistogramVec
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
	p.workloadLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: metrics.WorkloadLatencyMetricName,
		Help: "Latency histogram of workload",
		Buckets: prometheus.LinearBuckets(
			p.cfg.LatencyBucketStartMS,
			p.cfg.LatencyBucketWidthMS,
			p.cfg.LatencyBucketCount,
		),
	}, []string{
		metrics.PolicyNameLabel,
		metrics.PolicyHashLabel,
		metrics.ComponentIndexLabel,
		metrics.DecisionTypeLabel,
		metrics.WorkloadIndexLabel,
	})
	err := p.cfg.promRegistry.Register(p.workloadLatencyHistogram)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			// A histogram for that metric has been registered before. Use the old histogram from now on.
			p.workloadLatencyHistogram = are.ExistingCollector.(*prometheus.HistogramVec)
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
		checkResponse := otelcollector.GetCheckResponse(logRecord.Attributes())
		if checkResponse == nil {
			return errors.New("failed to get check_response from attributes")
		}
		p.addCheckResponseBasedLabels(logRecord.Attributes(), checkResponse)

		authzResponse := otelcollector.GetAuthzResponse(logRecord.Attributes())
		p.addAuthzResponseBasedLabels(logRecord.Attributes(), authzResponse)

		log.Trace().Msgf("CheckResponse: %v", checkResponse)
		return p.updateMetrics(logRecord.Attributes(), checkResponse)
	})
	return ld, err
}

// ConsumeTraces receives ptrace.Traces for consumption then returns updated traces with policy labels and metrics.
func (p *metricsProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	err := otelcollector.IterateSpans(td, func(span ptrace.Span) error {
		checkResponse := otelcollector.GetCheckResponse(span.Attributes())
		if checkResponse == nil {
			return errors.New("failed getting check response from attributes")
		}
		p.addCheckResponseBasedLabels(span.Attributes(), checkResponse)
		return p.updateMetrics(span.Attributes(), checkResponse)
	})
	return td, err
}

func (p *metricsProcessor) addAuthzResponseBasedLabels(attributes pcommon.Map, authzResponse *flowcontrolv1.AuthzResponse) {
	labels := map[string]pcommon.Value{
		otelcollector.AuthzStatusLabel: pcommon.NewValueString(authzResponse.GetStatus().String()),
	}
	for key, value := range labels {
		attributes.Insert(key, value)
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
func (p *metricsProcessor) addCheckResponseBasedLabels(attributes pcommon.Map, checkResponse *flowcontrolv1.CheckResponse) {
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
		attributes.Insert(key, value)
	}
	attributes.Remove(otelcollector.MarshalledCheckResponseLabel)
}

func (p *metricsProcessor) updateMetrics(
	attributes pcommon.Map,
	checkResponse *flowcontrolv1.CheckResponse,
) error {
	latencyLabel := getLatencyLabel(attributes)

	if latencyLabel == "" {
		log.Debug().Msg("Failed determining latency label")
		return nil
	}
	rawLatency, exists := attributes.Get(latencyLabel)
	if !exists {
		log.Debug().Str("label", latencyLabel).Msg("Label does not exist")
		return nil
	}
	latency, err := strconv.ParseFloat(rawLatency.StringVal(), 64)
	if err != nil {
		log.Debug().Str("rawLatency", rawLatency.AsString()).Msg("Could not parse raw latency to float")
		return nil
	}
	statusCode, exists := attributes.Get(otelcollector.StatusCodeLabel)
	if !exists {
		log.Debug().Str("label", otelcollector.StatusCodeLabel).Msg("Label does not exist")
		return nil
	}
	statusCodeStr := statusCode.StringVal()

	for _, decision := range checkResponse.LimiterDecisions {
		labels := map[string]string{
			metrics.PolicyNameLabel:     decision.PolicyName,
			metrics.PolicyHashLabel:     decision.PolicyHash,
			metrics.ComponentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
			metrics.DecisionTypeLabel:   checkResponse.DecisionType.String(),
		}
		log.Trace().Msgf("labels: %v", labels)

		workload := ""
		if cl := decision.GetConcurrencyLimiter(); cl != nil {
			workload = cl.GetWorkloadIndex()
		}
		err = p.updateMetricsForWorkload(labels, latency, workload)
		if err != nil {
			return err
		}
	}

	for _, fluxMeter := range checkResponse.FluxMeters {
		p.updateMetricsForFluxMeters(fluxMeter, checkResponse.DecisionType, statusCodeStr, latency)
	}

	return nil
}

func (p *metricsProcessor) updateMetricsForWorkload(labels map[string]string, latency float64, workload string) error {
	labels[metrics.WorkloadIndexLabel] = workload
	latencyHistogram, err := p.workloadLatencyHistogram.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return err
	}
	latencyHistogram.Observe(latency)

	return nil
}

func (p *metricsProcessor) updateMetricsForFluxMeters(
	fluxMeter *flowcontrolv1.FluxMeter,
	decisionType flowcontrolv1.DecisionType,
	statusCode string,
	latency float64,
) {
	fluxmeterHistogram := p.cfg.engine.GetFluxMeterHist(
		fluxMeter.GetFluxMeterName(),
		statusCode,
		decisionType,
	)
	if fluxmeterHistogram == nil {
		log.Debug().Str(metrics.FluxMeterNameLabel, fluxMeter.GetFluxMeterName()).
			Str(metrics.DecisionTypeLabel, decisionType.String()).
			Str(metrics.StatusCodeLabel, statusCode).
			Msg("Fluxmeter not found")
		return
	}
	fluxmeterHistogram.Observe(latency)
}

func getLatencyLabel(attributes pcommon.Map) string {
	controlPoint, exists := attributes.Get(otelcollector.ControlPointLabel)
	if !exists {
		log.Debug().Str("label", otelcollector.ControlPointLabel).Msg("Label does not exist")
		// This should not happen
		return ""
	}
	switch controlPoint.AsString() {
	case otelcollector.ControlPointFeature:
		return otelcollector.FeatureDurationLabel
	case otelcollector.ControlPointIngress, otelcollector.ControlPointEgress:
		return otelcollector.HTTPDurationLabel
	}
	return ""
}
