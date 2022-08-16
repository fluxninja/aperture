package config

import (
	"errors"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/knadh/koanf"
	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"

	"github.com/Henry-Sarabia/sliceconv"

	"github.com/fluxninja/aperture/pkg/log"
)

// Unmarshal unmarshals using _just_ bytes as source of truth (no env, no
// flags, no other overrides).
func Unmarshal(bytes []byte, output interface{}) error {
	un, err := KoanfUnmarshallerConstructor{}.NewKoanfUnmarshaller(bytes)
	if err != nil {
		return err
	}
	return un.Unmarshal(output)
}

// ConfigFormat specifies the type of the configuration format in string.
type ConfigFormat string

const (
	// YAML is one of config formats.
	YAML ConfigFormat = "yaml"
	// JSON is one of config formats.
	JSON ConfigFormat = "json"
)

// KoanfUnmarshallerConstructor holds fields to create an annotated instance of KoanfUnmarshaller.
type KoanfUnmarshallerConstructor struct {
	// Command line flags
	FlagSet *pflag.FlagSet
	// Optional Merge Config
	MergeConfig map[string]interface{}
	// Config format (yaml, json)
	ConfigFormat ConfigFormat
	// Enable AutomaticEnv
	EnableEnv bool
}

// NewKoanfUnmarshaller creates a new Unmarshaller instance that can be used to unmarshal configs.
func (constructor KoanfUnmarshallerConstructor) NewKoanfUnmarshaller(bytes []byte) (Unmarshaller, error) {
	k := koanf.New(DefaultKoanfDelim)

	// Precedence:
	// 1. MergeConfig
	// 2. Env
	// 3. Bytes
	// 4. Flags

	// On runtime, reloaded bytes config take precedence

	if constructor.FlagSet != nil {
		err := k.Load(posflag.Provider(constructor.FlagSet, k.Delim(), k), nil)
		if err != nil {
			log.Error().Err(err).Msg("failed to load flags")
			return nil, err
		}
	}

	// Default to YAML
	if constructor.ConfigFormat == "" {
		constructor.ConfigFormat = YAML
	}

	err := loadBytesProvider(k, bytes, constructor.ConfigFormat)
	if err != nil {
		log.Error().Err(err).Msg("failed to load")
		return nil, err
	}

	if constructor.EnableEnv {
		err = k.Load(env.Provider(EnvPrefix, k.Delim(), func(s string) string {
			return strings.TrimPrefix(s, EnvPrefix)
		}), nil)
		if err != nil {
			return nil, err
		}
	}

	if constructor.MergeConfig != nil {
		err = k.Load(confmap.Provider(constructor.MergeConfig, k.Delim()), nil)
		if err != nil {
			return nil, err
		}
	}

	log.Trace().Strs("keys", k.Keys()).Msg("All merged config")

	unmarshaller := &KoanfUnmarshaller{
		koanf:        k,
		bytes:        bytes,
		enableEnv:    constructor.EnableEnv,
		configFormat: constructor.ConfigFormat,
	}

	return unmarshaller, nil
}

// KoanfUnmarshaller backed by viper + validator + go-defaults.
type KoanfUnmarshaller struct {
	sync.Mutex
	koanf        *koanf.Koanf
	configFormat ConfigFormat
	bytes        []byte
	enableEnv    bool
}

// Make sure KoanfUnmarshaller complies with the Unmarshaller interface.
var _ Unmarshaller = &KoanfUnmarshaller{}

// Get returns an interface value for the given key path in the config map.
func (u *KoanfUnmarshaller) Get(key string) interface{} {
	u.Lock()
	defer u.Unlock()

	return u.koanf.Get(key)
}

// IsSet checks if the given key is set in the config map.
func (u *KoanfUnmarshaller) IsSet(key string) bool {
	u.Lock()
	defer u.Unlock()

	return u.koanf.Exists(key)
}

// Unmarshal unmarshals given i using the underlying koanf.
func (u *KoanfUnmarshaller) Unmarshal(i interface{}) error {
	return u.UnmarshalKey("", i)
}

// UnmarshalKey binds the given interface to the given key path in the config map and
// unmarshals given i using the underlying koanf with additional configuration.
func (u *KoanfUnmarshaller) UnmarshalKey(keyPath string, i interface{}) error {
	u.Lock()
	defer u.Unlock()

	if u.enableEnv {
		u.bindEnvsKey(keyPath, i)
	}

	unmarshallerConf := koanf.UnmarshalConf{
		Tag: "json",
		DecoderConfig: &mapstructure.DecoderConfig{
			Squash: true,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				jsonOverrideHookFunc(true),
			),
			Result:   i,
			Metadata: nil,
		},
	}

	// do an initial decode to instantiate pointers
	var err error
	err = u.koanf.UnmarshalWithConf(keyPath, i, unmarshallerConf)
	if err != nil {
		return err
	}

	// Set defaults to fill missing values/zero values
	SetDefaults(i)

	// Decode again to override any zero values
	unmarshallerConf.DecoderConfig.DecodeHook = mapstructure.ComposeDecodeHookFunc(jsonOverrideHookFunc(false))

	err = u.koanf.UnmarshalWithConf(keyPath, i, unmarshallerConf)
	if err != nil {
		return err
	}
	// Validate
	return ValidateStruct(i)
}

