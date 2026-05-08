# Evolution Gaming Provider Analysis

**Date:** 2026-05-08  
**Analyst:** @dev (Dex)  
**Status:** Complete  
**Branch:** `feature/casino-proxy-migration-1/1.2-analyze-evolution`

---

## Executive Summary

Evolution Gaming is integrated via a **transparent proxy pattern** using **HMAC-SHA256 signature validation**. The integration includes **5 endpoints** covering authentication, financial transactions (debit/credit), rollback/refund, and token refresh.

**Key Characteristics:**
- **Pattern:** Direct endpoint routing with transparent proxy
- **Authentication:** HMAC-SHA256 (Base64 encoded) - stronger than Pragmatic Play's MD5
- **Financial Operations:** Explicit debit (bet), credit (win), rollback (refund) with separate transaction IDs
- **Operator Context:** Embedded in token using format `{operator_slug}_{token}`
- **State Management:** Transaction state machine with PENDING → COMPLETED/FAILED → ROLLED_BACK states
- **Session Refresh:** Token lifecycle with refresh capability

---

## Integration Architecture

### Request Flow

```
Evolution Gaming Provider
     ↓
POST /v1/webhooks/evolution/{endpoint}
     ↓
EvolutionController::authentication|debit|credit|rollback|getNewToken
     ↓
EvolutionService::{method}()
     ↓
1. Validate signature (HMAC-SHA256)
2. Extract & sanitize operator context
3. Forward to operator URL
4. Enrich response with operator slug
5. Return to Evolution Gaming
```

### Key Components

**Controller:** `app/Http/Controllers/EvolutionController.php`
- 5 public methods handling each endpoint
- Delegates to EvolutionService
- Returns JSON responses

**Service:** `app/Services/EvolutionService.php`
- HMAC-SHA256 signature calculation
- Token sanitization (removes operator slug prefix)
- Operator context extraction
- HTTP forwarding to operator URL with retry (3 attempts, 100ms backoff)

**Token Format Handling:**
```php
// Incoming token from Evolution:
"my-operator_abc123xyz"

// After sanitization:
"abc123xyz"

// In responses (added back):
"my-operator_{new_token}"
```

---

## Endpoints (5 Total)

### 1. Authentication Endpoint

**Path:** `/v1/webhooks/evolution/authentication`

**Method:** `authentication(array $data, array $headers, string $rawContent): array`

**Purpose:** Authenticate player session and retrieve initial player information. Called when player launches an Evolution game.

**Request Fields:**
- `operatorId` (integer, required): Operator ID
- `token` (string, required): Format `{operator_slug}_{actual_token}`
- `platformId` (integer, required): Platform identifier
- `timestamp` (integer, required): Unix timestamp in milliseconds

**Response Fields:**
- `operatorId` (integer): Operator ID
- `uid` (integer): Player unique ID
- `nickName` (string): Player display name
- `token` (string): New session token with operator prefix
- `balance` (number): Player's cash balance
- `currency` (string): ISO 4217 currency code
- `errorCode` (integer): 0=success, >0=error
- `errorDescription` (string): Status message
- `timestamp` (integer): Response timestamp in milliseconds

**Implementation Notes:**
- Gateway forwards request to operator's URL
- Operator returns player profile
- Gateway enriches response with operator context by prepending operator slug to token

---

### 2. Debit Endpoint (Bet Placement)

**Path:** `/v1/webhooks/evolution/debit`

**Method:** `debit(array $data, array $headers, string $rawContent): array`

**Purpose:** Record a bet placed by player (debit operation). Subtracts from player's balance.

**Request Fields:**
- `operatorId` (integer): Operator ID
- `token` (string): Session token with operator prefix
- `uid` (integer): Player unique ID
- `transactionId` (string): Unique bet transaction ID
- `roundId` (string): Game round identifier
- `gameId` (string): Game identifier
- `tableId` (string): Table identifier (for live games)
- `currency` (string): ISO 4217 currency code
- `debitAmount` (number): Bet amount (minimum 2 decimal places)
- `isEndRound` (boolean): Whether this is the last bet in the round
- `timestamp` (integer): Transaction timestamp in milliseconds

