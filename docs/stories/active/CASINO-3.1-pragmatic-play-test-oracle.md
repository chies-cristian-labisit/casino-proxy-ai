# CASINO-3.1: Construir Test Oracle Java — Pragmatic Play

**Story ID:** CASINO-3.1  
**Epic:** CASINO-3 (Test Oracle — Business Rules Validation Suite)  
**Tipo:** Fase 3 de 5 — Test Oracle Implementation  
**Status:** Draft  
**Prioridade:** Alta (primeiro oracle — define padrão para os demais providers)  
**Atribuído a:** @dev (Java + WireMock)  
**Relacionado:** CASINO-2.2 (9 endpoint docs ✅ Done), CASINO-1.7 (12 regras BR-* ✅ Done)  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Construir a suite de testes Java (JUnit 5 + WireMock) que valida 100% das regras de negócio do Pragmatic Play contra o PHP legado. Esta é a **primeira suite do projeto** — define a arquitetura e os padrões que todos os outros providers replicarão.

**Objetivo:** Produzir `casino-proxy-test-oracle/pragmatic-play/` com 50+ testes cobrindo as 12 regras BR-* e os 9 endpoints, todos passando 100% contra o PHP.

---

## Contexto

### Por que esta Story?

A Fase 2 (CASINO-2.2) documentou *como o PHP faz* para cada um dos 9 endpoints do Pragmatic Play. Esta story **prova** esse comportamento via testes automáticos:

1. **Rastreabilidade:** Cada regra BR-* tem pelo menos 1 teste correspondente
2. **Critério de aceitação para CASINO-4:** A implementação Go deve passar estes mesmos testes
3. **Regressão:** Detecta mudanças inadvertidas no PHP durante a migração

### Como se Encaixa no Plano

```
Fase 1: Extrair regras ✅ (CASINO-1.7 — 12 regras BR-* documentadas)
Fase 2: Documentar endpoints ✅ (CASINO-2.2 — 9 endpoints documentados)
Fase 3: Test Oracle ← VOCÊ ESTÁ AQUI
  ├─ CASINO-3.1: Construir suite JUnit 5 + WireMock
  ├─ CASINO-3.2: Criar trace matrix YAML
  └─ CASINO-3.3: Validação gate (100% PHP PASS + @po aprova)
Fase 4 (CASINO-4): Go implementation usa oracle como acceptance criteria
```

### Arquitetura do Oracle

```
casino-proxy-test-oracle/
└─ pragmatic-play/
   ├─ pom.xml                              # Maven: JUnit 5 + WireMock + RestAssured
   ├─ src/
   │  ├─ test/java/
   │  │  └─ com/casino/oracle/pragmatic/
   │  │     ├─ PragmaticPlayOracleTest.java # Suite principal (50+ testes)
   │  │     ├─ rules/                       # 1 classe por regra BR-*
   │  │     └─ endpoints/                   # 1 classe por endpoint
   │  └─ test/resources/
   │     └─ wiremock/pragmatic-play/
   │        ├─ authenticate/                # Stubs JSON para /authenticate
   │        ├─ bet/                         # Stubs para /bet
   │        └─ ... (9 endpoints)
   └─ README.md                            # Como rodar localmente
```

### Regras a Cobrir (BR-*)

| Regra | Tipo | Endpoints Afetados |
|-------|------|-------------------|
| BR-GENERIC-ROUTING-VALIDATION-001 | Genérica | Todos (9) |
| BR-GENERIC-TENANT-EXTRACTION-001 | Genérica | Todos (9) |
| BR-GENERIC-OPERATOR-CACHING-001 | Genérica | Todos (9) |
| BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | PP | Todos (9) |
| BR-PRAGMATIC-PARAMETER-ORDER-001 | PP | Todos (9) |
| BR-GENERIC-CREDENTIAL-LOOKUP-001 | Genérica | Todos (9) |
| BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | Genérica | Todos (9) |
| BR-GENERIC-PROVIDER-INTEGRATION-001 | Genérica | Todos (9) |
| BR-GENERIC-RESPONSE-PASSTHROUGH-001 | Genérica | 8/9 (exceto authenticate) |
| BR-PRAGMATIC-BALANCE-DUAL-TOKEN-001 | PP | /balance apenas |
| BR-PRAGMATIC-AUTH-USERID-REPREFIX-001 (PP-007) | PP | /authenticate apenas |
| BR-PRAGMATIC-AUTH-TRANSFORM-EXCLUSIVE-001 (PP-012) | PP | /authenticate apenas |

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Módulo Maven criado em `casino-proxy-test-oracle/pragmatic-play/` com `pom.xml` válido (JUnit 5.x + WireMock 3.x + RestAssured)
- [ ] **AC-2:** WireMock stubs JSON criados para cada um dos 9 endpoints — cobrindo cenários de sucesso e pelo menos 3 cenários de erro por endpoint
- [ ] **AC-3:** Suite `PragmaticPlayOracleTest.java` com mínimo 50 testes — pelo menos 1 teste por regra BR-* + pelo menos 3 testes por endpoint
- [ ] **AC-4:** Testes executáveis localmente contra PHP: `mvn clean test -Dphp.base.url=http://localhost:8000`
- [ ] **AC-5:** 100% dos testes passam contra PHP legado (0 failures, 0 errors)
- [ ] **AC-6:** `README.md` documenta: pré-requisitos, como configurar WireMock, como apontar para PHP, como rodar suite completa
- [ ] **AC-7:** Caso especial authenticate coberto — teste verifica re-prefixação de `userId` quando `error==0` E ausência de re-prefixação quando `error!=0`
- [ ] **AC-8:** Arquivo de saída listado na File List desta story

