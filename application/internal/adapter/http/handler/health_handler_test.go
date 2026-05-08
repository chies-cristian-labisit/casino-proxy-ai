package handler_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/handler"
	"github.com/gofiber/fiber/v2"
)

func newTestApp(h *handler.HealthHandler) *fiber.App {
	app := fiber.New()
	app.Get("/liveness", h.Liveness)
	app.Get("/readiness", h.Readiness)
	return app
}

func TestLiveness_AlwaysOK(t *testing.T) {
	h := handler.NewHealthHandler(func(_ context.Context) error {
		return errors.New("checker should not be called for liveness")
	})
	app := newTestApp(h)

	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestReadiness_ReturnsOKWhenCheckerPasses(t *testing.T) {
	h := handler.NewHealthHandler(func(_ context.Context) error { return nil })
	app := newTestApp(h)

	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestReadiness_Returns503WhenCheckerFails(t *testing.T) {
	h := handler.NewHealthHandler(func(_ context.Context) error {
		return errors.New("db unreachable")
	})
	app := newTestApp(h)

	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Error("expected non-empty error body")
	}
}
