package flowcontrol

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	flowcontrolv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/flowcontrol/check/v1"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

var (
	// ErrCacheKeyEmpty is the error returned when the cache key is empty.
	ErrCacheKeyEmpty = errors.New("cache key cannot be empty")
	// ErrCacheControlPointEmpty is the error returned when the cache control point is empty.
	ErrCacheControlPointEmpty = errors.New("cache control_point cannot be empty")
	// ErrCacheNotReady is the error returned when the cache is not ready to be used.
	ErrCacheNotReady = errors.New("cache is not ready")
	// ErrCacheKeyNotFound is the error returned when the key is not found in the cache. This is copied from the internal olric package.
	ErrCacheKeyNotFound = errors.New("key not found")
)

// Cache for saving responses at flow end.
type Cache struct {
	dmapCache                  olric.DMap
	cacheLookupHitsTotal       *prometheus.CounterVec
	cacheMissesTotal           *prometheus.CounterVec
	cacheOperationResultsTotal *prometheus.CounterVec
}

// Cache implements iface.Cache.
var _ iface.Cache = (*Cache)(nil)

// NewCache creates a new cache.
func NewCache(dc *distcache.DistCache, lc fx.Lifecycle, pr *prometheus.Registry) (iface.Cache, error) {
	labels := []string{
		metrics.CacheTypeLabel,
		metrics.ControlPointLabel,
	}
	cacheLookupHitsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.CacheLookupHitsTotalMetricName,
		Help: "Cumulative number of cache lookup hits.",
	}, labels)
	cacheLookupMissesTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.CacheLookupMissesTotalMetricName,
		Help: "Cumulative number of cache lookup misses.",
	}, labels)
	labels = append(labels, []string{metrics.CacheOperationTypeLabel, metrics.CacheOperationStatusLabel}...)
	cacheOperationResultsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.CacheOperationResultsTotalMetricName,
		Help: "Cumulative number of operation statuses.",
	}, labels)
	cache := &Cache{
		cacheLookupHitsTotal:       cacheLookupHitsTotal,
		cacheMissesTotal:           cacheLookupMissesTotal,
		cacheOperationResultsTotal: cacheOperationResultsTotal,
	}
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			for _, m := range []prometheus.Collector{
				cacheLookupHitsTotal,
				cacheLookupMissesTotal,
				cacheOperationResultsTotal,
			} {
				err := pr.Register(m)
				if err != nil {
					if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
						return fmt.Errorf("unable to register cache metrics: %v", err)
					}
				}
			}
			dmapCache, err := dc.NewDMap("control_point_cache", olricconfig.DMap{}, true)
			if err != nil {
				return err
			}
			cache.dmapCache = dmapCache
			return nil
		},
	})
	return cache, nil
}

// get returns the value for the given key.
func (c *Cache) get(ctx context.Context, controlPoint string, cacheType iface.CacheType, key string) ([]byte, error) {
	err := c.ready()
	if err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrCacheKeyEmpty
	}
	cacheKey, err := formatCacheKey(controlPoint, cacheType, key)
	if err != nil {
		return nil, err
	}
	getResponse, err := c.dmapCache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	cachedBytes, err := getResponse.Byte()
	if err != nil {
		return nil, err
	}

	return cachedBytes, nil
}

// upsert inserts or updates the value for the given key.
func (c *Cache) upsert(ctx context.Context, controlPoint string, cacheType iface.CacheType, key string, value []byte, ttl time.Duration) error {
	err := c.ready()
	if err != nil {
		return err
	}
	if key == "" {
		return ErrCacheKeyEmpty
	}
	cacheKey, err := formatCacheKey(controlPoint, cacheType, key)
	if err != nil {
		return err
	}
	_, err = c.dmapCache.Put(ctx, cacheKey, value, olric.EX(ttl))
	return err
}

