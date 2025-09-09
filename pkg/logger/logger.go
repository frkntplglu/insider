package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

var global *Logger

func Init(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	global = &Logger{
		Logger: slog.New(handler),
	}
}

func Get() *Logger {
	if global == nil {
		Init(slog.LevelInfo)
	}
	return global
}

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Get().Error(msg, args...)
	os.Exit(1)
}