// Reload reloads the config using the underlying koanf.
func (u *KoanfUnmarshaller) Reload(bytes []byte) error {
	return loadBytesProvider(u.koanf, bytes, u.configFormat)
}

func loadBytesProvider(k *koanf.Koanf, bytes []byte, configFormat ConfigFormat) error {
	var err error
	if bytes != nil {
		// TODO: allow parser config
		if configFormat == YAML {
			err = k.Load(rawbytes.Provider(bytes), yaml.Parser())
		} else {
			err = k.Load(rawbytes.Provider(bytes), koanfjson.Parser())
		}
	}
	return err
}

func (u *KoanfUnmarshaller) bindEnvsKey(keyPrefix string, in interface{}, prev ...string) {
	ifv := reflect.ValueOf(in)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	keyVals := map[string]interface{}{}

	for i := 0; i < ifv.NumField(); i++ {
		fv := ifv.Field(i)
		if fv.Kind() == reflect.Ptr {
			if fv.IsZero() {
				fv = reflect.New(fv.Type().Elem()).Elem()
			} else {
				fv = fv.Elem()
			}
		}
		t := ifv.Type().Field(i)

		// Embedded struct?
		if t.Anonymous {
			if fv.CanInterface() {
				u.bindEnvsKey(keyPrefix, fv.Interface(), prev...)
			}
			continue
		} else if !t.IsExported() {
			continue
		}

		tv, ok := t.Tag.Lookup("json")
		if ok && tv != "" {
			// scrub omitmepty and string options
			vals := strings.Split(tv, ",")
			tv = vals[0]
			if tv == "-" {
				continue
			}
		} else {
			tv = t.Name
		}

		var key string
		switch fv.Kind() {
		case reflect.Struct:
			// Check for duration types and treat them like native types
			if fv.Type().String() == "config.Duration" || fv.Type().String() == "*durationpb.Duration" {
				break
			}
			// Check for timestamp types and treat them like native types
			if fv.Type().String() == "config.Timestamp" || fv.Type().String() == "*timestamp.Timestamp" {
				break
			}
			if fv.CanInterface() {
				u.bindEnvsKey(keyPrefix, fv.Interface(), append(prev, tv)...)
			}
			continue
		case reflect.Map:
			iter := fv.MapRange()
			for iter.Next() {
				if key, ok = iter.Key().Interface().(string); ok {
					u.bindEnvsKey(keyPrefix, iter.Value().Interface(), append(prev, tv, key)...)
				}
			}
			continue
		}
		// Load env
		key = strings.Join(append(prev, tv), ".")

		if keyPrefix != "" {
			key = keyPrefix + "." + key
		}

		env := strings.ReplaceAll(key, ".", "_")
		env = strings.ReplaceAll(env, "-", "_")
		env = strings.ToUpper(env)

		env = EnvPrefix + env

		val, ok := os.LookupEnv(env)
		if ok {
			var v interface{}
			var err error
			switch fv.Kind() {
			case reflect.Slice:
				sliceType := fv.Type().Elem()
				reg := regexp.MustCompile(`^\[(.*)\]$`)
				matchs := reg.FindStringSubmatch(val)
				if len(matchs) != 2 {
					return
				}
				if matchs[1] == "" {
					v, err = nil, errors.New("empty slice provided in env var")
				} else {
					sliceValues := strings.Split(matchs[1], ",")
					switch sliceType.Kind() {
					case reflect.Bool:
						v, err = sliceconv.Atob(sliceValues)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						v, err = sliceconv.Atoi(sliceValues)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						v, err = sliceconv.Atoi(sliceValues)
					case reflect.Float32, reflect.Float64:
						v, err = sliceconv.Atof(sliceValues)
					case reflect.Struct:
						switch fv.Type().String() {
						case "config.Duration", "*durationpb.Duration", "config.Timestamp", "*timestamp.Timestamp":
							v, err = val, nil
						default:
							v, err = nil, errors.New("unable to decode struct from env var")
						}
					}
				}
			case reflect.Bool:
				v, err = strconv.ParseBool(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, err = strconv.ParseInt(val, 10, fv.Type().Bits())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, err = strconv.ParseUint(val, 10, fv.Type().Bits())
			case reflect.Float32, reflect.Float64:
				v, err = strconv.ParseFloat(val, fv.Type().Bits())
			case reflect.String:
				v, err = val, nil
			case reflect.Struct:
				switch fv.Type().String() {
				case "config.Duration", "*durationpb.Duration", "config.Timestamp", "*timestamp.Timestamp":
					v, err = val, nil
				default:
					v, err = nil, errors.New("unable to decode struct from env var")
				}
			default:
				v, err = nil, errors.New("unable to decode")
			}
			if err != nil {
				log.Error().Err(err).Msg("unable to decode env var")
				continue
			}
			if v != nil {
				if tv, ok := t.Tag.Lookup("public"); ok && tv == "true" {
					log.Info().Str("env", env).Str("key", key).Interface("v", v).Msg("reading env var")
				}
				keyVals[key] = v
			}
		}
	}
	// load into koanf
	if err := u.koanf.Load(confmap.Provider(keyVals, "."), nil); err != nil {
		log.Error().Err(err).Msg("unable to load env vars into koanf")
	}
}

func jsonOverrideHookFunc(replaceSlice bool) mapstructure.DecodeHookFunc {
	return func(f reflect.Value, t reflect.Value) (interface{}, error) {
		log.Trace().Interface("From Kind", f.Kind().String()).Interface("To Kind", t.Kind().String()).Msg("UNMARSHAL")

		// Raw map
		data := f.Interface()
		// Struct with default values
		result := t.Addr().Interface()

		// First, encode existing struct to json and decode it to a map
		var mapStruct map[string]interface{}
		var b []byte
		var err error

		log.Trace().Interface("data", data).Interface("data type", f.Type().String()).Interface("result type", t.Type().String()).Interface("result", result).Interface("mapStruct", mapStruct).Msg("BEFORE")
		// use ptr to marshal json in case MarshalJSON is defined on a pointer receiver (e.g. in case of json generated code for proto.Message)
		b, err = json.Marshal(result)

		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &mapStruct)
		if err != nil {
			return nil, err
		}

		log.Trace().Interface("data", data).Interface("data type", f.Type().String()).Interface("result type", t.Type().String()).Interface("result", result).Interface("mapStruct", mapStruct).Msg("LOAD EXISTING STRUCT")

		// Now we merge the raw map into defaults map

		merge(data.(map[string]interface{}), mapStruct, replaceSlice)

		log.Trace().Interface("data", data).Interface("data type", f.Type().String()).Interface("result type", t.Type().String()).Interface("result", result).Interface("mapStruct", mapStruct).Msg("AFTER MERGE")

		// Now we finally decode the map into our struct
		result = reflect.New(t.Type()).Interface()
		b, err = json.Marshal(mapStruct)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, result)
		if err != nil {
			return nil, err
		}
		log.Trace().Interface("data", data).Interface("data type", f.Type().String()).Interface("result type", t.Type().String()).Interface("result", result).Interface("mapStruct", mapStruct).Msg("FINAL")

		return result, nil
	}
}

