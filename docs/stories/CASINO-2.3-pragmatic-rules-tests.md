# CASINO-2.3-pragmatic-rules: Testes das Regras BR-PRAGMATIC-* do Pragmatic Play

**Story ID:** CASINO-2.3-pragmatic-rules  
**Epic:** CASINO-2.3 (Pragmatic Play — Fase 3 Test Oracle)  
**Tipo:** Implementação de Testes — Regras de Negócio Pragmatic-Específicas  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (Dex)  
**Relacionado:** CASINO-2.3-generic-rules (pré-requisito), CASINO-2.2-authenticate, CASINO-2.2-bet  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Implementar 5 classes de teste JUnit 5 cobrindo as regras **BR-PRAGMATIC-*** exclusivas do Pragmatic Play. Estas regras afetam subconjuntos específicos de endpoints (não todos os 9), com destaque para o comportamento único do `/authenticate` (re-prefixação do `userId`) e o suporte dual-token do `/balance`.

**Objetivo:** Provar que os comportamentos exclusivos do Pragmatic Play — que o distinguem de outros providers — são replicados corretamente pelo sistema sob teste.

---

## Contexto

### Por que esta Story?

As regras BR-PRAGMATIC-* são os comportamentos que **diferenciam o Pragmatic Play dos demais providers**. Se a migração Go errar qualquer uma delas, funcionalidades críticas quebram: jogadores não conseguem autenticar (`userId` sem prefixo), ou transações falham por token incorretamente sanitizado.

### Regras BR-PRAGMATIC-* em Escopo

| # | Regra | Endpoints Afetados | O que valida |
|---|-------|-------------------|-------------|
| 1 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | `/balance`, `/bet`, `/refund`, `/result`, family | Prefixo tenant removido do token antes de enviar ao provider |
| 2 | BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001 | `/balance` | Aceita `token` OU `userId`; ambos ausentes → erro |
| 3 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001 | `/balance` | `token` processado antes de `userId` quando ambos presentes |
| 4 | BR-PRAGMATIC-AUTHENTICATE-USERID-REPREFIX-001 | `/authenticate` | `userId` re-prefixado com `operator_slug` quando `error==0`; inalterado quando `error!=0` |
| 5 | (implícito em PP-012) | `/authenticate` | Apenas `/authenticate` faz transform — todos os outros são passthrough |

> **Nota sobre regra 5:** O isolamento do comportamento de re-prefixação apenas para `/authenticate` não tem ID BR-* explícito nas regras extraídas, mas é um comportamento crítico validado aqui para garantir que nenhum outro endpoint aplique a transformação por engano.

> **Fonte:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Como se Encaixa no Plano

```
CASINO-2.3-setup          ✅ (pré-requisito)
CASINO-2.3-generic-rules  ✅ (pré-requisito)
CASINO-2.3-pragmatic-rules ← ESTA STORY
CASINO-2.3-endpoint-tests  (aguarda esta)
CASINO-2.3-ci-cd           (aguarda endpoint-tests)
```

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** `TokenSanitizationTest.java` — token `"myoperator_abc123"` enviado ao provider como `"abc123"` (prefixo removido); token sem prefixo permanece inalterado
- [ ] **AC-2:** `DualTokenSupportTest.java` — request com `token` funciona; request com `userId` funciona; request sem nenhum dos dois resulta em erro
- [ ] **AC-3:** `TokenSanitizationOrderTest.java` — quando `token` e `userId` estão ambos presentes, `token` é processado primeiro (ordem de sanitização validada)
- [ ] **AC-4:** `AuthenticateTransformTest.java` — quando `error==0`: `userId` na resposta recebe prefixo `"myoperator_"`; quando `error!=0`: `userId` retornado sem modificação
- [ ] **AC-5:** `PassthroughIsolationTest.java` — endpoint `/bet` com `error==0`: `userId` **não** é re-prefixado (confirma que transformação é exclusiva do `/authenticate`)
- [ ] **AC-6:** Todos os 5 arquivos em `src/test/java/com/casino/oracle/rules/`
- [ ] **AC-7:** `mvn clean test` com BUILD SUCCESS — todos os testes (incluindo os da story anterior) passando
- [ ] **AC-8:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-9:** Stubs WireMock adicionados em `PragmaticPlayMocks.java` para simular respostas com `error==0` e `error!=0`
- [ ] **AC-10:** Fixtures em `PragmaticPlayFixtures.java` para tokens com e sem prefixo, e payloads com `token` vs `userId`

