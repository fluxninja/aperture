package aperture

import (
	"errors"
	"time"

	flowcontrolhttp "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/checkhttp/v1"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// HTTPFlow is the interface that is returned to the user every time a CheckHTTP call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the workload is executed.
type HTTPFlow interface {
	Decision() FlowDecision
	ShouldRun() bool
	DisableFailOpen()
	SetStatus(status FlowStatus)
	End() error
	CheckResponse() *flowcontrolhttp.CheckHTTPResponse
}

type httpflow struct {
	span          trace.Span
	checkResponse *flowcontrolhttp.CheckHTTPResponse
	statusCode    FlowStatus
	ended         bool
	failOpen      bool
}

// newFlow creates a new flow with default field values.
func newHTTPFlow(span trace.Span) *httpflow {
	return &httpflow{
		span:          span,
		checkResponse: nil,
		statusCode:    OK,
		ended:         false,
		failOpen:      true,
	}
}

// Decision returns Aperture Agent's decision or information on Agent being unreachable.
func (f *httpflow) Decision() FlowDecision {
	if f.checkResponse == nil {
		return Unreachable
	} else if f.checkResponse.Status.Code == int32(codes.OK) {
		return Accepted
	} else {
		return Rejected
	}
}

// ShouldRun returns whether the Flow was allowed to run by Aperture Agent.
// By default, fail-open behavior is enabled. Use DisableFailOpen to disable it.
func (f *httpflow) ShouldRun() bool {
	decision := f.Decision()
	return decision == Accepted || (f.failOpen && decision == Unreachable)
}

// DisableFailOpen disables fail-open behavior for the flow.
func (f *httpflow) DisableFailOpen() {
	f.failOpen = false
}

// CheckResponse returns the response from the server.
func (f *httpflow) CheckResponse() *flowcontrolhttp.CheckHTTPResponse {
	return f.checkResponse
}

// SetStatus sets the status code of a flow.
// If not set explicitly, defaults to FlowStatus.OK.
func (f *httpflow) SetStatus(statusCode FlowStatus) {
	f.statusCode = statusCode
}

// End is used to end the flow, using the status code previously set using SetStatus method.
func (f *httpflow) End() error {
	if f.ended {
		return errors.New("flow already ended")
	}
	f.ended = true

	checkResponseStr := ""
	if dynamicmeta := f.checkResponse.GetDynamicMetadata(); dynamicmeta != nil {
		value := dynamicmeta.GetFields()[checkResponseLabel]
		if value != nil {
			if sv, ok := value.GetKind().(*structpb.Value_StringValue); ok {
				checkResponseStr = sv.StringValue
			} else {
				checkResponseBytes, err := protojson.Marshal(value)
				if err != nil {
					return err
				}
				checkResponseStr = string(checkResponseBytes)
			}
		}
	}

	f.span.SetAttributes(
		attribute.String(flowStatusLabel, f.statusCode.String()),
		attribute.String(checkResponseLabel, checkResponseStr),
		attribute.Int64(flowEndTimestampLabel, time.Now().UnixNano()),
	)
	f.span.End()
	return nil
}
