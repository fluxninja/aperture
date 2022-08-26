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
	Exists() bool
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
// Delim is the delimiter to use when specifying key paths,
// e.g., . For "parent.child.key" or / for "parent/child/key".
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

	return getKeys(reg)
}

func getKeys(reg *registry) []string {
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

// Exists returns true if the given path exists in the registry.
func (reg *registry) Exists() bool {
	return existsAtPath(reg, reg.path)
}

func existsAtPath(reg *registry, path string) bool {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	_, ok := existsInMap(reg.statusMap, reg.delim, path)
	return ok
}

// Get returns the *statusv1.GroupStatus from the provided path in registry.
// If the path does not exist, nil is returned.
func (reg *registry) Get() *statusv1.GroupStatus {
	return getAtPath(reg, reg.path)
}

func getAtPath(reg *registry, path string) *statusv1.GroupStatus {
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
func (reg *registry) Push(status *statusv1.Status) error {
	return pushAtPath(reg, reg.path, status)
}

func pushAtPath(reg *registry, path string, status *statusv1.Status) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	return pushToMap(reg.statusMap, reg.delim, path, status)
}

// Delete removes all nested values from registry.
func (reg *registry) Delete() error {
	return deleteAtPath(reg, reg.path)
}

func deleteAtPath(reg *registry, path string) error {
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
