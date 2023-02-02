// +kubebuilder:validation:Optional
package jobs

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/status"
)

// JobCallback is the callback function that is called after a job is executed.
type JobCallback func(context.Context) (proto.Message, error)

// Job interface and basic job implementation.
type Job interface {
	// Returns the name
	Name() string
	// Executes the job
	Execute(ctx context.Context) (proto.Message, error)
	// JobWatchers
	JobWatchers() JobWatchers
}

// JobInfo contains information such as run count, last run time, etc. for a Job.
type JobInfo struct {
	LastRunTime time.Time
	NextRunTime time.Time
	RunCount    int
}

// JobBase is the base job implementation.
type JobBase struct {
	JobName string
	JWS     JobWatchers
}

// Name returns the name of the job.
func (job JobBase) Name() string {
	return job.JobName
}

// JobWatchers returns the job watchers.
func (job JobBase) JobWatchers() JobWatchers {
	return job.JWS
}

// JobConfig is config for Job
// swagger:model
// +kubebuilder:object:generate=true
type JobConfig struct {
	// Initial delay to start the job. Zero value will schedule the job immediately. Negative value will wait for next scheduled interval.
	InitialDelay config.Duration `json:"initial_delay" default:"0s"`

	// Time period between job executions. Zero or negative value means that the job will never execute periodically.
	ExecutionPeriod config.Duration `json:"execution_period" default:"10s"`

	// Execution timeout
	ExecutionTimeout config.Duration `json:"execution_timeout" validate:"gte=0s" default:"5s"`

	// Sets whether the job is initially healthy
	InitiallyHealthy bool `json:"initially_healthy" default:"false"`
}

type jobExecutor struct {
	execLock sync.Mutex
	Job
	parentRegistry   status.Registry
	livenessRegistry status.Registry
	jg               *JobGroup
	job              *gocron.Job
	config           JobConfig
	jobTag           string
	running          bool
}

// Make sure jobExecutor complies with Job interface.
var _ Job = (*jobExecutor)(nil)

func newJobExecutor(job Job, jg *JobGroup, config JobConfig, parentRegistry status.Registry) *jobExecutor {
	executor := &jobExecutor{
		Job:            job,
		jg:             jg,
		job:            &gocron.Job{},
		config:         config,
		jobTag:         job.Name(),
		parentRegistry: parentRegistry,
		running:        false,
	}
	return executor
}

// Name returns the name of the Job that the executor is associated with.
func (executor *jobExecutor) Name() string {
	return executor.Job.Name()
}

// JobWatchers returns the job watchers for the Job that the executor is associated with.
func (executor *jobExecutor) JobWatchers() JobWatchers {
	return executor.Job.JobWatchers()
}

// Execute executes the Job that the executor is associated with.
func (executor *jobExecutor) Execute(ctx context.Context) (proto.Message, error) {
	return executor.Job.Execute(ctx)
}

func (executor *jobExecutor) doJob() {
	executor.execLock.Lock()
	defer executor.execLock.Unlock()

	if !executor.running {
		return
	}

	executionTimeout := executor.config.ExecutionTimeout.AsDuration()

	ctx, cancel := context.WithCancel(context.Background())
	if executionTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, executionTimeout)
	}
	defer cancel()

	now := time.Now()
	newTime := now.Add(executionTimeout).Add(time.Second * 1)
	newDuration := newTime.Sub(now)

	jobCh := make(chan bool, 1)

	panichandler.Go(func() {
		defer func() {
			jobCh <- true
		}()
		_, err := executor.jg.gt.execute(ctx, executor)
		if err != nil {
			executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("job", executor.Name()).Msg("job status unhealthy")
			return
		}
	})

	timerCh := make(chan bool, 1)
	timer := time.AfterFunc(newDuration, func() {
		timerCh <- true
	})

	for {
		select {
		case <-timerCh:
			s := status.NewStatus(wrapperspb.String("Timeout"), errors.New("job execution timeout"))
			executor.livenessRegistry.SetStatus(s)
			timer.Reset(time.Second * 1)
		case <-jobCh:
			s := status.NewStatus(wrapperspb.String("OK"), nil)
			executor.livenessRegistry.SetStatus(s)
			timer.Stop()
			return
		}
	}
}

func (executor *jobExecutor) start() {
	executor.execLock.Lock()
	defer executor.execLock.Unlock()

	if executor.running {
		return
	}

	executor.livenessRegistry = executor.parentRegistry.Child("executor", executor.Name())

	var scheduler *gocron.Scheduler
	if executor.config.ExecutionPeriod.AsDuration() > 0 {
		scheduler = executor.jg.scheduler.
			Every(executor.config.ExecutionPeriod.AsDuration())
	} else {
		scheduler = executor.jg.scheduler.
			Every(time.Duration(math.MaxInt64))
	}

	scheduler = scheduler.
		Tag(executor.jobTag).
		SingletonMode()

	initialDelay := executor.config.InitialDelay.AsDuration()

	if initialDelay > 0 {
		scheduler = scheduler.StartAt(time.Now().Add(initialDelay))
	} else if initialDelay == 0 {
		scheduler = scheduler.StartImmediately()
	} else {
		scheduler = scheduler.WaitForSchedule()
	}

	// Scheduler.Do checks executor parameter and if the job has not been run before, it will schedule the job.
	j, err := scheduler.Do(executor.doJob)
	if err != nil {
		executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("executor", executor.Name()).Msg("Unable to schedule the job")
		return
	}
	executor.job = j
	executor.running = true
}

func (executor *jobExecutor) stop() {
	executor.execLock.Lock()
	defer executor.execLock.Unlock()

	if !executor.running {
		return
	}

	err := executor.jg.scheduler.RemoveByTag(executor.jobTag)
	if err != nil {
		executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("executor", executor.Name()).Msg("Unable to remove job")
		return
	}
	executor.livenessRegistry.Detach()
	executor.running = false
}

func (executor *jobExecutor) trigger() {
	err := executor.jg.scheduler.RunByTag(executor.jobTag)
	if err != nil {
		executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("executor", executor.Name()).Msg("Unable to trigger job")
		return
	}
}

func (executor *jobExecutor) jobInfo() *JobInfo {
	return &JobInfo{
		LastRunTime: executor.job.LastRun(),
		NextRunTime: executor.job.NextRun(),
		RunCount:    executor.job.RunCount(),
	}
}
