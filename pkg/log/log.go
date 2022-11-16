package log

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/fluxninja/aperture/pkg/info"
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
	// FlushWait is the amount of time to wait for the buffer to flush.
	flushWait = 1000 * time.Millisecond
	// DefaultLevel sets info log level, InfoLevel, as default.
	defaultLevel = "info"
	// ServiceKey is a field key that are used with Service name value as a string to the logger context.
	serviceKey = "service"
	// ComponentKey is a field key that are used with Component name value as a string to the logger context.
	componentKey = "component"
	// Sampled is a field key that are used with Sampled value as a bool to the logger context.
	sampledKey = "sampled"
	// BugKey is a field key that is used for Bug() events.
	bugKey = "bug"
)

// Logger is wrapper around zerolog.Logger and io.writers.
type Logger struct {
	logger *zerolog.Logger
}

var global *Logger

// Always create a global logger instance.
func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		// short caller format
		short := strings.Split(file, "/")
		// return the last n elements of the file path
		n := 3
		if len(short) < n {
			n = len(short)
		}
		path := strings.Join(short[len(short)-n:], "/")
		return path + fmt.Sprintf("%d", line)
	}
	// Create global logger
	SetGlobalLogger(NewDefaultLogger())
}

// NewDefaultLogger creates a new default logger with default settings.
func NewDefaultLogger() *Logger {
	return NewLogger(os.Stderr, defaultLevel)
}

// SetGlobalLogger closes the previous global logger and sets given logger as a new global logger.
func SetGlobalLogger(lg *Logger) {
	global = lg
}

// SetStdLogger sets output for the standard logger.
func SetStdLogger(lg *Logger) {
	stdlog.SetFlags(0)
	stdlog.SetOutput(lg.logger)
}

// NewLogger creates a new logger.
func NewLogger(w io.Writer, levelString string) *Logger {
	level, err := zerolog.ParseLevel(levelString)
	if err != nil {
		log.Panic().Err(err).Str("level", level.String()).Msg("Unable to parse logger level")
	}

	zerolog := zerolog.New(w).Level(level).With().Timestamp().Caller().Str(serviceKey, info.Service).Logger()
	logger := &Logger{
		logger: &zerolog,
	}

	return logger
}

// WaitFlush waits a few ms to let the the buffer to flush.
func WaitFlush() {
	time.Sleep(flushWait)
}

