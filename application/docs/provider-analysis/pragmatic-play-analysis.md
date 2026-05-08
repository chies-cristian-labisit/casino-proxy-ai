# Pragmatic Play Provider Analysis

**Date:** 2026-05-08  
**Analyst:** @dev (Dex)  
**Status:** Complete  
**Branch:** `feature/casino-proxy-migration-1/1.1-analyze-pragmatic-play`

---

## Executive Summary

Pragmatic Play is integrated via a **generic webhook router** pattern. The gateway accepts dynamic endpoints (e.g., `/webhooks/pragmatic-play/{endpoint}`) and routes to corresponding service methods. The integration includes **9 endpoints** covering authentication, balance operations, betting, and game results.

**Key Characteristics:**
- **Pattern:** Generic routing with method dispatch
- **Authentication:** HMAC-MD5 signature validation
- **Operator Context:** Embedded in token/userId fields using format `{operator_slug}_{actual_value}`
- **Callback Mechanism:** Direct postback to operator URL for transparent provider integration

---

## Integration Architecture

### Request Flow

```
Pragmatic Play
     ↓
POST /v1/webhooks/pragmatic-play/{endpoint}
     ↓
PragmaticPlayController::__invoke()
     ↓
PragmaticPlayService::call($endpoint, $data)
     ↓
PragmaticPlayService::{$endpoint}()  (dynamic method dispatch)
     ↓
Validate signature → Extract operator → Forward to operator URL → Return response
```

### Key Components

**Controller:** `app/Http/Controllers/PragmaticPlayController.php`
- Single `__invoke()` method handles all endpoints
- Delegates to `PragmaticPlayService`

**Service:** `app/Services/PragmaticPlayService.php`
- Implements 9 methods (one per endpoint)
- Handles signature validation, operator extraction, and provider communication
- Uses `OperatorService` for operator context

**Operator Context:**
- Tokens arrive in format: `{operator_slug}_{token}`
- Service strips the slug to get original token
- Adds slug to userId in response

---

## Endpoints (9 Total)

### 1. Authentication Endpoint

**Path:** `/v1/webhooks/pragmatic-play/authenticate.html`

**Method:** `authenticate(array $data): array`

**Purpose:** Authenticate player session and retrieve initial player information.

**Request Fields:**
- `providerId` (string, required): "PragmaticPlay"
- `token` (string, required): Format `{operator_slug}_{actual_token}`
- `hash` (string, required): HMAC-MD5 signature

**Response Fields:**
- `userId` (string): Player ID with operator prefix added
- `currency` (string): Player's currency (e.g., "BRL")
- `cash` (number): Cash balance
- `bonus` (number): Bonus balance
- `country` (string): Player's country code
- `jurisdiction` (integer): Jurisdiction ID
- `betLimits` (object):
  - `defaultTotalBet` (number)
  - `minTotalBet` (number)
  - `maxTotalBet` (number)
- `error` (integer): 0 = success, non-zero = error
- `description` (string): Status message

**Implementation Notes:**
- Called when player launches a Pragmatic Play game
- Gateway forwards request to operator's URL with signature
- Operator returns player profile, gateway enriches with operator context

---

### 2. Balance Endpoint

**Path:** `/v1/webhooks/pragmatic-play/balance.html`

**Method:** `balance(array $data): array`

**Purpose:** Get current player balance (cash + bonus).

**Request Fields:**
- `providerId` (string, required)
- `token` (string, required): Session token
- `userId` (string, required): Format `{operator_slug}_{player_id}`
- `hash` (string, required)

**Response Fields:**
- `transactionId` (string): Transaction reference
- `currency` (string)
- `cash` (number)
- `bonus` (number)
- `usedPromo` (number): Amount used from promo balance
- `error` (integer)
- `description` (string)

**Implementation Notes:**
- Called frequently during gameplay
- Single signature validation - token and userId can be used interchangeably
- Must return current balance immediately

---

### 3-5. Betting Endpoints

#### Bet Endpoint

**Path:** `/v1/webhooks/pragmatic-play/bet.html`

**Method:** `bet(array $data): array`

