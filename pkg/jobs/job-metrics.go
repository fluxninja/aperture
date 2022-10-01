package jobs

import (
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// JobMetrics holds prometheus metrics related to jobs.
type JobMetrics struct {
	errorTotal     *prometheus.CounterVec
	executionTotal *prometheus.CounterVec
	latencySummary *prometheus.SummaryVec
}

func newJobMetrics() *JobMetrics {
	jobMetricsLabels := []string{metrics.JobNameLabel}
	return &JobMetrics{
		errorTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.JobErrorMetricName,
				Help: "The total number of job errors that occurred",
			}, jobMetricsLabels),
		executionTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.JobExecutionMetricName,
				Help: "The total number job executions that occurred",
			}, jobMetricsLabels),
		latencySummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       metrics.JobLatencyMetricName,
				Help:       "Latency summary of the job executions",
				Objectives: map[float64]float64{0: 0, 0.01: 0.001, 0.05: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 1: 0},
			}, jobMetricsLabels),
	}
}

func (jm *JobMetrics) allMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		jm.errorTotal,
		jm.executionTotal,
		jm.latencySummary,
	}
}
