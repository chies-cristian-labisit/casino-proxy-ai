package config

import (
	"log/slog"
	"os"
	"testing"
)

func TestParseLogLevelFromEnvVar(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		expected  slog.Level
		wantError bool
	}{
		{
			name:      "DEBUG level",
			envValue:  "DEBUG",
			expected:  slog.LevelDebug,
			wantError: false,
		},
		{
			name:      "INFO level",
			envValue:  "INFO",
			expected:  slog.LevelInfo,
			wantError: false,
		},
		{
			name:      "WARN level",
			envValue:  "WARN",
			expected:  slog.LevelWarn,
			wantError: false,
		},
		{
			name:      "ERROR level",
			envValue:  "ERROR",
			expected:  slog.LevelError,
			wantError: false,
		},
		{
			name:      "lowercase debug",
			envValue:  "debug",
			expected:  slog.LevelDebug,
			wantError: false,
		},
		{
			name:      "mixed case info",
			envValue:  "InFo",
			expected:  slog.LevelInfo,
			wantError: false,
		},
		{
			name:      "invalid level",
			envValue:  "TRACE",
			expected:  slog.LevelInfo,
			wantError: true,
		},
		{
			name:      "empty string defaults to INFO",
			envValue:  "",
			expected:  slog.LevelInfo,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env var
			if tt.envValue != "" {
				os.Setenv("LOG_LEVEL", tt.envValue)
				defer os.Unsetenv("LOG_LEVEL")
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			// Parse
			level, err := ParseLogLevel()

			// Check error
			if (err != nil) != tt.wantError {
				t.Errorf("ParseLogLevel() error = %v, wantError %v", err, tt.wantError)
			}

			// Check result
			if level != tt.expected {
				t.Errorf("ParseLogLevel() = %v, want %v", level, tt.expected)
			}
		})
	}
}

func TestParseLogLevelString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  slog.Level
		wantError bool
	}{
		{
			name:      "valid DEBUG",
			input:     "DEBUG",
			expected:  slog.LevelDebug,
			wantError: false,
		},
		{
			name:      "valid INFO",
			input:     "INFO",
			expected:  slog.LevelInfo,
			wantError: false,
		},
		{
			name:      "valid WARN",
			input:     "WARN",
			expected:  slog.LevelWarn,
			wantError: false,
		},
		{
			name:      "valid ERROR",
			input:     "ERROR",
			expected:  slog.LevelError,
			wantError: false,
		},
		{
			name:      "invalid INVALID",
			input:     "INVALID",
			expected:  slog.LevelInfo,
			wantError: true,
		},
		{
			name:      "with whitespace",
			input:     "  DEBUG  ",
			expected:  slog.LevelDebug,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevelString(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("parseLogLevelString() error = %v, wantError %v", err, tt.wantError)
			}

			if level != tt.expected {
				t.Errorf("parseLogLevelString() = %v, want %v", level, tt.expected)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
	}{
		{
			name:  "DEBUG level logger",
			level: slog.LevelDebug,
		},
		{
			name:  "INFO level logger",
			level: slog.LevelInfo,
		},
		{
			name:  "WARN level logger",
			level: slog.LevelWarn,
		},
		{
			name:  "ERROR level logger",
			level: slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.level)
			if logger == nil {
				t.Error("NewLogger() returned nil")
			}
		})
	}
}
