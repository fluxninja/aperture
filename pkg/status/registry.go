package status

import (
	"fmt"
	"strings"
	"sync"
	"time"

	statusv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/pkg/alerts"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Registry .
type Registry interface {
	GetStatus() *statusv1.Status
	SetStatus(*statusv1.Status)
	SetGroupStatus(*statusv1.GroupStatus)
	GetGroupStatus() *statusv1.GroupStatus
	Child(key, value string) Registry
	ChildIfExists(key, value string) Registry
	Parent() Registry
	Root() Registry
	Detach()
	Key() string
	Value() string
	URI() string
	HasError() bool
	GetLogger() *log.Logger
	GetAlerter() alerts.Alerter
}

var _ Registry = &registry{}

const (
	uriKey       = "status_uri"
	alertChannel = "status_registry"
	// Resolve timeout in seconds.
	alertResolveTimeout = 10
)

// registry implements Registry.
// Note: Please take locks from parent to child and not the other way around to avoid deadlocks.
type registry struct {
	mu       sync.RWMutex
	status   *statusv1.Status
	root     *registry
	parent   *registry
	children map[kv]*registry
	logger   *log.Logger
	key      string
	value    string
	uri      string
	alerter  alerts.Alerter
}

// helper struct for key-value keys in children map.
type kv struct {
	Key   string
	Value string
}

func toString(p kv) string {
	return fmt.Sprintf("%s:%s", p.Key, p.Value)
}

func fromString(hash string) kv {
	split := strings.Split(hash, ":")
	return kv{Key: split[0], Value: split[1]}
}

// NewRegistry creates a new Registry.
func NewRegistry(logger *log.Logger, alerter alerts.Alerter) Registry {
	labeledAlerter := alerter.WithLabels(map[string]string{uriKey: "/"})
	r := &registry{
		key:      "root",
		value:    "root",
		uri:      "",
		parent:   nil,
		status:   &statusv1.Status{},
		children: make(map[kv]*registry),
		logger:   logger,
		alerter:  labeledAlerter,
	}
	r.root = r
	return r
}

// Child creates a new Registry with the given key and value.
func (r *registry) Child(key, value string) Registry {
	r.mu.Lock()
	defer r.mu.Unlock()

	childKV := kv{Key: key, Value: value}
	uri := fmt.Sprintf("%s/%s/%s", r.uri, key, value)
	labeledAlerter := r.alerter.WithLabels(map[string]string{uriKey: uri, key: value})
	child, ok := r.children[childKV]
	if !ok {
		child = &registry{
			key:      key,
			value:    value,
			uri:      uri,
			parent:   r,
			root:     r.root,
			status:   &statusv1.Status{},
			children: make(map[kv]*registry),
			logger:   r.logger.WithStr(r.key, key),
			alerter:  labeledAlerter,
		}
		r.children[childKV] = child
	}
	return child
}

// ChildIfExists returns the child Registry with the given key and value if it exists.
func (r *registry) ChildIfExists(key, value string) Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	childKV := kv{Key: key, Value: value}
	child, ok := r.children[childKV]
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

	// We do not have Attach() so parent cannot change to other than nil.
	if r.parent != nil {
		// remove child from parent
		childKV := kv{Key: r.key, Value: r.value}
		if r.parent.children[childKV] == r {
			delete(r.parent.children, childKV)
		}
		// set parent to nil
		r.parent = nil
		r.logger = r.root.logger
		r.alerter = r.root.alerter
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

	if r.status != nil && r.status.Error != nil {
		r.alerter.AddAlert(r.createAlert(r.status.Error))
	}
}

func (r *registry) createAlert(err *statusv1.Status_Error) *alerts.Alert {
	resolve := time.Duration(time.Second * alertResolveTimeout)
	newAlert := alerts.NewAlert(
		alerts.WithName(err.String()),
		alerts.WithSeverity(alerts.ParseSeverity("info")),
		alerts.WithAlertChannels([]string{alertChannel}),
		alerts.WithResolveTimeout(resolve),
		alerts.WithGeneratorURL(
			fmt.Sprintf("http://%s%s", info.GetHostInfo().Hostname, r.uri),
		),
	)

	return newAlert
}

// SetGroupStatus sets the status of the Registry.
func (r *registry) SetGroupStatus(groupStatus *statusv1.GroupStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = groupStatus.Status
	for hash, gs := range groupStatus.Groups {
		kv := fromString(hash)
		r.Child(kv.Key, kv.Value).SetGroupStatus(gs)
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

	for kv, child := range r.children {
		groupStatus.Groups[toString(kv)] = child.GetGroupStatus()
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

// Value returns the value of the Registry that is registered with the parent.
func (r *registry) Value() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.value
}

// URI returns the uri of the Registry that is registered with the parent.
func (r *registry) URI() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.uri
}

// HasError returns true if the Registry has an error.
func (r *registry) HasError() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.status.GetError() != nil {
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

// GetAlerter returns the alerter of the Registry.
func (r *registry) GetAlerter() alerts.Alerter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.alerter
}
