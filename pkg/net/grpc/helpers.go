package grpc

import (
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/fluxninja/aperture/pkg/log"
)

// LogEvent wraps *zerolog.Event, so that sending an event also returns an grpc
// error
//
// LogEvent also forwards _some_ builder-functions (like Str()) to the
// underlying zerolog.Event. This is just a sugar, as you can always use
// WithEvent on the whole chain.
type LogEvent struct {
	e    *zerolog.Event
	code codes.Code
}

// LoggedError wraps zerolog.Event so that sending an event will also return a
// grpc error.
//
// GRPC error code should be provided using Code(), otherwise, the error will
// use codes.Unknown.
//
// Example:
//
//	return nil, grpc.LoggedError(log.Autosample().Warn()).
//	    Code(codes.InvalidArgument).
//	    Msg("missing frobnicator")
func LoggedError(e *zerolog.Event) LogEvent {
	return LogEvent{
		e:    e,
		code: codes.Unknown,
	}
}

// Bug is equivalent to WithLogEvent(log.Bug()).Code(codes.Internal)
//
// Example: return nil, grpc.Bug().Msg("impossible happened").
func Bug() LogEvent {
	// Note: not forwarding to log.Bug() as it might mess up with caller depth
	// in GetAutosampler()
	return LogEvent{
		e:    log.BugWithSampler(log.GetGlobalLogger(), log.GetAutosampler()),
		code: codes.Internal,
	}
}

// BugWithLogger is equivalent to WithLogEvent(logger.Bug()).Code(codes.Internal).
func BugWithLogger(lg *log.Logger) LogEvent {
	// Note: not forwarding to lg.Bug() as it might mess up with caller depth
	// in GetAutosampler()
	return LogEvent{
		e:    log.BugWithSampler(lg, log.GetAutosampler()),
		code: codes.Internal,
	}
}

// Code sets the gRPC code to be used in error returned via Msg() or Send()
//
// Additionally, it adds the code field to the event.
func (e LogEvent) Code(code codes.Code) LogEvent {
	e.e.Stringer("code", code)
	e.code = code
	return e
}

// Msg sends the *Event with msg added as the message field if not empty.
//
// Msg returns grpc error using given msg as message and code previously set
// with Code
//
// NOTICE: once this method is called, the LogEvent should be disposed.
// Calling Msg twice can have unexpected result.
func (e LogEvent) Msg(msg string) error {
	e.e.Msg(msg)
	return status.Error(e.code, msg)
}

// Send is equivalent to calling Msg("").
func (e LogEvent) Send() error {
	e.e.Msg("")
	return status.Error(e.code, "")
}

// *** forwarded builders ***

// Str adds the field key with val as a string to the *Event context.
func (e LogEvent) Str(key, val string) LogEvent {
	e.e.Str(key, val)
	return e
}

// Stringer adds the field key with val.String() (or null if val is nil)
// to the *Event context.
func (e LogEvent) Stringer(key string, val fmt.Stringer) LogEvent {
	e.e.Stringer(key, val)
	return e
}

// Err adds the field "error" with serialized err to the *Event context.
// If err is nil, no field is added.
//
// To customize the key name, change zerolog.ErrorFieldName.
//
// If Stack() has been called before and zerolog.ErrorStackMarshaler is defined,
// the err is passed to ErrorStackMarshaler and the result is appended to the
// zerolog.ErrorStackFieldName.
func (e LogEvent) Err(err error) LogEvent {
	e.e.Err(err)
	return e
}

// Bool adds the field key with val as a bool to the *Event context.
func (e LogEvent) Bool(key string, b bool) LogEvent {
	e.e.Bool(key, b)
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e LogEvent) Int(key string, i int) LogEvent {
	e.e.Int(key, i)
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e LogEvent) Float32(key string, f float32) LogEvent {
	e.e.Float32(key, f)
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e LogEvent) Float64(key string, f float64) LogEvent {
	e.e.Float64(key, f)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e LogEvent) Interface(key string, i interface{}) LogEvent {
	e.e.Interface(key, i)
	return e
}
