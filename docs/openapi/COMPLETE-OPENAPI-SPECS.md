# Complete OpenAPI Specifications - All Providers

**Date:** 2026-05-11  
**Status:** Draft - Ready for Validation  
**Includes:** 7 Gaming Providers + Internal APIs + System Endpoints

---

# TABLE OF CONTENTS

1. [Mancala Provider](#mancala-provider)
2. [Digitain RGS Provider](#digitain-rgs-provider)
3. [PG Soft Provider](#pg-soft-provider)
4. [Evoplay Provider](#evoplay-provider)
5. [Evolution Gaming Provider](#evolution-gaming-provider)
6. [OpenBox Provider](#openbox-provider)
7. [Alternar Provider](#alternar-provider)
8. [Internal Admin APIs](#internal-admin-apis)
9. [System Endpoints](#system-endpoints)

---

# MANCALA PROVIDER

**Endpoint:** `POST /v1/webhooks/mancala/{endpoint}`  
**Allowed Methods:** Balance, Credit, Debit, Refund  
**Authentication:** MD5 Hash  
**Token Format:** `{operator_slug}_{actual_token}`

## Hash Calculation

```
hash = MD5(
  endpoint + "/" +
  sessionId +
  [transactionGuid?] +
  [refundTransactionGuid?] +
  [roundGuid?] +
  [amount?] +
  secretToken
)
```

## Endpoints

### POST /v1/webhooks/mancala/Balance

Get player balance.

**Request:**
```json
{
  "SessionId": "operator_slug_session123",
  "ExtraData": "{\"operator_slug\": \"my_operator\"}",
  "Hash": "md5_hash_value"
}
```

**Response:**
```json
{
  "balance": 1500.50,
  "currency": "BRL"
}
```

---

### POST /v1/webhooks/mancala/Credit

Credit player account (win/payout).

**Request:**
```json
{
  "SessionId": "operator_slug_session123",
  "TransactionGuid": "guid_value",
  "Amount": 100.50,
  "RoundGuid": "round_guid",
  "ExtraData": "{\"operator_slug\": \"my_operator\"}",
  "Hash": "md5_hash_value"
}
```

**Response:**
```json
{
  "success": true,
  "balance": 1600.00,
  "transactionId": "txn_123"
}
```

---

### POST /v1/webhooks/mancala/Debit

Debit player account (bet placement).

**Request:**
```json
{
  "SessionId": "operator_slug_session123",
  "TransactionGuid": "guid_value",
  "Amount": 50.00,
  "RoundGuid": "round_guid",
  "ExtraData": "{\"operator_slug\": \"my_operator\"}",
  "Hash": "md5_hash_value"
}
```

**Response:**
```json
{
  "success": true,
  "balance": 1450.00,
  "transactionId": "txn_124"
}
```

---

### POST /v1/webhooks/mancala/Refund

Refund transaction.

**Request:**
```json
{
  "SessionId": "operator_slug_session123",
  "RefundTransactionGuid": "guid_value",
  "Amount": 50.00,
  "ExtraData": "{\"operator_slug\": \"my_operator\"}",
  "Hash": "md5_hash_value"
}
```

**Response:**
```json
{
  "success": true,
  "balance": 1500.00,
  "transactionId": "txn_125"
}
```

---

# DIGITAIN RGS PROVIDER

**Endpoint:** `POST /v1/webhooks/digitain-rgs/{endpoint}`  
**Allowed Methods:** authenticate, getbalance, bet, win, betwin, refund, amend, checktxstatus, charge, promowin, refreshtoken  
**Authentication:** HMAC-SHA256  
**Token Format:** `{operator_slug}_{actual_token}` OR playerId

## Hash Calculation

```
hash = HMAC-SHA256(
  timestamp + operatorId,
  secretKey
)
```

## Endpoints (Sample)

### POST /v1/webhooks/digitain-rgs/authenticate

Authenticate player.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value"
}
```

**Response (Success: errorCode === 1):**
```json
{
  "errorCode": 1,
  "playerId": "operator_slug_player456",
  "balance": 1000.00,
  "currency": "BRL"
}
```

**Response (Error):**
```json
{
  "errorCode": 2,
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value",
  "message": "Invalid token"
}
```

---

### POST /v1/webhooks/digitain-rgs/getbalance

Get player balance.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value"
}
```

**Response:**
```json
{
  "errorCode": 1,
  "playerId": "operator_slug_player456",
  "balance": 1000.00,
  "currency": "BRL"
}
```

---

### POST /v1/webhooks/digitain-rgs/bet

Place bet (debit).

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "bet_info",
      "amount": 50.00
    }
  ],
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value"
}
```

**Response:**
```json
{
  "errorCode": 1,
  "items": [
    {
      "info": "bet_info",
      "balance": 950.00,
      "transactionId": "txn_126"
    }
  ]
}
```

---

# PG SOFT PROVIDER

**Endpoint:** `POST /v1/webhooks/pgsoft/{method}`  
**Methods:** VerifySession, Cash/Get, Cash/TransferInOut, Cash/Adjustment  
**Authentication:** Provider signature validation  
**Token Format:** `{operator_slug}_{actual_token}`

## Endpoints

### POST /v1/webhooks/pgsoft/VerifySession

Verify player session.

**Request:**
```json
{
  "operator_player_session": "operator_slug_session123",
  "playerId": "player_456"
}
```

**Response:**
```json
{
  "error": null,
  "data": {
    "player_name": "operator_slug_john_doe",
    "balance": 1000.00,
    "currency": "BRL",
    "country": "BR"
  }
}
```

---

### POST /v1/webhooks/pgsoft/Cash/Get

Get player cash balance.

**Request:**
```json
{
  "operator_player_session": "operator_slug_session123",
  "playerId": "player_456"
}
```

**Response:**
```json
{
  "error": null,
  "balance": 1000.00,
  "currency": "BRL"
}
```

---

### POST /v1/webhooks/pgsoft/Cash/TransferInOut

Transfer credits in/out (bet, win).

**Request:**
```json
{
  "operator_player_session": "operator_slug_session123",
  "transactionId": "txn_127",
  "amount": 50.00,
  "transactionType": "bet|win",
  "gameId": "game_123",
  "roundId": "round_123"
}
```

**Response:**
```json
{
  "error": null,
  "balance": 950.00,
  "transactionId": "txn_127",
  "status": "success"
}
```

---

### POST /v1/webhooks/pgsoft/Cash/Adjustment

Adjust player balance (refund, correction).

**Request:**
```json
{
  "operator_player_session": "operator_slug_session123",
  "amount": 50.00,
  "reason": "refund|correction",
  "referenceId": "ref_123"
}
```

**Response:**
```json
{
  "error": null,
  "balance": 1000.00,
  "status": "success"
}
```

---

# EVOPLAY PROVIDER

**Endpoint:** `POST /v1/webhooks/evoplay`  
**Allowed Methods:** init, bet, win, refund, BalanceIncrease  
**Authentication:** Signature validation  
**Token Format:** `{operator_slug}_{actual_token}`

## Endpoints

### POST /v1/webhooks/evoplay

Generic Evoplay webhook.

**Request:**
```json
{
  "name": "init|bet|win|refund|BalanceIncrease",
  "token": "operator_slug_token123",
  "data": {
    "user_id": "user_456",
    "game_id": "game_789",
    "amount": 100.50,
    "transaction_id": "txn_128"
  }
}
```

**Response:**
```json
{
  "success": true,
  "balance": 1100.50,
  "transactionId": "txn_128"
}
```

---

# EVOLUTION GAMING PROVIDER

**Endpoint:** `POST /v1/webhooks/evolution/{method}`  
**Methods:** authentication, debit, credit, rollback, getNewToken  
**Authentication:** HMAC-SHA256 (base64 encoded)  
**Token Format:** `{operator_slug}_{actual_token}`  
**Signature Header:** `hash`

## Hash Calculation

```
signature = base64(HMAC-SHA256(
  rawContent (request body with operator prefix removed),
  secret_key
))
```

## Endpoints

### POST /v1/webhooks/evolution/authentication

Authenticate player.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "gameId": "game_789",
  "sessionId": "session_456"
}
```

**Headers:**
```
hash: base64_encoded_hmac_sha256
```

**Response:**
```json
{
  "playerId": "player_456",
  "balance": 1000.00,
  "currency": "BRL",
  "country": "BR",
  "status": "active"
}
```

---

### POST /v1/webhooks/evolution/debit

Debit player account (bet).

**Request:**
```json
{
  "token": "operator_slug_token123",
  "transactionId": "txn_129",
  "amount": 50.00,
  "gameId": "game_789"
}
```

**Response:**
```json
{
  "balance": 950.00,
  "status": "success"
}
```

---

### POST /v1/webhooks/evolution/credit

Credit player account (win).

**Request:**
```json
{
  "token": "operator_slug_token123",
  "transactionId": "txn_130",
  "amount": 100.00,
  "gameId": "game_789"
}
```

**Response:**
```json
{
  "balance": 1050.00,
  "status": "success"
}
```

---

### POST /v1/webhooks/evolution/rollback

Rollback transaction.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "originalTransactionId": "txn_129",
  "gameId": "game_789"
}
```

**Response:**
```json
{
  "balance": 1000.00,
  "status": "rolled_back"
}
```

---

### POST /v1/webhooks/evolution/getNewToken

Get new session token.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "gameId": "game_789"
}
```

**Response:**
```json
{
  "newToken": "operator_slug_newtoken456",
  "expiresAt": "2026-05-12T14:30:00Z"
}
```

---

# OPENBOX PROVIDER

**Endpoint:** `PUT /v1/webhooks/openbox/seamless/{method}`  
**Methods:** balance, bet, win, refund  
**HTTP Method:** PUT (non-standard)  
**Authentication:** HMAC-SHA256 signature  
**Token Format:** `{operator_slug}_{actual_token}`  
**Signature Format:** `{api_key}:{base64_signature}`

## Signature Calculation

```
queryString = sort_by_key(data) → "key1=value1&key2=value2&..."
urlQuery = queryString + ":" + secret_key + timestamp
hmacSha256 = HMAC-SHA256(urlQuery, secret_key)
base64_signature = rtrim(base64(hmacSha256), "=")
signature = api_key + ":" + base64_signature
```

## Endpoints

### PUT /v1/webhooks/openbox/seamless/balance

Get player balance.

**Request:**
```json
{
  "playerId": "operator_slug_player456",
  "timestamp": "1715424600"
}
```

**Headers:**
```
signature: api_key:base64_hmac_sha256
```

**Response:**
```json
{
  "balance": 1000.00,
  "currency": "BRL",
  "status": "success"
}
```

---

### PUT /v1/webhooks/openbox/seamless/bet

Place bet.

**Request:**
```json
{
  "playerId": "operator_slug_player456",
  "amount": 50.00,
  "gameId": "game_789",
  "roundId": "round_456",
  "timestamp": "1715424600"
}
```

**Response:**
```json
{
  "balance": 950.00,
  "transactionId": "txn_131",
  "status": "success"
}
```

---

### PUT /v1/webhooks/openbox/seamless/win

Record win.

**Request:**
```json
{
  "playerId": "operator_slug_player456",
  "amount": 100.00,
  "gameId": "game_789",
  "roundId": "round_456",
  "timestamp": "1715424600"
}
```

**Response:**
```json
{
  "balance": 1050.00,
  "transactionId": "txn_132",
  "status": "success"
}
```

---

### PUT /v1/webhooks/openbox/seamless/refund

Refund bet.

**Request:**
```json
{
  "playerId": "operator_slug_player456",
  "originalTransactionId": "txn_131",
  "amount": 50.00,
  "timestamp": "1715424600"
}
```

**Response:**
```json
{
  "balance": 1000.00,
  "status": "refunded"
}
```

---

# ALTERNAR PROVIDER

**Endpoint:** `POST /v1/webhooks/alternar/`  
**Description:** Redirect handler  
**Authentication:** Provider signature validation  

## Endpoints

### POST /v1/webhooks/alternar/

Handle Alternar redirect/webhook.

**Request:**
```json
{
  "playerId": "operator_slug_player456",
  "token": "operator_slug_token123",
  "action": "launch|logout",
  "gameId": "game_789",
  "signature": "hash_value"
}
```

**Response:**
```json
{
  "success": true,
  "redirectUrl": "https://game.alternar.com/...",
  "sessionToken": "token_value"
}
```

---

# INTERNAL ADMIN APIS

**Middleware:** `internalCheck` (private access only)  
**Base Path:** `/v1/internal/`

## POST /v1/internal/operator/store

Create new operator.

**Request:**
```json
{
  "slug": "new_operator",
  "name": "New Operator Name",
  "url": "https://api.newoperator.com"
}
```

**Response:**
```json
{
  "id": "uuid",
  "slug": "new_operator",
  "name": "New Operator Name",
  "status": "active"
}
```

---

## GET /v1/internal/operator/

List all operators.

**Response:**
```json
{
  "operators": [
    {
      "id": "uuid",
      "slug": "operator_slug",
      "name": "Operator Name",
      "url": "https://operator.api.com",
      "status": "active"
    }
  ]
}
```

---

## DELETE /v1/internal/operator/delete/{slug}

Delete operator.

**Response:**
```json
{
  "success": true,
  "message": "Operator deleted"
}
```

---

## GET /v1/internal/credentials/

List credentials (admin panel).

**Response:**
```json
{
  "credentials": [
    {
      "id": "uuid",
      "operator_id": "uuid",
      "name": "pragmatic",
      "key": "secret-key",
      "value": "***hidden***"
    }
  ]
}
```

---

# SYSTEM ENDPOINTS

## GET /v1/status

Health check endpoint.

**Response:**
```json
{
  "message": "OK"
}
```

---

## POST /v1/image

Upload operator image/branding.

**Request:** Multipart form-data

```
Content-Type: multipart/form-data
Body: {
  "file": <binary>,
  "operator_id": "uuid"
}
```

**Response:**
```json
{
  "url": "https://storage.casino-proxy.com/images/...",
  "filename": "image_name.jpg"
}
```

---

## POST /v1/entry

Player session entry point.

**Request:**
```json
{
  "operator_slug": "my_operator",
  "player_id": "player_456",
  "game_id": "game_789",
  "token": "session_token"
}
```

**Response:**
```json
{
  "success": true,
  "session": {
    "id": "session_123",
    "player_id": "player_456",
    "balance": 1000.00,
    "token": "session_token"
  }
}
```

---

# SUMMARY OF SPECIFICATIONS

| Provider | Endpoints | Auth Method | Token Format | Status |
|----------|-----------|-------------|--------------|--------|
| Pragmatic Play | 9 | MD5 Hash | `{slug}_{value}` | ✅ Existing |
| Mancala | 4 | MD5 Hash | `{slug}_{value}` | 📝 NEW |
| Digitain RGS | 11 | HMAC-SHA256 | `{slug}_{value}` | 📝 NEW |
| PG Soft | 4 | Signature | `{slug}_{value}` | 📝 NEW |
| Evoplay | 5 | Signature | `{slug}_{value}` | 📝 NEW |
| Evolution | 5 | HMAC-SHA256 | `{slug}_{value}` | 📝 NEW |
| OpenBox | 4 | HMAC-SHA256 | `{slug}_{value}` | 📝 NEW |
| Alternar | 1 | Signature | `{slug}_{value}` | 📝 NEW |
| **Internal Admin** | 3 | middleware | N/A | 📝 NEW |
| **System** | 3 | None | N/A | ✅ Partial |

**Total Endpoints:** 49 (41 gaming webhooks + 3 internal + 3 system + 2 auth)

---

**Status:** ⏸️ Awaiting Validation

Next: Separate into individual YAML files per provider + create internal-admin-spec.yaml
