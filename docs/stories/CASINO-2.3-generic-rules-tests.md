# CASINO-2.3-generic-rules: Testes das Regras BR-GENERIC-* do Pragmatic Play

**Story ID:** CASINO-2.3-generic-rules  
**Epic:** CASINO-2.3 (Pragmatic Play — Fase 3 Test Oracle)  
**Tipo:** Implementação de Testes — Regras de Negócio Genéricas  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (Dex)  
**Relacionado:** CASINO-2.3-setup (infraestrutura ✅ pré-requisito), CASINO-1.7 (regras BR-* extraídas)  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Implementar 7 classes de teste JUnit 5 cobrindo todas as regras **BR-GENERIC-*** do Pragmatic Play, usando WireMock para simular respostas do provider externo. Cada regra tem sua própria classe de teste com múltiplos cenários (happy path + error path).

**Objetivo:** Provar que o sistema sob teste (PHP hoje, Go amanhã) aplica corretamente as 7 regras genéricas que afetam todos os 9 endpoints do Pragmatic Play.

---

## Contexto

### Por que esta Story?

As regras BR-GENERIC-* são o backbone do Pragmatic Play — afetam todos os 9 endpoints. Se qualquer uma falhar na migração Go, 9 endpoints quebram. Esta story valida exatamente esse comportamento antes de qualquer implementação Go.

### Regras BR-GENERIC-* em Escopo

| # | Regra | O que valida |
|---|-------|-------------|
| 1 | BR-GENERIC-ROUTING-VALIDATION-001 | Endpoint inválido → 500; endpoint válido → roteado |
| 2 | BR-GENERIC-TENANT-EXTRACTION-001 | Token `operator_abc` → extrai `operator` como slug |
| 3 | BR-GENERIC-OPERATOR-CACHING-001 | Operador encontrado no DB; operador não encontrado → erro |
| 4 | BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | MD5 correto passa; MD5 inválido → rejeitado pelo provider |
| 5 | BR-GENERIC-CREDENTIAL-LOOKUP-001 | Secret-key encontrada; credencial ausente → erro |
| 6 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST correto ao provider; timeout → falha |
| 7 | BR-GENERIC-RESPONSE-PASSTHROUGH-001 | Response do provider chega inalterada ao cliente |
| 8 | BR-GENERIC-ERROR-HANDLING-001 | Endpoint inválido lança exception → 500 |

> **Fonte:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Como se Encaixa no Plano

```
CASINO-2.3-setup          ✅ (pré-requisito)
CASINO-2.3-generic-rules  ← ESTA STORY
CASINO-2.3-pragmatic-rules (aguarda esta)
CASINO-2.3-endpoint-tests  (aguarda pragmatic-rules)
CASINO-2.3-ci-cd           (aguarda endpoint-tests)
```

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** `RoutingValidationTest.java` — endpoint inválido retorna 500; endpoint válido é roteado corretamente
- [ ] **AC-2:** `TenantExtractionTest.java` — token `"myoperator_abc123"` extrai slug `"myoperator"`; token sem underscore resulta em erro
- [ ] **AC-3:** `OperatorCachingTest.java` — operador existente retorna dados; operador inexistente retorna erro
- [ ] **AC-4:** `AuthenticationHmacTest.java` — MD5 correto: provider aceita; MD5 incorreto: provider retorna error != 0
- [ ] **AC-5:** `CredentialLookupTest.java` — credencial `secret-key` presente: fluxo continua; credencial ausente: erro antes de chegar ao provider
- [ ] **AC-6:** `ProviderIntegrationTest.java` — HTTP POST chega ao WireMock stub; timeout do provider resulta em erro
- [ ] **AC-7:** `ResponsePassthroughTest.java` — JSON de resposta do provider (via WireMock) chega idêntico ao cliente, sem modificação
- [ ] **AC-8:** `ErrorHandlingTest.java` — endpoint `/v1/webhooks/pragmatic-play/invalid-endpoint` retorna HTTP 500
- [ ] **AC-9:** Todos os 8 arquivos em `src/test/java/com/casino/oracle/rules/`
- [ ] **AC-10:** `mvn clean test` com BUILD SUCCESS — todos os testes passando
- [ ] **AC-11:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-12:** Cada classe de teste tem `@DisplayName` descritivo em português (ex: `"Deve retornar 500 para endpoint inválido"`)
- [ ] **AC-13:** Stubs WireMock reutilizáveis via `PragmaticPlayMocks.java` (não inline em cada teste)

### Fora do Escopo

- ❌ Regras BR-PRAGMATIC-* (story CASINO-2.3-pragmatic-rules)
- ❌ Testes por endpoint (story CASINO-2.3-endpoint-tests)
- ❌ Stubs WireMock para todos os 9 endpoints — apenas os necessários para estas 8 regras
- ❌ Testar Evolution Gaming ou qualquer outro provider

