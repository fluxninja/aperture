package log

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapAdapter struct {
	logger *zerolog.Logger
}

func (z *zapAdapter) Write(p []byte) (int, error) {
	// parse zap's json log format and output it to zerolog
	var fields map[string]interface{}
	if err := json.Unmarshal(p, &fields); err != nil {
		return 0, err
	}

	// build context
	context := z.logger.With()

	// read level and remove it from fields
	level, ok := fields[zerolog.LevelFieldName]
	if !ok {
		return 0, errors.New("level not found")
	}
	if level == "dpanic" {
		level = "panic"
	}
	zerologLevel, err := zerolog.ParseLevel(level.(string))
	if err != nil {
		// not a well-known level so we user zerolog.NoLevel
		zerologLevel = zerolog.NoLevel
		// add level to context
		context = context.Str(zerolog.LevelFieldName, level.(string))
	}
	// remove level from fields
	delete(fields, zerolog.LevelFieldName)

	// read message and remove it from fields
	message, ok := fields[zerolog.MessageFieldName]
	if !ok {
		return 0, errors.New("message not found")
	}
	// remove message from fields
	delete(fields, zerolog.MessageFieldName)

	for k, v := range fields {
		context = context.Interface(k, v)
	}
	logger := context.Logger()

	logger.WithLevel(zerologLevel).CallerSkipFrame(4).Msg(message.(string))

	return len(p), nil
}

// NewZapAdapter returns a new zapcore Core logger interface that converts given logger's io.writer to WriteSyncer.
func NewZapAdapter(logger *Logger, component string) zapcore.Core {
	currentLevel := logger.GetLevel().String()
	// zap has no "trace" level, so we use "debug" instead
	if currentLevel == "trace" {
		currentLevel = "debug"
	}
	level, err := zapcore.ParseLevel(currentLevel)
	if err != nil {
		log.Panic().Err(err).Str("level", level.String()).Msg("Unable to parse logger level")
	}
	encodercfg := zap.NewProductionEncoderConfig()
	encodercfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encodercfg.TimeKey = zerolog.TimestampFieldName
	encodercfg.LevelKey = zerolog.LevelFieldName
	encodercfg.MessageKey = zerolog.MessageFieldName
	encodercfg.StacktraceKey = zerolog.ErrorStackFieldName

	encodercfg.CallerKey = zerolog.CallerFieldName
	encodercfg.EncodeCaller = zapcore.ShortCallerEncoder

	// duplicate zerolog logger and remove it is context
	childLogger := logger.logger.With().Logger()
	childLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		// reset context
		return zerolog.Context{}
	})
	adapter := &zapAdapter{
		logger: &childLogger,
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encodercfg),
		zapcore.AddSync(adapter),
		level,
	).With([]zap.Field{zap.String(componentKey, component)})
}
