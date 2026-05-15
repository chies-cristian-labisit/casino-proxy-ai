# CASINO-2.3-endpoint-tests: Testes de Integração dos 9 Endpoints do Pragmatic Play

**Story ID:** CASINO-2.3-endpoint-tests  
**Epic:** CASINO-2.3 (Pragmatic Play — Fase 3 Test Oracle)  
**Tipo:** Implementação de Testes — Integração por Endpoint  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (Dex)  
**Relacionado:** CASINO-2.3-pragmatic-rules (pré-requisito), CASINO-2.2-authenticate, CASINO-2.2-bet, CASINO-2.2-refund, CASINO-2.2-result, CASINO-2.2-bonusWin, CASINO-2.2-jackpotWin, CASINO-2.2-promoWin, CASINO-2.2-adjustment  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Implementar 9 classes de teste de integração — uma por endpoint do Pragmatic Play — com cenários completos de happy path e error path, usando WireMock stubs JSON. Esta story completa a cobertura de 50+ test cases exigida pelo epic CASINO-2.3.

**Objetivo:** Provar end-to-end que cada um dos 9 endpoints processa requests corretamente, lida com erros e respeita os contratos definidos na Fase 2 (Technical Documentation).

---

## Contexto

### Por que esta Story?

As stories anteriores testam **regras de negócio individuais**. Esta story testa **endpoints completos** — o fluxo de 8 fases de ponta a ponta para cada endpoint. É o nível de teste mais próximo do que o CASINO-2.5 (Validation Gate) vai executar contra o PHP legado.

### Grupos de Endpoints

Os 9 endpoints se dividem em 4 grupos de comportamento (documentados no CASINO-2.2 epic):

| Grupo | Endpoints | Padrão de Response |
|-------|-----------|-------------------|
| **Sessão** | `authenticate` | Response transformation (re-prefixa userId) |
| **Consulta** | `balance` | Passthrough + dual token |
| **Transação inline** | `bet`, `refund`, `adjustment` | Passthrough + userId |
| **handleResult() family** | `result`, `bonusWin`, `jackpotWin`, `promoWin` | Passthrough via handleResult() |

### Como se Encaixa no Plano

```
CASINO-2.3-setup           ✅ (pré-requisito)
CASINO-2.3-generic-rules   ✅ (pré-requisito)
CASINO-2.3-pragmatic-rules ✅ (pré-requisito)
CASINO-2.3-endpoint-tests  ← ESTA STORY
CASINO-2.3-ci-cd            (aguarda esta)
```

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** `AuthenticateEndpointTest.java` — 5 cenários: sucesso com userId prefixado, erro sem userId prefixado, token inválido, operador não encontrado, hash inválido
- [ ] **AC-2:** `BalanceEndpointTest.java` — 5 cenários: sucesso via token, sucesso via userId, sem token nem userId → erro, operador não encontrado, provider timeout
- [ ] **AC-3:** `BetEndpointTest.java` — 4 cenários: sucesso passthrough, hash inválido, operador não encontrado, provider retorna erro
- [ ] **AC-4:** `RefundEndpointTest.java` — 4 cenários: sucesso passthrough, transação não encontrada (erro provider), hash inválido, operador não encontrado
- [ ] **AC-5:** `ResultEndpointTest.java` — 4 cenários: sucesso passthrough via handleResult(), hash inválido, operador não encontrado, provider retorna erro
- [ ] **AC-6:** `BonusWinEndpointTest.java` — 4 cenários: sucesso passthrough, hash inválido, operador não encontrado, provider retorna erro
- [ ] **AC-7:** `JackpotWinEndpointTest.java` — 4 cenários: sucesso passthrough (evento alto valor), hash inválido, operador não encontrado, provider retorna erro
- [ ] **AC-8:** `PromoWinEndpointTest.java` — 4 cenários: sucesso passthrough, hash inválido, operador não encontrado, provider retorna erro
- [ ] **AC-9:** `AdjustmentEndpointTest.java` — 4 cenários: sucesso passthrough (iniciado pelo operador), hash inválido, operador não encontrado, provider retorna erro
- [ ] **AC-10:** 12 WireMock stub JSON criados em `src/main/resources/wiremock/pragmatic-play/`
- [ ] **AC-11:** Total acumulado: 50+ test cases (`mvn test` reporta ≥ 50 testes passando)
- [ ] **AC-12:** Todos os 9 arquivos em `src/test/java/com/casino/oracle/integration/`
- [ ] **AC-13:** `mvn clean test` com BUILD SUCCESS
- [ ] **AC-14:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-15:** Cada classe de teste tem `@Tag("integration")` para permitir execução seletiva com `mvn test -Dgroups=integration`
- [ ] **AC-16:** `PayloadBuilder.java` implementado com métodos para construir payloads válidos de cada endpoint

