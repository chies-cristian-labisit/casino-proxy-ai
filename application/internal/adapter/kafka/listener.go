package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"

	"github.com/cometagaming/ms-casino-go-v2/internal/logging"
	"github.com/cometagaming/ms-casino-go-v2/internal/usecase"
)

type kafkaReader interface {
	FetchMessage(ctx context.Context) (kafkago.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafkago.Message) error
	Close() error
}

type messagePayload struct {
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
}

type Listener struct {
	reader  kafkaReader
	workers int
	uc      *usecase.UpdateClientNameUseCase
	logger  *logging.StructuredLogger
}

func NewListener(reader kafkaReader, workers int, uc *usecase.UpdateClientNameUseCase) *Listener {
	logger := logging.New(slog.Default(), logging.LogFormatJSON)
	return &Listener{reader: reader, workers: workers, uc: uc, logger: logger}
}

func (l *Listener) Run(ctx context.Context) {
	msgCh := make(chan kafkago.Message, l.workers)
	var wg sync.WaitGroup
	for i := 0; i < l.workers; i++ {
		wg.Add(1)
		go l.worker(ctx, msgCh, &wg)
	}
	l.poll(ctx, msgCh)
	close(msgCh)
	wg.Wait()
	if err := l.reader.Close(); err != nil {
		l.logger.Error(ctx, "failed to close kafka reader", map[string]interface{}{"error": err.Error()})
	}
}

func (l *Listener) poll(ctx context.Context, msgCh chan<- kafkago.Message) {
	for {
		msg, err := l.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			l.logger.Warn(ctx, "failed to fetch kafka message", map[string]interface{}{"error": err.Error()})
			continue
		}
		select {
		case <-ctx.Done():
			return
		case msgCh <- msg:
		}
	}
}

func (l *Listener) worker(ctx context.Context, msgCh <-chan kafkago.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range msgCh {
		l.handle(ctx, msg)
	}
}

func (l *Listener) handle(ctx context.Context, msg kafkago.Message) {
	// Start a Datadog span for this message. When DD is disabled the tracer
	// returns a no-op span whose TraceID() is all-zeros (not ""), so we check
	// for both the empty string and the zero sentinel to trigger the UUID fallback.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.message.process")
	defer span.Finish()

	const ddNoopTraceID = "00000000000000000000000000000000"
	traceId := span.Context().TraceID()
	if traceId == "" || traceId == ddNoopTraceID {
		traceId = uuid.New().String()
	}
	msgCtx := context.WithValue(ctx, "traceId", traceId)

	var payload messagePayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		l.logger.Error(msgCtx, "failed to unmarshal kafka message", map[string]interface{}{
			"error": err.Error(),
		})
		// Commit to avoid blocking the partition on a permanently malformed message.
		if err := l.reader.CommitMessages(msgCtx, msg); err != nil {
			l.logger.Error(msgCtx, "failed to commit bad kafka message", map[string]interface{}{"error": err.Error()})
		}
		return
	}

	l.logger.Info(msgCtx, "kafka message received", payload)

	if err := l.uc.Execute(msgCtx, payload.CustomerCode, payload.CustomerName); err != nil {
		l.logger.Error(msgCtx, "failed to execute usecase", map[string]interface{}{
			"customer_code": payload.CustomerCode,
			"error":         err.Error(),
		})
		// Do not commit — let the message redeliver on consumer restart/rebalance.
		return
	}

	l.logger.Info(msgCtx, "kafka message processed successfully", payload)

	if err := l.reader.CommitMessages(msgCtx, msg); err != nil {
		l.logger.Error(msgCtx, "failed to commit kafka message", map[string]interface{}{"error": err.Error()})
	}
}
