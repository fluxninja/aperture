package aperture

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"

	checkv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
)

var (
	// ErrResultCacheResponseNil is returned when the result cache response is nil.
	ErrResultCacheResponseNil = errors.New("result cache response is nil")

	// ErrResultCacheKeyNotSet is returned when empty result cache key is provided by the caller during start flow.
	ErrResultCacheKeyNotSet = errors.New("result cache key not set")

	// ErrKeyMissingFromGlobalCacheResponse is returned when the global cache response does not contain the key.
	ErrKeyMissingFromGlobalCacheResponse = errors.New("key missing from global cache response")
)

// CacheEntry describes the properties of cache entry.
type CacheEntry struct {
	Value []byte
	TTL   time.Duration
}

// EndResponse is the response returned by the End method of the Flow interface.
type EndResponse struct {
	// FlowEndResponse is populated if the flow end request succeeded.
	FlowEndResponse *checkv1.FlowEndResponse

	// Error is populated if the flow end request failed.
	Error error
}

// Flow is the interface that is returned to the user every time a CheckHTTP call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the workload is executed.
type Flow interface {
	ShouldRun() bool
	SetStatus(status FlowStatus)
	ResultCache() KeyLookupResponse
	SetResultCache(ctx context.Context, cacheEntry CacheEntry, opts ...grpc.CallOption) KeyUpsertResponse
	DeleteResultCache(ctx context.Context, opts ...grpc.CallOption) KeyDeleteResponse
	GlobalCache(key string) KeyLookupResponse
	SetGlobalCache(ctx context.Context, key string, cacheEntry CacheEntry, opts ...grpc.CallOption) KeyUpsertResponse
	DeleteGlobalCache(ctx context.Context, key string, opts ...grpc.CallOption) KeyDeleteResponse
	Error() error
	Span() trace.Span
	End() EndResponse
	CheckResponse() *checkv1.CheckResponse
	RetryAfter() time.Duration
	HTTPResponseCode() int
}

type flow struct {
	flowControlClient checkv1.FlowControlServiceClient
	span              trace.Span
	err               error
	checkResponse     *checkv1.CheckResponse
	resultCacheKey    string
	globalCacheKeys   []string
	statusCode        FlowStatus
	ended             bool
	rampMode          bool
	callOptions       []grpc.CallOption
}

// flow implements the Flow interface.
var _ Flow = (*flow)(nil)

// newFlow creates a new flow with default field values.
func newFlow(
	flowControlClient checkv1.FlowControlServiceClient,
	span trace.Span,
	rampMode bool,
	resultCacheKey string,
	globalCacheKeys []string,
	callOptions []grpc.CallOption,
) *flow {
	return &flow{
		flowControlClient: flowControlClient,
		span:              span,
		checkResponse:     nil,
		statusCode:        OK,
		ended:             false,
		rampMode:          rampMode,
		resultCacheKey:    resultCacheKey,
		globalCacheKeys:   globalCacheKeys,
		callOptions:       callOptions,
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

// RetryAfter returns the retry-after duration.
func (f *flow) RetryAfter() time.Duration {
	if f.checkResponse == nil {
		return 0
	}
	return f.checkResponse.WaitTime.AsDuration()
}

// HTTPResponseCode returns the HTTP response code.
func (f *flow) HTTPResponseCode() int {
	// Mapping empty status code to 200 to a success
	if f.checkResponse.DeniedResponseStatusCode == checkv1.StatusCode_Empty {
		return 200
	}
	return int(f.checkResponse.DeniedResponseStatusCode)
}

// SetStatus sets the status code of a flow.
// If not set explicitly, defaults to FlowStatus.OK.
func (f *flow) SetStatus(statusCode FlowStatus) {
	f.statusCode = statusCode
}

// ResultCache returns the cached value for the flow.
func (f *flow) ResultCache() KeyLookupResponse {
	if f.err != nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, f.err)
	}
	if f.checkResponse == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("check response is nil"))
	}
	if !f.ShouldRun() {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("flow was rejected"))
	}
	if f.checkResponse.CacheLookupResponse == nil || f.checkResponse.CacheLookupResponse.GetResultCacheResponse() == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("result cache is nil"))
	}
	lookupResponse := f.checkResponse.CacheLookupResponse.GetResultCacheResponse()

	return newKeyLookupResponse(lookupResponse.Value, convertCacheLookupStatus(lookupResponse.LookupStatus), convertCacheError(lookupResponse.Error))
}

// SetResultCache sets the result cache entry for the flow.
func (f *flow) SetResultCache(ctx context.Context, cacheEntry CacheEntry, opts ...grpc.CallOption) KeyUpsertResponse {
	if f.resultCacheKey == "" {
		return newKeyUpsertResponse(ErrResultCacheKeyNotSet)
	}

	if f.checkResponse == nil {
		return newKeyDeleteResponse(f.err)
	}

	ttlProto := durationpb.New(cacheEntry.TTL)

	cacheUpsertResponse, err := f.flowControlClient.CacheUpsert(ctx, &checkv1.CacheUpsertRequest{
		ControlPoint: f.checkResponse.ControlPoint,
		ResultCacheEntry: &checkv1.CacheEntry{
			Key:   f.resultCacheKey,
			Value: cacheEntry.Value,
			Ttl:   ttlProto,
		},
	}, opts...)
	if err != nil {
		return newKeyUpsertResponse(err)
	}

	if cacheUpsertResponse.ResultCacheResponse == nil {
		return newKeyUpsertResponse(ErrResultCacheResponseNil)
	}

	return newKeyUpsertResponse(convertCacheError(cacheUpsertResponse.ResultCacheResponse.GetError()))
}

