package jobs

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// NewNoOpJob creates a job that does nothing.
func NewNoOpJob(name string) Job {
	return NewBasicJob(name, func(context.Context) (proto.Message, error) {
		return nil, nil
	})
}
