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
// This function can transform the contents before processing them further.
type TransformFunc func(key Key, bytes []byte, etype EventType) ([]byte, error)

// NotifierBase is the base type for all notifiers.
type NotifierBase struct {
	tf TransformFunc
	id string
}

func (idb *NotifierBase) getID() string {
	if idb.id == "" {
		idb.id = guuid.NewString()
	}
	return idb.id
}

func (idb *NotifierBase) setID(id string) {
	idb.id = id
}

// GetTransformFunc returns the transform function.
func (idb *NotifierBase) GetTransformFunc() TransformFunc {
	return idb.tf
}

// SetTransformFunc sets the transform function.
func (idb *NotifierBase) SetTransformFunc(tf TransformFunc) {
	idb.tf = tf
}

/////////////////////////////////////////////////////////////////////////////////////////////
// KeyNotifier
/////////////////////////////////////////////////////////////////////////////////////////////

// KeyNotifier is the interface that all key notifiers must implement.
type KeyNotifier interface {
	notifier
	GetKey() Key
	SetKey(Key)
	Notify(Event)
	inherit(key Key, pn PrefixNotifier)
}

func transformNotify(kn KeyNotifier, event Event) {
	transVal := event.Value
	var err error
	if tf := kn.GetTransformFunc(); tf != nil {
		transVal, err = tf(event.Key, event.Value, event.Type)
		if err != nil {
			return
		}
	}
	ev := Event{
		Key:   event.Key,
		Value: transVal,
		Type:  event.Type,
	}
	// dispatch notification
	kn.Notify(ev)
}

// KeyNotifierBase is the base type for all key notifiers.
type KeyNotifierBase struct {
	NotifierBase
	key Key
}

// GetKey returns the key.
func (knb *KeyNotifierBase) GetKey() Key {
	return knb.key
}

// SetKey sets the key.
func (knb *KeyNotifierBase) SetKey(key Key) {
	knb.key = key
}

func (knb *KeyNotifierBase) inherit(key Key, pn PrefixNotifier) {
	knb.SetKey(key)
	knb.setID(pn.getID())
	knb.SetTransformFunc(pn.GetTransformFunc())
}

type keyNotifiers []KeyNotifier

func (kns keyNotifiers) notify(event Event) {
	for _, kn := range kns {
		transformNotify(kn, event)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////
// PrefixNotifier
/////////////////////////////////////////////////////////////////////////////////////////////

// PrefixNotifier is the interface that all prefix notifiers must implement.
type PrefixNotifier interface {
	notifier
	GetPrefix() string
	GetKeyNotifier(key Key) KeyNotifier
}

// PrefixNotifierBase is the base type for all prefix notifiers.
type PrefixNotifierBase struct {
	NotifierBase
	Prefix string
}

// GetPrefix returns the prefix.
func (pnb *PrefixNotifierBase) GetPrefix() string {
	return pnb.Prefix
}

type prefixNotifiers []PrefixNotifier
