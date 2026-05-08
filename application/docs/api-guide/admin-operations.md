# Casino Proxy API - Admin Operations Guide

## Overview

This guide documents all administrative API endpoints for managing operators, handling system health checks, and uploading operator-related images. These endpoints are used for operator onboarding, management, and system monitoring.

---

## Quick Start

### 1. Create an Operator

```bash
curl -X POST https://api.casino-proxy.local/v1/internal/operator/store \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "my-operator",
    "name": "My Casino",
    "url": "https://my-casino.example.com/webhooks"
  }'
```

**Response:**
```json
{
  "id": 1,
  "slug": "my-operator",
  "name": "My Casino",
  "url": "https://my-casino.example.com/webhooks",
  "is_active": true,
  "created_at": "2026-05-08T12:00:00Z",
  "updated_at": "2026-05-08T12:00:00Z"
}
```

### 2. List All Operators

```bash
curl -X GET https://api.casino-proxy.local/v1/internal/operator/ \
  -H "Authorization: Bearer {token}"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "slug": "my-operator",
      "name": "My Casino",
      "url": "https://my-casino.example.com/webhooks",
      "is_active": true,
      "created_at": "2026-05-08T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "per_page": 50,
    "current_page": 1,
    "last_page": 1
  }
}
```

### 3. Delete an Operator

```bash
curl -X DELETE https://api.casino-proxy.local/v1/internal/operator/delete/my-operator \
  -H "Authorization: Bearer {token}"
```

**Response:**
```json
{
  "message": "Operator deleted successfully",
  "id": 1
}
```

---

## Authentication

### Admin Access Control

All admin API endpoints require the `internalCheck` middleware, which verifies:
1. Request from internal/authorized system
2. Valid bearer token (if applicable)
3. Administrator role/permissions

**Bearer Token Format:**
```
Authorization: Bearer {sanctum_token}
```

### Security Headers

Always include when making admin API requests:
```
Content-Type: application/json
Authorization: Bearer {token}
X-Requested-With: XMLHttpRequest (for CSRF protection)
```

---

## Operator Management API

### Create Operator

**Endpoint:** `POST /v1/internal/operator/store`

**Purpose:** Create a new operator/casino account in the system.

**Authentication:** Required (internalCheck middleware)

**Request Body:**
```json
{
  "slug": "string",        // Required: unique operator identifier (alphanumeric, lowercase)
  "name": "string",        // Required: display name
  "url": "string"          // Required: webhook callback URL (HTTPS)
}
```

**Field Validation Rules:**
| Field | Rules | Example |
|-------|-------|---------|
| slug | Unique, 3-50 chars, lowercase alphanumeric + hyphens, required | "my-operator" |
| name | 1-255 chars, required | "My Casino" |
| url | Valid HTTPS URL, required | "https://my-casino.example.com/webhooks" |

**Success Response (201 Created):**
```json
{
  "id": 1,
  "slug": "my-operator",
  "name": "My Casino",
  "url": "https://my-casino.example.com/webhooks",
  "is_active": true,
  "created_at": "2026-05-08T12:00:00Z",
  "updated_at": "2026-05-08T12:00:00Z"
}
```

**Error Responses:**

**409 Conflict** - Slug already exists:
```json
{
  "error": "DUPLICATE_SLUG",
  "message": "Operator with slug 'my-operator' already exists",
  "field": "slug"
}
```

**400 Bad Request** - Validation failed:
```json
{
  "error": "VALIDATION_ERROR",
  "message": "Validation failed",
  "errors": {
    "slug": ["Slug must be unique", "Slug must be lowercase alphanumeric"],
    "url": ["URL must be HTTPS"]
  }
}
```

**401 Unauthorized** - Missing/invalid token:
```json
{
  "error": "UNAUTHORIZED",
  "message": "Authentication required"
}
```

**503 Server Error** - Database issue:
```json
{
  "error": "SERVER_ERROR",
  "message": "Failed to create operator"
}
```

### List Operators

**Endpoint:** `GET /v1/internal/operator/`

**Purpose:** Retrieve all operators with optional filtering and pagination.

**Authentication:** Required (internalCheck middleware)

**Query Parameters:**
| Parameter | Type | Optional | Default | Description |
|-----------|------|----------|---------|-------------|
| page | integer | Yes | 1 | Page number for pagination |
| per_page | integer | Yes | 50 | Records per page (max 100) |
| search | string | Yes | — | Search by slug or name (partial match) |
| active_only | boolean | Yes | false | Only return active operators |

