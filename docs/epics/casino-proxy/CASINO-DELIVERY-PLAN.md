# Plano de Entrega — Casino Proxy Migration Roadmap

**Data:** 2026-05-11  
**Status:** 🚀 Em Execução — CASINO-2.1 (Pragmatic Play) em andamento  
**Próxima Revisão:** 2026-05-18

---

## 1. Visão Geral Executiva

### Estrutura de 3 Epics Sequenciais

```
CASINO-1: OpenAPI Documentation ✅ COMPLETO
   ↓ (bloqueia por documentação)
CASINO-2: Business Rules Discovery & Test Oracle 🚀 EM ANDAMENTO
   ↓ (bloqueia por teste de validação)
CASINO-3: Go Microservices Implementation ⏸️ AGUARDANDO
```

### Tabela de Resumo Executivo

| Métrica | CASINO-1 | CASINO-2 | CASINO-3 |
|---------|----------|----------|----------|
| **Objetivo** | Documentar 100% endpoints | Extrair & validar regras | Reconstruir em Go |
| **Stories** | 6 | 40 | 22 |
| **Status** | ✅ Completo | 🚀 Iniciando | ⏸️ Aguardando |
| **Depende De** | — | CASINO-1 | CASINO-2 |
| **Bloqueia** | — | CASINO-3 | — |
| **Artefatos** | OpenAPI specs | Rules + Tests + Oracle | Go services |
| **Validação** | @po | @po (por fase) | @qa + @devops |
| **Risco** | ✅ Baixo | 🟡 Médio | 🟡 Médio |

---

## 2. CASINO-1 — Documentação OpenAPI ✅ COMPLETO

### Objetivo
Documentar 100% dos endpoints e integrações em padrão OpenAPI 3.0, criando o contrato técnico que guia todas as fases subsequentes.

### Escopo — 6 Stories Entregues

| ID | Descrição | Status | Artefato |
|----|-----------|--------|----------|
| CASINO-1.1 | Analyze & Document Pragmatic Play | ✅ | `/docs/casino-proxy/openapi/pragmatic-play.yaml` |
| CASINO-1.2 | Analyze & Document Evolution Gaming | ✅ | `/docs/casino-proxy/openapi/evolution-gaming.yaml` |
| CASINO-1.3 | Analyze & Document PG Soft | ✅ | `/docs/casino-proxy/openapi/pg-soft.yaml` |
| CASINO-1.4 | Analyze & Document Remaining Providers | ✅ | `/docs/casino-proxy/openapi/mancala.yaml` + 3 others |
| CASINO-1.5 | Create Master OpenAPI Spec | ✅ | `/docs/casino-proxy/openapi/casino-proxy-master.yaml` |
| CASINO-1.6 | Document Admin API | ✅ | `/docs/casino-proxy/openapi/admin-api.yaml` |

### Artefatos Entregues

```
docs/casino-proxy/openapi/
├── pragmatic-play.yaml          ✅ 8 endpoints
├── evolution-gaming.yaml        ✅ 12 endpoints
├── pg-soft.yaml                 ✅ 7 endpoints
├── mancala.yaml                 ✅ 6 endpoints
├── digitain.yaml                ✅ 5 endpoints
├── evoplay.yaml                 ✅ 8 endpoints
├── openbox.yaml                 ✅ 4 endpoints
├── alternar.yaml                ✅ 5 endpoints
├── casino-proxy-master.yaml     ✅ (union de todos)
└── admin-api.yaml               ✅ 15 endpoints
```

### Definition of Done ✅

- [x] 8 provedores documentados
- [x] Todos endpoints com request/response schemas
- [x] Autenticação documentada (signatures, headers)
- [x] Error codes mapeados por provider
- [x] Master spec referencia todos os provedores
- [x] Admin API completa
- [x] Validação com @po
- [x] Pronto para merge em master

### Status Atual
**PRONTO PARA MERGE** — Nenhuma tarefa pendente. Branch `feature/openapi-documentation-viewers` pode ser mergeado em `master` com validação de @po.

---

## 3. CASINO-2 — Business Rules Discovery & Test Oracle 🚀 INICIANDO

### Objetivo
Extrair e documentar todas as regras de negócio do sistema PHP (1 provider por semana), construindo uma suite de testes que serve como "oráculo" — prova que o sistema Go será idêntico ao PHP.

### Por Que Esta Fase?
- **Risco:** Sem documentação de regras, é impossível saber se o novo sistema cumpre os requisitos
- **Validação:** Testes contra PHP = acceptance criteria para CASINO-3
- **Escalabilidade:** Templates reutilizáveis para os 8 provedores

