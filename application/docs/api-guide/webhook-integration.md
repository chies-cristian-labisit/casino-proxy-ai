# Casino Proxy API - Webhook Integration Guide

## Quick Start

### 1. Register as Operator

Create an operator account via the admin API:

```bash
curl -X POST https://api.casino-proxy.local/v1/internal/operator/store \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "my-operator",
    "name": "My Casino",
    "url": "https://my-casino.example.com/webhooks"
  }'
```

Response:
```json
{
  "id": 1,
  "slug": "my-operator",
  "name": "My Casino",
  "url": "https://my-casino.example.com/webhooks",
  "created_at": "2026-05-08T12:00:00Z"
}
```

### 2. Configure Provider Credentials

Store credentials for each gaming provider:

```json
{
  "operator_id": 1,
  "provider": "pragmatic_play",
  "api_key": "your_api_key",
  "secret_key": "your_secret_key"
}
```

### 3. Receive Webhooks

Gaming providers send webhooks to the gateway. The gateway validates, transforms, and forwards them to your operator URL.

**Request Flow:**
```
Gaming Provider 
  ↓
Casino Proxy Gateway (https://casino-proxy.local/v1/webhooks/{provider}/{endpoint})
  ├─ Validate signature
  ├─ Extract operator context
  ├─ Sanitize tokens
  └─ Forward to operator URL (https://my-casino.example.com/webhooks)
```

---

## Integration Patterns by Provider

### Pattern 1: Generic Routing (Pragmatic Play, Mancala, Digitain, Evoplay)

**Characteristics:**
- Single controller handles multiple endpoints
- Endpoint specified as parameter or field
- Operator context extracted from payload

**Gateway Path:**
```
POST /v1/webhooks/pragmatic-play/{endpoint}
POST /v1/webhooks/mancala/{endpoint}
POST /v1/webhooks/digitain-rgs/{endpoint}
POST /v1/webhooks/evoplay  (endpoint via 'name' field)
```

**Your Callback URL Structure:**
```
https://my-casino.example.com/webhooks/pragmatic-play/{endpoint}
https://my-casino.example.com/webhooks/mancala/{endpoint}
https://my-casino.example.com/webhooks/digitain-rgs/{endpoint}
https://my-casino.example.com/webhooks/evoplay
```

### Pattern 2: Specific Endpoints (Evolution Gaming, PG Soft)

**Characteristics:**
- Each operation has distinct endpoint
- No dynamic routing needed
- Clear endpoint paths

**Gateway Paths:**
```
POST /v1/webhooks/evolution/authentication
POST /v1/webhooks/evolution/debit
POST /v1/webhooks/evolution/credit
POST /v1/webhooks/evolution/rollback
POST /v1/webhooks/evolution/getNewToken

POST /v1/webhooks/pgsoft/VerifySession
POST /v1/webhooks/pgsoft/Cash/Get
POST /v1/webhooks/pgsoft/Cash/TransferInOut
POST /v1/webhooks/pgsoft/Cash/Adjustment
```

**Your Callback URLs:**
```
https://my-casino.example.com/webhooks/evolution/authentication
https://my-casino.example.com/webhooks/evolution/debit
... (one per endpoint)

https://my-casino.example.com/webhooks/pgsoft/VerifySession
... (one per endpoint)
```

### Pattern 3: PUT-Based Operations (OpenBox)

**Characteristics:**
- Uses PUT method (idempotent semantics)
- Signature in headers (not body)
- 4 core operations + round status

**Gateway Paths:**
```
PUT /v1/webhooks/openbox/seamless/balance
PUT /v1/webhooks/openbox/seamless/bet
PUT /v1/webhooks/openbox/seamless/win
PUT /v1/webhooks/openbox/seamless/refund
PUT /v1/webhooks/openbox/seamless/round_status
```

**Headers:**
```
Signature: {apiKey}:{base64url_hmac}
Timestamp: {unix_timestamp}
Content-Type: application/json; charset=utf-8
```

### Pattern 4: HTTP Proxy (Alternar)

**Characteristics:**
- Gateway acts as transparent HTTP proxy
- No signature validation at gateway
- Request forwarded as-is

**Gateway Path:**
```
POST /v1/webhooks/alternar
```

---

## Request Transformation

The gateway performs these transformations:

### 1. Token Sanitization

**Input from Provider:** `my-operator_abc123xyz`  
**Processing:**
1. Extract operator slug: `my-operator`
2. Verify operator exists in system
3. Forward with sanitized token to operator URL

