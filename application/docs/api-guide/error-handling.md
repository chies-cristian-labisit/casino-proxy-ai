# Casino Proxy API - Error Handling Guide

## Overview

This guide documents error scenarios across all 8 gaming providers, provider-specific error codes, recovery strategies, and best practices for handling failures in the gateway.

---

## HTTP Status Codes

### Gateway-Level Status Codes

| Code | Meaning | Action |
|------|---------|--------|
| **200** | Request successful | Process response normally |
| **400** | Bad request format | Check payload structure, required fields |
| **401** | Authentication failed | Verify credentials, signature, or operator context |
| **403** | Access forbidden | Operator not enabled or not authorized |
| **409** | Conflict (idempotency) | Transaction already processed, return cached response |
| **500** | Server error | Retry with exponential backoff |
| **502** | Bad gateway | Provider returned error, check provider logs |
| **503** | Service unavailable | Provider timeout, retry later |
| **504** | Gateway timeout | Provider not responding, escalate |

### Provider-Level Status Codes

Each provider can return provider-specific error codes within their response body.

---

## Error Response Format

### Standard Gateway Error

```json
{
  "error": "ERROR_CODE",
  "description": "Human-readable error message",
  "provider": "pragmatic_play",
  "operator_id": 1,
  "transaction_id": "txn_abc123",
  "timestamp": "2026-05-08T12:00:00Z"
}
```

### Provider-Specific Error (in response body)

Each provider's error response format varies. See sections below.

---

## Provider Error Codes & Recovery

### 1. Pragmatic Play

**HTTP Status:** Embedded in response body

| Error Code | Meaning | Recovery |
|-----------|---------|----------|
| INSUFFICIENT_BALANCE | Player lacks funds | Reject bet, return balance |
| SESSION_EXPIRED | Session token invalid | Request new session via authentication endpoint |
| INVALID_SESSION | Session doesn't exist | Create new session |
| DUPLICATE_TRANSACTION | Transaction ID already processed | Return cached response (idempotency) |
| INVALID_GAME | Game doesn't exist | Validate game_id before placing bet |
| OPERATOR_DISABLED | Operator not active | Enable operator, contact support |
| CURRENCY_MISMATCH | Request currency ≠ player currency | Convert or return error |

**Example Error Response:**
```json
{
  "status": "error",
  "code": "INSUFFICIENT_BALANCE",
  "message": "Player balance is insufficient for this bet",
  "balance": 25.50,
  "required": 50.00
}
```

**Recovery Strategy:**
```go
if response.Status == "error" {
    switch response.Code {
    case "INSUFFICIENT_BALANCE":
        // Log, reject transaction, return balance to player
    case "SESSION_EXPIRED":
        // Trigger re-authentication flow
    case "DUPLICATE_TRANSACTION":
        // Check transaction log, return original response
    default:
        // Log unknown error, escalate to support
    }
}
```

---

### 2. Evolution Gaming

**HTTP Status:** Implied by response structure

| Error Code | Meaning | Recovery |
|-----------|---------|----------|
| INTERNAL_ERROR | Server error | Retry (max 3 times, 100ms backoff) |
| INVALID_SIGNATURE | Signature validation failed | Check secret key, retry |
| INVALID_TOKEN | Token format incorrect | Verify token format, re-issue |
| INVALID_AMOUNT | Amount validation failed | Validate amount range before sending |
| INSUFFICIENT_FUNDS | Player lacks balance | Check balance endpoint first |
| OPERATOR_DISABLED | Operator blocked | Contact Evolution Gaming support |
| SESSION_TIMEOUT | Session expired | Initiate new session |

**Example Error Response:**
```json
{
  "status": 1,
  "code": 100,
  "description": "INTERNAL_ERROR",
  "digest": "signature_here"
}
```

**Recovery Strategy:**
```go
switch response.Code {
case 100: // INTERNAL_ERROR
    // Retry with exponential backoff
    backoff := 100 * math.Pow(2, retryCount) // 100ms, 200ms, 400ms...
    time.Sleep(time.Duration(backoff) * time.Millisecond)
case 102: // INVALID_SIGNATURE
    // Check gateway signature generation
    logSignatureDebug(request, response)
case 103: // INVALID_TOKEN
    // Validate token format, re-issue
default:
    // Log error, return to operator
}
```

