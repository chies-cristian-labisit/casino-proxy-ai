# Local Development Guide

This guide walks you through running your generated microservice locally. Infrastructure (PostgreSQL + Kafka) runs via Docker Compose; the application runs directly on your machine.

---

## Prerequisites

- Go 1.25.4+
- Docker & Docker Compose
- `make` (optional, for convenience commands)
- `psql` (optional, for database connectivity checks)

---

## Setup

### Step 1 — Resolve dependencies

```bash
cd application
go mod tidy
```

### Step 2 — Start infrastructure

```bash
docker-compose up -d
```

Verify both services are healthy before continuing:

```bash
docker-compose ps
```

Expected output:
```
CONTAINER ID   IMAGE                              STATUS
...            postgres:16-alpine                 Up (healthy)
...            confluentinc/confluent-local:7.5.0 Up (healthy)
```

### Step 3 — Load environment variables

**Linux / macOS:**

```bash
set -a
source .env.local
set +a
```

**Windows (PowerShell):**

```powershell
Get-Content .env.local | ForEach-Object {
    $name, $value = $_ -split '=';
    if ($name -and -not $name.StartsWith('#')) {
        [Environment]::SetEnvironmentVariable($name, $value)
    }
}
```

### Step 4 — Initialize the database

First time only — starts the app briefly so GORM auto-migrates the schema:

```bash
go run ./cmd/api &
sleep 5
pkill -f "go run"
```

### Step 5 — Run the application

```bash
make run
```

Or without Make:

```bash
go run ./cmd/api
```

### Step 6 — Verify

```bash
make smoke
```

Or manually:

```bash
curl http://localhost:8081/liveness   # → 200 OK
curl http://localhost:8081/readiness  # → 200 OK (DB reachable)
```

---

## Service URLs

| Service | Address |
|---------|---------|
| Application | `http://localhost:8081` |
| PostgreSQL | `localhost:5432` |
| Kafka | `localhost:9092` |

---

## Testing

### Run all tests

```bash
make test-all       # unit → integration → acceptance
```

### Run individually

```bash
make test           # Unit tests (./internal/...) — no Docker required
make integration    # Integration tests — PostgreSQL container only
make acceptance     # Acceptance tests — full stack (PostgreSQL + Kafka + Fiber)
```

Acceptance tests spin up their own containers via Testcontainers — `docker-compose` does **not** need to be running.

### Test a specific endpoint

```bash
# Get customer by code
curl http://localhost:8081/api/v1/customers/BR123456789
```

---

## Stopping Services

```bash
make down           # Stop and remove containers (keep volumes)
make clean          # Stop and remove containers + volumes (deletes all data)
```

Or with Docker Compose directly:

```bash
docker-compose stop          # Stop only (keep containers + volumes)
docker-compose down          # Remove containers (keep volumes)
docker-compose down -v       # Remove containers + volumes
```

---

## Development Workflow

1. **Edit code** → restart with `make run`
2. **Unit test** → `make test`
3. **Full stack test** → `make acceptance`
4. **Check health** → `make smoke`

---

## Story-Driven Development

This project uses **AIOX Story-Driven Development (SDD)** — a structured workflow where every feature or enhancement starts with a story file in `docs/stories/`.

### Creating a Story

To create a development story:

```bash
@sm *create-story
```

Or activate the Scrum Master agent:

```
@sm
*create-story
```

This creates a story file (e.g., `docs/stories/1.1.story.md`) with:
- **Description** — What needs to be built and why
- **Acceptance Criteria** — Testable checkboxes (`[ ]`) defining success
- **File List** — Track which files you modify
- **Status** — Lifecycle (Draft → Ready → InProgress → InReview → Done)
- **Decision Log** — Document architectural decisions

### Story Status & Workflow

```
Draft         → @po validates (10-point checklist)
    ↓
Ready         → @dev implements (feature branch: feat/epic-{N}/story-{N.M}-{slug})
    ↓
InProgress    → Work on acceptance criteria, commit with [Story N.M] tag
    ↓
InReview      → @qa tests and verifies completion
    ↓
Done          → @devops creates PR and merges to main
```

### Marking Progress

As you work, update the story file:

```markdown
## Acceptance Criteria

- [x] Implement customer endpoint      # ← Checked when done
- [x] Add database migration
- [ ] Add integration test              # ← Still todo
```

### Three-Layer Quality Gates

Every story passes through **3 quality gates** before merging:

#### Layer 1: Pre-Commit (Local)
**When:** Before `git commit`  
**What runs:** ESLint (code style), TypeScript (type checking)  
**Purpose:** Catch style and type errors early  
**How to pass:** Run locally before committing
```bash
npm run lint        # Fix linting issues
npm run typecheck   # Fix TypeScript errors
```

