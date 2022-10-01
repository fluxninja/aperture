package jobs

import (
	"github.com/fluxninja/aperture/pkg/log"
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

func (jg *JobGroup) registerJobMetrics(prometheusRegistry *prometheus.Registry) error {
	for _, m := range jg.metrics.allMetrics() {
		err := prometheusRegistry.Register(m)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to register job metric")
			return err
		}
	}
	return nil
}

func (jg *JobGroup) deregisterJobMetrics(prometheusRegistry *prometheus.Registry) error {
	for _, m := range jg.metrics.allMetrics() {
		deregistered := prometheusRegistry.Unregister(m)
		if !deregistered {
			log.Warn().Msg("Unable to deregister job metric")
		}
	}
	return nil
}
