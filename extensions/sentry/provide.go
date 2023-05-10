package sentry

import (
	"context"
	"io"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/getsentry/sentry-go"

	sentryconfig "github.com/fluxninja/aperture/v2/extensions/sentry/config"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/panichandler"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// sentryWriterConstructor holds fields to create an annotated instance of Sentry Writer.
type sentryWriterConstructor struct {
	// sentryWriter
	sentryWriter *sentryWriter
	// Name of sentry instance
	Name string
	// Config key
	ConfigKey string
	// Default Config
	DefaultConfig sentryconfig.SentryConfig
}

// annotate creates an annotated instance of SentryWriter.
func (constructor *sentryWriterConstructor) annotate() fx.Option {
	var group string
	if constructor.Name == "" {
		group = config.GroupTag("main-logger")
	} else {
		group = config.GroupTag(constructor.Name)
	}
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor.provideSentryWriter,
				fx.ResultTags(group),
			),
		),
		fx.Invoke(constructor.setupSentryWriter),
	)
}

func (constructor *sentryWriterConstructor) provideSentryWriter(
	unmarshaller config.Unmarshaller,
	lifecycle fx.Lifecycle,
) (io.Writer, error) {
	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.ConfigKey, &config); err != nil {
		log.Panic().Err(err).Msg("Unable to deserialize sentry config")
	}

	if config.Disabled {
		log.Info().Msg("Sentry crash report disabled")
		return nil, nil
	}

	sentryWriter, _ := newSentryWriter(config)

	constructor.sentryWriter = sentryWriter

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			sentry.CurrentHub().BindClient(sentryWriter.client)
			return nil
		},
		OnStop: func(_ context.Context) error {
			duration, _ := time.ParseDuration(SentryFlushWait)
			sentry.Flush(duration)
			return nil
		},
	})

	return sentryWriter, nil
}

func (constructor *sentryWriterConstructor) setupSentryWriter(registry status.Registry) {
	if constructor.sentryWriter != nil {
		constructor.sentryWriter.statusRegistry = registry
	}
}

// newSentryWriter creates a new SentryWriter instance with Sentry Client and registers panic handler.
func newSentryWriter(config sentryconfig.SentryConfig) (*sentryWriter, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              config.Dsn,
		Debug:            config.Debug,
		Environment:      config.Environment,
		Release:          info.Version,
		AttachStacktrace: config.AttachStacktrace,
		SampleRate:       config.SampleRate,
		TracesSampleRate: config.TracesSampleRate,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sentry client")
		return nil, err
	}

	reportLevels := []zerolog.Level{
		zerolog.DebugLevel,
		zerolog.InfoLevel,
		zerolog.WarnLevel,
		zerolog.ErrorLevel,
		zerolog.FatalLevel,
		zerolog.PanicLevel,
	}

	levels := make(map[zerolog.Level]struct{}, len(reportLevels))
	for _, level := range reportLevels {
		levels[level] = struct{}{}
	}

	crashWriter := newCrashWriter(logCountLimit)
	sentryWriter := &sentryWriter{
		client:      client,
		levels:      levels,
		crashWriter: crashWriter,
	}

	panichandler.RegisterPanicHandler(sentryWriter.sentryPanicHandler)
	return sentryWriter, nil
}
