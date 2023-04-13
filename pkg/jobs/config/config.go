// +kubebuilder:validation:Optional
package config

import "github.com/fluxninja/aperture/pkg/config"

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

// JobGroupConfig holds configuration for JobGroup.
// swagger:model
// +kubebuilder:object:generate=true
type JobGroupConfig struct {
	SchedulerConfig `json:",inline"`
}

// SchedulerConfig holds configuration for job Scheduler.
// swagger:model
// +kubebuilder:object:generate=true
type SchedulerConfig struct {
	// When true, the scheduler will run jobs synchronously,
	// waiting for each execution instance of the job to return
	// before starting the next execution. Running with this
	// option effectively serializes all job execution.
	BlockingExecution bool `json:"blocking_execution" default:"false"`

	// Limits how many jobs can be running at the same time. This is
	// useful when running resource intensive jobs and a precise start time is
	// not critical. 0 = no limit. If BlockingExecution is set, then WorkerLimit
	// is ignored.
	WorkerLimit int `json:"worker_limit" default:"0"`
}
