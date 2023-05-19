package ratetracker

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"

	"github.com/fluxninja/aperture/v2/pkg/distcache"
)

// DistCacheRateTracker implements Limiter.
type DistCacheRateTracker struct {
	mu         sync.RWMutex
	limitCheck RateLimitChecker
	dMap       *olric.DMap
	name       string
}

// NewDistCacheRateTracker creates a new instance of DistCacheRateTracker.
func NewDistCacheRateTracker(limitCheck RateLimitChecker, dc *distcache.DistCache, name string, ttl time.Duration) (RateTracker, error) {
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

	ol := &DistCacheRateTracker{
		name:       name,
		dMap:       dMap,
		limitCheck: limitCheck,
	}

	return ol, nil
}

// Name returns the name of the DistCacheRateTracker.
func (ol *DistCacheRateTracker) Name() string {
	return ol.name
}

// Close cleans up DMap held within the DistCacheRateTracker.
func (ol *DistCacheRateTracker) Close() error {
	ol.mu.Lock()
	defer ol.mu.Unlock()
	err := ol.dMap.Destroy()
	if err != nil {
		return err
	}
	return nil
}

// TakeIfAvailable increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// If an error occurred it returns true, 0 and 0 (fail open).
func (ol *DistCacheRateTracker) TakeIfAvailable(label string, n float64) (bool, float64, float64) {
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

	ok, remaining := ol.limitCheck.CheckRateLimit(label, newN)
	return ok, remaining, newN
}

// GetRateLimitChecker returns the RateLimitCheck of the DistCacheRateTracker.
func (ol *DistCacheRateTracker) GetRateLimitChecker() RateLimitChecker {
	return ol.limitCheck
}

// Make sure DistCacheRateTracker implements Limiter interface.
var _ RateTracker = (*DistCacheRateTracker)(nil)
