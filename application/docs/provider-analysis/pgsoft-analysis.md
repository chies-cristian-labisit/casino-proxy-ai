# PG Soft Provider Integration Analysis

## Overview

**Provider:** PG Soft (Philippine-based provider)  
**Integration Type:** Session-based HTTP webhooks  
**Authentication:** Operator credentials in request body (no HMAC signature)  
**API Style:** Namespace-style paths with cash operations  
**Endpoints:** 4 (VerifySession, GetCash, TransferInOut, Adjustment)  
**Key Distinction:** Lacks cryptographic signature validation; relies on operator credentials and session tokens

---

## Authentication & Session Protocol

### Session Management Model

Unlike Pragmatic Play (transaction-based) and Evolution Gaming (token-based state machine), PG Soft uses a **session-based protocol** where:

1. **Initial Session:** Operator sends `operator_player_session` in `{operator_slug}_{session_token}` format
2. **Validation:** Gateway verifies operator context and forwards to operator's callback URL
3. **No Signature:** PG Soft does NOT use HMAC-MD5 or HMAC-SHA256 validation
4. **Operator Credentials:** `operator_token` and `secret_key` included in every request body

### Operator Context Binding

All requests include operator identification:

```json
{
  "operator_token": "operator",
  "secret_key": "secret_key",
  "operator_player_session": "my-operator_abc123xyz",
  "player_name": "my-operator_player_001"
}
```

**Token Format:**
- `operator_token`: Static operator identifier (e.g., "operator")
- `secret_key`: Static secret (e.g., "secret_key")
- `operator_player_session`: Composite `{operator_slug}_{session_token}`
- `player_name`: Composite `{operator_slug}_{player_id}`

### Security Implications

**Strengths:**
- Operator token + secret_key two-factor approach
- Session token binding to specific player
- Operator scope isolation via slug prefix

**Weaknesses:**
- No cryptographic signature (rely on HTTPS only)
- Static credentials vulnerable if leaked
- No replay attack prevention (no nonce/timestamp validation in spec)
- No request sequence numbering

---

## Session Verification Flow

### VerifySession Endpoint

**Purpose:** Validate player session and retrieve player information

**Request:**
```json
{
  "operator_token": "operator",
  "secret_key": "secret_key",
  "operator_player_session": "my-operator_abc123xyz",
  "game_id": 123,
  "custom_parameter": "optional_value"
}
```

**Response (Success):**
```json
{
  "data": {
    "player_name": "my-operator_player_001",
    "nickname": "John Player",
    "currency": "BRL"
  },
  "error": null
}
```

**Response (Failure):**
```json
{
  "error": "Invalid session",
  "description": "Session token not found or expired"
}
```

**Implementation Notes:**

1. **Operator Context Extraction:**
   - Extract `operator_slug` from `operator_player_session` (prefix before `_`)
   - Validate `operator_token` + `secret_key` match known operators
   - Bind player to operator scope

2. **Session Validation:**
   - No token expiration mentioned in spec
   - Operator is responsible for session lifetime management
   - Gateway should validate operator owns the session

3. **Currency Context:**
   - Returned with player info (e.g., "BRL" for Brazilian Real)
   - Affects all subsequent balance/transfer operations
   - Must be consistent across calls for same player

### Key Data Points

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| operator_token | string | Yes | Operator ID |
| secret_key | string | Yes | Operator secret |
| operator_player_session | string | Yes | Session in format `{operator_slug}_{token}` |
| game_id | integer | No | Game context |
| player_name | string | No (response only) | Format: `{operator_slug}_{player_id}` |
| nickname | string | No | Player display name |
| currency | string | Yes (response) | ISO 4217 code |

---

## Cash Operations Flow

### 1. GetCash: Retrieve Current Balance

**Purpose:** Query player's current cash balance

**Request:**
```json
{
  "operator_token": "operator",
  "secret_key": "secret_key",
  "player_name": "my-operator_player_001",
  "operator_player_session": "my-operator_abc123xyz",
  "game_id": 123,
  "custom_parameter": "optional_value"
}
```

