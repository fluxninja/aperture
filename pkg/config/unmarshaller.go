package config

import "github.com/spf13/cast"

// Unmarshaller provides common interface for unmarshallers.
type Unmarshaller interface {
	// Check whether config key is present
	IsSet(string) bool
	// Unmarshal the key into a user provided config struct
	UnmarshalKey(string, interface{}) error
	// Get specific key from config -- use cast lib to convert to bool, string etc.
	Get(string) interface{}
	// Unmarshal entire config into a struct
	Unmarshal(interface{}) error
	// Reload config
	Reload(bytes []byte) error
	// Marshal the config into bytes
	Marshal() ([]byte, error)
}

// GetValue returns the value for the given key if config key is present.
func GetValue(unmarshaller Unmarshaller, key string, defaultVal interface{}) interface{} {
	val := defaultVal
	if key != "" && unmarshaller.IsSet(key) {
		val = unmarshaller.Get(key)
	}
	return val
}

// GetStringValue returns the string value for the given key.
func GetStringValue(unmarshaller Unmarshaller, key string, defaultVal string) string {
	return cast.ToString(GetValue(unmarshaller, key, defaultVal))
}

// GetIntValue returns the integer value for the given key.
func GetIntValue(unmarshaller Unmarshaller, key string, defaultVal int) int {
	return cast.ToInt(GetValue(unmarshaller, key, defaultVal))
}

// GetBoolValue returns the boolean value for the given key.
func GetBoolValue(unmarshaller Unmarshaller, key string, defaultVal bool) bool {
	return cast.ToBool(GetValue(unmarshaller, key, defaultVal))
}
