package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/middleware"
	"github.com/cometagaming/ms-casino-go-v2/internal/logging"
)

func TestTraceIdMiddlewareGeneratesNewTraceId(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.TraceIdMiddleware)

	app.Get("/test", func(c *fiber.Ctx) error {
		// Verify traceId is available in context.
		ctx := c.UserContext()
		traceId, ok := ctx.Value("traceId").(string)
		if !ok {
			t.Fatal("traceId not found in context")
		}
		if traceId == "" {
			t.Fatal("traceId is empty")
		}
		// Verify it looks like a UUID.
		if _, err := uuid.Parse(traceId); err != nil {
			t.Fatalf("Generated traceId is not a valid UUID: %v", err)
		}
		return c.JSON(fiber.Map{"traceId": traceId})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)

	body, _ := io.ReadAll(resp.Body)
	var result fiber.Map
	json.Unmarshal(body, &result)

	if resp.Header.Get("X-Trace-Id") == "" {
		t.Fatal("X-Trace-Id header not set in response")
	}
	if result["traceId"] == "" {
		t.Fatal("traceId not returned in response body")
	}
}

func TestTraceIdMiddlewareUsesProvidedTraceId(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.TraceIdMiddleware)

	expectedTraceId := "550e8400-e29b-41d4-a716-446655440000" // Valid UUID

	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		traceId, ok := ctx.Value("traceId").(string)
		if !ok || traceId != expectedTraceId {
			t.Fatalf("Expected traceId %q, got %q", expectedTraceId, traceId)
		}
		return c.JSON(fiber.Map{"traceId": traceId})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Trace-Id", expectedTraceId)

	resp, _ := app.Test(req)

	if resp.Header.Get("X-Trace-Id") != expectedTraceId {
		t.Fatalf("Expected header %q, got %q", expectedTraceId, resp.Header.Get("X-Trace-Id"))
	}
}

func TestTraceIdMiddlewareFallsBackForInvalidUUID(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.TraceIdMiddleware)

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Middleware must NOT return 400 — it should silently replace invalid values
	// with a newly generated UUID and proceed with HTTP 200.
	invalidTraceIds := []string{
		"not-a-uuid",
		"12345",
		"550e8400-e29b-41d4-a716", // Incomplete UUID
		"not-uuid-at-all-bad",
	}

	for _, invalidId := range invalidTraceIds {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Trace-Id", invalidId)

		resp, _ := app.Test(req)

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected HTTP 200 for invalid traceId %q (fallback), got %d", invalidId, resp.StatusCode)
		}
		generated := resp.Header.Get("X-Trace-Id")
		if generated == "" {
			t.Errorf("expected X-Trace-Id header in response for input %q", invalidId)
			continue
		}
		if _, err := uuid.Parse(generated); err != nil {
			t.Errorf("expected UUID fallback for input %q, got %q: %v", invalidId, generated, err)
		}
		if generated == invalidId {
			t.Errorf("expected invalid value %q to be replaced, but it was echoed back", invalidId)
		}
	}
}

func TestErrorHandlerReturnsConsistentFormat(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.TraceIdMiddleware)

	// Trigger ErrorHandler (panic recovery) with a panicking route handler.
	app.Get("/test", func(c *fiber.Ctx) error {
		panic("something went wrong")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected HTTP 500, got %d", resp.StatusCode)
	}

	var errorResp fiber.Map
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &errorResp)

	if errorResp["error"] == nil {
		t.Error("expected 'error' field in error response")
	}
	if errorResp["message"] == nil {
		t.Error("expected 'message' field in error response")
	}
	if errorResp["status"] == nil {
		t.Error("expected 'status' field in error response")
	}
}

func TestStructuredLoggerWithTraceIdContext(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := logging.New(baseLogger, logging.LogFormatJSON)

	traceId := "test-trace-456"
	ctx := logging.WithTraceId(context.Background(), traceId)

	fields := map[string]interface{}{
		"userId": 12345,
		"action": "customer.created",
	}

	logger.Info(ctx, "Customer created successfully", fields)

	if buf.Len() == 0 {
		t.Fatal("No log output captured")
	}
}