**Response:**
```json
{
  "data": {
    "currency_code": "BRL",
    "balance_amount": 1000.50,
    "updated_time": 1715177400123
  },
  "error": null
}
```

**Behavior:**
- Read-only operation (no state change)
- Operator is authoritative source of balance
- Timestamp indicates last balance update (milliseconds)
- No caching recommendations in spec

---

### 2. TransferInOut: Handle Bets & Wins

**Purpose:** Execute financial transfers (bet debit, win credit)

**Request (Bet Example - Debit):**
```json
{
  "operator_token": "operator",
  "secret_key": "secret_key",
  "operator_player_session": "my-operator_abc123xyz",
  "player_name": "my-operator_player_001",
  "game_id": 123,
  "parent_bet_id": "bet_001_123456",
  "bet_id": "bet_001_123456",
  "bet_type": 1,
  "currency_code": "BRL",
  "platform": 1,
  "create_time": 1715177400000,
  "updated_time": 1715177400000,
  "bet_amount": 50.00,
  "win_amount": 0.00,
  "transfer_amount": -50.00,
  "transaction_id": "txn_001_bet",
  "wallet_type": "main",
  "is_validate_bet": true,
  "is_feature": false,
  "is_adjustment": false,
  "is_end_round": false
}
```

**Response (Bet Success):**
```json
{
  "data": {
    "currency_code": "BRL",
    "balance_amount": 950.50,
    "updated_time": 1715177400124
  },
  "error": null
}
```

**Transfer Amount Logic:**
- **Negative value (-50.00):** Debit (bet placement)
- **Positive value (+100.00):** Credit (win payout)
- **Absolute value:** Amount to transfer
- **Sign:** Direction of flow

**Financial Precision:**
- Decimal places: 2 (cents in BRL, smallest unit)
- Float handling: 50.00 is exact (no rounding)
- Amount precision: Operator responsible for precision, gateway forwards as-is

**Bet State Lifecycle:**
```
Pending (is_validate_bet=true) 
  → Settled (is_validate_bet=false, is_end_round=true)
  → or Rolled Back (subsequent /Adjustment with negative reversal)
```

**Flags Interpretation:**
- `is_validate_bet`: Bet awaiting confirmation (allow partial reversal)
- `is_end_round`: Marks final bet of game round
- `is_feature`: Bonus feature game (affects payout calculation)
- `is_adjustment`: Manual adjustment flag (bonus, correction)

---

### 3. Adjustment: Manual Balance Adjustment

**Purpose:** Perform manual balance corrections or bonuses

**Request (Bonus Example):**
```json
{
  "operator_token": "operator",
  "secret_key": "secret_key",
  "player_name": "my-operator_player_001",
  "currency_code": "BRL",
  "transfer_amount": 100.00,
  "adjustment_id": 123,
  "adjustment_transaction_id": "adj_001_123",
  "adjustment_time": 1715177400000,
  "transaction_type": 901,
  "bet_type": 1
}
```

**Response:**
```json
{
  "data": {
    "currency_code": "BRL",
    "balance_amount": 1050.50,
    "updated_time": 1715177400125
  },
  "error": null
}
```

**Use Cases:**
- **Positive transfer_amount:** Bonus credit, customer service adjustment, promotion credit
- **Negative transfer_amount:** Penalty debit, reversal of previous adjustment, correction

**Transaction Types (Enum):**
| Code | Meaning | Notes |
|------|---------|-------|
| 901 | Manual Adjustment | Standard correction or bonus |
| Other | (provider-specific) | Not documented in spec |

**Audit Trail:**
- `adjustment_id`: Unique adjustment identifier
- `adjustment_transaction_id`: Transaction reference
- `adjustment_time`: Millisecond timestamp
- Operator responsible for logging reason/description

---

## Transaction State Management

### Bet Lifecycle with State Flags