// DeleteResultCache deletes the result cache entry for the flow.
func (f *flow) DeleteResultCache(ctx context.Context, opts ...grpc.CallOption) KeyDeleteResponse {
	if f.resultCacheKey == "" {
		return newKeyDeleteResponse(ErrResultCacheKeyNotSet)
	}

	if f.checkResponse == nil {
		return newKeyDeleteResponse(f.err)
	}

	cacheDeleteResponse, err := f.flowControlClient.CacheDelete(ctx, &checkv1.CacheDeleteRequest{
		ControlPoint:   f.checkResponse.ControlPoint,
		ResultCacheKey: f.resultCacheKey,
	}, opts...)
	if err != nil {
		return newKeyDeleteResponse(err)
	}

	if cacheDeleteResponse.ResultCacheResponse == nil {
		return newKeyDeleteResponse(ErrResultCacheResponseNil)
	}
	return newKeyDeleteResponse(convertCacheError(cacheDeleteResponse.ResultCacheResponse.Error))
}

// GlobalCache returns a global cache entry for the flow.
func (f *flow) GlobalCache(key string) KeyLookupResponse {
	if f.err != nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, f.err)
	}
	if f.checkResponse == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("check response is nil"))
	}
	if !f.ShouldRun() {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("flow was rejected"))
	}
	if f.checkResponse.CacheLookupResponse == nil || f.checkResponse.CacheLookupResponse.GetGlobalCacheResponses() == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("global cache is nil"))
	}
	lookupResponseMap := f.checkResponse.CacheLookupResponse.GetGlobalCacheResponses()
	lookupResponse, ok := lookupResponseMap[key]
	if !ok {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("unknown global cache key"))
	}

	return newKeyLookupResponse(lookupResponse.Value, convertCacheLookupStatus(lookupResponse.LookupStatus), convertCacheError(lookupResponse.Error))
}

// SetGlobalCache sets a global cache entry for the flow.
func (f *flow) SetGlobalCache(ctx context.Context, key string, cacheEntry CacheEntry, opts ...grpc.CallOption) KeyUpsertResponse {
	ttlProto := durationpb.New(cacheEntry.TTL)

	cacheUpsertResponse, err := f.flowControlClient.CacheUpsert(ctx, &checkv1.CacheUpsertRequest{
		GlobalCacheEntries: map[string]*checkv1.CacheEntry{
			key: {
				Value: cacheEntry.Value,
				Ttl:   ttlProto,
			},
		},
	}, opts...)
	if err != nil {
		return newKeyUpsertResponse(err)
	}

	upsertResponse, ok := cacheUpsertResponse.GlobalCacheResponses[key]
	if !ok {
		return newKeyUpsertResponse(ErrKeyMissingFromGlobalCacheResponse)
	}

	return newKeyUpsertResponse(convertCacheError(upsertResponse.Error))
}

// DeleteGlobalCache deletes a global cache entry for the flow.
func (f *flow) DeleteGlobalCache(ctx context.Context, key string, opts ...grpc.CallOption) KeyDeleteResponse {
	cacheDeleteResponse, err := f.flowControlClient.CacheDelete(ctx, &checkv1.CacheDeleteRequest{
		GlobalCacheKeys: []string{
			key,
		},
	}, opts...)
	if err != nil {
		return newKeyDeleteResponse(err)
	}

	deleteResponse, ok := cacheDeleteResponse.GlobalCacheResponses[key]
	if !ok {
		return newKeyDeleteResponse(ErrKeyMissingFromGlobalCacheResponse)
	}

	return newKeyDeleteResponse(convertCacheError(deleteResponse.Error))
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
func (f *flow) End() EndResponse {
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

	checkResponseJSONBytes, err := protojson.Marshal(f.checkResponse)
	if err != nil {
		return EndResponse{
			Error: err,
		}
	}
	f.span.SetAttributes(
		attribute.String(flowStatusLabel, f.statusCode.String()),
		attribute.String(checkResponseLabel, string(checkResponseJSONBytes)),
		attribute.Int64(flowEndTimestampLabel, time.Now().UnixNano()),
	)
	f.span.End()

	inflightRequests := []*checkv1.InflightRequestRef{}

	for _, decision := range f.checkResponse.GetLimiterDecisions() {
		if decision.GetConcurrencyLimiterInfo() != nil {
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
		ControlPoint:     f.checkResponse.ControlPoint,
		InflightRequests: inflightRequests,
	}, f.callOptions...)

	return EndResponse{
		FlowEndResponse: flowEndResponse,
		Error:           err,
	}
}