### Estrutura: 5 Fases por Provider

Cada um dos 8 provedores segue a mesma sequência (5 fases):

```
Fase 1: EXTRACT     → Ler código PHP, extrair regras
Fase 2: DOCUMENT    → Escrever markdown com fluxos e regras
Fase 3: TEST        → Construir suite de testes de integração
Fase 4: MATRIX      → YAML trace matrices (regra → spec → teste)
Fase 5: VALIDATE    → Testes PHP passam 100%
```

### 🔍 Modelo Detalhado: Pragmatic Play (Provider 1 de 8)

#### **Pragmatic Play — Overview**

- **Endpoints:** 9 (authenticate, balance, bet, refund, result, bonusWin, jackpotWin, promoWin, adjustment)
- **Integração:** HTTP POST + MD5 signature + tenant isolation
- **Stories:** CASINO-2.1 até CASINO-2.5 (5 fases)

---

### **Fase 1: EXTRACT — Extrair Regras de Negócio** 🔨

**Story:** CASINO-2.1  
**Status:** ✅ COMPLETO

#### O que foi feito

Análise de `legacy/casino-proxy/app/Services/PragmaticPlayService.php` + `OperatorService.php` extraiu **12 regras de negócio** com identificadores empresariais em formato **BR-[TYPE]-[ENDPOINT]-[CONCERN]-[SEQUENCE]**:

#### Regras Extraídas (Resumo)

| ID | Tipo | Descrição | Endpoints |
|----|------|-----------|-----------|
| BR-GENERIC-ROUTING-VALIDATION-001 | Routing | Dynamic endpoint routing via method resolution | Todos |
| BR-GENERIC-TENANT-EXTRACTION-001 | Auth | Tenant/operator extraction from token | Todos |
| BR-GENERIC-OPERATOR-CACHING-001 | Cache | 1-hour TTL caching de operador | Todos |
| BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | Auth | Token prefix removal antes de enviar ao provider | /balance, /bet, etc |
| BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | Security | MD5 hash generation para autenticação | Todos |
| BR-GENERIC-CREDENTIAL-LOOKUP-001 | Auth | Busca secret-key do operador | Todos |
| BR-GENERIC-PROVIDER-INTEGRATION-001 | Integration | HTTP POST ao provider com tenant URL | Todos |
| BR-GENERIC-ERROR-HANDLING-001 | Error | Captura de endpoint inválido (500) | Todos |
| BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001 | Feature | Suporte a `token` OU `userId` | /balance, /bet |
| BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001 | Order | Ordem de sanitização (token primeiro, userId depois) | /balance |
| BR-GENERIC-RESPONSE-PASSTHROUGH-001 | Response | Response re-passthrough sem transformação | /balance, /result |

#### Artefatos Entregues

**Arquivo:** `/docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` (552 linhas)

Estrutura:
```
- ID da Regra (BR-*)
- Descrição (1 frase)
- Contexto de Negócio (por que existe)
- Escopo (quais endpoints)
- Lógica de Decisão (pseudocódigo)
- Casos Extremos (validações)
- Código Fonte (linhas exatas em PHP)
- Dependências (relação com outras regras)
```

#### Definition of Done ✅

- [x] Todos 9 endpoints analisados
- [x] 12 regras extraídas e documentadas
- [x] Rastreabilidade até código-fonte (números de linha)
- [x] Nomenclatura BR-* aplicada
- [x] Matriz de dependências incluída
- [x] 5 questões abertas documentadas

---

### **Fase 2: DOCUMENT — Documentar Endpoints Técnicos**

**Story:** CASINO-2.2  
**Status:** ✅ COMPLETO

#### O que será feito

Para cada endpoint, documentar o fluxo técnico completo mostrando como as regras (BR-*) interagem.

#### 🔗 Modelo: `/balance` Endpoint

**Arquivo:** `/docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` (439 linhas)

**Estrutura do Documento:**

1. **Resumo Executivo**
   - O que faz: "Consulta saldo atual de um jogador"
   - Quando usado: "Cliente requisita saldo antes/depois de aposta"