**Response Fields:** Same as balance response (operatorId, uid, token, balance, transactionId, currency, errorCode, errorDescription, timestamp)

**Financial State:**
- **State Before:** Balance = X
- **Operation:** Subtract debitAmount
- **State After:** Balance = X - debitAmount (if successful)
- **Error Handling:** If insufficient balance, return errorCode=3

**Implementation Notes:**
- Amount is converted to float for precision
- Must validate player has sufficient balance
- transactionId should be unique per operator context
- Idempotency: Re-submission of same transactionId should return error or existing response

---

### 3. Credit Endpoint (Win Payout)

**Path:** `/v1/webhooks/evolution/credit`

**Method:** `credit(array $data, array $headers, string $rawContent): array`

**Purpose:** Record game result and credit winnings to player (credit operation). Adds to player's balance.

**Request Fields:**
- `operatorId` (integer): Operator ID
- `token` (string): Session token with operator prefix
- `uid` (integer): Player unique ID
- `transactionId` (string): Unique win transaction ID
- `debitTransactionId` (string): Reference to related bet transaction
- `roundId` (string): Game round identifier
- `gameId` (string): Game identifier
- `tableId` (string): Table identifier
- `currency` (string): ISO 4217 currency code
- `creditAmount` (number): Payout amount
- `returnReason` (integer): Return reason code
- `isEndRound` (boolean): Whether this ends the round
- `timestamp` (integer): Transaction timestamp in milliseconds

**Response Fields:** Same as debit response

**Financial State:**
- **State Before:** Balance = X
- **Operation:** Add creditAmount
- **State After:** Balance = X + creditAmount (if successful)

**Implementation Notes:**
- Amount is converted to float for precision
- Must support multiple credits per bet (partial payouts)
- transactionId must be unique
- debitTransactionId links this credit to the original bet

---

### 4. Rollback Endpoint (Refund/Reversal)

**Path:** `/v1/webhooks/evolution/rollback`

**Method:** `rollback(array $data, array $headers, string $rawContent): array`

**Purpose:** Reverse a previously placed bet (refund operation). Restores debit to player's account.

**Request Fields:**
- `operatorId` (integer): Operator ID
- `token` (string): Session token with operator prefix
- `uid` (integer): Player unique ID
- `transactionId` (string): Unique rollback transaction ID
- `roundId` (string): Game round identifier
- `gameId` (string): Game identifier
- `tableId` (string): Table identifier
- `currency` (string): ISO 4217 currency code
- `rollbackAmount` (number): Original debit amount to refund
- `timestamp` (integer): Transaction timestamp in milliseconds

**Response Fields:** Same as debit response

**Financial State:**
- **State Before:** Balance = X (after debit)
- **Operation:** Add rollbackAmount back
- **State After:** Balance = X + rollbackAmount

**Implementation Notes:**
- Must find original debit by reference (implicit in roundId + gameId)
- Check debit hasn't already been rolled back
- Rollback only applies if original bet exists and is valid
- Used when game is interrupted or player disconnects before results

---

### 5. GetNewToken Endpoint (Session Refresh)

**Path:** `/v1/webhooks/evolution/getNewToken`

**Method:** `getNewToken(array $data, array $headers, string $rawContent): array`

**Purpose:** Refresh player's session token. Allows extending session without re-authentication.

**Request Fields:**
- `operatorId` (integer): Operator ID
- `currentToken` (string): Current session token with operator prefix
- `uid` (integer): Player unique ID
- `gameId` (string): Current game ID
- `tableId` (string): Current table ID (optional)
- `timestamp` (integer): Request timestamp in milliseconds

