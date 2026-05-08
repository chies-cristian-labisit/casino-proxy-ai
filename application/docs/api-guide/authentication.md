# Casino Proxy API - Authentication Guide

## Overview

The Casino Proxy gateway supports multiple authentication methods depending on the gaming provider. This guide explains signature validation, token management, and security best practices for each integration pattern.

---

## Authentication Methods by Provider

### 1. HMAC-MD5 Signature (Pragmatic Play)

**Used By:** Pragmatic Play

**Signature Algorithm:**
```
signature = MD5(endpoint + "/" + SessionId + [TransactionGuid] + [RefundTransactionGuid] + [RoundGuid] + [Amount] + secret_token)
```

**Validation Steps:**
1. Extract `Hash` field from request body
2. Reconstruct payload: `endpoint/SessionId[+optional_fields]secret_token`
3. Compute MD5 hash of payload
4. Compare computed hash with provided hash (constant-time comparison)

**Go Implementation:**
```go
import "crypto/md5"

func validateMancalaSignature(payload map[string]interface{}, providedHash string, secretKey string) bool {
    endpoint := payload["endpoint"].(string)
    sessionId := payload["SessionId"].(string)
    
    body := endpoint + "/" + sessionId
    if txnGuid, ok := payload["TransactionGuid"]; ok {
        body += txnGuid.(string)
    }
    // ... add other optional fields
    body += secretKey
    
    computed := md5.Sum([]byte(body))
    computedHex := fmt.Sprintf("%x", computed)
    
    return subtle.ConstantTimeCompare([]byte(providedHash), []byte(computedHex)) == 1
}
```

---

### 2. HMAC-SHA256 Signature (Evolution Gaming)

**Used By:** Evolution Gaming

**Signature Algorithm:**
```
signature = HMAC-SHA256(concatenated_fields, secret_key)
Result: Base64-encoded
```

**Request/Response Fields:**
- Signature included in response (not request validation)
- Gateway generates signature for operator callback

**Validation Steps:**
1. Extract relevant fields from request
2. Concatenate in order
3. Compute HMAC-SHA256 with secret key
4. Compare with provided signature

**Go Implementation:**
```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
)

func generateEvolutionSignature(fields []string, secretKey string) string {
    message := strings.Join(fields, "")
    
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(message))
    
    signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
    return signature
}
```

---

### 3. HMAC-SHA256 with Base64URL (Digitain RGS)

**Used By:** Digitain

**Signature Algorithm:**
```
timestamp = now in format YmdHiu
message = timestamp + operatorId
signature = HMAC-SHA256(message, secret_key)
Result: Base64URL-encoded (URL-safe variant)
```

**Base64URL Encoding:**
- Replace `+` with `-`
- Replace `/` with `_`
- Strip padding (`=`)

**Go Implementation:**
```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "strings"
)

func generateDigitainSignature(timestamp string, operatorId string, secretKey string) string {
    message := timestamp + operatorId
    
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(message))
    
    // Standard Base64 encoding
    signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
    
    // Convert to Base64URL
    signature = strings.ReplaceAll(signature, "+", "-")
    signature = strings.ReplaceAll(signature, "/", "_")
    signature = strings.TrimRight(signature, "=")
    
    return signature
}
```

---

### 4. Complex MD5 with Field Traversal (Evoplay)

**Used By:** Evoplay

**Signature Algorithm:**
```
1. Exclude 'signature' field from calculation
2. Join all fields with '*' separator
3. For nested objects/arrays: recursively extract values, join with ':'
4. Append secret_key at end
5. Compute MD5
```

**Example:**
```
MD5("projectId*apiVersion*field1*nested_val1:nested_val2*amount*secretKey")
```

**Go Implementation:**
```go
func generateEvoplaySignature(payload map[string]interface{}, projectId string, secretKey string) string {
    parts := []string{projectId}
    
    // apiVersion would come from credentials
    parts = append(parts, apiVersion)
    
    // Process payload fields (excluding signature)
    for key, value := range payload {
        if key == "signature" {
            continue
        }
        
        if arr, ok := value.([]interface{}); ok {
            var vals []string
            for _, v := range arr {
                vals = append(vals, fmt.Sprintf("%v", v))
            }
            parts = append(parts, strings.Join(vals, ":"))
        } else {
            parts = append(parts, fmt.Sprintf("%v", value))
        }
    }
    
    parts = append(parts, secretKey)
    message := strings.Join(parts, "*")
    
    hash := md5.Sum([]byte(message))
    return fmt.Sprintf("%x", hash)
}
```

