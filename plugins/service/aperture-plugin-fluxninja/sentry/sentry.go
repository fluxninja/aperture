package sentry

import (
	"io"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/getsentry/sentry-go"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
)

const (
	// SentryConfigKey for sentry config.
	SentryConfigKey = "sentry"
)

// Configuration for Sentry
// swagger:model
type SentryConfig struct {
	// If DSN is not set, the client is effectively disabled
	// You can set test project's dsn to send log events: <https://7d7c8b13c9a44a31befe58199f59e4da@o574197.ingest.sentry.io/5758807>
	Dsn string `json:"dsn" default:""`
	// Environment
	Environment string `json:"environment" default:"production"`
	// Sample rate for sampling traces i.e. 0.0 to 1.0
	TracesSampleRate float64 `json:"traces_sample_rate" default:"0.2"`
	// Sample rate for event submission i.e. 0.0 to 1.0
	SampleRate float64 `json:"sample_rate" default:"1.0"`
	// Debug enables printing of Sentry SDK debug messages
	Debug bool `json:"debug" default:"true"`
	// Configure to generate and attach stacktraces to capturing message calls
	AttachStacktrace bool `json:"attach_stack_trace" default:"true"`
	// Sentry crash report disabled
	Disabled bool `json:"disabled" default:"false"`
}

// Construcotr.
type SentryWriterConstructor struct {
	// Name of sentry instance
	Name string
	// Config key
	Key string
	// Default Config
	DefaultConfig SentryConfig
}

// Annotate Fx provide.
func (constructor SentryWriterConstructor) Annotate() fx.Option {
	var group string
	if constructor.Name == "" {
		group = config.GroupTag("main-logger")
	} else {
		group = config.GroupTag(constructor.Name)
	}
	return fx.Provide(
		fx.Annotate(
			constructor.provideSentryWriter,
			fx.ResultTags(group),
		),
	)
}

func (constructor SentryWriterConstructor) provideSentryWriter(unmarshaller config.Unmarshaller) (io.Writer, error) {
	config := constructor.DefaultConfig

	if err := unmarshaller.UnmarshalKey(constructor.Key, &config); err != nil {
		log.Fatal().Err(err).Msg("Unable to deserialize sentry config")
	}

	if config.Disabled {
		log.Info().Msg("Sentry crash report disabled")
		return nil, nil
	}

	sentryWriter, _ := NewSentryWriter(config)

	return sentryWriter, nil
}

func NewSentryWriter(config SentryConfig) (*SentryWriter, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              config.Dsn,
		Debug:            config.Debug,
		Environment:      config.Environment,
		Release:          info.Version,
		AttachStacktrace: config.AttachStacktrace,
		SampleRate:       config.SampleRate,
		TracesSampleRate: config.TracesSampleRate,
		// TODO: Customize sampler if needed
		// TracesSampler enables to customize the sampleing of traces, overrides TracesSampleRate
		//
		// TracesSampler: sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
		// 	return sentry.SampledTrue
		// }),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sentry client")
		return nil, err
	}

	reportLevels := []zerolog.Level{
		zerolog.FatalLevel,
		zerolog.PanicLevel,
	}

	levels := make(map[zerolog.Level]struct{}, len(reportLevels))
	for _, level := range reportLevels {
		levels[level] = struct{}{}
	}

	sentryWriter := &SentryWriter{
		Client: client,
		Levels: levels,
	}

	return sentryWriter, nil
}
