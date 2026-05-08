package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
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
}

func NewListener(reader kafkaReader, workers int, uc *usecase.UpdateClientNameUseCase) *Listener {
	return &Listener{reader: reader, workers: workers, uc: uc}
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
		log.Printf("kafka: reader close: %v", err)
	}
}

func (l *Listener) poll(ctx context.Context, msgCh chan<- kafkago.Message) {
	for {
		msg, err := l.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("kafka: fetch error: %v", err)
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
	var payload messagePayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		log.Printf("kafka: unmarshal error: %v", err)
		// Commit to avoid blocking the partition on a permanently malformed message.
		if err := l.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka: commit error (bad message): %v", err)
		}
		return
	}
	if err := l.uc.Execute(ctx, payload.CustomerCode, payload.CustomerName); err != nil {
		log.Printf("kafka: execute error code=%s: %v", payload.CustomerCode, err)
		// Do not commit — let the message redeliver on consumer restart/rebalance.
		return
	}
	if err := l.reader.CommitMessages(ctx, msg); err != nil {
		log.Printf("kafka: commit error: %v", err)
	}
}
