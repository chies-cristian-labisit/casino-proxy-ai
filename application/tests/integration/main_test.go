package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cometagaming/ms-casino-go-v2/internal/infrastructure/database"
)

var (
	sharedDB   *gorm.DB
	sharedRepo *database.CustomerRepository
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
		fmt.Fprintf(os.Stderr, "skipping integration tests — could not start postgres container: %v\n", err)
		os.Exit(0)
	}

	pgDSN, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres connection string: %v\n", err)
		_ = pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "gorm open: %v\n", err)
		_ = pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	if err := database.Migrate(db); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		_ = pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	sharedDB = db
	sharedRepo = database.NewCustomerRepository(db)

	code := m.Run()
	_ = pgContainer.Terminate(ctx)
	os.Exit(code)
}
