package middleware

import (
	"context"
	"log/slog"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TraceIdMiddleware assigns a single correlation ID to every request.
//
// When Datadog APM is active (DD_ENABLED=true), the ID is taken from the active
// span's dd.trace_id — created upstream by fibertrace.Middleware() — so that
// logs, APM traces, and the X-Trace-Id response header all share one identifier.
//
// When Datadog is disabled (DD_ENABLED=false) there is no active span, so the
// middleware falls back to the X-Trace-Id request header (must be a valid UUID)
// or generates a fresh UUID.
func TraceIdMiddleware(c *fiber.Ctx) error {
	var traceId string

	// Prefer dd.trace_id from the span created by fibertrace.Middleware upstream.
	// The no-op span (DD disabled) returns all-zeros, which we treat as absent.
	const ddNoopTraceID = "00000000000000000000000000000000"
	if span, ok := tracer.SpanFromContext(c.UserContext()); ok {
		if id := span.Context().TraceID(); id != ddNoopTraceID {
			traceId = id
		}
	}

	// No active DD span — fall back to X-Trace-Id header or generated UUID.
	if traceId == "" {
		traceId = c.Get("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
		} else if _, err := uuid.Parse(traceId); err != nil {
			slog.Warn("invalid X-Trace-Id header, generating new trace ID",
				"provided_value", traceId,
				"error", err.Error(),
			)
			traceId = uuid.New().String()
		}
	}

	newCtx := context.WithValue(c.UserContext(), "traceId", traceId)
	c.SetUserContext(newCtx)
	c.Set("X-Trace-Id", traceId)

	return c.Next()
}
