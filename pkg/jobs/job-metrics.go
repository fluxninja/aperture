package jobs

import (
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
)

// provideGroupWatcherMetrics provides a GroupWatcherMetrics.
func provideGroupWatcherMetrics(registry *prometheus.Registry, router *mux.Router) *GroupWatcherMetrics {
	gwm, _ := NewGroupWatcherMetrics(registry)
	gwm.registerGroupWatcherMetrics(router)
	return gwm
}

// GroupWatcherMetrics reports metrics when jobs are registered, deregistered, scheduled, and completed in each job group.
type GroupWatcherMetrics struct {
	metricsRegistry *prometheus.Registry
	jobRegistered   prometheus.CounterVec
	jobScheduled    prometheus.CounterVec
	jobCompleted    prometheus.CounterVec
	jobLatency      prometheus.HistogramVec
}

// NewGroupWatcherMetrics creates a new GroupWatcherMetrics.
func NewGroupWatcherMetrics(registry *prometheus.Registry) (*GroupWatcherMetrics, error) {
	defaultGroupLabels := []string{"group_name"}
	groupWatcherMetrics := &GroupWatcherMetrics{
		metricsRegistry: registry,
		jobRegistered: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "group_job_registered_number",
				Help: "Current number of group job registered",
			}, defaultGroupLabels,
		),
		jobScheduled: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "group_job_scheduled_number",
				Help: "Current number of group job scheduled",
			}, defaultGroupLabels,
		),
		jobCompleted: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "group_job_completed_total",
				Help: "Total number of group job completed",
			}, []string{"group_name", "group_job_completed_healthy"},
		),
		jobLatency: *prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "group_job_latency_seconds",
				Help:    "The latency of the group jobs",
				Buckets: prometheus.DefBuckets,
			}, defaultGroupLabels,
		),
	}

	err := groupWatcherMetrics.metricsRegistry.Register(groupWatcherMetrics.jobRegistered)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to register group job watcher metrics")
		return nil, err
	}

	err = groupWatcherMetrics.metricsRegistry.Register(groupWatcherMetrics.jobScheduled)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to register group job watcher metrics")
		return nil, err
	}

	err = groupWatcherMetrics.metricsRegistry.Register(groupWatcherMetrics.jobCompleted)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to register group job watcher metrics")
		return nil, err
	}

	err = groupWatcherMetrics.metricsRegistry.Register(groupWatcherMetrics.jobLatency)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to register group job watcher metrics")
		return nil, err
	}

	return groupWatcherMetrics, nil
}

// OnJobRegistered increments 'jobRegistered' CounterVec with label 'name'.
func (g *GroupWatcherMetrics) OnJobRegistered(name string) {
	g.jobRegistered.WithLabelValues(name).Inc()
}

// OnJobDeregistered does nothing.
func (g *GroupWatcherMetrics) OnJobDeregistered(name string) {
}

// OnJobScheduled increments 'jobScheduled' CounterVec with label 'name'.
func (g *GroupWatcherMetrics) OnJobScheduled(name string) {
	g.jobScheduled.WithLabelValues(name).Inc()
}

// OnJobCompleted increments 'jobCompleted' CounterVec with label 'name' and 'healthy' label.
// It also records the latency of the job.
func (g *GroupWatcherMetrics) OnJobCompleted(name string, status *statusv1.Status, duration time.Duration) {
	g.jobCompleted.WithLabelValues(name, strconv.FormatBool(status.GetError().Message == "")).Inc()
	g.jobLatency.WithLabelValues(name).Observe(duration.Seconds())
}

func (g *GroupWatcherMetrics) registerGroupWatcherMetrics(router *mux.Router) {
	metrics.RegisterMetricsHandler(router, g.metricsRegistry)
}
