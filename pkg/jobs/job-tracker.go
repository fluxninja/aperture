package jobs

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	errInvalidJob  = errors.New("job is nil or invalid name is provided")
	errExistingJob = errors.New("job with same name already exists")
)

type jobTracker struct {
	job            Job
	statusRegistry status.Registry
}

func newJobTracker(job Job, statusRegistry status.Registry, trackerMetrics *JobMetrics) (*jobTracker, error) {
	reg := statusRegistry.Child(job.Name())

	metricLabels := make(prometheus.Labels)
	metricLabels[metrics.JobNameLabel] = job.Name()
	metricLabels[metrics.JobOkLabel] = "true"
	_, err := trackerMetrics.latencySummary.GetMetricWith(metricLabels)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create job ok latency summary metrics")
		return nil, err
	}

	metricLabels[metrics.JobOkLabel] = "false"
	_, err = trackerMetrics.latencySummary.GetMetricWith(metricLabels)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create job not ok latency summary metrics")
		return nil, err
	}

	jobTracker := &jobTracker{
		job:            job,
		statusRegistry: reg,
	}

	return jobTracker, nil
}

// Common groupTracker.
type groupTracker struct {
	mu                 sync.Mutex
	trackers           map[string]*jobTracker
	statusRegistry     status.Registry
	groupWatchers      GroupWatchers
	jobTrackerMetrics  *JobMetrics
	prometheusRegistry *prometheus.Registry
}

func newGroupTracker(gws GroupWatchers, statusRegistry status.Registry, prometheusRegistry *prometheus.Registry) *groupTracker {
	gt := &groupTracker{
		trackers:           make(map[string]*jobTracker),
		statusRegistry:     statusRegistry,
		groupWatchers:      gws,
		jobTrackerMetrics:  newJobMetrics(),
		prometheusRegistry: prometheusRegistry,
	}
	err := gt.jobTrackerMetrics.registerJobMetrics(gt.prometheusRegistry)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register job metrics")
	}
	return gt
}

func (gt *groupTracker) updateStatus(job Job, s *statusv1.Status) error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	// check whether this job still exists and hasn't been swapped with another job of the same name
	tracker, ok := gt.trackers[job.Name()]
	if !ok {
		return errInvalidJob
	}

	if tracker.job != job {
		return errExistingJob
	}

	tracker.statusRegistry.SetStatus(s)

	return nil
}

func (gt *groupTracker) registerJob(job Job) error {
	if job.Name() == "" {
		err := errInvalidJob
		return err
	}

	gt.mu.Lock()
	defer gt.mu.Unlock()

	_, ok := gt.trackers[job.Name()]
	if ok {
		return errExistingJob
	}

	tracker, err := newJobTracker(job, gt.statusRegistry, gt.jobTrackerMetrics)
	if err != nil {
		return err
	}

	gt.trackers[job.Name()] = tracker

	gt.groupWatchers.OnJobRegistered(job.Name())

	return nil
}

func (gt *groupTracker) deregisterJob(name string) (Job, error) {
	var ok bool
	var tracker *jobTracker

	gt.mu.Lock()
	defer gt.mu.Unlock()

	tracker, ok = gt.trackers[name]
	if !ok {
		return nil, errInvalidJob
	} else {
		err := gt.jobTrackerMetrics.removeMetrics(name)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove metrics")
			return nil, err
		}
	}

	delete(gt.trackers, name)
	gt.groupWatchers.OnJobDeregistered(name)

	tracker.statusRegistry.Detach()

	return tracker.job, nil
}

func (gt *groupTracker) reset() []Job {
	jobs := []Job{}

	gt.mu.Lock()
	defer gt.mu.Unlock()
	_ = gt.jobTrackerMetrics.unregisterJobMetrics(gt.prometheusRegistry)
	for _, tracker := range gt.trackers {
		job := tracker.job
		jobs = append(jobs, job)
		gt.groupWatchers.OnJobDeregistered(job.Name())
		err := gt.jobTrackerMetrics.removeMetrics(job.Name())
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove metrics")
		}
		tracker.statusRegistry.Detach()
	}
	gt.trackers = make(map[string]*jobTracker)

	return jobs
}

func (gt *groupTracker) isHealthy() bool {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	for _, tracker := range gt.trackers {
		if tracker.statusRegistry.GetStatus().GetError() != nil {
			return false
		}
	}
	return true
}

func (gt *groupTracker) results() (*statusv1.GroupStatus, bool) {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	return gt.statusRegistry.GetGroupStatus(), !gt.statusRegistry.HasError()
}

func (gt *groupTracker) getJobs() []Job {
	jobs := []Job{}

	gt.mu.Lock()
	defer gt.mu.Unlock()

	for _, tracker := range gt.trackers {
		job := tracker.job
		jobs = append(jobs, job)
	}

	return jobs
}

func (gt *groupTracker) execute(ctx context.Context, job Job) (proto.Message, error) {
	gt.groupWatchers.OnJobScheduled(job.Name())
	job.JobWatchers().OnJobScheduled()

	startTime := time.Now()
	details, err := job.Execute(ctx)
	if err != nil {
		return nil, err
	}
	duration := time.Since(startTime)

	s := status.NewStatus(details, err)
	err = gt.updateStatus(job, s)
	if err != nil {
		gt.statusRegistry.GetLogger().Error().Err(err).Str("job", job.Name()).Msg("Recently completed job has been removed from tracker and is not reporting results")
		return nil, err
	}

	jobStats := JobStats{Duration: duration}
	job.JobWatchers().OnJobCompleted(s, jobStats)
	gt.groupWatchers.OnJobCompleted(job.Name(), s, jobStats)

	return details, err
}
