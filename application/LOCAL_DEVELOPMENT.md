# Local Development Guide

> **Self-contained setup** — PostgreSQL and Kafka run locally via Docker Compose. No cloud accounts, no external services, no shared environments required. A working installation of **Go 1.25.4+** and **Docker** is all you need.

Two modes are supported:

| Mode | Prerequisites | Command | When to use |
|------|---------------|---------|-------------|
| **A — Hybrid** | Go + Docker | `make up` → `make run` | Active development (attach debugger, fast restart) |
| **B — Full Docker** | Docker only | `make up-full` | Run without installing Go locally; CI-equivalent |

---

## Prerequisites

| Tool | Mode A (Hybrid) | Mode B (Full Docker) |
|------|----------------|----------------------|
| Go 1.25.4+ | Required | Not required |
| Docker & Docker Compose | Required | Required |
| `make` | Optional | Optional |

> **No Go installed?** Use Mode B — `make up-full` builds and runs the application inside a container.

---

## Step 0 — Initialize AIOX Framework

Run once after cloning the project:

```bash
npm run setup
```

This installs the AIOX agent framework, configures git hooks, and verifies the environment. It is idempotent — safe to re-run at any time.

---

## Mode A — Hybrid: Infra in Docker, app runs locally

Use this mode while developing. Infra (PostgreSQL + Kafka) runs in Docker; you run the Go binary directly so you can attach a debugger and see live reloads.

### Step 1 — Resolve dependencies

```bash
cd application
go mod tidy
```

### Step 2 — Start infrastructure

```bash
make up
```

Verify both services are healthy:

```bash
docker-compose ps
```

Expected output:
```
NAME                        STATUS
ms-casino-go-v2-postgres   Up (healthy)
ms-casino-go-v2-kafka      Up (healthy)
```

### Step 3 — Run the application

```bash
make run
```

This sources `.env.local` and executes `go run ./cmd/api`. The app auto-migrates the database schema on first start.

### Step 4 — Verify

```bash
make smoke
```

Or manually:

```bash
curl http://localhost:8081/liveness   # → 200 OK
curl http://localhost:8081/readiness  # → 200 OK (DB reachable)
```

---

## Mode B — Full Docker: everything in containers

Use this mode to verify the Docker image, test CI-equivalent behaviour, or run a staging-like stack locally.

### Step 1 — Start full stack

```bash
make up-full
```

This builds the image from `Dockerfile` and starts PostgreSQL, Kafka, and the app container. The app connects to Kafka via the internal `kafka:29092` listener and to PostgreSQL via the `postgres` hostname — both resolved on the Docker network.

### Step 2 — Verify

```bash
make smoke
```

### Rebuild after code changes

```bash
make rebuild
```

### Tail app logs

```bash
make logs
```

---

## Service URLs

| Service | Address (host) |
|---------|----------------|
| Application | `http://localhost:8081` |
| PostgreSQL | `localhost:5432` |
| Kafka (host access) | `localhost:9092` |
| Kafka (container-to-container) | `kafka:29092` |

---

## Testing

### Run all tests

```bash
make test-all       # unit → integration → acceptance
```

### Run individually

```bash
make test           # Unit tests (./internal/...) — no Docker required
make integration    # Integration tests — PostgreSQL container only (Testcontainers)
make acceptance     # Acceptance tests — full stack (PostgreSQL + Kafka + Fiber)
```

Acceptance tests spin up their own containers via Testcontainers — `make up` does **not** need to be running for tests.

### Test a specific endpoint

```bash
curl http://localhost:8081/v2/customers/BR123456789
```

---

## Stopping Services

```bash
make down           # Stop containers (keep volumes)
make clean          # Stop containers + remove volumes (deletes all data)
```

---

## Development Workflow

### Hybrid (recommended)

```
make up             # start infra
make run            # start app locally
# edit code → Ctrl+C → make run
make test-all       # run tests
make down           # stop infra
```

### Full Docker

```
make up-full        # build image + start everything
# edit code → make rebuild
make smoke          # verify health
make down           # stop all
```

---

## Logging Configuration

### Setting the Log Level

Control verbosity with the `LOG_LEVEL` environment variable:

```bash
# Local development (default)
make run
# Uses LOG_LEVEL=INFO from .env.local

# Debug everything (very verbose)
LOG_LEVEL=DEBUG make run

# Only warnings and errors (production-like)
LOG_LEVEL=WARN make run

# Only errors (minimal output)
LOG_LEVEL=ERROR make run
```

### Log Format

Log format is automatically selected by `APP_ENV`:

#### Development Mode (`APP_ENV=development`)

Human-readable text format with timestamps and trace IDs. Useful for active debugging.