2. **Fluxo em 8 Fases** (Mermaid diagram + explicação)
   ```
   FASE 1: ROTEAMENTO
   └─ Valida que endpoint existe (BR-GENERIC-ROUTING-VALIDATION-001)
   
   FASE 2: EXTRAÇÃO TENANT
   └─ Parse token em operator_slug (BR-GENERIC-TENANT-EXTRACTION-001)
   
   FASE 3: LOOKUP OPERADOR
   └─ Query DB + cache 1h (BR-GENERIC-OPERATOR-CACHING-001)
   
   FASE 4: SANITIZAÇÃO
   └─ Remove prefixo tenant (BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001)
   
   FASE 5: LOOKUP CREDENCIAIS
   └─ Busca secret-key (BR-GENERIC-CREDENTIAL-LOOKUP-001)
   
   FASE 6: GERAÇÃO HASH
   └─ MD5(sorted_payload + secret) (BR-GENERIC-AUTHENTICATION-HMAC-MD5-001)
   
   FASE 7: HTTP POST
   └─ POST ao provider (BR-GENERIC-PROVIDER-INTEGRATION-001)
   
   FASE 8: RESPOSTA
   └─ Passthrough sem transformação (BR-GENERIC-RESPONSE-PASSTHROUGH-001)
   ```

3. **Matriz de Regras**
   - Quais 10 regras aplicam a este endpoint
   - Em qual fase cada uma executa
   - Impacto de cada regra

4. **5 Cenários de Erro**
   - Missing token/userId
   - Operator not found
   - Credentials missing
   - Provider timeout
   - Invalid hash
   - Cada um com validação explicada

5. **Exemplo Completo (Request → Response)**
   - Payload real com valores
   - Headers necessários
   - Response esperado

6. **Checklist de Segurança**
   - Tenant isolation validada
   - Hash authentication verificada
   - Operator & credential validation
   - Endpoint validation
   - HTTP method validation

#### Definition of Done

- [x] Fluxo 8-fases documentado
- [x] 10 regras mapeadas ao endpoint
- [x] 5 cenários de erro cobertos
- [x] Mermaid flowchart renderiza corretamente
- [x] Exemplo request/response funcional
- [x] Todos os rules aplicáveis listados
- [x] Security checklist completo
- [x] Pronto como template para os outros 8 endpoints

---

### **Fase 3: TEST — Construir Suite de Testes de Integração**

**Story:** CASINO-2.3  
**Status:** ⏳ PLANEJADO

#### O que será feito

Construir módulo agnóstico de testes de integração que valida cada uma das 12 regras de negócio contra o sistema que está sendo testado (PHP legado ou Go futuro).

#### Arquitetura — Casino Proxy Test Oracle

**Objetivo:** Framework reutilizável para testar qualquer implementação (PHP, Go, ou outra linguagem) contra as regras definidas em CASINO-2.

**Stack:**
- **Linguagem:** Java (JDK 21+)
- **Framework de testes:** JUnit 5
- **Mock de integrações externas:** WireMock 3.x
- **HTTP client:** RestAssured ou HttpClient nativo
- **Assertions:** AssertJ
- **Build:** Maven ou Gradle
- **CI/CD:** GitHub Actions / GitLab CI

#### Estrutura do Projeto

```
casino-proxy-test-oracle/
├── pom.xml                                      # Maven config
├── README.md                                    # Documentação
├── src/main/java/
│   └── com/casino/oracle/
│       ├── client/
│       │   ├── HttpClientFactory.java           # Factory para HTTP client
│       │   └── PayloadBuilder.java              # Builder para payloads
│       ├── mock/
│       │   ├── ProviderMockServer.java          # WireMock setup
│       │   ├── PragmaticPlayMocks.java          # Mocks específicos Pragmatic
│       │   ├── EvolutionGamingMocks.java        # Mocks específicos Evolution
│       │   └── (...outros providers)
│       ├── assertions/
│       │   ├── ResponseAssertions.java          # Custom assertions
│       │   ├── RuleAssertions.java              # Rule-specific assertions
│       │   └── SecurityAssertions.java          # Security checks
│       ├── data/
│       │   ├── Fixtures.java                    # Dados de teste
│       │   ├── PayloadExamples.java             # Exemplos de payloads reais
│       │   └── (...dados por provider)
│       └── config/
│           └── TestConfig.java                  # Config centralizada
│
├── src/test/java/
│   └── com/casino/oracle/
│       ├── integration/
│       │   ├── PragmaticPlayRulesTest.java      # Testa 12 regras vs Pragmatic
│       │   ├── PragmaticPlayEndpointsTest.java  # Testa 9 endpoints do Pragmatic
│       │   ├── EvolutionGamingRulesTest.java    # (próximo provider)
│       │   └── (...testes por provider)
│       └── rules/
│           ├── RoutingValidationTest.java       # BR-GENERIC-ROUTING-VALIDATION-001
│           ├── TenantExtractionTest.java        # BR-GENERIC-TENANT-EXTRACTION-001
│           ├── OperatorCachingTest.java         # BR-GENERIC-OPERATOR-CACHING-001
│           ├── TokenSanitizationTest.java       # BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001
│           ├── AuthenticationTest.java          # BR-GENERIC-AUTHENTICATION-HMAC-MD5-001
│           ├── CredentialLookupTest.java        # BR-GENERIC-CREDENTIAL-LOOKUP-001
│           ├── ProviderIntegrationTest.java     # BR-GENERIC-PROVIDER-INTEGRATION-001
│           ├── ErrorHandlingTest.java           # BR-GENERIC-ERROR-HANDLING-001
│           ├── DualTokenSupportTest.java        # BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001
│           ├── TokenSanitizationOrderTest.java  # BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001
│           └── ResponsePassthroughTest.java     # BR-GENERIC-RESPONSE-PASSTHROUGH-001
│
└── src/main/resources/
    ├── application.properties                    # Config
    ├── wiremock/
    │   ├── pragmatic-play/
    │   │   ├── balance-success-response.json
    │   │   ├── balance-error-response.json
    │   │   └── (...stubs por endpoint)
    │   ├── evolution-gaming/
    │   └── (...stubs por provider)
    └── test-data/
        ├── pragmatic-play-fixtures.json
        └── (...fixtures por provider)
```

