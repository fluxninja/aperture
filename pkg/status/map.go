package status

import (
	"strings"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
)

// kMap is a map of key paths to their corresponding possible traversal paths.
type kMap map[string][]string

func existsInMap(m map[string]*statusv1.GroupStatus, delim string, path string) ([]string, bool) {
	keyMap := generateKeyMap(m, nil, delim)
	keyPath, ok := keyMap[path]
	return keyPath, ok
}

// merge a into b.
func mergeMaps(a, b map[string]*statusv1.GroupStatus) {
	for key, val := range a {
		// If key does not exist in b, add it and continue
		bVal, ok := b[key]
		if !ok {
			b[key] = val
			continue
		}

		if val.Status != nil {
			b[key].Status = val.Status
		}

		// If key exists in b, merge the two
		if bVal.Groups != nil {
			mergeMaps(val.Groups, bVal.Groups)
		} else {
			b[key].Groups = val.Groups
		}
	}
}

func searchMap(m map[string]*statusv1.GroupStatus, p []string) *statusv1.GroupStatus {
	var status *statusv1.GroupStatus
	groups, ok := m[p[0]]
	if ok {
		if len(p) == 1 {
			status = groups
			return status
		}
		if len(groups.Groups) != 0 {
			return searchMap(groups.Groups, p[1:])
		}
	}
	return nil
}

func removeFromMap(m map[string]*statusv1.GroupStatus, p []string) {
	next, ok := m[p[0]]
	if ok {
		if len(p) == 1 {
			delete(m, p[0])
			return
		}
		next.Status = nil
		if next.Groups != nil {
			removeFromMap(next.Groups, p[1:])
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
					Status: nil,
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
