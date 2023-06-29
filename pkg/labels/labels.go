package labels

import (
	"sort"

	"golang.org/x/exp/maps"
)

// Labels describe labels map, where each label have a key and value.
//
// This is an interface so that callers are not forced to always convert to the map.
type Labels interface {
	Get(key string) (value string, exists bool)
	SortedKeys() []string
	Copy() PlainMap
}

// PlainMap implements Labels for a plain map.
//
// Nil map is ok.
type PlainMap map[string]string

// Get implements FlowLabels.
func (m PlainMap) Get(key string) (string, bool) {
	value, exists := m[key]
	return value, exists
}

// SortedKeys implements FlowLabels.
func (m PlainMap) SortedKeys() []string {
	keys := maps.Keys(m)
	sort.Strings(keys)
	return keys
}

// Copy implements FlowLabels.
func (m PlainMap) Copy() PlainMap {
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}