---

### 5. HMAC-SHA256 in Headers (OpenBox)

**Used By:** OpenBox

**Signature Format:**
```
Signature Header = "{apiKey}:{base64url(hmac)}"
Timestamp Header = {unix_timestamp}
```

**Signature Algorithm:**
```
1. Sort request fields by key (ksort)
2. Build query string: key1=value1&key2=value2&...
3. Build signature base: query_string + ":" + secretKey + timestamp
4. Compute HMAC-SHA256(signature_base, secretKey)
5. Base64URL encode result
6. Format: apiKey + ":" + encoded_result
```

**Go Implementation:**
```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "sort"
    "strings"
)

func generateOpenBoxSignature(data map[string]interface{}, apiKey string, secretKey string, timestamp int64) string {
    // Sort keys
    var keys []string
    for k := range data {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    
    // Build query string
    var parts []string
    for _, k := range keys {
        v := data[k]
        if b, ok := v.(bool); ok {
            v = map[bool]string{true: "true", false: "false"}[b]
        }
        parts = append(parts, fmt.Sprintf("%s=%v", k, v))
    }
    queryString := strings.Join(parts, "&")
    
    // Build signature base
    signatureBase := queryString + ":" + secretKey + fmt.Sprintf("%d", timestamp)
    
    // HMAC-SHA256
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(signatureBase))
    hmacResult := h.Sum(nil)
    
    // Base64URL encode
    signature := base64.URLEncoding.EncodeToString(hmacResult)
    signature = strings.TrimRight(signature, "=")
    
    return apiKey + ":" + signature
}
```

---

### 6. No Signature Validation (PG Soft)

**Used By:** PG Soft

**Security Model:**
- No cryptographic signature
- Authentication via operator credentials (token + secret_key)
- Validation done at transport layer (HTTPS/TLS)

**Credential Fields:**
- `operator_token`: Static operator identifier
- `secret_key`: Static secret (for future use)

**Validation Steps:**
1. Extract `operator_token` and `secret_key` from request
2. Look up operator in database
3. Verify credentials match
4. Validate operator owns the session token

**Go Implementation:**
```go
func validatePGSoftCredentials(token string, secretKey string) (*Operator, error) {
    operator := &Operator{}
    
    // Look up by token and secret_key
    err := db.Where("operator_token = ? AND secret_key = ?", token, secretKey).
        First(operator).Error
    
    if err != nil {
        return nil, fmt.Errorf("invalid credentials")
    }
    
    return operator, nil
}
```

---

### 7. HTTP Proxy (Alternar)

**Used By:** Alternar

**Security Model:**
- No gateway-level signature validation
- HTTPS forwarding to external endpoint
- Upstream service handles authentication

**No validation implementation required** - gateway acts as transparent HTTP proxy.

---

## Security Best Practices

### 1. Constant-Time Comparison

**Always use constant-time comparison for signatures to prevent timing attacks:**

```go
import "crypto/subtle"

// ✓ Correct
if subtle.ConstantTimeCompare([]byte(provided), []byte(computed)) != 1 {
    return errors.New("invalid signature")
}

// ✗ Wrong - vulnerable to timing attacks
if provided != computed {
    return errors.New("invalid signature")
}
```

### 2. Secure Credential Storage

**Credentials should be:**
- Encrypted at rest
- Never logged or printed
- Rotated regularly
- Stored in separate database or vault

```go
// Store encrypted
encrypted := encrypt(secretKey, masterKey)
db.Update("operators", map[string]interface{}{
    "secret_key": encrypted,
})
```

### 3. HTTPS/TLS Only

**All webhook endpoints must use HTTPS:**
- Certificate pinning recommended
- TLS 1.2 minimum
- Strong cipher suites only

### 4. Request Validation Order

**Always validate in this order:**
1. Check required fields present
2. Validate signature/authentication
3. Extract operator context
4. Validate operator permissions
5. Process payload

