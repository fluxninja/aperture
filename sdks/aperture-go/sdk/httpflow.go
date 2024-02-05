package aperture

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	checkv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	checkhttpv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/checkhttp/v1"
)

// HTTPFlow is the interface that is returned to the user every time a CheckHTTP call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the workload is executed.
type HTTPFlow interface {
	ShouldRun() bool
	SetStatus(status FlowStatus)
	Error() error
	Span() trace.Span
	End() EndResponse
	CheckResponse() *checkhttpv1.CheckHTTPResponse
}

type httpflow struct {
	span              trace.Span
	err               error
	checkResponse     *checkhttpv1.CheckHTTPResponse
	flowParams        FlowParams
	statusCode        FlowStatus
	ended             bool
	flowControlClient checkv1.FlowControlServiceClient
}

// newFlow creates a new flow with default field values.
func newHTTPFlow(span trace.Span, flowParams FlowParams, flowControlClient checkv1.FlowControlServiceClient) *httpflow {
	return &httpflow{
		span:              span,
		checkResponse:     nil,
		statusCode:        OK,
		ended:             false,
		flowParams:        flowParams,
		err:               nil,
		flowControlClient: flowControlClient,
	}
}

// ShouldRun returns whether the Flow was allowed to run by Aperture Agent.
// By default, fail-open behavior is enabled. Set rampMode to disable it.
func (f *httpflow) ShouldRun() bool {
	if (!f.flowParams.RampMode && f.checkResponse == nil) || (f.checkResponse.Status.Code == int32(code.Code_OK)) {
		return true
	} else {
		return false
	}
}

// CheckResponse returns the response from the server.
func (f *httpflow) CheckResponse() *checkhttpv1.CheckHTTPResponse {
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
func (f *httpflow) End() EndResponse {
	if f.ended {
		return EndResponse{
			Error: errors.New("flow already ended"),
		}
	}

	if f.checkResponse == nil {
		return EndResponse{
			Error: errors.New("check response is nil"),
		}
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
					return EndResponse{
						Error: err,
					}
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

	inflightRequests := make([]*checkv1.InflightRequestRef, len(f.checkResponse.GetCheckResponse().GetLimiterDecisions()))

	for _, decision := range f.checkResponse.GetCheckResponse().GetLimiterDecisions() {
		if decision.GetConcurrencyLimiterInfo() != nil {
			if decision.GetConcurrencyLimiterInfo().GetRequestId() == "" {
				continue
			}
			inflightRequest := &checkv1.InflightRequestRef{
				PolicyName:  decision.PolicyName,
				PolicyHash:  decision.PolicyHash,
				ComponentId: decision.ComponentId,
				Label:       decision.GetConcurrencyLimiterInfo().GetLabel(),
				RequestId:   decision.GetConcurrencyLimiterInfo().GetRequestId(),
			}
			if decision.GetConcurrencyLimiterInfo().GetTokensInfo() != nil {
				inflightRequest.Tokens = decision.GetConcurrencyLimiterInfo().GetTokensInfo().GetConsumed()
			}

			inflightRequests = append(inflightRequests, inflightRequest)
		}

		if decision.GetConcurrencySchedulerInfo() != nil {
			if decision.GetConcurrencySchedulerInfo().GetRequestId() == "" {
				continue
			}
			ref := &checkv1.InflightRequestRef{
				PolicyName:  decision.PolicyName,
				PolicyHash:  decision.PolicyHash,
				ComponentId: decision.ComponentId,
				Label:       decision.GetConcurrencySchedulerInfo().GetLabel(),
				RequestId:   decision.GetConcurrencySchedulerInfo().GetRequestId(),
			}
			if decision.GetConcurrencySchedulerInfo().GetTokensInfo() != nil {
				ref.Tokens = decision.GetConcurrencySchedulerInfo().GetTokensInfo().GetConsumed()
			}

			inflightRequests = append(inflightRequests, ref)
		}
	}

	if len(inflightRequests) == 0 {
		return EndResponse{}
	}

	flowEndResponse, err := f.flowControlClient.FlowEnd(context.Background(), &checkv1.FlowEndRequest{
		ControlPoint:     f.checkResponse.GetCheckResponse().ControlPoint,
		InflightRequests: inflightRequests,
	}, f.flowParams.CallOptions...)

	return EndResponse{
		FlowEndResponse: flowEndResponse,
		Error:           err,
	}
}
