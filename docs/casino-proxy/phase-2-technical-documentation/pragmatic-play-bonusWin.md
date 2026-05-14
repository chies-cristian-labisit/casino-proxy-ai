# Pragmatic Play `/bonusWin` Endpoint — Documentação Técnica

**Endpoint:** `POST /v1/webhooks/pragmatic-play/bonusWin`  
**Provider:** Pragmatic Play  
**Funcionalidade:** Registrar pagamento de bônus ganho pelo jogador  
**Status:** ✅ Documentação Fase 2  

> 🏗️ **Família handleResult() — 2º membro:** `/bonusWin` é o segundo membro da família handleResult(). O wrapper público `bonusWin()` (linhas ~99-102) delega para `handleResult('bonusWin', data)` — a mesma lógica privada compartilhada com `/result`, `/jackpotWin` e `/promoWin`. Referência canônica da família: `pragmatic-play-result.md`.

---

## 1. Resumo Executivo

O endpoint `/bonusWin` registra o pagamento de um bônus ganho pelo jogador. É um evento de negócio distinto do resultado de rodada padrão (`/result`): enquanto `/result` registra o desfecho de uma rodada, `/bonusWin` é disparado especificamente quando o jogador ativa um bônus (ex.: free spins, bônus de acumulador). Tecnicamente, compartilha 100% da implementação PHP com `/result` via `handleResult()`.

**Características:**
- ✅ Usa **apenas `userId`** como identificador
- ✅ **Passthrough** da resposta do provider — sem transformação
- ✅ Requer autenticação via hash MD5
- ✅ Multi-tenant com isolamento de operador
- ✅ **Sem regras exclusivas** — mesmas 9 regras genéricas do `/result`
- 🏗️ **Arquitetura:** thin wrapper `bonusWin()` → `handleResult('bonusWin', data)`

**Fonte PHP:**
- Wrapper: `PragmaticPlayService.php` — método `bonusWin()`, linhas ~99-102
- Lógica: `PragmaticPlayService.php` — método `handleResult()`, linhas ~161-175

---

## 2. Fluxo de Requisição (Request → Response)

```mermaid
graph TD
    A["<b>INPUT</b><br/>POST /v1/webhooks/pragmatic-play/bonusWin<br/>{ userId, amount, currency, bonusId, ... }"]

    A --> B["<b>FASE 1: ROTEAMENTO</b><br/>BR-GENERIC-ROUTING-VALIDATION-001<br/>BR-GENERIC-ERROR-HANDLING-001<br/>method_exists 'bonusWin'<br/>✅ válido"]
    B --> B_err["❌ Endpoint inválido<br/>Exception 500"]
    B --> B2["bonusWin() wrapper (~99-102)<br/>→ delega para handleResult('bonusWin', data)"]

    B2 --> C["<b>FASE 2: EXTRAÇÃO TENANT</b><br/>BR-GENERIC-TENANT-EXTRACTION-001<br/>userId = 'myoperator_user456'<br/>operator_slug = 'myoperator'<br/>📌 handleResult() linha ~161"]

    C --> D["<b>FASE 3: LOOKUP OPERADOR</b><br/>BR-GENERIC-OPERATOR-CACHING-001<br/>SELECT * FROM operators<br/>Cache TTL 1 hora"]
    D --> D_err["❌ Operador não encontrado<br/>ModelNotFoundException"]
    D --> E["<b>FASE 4: SANITIZAÇÃO userId</b><br/>BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001<br/>BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001<br/>userId = removeTenant(userId)<br/>'myoperator_user456' → 'user456'"]

    E --> F["<b>FASE 5: LOOKUP CREDENCIAIS</b><br/>BR-GENERIC-CREDENTIAL-LOOKUP-001<br/>WHERE name='pragmatic'<br/>AND key='secret-key'"]
    F --> F_err["❌ Credencial não encontrada<br/>NullPointerException"]
    F --> G["<b>FASE 6: GERAÇÃO HASH</b><br/>BR-GENERIC-AUTHENTICATION-HMAC-MD5-001<br/>MD5(sorted_payload + secret)"]

    G --> H["<b>FASE 7: HTTP POST</b><br/>BR-GENERIC-PROVIDER-INTEGRATION-001<br/>POST {operator.url}/pragmatic-play/bonusWin.html<br/>Content-Type: application/json"]
    H --> H_err["❌ Provider timeout<br/>Connection failed"]
    H --> I["<b>FASE 8: PASSTHROUGH RESPONSE</b><br/>Resposta retornada inalterada<br/>Sem re-prefixação, sem transformação"]

    I --> J["<b>OUTPUT</b><br/>HTTP 200 OK<br/>{ error: 0, ... } ou { error: N, ... }"]

    style A fill:#e1f5ff
    style B fill:#c8e6c9
    style B2 fill:#fff9c4
    style C fill:#c8e6c9
    style D fill:#c8e6c9
    style E fill:#c8e6c9
    style F fill:#c8e6c9
    style G fill:#c8e6c9
    style H fill:#c8e6c9
    style I fill:#c8e6c9
    style J fill:#e1f5ff
    style B_err fill:#ffcdd2
    style D_err fill:#ffcdd2
    style F_err fill:#ffcdd2
    style H_err fill:#ffcdd2
```