#### Como Funciona (Exemplo: /balance endpoint)

```java
@ExtendWith(WireMockExtension.class)
class PragmaticPlayBalanceTest {
    
    private HttpClient client;
    private ProviderMockServer mockServer;
    
    @BeforeEach
    void setup(WireMockRuntimeInfo wmRuntimeInfo) {
        mockServer = new ProviderMockServer(wmRuntimeInfo);
        client = HttpClientFactory.create("http://system-under-test:8080");
        
        // Setup WireMock para simular responses do Pragmatic Play
        mockServer.stubPragmaticPlayBalance(
            Fixtures.PRAGMATIC_BALANCE_SUCCESS_RESPONSE
        );
    }
    
    @Test
    void testBrGenericRoutingValidation001() {
        // BR-GENERIC-ROUTING-VALIDATION-001: Endpoint inválido → 500
        Response response = client.post("/webhooks/pragmatic-play/invalid-endpoint");
        
        assertThat(response.getStatusCode()).isEqualTo(500);
        assertThat(response.getBody()).contains("method_not_found");
    }
    
    @Test
    void testBrPragmaticBalanceTokenSanitization001() {
        // BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001: Sanitiza token
        String payloadWithToken = PayloadBuilder.balance()
            .withToken("operator1_abc123")
            .build();
        
        Response response = client.post("/webhooks/pragmatic-play/balance", payloadWithToken);
        
        // Verifica que a chamada ao mock (Pragmatic Play) recebeu token SEM prefixo
        mockServer.verify(
            postRequestedFor(urlEqualTo("/pragmatic-play/balance.html"))
                .withRequestBody(jsonPath("$.token", equalTo("abc123")))
        );
    }
    
    @Test
    void testBrGenericAuthenticationHmacMd5001() {
        // BR-GENERIC-AUTHENTICATION-HMAC-MD5-001: Hash válido
        String payload = PayloadBuilder.balance()
            .withOperator("operator1")
            .withSecretKey("my-secret")
            .build();
        
        String correctHash = HashGenerator.md5(payload, "my-secret");
        
        Response response = client.post(
            "/webhooks/pragmatic-play/balance",
            payload,
            headers("X-Signature", correctHash)
        );
        
        assertThat(response.getStatusCode()).isEqualTo(200);
    }
    
    @Test
    void testBalanceEndpointFullFlow() {
        // Teste de fluxo completo: Fase 1-8 descritas em balance.md
        String request = Fixtures.PRAGMATIC_BALANCE_SUCCESS_REQUEST;
        
        Response response = client.post("/webhooks/pragmatic-play/balance", request);
        
        // Fase 1: Roteamento ✓
        assertThat(response.getStatusCode()).isNotEqualTo(500);
        
        // Fase 2-4: Tenant extraction e sanitização ✓ (verificado em WireMock)
        mockServer.verify(
            postRequestedFor(urlContaining("/pragmatic-play/balance.html"))
        );
        
        // Fase 5-7: Lookup, hash, HTTP ✓ (mock respondeu)
        
        // Fase 8: Response passthrough ✓
        assertThat(response.getBody()).containsKeys("player_id", "balance");
    }
}
```

