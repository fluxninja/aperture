package jobs

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
)

// Make sure BasicJob complies with Job interface.
var _ Job = (*BasicJob)(nil)

// BasicJob is a basic job that every other job builds on.
type BasicJob struct {
	JobFunc JobCallback
	JobBase
}

// Name returns the name of the job.
func (job *BasicJob) Name() string {
	return job.JobBase.Name()
}

// GetJobFunc returns the function of the job.
func (job *BasicJob) GetJobFunc() (JobCallback, error) {
	return job.JobFunc, nil
}

// JobWatchers returns the job watchers.
func (job *BasicJob) JobWatchers() JobWatchers {
	return job.JobBase.JobWatchers()
}

// Execute executes the job.
func (job *BasicJob) Execute(ctx context.Context) (proto.Message, error) {
	if job.JobFunc == nil {
		return nil, errors.New("JobFunc not provided")
	}
	return job.JobFunc(ctx)
}
