# CASINO-2.3: Pragmatic Play — Fase 3 Test Oracle

**Epic ID:** CASINO-2.3  
**Tipo:** Fase 3 de 5 — Test Oracle (Pragmatic Play)  
**Status:** 🟡 Draft — Stories pendentes  
**Criado:** 2026-05-15  
**Atualizado:** 2026-05-15  
**Owner:** @sm (River)  
**Executor:** @dev (Dex)  
**Validação:** @po (Pax)  
**Depende De:** CASINO-2.2 (Fase 2 — Technical Documentation) ✅ Done  
**Bloqueia:** CASINO-2.4 (Fase 4 — Trace Matrix) / CASINO-2.5 (Fase 5 — Validation Gate)

---

## Objetivo

Construir o **Casino Proxy Test Oracle** — um framework de testes Java agnóstico de implementação que valida as 12 regras de negócio (BR-*) do Pragmatic Play contra qualquer sistema (PHP legado ou Go futuro) via HTTP, sem dependência de recursos externos reais.

```
Input:   docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md (12 regras BR-*)
         docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-*.md (9 endpoints)
Output:  casino-proxy-test-oracle/  (projeto Java — 50+ testes, WireMock stubs, CI/CD)
```

---

## Contexto

### Por que esta Fase?

A Fase 2 documentou *como o PHP faz*. A Fase 3 **prova** que o sistema (PHP hoje, Go amanhã) realmente faz isso. Sem o Test Oracle, a migração PHP → Go não tem critério de aceitação — é impossível garantir parity. Os mesmos testes que validam o PHP legado vão validar o Go sem nenhuma modificação de código.

### Princípio Core: Agnóstico de Implementação

Os testes validam **comportamento HTTP observável**, não linguagem de programação:

```
Sistema sob teste (PHP ou Go)  ←── HTTP requests ──  Test Oracle (Java)
         ↓                                                    ↑
   HTTP responses ───────────────────────────────────────────┘
         ↓
WireMock intercepta chamadas ao provider externo (Pragmatic Play)
retorna stubs pré-configurados
```

**Portabilidade:** Mudar de PHP para Go = apenas alterar `base_url` na config. Zero mudança nos testes.

---

## Stack Técnica

| Componente | Tecnologia | Versão |
|------------|-----------|--------|
| Linguagem | Java | JDK 21+ |
| Framework de testes | JUnit | 5.x |
| Mock de integrações externas | WireMock | 3.x |
| HTTP client | RestAssured | 5.x |
| Assertions | AssertJ | 3.x |
| Build | Maven | 3.9+ |
| CI/CD | GitHub Actions | — |

---

## Estrutura do Projeto

