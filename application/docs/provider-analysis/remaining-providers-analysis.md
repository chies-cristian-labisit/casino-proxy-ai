# Remaining Providers Integration Analysis (Mancala, Digitain, Evoplay, OpenBox, Alternar)

## Overview

**5 Remaining Providers Analyzed:**
1. **Mancala** - Generic routing, MD5 signature, 4 operations
2. **Digitain RGS** - Generic routing + batch processing, HMAC-SHA256, 11 endpoints
3. **Evoplay** - Single webhook endpoint, MD5 signature
4. **OpenBox** - PUT-based seamless operations, HMAC-SHA256, 4 operations
5. **Alternar** - Redirect handler, HTTP forwarding

---

## Provider Integration Patterns

### Pattern 1: Generic Routing Providers (Mancala, Digitain, Evoplay)

**Characteristic:** Dynamic endpoint routing via request parameters

```
REQUEST → Gateway → Dynamic Handler Selection → Operator URL
                        ↓
                   Validate Endpoint
                   Extract Operator
                   Sanitize Payload
                   Forward to Operator
```

**Key Features:**
- Single controller handles multiple logical endpoints
- Endpoint specified as path parameter or request field
- Operator context extracted from payload (token, ExtraData, playerId)
- Token/session sanitization: `{operator_slug}_{value}` format
- Signature validation varies by provider (MD5 or HMAC-SHA256)

---

## 1. Mancala Provider

### Integration Characteristics

**Type:** Generic routing webhook provider  
**Authentication:** MD5 signature in request payload  
**Allowed Endpoints:** 4 (Balance, Credit, Debit, Refund)  
**Operator Binding:** Via `ExtraData` JSON field  
**Signature Algorithm:** MD5 hash  

### Endpoint Structure

#### Request Format
```json
{
  "ExtraData": "{\"operator_slug\": \"my-operator\"}",
  "SessionId": "my-operator_abc123xyz",
  "TransactionGuid": "guid_value",
  "Amount": 50.00,
  "RoundGuid": "round_guid",
  "RefundTransactionGuid": "refund_guid"
}
```

**Allowed Endpoints:**
- `Balance` - Query player balance
- `Credit` - Add funds to player account
- `Debit` - Withdraw funds from player account
- `Refund` - Refund transaction

#### Response Format
```json
{
  "data": {},
  "error": null
}
```

### Mancala-Specific Behaviors

1. **ExtraData Parsing:**
   - Operator context passed as JSON string in `ExtraData` field
   - Gateway parses and extracts operator slug
   - Allows operator isolation

2. **Hash Calculation:**
   ```
   MD5(endpoint + "/" + SessionId + [TransactionGuid] + [RefundTransactionGuid] + [RoundGuid] + [Amount] + secret_token)
   ```
   - Optional fields included only if present
   - Order matters: endpoint → session → transaction → amount → secret

3. **Token Sanitization:**
   - `SessionId` stripped of `{operator_slug}_` prefix before forwarding
   - Ensures clean token format for operator callback

### Implementation Notes

- **Credential Format:** `secret_token` stored in operator credentials
- **Operator URL:** Base URL points to operator callback endpoint
- **Path Construction:** `{operator.url}/mancala/{endpoint}`
- **Forwarding:** Transparent proxy after hash validation

---

## 2. Digitain RGS Provider

### Integration Characteristics

**Type:** Generic routing RGS (Random Generator Server) integration  
**Authentication:** HMAC-SHA256 signature  
**Allowed Endpoints:** 11 (authenticate, getbalance, bet, win, betwin, refund, amend, checktxstatus, charge, promowin, refreshtoken)  
**Operator Binding:** Via `token` or `playerId` field  
**Batch Processing:** Supports multi-operator batches via `items` array  

### Endpoint Structure

#### Single Request
```json
{
  "timestamp": "20260508120000",
  "operatorId": "operator_id_123",
  "token": "my-operator_abc123xyz",
  "playerId": "my-operator_player_001",
  "info": "bet_info",
  "amount": 100.00,
  "currency": "BRL"
}
```

