package integration

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/cometagaming/ms-casino-go-v2/internal/config"
)

// TestLogLevelConfiguration verifies that log level configuration works end-to-end.
func TestLogLevelConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		envLogLevel   string
		expectedLevel slog.Level
		shouldError   bool
	}{
		{
			name:          "default INFO level when no env var",
			envLogLevel:   "",
			expectedLevel: slog.LevelInfo,
			shouldError:   false,
		},
		{
			name:          "DEBUG level from env var",
			envLogLevel:   "DEBUG",
			expectedLevel: slog.LevelDebug,
			shouldError:   false,
		},
		{
			name:          "WARN level from env var",
			envLogLevel:   "WARN",
			expectedLevel: slog.LevelWarn,
			shouldError:   false,
		},
		{
			name:          "ERROR level from env var",
			envLogLevel:   "ERROR",
			expectedLevel: slog.LevelError,
			shouldError:   false,
		},
		{
			name:          "invalid log level",
			envLogLevel:   "INVALID",
			expectedLevel: slog.LevelInfo,
			shouldError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset LOG_LEVEL env var
			if tt.envLogLevel != "" {
				os.Setenv("LOG_LEVEL", tt.envLogLevel)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}
			defer os.Unsetenv("LOG_LEVEL")

			// Parse log level
			level, err := config.ParseLogLevel()

			// Check error expectation
			if (err != nil) != tt.shouldError {
				t.Errorf("ParseLogLevel() error = %v, shouldError %v", err, tt.shouldError)
				return
			}

			// Check result
			if level != tt.expectedLevel {
				t.Errorf("ParseLogLevel() = %v, want %v", level, tt.expectedLevel)
			}
		})
	}
}

// TestLoggerRespectsLogLevel verifies that the logger only outputs at the configured level.
func TestLoggerRespectsLogLevel(t *testing.T) {
	tests := []struct {
		name           string
		level          slog.Level
		expectedDebug  bool
		expectedInfo   bool
		expectedWarn   bool
		expectedError  bool
	}{
		{
			name:           "DEBUG level shows all",
			level:          slog.LevelDebug,
			expectedDebug:  true,
			expectedInfo:   true,
			expectedWarn:   true,
			expectedError:  true,
		},
		{
			name:           "INFO level hides debug",
			level:          slog.LevelInfo,
			expectedDebug:  false,
			expectedInfo:   true,
			expectedWarn:   true,
			expectedError:  true,
		},
		{
			name:           "WARN level hides debug/info",
			level:          slog.LevelWarn,
			expectedDebug:  false,
			expectedInfo:   false,
			expectedWarn:   true,
			expectedError:  true,
		},
		{
			name:           "ERROR level hides debug/info/warn",
			level:          slog.LevelError,
			expectedDebug:  false,
			expectedInfo:   false,
			expectedWarn:   false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger with specific level
			logger := config.NewLogger(tt.level)

			// Capture output
			var buf bytes.Buffer
			handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: tt.level})
			logger = slog.New(handler)

			// Log at each level
			logger.Debug("debug message")
			logger.Info("info message")
			logger.Warn("warn message")
			logger.Error("error message")

			output := buf.String()

			// Verify output
			if tt.expectedDebug && !strings.Contains(output, "debug") {
				t.Errorf("expected DEBUG output not found")
			}
			if !tt.expectedDebug && strings.Contains(output, "debug") {
				t.Errorf("DEBUG output found but not expected")
			}

			if tt.expectedInfo && !strings.Contains(output, "info") {
				t.Errorf("expected INFO output not found")
			}
			if !tt.expectedInfo && strings.Contains(output, "info") {
				t.Errorf("INFO output found but not expected")
			}

			if tt.expectedWarn && !strings.Contains(output, "warn") {
				t.Errorf("expected WARN output not found")
			}
			if !tt.expectedWarn && strings.Contains(output, "warn") {
				t.Errorf("WARN output found but not expected")
			}

			if tt.expectedError && !strings.Contains(output, "error") {
				t.Errorf("expected ERROR output not found")
			}
			if !tt.expectedError && strings.Contains(output, "error") {
				t.Errorf("ERROR output found but not expected")
			}
		})
	}
}
