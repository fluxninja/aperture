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
	err                 error
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
func (cfn *UnmarshalKeyNotifier) Notify(event Event) {
	if cfn.Unmarshaller != nil {
		// Only update unmarshaller on write
		if event.Type == Write {
			log.Trace().Str("key", event.Key.String()).Msg("write event")
			// read config via unmarshaller
			err := cfn.Unmarshaller.Reload(event.Value)
			if err != nil {
				log.Error().Err(err).Str("key", event.Key.String()).Msg("error reading")
				return
			}
		}
		if cfn.UnmarshalNotifyFunc != nil {
			// notify
			cfn.UnmarshalNotifyFunc(event, cfn.Unmarshaller)
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
func (cdn *UnmarshalPrefixNotifier) GetKeyNotifier(key Key) KeyNotifier {
	keyNotifier := cdn.getUnmarshalKeyNotifier(key)
	return &keyNotifier
}

func (cdn *UnmarshalPrefixNotifier) getUnmarshalKeyNotifier(key Key) UnmarshalKeyNotifier {
	// create a new unmarshaller instance (bytes will be reloaded within key notifier
	unmarshaller, err := cdn.GetUnmarshallerFunc(nil)
	if err != nil {
		unmarshaller = nil
	}

	notifier := UnmarshalKeyNotifier{
		Unmarshaller:        unmarshaller,
		UnmarshalNotifyFunc: cdn.UnmarshalNotifyFunc,
		err:                 err,
	}

	return notifier
}
