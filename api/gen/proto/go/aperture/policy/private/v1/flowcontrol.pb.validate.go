// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: aperture/policy/private/v1/flowcontrol.proto

package privatev1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on LoadActuator with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoadActuator) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoadActuator with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoadActuatorMultiError, or
// nil if none found.
func (m *LoadActuator) ValidateAll() error {
	return m.validate(true)
}

func (m *LoadActuator) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetInPorts()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadActuatorValidationError{
					field:  "InPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadActuatorValidationError{
					field:  "InPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInPorts()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadActuatorValidationError{
				field:  "InPorts",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ComponentId

	if all {
		switch v := interface{}(m.GetDefaultConfig()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadActuatorValidationError{
					field:  "DefaultConfig",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadActuatorValidationError{
					field:  "DefaultConfig",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDefaultConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadActuatorValidationError{
				field:  "DefaultConfig",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for DynamicConfigKey

	// no validation rules for WorkloadLatencyBasedTokens

	if len(errors) > 0 {
		return LoadActuatorMultiError(errors)
	}

	return nil
}

// LoadActuatorMultiError is an error wrapping multiple validation errors
// returned by LoadActuator.ValidateAll() if the designated constraints aren't met.
type LoadActuatorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoadActuatorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoadActuatorMultiError) AllErrors() []error { return m }

// LoadActuatorValidationError is the validation error returned by
// LoadActuator.Validate if the designated constraints aren't met.
type LoadActuatorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoadActuatorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoadActuatorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoadActuatorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoadActuatorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoadActuatorValidationError) ErrorName() string { return "LoadActuatorValidationError" }

// Error satisfies the builtin error interface
func (e LoadActuatorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoadActuator.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoadActuatorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoadActuatorValidationError{}

// Validate checks the field values on LoadActuator_Ins with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *LoadActuator_Ins) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoadActuator_Ins with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LoadActuator_InsMultiError, or nil if none found.
func (m *LoadActuator_Ins) ValidateAll() error {
	return m.validate(true)
}

func (m *LoadActuator_Ins) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetLoadMultiplier()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadActuator_InsValidationError{
					field:  "LoadMultiplier",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadActuator_InsValidationError{
					field:  "LoadMultiplier",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLoadMultiplier()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadActuator_InsValidationError{
				field:  "LoadMultiplier",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return LoadActuator_InsMultiError(errors)
	}

	return nil
}

// LoadActuator_InsMultiError is an error wrapping multiple validation errors
// returned by LoadActuator_Ins.ValidateAll() if the designated constraints
// aren't met.
type LoadActuator_InsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoadActuator_InsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoadActuator_InsMultiError) AllErrors() []error { return m }

// LoadActuator_InsValidationError is the validation error returned by
// LoadActuator_Ins.Validate if the designated constraints aren't met.
type LoadActuator_InsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoadActuator_InsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoadActuator_InsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoadActuator_InsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoadActuator_InsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoadActuator_InsValidationError) ErrorName() string { return "LoadActuator_InsValidationError" }

// Error satisfies the builtin error interface
func (e LoadActuator_InsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoadActuator_Ins.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoadActuator_InsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoadActuator_InsValidationError{}
