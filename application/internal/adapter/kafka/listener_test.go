package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/cometagaming/casino-proxy-ai/internal/domain"
	"github.com/cometagaming/casino-proxy-ai/internal/infrastructure/idempotency"
	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
)

type mockKafkaReader struct {
	commitCalled bool
	commitErr    error
}

func (m *mockKafkaReader) FetchMessage(_ context.Context) (kafkago.Message, error) {
	return kafkago.Message{}, nil
}

func (m *mockKafkaReader) CommitMessages(_ context.Context, _ ...kafkago.Message) error {
	m.commitCalled = true
	return m.commitErr
}

func (m *mockKafkaReader) Close() error { return nil }

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

func makeListener(reader kafkaReader, repo usecase.CustomerRepository) *Listener {
	store := idempotency.NewMockIdempotencyStore()
	uc := usecase.NewUpdateClientNameUseCase(repo, store, 30*time.Second)
	return NewListener(reader, 1, uc)
}

func msg(code, name string) kafkago.Message {
	return kafkago.Message{
		Value: []byte(`{"customer_code":"` + code + `","customer_name":"` + name + `"}`),
	}
}

func TestHandle_HappyPath(t *testing.T) {
	reader := &mockKafkaReader{}
	repo := &mockCustomerRepository{
		customer: &domain.Customer{ID: 1, Code: "TX-001", Name: "Old"},
	}
	l := makeListener(reader, repo)

	l.handle(context.Background(), msg("TX-001", "New"))

	if !reader.commitCalled {
		t.Error("expected CommitMessages to be called on success")
	}
}

func TestHandle_UnmarshalError(t *testing.T) {
	reader := &mockKafkaReader{}
	repo := &mockCustomerRepository{}
	l := makeListener(reader, repo)

	l.handle(context.Background(), kafkago.Message{Value: []byte("not-json{")})

	if !reader.commitCalled {
		t.Error("expected CommitMessages to be called on unmarshal error (to avoid partition block)")
	}
}

func TestHandle_ExecuteError(t *testing.T) {
	reader := &mockKafkaReader{}
	repo := &mockCustomerRepository{
		getByCodeErr: errors.New("db down"),
	}
	l := makeListener(reader, repo)

	l.handle(context.Background(), msg("TX-002", "New"))

	if reader.commitCalled {
		t.Error("expected CommitMessages NOT to be called on Execute error")
	}
}
