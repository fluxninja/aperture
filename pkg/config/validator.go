package config

import (
	"fmt"
	"reflect"

	validator "github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

var globalValidate = getValidate()

func init() {
	if err := RegisterDeprecatedValidator(true); err != nil {
		log.Panic().Err(err).Msg("Failed to register deprecated validator")
	}
}

func getValidate() *validator.Validate {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(durationCustomTypeFunc, Duration{})
	validate.RegisterCustomTypeFunc(durationpbCustomTypeFunc, durationpb.Duration{})
	validate.RegisterCustomTypeFunc(timestampCustomTypeFunc, Time{})
	validate.RegisterCustomTypeFunc(timestamppbCustomTypeFunc, timestamppb.Timestamp{})
	return validate
}

// RegisterDeprecatedValidator a function to validate "deprecated". If deprecated constraint exists then fail the validation.
func RegisterDeprecatedValidator(warnOnly bool) error {
	return globalValidate.RegisterValidation("deprecated", func(fl validator.FieldLevel) bool {
		// check whether the field is set or not
		if fl.Field().IsZero() {
			return true
		}
		if warnOnly {
			log.Warn().Msgf("Field %s is deprecated", fl.Field().String())
		} else {
			log.Error().Msgf("Field %s is deprecated, allow deprecated fields and retry", fl.Field().String())
		}
		return warnOnly
	})
}

// ValidateStruct takes interface value and validates its fields of a struct.
func ValidateStruct(rawVal interface{}) error {
	// validate configuration
	err := globalValidate.Struct(rawVal)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Panic().Err(err).Msg("InvalidValidationError!")
		} else if _, ok := err.(validator.ValidationErrors); ok {
			for _, err := range err.(validator.ValidationErrors) {
				errorStr := fmt.Sprintf("ValidationError<"+
					"Namespace: %s"+
					"| Field: %s"+
					"| StructNamespace: %s"+
					"| StructField: %s"+
					"| Tag: %s"+
					"| ActualTag: %s"+
					"| Kind: %s"+
					"| Type: %s"+
					"| Value: %s"+
					"| Param: %s"+
					">",
					err.Namespace(),
					err.Field(),
					err.StructNamespace(),
					err.StructField(),
					err.Tag(),
					err.ActualTag(),
					err.Kind(),
					err.Type(),
					err.Value(),
					err.Param())
				log.Error().Err(err).Msg(errorStr)
			}
		}
	}
	return err
}

func durationCustomTypeFunc(field reflect.Value) interface{} {
	if value, ok := field.Interface().(Duration); ok {
		return value.AsDuration()
	}
	return nil
}

func durationpbCustomTypeFunc(field reflect.Value) interface{} {
	iface := field.Interface()
	switch iface.(type) {
	case durationpb.Duration:
		ptr := field.Addr().Interface()
		return ptr.(*durationpb.Duration).AsDuration()
	}
	return nil
}

func timestampCustomTypeFunc(field reflect.Value) interface{} {
	if value, ok := field.Interface().(Time); ok {
		return value.timestamp.AsTime()
	}
	return nil
}

func timestamppbCustomTypeFunc(field reflect.Value) interface{} {
	iface := field.Interface()
	switch iface.(type) {
	case timestamppb.Timestamp:
		ptr := field.Addr().Interface()
		return ptr.(*timestamppb.Timestamp).AsTime()
	}
	return nil
}
