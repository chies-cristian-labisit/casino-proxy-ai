.PHONY: up down clean run kill test integration acceptance test-all smoke help

PORT ?= 8081

help:
	@echo "Available targets:"
	@echo "  make up            - Start Docker Compose services (PostgreSQL + Kafka)"
	@echo "  make down          - Stop Docker Compose services"
	@echo "  make clean         - Stop and remove services + volumes"
	@echo "  make run           - Load .env.local and start the application"
	@echo "  make kill          - Kill any process on port $(PORT)"
	@echo "  make test          - Run unit tests (./internal/...)"
	@echo "  make integration   - Run integration tests (./tests/integration/...)"
	@echo "  make acceptance    - Run acceptance tests (./tests/acceptance/...)"
	@echo "  make test-all      - Run all tests: unit → integration → acceptance"
	@echo "  make smoke         - Quick health check: liveness + readiness"
	@echo ""

up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d
	@echo "✓ Services started. Run 'make smoke' to verify health."

down:
	@echo "Stopping Docker Compose services..."
	docker-compose down
	@echo "✓ Services stopped."

clean:
	@echo "Stopping and removing all services + volumes..."
	docker-compose down -v
	@echo "✓ Cleanup complete."

kill:
	@echo "Killing any process on port $(PORT)..."
	@lsof -i :$(PORT) -t | xargs -r kill -9 2>/dev/null || echo "No process found on port $(PORT)"
	@echo "✓ Port $(PORT) is free."

run: kill
	@echo "Loading environment and starting application..."
	@set -a; \
	source .env.local; \
	set +a; \
	go run ./cmd/api

test:
	@echo "Running unit tests..."
	go test ./internal/... -v -timeout 60s

integration:
	@echo "Running integration tests..."
	go test ./tests/integration/... -v -timeout 120s

acceptance:
	@echo "Running acceptance tests..."
	go test ./tests/acceptance/... -v -timeout 300s

test-all: test integration acceptance
	@echo "✓ All tests passed!"

smoke:
	@echo "Testing health endpoints..."
	@set -a; \
	source .env.local; \
	set +a; \
	echo "" && \
	echo "1. Liveness probe:" && \
	curl -s http://localhost:$$PORT/liveness | jq . 2>/dev/null || curl -s http://localhost:$$PORT/liveness && \
	echo "" && \
	echo "2. Readiness probe:" && \
	curl -s http://localhost:$$PORT/readiness | jq . 2>/dev/null || curl -s http://localhost:$$PORT/readiness && \
	echo "" && \
	echo "✓ Health checks complete"
