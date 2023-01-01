package infrastructure

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(message string)
	Error(err error)
	Fatal(err error)
	With(f Field) Logger
}

type logger struct {
	zapLogger *zap.Logger
}

type Field struct {
	Key   string
	Value string
}

func NewLogger() Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	zapConfig := zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stack_trace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zapLogger, _ := zapConfig.Build(zap.AddCallerSkip(1))

	logger := &logger{
		zapLogger: zapLogger,
	}

	return logger
}

func (l *logger) Info(message string) {
	l.zapLogger.Info(message)
}

func (l *logger) Error(err error) {
	l.zapLogger.Error(
		err.Error(),
		zap.Error(err),
	)
}

func (l *logger) Fatal(err error) {
	l.zapLogger.Fatal(
		err.Error(),
		zap.Error(err),
	)
}

func (l *logger) With(f Field) Logger {
	field := zap.String(f.Key, f.Value)
	return &logger{
		zapLogger: l.zapLogger.With(field),
	}
}
