package jobs

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/status"
)

var (
	errInvalidJob  = errors.New("job is nil or invalid name is provided")
	errExistingJob = errors.New("job with same name already exists")
)

type jobTracker struct {
	job            Job
	statusRegistry status.Registry
}

func newJobTracker(job Job, statusRegistry status.Registry) *jobTracker {
	reg := status.NewRegistry(statusRegistry, job.Name())
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
	name           string
	groupWatchers  GroupWatchers
}

func newGroupTracker(gws GroupWatchers, statusRegistry status.Registry, name string) *groupTracker {
	reg := status.NewRegistry(statusRegistry, name)
	return &groupTracker{
		name:           name,
		trackers:       make(map[string]*jobTracker),
		statusRegistry: reg,
		groupWatchers:  gws,
	}
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

	err := gt.statusRegistry.Push(s)
	if err != nil {
		return err
	}

	return nil
}

func (gt *groupTracker) registerJob(job Job) error {
	if job.Name() == "" {
		err := errInvalidJob
		return err
	}

	gt.mu.Lock()
	defer gt.mu.Unlock()

	s := status.NewStatus(nil, nil)

	_, ok := gt.trackers[job.Name()]
	if ok {
		return errExistingJob
	}

	tracker := newJobTracker(job, gt.statusRegistry)
	gt.trackers[job.Name()] = tracker

	err := gt.statusRegistry.Push(s)
	if err != nil {
		log.Error().Err(err).Msg("failed to push job status to registry")
		return err
	}
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

	err := gt.statusRegistry.Delete()
	if err != nil {
		log.Error().Err(err).Msg("failed to delete job status")
	}

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

		err := gt.statusRegistry.Delete()
		if err != nil {
			log.Error().Err(err).Msg("failed to delete job status")
		}
	}

	gt.trackers = make(map[string]*jobTracker)

	return jobs
}

// func (gt *groupTracker) isHealthy() bool {
// 	gt.mu.Lock()
// 	defer gt.mu.Unlock()

// 	for _, tracker := range gt.trackers {
// 		gs := tracker.statusRegistry.Get()
// 		if gs.Status.GetError().GetMessage() != "" {
// 			return false
// 		}
// 	}
// 	return true
// }

func (gt *groupTracker) results() (*statusv1.GroupStatus, bool) {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	gs := &statusv1.GroupStatus{
		Groups: make(map[string]*statusv1.GroupStatus, len(gt.trackers)),
	}

	healthy := true

	for name, tracker := range gt.trackers {
		tgs := tracker.statusRegistry.Get()
		if tgs == nil {
			log.Debug().Str("path", tracker.job.Name()).Msg("returned nil status")
			continue
		}

		gs.Groups[name] = tgs
		if tgs.Status.GetError().GetMessage() != "" {
			healthy = false
		}
	}

	return gs, healthy
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
		log.Error().Err(err).Str("job", job.Name()).Msg("Recently completed job has been removed from tracker and is not reporting results")
		return nil, err
	}

	jobStats := JobStats{Duration: duration}
	job.JobWatchers().OnJobCompleted(s, jobStats)
	gt.groupWatchers.OnJobCompleted(job.Name(), s, jobStats)

	return details, err
}
