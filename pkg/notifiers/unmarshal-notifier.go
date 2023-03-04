package notifiers

import (
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
)

// UnmarshalNotifyFunc is a function that is called when a config key is written.
type UnmarshalNotifyFunc func(Event, config.Unmarshaller)

// GetUnmarshallerFunc is a function that is called to create a new unmarshaller.
type GetUnmarshallerFunc func(bytes []byte) (config.Unmarshaller, error)

// UnmarshalKeyNotifier holds the state of a key notifier that updates config at a key using the provided unmarshaller.
type UnmarshalKeyNotifier struct {
	Unmarshaller        config.Unmarshaller
	UnmarshalNotifyFunc UnmarshalNotifyFunc
	KeyNotifierBase
}

// Make sure ConfigKeyNotifier implements KeyNotifier.
var _ KeyNotifier = (*UnmarshalKeyNotifier)(nil)

// NewUnmarshalKeyNotifier creates a new instance of ConfigKeyNotifier.
func NewUnmarshalKeyNotifier(key Key,
	unmarshaller config.Unmarshaller,
	unmarshalNotifyFunc UnmarshalNotifyFunc,
) *UnmarshalKeyNotifier {
	notifier := &UnmarshalKeyNotifier{
		KeyNotifierBase: KeyNotifierBase{
			key: key,
		},
		Unmarshaller:        unmarshaller,
		UnmarshalNotifyFunc: unmarshalNotifyFunc,
	}
	return notifier
}

// Notify provides an unmarshaller based on received event.
// It reloads the the bytes into unmarshaller before invoking callback.
func (ukn *UnmarshalKeyNotifier) Notify(event Event) {
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

// UnmarshalPrefixNotifier holds the state of a prefix notifier that updates config at a prefix using the provided unmarshaller.
type UnmarshalPrefixNotifier struct {
	UnmarshalNotifyFunc UnmarshalNotifyFunc
	GetUnmarshallerFunc GetUnmarshallerFunc
	PrefixNotifierBase
}

// Make sure UnmarshalPrefixNotifier implements PrefixNotifier.
var _ PrefixNotifier = (*UnmarshalPrefixNotifier)(nil)

// GetKeyNotifier returns a new key notifier for the given key.
func (upn *UnmarshalPrefixNotifier) GetKeyNotifier(key Key) (KeyNotifier, error) {
	keyNotifier, err := upn.GetUnmarshalKeyNotifier(key)
	return &keyNotifier, err
}

// GetUnmarshalKeyNotifier returns a new unmarshal key notifier for the given key.
func (upn *UnmarshalPrefixNotifier) GetUnmarshalKeyNotifier(key Key) (UnmarshalKeyNotifier, error) {
	// create a new unmarshaller instance (bytes will be reloaded within key notifier
	unmarshaller, err := upn.GetUnmarshallerFunc(nil)
	if err != nil {
		return UnmarshalKeyNotifier{}, err
	}

	notifier := UnmarshalKeyNotifier{
		Unmarshaller:        unmarshaller,
		UnmarshalNotifyFunc: upn.UnmarshalNotifyFunc,
	}

	return notifier, nil
}