#### WireMock Integration (Exemplo)

```java
@Component
class ProviderMockServer {
    private WireMockServer mockServer;
    
    public ProviderMockServer(WireMockRuntimeInfo wmRuntimeInfo) {
        this.mockServer = wmRuntimeInfo.getWireMock();
    }
    
    public void stubPragmaticPlayBalance(String responseBody) {
        mockServer.stubFor(
            post(urlEqualTo("/pragmatic-play/balance.html"))
                .withHeader("Content-Type", containing("application/json"))
                .willReturn(aResponse()
                    .withStatus(200)
                    .withHeader("Content-Type", "application/json")
                    .withBody(responseBody))
        );
    }
    
    public void stubPragmaticPlayBalanceError(String errorCode) {
        mockServer.stubFor(
            post(urlEqualTo("/pragmatic-play/balance.html"))
                .willReturn(aResponse()
                    .withStatus(400)
                    .withBody("{\"error\": \"" + errorCode + "\"}")
                )
        );
    }
    
    public void verify(RequestPatternBuilder pattern) {
        mockServer.verify(pattern);
    }
}
```

#### Agnóstico a Linguagem

**Como funciona:**

1. **Sistema sob teste** pode estar em qualquer linguagem (PHP, Go, Java, Python)
2. **Test Oracle** faz HTTP requests contra qualquer servidor (via configurable `base_url`)
3. **WireMock** simula respostas do provider externo (Pragmatic Play, Evolution, etc.)
4. **Assertions** são agnósticas: validam comportamento, não implementação

**Exemplo — Testando Go:**

```java
// Mesmo código de teste funciona contra Go quando migrar
HttpClient client = HttpClientFactory.create("http://casino-proxy-go:8080");
// ... resto do teste é idêntico
```

#### Definition of Done

- [ ] Estrutura de projeto Maven/Gradle criada em `casino-proxy-test-oracle/`
- [ ] WireMock integrado e configurado
- [ ] 50+ test cases: 1 por regra BR-* + 40+ por endpoint
- [ ] Fixtures realistas para Pragmatic Play (9 endpoints)
- [ ] Testes passam 100% contra PHP legado
- [ ] Documentação README com:
  - Como rodar tests
  - Como adicionar novo provider
  - Como adicionar novo endpoint
- [ ] CI/CD pipeline configurado
- [ ] Testes reutilizáveis para próximas linguagens (Go)

---

### **Fase 4: MATRIX — Criar Matrizes de Rastreamento YAML**

**Story:** CASINO-2.4  
**Status:** ⏳ PLANEJADO

#### O que será feito

Criar matrizes YAML que rastreiam cada regra até:
- OpenAPI spec (CASINO-1)
- Código-fonte PHP
- Teste de validação (Fase 3)
- Implementation Go (CASINO-3)

#### Estrutura de Matrizes

```
docs/casino-proxy/trace-matrices/pragmatic-play-trace-matrix.yaml

rules:
  - id: BR-GENERIC-ROUTING-VALIDATION-001
    name: "Dynamic Endpoint Routing via Method Resolution"
    
    openapi_reference: "casino-proxy-master.yaml#/components/schemas/WebhookRequest"
    
    php_code:
      file: "legacy/casino-proxy/app/Services/PragmaticPlayService.php"
      lines: "18-25"
    
    test_coverage:
      file: "tests/php-oracle/PragmaticPlayControllerTest.php"
      test_ids: ["test_invalid_endpoint_returns_500", "test_valid_endpoint_routes_correctly"]
    
    go_implementation:
      status: "pending"
      file: "services/pragmatic-play/handlers/routing.go"
      lines: "TBD"
    
    severity: "CRITICAL"
    validation_gate: "MUST PASS in Phase 5"
```

#### Artefatos Esperados

```
docs/casino-proxy/trace-matrices/
├── pragmatic-play-trace-matrix.yaml
├── evolution-gaming-trace-matrix.yaml
├── pg-soft-trace-matrix.yaml
└── ... (8 providers totais)
```

#### Definition of Done (esperado)

- [ ] Cada uma das 12 regras rastreada em 4 camadas
- [ ] Referencias exatas (arquivo + linhas)
- [ ] Links bidirecional (spec ↔ código ↔ teste ↔ implementação)
- [ ] YAML válido, sem erros de sintaxe
- [ ] Documentação de como atualizar matrices

---

### **Fase 5: VALIDATE — Validação 100% de Testes PHP**

**Story:** CASINO-2.5  
**Status:** ⏳ PLANEJADO