### Explicação das Fases

| Fase | Nome | Regra | Descrição |
|------|------|-------|-----------|
| 1 | Roteamento | BR-GENERIC-ROUTING-VALIDATION-001 + BR-GENERIC-ERROR-HANDLING-001 | `method_exists($service, 'bonusWin')` → válido. Wrapper `bonusWin()` delega imediatamente para `handleResult('bonusWin', $data)`. |
| 2 | Extração Tenant | BR-GENERIC-TENANT-EXTRACTION-001 | Executado **dentro de `handleResult()`** (~linha 161). `userId.split('_')[0]` → `operator_slug`. |
| 3 | Lookup Operador | BR-GENERIC-OPERATOR-CACHING-001 | `OperatorService::get(userId)` com cache Redis TTL 1h. |
| 4 | Sanitização | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 + ORDER-001 | `removeTenant(userId)` remove o prefixo `operator_slug_`. Campo único. |
| 5 | Lookup Credenciais | BR-GENERIC-CREDENTIAL-LOOKUP-001 | `credentials.where('name','pragmatic').where('key','secret-key').first()->value` |
| 6 | Geração Hash | BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | `MD5(ksort(payload) + '&hash=' + secret)` |
| 7 | HTTP POST | BR-GENERIC-PROVIDER-INTEGRATION-001 | `postJson("{operator.url}/pragmatic-play/bonusWin.html", payload)` — URL construída com argumento `'bonusWin'`. |
| 8 | Passthrough | — | Resposta do provider retornada **sem nenhuma modificação**. |

---

## 3. Matriz de Regras Aplicáveis

| # | Regra | Descrição | Fase | Exclusiva? |
|---|-------|-----------|------|------------|
| 1 | **BR-GENERIC-ROUTING-VALIDATION-001** | Dynamic Endpoint Routing | 1 | Não |
| 2 | **BR-GENERIC-ERROR-HANDLING-001** | Unknown endpoint → Exception 500 | 1 (guard) | Não |
| 3 | **BR-GENERIC-TENANT-EXTRACTION-001** | Extrair `operator_slug` do `userId` | 2 | Não |
| 4 | **BR-GENERIC-OPERATOR-CACHING-001** | Operator lookup com cache 1h | 3 | Não |
| 5 | **BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001** | Remover prefixo tenant do `userId` | 4 | Não |
| 6 | **BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001** | Sanitização de `userId` (campo único) | 4 | Não |
| 7 | **BR-GENERIC-CREDENTIAL-LOOKUP-001** | Buscar `secret-key` do operador | 5 | Não |
| 8 | **BR-GENERIC-AUTHENTICATION-HMAC-MD5-001** | Gerar hash MD5 (sort + concat + md5) | 6 | Não |
| 9 | **BR-GENERIC-PROVIDER-INTEGRATION-001** | HTTP POST para `{tenant_url}/pragmatic-play/bonusWin.html` | 7 | Não |

