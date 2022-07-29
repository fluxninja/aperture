package metricsprocessor

import (
	"context"
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
		decisions := getDecisions(logRecord.Attributes())
		fluxMeters := getFluxMeterIDs(logRecord.Attributes())
		p.addPolicyLabels(logRecord.Attributes(), decisions)
		return p.updateMetrics(logRecord.Attributes(), decisions, fluxMeters)
	})
	return ld, err
}

// ConsumeTraces receives ptrace.Traces for consumption then returns updated traces with policy labels and metrics.
func (p *metricsProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	err := otelcollector.IterateSpans(td, func(span ptrace.Span) error {
		decisions := getDecisions(span.Attributes())
		fluxMeters := getFluxMeterIDs(span.Attributes())
		p.addPolicyLabels(span.Attributes(), decisions)
		return p.updateMetrics(span.Attributes(), decisions, fluxMeters)
	})
	return td, err
}

func (p *metricsProcessor) addPolicyLabels(attributes pcommon.Map, decisions []*flowcontrolv1.LimiterDecision) {
	matched, dropped := getIDs(decisions)
	attributes.InsertString(otelcollector.PoliciesMatchedLabel, matched)
	attributes.InsertString(otelcollector.PoliciesDroppedLabel, dropped)
	attributes.Remove(otelcollector.LimiterDecisionsLabel)
}

func (p *metricsProcessor) updateMetrics(attributes pcommon.Map, decisions []*flowcontrolv1.LimiterDecision, fluxMeters []string) error {
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

	for _, decision := range decisions {
		labels := map[string]string{
			policyNameLabel:     decision.PolicyName,
			policyHashLabel:     decision.PolicyHash,
			componentIndexLabel: fmt.Sprintf("%d", decision.ComponentIndex),
			droppedMetricsLabel: fmt.Sprintf("%t", decision.Dropped),
		}

		workload := decision.GetConcurrencyLimiter().Workload
		err = p.updateMetricsForWorkload(labels, latency, workload)
		if err != nil {
			return err
		}
	}

	for _, fluxMeter := range fluxMeters {
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

func (p *metricsProcessor) updateMetricsForFluxMeters(fluxMeter string, latency float64) {
	fluxmeterHistogram := p.cfg.engine.GetFluxMeterHist(fluxMeter)
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

func getDecisions(attributes pcommon.Map) []*flowcontrolv1.LimiterDecision {
	var decisions []*flowcontrolv1.LimiterDecision
	rawPolicies, exists := attributes.Get(otelcollector.LimiterDecisionsLabel)
	if !exists {
		log.Debug().Str("label", otelcollector.LimiterDecisionsLabel).Msg("Label does not exist")
		return decisions
	}
	if !otelcollector.UnmarshalStringVal(rawPolicies, otelcollector.LimiterDecisionsLabel, &decisions) {
		log.Debug().Str("label", otelcollector.LimiterDecisionsLabel).Msg("Label is not a string")
	}
	return decisions
}

func getFluxMeterIDs(attributes pcommon.Map) []string {
	rawFluxMeters, exists := attributes.Get(otelcollector.FluxMetersLabel)
	if !exists {
		log.Debug().Str("label", otelcollector.FluxMetersLabel).Msg("Label does not exist")
		return nil
	}

	// flux meters are either send properly as a slice of strings (when sent
	// via sdk) or as json-encoded array
	var fluxMeters []string
	if otelcollector.UnmarshalStringVal(rawFluxMeters, otelcollector.FluxMetersLabel, &fluxMeters) {
		return fluxMeters
	}

	slice := rawFluxMeters.SliceVal()
	fluxMeters = make([]string, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		fluxMeters = append(fluxMeters, slice.At(i).StringVal())
	}

	return fluxMeters
}

func getIDs(decisions []*flowcontrolv1.LimiterDecision) (string, string) {
	matched := make([]string, len(decisions))
	dropped := make([]string, len(decisions))
	for i, decision := range decisions {
		matched[i] = decision.PolicyName
		if decision.Dropped {
			dropped[i] = decision.PolicyName
		}
	}
	return strings.Join(matched, ","), strings.Join(dropped, ",")
}