---

## Detalhes Técnicos / Dev Notes

### Padrão de Cada Classe de Teste

```java
@ExtendWith(...)
class RoutingValidationTest {

    static WireMockServer wireMock;
    
    @BeforeAll
    static void setup() {
        wireMock = ProviderMockServer.start();
    }
    
    @AfterAll
    static void teardown() {
        wireMock.stop();
    }
    
    @Test
    @DisplayName("Deve retornar 500 para endpoint inválido")
    void invalidEndpointReturns500() {
        // given: request para /v1/webhooks/pragmatic-play/invalid-endpoint
        // when: sistema sob teste processa
        // then: HTTP 500
    }
    
    @Test
    @DisplayName("Deve rotear corretamente para endpoint válido")
    void validEndpointIsRouted() {
        // given: stub WireMock para /authenticate.html
        // when: request para /v1/webhooks/pragmatic-play/authenticate
        // then: HTTP 200, WireMock recebeu a chamada
    }
}
```

### Stubs WireMock para estas Regras

Criar stubs mínimos em `PragmaticPlayMocks.java`:
```java
// Stub genérico para simular provider aceitando request
public static void stubProviderAccepts(WireMockServer server, String endpoint) {
    server.stubFor(post(urlEqualTo("/pragmatic-play/" + endpoint + ".html"))
        .willReturn(aResponse().withStatus(200)
            .withBody("{\"error\":0}")));
}

// Stub para simular provider rejeitando (hash inválido)
public static void stubProviderRejects(WireMockServer server, String endpoint) {
    server.stubFor(post(urlEqualTo("/pragmatic-play/" + endpoint + ".html"))
        .willReturn(aResponse().withStatus(200)
            .withBody("{\"error\":9,\"description\":\"Invalid hash\"}")));
}
```

### Fixtures Necessárias em `PragmaticPlayFixtures.java`

```java
// Token válido com prefixo tenant
public static final String VALID_TOKEN = "myoperator_abc123def456";
public static final String OPERATOR_SLUG = "myoperator";

// Token sem prefixo (inválido)
public static final String INVALID_TOKEN_NO_PREFIX = "abc123def456";

// Payload mínimo para autenticação
public static Map<String, Object> minimalAuthPayload() {
    return Map.of("token", VALID_TOKEN, "hash", "placeholder");
}
```

### Referências de Regras

```
BR-GENERIC-ROUTING-VALIDATION-001:
  Fonte: pragmatic-play-rules.md#routing-validation
  PHP:   PragmaticPlayService.php:18-25

BR-GENERIC-TENANT-EXTRACTION-001:
  Fonte: pragmatic-play-rules.md#tenant-extraction
  PHP:   OperatorService.php:20-34 (método get())

BR-GENERIC-OPERATOR-CACHING-001:
  Fonte: pragmatic-play-rules.md#operator-caching
  PHP:   OperatorService.php:20-34 (cache 1h)

BR-GENERIC-AUTHENTICATION-HMAC-MD5-001:
  Fonte: pragmatic-play-rules.md#hmac-md5
  PHP:   PragmaticPlayService.php:142-152

BR-GENERIC-CREDENTIAL-LOOKUP-001:
  Fonte: pragmatic-play-rules.md#credential-lookup
  PHP:   PragmaticPlayService.php (credentials.where)

BR-GENERIC-PROVIDER-INTEGRATION-001:
  Fonte: pragmatic-play-rules.md#provider-integration
  PHP:   BaseService.php:16-22 (postJson)

BR-GENERIC-ERROR-HANDLING-001:
  Fonte: pragmatic-play-rules.md#error-handling
  PHP:   PragmaticPlayService.php:18-25 (exception)

BR-GENERIC-RESPONSE-PASSTHROUGH-001:
  Fonte: pragmatic-play-rules.md#response-passthrough
  PHP:   PragmaticPlayService.php (return response)
```

---

## Tasks / Subtasks

