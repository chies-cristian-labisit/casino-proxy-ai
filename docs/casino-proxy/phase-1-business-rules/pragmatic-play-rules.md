# Pragmatic Play — Regras de Lógica de Negócio Extraídas

**Status:** Extração Completa  
**Data de Extração:** 2026-05-11  
**Fonte:** `legacy/casino-proxy/app/Services/PragmaticPlayService.php` + testes em `tests/Feature/PragmaticPlayControllerTest.php`  
**Total de Regras:** 12  

---

## Resumo Executivo

O serviço Pragmatic Play implementa um **proxy de webhook para casino online**, interceptando requisições da Pragmatic Play para:
1. Extrair o operador/tenant da requisição
2. Sanitizar tokens (remover prefixo de tenant)
3. Gerar assinatura MD5 de autenticação
4. Encaminhar para API do provider
5. Re-prefixar respostas com tenant

Todas as 9 endpoints compartilham a mesma lógica base, com variações mínimas.

---

## Regras Extraídas

### Regra PP-001: Dynamic Endpoint Routing via Method Resolution

**ID:** PP-001  
**Descrição:** O serviço roteia requisições para métodos específicos usando dynamic method calling baseado no parâmetro `endpoint` da requisição.

**Contexto de Negócio:**  
Permite que um único controller suporte múltiplos endpoints sem duplicação de código. O parâmetro `{endpoint}` vindo da URL é convertido em chamada de método PHP.

**Escopo:** Todos os endpoints  

**Lógica de Decisão:**
```
SE request.endpoint existe como método na classe
  ENTÃO chamar $this->{método}($data)
  SENÃO lançar exceção 500 "Endpoint {endpoint} was not found"
```

**Casos Extremos:**
- Endpoint inválido/desconhecido → Exceção com status 500
- Método não existe na classe → Exceção capturada em log

**Código Fonte:**  
`PragmaticPlayService.php:16-25` (método `call()`)

**Dependências:** Nenhuma (primeiro check no handler)

**Validação Automática:**  
Teste: `test_unknown_method_to_pragmatic_play_endpoint()` — verifica que endpoint desconhecido retorna 500

---

### Regra PP-002: Tenant/Operator Extraction from Token

**ID:** PP-002  
**Descrição:** O `token` (ou `userId`) é parseado usando delimiter `_` para extrair o slug do operador (tenant).

**Contexto de Negócio:**  
Implementa multi-tenancy: requisições vêm com token prefixado com slug do operador (formato: `{operator_slug}_{actual_token}`). Esse prefixo identifica qual operador está fazendo a requisição.

**Escopo:** Todos os endpoints (via `authenticate`, `balance`, `bet`, `refund`, `result`, `bonusWin`, `jackpotWin`, `promoWin`, `adjustment`)

**Lógica de Decisão:**
```
token_parts = token.split('_')
operator_slug = token_parts[0:-1].join('_')  # tudo menos última parte
actual_token = token_parts[-1]  # última parte
```

**Exemplo:**
- Input: `"myoperator_xyz123abc"`
- Output: `operator_slug = "myoperator"`, `actual_token = "xyz123abc"`

**Casos Extremos:**
- Token sem underscore → Thrown exception "Não foi possível encontrar um operator na string {operator}"
- Token com múltiplos underscores (ex: `"my_operator_token"`) → `operator_slug = "my_operator"`, `actual_token = "token"`

**Código Fonte:**  
`OperatorService.php:20-28` (método `get()`)  
`PragmaticPlayService.php:132-137` (método `removeTenant()`)

**Dependências:** Nenhuma (primeiro check no handler)

---

### Regra PP-003: Tenant Lookup with 1-Hour Caching

**ID:** PP-003  
**Descrição:** O operador é recuperado do banco usando o `operator_slug` e cacheado por 1 hora em memória.

**Contexto de Negócio:**  
Reduz queries ao banco de dados. Operadores são entidades estáticas durante a sessão, então cache é seguro.

**Escopo:** Todos os endpoints

**Lógica de Decisão:**
```
cache_key = 'operator_' + operator_slug
SE cache[cache_key] existe
  ENTÃO retornar cached value
  SENÃO 
    operador = SELECT * FROM operators WHERE slug = operator_slug WITH credentials
    cache[cache_key] = operador (com TTL 60 minutos)
    retornar operador
```