```
1. BET PLACEMENT (TransferInOut)
   ├─ transfer_amount: -50.00 (negative)
   ├─ is_validate_bet: true (pending confirmation)
   ├─ is_end_round: false (round ongoing)
   └─ Balance: 1000.00 → 950.00

2. GAME RESOLUTION (TransferInOut)
   ├─ transfer_amount: +100.00 (win)
   ├─ is_validate_bet: false (confirmed)
   ├─ is_end_round: true (round complete)
   └─ Balance: 950.00 → 1050.00

3. [OPTIONAL] ROLLBACK (Adjustment)
   ├─ transfer_amount: -100.00 (reversal)
   ├─ Reason: Round cancellation, dispute
   └─ Balance: 1050.00 → 950.00
```

### Idempotency & Duplicate Handling

**No idempotency key documented.** Implementation strategy:

1. **Unique transaction_id per bet:**
   - PG Soft should send same `transaction_id` for duplicate requests
   - Operator (gateway) should detect via: `{operator_id}_{transaction_id}` unique constraint

2. **For Adjustments:**
   - Use `adjustment_transaction_id` + `operator_id` as unique key
   - Prevent double-crediting if PG Soft retries

3. **Duplicate Window:**
   - Not specified by provider
   - Recommend: 24 hours (match bet window)

---

## Concurrency & Error Scenarios

### Race Conditions

**Scenario 1: Simultaneous Bets**
```
Request 1: Bet 50.00 (Player balance: 1000.00)
Request 2: Bet 50.00 (Player balance: 1000.00) [simultaneous]
────────────────────────────────────────────────
Expected: Req1 succeeds (950.00), Req2 fails (insufficient)
Risk: Both succeed if operator uses optimistic locking
```

**Resolution:**
- Operator must use pessimistic locking (row-level locks)
- Database transaction must be SERIALIZABLE
- Timestamp-based optimistic locking not sufficient for financial operations

### Error Scenarios

| Scenario | HTTP Status | Error Message | Recovery |
|----------|------------|---------------|----------|
| Insufficient balance | 400 | "Insufficient balance" | Reject bet, return to player |
| Invalid player | 400 | "Player not found" | Verify operator + session |
| Invalid operator | 400 | "Operator validation failed" | Check credentials |
| Session expired | 400 | "Session expired" | Force re-authentication |
| Invalid currency | 400 | "Currency mismatch" | Return operator currency |
| Database error | 500 | "Server error" | Retry with exponential backoff |

### Provider Quirks

1. **No signature validation:** Vulnerable to man-in-the-middle unless HTTPS enforced
2. **Static credentials:** No token refresh mechanism (unlike Evolution Gaming)
3. **Session binding:** No explicit timeout documented (operator manages lifetime)
4. **Amount precision:** Always 2 decimal places (no sub-cent handling)
5. **Multi-currency:** Per-player currency, not per-operator (BRL for player = BRL for all operations)

---

## Data Model

### Request/Response Schemas

**VerifySessionRequest:**
```typescript
interface VerifySessionRequest {
  operator_token: string;           // e.g., "operator"
  secret_key: string;               // e.g., "secret_key"
  operator_player_session: string;  // Format: {operator_slug}_{session_token}
  game_id?: integer;
  custom_parameter?: string;
}

interface VerifySessionResponse {
  data: {
    player_name: string;      // {operator_slug}_{player_id}
    nickname: string;
    currency: string;         // ISO 4217 (BRL, USD, etc.)
  };
  error: null;
}
```

**CashResponse:**
```typescript
interface CashResponse {
  data: {
    currency_code: string;    // ISO 4217
    balance_amount: number;   // Decimal with 2 places
    updated_time: number;     // Unix ms timestamp
  };
  error: null;
}
```

