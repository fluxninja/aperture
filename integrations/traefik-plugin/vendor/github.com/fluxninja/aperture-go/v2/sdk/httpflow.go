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
	Accepted() bool
	End(status FlowStatus) error
	CheckResponse() *flowcontrolhttp.CheckHTTPResponse
}

type httpflow struct {
	span          trace.Span
	checkResponse *flowcontrolhttp.CheckHTTPResponse
	ended         bool
}

// Accepted returns whether the Flow was accepted by Aperture Agent.
func (f *httpflow) Accepted() bool {
	if f.checkResponse == nil {
		return true
	}
	if f.checkResponse.Status.Code == int32(codes.OK) {
		return true
	}
	return false
}

// CheckResponse returns the response from the server.
func (f *httpflow) CheckResponse() *flowcontrolhttp.CheckHTTPResponse {
	return f.checkResponse
}

// End is used to end the flow, the user will have to pass a status code and an error description which will define the state and result of the flow.
func (f *httpflow) End(statusCode FlowStatus) error {
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
		attribute.String(flowStatusLabel, statusCode.String()),
		attribute.String(checkResponseLabel, checkResponseStr),
		attribute.Int64(flowEndTimestampLabel, time.Now().UnixNano()),
	)
	f.span.End()
	return nil
}
