package log

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/panichandler"
)

const (
	// DebugLevel defines debug log level.
	DebugLevel = zerolog.DebugLevel
	// InfoLevel defines info log level.
	InfoLevel = zerolog.InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel = zerolog.WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel = zerolog.ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel = zerolog.FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel = zerolog.PanicLevel
	// NoLevel defines an absent log level.
	NoLevel = zerolog.NoLevel
	// Disabled disables the logger.
	Disabled = zerolog.Disabled
	// TraceLevel defines trace log level.
	TraceLevel = zerolog.TraceLevel
)

const (
	// DiodeFlushWait is the amount of time to wait for diode buffer to flush.
	diodeFlushWait = 1000 * time.Millisecond
	// DefaultLevel sets info log level, InfoLevel, as default.
	defaultLevel = "info"
	// ServiceKey is a field key that are used with Service name value as a string to the logger context.
	serviceKey = "service"
)

// Logger is wrapper around zerolog.Logger and io.writers.
type Logger struct {
	logger *zerolog.Logger
	w      io.Writer
	// flag to detect valid Logger instance
	valid bool
}

var global Logger

// Always create a global logger instance.
func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// Create global logger
	SetGlobalLogger(NewDefaultLogger())
}

// NewDefaultLogger creates a new default logger with default settings.
func NewDefaultLogger() Logger {
	return NewLogger(os.Stderr, false, defaultLevel)
}

// SetGlobalLogger closes the previous global logger and sets given logger as a new global logger.
func SetGlobalLogger(lg Logger) {
	global.Close()
	global = lg
}

// SetStdLogger sets output for the standard logger.
func SetStdLogger(lg Logger) {
	stdlog.SetFlags(0)
	stdlog.SetOutput(lg.logger)
}

// NewLogger creates a new logger by wrapping io.Writer in a diode. Make sure to call Logger.Close() once done with this logger.
func NewLogger(w io.Writer, useDiode bool, levelString string) Logger {
	level, err := zerolog.ParseLevel(levelString)
	if err != nil {
		log.Panic().Err(err).Str("level", level.String()).Msg("Unable to parse logger level")
	}

	var wr io.Writer

	if useDiode {
		// Use diode writer
		wr = diode.NewWriter(w, 1000, 0, func(missed int) {
			Printf("Dropped %d messages", missed)
		})
	} else {
		wr = w
	}
	zerolog := zerolog.New(wr).Level(level).With().Timestamp().Caller().Str(serviceKey, info.Service).Logger()
	logger := Logger{
		logger: &zerolog,
		w:      wr,
		valid:  true,
	}

	return logger
}

// Close closes all the underlying diode Writer when there are valid Logger instances.
func (lg Logger) Close() {
	if lg.valid {
		if dw, ok := lg.w.(diode.Writer); ok {
			closeDiodeWriter(dw)
		}
	}
}

// closeDiodeWriter.
func closeDiodeWriter(dw diode.Writer) {
	log.Info().Msg("Closing DiodeWriter after a delay!")
	panichandler.Go(func() {
		WaitFlush()
		_ = dw.Close()
	})
}

// WaitFlush waits a few ms to let the diode buffer to flush.
func WaitFlush() {
	time.Sleep(diodeFlushWait)
}

/* Wrappers around zerolog */

// SetGlobalLevel sets the global log level with given level.
func SetGlobalLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// SetGlobalLevelString parses given levelString and sets the global log level.
func SetGlobalLevelString(levelString string) error {
	level, err := zerolog.ParseLevel(levelString)
	if err != nil {
		return err
	}

	SetGlobalLevel(level)
	return nil
}

// GetGlobalLogger returns the global logger.
func GetGlobalLogger() Logger {
	return global
}

// Component enables the global logger to chain loggers with additional context, component name.
func Component(component string) Logger {
	return global.Component(component)
}

// Component enables the current logger to chain loggers with additional context, component name.
func (lg *Logger) Component(component string) Logger {
	zerolog := lg.logger.With().Str("component", component).Logger()
	return Logger{
		logger: &zerolog,
		w:      lg.w,
		valid:  lg.valid,
	}
}

// Zerolog returns underlying zerolog logger.
func (lg *Logger) Zerolog() *zerolog.Logger {
	return lg.logger
}

// Output duplicates the current logger and sets w as its output.
func (lg Logger) Output(w io.Writer) Logger {
	zerolog := lg.logger.Output(w)
	return Logger{
		logger: &zerolog,
		w:      lg.w,
		valid:  lg.valid,
	}
}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) Logger {
	return global.Output(w)
}

// With creates a child logger of the current logger with the field added to its context.
func (lg Logger) With() zerolog.Context {
	return lg.logger.With()
}

// With creates a child logger of the global logger with the field added to its context.
func With() zerolog.Context {
	return global.With()
}

// Level creates a child logger of the current logger with the minimum accepted level set to level.
func (lg Logger) Level(level zerolog.Level) Logger {
	zerolog := lg.logger.Level(level)
	return Logger{
		logger: &zerolog,
		w:      lg.w,
		valid:  lg.valid,
	}
}

// Level creates a child logger of the global logger with the minimum accepted level set to level.
func Level(level zerolog.Level) Logger {
	return global.Level(level)
}

// Sample returns the current logger with the s sampler.
func (lg Logger) Sample(sampler zerolog.Sampler) Logger {
	zerolog := lg.logger.Sample(sampler)
	return Logger{
		logger: &zerolog,
		w:      lg.w,
		valid:  lg.valid,
	}
}

// Sample returns the global logger with the s sampler.
func Sample(sampler zerolog.Sampler) Logger {
	return global.Sample(sampler)
}

