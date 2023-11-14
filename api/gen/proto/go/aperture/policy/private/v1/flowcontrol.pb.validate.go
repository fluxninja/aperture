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

	checkv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
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

	_ = checkv1.StatusCode(0)
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

	// no validation rules for LoadSchedulerComponentId

	for idx, item := range m.GetSelectors() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, LoadActuatorValidationError{
						field:  fmt.Sprintf("Selectors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, LoadActuatorValidationError{
						field:  fmt.Sprintf("Selectors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return LoadActuatorValidationError{
					field:  fmt.Sprintf("Selectors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

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

// Validate checks the field values on RateLimiter with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *RateLimiter) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in RateLimiterMultiError, or
// nil if none found.
func (m *RateLimiter) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetInPorts()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "InPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "InPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInPorts()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiterValidationError{
				field:  "InPorts",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetParameters()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "Parameters",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "Parameters",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetParameters()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiterValidationError{
				field:  "Parameters",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetSelectors() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RateLimiterValidationError{
						field:  fmt.Sprintf("Selectors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RateLimiterValidationError{
						field:  fmt.Sprintf("Selectors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RateLimiterValidationError{
					field:  fmt.Sprintf("Selectors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetRequestParameters()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "RequestParameters",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "RequestParameters",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRequestParameters()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiterValidationError{
				field:  "RequestParameters",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ParentComponentId

	if len(errors) > 0 {
		return RateLimiterMultiError(errors)
	}

	return nil
}

// RateLimiterMultiError is an error wrapping multiple validation errors
// returned by RateLimiter.ValidateAll() if the designated constraints aren't met.
type RateLimiterMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiterMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiterMultiError) AllErrors() []error { return m }

// RateLimiterValidationError is the validation error returned by
// RateLimiter.Validate if the designated constraints aren't met.
type RateLimiterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiterValidationError) ErrorName() string { return "RateLimiterValidationError" }

// Error satisfies the builtin error interface
func (e RateLimiterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiterValidationError{}

// Validate checks the field values on QuotaScheduler with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *QuotaScheduler) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QuotaScheduler with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in QuotaSchedulerMultiError,
// or nil if none found.
func (m *QuotaScheduler) ValidateAll() error {
	return m.validate(true)
}

func (m *QuotaScheduler) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetInPorts()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, QuotaSchedulerValidationError{
					field:  "InPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, QuotaSchedulerValidationError{
					field:  "InPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInPorts()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return QuotaSchedulerValidationError{
				field:  "InPorts",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetSelectors() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, QuotaSchedulerValidationError{
						field:  fmt.Sprintf("Selectors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, QuotaSchedulerValidationError{
						field:  fmt.Sprintf("Selectors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return QuotaSchedulerValidationError{
					field:  fmt.Sprintf("Selectors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetRateLimiter()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, QuotaSchedulerValidationError{
					field:  "RateLimiter",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, QuotaSchedulerValidationError{
					field:  "RateLimiter",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRateLimiter()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return QuotaSchedulerValidationError{
				field:  "RateLimiter",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetScheduler()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, QuotaSchedulerValidationError{
					field:  "Scheduler",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, QuotaSchedulerValidationError{
					field:  "Scheduler",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetScheduler()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return QuotaSchedulerValidationError{
				field:  "Scheduler",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ParentComponentId

	if len(errors) > 0 {
		return QuotaSchedulerMultiError(errors)
	}

	return nil
}

// QuotaSchedulerMultiError is an error wrapping multiple validation errors
// returned by QuotaScheduler.ValidateAll() if the designated constraints
// aren't met.
type QuotaSchedulerMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QuotaSchedulerMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QuotaSchedulerMultiError) AllErrors() []error { return m }

// QuotaSchedulerValidationError is the validation error returned by
// QuotaScheduler.Validate if the designated constraints aren't met.
type QuotaSchedulerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QuotaSchedulerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QuotaSchedulerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QuotaSchedulerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QuotaSchedulerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QuotaSchedulerValidationError) ErrorName() string { return "QuotaSchedulerValidationError" }

// Error satisfies the builtin error interface
func (e QuotaSchedulerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQuotaScheduler.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QuotaSchedulerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QuotaSchedulerValidationError{}

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

// Validate checks the field values on RateLimiter_Parameters with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RateLimiter_Parameters) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter_Parameters with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimiter_ParametersMultiError, or nil if none found.
func (m *RateLimiter_Parameters) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter_Parameters) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for LimitByLabelKey

	if all {
		switch v := interface{}(m.GetInterval()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_ParametersValidationError{
					field:  "Interval",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_ParametersValidationError{
					field:  "Interval",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_ParametersValidationError{
				field:  "Interval",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ContinuousFill

	if all {
		switch v := interface{}(m.GetMaxIdleTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_ParametersValidationError{
					field:  "MaxIdleTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_ParametersValidationError{
					field:  "MaxIdleTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMaxIdleTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_ParametersValidationError{
				field:  "MaxIdleTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetLazySync()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_ParametersValidationError{
					field:  "LazySync",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_ParametersValidationError{
					field:  "LazySync",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLazySync()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_ParametersValidationError{
				field:  "LazySync",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for DelayInitialFill

	if len(errors) > 0 {
		return RateLimiter_ParametersMultiError(errors)
	}

	return nil
}

// RateLimiter_ParametersMultiError is an error wrapping multiple validation
// errors returned by RateLimiter_Parameters.ValidateAll() if the designated
// constraints aren't met.
type RateLimiter_ParametersMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiter_ParametersMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiter_ParametersMultiError) AllErrors() []error { return m }

// RateLimiter_ParametersValidationError is the validation error returned by
// RateLimiter_Parameters.Validate if the designated constraints aren't met.
type RateLimiter_ParametersValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiter_ParametersValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiter_ParametersValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiter_ParametersValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiter_ParametersValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiter_ParametersValidationError) ErrorName() string {
	return "RateLimiter_ParametersValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimiter_ParametersValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter_Parameters.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiter_ParametersValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiter_ParametersValidationError{}

// Validate checks the field values on RateLimiter_RequestParameters with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RateLimiter_RequestParameters) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter_RequestParameters with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// RateLimiter_RequestParametersMultiError, or nil if none found.
func (m *RateLimiter_RequestParameters) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter_RequestParameters) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TokensLabelKey

	// no validation rules for DeniedResponseStatusCode

	if len(errors) > 0 {
		return RateLimiter_RequestParametersMultiError(errors)
	}

	return nil
}

// RateLimiter_RequestParametersMultiError is an error wrapping multiple
// validation errors returned by RateLimiter_RequestParameters.ValidateAll()
// if the designated constraints aren't met.
type RateLimiter_RequestParametersMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiter_RequestParametersMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiter_RequestParametersMultiError) AllErrors() []error { return m }

// RateLimiter_RequestParametersValidationError is the validation error
// returned by RateLimiter_RequestParameters.Validate if the designated
// constraints aren't met.
type RateLimiter_RequestParametersValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiter_RequestParametersValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiter_RequestParametersValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiter_RequestParametersValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiter_RequestParametersValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiter_RequestParametersValidationError) ErrorName() string {
	return "RateLimiter_RequestParametersValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimiter_RequestParametersValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter_RequestParameters.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiter_RequestParametersValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiter_RequestParametersValidationError{}

// Validate checks the field values on RateLimiter_Ins with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RateLimiter_Ins) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter_Ins with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimiter_InsMultiError, or nil if none found.
func (m *RateLimiter_Ins) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter_Ins) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetBucketCapacity()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_InsValidationError{
					field:  "BucketCapacity",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_InsValidationError{
					field:  "BucketCapacity",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetBucketCapacity()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_InsValidationError{
				field:  "BucketCapacity",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetFillAmount()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_InsValidationError{
					field:  "FillAmount",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_InsValidationError{
					field:  "FillAmount",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetFillAmount()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_InsValidationError{
				field:  "FillAmount",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetPassThrough()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_InsValidationError{
					field:  "PassThrough",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_InsValidationError{
					field:  "PassThrough",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPassThrough()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_InsValidationError{
				field:  "PassThrough",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return RateLimiter_InsMultiError(errors)
	}

	return nil
}

// RateLimiter_InsMultiError is an error wrapping multiple validation errors
// returned by RateLimiter_Ins.ValidateAll() if the designated constraints
// aren't met.
type RateLimiter_InsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiter_InsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiter_InsMultiError) AllErrors() []error { return m }

// RateLimiter_InsValidationError is the validation error returned by
// RateLimiter_Ins.Validate if the designated constraints aren't met.
type RateLimiter_InsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiter_InsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiter_InsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiter_InsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiter_InsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiter_InsValidationError) ErrorName() string { return "RateLimiter_InsValidationError" }

// Error satisfies the builtin error interface
func (e RateLimiter_InsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter_Ins.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiter_InsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiter_InsValidationError{}

// Validate checks the field values on RateLimiter_Parameters_LazySync with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RateLimiter_Parameters_LazySync) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter_Parameters_LazySync with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// RateLimiter_Parameters_LazySyncMultiError, or nil if none found.
func (m *RateLimiter_Parameters_LazySync) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter_Parameters_LazySync) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Enabled

	// no validation rules for NumSync

	if len(errors) > 0 {
		return RateLimiter_Parameters_LazySyncMultiError(errors)
	}

	return nil
}

// RateLimiter_Parameters_LazySyncMultiError is an error wrapping multiple
// validation errors returned by RateLimiter_Parameters_LazySync.ValidateAll()
// if the designated constraints aren't met.
type RateLimiter_Parameters_LazySyncMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiter_Parameters_LazySyncMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiter_Parameters_LazySyncMultiError) AllErrors() []error { return m }

// RateLimiter_Parameters_LazySyncValidationError is the validation error
// returned by RateLimiter_Parameters_LazySync.Validate if the designated
// constraints aren't met.
type RateLimiter_Parameters_LazySyncValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiter_Parameters_LazySyncValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiter_Parameters_LazySyncValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiter_Parameters_LazySyncValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiter_Parameters_LazySyncValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiter_Parameters_LazySyncValidationError) ErrorName() string {
	return "RateLimiter_Parameters_LazySyncValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimiter_Parameters_LazySyncValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter_Parameters_LazySync.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiter_Parameters_LazySyncValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiter_Parameters_LazySyncValidationError{}
