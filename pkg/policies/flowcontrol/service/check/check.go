package check

import (
	"context"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	"github.com/fluxninja/aperture/v2/pkg/labels"
	otelconsts "github.com/fluxninja/aperture/v2/pkg/otelcollector/consts"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	servicegetter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/service-getter"
)

const (
	// CacheNotEnabled is the error message when cache is not enabled.
	CacheNotEnabled = "cache is not enabled"
)

//go:generate mockgen -source=check.go -destination=../../../mocks/mock_check.go -package=mocks

// Handler implements the flowcontrol.v1 Service
//
// It also accepts a pointer to an EntityCache for services lookup.
type Handler struct {
	flowcontrolv1.UnimplementedFlowControlServiceServer
	serviceGetter servicegetter.ServiceGetter
	metrics       Metrics
	engine        iface.Engine
	cache         iface.Cache
}

// NewHandler creates a flowcontrol Handler.
func NewHandler(
	serviceGetter servicegetter.ServiceGetter,
	metrics Metrics,
	engine iface.Engine,
	cache iface.Cache,
) *Handler {
	return &Handler{
		serviceGetter: serviceGetter,
		metrics:       metrics,
		engine:        engine,
		cache:         cache,
	}
}

// HandlerWithValues implements the flowcontrol.v1 service using collected inferred values.
type HandlerWithValues interface {
	CheckRequest(
		context.Context,
		iface.RequestContext,
	) *flowcontrolv1.CheckResponse
}

// CheckRequest makes decision using collected inferred fields from authz or Handler.
func (h *Handler) CheckRequest(ctx context.Context, requestContext iface.RequestContext) *flowcontrolv1.CheckResponse {
	checkResponse := h.engine.ProcessRequest(ctx, requestContext)
	h.metrics.CheckResponse(checkResponse.DecisionType, checkResponse.GetRejectReason(), h.engine.GetAgentInfo())
	return checkResponse
}

// Check is the Check method of Flow Control service returns the allow/deny decisions of
// whether to accept the traffic after running the algorithms.
func (h *Handler) Check(ctx context.Context, req *flowcontrolv1.CheckRequest) (*flowcontrolv1.CheckResponse, error) {
	// record the start time of the request
	start := time.Now()

	// make sure labels are not nil, as we append control type label later
	if req.Labels == nil {
		req.Labels = make(map[string]string)
	}

	services := h.serviceGetter.ServicesFromContext(ctx)

	// CheckRequest already pushes result to metrics
	resp := h.CheckRequest(
		ctx,
		iface.RequestContext{
			FlowLabels:     labels.PlainMap(req.Labels),
			ControlPoint:   req.ControlPoint,
			Services:       services,
			RampMode:       req.RampMode,
			ResultCacheKey: req.ResultCacheKey,
			StateCacheKeys: req.StateCacheKeys,
		},
	)
	end := time.Now()
	resp.Start = timestamppb.New(start)
	resp.End = timestamppb.New(end)
	resp.TelemetryFlowLabels = req.Labels
	// add control point type
	resp.TelemetryFlowLabels[otelconsts.ApertureControlPointTypeLabel] = otelconsts.FeatureControlPoint
	return resp, nil
}

