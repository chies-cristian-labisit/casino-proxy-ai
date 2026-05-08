package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/handler"
	"github.com/cometagaming/casino-proxy-ai/internal/domain"
	"github.com/gofiber/fiber/v2"
)

// mockCustomerRepo satisfies usecase.CustomerRepository for handler tests.
type mockCustomerRepo struct {
	customer *domain.Customer
	err      error
}

func (m *mockCustomerRepo) GetByCode(_ context.Context, _ string) (*domain.Customer, error) {
	return m.customer, m.err
}

func (m *mockCustomerRepo) Save(_ context.Context, _ *domain.Customer) error {
	return nil
}

func newCustomerTestApp(h *handler.CustomerHandler) *fiber.App {
	app := fiber.New()
	app.Get("/api/v1/customers/:idTx", h.GetByIdTx)
	return app
}

func TestGetByIdTx_ReturnsCustomerJSON(t *testing.T) {
	repo := &mockCustomerRepo{
		customer: &domain.Customer{ID: 1, Code: "TX-001", Name: "Alice"},
	}
	h := handler.NewCustomerHandler(repo)
	app := newCustomerTestApp(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/TX-001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetByIdTx_Returns404OnNotFound(t *testing.T) {
	repo := &mockCustomerRepo{err: domain.ErrCustomerNotFound}
	h := handler.NewCustomerHandler(repo)
	app := newCustomerTestApp(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/MISSING", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetByIdTx_Returns500OnUnexpectedError(t *testing.T) {
	repo := &mockCustomerRepo{err: errors.New("db timeout")}
	h := handler.NewCustomerHandler(repo)
	app := newCustomerTestApp(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/TX-001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
