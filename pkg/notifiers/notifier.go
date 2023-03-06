package notifiers

import (
	guuid "github.com/google/uuid"
)

// Key is the key that is used to identify the notifier.
type Key string

// String returns the string representation of the notifier.
func (key Key) String() string {
	return string(key)
}

type notifier interface {
	getID() string
	setID(string)
	GetTransformFunc() TransformFunc
	SetTransformFunc(TransformFunc)
}

// TransformFunc is the callback signature that other notifiers (like NewPrefixToEtcdNotifier and NewKeyToEtcdNotifier) can use to receive notification before contents are processed.
// This function can transform the key and contents before processing them further.
type TransformFunc func(key Key, bytes []byte, etype EventType) (Key, []byte, error)

// notifierBase is the base type for all notifiers.
type notifierBase struct {
	tf TransformFunc
	id string
}

func (idb *notifierBase) getID() string {
	if idb.id == "" {
		idb.id = guuid.NewString()
	}
	return idb.id
}

func (idb *notifierBase) setID(id string) {
	idb.id = id
}

// GetTransformFunc returns the transform function.
func (idb *notifierBase) GetTransformFunc() TransformFunc {
	return idb.tf
}

// SetTransformFunc sets the transform function.
func (idb *notifierBase) SetTransformFunc(tf TransformFunc) {
	idb.tf = tf
}

/////////////////////////////////////////////////////////////////////////////////////////////
// KeyNotifier
/////////////////////////////////////////////////////////////////////////////////////////////

// KeyBase is interface for key.
type KeyBase interface {
	notifier
	GetKey() Key
	inherit(key Key, pn PrefixNotifier)
}

// keyBase is the base type for all key notifiers.
type keyBase struct {
	notifierBase
	key Key
}

// NewKeyBase creates a new key notifier.
func NewKeyBase(key Key) *keyBase {
	return &keyBase{
		key: key,
	}
}

// GetKey returns the key.
func (knb *keyBase) GetKey() Key {
	return knb.key
}

func (knb *keyBase) inherit(key Key, pn PrefixNotifier) {
	knb.key = key
	knb.setID(pn.getID())
	knb.SetTransformFunc(pn.GetTransformFunc())
}

// KeyNotifier is the interface that all key notifiers must implement.
type KeyNotifier interface {
	KeyBase
	Notify(Event)
}

type keyNotifiers []KeyNotifier

func (kns keyNotifiers) notify(event Event) {
	for _, kn := range kns {
		transformNotify(kn, event)
	}
}

func transformNotify(kn KeyNotifier, event Event) {
	transKey := event.Key
	transVal := event.Value
	var err error
	if tf := kn.GetTransformFunc(); tf != nil {
		transKey, transVal, err = tf(event.Key, event.Value, event.Type)
		if err != nil {
			return
		}
	}
	ev := Event{
		Key:   transKey,
		Value: transVal,
		Type:  event.Type,
	}
	// dispatch notification
	kn.Notify(ev)
}

/////////////////////////////////////////////////////////////////////////////////////////////
// PrefixNotifier
/////////////////////////////////////////////////////////////////////////////////////////////

// PrefixBase is the base type for all prefix notifiers.
type PrefixBase interface {
	notifier
	GetPrefix() string
}

// prefixBase is the base type for all prefix notifiers.
type prefixBase struct {
	prefix string
	notifierBase
}

// NewPrefixBase creates a new prefix notifier.
func NewPrefixBase(prefix string) *prefixBase {
	return &prefixBase{
		prefix: prefix,
	}
}

// GetPrefix returns the prefix.
func (pnb *prefixBase) GetPrefix() string {
	return pnb.prefix
}

// PrefixNotifier is the interface that all prefix notifiers must implement.
type PrefixNotifier interface {
	PrefixBase
	GetKeyNotifier(key Key) (KeyNotifier, error)
}

type prefixNotifiers []PrefixNotifier