**TransferInOutRequest:**
```typescript
interface TransferInOutRequest {
  operator_token: string;
  secret_key: string;
  operator_player_session: string;
  player_name: string;
  game_id: integer;
  bet_id: string;                   // Unique per bet
  parent_bet_id?: string;           // For bet chains
  bet_type: integer;
  currency_code: string;
  platform: integer;
  create_time: number;              // Unix ms
  updated_time: number;             // Unix ms
  bet_amount: number;               // Always positive
  win_amount: number;               // 0 for bet, >0 for win
  transfer_amount: number;          // Negative=debit, positive=credit
  transaction_id: string;           // Unique per transaction
  wallet_type?: string;
  is_validate_bet?: boolean;
  is_feature?: boolean;
  is_adjustment?: boolean;
  is_end_round?: boolean;
}
```

**AdjustmentRequest:**
```typescript
interface AdjustmentRequest {
  operator_token: string;
  secret_key: string;
  player_name: string;
  currency_code: string;
  transfer_amount: number;          // Positive=credit, negative=debit
  adjustment_id: integer;
  adjustment_transaction_id: string;
  adjustment_time: number;          // Unix ms
  transaction_type: integer;        // 901 = manual adjustment
  bet_type: integer;
}
```

---

## Go Implementation Roadmap

### Database Schema

```sql
-- Operator context
CREATE TABLE operators (
  id SERIAL PRIMARY KEY,
  operator_token VARCHAR(100) UNIQUE NOT NULL,
  secret_key VARCHAR(100) NOT NULL,
  operator_slug VARCHAR(50) UNIQUE NOT NULL,
  callback_url VARCHAR(500) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Player sessions
CREATE TABLE pgsoft_sessions (
  id SERIAL PRIMARY KEY,
  operator_id INTEGER NOT NULL REFERENCES operators(id),
  player_id VARCHAR(100) NOT NULL,
  session_token VARCHAR(255) NOT NULL,
  player_name VARCHAR(100) NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(operator_id, player_id, session_token)
);

-- Balance ledger (operator-side)
CREATE TABLE pgsoft_balances (
  id SERIAL PRIMARY KEY,
  operator_id INTEGER NOT NULL REFERENCES operators(id),
  player_id VARCHAR(100) NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  balance_amount DECIMAL(15, 2) NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(operator_id, player_id)
);

-- Transaction log (audit trail)
CREATE TABLE pgsoft_transactions (
  id SERIAL PRIMARY KEY,
  operator_id INTEGER NOT NULL REFERENCES operators(id),
  transaction_id VARCHAR(255) NOT NULL,
  transaction_type VARCHAR(20) NOT NULL,  -- 'bet', 'win', 'adjustment'
  player_id VARCHAR(100) NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  transfer_amount DECIMAL(15, 2) NOT NULL,
  balance_before DECIMAL(15, 2) NOT NULL,
  balance_after DECIMAL(15, 2) NOT NULL,
  flags JSONB,  -- { is_validate_bet, is_feature, is_end_round, etc. }
  status VARCHAR(20) DEFAULT 'pending',  -- 'pending', 'completed', 'failed'
  error_message VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(operator_id, transaction_id)
);
```

### Concurrency Control

```go
// Pessimistic locking for bet placement
func (s *PGSoftService) TransferInOut(ctx context.Context, req *TransferRequest) (*CashResponse, error) {
    tx := db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,  // Strongest isolation
    })
    
    // Lock row for update (pessimistic locking)
    var balance decimal.Decimal
    err := tx.WithContext(ctx).
        Clauses(clause.Locking{Strength: "UPDATE"}).
        Model(&PGSoftBalance{}).
        Where("operator_id = ? AND player_id = ?", operatorID, playerID).
        Select("balance_amount").
        Scan(&balance).Error
    
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("lock acquisition failed: %w", err)
    }
    
    // Check balance
    if req.TransferAmount < 0 && balance.Add(req.TransferAmount).IsNegative() {
        tx.Rollback()
        return nil, fmt.Errorf("insufficient balance")
    }
    
    // Update balance
    newBalance := balance.Add(req.TransferAmount)
    err = tx.Model(&PGSoftBalance{}).
        Where("operator_id = ? AND player_id = ?", operatorID, playerID).
        Update("balance_amount", newBalance).Error
    
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("balance update failed: %w", err)
    }
    
    // Log transaction
    err = tx.Create(&PGSoftTransaction{
        OperatorID:      operatorID,
        TransactionID:   req.TransactionID,
        PlayerID:        req.PlayerName,
        TransferAmount:  req.TransferAmount,
        BalanceBefore:   balance,
        BalanceAfter:    newBalance,
        Status:          "completed",
    }).Error
    
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("transaction log failed: %w", err)
    }
    
    tx.Commit()
    
    return &CashResponse{
        Data: &CashData{
            CurrencyCode:  req.CurrencyCode,
            BalanceAmount: newBalance.InexactFloat64(),
            UpdatedTime:   time.Now().UnixMilli(),
        },
        Error: nil,
    }, nil
}
```

