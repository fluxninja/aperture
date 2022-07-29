package sentry

import (
	"io"
	"time"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"

	"github.com/fluxninja/aperture/pkg/log"
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
	Client *sentry.Client
	Levels map[zerolog.Level]struct{}
}

func (s *SentryWriter) Write(data []byte) (int, error) {
	event, ok := s.parseLogEvent(data)
	if ok {
		s.Client.CaptureEvent(event, nil, nil)
		if event.Level == sentry.LevelFatal {
			_ = s.Close()
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

	err = jsonparser.ObjectEach(data, func(key, value []byte, vt jsonparser.ValueType, offset int) error {
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
