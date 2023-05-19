package basic

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"

	"github.com/fluxninja/aperture/v2/pkg/distcache"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/flowcontrol/rate-limiter"
)

// BasicRateLimiter implements Limiter.
type BasicRateLimiter struct {
	mu    sync.RWMutex
	dMap  *olric.DMap
	name  string
	limit float64
}

// NewBasicRateLimiter creates a new instance of DistCacheRateTracker.
func NewBasicRateLimiter(dc *distcache.DistCache, name string, ttl time.Duration) (*BasicRateLimiter, error) {
	dc.Mutex.Lock()
	defer dc.Mutex.Unlock()
	dmapConfig := config.DMap{
		TTLDuration: ttl,
		Functions: map[string]config.Function{
			addFunction: add,
		},
	}
	dc.AddDMapCustomConfig(name, dmapConfig)
	dMap, err := dc.Olric.NewDMap(name)
	if err != nil {
		return nil, err
	}
	dc.RemoveDMapCustomConfig(name)

	ol := &BasicRateLimiter{
		name:  name,
		dMap:  dMap,
		limit: -1,
	}

	return ol, nil
}

// Name returns the name of the DistCacheRateTracker.
func (ol *BasicRateLimiter) Name() string {
	return ol.name
}

// Close cleans up DMap held within the DistCacheRateTracker.
func (ol *BasicRateLimiter) Close() error {
	ol.mu.Lock()
	defer ol.mu.Unlock()
	err := ol.dMap.Destroy()
	if err != nil {
		return err
	}
	return nil
}

// SetRateLimit sets the limit.
func (ol *BasicRateLimiter) SetRateLimit(limit float64) {
	ol.mu.Lock()
	defer ol.mu.Unlock()
	ol.limit = limit
}

// GetRateLimit returns the limit.
func (ol *BasicRateLimiter) GetRateLimit() float64 {
	ol.mu.RLock()
	defer ol.mu.RUnlock()
	return ol.limit
}

// TakeIfAvailable increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// If an error occurred it returns true, 0 and 0 (fail open).
func (ol *BasicRateLimiter) TakeIfAvailable(label string, n float64) (bool, float64, float64) {
	ol.mu.RLock()
	defer ol.mu.RUnlock()

	// marshal n as gob
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(n)
	if err != nil {
		return true, 0, 0
	}
	deltaBytes := buf.Bytes()

	resultBytes, err := ol.dMap.Function(label, addFunction, deltaBytes)
	if err != nil {
		return true, 0, 0
	}

	// unmarshal result as float64
	var newN float64
	buf = bytes.NewBuffer(resultBytes)
	err = gob.NewDecoder(buf).Decode(&newN)
	if err != nil {
		return true, 0, 0
	}

	ok, remaining := ratelimiter.CheckRateLimit(newN, ol.limit)
	return ok, remaining, newN
}

// Make sure DistCacheRateTracker implements Limiter interface.
var _ ratelimiter.RateLimiter = (*BasicRateLimiter)(nil)