```
2026-05-15 10:30:45.123 [a1b2c3d4] INFO  listening on server addr=:8081 env=development
2026-05-15 10:30:46.234 [b2c3d4e5] INFO  kafka message received {customerCode: BR123456789, customerName: João Silva}
2026-05-15 10:30:46.456 [b2c3d4e5] INFO  kafka message processed successfully {customerCode: BR123456789}
2026-05-15 10:30:47.567 [c3d4e5f6] WARN  failed to fetch kafka message error="context deadline exceeded"
```

#### Production Mode (`APP_ENV=production`)

Structured JSON format (suitable for log aggregators like Datadog, ELK, Splunk). Note: `traceId` and `dd.trace_id` are the same value when Datadog APM is enabled.

```json
{"timestamp":"2026-05-15T10:30:45.123Z","level":"INFO","message":"listening on server","traceId":"a1b2c3d4","context":{"addr":":8081","env":"production"},"dd.trace_id":"a1b2c3d4","dd.span_id":"987654321"}
{"timestamp":"2026-05-15T10:30:46.234Z","level":"INFO","message":"kafka message received","traceId":"b2c3d4e5","context":{"customerCode":"BR123456789","customerName":"João Silva"},"dd.trace_id":"b2c3d4e5","dd.span_id":"123456789"}
{"timestamp":"2026-05-15T10:30:46.456Z","level":"INFO","message":"kafka message processed successfully","traceId":"b2c3d4e5","context":{"customerCode":"BR123456789"},"dd.trace_id":"b2c3d4e5","dd.span_id":"123456790"}
{"timestamp":"2026-05-15T10:30:47.567Z","level":"WARN","message":"failed to fetch kafka message","traceId":"c3d4e5f6","context":{"error":"context deadline exceeded"},"dd.trace_id":"c3d4e5f6","dd.span_id":"234567891"}
```

### Log Level Examples

**DEBUG level** — Shows all messages (very verbose, development only):
```
2026-05-15 10:30:45.123 [a1b2c3d4] DEBUG starting http server...
2026-05-15 10:30:45.234 [a1b2c3d4] DEBUG loading database config...
2026-05-15 10:30:45.345 [a1b2c3d4] INFO  listening on server addr=:8081
```

**INFO level** — General info + warnings + errors (default, recommended):
```
2026-05-15 10:30:45.123 [a1b2c3d4] INFO  listening on server addr=:8081
2026-05-15 10:30:46.234 [b2c3d4e5] INFO  kafka message received {...}
2026-05-15 10:30:47.567 [c3d4e5f6] WARN  failed to fetch kafka message error="..."
```

**WARN level** — Warnings + errors only:
```
2026-05-15 10:30:47.567 [c3d4e5f6] WARN  failed to fetch kafka message error="context deadline exceeded"
2026-05-15 10:30:48.678 [d4e5f6g7] ERROR failed to execute usecase error="customer not found"
```

**ERROR level** — Errors only (minimal output, production-safe):
```
2026-05-15 10:30:48.678 [d4e5f6g7] ERROR failed to execute usecase error="customer not found"
2026-05-15 10:30:49.789 [e5f6g7h8] ERROR failed to commit kafka message error="broker unavailable"
```

---

## Architecture

This service follows **Clean / Hexagonal Architecture**. Dependencies point inward: the domain has no external imports; infrastructure adapters implement interfaces defined by the use case layer.

```mermaid
graph TD
    subgraph Adapters["Adapters (outer)"]
        HTTP["HTTP · Fiber\nHealthHandler · CustomerHandler · Router"]
        KafkaL["Kafka Listener · kafka-go\nWorker Pool"]
    end

    subgraph UC["Use Cases"]
        UCS["UpdateClientNameUseCase\n──────────────────────\nPorts (interfaces)\nCustomerRepository · IdempotencyStore"]
    end

    subgraph Domain["Domain (inner)"]
        DOM["Customer · UpdateName()\n──────────────────────\nErrCustomerNotFound\nErrInvalidCode · ErrInvalidName"]
    end

    subgraph Infra["Infrastructure"]
        DB["CustomerRepository\nGORM + PostgreSQL"]
        IDEM["MockIdempotencyStore\nreplace with Redis in production"]
    end

    HTTP -->|calls| UCS
    KafkaL -->|calls| UCS
    UCS -->|uses| DOM
    DB -. implements .-> UCS
    IDEM -. implements .-> UCS
```

**Layer order (inner → outer):** `domain → usecase → adapter → infrastructure`

### Directory Structure

```
application/
├── cmd/
│   └── api/
│       └── main.go              # Entry point: wiring, graceful shutdown
├── internal/
│   ├── domain/
│   │   └── customer.go          # Customer entity + sentinel errors
│   ├── usecase/
│   │   ├── ports.go             # CustomerRepository + IdempotencyStore interfaces
│   │   └── update_client_name.go # Core business logic
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── handler/         # HealthHandler, CustomerHandler
│   │   │   └── router/          # Route registration
│   │   └── kafka/
│   │       └── listener.go      # Kafka consumer worker pool
│   ├── infrastructure/
│   │   ├── database/            # CustomerRepository (GORM)
│   │   └── idempotency/         # MockIdempotencyStore
│   ├── config/
│   │   └── config.go            # Env var loading
│   └── logging/
│       └── structured_logger.go # Structured logging with traceId
├── tests/
│   ├── integration/             # Postgres container, no Kafka
│   └── acceptance/              # Full stack: Postgres + Kafka + Fiber
├── k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
└── .env.example
```