#### Batch Request (Multiple Operators)
```json
{
  "timestamp": "20260508120000",
  "operatorId": "operator_id_123",
  "token": "my-operator_abc123xyz",
  "items": [
    {"playerId": "player_1", "amount": 50.00},
    {"playerId": "player_2", "amount": 75.00}
  ]
}
```

#### Response Format
```json
{
  "timestamp": "20260508120000",
  "signature": "hmac_sha256_signature",
  "errorCode": 1,
  "items": [
    {
      "info": "bet_info",
      "errorCode": 1,
      "metadata": "",
      "balance": 1000.00
    }
  ]
}
```

### Digitain-Specific Behaviors

1. **Multi-Operator Batch Processing:**
   - Single request can contain items from multiple operators
   - Gateway groups items by operator
   - Sends separate requests per operator to their callbacks
   - Aggregates responses back into single batch response
   - Powerful for reducing webhook spam from Digitain

2. **Token/PlayerId Handling:**
   - Operator context extracted from `token` or `playerId` field
   - Response includes operator slug prefix for reconstruction
   - Sanitization: Strip `{operator_slug}_` on request, add on response

3. **Hash Validation:**
   ```
   HMAC-SHA256(timestamp + operatorId, secretKey)
   ```
   - Timestamp in format: YmdHiu (20260508120000)
   - Result Base64-encoded in signature field

4. **Special Endpoint Handling:**
   - `checktxstatus`: Returns mock response (bypasses operator)
   - Used for transaction status checks without provider call

### Implementation Notes

- **Credentials:** operatorId, secret-key (HMAC key)
- **Operator URL:** Base URL points to operator callback endpoint
- **Path Construction:** `{operator.url}/digitain-rgs/{endpoint}`
- **Batch Intelligence:** Reduces webhook volume when operators send multi-item batches

---

## 3. Evoplay Provider

### Integration Characteristics

**Type:** Single webhook endpoint (generic event handling)  
**Authentication:** MD5 signature  
**Allowed Endpoints:** 5 (init, bet, win, refund, BalanceIncrease)  
**Operator Binding:** Via `token` or `data.user_id` field  
**Signature Algorithm:** Complex MD5 with recursive field traversal  

### Endpoint Structure

#### Request Format
```json
{
  "token": "my-operator_abc123xyz",
  "name": "bet",
  "data": {
    "user_id": "my-operator_player_001",
    "amount": 50.00,
    "currency": "BRL",
    "info": "custom_data"
  },
  "signature": "md5_signature"
}
```

**Allowed Operations (via 'name' field):**
- `init` - Initialize game session
- `bet` - Bet placement
- `win` - Win/payout
- `refund` - Transaction refund
- `BalanceIncrease` - Balance increase (promotional credit)

#### Response Format
```json
{
  "balance": 1000.00,
  "currency": "BRL",
  "error": null
}
```

### Evoplay-Specific Behaviors

1. **Endpoint Selection via 'name' Field:**
   - No URL path parameters for endpoint selection
   - Single webhook endpoint (`/v1/webhooks/evoplay`)
   - Operation type specified in JSON payload `name` field
   - Simplifies provider integration

2. **Complex MD5 Signature Calculation:**
   ```
   MD5(projectId * apiVersion * field1 * field2_recursive * ... * secretKey)
   ```
   - Fields joined with `*` separator
   - For nested arrays: recursively extract all values, join with `:`
   - Signature field itself excluded from calculation
   - Result is final signature

3. **Token Sanitization:**
   - `token` field: Strip `{operator_slug}_` prefix
   - `data.user_id` field: Strip `{operator_slug}_` prefix

### Implementation Notes

- **Credentials:** project_id, secret_key, api_version (required)
- **Operator URL:** Base URL points to operator callback endpoint
- **Path Construction:** `{operator.url}/evoplay` (single endpoint)
- **Signature Type:** MD5 (unique algorithm, recursive field handling)
- **No Signature Validation:** Gateway doesn't validate incoming signatures, only generates outgoing

---

## 4. OpenBox Provider (PUT-based)