#### O que será feito

Executar suite completa de testes contra PHP e confirmar que 100% passam. Gate para liberar CASINO-3.

#### Validação Checklist

```
Pragmatic Play Validation Gate
✅ Fase 1: 12 regras extraídas e documentadas
✅ Fase 2: /balance endpoint documentado (template pronto)
⏳ Fase 3: 50+ testes implementados
⏳ Fase 4: Matrizes YAML preenchidas
⏳ Fase 5: phpunit --filter=PragmaticPlay 100% PASS

Gate Decision: GO/NO-GO para CASINO-3
└─ SE 100% testes passam → GO → @po aprova → libera proximo provider
└─ SE falhas encontradas → Fix → Re-run → Loop até GO
```

#### Artefatos Esperados

```
docs/casino-proxy/validation-gates/
├── pragmatic-play-validation-report.md
│  ├─ Test run summary (50 testes = 50 PASS)
│  ├─ Rule coverage checklist
│  ├─ Endpoint coverage checklist
│  ├─ Error scenario coverage
│  └─ Sign-off (✅ @po approves)
└── evolution-gaming-validation-report.md
```

---

### Backlog Completo: CASINO-2 (40 Stories)

**8 provedores × 5 fases = 40 stories**

| Provider | Fase 1 (Extract) | Fase 2 (Document) | Fase 3 (Test) | Fase 4 (Matrix) | Fase 5 (Validate) | Status |
|----------|------------------|------------------|---------------|-----------------|-------------------|--------|
| **Pragmatic Play** | CASINO-2.1 ✅ | CASINO-2.2 ✅ | CASINO-2.3 ⏳ | CASINO-2.4 ⏳ | CASINO-2.5 ⏳ | Em Andamento |
| **Evolution Gaming** | CASINO-2.6 | CASINO-2.7 | CASINO-2.8 | CASINO-2.9 | CASINO-2.10 | Planejado |
| **PG Soft** | CASINO-2.11 | CASINO-2.12 | CASINO-2.13 | CASINO-2.14 | CASINO-2.15 | Planejado |
| **Mancala** | CASINO-2.16 | CASINO-2.17 | CASINO-2.18 | CASINO-2.19 | CASINO-2.20 | Planejado |
| **Digitain** | CASINO-2.21 | CASINO-2.22 | CASINO-2.23 | CASINO-2.24 | CASINO-2.25 | Planejado |
| **Evoplay** | CASINO-2.26 | CASINO-2.27 | CASINO-2.28 | CASINO-2.29 | CASINO-2.30 | Planejado |
| **OpenBox** | CASINO-2.31 | CASINO-2.32 | CASINO-2.33 | CASINO-2.34 | CASINO-2.35 | Planejado |
| **Alternar** | CASINO-2.36 | CASINO-2.37 | CASINO-2.38 | CASINO-2.39 | CASINO-2.40 | Planejado |

---

### Timeline de CASINO-2

```
Semana 1 (2026-05-11 a 2026-05-18)
├─ Pragmatic Play: Fases 1-2 ✅ completas
└─ Pragmatic Play: Fase 3 em andamento

Semana 2 (2026-05-18 a 2026-05-25)
├─ Pragmatic Play: Fases 3-4 finalizadas
├─ Pragmatic Play: Fase 5 ✅ gate passa
└─ @po aprova, começa Evolution Gaming

Semanas 3-6
├─ Evolution Gaming (semana 2)
├─ PG Soft (semana 3)
├─ Mancala (semana 4)
├─ Digitain, Evoplay, OpenBox, Alternar (semanas 5-6)
└─ @po valida cada provider antes de próximo

Final de Semana 6 (2026-06-22)
└─ CASINO-2 100% completo → Desbloqueia CASINO-3
```

---

## 4. CASINO-3 — Go Microservices Implementation ⏸️ AGUARDANDO

### Objetivo
Implementar serviços Go que replicam 100% o comportamento do PHP, com infraestrutura moderna e escalabilidade.

### Status
🔴 **BLOQUEADO por CASINO-2** — Não pode começar implementação sem oráculo de testes PHP

### Estrutura: 4 Fases

| Fase | Stories | O que é Feito |
|------|---------|--------------|
| **Fase 0: IaC** | CASINO-3.0 | Escolhe Terraform/CloudFormation, setup inicial |
| **Fase 2: Architecture** | CASINO-3.1-3.4 | Design microservices + DB schema (PostgreSQL) |
| **Fase 3: Implementation** | CASINO-3.5-3.11 | Implementa 8 serviços Go + gateway + admin API |
| **Fase 4: Migration** | CASINO-3.12-3.16 | Migração dados + dual-write + tráfego gradual |
| **Fase 5: Deploy** | CASINO-3.17-3.22 | Deploy híbrido + cutover final + decommission PHP |