---

### 3. PG Soft

**HTTP Status:** 200 (success) or 400+ (validation failure)

| Error Code | Meaning | Recovery |
|-----------|---------|----------|
| 0 | Success | Process normally |
| 1 | Invalid session | Create new session |
| 2 | Insufficient funds | Check balance first |
| 3 | Invalid currency | Use correct currency for player |
| 4 | Duplicate transaction | Return cached response |
| 5 | Invalid player | Verify player exists in database |
| 6 | Maintenance | Wait and retry (max 3 times) |

**Example Error Response:**
```json
{
  "status": 2,
  "message": "Insufficient funds for this transaction",
  "balance": 10.00,
  "currency": "BRL"
}
```

**Recovery Strategy:**
```go
if response.Status != 0 {
    switch response.Status {
    case 1:
        // Trigger session re-creation
    case 2:
        // Return insufficient funds error to operator
    case 4:
        // Check transaction log, return original response
    case 6:
        // Wait 30s, retry (max 3 times)
    default:
        // Log and escalate
    }
}
```

---

### 4. Mancala

**HTTP Status:** 200 (response) or 400+ (validation)

| Error Code | Meaning | Recovery |
|-----------|---------|----------|
| ERR_INVALID_SESSION | Session not found | Issue new session token |
| ERR_INSUFFICIENT_BALANCE | Insufficient funds | Check balance, reject bet |
| ERR_DUPLICATE_TRANSACTION | Already processed | Return cached response |
| ERR_INVALID_SIGNATURE | Signature mismatch | Verify secret key |
| ERR_GAME_NOT_FOUND | Invalid game ID | Validate against provider game list |
| ERR_PLAYER_BLOCKED | Player account disabled | Contact support |
| ERR_CURRENCY_INVALID | Currency not supported | Use supported currency |

**Example Error Response:**
```json
{
  "sessionId": "op1_session123",
  "success": false,
  "error": "ERR_INSUFFICIENT_BALANCE",
  "balance": 5.00,
  "required": 25.00
}
```

**Recovery Strategy:**
```go
if !response.Success {
    switch response.Error {
    case "ERR_INSUFFICIENT_BALANCE":
        // Reject operation, log balance mismatch
    case "ERR_DUPLICATE_TRANSACTION":
        // Query transaction log, return original
    case "ERR_INVALID_SIGNATURE":
        // Debug signature calculation, check secret
    default:
        // Log error, escalate
    }
}
```

---

### 5. Digitain

**HTTP Status:** 200 (response) with individual item status codes

| Error Code | Meaning | Recovery |
|-----------|---------|----------|
| 0 | Success | Process normally |
| 1 | Invalid operator | Verify operator_id |
| 2 | Invalid token | Re-issue session token |
| 3 | Insufficient balance | Check player balance |
| 4 | Duplicate transaction | Return cached response |
| 5 | Invalid amount | Validate amount before sending |
| 6 | Operation timeout | Retry with backoff |
| 100 | Invalid signature | Check secret key |

**Example Error Response (per item):**
```json
{
  "operatorId": "op123",
  "timestamp": "20260508120000",
  "items": [
    {
      "transactionId": "txn1",
      "status": 3,
      "statusText": "Insufficient balance"
    }
  ]
}
```

**Recovery Strategy (Batch Processing):**
```go
// Process per-item status codes
for _, item := range response.Items {
    switch item.Status {
    case 0:
        // Success, log transaction
    case 3:
        // Insufficient balance, log and skip
    case 4:
        // Duplicate, return cached response for this transaction
    case 6:
        // Timeout, add to retry queue
    default:
        // Log unknown status, add to failed items
    }
}
```

---

### 6. Evoplay

**HTTP Status:** 200 (response) with error in body

| Error Code | Meaning | Recovery |
|-----------|---------|----------|
| success | Operation succeeded | Process normally |
| invalid_signature | Signature validation failed | Verify secret key and field traversal |
| invalid_token | Token format incorrect | Re-issue token |
| insufficient_balance | Insufficient funds | Check balance first |
| duplicate_transaction | Already processed | Return cached response |
| invalid_amount | Amount validation failed | Validate amount rules |
| player_not_found | Player doesn't exist | Create player account |
| maintenance | Provider maintenance | Retry later (max 3 times, 30s intervals) |