### Decimal Precision Handling

```go
import "github.com/shopspring/decimal"

// Always use decimal.Decimal for financial amounts
type TransferRequest struct {
    TransferAmount decimal.Decimal  // Not float64
    BetAmount      decimal.Decimal
    WinAmount      decimal.Decimal
}

// Response uses float64 only for JSON serialization
type CashData struct {
    BalanceAmount float64 `json:"balance_amount"`  // decimal.Decimal → float64 at boundary
}

// Conversion
func (cd *CashData) MarshalJSON() ([]byte, error) {
    type Alias CashData
    return json.Marshal(&struct {
        *Alias
        BalanceAmount string `json:"balance_amount"`
    }{
        Alias:         (*Alias)(cd),
        BalanceAmount: cd.Balance.String(),  // Use string for JSON
    })
}
```

### Operator Context Extraction

```go
func extractOperatorContext(operatorToken, secretKey, sessionToken string) (*OperatorContext, error) {
    // Lookup operator by token + secret
    operator := &Operator{}
    err := db.Where("operator_token = ? AND secret_key = ?", operatorToken, secretKey).
        First(operator).Error
    
    if err != nil {
        return nil, fmt.Errorf("invalid operator credentials")
    }
    
    // Extract operator slug from session token
    parts := strings.Split(sessionToken, "_")
    if len(parts) != 2 {
        return nil, fmt.Errorf("invalid session format: expected {slug}_{token}")
    }
    
    operatorSlug, sessionVal := parts[0], parts[1]
    
    // Verify slug matches operator
    if operator.OperatorSlug != operatorSlug {
        return nil, fmt.Errorf("session slug mismatch: %s != %s", operatorSlug, operator.OperatorSlug)
    }
    
    return &OperatorContext{
        OperatorID:   operator.ID,
        OperatorSlug: operatorSlug,
        SessionToken: sessionVal,
    }, nil
}
```

### Idempotency Handling

```go
func (s *PGSoftService) TransferInOutIdempotent(ctx context.Context, req *TransferRequest) (*CashResponse, error) {
    // Check if transaction already processed
    existing := &PGSoftTransaction{}
    err := db.Where("operator_id = ? AND transaction_id = ?", operatorID, req.TransactionID).
        First(existing).Error
    
    if err == nil {
        // Transaction already exists, return cached response
        if existing.Status == "completed" {
            return &CashResponse{
                Data: &CashData{
                    BalanceAmount: existing.BalanceAfter.InexactFloat64(),
                    UpdatedTime:   existing.CreatedAt.UnixMilli(),
                },
            }, nil
        } else if existing.Status == "failed" {
            return nil, fmt.Errorf("previous attempt failed: %s", existing.ErrorMessage)
        }
    }
    
    // Process new transaction
    return s.TransferInOut(ctx, req)
}
```

---

## Known Limitations & Quirks

