package aperture

import (
	"errors"
	"time"

	flowcontrol "github.com/fluxninja/aperture-go/gen/proto/flowcontrol/v1"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
)

// FlowStatus represents status of feature execution.
type FlowStatus uint8

// User passes a code to indicate status of feature execution.
//
//go:generate enumer -type=FlowStatus -output=flow-status-string.go
const (
	// OK indicates successful feature execution.
	OK FlowStatus = iota
	// Error indicate error on feature execution.
	Error
)

// Flow is the interface that is returned to the user every time a Check call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the feature is executed.
type Flow interface {
	Accepted() bool
	End(status FlowStatus) error
	CheckResponse() *flowcontrol.CheckResponse
}

type flow struct {
	span          trace.Span
	checkResponse *flowcontrol.CheckResponse
	ended         bool
}

// Accepted returns whether the Flow was accepted by Aperture Agent.
func (f *flow) Accepted() bool {
	if f.checkResponse == nil {
		return true
	}
	if f.checkResponse.DecisionType == flowcontrol.CheckResponse_DECISION_TYPE_ACCEPTED {
		return true
	}
	return false
}

// CheckResponse returns the response from the server.
func (f *flow) CheckResponse() *flowcontrol.CheckResponse {
	return f.checkResponse
}

// End is used to end the flow, the user will have to pass a status code and an error description which will define the state and result of the flow.
func (f *flow) End(statusCode FlowStatus) error {
	if f.ended {
		return errors.New("flow already ended")
	}
	f.ended = true

	checkResponseJSONBytes, err := protojson.Marshal(f.checkResponse)
	if err != nil {
		return err
	}
	f.span.SetAttributes(
		attribute.String(featureStatusLabel, statusCode.String()),
		attribute.String(checkResponseLabel, string(checkResponseJSONBytes)),
		attribute.Int64(flowEndTimestampLabel, time.Now().UnixNano()),
	)
	f.span.End()
	return nil
}