func TestStructuredLoggerWithStructObject(t *testing.T) {
	type CustomerResponse struct {
		ID   uint   `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}

	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := logging.New(baseLogger, logging.LogFormatJSON)

	traceId := "test-trace-struct"
	ctx := logging.WithTraceId(context.Background(), traceId)

	response := CustomerResponse{
		ID:   12345,
		Code: "CODE123",
		Name: "John Doe",
	}

	logger.Info(ctx, "Customer created", response)

	if buf.Len() == 0 {
		t.Fatal("No log output captured for struct object")
	}
}

func TestStructuredLoggerTextFormatDevelopment(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	baseLogger := slog.New(handler)

	logger := logging.New(baseLogger, logging.LogFormatText)

	traceId := "dev-trace-789"
	ctx := logging.WithTraceId(context.Background(), traceId)

	logger.Info(ctx, "Development log message", map[string]interface{}{
		"userId": 987,
	})

	output := buf.String()
	if output == "" {
		t.Fatal("No log output for text format")
	}

	// The output should be human-readable (check for expected pattern).
	// Format: timestamp [traceIdShort] LEVEL message {context}
	if !bytes.Contains(buf.Bytes(), []byte("INFO")) {
		t.Error("Expected log level INFO in output")
	}
}

func TestMultipleRequestsHaveDifferentTraceIds(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.ErrorHandler)
	app.Use(middleware.TraceIdMiddleware)

	traceIds := make([]string, 3)
	idx := 0

	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		traceId := ctx.Value("traceId").(string)
		if idx < len(traceIds) {
			traceIds[idx] = traceId
			idx++
		}
		return c.JSON(fiber.Map{"traceId": traceId})
	})

	// Make 3 requests without providing traceId header.
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		app.Test(req)
	}

	// Verify each request got a different traceId.
	if traceIds[0] == traceIds[1] || traceIds[1] == traceIds[2] {
		t.Error("Expected different traceIds for each request")
	}

	// Verify all are valid UUIDs.
	for _, id := range traceIds {
		if _, err := uuid.Parse(id); err != nil {
			t.Errorf("Invalid UUID traceId: %v", err)
		}
	}
}

func TestStructuredLoggerAllLevels(t *testing.T) {
	levels := []struct {
		fn    func(context.Context, string, map[string]interface{})
		level string
	}{
		{
			fn:    func(ctx context.Context, msg string, fields map[string]interface{}) {
				buf := &bytes.Buffer{}
				h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
				logger := logging.New(slog.New(h), logging.LogFormatJSON)
				logger.Debug(ctx, msg, fields)
			},
			level: "DEBUG",
		},
		{
			fn:    func(ctx context.Context, msg string, fields map[string]interface{}) {
				buf := &bytes.Buffer{}
				h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
				logger := logging.New(slog.New(h), logging.LogFormatJSON)
				logger.Info(ctx, msg, fields)
			},
			level: "INFO",
		},
		{
			fn:    func(ctx context.Context, msg string, fields map[string]interface{}) {
				buf := &bytes.Buffer{}
				h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
				logger := logging.New(slog.New(h), logging.LogFormatJSON)
				logger.Warn(ctx, msg, fields)
			},
			level: "WARN",
		},
		{
			fn:    func(ctx context.Context, msg string, fields map[string]interface{}) {
				buf := &bytes.Buffer{}
				h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
				logger := logging.New(slog.New(h), logging.LogFormatJSON)
				logger.Error(ctx, msg, fields)
			},
			level: "ERROR",
		},
	}

	for _, tt := range levels {
		t.Run(tt.level, func(t *testing.T) {
			ctx := logging.WithTraceId(context.Background(), "test-trace")
			// Call the log function (it shouldn't panic).
			tt.fn(ctx, "Test message", map[string]interface{}{"key": "value"})
		})
	}
}
