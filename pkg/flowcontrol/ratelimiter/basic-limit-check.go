package ratelimiter

import (
	"sync"
)

// Make sure BasicLimitCheck implements LimitCheck.
var _ RateLimitCheck = &BasicRateLimitCheck{}

// BasicRateLimitCheck implements LimitCheck.
type BasicRateLimitCheck struct {
	lock      sync.RWMutex
	overrides map[string]float64
	limit     int
}

// NewBasicRateLimitCheck creates a new instance of BasicLimitCheck.
func NewBasicRateLimitCheck() *BasicRateLimitCheck {
	return &BasicRateLimitCheck{
		limit:     -1,
		overrides: make(map[string]float64),
	}
}

// AddOverride sets the limit for a specific label.
func (l *BasicRateLimitCheck) AddOverride(label string, scaleFactor float64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.overrides[label] = scaleFactor
}

// RemoveOverride removes the limit for a specific label.
func (l *BasicRateLimitCheck) RemoveOverride(label string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.overrides, label)
}

// CheckRateLimit checks the limit for a specific label and the remaining limit. If limit is exceeded then we return false and 0 as remaining limit.
func (l *BasicRateLimitCheck) CheckRateLimit(label string, count int) (bool, int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	limit := l.GetLabelRateLimit(label)
	// limit < 0 means that there is no limit
	if limit < 0 {
		return true, -1
	}
	if count > limit {
		return false, 0
	}
	return true, limit - count
}

// SetRateLimit sets the limit.
func (l *BasicRateLimitCheck) SetRateLimit(limit int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.limit = limit
}

// GetRateLimit returns the limit.
func (l *BasicRateLimitCheck) GetRateLimit() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.limit
}

// GetLabelRateLimit returns the limit for a specific label.
func (l *BasicRateLimitCheck) GetLabelRateLimit(label string) int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if scaleFactor, ok := l.overrides[label]; ok {
		return int(float64(l.limit) * scaleFactor)
	}
	return l.limit
}
