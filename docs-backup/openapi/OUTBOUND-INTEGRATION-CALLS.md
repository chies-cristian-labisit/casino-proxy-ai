# Outbound Integration Calls - What Casino Proxy Calls

**Date:** 2026-05-11  
**Purpose:** Document all calls made BY Casino Proxy TO external provider APIs  
**Type:** Internal Integration Documentation (NOT Public API Endpoints)

---

## Overview

These are NOT endpoints that receive requests. These are endpoints that the Casino Proxy **makes requests to** when processing provider webhooks.

**Pattern:** When Casino Proxy receives webhook from Provider A, it calls Provider A's internal API to perform operations (balance check, transaction, etc).

---

# PRAGMATIC PLAY OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/pragmatic-play/`  
**Method:** POST  
**Content-Type:** application/json  
**Source:** PragmaticPlayService.php:37, 59, 74, 89, 124, 165

## POST /pragmatic-play/authenticate.html

**Trigger:** When Casino Proxy receives authenticate webhook  
**Purpose:** Authenticate player session with Pragmatic Play backend  

**Request Body (sent to Pragmatic Play):**
```json
{
  "providerId": "PragmaticPlay",
  "token": "abc123def456",  // Sanitized (tenant prefix removed)
  "hash": "md5_hash_value"
}
```

**Expected Response (from Pragmatic Play):**
```json
{
  "userId": "123456",
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

**Casino Proxy Processing:**
1. Sanitize token (remove tenant prefix)
2. Generate MD5 hash
3. POST to Pragmatic Play
4. If error === 0: Prefix userId back with tenant_slug
5. Return to client

---

## POST /pragmatic-play/balance.html

**Trigger:** When Casino Proxy receives balance webhook  

**Request Body:**
```json
{
  "providerId": "PragmaticPlay",
  "token": "abc123def456",  // Optional
  "userId": "user456",       // Optional (can use either)
  "hash": "md5_hash_value"
}
```

**Expected Response:**
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

## POST /pragmatic-play/bet.html

**Trigger:** Bet placement webhook  

**Request Body:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "bet_ref_123",
  "gameId": "game_789",
  "roundId": "round_456",
  "roundDetails": "game_details",
  "amount": 50.00,
  "timestamp": "2026-05-11T14:30:00Z",
  "userId": "user456",
  "bonusCode": "BONUS2024",  // Optional
  "hash": "md5_hash_value"
}
```

**Expected Response:**
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

## POST /pragmatic-play/refund.html

**Trigger:** Bet refund/cancellation  

**Request Body:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "bet_ref_123",
  "userId": "user456",
  "hash": "md5_hash_value"
}
```

**Expected Response:**
```json
{
  "transactionId": "124",
  "error": 0,
  "description": "Success"
}
```

---

## POST /pragmatic-play/result.html

**Trigger:** Bet result notification  

**Request Body:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "result_ref_123",
  "gameId": "game_789",
  "roundId": "round_456",
  "amount": 100.00,
  "timestamp": "2026-05-11T14:35:00Z",
  "userId": "user456",
  "roundDetails": "game_details",  // Optional
  "promoWinAmount": 10.00,          // Optional
  "promoWinReference": "promo_ref", // Optional
  "promoCampaignID": "campaign_123", // Optional
  "promoCampaignType": "type",      // Optional
  "bonusCode": "BONUS2024",         // Optional
  "hash": "md5_hash_value"
}
```

**Expected Response:**
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

## POST /pragmatic-play/bonusWin.html

**Trigger:** Bonus win notification  
**Structure:** Same as result.html

---

## POST /pragmatic-play/jackpotWin.html

**Trigger:** Jackpot win notification  
**Structure:** Same as result.html

---

## POST /pragmatic-play/promoWin.html

**Trigger:** Promotional win notification  
**Structure:** Same as result.html

---

## POST /pragmatic-play/adjustment.html

**Trigger:** Account adjustment (correction, admin action)  

**Request Body:**
```json
{
  "providerId": "PragmaticPlay",
  "reference": "adj_ref_123",
  "gameId": "game_789",
  "roundId": "round_456",
  "amount": 50.00,
  "userId": "user456",
  "hash": "md5_hash_value"
}
```

**Expected Response:**
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

# MANCALA OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/mancala/`  
**Allowed Endpoints:** Balance, Credit, Debit, Refund  

### POST /mancala/Balance

**Request:**
```json
{
  "SessionId": "session123",  // Sanitized
  "Hash": "md5_hash"
}
```

---

### POST /mancala/Credit

**Request:**
```json
{
  "SessionId": "session123",
  "TransactionGuid": "txn_guid",
  "Amount": 100.50,
  "RoundGuid": "round_guid",
  "Hash": "md5_hash"
}
```

---

### POST /mancala/Debit

**Request:**
```json
{
  "SessionId": "session123",
  "TransactionGuid": "txn_guid",
  "Amount": 50.00,
  "RoundGuid": "round_guid",
  "Hash": "md5_hash"
}
```

---

### POST /mancala/Refund

