package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/handler"
	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/router"
	"github.com/cometagaming/casino-proxy-ai/internal/domain"
)

type stubCustomerRepo struct{}

func (s *stubCustomerRepo) GetByCode(_ context.Context, _ string) (*domain.Customer, error) {
	return &domain.Customer{Code: "TX-001", Name: "Test"}, nil
}

func (s *stubCustomerRepo) Save(_ context.Context, _ *domain.Customer) error {
	return nil
}

func newTestApp() *fiber.App {
	app := fiber.New()
	healthH := handler.NewHealthHandler(func(_ context.Context) error { return nil })
	customerH := handler.NewCustomerHandler(&stubCustomerRepo{})
	router.Setup(app, healthH, customerH)
	return app
}

func TestSetup_RegistersLivenessRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	resp, err := newTestApp().Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		t.Error("/liveness route not registered")
	}
}

func TestSetup_RegistersReadinessRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	resp, err := newTestApp().Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		t.Error("/readiness route not registered")
	}
}

func TestSetup_RegistersGetCustomerRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/TX-001", nil)
	resp, err := newTestApp().Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		t.Error("/api/v1/customers/:idTx route not registered")
	}
}

func TestSetup_Returns404ForUnknownRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	resp, err := newTestApp().Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for unknown route, got %d", resp.StatusCode)
	}
}