// delete deletes the value for the given key.
func (c *Cache) delete(ctx context.Context, controlPoint string, cacheType iface.CacheType, key string) error {
	err := c.ready()
	if err != nil {
		return err
	}
	if key == "" {
		return ErrCacheKeyEmpty
	}
	cacheKey, err := formatCacheKey(controlPoint, cacheType, key)
	if err != nil {
		return err
	}
	_, err = c.dmapCache.Delete(ctx, cacheKey)
	return err
}

// ready returns nil if the cache is ready to be used.
func (c *Cache) ready() error {
	if c.dmapCache == nil {
		return ErrCacheNotReady
	}
	return nil
}

// Lookup looks up the cache for the given CacheLookupRequest.
func (c *Cache) Lookup(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) *flowcontrolv1.CacheLookupResponse {
	response, wgResult, wgGlobal := c.LookupWait(ctx, request)
	wgResult.Wait()
	wgGlobal.Wait()
	c.reportLookupMetrics(request, response)
	return response
}

// LookupWait looks up the cache for the given CacheLookupRequest. It does not wait for the response.
func (c *Cache) LookupWait(ctx context.Context, request *flowcontrolv1.CacheLookupRequest) (*flowcontrolv1.CacheLookupResponse, *sync.WaitGroup, *sync.WaitGroup) {
	// define wait groups
	var wgResult, wgGlobal sync.WaitGroup
	response := &flowcontrolv1.CacheLookupResponse{
		GlobalCacheResponses: make(map[string]*flowcontrolv1.KeyLookupResponse),
	}
	if request == nil {
		return response, &wgResult, &wgGlobal
	}

	type Lookup struct {
		lookupResponse *flowcontrolv1.KeyLookupResponse
		key            string
		cacheType      iface.CacheType
		wg             *sync.WaitGroup
	}

	execLookup := func(lookup *Lookup) func() {
		return func() {
			defer lookup.wg.Done()
			lookupResponse := lookup.lookupResponse
			cachedBytes, err := c.get(ctx, request.ControlPoint, lookup.cacheType, lookup.key)
			if err == nil {
				lookupResponse.Value = cachedBytes
				lookupResponse.LookupStatus = flowcontrolv1.CacheLookupStatus_HIT
				lookupResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
				return
			}
			lookupResponse.LookupStatus = flowcontrolv1.CacheLookupStatus_MISS
			if err.Error() == ErrCacheKeyNotFound.Error() {
				lookupResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
			} else {
				lookupResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				lookupResponse.Error = err.Error()
			}
		}
	}
	var lookups []*Lookup
	// define a lookup struct to hold the cache key and the cached value
	if request.ResultCacheKey != "" {
		response.ResultCacheResponse = &flowcontrolv1.KeyLookupResponse{}
		lookups = append(lookups, &Lookup{
			key:            request.ResultCacheKey,
			lookupResponse: response.ResultCacheResponse,
			cacheType:      iface.Result,
			wg:             &wgResult,
		})
	}
	for _, globalCacheKey := range request.GlobalCacheKeys {
		lookupResponse := &flowcontrolv1.KeyLookupResponse{}
		response.GlobalCacheResponses[globalCacheKey] = lookupResponse
		lookups = append(lookups, &Lookup{
			key:            globalCacheKey,
			lookupResponse: lookupResponse,
			cacheType:      iface.Global,
			wg:             &wgGlobal,
		})
	}

	for i, lookup := range lookups {
		lookup.wg.Add(1)
		if i == len(lookups)-1 {
			execLookup(lookup)()
			continue
		}
		panichandler.Go(execLookup(lookup))
	}
	return response, &wgResult, &wgGlobal
}

