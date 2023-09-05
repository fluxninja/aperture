package backgroundscheduler

import (
	"context"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/jobs"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/status"
	"go.uber.org/fx"
)

var circuitBackgroundJobGroupTag = iface.PoliciesRoot + "circuit_background_jobs"

// Module returns fx options for PromQL in the main app.
func Module() fx.Option {
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

// ModuleForPolicyApp returns fx options for Scheduler in the policy app. Invoked only once per policy.
func ModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	provideScheduler := func(jobGroup *jobs.JobGroup, lifecycle fx.Lifecycle) (*scheduler, error) {
		// Create this as a singleton at the policy/circuit level
		scheduler := &scheduler{
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
		// Do not want this jobs to timeout
		executionTimeout := config.MakeDuration(0)
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
