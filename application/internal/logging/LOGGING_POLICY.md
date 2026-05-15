# Logging Policy — TEMPL-001.4

**Effective:** 2026-05-13

## Overview

This document defines the logging level policy for all services using the StructuredLogger. The goal is consistent observability and appropriate error visibility: never hide errors at DEBUG level.

---

## Log Level Rules

| Level | When to Use | Examples |
|-------|------------|----------|
| **INFO** | Request/response lifecycle, successful operations | "customer lookup request received", "kafka message received", "order processed successfully" |
| **WARNING** | Recoverable errors, degraded states, validation failures | "customer not found (404)", "kafka message unmarshal failed", "retry after backoff" |
| **ERROR** | Unrecoverable failures, system errors | "database connection lost", "failed to save to database", "http server error" |
| **DEBUG** | Infrastructure details, internal operations (not errors) | "executing database query", "cache hit", "middleware processing" |

**CRITICAL RULE:** Never hide errors at DEBUG level. All error/failure logs must be WARNING or ERROR.

---

## HTTP Handlers

**Rule:** Log incoming requests and final responses at **INFO** level. Errors at **WARNING/ERROR**.

```go
func (h *CustomerHandler) GetByIdTx(c *fiber.Ctx) error {
    ctx := c.UserContext()
    idTx := c.Params("idTx")

    // 1. Log incoming REQUEST (INFO level) - use map for simple contextual data
    h.logger.Info(ctx, "customer lookup request received", map[string]interface{}{
        "code": idTx,
    })

    // 2. Process request
    customer, err := h.repo.GetByCode(ctx, idTx)
    
    // 3. Handle errors at WARNING/ERROR level (NOT DEBUG!)
    if err != nil {
        if errors.Is(err, domain.ErrCustomerNotFound) {
            // Not found is WARNING level (noteworthy but expected)
            h.logger.Warn(ctx, "customer not found", map[string]interface{}{
                "code": idTx,
            })
            return c.Status(fiber.StatusNotFound).JSON(...)
        }
        // Database errors are ERROR level (actual failures)
        h.logger.Error(ctx, "failed to fetch customer from repository", map[string]interface{}{
            "code": idTx,
            "error": err.Error(),
        })
        return c.Status(fiber.StatusInternalServerError).JSON(...)
    }

    // 4. Log successful RESPONSE (INFO level)
    response := customerResponse{...}
    h.logger.Info(ctx, "customer lookup response", response)

    return c.Status(fiber.StatusOK).JSON(response)
}
```

---

## Kafka Listeners

**Rule:** Log received messages and processing results at **INFO** level. All errors at **WARNING/ERROR**.

```go
func (l *Listener) handle(ctx context.Context, msg kafkago.Message) {
    // Extract traceId (set up in context)
    traceId, err := extractOrGenerateTraceId(msg)
    if err != nil {
        // Invalid header format is WARNING level
        l.logger.Warn(ctx, "invalid trace id", map[string]interface{}{
            "error": err.Error(),
        })
        return
    }

    msgCtx := context.WithValue(ctx, "traceId", traceId)

    // Parse message
    var payload messagePayload
    if err := json.Unmarshal(msg.Value, &payload); err != nil {
        // Unmarshal failure is ERROR level
        l.logger.Error(ctx, "failed to unmarshal message", map[string]interface{}{
            "error": err.Error(),
        })
        return
    }

    // 1. Log INCOMING MESSAGE (INFO level)
    l.logger.Info(msgCtx, "kafka message received", payload)

    // 2. Process message
    if err := l.uc.Execute(msgCtx, payload.CustomerCode, payload.CustomerName); err != nil {
        // Processing failure is ERROR level
        l.logger.Error(ctx, "processing failed", map[string]interface{}{
            "error": err.Error(),
        })
        return
    }

    // 3. Log SUCCESSFUL PROCESSING (INFO level)
    l.logger.Info(msgCtx, "kafka message processed successfully", payload)
}
```

---

## Application Lifecycle

