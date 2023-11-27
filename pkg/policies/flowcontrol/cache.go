package flowcontrol

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
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
	dmapCache olric.DMap
}

// Cache implements iface.Cache.
var _ iface.Cache = (*Cache)(nil)

// NewCache creates a new cache.
func NewCache(dc *distcache.DistCache, lc fx.Lifecycle) (iface.Cache, error) {
	cache := &Cache{}
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			dmapCache, err := dc.NewDMap("control_point_cache", olricconfig.DMap{})
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
	if controlPoint == "" {
		return nil, ErrCacheControlPointEmpty
	}
	cacheKey := formatCacheKey(controlPoint, cacheType, key)
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
	if controlPoint == "" {
		return ErrCacheControlPointEmpty
	}
	cacheKey := formatCacheKey(controlPoint, cacheType, key)
	return c.dmapCache.Put(ctx, cacheKey, value, olric.EX(ttl))
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
	if controlPoint == "" {
		return ErrCacheControlPointEmpty
	}
	cacheKey := formatCacheKey(controlPoint, cacheType, key)
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
	response := &flowcontrolv1.CacheLookupResponse{
		StateCacheResponses: make(map[string]*flowcontrolv1.KeyLookupResponse),
	}
	if request == nil {
		return response
	}

	type Lookup struct {
		lookupResponse *flowcontrolv1.KeyLookupResponse
		key            string
		cacheType      iface.CacheType
	}

	// define a wait group to wait for all the lookups to complete
	var wg sync.WaitGroup

	execLookup := func(lookup *Lookup) func() {
		return func() {
			defer wg.Done()
			lookupResponse := lookup.lookupResponse
			cachedBytes, err := c.get(ctx, request.ControlPoint, lookup.cacheType, lookup.key)
			if err == nil {
				lookupResponse.Value = cachedBytes
				lookupResponse.LookupStatus = flowcontrolv1.CacheLookupStatus_HIT
				lookupResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
				return
			}
			lookupResponse.LookupStatus = flowcontrolv1.CacheLookupStatus_MISS
			log.Info().Err(err).Msg("error")
			if err.Error() == ErrCacheKeyNotFound.Error() {
				log.Info().Msg("key not found")
				lookupResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_SUCCESS
			} else {
				log.Info().Msg("some other error")
				lookupResponse.OperationStatus = flowcontrolv1.CacheOperationStatus_ERROR
				lookupResponse.Error = err.Error()
			}
		}
	}
	var lookups []*Lookup
	// define a lookup struct to hold the cache key and the cached value
	if request.ResultCacheKey != "" {
		response.ResultCacheResponse = &flowcontrolv1.KeyLookupResponse{
			Key: request.ResultCacheKey,
		}
		lookups = append(lookups, &Lookup{
			key:            request.ResultCacheKey,
			lookupResponse: response.ResultCacheResponse,
			cacheType:      iface.Result,
		})
	}
	for _, stateCacheKey := range request.StateCacheKeys {
		lookupResponse := &flowcontrolv1.KeyLookupResponse{
			Key: stateCacheKey,
		}
		response.StateCacheResponses[stateCacheKey] = lookupResponse
		lookups = append(lookups, &Lookup{
			key:            stateCacheKey,
			lookupResponse: lookupResponse,
			cacheType:      iface.State,
		})
	}

	for i, lookup := range lookups {
		wg.Add(1)
		if i == len(lookups)-1 {
			execLookup(lookup)()
			continue
		}
		panichandler.Go(execLookup(lookup))
	}
	wg.Wait()
	return response
}

// Upsert upserts the cache for the given CacheUpsertRequest.
func (c *Cache) Upsert(ctx context.Context, req *flowcontrolv1.CacheUpsertRequest) *flowcontrolv1.CacheUpsertResponse {
	response := &flowcontrolv1.CacheUpsertResponse{
		StateCacheResponses: make(map[string]*flowcontrolv1.KeyUpsertResponse),
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

	return response
}

// Delete deletes the cache for the given CacheDeleteRequest.
func (c *Cache) Delete(ctx context.Context, req *flowcontrolv1.CacheDeleteRequest) *flowcontrolv1.CacheDeleteResponse {
	response := &flowcontrolv1.CacheDeleteResponse{
		StateCacheResponses: make(map[string]*flowcontrolv1.KeyDeleteResponse),
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

	return response
}

// formatCacheKey returns the cache key for the given control point and key.
func formatCacheKey(controlPoint string, cacheType iface.CacheType, key string) string {
	return "@controlpoint:" + controlPoint + "/type:" + cacheType.String() + "/key:" + key
}
