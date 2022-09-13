package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fluxninja/aperture/pkg/log"
)

// Most of the code ported from github.com/mcuadros/go-defaults

type fieldData struct {
	Parent   *fieldData
	TagValue string
	Value    reflect.Value
	Field    reflect.StructField
}

type fillerFunc func(field *fieldData)

// filler contains all the functions to fill any struct field with any type allowing to define function by Kind, Type of field name.
type filler struct {
	FuncByName map[string]fillerFunc
	FuncByType map[typeHash]fillerFunc
	FuncByKind map[reflect.Kind]fillerFunc
	Tag        string
}

// fill apply all the functions contained on Filler, setting all the possible
// values.
func (f *filler) fill(variable interface{}) {
	fields := f.getFields(variable)
	f.setDefaultValues(fields)
}

func (f *filler) getFields(variable interface{}) []*fieldData {
	valueObject := reflect.ValueOf(variable).Elem()

	return f.getFieldsFromValue(valueObject, nil)
}

func (f *filler) getFieldsFromValue(valueObject reflect.Value, parent *fieldData) []*fieldData {
	typeObject := valueObject.Type()

	count := valueObject.NumField()
	var results []*fieldData
	for i := 0; i < count; i++ {
		value := valueObject.Field(i)
		field := typeObject.Field(i)

		if value.CanSet() {
			results = append(results, &fieldData{
				Value:    value,
				Field:    field,
				TagValue: field.Tag.Get(f.Tag),
				Parent:   parent,
			})
		}
	}

	return results
}

func (f *filler) setDefaultValues(fields []*fieldData) {
	for _, field := range fields {
		if f.isEmpty(field) {
			f.setDefaultValue(field)
		}
	}
}

func (f *filler) isEmpty(field *fieldData) bool {
	switch field.Value.Kind() {
	case reflect.Bool:
		return !field.Value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Value.Int() == 0
	case reflect.Float32, reflect.Float64:
		return field.Value.Float() == .0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Value.Uint() == 0
	case reflect.Slice:
		switch field.Value.Type().Elem().Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Interface:
			// always assume the structs, ptrs, interfaces in the slice is empty and can be filled
			// the actually filling logic should take care of the rest
			return true
		default:
			return field.Value.Len() == 0
		}
	case reflect.String:
		return field.Value.String() == ""
	}
	return true
}

func (f *filler) setDefaultValue(field *fieldData) {
	getters := []func(field *fieldData) fillerFunc{
		f.getFunctionByName,
		f.getFunctionByType,
		f.getFunctionByKind,
	}

	for _, getter := range getters {
		filler := getter(field)
		if filler != nil {
			filler(field)
			return
		}
	}
}

func (f *filler) getFunctionByName(field *fieldData) fillerFunc {
	if f, ok := f.FuncByName[field.Field.Name]; ok {
		return f
	}

	return nil
}

func (f *filler) getFunctionByType(field *fieldData) fillerFunc {
	if f, ok := f.FuncByType[getTypeHash(field.Field.Type)]; ok {
		return f
	}

	return nil
}

func (f *filler) getFunctionByKind(field *fieldData) fillerFunc {
	if f, ok := f.FuncByKind[field.Field.Type.Kind()]; ok {
		return f
	}

	return nil
}

// typeHash is a string representing a reflect.Type following the next pattern:
// <package.name>.<type.name>.
type typeHash string

// getTypeHash returns the TypeHash for a given reflect.Type.
func getTypeHash(t reflect.Type) typeHash {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return typeHash(fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()))
}

// SetDefaults applies the default values to the struct object, the struct type must have
// the StructTag with name "default" and the directed value.
// Example usage:
//
//	type ExampleBasic struct {
//	    Foo bool   `default:"true"`
//	    Bar string `default:"33"`
//	    Qux int8
//	    Dur time.Duration `default:"2m3s"`
//	}
//
//	 foo := &ExampleBasic{}
//	 SetDefaults(foo)
func SetDefaults(variable interface{}) {
	kind := reflect.ValueOf(variable).Elem().Kind()
	if kind != reflect.Struct {
		// Panic early with clear message to avoid cryptic errors from inner functions.
		// Note: Perhaps we can lift this limitation.
		panic(fmt.Sprintf("SetDefaults can be used only on structs, got %v", kind))
	}
	getDefaultFiller().fill(variable)
}

var defaultFiller *filler = nil

func init() {
	defaultFiller = getDefaultFiller()
}

func getDefaultFiller() *filler {
	if defaultFiller == nil {
		defaultFiller = newDefaultFiller()
	}

	return defaultFiller
}

