package runtime

import (
	"context"
	"errors"
	"time"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/status"
	"go.uber.org/fx"
)

var circuitBackgroundJobGroupTag = iface.PoliciesRoot + "circuit_background_jobs"

// BackgroundSchedulerModule returns fx options for PromQL in the main app.
func BackgroundSchedulerModule() fx.Option {
	return fx.Options(
		jobs.JobGroupConstructor{Name: circuitBackgroundJobGroupTag, Key: iface.PoliciesRoot + ".promql_jobs_scheduler"}.Annotate(),
		fx.Provide(fx.Annotate(
			provideFxOptionsFunc,
			fx.ParamTags(config.NameTag(circuitBackgroundJobGroupTag)),
			fx.ResultTags(iface.FxOptionsFuncTag),
		)),
	)
}

func provideFxOptionsFunc(promQLJobGroup *jobs.JobGroup) notifiers.FxOptionsFunc {
	return func(key notifiers.Key, _ config.Unmarshaller, _ status.Registry) (fx.Option, error) {
		return fx.Supply(fx.Annotated{Name: circuitBackgroundJobGroupTag, Target: promQLJobGroup}), nil
	}
}

// BackgroundSchedulerModuleForPolicyApp returns fx options for Scheduler in the policy app. Invoked only once per policy.
func BackgroundSchedulerModuleForPolicyApp(circuitAPI CircuitSuperAPI) fx.Option {
	provideScheduler := func(jobGroup *jobs.JobGroup, lifecycle fx.Lifecycle) (*backgroundScheduler, error) {
		// Create this as a singleton at the policy/circuit level
		scheduler := &backgroundScheduler{
			circuitAPI:   circuitAPI,
			inflightJobs: make(jobsMap),
			pendingJobs:  make(jobsMap),
			jobGroup:     jobGroup,
		}
		// Register TickEndCallback
		circuitAPI.RegisterTickEndCallback(scheduler.onTickEnd)

		var jws []jobs.JobWatcher
		jws = append(jws, scheduler)

		// Create backgroundMultiJob for running background jobs in this circuit
		backgroundMultiJob := jobs.NewMultiJob(jobGroup.GetStatusRegistry().Child("policy", circuitAPI.GetPolicyName()), jws, nil)
		scheduler.multiJob = backgroundMultiJob

		executionPeriod := config.MakeDuration(-1)
		// Execution timeout for background jobs
		executionTimeout := config.MakeDuration(time.Second * 60)
		jobConfig := jobs.JobConfig{
			InitiallyHealthy: true,
			ExecutionPeriod:  executionPeriod,
			ExecutionTimeout: executionTimeout,
		}

		// Lifecycle hooks to register and deregister this circuit's backgroundMultiJob in jobGroup
		lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				// Register multi job with job group
				err := jobGroup.RegisterJob(backgroundMultiJob, jobConfig)
				return err
			},
			OnStop: func(_ context.Context) error {
				// Deregister multi job from job group
				err := jobGroup.DeregisterJob(backgroundMultiJob.Name())
				return err
			},
		})
		return scheduler, nil
	}

	return fx.Options(
		fx.Provide(fx.Annotate(
			provideScheduler,
			fx.ParamTags(config.NameTag(circuitBackgroundJobGroupTag)),
		)),
	)
}

type jobsMap map[string]BackgroundJob

type backgroundScheduler struct {
	// CircuitAPI
	circuitAPI CircuitSuperAPI
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

// Make sure scheduler complies with the jobs.JobsWatcher interface.
var _ jobs.JobWatcher = (*backgroundScheduler)(nil)

// ScheduleJob schedules a job using the background scheduler.
func (scheduler *backgroundScheduler) scheduleJob(job BackgroundJob) {
	scheduler.pendingJobs[job.GetJob().Name()] = job
}

// OnJobScheduled is called when the scheduler.multiJob is scheduled.
func (scheduler *backgroundScheduler) OnJobScheduled() {
}

// OnJobCompleted is called when the scheduler.promMultiJob is completed.
func (scheduler *backgroundScheduler) OnJobCompleted(_ *statusv1.Status, _ jobs.JobStats) {
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

func (scheduler *backgroundScheduler) onTickEnd(_ CircuitAPI) (err error) {
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
