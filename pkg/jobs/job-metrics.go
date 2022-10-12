package jobs

import (
	"fmt"

	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// JobMetrics holds prometheus metrics related to job.
type JobMetrics struct {
	latencySummary *prometheus.SummaryVec
}

// NewJobMetrics returns JobMetrics with prometheus summary metrics.
func NewJobMetrics() JobMetrics {
	jobMetricsLabels := []string{metrics.JobNameLabel, metrics.JobOkLabel}
	return JobMetrics{
		latencySummary: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: metrics.JobLatencyMetricName,
			Help: "Latency summary of the job execution",
		}, jobMetricsLabels),
	}
}

func (jm *JobMetrics) registerJobMetrics(reg *prometheus.Registry) error {
	err := reg.Register(jm.latencySummary)
	if err != nil {
		// Ignore already registered error, as this is not harmful. Metrics may
		// be registered by other running jobs.
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return fmt.Errorf("failed to register job metrics: %w", err)
		}
	}
	return nil
}

func (jm *JobMetrics) unregisterJobMetrics(reg *prometheus.Registry) error {
	if !reg.Unregister(jm.latencySummary) {
		return fmt.Errorf("failed to unregister %s metrics", metrics.JobLatencyMetricName)
	}
	return nil
}

func (jm *JobMetrics) removeMetrics(nameLabel string) error {
	label := prometheus.Labels{
		metrics.JobNameLabel: nameLabel,
	}
	deleted := jm.latencySummary.DeletePartialMatch(label)
	if deleted == 0 {
		return fmt.Errorf("failed to remove job metrics, %s", nameLabel)
	}
	return nil
}
