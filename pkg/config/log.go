package config

import (
	"context"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/fluxninja/lumberjack"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
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
		LoggerConstructor{ConfigKey: configKey}.Annotate(),
		fx.Invoke(log.SetGlobalLogger),
		fx.Invoke(log.SetStdLogger),
		fx.WithLogger(WithApertureLogger()),
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

	// Additional log writers
	Writers []LogWriterConfig `json:"writers" validate:"omitempty,dive,omitempty"`

	// Base LogWriterConfig
	LogWriterConfig `json:",inline"`

	// Use non-blocking log writer (can lose logs at high throughput)
	NonBlocking bool `json:"non_blocking" default:"true"`

	// Additional log writer: pretty console (stdout) logging (not recommended for prod environments)
	PrettyConsole bool `json:"pretty_console" default:"false"`
}

// LogWriterConfig holds configuration for a log writer.
// swagger:model
// +kubebuilder:object:generate=true
type LogWriterConfig struct {
	// Output file for logs. Keywords allowed - ["stderr", "default"]. "default" maps to `/var/log/fluxninja/<service>.log`
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

func (constructor LoggerConstructor) provideLogger(w []io.Writer,
	unmarshaller Unmarshaller,
	lifecycle fx.Lifecycle,
) (log.Logger, error) {
	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Panic().Err(err).Msg("Unable to deserialize log configuration!")
	}
	logger, writers := NewLogger(config)
	// append additional writers provided via Fx
	writers = append(writers, w...)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			panichandler.Go(func() {
				logger.Close()
				log.WaitFlush()
				for _, writer := range writers {
					if closer, ok := writer.(io.Closer); ok {
						_ = closer.Close()
					}
				}
			})
			return nil
		},
	})

	return logger, nil
}

// NewLogger creates a new instance of logger and writers with the given configuration.
func NewLogger(config LogConfig) (log.Logger, []io.Writer) {
	var writers []io.Writer
	// append file writers
	for _, writerConfig := range config.Writers {
		var writer io.Writer
		if config.File != "" {
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
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	multi := zerolog.MultiLevelWriter(writers...)

	logger := log.NewLogger(multi, config.NonBlocking, strings.ToLower(config.LogLevel))

	logger.Info().Msg("Configured logger")

	return logger, writers
}
