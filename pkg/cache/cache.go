package cache

import (
	"sync"
)

// Provide cache.
func Provide[K comparable]() *Cache[K] {
	return NewCache[K]()
}

// NewCache returns new instance of cache.
func NewCache[K comparable]() *Cache[K] {
	return &Cache[K]{
		cache: map[K]struct{}{},
	}
}

// Cache is a set of unique keys of type K.
type Cache[K comparable] struct {
	sync.Mutex
	cache map[K]struct{}
}

// Put inserts k into the cache.
func (c *Cache[K]) Put(k K) {
	c.Lock()
	defer c.Unlock()
	c.cache[k] = struct{}{}
}

// GetAllAndClear returns the current state of cache and clears the cache.
func (c *Cache[K]) GetAllAndClear() map[K]struct{} {
	c.Lock()
	defer c.Unlock()
	result := c.cache
	c.cache = map[K]struct{}{}
	return result
}