```
casino-proxy-test-oracle/
├── pom.xml
├── README.md
├── src/main/java/com/casino/oracle/
│   ├── client/
│   │   ├── HttpClientFactory.java
│   │   └── PayloadBuilder.java
│   ├── mock/
│   │   ├── ProviderMockServer.java
│   │   └── PragmaticPlayMocks.java
│   ├── assertions/
│   │   ├── ResponseAssertions.java
│   │   ├── RuleAssertions.java
│   │   └── SecurityAssertions.java
│   ├── data/
│   │   ├── Fixtures.java
│   │   └── PragmaticPlayFixtures.java
│   └── config/
│       └── TestConfig.java
├── src/test/java/com/casino/oracle/
│   ├── rules/
│   │   ├── RoutingValidationTest.java          # BR-GENERIC-ROUTING-VALIDATION-001
│   │   ├── TenantExtractionTest.java           # BR-GENERIC-TENANT-EXTRACTION-001
│   │   ├── OperatorCachingTest.java            # BR-GENERIC-OPERATOR-CACHING-001
│   │   ├── TokenSanitizationTest.java          # BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001
│   │   ├── AuthenticationHmacTest.java         # BR-GENERIC-AUTHENTICATION-HMAC-MD5-001
│   │   ├── CredentialLookupTest.java           # BR-GENERIC-CREDENTIAL-LOOKUP-001
│   │   ├── ProviderIntegrationTest.java        # BR-GENERIC-PROVIDER-INTEGRATION-001
│   │   ├── ErrorHandlingTest.java              # BR-GENERIC-ERROR-HANDLING-001
│   │   ├── DualTokenSupportTest.java           # BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001
│   │   ├── TokenSanitizationOrderTest.java     # BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001
│   │   ├── ResponsePassthroughTest.java        # BR-GENERIC-RESPONSE-PASSTHROUGH-001
│   │   └── AuthenticateTransformTest.java      # BR-PRAGMATIC-AUTHENTICATE-USERID-REPREFIX-001
│   └── integration/
│       ├── AuthenticateEndpointTest.java       # /authenticate — 5+ cenários
│       ├── BalanceEndpointTest.java            # /balance — 5+ cenários
│       ├── BetEndpointTest.java                # /bet — 5+ cenários
│       ├── RefundEndpointTest.java             # /refund — 4+ cenários
│       ├── ResultEndpointTest.java             # /result — 4+ cenários
│       ├── BonusWinEndpointTest.java           # /bonusWin — 4+ cenários
│       ├── JackpotWinEndpointTest.java         # /jackpotWin — 4+ cenários
│       ├── PromoWinEndpointTest.java           # /promoWin — 4+ cenários
│       └── AdjustmentEndpointTest.java        # /adjustment — 4+ cenários
└── src/main/resources/
    ├── application.properties
    └── wiremock/pragmatic-play/
        ├── authenticate-success.json
        ├── authenticate-error.json
        ├── balance-success.json
        ├── balance-error.json
        ├── bet-success.json
        ├── bet-error.json
        ├── refund-success.json
        ├── result-success.json
        ├── bonuswin-success.json
        ├── jackpotwin-success.json
        ├── promowin-success.json
        └── adjustment-success.json
```

---

## Cobertura de Testes

### 12 Regras BR-* (1 classe de teste por regra)

| Regra | Classe | Cenários mínimos |
|-------|--------|-----------------|
| BR-GENERIC-ROUTING-VALIDATION-001 | `RoutingValidationTest` | endpoint inválido → 500, endpoint válido → roteado |
| BR-GENERIC-TENANT-EXTRACTION-001 | `TenantExtractionTest` | token com prefixo correto, token malformado |
| BR-GENERIC-OPERATOR-CACHING-001 | `OperatorCachingTest` | operador existente, operador não encontrado |
| BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | `TokenSanitizationTest` | prefixo removido antes de enviar ao provider |
| BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | `AuthenticationHmacTest` | MD5 correto, MD5 inválido → 403 |
| BR-GENERIC-CREDENTIAL-LOOKUP-001 | `CredentialLookupTest` | credencial existente, credencial ausente |
| BR-GENERIC-PROVIDER-INTEGRATION-001 | `ProviderIntegrationTest` | POST correto ao provider, timeout |
| BR-GENERIC-ERROR-HANDLING-001 | `ErrorHandlingTest` | endpoint inválido retorna 500 |
| BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001 | `DualTokenSupportTest` | via `token`, via `userId`, ambos ausentes |
| BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001 | `TokenSanitizationOrderTest` | `token` processado antes de `userId` |
| BR-GENERIC-RESPONSE-PASSTHROUGH-001 | `ResponsePassthroughTest` | response do provider chega sem transformação |
| BR-PRAGMATIC-AUTHENTICATE-USERID-REPREFIX-001 | `AuthenticateTransformTest` | `userId` re-prefixado quando `error==0`, não alterado quando `error!=0` |

### 9 Endpoints (cenários de integração por endpoint)