**Response Fields:**
- `operatorId` (integer): Operator ID
- `token` (string): New session token with operator prefix
- `balance` (number): Current player balance
- `uid` (integer): Player unique ID
- `errorCode` (integer): 0=success, >0=error
- `errorDescription` (string): Status message
- `timestamp` (integer): Response timestamp in milliseconds

**Implementation Notes:**
- Extends session without requiring re-authentication
- Returns current balance along with new token
- Useful for long-running games (live dealer)
- New token replaces old one in player's session

---

## Authentication & Signature Validation

### Algorithm: HMAC-SHA256 (Base64 Encoded)

**Step-by-step:**

1. Take request body as **raw JSON string** (not parsed)
2. Calculate HMAC-SHA256 hash using secret key
3. Base64 encode the result
4. Pass as `hash` header in request to operator

**Example:**

```
Request body:
{"operatorId":1,"token":"my-operator_abc123","timestamp":1715177400123}

Secret key: "my_secret_key"

Step 1: JSON string (already in correct format)
Step 2: HMAC-SHA256 = binary_hash_result
Step 3: Base64 = "xY9pQ2r3Sz5uV6wX7yZ8="
Step 4: Header = hash: "xY9pQ2r3Sz5uV6wX7yZ8="
```

**PHP Implementation:**
```php
public function getSignature(string $message, string $secret_key)
{
    return base64_encode(hash_hmac('sha256', $message, $secret_key, true));
}
```

**Validation Process:**
```php
// Incoming request hash
$requestHash = $headers['hash'];

// Calculated hash using raw body and secret
$expectedHash = $this->getSignature($rawContent, $secretKey);

// Verify
if ($requestHash !== $expectedHash) {
    return ['errorCode' => 1, 'errorDescription' => 'Invalid signature'];
}
```

### Credential Management

- Secret key stored in `credentials` table
- Linked to operator and provider name ('evolution')
- Query: `$operator->credentials()->where('name', 'evolution')->where('key', 'secret_key')->first()->value`

---

## Transaction State Machine

```
┌──────────────────────────────────────────────────────┐
│ Player Starts Game                                   │
│ AUTHENTICATION → Balance retrieved                   │
└─────────────────────┬────────────────────────────────┘
                      │
         ┌────────────▼──────────────┐
         │  TRANSACTION READY        │
         │  (balance confirmed)      │
         └────────────┬──────────────┘
                      │
        ┌─────────────▼──────────────┐
        │ BET PLACED (DEBIT)         │
        │ State: PENDING             │
        │ Balance -= debitAmount     │
        └────────┬──────────┬────────┘
                 │          │
         ┌───────▼──┐   ┌──▼──────────┐
         │ FAIL     │   │ SUCCESS     │
         │ Revert   │   │ COMPLETED   │
         │ balance  │   │             │
         └────┬─────┘   └──┬──────┬───┘
              │            │      │
              │    ┌───────▼─┐    │
              │    │ CREDIT  │    │
              │    │(PAYOUT) │    │
              │    │ PENDING │    │
              │    └──┬──────┘    │
              │       │           │
              │  ┌────▼────┐      │
              │  │ RESULT  │      │
              │  │COMPLETED│      │
              │  └─────────┘      │
              │                   │
         ┌────▼───────────────────▼────┐
         │ GAME END or DISCONNECT      │
         │                             │
         │ If no result: ROLLBACK      │
         │ Return bet amount to player │
         └────────────────────────────┘
```

### State Transitions

| From | To | Operation | Condition |
|------|----|-----------|-----------  |
| READY | PENDING | Debit | Bet placed |
| PENDING | COMPLETED | (wait) | Debit successful |
| PENDING | ROLLED_BACK | Rollback | Bet failed/interrupted |
| COMPLETED | COMPLETED | Credit | Payout |
| COMPLETED | ROLLED_BACK | Rollback | Game result cancelled |

---

## Transaction Flow Examples

### Complete Bet → Win Flow