### Fora do Escopo

- ❌ Testes de performance ou carga
- ❌ Testes contra PHP real (isso é CASINO-2.5)
- ❌ Outros providers além de Pragmatic Play
- ❌ Testes do Admin API ou endpoints internos

---

## Detalhes Técnicos / Dev Notes

### Padrão de Classe de Integração

```java
@Tag("integration")
class AuthenticateEndpointTest {

    static WireMockServer wireMock;

    @BeforeAll
    static void setup() {
        wireMock = ProviderMockServer.start();
        PragmaticPlayMocks.stubAuthenticateSuccess(wireMock, "12345");
        PragmaticPlayMocks.stubAuthenticateError(wireMock);
    }

    @AfterAll
    static void teardown() { wireMock.stop(); }

    @Test
    @DisplayName("Deve re-prefixar userId na resposta quando autenticação bem sucedida")
    void authenticateSuccessReturnsWithPrefixedUserId() {
        given()
            .baseUri(TestConfig.getBaseUrl())
            .contentType("application/json")
            .body(PragmaticPlayFixtures.authenticatePayload())
        .when()
            .post("/v1/webhooks/pragmatic-play/authenticate")
        .then()
            .statusCode(200)
            .body("error", equalTo(0))
            .body("userId", startsWith(PragmaticPlayFixtures.OPERATOR_SLUG + "_"));
    }
    
    // ... 4 outros testes
}
```

### 12 WireMock Stubs JSON a Criar

```
src/main/resources/wiremock/pragmatic-play/
├── authenticate-success.json      {"error":0,"userId":"12345"}
├── authenticate-error.json        {"error":9,"description":"Auth failed"}
├── balance-success.json           {"error":0,"balance":10000,"currency":"BRL"}
├── balance-error.json             {"error":9,"description":"Player not found"}
├── bet-success.json               {"error":0,"balance":9500,"currency":"BRL"}
├── bet-error.json                 {"error":9,"description":"Insufficient funds"}
├── refund-success.json            {"error":0,"balance":10000,"currency":"BRL"}
├── result-success.json            {"error":0,"balance":10500,"currency":"BRL"}
├── bonuswin-success.json          {"error":0,"balance":15000,"currency":"BRL"}
├── jackpotwin-success.json        {"error":0,"balance":100000,"currency":"BRL"}
├── promowin-success.json          {"error":0,"balance":12000,"currency":"BRL"}
└── adjustment-success.json        {"error":0,"balance":10000,"currency":"BRL"}
```

### PayloadBuilder.java — Métodos a Implementar

```java
public class PayloadBuilder {
    public static Map<String, Object> authenticatePayload(String token) { ... }
    public static Map<String, Object> balancePayloadWithToken(String token) { ... }
    public static Map<String, Object> balancePayloadWithUserId(String userId) { ... }
    public static Map<String, Object> betPayload(String userId, long amount) { ... }
    public static Map<String, Object> refundPayload(String userId, String txId) { ... }
    public static Map<String, Object> resultPayload(String userId, long win) { ... }
    public static Map<String, Object> handleResultPayload(String userId, String type) { ... }
    public static Map<String, Object> adjustmentPayload(String userId, long amount) { ... }
}
```

