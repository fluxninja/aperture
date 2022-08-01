package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZerologAdapter returns a new zapcore Core logger interface that converts given logger's io.writer to WriteSyncer.
func NewZerologAdapter(logger Logger) zapcore.Core {
	currentLevel := zerolog.GlobalLevel().String()
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

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encodercfg),
		zapcore.AddSync(logger.w),
		level,
	).With([]zap.Field{zap.String(serviceKey, "otelcollector")})
}