**Request Examples:**
```bash
# List all operators (first 50)
GET /v1/internal/operator/

# List with pagination
GET /v1/internal/operator/?page=2&per_page=25

# Search by name
GET /v1/internal/operator/?search=casino

# Only active operators
GET /v1/internal/operator/?active_only=true
```

**Success Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "slug": "my-operator",
      "name": "My Casino",
      "url": "https://my-casino.example.com/webhooks",
      "is_active": true,
      "created_at": "2026-05-08T12:00:00Z",
      "updated_at": "2026-05-08T12:00:00Z"
    },
    {
      "id": 2,
      "slug": "another-operator",
      "name": "Another Casino",
      "url": "https://another-casino.example.com/webhooks",
      "is_active": true,
      "created_at": "2026-05-08T12:30:00Z",
      "updated_at": "2026-05-08T12:30:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "per_page": 50,
    "current_page": 1,
    "last_page": 1,
    "from": 1,
    "to": 2
  }
}
```

**Error Responses:**

**400 Bad Request** - Invalid pagination:
```json
{
  "error": "INVALID_PARAMETER",
  "message": "per_page must be <= 100"
}
```

**401 Unauthorized** - Missing authentication:
```json
{
  "error": "UNAUTHORIZED",
  "message": "Authentication required"
}
```

### Delete Operator

**Endpoint:** `DELETE /v1/internal/operator/delete/{slug}`

**Purpose:** Remove an operator from the system.

**Authentication:** Required (internalCheck middleware)

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| slug | string | Operator slug (must exist) |

**Request Example:**
```bash
DELETE /v1/internal/operator/delete/my-operator
```

**Success Response (200 OK):**
```json
{
  "message": "Operator deleted successfully",
  "id": 1,
  "slug": "my-operator"
}
```

**Error Responses:**

**404 Not Found** - Operator doesn't exist:
```json
{
  "error": "NOT_FOUND",
  "message": "Operator with slug 'my-operator' not found"
}
```

**400 Bad Request** - Cannot delete (cascade rules):
```json
{
  "error": "CONSTRAINT_VIOLATION",
  "message": "Cannot delete operator: still has active transactions",
  "details": {
    "active_transactions": 5,
    "action": "Mark operator as inactive instead"
  }
}
```

**401 Unauthorized** - Missing authentication:
```json
{
  "error": "UNAUTHORIZED",
  "message": "Authentication required"
}
```

---

## System Health Check API

### Health Status

**Endpoint:** `GET /v1/status`

**Purpose:** Check system health and readiness.

**Authentication:** Not required

**Request Example:**
```bash
curl -X GET https://api.casino-proxy.local/v1/status
```

**Success Response (200 OK):**
```json
{
  "message": "OK",
  "status": "healthy",
  "timestamp": "2026-05-08T12:00:00Z"
}
```

**Extended Health Check (Optional):**
```json
{
  "message": "OK",
  "status": "healthy",
  "timestamp": "2026-05-08T12:00:00Z",
  "services": {
    "database": "healthy",
    "cache": "healthy",
    "webhooks": "healthy"
  },
  "uptime_seconds": 86400,
  "version": "1.0.0"
}
```

**Unhealthy Response (503 Service Unavailable):**
```json
{
  "message": "Service unhealthy",
  "status": "unhealthy",
  "timestamp": "2026-05-08T12:00:00Z",
  "services": {
    "database": "unavailable",
    "cache": "healthy"
  }
}
```

---

## Session Entry API

### Admin Dashboard Entry Point

**Endpoint:** `POST /v1/entry`

**Purpose:** Create session for admin dashboard access.

**Authentication:** Optional (depends on implementation)

**Request Body:**
```json
{
  "email": "string",        // Admin email
  "password": "string",     // Admin password
  "remember": "boolean"     // Optional: remember login
}
```

**Success Response (200 OK):**
```json
{
  "message": "Authenticated",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "name": "Admin User"
  },
  "token": "token_string",  // Optional: Sanctum token
  "redirect": "/dashboard"  // Optional: redirect URL
}
```

**Error Responses:**

**401 Unauthorized** - Invalid credentials:
```json
{
  "error": "INVALID_CREDENTIALS",
  "message": "Invalid email or password"
}
```

**429 Too Many Requests** - Too many login attempts:
```json
{
  "error": "TOO_MANY_ATTEMPTS",
  "message": "Too many login attempts. Please try again in 15 minutes."
}
```

---

## Image Upload API

### Upload Operator Image

**Endpoint:** `POST /v1/image`

**Purpose:** Upload an image file for operator (logo, banner, etc.)

**Authentication:** Required (internalCheck middleware)

**Content-Type:** `multipart/form-data`

**Request Parameters:**
| Field | Type | Optional | Description |
|-------|------|----------|-------------|
| image | file | Required | Image file (see accepted types below) |
| type | string | Optional | Image type/category (e.g., "logo", "banner") |
| operator_slug | string | Optional | Operator slug to associate image |

**Accepted File Types:**
- JPEG/JPG (`.jpg`, `.jpeg`)
- PNG (`.png`)
- WebP (`.webp`)
- GIF (`.gif`)

**File Size Limits:**
- Maximum: 5 MB per file
- Minimum: 100 bytes

**Request Example:**
```bash
curl -X POST https://api.casino-proxy.local/v1/image \
  -H "Authorization: Bearer {token}" \
  -F "image=@/path/to/logo.png" \
  -F "type=logo" \
  -F "operator_slug=my-operator"
