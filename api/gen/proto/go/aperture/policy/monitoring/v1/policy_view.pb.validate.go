// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: aperture/policy/monitoring/v1/policy_view.proto

package monitoringv1

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

// Validate checks the field values on PortView with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PortView) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PortView with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PortViewMultiError, or nil
// if none found.
func (m *PortView) ValidateAll() error {
	return m.validate(true)
}

func (m *PortView) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PortName

	// no validation rules for Looped

	// no validation rules for SubCircuitId

	switch v := m.Value.(type) {
	case *PortView_SignalName:
		if v == nil {
			err := PortViewValidationError{
				field:  "Value",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		// no validation rules for SignalName
	case *PortView_ConstantValue:
		if v == nil {
			err := PortViewValidationError{
				field:  "Value",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		// no validation rules for ConstantValue
	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return PortViewMultiError(errors)
	}

	return nil
}

// PortViewMultiError is an error wrapping multiple validation errors returned
// by PortView.ValidateAll() if the designated constraints aren't met.
type PortViewMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PortViewMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PortViewMultiError) AllErrors() []error { return m }

// PortViewValidationError is the validation error returned by
// PortView.Validate if the designated constraints aren't met.
type PortViewValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PortViewValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PortViewValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PortViewValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PortViewValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PortViewValidationError) ErrorName() string { return "PortViewValidationError" }

// Error satisfies the builtin error interface
func (e PortViewValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPortView.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PortViewValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PortViewValidationError{}

// Validate checks the field values on ComponentView with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ComponentView) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ComponentView with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ComponentViewMultiError, or
// nil if none found.
func (m *ComponentView) ValidateAll() error {
	return m.validate(true)
}

func (m *ComponentView) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ComponentId

	// no validation rules for ComponentName

	// no validation rules for ComponentType

	// no validation rules for ComponentDescription

	if all {
		switch v := interface{}(m.GetComponent()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ComponentViewValidationError{
					field:  "Component",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ComponentViewValidationError{
					field:  "Component",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetComponent()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ComponentViewValidationError{
				field:  "Component",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetInPorts() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ComponentViewValidationError{
						field:  fmt.Sprintf("InPorts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ComponentViewValidationError{
						field:  fmt.Sprintf("InPorts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ComponentViewValidationError{
					field:  fmt.Sprintf("InPorts[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetOutPorts() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ComponentViewValidationError{
						field:  fmt.Sprintf("OutPorts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ComponentViewValidationError{
						field:  fmt.Sprintf("OutPorts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ComponentViewValidationError{
					field:  fmt.Sprintf("OutPorts[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ComponentViewMultiError(errors)
	}

	return nil
}

// ComponentViewMultiError is an error wrapping multiple validation errors
// returned by ComponentView.ValidateAll() if the designated constraints
// aren't met.
type ComponentViewMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ComponentViewMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ComponentViewMultiError) AllErrors() []error { return m }

// ComponentViewValidationError is the validation error returned by
// ComponentView.Validate if the designated constraints aren't met.
type ComponentViewValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ComponentViewValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ComponentViewValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ComponentViewValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ComponentViewValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ComponentViewValidationError) ErrorName() string { return "ComponentViewValidationError" }

// Error satisfies the builtin error interface
func (e ComponentViewValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sComponentView.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ComponentViewValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ComponentViewValidationError{}

// Validate checks the field values on SourceTarget with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SourceTarget) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SourceTarget with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SourceTargetMultiError, or
// nil if none found.
func (m *SourceTarget) ValidateAll() error {
	return m.validate(true)
}

func (m *SourceTarget) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ComponentId

	// no validation rules for PortName

	if len(errors) > 0 {
		return SourceTargetMultiError(errors)
	}

	return nil
}

// SourceTargetMultiError is an error wrapping multiple validation errors
// returned by SourceTarget.ValidateAll() if the designated constraints aren't met.
type SourceTargetMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SourceTargetMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SourceTargetMultiError) AllErrors() []error { return m }

// SourceTargetValidationError is the validation error returned by
// SourceTarget.Validate if the designated constraints aren't met.
type SourceTargetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SourceTargetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SourceTargetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SourceTargetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SourceTargetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SourceTargetValidationError) ErrorName() string { return "SourceTargetValidationError" }

// Error satisfies the builtin error interface
func (e SourceTargetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSourceTarget.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SourceTargetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SourceTargetValidationError{}

// Validate checks the field values on Link with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Link) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Link with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in LinkMultiError, or nil if none found.
func (m *Link) ValidateAll() error {
	return m.validate(true)
}

func (m *Link) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetSource()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LinkValidationError{
					field:  "Source",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LinkValidationError{
					field:  "Source",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSource()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LinkValidationError{
				field:  "Source",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetTarget()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, LinkValidationError{
					field:  "Target",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, LinkValidationError{
					field:  "Target",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTarget()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LinkValidationError{
				field:  "Target",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Looped

	// no validation rules for SubCircuitId

	switch v := m.Value.(type) {
	case *Link_SignalName:
		if v == nil {
			err := LinkValidationError{
				field:  "Value",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		// no validation rules for SignalName
	case *Link_ConstantValue:
		if v == nil {
			err := LinkValidationError{
				field:  "Value",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
		// no validation rules for ConstantValue
	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return LinkMultiError(errors)
	}

	return nil
}

// LinkMultiError is an error wrapping multiple validation errors returned by
// Link.ValidateAll() if the designated constraints aren't met.
type LinkMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LinkMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LinkMultiError) AllErrors() []error { return m }

// LinkValidationError is the validation error returned by Link.Validate if the
// designated constraints aren't met.
type LinkValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LinkValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LinkValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LinkValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LinkValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LinkValidationError) ErrorName() string { return "LinkValidationError" }

// Error satisfies the builtin error interface
func (e LinkValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLink.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LinkValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LinkValidationError{}

// Validate checks the field values on Graph with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Graph) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Graph with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in GraphMultiError, or nil if none found.
func (m *Graph) ValidateAll() error {
	return m.validate(true)
}

func (m *Graph) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetInternalComponents() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("InternalComponents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("InternalComponents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GraphValidationError{
					field:  fmt.Sprintf("InternalComponents[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetExternalComponents() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("ExternalComponents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("ExternalComponents[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GraphValidationError{
					field:  fmt.Sprintf("ExternalComponents[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetInternalLinks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("InternalLinks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("InternalLinks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GraphValidationError{
					field:  fmt.Sprintf("InternalLinks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetExternalLinks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("ExternalLinks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GraphValidationError{
						field:  fmt.Sprintf("ExternalLinks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GraphValidationError{
					field:  fmt.Sprintf("ExternalLinks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GraphMultiError(errors)
	}

	return nil
}

// GraphMultiError is an error wrapping multiple validation errors returned by
// Graph.ValidateAll() if the designated constraints aren't met.
type GraphMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GraphMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GraphMultiError) AllErrors() []error { return m }

// GraphValidationError is the validation error returned by Graph.Validate if
// the designated constraints aren't met.
type GraphValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GraphValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GraphValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GraphValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GraphValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GraphValidationError) ErrorName() string { return "GraphValidationError" }

// Error satisfies the builtin error interface
func (e GraphValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGraph.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GraphValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GraphValidationError{}

// Validate checks the field values on CircuitView with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CircuitView) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CircuitView with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CircuitViewMultiError, or
// nil if none found.
func (m *CircuitView) ValidateAll() error {
	return m.validate(true)
}

func (m *CircuitView) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTree()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CircuitViewValidationError{
					field:  "Tree",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CircuitViewValidationError{
					field:  "Tree",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTree()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CircuitViewValidationError{
				field:  "Tree",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CircuitViewMultiError(errors)
	}

	return nil
}

// CircuitViewMultiError is an error wrapping multiple validation errors
// returned by CircuitView.ValidateAll() if the designated constraints aren't met.
type CircuitViewMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CircuitViewMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CircuitViewMultiError) AllErrors() []error { return m }

// CircuitViewValidationError is the validation error returned by
// CircuitView.Validate if the designated constraints aren't met.
type CircuitViewValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CircuitViewValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CircuitViewValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CircuitViewValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CircuitViewValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CircuitViewValidationError) ErrorName() string { return "CircuitViewValidationError" }

// Error satisfies the builtin error interface
func (e CircuitViewValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCircuitView.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CircuitViewValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CircuitViewValidationError{}

// Validate checks the field values on Tree with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Tree) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Tree with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in TreeMultiError, or nil if none found.
func (m *Tree) ValidateAll() error {
	return m.validate(true)
}

func (m *Tree) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetNode()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TreeValidationError{
					field:  "Node",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TreeValidationError{
					field:  "Node",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetNode()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TreeValidationError{
				field:  "Node",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetGraph()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TreeValidationError{
					field:  "Graph",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TreeValidationError{
					field:  "Graph",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGraph()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TreeValidationError{
				field:  "Graph",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetChildren() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TreeValidationError{
						field:  fmt.Sprintf("Children[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TreeValidationError{
						field:  fmt.Sprintf("Children[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TreeValidationError{
					field:  fmt.Sprintf("Children[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetActuators() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TreeValidationError{
						field:  fmt.Sprintf("Actuators[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TreeValidationError{
						field:  fmt.Sprintf("Actuators[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TreeValidationError{
					field:  fmt.Sprintf("Actuators[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return TreeMultiError(errors)
	}

	return nil
}

// TreeMultiError is an error wrapping multiple validation errors returned by
// Tree.ValidateAll() if the designated constraints aren't met.
type TreeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TreeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TreeMultiError) AllErrors() []error { return m }

// TreeValidationError is the validation error returned by Tree.Validate if the
// designated constraints aren't met.
type TreeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TreeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TreeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TreeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TreeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TreeValidationError) ErrorName() string { return "TreeValidationError" }

// Error satisfies the builtin error interface
func (e TreeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTree.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TreeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TreeValidationError{}
