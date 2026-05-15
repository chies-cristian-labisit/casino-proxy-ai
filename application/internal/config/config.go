package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Port               string
	LogLevel           slog.Level
	DatabaseURL        string
	KafkaBrokers       string
	KafkaTopic         string
	KafkaGroupID       string
	KafkaWorkers       int
	IdempotencyLockTTL int
	AppEnv             string
	Datadog            *DatadogConfig
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	var missing []string
	if dbURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if kafkaBrokers == "" {
		missing = append(missing, "KAFKA_BROKERS")
	}
	if kafkaTopic == "" {
		missing = append(missing, "KAFKA_TOPIC")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	logLevel, err := ParseLogLevel()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:               getEnvOrDefault("PORT", "8081"),
		LogLevel:           logLevel,
		DatabaseURL:        dbURL,
		KafkaBrokers:       kafkaBrokers,
		KafkaTopic:         kafkaTopic,
		KafkaGroupID:       getEnvOrDefault("KAFKA_GROUP_ID", "ms-casino-go-v2"),
		KafkaWorkers:       mustAtoi(getEnvOrDefault("KAFKA_WORKERS", "5")),
		IdempotencyLockTTL: mustAtoi(getEnvOrDefault("IDEMPOTENCY_LOCK_TTL", "30")),
		AppEnv:             getEnvOrDefault("APP_ENV", "development"),
		Datadog:            loadDatadogConfig(),
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
