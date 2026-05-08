package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/handler"
	"github.com/cometagaming/casino-proxy-ai/internal/adapter/http/router"
	kafkaadapter "github.com/cometagaming/casino-proxy-ai/internal/adapter/kafka"
	"github.com/cometagaming/casino-proxy-ai/internal/config"
	"github.com/cometagaming/casino-proxy-ai/internal/infrastructure/database"
	"github.com/cometagaming/casino-proxy-ai/internal/infrastructure/idempotency"
	"github.com/cometagaming/casino-proxy-ai/internal/usecase"
)

func main() {
	// 1. Load configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// 2. Connect to Aurora PostgreSQL via GORM (pgx driver).
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("database: connect: %v", err)
	}

	// 3. Run schema migration for customerRecord model.
	if err := database.Migrate(db); err != nil {
		log.Fatalf("database: migrate: %v", err)
	}

	// 4. Create in-memory idempotency store (replace with Redis in production).
	store := idempotency.NewMockIdempotencyStore()

	// 5. Create CustomerRepository.
	repo := database.NewCustomerRepository(db)

	// 6. Create UpdateClientNameUseCase; lockTTL from IDEMPOTENCY_LOCK_TTL.
	lockTTL := time.Duration(cfg.IdempotencyLockTTL) * time.Second
	uc := usecase.NewUpdateClientNameUseCase(repo, store, lockTTL)

	// 7. Create Kafka reader and start listener in a separate goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     strings.Split(cfg.KafkaBrokers, ","),
		GroupID:     cfg.KafkaGroupID,
		Topic:       cfg.KafkaTopic,
		StartOffset: kafkago.FirstOffset,
	})
	listener := kafkaadapter.NewListener(reader, cfg.KafkaWorkers, uc)
	go listener.Run(ctx)

	// 8. Build readiness checker: DB ping.
	readiness := func(rCtx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(rCtx)
	}

	// 9. Wire HTTP handlers and router.
	healthH := handler.NewHealthHandler(readiness)
	customerH := handler.NewCustomerHandler(repo)

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	router.Setup(app, healthH, customerH)

	// 10. Graceful shutdown: wait for SIGTERM / SIGINT.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		log.Println("shutting down…")
		cancel() // stop Kafka listener
		if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
			log.Printf("http shutdown: %v", err)
		}
	}()

	addr := ":" + cfg.Port
	log.Printf("listening on %s (env=%s)", addr, cfg.AppEnv)
	if err := app.Listen(addr); err != nil {
		log.Printf("http: %v", err)
	}
}