**Example Error Response:**
```json
{
  "error": "insufficient_balance",
  "message": "Player does not have sufficient balance",
  "balance": 10.00,
  "transactionId": "txn_abc123"
}
```

**Recovery Strategy:**
```go
if response.Error != "success" {
    switch response.Error {
    case "invalid_signature":
        // Debug signature calculation
        // Check field traversal recursion
    case "insufficient_balance":
        // Reject operation, return balance
    case "maintenance":
        // Queue for retry (30s intervals, max 3)
    default:
        // Log and escalate
    }
}
```

---

### 7. OpenBox

**HTTP Status:** 200 (response) with status code

| Status Code | Meaning | Recovery |
|-----------|---------|----------|
| 0 | Success | Process normally |
| 1 | Invalid signature | Verify HMAC calculation and timestamp |
| 2 | Invalid timestamp | Check system time synchronization |
| 3 | Insufficient balance | Check balance before operation |
| 4 | Invalid player | Verify player exists |
| 5 | Duplicate operation | Return cached response |
| 6 | Session expired | Re-authenticate |
| 100 | Server error | Retry with exponential backoff |

**Example Error Response:**
```json
{
  "status": 3,
  "message": "Insufficient balance",
  "balance": 5.50,
  "roundId": "round_abc"
}
```

**Recovery Strategy:**
```go
switch response.Status {
case 0:
    // Success
case 1:
    // Invalid signature - debug HMAC calculation
    logSignatureDebug(request, expectedSignature)
case 2:
    // Timestamp issue - verify NTP sync
case 5:
    // Duplicate - check transaction log
case 100:
    // Server error - retry with backoff
default:
    // Log and escalate
}
```

---

### 8. Alternar

**HTTP Status:** Passed through from provider

| Status Code | Meaning | Recovery |
|-----------|---------|----------|
| 200 | Success | Process normally |
| 400 | Bad request | Check request format |
| 401 | Unauthorized | Verify provider credentials |
| 500 | Provider error | Retry with backoff |

**No special error handling** - gateway passes through provider response as-is.

---

## Common Error Scenarios & Recovery

### Scenario 1: Signature Validation Failure

**Symptoms:**
- Error code: `INVALID_SIGNATURE` or similar
- Recurs across multiple requests

**Root Causes:**
1. Secret key mismatch
2. Field ordering incorrect
3. Optional fields included/excluded incorrectly
4. Encoding mismatch (Base64 vs Base64URL)

**Debugging Steps:**
```go
func debugSignatureFailure(provider string, request map[string]interface{}, expectedSig string) {
    // Step 1: Verify secret key
    log.Printf("Provider: %s", provider)
    log.Printf("Secret Key Hash: %x", md5.Sum([]byte(secretKey)))
    
    // Step 2: Reconstruct payload exactly
    payload := reconstructPayload(provider, request)
    log.Printf("Payload: %s", payload)
    
    // Step 3: Compute signature step-by-step
    computed := computeSignature(provider, payload, secretKey)
    log.Printf("Expected: %s", expectedSig)
    log.Printf("Computed: %s", computed)
    
    // Step 4: Check for encoding issues
    if provider == "openbox" || provider == "digitain" {
        decodedExpected := decodeBase64URL(expectedSig)
        decodedComputed := decodeBase64URL(computed)
        log.Printf("Decoded match: %v", decodedExpected == decodedComputed)
    }
}
```

---

### Scenario 2: Duplicate Transaction Detection

**Symptoms:**
- Error: `DUPLICATE_TRANSACTION` or `409 Conflict`
- Same transaction ID appears twice

**Recovery:**
```go
func handleDuplicateTransaction(txnId string) (response interface{}, error error) {
    // 1. Check transaction log
    existing := db.Where("transaction_id = ?", txnId).First(&TransactionLog{})
    if existing != nil {
        // 2. Return original response
        log.Printf("Duplicate detected: %s, returning cached response", txnId)
        return existing.Response, nil
    }
    
    // 3. If not in log but provider says duplicate, something is wrong
    return nil, fmt.Errorf("provider reported duplicate but not in system: %s", txnId)
}
```

---

### Scenario 3: Balance Mismatch

**Symptoms:**
- Error: `INSUFFICIENT_BALANCE` when balance appears sufficient
- Recurs for same player

