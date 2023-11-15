package flowcontrol

import (
	"context"
	"errors"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
	"github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/iface"
)

const (
	// ErrCacheKeyEmpty is the error returned when the cache key is empty.
	ErrCacheKeyEmpty = "cache key cannot be empty"
	// ErrCacheControlPointEmpty is the error returned when the cache control point is empty.
	ErrCacheControlPointEmpty = "cache control_point cannot be empty"
)

// Cache for saving responses at flow end.
type Cache struct {
	dmapCache olric.DMap
}

// Cache implements iface.Cache.
var _ iface.Cache = (*Cache)(nil)

// NewCache creates a new cache.
func NewCache(dc *distcache.DistCache) (iface.Cache, error) {
	dmapCache, err := dc.NewDMap("control_point_cache", olricconfig.DMap{})
	if err != nil {
		return nil, err
	}
	return &Cache{dmapCache: dmapCache}, nil
}

// Get returns the value for the given key.
func (c *Cache) Get(ctx context.Context, controlPoint, key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New(ErrCacheKeyEmpty)
	}
	if controlPoint == "" {
		return nil, errors.New(ErrCacheControlPointEmpty)
	}
	cacheKey := formatCacheKey(controlPoint, key)
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

// Upsert inserts or updates the value for the given key.
func (c *Cache) Upsert(ctx context.Context, controlPoint, key string, value []byte) error {
	if key == "" {
		return errors.New(ErrCacheKeyEmpty)
	}
	if controlPoint == "" {
		return errors.New(ErrCacheControlPointEmpty)
	}
	cacheKey := formatCacheKey(controlPoint, key)
	return c.dmapCache.Put(ctx, cacheKey, value)
}

// Delete deletes the value for the given key.
func (c *Cache) Delete(ctx context.Context, controlPoint, key string) error {
	if key == "" {
		return errors.New(ErrCacheKeyEmpty)
	}
	if controlPoint == "" {
		return errors.New(ErrCacheControlPointEmpty)
	}
	cacheKey := formatCacheKey(controlPoint, key)
	_, err := c.dmapCache.Delete(ctx, cacheKey)
	return err
}

// formatCacheKey returns the cache key for the given control point and key.
func formatCacheKey(controlPoint, key string) string {
	return "@controlpoint:" + controlPoint + "/key:" + key
}
