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

// GetAll returns a copy of current state of cache.
//
// State is returned as a list in an arbitrary order.
func (c *Cache[K]) GetAll() []K {
	c.Lock()
	defer c.Unlock()
	items := make([]K, 0, len(c.cache))
	for elem := range c.cache {
		items = append(items, elem)
	}
	return items
}