| Endpoint | Classe | Cenários mínimos |
|----------|--------|-----------------|
| `/authenticate` | `AuthenticateEndpointTest` | success + transform, error sem transform, token inválido, operador não encontrado, hash inválido |
| `/balance` | `BalanceEndpointTest` | via token, via userId, ambos ausentes, operador não encontrado, provider timeout |
| `/bet` | `BetEndpointTest` | success passthrough, hash inválido, operador não encontrado, provider erro |
| `/refund` | `RefundEndpointTest` | success passthrough, transação não encontrada, hash inválido, operador não encontrado |
| `/result` | `ResultEndpointTest` | success passthrough (handleResult), hash inválido, operador não encontrado, provider erro |
| `/bonusWin` | `BonusWinEndpointTest` | success passthrough, hash inválido, operador não encontrado, provider erro |
| `/jackpotWin` | `JackpotWinEndpointTest` | success passthrough (high-value audit), hash inválido, operador não encontrado |
| `/promoWin` | `PromoWinEndpointTest` | success passthrough, hash inválido, operador não encontrado, provider erro |
| `/adjustment` | `AdjustmentEndpointTest` | success passthrough (admin-initiated), hash inválido, operador não encontrado, provider erro |

**Total esperado: 50+ test cases**

---

## Stories — Planejadas

### Visão Geral

| Métrica | Valor |
|---------|-------|
| Total de stories | 5 |
| Stories criadas | 0 / 5 ⏳ |
| Stories validadas (@po) | 0 / 5 ⏳ |
| Projeto Test Oracle criado | ❌ Pendente |
| Fase 3 completa | ❌ Pendente |

---

### Story 1 — Setup do Projeto Test Oracle

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.3-setup |
| **Arquivo Story** | `docs/stories/CASINO-2.3-setup-test-oracle.md` |
| **Entregável** | Projeto Maven em `casino-proxy-test-oracle/` com WireMock, JUnit 5, RestAssured, AssertJ configurados e health-check test verde |
| **Status Story** | ⏳ A criar |
| **Estimativa @dev** | 3-4 horas |
| **Dependência** | CASINO-2.2 ✅ |

---

### Story 2 — Testes das Regras BR-GENERIC-* (7 regras)

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.3-generic-rules |
| **Arquivo Story** | `docs/stories/CASINO-2.3-generic-rules-tests.md` |
| **Entregável** | 7 classes de teste cobrindo todas as regras BR-GENERIC-* com stubs WireMock correspondentes |
| **Status Story** | ⏳ A criar |
| **Estimativa @dev** | 4-6 horas |
| **Dependência** | CASINO-2.3-setup |

---

### Story 3 — Testes das Regras BR-PRAGMATIC-* (4 regras + 1 nova)

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.3-pragmatic-rules |
| **Arquivo Story** | `docs/stories/CASINO-2.3-pragmatic-rules-tests.md` |
| **Entregável** | 5 classes de teste cobrindo regras exclusivas do Pragmatic Play, incluindo dual-token e authenticate transform |
| **Status Story** | ⏳ A criar |
| **Estimativa @dev** | 3-4 horas |
| **Dependência** | CASINO-2.3-generic-rules |

---

### Story 4 — Testes de Integração dos 9 Endpoints

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.3-endpoint-tests |
| **Arquivo Story** | `docs/stories/CASINO-2.3-endpoint-tests.md` |
| **Entregável** | 9 classes de integração + 12 WireMock stubs JSON cobrindo todos os endpoints com cenários happy path + error path |
| **Status Story** | ⏳ A criar |
| **Estimativa @dev** | 6-8 horas |
| **Dependência** | CASINO-2.3-pragmatic-rules |

---

### Story 5 — CI/CD Pipeline e Documentação README

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.3-ci-cd |
| **Arquivo Story** | `docs/stories/CASINO-2.3-ci-cd.md` |
| **Entregável** | `.github/workflows/test-oracle.yml` + `casino-proxy-test-oracle/README.md` com instruções de execução, adição de providers e adição de endpoints |
| **Status Story** | ⏳ A criar |
| **Estimativa @dev** | 2-3 horas |
| **Dependência** | CASINO-2.3-endpoint-tests |

---

## Kanban de Execução

