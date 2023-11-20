package aperture

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"

	checkv1 "github.com/fluxninja/aperture-go/v2/gen/proto/flowcontrol/check/v1"
)

// Flow is the interface that is returned to the user every time a CheckHTTP call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the workload is executed.
type Flow interface {
	ShouldRun() bool
	SetStatus(status FlowStatus)
	CachedValue() GetCachedValueResponse
	SetCachedValue(ctx context.Context, value []byte, ttl time.Duration) SetCachedValueResponse
	DeleteCachedValue(ctx context.Context) DeleteCachedValueResponse
	Error() error
	Span() trace.Span
	End() error
	CheckResponse() *checkv1.CheckResponse
}

type flow struct {
	flowControlClient checkv1.FlowControlServiceClient
	span              trace.Span
	err               error
	checkResponse     *checkv1.CheckResponse
	statusCode        FlowStatus
	ended             bool
	rampMode          bool
	cacheKey          string
}

// newFlow creates a new flow with default field values.
func newFlow(flowControlClient checkv1.FlowControlServiceClient, span trace.Span, rampMode bool, cacheKey string) *flow {
	return &flow{
		flowControlClient: flowControlClient,
		span:              span,
		checkResponse:     nil,
		statusCode:        OK,
		ended:             false,
		rampMode:          rampMode,
		cacheKey:          cacheKey,
	}
}

// ShouldRun returns whether the Flow was allowed to run by Aperture Agent.
// By default, fail-open behavior is enabled. Set rampMode to disable it.
func (f *flow) ShouldRun() bool {
	return (!f.rampMode && f.checkResponse == nil) || (f.checkResponse.DecisionType == checkv1.CheckResponse_DECISION_TYPE_ACCEPTED)
}

// CheckResponse returns the response from the server.
func (f *flow) CheckResponse() *checkv1.CheckResponse {
	return f.checkResponse
}

// SetStatus sets the status code of a flow.
// If not set explicitly, defaults to FlowStatus.OK.
func (f *flow) SetStatus(statusCode FlowStatus) {
	f.statusCode = statusCode
}

// CachedValue returns the cached value for the flow.
func (f *flow) CachedValue() GetCachedValueResponse {
	if f.err != nil {
		return newGetCachedValueResponse(nil, LookupStatusMiss, OperationStatusError, f.err)
	}
	if f.checkResponse == nil {
		return newGetCachedValueResponse(nil, LookupStatusMiss, OperationStatusError, errors.New("check response is nil"))
	}
	cachedValue := f.checkResponse.GetCachedValue()
	if cachedValue == nil {
		return newGetCachedValueResponse(nil, LookupStatusMiss, OperationStatusError, errors.New("cached value is nil"))
	}

	return newGetCachedValueResponse(cachedValue.Value, convertCacheLookupStatus(cachedValue.LookupStatus), convertCacheOperationStatus(cachedValue.OperationStatus), nil)
}

// SetCachedValue sets the cached value for the flow.
func (f *flow) SetCachedValue(ctx context.Context, value []byte, ttl time.Duration) SetCachedValueResponse {
	if f.cacheKey == "" {
		return newSetCachedValueResponse(OperationStatusError, ErrCacheKeyNotSet)
	}

	ttlProto := durationpb.New(ttl)

	cacheUpsertResponse, err := f.flowControlClient.CacheUpsert(ctx, &checkv1.CacheUpsertRequest{
		ControlPoint: f.checkResponse.ControlPoint,
		Key:          f.cacheKey,
		Value:        value,
		Ttl:          ttlProto,
	})
	if err != nil {
		return newSetCachedValueResponse(OperationStatusError, err)
	}

	return newSetCachedValueResponse(convertCacheOperationStatus(cacheUpsertResponse.GetOperationStatus()), convertCacheError(cacheUpsertResponse.GetError()))
}

// DeleteCachedValue deletes the cached value for the flow.
func (f *flow) DeleteCachedValue(ctx context.Context) DeleteCachedValueResponse {
	if f.cacheKey == "" {
		return newDeleteCachedValueResponse(OperationStatusError, ErrCacheKeyNotSet)
	}

	cacheDeleteResponse, err := f.flowControlClient.CacheDelete(ctx, &checkv1.CacheDeleteRequest{
		ControlPoint: f.checkResponse.ControlPoint,
		Key:          f.cacheKey,
	})
	if err != nil {
		return newDeleteCachedValueResponse(OperationStatusError, err)
	}

	return newDeleteCachedValueResponse(convertCacheOperationStatus(cacheDeleteResponse.GetOperationStatus()), convertCacheError(cacheDeleteResponse.GetError()))
}

// Error returns the error that occurred during the flow.
func (f *flow) Error() error {
	return f.err
}

// Span returns the span associated with the flow.
func (f *flow) Span() trace.Span {
	return f.span
}

// End is used to end the flow, using the status code previously set using SetStatus method.
func (f *flow) End() error {
	if f.ended {
		return errors.New("flow already ended")
	}
	f.ended = true

	checkResponseJSONBytes, err := protojson.Marshal(f.checkResponse)
	if err != nil {
		return err
	}
	f.span.SetAttributes(
		attribute.String(flowStatusLabel, f.statusCode.String()),
		attribute.String(checkResponseLabel, string(checkResponseJSONBytes)),
		attribute.Int64(flowEndTimestampLabel, time.Now().UnixNano()),
	)
	f.span.End()
	return nil
}