**Purpose:** Record a bet placed by player (debit operation).

**Request Fields:**
- `providerId` (string, required)
- `reference` (string, required): Unique bet reference
- `gameId` (string, required): Game identifier
- `roundId` (string, required): Game round identifier
- `roundDetails` (string, required): JSON with round-specific data
- `amount` (number, required): Bet amount
- `timestamp` (string, required): ISO 8601 timestamp
- `bonusCode` (string, optional): Applied bonus code
- `userId` (string, required)
- `hash` (string, required)

**Response Fields:** Same as balance response

**Implementation Notes:**
- Must validate player has sufficient balance
- Reference should be unique within operator context
- Idempotency: Re-submission of same reference should return error or existing response

#### Refund Endpoint

**Path:** `/v1/webhooks/pragmatic-play/refund.html`

**Method:** `refund(array $data): array`

**Purpose:** Refund a previously placed bet (credit operation).

**Request Fields:**
- `providerId` (string)
- `reference` (string, required): Original bet reference to refund
- `userId` (string, required)
- `hash` (string)

**Response Fields:**
- `transactionId` (string): Refund transaction ID
- `error` (integer)
- `description` (string)

**Implementation Notes:**
- Must find original bet by reference
- Check bet hasn't already been refunded
- Refund only applies if bet exists and is valid

---

### 6-9. Game Result Endpoints

All result endpoints follow similar pattern. Differ by endpoint name (result, bonusWin, jackpotWin, promoWin).

**Common Path Pattern:** `/v1/webhooks/pragmatic-play/{result_type}.html`

**Common Method Pattern:** `{result_type}(array $data): array`

**Purpose:** Record game result and credit winnings to player.

**Request Fields (Same for all):**
- `providerId` (string)
- `reference` (string, required): Unique result reference
- `gameId` (string)
- `roundId` (string)
- `roundDetails` (string, optional): JSON details
- `amount` (number, required): Win amount
- `timestamp` (string)
- `promoWinAmount` (number, optional): Promo win amount
- `promoWinReference` (string, optional)
- `promoCampaignID` (string, optional)
- `promoCampaignType` (string, optional)
- `bonusCode` (string, optional)
- `userId` (string)
- `hash` (string)

**Result Types:**
1. **result** - Regular game win
2. **bonusWin** - Bonus round win
3. **jackpotWin** - Jackpot win (may trigger special handling)
4. **promoWin** - Promotional/campaign win

**Implementation Notes:**
- All use same underlying handler: `handleResult()`
- Functionally identical, difference is semantic (for reporting/analytics)
- Must support promo-specific metadata

---

### 10. Adjustment Endpoint

**Path:** `/v1/webhooks/pragmatic-play/adjustment.html`

**Method:** `adjustment(array $data): array`

**Purpose:** Manual balance adjustment (for corrections, bonuses, customer service).

**Request Fields:**
- `providerId` (string)
- `reference` (string, required): Unique adjustment reference
- `gameId` (string, required)
- `roundId` (string, required)
- `amount` (number, required): Can be positive or negative
- `userId` (string)
- `hash` (string)

**Response Fields:** Same as balance response

**Implementation Notes:**
- Used for manual corrections, promotional bonuses
- Amount can be negative (debit) or positive (credit)
- Requires audit trail and operator approval (in production)

---

## Authentication & Signature Validation

### Algorithm: HMAC-MD5

**Step-by-step:**

1. Remove `hash` field from payload
2. Sort payload keys alphabetically
3. Build query string: `key1=value1&key2=value2&...{secret_key}`
4. URL-decode the entire string
5. Calculate MD5 hash

**Example:**

```php
Input payload:
{
  "providerId": "PragmaticPlay",
  "token": "op_abc123",
  "hash": "xyz"
}

Secret: "my_secret_key"

Step 1: Remove hash
{
  "providerId": "PragmaticPlay",
  "token": "op_abc123"
}

Step 2: Sort alphabetically (already sorted)

Step 3: Build query string
"providerId=PragmaticPlay&token=op_abc123my_secret_key"

Step 4: URL decode (no change in this example)
"providerId=PragmaticPlay&token=op_abc123my_secret_key"

Step 5: MD5
md5("providerId=PragmaticPlay&token=op_abc123my_secret_key") = "5d41402abc4b2a76b9719d911017c592"
```

