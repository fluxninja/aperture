// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: aperture/policy/language/v1/telemetry_collector.proto

package languagev1

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

// Validate checks the field values on TelemetryCollector with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TelemetryCollector) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TelemetryCollector with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TelemetryCollectorMultiError, or nil if none found.
func (m *TelemetryCollector) ValidateAll() error {
	return m.validate(true)
}

func (m *TelemetryCollector) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AgentGroup

	{
		sorted_keys := make([]string, len(m.GetInfraMeters()))
		i := 0
		for key := range m.GetInfraMeters() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetInfraMeters()[key]
			_ = val

			// no validation rules for InfraMeters[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, TelemetryCollectorValidationError{
							field:  fmt.Sprintf("InfraMeters[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, TelemetryCollectorValidationError{
							field:  fmt.Sprintf("InfraMeters[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return TelemetryCollectorValidationError{
						field:  fmt.Sprintf("InfraMeters[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if len(errors) > 0 {
		return TelemetryCollectorMultiError(errors)
	}

	return nil
}

// TelemetryCollectorMultiError is an error wrapping multiple validation errors
// returned by TelemetryCollector.ValidateAll() if the designated constraints
// aren't met.
type TelemetryCollectorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TelemetryCollectorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TelemetryCollectorMultiError) AllErrors() []error { return m }

// TelemetryCollectorValidationError is the validation error returned by
// TelemetryCollector.Validate if the designated constraints aren't met.
type TelemetryCollectorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TelemetryCollectorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TelemetryCollectorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TelemetryCollectorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TelemetryCollectorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TelemetryCollectorValidationError) ErrorName() string {
	return "TelemetryCollectorValidationError"
}

// Error satisfies the builtin error interface
func (e TelemetryCollectorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTelemetryCollector.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TelemetryCollectorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TelemetryCollectorValidationError{}

// Validate checks the field values on InfraMeter with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *InfraMeter) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InfraMeter with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in InfraMeterMultiError, or
// nil if none found.
func (m *InfraMeter) ValidateAll() error {
	return m.validate(true)
}

func (m *InfraMeter) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	{
		sorted_keys := make([]string, len(m.GetReceivers()))
		i := 0
		for key := range m.GetReceivers() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetReceivers()[key]
			_ = val

			// no validation rules for Receivers[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, InfraMeterValidationError{
							field:  fmt.Sprintf("Receivers[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, InfraMeterValidationError{
							field:  fmt.Sprintf("Receivers[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return InfraMeterValidationError{
						field:  fmt.Sprintf("Receivers[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	{
		sorted_keys := make([]string, len(m.GetProcessors()))
		i := 0
		for key := range m.GetProcessors() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetProcessors()[key]
			_ = val

			// no validation rules for Processors[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, InfraMeterValidationError{
							field:  fmt.Sprintf("Processors[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, InfraMeterValidationError{
							field:  fmt.Sprintf("Processors[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return InfraMeterValidationError{
						field:  fmt.Sprintf("Processors[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if all {
		switch v := interface{}(m.GetPipeline()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, InfraMeterValidationError{
					field:  "Pipeline",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, InfraMeterValidationError{
					field:  "Pipeline",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPipeline()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InfraMeterValidationError{
				field:  "Pipeline",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for PerAgentGroup

	// no validation rules for AgentGroup

	if len(errors) > 0 {
		return InfraMeterMultiError(errors)
	}

	return nil
}

// InfraMeterMultiError is an error wrapping multiple validation errors
// returned by InfraMeter.ValidateAll() if the designated constraints aren't met.
type InfraMeterMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InfraMeterMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InfraMeterMultiError) AllErrors() []error { return m }

// InfraMeterValidationError is the validation error returned by
// InfraMeter.Validate if the designated constraints aren't met.
type InfraMeterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InfraMeterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InfraMeterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InfraMeterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InfraMeterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InfraMeterValidationError) ErrorName() string { return "InfraMeterValidationError" }

// Error satisfies the builtin error interface
func (e InfraMeterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInfraMeter.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InfraMeterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InfraMeterValidationError{}

// Validate checks the field values on InfraMeter_MetricsPipeline with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InfraMeter_MetricsPipeline) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InfraMeter_MetricsPipeline with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InfraMeter_MetricsPipelineMultiError, or nil if none found.
func (m *InfraMeter_MetricsPipeline) ValidateAll() error {
	return m.validate(true)
}

func (m *InfraMeter_MetricsPipeline) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return InfraMeter_MetricsPipelineMultiError(errors)
	}

	return nil
}

// InfraMeter_MetricsPipelineMultiError is an error wrapping multiple
// validation errors returned by InfraMeter_MetricsPipeline.ValidateAll() if
// the designated constraints aren't met.
type InfraMeter_MetricsPipelineMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InfraMeter_MetricsPipelineMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InfraMeter_MetricsPipelineMultiError) AllErrors() []error { return m }

// InfraMeter_MetricsPipelineValidationError is the validation error returned
// by InfraMeter_MetricsPipeline.Validate if the designated constraints aren't met.
type InfraMeter_MetricsPipelineValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InfraMeter_MetricsPipelineValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InfraMeter_MetricsPipelineValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InfraMeter_MetricsPipelineValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InfraMeter_MetricsPipelineValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InfraMeter_MetricsPipelineValidationError) ErrorName() string {
	return "InfraMeter_MetricsPipelineValidationError"
}

// Error satisfies the builtin error interface
func (e InfraMeter_MetricsPipelineValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInfraMeter_MetricsPipeline.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InfraMeter_MetricsPipelineValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InfraMeter_MetricsPipelineValidationError{}
