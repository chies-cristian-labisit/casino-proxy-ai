package usecase

import (
	"context"
	"time"

	"github.com/cometagaming/ms-casino-go-v2/internal/domain"
)

type CustomerRepository interface {
	GetByCode(ctx context.Context, code string) (*domain.Customer, error)
	Save(ctx context.Context, customer *domain.Customer) error
}

type IdempotencyStore interface {
	// AcquireLock attempts to set a PENDING lock for the given key.
	// Returns true if the lock was acquired (first-time processing).
	// Returns false if the key already exists (duplicate message).
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// SetStatus updates the status of an existing lock key (e.g., COMPLETED).
	SetStatus(ctx context.Context, key string, status string) error

	// DeleteKey removes the lock key to allow retry on failure.
	DeleteKey(ctx context.Context, key string) error
}
