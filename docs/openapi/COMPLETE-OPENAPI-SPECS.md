# Complete OpenAPI Specifications - All Providers

**Date:** 2026-05-11  
**Status:** Draft - Ready for Validation  
**Includes:** 7 Gaming Providers + Internal APIs + System Endpoints

---

# TABLE OF CONTENTS

1. [Pragmatic Play Provider](#pragmatic-play-provider)
2. [Mancala Provider](#mancala-provider)
3. [Digitain RGS Provider](#digitain-rgs-provider)
4. [PG Soft Provider](#pg-soft-provider)
5. [Evoplay Provider](#evoplay-provider)
6. [Evolution Gaming Provider](#evolution-gaming-provider)
7. [OpenBox Provider](#openbox-provider)
8. [Alternar Provider](#alternar-provider)
9. [Internal Admin APIs](#internal-admin-apis)
10. [System Endpoints](#system-endpoints)

---

# PRAGMATIC PLAY PROVIDER

**Endpoint:** `POST /v1/webhooks/pragmatic-play/{endpoint}`  
**Allowed Methods:** authenticate, balance, bet, refund, result, bonusWin, jackpotWin, promoWin, adjustment  
**Authentication:** MD5 Hash  
**Token Format:** `{operator_slug}_{actual_token}`

## Hash Calculation

```
hash = MD5(
  HTTP_QUERY_STRING (url-encoded params) +
  secret_key
)
```

**Process:**
1. Remove hash from payload
2. Sort payload by keys (ksort)
3. Encode to URL query string
4. Append secret key
5. Generate MD5

## Endpoints

### POST /v1/webhooks/pragmatic-play/authenticate

Authenticate player session with Pragmatic Play backend.

**Request:**
```json
{
  "providerId": "PragmaticPlay",
  "token": "operator_slug_abc123def456",
  "hash": "md5_hash_value"
}
```

**Response (Success: error === 0):**
```json
{
  "userId": "player123",
  "currency": "BRL",
  "cash": 1000.50,
  "bonus": 50.00,
  "country": "BR",
  "jurisdiction": 99,
  "betLimits": {
    "defaultTotalBet": 2000,
    "minTotalBet": 200,
    "maxTotalBet": 2000
  },
  "error": 0,
  "description": "Success"
}
```

**Casino Proxy Processing:** Prefixes userId back with `{operator_slug}_` if error === 0

---

### POST /v1/webhooks/pragmatic-play/balance

Get player balance.

**Request:**
```json
{
  "providerId": "PragmaticPlay",
  "token": "operator_slug_abc123def456",
  "userId": "operator_slug_user456",
  "hash": "md5_hash_value"
}
```

**Response (Success: error === 0):**
```json
{
  "transactionId": "123",
  "currency": "BRL",
  "cash": 1000.00,
  "bonus": 50.00,
  "usedPromo": 0,
  "error": 0,
  "description": "Success"
}
```

---

### POST /v1/webhooks/pragmatic-play/bet

Place bet (debit player account).

**Request:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "bet_ref_123",
  "gameId": "game_789",
  "roundId": "round_456",
  "roundDetails": "game_details",
  "amount": 50.00,
  "timestamp": "2026-05-11T14:30:00Z",
  "userId": "operator_slug_user456",
  "bonusCode": "BONUS2024",
  "hash": "md5_hash_value"
}
```

**Response (Success: error === 0):**
```json
{
  "transactionId": "123",
  "currency": "BRL",
  "cash": 950.00,
  "bonus": 50.00,
  "usedPromo": 0,
  "error": 0,
  "description": "Success"
}
```

---

### POST /v1/webhooks/pragmatic-play/refund

Refund/cancel bet.

**Request:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "bet_ref_123",
  "userId": "operator_slug_user456",
  "hash": "md5_hash_value"
}
```

**Response (Success: error === 0):**
```json
{
  "transactionId": "124",
  "error": 0,
  "description": "Success"
}
```

---

### POST /v1/webhooks/pragmatic-play/result

Bet result notification (win/loss).

**Request:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "result_ref_123",
  "gameId": "game_789",
  "roundId": "round_456",
  "amount": 100.00,
  "timestamp": "2026-05-11T14:35:00Z",
  "userId": "operator_slug_user456",
  "roundDetails": "game_details",
  "promoWinAmount": 10.00,
  "promoWinReference": "promo_ref",
  "promoCampaignID": "campaign_123",
  "promoCampaignType": "type",
  "bonusCode": "BONUS2024",
  "hash": "md5_hash_value"
}
```

**Response (Success: error === 0):**
```json
{
  "transactionId": "125",
  "currency": "BRL",
  "cash": 1050.00,
  "bonus": 50.00,
  "error": 0,
  "description": "Success"
}
```

---

### POST /v1/webhooks/pragmatic-play/bonusWin

Bonus win notification.

**Request:** Same structure as `/result` endpoint

**Response:** Same structure as `/result` endpoint

---

### POST /v1/webhooks/pragmatic-play/jackpotWin

Jackpot win notification.

**Request:** Same structure as `/result` endpoint

**Response:** Same structure as `/result` endpoint

---

### POST /v1/webhooks/pragmatic-play/promoWin

Promotional win notification.

**Request:** Same structure as `/result` endpoint

**Response:** Same structure as `/result` endpoint

---

### POST /v1/webhooks/pragmatic-play/adjustment

Account adjustment (correction, admin action).

**Request:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "adj_ref_123",
  "gameId": "game_789",
  "roundId": "round_456",
  "amount": 50.00,
  "userId": "operator_slug_user456",
  "hash": "md5_hash_value"
}
```

**Response (Success: error === 0):**
```json
{
  "transactionId": "126",
  "currency": "BRL",
  "cash": 1100.00,
  "bonus": 50.00,
  "error": 0,
  "description": "Success"
}
```

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

### POST /v1/webhooks/digitain-rgs/win

Record win/payout.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "win_info",
      "amount": 100.00
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
      "info": "win_info",
      "balance": 1050.00,
      "transactionId": "txn_127"
    }
  ]
}
```

---

### POST /v1/webhooks/digitain-rgs/betwin

Combined bet and win transaction.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "betwin_info",
      "betAmount": 50.00,
      "winAmount": 150.00
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
      "info": "betwin_info",
      "balance": 1150.00,
      "transactionId": "txn_128"
    }
  ]
}
```

---

### POST /v1/webhooks/digitain-rgs/refund

Refund/cancel transaction.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "refund_info",
      "transactionId": "txn_126",
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
      "info": "refund_info",
      "balance": 1000.00,
      "transactionId": "txn_129"
    }
  ]
}
```

---

### POST /v1/webhooks/digitain-rgs/amend

Amend/modify existing transaction.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "amend_info",
      "transactionId": "txn_126",
      "amount": 75.00
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
      "info": "amend_info",
      "balance": 975.00,
      "transactionId": "txn_129"
    }
  ]
}
```

---

### POST /v1/webhooks/digitain-rgs/checktxstatus

Check transaction status.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "providerTxId": "provider_tx_id",
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value"
}
```

**Response:**
```json
{
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value",
  "txStatus": true,
  "txCreationDate": "2026-05-11T14:30:00Z",
  "externalTxId": "provider_tx_id",
  "currencyId": "USD",
  "errorCode": 1
}
```

---

### POST /v1/webhooks/digitain-rgs/charge

Charge player (debit for multiple items).

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "charge_info",
      "amount": 100.00
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
      "info": "charge_info",
      "balance": 900.00,
      "transactionId": "txn_130"
    }
  ]
}
```

---

### POST /v1/webhooks/digitain-rgs/promowin

Promotional win notification.

**Request:**
```json
{
  "token": "operator_slug_token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "promowin_info",
      "amount": 50.00,
      "campaignId": "campaign_123"
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
      "info": "promowin_info",
      "balance": 950.00,
      "transactionId": "txn_131"
    }
  ]
}
```

---

### POST /v1/webhooks/digitain-rgs/refreshtoken

Refresh/renew player token.

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
  "token": "operator_slug_new_token456",
  "timestamp": "20260511143022",
  "signature": "hmac_sha256_value"
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
| Pragmatic Play | 9 | MD5 Hash | `{slug}_{value}` | ✅ Complete |
| Mancala | 4 | MD5 Hash | `{slug}_{value}` | ✅ Complete |
| Digitain RGS | 11 | HMAC-SHA256 | `{slug}_{value}` | ✅ Complete |
| PG Soft | 4 | Signature | `{slug}_{value}` | ✅ Complete |
| Evoplay | 5 | Signature | `{slug}_{value}` | ✅ Complete |
| Evolution | 5 | HMAC-SHA256 | `{slug}_{value}` | ✅ Complete |
| OpenBox | 4 | HMAC-SHA256 | `{slug}_{value}` | ✅ Complete |
| Alternar | 1 | Signature | `{slug}_{value}` | ✅ Complete |
| **Internal Admin** | 4 | middleware | N/A | ✅ Complete |
| **System** | 3 | None | N/A | ✅ Complete |

**Total Endpoints:** 50 (41 gaming webhooks + 4 internal + 3 system + 2 auth)

**Updated:** All 8 missing Digitain endpoints added (win, betwin, refund, amend, checktxstatus, charge, promowin, refreshtoken)  
**Updated:** Complete Pragmatic Play provider section added with all 9 endpoints

---

**Status:** ✅ Ready for YAML Spec Generation

Next: Separate into individual YAML files per provider + create internal-admin-spec.yaml
