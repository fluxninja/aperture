package jobs

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

var (
	errInvalidJob  = errors.New("job is nil or invalid name is provided")
	errExistingJob = errors.New("job with same name already exists")
)

// JobInfo contains information such as run count, last run time, etc. for a Job.
type JobInfo struct {
	LastExecuteTime time.Time
	ExecuteCount    int
}

type jobTracker struct {
	job            Job
	statusRegistry status.Registry
	jobInfo        JobInfo
}

func newJobTracker(job Job, statusRegistry status.Registry) *jobTracker {
	reg := statusRegistry.Child("job-tracker", job.Name())
	return &jobTracker{
		job:            job,
		statusRegistry: reg,
	}
}

// Common groupTracker.
type groupTracker struct {
	mu             sync.Mutex
	trackers       map[string]*jobTracker
	statusRegistry status.Registry
	groupWatchers  GroupWatchers
}

func newGroupTracker(gws GroupWatchers, statusRegistry status.Registry) *groupTracker {
	gt := &groupTracker{
		trackers:       make(map[string]*jobTracker),
		statusRegistry: statusRegistry,
		groupWatchers:  gws,
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

	tracker.statusRegistry.SetStatus(s, nil)
	return nil
}

func (gt *groupTracker) registerJob(job Job) error {
	if job.Name() == "" {
		return errInvalidJob
	}

	gt.mu.Lock()
	defer gt.mu.Unlock()

	_, ok := gt.trackers[job.Name()]
	if ok {
		return errExistingJob
	}

	tracker := newJobTracker(job, gt.statusRegistry)
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

	for _, tracker := range gt.trackers {
		job := tracker.job
		jobs = append(jobs, job)
		gt.groupWatchers.OnJobDeregistered(job.Name())

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

func (gt *groupTracker) updateJobInfo(job Job, startTime time.Time) error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	tracker, ok := gt.trackers[job.Name()]
	if !ok {
		return errInvalidJob
	}

	tracker.jobInfo.LastExecuteTime = startTime
	tracker.jobInfo.ExecuteCount++
	return nil
}

func (gt *groupTracker) execute(ctx context.Context, job Job) (proto.Message, error) {
	gt.groupWatchers.OnJobScheduled(job.Name())
	job.JobWatchers().OnJobScheduled()

	startTime := time.Now()

	err := gt.updateJobInfo(job, startTime)
	if err != nil {
		return nil, err
	}

	details, err := job.Execute(ctx)
	if err != nil {
		return nil, err
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)

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
