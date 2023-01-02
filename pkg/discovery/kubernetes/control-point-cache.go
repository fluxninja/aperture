package kubernetes

// ControlPointCache is a cache of Kubernetes control points.
type ControlPointCache struct {
	// Set of unique ControlPoints
	ControlPoints map[ControlPoint]struct{}
}

// A ControlPoint is identified by Group, Version, Type, Namespace and Name.
type ControlPoint struct {
	Group     string
	Version   string
	Type      string
	Namespace string
	Name      string
}

// newControlPointCache returns a new ControlPointCache.
func newControlPointCache() *ControlPointCache {
	return &ControlPointCache{
		ControlPoints: make(map[ControlPoint]struct{}),
	}
}
