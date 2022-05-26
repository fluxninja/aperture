package metricsprocessor

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"

	policydecisionsv1 "aperture.tech/aperture/api/gen/proto/go/aperture/policy/decisions/v1"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/otelcollector"
)

const (
	latencyHistogramMetricName = "request_latency_ms"

	metricIDLabel                = "metric_id"
	droppedMetricsLabel          = "dropped"
	workloadKeyNameMetricsLabel  = "workload_key_name"
	workloadKeyValueMetricsLabel = "workload_key_value"
)

type metricsProcessor struct {
	cfg                     *Config
	requestLatencyHistogram *prometheus.HistogramVec
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
	p.requestLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: latencyHistogramMetricName,
		Help: "Latency of requests histogram",
		Buckets: prometheus.LinearBuckets(
			p.cfg.LatencyBucketStartMS,
			p.cfg.LatencyBucketWidthMS,
			p.cfg.LatencyBucketCount,
		),
	}, []string{
		metricIDLabel,
		droppedMetricsLabel,
		workloadKeyNameMetricsLabel,
		workloadKeyValueMetricsLabel,
	})
	err := p.cfg.promRegistry.Register(p.requestLatencyHistogram)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			// A histogram for that metric has been registered before. Use the old histogram from now on.
			p.requestLatencyHistogram = are.ExistingCollector.(*prometheus.HistogramVec)
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
		policies := getPolicies(logRecord.Attributes())
		fluxmeters := getFluxMeterIDs(logRecord.Attributes())
		p.addPolicyLabels(logRecord.Attributes(), policies)
		return p.updateMetrics(logRecord.Attributes(), policies, fluxmeters)
	})
	return ld, err
}

// ConsumeTraces receives ptrace.Traces for consumption then returns updated traces with policy labels and metrics.
func (p *metricsProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	err := otelcollector.IterateSpans(td, func(span ptrace.Span) error {
		policies := getPolicies(span.Attributes())
		fluxmeters := getFluxMeterIDs(span.Attributes())
		p.addPolicyLabels(span.Attributes(), policies)
		return p.updateMetrics(span.Attributes(), policies, fluxmeters)
	})
	return td, err
}

func (p *metricsProcessor) addPolicyLabels(
	attributes pcommon.Map,
	policies []policy,
) {
	matched, dropped := getIDs(policies)
	attributes.InsertString(otelcollector.PoliciesMatchedLabel, matched)
	attributes.InsertString(otelcollector.PoliciesDroppedLabel, dropped)
	attributes.Remove(otelcollector.PoliciesLabel)
}

func (p *metricsProcessor) updateMetrics(
	attributes pcommon.Map,
	policies []policy,
	fluxmeters []string,
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

	for _, policy := range policies {
		labels := map[string]string{
			metricIDLabel:       policy.ID,
			droppedMetricsLabel: fmt.Sprintf("%t", policy.Dropped),
		}

		workload := p.extractWorkloadDesc(policy.Workload)
		err = p.updateMetricsForWorkload(labels, latency, workload)
		if err != nil {
			return err
		}
	}

	for _, fluxmeterID := range fluxmeters {
		p.updateMetricsForFluxMeters(fluxmeterID, latency)
	}

	return nil
}

func (p *metricsProcessor) extractWorkloadDesc(workload string) *policydecisionsv1.WorkloadDesc {
	workloadDesc := &policydecisionsv1.WorkloadDesc{
		WorkloadKey:   "default_workload_key",
		WorkloadValue: "default_workload_value",
	}
	re := regexp.MustCompile(`\"(.*)\"`)

	arr := strings.Split(workload, " ")
	for _, elem := range arr {
		if strings.Contains(elem, "workload_key") {
			match := re.FindStringSubmatch(elem)
			if len(match) > 1 {
				workloadDesc.WorkloadKey = match[1]
			}
		}
		if strings.Contains(elem, "workload_value") {
			match := re.FindStringSubmatch(elem)
			if len(match) > 1 {
				workloadDesc.WorkloadValue = match[1]
			}
		}
	}
	return workloadDesc
}

func (p *metricsProcessor) updateMetricsForWorkload(
	labels map[string]string,
	latency float64,
	workload *policydecisionsv1.WorkloadDesc,
) error {
	labels[workloadKeyNameMetricsLabel] = workload.WorkloadKey
	labels[workloadKeyValueMetricsLabel] = workload.WorkloadValue
	latencyHistogram, err := p.requestLatencyHistogram.GetMetricWith(labels)
	if err != nil {
		log.Warn().Err(err).Msg("Getting latency histogram")
		return err
	}
	latencyHistogram.Observe(latency)

	return nil
}

func (p *metricsProcessor) updateMetricsForFluxMeters(fluxmeterID string, latency float64) {
	fluxmeterHistogram := p.cfg.engine.GetFluxMeterHist(fluxmeterID)
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

func getPolicies(attributes pcommon.Map) (policies []policy) {
	rawPolicies, exists := attributes.Get(otelcollector.PoliciesLabel)
	if !exists {
		log.Debug().Str("label", otelcollector.PoliciesLabel).Msg("Label does not exist")
		return
	}
	if !otelcollector.UnmarshalStringVal(rawPolicies, otelcollector.PoliciesLabel, &policies) {
		log.Debug().Str("label", otelcollector.PoliciesLabel).Msg("Label is not a string")
	}
	return
}

func getFluxMeterIDs(attributes pcommon.Map) []string {
	rawFluxMeters, exists := attributes.Get(otelcollector.FluxMeterIDsLabel)
	if !exists {
		log.Debug().Str("label", otelcollector.FluxMeterIDsLabel).Msg("Label does not exist")
		return nil
	}

	// flux meters are either send properly as a slice of strings (when sent
	// via sdk) or as json-encoded array
	var fluxMeters []string
	if otelcollector.UnmarshalStringVal(rawFluxMeters, otelcollector.FluxMeterIDsLabel, &fluxMeters) {
		return fluxMeters
	}

	slice := rawFluxMeters.SliceVal()
	fluxMeters = make([]string, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		fluxMeters = append(fluxMeters, slice.At(i).StringVal())
	}

	return fluxMeters
}

type policy struct {
	Workload string `json:"workload"`
	ID       string `json:"id"`
	Dropped  bool   `json:"dropped"`
}

func getIDs(policies []policy) (string, string) {
	matched := make([]string, len(policies))
	dropped := make([]string, len(policies))
	for i, p := range policies {
		matched[i] = p.ID
		if p.Dropped {
			dropped[i] = p.ID
		}
	}
	return strings.Join(matched, ","), strings.Join(dropped, ",")
}