### Fora do Escopo

- ❌ Testes de integração por endpoint (story CASINO-2.3-endpoint-tests)
- ❌ Regras BR-GENERIC-* (story CASINO-2.3-generic-rules)
- ❌ Outros providers (Evolution Gaming, PG Soft, etc.)
- ❌ Implementar a lógica — apenas testar o comportamento via HTTP

---

## Detalhes Técnicos / Dev Notes

### Regra mais crítica: AuthenticateTransformTest

```java
@Test
@DisplayName("Deve re-prefixar userId quando error==0")
void reпрефixesUserIdOnSuccess() {
    // given: WireMock stub retorna {"error":0,"userId":"12345"}
    // when: POST /v1/webhooks/pragmatic-play/authenticate com token válido
    // then: response.userId == "myoperator_12345"
}

@Test
@DisplayName("Não deve modificar userId quando error!=0")
void doesNotModifyUserIdOnError() {
    // given: WireMock stub retorna {"error":9,"userId":"12345"}
    // when: POST /v1/webhooks/pragmatic-play/authenticate com token válido
    // then: response.userId == "12345" (sem prefixo)
}
```

### Stubs WireMock a Adicionar em PragmaticPlayMocks.java

```java
// Simula provider retornando sucesso com userId
public static void stubAuthenticateSuccess(WireMockServer server, String userId) {
    server.stubFor(post(urlEqualTo("/pragmatic-play/authenticate.html"))
        .willReturn(aResponse().withStatus(200)
            .withHeader("Content-Type", "application/json")
            .withBody("{\"error\":0,\"userId\":\"" + userId + "\"}")));
}

// Simula provider retornando erro de autenticação
public static void stubAuthenticateError(WireMockServer server) {
    server.stubFor(post(urlEqualTo("/pragmatic-play/authenticate.html"))
        .willReturn(aResponse().withStatus(200)
            .withHeader("Content-Type", "application/json")
            .withBody("{\"error\":9,\"description\":\"Authentication failed\"}")));
}
```

### Fixtures a Adicionar em PragmaticPlayFixtures.java

```java
// Tokens para sanitização
public static final String TOKEN_WITH_PREFIX = "myoperator_abc123def456";
public static final String TOKEN_WITHOUT_PREFIX = "abc123def456";

// UserId do provider (antes da re-prefixação)
public static final String PROVIDER_USER_ID = "12345";
public static final String EXPECTED_PREFIXED_USER_ID = "myoperator_12345";

// Payload com userId em vez de token
public static Map<String, Object> payloadWithUserId() {
    return Map.of("userId", PROVIDER_USER_ID, "hash", "placeholder");
}
```

### Referências de Regras no Documento Phase-1

```
BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001:
  Fonte: pragmatic-play-rules.md#token-sanitization
  PHP:   PragmaticPlayService.php:132-137 (removeTenant())

BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001:
  Fonte: pragmatic-play-rules.md#dual-token
  PHP:   PragmaticPlayService.php:balance() (token OU userId)

BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001:
  Fonte: pragmatic-play-rules.md#sanitization-order
  PHP:   PragmaticPlayService.php:balance() (token primeiro)

BR-PRAGMATIC-AUTHENTICATE-USERID-REPREFIX-001:
  Fonte: pragmatic-play-rules.md#authenticate-transform
  PHP:   PragmaticPlayService.php:40-42 (re-prefixação)
```

---

## Tasks / Subtasks