### Deveria Ter

- [ ] **AC-9:** Testes organizados por grupo de endpoint (Sessão, Consulta, Transação inline, handleResult family)
- [ ] **AC-10:** Configuração parametrizável via `application-test.properties` — URL base do PHP, tenant, operator slug

### Fora do Escopo

- ❌ Testes contra Go (CASINO-4)
- ❌ Trace Matrix YAML (CASINO-3.2)
- ❌ Validation report / gate (CASINO-3.3)
- ❌ Outros providers (Evolution, PG Soft, etc.) — stories separadas

---

## Detalhes Técnicos

### Stack

| Dependência | Versão | Papel |
|-------------|--------|-------|
| JUnit 5 | 5.10+ | Framework de testes |
| WireMock | 3.x | Mock do provider externo (Pragmatic Play) |
| RestAssured | 5.x | Assertions HTTP |
| Maven | 3.9+ | Build tool |

### Exemplo de Teste (AC-7 — authenticate re-prefix)

```java
@Test
void authenticate_successResponse_reprefixesUserId() {
    // Given: WireMock stub retorna userId "12345" com error=0
    stubFor(post("/authenticate.html")
        .willReturn(okJson("{\"userId\":\"12345\",\"error\":0}")));

    // When: chamamos o PHP proxy /authenticate
    Response response = given()
        .formParam("token", "test-token")
        .post(phpBaseUrl + "/v1/webhooks/pragmatic-play/authenticate");

    // Then: userId deve estar re-prefixado
    assertThat(response.jsonPath().getString("userId"))
        .isEqualTo("myoperator_12345");
}

@Test
void authenticate_errorResponse_doesNotReprefixUserId() {
    // Given: error!=0
    stubFor(post("/authenticate.html")
        .willReturn(okJson("{\"userId\":\"12345\",\"error\":1}")));

    // When/Then: userId NÃO é re-prefixado
    Response response = given()
        .formParam("token", "test-token")
        .post(phpBaseUrl + "/v1/webhooks/pragmatic-play/authenticate");

    assertThat(response.jsonPath().getString("userId"))
        .isEqualTo("12345");
}
```

---

## File List

| Arquivo | Ação | Status |
|---------|------|--------|
| `casino-proxy-test-oracle/pragmatic-play/pom.xml` | Criar | ⏳ |
| `casino-proxy-test-oracle/pragmatic-play/src/test/java/.../PragmaticPlayOracleTest.java` | Criar | ⏳ |
| `casino-proxy-test-oracle/pragmatic-play/src/test/resources/wiremock/pragmatic-play/` | Criar | ⏳ |
| `casino-proxy-test-oracle/pragmatic-play/README.md` | Criar | ⏳ |

---

## Notas de Implementação

- **Abordagem recomendada:** Implementar por grupo de endpoint (authenticate primeiro — mais complexo; depois bet como padrão canônico; depois handleResult family em batch)
- **WireMock:** Rodar WireMock em modo standalone ou embedded no JUnit — decidir baseado em ambiente CI disponível
- **PHP env:** Oracle requer PHP legado rodando localmente ou em staging — coordenar com @devops

---

## Estimativa

**8–12 horas** (@dev Java):
- Setup Maven/WireMock/JUnit: 1–2h
- Stubs WireMock (9 endpoints × cenários): 2–3h
- Testes authenticate (mais complexo): 2h
- Testes demais endpoints: 3–4h
- README + execução final: 1h

---

## Dependências & Bloqueadores

**Precisa:**
- CASINO-2.2 (9 endpoint docs) ✅ Done — referência para construir stubs e testes
- CASINO-1.7 (12 regras BR-*) ✅ Done — referência para AC por regra
- Acesso ao PHP legado (ambiente staging ou local)

**Bloqueia:**
- CASINO-3.2 (trace matrix não pode ter `test_coverage` sem os test IDs desta story)
- CASINO-3.3 (validation gate)

---

## Change Log

| Data | Agente | Alteração |
|------|--------|-----------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
