package cache

import (
	"sync"
	"time"

	"go.uber.org/fx"
)

// Provide cache.
func Provide[K comparable](lc fx.Lifecycle) *Cache[K] {
	cache := NewCache[K]()
	lc.Append(fx.StartStopHook(cache.Start, cache.Stop))
	return cache
}

// NewCache returns new instance of cache.
func NewCache[K comparable]() *Cache[K] {
	return &Cache[K]{
		current: make(map[K]struct{}),
	}
}

// Cache is a set of unique keys of type K.
//
// Cache implements epoch-based expiry. Each PERIOD current keys are marked as
// stale and stale keys are discarded.  This means that effective per-key
// expiry time is in between PERIOD and 2*PERIOD.
type Cache[K comparable] struct {
	prev    map[K]struct{} // keys from previous epoch
	current map[K]struct{} // keys from current epoch
	mutex   sync.Mutex     // protects prev & current maps
	stopCh  chan struct{}
	ticker  *time.Ticker
}

const (
	gcPeriod = 30 * time.Second
)

// Start starts the periodic epoch rotation in background.
func (c *Cache[K]) Start() {
	c.ticker = time.NewTicker(gcPeriod)
	c.stopCh = make(chan struct{})
	go func() {
		for {
			select {
			case <-c.stopCh:
				return
			case <-c.ticker.C:
				c.gc()
			}
		}
	}()
}

// Stop stops all the background tasks of the cache.
func (c *Cache[K]) Stop() {
	c.ticker.Stop()
	close(c.stopCh)
}

// Put inserts k into the cache.
func (c *Cache[K]) Put(k K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.current[k] = struct{}{}
}

// Contains checks if map contains given key.
func (c *Cache[K]) Contains(k K) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, inPrev := c.prev[k]; inPrev {
		return true
	}

	_, inCurrent := c.current[k]
	return inCurrent
}

func (c *Cache[K]) gc() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.prev = c.current
	c.current = make(map[K]struct{}, 2*len(c.current))
}

// NewEpochForTest triggers manual epoch rotation.
//
// Should be used only in tests.
func (c *Cache[K]) NewEpochForTest() { c.gc() }

// GetAll returns a copy of current state of cache.
//
// State is returned as a list in an arbitrary order.
func (c *Cache[K]) GetAll() []K {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := make([]K, 0, len(c.prev)+len(c.current))
	for k := range c.current {
		keys = append(keys, k)
	}
	for k := range c.prev {
		if _, inCurrent := c.current[k]; inCurrent {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}
