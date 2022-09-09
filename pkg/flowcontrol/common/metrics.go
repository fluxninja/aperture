package common

import (
	"github.com/prometheus/client_golang/prometheus"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
)

// Metrics is used for collecting metrics about Aperture flowcontrol.
type Metrics interface {
	// CheckResponse collects metrics about Aperture Check call with DecisionType and Reason.
	CheckResponse(flowcontrolv1.DecisionType, *flowcontrolv1.DecisionReason)
}

// NopMetrics is a no-op implementation of Metrics.
type NopMetrics struct{}

// CheckResponse is no-op method for NopMetrics.
func (NopMetrics) CheckResponse(flowcontrolv1.DecisionType, *flowcontrolv1.DecisionReason) {
}

// PrometheusMetrics stores collected metrics.
type PrometheusMetrics struct {
	registry *prometheus.Registry

	// Flow control service metrics
	// TODO: 3 gauges for 3 types of flowcontrol decisions
	checkReceivedTotal prometheus.Counter
	checkDecision      prometheus.CounterVec
	errorReason        prometheus.CounterVec
	rejectReason       prometheus.CounterVec
}

func (pm *PrometheusMetrics) allMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		pm.checkReceivedTotal,
		pm.checkDecision,
		pm.errorReason,
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
			}, []string{metrics.FlowControlCheckDecisionTypeLabel},
		),
		errorReason: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.FlowControlErrorReasonMetricName,
				Help: "Number of error reasons other than unspecified",
			}, []string{metrics.FlowControlCheckErrorReasonLabel},
		),
		rejectReason: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.FlowControlRejectReasonMetricName,
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

// CheckResponse collects metrics about Aperture Check call with DecisionType, Reason.
func (pm *PrometheusMetrics) CheckResponse(decision flowcontrolv1.DecisionType, reason *flowcontrolv1.DecisionReason) {
	pm.checkReceivedTotal.Inc()
	pm.checkDecision.With(prometheus.Labels{metrics.FlowControlCheckDecisionTypeLabel: decision.Enum().String()}).Inc()
	if reason != nil {
		if reason.GetErrorReason() != flowcontrolv1.DecisionReason_ERROR_REASON_UNSPECIFIED {
			pm.errorReason.With(prometheus.Labels{metrics.FlowControlCheckErrorReasonLabel: reason.GetErrorReason().Enum().String()}).Inc()
		}
		if reason.GetRejectReason() != flowcontrolv1.DecisionReason_REJECT_REASON_UNSPECIFIED {
			pm.rejectReason.With(prometheus.Labels{metrics.FlowControlCheckRejectReasonLabel: reason.GetRejectReason().Enum().String()}).Inc()
		}
	}
	// TODO: update fluxmeter metrics
}
