package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/handler"
	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/router"
	"github.com/cometagaming/ms-casino-go-v2/internal/domain"
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

func newCustomerTestApp(repo *mockCustomerRepo) *fiber.App {
	app := fiber.New()
	healthH := handler.NewHealthHandler(func(_ context.Context) error { return nil })
	customerH := handler.NewCustomerHandler(repo)
	router.Setup(app, healthH, customerH)
	return app
}

func TestGetByIdTx_ReturnsCustomerJSON(t *testing.T) {
	repo := &mockCustomerRepo{
		customer: &domain.Customer{ID: 1, Code: "TX-001", Name: "Alice"},
	}
	app := newCustomerTestApp(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/customers/TX-001", nil)
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
	app := newCustomerTestApp(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/customers/MISSING", nil)
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
	app := newCustomerTestApp(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/customers/TX-001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
