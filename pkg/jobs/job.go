// +kubebuilder:validation:Optional
package jobs

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/status"
	"github.com/reugn/go-quartz/quartz"
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
	job              quartz.Job
	config           JobConfig
	running          bool
}

// Make sure jobExecutor complies with Job interface.
var _ Job = (*jobExecutor)(nil)

func newJobExecutor(job Job, jg *JobGroup, config JobConfig, parentRegistry status.Registry) *jobExecutor {
	executor := &jobExecutor{
		Job:            job,
		jg:             jg,
		config:         config,
		parentRegistry: parentRegistry,
		running:        false,
	}
	executor.job = quartz.NewIsolatedJob(quartz.NewFunctionJobWithDesc(job.Name(), executor.doJob))
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

func (executor *jobExecutor) doJob(ctx context.Context) (proto.Message, error) {
	executor.execLock.Lock()
	defer executor.execLock.Unlock()

	if !executor.running {
		return nil, fmt.Errorf("job %s is not running", executor.Name())
	}
	now := time.Now()

	executionTimeout := executor.config.ExecutionTimeout.AsDuration()

	ctx, cancel := context.WithCancel(context.Background())
	if executionTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, executionTimeout)
	}
	defer cancel()

	newTime := now.Add(executionTimeout).Add(time.Second * 1)
	newDuration := newTime.Sub(now)

	jobCh := make(chan bool, 1)

	var msg proto.Message
	var err error
	panichandler.Go(func() {
		defer func() {
			jobCh <- true
		}()
		msg, err = executor.jg.gt.execute(ctx, executor)
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
			return msg, err
		}
	}
}

func (executor *jobExecutor) start() {
	executor.execLock.Lock()
	defer executor.execLock.Unlock()

	if executor.running {
		return
	}
	executor.running = true

	executor.livenessRegistry = executor.parentRegistry.Child("executor", executor.Name())

	execPeriod := executor.config.ExecutionPeriod.AsDuration()
	if execPeriod < 0 {
		// no need to schedule a job that will never run periodically
		return
	}

	trigger := quartz.NewSimpleTrigger(execPeriod)

	err := executor.jg.scheduler.ScheduleJob(context.TODO(), executor.job, trigger)
	if err != nil {
		executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("executor", executor.Name()).Msg("Unable to schedule the job")
		return
	}
}

func (executor *jobExecutor) stop() {
	executor.execLock.Lock()
	defer executor.execLock.Unlock()

	if !executor.running {
		return
	}

	// no need to DeleteJob if the job is not scheduled periodically.
	if executor.config.ExecutionPeriod.AsDuration() > 0 {
		err := executor.jg.scheduler.DeleteJob(executor.job.Key())
		if err != nil {
			executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("executor", executor.Name()).Msg("Unable to remove job")
			return
		}
	}
	executor.livenessRegistry.Detach()
	executor.running = false
}

func (executor *jobExecutor) trigger(delay time.Duration) {
	trigger := quartz.NewRunOnceTrigger(delay)
	err := executor.jg.scheduler.ScheduleJob(context.TODO(), executor.job, trigger)
	if err != nil {
		executor.jg.gt.statusRegistry.GetLogger().Error().Err(err).Str("executor", executor.Name()).Msg("Unable to trigger the job")
		return
	}
}