// CacheUpsert is the CacheUpsert method of Flow Control service updates the cache with the given key and value.
func (h *Handler) CacheUpsert(ctx context.Context, req *flowcontrolv1.CacheUpsertRequest) (*flowcontrolv1.CacheUpsertResponse, error) {
	type UpsertRequest struct {
		key            string
		entry          *flowcontrolv1.CacheEntry
		upsertResponse *flowcontrolv1.KeyUpsertResponse
		cacheType      iface.CacheType
	}

	wg := sync.WaitGroup{}

	execCacheUpsert := func(upsertRequest *UpsertRequest) func() {
		return func() {
			defer wg.Done()
			if h.cache == nil {
				upsertRequest.upsertResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				upsertRequest.upsertResponse.Error = CacheNotEnabled
				return
			}
			cacheType := upsertRequest.cacheType
			key := upsertRequest.key
			value := upsertRequest.entry.Value
			ttl := upsertRequest.entry.Ttl.AsDuration()
			err := h.cache.Upsert(ctx, req.ControlPoint, cacheType, key, value, ttl)
			if err != nil {
				upsertRequest.upsertResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				upsertRequest.upsertResponse.Error = err.Error()
			} else {
				upsertRequest.upsertResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
			}
		}
	}

	response := &flowcontrolv1.CacheUpsertResponse{}

	var upsertRequests []*UpsertRequest
	if req.ResultCacheEntry != nil && req.ResultCacheEntry.Key != "" {
		wg.Add(1)
		upsertRequests = append(upsertRequests, &UpsertRequest{
			key:            req.ResultCacheEntry.Key,
			cacheType:      iface.Result,
			entry:          req.ResultCacheEntry,
			upsertResponse: response.ResultCacheResponse,
		})
	}

	// iterate over the state cache entries map
	for key, stateCacheEntry := range req.StateCacheEntries {
		if key == "" {
			continue
		}
		wg.Add(1)
		upsertResponse := &flowcontrolv1.KeyUpsertResponse{}
		response.StateCacheResponses[key] = upsertResponse
		// set the state cache entry key
		upsertRequests = append(upsertRequests, &UpsertRequest{
			key:            key,
			cacheType:      iface.State,
			entry:          stateCacheEntry,
			upsertResponse: upsertResponse,
		})
	}

	for i, upsertRequest := range upsertRequests {
		if i == len(upsertRequests)-1 {
			execCacheUpsert(upsertRequest)()
			continue
		}
		go execCacheUpsert(upsertRequest)()
	}
	wg.Wait()

	return response, nil
}

// CacheDelete is the CacheDelete method of Flow Control service deletes the cache entry with the given key.
func (h *Handler) CacheDelete(ctx context.Context, req *flowcontrolv1.CacheDeleteRequest) (*flowcontrolv1.CacheDeleteResponse, error) {
	type DeleteRequest struct {
		deleteResponse *flowcontrolv1.KeyDeleteResponse
		key            string
		cacheType      iface.CacheType
	}

	wg := sync.WaitGroup{}

	execCacheDelete := func(deleteRequest *DeleteRequest) func() {
		return func() {
			defer wg.Done()
			if h.cache == nil {
				deleteRequest.deleteResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				deleteRequest.deleteResponse.Error = CacheNotEnabled
				return
			}
			cacheType := deleteRequest.cacheType
			key := deleteRequest.key
			err := h.cache.Delete(ctx, req.ControlPoint, cacheType, key)
			if err != nil {
				deleteRequest.deleteResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				deleteRequest.deleteResponse.Error = err.Error()
			} else {
				deleteRequest.deleteResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
			}
		}
	}

	response := &flowcontrolv1.CacheDeleteResponse{}

	var deleteRequests []*DeleteRequest
	if req.ResultCacheKey != "" {
		wg.Add(1)
		deleteRequests = append(deleteRequests, &DeleteRequest{
			cacheType:      iface.Result,
			key:            req.ResultCacheKey,
			deleteResponse: response.ResultCacheResponse,
		})
	}

	for _, stateCacheKey := range req.StateCacheKeys {
		if stateCacheKey == "" {
			continue
		}
		wg.Add(1)
		deleteResponse := &flowcontrolv1.KeyDeleteResponse{}
		response.StateCacheResponses[stateCacheKey] = deleteResponse
		deleteRequests = append(deleteRequests, &DeleteRequest{
			cacheType:      iface.State,
			key:            stateCacheKey,
			deleteResponse: deleteResponse,
		})
	}

	for i, deleteRequest := range deleteRequests {
		if i == len(deleteRequests)-1 {
			execCacheDelete(deleteRequest)()
			continue
		}
		go execCacheDelete(deleteRequest)()
	}
	wg.Wait()

	return response, nil
}
