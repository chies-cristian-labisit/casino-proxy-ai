package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	fibertrace "github.com/DataDog/dd-trace-go/contrib/gofiber/fiber.v2/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// newTestApp builds a Fiber app with TraceIdMiddleware.
// withDD=true registers fibertrace.Middleware() first, simulating production order.
func newTestApp(withDD bool) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	if withDD {
		app.Use(fibertrace.Middleware())
	}
	app.Use(TraceIdMiddleware)
	return app
}

func doGet(t *testing.T, app *fiber.App, headers map[string]string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp
}

// TestTraceIdMiddleware_UsesDDTraceID_WhenSpanActive verifies that when Datadog APM
// is active, the middleware uses dd.trace_id as the single correlation ID — stored
// in context AND echoed in X-Trace-Id response header, both matching the span.
func TestTraceIdMiddleware_UsesDDTraceID_WhenSpanActive(t *testing.T) {
	tracer.Start(tracer.WithService("test"))
	defer tracer.Stop()

	app := newTestApp(true)
	var capturedCtx context.Context
	app.Get("/test", func(c *fiber.Ctx) error {
		capturedCtx = c.UserContext()
		return c.SendString("ok")
	})

	resp := doGet(t, app, nil)

	headerID := resp.Header.Get("X-Trace-Id")
	if headerID == "" {
		t.Fatal("expected X-Trace-Id response header when DD span is active")
	}

	// Context and response header must carry the same ID.
	ctxID, _ := capturedCtx.Value("traceId").(string)
	if ctxID != headerID {
		t.Errorf("context traceId %q != X-Trace-Id header %q", ctxID, headerID)
	}

	// The ID must equal the actual dd.trace_id of the span that was in context.
	span, ok := tracer.SpanFromContext(capturedCtx)
	if !ok {
		t.Fatal("expected a DD span in the request context")
	}
	if want := span.Context().TraceID(); ctxID != want {
		t.Errorf("traceId %q does not match dd.trace_id %q from span", ctxID, want)
	}
}

// TestTraceIdMiddleware_GeneratesUUID_WhenNoSpan verifies that when DD is disabled
// (no active span), a fresh UUID is generated and set as the correlation ID.
func TestTraceIdMiddleware_GeneratesUUID_WhenNoSpan(t *testing.T) {
	app := newTestApp(false)
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("ok") })

	resp := doGet(t, app, nil)

	traceId := resp.Header.Get("X-Trace-Id")
	if traceId == "" {
		t.Fatal("expected X-Trace-Id to be generated when no DD span")
	}
	if _, err := uuid.Parse(traceId); err != nil {
		t.Errorf("expected UUID fallback when no span, got %q: %v", traceId, err)
	}
}

// TestTraceIdMiddleware_EchoesUUIDHeader_WhenNoSpan verifies that a valid UUID
// provided in the X-Trace-Id request header is accepted and echoed back.
func TestTraceIdMiddleware_EchoesUUIDHeader_WhenNoSpan(t *testing.T) {
	app := newTestApp(false)
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("ok") })

	clientID := uuid.New().String()
	resp := doGet(t, app, map[string]string{"X-Trace-Id": clientID})

	if got := resp.Header.Get("X-Trace-Id"); got != clientID {
		t.Errorf("expected X-Trace-Id %q echoed back, got %q", clientID, got)
	}
}

// TestTraceIdMiddleware_FallsBackToUUID_WhenInvalidHeader verifies that a
// non-UUID value in X-Trace-Id is silently replaced by a generated UUID (no
// HTTP 400) when DD is disabled. The warning is logged to slog but the request
// proceeds normally with HTTP 200.
func TestTraceIdMiddleware_FallsBackToUUID_WhenInvalidHeader(t *testing.T) {
	app := newTestApp(false)
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("ok") })

	resp := doGet(t, app, map[string]string{"X-Trace-Id": "not-a-uuid"})

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected HTTP 200 for invalid X-Trace-Id (fallback), got %d", resp.StatusCode)
	}
	traceId := resp.Header.Get("X-Trace-Id")
	if traceId == "" {
		t.Fatal("expected X-Trace-Id header in response")
	}
	if _, err := uuid.Parse(traceId); err != nil {
		t.Errorf("expected a valid UUID fallback, got %q: %v", traceId, err)
	}
	if traceId == "not-a-uuid" {
		t.Error("expected invalid header to be replaced, but original value was echoed")
	}
}

// TestTraceIdMiddleware_DifferentRequestsDifferentIDs verifies that consecutive
// requests without a header each receive a unique generated ID.
func TestTraceIdMiddleware_DifferentRequestsDifferentIDs(t *testing.T) {
	app := newTestApp(false)
	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("ok") })

	ids := make([]string, 3)
	for i := range ids {
		ids[i] = doGet(t, app, nil).Header.Get("X-Trace-Id")
	}

	if ids[0] == ids[1] || ids[1] == ids[2] {
		t.Errorf("expected unique trace IDs per request, got %v", ids)
	}
}