**Your Callback Receives:** `abc123xyz` (clean token)

```go
// In your application
func handlePragmaticPlayWebhook(req *http.Request) {
    payload := parseJSON(req.Body)
    
    // Token is already sanitized (operator prefix removed)
    sessionId := payload["SessionId"]  // e.g., "abc123xyz"
    
    // You can look up the session without operator prefix
    session := db.Where("token = ?", sessionId).First()
}
```

### 2. Signature Regeneration (for some providers)

For providers that validate signatures, the gateway:
1. Validates provider's signature
2. Regenerates signature with your operator credentials
3. Includes new signature in request to your URL

**Example (OpenBox):**
```
Gateway receives:
  Signature: provider_api_key:signature123
  
Gateway validates it, then regenerates:
  Signature: your_api_key:signature456
  
Your system validates with YOUR credentials
```

### 3. Response Passthrough

The gateway forwards your response back to the provider with:
- Your status code
- Your response body
- No modification

---

## Webhook Examples

### Example 1: Pragmatic Play Bet Placement

**Provider → Gateway:**
```json
{
  "sessionId": "my-operator_session123",
  "userId": "my-operator_user456",
  "gameId": 42,
  "betAmount": 50.00,
  "hash": "abc123def456"
}
```

**Gateway → Your Callback:**
```json
{
  "sessionId": "session123",
  "userId": "user456",
  "gameId": 42,
  "betAmount": 50.00,
  "hash": "xyz789abc123"
}
```

**Your Response → Gateway → Provider:**
```json
{
  "balance": 950.00,
  "success": true
}
```

### Example 2: Evolution Gaming Debit (Bet)

**Provider → Gateway:**
```json
{
  "token": "my-operator_token123",
  "roundId": "round_abc",
  "amount": 50.00,
  "gameId": 7,
  "digest": "sha256_signature"
}
```

**Gateway → Your Callback:**
```json
{
  "token": "my-operator_token123",
  "roundId": "round_abc",
  "amount": 50.00,
  "gameId": 7,
  "digest": "your_sha256_signature"
}
```

**Your Response:**
```json
{
  "balance": 950.00,
  "currency": "BRL"
}
```

### Example 3: OpenBox Balance Query (PUT)

**Provider → Gateway:**
```
PUT /v1/webhooks/openbox/seamless/balance
Headers:
  Signature: provider_api_key:signature123
  Timestamp: 1715177400

Body:
{
  "player": "my-operator_player123",
  "roundId": "round_abc"
}
```

**Gateway → Your Callback:**
```
PUT https://my-casino.example.com/webhooks/openbox/seamless/balance
Headers:
  Signature: your_api_key:signature456
  Timestamp: 1715177400

Body:
{
  "player": "my-operator_player123",
  "roundId": "round_abc"
}
```

---

## Batch Processing (Digitain)

Digitain supports multi-operator batch processing:

**Provider → Gateway:**
```json
{
  "timestamp": "20260508120000",
  "operatorId": "op123",
  "items": [
    {"token": "op1_token1", "amount": 50.00},
    {"token": "op2_token2", "amount": 75.00}
  ]
}
```

**Gateway Processing:**
1. Groups items by operator
2. Sends separate requests per operator:
   - Request 1 to op1 webhook with op1's items
   - Request 2 to op2 webhook with op2's items
3. Aggregates responses into single batch response

**Your Callback (for op1) receives:**
```json
{
  "timestamp": "20260508120000",
  "operatorId": "op123",
  "items": [
    {"token": "op1_token1", "amount": 50.00}
  ]
}
```

---

## Error Handling in Webhooks

### Your Response Status Codes

- **200 OK:** Operation processed successfully
- **400 Bad Request:** Invalid request format or validation failure
- **409 Conflict:** Transaction already processed (idempotency)
- **500 Server Error:** Internal server error

**Example Error Response:**
```json
{
  "error": "INSUFFICIENT_BALANCE",
  "description": "Player does not have sufficient balance for this transaction",
  "balance": 25.00
}
```

### Idempotency

All requests include transaction IDs. Your system should:
1. Store transaction ID with processing result
2. Check if transaction already processed
3. Return same response for duplicate requests

```go
// Pseudocode
func handleBet(req *http.Request) {
    txnId := req.JSON["transactionId"]
    
    // Check if already processed
    existing := db.Where("transaction_id = ?", txnId).First()
    if existing != nil {
        return existing.response  // Return original response
    }
    
    // Process new transaction
    result := processBet(req.JSON)
    
    // Store for future deduplication
    db.Create(&TransactionLog{
        TransactionID: txnId,
        Response: result,
    })
    
    return result
}
```

