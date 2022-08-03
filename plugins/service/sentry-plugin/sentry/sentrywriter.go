package sentry

import (
	"encoding/json"
	"io"
	"time"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"github.com/fluxninja/aperture/pkg/status"
)

var _ = io.WriteCloser(new(SentryWriter))

const (
	SentryFlushWait = "5000ms"
)

var zerologToSentryLevel = map[zerolog.Level]sentry.Level{
	log.DebugLevel: sentry.LevelDebug,
	log.InfoLevel:  sentry.LevelInfo,
	log.WarnLevel:  sentry.LevelWarning,
	log.ErrorLevel: sentry.LevelError,
	log.FatalLevel: sentry.LevelFatal,
	log.PanicLevel: sentry.LevelFatal,
}

type SentryWriter struct {
	Client         *sentry.Client
	Levels         map[zerolog.Level]struct{}
	CrashWriter    *CrashWriter
	StatusRegistry *status.Registry
}

func (s *SentryWriter) Write(data []byte) (int, error) {
	event, ok := s.parseLogEvent(data)

	if ok {
		if event.Level == sentry.LevelFatal {
			s.Client.CaptureEvent(event, nil, nil)
			_ = s.Close()
		} else {
			_, _ = s.CrashWriter.Write(data)
		}
	}
	return len(data), nil
}

func (s *SentryWriter) parseLogEvent(data []byte) (*sentry.Event, bool) {
	levelStr, err := jsonparser.GetUnsafeString(data, zerolog.LevelFieldName)
	if err != nil {
		return nil, false
	}

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		return nil, false
	}

	_, enabled := s.Levels[level]
	if !enabled {
		return nil, false
	}

	sentryLevel, ok := zerologToSentryLevel[level]
	if !ok {
		return nil, false
	}

	event := sentry.Event{
		Timestamp: time.Now(),
		Level:     sentryLevel,
		Logger:    "zerolog",
	}

	err = jsonparser.ObjectEach(data, func(key, value []byte, _ jsonparser.ValueType, _ int) error {
		switch string(key) {
		case zerolog.MessageFieldName:
			event.Message = bytesToStrUnsafe(value)
		case zerolog.ErrorFieldName:
			event.Exception = append(event.Exception, sentry.Exception{
				Value: bytesToStrUnsafe(value),
				// TODO: Create stacktrace which holds information about the frames of the stack if needed
				Stacktrace: nil,
			})
		}
		return nil
	})

	if err != nil {
		return nil, false
	}

	return &event, true
}

func (s *SentryWriter) Close() error {
	duration, _ := time.ParseDuration(SentryFlushWait)
	s.Client.Flush(duration)
	return nil
}

func bytesToStrUnsafe(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

// func RegisterSentryPanicHandler() {}

func (s *SentryWriter) SentryPanicHandler(e interface{}, _ panichandler.Callstack) {
	duration, _ := time.ParseDuration(SentryFlushWait)

	// Crash Log
	var crashLogData map[string]interface{}
	crashLog := s.CrashWriter.GetCrashLog()
	err := json.Unmarshal(crashLog, &crashLogData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal crash log")
	}

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "Crash report",
		Level:    sentry.LevelInfo,
		Data:     crashLogData,
	})

	// Dump Status Registry
	var statusData map[string]interface{}
	status := s.StatusRegistry.Get("")
	if status != nil {
		groupStatus, err := json.MarshalIndent(status, "", " ")
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal status registry")
		}

		err = json.Unmarshal(groupStatus, &statusData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal status registry")
		}

		sentry.AddBreadcrumb(&sentry.Breadcrumb{
			Category: "Status Registry",
			Level:    sentry.LevelInfo,
			Data:     statusData,
		})
	}

	// Service Version Information
	var versionData map[string]interface{}
	versionInfo := info.GetVersionInfo()
	if versionInfo != nil {
		vInfo, err := json.MarshalIndent(versionInfo, "", " ")
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal version information")
		}

		err = json.Unmarshal(vInfo, &versionData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal version information")
		}

		sentry.AddBreadcrumb(&sentry.Breadcrumb{
			Category: "Version Information",
			Level:    sentry.LevelInfo,
			Data:     versionData,
		})
	}

	sentry.CurrentHub().Recover(e)
	sentry.Flush(duration)
}