```

**Success Response (200 OK):**
```json
{
  "message": "Image uploaded successfully",
  "image": {
    "id": 1,
    "filename": "logo-1715241600.png",
    "url": "https://cdn.casino-proxy.local/images/logo-1715241600.png",
    "size": 25600,
    "type": "logo",
    "operator_slug": "my-operator",
    "uploaded_at": "2026-05-08T12:00:00Z"
  }
}
```

**Error Responses:**

**400 Bad Request** - Invalid file type:
```json
{
  "error": "INVALID_FILE_TYPE",
  "message": "File type not accepted",
  "accepted_types": ["jpg", "jpeg", "png", "webp", "gif"]
}
```

**413 Payload Too Large** - File exceeds size limit:
```json
{
  "error": "FILE_TOO_LARGE",
  "message": "File size exceeds maximum allowed (5 MB)",
  "max_size_mb": 5,
  "file_size_mb": 12
}
```

**415 Unsupported Media Type** - No file provided:
```json
{
  "error": "MISSING_FILE",
  "message": "No file provided in 'image' field"
}
```

**401 Unauthorized** - Missing authentication:
```json
{
  "error": "UNAUTHORIZED",
  "message": "Authentication required"
}
```

**507 Insufficient Storage** - Storage full:
```json
{
  "error": "STORAGE_FULL",
  "message": "Image storage is full"
}
```

---

## Common Workflows

### Workflow 1: Onboard a New Operator

1. **Create Operator Account**
   ```bash
   POST /v1/internal/operator/store
   {
     "slug": "new-casino",
     "name": "New Casino",
     "url": "https://new-casino.example.com/webhooks"
   }
   ```

2. **Upload Operator Logo**
   ```bash
   POST /v1/image
   -F "image=@logo.png"
   -F "type=logo"
   -F "operator_slug=new-casino"
   ```

3. **Verify in Dashboard**
   - Navigate to operator list
   - Confirm operator appears with logo

### Workflow 2: Manage Operator Webhook URL

**Issue:** Operator changes webhook URL

**Solution:**
1. Delete old operator
   ```bash
   DELETE /v1/internal/operator/delete/old-slug
   ```

2. Create new operator with new URL
   ```bash
   POST /v1/internal/operator/store
   {
     "slug": "old-slug",
     "name": "Casino Name",
     "url": "https://new-webhook-url.example.com/webhooks"
   }
   ```

**Note:** Consider adding an Update endpoint in future for better UX.

### Workflow 3: Disable an Operator Temporarily

**Issue:** Operator needs temporary suspension

**Solution:**
1. List all operators
   ```bash
   GET /v1/internal/operator/?search=casino-name
   ```

2. Check system status
   ```bash
   GET /v1/status
   ```

3. Plan deletion timing (avoid peak traffic)

4. Delete operator
   ```bash
   DELETE /v1/internal/operator/delete/casino-slug
   ```

---

## Error Handling Best Practices

### Standard HTTP Status Codes

| Code | Meaning | Retry? | Action |
|------|---------|--------|--------|
| 200 | Success | No | Process response |
| 201 | Created | No | Resource created successfully |
| 400 | Bad Request | No | Fix request parameters/body |
| 401 | Unauthorized | No | Check authentication token |
| 403 | Forbidden | No | Check permissions/authorization |
| 404 | Not Found | No | Verify resource slug/id |
| 409 | Conflict | No (idempotent) | Handle duplicate or constraint |
| 413 | Payload Too Large | No | Reduce file size |
| 500 | Server Error | Yes (with backoff) | Retry after delay |
| 503 | Service Unavailable | Yes (with backoff) | Retry after delay |

### Handling Validation Errors

**Example 1: Field-level validation**
```json
{
  "error": "VALIDATION_ERROR",
  "errors": {
    "slug": ["Slug already exists", "Slug must be lowercase"],
    "url": ["URL must be HTTPS"]
  }
}
```

**Recovery:**
```python
for field, messages in errors.items():
    for message in messages:
        log(f"Field {field}: {message}")
        # Show to user or fix programmatically