**Request:**
```json
{
  "SessionId": "session123",
  "RefundTransactionGuid": "txn_guid",
  "Amount": 50.00,
  "Hash": "md5_hash"
}
```

---

# DIGITAIN RGS OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/digitain-rgs/`  
**Allowed Endpoints:** authenticate, getbalance, bet, win, betwin, refund, amend, checktxstatus, charge, promowin, refreshtoken

### POST /digitain-rgs/authenticate

**Request:**
```json
{
  "token": "token123",  // Sanitized
  "operatorId": "op_123",
  "timestamp": "20260511143022",
  "signature": "hmac_sha256"
}
```

---

### POST /digitain-rgs/getbalance

**Request:**
```json
{
  "token": "token123",
  "operatorId": "op_123",
  "timestamp": "20260511143022",
  "signature": "hmac_sha256"
}
```

---

### POST /digitain-rgs/bet

**Request:**
```json
{
  "token": "token123",
  "operatorId": "op_123",
  "items": [
    {
      "info": "bet_info",
      "amount": 50.00
    }
  ],
  "timestamp": "20260511143022",
  "signature": "hmac_sha256"
}
```

**Note:** May batch multiple items into single request

---

# PG SOFT OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/pgsoft/`

### POST /pgsoft/VerifySession

**Request:**
```json
{
  "operator_player_session": "session123",  // Sanitized
  "playerId": "player_456"
}
```

---

### POST /pgsoft/Cash/Get

**Request:**
```json
{
  "operator_player_session": "session123",
  "playerId": "player_456"
}
```

---

### POST /pgsoft/Cash/TransferInOut

**Request:**
```json
{
  "operator_player_session": "session123",
  "transactionId": "txn_127",
  "amount": 50.00,
  "transactionType": "bet|win",
  "gameId": "game_123",
  "roundId": "round_123"
}
```

---

### POST /pgsoft/Cash/Adjustment

**Request:**
```json
{
  "operator_player_session": "session123",
  "amount": 50.00,
  "reason": "refund|correction",
  "referenceId": "ref_123"
}
```

---

# EVOPLAY OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/evoplay`  
**Method:** POST  
**Single Endpoint:** `/evoplay` (no method in path, method specified in request body)

### POST /evoplay

**Request:**
```json
{
  "name": "init|bet|win|refund|BalanceIncrease",
  "token": "token123",  // Sanitized
  "data": {
    "user_id": "user_456",
    "game_id": "game_789",
    "amount": 100.50,
    "transaction_id": "txn_128"
  },
  "signature": "signature_value"
}
```

---

# EVOLUTION GAMING OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/evolution/`  
**Authentication:** HMAC-SHA256 in `hash` header

### POST /evolution/authentication

**Request:**
```json
{
  "token": "token123",  // Sanitized
  "gameId": "game_789",
  "sessionId": "session_456"
}
```

**Headers:**
```
hash: base64_encoded_hmac_sha256
```

---

### POST /evolution/debit

**Request:**
```json
{
  "token": "token123",
  "transactionId": "txn_129",
  "amount": 50.00,
  "gameId": "game_789"
}
```

---

### POST /evolution/credit

**Request:**
```json
{
  "token": "token123",
  "transactionId": "txn_130",
  "amount": 100.00,
  "gameId": "game_789"
}
```

---

# OPENBOX OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/openbox/seamless/`  
**Method:** PUT (non-standard)  
**Authentication:** HMAC-SHA256 signature in `signature` header

### PUT /openbox/seamless/balance

**Request:**
```json
{
  "playerId": "player_456",
  "timestamp": "1715424600"
}
```

---

### PUT /openbox/seamless/bet

**Request:**
```json
{
  "playerId": "player_456",
  "amount": 50.00,
  "gameId": "game_789",
  "roundId": "round_456",
  "timestamp": "1715424600"
}
```

---

### PUT /openbox/seamless/win

**Request:**
```json
{
  "playerId": "player_456",
  "amount": 100.00,
  "gameId": "game_789",
  "roundId": "round_456",
  "timestamp": "1715424600"
}
```

---

### PUT /openbox/seamless/refund

**Request:**
```json
{
  "playerId": "player_456",
  "originalTransactionId": "txn_131",
  "amount": 50.00,
  "timestamp": "1715424600"
}
```

---

# ALTERNAR OUTBOUND CALLS

**Base URL:** `{$tenant['url']}/alternar/`

### POST /alternar/

**Request:**
```json
{
  "playerId": "player_456",
  "token": "token123",
  "action": "launch|logout",
  "gameId": "game_789",
  "signature": "hash_value"
}
```

---

# Summary

| Provider | Outbound Endpoints | Auth |
|----------|-------------------|------|
| Pragmatic Play | 9 | MD5 Hash |
| Mancala | 4 | MD5 Hash |
| Digitain RGS | 11 | HMAC-SHA256 |
| PG Soft | 4 | Provider Signature |
| Evoplay | 1 | Signature |
| Evolution | 5+ | HMAC-SHA256 header |
| OpenBox | 4 | HMAC-SHA256 signature |
| Alternar | 1 | Signature |
| **Total** | **39** | Various |

---

**Status:** ⏸️ Ready for Validation
