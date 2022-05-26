package status

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/mitchellh/copystructure"
	"google.golang.org/protobuf/types/known/emptypb"

	statusv1 "aperture.tech/aperture/api/gen/proto/go/aperture/common/status/v1"
	"aperture.tech/aperture/pkg/log"
)

const defaultDelim = "."

// Registry holds results.
type Registry struct {
	statusv1.UnimplementedStatusServiceServer
	mu            sync.Mutex
	statusMap     map[string]*statusv1.GroupStatus
	statusMapFlat map[string]*statusv1.GroupStatus
	keyMap        KeyMap
	delim         string
}

// KeyMap is a map of key paths to their correspodning possible traversal paths.
type KeyMap map[string][]string

// NewRegistry returns a new instance of Registry.
// Delim is the delimiter to use when specifying key paths,
// e.g., . For "parent.child.key" or / for "parent/child/key".
func NewRegistry(delim string) *Registry {
	if delim == "" {
		delim = defaultDelim
	}

	return &Registry{
		statusMap:     make(map[string]*statusv1.GroupStatus),
		statusMapFlat: make(map[string]*statusv1.GroupStatus),
		keyMap:        make(KeyMap),
		delim:         delim,
	}
}

// Delim returns the delimiter used by the Registry.
func (reg *Registry) Delim() string {
	return reg.delim
}

// GetGroupStatus returns the group status for the requested group in the Registry.
func (reg *Registry) GetGroupStatus(ctx context.Context, req *statusv1.GroupStatusRequest) (*statusv1.GroupStatus, error) {
	log.Trace().Str("group", req.Group).Msg("Received request on Get job handler")

	status := reg.Get(req.Group)

	return status, nil
}

// GetGroups returns the groups from the keys in the Registry.
func (reg *Registry) GetGroups(ctx context.Context, req *emptypb.Empty) (*statusv1.Groups, error) {
	log.Trace().Msg("Received request on GetGroups jobs handler")

	groups := reg.Keys()

	response := &statusv1.Groups{
		Groups: groups,
	}

	return response, nil
}

// Exists returns true if the given key path exists in the result map.
func (reg *Registry) Exists(path string) bool {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	_, ok := reg.keyMap[path]
	return ok
}

func (reg *Registry) getAll() *statusv1.GroupStatus {
	statusMapCopy, _ := copystructure.Copy(reg.statusMap)
	statusMap, ok := statusMapCopy.(map[string]*statusv1.GroupStatus)
	if !ok {
		statusMap = nil
	}

	return &statusv1.GroupStatus{
		Status: &statusv1.Status{},
		Groups: statusMap,
	}
}

// GetAllFlat returns entire flattened map[string]*Results in the Registry.
func (reg *Registry) GetAllFlat() (map[string]*statusv1.GroupStatus, error) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	// var sm map[string]*statusv1.GroupStatus

	out, _ := copystructure.Copy(reg.statusMapFlat)
	sm, ok := out.(map[string]*statusv1.GroupStatus)
	if !ok {
		return nil, fmt.Errorf("could not deep copy status flat map")
	}
	return sm, nil
}

// Get returns the map[string]*Results of a given key path
// in the Registry. If the key path does not exist, nil is returned.
func (reg *Registry) Get(path string) *statusv1.GroupStatus {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if path == "" {
		return reg.getAll()
	}

	p, ok := reg.keyMap[path]
	if !ok {
		return nil
	}
	statusMap := search(reg.statusMap, p)

	out, _ := copystructure.Copy(statusMap)
	sm, ok := out.(*statusv1.GroupStatus)
	if !ok {
		return nil
	}
	return sm
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

// Delete removes all nested values from a given path.
// Clears all keys/values if no path is specified.
// Every empty, key on the path, is recursively deleted.
func (reg *Registry) Delete(path string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	// If no path is provided, empty the whole registry.
	if path == "" {
		reg.statusMap = make(map[string]*statusv1.GroupStatus)
		reg.statusMapFlat = make(map[string]*statusv1.GroupStatus)
		reg.keyMap = make(KeyMap)
		return
	}

	p, ok := reg.keyMap[path]
	if !ok {
		return
	}
	removeAtPath(reg.statusMap, p)

	// Update the flattened maps too.
	reg.statusMapFlat, reg.keyMap = flattenMap(reg.statusMap, nil, reg.delim)
	reg.keyMap = populateKeyParts(reg.keyMap, reg.delim)
}

func removeAtPath(m map[string]*statusv1.GroupStatus, p []string) {
	next, ok := m[p[0]]
	if ok {
		if len(p) == 1 {
			delete(m, p[0])
			return
		}
		next.Status = nil
		if next.Groups != nil {
			removeAtPath(next.Groups, p[1:])
			// Delete map if it has no keys.
			if len(next.Groups) == 0 {
				delete(m, p[0])
			}
		}
	}
}

func flattenMap(m map[string]*statusv1.GroupStatus, keys []string, delim string) (map[string]*statusv1.GroupStatus, KeyMap) {
	flatMap := make(map[string]*statusv1.GroupStatus)
	keyMap := make(KeyMap)

	flatten(m, keys, delim, flatMap, keyMap)
	return flatMap, keyMap
}

func flatten(m map[string]*statusv1.GroupStatus, keys []string, delim string, flatMap map[string]*statusv1.GroupStatus, keyMap KeyMap) {
	for k, v := range m {
		// Copy the incoming key paths into a fresh list
		// and append the current key in the iteration.
		kp := make([]string, 0, len(keys)+1)
		kp = append(kp, keys...)
		kp = append(kp, k)

		if v.Groups != nil {
			if len(v.Groups) == 0 {
				newKey := strings.Join(kp, delim)
				flatMap[newKey] = v
				keyMap[newKey] = kp
				continue
			}

			// There is more to flatten underneath
			flatten(v.Groups, kp, delim, flatMap, keyMap)
		} else {
			newKey := strings.Join(kp, delim)
			flatMap[newKey] = v
			keyMap[newKey] = kp
		}
	}
}

// Keys returns all the keys stored in the Registry keyMap in order.
func (reg *Registry) Keys() []string {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	out := make([]string, 0, len(reg.keyMap))
	for k := range reg.keyMap {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Push adds a new result to the provided path.
func (reg *Registry) Push(path string, status *statusv1.Status) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if path == "" {
		return errors.New("path doesn't exist")
	}

	gs := &statusv1.GroupStatus{
		Status: status,
		Groups: nil,
	}

	m := unflattenMap(
		map[string]*statusv1.GroupStatus{
			path: gs,
		},
		reg.delim,
	)

	return reg.merge(m)
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

func (reg *Registry) merge(m map[string]*statusv1.GroupStatus) error {
	mergeMaps(m, reg.statusMap)

	// Update the flattened maps too.
	reg.statusMapFlat, reg.keyMap = flattenMap(reg.statusMap, nil, reg.delim)
	reg.keyMap = populateKeyParts(reg.keyMap, reg.delim)

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

// populateKeyParts iterates a key map and generates all possible
// traversal paths. For instance, `parent.child.key` generates
// `parent`, and `parent.child`.
func populateKeyParts(m KeyMap, delim string) KeyMap {
	out := make(KeyMap, len(m))
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