- [ ] **T-1:** Ler `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — focar nas seções das 4 regras BR-PRAGMATIC-* e no comportamento de re-prefixação do authenticate
- [ ] **T-2:** Ler `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` — internalizar fluxo Fase 8 (re-prefixação) e cenários de erro documentados
- [ ] **T-3:** Adicionar stubs `stubAuthenticateSuccess()` e `stubAuthenticateError()` em `PragmaticPlayMocks.java`
- [ ] **T-4:** Adicionar fixtures `TOKEN_WITH_PREFIX`, `PROVIDER_USER_ID`, `EXPECTED_PREFIXED_USER_ID`, `payloadWithUserId()` em `PragmaticPlayFixtures.java`
- [ ] **T-5:** Criar `TokenSanitizationTest.java` — 2 testes: prefixo removido antes do provider, sem prefixo permanece inalterado
- [ ] **T-6:** Criar `DualTokenSupportTest.java` — 3 testes: via token, via userId, sem nenhum → erro
- [ ] **T-7:** Criar `TokenSanitizationOrderTest.java` — 1 teste: token processado antes de userId quando ambos presentes
- [ ] **T-8:** Criar `AuthenticateTransformTest.java` — 2 testes: userId prefixado em sucesso, userId inalterado em erro
- [ ] **T-9:** Criar `PassthroughIsolationTest.java` — 1 teste: /bet não re-prefixa userId mesmo quando error==0
- [ ] **T-10:** Executar `mvn clean test` — confirmar todos os testes (incluindo da story anterior) passando
- [ ] **T-11:** Atualizar File List desta story

---

## CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Test`
- Complexidade: Medium (comportamentos exclusivos, lógica condicional)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @qa

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): `mvn clean test` BUILD SUCCESS
- [ ] Pre-PR (@devops): `PassthroughIsolationTest` presente (garante que transformação não vaza para outros endpoints)

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
- `AuthenticateTransformTest`: cenário `error==0` E `error!=0` ambos cobertos
- `PassthroughIsolationTest`: teste de isolamento é obrigatório — risco alto se omitido
- Stubs WireMock em `PragmaticPlayMocks`, não inline
- Fixtures em `PragmaticPlayFixtures`, não hardcoded

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `src/test/java/com/casino/oracle/rules/TokenSanitizationTest.java` | Testa BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/DualTokenSupportTest.java` | Testa BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/TokenSanitizationOrderTest.java` | Testa BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001 | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/AuthenticateTransformTest.java` | Testa re-prefixação userId no /authenticate | ⏳ A Criar |
| `src/test/java/com/casino/oracle/rules/PassthroughIsolationTest.java` | Garante transformação isolada no /authenticate | ⏳ A Criar |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/rules/TokenSanitizationTest.java` | Testa sanitização | ⏳ A Criar |
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/rules/DualTokenSupportTest.java` | Testa dual token | ⏳ A Criar |
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/rules/TokenSanitizationOrderTest.java` | Testa ordem sanitização | ⏳ A Criar |
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/rules/AuthenticateTransformTest.java` | Testa re-prefixação | ⏳ A Criar |
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/rules/PassthroughIsolationTest.java` | Testa isolamento da transformação | ⏳ A Criar |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das regras | ✅ Existe |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` | Doc do endpoint mais crítico | ⏳ CASINO-2.2 |

---

## Definição de Pronto

- [ ] 5 classes de teste criadas (4 regras BR-PRAGMATIC-* + 1 isolamento)
- [ ] `AuthenticateTransformTest` cobre ambos os cenários: `error==0` (com prefixo) e `error!=0` (sem prefixo)
- [ ] `PassthroughIsolationTest` confirma que `/bet` não re-prefixa userId
- [ ] Stubs e fixtures atualizados sem quebrar testes da story anterior
- [ ] `mvn clean test` BUILD SUCCESS — todos os testes passando
- [ ] File List atualizada
- [ ] Pronto para validação @po antes de CASINO-2.3-endpoint-tests iniciar

---

## Estratégia de Teste

**Validação desta story:** `mvn clean test` — todas as 13 classes de teste passando (8 da generic-rules + 5 desta story).  
**Validação de @po:** Confirmar que `AuthenticateTransformTest` e `PassthroughIsolationTest` cobrem o comportamento mais crítico e diferenciado do Pragmatic Play.  
**Próxima Story:** CASINO-2.3-endpoint-tests (9 endpoints, 40+ cenários de integração).

---

## Notas

- **Criado:** 2026-05-15
- **Estimado:** 3-4 horas
- **Depende De:** CASINO-2.3-generic-rules (classes base, stubs e fixtures já existem)
- **Bloqueia:** CASINO-2.3-endpoint-tests

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
| 2026-05-15 | @po (Pax) | Validação GO (8/10) — Status: Draft → Ready. Corrigido: CodeRabbit mode standard → light. Should-fix: adicionar seção Riscos. |
