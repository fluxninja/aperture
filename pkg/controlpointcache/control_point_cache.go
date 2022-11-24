package controlpointcache

import (
	"sync"
)

// Provide provides control point cache.
func Provide() *ControlPointCache {
	return NewControlPointCache()
}

// NewControlPointCache returns new instance of control point cache.
func NewControlPointCache() *ControlPointCache {
	return &ControlPointCache{
		controlPoints: map[ControlPoint]struct{}{},
	}
}

// ControlPointCache keeps information about control points their services.
type ControlPointCache struct {
	sync.Mutex
	controlPoints map[ControlPoint]struct{}
}

// Put inserts control point with given name and service.
func (c *ControlPointCache) Put(name, service string) {
	c.Lock()
	defer c.Unlock()
	c.controlPoints[ControlPoint{Name: name, Service: service}] = struct{}{}
}

// GetAllAndClear returns the current state of cache and clears the cache.
func (c *ControlPointCache) GetAllAndClear() map[ControlPoint]struct{} {
	c.Lock()
	defer c.Unlock()
	result := c.controlPoints
	c.controlPoints = map[ControlPoint]struct{}{}
	return result
}

// ControlPoint represents control point kept in the cache.
type ControlPoint struct {
	Name    string
	Service string
}