```
1. AUTHENTICATE
   Request: authentication {token, operatorId, platformId}
   Response: {uid, token, balance, currency}
   State: READY (balance confirmed)

2. PLACE BET (DEBIT)
   Request: debit {token, uid, transactionId, debitAmount, roundId}
   Response: {token, balance -= debitAmount, transactionId}
   State: PENDING → COMPLETED

3. GAME PLAYS (on Evolution servers)
   Evolution determines outcome

4. PUBLISH RESULT (CREDIT)
   Request: credit {token, uid, transactionId, creditAmount, debitTransactionId}
   Response: {token, balance += creditAmount, transactionId}
   State: COMPLETED

5. GET FINAL BALANCE (IMPLICIT)
   Player's balance now reflects: original - bet + payout
```

### Interrupted Game Flow (Rollback)

```
1. AUTHENTICATE
   State: READY

2. PLACE BET (DEBIT)
   Balance = 1000
   Debit 50
   Balance = 950
   State: COMPLETED

3. GAME INTERRUPTED (connection lost)

4. ROLLBACK (REFUND)
   Request: rollback {token, uid, transactionId, rollbackAmount=50}
   Response: {token, balance += 50, transactionId}
   Balance = 1000
   State: ROLLED_BACK
```

### Token Refresh Flow

```
1. AUTHENTICATE
   token = "old-token-xyz"

2. LONG GAME SESSION
   After time-limit approaching

3. GET NEW TOKEN
   Request: getNewToken {currentToken, operatorId}
   Response: {token: "new-token-abc", balance}
   Old token: invalidated
   New token: replaces old one
```

---

## Concurrency & Idempotency

### Race Condition Scenarios

**Scenario 1: Double Debit**
```
Thread 1: POST debit txn_001 amount=50
Thread 2: POST debit txn_001 amount=50  (duplicate)

Expected: First succeeds, second returns error (duplicate txn)
Implementation: Check transactionId uniqueness per operator
```

**Scenario 2: Debit + Rollback Race**
```
Thread 1: POST debit txn_001 amount=50
Thread 2: POST rollback txn_001 amount=50 (simultaneous)

Expected: One completes, other respects state
Implementation: Use database transaction isolation (SERIALIZABLE)
```

**Scenario 3: Multiple Credits per Bet**
```
Debit: txn_001, amount 50
Credit 1: txn_win_1, debitTransactionId=txn_001, amount=100
Credit 2: txn_win_2, debitTransactionId=txn_001, amount=50 (partial payout)

Expected: Both credits applied (partial payouts allowed)
Implementation: Multiple credits allowed if referenced to same debit
```

### Idempotency Strategy

- **transactionId** must be unique per operator per provider
- Duplicate transactionId within reasonable timeframe (e.g., 24 hours) should return:
  - SUCCESS with previous response (safe idempotent)
  - Or ERROR indicating duplicate (requires client retry)
- Operator responsible for tracking and handling duplicates

---

## Data Precision & Constraints

### Amount Fields (Financial Precision)

**Decimal Precision:** Minimum 2 decimal places (cents)

**Example Values:**
- `1000.50` BRL
- `100.00` EUR
- `50.00` USD

**PHP Implementation (Critical):**
```php
public function formatRequestValue($string, $value)
{
    $data = json_decode($string, true);
    $data[$value] = (float) $data[$value];  // ← Convert to float
    return json_encode($data);              // ← Re-encode for signature
}
```

**Why This Matters:**
- JSON may represent amounts as strings
- Must convert to float for precise calculation
- Re-encode before signature calculation to ensure consistency

**Validation:**
- Amounts > 0 (no zero or negative amounts)
- Sufficient balance before debit placement
- Reasonable bet limits per operator (defined by operator)

### Currency

- ISO 4217 code (e.g., "BRL", "EUR", "USD")
- Player's account currency set at authentication
- All transactions must use same currency
- Multi-currency not supported in this integration

### Timestamps

