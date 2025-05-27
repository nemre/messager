package logger

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

type Logger interface {
	Debug(message string, arguments ...any)
	Info(message string, arguments ...any)
	Warning(message string, err error, arguments ...any)
	Error(message string, err error, arguments ...any)
	Fatal(message string, err error, arguments ...any)
	FatalWithoutExit(message string, err error, arguments ...any)
}

type logger struct {
	stdoutLogger *slog.Logger
	stderrLogger *slog.Logger
}

func New() Logger {
	handlerOptions := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(_ []string, attribute slog.Attr) slog.Attr {
			value := reflect.ValueOf(attribute.Value.Any())

			if !value.IsValid() || value.IsZero() {
				return slog.Attr{
					Key:   "",
					Value: slog.Value{},
				}
			}

			if attribute.Key == slog.MessageKey {
				attribute.Key = "message"
			}

			if attribute.Key != slog.LevelKey {
				return attribute
			}

			level, ok := attribute.Value.Any().(slog.Level)
			if !ok {
				return attribute
			}

			if level == slog.LevelWarn {
				attribute.Value = slog.StringValue("warning")
			}

			if level == slog.LevelError+1 {
				attribute.Value = slog.StringValue("fatal")
			}

			attribute.Value = slog.StringValue(strings.ToLower(attribute.Value.String()))

			return attribute
		},
	}

	return &logger{
		stdoutLogger: slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions)),
		stderrLogger: slog.New(slog.NewJSONHandler(os.Stderr, &handlerOptions)),
	}
}

func (l *logger) Debug(message string, arguments ...any) {
	l.stdoutLogger.Debug(message, arguments...)
}

func (l *logger) Info(message string, arguments ...any) {
	l.stdoutLogger.Info(message, arguments...)
}

func (l *logger) Warning(message string, err error, arguments ...any) {
	arguments = append(arguments, "error", err)
	l.stderrLogger.Warn(message, arguments...)
}

func (l *logger) Error(message string, err error, arguments ...any) {
	arguments = append(arguments, "error", err)
	l.stderrLogger.Error(message, arguments...)
}

func (l *logger) Fatal(message string, err error, arguments ...any) {
	arguments = append(arguments, "error", err)
	l.stderrLogger.Log(context.Background(), slog.LevelError+1, message, arguments...)

	os.Exit(1)
}

func (l *logger) FatalWithoutExit(message string, err error, arguments ...any) {
	arguments = append(arguments, "error", err)
	l.stderrLogger.Log(context.Background(), slog.LevelError+1, message, arguments...)
}