### Integration Characteristics

**Type:** PUT-based seamless operations (unusual pattern)  
**Authentication:** HMAC-SHA256 signature in headers  
**Allowed Operations:** 4 (balance, bet, win, refund) + roundStatus  
**HTTP Method:** PUT (not POST) - indicates idempotency semantics  
**Signature Format:** Base64URL-encoded HMAC with API key prefix  

### Endpoint Structure

#### Request Format (All Operations)
```json
{
  "player": "my-operator_abc123xyz",
  "roundId": "round_123",
  "amount": 50.00,
  "currency": "BRL",
  "transactionId": "txn_001"
}
```

**Operations:**
- `PUT /v1/webhooks/openbox/seamless/balance` - Get balance
- `PUT /v1/webhooks/openbox/seamless/bet` - Place bet
- `PUT /v1/webhooks/openbox/seamless/win` - Process win
- `PUT /v1/webhooks/openbox/seamless/refund` - Refund operation
- `PUT /v1/webhooks/openbox/seamless/round_status` - Round status query

**Headers:**
```
Signature: {apiKey}:{base64url(hmac)}
Timestamp: {unix_timestamp}
Content-Type: application/json; charset=utf-8
```

#### Response Format
```json
{
  "balance": 1000.00,
  "currency": "BRL",
  "status": "success"
}
```

### OpenBox-Specific Behaviors

1. **Non-Standard PUT Method:**
   - Unusual use of PUT for state-changing operations (normally POST)
   - Suggests idempotent semantics: PUT same request = safe to retry
   - Simplifies deduplication logic (PUT naturally idempotent)

2. **Signature Generation:**
   ```
   sorted_data = ksort(data)
   query_string = key1=value1&key2=value2&...
   signature_base = query_string + ":" + secretKey + timestamp
   hmac_result = HMAC-SHA256(signature_base, secretKey)
   final_signature = apiKey + ":" + base64url(hmac_result)
   ```
   - Data sorted by key (ksort) before signature
   - Boolean values converted to "true"/"false" strings
   - Signature includes apiKey prefix

3. **Token Sanitization:**
   - `player` field: Strip `{operator_slug}_` prefix
   - After sanitization, regenerate signature for operator

4. **Header-Based Signature:**
   - Signature passed in `Signature` header (not request body)
   - Timestamp passed in `Timestamp` header
   - Enables HTTP-level signature verification

### Implementation Notes

- **Credentials:** api_key, secret_key (both required)
- **Operator URL:** Base URL points to operator callback endpoint
- **Path Construction:** `{operator.url}/openbox/seamless/{operation}`
- **HTTP Method:** PUT (idempotent)
- **Signature Type:** HMAC-SHA256 with Base64URL encoding
- **Header Validation:** Signature computed fresh for each operator request

---

## 5. Alternar Provider (Redirect Handler)

### Integration Characteristics

**Type:** Redirect handler (forwards to external system)  
**HTTP Method:** POST  
**Authentication:** Not cryptographic (HTTP forwarding only)  
**Pattern:** Proxy to external webhook endpoint  

### Endpoint Structure

#### Request Format
```json
{
  "player": "player_id",
  "roundId": "round_123",
  "amount": 50.00,
  "operatorId": "operator_id",
  "sessionToken": "session_token"
}
```

**Single Endpoint:**
- `POST /v1/webhooks/alternar/` - Generic redirect

#### Response Format
```json
{
  "status": 200,
  "data": {}
}
```

### Alternar-Specific Behaviors

1. **HTTP Forwarding Pattern:**
   - Request received at gateway
   - Request forwarded to external URL: `https://hypezbet-backoffice.cometagaming.com/api/v1/webhooks/alternar/`
   - Response status code passed through from upstream
   - Acts as transparent HTTP proxy

2. **Retry Logic:**
   - Automatic retry: 3 attempts, 100ms backoff
   - Ensures reliability for network failures
   - Complies with HTTP retry semantics

3. **No Token Sanitization:**
   - Request forwarded as-is to external system
   - No operator context extraction/management
   - Suggests Alternar integration handles its own operator routing