- Unix timestamp in **milliseconds** (not seconds)
- Example: `1715177400123` (13 digits)
- Used for audit trail and replay protection
- Server can validate against current time (allow ±30 seconds)

---

## Error Handling

### Error Code Convention

- `errorCode: 0` = Success
- `errorCode: 1` = Generic error
- `errorCode: 2` = Invalid request
- `errorCode: 3` = Insufficient balance
- `errorCode: 4` = Invalid player
- `errorCode: 5` = Invalid operator

### Common Error Scenarios

| Scenario | Error Code | Description |
|----------|-----------|---|
| Invalid signature | 1 | Invalid signature |
| Unknown endpoint | 1 | Endpoint not found |
| Missing field | 2 | Missing required field |
| Insufficient balance | 3 | Insufficient balance for bet |
| Invalid player | 4 | Player not found |
| Invalid operator | 5 | Operator not found or access denied |
| Duplicate transaction | 1 | Duplicate transaction reference |
| Server error | 1 | Internal server error |

### HTTP Status Code Mapping

| Scenario | HTTP Status |
|----------|------------|
| Success | 200 OK |
| Invalid signature | 401 Unauthorized |
| Invalid request | 400 Bad Request |
| Insufficient balance | 409 Conflict |
| Player not found | 404 Not Found |
| Operator error | 403 Forbidden |
| Server error | 500 Internal Server Error |

---

## Gateway to Operator Communication

### Transparency Pattern

The gateway acts as a **transparent proxy**:

1. Gateway validates HMAC-SHA256 signature using provider's secret key
2. Gateway extracts operator context from token
3. Gateway **forwards request to operator's URL** (stored in Operator model)
4. **Operator maintains actual balance** (not gateway)
5. Operator returns updated balance after transaction
6. Gateway enriches response with operator slug on token
7. Gateway returns response to Evolution Gaming

**Example Operator URL:** `https://operator-api.example.com/evolution`

### Request Forwarding

```
Evolution → Gateway
           ↓
           Verify signature
           Extract operator
           ↓
           Operator → {operator.url}/evolution/{endpoint}
           ↓
           Operator returns response
           ↓
           Gateway enriches token
           ↓
           Gateway → Evolution
```

### Implications for Go Implementation

- **No balance storage in Casino Proxy**: Operator is source of truth
- **Network call per webhook**: Every request involves operator communication
- **Operator availability critical**: If operator URL is down, webhooks fail
- **Retry strategy essential**: 3 retries with backoff (current: 100ms per retry)
- **Timeout handling**: Must handle operator timeouts gracefully
- **Response validation**: Validate operator response format matches spec

---

## Known Limitations & Quirks

### 1. HMAC-SHA256 (Stronger than Pragmatic Play)

- SHA256 is cryptographically stronger than MD5 (Pragmatic Play)
- Must implement exactly as specified
- Base64 encoding required (not hex like some providers)

### 2. Token Lifecycle

- Tokens have implicit TTL (managed by operator, not specified in contract)
- GetNewToken allows session extension
- Old token invalidated when new token issued (implicit)

### 3. Amount Precision Requirements

- Must handle float conversion carefully (see formatRequestValue)
- PHP's `json_encode` may truncate floats without explicit conversion
- Go implementation must use `decimal.Decimal` for financial precision

### 4. Operator URL Single Endpoint

- All Evolution endpoints forward to same operator URL path base
- Operator distinguishes by endpoint name in path or method
- Example: `{operator.url}/evolution/debit`, `{operator.url}/evolution/credit`

### 5. Transaction State Persistence

- No explicit state machine in webhook spec
- State must be tracked by operator or gateway
- Go implementation should use database transactions for atomicity

### 6. Multi-Round Betting

- `isEndRound` flag indicates round completion
- Multiple debits allowed per round (split bets)
- Multiple credits allowed per round (partial payouts)
- Rollback reverses single debit, not entire round

---

## Testing & Test Cases

### From Legacy Tests (`tests/Feature/EvolutionServiceTest.php`)

