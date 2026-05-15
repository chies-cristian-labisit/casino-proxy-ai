package acceptance_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	kafkatc "github.com/testcontainers/testcontainers-go/modules/kafka"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/handler"
	"github.com/cometagaming/ms-casino-go-v2/internal/adapter/http/router"
	kafkaadapter "github.com/cometagaming/ms-casino-go-v2/internal/adapter/kafka"
	"github.com/cometagaming/ms-casino-go-v2/internal/infrastructure/database"
	"github.com/cometagaming/ms-casino-go-v2/internal/infrastructure/idempotency"
	"github.com/cometagaming/ms-casino-go-v2/internal/usecase"
	"testing"
)

const testTopic = "events"

var (
	sharedDB          *gorm.DB
	sharedRepo        *database.CustomerRepository
	sharedApp         *fiber.App
	sharedKafkaWriter *kafkago.Writer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "skipping acceptance tests — could not start postgres container: %v\n", err)
		os.Exit(0)
	}
	defer pgContainer.Terminate(ctx) //nolint:errcheck

	pgDSN, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres connection string: %v\n", err)
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "gorm open: %v\n", err)
		os.Exit(1)
	}
	if err := database.Migrate(db); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}
	sharedDB = db
	sharedRepo = database.NewCustomerRepository(db)

	kafkaContainer, err := kafkatc.Run(ctx, "confluentinc/confluent-local:7.5.0", kafkatc.WithClusterID("test-cluster"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "skipping acceptance tests — could not start kafka container: %v\n", err)
		os.Exit(0)
	}
	defer kafkaContainer.Terminate(ctx) //nolint:errcheck

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kafka brokers: %v\n", err)
		os.Exit(1)
	}

	// Create the topic before the consumer starts.
	conn, err := kafkago.Dial("tcp", brokers[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "kafka dial: %v\n", err)
		os.Exit(1)
	}
	if err := conn.CreateTopics(kafkago.TopicConfig{
		Topic:             testTopic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "kafka create topic: %v\n", err)
		conn.Close()
		os.Exit(1)
	}
	conn.Close()

	sharedKafkaWriter = kafkago.NewWriter(kafkago.WriterConfig{
		Brokers: brokers,
		Topic:   testTopic,
	})

	store := idempotency.NewMockIdempotencyStore()
	lockTTL := 30 * time.Second
	uc := usecase.NewUpdateClientNameUseCase(sharedRepo, store, lockTTL)

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     "acceptance-test",
		Topic:       testTopic,
		StartOffset: kafkago.FirstOffset,
		MaxWait:     500 * time.Millisecond,
	})

	kafkaCtx, kafkaCancel := context.WithCancel(context.Background())
	listener := kafkaadapter.NewListener(reader, 1, uc)
	go listener.Run(kafkaCtx)

	readiness := func(rCtx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(rCtx)
	}
	healthH := handler.NewHealthHandler(readiness)
	customerH := handler.NewCustomerHandler(sharedRepo)
	app := fiber.New()
	router.Setup(app, healthH, customerH)
	sharedApp = app

	code := m.Run()
	_ = sharedKafkaWriter.Close()
	kafkaCancel()
	os.Exit(code)
}
