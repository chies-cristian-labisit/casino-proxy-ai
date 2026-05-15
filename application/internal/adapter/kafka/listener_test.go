package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"

	"github.com/cometagaming/ms-casino-go-v2/internal/domain"
	"github.com/cometagaming/ms-casino-go-v2/internal/infrastructure/idempotency"
	"github.com/cometagaming/ms-casino-go-v2/internal/usecase"
)

// ── mocks ──────────────────────────────────────────────────────────────────────

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

// contextCapturingRepo wraps mockCustomerRepository and records the context
// received by GetByCode so tests can inspect trace IDs propagated into it.
type contextCapturingRepo struct {
	mockCustomerRepository
	capturedCtx context.Context
}

func (r *contextCapturingRepo) GetByCode(ctx context.Context, code string) (*domain.Customer, error) {
	r.capturedCtx = ctx
	return r.mockCustomerRepository.GetByCode(ctx, code)
}

// ── helpers ────────────────────────────────────────────────────────────────────

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

// ── behaviour tests ────────────────────────────────────────────────────────────

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

// ── trace ID tests ─────────────────────────────────────────────────────────────

// TestHandle_TraceId_IsDDTraceID_WhenTracerActive verifies that when the Datadog
// tracer is running, handle() stores dd.trace_id as the traceId in context —
// giving Kafka logs the same single ID as APM spans.
func TestHandle_TraceId_IsDDTraceID_WhenTracerActive(t *testing.T) {
	tracer.Start(tracer.WithService("test"))
	defer tracer.Stop()

	repo := &contextCapturingRepo{
		mockCustomerRepository: mockCustomerRepository{
			customer: &domain.Customer{ID: 1, Code: "TX-001", Name: "Old"},
		},
	}
	l := makeListener(&mockKafkaReader{}, repo)

	l.handle(context.Background(), msg("TX-001", "New"))

	traceId, ok := repo.capturedCtx.Value("traceId").(string)
	if !ok || traceId == "" {
		t.Fatal("expected non-empty traceId in context when DD tracer is active")
	}
	// The traceId must equal the actual dd.trace_id stored in the message context span.
	span, spanOK := tracer.SpanFromContext(repo.capturedCtx)
	if !spanOK {
		t.Fatal("expected a DD span in the message context when tracer is active")
	}
	if want := span.Context().TraceID(); traceId != want {
		t.Errorf("traceId %q does not match dd.trace_id %q from span", traceId, want)
	}
}

// TestHandle_TraceId_IsUUID_WhenTracerDisabled verifies that when Datadog is
// disabled (no active tracer), handle() falls back to a random UUID as the
// traceId — not the all-zeros sentinel that the DD no-op span returns.
func TestHandle_TraceId_IsUUID_WhenTracerDisabled(t *testing.T) {
	repo := &contextCapturingRepo{
		mockCustomerRepository: mockCustomerRepository{
			customer: &domain.Customer{ID: 1, Code: "TX-001", Name: "Old"},
		},
	}
	l := makeListener(&mockKafkaReader{}, repo)

	l.handle(context.Background(), msg("TX-001", "New"))

	traceId, ok := repo.capturedCtx.Value("traceId").(string)
	if !ok || traceId == "" {
		t.Fatal("expected non-empty traceId in context when DD tracer is disabled")
	}
	// Must be a valid UUID format.
	if _, err := uuid.Parse(traceId); err != nil {
		t.Errorf("expected UUID fallback when tracer is disabled, got %q: %v", traceId, err)
	}
	// Must NOT be the all-zeros no-op trace ID returned by the DD disabled span.
	const ddNoopTraceID = "00000000000000000000000000000000"
	if traceId == ddNoopTraceID {
		t.Errorf("expected a random UUID but got the DD no-op sentinel %q", traceId)
	}
}