```go
func validateWebhookRequest(req *http.Request) (*Operator, error) {
    // 1. Parse and validate structure
    payload := parsePayload(req)
    if err := payload.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Validate signature
    if !validateSignature(payload, secretKey) {
        return nil, errors.New("invalid signature")
    }
    
    // 3. Extract and validate operator
    operator, err := getOperator(payload)
    if err != nil {
        return nil, err
    }
    
    // 4. Additional permission checks
    if !operator.IsEnabled {
        return nil, errors.New("operator disabled")
    }
    
    return operator, nil
}
```

---

## Token Sanitization

All providers use the operator context prefix format: `{operator_slug}_{actual_value}`

**Sanitization Pattern:**
```go
func extractOperatorContext(token string) (slug string, actualValue string, err error) {
    parts := strings.Split(token, "_")
    if len(parts) != 2 {
        return "", "", errors.New("invalid token format")
    }
    return parts[0], parts[1], nil
}

func reconstructToken(operator *Operator, actualValue string) string {
    return operator.Slug + "_" + actualValue
}
```

This ensures:
- Operator isolation (different operators can't access each other's players)
- Clean token format for callbacks
- Prevents token leakage across operators

---

## Testing Authentication

**Unit Tests for Signature Validation:**

```go
func TestMancalaSignatureValidation(t *testing.T) {
    payload := map[string]interface{}{
        "SessionId": "my-op_abc123",
        "Amount": 50.00,
    }
    secretKey := "test_secret"
    endpoint := "Debit"
    
    // Compute expected signature
    body := endpoint + "/" + payload["SessionId"].(string) + payload["Amount"].(string) + secretKey
    expectedHash := md5.Sum([]byte(body))
    expectedHex := fmt.Sprintf("%x", expectedHash)
    
    // Test validation
    valid := validateMancalaSignature(payload, expectedHex, secretKey)
    assert.True(t, valid)
    
    // Test invalid signature
    invalid := validateMancalaSignature(payload, "invalid", secretKey)
    assert.False(t, invalid)
}
```

---

## Credential Management

**Rotating Provider Credentials:**

1. Generate new secret key
2. Store both old and new (with expiration)
3. Validate using both during transition period
4. Remove old credential after grace period
5. Update provider configuration

```go
type OperatorCredential struct {
    ID        int
    OperatorID int
    Provider  string
    Key       string
    Value     string // encrypted
    IsActive  bool
    ExpiresAt *time.Time
    CreatedAt time.Time
}

// Accept both old and new during rotation
func validateWithMultipleKeys(signature string, newKey string, oldKey string) bool {
    return validateSignature(signature, newKey) || validateSignature(signature, oldKey)
}
```

---

## Summary Table

| Provider | Auth Method | Location | Algorithm | Status |
|----------|-------------|----------|-----------|--------|
| Pragmatic Play | MD5 Signature | Body | MD5 | Active |
| Evolution Gaming | HMAC-SHA256 | Response | SHA256 | Active |
| PG Soft | Credentials | Body | Token+Secret | Active |
| Mancala | MD5 Signature | Body | MD5 | Active |
| Digitain | HMAC-SHA256 | Response | SHA256+Base64URL | Active |
| Evoplay | MD5 Signature | Body | MD5+Recursive | Active |
| OpenBox | HMAC-SHA256 | Headers | SHA256+Base64URL | Active |
| Alternar | HTTP Proxy | N/A | HTTPS only | Active |

---

## Troubleshooting

### Signature Validation Fails

**Check:**
1. Is secret key correct? (Compare with provider config)
2. Are optional fields included in the signature? (Check payload structure)
3. Is the signature algorithm correct? (Check provider documentation)
4. Is timestamp within acceptable range? (For time-based signatures)

### Token Extraction Issues

**Check:**
1. Does token follow `{operator_slug}_{value}` format?
2. Is operator slug in system? (Check operator table)
3. Is operator enabled? (Check operator.is_active flag)

### Recurring Signature Failures

**Solutions:**
1. Verify provider credentials are up to date
2. Check for clock skew (time synchronization between systems)
3. Review provider signature calculation documentation
4. Enable debug logging for signature calculation
5. Contact provider technical support