> **Fase 8:** Passthrough direto — sem regra adicional. Resposta do provider retornada inalterada.  
> **Fonte das regras:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

---

## 4. Casos de Erro e Tratamento

### 4.1 `userId` Faltando no Payload

**Entrada:**
```json
{ "amount": 5.00, "currency": "BRL", "bonusId": "bonus_xyz" }
```

**Falha em:** Fase 2 — `handleResult()` linha ~161, `$data['userId']` é null

**Saída:**
```
Exception: Não foi possível encontrar um operator na string {null}
HTTP 500 Internal Server Error
```

---

### 4.2 `userId` sem Underscore (Formato Inválido)

**Entrada:**
```json
{ "userId": "semseparador", "amount": 5.00, "currency": "BRL" }
```

**Falha em:** Fase 2 — parse do `operator_slug` falha

**Saída:**
```
Exception: Não foi possível encontrar um operator na string semseparador
HTTP 500 Internal Server Error
```

---

### 4.3 Operador Não Encontrado

**Entrada:**
```json
{ "userId": "operadorinexistente_user123", "amount": 5.00, "currency": "BRL" }
```

**Falha em:** Fase 3 — `firstOrFail()` lança exceção

**Saída:**
```
Exception: No query results for model [App\Models\Operator]
HTTP 500 Internal Server Error
```

---

### 4.4 Credencial Pragmatic Faltando

**Falha em:** Fase 5 — `credentials->first()` retorna null

**Saída:**
```
Exception: Call to a member function value() on null
HTTP 500 Internal Server Error
```

---

### 4.5 Provider Timeout

**Falha em:** Fase 7 — `postJson()` sem retry (BaseService:19)

**Saída:**
```
Exception: Connection timeout / cURL error
HTTP 500 Internal Server Error
```

---

### 4.6 Bônus Rejeitado pelo Provider (`error != 0`)

**Provider responde:**
```json
{ "error": 2, "description": "Bonus already claimed" }
```

**Comportamento em Fase 8:** Passthrough inalterado

**Saída para o cliente:**
```json
{ "error": 2, "description": "Bonus already claimed" }
```

---

## 5. Exemplo Completo: Request → Response

### 5.1 Caso de Sucesso

**Cliente envia:**
```bash
curl -X POST http://localhost:8080/v1/webhooks/pragmatic-play/bonusWin \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "myoperator_user456",
    "amount": 5.00,
    "currency": "BRL",
    "gameId": "vs20doghouse",
    "roundId": "round_abc789",
    "bonusId": "bonus_freespin_001"
  }'
```

**Processamento interno:**

| Fase | Operação | Input | Output |
|------|----------|-------|--------|
| 1 | Routing + Delegação | endpoint="bonusWin" | `bonusWin()` wrapper → `handleResult('bonusWin', data)` |
| 2 | Tenant Extraction | userId="myoperator_user456" | operator_slug="myoperator" |
| 3 | Operator Lookup | slug="myoperator" | Operador + credentials (cache TTL 1h) |
| 4 | Sanitização | userId="myoperator_user456" | userId="user456" |
| 5 | Credencial | operador.credentials | secret="my_pp_secret_key" |
| 6 | Hash MD5 | sorted payload + secret | hash="d4e5f6a1b2c3..." |
| 7 | HTTP POST | `{url}/pragmatic-play/bonusWin.html` | provider response recebida |
| 8 | **Passthrough** | response do provider | retornada inalterada |

