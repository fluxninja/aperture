package backgroundscheduler

import (
	"errors"
	"time"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Scheduler is an interface that allows background jobs to be scheduled in a Circuit.
type Scheduler interface {
	RegisterJob(job Job)
}

type jobsMap map[string]Job

type scheduler struct {
	// CircuitAPI
	circuitAPI runtime.CircuitAPI
	// inflightJobs contains a Job Result Broker for each job in the multi job
	inflightJobs jobsMap
	// pendingJobs contains a Job Result Broker for each job in the multi job
	pendingJobs jobsMap
	// Prom Multi Job
	multiJob *jobs.MultiJob
	// Job group
	jobGroup *jobs.JobGroup
	// Query job state
	jobRunning bool
}

// Make sure scheduler complies with the Scheduler interface.
var _ Scheduler = (*scheduler)(nil)

// Make sure scheduler complies with the jobs.JobsWatcher interface.
var _ jobs.JobWatcher = (*scheduler)(nil)

// RegisterJob registers a job with the scheduler.
func (scheduler *scheduler) RegisterJob(job Job) {
	scheduler.pendingJobs[job.GetJob().Name()] = job
}

// OnJobScheduled is called when the scheduler.multiJob is scheduled.
func (scheduler *scheduler) OnJobScheduled() {
}

// OnJobCompleted is called when the scheduler.promMultiJob is completed.
func (scheduler *scheduler) OnJobCompleted(_ *statusv1.Status, _ jobs.JobStats) {
	// Take circuit execution lock
	scheduler.circuitAPI.LockExecution()
	defer scheduler.circuitAPI.UnlockExecution()

	// Provide results via callbacks
	for _, jobResBroker := range scheduler.inflightJobs {
		jobResBroker.NotifyCompletion()
	}
	// Reset inflightJobs
	scheduler.inflightJobs = make(jobsMap)
	scheduler.jobRunning = false
}

func (scheduler *scheduler) onTickEnd(_ runtime.TickInfo) (err error) {
	logger := scheduler.circuitAPI.GetStatusRegistry().GetLogger()
	// Already under circuit execution lock
	// Launch job only if previous one is completed
	if scheduler.jobRunning {
		err = errors.New("previous job is still running")
	} else {
		scheduler.jobRunning = true
		// Remove all the previous jobs in the multi job
		scheduler.multiJob.DeregisterAll()
		// Add all the pendingJobs to the multijob and trigger it
		for _, jobResBroker := range scheduler.pendingJobs {
			job := jobResBroker.GetJob()
			err = scheduler.multiJob.RegisterJob(job)
			if err != nil {
				logger.Error().Err(err).Str("job", job.Name()).Msg("Error registering job")
				return err
			}
		}
		// Move pendingJobs to inflightJobs
		scheduler.inflightJobs = scheduler.pendingJobs
		// Clear pendingJobs for future ticks
		scheduler.pendingJobs = make(jobsMap)
		// Trigger the multi job
		scheduler.jobGroup.TriggerJob(scheduler.multiJob.Name(), time.Duration(0))
	}
	return err
}
