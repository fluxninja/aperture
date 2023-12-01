package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/reugn/go-quartz/quartz"
	"go.uber.org/fx"

	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	jobsconfig "github.com/fluxninja/aperture/v2/pkg/jobs/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

const (
	schedulerConfigKey = "scheduler"
)

// JobGroupConfig is reexported as it is commonly used when importing jobs package.
type JobGroupConfig = jobsconfig.JobGroupConfig

// JobGroupConstructor holds fields to create annotated instances of JobGroup.
type JobGroupConstructor struct {
	// Name of the job group - config key is <name> and statuses are updated under <name>.<job>
	Name string
	// Config key --  if it is empty then it is <name>.scheduler
	Key           string
	GW            GroupWatchers
	DefaultConfig jobsconfig.JobGroupConfig
}

// Annotate provides annotated instances of JobGroup.
func (jgc JobGroupConstructor) Annotate() fx.Option {
	groupTag := config.GroupTag(jgc.Name)
	nameTag := config.NameTag(jgc.Name)
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				jgc.provideJobGroup,
				fx.ParamTags(groupTag),
				fx.ResultTags(nameTag),
			),
		),
	)
}

func (jgc JobGroupConstructor) provideJobGroup(
	gw GroupWatchers,
	registry status.Registry,
	unmarshaller config.Unmarshaller,
	lifecycle fx.Lifecycle,
) (*JobGroup, error) {
	schedulerConfig := jgc.DefaultConfig

	key := jgc.Key
	if key == "" {
		key = jgc.Name + "." + schedulerConfigKey
	}

	if err := unmarshaller.UnmarshalKey(key, &schedulerConfig); err != nil {
		log.Panic().Err(err).Msg("Unable to deserialize JobGroup configuration!")
	}

	gwAll := GroupWatchers{}
	if len(jgc.GW) > 0 || len(gw) > 0 {
		gwAll = append(gwAll, jgc.GW...)
		gwAll = append(gwAll, gw...)
	}
	reg := registry.Child("job-group", jgc.Name)

	jg, err := NewJobGroup(reg, schedulerConfig, gwAll)
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return jg.Start()
		},
		OnStop: func(_ context.Context) error {
			defer reg.Detach()
			return jg.Stop()
		},
	})

	return jg, nil
}

var errInitialResult = errors.New("job hasn't been scheduled yet")

// JobGroup tracks a group of jobs.
// It is responsible for scheduling jobs and keeping track of their statuses.
type JobGroup struct {
	scheduler        *quartz.StdScheduler
	gt               *groupTracker
	livenessRegistry status.Registry
}

// NewJobGroup creates a new JobGroup.
func NewJobGroup(
	statusRegistry status.Registry,
	config jobsconfig.JobGroupConfig,
	gws GroupWatchers,
) (*JobGroup, error) {
	scheduler := quartz.NewStdSchedulerWithOptions(quartz.StdSchedulerOptions{
		BlockingExecution: config.BlockingExecution,
		WorkerLimit:       config.WorkerLimit,
	})

	jg := &JobGroup{
		scheduler: scheduler,
		gt:        newGroupTracker(gws, statusRegistry),
	}

	return jg, nil
}

// Start starts the JobGroup.
func (jg *JobGroup) Start() error {
	jg.livenessRegistry = jg.gt.statusRegistry.Root().
		Child("system", "liveness").
		Child("job-group", "job_groups").
		Child(jg.gt.statusRegistry.Key(), jg.gt.statusRegistry.Value())
	jg.scheduler.Start(context.Background())
	return nil
}

// Stop stops the JobGroup.
func (jg *JobGroup) Stop() error {
	jg.DeregisterAll()
	jg.scheduler.Stop()
	jg.scheduler.Wait(context.Background())
	jg.livenessRegistry.Detach()
	return nil
}

// RegisterJob registers a new Job in a JobGroup.
// It returns an error if the job is already registered.
// It also starts the job's executor.
func (jg *JobGroup) RegisterJob(job Job, config jobsconfig.JobConfig) error {
	var initialErr error
	if !config.InitiallyHealthy {
		initialErr = errInitialResult
	}

	executor := newJobExecutor(job, jg, config, jg.livenessRegistry)
	// add to the tracker
	err := jg.gt.registerJob(executor)
	if err != nil {
		return err
	}

	// set initial status
	err = jg.gt.updateStatus(executor, status.NewStatus(nil, initialErr))
	if err != nil {
		jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("job", job.Name()).Msg("Unable to update status of job")
		return err
	}

	// start the executor
	executor.start()

	return nil
}

// DeregisterJob deregisters a Job from the JobGroup.
// It returns an error if the job is not registered.
// It also stops the job's executor.
func (jg *JobGroup) DeregisterJob(name string) error {
	job, err := jg.gt.deregisterJob(name)
	if err != nil {
		return err
	}
	if executor, ok := job.(*jobExecutor); ok {
		executor.stop()
	}
	return nil
}

// DeregisterAll deregisters all Jobs from the JobGroup.
func (jg *JobGroup) DeregisterAll() {
	jobs := jg.gt.reset()
	for _, job := range jobs {
		if executor, ok := job.(*jobExecutor); ok {
			executor.stop()
		}
	}
}

// TriggerJob triggers a Job in the JobGroup.
func (jg *JobGroup) TriggerJob(name string, delay time.Duration) {
	jg.gt.mu.Lock()
	defer jg.gt.mu.Unlock()

	tracker, ok := jg.gt.trackers[name]
	if ok {
		if executor, ok := tracker.job.(*jobExecutor); ok {
			executor.trigger(delay)
		}
	}
}

// JobInfo returns the information related to a job with given name.
func (jg *JobGroup) JobInfo(name string) (JobInfo, error) {
	jg.gt.mu.Lock()
	defer jg.gt.mu.Unlock()

	tracker, ok := jg.gt.trackers[name]
	if ok {
		return tracker.jobInfo, nil
	}
	return JobInfo{}, fmt.Errorf("job %s not found", name)
}

// IsHealthy returns true if the job is healthy.
func (jg *JobGroup) IsHealthy() bool {
	return jg.gt.isHealthy()
}

// Results returns the results of all jobs in the JobGroup.
func (jg *JobGroup) Results() (*statusv1.GroupStatus, bool) {
	return jg.gt.results()
}

// GetStatusRegistry returns the registry of the JobGroup.
func (jg *JobGroup) GetStatusRegistry() status.Registry {
	return jg.gt.statusRegistry
}