```

**Example 2: Business logic error**
```json
{
  "error": "CONSTRAINT_VIOLATION",
  "message": "Cannot delete: operator has active transactions",
  "details": {
    "operator_id": 1,
    "active_transactions": 5
  }
}
```

**Recovery:**
```
1. Ask user to wait for transactions to complete
2. Or: Mark operator as inactive instead of deleting
3. Or: Archive operator without deleting
```

### Handling Server Errors

**Pattern:**
```python
import time

def create_operator_with_retry(slug, name, url, max_retries=3):
    for attempt in range(max_retries):
        try:
            response = requests.post(
                "/v1/internal/operator/store",
                json={"slug": slug, "name": name, "url": url},
                timeout=10
            )
            
            if response.status_code == 201:
                return response.json()
            elif response.status_code >= 500:
                # Server error, retry with backoff
                wait_seconds = 2 ** attempt  # 1s, 2s, 4s
                time.sleep(wait_seconds)
            else:
                # Client error, don't retry
                raise Exception(f"Error {response.status_code}: {response.text}")
                
        except requests.RequestException as e:
            # Connection error, retry
            if attempt < max_retries - 1:
                time.sleep(2 ** attempt)
            else:
                raise
    
    raise Exception(f"Failed after {max_retries} attempts")
```

---

## Monitoring Admin API

### Key Metrics

| Metric | Alert Threshold | Action |
|--------|-----------------|--------|
| Error rate | > 5% | Check authentication, API server |
| Response time p95 | > 1s | Check database, check load |
| Unauthorized (401) | > 10/min | Check token generation |
| Validation errors (400) | > 20/min | Check client, review API changes |
| Duplicate slugs (409) | > 1/min | Check client logic, fix naming |

### Logging Template

```go
type AdminAPILog struct {
    Timestamp    time.Time
    Method       string
    Endpoint     string
    Status       int
    Error        string
    RequestBody  string  // Sanitized
    ResponseTime int     // ms
    UserID       int
    OperatorID   int
}

func logAdminAPICall(method, endpoint string, status int, duration time.Duration) {
    log := AdminAPILog{
        Timestamp:    time.Now(),
        Method:       method,
        Endpoint:     endpoint,
        Status:       status,
        ResponseTime: int(duration.Milliseconds()),
    }
    
    db.Create(&log)
    
    if status >= 500 {
        alertSlack(log)
    }
}
```

---

## Testing Admin Endpoints

### Unit Test Example

```go
func TestCreateOperator(t *testing.T) {
    body := []byte(`{
        "slug": "test-op",
        "name": "Test Operator",
        "url": "https://test.example.com/webhooks"
    }`)
    
    req := httptest.NewRequest("POST", "/v1/internal/operator/store", strings.NewReader(string(body)))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer test_token")
    
    w := httptest.NewRecorder()
    handler(w, req)
    
    assert.Equal(t, 201, w.Code)
    
    var response OperatorResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "test-op", response.Slug)
}
```

### Integration Test Example

```go
func TestOperatorLifecycle(t *testing.T) {
    // Create
    operator := createOperator("lifecycle-test", "Lifecycle Test")
    assert.NotZero(t, operator.ID)
    
    // List
    operators := listOperators()
    assert.Len(t, operators, 1)
    assert.Equal(t, "lifecycle-test", operators[0].Slug)
    
    // Delete
    deleteOperator("lifecycle-test")
    
    // Verify deleted
    operators = listOperators()
    assert.Len(t, operators, 0)
}
```

---

## Summary

| Endpoint | Method | Purpose | Auth |
|----------|--------|---------|------|
| `/v1/internal/operator/` | GET | List operators | Required |
| `/v1/internal/operator/store` | POST | Create operator | Required |
| `/v1/internal/operator/delete/{slug}` | DELETE | Delete operator | Required |
| `/v1/status` | GET | Health check | Not required |
| `/v1/entry` | POST | Dashboard entry | Optional |
| `/v1/image` | POST | Upload image | Required |
