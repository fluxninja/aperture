// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: aperture/policy/sync/v1/load_scheduler.proto

package syncv1

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

// Validate checks the field values on LoadSchedulerWrapper with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *LoadSchedulerWrapper) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoadSchedulerWrapper with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LoadSchedulerWrapperMultiError, or nil if none found.
func (m *LoadSchedulerWrapper) ValidateAll() error {
	return m.validate(true)
}

func (m *LoadSchedulerWrapper) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCommonAttributes()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadSchedulerWrapperValidationError{
					field:  "CommonAttributes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadSchedulerWrapperValidationError{
					field:  "CommonAttributes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommonAttributes()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadSchedulerWrapperValidationError{
				field:  "CommonAttributes",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetLoadScheduler()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadSchedulerWrapperValidationError{
					field:  "LoadScheduler",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadSchedulerWrapperValidationError{
					field:  "LoadScheduler",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLoadScheduler()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadSchedulerWrapperValidationError{
				field:  "LoadScheduler",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return LoadSchedulerWrapperMultiError(errors)
	}

	return nil
}

// LoadSchedulerWrapperMultiError is an error wrapping multiple validation
// errors returned by LoadSchedulerWrapper.ValidateAll() if the designated
// constraints aren't met.
type LoadSchedulerWrapperMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoadSchedulerWrapperMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoadSchedulerWrapperMultiError) AllErrors() []error { return m }

// LoadSchedulerWrapperValidationError is the validation error returned by
// LoadSchedulerWrapper.Validate if the designated constraints aren't met.
type LoadSchedulerWrapperValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoadSchedulerWrapperValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoadSchedulerWrapperValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoadSchedulerWrapperValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoadSchedulerWrapperValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoadSchedulerWrapperValidationError) ErrorName() string {
	return "LoadSchedulerWrapperValidationError"
}

// Error satisfies the builtin error interface
func (e LoadSchedulerWrapperValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoadSchedulerWrapper.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoadSchedulerWrapperValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoadSchedulerWrapperValidationError{}

// Validate checks the field values on LoadDecisionWrapper with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *LoadDecisionWrapper) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoadDecisionWrapper with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LoadDecisionWrapperMultiError, or nil if none found.
func (m *LoadDecisionWrapper) ValidateAll() error {
	return m.validate(true)
}

func (m *LoadDecisionWrapper) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCommonAttributes()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadDecisionWrapperValidationError{
					field:  "CommonAttributes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadDecisionWrapperValidationError{
					field:  "CommonAttributes",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommonAttributes()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadDecisionWrapperValidationError{
				field:  "CommonAttributes",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetLoadDecision()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadDecisionWrapperValidationError{
					field:  "LoadDecision",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadDecisionWrapperValidationError{
					field:  "LoadDecision",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLoadDecision()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadDecisionWrapperValidationError{
				field:  "LoadDecision",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return LoadDecisionWrapperMultiError(errors)
	}

	return nil
}

// LoadDecisionWrapperMultiError is an error wrapping multiple validation
// errors returned by LoadDecisionWrapper.ValidateAll() if the designated
// constraints aren't met.
type LoadDecisionWrapperMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoadDecisionWrapperMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoadDecisionWrapperMultiError) AllErrors() []error { return m }

// LoadDecisionWrapperValidationError is the validation error returned by
// LoadDecisionWrapper.Validate if the designated constraints aren't met.
type LoadDecisionWrapperValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoadDecisionWrapperValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoadDecisionWrapperValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoadDecisionWrapperValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoadDecisionWrapperValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoadDecisionWrapperValidationError) ErrorName() string {
	return "LoadDecisionWrapperValidationError"
}

// Error satisfies the builtin error interface
func (e LoadDecisionWrapperValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoadDecisionWrapper.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoadDecisionWrapperValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoadDecisionWrapperValidationError{}

// Validate checks the field values on LoadDecision with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoadDecision) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoadDecision with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoadDecisionMultiError, or
// nil if none found.
func (m *LoadDecision) ValidateAll() error {
	return m.validate(true)
}

func (m *LoadDecision) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTickInfo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadDecisionValidationError{
					field:  "TickInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadDecisionValidationError{
					field:  "TickInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTickInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadDecisionValidationError{
				field:  "TickInfo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for LoadMultiplier

	// no validation rules for PassThrough

	// no validation rules for TokensByWorkloadIndex

	if all {
		switch v := interface{}(m.GetValidUntil()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LoadDecisionValidationError{
					field:  "ValidUntil",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LoadDecisionValidationError{
					field:  "ValidUntil",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetValidUntil()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LoadDecisionValidationError{
				field:  "ValidUntil",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return LoadDecisionMultiError(errors)
	}

	return nil
}

// LoadDecisionMultiError is an error wrapping multiple validation errors
// returned by LoadDecision.ValidateAll() if the designated constraints aren't met.
type LoadDecisionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoadDecisionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoadDecisionMultiError) AllErrors() []error { return m }

// LoadDecisionValidationError is the validation error returned by
// LoadDecision.Validate if the designated constraints aren't met.
type LoadDecisionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoadDecisionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoadDecisionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoadDecisionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoadDecisionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoadDecisionValidationError) ErrorName() string { return "LoadDecisionValidationError" }

// Error satisfies the builtin error interface
func (e LoadDecisionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoadDecision.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoadDecisionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoadDecisionValidationError{}
