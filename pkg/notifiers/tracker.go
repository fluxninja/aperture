package notifiers

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
)

//go:generate mockgen -source tracker.go -package mocks -aux_files github.com/fluxninja/aperture/pkg/notifiers=./watcher.go -destination ../mocks/mock-trackers.go

// Per key trackers.
type keyTracker struct {
	key       Key
	value     []byte
	notifiers keyNotifiers
	validKey  bool
}

func newTracker(key Key) *keyTracker {
	tracker := &keyTracker{
		key:       key,
		notifiers: make(keyNotifiers, 0),
		value:     nil,
	}
	return tracker
}

func (tracker *keyTracker) notify(eventType EventType, value []byte) {
	if !bytes.Equal(value, tracker.value) {
		// send last value if the event is the remove event and value is nil
		if eventType == Remove && value == nil {
			value = tracker.value
		} else {
			tracker.value = value
		}
		tracker.notifiers.notify(Event{
			Type:  eventType,
			Key:   tracker.key,
			Value: value,
		})
	}
	switch eventType {
	case Write:
		tracker.validKey = true
	case Remove:
		tracker.validKey = false
		tracker.value = nil
	}
}

func (tracker *keyTracker) addKeyNotifier(notifier KeyNotifier) {
	// check existing notifier
	for _, n := range tracker.notifiers {
		if n.getID() == notifier.getID() {
			// already exists
			return
		}
	}
	tracker.notifiers = append(tracker.notifiers, notifier)
	if tracker.isValidKey() {
		transformNotify(notifier, Event{
			Type:  Write,
			Key:   tracker.key,
			Value: tracker.value,
		})
	}
}

func (tracker *keyTracker) removeKeyNotifier(id string) {
	for i, n := range tracker.notifiers {
		if n.getID() == id {
			transformNotify(n, Event{
				Type:  Remove,
				Key:   tracker.key,
				Value: nil,
			})
			// remove the key notifier
			tracker.notifiers[i] = tracker.notifiers[len(tracker.notifiers)-1]
			tracker.notifiers = tracker.notifiers[:len(tracker.notifiers)-1]
			return
		}
	}
}

func (tracker *keyTracker) isValidKey() bool {
	return tracker.validKey
}

func (tracker *keyTracker) getKeyNotifiers() keyNotifiers {
	return tracker.notifiers
}

func (tracker *keyTracker) String() string {
	notifiers := []string{}
	for _, n := range tracker.notifiers {
		notifiers = append(notifiers, n.GetKey().String())
	}
	tStr := fmt.Sprintf("key: %s, value: %s, notifiers: %+v", tracker.key, tracker.value, notifiers)
	return tStr
}

////////////////////////////////////////////////////////////////////////////////
// Trackers
////////////////////////////////////////////////////////////////////////////////

const (
	add = iota
	remove
	update
	purge
	stop
)

type notifierOp struct {
	keyNotifier    KeyNotifier
	prefixNotifier PrefixNotifier
	updateKey      Key
	updateFunc     UpdateValueFunc
	purgePrefix    string
	op             int
}

// UpdateValueFunc is a function that can be used to update the value of an existing tracker entry.
type UpdateValueFunc func(oldValue []byte) (EventType, []byte)

// EventWriter can be used to inject events into a tracker collection.
type EventWriter interface {
	WriteEvent(key Key, value []byte)
	RemoveEvent(key Key)
	Purge(prefix string)
	UpdateValue(key Key, updateFunc UpdateValueFunc)
}

// Trackers is the interface of a tracker collection.
type Trackers interface {
	Watcher
	EventWriter
}

// DefaultTrackers is a collection of key trackers.
type DefaultTrackers struct {
	waitGroup        sync.WaitGroup
	ctx              context.Context
	trackers         map[Key]*keyTracker
	notifiersChannel chan notifierOp
	eventsChannel    chan Event
	cancel           context.CancelFunc
	prefixNotifiers  prefixNotifiers
}

// Make sure Trackers implements Watcher interface.
var _ Watcher = &DefaultTrackers{}

// Make sure Trackers implements Trackers interface.
var _ Trackers = &DefaultTrackers{}