// GetPrettyConsoleWriter returns a pretty console writer.
func GetPrettyConsoleWriter() io.Writer {
	output := zerolog.NewConsoleWriter()
	return output
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
func GetGlobalLogger() *Logger {
	return global
}

// WithComponent enables the global logger to chain loggers with additional context, component name.
func WithComponent(component string) *Logger {
	return global.WithComponent(component)
}

// WithComponent enables the current logger to chain loggers with additional context, component name.
func (lg *Logger) WithComponent(component string) *Logger {
	zerolog := lg.logger.With().Str(componentKey, component).Logger()
	return &Logger{
		logger: &zerolog,
	}
}

// WithInterface adds an interface to the logger context.
func (lg *Logger) WithInterface(key string, value interface{}) *Logger {
	zerolog := lg.logger.With().Interface(key, value).Logger()
	return &Logger{
		logger: &zerolog,
	}
}

// WithInterface adds an interface to the global logger context.
func WithInterface(key string, value interface{}) *Logger {
	return global.WithInterface(key, value)
}

// WithStr adds a string to the logger context.
func (lg *Logger) WithStr(key string, value string) *Logger {
	zerolog := lg.logger.With().Str(key, value).Logger()
	return &Logger{
		logger: &zerolog,
	}
}

// WithStr adds a string to the global logger context.
func WithStr(key string, value string) *Logger {
	return global.WithStr(key, value)
}

// WithBool adds a bool to the logger context.
func (lg *Logger) WithBool(key string, value bool) *Logger {
	zerolog := lg.logger.With().Bool(key, value).Logger()
	return &Logger{
		logger: &zerolog,
	}
}

// WithBool adds a bool to the global logger context.
func WithBool(key string, value bool) *Logger {
	return global.WithBool(key, value)
}

// GetZerolog returns underlying zerolog logger.
func (lg *Logger) GetZerolog() *zerolog.Logger {
	return lg.logger
}

// NewFromZerolog creates the logger from zerolog instance.
func (lg *Logger) NewFromZerolog(logger *zerolog.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Output duplicates the current logger and sets w as its output.
func (lg *Logger) Output(w io.Writer) *Logger {
	zerolog := lg.logger.Output(w)
	return &Logger{
		logger: &zerolog,
	}
}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) *Logger {
	return global.Output(w)
}

// With creates a child logger of the current logger with the field added to its context.
func (lg *Logger) With() zerolog.Context {
	return lg.logger.With()
}

// With creates a child logger of the global logger with the field added to its context.
func With() zerolog.Context {
	return global.With()
}

// Level creates a child logger of the current logger with the minimum accepted level set to level.
func (lg *Logger) Level(level zerolog.Level) *Logger {
	zerolog := lg.logger.Level(level)
	return &Logger{
		logger: &zerolog,
	}
}

// Level creates a child logger of the global logger with the minimum accepted level set to level.
func Level(level zerolog.Level) *Logger {
	return global.Level(level)
}

// GetLevel returns the current logger level.
func (lg *Logger) GetLevel() zerolog.Level {
	return lg.logger.GetLevel()
}

// GetLevel returns the global logger level.
func GetLevel() zerolog.Level {
	return global.GetLevel()
}

// Sample returns the current logger with the s sampler.
func (lg *Logger) Sample(sampler zerolog.Sampler) *Logger {
	zerolog := lg.WithBool(sampledKey, true).logger.Sample(sampler)
	return &Logger{
		logger: &zerolog,
	}
}

// Sample returns the global logger with the s sampler.
func Sample(sampler zerolog.Sampler) *Logger {
	return global.Sample(sampler)
}

// Autosample returns the current logger with sampler based on caller location.
//
// Sampler will be created using NewRatelimitingSampler()
//
// This is basically shorthand for:
//
// ```go
// var mySampler = NewRatelimitingSampler()
// ...
// Logger.Sample(mySampler)
// ```
//
// The "auto" part has a slight runtime cost though, so the full should be
// preferred for cases where performance matters, like on datapath.
func (lg *Logger) Autosample() *Logger {
	return lg.Sample(getAutosampler())
}

// Autosample returns the global logger with sampler based on caller location.
//
// See Logger.Autosample().
func Autosample() *Logger {
	// Note: not calling global.Autosample() as it might mess up with caller
	// depth in getAutosampler().
	return Sample(getAutosampler())
}

// Bug starts a new message with "bug" level
//
// "Bug" is the same level as "warn", but it's intended for programmer's
// errors, where normally you'd want to use "panic", but:
// * error is not affecting the service as-a-whole,
// * there's reasonable way to continue,
// * restarting service won't fix the error.
//
// You might want to use Bug() for cases like hitting "impossible" case of a
// switch in an rpc call, or similar.
//
// Additionally, every callsite of Bug will be automatically ratelimited (like
// with Autosample).
//
// Also, automatic bug-reporting may be integrated here.
//
// You must call Msg on the returned event in order to send the event.
func (lg *Logger) Bug() *zerolog.Event {
	return bugWithSampler(lg, getAutosampler())
}

// Bug starts a new message with "bug" level
//
// See Logger.Bug()
//
// You must call Msg on the returned event in order to send the event.
func Bug() *zerolog.Event {
	// Note: not calling global.Bug() as it might mess up with caller depth in
	// getAutosampler().
	return bugWithSampler(global, getAutosampler())
}

// Hook returns the current logger with the h hook.
func (lg *Logger) Hook(hook zerolog.Hook) *Logger {
	zerolog := lg.logger.Hook(hook)
	return &Logger{
		logger: &zerolog,
	}
}

// Hook returns the global logger with the h hook.
func Hook(hook zerolog.Hook) *Logger {
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

// Write implements io.Writer interface.
func (lg *Logger) Write(p []byte) (n int, err error) {
	return lg.logger.Write(p)
}

// Write implements io.Writer interface.
func Write(p []byte) (n int, err error) {
	return global.logger.Write(p)
}
