package logger

import (
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return &Logger{Logger: logger}
}

func NewDiscardLogger() *Logger {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return &Logger{Logger: logger}
}