```
STORIES (Backlog)       READY               IN PROGRESS         DONE
────────────────────    ─────────────────   ────────────────    ────────────────
setup                   —                   —                   —
generic-rules           —                   —                   —
pragmatic-rules         —                   —                   —
endpoint-tests          —                   —                   —
ci-cd                   —                   —                   —
```

### Ordem de Implementação (@dev)

```
1. setup           (estrutura Maven + WireMock base; 3-4h)
2. generic-rules   (7 regras genéricas com stubs; 4-6h)
3. pragmatic-rules (5 regras Pragmatic-específicas; 3-4h)
4. endpoint-tests  (9 endpoints + 12 stubs JSON; 6-8h)
5. ci-cd           (GitHub Actions + README; 2-3h)
```

**Total estimado:** 18-25 horas de implementação @dev

---

## Artefatos

### Input (já existem)

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | 12 regras BR-* com rastreabilidade PHP | ✅ Completo |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` | Template canônico fluxo 8-fases | ✅ Completo |
| `docs/stories/CASINO-2.2-*.md` (8 stories) | Documentação técnica por endpoint | ✅ Ready |

### Output (a criar pelo @dev)

| Artefato | Localização | Status |
|----------|-------------|--------|
| Projeto Maven base | `casino-proxy-test-oracle/pom.xml` | ⏳ Pendente |
| 12 classes de teste BR-* | `src/test/java/.../rules/` | ⏳ Pendente |
| 9 classes de integração | `src/test/java/.../integration/` | ⏳ Pendente |
| 12 WireMock stubs | `src/main/resources/wiremock/pragmatic-play/` | ⏳ Pendente |
| GitHub Actions workflow | `.github/workflows/test-oracle.yml` | ⏳ Pendente |
| README | `casino-proxy-test-oracle/README.md` | ⏳ Pendente |

---

## Definition of Done — Epic CASINO-2.3

- [ ] Projeto Maven `casino-proxy-test-oracle/` criado e buildando com `mvn clean test`
- [ ] 12 classes de teste cobrindo todas as regras BR-* (1 por regra)
- [ ] 9 classes de teste de integração cobrindo todos os endpoints
- [ ] 50+ test cases no total
- [ ] Todos os stubs WireMock configurados (12 arquivos JSON)
- [ ] Testes são agnósticos de implementação (`base_url` configurável via `application.properties`)
- [ ] README documenta: como rodar, como adicionar provider, como adicionar endpoint
- [ ] GitHub Actions pipeline configurado e passando
- [ ] @po revisa e aprova antes de CASINO-2.4 iniciar
- [ ] CASINO-2.4 (Trace Matrix) desbloqueado

---

## Riscos

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| PHP não acessível para confirmar comportamento real | Alta | Médio | Regras BR-* e docs Phase 2 são suficientes para testes; confirmar edge cases em CASINO-2.5 |
| WireMock stubs divergem do comportamento real do provider | Média | Alto | Testes devem focar em comportamento do proxy, não do provider; stubs simulam respostas mínimas necessárias |
| Escopo creep (adicionar mais providers antes de concluir PP) | Baixa | Alto | Epic limita escopo a Pragmatic Play; Evolution Gaming aguarda CASINO-2.4-2.5 |
| Complexidade do setup Java em ambiente Windows | Baixa | Médio | `pom.xml` especifica JDK 21+; @dev documenta setup no README |

---

## Próximas Fases (após CASINO-2.3)

| Epic | Fase | O que é | Desbloqueado por |
|------|------|---------|-----------------|
| CASINO-2.4 | Fase 4 — Trace Matrix | YAML rastreando cada BR-* por 4 camadas (spec → PHP → teste → Go) | CASINO-2.3 ✅ |
| CASINO-2.5 | Fase 5 — Validation Gate | Executar suite completa contra PHP legado, 100% pass = GO para Evolution Gaming | CASINO-2.4 ✅ |

---

## Histórico de Execução

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Epic CASINO-2.3 criado — 5 stories planejadas |
