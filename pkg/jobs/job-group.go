package jobs

import (
	"context"
	"errors"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/fx"

	statusv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/status/v1"
	"aperture.tech/aperture/pkg/config"
	"aperture.tech/aperture/pkg/log"
	"aperture.tech/aperture/pkg/status"
)

const (
	schedulerConfigKey = "scheduler"
)

// SchedulerMode configures the scheduler's behavior when concurrency limit is applied.
type SchedulerMode int8

const (
	// RescheduleMode - the default is that if a limit on maximum
	// concurrent jobs is set and the limit is reached, a job will
	// skip it's run and try again on the next occurrence in the schedule.
	RescheduleMode SchedulerMode = iota
	// WaitMode - if a limit on maximum concurrent jobs is set
	// and the limit is reached, a job will wait to try and run
	// until a spot in the limit is freed up.
	//
	// Note: this mode can produce unpredictable results as
	// job execution order isn't guaranteed. For example, a job that
	// executes frequently may pile up in the wait queue and be executed
	// many times back to back when the queue opens.
	WaitMode
)

// JobGroupConfig holds configuration for JobGroup.
// swagger:model
type JobGroupConfig struct {
	SchedulerConfig
}

// SchedulerConfig holds configuration for job Scheduler.
// swagger:model
type SchedulerConfig struct {
	// Limits how many jobs can be running at the same time. This is useful when running resource intensive jobs and a precise start time is not critical. 0 = no limit.
	MaxConcurrentJobs int `json:"max_concurrent_jobs" validate:"gte=0" default:"0"`
}

// JobGroupConstructor holds fields to create annotated instances of JobGroup.
type JobGroupConstructor struct {
	Group         string
	GW            GroupWatchers
	DefaultConfig JobGroupConfig
	SchedulerMode SchedulerMode
}

// Annotate provides annotated instances of GroupWatcherMetrics and JobGroup.
func (jgc JobGroupConstructor) Annotate() fx.Option {
	groupTag := config.GroupTag(jgc.Group)
	nameTag := config.NameTag(jgc.Group)
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				provideGroupWatcherMetrics,
				fx.ResultTags(groupTag),
			),
		),
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
	registry *status.Registry,
	unmarshaller config.Unmarshaller,
	lifecycle fx.Lifecycle,
) (*JobGroup, error) {
	config := jgc.DefaultConfig

	key := jgc.Group + "." + schedulerConfigKey

	if err := unmarshaller.UnmarshalKey(key, &config); err != nil {
		log.Fatal().Err(err).Msg("Unable to deserialize JobGroup configuration!")
	}

	scheduler := gocron.NewScheduler(time.Local)

	scheduler.TagsUnique()

	if config.MaxConcurrentJobs > 0 {
		switch jgc.SchedulerMode {
		case RescheduleMode:
			scheduler.SetMaxConcurrentJobs(config.MaxConcurrentJobs, gocron.RescheduleMode)
		case WaitMode:
			scheduler.SetMaxConcurrentJobs(config.MaxConcurrentJobs, gocron.WaitMode)
		}
	}

	// always singleton
	scheduler.SingletonModeAll()

	gwAll := GroupWatchers{}

	if len(jgc.GW) > 0 || len(gw) > 0 {
		gwAll = append(gwAll, jgc.GW...)
		gwAll = append(gwAll, gw...)
	}

	jg, err := NewJobGroup(jgc.Group, registry, config.MaxConcurrentJobs, jgc.SchedulerMode, gwAll)
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return jg.Start()
		},
		OnStop: func(_ context.Context) error {
			return jg.Stop()
		},
	})

	return jg, nil
}

var errInitialResult = errors.New("job hasn't been scheduled yet")

// JobGroup tracks a group of jobs.
// It is responsible for scheduling jobs and keeping track of their statuses.
type JobGroup struct {
	registry  *status.Registry
	scheduler *gocron.Scheduler
	gt        *groupTracker
	name      string
}

// NewJobGroup creates a new JobGroup.
func NewJobGroup(name string,
	registry *status.Registry,
	maxConcurrentJobs int,
	schedulerMode SchedulerMode,
	gws GroupWatchers,
) (*JobGroup, error) {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.TagsUnique()

	if maxConcurrentJobs > 0 {
		switch schedulerMode {
		case RescheduleMode:
			scheduler.SetMaxConcurrentJobs(maxConcurrentJobs, gocron.RescheduleMode)
		case WaitMode:
			scheduler.SetMaxConcurrentJobs(maxConcurrentJobs, gocron.WaitMode)
		}
	}

	// always singleton
	scheduler.SingletonModeAll()
	jg := &JobGroup{
		name:      name,
		registry:  registry,
		scheduler: scheduler,
		gt:        newGroupTracker(gws, registry, name),
	}

	return jg, nil
}

// Start starts the JobGroup.
func (jg *JobGroup) Start() error {
	jg.scheduler.StartAsync()
	return nil
}

// Stop stops the JobGroup.
func (jg *JobGroup) Stop() error {
	jg.DeregisterAll()
	jg.scheduler.Stop()
	return nil
}

// GroupName returns the name of the JobGroup.
func (jg *JobGroup) GroupName() string {
	return jg.name
}

// RegisterJob registers a new Job in a JobGroup.
// It returns an error if the job is already registered.
// It also starts the job's executor.
func (jg *JobGroup) RegisterJob(job Job, config JobConfig) error {
	var initialErr error
	if !config.InitiallyHealthy {
		initialErr = errInitialResult
	}

	executor := newJobExecutor(job, jg, config)
	// setPanicHandler sets the global panicHandler to the given executor.
	executor.setPanicHandler(panicHandlerFunc(func(jobName string, recoverData interface{}) {
		log.Panic().Str("job", jobName).Interface("recoverData", recoverData).Msg("Panic recovered during job execution")
	}))
	// add to the tracker
	err := jg.gt.registerJob(executor)
	if err != nil {
		return err
	}

	// set initial status
	err = jg.gt.updateStatus(executor, status.NewStatus(nil, initialErr))
	if err != nil {
		log.Error().Err(err).Str("job", job.Name()).Str("jgName", jg.name).Msg("Unable to update status of job")
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
func (jg *JobGroup) TriggerJob(name string) {
	jg.gt.mu.Lock()
	defer jg.gt.mu.Unlock()

	tracker, ok := jg.gt.trackers[name]
	if ok {
		if executor, ok := tracker.job.(*jobExecutor); ok {
			executor.trigger()
		}
	}
}

// JobInfo returns the information related to a job with given name.
func (jg *JobGroup) JobInfo(name string) *JobInfo {
	jg.gt.mu.Lock()
	defer jg.gt.mu.Unlock()

	tracker, ok := jg.gt.trackers[name]
	if ok {
		if executor, ok := tracker.job.(*jobExecutor); ok {
			return executor.jobInfo()
		}
	}
	return nil
}

// IsHealthy returns true if the job is healthy.
func (jg *JobGroup) IsHealthy() bool {
	return jg.gt.isHealthy()
}

// Results returns the results of all jobs in the JobGroup.
func (jg *JobGroup) Results() (*statusv1.GroupStatus, bool) {
	return jg.gt.results()
}