### Retry Logic

The gateway retries failed requests with exponential backoff:
- Retry count: 3
- Initial backoff: 100ms
- Max delay: 30s

Your endpoint should handle:
- Idempotent processing (same request = same result)
- Reasonable timeout (30s)
- Proper error responses

---

## Testing Webhooks

### Local Development Setup

1. **Install ngrok or similar tunnel:**
```bash
ngrok http 3000  # Exposes localhost:3000 to internet
```

2. **Update operator URL:**
```bash
curl -X POST http://localhost:3000/v1/internal/operator/store \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "test-op",
    "name": "Test Operator",
    "url": "https://abc123.ngrok.io/webhooks"
  }'
```

3. **Create test endpoints in your application:**
```go
func handleTestWebhook(w http.ResponseWriter, r *http.Request) {
    // Log incoming webhook
    log.Printf("Received: %s %s", r.Method, r.URL)
    
    // Parse and validate
    payload := parsePayload(r)
    log.Printf("Payload: %+v", payload)
    
    // Return test response
    json.NewEncoder(w).Encode(map[string]interface{}{
        "balance": 1000.00,
        "status": "ok",
    })
}
```

### Using curl for Testing

```bash
# Test Pragmatic Play webhook
curl -X POST http://localhost:3000/v1/webhooks/pragmatic-play/balance \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "test-op_session123",
    "hash": "test_hash"
  }'

# Test PG Soft webhook
curl -X POST http://localhost:3000/v1/webhooks/pgsoft/Cash/Get \
  -H "Content-Type: application/json" \
  -d '{
    "operator_token": "test_token",
    "secret_key": "test_secret",
    "player_name": "test-op_player1"
  }'

# Test OpenBox webhook (PUT)
curl -X PUT http://localhost:3000/v1/webhooks/openbox/seamless/balance \
  -H "Content-Type: application/json" \
  -H "Signature: api_key:signature" \
  -H "Timestamp: 1715177400" \
  -d '{
    "player": "test-op_player1",
    "roundId": "round123"
  }'
```

### Webhook Payload Validation

Always validate webhook payloads:

```go
func validateWebhookPayload(provider string, payload map[string]interface{}) error {
    // Check required fields
    requiredFields := getRequiredFields(provider)
    for _, field := range requiredFields {
        if _, ok := payload[field]; !ok {
            return fmt.Errorf("missing required field: %s", field)
        }
    }
    
    // Validate data types
    if sessionId, ok := payload["sessionId"].(string); !ok {
        return errors.New("sessionId must be string")
    }
    
    // Validate operator context
    operator, err := extractOperator(payload)
    if err != nil {
        return err
    }
    
    // Validate signature if present
    if sig, ok := payload["hash"].(string); ok {
        if !validateSignature(payload, sig) {
            return errors.New("invalid signature")
        }
    }
    
    return nil
}
```

---

## Rate Limiting

No gateway-level rate limiting is enforced. Implement in your application:

```go
// Redis-based rate limiting
func RateLimitMiddleware(maxPerSecond int) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            operator := extractOperator(r)
            key := fmt.Sprintf("rate:%s:%d", operator, time.Now().Unix())
            
            count, _ := redis.Incr(key)
            if count > maxPerSecond {
                w.WriteHeader(http.StatusTooManyRequests)
                return
            }
            
            redis.Expire(key, 1*time.Second)
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Monitoring & Logging

Log all webhook activity:

```go
type WebhookLog struct {
    ID            int
    Provider      string
    Endpoint      string
    OperatorID    int
    TransactionID string
    Status        int
    RequestBody   string // JSON
    ResponseBody  string // JSON
    ProcessingTime int    // ms
    CreatedAt     time.Time
}

func logWebhook(provider string, endpoint string, req, res []byte, duration time.Duration) {
    log := WebhookLog{
        Provider:     provider,
        Endpoint:     endpoint,
        RequestBody:  string(req),
        ResponseBody: string(res),
        ProcessingTime: int(duration.Milliseconds()),
        CreatedAt:    time.Now(),
    }
    db.Create(&log)
}
```

Monitor key metrics:
- Webhook latency (p50, p95, p99)
- Error rate by provider
- Transaction success rate
- Duplicate request ratio