// Hook returns the current logger with the h hook.
func (lg Logger) Hook(hook zerolog.Hook) Logger {
	zerolog := lg.logger.Hook(hook)
	return Logger{
		logger: &zerolog,
		w:      lg.w,
		valid:  lg.valid,
	}
}

// Hook returns the global logger with the h hook.
func Hook(hook zerolog.Hook) Logger {
	return global.Hook(hook)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Trace() *zerolog.Event {
	return lg.logger.Trace()
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return global.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Debug() *zerolog.Event {
	return lg.logger.Debug()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return global.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Info() *zerolog.Event {
	return lg.logger.Info()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return global.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Warn() *zerolog.Event {
	return lg.logger.Warn()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return global.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Error() *zerolog.Event {
	return lg.logger.Error()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return global.Error()
}

// Fatal starts a new message with fatal level. This is an alias for Panic.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Fatal() *zerolog.Event {
	return lg.logger.Panic()
}

// Fatal starts a new message with fatal level.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return global.Fatal()
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine and
// invokes any registered PanicHandler.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Panic() *zerolog.Event {
	return lg.logger.Panic()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return global.Panic()
}

// WithLevel starts a new message with level. Unlike Fatal and Panic
// methods, WithLevel does not terminate the program or stop the ordinary
// flow of a goroutine when used with their respective levels.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return lg.logger.WithLevel(level)
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return global.WithLevel(level)
}

// Log starts a new message with no level. This is equivalent to using lg.WithLevel(NoLevel).
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Log() *zerolog.Event {
	return lg.logger.Log()
}

// Log starts a new message with no level. This is equivalent to using WithLevel(NoLevel).
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return global.Log()
}

func printfEvent(e *zerolog.Event, format string, v ...interface{}) {
	if e.Enabled() {
		e.CallerSkipFrame(2).Msg(fmt.Sprintf(format, v...))
	}
}

func printEvent(e *zerolog.Event, v ...interface{}) {
	if e.Enabled() {
		e.CallerSkipFrame(2).Msg(fmt.Sprint(v...))
	}
}

func printlnEvent(e *zerolog.Event, v ...interface{}) {
	if e.Enabled() {
		e.CallerSkipFrame(2).Msg(fmt.Sprintln(v...))
	}
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (lg *Logger) Print(v ...interface{}) {
	printEvent(lg.logger.Debug(), v...)
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	printEvent(global.logger.Debug(), v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Printf(format string, v ...interface{}) {
	printfEvent(lg.logger.Debug(), format, v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	printfEvent(global.logger.Debug(), format, v...)
}

// Println sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Println(v ...interface{}) {
	printlnEvent(lg.logger.Debug(), v...)
}

// Println sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	printlnEvent(global.logger.Debug(), v...)
}

// more adapters

/* Fatal */

// Fatalf sends a log event using fatal level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Fatalf(format string, v ...interface{}) {
	printfEvent(lg.logger.Fatal(), format, v...)
}

// Fatalf sends a log event using fatal level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	printfEvent(global.logger.Fatal(), format, v...)
}

// Fatalln sends a log event using fatal level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Fatalln(v ...interface{}) {
	printlnEvent(lg.logger.Fatal(), v...)
}

// Fatalln sends a log event using fatal level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Fatalln(v ...interface{}) {
	printlnEvent(global.logger.Fatal(), v...)
}

/* Panic */

// Panicf sends a log event using panic level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Panicf(format string, v ...interface{}) {
	printfEvent(lg.logger.Panic(), format, v...)
}

// Panicf sends a log event using panic level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Panicf(format string, v ...interface{}) {
	printfEvent(global.logger.Panic(), format, v...)
}

// Panicln sends a log event using panic level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Panicln(v ...interface{}) {
	printlnEvent(lg.logger.Panic(), v...)
}

// Panicln sends a log event using panic level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Panicln(v ...interface{}) {
	printlnEvent(global.logger.Panic(), v...)
}

/* Debug */

// Debugf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Debugf(format string, v ...interface{}) {
	printfEvent(lg.logger.Debug(), format, v...)
}

// Debugf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	printfEvent(global.logger.Debug(), format, v...)
}

// Debugln sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Debugln(v ...interface{}) {
	printlnEvent(lg.logger.Debug(), v...)
}

// Debugln sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	printlnEvent(global.logger.Debug(), v...)
}

/* Info */

// Infof sends a log event using info level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Infof(format string, v ...interface{}) {
	printfEvent(lg.logger.Info(), format, v...)
}

// Infof sends a log event using info level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	printfEvent(global.logger.Info(), format, v...)
}

// Infoln sends a log event using info level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Infoln(v ...interface{}) {
	printlnEvent(lg.logger.Info(), v...)
}

// Infoln sends a log event using info level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	printlnEvent(global.logger.Info(), v...)
}

/* Warn */

// Warnf sends a log event using warn level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Warnf(format string, v ...interface{}) {
	printfEvent(lg.logger.Warn(), format, v...)
}

// Warnf sends a log event using warn level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	printfEvent(global.logger.Warn(), format, v...)
}

// Warnln sends a log event using warn level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Warnln(v ...interface{}) {
	printlnEvent(lg.logger.Warn(), v...)
}

// Warnln sends a log event using warn level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Warnln(v ...interface{}) {
	printlnEvent(global.logger.Warn(), v...)
}

/* Error */

// Errorf sends a log event using error level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (lg *Logger) Errorf(format string, v ...interface{}) {
	printfEvent(lg.logger.Error(), format, v...)
}

// Errorf sends a log event using error level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	printfEvent(global.logger.Error(), format, v...)
}

// Errorln sends a log event using error level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func (lg *Logger) Errorln(v ...interface{}) {
	printlnEvent(lg.logger.Error(), v...)
}

// Errorln sends a log event using error level and no extra field.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	printlnEvent(global.logger.Error(), v...)
}
