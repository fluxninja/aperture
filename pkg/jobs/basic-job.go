package jobs

import (
	"context"
	"errors"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"google.golang.org/protobuf/proto"
)

// Make sure basicJob complies with Job interface.
var _ Job = (*basicJob)(nil)

// basicJob is a basic job that every other job builds on.
type basicJob struct {
	JobFunc JobCallback
	JobBase
}

// NewBasicJob is a constructor for basicJob struct.
func NewBasicJob(name string, jobFunc JobCallback) Job {
	job := &basicJob{
		JobBase: JobBase{JobName: name},
		JobFunc: jobFunc,
	}
	if jobFunc == nil {
		log.Warn().Msgf("Provided job callback is nil for job: %s", name)
	}
	return job
}

// Name returns the name of the job.
func (job *basicJob) Name() string {
	return job.JobBase.Name()
}

// JobWatchers returns the job watchers.
func (job *basicJob) JobWatchers() JobWatchers {
	return job.JobBase.JobWatchers()
}

// Execute executes the job.
func (job *basicJob) Execute(ctx context.Context) (proto.Message, error) {
	if job.JobFunc == nil {
		return nil, errors.New("JobFunc not provided")
	}
	return job.JobFunc(ctx)
}