### Signature Validation Process

**PHP Implementation:**
```php
public function generateHashCode(array $payload, string $secretKey): string
{
    unset($payload['hash']);
    ksort($payload);
    $urlQuery = http_build_query($payload) . $secretKey;
    $urlQuery = rawurldecode($urlQuery);
    return md5($urlQuery);
}
```

**Validation:**
```php
// Incoming request hash
$requestHash = $data['hash'];

// Calculated hash
$expectedHash = $this->generateHashCode($data, $secretKey);

// Verify
if ($requestHash !== $expectedHash) {
    // Invalid signature - reject request
    return ['error' => 1, 'description' => 'Invalid signature'];
}
```

### Credential Management

- Secret key stored in `credentials` table
- Linked to operator and provider name
- Query: `$operator->credentials()->where('name', 'pragmatic')->where('key', 'secret-key')->first()->value`

---

## Operator Context & Multi-Tenancy

### Operator Binding

Pragmatic Play tokens embed operator context using a `{operator_slug}_{value}` format:

**Request:**
```json
{
  "token": "my-operator_abc123xyz",
  "userId": "my-operator_user_9876"
}
```

**Processing:**
1. Gateway receives request
2. Extracts operator slug from token: `my-operator`
3. Looks up operator and credentials
4. Strips slug to get actual token: `abc123xyz`
5. Forwards request to operator with clean token
6. Receives response and adds operator slug to userId
7. Returns to Pragmatic Play

### Multi-Operator Support

- Each operator has separate credentials (secret key)
- Same game can be played by players from different operators
- Balance and transactions isolated per operator
- No cross-operator balance transfers possible

---

## Data Precision & Constraints

### Amount Fields

**Decimal Precision:** Minimum 2 decimal places (cents)

**Example Values:**
- `1000.50` BRL
- `100.00` EUR
- `50.00` USD

**Validation:**
- Amounts > 0 (no zero or negative bets)
- Sufficient balance before bet placement
- Reasonable bet limits per operator

### Currency

- ISO 4217 code (e.g., "BRL", "EUR", "USD")
- Player's account currency set at authentication
- All balance operations in same currency

### Timestamps

- ISO 8601 format: `2026-05-08T10:30:00Z`
- Server must track timestamp for audit trail
- Can validate against server time for replay protection

---

## Error Handling

### Error Code Convention

- `error: 0` = Success
- `error: 1` = Generic error
- `error: 2` = Invalid request
- `error: 3` = Insufficient balance
- `error: 4` = Invalid player
- `error: 5` = Invalid operator

### Common Error Scenarios

| Scenario | HTTP Status | Error Code | Description |
|----------|------------|-----------|---|
| Invalid signature | 401 | 1 | Invalid signature |
| Unknown endpoint | 404 | 1 | Endpoint not found |
| Missing field | 400 | 2 | Missing required field |
| Insufficient balance | 409 | 3 | Insufficient balance for bet |
| Invalid player | 404 | 4 | Player not found |
| Invalid operator | 403 | 5 | Operator not found or access denied |
| Duplicate reference | 409 | 1 | Duplicate transaction reference |
| Server error | 500 | 1 | Internal server error |

---

## Transaction Flow Examples

### Complete Bet → Win Flow

```
1. Authenticate
   Request: authenticate.html + token
   Response: userId, balance, betLimits
   
2. Place Bet
   Request: bet.html + reference + amount + roundDetails
   Response: Updated balance (cash reduced)
   
3. Game Plays (on Pragmatic Play servers)
   
4. Publish Result
   Request: result.html + reference + win amount
   Response: Updated balance (cash increased)
   
5. Get Final Balance
   Request: balance.html
   Response: Current cash + bonus
```

### Refund Flow (Game Interrupted)

```
1. Place Bet
   Request: bet.html
   Response: Balance reduced
   
2. Game Interrupted
   Request: refund.html + original bet reference
   Response: Balance restored
```

