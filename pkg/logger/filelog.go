package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lumlog *lumberjack.Logger

func initLumberJack() {
	lumlog = &lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    50, // megabytes
		MaxBackups: 30,
		MaxAge:     28, // days

	}
}

func lumberjackZapHook(e zapcore.Entry) error {
	lumlog.Write([]byte(fmt.Sprintf("%+v \n", e)))
	return nil
}
