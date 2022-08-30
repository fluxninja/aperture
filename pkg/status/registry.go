package status

import (
	"errors"
	"fmt"
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
	Path() string
	Keys() []string
	Get() *statusv1.GroupStatus
	GetAllFlat() (map[string]*statusv1.GroupStatus, error)
	Push(*statusv1.Status) error
	Delete() error
}

// registry implements Registry.
type registry struct {
	mu        sync.Mutex
	statusMap map[string]*statusv1.GroupStatus
	delim     string
	path      string
}

var _ Registry = (*registry)(nil)

// NewRegistry returns a new instance of Registry.
// It is possible to chain a Registry with another.
func NewRegistry(reg Registry, path string) Registry {
	r := &registry{
		delim: defaultDelim,
	}

	if reg == nil {
		r.statusMap = make(map[string]*statusv1.GroupStatus)
		r.path = path
		return r
	}

	typedReg := reg.(*registry)

	r.statusMap = typedReg.statusMap
	fullPath := typedReg.Path()
	if fullPath == "" {
		fullPath = path
	} else {
		fullPath = fmt.Sprintf("%s%s%s", fullPath, reg.Delim(), path)
	}
	r.path = fullPath

	return r
}

// Delim returns the default delimiter used in registry.
func (reg *registry) Delim() string {
	return reg.delim
}

// Path returns the path used in Registry.
// In the case of registry, it is always an empty string.
func (reg *registry) Path() string {
	return reg.path
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

	prefixedKeys := make([]string, 0)
	for _, key := range keys {
		if strings.HasPrefix(key, reg.path) {
			prefixedKeys = append(prefixedKeys, key)
		}
	}

	sort.Strings(prefixedKeys)
	return prefixedKeys
}

// Get returns the *statusv1.GroupStatus from the provided path in registry.
// If the path does not exist, nil is returned.
func (reg *registry) Get() *statusv1.GroupStatus {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	paths := strings.Split(reg.path, reg.delim)
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
func (reg *registry) Push(status *statusv1.Status) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if reg.path == "" {
		return errors.New("path doesn't exist")
	}

	gs := &statusv1.GroupStatus{
		Status: status,
		Groups: nil,
	}

	um := unflattenMap(
		map[string]*statusv1.GroupStatus{
			reg.path: gs,
		},
		reg.delim,
	)

	mergeMaps(um, reg.statusMap)
	return nil
}

// Delete removes all nested values from registry.
func (reg *registry) Delete() error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	// If no path is provided, empty the whole registry.
	if reg.path == "" {
		reg.statusMap = make(map[string]*statusv1.GroupStatus)
		return nil
	}

	keyMap := generateKeyMap(reg.statusMap, nil, reg.delim)
	keyPath := keyMap[reg.path]

	removeFromMap(reg.statusMap, keyPath)
	return nil
}
