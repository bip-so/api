package logger

import (
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/sentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

type BipLogger interface {
	Info(msg string, args ...zapcore.Field)
	Debug(msg string, args ...zapcore.Field)
	Warn(msg string, args ...zapcore.Field)
	Error(msg string, args ...zapcore.Field)
	Panic(msg string, args ...zapcore.Field)
	Fatal(msg string, args ...zapcore.Field)
}

func InitLogger() {
	var err error
	if configs.IsLive() {
		logger, err = zap.NewProduction(sentry.SentryOptions())
	} else if configs.IsDev() {
		logger, err = zap.NewDevelopment(sentry.SentryOptions())
	} else {
		initLumberJack()
		logger, err = zap.NewDevelopment(zap.Hooks(lumberjackZapHook))
	}
	if err != nil {
		panic(err)
	}
}

func Info(msg string, args ...zapcore.Field) {
	logger.Info(msg, args...)
}

func Debug(msg string, args ...zapcore.Field) {
	logger.Debug(msg, args...)
}

func Warn(msg string, args ...zapcore.Field) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...zapcore.Field) {
	logger.Error(msg, args...)
}

func Panic(msg string, args ...zapcore.Field) {
	logger.Panic(msg, args...)
}

func Fatal(msg string, args ...zapcore.Field) {
	logger.Fatal(msg, args...)
}
