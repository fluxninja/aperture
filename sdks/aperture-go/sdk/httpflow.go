package aperture

import (
	"errors"
	"time"

	checkhttpproto "buf.build/gen/go/fluxninja/aperture/protocolbuffers/go/aperture/flowcontrol/checkhttp/v1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// HTTPFlow is the interface that is returned to the user every time a CheckHTTP call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the workload is executed.
type HTTPFlow interface {
	ShouldRun() bool
	SetStatus(status FlowStatus)
	Error() error
	Span() trace.Span
	End() error
	CheckResponse() *checkhttpproto.CheckHTTPResponse
}

type httpflow struct {
	span          trace.Span
	checkResponse *checkhttpproto.CheckHTTPResponse
	statusCode    FlowStatus
	ended         bool
	failOpen      bool
	err           error
}

// newFlow creates a new flow with default field values.
func newHTTPFlow(span trace.Span, failOpen bool) *httpflow {
	return &httpflow{
		span:          span,
		checkResponse: nil,
		statusCode:    OK,
		ended:         false,
		failOpen:      failOpen,
		err:           nil,
	}
}

// ShouldRun returns whether the Flow was allowed to run by Aperture Agent.
// By default, fail-open behavior is enabled. Use DisableFailOpen to disable it.
func (f *httpflow) ShouldRun() bool {
	if (f.failOpen && f.checkResponse == nil) || (f.checkResponse.Status.Code == int32(code.Code_OK)) {
		return true
	} else {
		return false
	}
}

// CheckResponse returns the response from the server.
func (f *httpflow) CheckResponse() *checkhttpproto.CheckHTTPResponse {
	return f.checkResponse
}

// SetStatus sets the status code of a flow.
// If not set explicitly, defaults to FlowStatus.OK.
func (f *httpflow) SetStatus(statusCode FlowStatus) {
	f.statusCode = statusCode
}

// Error returns the error that occurred during the flow.
func (f *httpflow) Error() error {
	return f.err
}

// Span returns the span associated with the flow.
func (f *httpflow) Span() trace.Span {
	return f.span
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
