package config

import (
	"context"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	panichandler "github.com/fluxninja/aperture/v2/pkg/panic-handler"
	"github.com/fluxninja/lumberjack"
)

// defaultLogFilePath is the default path for the log files to be stored.
var defaultLogFilePath = path.Join(DefaultLogDirectory, info.Service+".log")

const (
	configKey   = "log"
	stdOutFile  = "stdout"
	stdErrFile  = "stderr"
	defaultFile = "default"
)

// LogModule is a fx module that provides a logger and invokes setting global and standard loggers.
func LogModule() fx.Option {
	return fx.Options(
		LoggerConstructor{ConfigKey: configKey, IsGlobal: true}.Annotate(),
		fx.WithLogger(FxLogger()),
	)
}

// swagger:operation POST /log common-configuration Log
// ---
// x-fn-config-env: true
// parameters:
// - in: body
//   schema:
//     "$ref": "#/definitions/LogConfig"

// LogConfig holds configuration for a logger and log writers.
// swagger:model
// +kubebuilder:object:generate=true
type LogConfig struct {
	// Log level
	LogLevel string `json:"level" validate:"oneof=debug DEBUG info INFO warn WARN error ERROR fatal FATAL panic PANIC trace TRACE disabled DISABLED" default:"info"`

	// Log writers
	Writers []LogWriterConfig `json:"writers,omitempty" validate:"omitempty,dive,omitempty"`

	// Use non-blocking log writer (can lose logs at high throughput)
	NonBlocking bool `json:"non_blocking" default:"true"`

	// Additional log writer: pretty console (`stdout`) logging (not recommended for prod environments)
	PrettyConsole bool `json:"pretty_console" default:"false"`
}

// LogWriterConfig holds configuration for a log writer.
// swagger:model
// +kubebuilder:object:generate=true
type LogWriterConfig struct {
	// Output file for logs. Keywords allowed - [`stderr`, `default`]. `default` maps to `/var/log/fluxninja/<service>.log`
	File string `json:"file" default:"stderr"`
	// Log file max size in MB
	MaxSize int `json:"max_size" validate:"gte=0" default:"50"`
	// Max log file backups
	MaxBackups int `json:"max_backups" validate:"gte=0" default:"3"`
	// Max age in days for log files
	MaxAge int `json:"max_age" validate:"gte=0" default:"7"`
	// Compress
	Compress bool `json:"compress" default:"false"`
}

// LoggerConstructor holds fields used to create an annotated instance of a logger.
type LoggerConstructor struct {
	// Name of logger instance
	Name string
	// Config key
	ConfigKey string
	// Default Config
	DefaultConfig LogConfig
	// Global logger
	IsGlobal bool
}

// Annotate creates an annotated instance of loggers which can be used to create multiple loggers.
func (constructor LoggerConstructor) Annotate() fx.Option {
	var name, group string
	name = NameTag(constructor.Name)

	if constructor.Name == "" {
		group = GroupTag("main-logger")
	} else {
		group = GroupTag(constructor.Name)
	}

	return fx.Provide(
		fx.Annotate(
			constructor.provideLogger,
			fx.ParamTags(group),
			fx.ResultTags(name),
		),
	)
}

// NewLogger creates a new logger instances with the provided config, and a list of additional writers.
func NewLogger(config LogConfig, w []io.Writer, isGlobalLogger bool) (*log.Logger, io.Writer) {
	var writers []io.Writer
	// append file writers
	for _, writerConfig := range config.Writers {
		var writer io.Writer
		if writerConfig.File != "" {
			switch writerConfig.File {
			case stdErrFile:
				writer = os.Stderr
			case stdOutFile:
				writer = os.Stdout
			default:
				if writerConfig.File == defaultFile {
					writerConfig.File = defaultLogFilePath
				}
				lj := &lumberjack.Logger{
					Filename:   writerConfig.File,
					MaxSize:    writerConfig.MaxSize,
					MaxBackups: writerConfig.MaxBackups,
					MaxAge:     writerConfig.MaxAge,
					Compress:   writerConfig.Compress,
				}
				// Set finalizer to automatically close file writers
				runtime.SetFinalizer(lj, func(lj *lumberjack.Logger) {
					log.Debug().Msg("Closing lumberjack file writer")
					_ = lj.Close()
				})
				writer = lj
			}
			writers = append(writers, writer)
		}
	}

	if config.PrettyConsole {
		writers = append(writers, log.GetPrettyConsoleWriter())
	}

	// append writers provided via Fx if they are not nil
	for _, w := range w {
		if w != nil {
			writers = append(writers, w)
		}
	}

	multi := zerolog.MultiLevelWriter(writers...)

	var wr io.Writer

	if config.NonBlocking {
		// Use diode writer
		dr := diode.NewWriter(multi, 1000, 0, func(missed int) {
			log.Printf("Dropped %d messages", missed)
		})
		wr = dr
	} else {
		// use sync writer
		wr = zerolog.SyncWriter(multi)
	}

	logger := log.NewLoggerWithWriters(wr, writers, strings.ToLower(config.LogLevel))

	if isGlobalLogger {
		// set global logger
		log.SetGlobalLogger(logger)
		// set standard loggers
		log.SetStdLogger(logger)
	}

	return logger, wr
}

func (constructor LoggerConstructor) provideLogger(w []io.Writer,
	unmarshaller Unmarshaller,
	lifecycle fx.Lifecycle,
) (*log.Logger, error) {
	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Panic().Err(err).Msg("Unable to deserialize log configuration!")
	}

	logger, writer := NewLogger(config, w, constructor.IsGlobal)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			panichandler.Go(func() {
				log.WaitFlush()
				// close diode writer
				if config.NonBlocking {
					if dr, ok := writer.(diode.Writer); ok {
						dr.Close()
					}
				}
				logger.CloseWriters()
			})
			return nil
		},
	})

	return logger, nil
}

// FxLogger overrides fx default logger.
func FxLogger() func(lg *log.Logger) fxevent.Logger {
	return func(lg *log.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: zap.New(log.NewZapAdapter(lg, "fx"))}
	}
}