**Casos Extremos:**
- Operador não existe no banco → `firstOrFail()` lança `ModelNotFoundException`
- Cache expirado → Refetch do banco
- Múltiplos underscores no slug → Handled correctly (exemplo: `operator_slug = "my_op_co"` é valid)

**Código Fonte:**  
`OperatorService.php:30-34` (método `get()` — cache logic)

**Dependências:** PP-002 (deve extrair operator_slug primeiro)

**Performance Note:** Cache TTL = 3600 segundos (1 hora)

---

### Regra PP-004: Tenant Prefix Removal (Sanitization)

**ID:** PP-004  
**Descrição:** Antes de enviar requisição ao provider, todos os tokens no payload são sanitizados removendo o prefixo de tenant.

**Contexto de Negócio:**  
O provider (Pragmatic Play) não conhece sobre nosso sistema de multi-tenancy. Requisições devem conter apenas o token real, não o token prefixado.

**Escopo:** Todos os endpoints

**Lógica de Decisão:**
```
PARA CADA campo em ['token', 'userId'] no payload:
  SE payload[campo] contém '_'
    ENTÃO payload[campo] = removeTenant(payload[campo])
```

Exemplo:
- Input: `{ token: "myop_abc123", userId: "myop_user456" }`
- Output: `{ token: "abc123", userId: "user456" }`

**Casos Extremos:**
- Campo faltando no payload → Ignorado (campo optional)
- Campo sem underscore → Não modificado
- Campo é null/empty → Handled gracefully

**Código Fonte:**  
`PragmaticPlayService.php:132-137` (método `removeTenant()`)  
Chamado em: `authenticate()` linha 33, `balance()` linhas 54-55, `bet()` linha 70, etc.

**Dependências:** PP-002 (parsing logic é a mesma)

---

### Regra PP-005: MD5 Hash Generation for Request Authentication

**ID:** PP-005  
**Descrição:** Todo payload deve conter um campo `hash` calculado como MD5 do payload + secret-key do operador.

**Contexto de Negócio:**  
Pragmatic Play valida que requisições vêm de um cliente legítimo usando HMAC-based authentication. O hash prova posse da secret-key.

**Escopo:** Todos os endpoints (exceto o controller que apenas delega)

**Lógica de Decisão:**
```
FUNÇÃO generateHashCode(payload, secret_key):
  1. Remover campo 'hash' do payload (se existente)
  2. Fazer sort alphabético por chave do payload
  3. Construir URL query string: key1=val1&key2=val2&...
  4. Concatenar secret_key no final: "query_string" + secret_key
  5. URL-decode a string resultante
  6. Calcular MD5 da string
  7. Retornar MD5 em hexadecimal
```

**Exemplo:**
```
Payload: { providerId: "PragmaticPlay", token: "abc123", userId: "user456" }
Secret: "secret-key-value"

1. Sorted: providerId, token, userId
2. Query: "providerId=PragmaticPlay&token=abc123&userId=user456"
3. Concat: "providerId=PragmaticPlay&token=abc123&userId=user456secret-key-value"
4. URL-decode: (no change in this example)
5. MD5: "hash_value"
```

**Casos Extremos:**
- Hash field já presente no payload → Removido antes de cálculo
- Chaves com caracteres especiais → URL-encoded e depois decode
- Empty payload → MD5 de apenas secret-key

**Código Fonte:**  
`PragmaticPlayService.php:142-152` (método `generateHashCode()`)

**Dependências:** Nenhuma (puro cálculo)

**Validação:** Todos os testes verificam que hash está correto antes de envio

---

### Regra PP-006: Provider-Specific Credential Lookup

**ID:** PP-006  
**Descrição:** Cada operador possui credenciais armazenadas com `name='pragmatic'` e `key='secret-key'`, que são recuperadas do banco e usadas para hash generation.

**Contexto de Negócio:**  
Cada operador tem seu próprio secret-key fornecido pela Pragmatic Play. Sistema suporta múltiplos providers com credenciais diferentes.

**Escopo:** Todos os endpoints

**Lógica de Decisão:**
```
credential = operador.credentials()
  .where('name', 'pragmatic')
  .where('key', 'secret-key')
  .first()
  .value

SE credential não existir
  ENTÃO ERRO (implicit null reference exception)
```

**Casos Extremos:**
- Operador sem credenciais pragmatic → Exceção
- Múltiplas credenciais pragmatic → Primeira (`.first()`) é usada
- Credencial vazia → Usada como-é (pode causar hash inválido)