4. **Status Code Passthrough:**
   - Uses response status code from upstream (default 400 if missing)
   - Allows upstream system to control HTTP semantics

### Implementation Notes

- **Credentials:** None documented (may use external authentication)
- **Target URL:** Hardcoded to external webhook endpoint
- **HTTP Method:** POST
- **Pattern Type:** HTTP proxy/redirect (unusual for webhook handling)
- **Error Handling:** Status code passthrough

---

## Provider Comparison Matrix

| Feature | Mancala | Digitain | Evoplay | OpenBox | Alternar |
|---------|---------|----------|---------|---------|----------|
| **Auth Type** | MD5 sig | HMAC-SHA256 | MD5 sig | HMAC-SHA256 | HTTP proxy |
| **Routing** | Dynamic {endpoint} | Dynamic {endpoint} | Single endpoint | 4 PUT methods | Single endpoint |
| **Batch Processing** | No | Yes (items array) | No | No | No |
| **HTTP Method** | POST | POST | POST | PUT | POST |
| **Operator Binding** | ExtraData JSON | token/playerId | token/user_id | player | (external) |
| **Endpoints** | 4 | 11 | 5 | 4+1 | 1 |
| **Signature Location** | Body field | Response field | Body field | Headers | N/A |
| **Idempotency** | Via TxID | Via items | Implicit | PUT semantics | Implicit |
| **Complexity** | Low | Medium-High | Medium | Low-Medium | Very Low |

---

## Common Patterns Across Remaining Providers

### 1. Generic Routing Pattern (Mancala, Digitain, Evoplay)

**Benefit:** Flexible endpoint handling, single controller handles multiple operations

**Pattern:**
```
Single Controller __invoke() method
  ↓
Extract operation/endpoint from request
  ↓
Validate against allowedEndpoints array
  ↓
Extract operator from payload
  ↓
Sanitize tokens: {operator_slug}_{value} format
  ↓
Forward to operator URL
```

**Go Implementation:** Use a service registry pattern with endpoint dispatch

```go
type GenericRouter struct {
    allowedEndpoints map[string]func(context.Context, []byte) (interface{}, error)
}

func (r *GenericRouter) Route(ctx context.Context, endpoint string, payload []byte) (interface{}, error) {
    handler, ok := r.allowedEndpoints[endpoint]
    if !ok {
        return nil, fmt.Errorf("endpoint not found: %s", endpoint)
    }
    return handler(ctx, payload)
}
```

### 2. Token Sanitization Pattern

**Applied By:** All 5 providers

**Logic:**
1. Extract operator slug from token prefix
2. Validate operator exists
3. Strip prefix for upstream forwarding
4. Add prefix back in response if needed

**Go Implementation:**
```go
func SanitizeToken(token string) (operatorSlug, actualToken string, err error) {
    parts := strings.Split(token, "_")
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid token format")
    }
    return parts[0], parts[1], nil
}
```

### 3. Signature Validation Pattern

**Applied By:** Mancala (MD5), Digitain (HMAC-SHA256), Evoplay (MD5), OpenBox (HMAC-SHA256)

**Flow:**
1. Extract signature from request (body field or header)
2. Compute expected signature from payload + secret
3. Compare (constant-time comparison for security)
4. Reject if mismatch

**Go Implementation:**
```go
func ValidateSignature(payload []byte, providedSig string, secretKey string, algorithm string) bool {
    expected := GenerateSignature(payload, secretKey, algorithm)
    // Use subtle.ConstantTimeCompare to prevent timing attacks
    return subtle.ConstantTimeCompare([]byte(expected), []byte(providedSig)) == 1
}
```

### 4. Operator Context Extraction

**Extraction Points:**
- **Mancala:** `ExtraData` JSON field
- **Digitain:** `token` or `playerId` field
- **Evoplay:** `token` or `data.user_id` field
- **OpenBox:** `player` field
- **Alternar:** None (external handling)

