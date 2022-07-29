package status

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	statusv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/status/v1"
)

// NewStatus creates a new instance of Status to be pushed into status registry. Use this function for creating status instead of by hand.
// It can either have a detail message or a detail error but not both. This is enforced by first checking for detail message to not be nil.
func NewStatus(d proto.Message, e error) *statusv1.Status {
	s := &statusv1.Status{
		Timestamp: timestamppb.Now(),
	}

	if d != nil {
		detailsMessage, err := anypb.New(d)
		if err != nil {
			return nil
		}
		s.Details = &statusv1.Status_Message{
			Message: detailsMessage,
		}
		return s
	}

	errorDetails := NewErrorDetails(e)
	s.Details = &statusv1.Status_Error{
		Error: errorDetails,
	}

	return s
}

// NewErrorDetails is a helper function to create a new instance of ErrorDetails.
// This recursively fills the cause field from the provided error.
func NewErrorDetails(e error) *statusv1.ErrorDetails {
	errorDetails := &statusv1.ErrorDetails{}

	if e == nil {
		return errorDetails
	}

	errorDetails.Message = e.Error()

	cause := errors.Cause(e)
	if cause != nil {
		if cause != e {
			errorDetails.Cause = NewErrorDetails(cause)
		}
	}

	return errorDetails
}