func newDefaultFiller() *filler {
	filler := &filler{}

	var setPtrDefaults func(val reflect.Value)
	setPtrDefaults = func(val reflect.Value) {
		// handle all non-nil pointers
		if val.IsNil() {
			return
		}
		if val.Elem().Kind() == reflect.Struct {
			SetDefaults(val.Interface())
		} else if val.Elem().Kind() == reflect.Ptr || val.Elem().Kind() == reflect.Interface {
			setPtrDefaults(val.Elem())
		}
	}

	funcs := make(map[reflect.Kind]fillerFunc, 0)
	funcs[reflect.Bool] = func(field *fieldData) {
		value, _ := strconv.ParseBool(field.TagValue)
		field.Value.SetBool(value)
	}

	funcs[reflect.Int] = func(field *fieldData) {
		value, _ := strconv.ParseInt(field.TagValue, 10, 64)
		field.Value.SetInt(value)
	}

	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]
	funcs[reflect.Int64] = func(field *fieldData) {
		if field.Field.Type == reflect.TypeOf(time.Second) {
			value, _ := time.ParseDuration(field.TagValue)
			field.Value.Set(reflect.ValueOf(value))
		} else {
			value, _ := strconv.ParseInt(field.TagValue, 10, 64)
			field.Value.SetInt(value)
		}
	}

	funcs[reflect.Float32] = func(field *fieldData) {
		value, _ := strconv.ParseFloat(field.TagValue, 64)
		field.Value.SetFloat(value)
	}

	funcs[reflect.Float64] = funcs[reflect.Float32]

	funcs[reflect.Uint] = func(field *fieldData) {
		value, _ := strconv.ParseUint(field.TagValue, 10, 64)
		field.Value.SetUint(value)
	}

	funcs[reflect.Uint8] = funcs[reflect.Uint]
	funcs[reflect.Uint16] = funcs[reflect.Uint]
	funcs[reflect.Uint32] = funcs[reflect.Uint]
	funcs[reflect.Uint64] = funcs[reflect.Uint]

	funcs[reflect.String] = func(field *fieldData) {
		tagValue := parseDateTimeString(field.TagValue)
		field.Value.SetString(tagValue)
	}

	funcs[reflect.Struct] = func(field *fieldData) {
		fields := filler.getFieldsFromValue(field.Value, nil)
		filler.setDefaultValues(fields)
	}

	funcs[reflect.Ptr] = func(field *fieldData) {
		val := field.Value
		setPtrDefaults(val)
	}

	funcs[reflect.Interface] = func(field *fieldData) {
		val := field.Value
		setPtrDefaults(val)
	}

	types := make(map[typeHash]fillerFunc, 5)
	types["time.Duration"] = func(field *fieldData) {
		d, _ := time.ParseDuration(field.TagValue)
		field.Value.Set(reflect.ValueOf(d))
	}

	types[getTypeHash(reflect.TypeOf(Duration{}))] = func(field *fieldData) {
		if value, ok := field.Value.Interface().(Duration); ok {
			if value.AsDuration() != 0 {
				return
			}
			if field.TagValue != "" {
				durationJSON, _ := json.Marshal(field.TagValue)
				dur := Duration{}
				dur.duration = durationpb.New(0)
				err := json.Unmarshal(durationJSON, &dur)
				if err != nil {
					log.Error().Err(err).Msg("Unable to unmarshal default duration")
				}
				field.Value.Set(reflect.ValueOf(dur))
			}
		}
	}

	types[getTypeHash(reflect.TypeOf(durationpb.New(0)))] = func(field *fieldData) {
		iface := field.Value.Interface()
		if value, ok := iface.(*durationpb.Duration); ok {
			if value.AsDuration() != 0 {
				return
			}
			if field.TagValue != "" {
				durationJSON, _ := json.Marshal(field.TagValue)
				dur := durationpb.New(0)
				err := protojson.Unmarshal(durationJSON, dur)
				if err != nil {
					log.Error().Err(err).Msg("Unable to unmarshal duration default")
				}
				field.Value.Set(reflect.ValueOf(dur))
			}
		}
	}

	types[getTypeHash(reflect.TypeOf(Timestamp{}))] = func(field *fieldData) {
		if value, ok := field.Value.Interface().(Timestamp); ok {
			if field.TagValue != "" {
				nullTime := time.Time{}
				if value.Timestamp.AsTime() != nullTime {
					return
				}
				timestampJSON, _ := json.Marshal(field.TagValue)
				t := &Timestamp{
					Timestamp: timestamppb.Now(),
				}
				err := json.Unmarshal(timestampJSON, t)
				if err != nil {
					log.Error().Err(err).Msg("Unable to unmarshal default config.Timestamp")
				}
				field.Value.Set(reflect.ValueOf(t))
			}
		}
	}

	types[getTypeHash(reflect.TypeOf(timestamppb.Now()))] = func(field *fieldData) {
		iface := field.Value.Interface()
		if value, ok := iface.(*timestamppb.Timestamp); ok {
			nullTime := time.Time{}
			if value.AsTime() != nullTime {
				return
			}
			if field.TagValue != "" {
				timestampJSON, _ := json.Marshal(field.TagValue)
				t := &Timestamp{
					Timestamp: timestamppb.Now(),
				}
				err := json.Unmarshal(timestampJSON, t)
				if err != nil {
					log.Error().Err(err).Msg("Unable to unmarshal default timestamppb")
				}
				field.Value.Set(reflect.ValueOf(t))
			}
		}
	}

	funcs[reflect.Slice] = func(field *fieldData) {
		k := field.Value.Type().Elem().Kind()
		switch k {
		case reflect.Uint8:
			if field.Value.Bytes() != nil {
				return
			}
			field.Value.SetBytes([]byte(field.TagValue))
		case reflect.Struct:
			count := field.Value.Len()
			for i := 0; i < count; i++ {
				fields := filler.getFieldsFromValue(field.Value.Index(i), nil)
				filler.setDefaultValues(fields)
			}
		case reflect.Ptr, reflect.Interface:
			count := field.Value.Len()
			for i := 0; i < count; i++ {
				val := field.Value.Index(i)
				setPtrDefaults(val)
			}
			// FIXME: also descend into slice of slices
			// FIXME: also descend into slice of maps
		default:
			reg := regexp.MustCompile(`^\[(.*)\]$`)
			matches := reg.FindStringSubmatch(field.TagValue)
			if len(matches) != 2 {
				return
			}
			if matches[1] == "" {
				field.Value.Set(reflect.MakeSlice(field.Value.Type(), 0, 0))
			} else {
				defaultValue := strings.Split(matches[1], ",")
				result := reflect.MakeSlice(field.Value.Type(), len(defaultValue), len(defaultValue))
				for i := 0; i < len(defaultValue); i++ {
					itemValue := result.Index(i)
					item := &fieldData{
						Value:    itemValue,
						Field:    reflect.StructField{},
						TagValue: defaultValue[i],
						Parent:   nil,
					}
					funcs[k](item)
				}
				field.Value.Set(result)
			}
		}
	}
	funcs[reflect.Map] = func(field *fieldData) {
		switch field.Field.Type.Elem().Kind() {
		case reflect.Struct:
			keys := field.Value.MapKeys()
			for _, key := range keys {
				// We need to pull the value from map, copy to a settable
				// location, fill in defaults, and then put it back into map.
				value := field.Value.MapIndex(key)
				tmp := reflect.New(value.Type()).Elem()
				tmp.Set(value)
				fields := filler.getFieldsFromValue(tmp, nil)
				filler.setDefaultValues(fields)
				field.Value.SetMapIndex(key, tmp)
			}
		case reflect.Ptr, reflect.Interface:
			iter := field.Value.MapRange()
			for iter.Next() {
				setPtrDefaults(iter.Value())
			}
		case reflect.Slice:
			// FIXME: also descend into map of slices
		case reflect.Map:
			// FIXME: also descend into map of maps
		}
		// FIXME: handle actual default Tags on maps?
	}
	filler.FuncByKind = funcs
	filler.FuncByType = types
	filler.Tag = "default"
	return filler
}

func parseDateTimeString(data string) string {
	pattern := regexp.MustCompile(`\{\{(\w+\:(?:-|)\d*,(?:-|)\d*,(?:-|)\d*)\}\}`)
	matches := pattern.FindAllStringSubmatch(data, -1) // matches is [][]string
	for _, match := range matches {

		tags := strings.Split(match[1], ":")
		if len(tags) == 2 {

			valueStrings := strings.Split(tags[1], ",")
			if len(valueStrings) == 3 {
				var values [3]int
				for key, valueString := range valueStrings {
					num, _ := strconv.ParseInt(valueString, 10, 64)
					values[key] = int(num)
				}

				switch tags[0] {

				case "date":
					str := time.Now().AddDate(values[0], values[1], values[2]).Format("2006-01-02")
					data = strings.ReplaceAll(data, match[0], str)
				case "time":
					str := time.Now().Add((time.Duration(values[0]) * time.Hour) +
						(time.Duration(values[1]) * time.Minute) +
						(time.Duration(values[2]) * time.Second)).Format("15:04:05")
					data = strings.ReplaceAll(data, match[0], str)
				}
			}
		}

	}
	return data
}
