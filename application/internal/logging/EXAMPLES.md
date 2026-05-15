# Structured Logger Usage Examples

## Quick Reference

The StructuredLogger accepts any serializable type for context: maps, structs, slices, primitives.

---

## Example 1: Simple Map Context

```go
logger.Info(ctx, "User login", map[string]interface{}{
    "userId": 12345,
    "ip": "192.168.1.1",
})
```

**JSON Output:**
```json
{
  "timestamp": "2026-05-13T10:30:45.123Z",
  "level": "INFO",
  "traceId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "message": "User login",
  "context": {
    "userId": 12345,
    "ip": "192.168.1.1"
  }
}
```

---

## Example 2: Struct Object (Recommended)

```go
type CustomerResponse struct {
    ID   uint   `json:"id"`
    Code string `json:"code"`
    Name string `json:"name"`
}

response := CustomerResponse{
    ID:   12345,
    Code: "CUST-001",
    Name: "John Doe",
}

logger.Info(ctx, "Customer created", response)
```

**JSON Output:**
```json
{
  "timestamp": "2026-05-13T10:30:45.123Z",
  "level": "INFO",
  "traceId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "message": "Customer created",
  "context": {
    "id": 12345,
    "code": "CUST-001",
    "name": "John Doe"
  }
}
```

**Why this is better:**
- Type-safe (struct fields with proper types)
- Respects `json` struct tags (customization, omitempty, renaming)
- Self-documenting (clear what fields are logged)

---

## Example 3: Slice of Objects

```go
customers := []CustomerResponse{
    {ID: 1, Code: "C001", Name: "Alice"},
    {ID: 2, Code: "C002", Name: "Bob"},
}

logger.Info(ctx, "Bulk import complete", customers)
```

**JSON Output:**
```json
{
  "context": [
    {"id": 1, "code": "C001", "name": "Alice"},
    {"id": 2, "code": "C002", "name": "Bob"}
  ]
}
```

---

## Example 4: Nil Context (No Additional Data)

```go
logger.Info(ctx, "Server started", nil)
```

**JSON Output:**
```json
{
  "timestamp": "2026-05-13T10:30:45.123Z",
  "level": "INFO",
  "traceId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "message": "Server started"
}
```

Note: `context` field is omitted when nil (respects `omitempty` tag).

---

## Example 5: HTTP Handler with Response Object

```go
func (h *CustomerHandler) GetByIdTx(c *fiber.Ctx) error {
    ctx := c.UserContext()
    idTx := c.Params("idTx")

    customer, err := h.repo.GetByCode(ctx, idTx)
    if err != nil {
        h.logger.Error(ctx, "failed to fetch customer", map[string]interface{}{
            "code": idTx,
            "error": err.Error(),
        })
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
    }

    response := customerResponse{
        ID:   customer.ID,
        Code: customer.Code,
        Name: customer.Name,
    }

    // Log the successful response
    h.logger.Info(ctx, "customer fetched", response)

    return c.Status(fiber.StatusOK).JSON(response)
}
```

**JSON Output:**
```json
{
  "timestamp": "2026-05-13T10:30:45.123Z",
  "level": "INFO",
  "traceId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "message": "customer fetched",
  "context": {
    "id": 12345,
    "code": "CUST-001",
    "name": "John Doe"
  }
}
```

---

## API Methods

All methods follow the same pattern: `(ctx, message, context)`

```go
// Debug level
logger.Debug(ctx, "Debug info", object)

// Info level
logger.Info(ctx, "Information", object)

// Warn level
logger.Warn(ctx, "Warning message", object)

// Error level
logger.Error(ctx, "Error occurred", object)
```

---

## Struct Tags (JSON Customization)

Use standard Go `json` tags to customize field names, omit empty fields, etc:

```go
type UserEvent struct {
    UserID       uint      `json:"user_id"`              // Rename field
    Action       string    `json:"action"`
    Timestamp    time.Time `json:"timestamp,omitempty"` // Omit if empty
    Metadata     map[string]string `json:"-"`           // Exclude from JSON
}
```

---

## Datadog Integration

These logs are production-ready for Datadog:
- TraceId enables distributed tracing across services
- JSON format is immediately parseable by Datadog intake
- Context field preserves structured data for faceted searching

**Datadog faceted search example:**
```
@traceId:a1b2c3d4-e5f6-7890-abcd-ef1234567890
@context.userId:12345
@level:ERROR
```

---

## Development Output (Text Format)

Set `LOG_FORMAT=text` or `APP_ENV != production` to get human-readable logs:

```
2026-05-13 10:30:45.123 [a1b2c3d4] INFO  Customer created {id: 12345, code: CUST-001, name: John Doe}
```

---

## Best Practices

1. **Prefer structs over maps** — Type-safe and self-documenting
2. **Use context from c.UserContext()** — Ensures traceId is available
3. **Log on error paths** — Include error details in context
4. **Log on business events** — Customer created, order placed, etc.
5. **Avoid logging sensitive data** — Passwords, tokens, PII (unless redacted)