**Root Causes:**
1. Pending transactions not deducted
2. Currency mismatch (e.g., cents vs. dollars)
3. Decimal precision loss

**Verification:**
```go
func verifyBalance(playerId string, playerBalance decimal.Decimal) error {
    // 1. Check local balance
    localBalance := getLocalBalance(playerId)
    
    // 2. Check pending transactions
    pending := db.Where("player_id = ? AND status = ?", playerId, "pending").
        Sum("amount")
    
    available := localBalance.Sub(pending)
    
    // 3. Verify precision (use decimal.Decimal, not float64)
    if available.Compare(playerBalance) < 0 {
        log.Warnf("Balance mismatch for player %s: local=%v, available=%v, provider=%v",
            playerId, localBalance, available, playerBalance)
        return fmt.Errorf("balance mismatch detected")
    }
    
    return nil
}
```

---

### Scenario 4: Timeout & Retry

**Symptoms:**
- Error: `504 Gateway Timeout`
- Recurs intermittently

**Recovery:**
```go
func requestWithRetry(provider string, request interface{}, maxRetries int) (response interface{}, error error) {
    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        // Exponential backoff: 100ms, 200ms, 400ms, ...
        if attempt > 0 {
            backoff := time.Duration(100 * math.Pow(2, float64(attempt-1))) * time.Millisecond
            time.Sleep(backoff)
        }
        
        // Make request with timeout
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        response, err := sendRequest(ctx, provider, request)
        
        // Check if retryable
        if err == nil {
            return response, nil
        }
        
        if !isRetryable(err) {
            return nil, err
        }
        
        lastErr = err
        log.Warnf("Attempt %d failed for %s: %v, retrying...", attempt+1, provider, err)
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetryable(err error) bool {
    // Retryable: timeout, connection error, 5xx
    // Not retryable: 4xx, signature failure, validation
    code := getHTTPStatusCode(err)
    return code >= 500 || code == 0 // 0 = connection error
}
```

---

### Scenario 5: Operator Context Extraction Failure

**Symptoms:**
- Error: Invalid operator context
- Token format incorrect

**Root Causes:**
1. Token missing `_` separator
2. Operator slug not in system
3. Operator disabled

**Verification:**
```go
func validateOperatorContext(token string) (*Operator, error) {
    // 1. Check format
    parts := strings.Split(token, "_")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid token format: expected {slug}_{value}, got %s", token)
    }
    
    slug := parts[0]
    
    // 2. Check operator exists
    operator := &Operator{}
    err := db.Where("slug = ?", slug).First(operator).Error
    if err != nil {
        return nil, fmt.Errorf("operator not found: %s", slug)
    }
    
    // 3. Check operator enabled
    if !operator.IsEnabled {
        return nil, fmt.Errorf("operator disabled: %s", slug)
    }
    
    return operator, nil
}
```

---

## Monitoring & Alerting

### Key Metrics to Track

| Metric | Alert Threshold | Action |
|--------|-----------------|--------|
| Error rate by provider | > 5% | Page oncall, check provider status |
| Signature validation failures | > 1 per minute | Debug signature calculation |
| Timeout rate | > 2% | Check provider performance |
| Duplicate transaction rate | > 0.5% | Investigate idempotency tracking |
| Balance mismatch rate | > 0.1% | Audit transaction log |

### Logging Template

```go
func logError(provider string, endpoint string, err error, req, res interface{}) {
    errorLog := struct {
        Timestamp    time.Time
        Provider     string
        Endpoint     string
        Error        string
        ErrorType    string
        RequestID    string
        OperatorID   int
        Status       int
        RetryCount   int
        Request      interface{}
        Response     interface{}
    }{
        Timestamp:  time.Now(),
        Provider:   provider,
        Endpoint:   endpoint,
        Error:      err.Error(),
        ErrorType:  getErrorType(err),
        RequestID:  uuid.New().String(),
        OperatorID: extractOperatorID(req),
        Status:     getStatusCode(res),
        RetryCount: getRetryCount(req),
        Request:    sanitizeRequest(req),  // Don't log credentials
        Response:   sanitizeResponse(res),
    }
    
    db.Create(&errorLog)
    
    // Alert if critical
    if isCritical(errorLog) {
        alertSlack(errorLog)
    }
}
```

