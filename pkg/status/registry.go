package status

import (
	"sync"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// Registry .
type Registry interface {
	GetStatus() *statusv1.Status
	SetStatus(*statusv1.Status)
	SetGroupStatus(*statusv1.GroupStatus)
	GetGroupStatus() *statusv1.GroupStatus
	Child(key string) Registry
	ChildIfExists(key string) Registry
	Parent() Registry
	Root() Registry
	Detach()
	Key() string
	HasError() bool
	GetLogger() *log.Logger
}

var _ Registry = &registry{}

// registry implements Registry.
// Note: Please take locks from parent to child and not the other way around to avoid deadlocks.
type registry struct {
	mu       sync.RWMutex
	status   *statusv1.Status
	root     *registry
	parent   *registry
	children map[string]*registry
	logger   *log.Logger
	key      string
}

// NewRegistry creates a new Registry.
func NewRegistry(logger *log.Logger) Registry {
	r := &registry{
		key:      "root",
		parent:   nil,
		status:   &statusv1.Status{},
		children: make(map[string]*registry),
		logger:   logger,
	}
	r.root = r
	return r
}

// Child creates a new Registry with the given key.
func (r *registry) Child(key string) Registry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var child *registry
	var ok bool
	child, ok = r.children[key]
	if !ok {
		child = &registry{
			key:      key,
			parent:   r,
			root:     r.root,
			status:   &statusv1.Status{},
			children: make(map[string]*registry),
			logger:   r.logger.WithStr(r.key, key),
		}
		r.children[key] = child
	}
	return child
}

// ChildIfExists returns the child Registry with the given key if it exists.
func (r *registry) ChildIfExists(key string) Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	child, ok := r.children[key]
	if !ok {
		return nil
	}
	return child
}

// Parent returns the parent Registry.
func (r *registry) Parent() Registry {
	return r.getParent()
}

func (r *registry) getParent() *registry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.parent
}

// Detach detaches the child from the parent to become root.
func (r *registry) Detach() {
	parent := r.getParent()
	if parent == nil {
		return
	}

	// lock parent
	parent.mu.Lock()
	defer parent.mu.Unlock()
	// lock child
	r.mu.Lock()
	defer r.mu.Unlock()

	// We don't have Attach() so parent can't change to other than nil.
	if r.parent != nil {
		// remove child from parent
		if r.parent.children[r.key] == r {
			delete(r.parent.children, r.key)
		}
		// set parent to nil
		r.parent = nil
		r.logger = r.root.logger
		r.root = r
	}
}

// GetStatus returns the status of the Registry.
func (r *registry) GetStatus() *statusv1.Status {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.status != nil {
		return r.status
	} else {
		return &statusv1.Status{}
	}
}

// SetStatus sets the status of the Registry.
func (r *registry) SetStatus(status *statusv1.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = status
}

// SetGroupStatus sets the status of the Registry.
func (r *registry) SetGroupStatus(groupStatus *statusv1.GroupStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = groupStatus.Status
	for key, gs := range groupStatus.Groups {
		r.Child(key).SetGroupStatus(gs)
	}
}

// GetGroupStatus returns the status of the Registry.
func (r *registry) GetGroupStatus() *statusv1.GroupStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groupStatus := &statusv1.GroupStatus{
		Status: r.status,
		Groups: make(map[string]*statusv1.GroupStatus),
	}

	for _, child := range r.children {
		groupStatus.Groups[child.key] = child.GetGroupStatus()
	}
	return groupStatus
}

// Root returns the top-level Registry.
func (r *registry) Root() Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.root
}

// Key returns the key of the Registry that is registered with the parent.
func (r *registry) Key() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.key
}

// HasError returns true if the Registry has an error.
func (r *registry) HasError() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.status.GetError().GetMessage() != "" {
		return true
	}
	for _, child := range r.children {
		if child.HasError() {
			return true
		}
	}
	return false
}

// GetLogger returns the logger of the Registry.
func (r *registry) GetLogger() *log.Logger {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.logger
}
