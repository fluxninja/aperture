package iface

import "context"

// Cache is an interface for the cache.
type Cache interface {
	// Get returns the value for the given key.
	Get(ctx context.Context, controlPoint, key string) ([]byte, error)
	// Upsert inserts or updates the value for the given key.
	Upsert(ctx context.Context, controlPoint, key string, value []byte) error
	// Delete deletes the value for the given key.
	Delete(ctx context.Context, controlPoint, key string) error
}
