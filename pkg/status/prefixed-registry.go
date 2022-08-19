package status

import (
	"strings"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
)

type registryPrefix struct {
	Registry
	prefix string
}

var _ Registry = (*registryPrefix)(nil)

// NewRegistryPrefix wraps a Registry instance so that all requests are prefixed with a given string.
func NewRegistryPrefix(reg Registry, prefix string) Registry {
	return &registryPrefix{reg, prefix}
}

// Delim returns the delimiter used by the underlying Registry.
func (reg *registryPrefix) Delim() string {
	return reg.Registry.Delim()
}

// Prefix returns the prefix string used in Registry.
// In the case of registryPrefix, it is an appended string of prefixes all the way down to the underlying Registry.
func (reg *registryPrefix) Prefix() string {
	if reg.Registry.Prefix() == "" {
		return reg.prefix
	}
	return strings.Join([]string{reg.Registry.Prefix(), reg.prefix}, reg.Delim())
}

// Keys returns all the keys stored in the Registry in order.
// In the case of registryPrefix, it returns only the keys that belong to it.
func (reg *registryPrefix) Keys() []string {
	keys := reg.Registry.Keys()
	prefixedKeys := make([]string, 0)
	for _, key := range keys {
		if strings.HasPrefix(key, reg.Prefix()) {
			prefixedKeys = append(prefixedKeys, key)
		}
	}
	return prefixedKeys
}

// Exists returns true if the given path exists in the registryPrefix.
func (reg *registryPrefix) Exists(path string) bool {
	prefixedPath := reg.prefix
	if path != "" {
		prefixedPath = strings.Join([]string{prefixedPath, path}, reg.Delim())
	}
	return reg.Registry.Exists(prefixedPath)
}

// Get returns the *statusv1.GroupStatus from the provided path in registryPrefix.
// If the path does not exist, nil is returned.
func (reg *registryPrefix) Get(path string) *statusv1.GroupStatus {
	prefixedPath := reg.prefix
	if path != "" {
		prefixedPath = strings.Join([]string{prefixedPath, path}, reg.Delim())
	}
	return reg.Registry.Get(prefixedPath)
}

// GetAllFlat returns entire flattened map[string]*statusv1.GroupStatus of registryPrefix.
func (reg *registryPrefix) GetAllFlat() (map[string]*statusv1.GroupStatus, error) {
	return nil, nil
}

// Push adds a new result to the provided path in registryPrefix.
func (reg *registryPrefix) Push(path string, status *statusv1.Status) error {
	prefixedPath := reg.prefix
	if path != "" {
		prefixedPath = strings.Join([]string{prefixedPath, path}, reg.Delim())
	}

	return reg.Registry.Push(prefixedPath, status)
}

// Delete removes all nested values from a given path in registryPrefix.
// Clears all keys/values if no path is specified.
// Every empty, key on the path, is recursively deleted.
func (reg *registryPrefix) Delete(path string) error {
	prefixedPath := reg.prefix
	if path != "" {
		prefixedPath = strings.Join([]string{prefixedPath, path}, reg.Delim())
	}
	return reg.Registry.Delete(prefixedPath)
}