func merge(a, b map[string]interface{}, replaceSlice bool) {
	for key, val := range a {
		// Does the key exist in the target map?
		// If no, add it and move on.
		bVal, ok := b[key]
		if !ok {
			log.Trace().Str("key", key).Msg("override")
			b[key] = val
			continue
		}

		if !replaceSlice {
			if sliceA, ok := val.([]interface{}); ok {
				if sliceB, ok := val.([]interface{}); ok {
					sliceMapA := []map[string]interface{}{}
					sliceMapB := []map[string]interface{}{}
					for _, v := range sliceA {
						if m, ok := v.(map[string]interface{}); ok {
							sliceMapA = append(sliceMapA, m)
						}
					}
					for _, v := range sliceB {
						if m, ok := v.(map[string]interface{}); ok {
							sliceMapB = append(sliceMapB, m)
						}
					}
					for i, sliceValA := range sliceMapA {
						if i < len(sliceB) {
							log.Trace().Str("key", key).Int("index", i).Msg("merging")
							// merge
							merge(sliceValA, sliceMapB[i], replaceSlice)
						} else {
							log.Trace().Str("key", key).Int("index", i).Msg("overriding")
							// append
							sliceMapB = append(sliceMapB, sliceValA)
						}
					}
					continue
				}
			}
		}

		// If the incoming val is not a map, do a direct merge.
		if _, ok := val.(map[string]interface{}); !ok {
			log.Trace().Str("key", key).Msg("override")
			b[key] = val
			continue
		}

		// The source key and target keys are both maps. Merge them.
		switch v := bVal.(type) {
		case map[string]interface{}:
			log.Trace().Str("key", key).Msg("merge")
			merge(val.(map[string]interface{}), v, replaceSlice)
		default:
			log.Trace().Str("key", key).Msg("override")
			b[key] = val
		}
	}
}

var json = jsoniter.Config{
	// encoding/json compat flags:
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,

	// Allowing alternate cases causes weird issues with our merge logic.
	// Disallow alternative casing.
	CaseSensitive: true,

	// Error on typos.
	// (Currently disabled, as we sometimes marshal the same config into two structs)
	// DisallowUnknownFields: true,
}.Froze()