**Test Cases Available:**
1. **Authentication Success** - Valid authentication flow with token embedding
2. **Authentication Failure** - Missing operator slug (error handling)
3. **Debit Success** - Bet placement with amount precision
4. **Debit Validation** - Balance check and amount handling
5. **Credit Success** - Win payout with amount precision
6. **Rollback Success** - Refund processing with state reversal
7. **GetNewToken Success** - Token refresh mechanism

**Test Fixtures:**
- Operator model with credentials
- Secret key: `"evolution"` (for tests)
- Sample payloads with all field variations
- HTTP mocking for operator URL calls

**Key Test Assertion Pattern:**
```php
Http::assertSent(function ($request) use ($payload, $operator) {
    $requestData = $request->data();
    // Verify operator slug was stripped
    return $operator->slug.'_'.$requestData['token'] === $payload['token'];
});
```

---

## Go Implementation Roadmap

### Critical Implementation Points

1. **HMAC-SHA256 Signature**
   - Use `crypto/hmac` and `crypto/sha256`
   - **Base64 encode** result (not hex)
   - Validate incoming signature in request header

2. **Operator Context Extraction**
   - Parse `{operator_slug}_{token}` format
   - Use `strings.Split()` to extract slug
   - Validate operator exists in database
   - Retrieve secret key for signature validation

3. **Financial Precision**
   - Use `decimal.Decimal` (github.com/shopspring/decimal) for amounts
   - Never use `float64` for money
   - Properly handle JSON unmarshaling of float strings

4. **Provider Communication**
   - Forward requests to operator URL
   - Include signature header (`hash`)
   - Implement retry logic (3 retries, 100ms backoff)
   - Handle operator timeout gracefully

5. **Transaction Atomicity**
   - Use database transactions for financial operations
   - SERIALIZABLE isolation level for race condition prevention
   - Idempotent handling: check transactionId uniqueness

6. **State Management**
   - Track transaction state per debit
   - Handle concurrent debit → rollback scenarios
   - Support multiple credits per debit

7. **Logging & Audit Trail**
   - Log all requests and responses
   - Include operator slug and player ID
   - Transaction reference for traceability
   - Amount values with full precision

### Recommended Go Libraries

- **HTTP Client:** `net/http` (standard library)
- **Signature:** `crypto/hmac` + `crypto/sha256` + `encoding/base64`
- **Financial:** `github.com/shopspring/decimal` (NOT `float64`)
- **Database:** GORM with transaction support (`gorm.io/gorm`)
- **Logging:** Structured JSON logging (`go.uber.org/zap`)
- **Testing:** `testing` + `net/http/httptest` (standard library)

### Database Schema Considerations

```sql
-- Transactions table
CREATE TABLE transactions (
  id BIGSERIAL PRIMARY KEY,
  operator_id INT NOT NULL,
  provider_id INT NOT NULL,
  player_id INT NOT NULL,
  transaction_id VARCHAR(255) NOT NULL,  -- Evolution's transactionId
  type VARCHAR(50) NOT NULL,              -- 'debit', 'credit', 'rollback'
  amount DECIMAL(19,2) NOT NULL,
  currency VARCHAR(3) NOT NULL,
  state VARCHAR(50) NOT NULL,             -- 'pending', 'completed', 'failed', 'rolled_back'
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE(operator_id, transaction_id)     -- Prevent duplicates per operator
);
```

---

## References

**Source Files:**
- Controller: `legacy/casino-proxy/app/Http/Controllers/EvolutionController.php`
- Service: `legacy/casino-proxy/app/Services/EvolutionService.php`
- Tests: `legacy/casino-proxy/tests/Feature/EvolutionServiceTest.php`

**External:**
- Evolution Gaming Integration Guide (if available from provider)
- Legacy PHP codebase structure and conventions

---

**Analysis Complete** ✅  
**Next Steps:** Create OpenAPI spec (done) and proceed to Story 1.3 (PG Soft)

