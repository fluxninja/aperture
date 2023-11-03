package runtime

import "github.com/fluxninja/aperture/v2/pkg/jobs"

// BackgroundJob is an interface that must be implemented by background jobs to get scheduled as a background job in the Circuit.
type BackgroundJob interface {
	GetJob() jobs.Job
	NotifyCompletion()
}
