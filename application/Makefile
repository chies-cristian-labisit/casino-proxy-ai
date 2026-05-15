.PHONY: up down clean run kill test integration acceptance test-all smoke \
        up-full rebuild logs help

PORT ?= 8081

help:
	@echo "Local development targets:"
	@echo ""
	@echo "  Infrastructure (Docker):"
	@echo "  make up            - Start PostgreSQL + Kafka only"
	@echo "  make up-full       - Start full stack (PostgreSQL + Kafka + app container)"
	@echo "  make rebuild       - Rebuild and restart app container"
	@echo "  make logs          - Tail app container logs"
	@echo "  make down          - Stop all containers"
	@echo "  make clean         - Stop containers and remove volumes"
	@echo ""
	@echo "  Application (local Go):"
	@echo "  make run           - Load .env.local and run app with 'go run'"
	@echo "  make kill          - Kill any process on port $(PORT)"
	@echo ""
	@echo "  Testing:"
	@echo "  make test          - Unit tests (./internal/...)"
	@echo "  make integration   - Integration tests (./tests/integration/...)"
	@echo "  make acceptance    - Acceptance tests (./tests/acceptance/...)"
	@echo "  make test-all      - All tests in order: unit → integration → acceptance"
	@echo "  make smoke         - Quick health check: liveness + readiness"
	@echo ""

# ── Infrastructure ─────────────────────────────────────────────────────────────

up:
	@echo "Starting infrastructure (PostgreSQL + Kafka)..."
	docker-compose up -d postgres kafka
	@echo "✓ Infrastructure ready. Run 'make run' (local) or 'make up-full' (Docker)."

up-full:
	@echo "Starting full stack (PostgreSQL + Kafka + app)..."
	docker-compose --profile full up -d
	@echo "✓ Full stack running. Run 'make smoke' to verify health."

rebuild:
	@echo "Rebuilding and restarting app container..."
	docker-compose --profile full up -d --build app
	@echo "✓ App container rebuilt. Run 'make smoke' to verify."

logs:
	docker-compose logs -f app

down:
	@echo "Stopping all containers..."
	docker-compose --profile full down
	@echo "✓ Containers stopped."

clean:
	@echo "Stopping containers and removing volumes..."
	docker-compose --profile full down -v
	@echo "✓ Cleanup complete."

# ── Local application ───────────────────────────────────────────────────────────

kill:
	@echo "Killing any process on port $(PORT)..."
	@lsof -i :$(PORT) -t | xargs -r kill -9 2>/dev/null || echo "No process found on port $(PORT)"
	@echo "✓ Port $(PORT) is free."

run: kill
	@echo "Loading .env.local and starting application..."
	@set -a; \
	source ./.env.local; \
	set +a; \
	go run ./cmd/api

# ── Tests ───────────────────────────────────────────────────────────────────────

test:
	@echo "Running unit tests..."
	CGO_ENABLED=0 go test ./internal/... -v -timeout 60s

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
	source ./.env.local; \
	set +a; \
	echo "" && \
	echo "1. Liveness probe:" && \
	curl -s http://localhost:$$PORT/liveness | jq . 2>/dev/null || curl -s http://localhost:$$PORT/liveness && \
	echo "" && \
	echo "2. Readiness probe:" && \
	curl -s http://localhost:$$PORT/readiness | jq . 2>/dev/null || curl -s http://localhost:$$PORT/readiness && \
	echo "" && \
	echo "✓ Health checks complete"
