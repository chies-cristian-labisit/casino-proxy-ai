package idempotency

import (
	"context"
	"sync"
	"time"

	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
)

var _ usecase.IdempotencyStore = (*MockIdempotencyStore)(nil)

type MockIdempotencyStore struct {
	mu    sync.Mutex
	store map[string]string
}

func NewMockIdempotencyStore() *MockIdempotencyStore {
	return &MockIdempotencyStore{store: make(map[string]string)}
}

func (m *MockIdempotencyStore) AcquireLock(_ context.Context, key string, _ time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.store[key]; exists {
		return false, nil
	}
	m.store[key] = "PENDING"
	return true, nil
}

func (m *MockIdempotencyStore) SetStatus(_ context.Context, key string, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = status
	return nil
}

func (m *MockIdempotencyStore) DeleteKey(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
	return nil
}
