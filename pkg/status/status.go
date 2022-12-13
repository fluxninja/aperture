package status

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/status/v1"
)

// NewStatus creates a new instance of Status to be pushed into status registry. Use this function for creating status instead of by hand.
func NewStatus(d proto.Message, e error) *statusv1.Status {
	s := &statusv1.Status{
		Timestamp: timestamppb.Now(),
	}

	if d != nil {
		messageAny, err := anypb.New(d)
		if err != nil {
			return nil
		}
		s.Message = messageAny
	}

	if e != nil {
		s.Error = &statusv1.Status_Error{
			Message: e.Error(),
		}
	}

	return s
}