**Provider responde:**
```json
{
  "error": 0,
  "description": "Success",
  "transactionId": "txn_bonus_001",
  "currency": "BRL",
  "cash": 1480.50,
  "bonus": 5.00
}
```

**Casino Proxy retorna (passthrough — inalterado):**
```bash
HTTP 200 OK
Content-Type: application/json

{
  "error": 0,
  "description": "Success",
  "transactionId": "txn_bonus_001",
  "currency": "BRL",
  "cash": 1480.50,
  "bonus": 5.00
}
```

---

## 6. Contexto de Negócio: `/bonusWin` vs `/result`

| Aspecto | `/result` | `/bonusWin` |
|---------|-----------|------------|
| **Evento** | Resultado padrão de uma rodada | Pagamento de bônus ativado pelo jogador |
| **Quando é chamado** | Toda rodada concluída com sucesso | Quando o jogador ativa/resolve um bônus |
| **Exemplos** | Rodada de slot concluída, resultado de aposta | Free spins ganhos, bônus de acumulador |
| **Frequência** | Alta (toda rodada) | Baixa (bônus são eventos especiais) |
| **Impacto no saldo** | Credita ganhos da rodada | Credita valor do bônus |
| **Regras do operador** | Contrato padrão | Podem incluir limites e regras de bônus específicas |
| **Implementação PHP** | `result()` → `handleResult('result')` | `bonusWin()` → `handleResult('bonusWin')` |
| **URL de destino** | `.../result.html` | `.../bonusWin.html` |

> **Nota para implementação Go:** Apesar da implementação PHP idêntica, o handler Go deve tratar `/bonusWin` como rota independente — o provider envia requisições separadas para cada endpoint e a lógica de negócio do operador pode diferir.

---

## 7. Checklist de Segurança

| Validação | Implementada | Regra | Severidade |
|-----------|-------------|-------|------------|
| Tenant isolation (prefixo no userId) | ✅ | BR-GENERIC-TENANT-EXTRACTION-001 | CRÍTICA |
| Sanitização do userId antes de envio ao provider | ✅ | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | CRÍTICA |
| Hash authentication (MD5) | ✅ | BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | CRÍTICA |
| Credencial por operador (secret-key isolado) | ✅ | BR-GENERIC-CREDENTIAL-LOOKUP-001 | CRÍTICA |
| Validação de endpoint (routing guard) | ✅ | BR-GENERIC-ERROR-HANDLING-001 | MÉDIA |
| HTTP method (POST only) | ✅ | routes/api.php | MÉDIA |

---

## 8. Limites e Restrições

| Restrição | Limite / Comportamento | Impacto |
|-----------|----------------------|---------|
| Identificador de entrada | Apenas `userId` (sem `token`) | Clientes devem sempre enviar `userId` |
| Formato do `userId` | Deve conter `_` como delimitador | `userId` sem `_` causa erro 500 |
| Response | Passthrough direto — sem transformação | O Casino Proxy não modifica o resultado do provider |
| Cache de operador | TTL 1 hora | Mudanças no operador levam até 1h para refletir |
| Retry automático | Desabilitado (BaseService:19) | Timeout do provider = falha imediata |
| Hash algorithm | MD5 | Compatibilidade com protocolo Pragmatic Play |

---

## 9. Referências

| Arquivo | Propósito |
|---------|-----------|
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php:99-102` | Wrapper `bonusWin()` |
| `PragmaticPlayService.php:161-175` | `handleResult()` — lógica compartilhada |
| `PragmaticPlayService.php:132-137` | Método `removeTenant()` |
| `PragmaticPlayService.php:142-152` | Método `generateHashCode()` |
| `OperatorService.php:20-34` | Método `get()` (tenant extraction + cache) |
| `BaseService.php:16-22` | Método `postJson()` |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das regras BR-* |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` | Referência canônica da família handleResult() |

---

**Status:** ✅ Documentação Técnica Completa — Pronta para @qa review
