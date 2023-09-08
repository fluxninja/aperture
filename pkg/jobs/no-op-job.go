package jobs

import (
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NewNoOpJob creates a job that does nothing.
func NewNoOpJob(name string) Job {
	return NewBasicJob(name, func(context.Context) (proto.Message, error) {
		return &emptypb.Empty{}, nil
	})
}