### Contagem de Test Cases

| Endpoint | Cenários | Acumulado |
|----------|----------|-----------|
| Regras BR-GENERIC-* (story anterior) | 16 | 16 |
| Regras BR-PRAGMATIC-* (story anterior) | 9 | 25 |
| `/authenticate` | 5 | 30 |
| `/balance` | 5 | 35 |
| `/bet` | 4 | 39 |
| `/refund` | 4 | 43 |
| `/result` | 4 | 47 |
| `/bonusWin` | 4 | 51 |
| `/jackpotWin` | 4 | 55 |
| `/promoWin` | 4 | 59 |
| `/adjustment` | 4 | **63 total** |

> ✅ Meta de 50+ atingida

### Documentação de Referência por Endpoint

| Endpoint | Doc Fase 2 |
|----------|-----------|
| `/authenticate` | `pragmatic-play-authenticate.md` (story CASINO-2.2) |
| `/balance` | `pragmatic-play-balance.md` (template canônico) |
| `/bet` | `pragmatic-play-bet.md` (story CASINO-2.2) |
| `/refund` | `pragmatic-play-refund.md` (story CASINO-2.2) |
| `/result` | `pragmatic-play-result.md` (story CASINO-2.2) |
| `/bonusWin` | `pragmatic-play-bonusWin.md` (story CASINO-2.2) |
| `/jackpotWin` | `pragmatic-play-jackpotWin.md` (story CASINO-2.2) |
| `/promoWin` | `pragmatic-play-promoWin.md` (story CASINO-2.2) |
| `/adjustment` | `pragmatic-play-adjustment.md` (story CASINO-2.2) |

> **Nota:** Os documentos Fase 2 para os 8 endpoints (exceto balance) ainda estão em Ready — @dev CASINO-2.2 precisa criá-los. Se não estiverem disponíveis ao iniciar esta story, usar `pragmatic-play-balance.md` como template e as regras BR-* como guia.

---

## Tasks / Subtasks

- [ ] **T-1:** Verificar se docs Fase 2 dos 8 endpoints existem em `docs/casino-proxy/phase-2-technical-documentation/` — usar como referência de fluxo e cenários de erro para cada endpoint
- [ ] **T-2:** Implementar `PayloadBuilder.java` com métodos para todos os 9 endpoints
- [ ] **T-3:** Criar 12 arquivos JSON de stubs WireMock em `src/main/resources/wiremock/pragmatic-play/`
- [ ] **T-4:** Atualizar `PragmaticPlayMocks.java` para carregar stubs JSON do classpath (em vez de inline)
- [ ] **T-5:** Criar `AuthenticateEndpointTest.java` — 5 cenários
- [ ] **T-6:** Criar `BalanceEndpointTest.java` — 5 cenários (incluindo dual token)
- [ ] **T-7:** Criar `BetEndpointTest.java` — 4 cenários
- [ ] **T-8:** Criar `RefundEndpointTest.java` — 4 cenários
- [ ] **T-9:** Criar `ResultEndpointTest.java` — 4 cenários (handleResult pattern)
- [ ] **T-10:** Criar `BonusWinEndpointTest.java` — 4 cenários
- [ ] **T-11:** Criar `JackpotWinEndpointTest.java` — 4 cenários
- [ ] **T-12:** Criar `PromoWinEndpointTest.java` — 4 cenários
- [ ] **T-13:** Criar `AdjustmentEndpointTest.java` — 4 cenários
- [ ] **T-14:** Executar `mvn clean test` — confirmar BUILD SUCCESS e ≥ 50 test cases passando
- [ ] **T-15:** Atualizar File List desta story

---

## CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Test`
- Complexidade: High (9 classes, 40+ cenários, 12 stubs JSON)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @qa

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): `mvn clean test` BUILD SUCCESS com ≥ 50 test cases
- [ ] Pre-PR (@devops): Verificar que 12 JSONs de stubs existem em `wiremock/pragmatic-play/`

**Self-Healing Configuration:**
```yaml
mode: light
max_iterations: 2
severity_filter: [CRITICAL, HIGH]
behavior:
  CRITICAL: auto_fix
  HIGH: document_as_debt
```

**Focus Areas (Integration Tests):**
- `AuthenticateEndpointTest`: cenário success COM prefixo E error SEM prefixo — ambos obrigatórios
- `BalanceEndpointTest`: 3 cenários de token obrigatórios (token, userId, nenhum)
- Stubs JSON em classpath (`wiremock/pragmatic-play/`), não inline
- `@Tag("integration")` em todas as 9 classes

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `src/test/java/com/casino/oracle/integration/AuthenticateEndpointTest.java` | 5 cenários /authenticate | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/BalanceEndpointTest.java` | 5 cenários /balance | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/BetEndpointTest.java` | 4 cenários /bet | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/RefundEndpointTest.java` | 4 cenários /refund | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/ResultEndpointTest.java` | 4 cenários /result | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/BonusWinEndpointTest.java` | 4 cenários /bonusWin | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/JackpotWinEndpointTest.java` | 4 cenários /jackpotWin | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/PromoWinEndpointTest.java` | 4 cenários /promoWin | ⏳ A Criar |
| `src/test/java/com/casino/oracle/integration/AdjustmentEndpointTest.java` | 4 cenários /adjustment | ⏳ A Criar |
| `src/main/resources/wiremock/pragmatic-play/*.json` (12 arquivos) | WireMock stubs | ⏳ A Criar |
| `src/main/java/com/casino/oracle/client/PayloadBuilder.java` | Builder de payloads | ⏳ A Criar |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/integration/*.java` (9 arquivos) | Testes por endpoint | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/resources/wiremock/pragmatic-play/*.json` (12 arquivos) | Stubs WireMock | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/client/PayloadBuilder.java` | Payload builder | ⏳ A Criar |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-*.md` (9 docs) | Referência de contratos | ⏳ CASINO-2.2 |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Regras BR-* | ✅ Existe |

---

## Definição de Pronto

- [ ] 9 classes de integração criadas (1 por endpoint)
- [ ] 12 stubs WireMock JSON criados em `wiremock/pragmatic-play/`
- [ ] `PayloadBuilder.java` implementado com métodos para os 9 endpoints
- [ ] `mvn clean test` BUILD SUCCESS com ≥ 50 test cases reportados
- [ ] `AuthenticateEndpointTest` cobre ambos os cenários de userId (com e sem prefixo)
- [ ] `BalanceEndpointTest` cobre os 3 casos de token (token, userId, nenhum)
- [ ] File List atualizada
- [ ] Pronto para validação @po antes de CASINO-2.3-ci-cd iniciar

---

## Estratégia de Teste

**Validação desta story:** `mvn clean test` — ≥ 50 testes passando, BUILD SUCCESS.  
**Validação de @po:** Confirmar que o total de testes cobre todos os 9 endpoints com happy path e error path.  
**Próxima Story:** CASINO-2.3-ci-cd (GitHub Actions pipeline + README).

---

## Notas

- **Criado:** 2026-05-15
- **Estimado:** 6-8 horas (maior story do epic — 9 classes + 12 stubs)
- **Depende De:** CASINO-2.3-pragmatic-rules (stubs e fixtures base já existem)
- **Bloqueia:** CASINO-2.3-ci-cd
- **Dependência soft:** Docs Fase 2 (CASINO-2.2 stories) idealmente completos antes desta story iniciar

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
| 2026-05-15 | @po (Pax) | Validação GO (8/10) — Status: Draft → Ready. Corrigido: CodeRabbit mode standard → light. Should-fix: adicionar seção Riscos. |