- [ ] **T-1:** Ler `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — focar nas 8 regras BR-GENERIC-* (lógica de decisão e casos extremos de cada uma)
- [ ] **T-2:** Criar `PragmaticPlayFixtures.java` com `VALID_TOKEN`, `OPERATOR_SLUG`, `INVALID_TOKEN_NO_PREFIX` e `minimalAuthPayload()`
- [ ] **T-3:** Criar `PragmaticPlayMocks.java` com `stubProviderAccepts()` e `stubProviderRejects()` reutilizáveis
- [ ] **T-4:** Criar `RoutingValidationTest.java` — 2 testes: endpoint inválido → 500, endpoint válido → roteado
- [ ] **T-5:** Criar `TenantExtractionTest.java` — 2 testes: token válido extrai slug, token sem `_` resulta em erro
- [ ] **T-6:** Criar `OperatorCachingTest.java` — 2 testes: operador existe retorna dados, operador inexistente retorna erro
- [ ] **T-7:** Criar `AuthenticationHmacTest.java` — 2 testes: MD5 correto aceito, MD5 incorreto rejeitado
- [ ] **T-8:** Criar `CredentialLookupTest.java` — 2 testes: credencial presente fluxo continua, credencial ausente erro
- [ ] **T-9:** Criar `ProviderIntegrationTest.java` — 2 testes: POST chega ao WireMock, timeout resulta em erro
- [ ] **T-10:** Criar `ResponsePassthroughTest.java` — 1 teste: JSON do WireMock chega idêntico ao cliente
- [ ] **T-11:** Criar `ErrorHandlingTest.java` — 1 teste: endpoint inválido retorna 500
- [ ] **T-12:** Executar `mvn clean test` — confirmar todos os 8 novos testes passam
- [ ] **T-13:** Atualizar File List desta story

---

## CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Test`
- Complexidade: Medium (8 classes, padrões WireMock)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @qa

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): `mvn clean test` com BUILD SUCCESS
- [ ] Pre-PR (@devops): Nenhum URL hardcoded; stubs WireMock em `PragmaticPlayMocks`, não inline

**Self-Healing Configuration:**
```yaml
mode: light
max_iterations: 2
severity_filter: [CRITICAL, HIGH]
behavior:
  CRITICAL: auto_fix
  HIGH: document_as_debt
```

**Focus Areas (Tests):**
- Cobertura: 1 classe por regra (sem consolidar em mega-classe)
- WireMock: stubs centralizados em `PragmaticPlayMocks`, não inline nos testes
- Fixtures: dados de teste em `PragmaticPlayFixtures`, não hardcoded
- `@DisplayName` descritivos em português

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `src/test/java/com/casino/oracle/rules/RoutingValidationTest.java` | Testa BR-GENERIC-ROUTING-VALIDATION-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/TenantExtractionTest.java` | Testa BR-GENERIC-TENANT-EXTRACTION-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/OperatorCachingTest.java` | Testa BR-GENERIC-OPERATOR-CACHING-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/AuthenticationHmacTest.java` | Testa BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/CredentialLookupTest.java` | Testa BR-GENERIC-CREDENTIAL-LOOKUP-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/ProviderIntegrationTest.java` | Testa BR-GENERIC-PROVIDER-INTEGRATION-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/ResponsePassthroughTest.java` | Testa BR-GENERIC-RESPONSE-PASSTHROUGH-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/ErrorHandlingTest.java` | Testa BR-GENERIC-ERROR-HANDLING-001 | ⏳ A Criar |
| `src/main/java/com/casino/oracle/mock/PragmaticPlayMocks.java` | Stubs WireMock reutilizáveis | ⏳ A Criar |
| `src/main/java/com/casino/oracle/data/PragmaticPlayFixtures.java` | Fixtures de dados de teste | ⏳ A Criar |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/rules/*.java` (8 arquivos) | Testes das regras genéricas | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/mock/PragmaticPlayMocks.java` | Stubs centralizados | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/data/PragmaticPlayFixtures.java` | Dados de teste | ⏳ A Criar |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das regras | ✅ Existe |
| `casino-proxy-test-oracle/` (projeto base) | Infraestrutura | ⏳ CASINO-2.3-setup |

---

## Definição de Pronto

- [ ] 8 classes de teste criadas (1 por regra BR-GENERIC-*)
- [ ] `PragmaticPlayMocks.java` com stubs reutilizáveis
- [ ] `PragmaticPlayFixtures.java` com dados de teste
- [ ] `mvn clean test` retorna BUILD SUCCESS — todos os testes passando
- [ ] Nenhuma URL ou dado hardcoded nos testes (tudo via Fixtures/Mocks/TestConfig)
- [ ] File List atualizada
- [ ] Pronto para validação @po antes de CASINO-2.3-pragmatic-rules iniciar

---

## Estratégia de Teste

**Validação desta story:** `mvn clean test` — todos os testes da pasta `rules/` passando.  
**Validação de @po:** Confirmar que cada classe testa exatamente a regra BR-* correspondente (rastreabilidade).  
**Próxima Story:** CASINO-2.3-pragmatic-rules (4 regras BR-PRAGMATIC-* + 1 nova).

---

## Notas

- **Criado:** 2026-05-15
- **Estimado:** 4-6 horas
- **Depende De:** CASINO-2.3-setup (projeto Maven base)
- **Bloqueia:** CASINO-2.3-pragmatic-rules

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
| 2026-05-15 | @po (Pax) | Validação GO (8/10) — Status: Draft → Ready. Corrigido: CodeRabbit mode standard → light. Should-fix: adicionar seção Riscos. |