### Key Data Flows

#### HTTP — GET /api/v2/customers/:idTx

```
Client Request
    ↓
HTTP Router (Fiber)
    ↓
CustomerHandler.GetByIdTx()
    ↓
CustomerRepository.GetByCode(ctx, code)
    ↓
PostgreSQL SELECT
    ↓
Return Customer or ErrCustomerNotFound
    ↓
HTTP 200 JSON / 404 Error
```

#### Async — Kafka Customer Name Update

```
Kafka Message (customer_code, new_name)
    ↓
Listener.FetchMessage()
    ↓
UpdateClientNameUseCase.Execute()
    ↓
IdempotencyStore.AcquireLock()
    ├─ Duplicate? → Skip silently, commit offset
    └─ First time? → Process:
        ├─ Repository.GetByCode()
        ├─ Customer.UpdateName()
        ├─ Repository.Save()
        └─ Commit offset
```

#### Idempotency Design

`IDEMPOTENCY_LOCK_TTL` (default: 30 seconds) controls the lock duration. When a message is received, a lock is acquired for that duration. If the same message is redelivered (Kafka at-least-once), `AcquireLock` returns `false` and the duplicate is silently dropped. If processing fails, the lock is released so the message can be retried.

---

## Environment Variables Reference

`.env.local` is used by `make run` and `make smoke`. It points to the Docker-hosted infrastructure via `localhost`:

| Variable | Local Value | Purpose |
|----------|------------|---------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/go_backend?sslmode=disable` | PostgreSQL connection string |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker address (host access for `make run`) |
| `KAFKA_TOPIC` | `events` | Kafka topic to consume customer updates from |
| `KAFKA_GROUP_ID` | `ms-casino-go-v2` | Kafka consumer group ID (default: project slug) |
| `KAFKA_WORKERS` | `5` | Number of parallel workers processing Kafka messages |
| `PORT` | `8081` | HTTP server listen port |
| `APP_ENV` | `development` | Controls log format: `development` (text) or `production` (JSON) |
| `LOG_LEVEL` | `INFO` | Log verbosity: `DEBUG`, `INFO`, `WARN`, or `ERROR` |
| `IDEMPOTENCY_LOCK_TTL` | `30` | Lock duration (seconds) for duplicate message detection |

When running in **Full Docker mode**, the app container uses the internal addresses (`kafka:29092`, `postgres:5432`) configured directly in `docker-compose.yml`.

See `.env.example` for all available variables and their descriptions.

---

## Story-Driven Development

This project uses **AIOX Story-Driven Development (SDD)** — a structured workflow where every feature or enhancement starts with a story file in `docs/stories/`.

### Creating a Story

```bash
@sm *create-story
```

This creates a story file (e.g., `docs/stories/active/1.1.story.md`) with:
- **Description** — What needs to be built and why
- **Acceptance Criteria** — Testable checkboxes (`[ ]`) defining success
- **File List** — Track which files you modify
- **Status** — Lifecycle (Draft → Ready → InProgress → InReview → Done)

### Story Status Workflow

```
Draft         → @po validates (10-point checklist)
    ↓
Ready         → @dev implements (feature branch: feat/epic-{N}/story-{N.M}-{slug})
    ↓
InProgress    → work on acceptance criteria, commit with [Story N.M] tag
    ↓
InReview      → @qa tests and verifies completion
    ↓
Done          → @devops creates PR and merges to main
```

### Quick Checklist Before Pushing

- [ ] All acceptance criteria checkboxes marked (or deferred with note)
- [ ] File List updated with modified files
- [ ] Story Status updated to `InProgress` or `InReview`
- [ ] Commit messages use format: `[Story N.M] feat: description`
- [ ] Local tests pass: `make test-all`

---

## Troubleshooting

### PostgreSQL connection refused

```bash
docker-compose logs postgres
docker-compose restart postgres
```

### Kafka not responding

```bash
docker-compose logs kafka
```

### Application won't start

```bash
# Check environment is loaded
echo $DATABASE_URL

# Free occupied port
fuser -k 8081/tcp
```

### App container won't connect to Kafka

The app container connects to `kafka:29092` (internal listener). If you see connection refused, check that the `INTERNAL` listener is configured in `docker-compose.yml` and that the kafka healthcheck has passed before the app starts.

```bash
docker-compose logs kafka | grep "INTERNAL"
docker-compose ps   # kafka should show (healthy)
```
