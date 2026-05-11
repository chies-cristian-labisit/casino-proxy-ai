# Complete Endpoints Inventory - Casino Proxy

**Data:** 2026-05-11  
**Status:** Draft - Awaiting Review  
**Source:** routes/api.php + all Controllers  
**Scope:** ALL endpoints (public, private, webhooks, internal)

---

## 📊 Summary

| Categoria | Público | Privado | Total |
|-----------|---------|---------|-------|
| **Player/Session** | 1 | 0 | 1 |
| **Gaming Provider Webhooks** | 8 | 0 | 8 |
| **Operator Management** | 0 | 3 | 3 |
| **Credentials Management** | 0 | 1 | 1 |
| **System Health** | 1 | 0 | 1 |
| **Media/Images** | 1 | 0 | 1 |
| **Authentication (Dashboard)** | 2 | 2 | 4 |
| **TOTAL** | **13** | **6** | **19** |

---

# 🔓 PUBLIC ENDPOINTS (Exposed APIs)

## 1️⃣ Player/Session Entry

### POST /v1/entry
**Method:** POST  
**Controller:** OperatorController@entry  
**Route File:** routes/api.php:17  
**Description:** Player session entry point (frontend entry)  
**Status:** Public  
**Auth:** None (external)  

**Payload:**
```json
{
  "operator_id": "string",
  "player_id": "string",
  "session_token": "string"
}
```

**Response:** (varies by provider)

---

## 2️⃣ System Health Check

### GET /v1/status
**Method:** GET  
**Route File:** routes/api.php:23  
**Description:** Health check endpoint  
**Status:** Public  
**Auth:** None  

**Response:**
```json
{
  "message": "OK"
}
```

---

## 3️⃣ Media Upload

### POST /v1/image
**Method:** POST  
**Controller:** ImageController@index  
**Route File:** routes/api.php:25  
**Description:** Image upload (operator branding, assets)  
**Status:** Public  
**Auth:** None specified (assumes operator context in body)  

**Payload:** Multipart form-data with image file  
**Response:** Image URL or reference  

---

## 4️⃣ Gaming Provider Webhooks

All provider webhooks are **PUBLIC** (called by external gaming providers).

### A. Pragmatic Play

#### POST /v1/webhooks/pragmatic-play/{endpoint}
**Method:** POST  
**Controller:** PragmaticPlayController@__invoke  
**Route File:** routes/api.php:29  
**Description:** Generic webhook gateway for Pragmatic Play  
**Status:** Public (called by Pragmatic Play)  
**Auth:** HMAC-MD5 hash validation  

**Endpoints (handled dynamically):**
- `/authenticate.html` → calls `PragmaticPlayService->call('authenticate', data)`
- `/balance.html` → calls `PragmaticPlayService->call('balance', data)`
- `/bet.html` → calls `PragmaticPlayService->call('bet', data)`
- `/refund.html` → calls `PragmaticPlayService->call('refund', data)`
- `/result.html` → calls `PragmaticPlayService->call('result', data)`
- `/bonusWin.html` → calls `PragmaticPlayService->call('bonusWin', data)`
- `/jackpotWin.html` → calls `PragmaticPlayService->call('jackpotWin', data)`
- `/promoWin.html` → calls `PragmaticPlayService->call('promoWin', data)`
- `/adjustment.html` → calls `PragmaticPlayService->call('adjustment', data)`

**Internal Outbound Calls (made BY Casino Proxy):**
```
POST {$tenant['url']}/pragmatic-play/authenticate.html
POST {$tenant['url']}/pragmatic-play/balance.html
POST {$tenant['url']}/pragmatic-play/bet.html
POST {$tenant['url']}/pragmatic-play/refund.html
POST {$tenant['url']}/pragmatic-play/result.html
POST {$tenant['url']}/pragmatic-play/bonusWin.html
POST {$tenant['url']}/pragmatic-play/jackpotWin.html
POST {$tenant['url']}/pragmatic-play/promoWin.html
POST {$tenant['url']}/pragmatic-play/adjustment.html
```

**Documentation:** ✅ Swagger (pragmatic-play-spec.yaml)  
**Note:** ⚠️ Outbound endpoints NOT in Swagger (internal integration)

---

### B. Mancala

#### POST /v1/webhooks/mancala/{endpoint}
**Method:** POST  
**Controller:** MancalaController@__invoke  
**Route File:** routes/api.php:30  
**Description:** Generic webhook gateway for Mancala  
**Status:** Public  
**Auth:** Provider signature validation  

**Endpoints (dynamic):**
- Generic routing: `/{endpoint}` → calls `MancalaService->call(endpoint, data)`

**Documentation:** ❌ NOT in Swagger  
**Note:** Dynamic endpoint handling (similar to Pragmatic Play)