### Como CASINO-2 Tests Guiam CASINO-3

```
CASINO-2 Artefato                  →  CASINO-3 Acceptance Criteria
─────────────────────────────────────────────────────────────────
rules.md (12 regras)               →  Cada regra deve ser implementada identicamente
/balance.md (fluxo 8-fases)        →  Go handler deve executar mesmas 8 fases
phpunit suite (50+ testes)         →  Go implementation roda testes PHP localmente
trace-matrices (rastreamento)      →  Cada regra mapeada de PHP → Go
validation report (100% pass)      →  Go tests devem igualar PHP tests 100%
```

### Definição de Pronto (CASINO-3)

- [ ] Fase 0: IaC testada, deploy bem-sucedido com health checks
- [ ] Fase 2: Arquitetura desenhada, DB schema validado
- [ ] Fase 3: Go services passam em CASINO-2 tests (parity 100%)
- [ ] Fase 4: Migração de dados zero-loss, dual-write funcionando
- [ ] Fase 5: Cutover completo, PHP decommissioned, performance ≥ PHP
- [ ] SLA 99.9% uptime atingido em produção

---

## 5. Template — Como Replicar para Outros Providers

### Checklist de 5 Fases (Reutilizável)

Use este checklist para cada um dos 7 providers restantes (Evolution Gaming até Alternar).

#### Fase 1: EXTRACT

**Entrada:** Código PHP handler do provider  
**Saída:** Documento `phase-1-business-rules/{provider}-rules.md`

```markdown
# Regras de Lógica de Negócio — {PROVIDER}

## Sumário
- Total de endpoints: X
- Total de regras extraídas: Y
- Nomenclatura: BR-[TYPE]-[ENDPOINT]-[CONCERN]-[SEQ]

## Regras Extraídas
[Lista de 10-15 regras com ID, descrição, contexto, lógica, código-fonte]

## Matriz de Dependências
[Quais regras dependem de quais]
```

#### Fase 2: DOCUMENT

**Entrada:** Rules extraídas em Fase 1  
**Saída:** Documento `phase-2-technical-documentation/{provider}-{endpoint}.md`

```markdown
# Endpoint Documentation: /{ENDPOINT}

## Resumo
[O que faz, quando usado]

## Fluxo em N Fases
[Mermaid diagram + 8-fase breakdown]

## Regras Aplicáveis
[Tabela: qual regra em qual fase, impacto]

## Cenários de Erro
[5+ casos com validação]

## Exemplo Completo
[Request → Response real]

## Security Checklist
[Validações de segurança]
```

**Para Pragmatic Play:** `/balance`, `/bet`, `/authenticate`, ... (9 endpoints)

#### Fase 3: TEST

**Entrada:** Documentação + Rules  
**Saída:** Módulo Java `casino-proxy-test-oracle/` com testes agnósticos

Integra-se com o framework Java `casino-proxy-test-oracle/`:
- Testes reutilizáveis contra PHP legado
- Mesmos testes funcionam contra Go futuro (sem mudança de código)
- WireMock simula respostas do provider externo

#### Fase 4: MATRIX

**Entrada:** Rules + Code + Tests  
**Saída:** `trace-matrices/{provider}-trace-matrix.yaml`

```yaml
rules:
  - id: BR-GENERIC-ROUTING-VALIDATION-001
    openapi_reference: "casino-proxy-master.yaml#/paths/~1{endpoint}"
    php_code:
      file: "legacy/casino-proxy/app/Services/..."
      lines: "X-Y"
    test_coverage:
      file: "tests/php-oracle/..."
      test_ids: ["test_..."]
    go_implementation:
      status: "pending"
      file: "services/.../..."
      lines: "TBD"
```

#### Fase 5: VALIDATE

**Entrada:** Suite completa de testes Java + Sistema sob teste (PHP ou Go)  
**Saída:** Relatório de validação + Gate GO/NO-GO

```bash
# Executar suite completa de testes
$ cd casino-proxy-test-oracle/
$ mvn clean test -Dtest=PragmaticPlayRulesTest,PragmaticPlayEndpointsTest

# Resultado esperado
BUILD SUCCESS (50+ tests)

# Se 100% pass → Gate: GO → @po aprova → Libera próximo provider
# Se falhas → Fix → Re-run → Loop até GO
```

