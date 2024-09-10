package logger

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
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
			AddSource: false, // This adds the filename and line number
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
		AddSource: false, // Add source for filename and line number
	}))
}

func getCallerInfo() string {
	// We skip 2 frames to get the caller of the function using the logger (skips this function and the logger method itself)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown source"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// GetLogger returns the logger instance
func GetLogger() *Logger {
	return l
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keysAndValues ...any) {
	callerInfo := getCallerInfo()
	l.logger.Debug(fmt.Sprintf("[%s] %s", callerInfo, msg), keysAndValues...)
}

// Info logs an info message
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	callerInfo := getCallerInfo()
	l.logger.Info(fmt.Sprintf("[%s] %s", callerInfo, msg), keysAndValues...)
}

// Error logs an error message with the correct source
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	callerInfo := getCallerInfo()
	l.logger.Error(fmt.Sprintf("[%s] %s", callerInfo, msg), keysAndValues...)
}
