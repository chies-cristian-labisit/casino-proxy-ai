package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

// LogFormat represents the output format for logs.
type LogFormat string

const (
	// LogFormatJSON outputs structured logs as JSON (suitable for log aggregators).
	LogFormatJSON LogFormat = "json"
	// LogFormatText outputs human-readable logs for development.
	LogFormatText LogFormat = "text"
)

// StructuredLogger is a wrapper around slog.Logger that adds traceId and context fields.
type StructuredLogger struct {
	logger *slog.Logger
	format LogFormat
}

// LogEntry represents a single structured log entry with all required fields.
type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	TraceId   string      `json:"traceId"`
	Context   interface{} `json:"context,omitempty"`
	DDTraceID string      `json:"dd.trace_id,omitempty"`
	DDSpanID  string      `json:"dd.span_id,omitempty"`
}

// New creates a new StructuredLogger with the given slog.Logger and format.
func New(logger *slog.Logger, format LogFormat) *StructuredLogger {
	return &StructuredLogger{
		logger: logger,
		format: format,
	}
}

// Debug logs a debug-level message with optional context.
// context can be: map[string]interface{}, struct, slice, or any serializable type.
func (sl *StructuredLogger) Debug(ctx context.Context, msg string, context interface{}) {
	entry := sl.buildLogEntry(ctx, slog.LevelDebug, msg, context)
	sl.logEntry(entry, slog.LevelDebug)
}

// Info logs an info-level message with optional context.
// context can be: map[string]interface{}, struct, slice, or any serializable type.
func (sl *StructuredLogger) Info(ctx context.Context, msg string, context interface{}) {
	entry := sl.buildLogEntry(ctx, slog.LevelInfo, msg, context)
	sl.logEntry(entry, slog.LevelInfo)
}

// Warn logs a warn-level message with optional context.
// context can be: map[string]interface{}, struct, slice, or any serializable type.
func (sl *StructuredLogger) Warn(ctx context.Context, msg string, context interface{}) {
	entry := sl.buildLogEntry(ctx, slog.LevelWarn, msg, context)
	sl.logEntry(entry, slog.LevelWarn)
}

// Error logs an error-level message with optional context.
// context can be: map[string]interface{}, struct, slice, or any serializable type.
func (sl *StructuredLogger) Error(ctx context.Context, msg string, context interface{}) {
	entry := sl.buildLogEntry(ctx, slog.LevelError, msg, context)
	sl.logEntry(entry, slog.LevelError)
}

// buildLogEntry constructs a LogEntry from context, level, message, and fields.
func (sl *StructuredLogger) buildLogEntry(ctx context.Context, level slog.Level, msg string, context interface{}) LogEntry {
	traceId := extractTraceId(ctx)
	ddTraceID, ddSpanID := extractDDSpanIDs(ctx)
	return LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level.String(),
		Message:   msg,
		TraceId:   traceId,
		Context:   context,
		DDTraceID: ddTraceID,
		DDSpanID:  ddSpanID,
	}
}

// extractDDSpanIDs returns Datadog trace and span IDs from the active span in ctx.
// Returns empty strings when no span is active (e.g. DD_ENABLED=false).
func extractDDSpanIDs(ctx context.Context) (traceID, spanID string) {
	if ctx == nil {
		return "", ""
	}
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return "", ""
	}
	return span.Context().TraceID(),
		fmt.Sprintf("%d", span.Context().SpanID())
}

// logEntry outputs the LogEntry in the configured format.
func (sl *StructuredLogger) logEntry(entry LogEntry, level slog.Level) {
	switch sl.format {
	case LogFormatJSON:
		sl.logJSON(entry, level)
	case LogFormatText:
		sl.logText(entry, level)
	default:
		sl.logJSON(entry, level)
	}
}

// logJSON outputs the LogEntry as JSON.
func (sl *StructuredLogger) logJSON(entry LogEntry, level slog.Level) {
	data, err := json.Marshal(entry)
	if err != nil {
		sl.logger.Error("failed to marshal log entry", "error", err)
		return
	}
	sl.logger.Log(context.Background(), level, string(data))
}

// logText outputs the LogEntry in human-readable format.
// Format: 2026-05-13 10:30:45.123 [a1b2c3d4] INFO  User created successfully {userId: 12345, email: user@example.com}
func (sl *StructuredLogger) logText(entry LogEntry, level slog.Level) {
	timestamp := entry.Timestamp
	// Extract short timestamp: 2026-05-13 10:30:45.123
	if len(timestamp) >= 23 {
		timestamp = timestamp[:23]
	}

	// Extract short traceId: first 8 chars
	traceIdShort := entry.TraceId
	if len(traceIdShort) > 8 {
		traceIdShort = traceIdShort[:8]
	}

	// Format context as key=value pairs
	contextStr := ""
	if entry.Context != nil {
		contextStr = fmt.Sprintf(" %v", entry.Context)
	}

	// Build final message
	msg := fmt.Sprintf("%s [%s] %-5s %s%s", timestamp, traceIdShort, entry.Level, entry.Message, contextStr)
	sl.logger.Log(context.Background(), level, msg)
}

// extractTraceId extracts the traceId from the context.
// Returns the traceId value or an empty string if not found.
func extractTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId, ok := ctx.Value("traceId").(string)
	if !ok {
		return ""
	}
	return traceId
}

// WithTraceId returns a new context with the traceId set.
func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, "traceId", traceId)
}
