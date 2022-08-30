package jobs

import (
	"context"
	"errors"
	"strings"
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
	job Job
}

func newJobTracker(job Job) *jobTracker {
	return &jobTracker{
		job: job,
	}
}

// Common groupTracker.
type groupTracker struct {
	mu            sync.Mutex
	trackers      map[string]*jobTracker
	registry      *status.Registry
	name          string
	groupWatchers GroupWatchers
}

func newGroupTracker(gws GroupWatchers, registry *status.Registry, name string) *groupTracker {
	return &groupTracker{
		name:          name,
		trackers:      make(map[string]*jobTracker),
		registry:      registry,
		groupWatchers: gws,
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

	regPath := gt.getStatusRegPath(job.Name())
	err := gt.registry.Push(regPath, s)
	if err != nil {
		return err
	}

	return nil
}

func (gt *groupTracker) getStatusRegPath(jobName string) string {
	return strings.Join([]string{gt.name, jobName}, gt.registry.Delim())
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

	tracker := newJobTracker(job)
	gt.trackers[job.Name()] = tracker

	regPath := gt.getStatusRegPath(job.Name())
	err := gt.registry.Push(regPath, s)
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

	regPath := gt.getStatusRegPath(name)
	gt.registry.Delete(regPath)

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

		regPath := gt.getStatusRegPath(job.Name())
		gt.registry.Delete(regPath)
	}

	gt.trackers = make(map[string]*jobTracker)

	return jobs
}

func (gt *groupTracker) isHealthy() bool {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	for _, tracker := range gt.trackers {
		regPath := gt.getStatusRegPath(tracker.job.Name())
		gs := gt.registry.Get(regPath)
		if gs == nil {
			log.Debug().Str("path", regPath).Msg("returned nil status")
			continue
		}

		if gs.Status.GetError().GetMessage() != "" {
			return false
		}
	}
	return true
}

func (gt *groupTracker) results() (*statusv1.GroupStatus, bool) {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	gs := &statusv1.GroupStatus{
		Groups: make(map[string]*statusv1.GroupStatus, len(gt.trackers)),
	}

	healthy := true

	for name, tracker := range gt.trackers {
		regPath := gt.getStatusRegPath(tracker.job.Name())
		tgs := gt.registry.Get(regPath)
		if tgs == nil {
			log.Debug().Str("path", regPath).Msg("returned nil status")
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