#### Layer 2: Pre-Push (Local)
**When:** Before `git push origin feat/...`  
**What runs:** Story validation (checkboxes, file list, status)  
**Purpose:** Ensure story metadata is complete before pushing  
**How to pass:** Fill out story file completely
```bash
# Pre-push hook validates:
# - All checkboxes filled (or intentionally deferred)
# - File List is updated with modified files
# - Status is InProgress or InReview
# - Story ID matches branch name
```

#### Layer 3: CI/CD (GitHub Actions)
**When:** After `git push`, during PR review  
**What runs:** Full test suite (unit, integration, acceptance)  
**Purpose:** Verify functionality on clean environment  
**How to pass:** All tests pass
```bash
# Runs automatically on GitHub
# Check PR status for test results
```

**If any gate fails:**
- **Layer 1:** Fix linting/types locally, re-commit
- **Layer 2:** Update story file, re-push
- **Layer 3:** @dev fixes issues, re-pushes to feature branch

---

### Example Story: Add a New API Endpoint

Here's a complete example story (you'll create similar ones):

```markdown
# Story 2.5: Add GET Customer by Email Endpoint

**Epic:** Customer Service API Extensions (Phase 2)  
**Status:** InProgress  
**Assigned to:** @dev  

## Description

Add a new REST endpoint to fetch customer by email address.
Needed for user profile lookups and customer onboarding flows.

## Acceptance Criteria

- [x] Create GET /api/v1/customers/by-email/{email} endpoint
- [x] Add database query by email_column (indexed)
- [x] Return 200 with customer JSON if found
- [x] Return 404 with error JSON if not found
- [x] Validate email format before querying
- [ ] Add integration test (Postgres container)
- [ ] Add acceptance test (full stack: Postgres + Kafka + Fiber)

## File List

- `application/internal/adapter/http/handler/customer_handler.go` ✅ (added handler)
- `application/internal/usecase/ports.go` ✅ (added interface)
- `application/internal/infrastructure/database/customer_repository.go` ✅ (added query)
- `application/tests/integration/customer_repository_test.go` (pending)
- `application/tests/acceptance/customer_endpoint_test.go` (pending)

## Decision Log

- **2026-04-30:** Used email_column as unique identifier (existing schema)
- **2026-04-30:** Added format validation in handler (no external libs)

## Technical Notes

Email query uses GORM `Where("email = ?", email)` with indexed column.
Response format matches existing customer endpoint for consistency.
```

---

### Agent Commands Reference

**@sm (Scrum Master — River)**
- `*create-story` — Create a new development story
- `*help` — Show all Scrum Master commands

**@dev (Developer — Dex)**
- Activate: `@dev` then use `*help` to see available commands
- Use for: Code implementation, bug fixes, refactoring
- Submits work to @qa when acceptance criteria met

**@qa (QA Lead — Quinn)**
- `*qa-gate {storyId}` — Test a completed story
- Verdict: PASS (ready to merge) / CONCERNS (proceed with note) / FAIL (needs fixes)
- If FAIL: @dev fixes issues, then @qa re-tests

**@devops (DevOps — Gage)**
- Exclusive authority: `*push` (creates PR and merges to main)
- Only agent who can push to main branch

---

### Quick Checklist Before Pushing

Before `git push origin feat/epic-{N}/story-{N.M}-{slug}`:

- [ ] All acceptance criteria checkboxes marked (or deferred with note)
- [ ] File List updated with modified files
- [ ] Story Status updated to `InProgress` or `InReview`
- [ ] Commit messages use format: `[Story N.M] feat: description`
- [ ] Local tests pass: `make test-all`
- [ ] No linting errors: `npm run lint` (if applicable)
- [ ] No type errors: `npm run typecheck` (if applicable)

---

## Troubleshooting

### PostgreSQL connection refused

```bash
docker-compose logs postgres
docker-compose restart postgres
psql -h localhost -U postgres -d go_backend -c "SELECT 1"
```

### Kafka not responding

```bash
docker-compose logs kafka
docker-compose exec kafka kafka-broker-api-versions.sh --bootstrap-server localhost:9092
```

### Application won't start

```bash
# Confirm environment is loaded
echo $DATABASE_URL   # should show postgres connection string

# Free ports if occupied
fuser -k 8081/tcp   # Application
fuser -k 5432/tcp                 # PostgreSQL
fuser -k 9092/tcp                 # Kafka
```

---

## Environment Variables Reference

See `.env.example` for all available variables. `.env.local` provides defaults for local development.

| Variable | Local Value | Purpose |
|----------|------------|---------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/go_backend?sslmode=disable` | PostgreSQL connection |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker address |
| `KAFKA_TOPIC` | `events` | Topic for customer update events |
| `PORT` | `8081` | HTTP server port |
| `APP_ENV` | `development` | Enables development-mode logging |
