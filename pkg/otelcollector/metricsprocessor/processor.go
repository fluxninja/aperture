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
	"github.com/fluxninja/aperture/pkg/otelcollector"
	"github.com/fluxninja/aperture/pkg/paths"
)

const (
	workloadLatencyMetricName = "workload_latency_ms"

	policyNameLabel     = "policy_name"
	policyHashLabel     = "policy_hash"
	componentIndexLabel = "component_index"
	droppedMetricsLabel = "dropped"
	workloadIndexLabel  = "workload_index"
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
		Name: workloadLatencyMetricName,
		Help: "Latency histogram of workload",
		Buckets: prometheus.LinearBuckets(
			p.cfg.LatencyBucketStartMS,
			p.cfg.LatencyBucketWidthMS,
			p.cfg.LatencyBucketCount,
		),
	}, []string{
		policyNameLabel,
		policyHashLabel,
		componentIndexLabel,
		droppedMetricsLabel,
		workloadIndexLabel,
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
				fmt.Sprintf("policy_name:%v", decision.GetPolicyName()),
				fmt.Sprintf("component_index:%v", decision.GetComponentIndex()),
				fmt.Sprintf("policy_hash:%v", decision.GetPolicyHash()),
			}
			value := strings.Join(rawValue, ",")
			labels[otelcollector.RateLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			if decision.Dropped {
				labels[otelcollector.DroppingRateLimitersLabel].SliceVal().AppendEmpty().SetStringVal(value)
			}
		}
		if cl := decision.GetConcurrencyLimiter(); cl != nil {
			rawValue := []string{
				fmt.Sprintf("policy_name:%v", decision.GetPolicyName()),
				fmt.Sprintf("component_index:%v", decision.GetComponentIndex()),
				fmt.Sprintf("workload_index:%v", cl.GetWorkload()),
				fmt.Sprintf("policy_hash:%v", decision.GetPolicyHash()),
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
			fmt.Sprintf("policy_name:%v", fluxMeter.GetPolicyName()),
			fmt.Sprintf("flux_meter_name:%v", fluxMeter.GetFluxMeterName()),
			fmt.Sprintf("policy_hash:%v", fluxMeter.GetPolicyHash()),
		}
		value := strings.Join(rawValue, ",")
		labels[otelcollector.FluxMetersLabel].SliceVal().AppendEmpty().SetStringVal(value)
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

	for _, decision := range checkResponse.LimiterDecisions {
		labels := map[string]string{
			policyNameLabel:     decision.PolicyName,
			policyHashLabel:     decision.PolicyHash,
			componentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
			droppedMetricsLabel: fmt.Sprintf("%t", decision.Dropped),
		}

		workload := ""
		if cl := decision.GetConcurrencyLimiter(); cl != nil {
			workload = cl.GetWorkload()
		}
		err = p.updateMetricsForWorkload(labels, latency, workload)
		if err != nil {
			return err
		}
	}

	for _, fluxMeter := range checkResponse.FluxMeters {
		p.updateMetricsForFluxMeters(fluxMeter, latency)
	}

	return nil
}

func (p *metricsProcessor) updateMetricsForWorkload(labels map[string]string, latency float64, workload string) error {
	labels[workloadIndexLabel] = workload
	latencyHistogram, err := p.workloadLatencyHistogram.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return err
	}
	latencyHistogram.Observe(latency)

	return nil
}

func (p *metricsProcessor) updateMetricsForFluxMeters(fluxMeter *flowcontrolv1.FluxMeter, latency float64) {
	fluxMeterID := paths.MetricIDForFluxMeterExpanded(
		fluxMeter.GetAgentGroup(),
		fluxMeter.GetPolicyName(),
		fluxMeter.GetFluxMeterName(),
		fluxMeter.GetPolicyHash(),
	)
	fluxmeterHistogram := p.cfg.engine.GetFluxMeterHist(fluxMeterID)
	if fluxmeterHistogram == nil {
		log.Debug().Str("fluxMeterID", fluxMeterID).Msg("Fluxmeter not found")
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