// NewDefaultTrackers creates a new instance of Trackers.
func NewDefaultTrackers() *DefaultTrackers {
	t := &DefaultTrackers{
		trackers:         make(map[Key]*keyTracker),
		notifiersChannel: make(chan notifierOp),
		eventsChannel:    make(chan Event),
		prefixNotifiers:  make(prefixNotifiers, 0),
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	return t
}

func (t *DefaultTrackers) getTracker(key Key) (*keyTracker, bool) {
	tracker, ok := t.trackers[key]
	if !ok {
		tracker = newTracker(key)
		t.trackers[key] = tracker
	}
	return tracker, ok
}

func (t *DefaultTrackers) getKeys() []Key {
	keys := make([]Key, 0)
	for key := range t.trackers {
		keys = append(keys, key)
	}
	return keys
}

// WriteEvent sends a Write event with the given key and value to the underlying event channel.
func (t *DefaultTrackers) WriteEvent(key Key, value []byte) {
	t.eventsChannel <- Event{
		Type:  Write,
		Key:   key,
		Value: value,
	}
}

// RemoveEvent sends a Remove event with the given key and value to the underlying event channel.
func (t *DefaultTrackers) RemoveEvent(key Key) {
	t.eventsChannel <- Event{
		Type:  Remove,
		Key:   key,
		Value: nil,
	}
}

// AddKeyNotifier is a convenience function to add a key notifier to the underlying trackers.
// If the key of the given notifier is already tracked, the notifier will be added to the existing tracker.
func (t *DefaultTrackers) AddKeyNotifier(notifier KeyNotifier) error {
	op := notifierOp{
		op:             add,
		keyNotifier:    notifier,
		prefixNotifier: nil,
	}
	t.notifiersChannel <- op
	return nil
}

func (t *DefaultTrackers) addKeyNotifier(n KeyNotifier) {
	key := n.GetKey()
	tracker, _ := t.getTracker(key)
	tracker.addKeyNotifier(n)
}

// RemoveKeyNotifier is a convenience function to remove a key notifier from the underlying trackers.
// If the key of the given notifier is not tracked, the notifier will be ignored.
func (t *DefaultTrackers) RemoveKeyNotifier(notifier KeyNotifier) error {
	op := notifierOp{
		op:             remove,
		keyNotifier:    notifier,
		prefixNotifier: nil,
	}
	t.notifiersChannel <- op
	return nil
}

func (t *DefaultTrackers) removeKeyNotifier(key Key, id string) {
	tracker, _ := t.getTracker(key)
	tracker.removeKeyNotifier(id)
	// if tracker has no notifiers, remove it
	if len(tracker.getKeyNotifiers()) == 0 && !tracker.isValidKey() {
		delete(t.trackers, key)
	}
}

// AddPrefixNotifier is a convenience function to add a prefix notifier to the underlying trackers.
// Internally, a key notifier is added for each key under the given prefix.
// If the prefix of the given notifier is already tracked, the notifier will be added to the existing tracker.
func (t *DefaultTrackers) AddPrefixNotifier(notifier PrefixNotifier) error {
	op := notifierOp{
		op:             add,
		prefixNotifier: notifier,
		keyNotifier:    nil,
	}
	t.notifiersChannel <- op
	return nil
}

func (t *DefaultTrackers) addPrefixNotifier(notifier PrefixNotifier) {
	t.prefixNotifiers = append(t.prefixNotifiers, notifier)
	// add to existing trackers
	for _, key := range t.getKeys() {
		if strings.HasPrefix(key.String(), notifier.GetPrefix()) {
			kn, err := notifier.GetKeyNotifier(key)
			if err != nil {
				continue
			}
			kn.inherit(key, notifier)
			t.addKeyNotifier(kn)
		}
	}
}

// RemovePrefixNotifier is a convenience function to remove a prefix notifier from the underlying trackers.
// Internally, a key notifier is removed for each key under the given prefix.
// If the prefix of the given notifier is not tracked, the notifier will be ignored.
func (t *DefaultTrackers) RemovePrefixNotifier(notifier PrefixNotifier) error {
	op := notifierOp{
		op:             remove,
		prefixNotifier: notifier,
		keyNotifier:    nil,
	}
	t.notifiersChannel <- op
	return nil
}

func (t *DefaultTrackers) removePrefixNotifier(notifier PrefixNotifier) {
	id := notifier.getID()
	for i, notifier := range t.prefixNotifiers {
		if notifier.getID() == id {
			t.prefixNotifiers = append(t.prefixNotifiers[:i], t.prefixNotifiers[i+1:]...)
			break
		}
	}
	// remove from trackers by iterating over all tracker keys
	for _, key := range t.getKeys() {
		if strings.HasPrefix(key.String(), notifier.GetPrefix()) {
			t.removeKeyNotifier(key, id)
		}
	}
}

// Purge is a convenience function to purge all trackers.
// This will remove all key notifiers and prefix notifiers.
func (t *DefaultTrackers) Purge(prefix string) {
	t.notifiersChannel <- notifierOp{
		op:          purge,
		purgePrefix: prefix,
	}
}

func (t *DefaultTrackers) purge(prefix string) {
	for key, tracker := range t.trackers {
		// if key is not a prefix of the purge prefix, skip it
		if !strings.HasPrefix(key.String(), prefix) {
			continue
		}
		// remove all prefix notifiers
		for _, pn := range t.prefixNotifiers {
			t.removeKeyNotifier(key, pn.getID())
		}
		tracker.notify(Remove, nil)
	}
}

// UpdateValue returns the current value tracked by a key.
func (t *DefaultTrackers) UpdateValue(key Key, updateFunc UpdateValueFunc) {
	t.notifiersChannel <- notifierOp{
		op:         update,
		updateKey:  key,
		updateFunc: updateFunc,
	}
}

func (t *DefaultTrackers) updateValue(key Key, updateFunc UpdateValueFunc) {
	tracker, _ := t.getTracker(key)
	eventType, newValue := updateFunc(tracker.value)
	event := Event{
		Type:  eventType,
		Key:   key,
		Value: newValue,
	}
	switch eventType {
	case Write:
		t.writeEvent(tracker, event)
	case Remove:
		t.removeEvent(tracker, event)
	}
}

func (t *DefaultTrackers) writeEvent(tracker *keyTracker, event Event) {
	valid := tracker.isValidKey()
	tracker.notify(Write, event.Value)
	// if the key was not valid earlier, then this is a create event
	if !valid {
		for _, pn := range t.prefixNotifiers {
			if strings.HasPrefix(event.Key.String(), pn.GetPrefix()) {
				n, err := pn.GetKeyNotifier(event.Key)
				if err != nil {
					continue
				}
				n.inherit(event.Key, pn)
				tracker.addKeyNotifier(n)
			}
		}
	}
}

func (t *DefaultTrackers) removeEvent(tracker *keyTracker, event Event) {
	for _, n := range t.prefixNotifiers {
		tracker.removeKeyNotifier(n.getID())
	}
	tracker.notify(Remove, nil)
	if len(tracker.getKeyNotifiers()) == 0 {
		delete(t.trackers, event.Key)
	}
}

// Start opens the underlying event channel and starts the event loop.
// See AddKeyNotifier, AddPrefixNotifier, RemoveKeyNotifier, RemovePrefixNotifier, and Purge for more information.
func (t *DefaultTrackers) Start() error {
	t.waitGroup.Add(1)
	panichandler.Go(func() {
		defer t.waitGroup.Done()
	OUTER:
		for {
			select {
			case op := <-t.notifiersChannel:
				switch op.op {
				case add:
					if op.keyNotifier != nil {
						t.addKeyNotifier(op.keyNotifier)
					} else if op.prefixNotifier != nil {
						t.addPrefixNotifier(op.prefixNotifier)
					}
				case remove:
					if op.keyNotifier != nil {
						t.removeKeyNotifier(op.keyNotifier.GetKey(), op.keyNotifier.getID())
					} else if op.prefixNotifier != nil {
						t.removePrefixNotifier(op.prefixNotifier)
					}
				case update:
					t.updateValue(op.updateKey, op.updateFunc)
				case purge:
					t.purge(op.purgePrefix)
				case stop:
					t.stop()
				}
			case event := <-t.eventsChannel:
				tracker, _ := t.getTracker(event.Key)
				switch event.Type {
				case Write:
					t.writeEvent(tracker, event)
				case Remove:
					t.removeEvent(tracker, event)
				}
			case <-t.ctx.Done():
				break OUTER
			}
		}
		if len(t.prefixNotifiers) > 0 {
			log.Warn().Msg("non-zero prefix notifiers detected on notifier shutdown")
		}
		if len(t.trackers) > 0 {
			// loop through trackers
			for _, tracker := range t.trackers {
				if len(tracker.getKeyNotifiers()) > 0 {
					log.Warn().Interface("key", tracker.key).Msg("dangling notifier detected on shutdown")
				}
			}
		}
	})
	return nil
}

// Stop closes all channels and waits for the goroutine to finish.
func (t *DefaultTrackers) Stop() error {
	op := notifierOp{
		op: stop,
	}
	t.notifiersChannel <- op
	t.waitGroup.Wait()
	return nil
}

func (t *DefaultTrackers) stop() {
	t.cancel()
}

// NewPrefixedEventWriter returns an event writer which keys will be
// automatically prefixed with given prefix.
//
// It's recommended that prefix ends up some kind of delimiter, like `.` or `/`.
func NewPrefixedEventWriter(prefix string, ew EventWriter) EventWriter {
	return &prefixedEventWriter{
		prefix: prefix,
		parent: ew,
	}
}

type prefixedEventWriter struct {
	parent EventWriter
	prefix string
}

// WriteEvent implements EventWriter interface.
func (ew *prefixedEventWriter) WriteEvent(key Key, value []byte) {
	ew.parent.WriteEvent(Key(ew.prefix+string(key)), value)
}

// RemoveEvent implements EventWriter interface.
func (ew *prefixedEventWriter) RemoveEvent(key Key) {
	ew.parent.RemoveEvent(Key(ew.prefix + string(key)))
}

// Purge implements EventWriter interface.
func (ew *prefixedEventWriter) Purge(prefix string) {
	ew.parent.Purge(ew.prefix + prefix)
}

// UpdateValue implements EventWriter interface.
func (ew *prefixedEventWriter) UpdateValue(key Key, updateFunc UpdateValueFunc) {
	ew.parent.UpdateValue(Key(ew.prefix+string(key)), updateFunc)
}
