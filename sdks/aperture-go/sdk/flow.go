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

var (
	// ErrResultCacheResponseNil is returned when the result cache response is nil.
	ErrResultCacheResponseNil = errors.New("result cache response is nil")

	// ErrResultCacheKeyNotSet is returned when empty result cache key is provided by the caller during start flow.
	ErrResultCacheKeyNotSet = errors.New("result cache key not set")

	// ErrKeyMissingFromStateCacheResponse is returned when the state cache response does not contain the key.
	ErrKeyMissingFromStateCacheResponse = errors.New("key missing from state cache response")
)

// CacheEntry describes the properties of cache entry.
type CacheEntry struct {
	value []byte
	ttl   time.Duration
}

// Flow is the interface that is returned to the user every time a CheckHTTP call through ApertureClient is made.
// The user can check the status of the check call, response from the server, and end the flow once the workload is executed.
type Flow interface {
	ShouldRun() bool
	SetStatus(status FlowStatus)
	ResultCache() KeyLookupResponse
	SetResultCache(ctx context.Context, cacheEntry CacheEntry) KeyUpsertResponse
	DeleteResultCache(ctx context.Context) KeyDeleteResponse
	StateCache(key string) KeyLookupResponse
	SetStateCache(ctx context.Context, key string, cacheEntry CacheEntry) KeyUpsertResponse
	DeleteStateCache(ctx context.Context, key string) KeyDeleteResponse
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
	resultCacheKey    string
	stateCacheKeys    []string
	statusCode        FlowStatus
	ended             bool
	rampMode          bool
}

// flow implements the Flow interface.
var _ Flow = (*flow)(nil)

// newFlow creates a new flow with default field values.
func newFlow(flowControlClient checkv1.FlowControlServiceClient, span trace.Span, rampMode bool, resultCacheKey string, stateCacheKeys []string) *flow {
	return &flow{
		flowControlClient: flowControlClient,
		span:              span,
		checkResponse:     nil,
		statusCode:        OK,
		ended:             false,
		rampMode:          rampMode,
		resultCacheKey:    resultCacheKey,
		stateCacheKeys:    stateCacheKeys,
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

// ResultCache returns the cached value for the flow.
func (f *flow) ResultCache() KeyLookupResponse {
	if f.err != nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, f.err)
	}
	if f.checkResponse == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("check response is nil"))
	}
	if f.checkResponse.CacheLookupResponse == nil || f.checkResponse.CacheLookupResponse.GetResultCacheResponse() == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("result cache is nil"))
	}
	lookupResponse := f.checkResponse.CacheLookupResponse.GetResultCacheResponse()

	return newKeyLookupResponse(lookupResponse.Value, convertCacheLookupStatus(lookupResponse.LookupStatus), convertCacheError(lookupResponse.Error))
}

// SetResultCache sets the result cache entry for the flow.
func (f *flow) SetResultCache(ctx context.Context, cacheEntry CacheEntry) KeyUpsertResponse {
	if f.resultCacheKey == "" {
		return newKeyUpsertResponse(ErrResultCacheKeyNotSet)
	}

	ttlProto := durationpb.New(cacheEntry.ttl)

	cacheUpsertResponse, err := f.flowControlClient.CacheUpsert(ctx, &checkv1.CacheUpsertRequest{
		ControlPoint: f.checkResponse.ControlPoint,
		ResultCacheEntry: &checkv1.CacheEntry{
			Key:   f.resultCacheKey,
			Value: cacheEntry.value,
			Ttl:   ttlProto,
		},
	})
	if err != nil {
		return newKeyUpsertResponse(err)
	}

	if cacheUpsertResponse.ResultCacheResponse == nil {
		return newKeyUpsertResponse(ErrResultCacheResponseNil)
	}

	return newKeyUpsertResponse(convertCacheError(cacheUpsertResponse.ResultCacheResponse.GetError()))
}

// DeleteResultCache deletes the result cache entry for the flow.
func (f *flow) DeleteResultCache(ctx context.Context) KeyDeleteResponse {
	if f.resultCacheKey == "" {
		return newKeyDeleteResponse(ErrResultCacheKeyNotSet)
	}

	cacheDeleteResponse, err := f.flowControlClient.CacheDelete(ctx, &checkv1.CacheDeleteRequest{
		ControlPoint:   f.checkResponse.ControlPoint,
		ResultCacheKey: f.resultCacheKey,
	})
	if err != nil {
		return newKeyDeleteResponse(err)
	}

	if cacheDeleteResponse.ResultCacheResponse == nil {
		return newKeyDeleteResponse(ErrResultCacheResponseNil)
	}
	return newKeyDeleteResponse(convertCacheError(cacheDeleteResponse.ResultCacheResponse.Error))
}

// StateCache returns a state cache entry for the flow.
func (f *flow) StateCache(key string) KeyLookupResponse {
	if f.err != nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, f.err)
	}
	if f.checkResponse == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("check response is nil"))
	}
	if f.checkResponse.CacheLookupResponse == nil || f.checkResponse.CacheLookupResponse.GetStateCacheResponses() == nil {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("state cache is nil"))
	}
	lookupResponseMap := f.checkResponse.CacheLookupResponse.GetStateCacheResponses()
	lookupResponse, ok := lookupResponseMap[key]
	if !ok {
		return newKeyLookupResponse(nil, LookupStatusMiss, errors.New("unknown state cache key"))
	}

	return newKeyLookupResponse(lookupResponse.Value, convertCacheLookupStatus(lookupResponse.LookupStatus), convertCacheError(lookupResponse.Error))
}

// SetStateCache sets a state cache entry for the flow.
func (f *flow) SetStateCache(ctx context.Context, key string, cacheEntry CacheEntry) KeyUpsertResponse {
	ttlProto := durationpb.New(cacheEntry.ttl)

	cacheUpsertResponse, err := f.flowControlClient.CacheUpsert(ctx, &checkv1.CacheUpsertRequest{
		ControlPoint: f.checkResponse.ControlPoint,
		StateCacheEntries: map[string]*checkv1.CacheEntry{
			key: {
				Value: cacheEntry.value,
				Ttl:   ttlProto,
			},
		},
	})
	if err != nil {
		return newKeyUpsertResponse(err)
	}

	upsertResponse, ok := cacheUpsertResponse.StateCacheResponses[key]
	if !ok {
		return newKeyUpsertResponse(ErrKeyMissingFromStateCacheResponse)
	}

	return newKeyUpsertResponse(convertCacheError(upsertResponse.Error))
}

// DeleteStateCache deletes a state cache entry for the flow.
func (f *flow) DeleteStateCache(ctx context.Context, key string) KeyDeleteResponse {
	cacheDeleteResponse, err := f.flowControlClient.CacheDelete(ctx, &checkv1.CacheDeleteRequest{
		ControlPoint: f.checkResponse.ControlPoint,
		StateCacheKeys: []string{
			key,
		},
	})
	if err != nil {
		return newKeyDeleteResponse(err)
	}

	deleteResponse, ok := cacheDeleteResponse.StateCacheResponses[key]
	if !ok {
		return newKeyDeleteResponse(ErrKeyMissingFromStateCacheResponse)
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
