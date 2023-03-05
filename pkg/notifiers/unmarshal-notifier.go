package notifiers

import (
	"fmt"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

// UnmarshalNotifyFunc is a function that is called when a config key is written.
type UnmarshalNotifyFunc func(Event, config.Unmarshaller)

// GetUnmarshallerFunc is a function that is called to create a new unmarshaller.
type GetUnmarshallerFunc func(bytes []byte) (config.Unmarshaller, error)

// unmarshalKeyNotifier holds the state of a key notifier that updates config at a key using the provided unmarshaller.
type unmarshalKeyNotifier struct {
	Unmarshaller        config.Unmarshaller
	UnmarshalNotifyFunc UnmarshalNotifyFunc
	KeyBase
}

// Make sure ConfigKeyNotifier implements KeyNotifier.
var _ KeyNotifier = (*unmarshalKeyNotifier)(nil)

// NewUnmarshalKeyNotifier creates a new instance of ConfigKeyNotifier.
func NewUnmarshalKeyNotifier(key Key,
	unmarshaller config.Unmarshaller,
	unmarshalNotifyFunc UnmarshalNotifyFunc,
) (*unmarshalKeyNotifier, error) {
	return &unmarshalKeyNotifier{
		KeyBase:             NewKeyBase(key),
		Unmarshaller:        unmarshaller,
		UnmarshalNotifyFunc: unmarshalNotifyFunc,
	}, nil
}

// Notify provides an unmarshaller based on received event.
// It reloads the the bytes into unmarshaller before invoking callback.
func (ukn *unmarshalKeyNotifier) Notify(event Event) {
	if ukn.Unmarshaller != nil {
		// Only update unmarshaller on write
		if event.Type == Write {
			log.Trace().Str("key", event.Key.String()).Msg("write event")
			// read config via unmarshaller
			err := ukn.Unmarshaller.Reload(event.Value)
			if err != nil {
				log.Error().Err(err).Str("key", event.Key.String()).Msg("error reading")
				return
			}
		}
		if ukn.UnmarshalNotifyFunc != nil {
			// notify
			ukn.UnmarshalNotifyFunc(event, ukn.Unmarshaller)
		}
	}
}

// unmarshalPrefixNotifier holds the state of a prefix notifier that updates config at a prefix using the provided unmarshaller.
type unmarshalPrefixNotifier struct {
	unmarshalNotifyFunc UnmarshalNotifyFunc
	getUnmarshallerFunc GetUnmarshallerFunc
	PrefixBase
}

// Make sure UnmarshalPrefixNotifier implements PrefixNotifier.
var _ PrefixNotifier = (*unmarshalPrefixNotifier)(nil)

// GetKeyNotifier returns a new key notifier for the given key.
func (upn *unmarshalPrefixNotifier) GetKeyNotifier(key Key) (KeyNotifier, error) {
	// create a new unmarshaller instance
	unmarshaller, err := upn.getUnmarshallerFunc(nil)
	if err != nil {
		return nil, err
	}
	return NewUnmarshalKeyNotifier(key, unmarshaller, upn.unmarshalNotifyFunc)
}

// NewUnmarshalPrefixNotifier returns a new instance of UnmarshalPrefixNotifier.
func NewUnmarshalPrefixNotifier(prefix string,
	unmarshalNotifyFunc UnmarshalNotifyFunc,
	getUnmarshallerFunc GetUnmarshallerFunc,
) (*unmarshalPrefixNotifier, error) {
	if getUnmarshallerFunc == nil {
		return nil, fmt.Errorf("getUnmarshallerFunc cannot be nil")
	}
	notifier := &unmarshalPrefixNotifier{
		PrefixBase:          NewPrefixBase(prefix),
		unmarshalNotifyFunc: unmarshalNotifyFunc,
		getUnmarshallerFunc: getUnmarshallerFunc,
	}
	return notifier, nil
}
