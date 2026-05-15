package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type LogLevel slog.Level

const (
	DebugLevel LogLevel = LogLevel(slog.LevelDebug)
	InfoLevel  LogLevel = LogLevel(slog.LevelInfo)
	WarnLevel  LogLevel = LogLevel(slog.LevelWarn)
	ErrorLevel LogLevel = LogLevel(slog.LevelError)
)

// ParseLogLevel reads LOG_LEVEL env var or config file, returns slog.Level
// Precedence: LOG_LEVEL env var > log_level config file > default (INFO)
// Returns error if log level is invalid.
func ParseLogLevel() (slog.Level, error) {
	// 1. Check environment variable first (highest precedence)
	envLevel := os.Getenv("LOG_LEVEL")
	if envLevel != "" {
		return parseLogLevelString(envLevel)
	}

	// 2. Check config file (if we add config file support later)
	// For now, only env var and default

	// 3. Default: INFO
	return slog.LevelInfo, nil
}

func parseLogLevelString(levelStr string) (slog.Level, error) {
	levelStr = strings.ToUpper(strings.TrimSpace(levelStr))
	switch levelStr {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid LOG_LEVEL: %s (must be DEBUG, INFO, WARN, or ERROR)", levelStr)
	}
}

// NewLogger creates a structured logger with the configured log level.
// Uses JSON format for production, text format for development.
// Automatically detects format from APP_ENV environment variable.
func NewLogger(level slog.Level) *slog.Logger {
	appEnv := os.Getenv("APP_ENV")
	var handler slog.Handler

	if appEnv == "production" {
		// JSON format for production (parseable by log aggregators).
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Text format for development (human-readable).
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}
