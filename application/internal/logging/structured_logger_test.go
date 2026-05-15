package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

func TestStructuredLoggerJSONFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := New(baseLogger, LogFormatJSON)
	ctx := WithTraceId(context.Background(), "test-trace-123")

	fields := map[string]interface{}{
		"userId": 12345,
		"email":  "user@example.com",
	}

	logger.Info(ctx, "User created successfully", fields)

	// Verify something was logged.
	if buf.Len() == 0 {
		t.Error("Expected log output but got none")
	}
}

func TestStructuredLoggerWithStruct(t *testing.T) {
	type Customer struct {
		ID   uint   `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}

	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := New(baseLogger, LogFormatJSON)
	ctx := WithTraceId(context.Background(), "test-trace-struct")

	customer := Customer{
		ID:   12345,
		Code: "CODE123",
		Name: "John Doe",
	}

	logger.Info(ctx, "Customer created", customer)

	// Verify the output contains the customer data.
	if buf.Len() == 0 {
		t.Error("Expected log output but got none")
	}
}

func TestStructuredLoggerWithMap(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := New(baseLogger, LogFormatJSON)
	ctx := WithTraceId(context.Background(), "test-trace-map")

	context := map[string]interface{}{
		"userId": 12345,
		"action": "created",
	}

	logger.Info(ctx, "User action logged", context)

	// Verify the output.
	if buf.Len() == 0 {
		t.Error("Expected log output but got none")
	}
}

func TestTraceIdExtraction(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "traceId present in context",
			ctx:      WithTraceId(context.Background(), "test-trace-456"),
			expected: "test-trace-456",
		},
		{
			name:     "traceId missing from context",
			ctx:      context.Background(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTraceId(tt.ctx)
			if result != tt.expected {
				t.Errorf("Expected traceId %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWithTraceId(t *testing.T) {
	ctx := context.Background()
	traceId := "unique-trace-789"

	newCtx := WithTraceId(ctx, traceId)
	extracted := extractTraceId(newCtx)

	if extracted != traceId {
		t.Errorf("WithTraceId failed: expected %q, got %q", traceId, extracted)
	}
}

func TestLogLevelStrings(t *testing.T) {
	levels := []struct {
		level  slog.Level
		name   string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARN"},
		{slog.LevelError, "ERROR"},
	}

	for _, tt := range levels {
		if tt.level.String() != tt.name {
			t.Errorf("Level string mismatch: expected %q, got %q", tt.name, tt.level.String())
		}
	}
}

func TestExtractDDSpanIDs_NoActiveSpan(t *testing.T) {
	traceID, spanID := extractDDSpanIDs(context.Background())
	if traceID != "" || spanID != "" {
		t.Errorf("expected empty IDs when no span active, got traceID=%q spanID=%q", traceID, spanID)
	}
}

func TestExtractDDSpanIDs_NilContext(t *testing.T) {
	traceID, spanID := extractDDSpanIDs(nil)
	if traceID != "" || spanID != "" {
		t.Errorf("expected empty IDs for nil context, got traceID=%q spanID=%q", traceID, spanID)
	}
}

func TestExtractDDSpanIDs_ActiveSpan(t *testing.T) {
	tracer.Start(tracer.WithService("test"))
	defer tracer.Stop()

	span, ctx := tracer.StartSpanFromContext(context.Background(), "test.op")
	defer span.Finish()

	traceID, spanID := extractDDSpanIDs(ctx)
	if traceID == "" {
		t.Error("expected non-empty dd.trace_id when span is active")
	}
	if spanID == "" {
		t.Error("expected non-empty dd.span_id when span is active")
	}
}

func TestStructuredLogger_DDFieldsInjectedWhenSpanActive(t *testing.T) {
	tracer.Start(tracer.WithService("test"))
	defer tracer.Stop()

	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := New(slog.New(handler), LogFormatJSON)

	span, ctx := tracer.StartSpanFromContext(context.Background(), "test.op")
	defer span.Finish()

	logger.Info(ctx, "test message", nil)

	out := buf.String()
	if !strings.Contains(out, "dd.trace_id") {
		t.Error("expected dd.trace_id in log output when span is active")
	}
	if !strings.Contains(out, "dd.span_id") {
		t.Error("expected dd.span_id in log output when span is active")
	}
}

func TestStructuredLogger_DDFieldsAbsentWhenNoSpan(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := New(slog.New(handler), LogFormatJSON)

	logger.Info(context.Background(), "test message", nil)

	out := buf.String()
	if strings.Contains(out, "dd.trace_id") {
		t.Error("dd.trace_id should be absent in log output when no span is active")
	}
	if strings.Contains(out, "dd.span_id") {
		t.Error("dd.span_id should be absent in log output when no span is active")
	}
}

func TestStructuredLoggerNilContext(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := New(baseLogger, LogFormatJSON)

	// Should not panic with nil context.
	logger.Info(nil, "Test message", nil)

	// Verify something was logged.
	if buf.Len() == 0 {
		t.Error("Expected log output but got none")
	}
}

func TestStructuredLoggerAllLogLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := New(baseLogger, LogFormatJSON)
	ctx := WithTraceId(context.Background(), "test-trace")

	// Test each method doesn't panic.
	logger.Debug(ctx, "Debug message", nil)
	logger.Info(ctx, "Info message", nil)
	logger.Warn(ctx, "Warn message", nil)
	logger.Error(ctx, "Error message", nil)

	// Verify logs were written.
	if buf.Len() == 0 {
		t.Error("Expected log output but got none")
	}
}