---

## 6. Timeline Integrada (Visão Completa)

### Fluxo de Execução

```
Ciclo CASINO-2: Pragmatic Play → Evolution Gaming → PG Soft → (Mancala, Digitain, Evoplay, OpenBox, Alternar)

Pragmatic Play (Provider 1)
├─ CASINO-2.1: Extract ✅
├─ CASINO-2.2: Document ✅
├─ CASINO-2.3: Test (em andamento)
├─ CASINO-2.4: Matrix (planejado)
├─ CASINO-2.5: Validate (planejado)
└─ Gate: GO → @po aprova → próximo provider

Evolution Gaming (Provider 2)
├─ CASINO-2.6: Extract (espera aprovação Pragmatic)
├─ CASINO-2.7: Document
├─ CASINO-2.8: Test
├─ CASINO-2.9: Matrix
├─ CASINO-2.10: Validate
└─ Gate: GO → @po aprova → próximo provider

(Repetir para PG Soft, Mancala, Digitain, Evoplay, OpenBox, Alternar)

Final: CASINO-2 100% Completo
└─ Desbloqueia CASINO-3: Implementação Go em paralelo

CASINO-3 Timeline
├─ Fase 0: IaC design + setup (em paralelo com CASINO-2)
├─ Fase 2: Architecture microservices + DB schema
├─ Fase 3: Implementation (8 serviços Go + gateway + admin)
├─ Fase 4: Database migration + dual-write
└─ Fase 5: Deployment híbrido + cutover final → PHP decommissioned
```

---

## 7. Métricas de Sucesso

### CASINO-1 ✅
- [x] 8 provedores documentados
- [x] 100% endpoints em OpenAPI
- [x] 0 gaps entre code e spec

### CASINO-2 (em progresso)
- [ ] 40 stories completas
- [ ] 100% regras de negócio documentadas
- [ ] 200+ testes passando (50+ por provider média)
- [ ] 0 regras perdidas na tradução PHP → Go
- [ ] PO approves cada provider antes de próximo

### CASINO-3 (aguardando)
- [ ] 8 serviços Go em produção
- [ ] Go tests = PHP tests (100% parity)
- [ ] 0 downtime durante migração
- [ ] Performance ≥ PHP baseline
- [ ] Infrastructure 100% via IaC
- [ ] 99.9% uptime SLA
- [ ] PHP completamente decommissioned

---

## 8. Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|--------|-----------|
| Provider toma mais tempo que estimado | 🟡 Médio | 🟡 Atraso cascata | Fases paralelas CASINO-3, buffer de 1 semana |
| Teste PHP descobre gaps em código | 🟡 Médio | 🟡 Rework necessário | Code review rigoroso em Fase 3 |
| Go implementation encontra edge cases | 🟡 Médio | 🟡 Delays | CASINO-2 tests cobrem 99% dos casos |
| Infraestrutura de IaC tem issue | 🔴 Baixo | 🔴 Crítico | Fase 0 de CASINO-3 com expert review |
| Performance Go < PHP | 🔴 Muito Baixo | 🔴 Crítico | Profiling + optimization em Fase 3 |

---

## 9. Glossário

| Termo | Definição |
|-------|-----------|
| **BR-*** | Business Rule ID (formato empresarial) |
| **Oracle** | Suite de testes PHP que valida paridade Go |
| **Parity** | Comportamento idêntico entre PHP e Go |
| **Trace Matrix** | Documento YAML que rastreia regra → spec → código → teste |
| **Gate** | Critério de saída de uma fase (ex: 100% testes passam) |
| **Provider** | Provedor de jogos (Pragmatic Play, Evolution, PG Soft, etc) |
| **Tenant** | Isolamento multi-tenant via prefixo em token |
| **IaC** | Infrastructure as Code (Terraform / CloudFormation) |

---

## 10. Documentos Relacionados

- **Epic Reorganization Plan:** `docs/epics/casino-proxy/EPIC-REORGANIZATION-PLAN.md`
- **Phase 1 Rules (Pragmatic Play):** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`
- **Phase 2 Documentation (/balance):** `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md`
- **CASINO-1 Master Spec:** `docs/casino-proxy/openapi/casino-proxy-master.yaml`

---

**Última Atualização:** 2026-05-11  
**Próxima Revisão:** 2026-05-18 (após Pragmatic Play Fase 2 completa)  
**Proprietário:** @dev (@architect para review)  
**Status:** 🚀 EM EXECUÇÃO
