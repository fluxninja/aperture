package ratetracker

import (
	"sync"
)

// Make sure BasicLimitCheck implements LimitCheck.
var _ RateLimitChecker = &BasicRateLimitChecker{}

// BasicRateLimitChecker implements LimitCheck.
type BasicRateLimitChecker struct {
	lock      sync.RWMutex
	overrides map[string]float64
	limit     int
}

// NewBasicRateLimitChecker creates a new instance of BasicLimitCheck.
func NewBasicRateLimitChecker() *BasicRateLimitChecker {
	return &BasicRateLimitChecker{
		limit:     -1,
		overrides: make(map[string]float64),
	}
}

// AddOverride sets the limit for a specific label.
func (l *BasicRateLimitChecker) AddOverride(label string, scaleFactor float64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.overrides[label] = scaleFactor
}

// RemoveOverride removes the limit for a specific label.
func (l *BasicRateLimitChecker) RemoveOverride(label string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.overrides, label)
}

// CheckRateLimit checks the limit for a specific label and the remaining limit. If limit is exceeded then we return false and 0 as remaining limit.
func (l *BasicRateLimitChecker) CheckRateLimit(label string, count int) (bool, int) {
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
func (l *BasicRateLimitChecker) SetRateLimit(limit int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.limit = limit
}

// GetRateLimit returns the limit.
func (l *BasicRateLimitChecker) GetRateLimit() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.limit
}

// GetLabelRateLimit returns the limit for a specific label.
func (l *BasicRateLimitChecker) GetLabelRateLimit(label string) int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if scaleFactor, ok := l.overrides[label]; ok {
		return int(float64(l.limit) * scaleFactor)
	}
	return l.limit
}
