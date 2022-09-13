package config

import (
	"go.uber.org/fx/fxevent"

	"github.com/fluxninja/aperture/pkg/log"
)

type logger struct {
	Logger log.Logger
}

// LogEvent is the event fired when a log event is written.
func (lo *logger) LogEvent(event fxevent.Event) {
	l := lo.Logger.Component("FX")
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Trace().
			Str("fxevent", "ON_START").
			Str("function", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Error().Err(e.Err).
				Str("fxevent", "ON_START").
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Dur("runtime", e.Runtime).
				Msg("OnStart call failed")
		} else {
			l.Trace().
				Str("fxevent", "ON_START").
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Dur("runtime", e.Runtime).
				Msg("OnStart call executed")
		}
	case *fxevent.OnStopExecuting:
		l.Trace().
			Str("fxevent", "ON_STOP").
			Str("function", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Error().Err(e.Err).
				Str("fxevent", "ON_STOP").
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Dur("runtime", e.Runtime).
				Msg("OnStop call failed")
		} else {
			l.Trace().
				Str("fxevent", "ON_STOP").
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Dur("runtime", e.Runtime).
				Msg("OnStop call executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Error().Err(e.Err).
				Str("fxevent", "SUPPLY").
				Str("type", e.TypeName).
				Msg("Failed to supply")
		} else {
			l.Trace().
				Str("fxevent", "SUPPLY").
				Str("type", e.TypeName).
				Msg("Supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Trace().
				Str("fxevent", "PROVIDE").
				Str("type", rtype).
				Str("constructor", e.ConstructorName).
				Msg("Constructor provided")
		}
		if e.Err != nil {
			l.Error().Err(e.Err).
				Str("fxevent", "PROVIDE").
				Msg("Error after options were applied")
		}
	case *fxevent.Invoking:
		l.Trace().
			Str("fxevent", "INVOKE").
			Str("function", e.FunctionName).
			Msg("Invoking function")

	case *fxevent.Invoked:
		if e.Err != nil {
			log.Error().Err(e.Err).
				Str("fxevent", "INVOKE").
				Str("function name", e.FunctionName).
				Str("trace", e.Trace).
				Msg("Failed invoke")
		}
	case *fxevent.Stopping:
		l.Trace().
			Str("fxevent", "STOPPING").
			Msg("Stopping")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Error().
				Err(e.Err).
				Str("fxevent", "STOPPED").
				Msg("Failed to stop cleanly")
		}
	case *fxevent.RollingBack:
		l.Error().
			Err(e.StartErr).
			Str("fxevent", "ROLLING_BACK").
			Msg("Start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Error().
				Err(e.Err).
				Str("fxevent", "ROLLING_BACK").
				Msg("Couldn't roll back cleanly")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Error().
				Err(e.Err).
				Str("fxevent", "STARTING").
				Msg("Failed to start")
		} else {
			l.Trace().
				Str("fxevent", "STARTING").
				Msg("Running")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Error().
				Err(e.Err).
				Str("fxevent", "LOGGER").
				Msg("Failed to initialize custom logger")
		} else {
			l.Trace().
				Str("fxevent", "LOGGER").
				Str("name", e.ConstructorName).
				Msg("Initialized custom logger")
		}
	}
}

// WithApertureLogger overrides fx default logger.
func WithApertureLogger() func(lg log.Logger) fxevent.Logger {
	return func(lg log.Logger) fxevent.Logger {
		return &logger{Logger: lg}
	}
}
