package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

var l *Logger

// InitLogger initializes the logger with a text handler
func InitLogger() {
	l = &Logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true, // This adds the filename and line number
		})),
	}
}

func SetLogLevel(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	default:
		logLevel = slog.LevelInfo
	}

	l.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true, // Add source for filename and line number
	}))
}

// GetLogger returns the logger instance
func GetLogger() *Logger {
	return l
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, keysAndValues...)
}

// Info logs an info message
func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

// Error logs an error message
func (l *Logger) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, keysAndValues...)
}
