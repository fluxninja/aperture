//go:generate enumer -type=EventType -output=event-type-string.go
package notifiers

import (
	"fmt"
)

// EventType is the type of event.
type EventType uint8

const (
	// Write represents that the watched entity has been written to.
	// In case of a file, it could have been created, modified, or symlinked.
	// In case of an etcd entry, it could have been added or updated.
	Write EventType = 1 << iota
	// Remove represents that the watched entity has been removed.
	Remove
)

// Event is the event that is passed to the notifier.
type Event struct {
	Key
	Value []byte
	Type  EventType
}

// String returns the string representation of the event.
func (event Event) String() string {
	return fmt.Sprintf("Event<"+"EventType: %s "+"| %s"+">", event.Type.String(), event.Key)
}
