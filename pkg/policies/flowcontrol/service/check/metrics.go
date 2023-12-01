package check

import (
	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
)

// Metrics is used for collecting metrics about Aperture flowcontrol.
type Metrics interface {
	// CheckResponse collects metrics about Aperture Check call with DecisionType and Reason.
	CheckResponse(
		decision flowcontrolv1.CheckResponse_DecisionType,
		rejectReason flowcontrolv1.CheckResponse_RejectReason,
		agentInfo *agentinfo.AgentInfo,
	)
}

// NopMetrics is a no-op implementation of Metrics.
type NopMetrics struct{}

// Ensure NopMetrics implements Metrics interface.
var _ Metrics = NopMetrics{}

// CheckResponse is no-op method for NopMetrics.
func (NopMetrics) CheckResponse(
	decision flowcontrolv1.CheckResponse_DecisionType,
	rejectReason flowcontrolv1.CheckResponse_RejectReason,
	agentInfo *agentinfo.AgentInfo,
) {
}

// PrometheusMetrics stores collected metrics.
type PrometheusMetrics struct {
	registry *prometheus.Registry

	// Flow control service metrics
	// TODO: 3 gauges for 3 types of flowcontrol decisions
	checkReceivedTotal prometheus.Counter
	checkDecision      prometheus.CounterVec
	rejectReason       prometheus.CounterVec
}

// Ensure PrometheusMetrics implements Metrics interface.
var _ Metrics = &PrometheusMetrics{}

func (pm *PrometheusMetrics) allMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		pm.checkReceivedTotal,
		pm.checkDecision,
		pm.rejectReason,
	}
}

// NewPrometheusMetrics creates a Prometheus metrics collector.
func NewPrometheusMetrics(registry *prometheus.Registry) (*PrometheusMetrics, error) {
	pm := &PrometheusMetrics{
		registry: registry,
		checkReceivedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: metrics.FlowControlRequestsMetricName,
				Help: "Total number of aperture check requests handled",
			},
		),
		checkDecision: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.FlowControlDecisionsMetricName,
				Help: "Number of aperture check decisions",
			}, []string{metrics.FlowControlCheckDecisionTypeLabel, metrics.AgentGroupLabel},
		),
		rejectReason: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.FlowControlRejectReasonsMetricName,
				Help: "Number of reject reasons other than unspecified",
			}, []string{metrics.FlowControlCheckRejectReasonLabel},
		),
	}

	for _, m := range pm.allMetrics() {
		err := pm.registry.Register(m)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to register metric")
			return nil, err
		}
	}

	return pm, nil
}

// CheckResponse collects metrics about Aperture Check call with DecisionType, RejectReason, Error.
func (pm *PrometheusMetrics) CheckResponse(
	decision flowcontrolv1.CheckResponse_DecisionType,
	rejectReason flowcontrolv1.CheckResponse_RejectReason,
	agentInfo *agentinfo.AgentInfo,
) {
	pm.checkReceivedTotal.Inc()
	pm.checkDecision.With(prometheus.Labels{
		metrics.FlowControlCheckDecisionTypeLabel: decision.Enum().String(),
		metrics.AgentGroupLabel:                   agentInfo.GetAgentGroup(),
	}).Inc()
	if rejectReason != flowcontrolv1.CheckResponse_REJECT_REASON_NONE {
		pm.rejectReason.With(prometheus.Labels{metrics.FlowControlCheckRejectReasonLabel: rejectReason.Enum().String()}).Inc()
	}
}