**Go Implementation:**
```go
func ExtractOperatorContext(payload map[string]interface{}) (operatorSlug string, err error) {
    // Try different field names per provider
    if extraData, ok := payload["ExtraData"].(string); ok {
        // Parse JSON and extract slug
    }
    // ... other extraction logic
}
```

---

## Go Implementation Roadmap

### Phase 1: Foundation
1. Base provider service interface
2. Token sanitization utilities
3. Signature validation library (MD5, HMAC-SHA256)
4. Operator context extraction

### Phase 2: Generic Router Implementation
1. Generic routing handler for Mancala, Digitain, Evoplay
2. Service registry for endpoint dispatch
3. Request/response transformation pipelines

### Phase 3: Special Handlers
1. OpenBox PUT handler with signature headers
2. Digitain batch processor for multi-operator items
3. Alternar HTTP proxy handler

### Phase 4: Testing
1. Unit tests for signature validation
2. Integration tests with mock operator callbacks
3. Concurrency tests for batch processing
4. Error handling scenarios

### Phase 5: Optimization
1. Connection pooling for operator callbacks
2. Batch aggregation for multi-operator requests
3. Circuit breaker for operator endpoint failures
4. Caching for operator credentials

---

## Known Limitations & Quirks

### Mancala
- ExtraData parsing required (non-standard operator binding)
- MD5 is weak cryptographic hash (security consideration)

### Digitain
- Batch processing adds complexity to request/response handling
- checktxstatus endpoint bypasses operator (special case handling)
- Multiple credential fields (operatorId, secret-key)

### Evoplay
- Complex MD5 algorithm with recursive field traversal (unique implementation)
- Signature not validated by gateway (only generated for forwarding)
- projectId, apiVersion, secretKey all required

### OpenBox
- Non-standard PUT method for state-changing operations (unexpected pattern)
- Base64URL encoding (different from standard Base64)
- Signature validation via headers (not body)

### Alternar
- Hardcoded external URL (not flexible per operator)
- No cryptographic validation (relies on HTTPS/TLS only)
- Acts as HTTP proxy (unusual for webhook gateway)

---

## File References (Legacy PHP Implementation)

**Controllers:**
- `app/Http/Controllers/MancalaController.php` (24 lines)
- `app/Http/Controllers/DigitainController.php` (26 lines)
- `app/Http/Controllers/EvoplayController.php` (22 lines)
- `app/Http/Controllers/OpenBoxController.php` (59 lines)
- `app/Http/Controllers/AlternarRedirectController.php` (20 lines)

**Services:**
- `app/Services/MancalaService.php` (75 lines, 4 endpoints, MD5 signature)
- `app/Services/DigitainService.php` (208 lines, 11 endpoints, batch processing, HMAC-SHA256)
- `app/Services/EvoplayService.php` (103 lines, 5 endpoints, complex MD5)
- `app/Services/OpenBoxService.php` (187 lines, 4+1 operations, HMAC-SHA256, PUT headers)

**Tests (inferred from pattern):**
- No dedicated test files found for remaining providers in current scan

---

## Summary

The 5 remaining providers showcase diverse integration patterns:

1. **Mancala, Digitain, Evoplay** - Generic routing allows flexible endpoint handling
2. **OpenBox** - PUT method indicates idempotent semantics, unusual pattern
3. **Alternar** - HTTP proxy redirect, simplest pattern

**Complexity Ranking:**
- **Simplest:** Alternar (HTTP proxy)
- **Simple:** Mancala (4 operations, MD5)
- **Medium:** OpenBox (4 operations, HMAC-SHA256, PUT headers), Evoplay (5 operations, complex MD5)
- **Most Complex:** Digitain (11 endpoints, batch processing, HMAC-SHA256)

**Go Implementation Priority:**
1. Generic routing infrastructure (benefits 3 providers)
2. Signature validation library (benefits 4 providers)
3. Token sanitization utilities (benefits all)
4. Batch processor for Digitain
5. Header-based signature for OpenBox

With all 5 remaining providers analyzed, combined with the 3 previously documented (Pragmatic Play, Evolution Gaming, PG Soft), the casino proxy migration now has **complete provider intelligence** for all 8 integrations.