---

## Best Practices

### 1. Always Validate Before Processing
```go
// ✓ Correct order
1. Validate request structure
2. Validate signature
3. Extract operator context
4. Check operator permissions
5. Process payload
```

### 2. Use Idempotency Keys
```go
// Every request must include unique ID
type Request struct {
    TransactionID string    // Unique per transaction
    RequestID     string    // Unique per request attempt
    // ...
}

// Log both to enable retry deduplication
```

### 3. Implement Circuit Breaker
```go
func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.IsOpen() {
        return fmt.Errorf("circuit breaker open")
    }
    
    err := fn()
    
    if err != nil {
        cb.RecordFailure()
        if cb.FailureCount() > cb.Threshold {
            cb.Open()
        }
    } else {
        cb.RecordSuccess()
    }
    
    return err
}
```

### 4. Document Provider-Specific Quirks
```go
// Example: PG Soft requires exact field order
// Example: OpenBox uses Base64URL (not standard Base64)
// Example: Digitain batch responses include per-item status
type ProviderQuirks struct {
    Name             string
    RequireFieldOrder bool
    Base64URLEncoding bool
    BatchSupport     bool
    SignatureLocation string
}
```

---

## Recovery Decision Tree

```
Error Received
├─ HTTP 4xx?
│  ├─ 400: Bad request → Validate payload structure
│  ├─ 401: Auth failed → Check signature, credentials
│  ├─ 403: Forbidden → Check operator enabled, permissions
│  └─ 409: Conflict → Check transaction log, return cached
│
├─ HTTP 5xx?
│  ├─ 500: Server error → Retry (exponential backoff)
│  ├─ 502: Bad gateway → Check provider status, retry
│  ├─ 503: Unavailable → Retry with longer backoff (max 30s)
│  └─ 504: Timeout → Retry (max 3 times, 30s timeout)
│
└─ Provider-specific error?
   ├─ DUPLICATE_TRANSACTION → Return cached response
   ├─ INSUFFICIENT_BALANCE → Reject operation
   ├─ SESSION_EXPIRED → Re-authenticate
   ├─ SIGNATURE_FAILED → Debug signature calculation
   └─ Other → Log, escalate to support
```

---

## Testing Error Scenarios

### Unit Test Template

```go
func TestErrorScenario_SignatureFailure(t *testing.T) {
    // Setup
    request := createTestRequest()
    secretKey := "test_secret"
    
    // Test with wrong secret
    wrongSecret := "wrong_secret"
    valid := validateSignature(request, "expected_sig", wrongSecret)
    
    assert.False(t, valid, "should reject with wrong secret")
}

func TestErrorScenario_DuplicateTransaction(t *testing.T) {
    // Setup
    txnId := "txn_abc123"
    
    // Insert first transaction
    firstResponse := processTransaction(txnId)
    assert.NotNil(t, firstResponse)
    
    // Try same transaction again
    secondResponse, err := processTransaction(txnId)
    assert.NoError(t, err)
    assert.Equal(t, firstResponse, secondResponse, "should return cached response")
}

func TestErrorScenario_BalanceMismatch(t *testing.T) {
    // Setup
    playerId := "player123"
    localBalance := decimal.NewFromString("100.00")
    
    // Provider reports higher balance
    providerBalance := decimal.NewFromString("150.00")
    
    // Should detect mismatch
    mismatch := checkBalance(playerId, localBalance, providerBalance)
    assert.True(t, mismatch, "should detect balance mismatch")
}
```

---

## Summary Table

| Provider | Common Errors | Recovery |
|----------|---------------|----------|
| Pragmatic Play | INSUFFICIENT_BALANCE, SESSION_EXPIRED | Check balance, re-auth |
| Evolution Gaming | INTERNAL_ERROR, INVALID_SIGNATURE | Retry, debug signature |
| PG Soft | Invalid session, Insufficient funds | Session re-creation |
| Mancala | Signature mismatch, Balance issues | Debug signature, verify balance |
| Digitain | Per-item status codes, batch errors | Process per item, batch retry |
| Evoplay | Signature/field traversal issues | Debug field recursion |
| OpenBox | Timestamp sync, signature validation | Check NTP, debug HMAC |
| Alternar | Provider-specific (transparent proxy) | Return provider response |
