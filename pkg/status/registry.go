package status

import (
	"sync"

	statusv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/common/status/v1"
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

// Detach detaches the child from the parent to become root.
func (r *registry) Detach() {
	if r.parent == nil {
		return
	}
	// lock parent
	r.parent.mu.Lock()
	defer func() {
		r.parent.mu.Unlock()
		r.logger = r.root.logger
		r.parent = nil
		r.root = r
	}()
	if r.parent.children[r.key] == r {
		delete(r.parent.children, r.key)
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

// Parent returns the parent Registry.
func (r *registry) Parent() Registry {
	return r.parent
}

// Root returns the top-level Registry.
func (r *registry) Root() Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.root
}

// Key returns the key of the Registry that is registered with the parent.
func (r *registry) Key() string {
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
	return r.logger
}