---

### C. Digitain RGS

#### POST /v1/webhooks/digitain-rgs/{endpoint}
**Method:** POST  
**Controller:** DigitainController@__invoke  
**Route File:** routes/api.php:31  
**Description:** Generic webhook gateway for Digitain RGS  
**Status:** Public  
**Auth:** Provider signature validation  

**Endpoints (dynamic):**
- Generic routing: `/{endpoint}` → calls `DigitainService->call(endpoint, data)`

**Documentation:** ❌ NOT in Swagger  
**Note:** Dynamic endpoint handling

---

### D. PG Soft

#### POST /v1/webhooks/pgsoft/VerifySession
**Method:** POST  
**Controller:** PgSoftController@verifySession  
**Route File:** routes/api.php:32-35  
**Description:** PG Soft session verification  
**Status:** Public  
**Auth:** Provider signature validation  

**Response:** Player info + balance

---

#### POST /v1/webhooks/pgsoft/Cash/Get
**Method:** POST  
**Controller:** PgSoftController@getCash  
**Route File:** routes/api.php:33  
**Description:** Get player cash balance  
**Status:** Public  

---

#### POST /v1/webhooks/pgsoft/Cash/TransferInOut
**Method:** POST  
**Controller:** PgSoftController@transferInOut  
**Route File:** routes/api.php:34  
**Description:** Transfer credits in/out (bets, wins)  
**Status:** Public  

---

#### POST /v1/webhooks/pgsoft/Cash/Adjustment
**Method:** POST  
**Controller:** PgSoftController@adjustment  
**Route File:** routes/api.php:35  
**Description:** Adjustment/correction of player account  
**Status:** Public  

**Documentation:** ❌ NOT in Swagger (separate spec needed)

---

### E. Evoplay

#### POST /v1/webhooks/evoplay
**Method:** POST  
**Controller:** EvoplayController@__invoke  
**Route File:** routes/api.php:36  
**Description:** Evoplay webhook gateway  
**Status:** Public  
**Auth:** Provider signature validation  

**Endpoints (dynamic):**
- Single endpoint: calls `EvoplayService->__invoke(request)`

**Documentation:** ❌ NOT in Swagger

---

### F. Evolution Gaming

#### POST /v1/webhooks/evolution/authentication
**Method:** POST  
**Controller:** EvolutionController@authentication  
**Route File:** routes/api.php:37-41  
**Description:** Player authentication  
**Status:** Public  

---

#### POST /v1/webhooks/evolution/debit
**Method:** POST  
**Controller:** EvolutionController@debit  
**Description:** Debit player account (bet placement)  
**Status:** Public  

---

#### POST /v1/webhooks/evolution/credit
**Method:** POST  
**Controller:** EvolutionController@credit  
**Description:** Credit player account (win/payout)  
**Status:** Public  

---

#### POST /v1/webhooks/evolution/rollback
**Method:** POST  
**Controller:** EvolutionController@rollback  
**Description:** Rollback transaction  
**Status:** Public  

---

#### POST /v1/webhooks/evolution/getNewToken
**Method:** POST  
**Controller:** EvolutionController@getNewToken  
**Description:** Request new session token  
**Status:** Public  

**Documentation:** ❌ NOT in Swagger (separate spec needed)

---

### G. OpenBox

#### PUT /v1/webhooks/openbox/seamless/balance
**Method:** PUT  
**Controller:** OpenBoxController@balance  
**Route File:** routes/api.php:42-45  
**Description:** Get player balance  
**Status:** Public  
**HTTP Method:** PUT (non-standard)  

---

#### PUT /v1/webhooks/openbox/seamless/bet
**Method:** PUT  
**Controller:** OpenBoxController@bet  
**Description:** Place bet  
**Status:** Public  

---

#### PUT /v1/webhooks/openbox/seamless/win
**Method:** PUT  
**Controller:** OpenBoxController@win  
**Description:** Record win/payout  
**Status:** Public  

---

#### PUT /v1/webhooks/openbox/seamless/refund
**Method:** PUT  
**Controller:** OpenBoxController@refund  
**Description:** Refund transaction  
**Status:** Public  

**Documentation:** ❌ NOT in Swagger (separate spec needed)

---

### H. Alternar

#### POST /v1/webhooks/alternar/
**Method:** POST  
**Controller:** AlternarRedirectController@index  
**Route File:** routes/api.php:46-47  
**Description:** Alternar redirect handler  
**Status:** Public  

**Documentation:** ❌ NOT in Swagger

---

# 🔒 PRIVATE/INTERNAL ENDPOINTS (Admin/Dashboard)

All private endpoints require **middleware: ['internalCheck']**  
Reference: routes/api.php:19

---

## 1️⃣ Operator Management (Internal)

