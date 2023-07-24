// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: aperture/policy/language/v1/query.proto

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

// Validate checks the field values on Query with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Query) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Query with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in QueryMultiError, or nil if none found.
func (m *Query) ValidateAll() error {
	return m.validate(true)
}

func (m *Query) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	switch v := m.Component.(type) {
	case *Query_Promql:
		if v == nil {
			err := QueryValidationError{
				field:  "Component",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetPromql()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, QueryValidationError{
						field:  "Promql",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, QueryValidationError{
						field:  "Promql",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetPromql()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return QueryValidationError{
					field:  "Promql",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return QueryMultiError(errors)
	}

	return nil
}

// QueryMultiError is an error wrapping multiple validation errors returned by
// Query.ValidateAll() if the designated constraints aren't met.
type QueryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueryMultiError) AllErrors() []error { return m }

// QueryValidationError is the validation error returned by Query.Validate if
// the designated constraints aren't met.
type QueryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueryValidationError) ErrorName() string { return "QueryValidationError" }

// Error satisfies the builtin error interface
func (e QueryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQuery.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueryValidationError{}

// Validate checks the field values on PromQL with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PromQL) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PromQL with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PromQLMultiError, or nil if none found.
func (m *PromQL) ValidateAll() error {
	return m.validate(true)
}

func (m *PromQL) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetOutPorts()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PromQLValidationError{
					field:  "OutPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PromQLValidationError{
					field:  "OutPorts",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOutPorts()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PromQLValidationError{
				field:  "OutPorts",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for QueryString

	if all {
		switch v := interface{}(m.GetEvaluationInterval()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PromQLValidationError{
					field:  "EvaluationInterval",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PromQLValidationError{
					field:  "EvaluationInterval",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEvaluationInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PromQLValidationError{
				field:  "EvaluationInterval",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PromQLMultiError(errors)
	}

	return nil
}

// PromQLMultiError is an error wrapping multiple validation errors returned by
// PromQL.ValidateAll() if the designated constraints aren't met.
type PromQLMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PromQLMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PromQLMultiError) AllErrors() []error { return m }

// PromQLValidationError is the validation error returned by PromQL.Validate if
// the designated constraints aren't met.
type PromQLValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PromQLValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PromQLValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PromQLValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PromQLValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PromQLValidationError) ErrorName() string { return "PromQLValidationError" }

// Error satisfies the builtin error interface
func (e PromQLValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPromQL.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PromQLValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PromQLValidationError{}

// Validate checks the field values on PromQL_Outs with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PromQL_Outs) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PromQL_Outs with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PromQL_OutsMultiError, or
// nil if none found.
func (m *PromQL_Outs) ValidateAll() error {
	return m.validate(true)
}

func (m *PromQL_Outs) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetOutput()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PromQL_OutsValidationError{
					field:  "Output",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PromQL_OutsValidationError{
					field:  "Output",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOutput()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PromQL_OutsValidationError{
				field:  "Output",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PromQL_OutsMultiError(errors)
	}

	return nil
}

// PromQL_OutsMultiError is an error wrapping multiple validation errors
// returned by PromQL_Outs.ValidateAll() if the designated constraints aren't met.
type PromQL_OutsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PromQL_OutsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PromQL_OutsMultiError) AllErrors() []error { return m }

// PromQL_OutsValidationError is the validation error returned by
// PromQL_Outs.Validate if the designated constraints aren't met.
type PromQL_OutsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PromQL_OutsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PromQL_OutsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PromQL_OutsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PromQL_OutsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PromQL_OutsValidationError) ErrorName() string { return "PromQL_OutsValidationError" }

// Error satisfies the builtin error interface
func (e PromQL_OutsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPromQL_Outs.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PromQL_OutsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PromQL_OutsValidationError{}