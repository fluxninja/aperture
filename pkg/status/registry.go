package status

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/mitchellh/copystructure"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
)

const defaultDelim = "."

// Registry .
type Registry interface {
	Delim() string
	Prefix() string
	Keys() []string
	Exists(string) bool
	Get(string) *statusv1.GroupStatus
	GetAllFlat() (map[string]*statusv1.GroupStatus, error)
	Push(string, *statusv1.Status) error
	Delete(string) error
}

// registry implements Registry.
type registry struct {
	mu        sync.Mutex
	statusMap map[string]*statusv1.GroupStatus
	delim     string
}

var _ Registry = (*registry)(nil)

// NewRegistry returns a new instance of Registry.
// Delim is the delimiter to use when specifying key paths,
// e.g., . For "parent.child.key" or / for "parent/child/key".
func NewRegistry(delim string) Registry {
	if delim == "" {
		delim = defaultDelim
	}

	return &registry{
		statusMap: make(map[string]*statusv1.GroupStatus),
		delim:     delim,
	}
}

// Delim returns the delimiter used by the Registry.
func (reg *registry) Delim() string {
	return reg.delim
}

// Prefix returns the prefix string used in Registry.
// In the case of registry, it is always an empty string.
func (reg *registry) Prefix() string {
	return ""
}

// Keys returns all the keys stored in the Registry in order.
func (reg *registry) Keys() []string {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	keyMap := generateKeyMap(reg.statusMap, nil, reg.delim)

	keys := make([]string, 0, len(keyMap))
	for key := range keyMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Exists returns true if the given path exists in the registry.
func (reg *registry) Exists(path string) bool {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	_, ok := existsInMap(reg.statusMap, reg.delim, path)
	return ok
}

// Get returns the *statusv1.GroupStatus from the provided path in registry.
// If the path does not exist, nil is returned.
func (reg *registry) Get(path string) *statusv1.GroupStatus {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	paths := strings.Split(path, reg.delim)
	gs := searchMap(reg.statusMap, paths)

	out, _ := copystructure.Copy(gs)
	sm, ok := out.(*statusv1.GroupStatus)
	if !ok {
		return nil
	}
	return sm
}

// GetAllFlat returns entire flattened map[string]*statusv1.GroupStatus of registry.
func (reg *registry) GetAllFlat() (map[string]*statusv1.GroupStatus, error) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	sm := make(map[string]*statusv1.GroupStatus)

	keyMap := generateKeyMap(reg.statusMap, nil, reg.delim)
	for k := range keyMap {
		paths := strings.Split(k, reg.delim)
		gs := searchMap(reg.statusMap, paths)
		sm[k] = gs
	}

	return sm, nil
}

// Push adds a new result to the provided path.
func (reg *registry) Push(path string, status *statusv1.Status) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	return pushToMap(reg.statusMap, reg.delim, path, status)
}

// Delete removes all nested values from a given path.
// Clears all keys/values if no path is specified.
// Every empty, key on the path, is recursively deleted.
func (reg *registry) Delete(path string) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	// If no path is provided, empty the whole registry.
	if path == "" {
		reg.statusMap = make(map[string]*statusv1.GroupStatus)
		return nil
	}

	keyPath, ok := existsInMap(reg.statusMap, reg.delim, path)
	if !ok {
		return errors.New("path doesn't exist")
	}

	removeFromMap(reg.statusMap, keyPath)
	return nil
}