// Upsert upserts the cache for the given CacheUpsertRequest.
func (c *Cache) Upsert(ctx context.Context, req *flowcontrolv1.CacheUpsertRequest) *flowcontrolv1.CacheUpsertResponse {
	response := &flowcontrolv1.CacheUpsertResponse{
		GlobalCacheResponses: make(map[string]*flowcontrolv1.KeyUpsertResponse),
	}
	if req == nil {
		return response
	}

	type UpsertRequest struct {
		entry          *flowcontrolv1.CacheEntry
		upsertResponse *flowcontrolv1.KeyUpsertResponse
		key            string
		cacheType      iface.CacheType
	}

	wg := sync.WaitGroup{}

	execCacheUpsert := func(upsertRequest *UpsertRequest) func() {
		return func() {
			defer wg.Done()
			cacheType := upsertRequest.cacheType
			key := upsertRequest.key
			value := upsertRequest.entry.Value
			ttl := upsertRequest.entry.Ttl.AsDuration()
			// Default and Cap TTL to one week
			if ttl == 0 || (ttl > time.Hour*24*7) {
				ttl = time.Hour * 24 * 7
			}
			err := c.upsert(ctx, req.ControlPoint, cacheType, key, value, ttl)
			if err != nil {
				upsertRequest.upsertResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				upsertRequest.upsertResponse.Error = err.Error()
			} else {
				upsertRequest.upsertResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
			}
		}
	}

	var upsertRequests []*UpsertRequest
	if req.ResultCacheEntry != nil && req.ResultCacheEntry.Key != "" {
		wg.Add(1)
		upsertResponse := &flowcontrolv1.KeyUpsertResponse{}
		response.ResultCacheResponse = upsertResponse
		upsertRequests = append(upsertRequests, &UpsertRequest{
			key:            req.ResultCacheEntry.Key,
			cacheType:      iface.Result,
			entry:          req.ResultCacheEntry,
			upsertResponse: upsertResponse,
		})
	}

	// iterate over the state cache entries map
	for key, globalCacheEntry := range req.GlobalCacheEntries {
		if key == "" {
			continue
		}
		wg.Add(1)
		upsertResponse := &flowcontrolv1.KeyUpsertResponse{}
		response.GlobalCacheResponses[key] = upsertResponse
		// set the state cache entry key
		upsertRequests = append(upsertRequests, &UpsertRequest{
			key:            key,
			cacheType:      iface.Global,
			entry:          globalCacheEntry,
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

	return response
}

// Delete deletes the cache for the given CacheDeleteRequest.
func (c *Cache) Delete(ctx context.Context, req *flowcontrolv1.CacheDeleteRequest) *flowcontrolv1.CacheDeleteResponse {
	response := &flowcontrolv1.CacheDeleteResponse{
		GlobalCacheResponses: make(map[string]*flowcontrolv1.KeyDeleteResponse),
	}

	if req == nil {
		return response
	}

	type DeleteRequest struct {
		deleteResponse *flowcontrolv1.KeyDeleteResponse
		key            string
		cacheType      iface.CacheType
	}

	wg := sync.WaitGroup{}

	execCacheDelete := func(deleteRequest *DeleteRequest) func() {
		return func() {
			defer wg.Done()
			cacheType := deleteRequest.cacheType
			key := deleteRequest.key
			err := c.delete(ctx, req.ControlPoint, cacheType, key)
			if err != nil {
				deleteRequest.deleteResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				deleteRequest.deleteResponse.Error = err.Error()
			} else {
				deleteRequest.deleteResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
			}
		}
	}

	var deleteRequests []*DeleteRequest
	if req.ResultCacheKey != "" {
		wg.Add(1)
		deleteResponse := &flowcontrolv1.KeyDeleteResponse{}
		response.ResultCacheResponse = deleteResponse
		deleteRequests = append(deleteRequests, &DeleteRequest{
			cacheType:      iface.Result,
			key:            req.ResultCacheKey,
			deleteResponse: deleteResponse,
		})
	}

	for _, globalCacheKey := range req.GlobalCacheKeys {
		if globalCacheKey == "" {
			continue
		}
		wg.Add(1)
		deleteResponse := &flowcontrolv1.KeyDeleteResponse{}
		response.GlobalCacheResponses[globalCacheKey] = deleteResponse
		deleteRequests = append(deleteRequests, &DeleteRequest{
			cacheType:      iface.Global,
			key:            globalCacheKey,
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

	c.reportDeleteMetrics(req, response)
	return response
}

func (c *Cache) reportLookupMetrics(
	request *flowcontrolv1.CacheLookupRequest,
	response *flowcontrolv1.CacheLookupResponse,
) {
	if response.ResultCacheResponse != nil {
		c.reportLookupMetricsForResponse(
			metrics.CacheTypeResult,
			request.ControlPoint,
			response.ResultCacheResponse,
		)
		c.reportOperationMetricsForResponse(
			metrics.CacheOperationTypeLookup,
			metrics.CacheTypeResult,
			request.ControlPoint,
			response.ResultCacheResponse.OperationStatus,
		)
	}
	for _, resp := range response.GlobalCacheResponses {
		c.reportLookupMetricsForResponse(
			metrics.CacheTypeGlobal,
			request.ControlPoint,
			resp,
		)
		c.reportOperationMetricsForResponse(
			metrics.CacheOperationTypeLookup,
			metrics.CacheTypeGlobal,
			request.ControlPoint,
			resp.OperationStatus,
		)
	}
}

func (c *Cache) reportDeleteMetrics(
	request *flowcontrolv1.CacheDeleteRequest,
	response *flowcontrolv1.CacheDeleteResponse,
) {
	if response.ResultCacheResponse != nil {
		c.reportOperationMetricsForResponse(
			metrics.CacheOperationTypeDelete,
			metrics.CacheTypeResult,
			request.ControlPoint,
			response.ResultCacheResponse.OperationStatus,
		)
	}
	for _, resp := range response.GlobalCacheResponses {
		c.reportOperationMetricsForResponse(
			metrics.CacheOperationTypeDelete,
			metrics.CacheTypeGlobal,
			request.ControlPoint,
			resp.OperationStatus,
		)
	}
}

func (c *Cache) reportLookupMetricsForResponse(cacheType, controlPoint string, response *flowcontrolv1.KeyLookupResponse) {
	labels := prometheus.Labels(map[string]string{
		metrics.CacheTypeLabel:    cacheType,
		metrics.ControlPointLabel: controlPoint,
	})
	switch response.LookupStatus {
	case flowcontrolv1.CacheLookupStatus_HIT:
		hitCounter, err := c.cacheLookupHitsTotal.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract cache hit count counter metric: %v", err)
			return
		}
		hitCounter.Inc()
	case flowcontrolv1.CacheLookupStatus_MISS:
		missCounter, err := c.cacheMissesTotal.GetMetricWith(labels)
		if err != nil {
			log.Debug().Msgf("Could not extract cache miss count counter metric: %v", err)
			return
		}
		missCounter.Inc()
	}
}

func (c *Cache) reportOperationMetricsForResponse(operationType, cacheType, controlPoint string, operationStatus flowcontrolv1.CacheOperationStatus) {
	labels := prometheus.Labels(map[string]string{
		metrics.CacheTypeLabel:          cacheType,
		metrics.ControlPointLabel:       controlPoint,
		metrics.CacheOperationTypeLabel: operationType,
	})
	switch operationStatus {
	case flowcontrolv1.CacheOperationStatus_SUCCESS:
		labels[metrics.CacheOperationStatusLabel] = metrics.CacheOperationStatusSuccess
	case flowcontrolv1.CacheOperationStatus_ERROR:
		labels[metrics.CacheOperationStatusLabel] = metrics.CacheOperationStatusError
	}
	statusCounter, err := c.cacheOperationResultsTotal.GetMetricWith(labels)
	if err != nil {
		log.Debug().Msgf("Could not extract cache operation status metric: %v", err)
		return
	}
	statusCounter.Inc()
}

// formatCacheKey returns the cache key for the given control point and key.
func formatCacheKey(controlPoint string, cacheType iface.CacheType, key string) (string, error) {
	if cacheType == iface.Result {
		if controlPoint == "" {
			return "", ErrCacheControlPointEmpty
		}
		return "@controlpoint:" + controlPoint + "/type:bytes" + "/key:" + key, nil
	} else {
		return "@global" + "/type:bytes" + "/key:" + key, nil
	}
}
