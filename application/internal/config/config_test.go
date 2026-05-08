package config

import (
	"strings"
	"testing"
)

func TestLoad_MissingAllRequired(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("KAFKA_BROKERS", "")
	t.Setenv("KAFKA_TOPIC", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when all required vars are missing, got nil")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("error should mention DATABASE_URL, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "KAFKA_BROKERS") {
		t.Errorf("error should mention KAFKA_BROKERS, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "KAFKA_TOPIC") {
		t.Errorf("error should mention KAFKA_TOPIC, got: %s", err.Error())
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("KAFKA_BROKERS", "localhost:9092")
	t.Setenv("KAFKA_TOPIC", "customer-updates")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing, got nil")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("error should mention DATABASE_URL, got: %s", err.Error())
	}
}

func TestLoad_MissingKafkaBrokers(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("KAFKA_BROKERS", "")
	t.Setenv("KAFKA_TOPIC", "customer-updates")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when KAFKA_BROKERS is missing, got nil")
	}
	if !strings.Contains(err.Error(), "KAFKA_BROKERS") {
		t.Errorf("error should mention KAFKA_BROKERS, got: %s", err.Error())
	}
}

func TestLoad_MissingKafkaTopic(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("KAFKA_BROKERS", "localhost:9092")
	t.Setenv("KAFKA_TOPIC", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when KAFKA_TOPIC is missing, got nil")
	}
	if !strings.Contains(err.Error(), "KAFKA_TOPIC") {
		t.Errorf("error should mention KAFKA_TOPIC, got: %s", err.Error())
	}
}

func TestLoad_AllRequiredVarsSet_Defaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("KAFKA_BROKERS", "localhost:9092")
	t.Setenv("KAFKA_TOPIC", "customer-updates")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "8081" {
		t.Errorf("expected default Port=8081, got %s", cfg.Port)
	}
	if cfg.KafkaGroupID != "casino-proxy-ai" {
		t.Errorf("expected default KafkaGroupID=casino-proxy-ai, got %s", cfg.KafkaGroupID)
	}
	if cfg.KafkaWorkers != 5 {
		t.Errorf("expected default KafkaWorkers=5, got %d", cfg.KafkaWorkers)
	}
	if cfg.IdempotencyLockTTL != 30 {
		t.Errorf("expected default IdempotencyLockTTL=30, got %d", cfg.IdempotencyLockTTL)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("expected default AppEnv=development, got %s", cfg.AppEnv)
	}
}

func TestLoad_AllVarsExplicitlySet(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://prod-host/mydb")
	t.Setenv("KAFKA_BROKERS", "broker1:9092,broker2:9092")
	t.Setenv("KAFKA_TOPIC", "prod-customer-updates")
	t.Setenv("PORT", "9090")
	t.Setenv("KAFKA_GROUP_ID", "my-group")
	t.Setenv("KAFKA_WORKERS", "10")
	t.Setenv("IDEMPOTENCY_LOCK_TTL", "60")
	t.Setenv("APP_ENV", "production")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("expected Port=9090, got %s", cfg.Port)
	}
	if cfg.KafkaBrokers != "broker1:9092,broker2:9092" {
		t.Errorf("expected KafkaBrokers=broker1:9092,broker2:9092, got %s", cfg.KafkaBrokers)
	}
	if cfg.KafkaTopic != "prod-customer-updates" {
		t.Errorf("expected KafkaTopic=prod-customer-updates, got %s", cfg.KafkaTopic)
	}
	if cfg.KafkaGroupID != "my-group" {
		t.Errorf("expected KafkaGroupID=my-group, got %s", cfg.KafkaGroupID)
	}
	if cfg.KafkaWorkers != 10 {
		t.Errorf("expected KafkaWorkers=10, got %d", cfg.KafkaWorkers)
	}
	if cfg.IdempotencyLockTTL != 60 {
		t.Errorf("expected IdempotencyLockTTL=60, got %d", cfg.IdempotencyLockTTL)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("expected AppEnv=production, got %s", cfg.AppEnv)
	}
}
