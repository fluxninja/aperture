// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: aperture/autoscale/kubernetes/controlpoints/v1/controlpoints.proto

package controlpointsv1

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

// Validate checks the field values on AutoScaleKubernetesControlPoints with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *AutoScaleKubernetesControlPoints) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AutoScaleKubernetesControlPoints with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// AutoScaleKubernetesControlPointsMultiError, or nil if none found.
func (m *AutoScaleKubernetesControlPoints) ValidateAll() error {
	return m.validate(true)
}

func (m *AutoScaleKubernetesControlPoints) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetAutoScaleKubernetesControlPoints() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, AutoScaleKubernetesControlPointsValidationError{
						field:  fmt.Sprintf("AutoScaleKubernetesControlPoints[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, AutoScaleKubernetesControlPointsValidationError{
						field:  fmt.Sprintf("AutoScaleKubernetesControlPoints[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AutoScaleKubernetesControlPointsValidationError{
					field:  fmt.Sprintf("AutoScaleKubernetesControlPoints[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return AutoScaleKubernetesControlPointsMultiError(errors)
	}

	return nil
}

// AutoScaleKubernetesControlPointsMultiError is an error wrapping multiple
// validation errors returned by
// AutoScaleKubernetesControlPoints.ValidateAll() if the designated
// constraints aren't met.
type AutoScaleKubernetesControlPointsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AutoScaleKubernetesControlPointsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AutoScaleKubernetesControlPointsMultiError) AllErrors() []error { return m }

// AutoScaleKubernetesControlPointsValidationError is the validation error
// returned by AutoScaleKubernetesControlPoints.Validate if the designated
// constraints aren't met.
type AutoScaleKubernetesControlPointsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AutoScaleKubernetesControlPointsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AutoScaleKubernetesControlPointsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AutoScaleKubernetesControlPointsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AutoScaleKubernetesControlPointsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AutoScaleKubernetesControlPointsValidationError) ErrorName() string {
	return "AutoScaleKubernetesControlPointsValidationError"
}

// Error satisfies the builtin error interface
func (e AutoScaleKubernetesControlPointsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAutoScaleKubernetesControlPoints.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AutoScaleKubernetesControlPointsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AutoScaleKubernetesControlPointsValidationError{}

// Validate checks the field values on AutoScaleKubernetesControlPoint with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AutoScaleKubernetesControlPoint) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AutoScaleKubernetesControlPoint with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// AutoScaleKubernetesControlPointMultiError, or nil if none found.
func (m *AutoScaleKubernetesControlPoint) ValidateAll() error {
	return m.validate(true)
}

func (m *AutoScaleKubernetesControlPoint) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ApiVersion

	// no validation rules for Kind

	// no validation rules for Namespace

	// no validation rules for Name

	if len(errors) > 0 {
		return AutoScaleKubernetesControlPointMultiError(errors)
	}

	return nil
}

// AutoScaleKubernetesControlPointMultiError is an error wrapping multiple
// validation errors returned by AutoScaleKubernetesControlPoint.ValidateAll()
// if the designated constraints aren't met.
type AutoScaleKubernetesControlPointMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AutoScaleKubernetesControlPointMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AutoScaleKubernetesControlPointMultiError) AllErrors() []error { return m }

// AutoScaleKubernetesControlPointValidationError is the validation error
// returned by AutoScaleKubernetesControlPoint.Validate if the designated
// constraints aren't met.
type AutoScaleKubernetesControlPointValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AutoScaleKubernetesControlPointValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AutoScaleKubernetesControlPointValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AutoScaleKubernetesControlPointValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AutoScaleKubernetesControlPointValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AutoScaleKubernetesControlPointValidationError) ErrorName() string {
	return "AutoScaleKubernetesControlPointValidationError"
}

// Error satisfies the builtin error interface
func (e AutoScaleKubernetesControlPointValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAutoScaleKubernetesControlPoint.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AutoScaleKubernetesControlPointValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AutoScaleKubernetesControlPointValidationError{}