1. **No Cryptographic Signature:** Vulnerable unless strict HTTPS + certificate pinning enforced
2. **No Token Refresh:** Unlike Evolution Gaming's GetNewToken
3. **No Explicit Session Timeout:** Operator manages lifetime
4. **No Replay Prevention:** No nonce/request ID validation
5. **Static Credentials:** No rotation mechanism documented
6. **Amount Precision:** Fixed 2 decimals (no support for 3rd decimal place)
7. **No Batch Operations:** Each bet/win requires separate request
8. **Currency per Player:** Not per-operator (affects multi-currency operators)

---

## Provider Comparison Matrix

| Feature | Pragmatic Play | Evolution Gaming | PG Soft |
|---------|---|---|---|
| **Auth Method** | HMAC-MD5 signature | HMAC-SHA256 signature | Operator token + secret_key |
| **Signature Validation** | Query string base64 | JSON body base64 | None |
| **Session Management** | Transaction-based | Token-based state machine | Session-based |
| **Token Refresh** | No | Yes (GetNewToken) | No |
| **Endpoints** | 9 | 5 | 4 |
| **Amount Precision** | 2 decimals | 2 decimals | 2 decimals |
| **Concurrency Control** | Implicit (ORM) | Explicit state machine | Requires pessimistic locking |
| **Idempotency** | Transaction ID | Built-in state machine | Transaction ID (recommend) |
| **Operator Isolation** | slug prefix | slug prefix | slug prefix |
| **Response Format** | Multiple formats | {data, error} | {data, error} |

---

## Migration Risks & Considerations

### High Priority
1. **Security:** Implement HTTPS + certificate pinning (no signature)
2. **Concurrency:** Use SERIALIZABLE isolation + pessimistic locking
3. **Precision:** Use decimal.Decimal library (not float64)
4. **Idempotency:** Implement transaction ID tracking per operator

### Medium Priority
1. **Operator Credential Management:** Secure storage of secret_key
2. **Audit Trail:** Log all transactions for reconciliation
3. **Error Handling:** Implement retry logic with exponential backoff
4. **Testing:** Mock operator callbacks for integration tests

### Low Priority
1. **Caching:** Consider caching GetCash results (operator permission-dependent)
2. **Webhooks:** Confirm PG Soft will send to operator callback URLs
3. **Monitoring:** Alert on error_rate > 5% or balance inconsistencies

---

## Implementation Checklist for Go Microservice

- [ ] GORM models for operators, sessions, balances, transactions
- [ ] HMAC validation removed (use HTTPS/TLS only)
- [ ] Operator context extraction from token + secret_key
- [ ] VerifySession handler with player lookup
- [ ] GetCash handler with balance query
- [ ] TransferInOut handler with pessimistic locking
- [ ] Adjustment handler with reversal logic
- [ ] Idempotency middleware (check transaction_id before processing)
- [ ] Decimal precision using shopspring/decimal
- [ ] Error response formatting ({error: "...", description: "..."})
- [ ] Transaction logging for audit trail
- [ ] Operator callback URL forwarding (transparent proxy)
- [ ] Unit tests for each endpoint
- [ ] Integration tests with real PostgreSQL
- [ ] Load tests for concurrency (5+ simultaneous bets)
- [ ] Monitoring & alerting

---

## File References (Legacy PHP Implementation)

- **Controller:** `legacy/casino-proxy/app/Http/Controllers/PgSoftController.php` (47 lines)
- **Service:** `legacy/casino-proxy/app/Services/PgSoftService.php` (105 lines)
- **Tests:** `legacy/casino-proxy/tests/Feature/PgSoftControllerTest.php` (273 lines, 4 test cases)

---

## Summary

PG Soft's integration is the **simplest of the three providers** analyzed:

- **Simplicity:** 4 endpoints, no HMAC validation, straightforward request/response
- **Trade-off:** Relies entirely on HTTPS for security (no signature validation)
- **Financial Logic:** Straightforward debit/credit (no complex state machine like Evolution)
- **Go Implementation:** Pessimistic locking + decimal precision are the only complex requirements

The transparent proxy pattern applies: gateway validates operator + session, then forwards to operator's callback URL. PG Soft itself does not store balances; the operator does.
