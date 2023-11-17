package flowcontrol

import (
	"context"
	"errors"
	"time"

	"github.com/buraksezer/olric"
	olricconfig "github.com/buraksezer/olric/config"
	distcache "github.com/fluxninja/aperture/v2/pkg/dist-cache"
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

// Get returns the value for the given key.
func (c *Cache) Get(ctx context.Context, controlPoint, key string) ([]byte, error) {
	err := c.Ready()
	if err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrCacheKeyEmpty
	}
	if controlPoint == "" {
		return nil, ErrCacheControlPointEmpty
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
func (c *Cache) Upsert(ctx context.Context, controlPoint, key string, value []byte, ttl time.Duration) error {
	err := c.Ready()
	if err != nil {
		return err
	}
	if key == "" {
		return ErrCacheKeyEmpty
	}
	if controlPoint == "" {
		return ErrCacheControlPointEmpty
	}
	cacheKey := formatCacheKey(controlPoint, key)
	return c.dmapCache.Put(ctx, cacheKey, value, olric.EX(ttl))
}

// Delete deletes the value for the given key.
func (c *Cache) Delete(ctx context.Context, controlPoint, key string) error {
	err := c.Ready()
	if err != nil {
		return err
	}
	if key == "" {
		return ErrCacheKeyEmpty
	}
	if controlPoint == "" {
		return ErrCacheControlPointEmpty
	}
	cacheKey := formatCacheKey(controlPoint, key)
	_, err = c.dmapCache.Delete(ctx, cacheKey)
	return err
}

// Ready returns nil if the cache is ready to be used.
func (c *Cache) Ready() error {
	if c.dmapCache == nil {
		return ErrCacheNotReady
	}
	return nil
}

// formatCacheKey returns the cache key for the given control point and key.
func formatCacheKey(controlPoint, key string) string {
	return "@controlpoint:" + controlPoint + "/key:" + key
}
