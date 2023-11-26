package iface

import (
	"context"
	"time"
)

// CacheType is the type of cache.
type CacheType int

//go:generate enumer -type=CacheType -transform=lower -output=cache-type-string.go
const (
	// Result is the type of cache for saving results.
	Result CacheType = iota
	// CacheTypeState is the type of cache for saving state.
	State
)

// Cache is an interface for the cache.
type Cache interface {
	// Get returns the value for the given key.
	Get(ctx context.Context, controlPoint string, cacheType CacheType, key string) ([]byte, error)
	// Upsert inserts or updates the value for the given key.
	Upsert(ctx context.Context, controlPoint string, cacheType CacheType, key string, value []byte, ttl time.Duration) error
	// Delete deletes the value for the given key.
	Delete(ctx context.Context, controlPoint string, cacheType CacheType, key string) error
}
