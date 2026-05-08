package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cometagaming/casino-proxy-ai/internal/domain"
	"github.com/cometagaming/casino-proxy-ai/internal/infrastructure/idempotency"
	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
)

// mockCustomerRepository satisfies usecase.CustomerRepository for unit tests.
type mockCustomerRepository struct {
	customer     *domain.Customer
	getByCodeErr error
	saveErr      error
}

func (m *mockCustomerRepository) GetByCode(_ context.Context, _ string) (*domain.Customer, error) {
	return m.customer, m.getByCodeErr
}

func (m *mockCustomerRepository) Save(_ context.Context, _ *domain.Customer) error {
	return m.saveErr
}

// errIdempotencyStore always returns an error from AcquireLock (simulates infra failure).
type errIdempotencyStore struct{ err error }

func (e *errIdempotencyStore) AcquireLock(_ context.Context, _ string, _ time.Duration) (bool, error) {
	return false, e.err
}
func (e *errIdempotencyStore) SetStatus(_ context.Context, _ string, _ string) error { return nil }
func (e *errIdempotencyStore) DeleteKey(_ context.Context, _ string) error           { return nil }

func TestExecute(t *testing.T) {
	ctx := context.Background()
	const key = "TX-001"
	const ttl = 30 * time.Second
	lockErr := errors.New("redis down")

	t.Run("happy path", func(t *testing.T) {
		store := idempotency.NewMockIdempotencyStore()
		repo := &mockCustomerRepository{customer: &domain.Customer{ID: 1, Code: key, Name: "Old"}}
		uc := usecase.NewUpdateClientNameUseCase(repo, store, ttl)

		if err := uc.Execute(ctx, key, "New"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// key still held (COMPLETED) — re-acquire must fail
		ok, _ := store.AcquireLock(ctx, key, ttl)
		if ok {
			t.Error("expected lock to remain held after happy path")
		}
	})

	t.Run("duplicate message — AcquireLock returns false", func(t *testing.T) {
		store := idempotency.NewMockIdempotencyStore()
		_, _ = store.AcquireLock(ctx, key, ttl) // pre-acquire → second call returns false
		repo := &mockCustomerRepository{}       // GetByCode must never be called
		uc := usecase.NewUpdateClientNameUseCase(repo, store, ttl)

		if err := uc.Execute(ctx, key, "New"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("AcquireLock error — GetByCode never called", func(t *testing.T) {
		store := &errIdempotencyStore{err: lockErr}
		repo := &mockCustomerRepository{} // must not be called
		uc := usecase.NewUpdateClientNameUseCase(repo, store, ttl)

		err := uc.Execute(ctx, key, "New")
		if err == nil || err.Error() != lockErr.Error() {
			t.Fatalf("expected %v, got %v", lockErr, err)
		}
	})

	t.Run("GetByCode error — DeleteKey called", func(t *testing.T) {
		store := idempotency.NewMockIdempotencyStore()
		repo := &mockCustomerRepository{getByCodeErr: domain.ErrCustomerNotFound}
		uc := usecase.NewUpdateClientNameUseCase(repo, store, ttl)

		err := uc.Execute(ctx, key, "New")
		if !errors.Is(err, domain.ErrCustomerNotFound) {
			t.Fatalf("expected ErrCustomerNotFound, got %v", err)
		}
		// DeleteKey ran — key must be re-acquirable
		ok, _ := store.AcquireLock(ctx, key, ttl)
		if !ok {
			t.Error("expected lock re-acquirable after GetByCode error (DeleteKey must have run)")
		}
	})

	t.Run("UpdateName error — DeleteKey called", func(t *testing.T) {
		store := idempotency.NewMockIdempotencyStore()
		repo := &mockCustomerRepository{customer: &domain.Customer{ID: 1, Code: key, Name: "Old"}}
		uc := usecase.NewUpdateClientNameUseCase(repo, store, ttl)

		err := uc.Execute(ctx, key, "") // empty name → ErrInvalidName
		if !errors.Is(err, domain.ErrInvalidName) {
			t.Fatalf("expected ErrInvalidName, got %v", err)
		}
		ok, _ := store.AcquireLock(ctx, key, ttl)
		if !ok {
			t.Error("expected lock re-acquirable after UpdateName error (DeleteKey must have run)")
		}
	})

	t.Run("Save error — DeleteKey called", func(t *testing.T) {
		store := idempotency.NewMockIdempotencyStore()
		dbErr := errors.New("db timeout")
		repo := &mockCustomerRepository{
			customer: &domain.Customer{ID: 1, Code: key, Name: "Old"},
			saveErr:  dbErr,
		}
		uc := usecase.NewUpdateClientNameUseCase(repo, store, ttl)

		err := uc.Execute(ctx, key, "New")
		if err == nil || err.Error() != dbErr.Error() {
			t.Fatalf("expected %v, got %v", dbErr, err)
		}
		ok, _ := store.AcquireLock(ctx, key, ttl)
		if !ok {
			t.Error("expected lock re-acquirable after Save error (DeleteKey must have run)")
		}
	})
}
