package backgroundscheduler

import "github.com/fluxninja/aperture/v2/pkg/jobs"

// Job is an interface that must be implemented by background jobs to get scheduled as a background job in the Circuit.
type Job interface {
	GetJob() jobs.Job
	NotifyCompletion()
}