### GET /v1/internal/operator/
**Method:** GET  
**Controller:** OperatorController@list  
**Route File:** routes/api.php:20  
**Description:** List all operators  
**Status:** Private (internal only)  
**Auth:** internalCheck middleware  
**Request Class:** ListOperatorRequest  

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

### POST /v1/internal/operator/store
**Method:** POST  
**Controller:** OperatorController@store  
**Route File:** routes/api.php:21  
**Description:** Create new operator  
**Status:** Private  
**Auth:** internalCheck middleware  
**Request Class:** StoreOperatorRequest  

**Payload:**
```json
{
  "slug": "new_operator",
  "name": "New Operator Name",
  "url": "https://api.newoperator.com"
}
```

---

### DELETE /v1/internal/operator/delete/{slug}
**Method:** DELETE  
**Controller:** OperatorController@delete  
**Route File:** routes/api.php:22  
**Description:** Delete operator  
**Status:** Private  
**Auth:** internalCheck middleware  
**Path Parameter:** operator:slug (Model binding)

---

## 2️⃣ Credential Management (Internal)

### GET /v1/internal/credentials/
**Method:** GET  
**Controller:** CredentialController@index  
**Route File:** (NOT FOUND - likely missing or not wired)  
**Description:** List credentials (likely)  
**Status:** Private  
**Auth:** internalCheck middleware  

**Note:** ⚠️ Route NOT found in routes/api.php - possibly missing or in different location

---

## 3️⃣ Authentication / Dashboard (Mixed)

### GET /auth/login
**Method:** GET  
**Controller:** AuthController@showLogin  
**Route File:** routes/web.php (NOT api.php)  
**Description:** Show login form (web page)  
**Status:** Public (frontend)  

---

### POST /auth/login
**Method:** POST  
**Controller:** AuthController@login  
**Route File:** routes/web.php  
**Description:** Process login  
**Status:** Public (frontend)  

---

### POST /auth/logout
**Method:** POST  
**Controller:** AuthController@logout  
**Route File:** routes/web.php  
**Description:** Process logout  
**Status:** Public (frontend)  

---

### GET /auth
**Method:** GET  
**Controller:** AuthController@index  
**Route File:** routes/web.php  
**Description:** Dashboard (likely)  
**Status:** Private (requires auth)  

---

# 🔄 OUTBOUND/INTERNAL CALLS (Not API Endpoints)

These are calls made BY Casino Proxy TO external systems. NOT endpoints that receive requests.

## Pragmatic Play Internal APIs

Constructed in: PragmaticPlayService.php:37, 59, 74, 89, 124, 165

```
POST {$tenant['url']}/pragmatic-play/authenticate.html
POST {$tenant['url']}/pragmatic-play/balance.html
POST {$tenant['url']}/pragmatic-play/bet.html
POST {$tenant['url']}/pragmatic-play/refund.html
POST {$tenant['url']}/pragmatic-play/result.html
POST {$tenant['url']}/pragmatic-play/bonusWin.html
POST {$tenant['url']}/pragmatic-play/jackpotWin.html
POST {$tenant['url']}/pragmatic-play/promoWin.html
POST {$tenant['url']}/pragmatic-play/adjustment.html
```

**Note:** Each service (Mancala, Digitain, etc) similarly makes outbound calls. Need to extract from each service.

---

# 📋 Documentation Status

## Fully Documented (OpenAPI)
- ✅ Pragmatic Play webhooks (pragmatic-play-spec.yaml)

## NOT Documented (OpenAPI)
- ❌ Mancala webhooks
- ❌ Digitain RGS webhooks
- ❌ PG Soft webhooks
- ❌ Evoplay webhooks
- ❌ Evolution Gaming webhooks
- ❌ OpenBox webhooks
- ❌ Alternar webhooks
- ❌ Player entry (/v1/entry)
- ❌ Image upload (/v1/image)
- ❌ Operator management (internal)
- ❌ Credentials management (internal)
- ❌ Authentication/Dashboard (web routes)

**Total Undocumented:** 13 endpoints + 6 private endpoints + outbound calls

---

# 🎯 Next Steps

1. **Create OpenAPI specs for each provider:**
   - mancala-spec.yaml
   - digitain-rgs-spec.yaml
   - pgsoft-spec.yaml
   - evoplay-spec.yaml
   - evolution-spec.yaml
   - openbox-spec.yaml
   - alternar-spec.yaml

2. **Create internal API spec:**
   - internal-admin-spec.yaml (operator, credentials, etc)

3. **Document outbound calls:**
   - provider-integration-internal.md (what Casino Proxy calls)

4. **Validate authentication methods:**
   - Which use HMAC, which use API keys, etc

5. **Map error codes per provider**

---

**Status:** ⏸️ Awaiting Review