---

## Gateway to Operator Communication

### Transparency Pattern

The gateway acts as a **transparent proxy**. Instead of storing balance itself:

1. Gateway validates signature
2. Gateway extracts operator context
3. Gateway forwards request to **operator's URL** (stored in Operator model)
4. Operator maintains actual balance
5. Operator returns updated balance
6. Gateway enriches response with operator slug
7. Gateway returns to Pragmatic Play

**Example Operator URL:** `https://operator-api.example.com/pragmatic-play`

### Implications for Go Implementation

- No need to store balance in Casino Proxy database
- Operator URL must be stored and maintained
- Network call to operator for every webhook
- Must handle operator timeout/failure gracefully
- Response format must match exactly

---

## Known Limitations & Quirks

### 1. Signature Algorithm Uses MD5

- MD5 is cryptographically weak but used by Pragmatic Play spec
- Must implement exactly as specified
- Cannot substitute with SHA-256 without provider changes

### 2. Operator URL Callback Model

- Casino Proxy doesn't manage balance, operator does
- Operator must be available for every webhook
- Single operator URL per operator (all games callback to same endpoint)
- No load balancing across multiple operator instances at gateway level

### 3. Token Format Embedding

- Operator slug in token requires coordination between Casino Proxy and operator
- Operator must format tokens with slug prefix
- Changing operator slug breaks existing sessions

### 4. No Idempotency Headers

- No standard idempotency key in Pragmatic Play spec
- Must detect duplicates manually (by reference)
- Reference should be globally unique

### 5. Result Types are Semantic Only

- `result`, `bonusWin`, `jackpotWin`, `promoWin` all use same handler
- No functional difference in Casino Proxy
- Difference is for reporting/analytics by the operator

---

## Testing & Test Cases

### From Legacy Tests (`tests/Feature/PragmaticPlayControllerTest.php`)

**Test Cases Available:**

1. **Unknown Endpoint** - Validates error handling for invalid endpoints
2. **Authenticate Success** - Valid authentication flow
3. **Balance Query** - Balance retrieval
4. **Bet Placement** - Successful bet
5. **Refund** - Bet refund
6. **Result** - Game result/win
7. **Bonus Win** - Bonus round win
8. **Jackpot Win** - Jackpot win
9. **Promo Win** - Promo win
10. **Adjustment** - Balance adjustment

**Test Fixtures:**
- Operator model with credentials
- Test secret key: `"test"`
- Sample payloads with all field variations
- HTTP mocking for operator URL calls

---

## Go Implementation Roadmap

### Critical Points

1. **Signature Validation**
   - Implement HMAC-MD5 exactly (crypto/md5 package)
   - URL query encoding + rawurldecode + sort

2. **Operator Context Extraction**
   - Parse `operator_slug_{value}` format
   - Validate operator exists in database
   - Retrieve secret key for signature validation

3. **Provider Communication**
   - Forward requests to operator URL
   - Timeout handling
   - Error response handling

4. **Transaction Atomicity**
   - If operator request fails, how do we handle partial state?
   - Retry strategy
   - Fallback responses

5. **Logging & Audit Trail**
   - Log all requests and responses
   - Include operator slug and player ID
   - Transaction reference for traceability

### Recommended Go Libraries

- **HTTP Client:** `net/http` (standard library)
- **Signature:** `crypto/md5` + custom query encoding
- **Database:** GORM with transaction support
- **Logging:** Structured JSON logging (zap or slog)

---

## References

**Source Files:**
- Controller: `legacy/casino-proxy/app/Http/Controllers/PragmaticPlayController.php`
- Service: `legacy/casino-proxy/app/Services/PragmaticPlayService.php`
- Enum: `legacy/casino-proxy/app/Enums/PragmaticPlayStatusEnum.php`
- Tests: `legacy/casino-proxy/tests/Feature/PragmaticPlayControllerTest.php`

**External:**
- Pragmatic Play Provider Documentation (if available)

---

**Analysis Complete** ✅
**Next Steps:** Create OpenAPI spec (done) and move to next provider (Evolution Gaming - Story 1.2)
