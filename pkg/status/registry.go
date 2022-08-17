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

// Registry holds results.
type Registry struct {
	mu        sync.Mutex
	statusMap map[string]*statusv1.GroupStatus
	delim     string
}

// kMap is a map of key paths to their corresponding possible traversal paths.
type kMap map[string][]string

// NewRegistry returns a new instance of Registry.
// Delim is the delimiter to use when specifying key paths,
// e.g., . For "parent.child.key" or / for "parent/child/key".
func NewRegistry(delim string) *Registry {
	if delim == "" {
		delim = defaultDelim
	}

	return &Registry{
		statusMap: make(map[string]*statusv1.GroupStatus),
		delim:     delim,
	}
}

// Delim returns the delimiter used by the Registry.
func (reg *Registry) Delim() string {
	return reg.delim
}

// Keys returns all the keys stored in the Registry keyMap in order.
func (reg *Registry) Keys() []string {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	keyMap := generateKeyMap(reg.statusMap, nil, reg.delim)

	out := make([]string, 0, len(keyMap))
	for k := range keyMap {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// GetAllFlat returns entire flattened map[string]*Results in the Registry.
func (reg *Registry) GetAllFlat() (map[string]*statusv1.GroupStatus, error) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	sm := make((map[string]*statusv1.GroupStatus))

	keyMap := generateKeyMap(reg.statusMap, nil, reg.delim)
	for k := range keyMap {
		paths := strings.Split(k, reg.delim)
		gs := search(reg.statusMap, paths)
		sm[k] = gs
	}

	return sm, nil
}

// RegistryPathBuilder holds contextual fields for performing further actions.
type RegistryPathBuilder struct {
	*Registry
	path string
}

// At creates a new RegistryPathBuilder with path information.
func (reg *Registry) At(paths ...string) *RegistryPathBuilder {
	b := &RegistryPathBuilder{}
	b.Registry = reg
	b.path = strings.Join(paths, reg.delim)

	return b
}

// Exists returns true if the given key path exists in the result map.
func (b *RegistryPathBuilder) Exists() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	keyMap := generateKeyMap(b.statusMap, nil, b.delim)
	_, ok := keyMap[b.path]
	return ok
}

// Get returns the map[string]*Results of a given key path
// in the Registry. If the key path does not exist, nil is returned.
func (b *RegistryPathBuilder) Get() *statusv1.GroupStatus {
	b.mu.Lock()
	defer b.mu.Unlock()

	paths := strings.Split(b.path, b.delim)
	gs := search(b.statusMap, paths)

	out, _ := copystructure.Copy(gs)
	sm, ok := out.(*statusv1.GroupStatus)
	if !ok {
		return nil
	}
	return sm
}

// Push adds a new result to the provided path.
func (b *RegistryPathBuilder) Push(status *statusv1.Status) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.path == "" {
		return errors.New("path doesn't exist")
	}

	gs := &statusv1.GroupStatus{
		Status: status,
		Groups: nil,
	}

	m := unflattenMap(
		map[string]*statusv1.GroupStatus{
			b.path: gs,
		},
		b.delim,
	)

	return b.merge(m)
}

// Delete removes all nested values from a given path.
// Clears all keys/values if no path is specified.
// Every empty, key on the path, is recursively deleted.
func (b *RegistryPathBuilder) Delete() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// If no path is provided, empty the whole registry.
	if b.path == "" {
		b.statusMap = make(map[string]*statusv1.GroupStatus)
		return
	}

	keyMap := generateKeyMap(b.statusMap, nil, b.delim)
	p, ok := keyMap[b.path]
	if !ok {
		return
	}

	remove(b.statusMap, p)
}

func (reg *Registry) merge(m map[string]*statusv1.GroupStatus) error {
	mergeMaps(m, reg.statusMap)

	return nil
}

// merge a into b.
func mergeMaps(a, b map[string]*statusv1.GroupStatus) {
	for ak, av := range a {
		// If key does not exist in b, add it and continue
		bv, ok := b[ak]
		if !ok {
			b[ak] = av
			continue
		}

		bv.Status = av.Status

		// If key exists in b, merge the two
		if bv.Groups != nil {
			mergeMaps(av.Groups, bv.Groups)
		} else {
			b[ak].Groups = av.Groups
		}
	}
}

func search(m map[string]*statusv1.GroupStatus, p []string) *statusv1.GroupStatus {
	var status *statusv1.GroupStatus
	groups, ok := m[p[0]]
	if ok {
		if len(p) == 1 {
			status = groups
			return status
		}
		if len(groups.Groups) != 0 {
			return search(groups.Groups, p[1:])
		}
	}
	return nil
}

func remove(m map[string]*statusv1.GroupStatus, p []string) {
	next, ok := m[p[0]]
	if ok {
		if len(p) == 1 {
			delete(m, p[0])
			return
		}
		next.Status = nil
		if next.Groups != nil {
			remove(next.Groups, p[1:])
			// Delete map if it has no keys.
			if len(next.Groups) == 0 {
				delete(m, p[0])
			}
		}
	}
}

func unflattenMap(m map[string]*statusv1.GroupStatus, delim string) map[string]*statusv1.GroupStatus {
	out := make(map[string]*statusv1.GroupStatus)

	for k, v := range m {
		keys := strings.Split(k, delim)
		next := out

		// Iterate through key parts, for eg:, parent.child.key
		// will be ["parent", "child", "key"]
		for _, key := range keys[:len(keys)-1] {
			sub, ok := next[key]
			if !ok {
				sub = &statusv1.GroupStatus{
					Status: &statusv1.Status{},
					Groups: make(map[string]*statusv1.GroupStatus),
				}
				next[key] = sub
			}

			if sub.Groups != nil {
				next = sub.Groups
			}
		}

		next[keys[len(keys)-1]] = v
	}
	return out
}

func generateKeyMap(m map[string]*statusv1.GroupStatus, keys []string, delim string) kMap {
	keyMap := make(kMap)

	updateKeyMap(m, keys, delim, keyMap)
	return populateKeyParts(keyMap, delim)
}

func updateKeyMap(m map[string]*statusv1.GroupStatus, keys []string, delim string, keyMap kMap) {
	for k, v := range m {
		// Copy the incoming key paths into a fresh list
		// and append the current key in the iteration.
		kp := make([]string, 0, len(keys)+1)
		kp = append(kp, keys...)
		kp = append(kp, k)

		if v.Groups != nil {
			if len(v.Groups) == 0 {
				newKey := strings.Join(kp, delim)
				keyMap[newKey] = kp
				continue
			}
			// There is more to flatten underneath
			updateKeyMap(v.Groups, kp, delim, keyMap)
		} else {
			newKey := strings.Join(kp, delim)
			keyMap[newKey] = kp
		}
	}
}

// populateKeyParts iterates a key map and generates all possible
// traversal paths. For instance, `parent.child.key` generates
// `parent`, and `parent.child`.
func populateKeyParts(m kMap, delim string) kMap {
	out := make(kMap, len(m))
	for _, parts := range m {
		var nk string
		for i := range parts {
			if i == 0 {
				nk = parts[i]
			} else {
				nk += delim + parts[i]
			}
			if _, ok := out[nk]; ok {
				continue
			}
			out[nk] = make([]string, i+1)
			copy(out[nk], parts[0:i+1])
		}
	}
	return out
}