**Rule:** Startup/shutdown at **INFO**. Errors at **ERROR**.

```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    baseLogger.Info("listening on server", "addr", addr, "env", cfg.AppEnv)
    
    if err := app.Listen(addr); err != nil {
        baseLogger.Error("http server error", "error", err)
    }

    baseLogger.Info("shutting down")
}
```

---

## Error vs Infrastructure Logging

**When to use each level:**

| Situation | Level | Example |
|-----------|-------|---------|
| Operation failed | ERROR | "failed to fetch from database" |
| Validation failed | WARNING | "invalid email format" |
| Not found (expected) | WARNING | "customer not found" |
| Infrastructure event | DEBUG | "connection pool size: 10" |
| Internal operation | DEBUG | "executing query" |

---

## Logging with Struct Objects

**Rule:** Log using flexible types. Avoid creating objects just for logging.

```go
// ✅ GOOD: Log the actual response object
response := customerResponse{ID: 1, Code: "CUST-001", Name: "John"}
h.logger.Info(ctx, "customer created", response)

// ✅ GOOD: Use a map for simple contextual data
h.logger.Info(ctx, "customer lookup request received", map[string]interface{}{
    "code": idTx,
})

// ✗ AVOID: Creating an object just to log it
request := customerRequest{Code: idTx}  // ← Unnecessary object creation
h.logger.Info(ctx, "customer lookup request received", request)
```

---

## TraceId Propagation

**All logs must include traceId** through the request context:

```go
// HTTP: TraceId automatically set by middleware
ctx := c.UserContext() // contains traceId
h.logger.Info(ctx, "message", object) // traceId included automatically

// Kafka: TraceId manually set
msgCtx := context.WithValue(ctx, "traceId", traceId)
l.logger.Info(msgCtx, "message", payload) // traceId included automatically
```

---

## Common Patterns

### Pattern 1: Request → Response

```go
// Incoming
h.logger.Info(ctx, "request received", incomingStruct)
// ... process ...
// Outgoing
h.logger.Info(ctx, "response sent", outgoingStruct)
```

### Pattern 2: Async Processing with Errors

```go
// Start
h.logger.Info(ctx, "async job started", jobRequest)
// ... process ...
// Error (debug level)
if err != nil {
    h.logger.Debug(ctx, "async job failed", map[string]interface{}{
        "error": err.Error(),
    })
}
// Success
h.logger.Info(ctx, "async job completed", result)
```

### Pattern 3: Integration/External Calls

```go
// Before call
h.logger.Debug(ctx, "calling external service", map[string]interface{}{
    "service": "datadog",
    "endpoint": "/v1/logs",
})
// After success
h.logger.Debug(ctx, "external service response received", response)
// On error
h.logger.Debug(ctx, "external service error", map[string]interface{}{
    "error": err.Error(),
    "status": httpStatus,
})
```

---

## Datadog Integration

These logs are production-ready for Datadog:

```
# Search by level
status:INFO AND @message:"request received"

# Trace across services
@traceId:a1b2c3d4-e5f6-7890-abcd-ef1234567890

# Search in context fields
@context.customer_code:CUST-001
@context.error:*connection*
```

---

## Examples by Service Type

### REST API Handler
- INFO: Request received, Response sent, State changes
- DEBUG: Database errors, Validation failures, Not found (404)
- ERROR: Unhandled exceptions, 5xx errors

### Message Consumer (Kafka/SQS)
- INFO: Message received, Processing succeeded
- DEBUG: Parsing failed, Processing failed, Retry logic
- ERROR: Consumer crash, Topic unavailable

### Scheduled Job
- INFO: Job started, Job completed
- DEBUG: Step details, Resource checks
- ERROR: Job failed permanently

---

## Policy Enforcement

This policy is enforced through:
1. **Code review** — Verify log levels before merge
2. **CodeRabbit** — Automated checks for log level consistency
3. **Documentation** — This file serves as the reference

---

**Questions or clarifications?** Refer to `EXAMPLES.md` for concrete usage patterns.
