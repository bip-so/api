package sentry

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func sentryLevel(lvl zapcore.Level) sentry.Level {
	switch lvl {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel:
		return sentry.LevelFatal
	case zapcore.PanicLevel:
		return sentry.LevelFatal
	case zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelFatal
	}
}

// Sentry specific options. Whenever zap log is logged this method will be trigerred and log event is sent to sentry.
func SentryOptions() zap.Option {
	sampleOptions := zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.InfoLevel {
			return nil
		}
		event := sentry.NewEvent()
		event.Message = fmt.Sprintf("%s, Line No: %d :: %s", entry.Caller.File, entry.Caller.Line, entry.Message)
		event.Timestamp = entry.Time
		event.Level = sentryLevel(entry.Level)
		event.Platform = "go"
		trace := sentry.NewStacktrace()
		if trace != nil {
			event.Exception = []sentry.Exception{{
				Type:       entry.Message,
				Value:      entry.Caller.TrimmedPath(),
				Stacktrace: trace,
			}}
		}
		sentry.CaptureEvent(event)
		if entry.Level > zapcore.ErrorLevel {
			defer sentry.Flush(2 * time.Second)
		}
		return nil
	})
	return sampleOptions
}
