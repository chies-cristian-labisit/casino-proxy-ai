package idempotency

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAcquireLock(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*MockIdempotencyStore)
		key      string
		wantLock bool
	}{
		{
			name:     "first call on new key returns true",
			setup:    func(_ *MockIdempotencyStore) {},
			key:      "msg-001",
			wantLock: true,
		},
		{
			name: "second call on same key returns false",
			setup: func(s *MockIdempotencyStore) {
				_, _ = s.AcquireLock(context.Background(), "msg-002", time.Second)
			},
			key:      "msg-002",
			wantLock: false,
		},
		{
			name: "different keys are independent",
			setup: func(s *MockIdempotencyStore) {
				_, _ = s.AcquireLock(context.Background(), "msg-003", time.Second)
			},
			key:      "msg-004",
			wantLock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMockIdempotencyStore()
			tt.setup(s)
			got, err := s.AcquireLock(context.Background(), tt.key, time.Second)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantLock {
				t.Errorf("AcquireLock(%q) = %v, want %v", tt.key, got, tt.wantLock)
			}
		})
	}
}

func TestSetStatus(t *testing.T) {
	s := NewMockIdempotencyStore()
	ctx := context.Background()

	if err := s.SetStatus(ctx, "msg-010", "PENDING"); err != nil {
		t.Fatalf("SetStatus returned error: %v", err)
	}
	if s.store["msg-010"] != "PENDING" {
		t.Errorf("expected PENDING, got %q", s.store["msg-010"])
	}

	if err := s.SetStatus(ctx, "msg-010", "COMPLETED"); err != nil {
		t.Fatalf("SetStatus overwrite returned error: %v", err)
	}
	if s.store["msg-010"] != "COMPLETED" {
		t.Errorf("expected COMPLETED after overwrite, got %q", s.store["msg-010"])
	}
}

func TestDeleteKey(t *testing.T) {
	s := NewMockIdempotencyStore()
	ctx := context.Background()

	// acquire then delete — re-acquire must succeed
	_, _ = s.AcquireLock(ctx, "msg-020", time.Second)

	if err := s.DeleteKey(ctx, "msg-020"); err != nil {
		t.Fatalf("DeleteKey returned error: %v", err)
	}

	got, err := s.AcquireLock(ctx, "msg-020", time.Second)
	if err != nil {
		t.Fatalf("AcquireLock after delete returned error: %v", err)
	}
	if !got {
		t.Error("expected AcquireLock to return true after DeleteKey, got false")
	}

	// delete absent key — must be a no-op (no error)
	if err := s.DeleteKey(ctx, "never-existed"); err != nil {
		t.Errorf("DeleteKey on absent key returned error: %v", err)
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := NewMockIdempotencyStore()
	ctx := context.Background()

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(n int) {
			defer wg.Done()
			key := "concurrent-key"
			_, _ = s.AcquireLock(ctx, key, time.Second)
			_ = s.SetStatus(ctx, key, "PROCESSING")
			if n%2 == 0 {
				_ = s.DeleteKey(ctx, key)
			}
		}(i)
	}

	wg.Wait()
}