**Código Fonte:**  
`PragmaticPlayService.php:34, 56, 71, 86, 121, 162` (em cada método, pattern:
```php
$tenant->credentials()->where('name', 'pragmatic')->where('key', 'secret-key')->first()->value
```

**Dependências:** PP-003 (operador deve ser recuperado primeiro)

---

### Regra PP-007: Response Re-Prefixing (Authenticate Only)

**ID:** PP-007  
**Descrição:** Na resposta do endpoint `authenticate`, o campo `userId` é re-prefixado com o slug do operador (formato: `{operator_slug}_{userId}`), mas **apenas se o erro for 0** (sucesso).

**Contexto de Negócio:**  
O cliente (frontend) espera receber `userId` prefixado para rastrear qual operador/sessão cada requisição pertence. Apenas sucesso é re-prefixado porque em caso de erro, não há userId válido.

**Escopo:** Apenas endpoint `authenticate`

**Lógica de Decisão:**
```
response = http_post(provider_url, payload)

SE response['error'] == 0  # sucesso
  ENTÃO response['userId'] = operator_slug + '_' + response['userId']
  SENÃO deixar response inalterado

retornar response
```

**Exemplo:**
```
Entrada do Provider: { userId: "12345", error: 0, ... }
Output: { userId: "myoperator_12345", error: 0, ... }
```

**Casos Extremos:**
- error != 0 (falha) → userId NOT re-prefixed
- response['userId'] missing → String concatenation fails (edge case não tratado, poderia causar erro)
- error field missing → Comparação com 0 falha (type coercion)

**Código Fonte:**  
`PragmaticPlayService.php:40-42` (método `authenticate()`)

**Dependências:** PP-002, PP-003, PP-004, PP-005 (tudo precisa vir antes)

**Nota de Risco:** Outros endpoints (balance, bet, etc.) **NÃO** re-prefixam userId — apenas authenticate faz isso.

---

### Regra PP-008: HTTP POST Integration with Tenant-Specific Base URL

**ID:** PP-008  
**Descrição:** Todos os endpoints enviam requisição HTTP POST para URL construída como `{tenant.url}/pragmatic-play/{endpoint}.html`, com payload JSON.

**Contexto de Negócio:**  
Cada operador tem seu próprio backend (ou tunnel/proxy) fornecido como `operator.url`. Requisições são enviadas para esse backend específico, isolando dados por tenant.

**Escopo:** Todos os endpoints

**Lógica de Decisão:**
```
base_url = tenant['url']
endpoint_path = '/pragmatic-play/' + endpoint + '.html'
full_url = base_url + endpoint_path

response = HTTP.POST(full_url, payload)
retornar response.json()
```

**URLs Construídas:**
- `authenticate` → `{tenant_url}/pragmatic-play/authenticate.html`
- `balance` → `{tenant_url}/pragmatic-play/balance.html`
- `bet` → `{tenant_url}/pragmatic-play/bet.html`
- `refund` → `{tenant_url}/pragmatic-play/refund.html`
- `result` → `{tenant_url}/pragmatic-play/result.html`
- `bonusWin` → `{tenant_url}/pragmatic-play/bonusWin.html`
- `jackpotWin` → `{tenant_url}/pragmatic-play/jackpotWin.html`
- `promoWin` → `{tenant_url}/pragmatic-play/promoWin.html`
- `adjustment` → `{tenant_url}/pragmatic-play/adjustment.html`

**Casos Extremos:**
- tenant_url vazio → HTTP error
- HTTP timeout → exception (não tratado, retry foi comentado em BaseService linha 19)
- Provider retorna non-JSON → json() parsing falha

**Código Fonte:**  
`BaseService.php:16-22` (método `postJson()`)  
`PragmaticPlayService.php:37, 59, 74, 89, 165` (em cada método endpoint)

**Dependências:** PP-001 (routing), PP-002 (tenant extraction), PP-004 (token sanitization), PP-005 (hash generation)

**Performance Note:** Retry foi removido (comentado), então falhas não são retentadas automaticamente

---

### Regra PP-009: Error Handling on Unknown Endpoint

**ID:** PP-009  
**Descrição:** Quando cliente invoca um endpoint que não existe como método na classe, uma exceção é lançada com mensagem "Endpoint {endpoint} was not found." e HTTP status 500.

**Contexto de Negócio:**  
Protege contra requisições malformadas ou exploração. Garante que apenas endpoints conhecidos sejam processados.

**Escopo:** Aplicável a todos os endpoints (é a guarda do router)

**Lógica de Decisão:**
```
SE method_exists(this, endpoint)
  ENTÃO chamar método
  SENÃO
    log.error("Endpoint {endpoint} was not found.")
    throw new Exception("Endpoint {endpoint} was not found.", 500)
```

**Casos Extremos:**
- Endpoint é null → `method_exists()` retorna false
- Endpoint contém caracteres inválidos → `method_exists()` retorna false
- Endpoint é "call" (método público) → Aceito (não protegido)

**Código Fonte:**  
`PragmaticPlayService.php:18-22` (método `call()`)

**Dependências:** Nenhuma (é a primeira validação)

**Validação:** Teste `test_unknown_method_to_pragmatic_play_endpoint()` confirma status 500

---

### Regra PP-010: Balance Endpoint Dual Token Support

**ID:** PP-010  
**Descrição:** O endpoint `balance` aceita **tanto `token` quanto `userId`** para identificar o operador. Se `token` estiver presente, usa-o; caso contrário, usa `userId`.

**Contexto de Negócio:**  
Flexibilidade no cliente. O cliente pode fazer requisição usando token (mais seguro) ou userId (já conhecido da sessão anterior). Ambos carregam o tenant prefix.

**Escopo:** Apenas endpoint `balance`

**Lógica de Decisão:**
```
operador = get($data['token'] ?? $data['userId'])
# PHP: usa 'token' se presente, caso contrário 'userId'

PARA CADA campo em ['token', 'userId']:
  SE $data[campo] existe
    ENTÃO $data[campo] = removeTenant($data[campo])
```

**Exemplo de Requisições Válidas:**
- `{ token: "op_abc123", ... }` → usa token
- `{ userId: "op_user456", ... }` → usa userId
- `{ token: "op_abc123", userId: "op_user456", ... }` → usa token (precedência)
- `{ }` → erro (nenhum identificador)

**Casos Extremos:**
- Ambos faltando → Exceção ao chamar `get(null)`
- Token e userId com operadores diferentes → Usa token (operador mismatch não é verificado)

**Código Fonte:**  
`PragmaticPlayService.php:48-62` (método `balance()`)

**Dependências:** PP-002 (tenant extraction), PP-003 (lookup), PP-004 (sanitization), PP-005 (hash), PP-008 (HTTP POST)

---

### Regra PP-011: Token Sanitization Order

**ID:** PP-011  
**Descrição:** Quando um endpoint tem múltiplos campos de token (`token` e `userId`), eles são sanitizados **nessa ordem específica**: `token` primeiro, depois `userId`.

**Contexto de Negócio:**  
Garante consistência. A ordem não afeta o resultado (ambos são sanitizados), mas é importante para auditoria e debugging.

**Escopo:** Endpoints que têm múltiplos tokens: `balance`, `bet`, `refund`, `result`, `bonusWin`, `jackpotWin`, `promoWin`

**Lógica de Decisão:**
```
# balance exemplo
$data['token'] = removeTenant($data['token'] ?? $data['userId']);  # token primeiro
$data['userId'] = removeTenant($data['userId']);  # depois userId
```

**Código Fonte:**  
`PragmaticPlayService.php:54-55` (balance)  
`PragmaticPlayService.php:70` (bet — apenas userId)  
`PragmaticPlayService.php:85` (refund — apenas userId)  
`PragmaticPlayService.php:161` (handleResult — apenas userId)  
`PragmaticPlayService.php:120` (adjustment — apenas userId)

**Dependências:** PP-004 (removeTenant logic)

---

### Regra PP-012: Authenticate Method Only Returns Prefixed UserId

**ID:** PP-012  
**Descrição:** **Apenas o endpoint `authenticate` re-prefixar o userId na resposta.** Todos os outros endpoints (`balance`, `bet`, `refund`, `result`, `bonusWin`, `jackpotWin`, `promoWin`, `adjustment`) **retornam a resposta inalterada** do provider.

**Contexto de Negócio:**  
Diferentes endpoints têm diferentes contratos de resposta. Authenticate é especial porque inicializa a sessão e deve fornecer userId prefixado. Outros endpoints retornam dados de operação (saldo, transação, etc.) que não precisam de prefixo.

**Escopo:** Todos os endpoints, com tratamento diferenciado para `authenticate`

**Lógica de Decisão:**
```
// authenticate (SPECIAL)
$response = postJson(url, data)
IF response['error'] == 0:
  THEN response['userId'] = operator_slug + '_' + response['userId']
RETURN response

// todos os outros endpoints (NORMAL)
$response = postJson(url, data)
RETURN response  # inalterado
```

**Exemplos:**

Authenticate:
```json
// Input Provider
{ "userId": "12345", "error": 0, "cash": 100 }

// Output Casino Proxy
{ "userId": "myop_12345", "error": 0, "cash": 100 }
```

Balance:
```json
// Input Provider
{ "transactionId": "xyz", "cash": 500, "error": 0 }

// Output Casino Proxy (UNCHANGED)
{ "transactionId": "xyz", "cash": 500, "error": 0 }
```

**Casos Extremos:**
- authenticate com error != 0 → userId NOT prefixed
- Outro endpoint com userId no response → Retornado inalterado (não prefixado)

**Código Fonte:**  
`PragmaticPlayService.php:40-42` (apenas em `authenticate()`)  
Todos os outros métodos (linhas 48-62, 64-77, 79-92, 94-166) não fazem re-prefixing

**Dependências:** PP-007 is a subset of this rule

---

## Matriz de Dependências entre Regras

```
PP-001 (Dynamic Routing)
  ↓
PP-002 (Tenant Extraction)
  ↓
PP-003 (Operator Lookup + Cache)
  ├→ PP-006 (Credential Lookup)
  └→ PP-004 (Token Sanitization)
      ├→ PP-005 (Hash Generation)
      └→ PP-008 (HTTP POST Integration)

PP-007 (Response Re-Prefixing) — only for authenticate
PP-009 (Error Handling) — first check
PP-010 (Dual Token Support) — only for balance
PP-011 (Sanitization Order) — implementation detail
PP-012 (Authenticate Special) — business rule
```

---

## Endpoints e Suas Regras Aplicáveis

| Endpoint | Métodos Aplicáveis | Regras | Notas |
|----------|-------------------|--------|-------|
| `authenticate` | POST | PP-001,002,003,004,005,006,008,009,012 | Usa `token`, re-prefixar userId se sucesso |
| `balance` | POST | PP-001,002,003,004,005,006,008,009,010,011 | Dual token support (token ou userId) |
| `bet` | POST | PP-001,002,003,004,005,006,008,009,011 | Usa `userId` |
| `refund` | POST | PP-001,002,003,004,005,006,008,009,011 | Usa `userId` |
| `result` | POST | PP-001,002,003,004,005,006,008,009,011 | Delega para handleResult |
| `bonusWin` | POST | PP-001,002,003,004,005,006,008,009,011 | Delega para handleResult |
| `jackpotWin` | POST | PP-001,002,003,004,005,006,008,009,011 | Delega para handleResult |
| `promoWin` | POST | PP-001,002,003,004,005,006,008,009,011 | Delega para handleResult |
| `adjustment` | POST | PP-001,002,003,004,005,006,008,009,011 | Usa `userId` |

---

## Questões Abertas e Ambiguidades

1. **Race condition em cache invalidation:** Se operador for deletado/modificado, cache de 1 hora causará dados desatualizado. Sem mecanismo de invalidação.

2. **Retry ausente:** HTTP retry foi comentado em `BaseService:19`. Se provider está temporariamente indisponível, requisição falha imediatamente.

3. **Error handling inconsistente:** Alguns erros lançam exceção (unknown endpoint, missing operator), outros retornam resposta do provider (error codes). Sem padronização.

4. **UserId em resposta:** Regra PP-012 assume que resposta de authenticate sempre tem `userId`. Se provider não incluir, causará erro.

5. **Underscore em operador slug:** Suporta múltiplos underscores (ex: `"my_op_co"`), mas parsing usa último underscore como delimiter. Ambiguidade se slug contém underscore.

---

## Próximas Fases

- **Fase 2:** Documentar essas regras em formato de documentação técnica (markdown com diagrama)
- **Fase 3:** Construir suite de testes de integração que valida cada regra
- **Fase 4:** Criar matriz YAML de rastreamento para regras críticas
- **Fase 5:** Validar que testes PHP passam 100%
- **Implementação Go:** Implementar handlers Go que seguem essas regras exatamente

---

## Referências

- **Código:** `legacy/casino-proxy/app/Services/PragmaticPlayService.php`
- **Testes:** `legacy/casino-proxy/tests/Feature/PragmaticPlayControllerTest.php`
- **Models:** `app/Models/Operator.php`, `app/Models/Credential.php`
- **Parent Service:** `BaseService.php`
- **Operator Service:** `OperatorService.php`

