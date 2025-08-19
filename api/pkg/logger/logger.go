package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// Init initializes the logger and sets it as the default logger.
func Init() {
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(defaultLogger)
}

// GetLogger returns the default logger instance.
func GetLogger() *slog.Logger {
	return defaultLogger
}
