package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/handler"
	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/router"
	kafkaadapter "github.com/cometagaming/ms-casino-go-v2/internal/adapter/kafka"
	"github.com/cometagaming/ms-casino-go-v2/internal/config"
	"github.com/cometagaming/ms-casino-go-v2/internal/infrastructure/database"
	"github.com/cometagaming/ms-casino-go-v2/internal/infrastructure/idempotency"
	"github.com/cometagaming/ms-casino-go-v2/internal/observability"
	"github.com/cometagaming/ms-casino-go-v2/internal/usecase"
)

func main() {
	// 1. Load configuration from environment variables.
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize logger with configured log level (must be first)
	baseLogger := config.NewLogger(cfg.LogLevel)
	slog.SetDefault(baseLogger)

	// Initialize Datadog APM tracer (no-op when DD_ENABLED=false).
	stopTracer := observability.InitTracer(cfg.Datadog)
	defer stopTracer()

	// 2. Connect to Aurora PostgreSQL via GORM (pgx driver).
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		baseLogger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// 3. Run schema migration for customerRecord model.
	if err := database.Migrate(db); err != nil {
		baseLogger.Error("failed to run migration", "error", err)
		os.Exit(1)
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
	app.Use(observability.FiberMiddleware())
	router.Setup(app, healthH, customerH)

	// 10. Graceful shutdown: wait for SIGTERM / SIGINT.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		baseLogger.Info("shutting down")
		cancel() // stop Kafka listener
		if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
			baseLogger.Error("failed to shutdown http server", "error", err)
		}
	}()

	addr := ":" + cfg.Port
	baseLogger.Info("listening on server", "addr", addr, "env", cfg.AppEnv)
	if err := app.Listen(addr); err != nil {
		baseLogger.Error("http server error", "error", err)
	}
}
