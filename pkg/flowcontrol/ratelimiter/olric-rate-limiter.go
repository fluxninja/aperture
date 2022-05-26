package ratelimiter

import (
	"sync"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"

	"aperture.tech/aperture/pkg/distcache"
)

// OlricRateLimiter implements Limiter.
type OlricRateLimiter struct {
	mu         sync.RWMutex
	limitCheck RateLimitCheck
	dMap       *olric.DMap
	name       string
}

// NewOlricRateLimiter creates a new instance of OlricRateLimiter.
func NewOlricRateLimiter(limitCheck RateLimitCheck, dc *distcache.DistCache, name string, ttl time.Duration) (RateLimiter, error) {
	dmapConfig := config.DMap{
		TTLDuration: ttl,
	}

	dc.Mutex.Lock()
	defer dc.Mutex.Unlock()
	dc.AddDMapCustomConfig(name, dmapConfig)
	dMap, err := dc.Olric.NewDMap(name)
	if err != nil {
		return nil, err
	}
	dc.RemoveDMapCustomConfig(name)

	ol := &OlricRateLimiter{
		name:       name,
		dMap:       dMap,
		limitCheck: limitCheck,
	}

	return ol, nil
}

// Name returns the name of the OlricRateLimiter.
func (ol *OlricRateLimiter) Name() string {
	return ol.name
}

// Close cleans up DMap held within the OlricRateLimiter.
func (ol *OlricRateLimiter) Close() error {
	ol.mu.Lock()
	defer ol.mu.Unlock()
	err := ol.dMap.Destroy()
	if err != nil {
		return err
	}
	return nil
}

// Take is a wrapper for TakeN(label, 1).
func (ol *OlricRateLimiter) Take(label string) (bool, int, int) {
	return ol.TakeN(label, 1)
}

// TakeN increments value in label by n and returns whether n events should be allowed along with the remaining value (limit - new n) after increment and the current count for the label.
// If an error occurred it returns true, 0 and 0 (fail open).
func (ol *OlricRateLimiter) TakeN(label string, n int) (bool, int, int) {
	ol.mu.RLock()
	defer ol.mu.RUnlock()
	newN, err := ol.dMap.Incr(label, n)
	if err != nil {
		return true, 0, 0
	}
	ok, remaining := ol.limitCheck.CheckRateLimit(label, newN)
	return ok, remaining, newN
}

// GetRateLimitCheck returns the RateLimitCheck of the OlricRateLimiter.
func (ol *OlricRateLimiter) GetRateLimitCheck() RateLimitCheck {
	return ol.limitCheck
}

// Make sure OlricRateLimiter implements Limiter interface.
var _ RateLimiter = (*OlricRateLimiter)(nil)
