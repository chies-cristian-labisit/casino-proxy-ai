# CASINO-3: Test Oracle — Business Rules Validation Suite

**Epic ID:** CASINO-3  
**Tipo:** Fase 3 — Test Oracle (Language-Agnostic Validation)  
**Status:** 📝 Draft — aguardando refinamento  
**Criado:** 2026-05-15  
**Atualizado:** 2026-05-15  
**Owner:** @sm (River)  
**Executor:** @dev (Java/WireMock)  
**Validação:** @po (Pax) — gate por provider  
**Depende De:** CASINO-2 (Fase 2 — Technical Documentation) — por provider  
**Bloqueia:** CASINO-4 (Go Microservices Implementation)

---

## Objetivo

Construir um **Test Oracle agnóstico de linguagem** — uma suite Java (JUnit 5 + WireMock) que valida 100% das regras de negócio (BR-*) de cada provider contra o PHP legado. Os mesmos testes rodam sem modificação contra a implementação Go (CASINO-4), provando paridade comportamental.

```
Input:   docs/architecture/casino-proxy/phase-2-technical-documentation/{provider}-{endpoint}.md
         docs/architecture/casino-proxy/phase-1-business-rules/{provider}-rules.md
Output:  casino-proxy-test-oracle/{provider}/ (JUnit 5 + WireMock suite)
         docs/architecture/casino-proxy/trace-matrices/{provider}-trace-matrix.yaml
         docs/architecture/casino-proxy/validation-gates/{provider}-validation-report.md
```

---

## Contexto

### Por que este Epic?

A Fase 2 (CASINO-2) documenta *como o PHP faz*. O CASINO-3 **prova** que o comportamento está correto — e cria o critério de aceitação que o CASINO-4 (Go) deve satisfazer. Sem o oracle:
- O time Go não tem como verificar paridade com o PHP
- Edge cases documentados em Fase 2 ficam sem cobertura de teste
- A migração depende de testes manuais frágeis

### Estrutura dos Outputs por Provider

1. **Suite Java** (`casino-proxy-test-oracle/{provider}/`)
   - `src/test/java/.../{Provider}OracleTest.java` — testes JUnit 5
   - `src/test/resources/wiremock/{provider}/` — stubs para cada endpoint
   - `pom.xml` — dependências Maven

2. **Trace Matrix** (`docs/architecture/casino-proxy/trace-matrices/{provider}-trace-matrix.yaml`)
   - Mapeia cada BR-* rule → spec OpenAPI → código PHP (arquivo + linhas) → test IDs → Go impl (TBD)

3. **Validation Report** (`docs/architecture/casino-proxy/validation-gates/{provider}-validation-report.md`)
   - Resultado do `mvn clean test`: X tests, 0 failures
   - Checklist de cobertura de regras e endpoints
   - Sign-off @po para liberar próximo provider

### Por que Java (e não Go)?

O oracle é construído em Java para ser **agnóstico de linguagem**:
- Testa via HTTP — funciona contra qualquer implementação (PHP, Go, futuro)
- WireMock simula responses do provider externo de forma determinística
- JUnit 5 + Maven são o padrão do time de QA

---

## 3 Fases por Provider

| Fase | Story Pattern | O que Produz | Dependência |
|------|--------------|--------------|-------------|
| **Fase 3: Test** | CASINO-3.{N} | Suite JUnit 5 + WireMock stubs | CASINO-2 Fase 2 do provider ✅ |
| **Fase 4: Matrix** | CASINO-3.{N+1} | Trace matrix YAML (12+ regras mapeadas) | Fase 3 do provider |
| **Fase 5: Validate** | CASINO-3.{N+2} | Validation report + @po gate | Fases 3+4 do provider |

---

## Backlog Completo — 24 Stories

| Provider | Fase 3 (Test) | Fase 4 (Matrix) | Fase 5 (Validate) | Status |
|----------|--------------|-----------------|-------------------|--------|
| **Pragmatic Play** | CASINO-3.1 | CASINO-3.2 | CASINO-3.3 | 📝 Draft |
| **Evolution Gaming** | CASINO-3.4 | CASINO-3.5 | CASINO-3.6 | ⏳ Planejado |
| **PG Soft** | CASINO-3.7 | CASINO-3.8 | CASINO-3.9 | ⏳ Planejado |
| **Mancala** | CASINO-3.10 | CASINO-3.11 | CASINO-3.12 | ⏳ Planejado |
| **Digitain** | CASINO-3.13 | CASINO-3.14 | CASINO-3.15 | ⏳ Planejado |
| **Evoplay** | CASINO-3.16 | CASINO-3.17 | CASINO-3.18 | ⏳ Planejado |
| **OpenBox** | CASINO-3.19 | CASINO-3.20 | CASINO-3.21 | ⏳ Planejado |
| **Alternar** | CASINO-3.22 | CASINO-3.23 | CASINO-3.24 | ⏳ Planejado |

> **Nota:** Apenas Pragmatic Play (CASINO-3.1–3.3) tem stories criadas. Os demais providers terão stories criadas pelo @sm conforme CASINO-2 completa a Fase 2 de cada provider.

---

## Stories — Pragmatic Play (Primeiras a executar)

### CASINO-3.1 — Pragmatic Play Test Oracle

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-3.1 |
| **Arquivo Story** | `docs/stories/active/CASINO-3.1-pragmatic-play-test-oracle.md` |
| **Output** | `casino-proxy-test-oracle/pragmatic-play/` |
| **Status** | 📝 Draft |
| **Validação @po** | ⏳ Pendente |
| **Diferencial** | Primeiro oracle do projeto — define padrão WireMock/JUnit que outros providers replicarão |
| **Estimativa** | 8–12h |
| **Dependência** | CASINO-2.2 (9 endpoint docs) ✅ Done |

---

### CASINO-3.2 — Pragmatic Play Trace Matrix

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-3.2 |
| **Arquivo Story** | `docs/stories/active/CASINO-3.2-pragmatic-play-trace-matrix.md` |
| **Output** | `docs/architecture/casino-proxy/trace-matrices/pragmatic-play-trace-matrix.yaml` |
| **Status** | 📝 Draft |
| **Validação @po** | ⏳ Pendente |
| **Diferencial** | Primeiro trace matrix — liga todas as camadas: spec → regra → código PHP → teste → Go impl (TBD) |
| **Estimativa** | 2–4h |
| **Dependência** | CASINO-3.1 (test IDs devem existir para preencher `test_coverage`) |

---

### CASINO-3.3 — Pragmatic Play Validation Gate

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-3.3 |
| **Arquivo Story** | `docs/stories/active/CASINO-3.3-pragmatic-play-validate.md` |
| **Output** | `docs/architecture/casino-proxy/validation-gates/pragmatic-play-validation-report.md` |
| **Status** | 📝 Draft |
| **Validação @po** | ⏳ Pendente — GO libera Evolution Gaming |
| **Diferencial** | Gate: 100% PHP pass → @po aprova → desbloqueia CASINO-3.4 (Evolution Gaming) |
| **Estimativa** | 2–4h |
| **Dependência** | CASINO-3.1 + CASINO-3.2 |

---

## Kanban de Execução

```
PLANEJADO              DRAFT               IN PROGRESS         DONE
──────────────────    ──────────────────   ──────────────────  ──────────────
Evolution 3.4-3.6     PP: 3.1, 3.2, 3.3   —                  —
PG Soft 3.7-3.9
Mancala 3.10-3.12
Digitain 3.13-3.15
Evoplay 3.16-3.18
OpenBox 3.19-3.21
Alternar 3.22-3.24
```

---

## Artefatos

### Input (CASINO-2 — já existem para Pragmatic Play)

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | 12 regras BR-* | ✅ Completo |
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-*.md` | 9 endpoint docs | ✅ Completo (9/9) |

### Output (CASINO-3 — a criar)

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `casino-proxy-test-oracle/pragmatic-play/` | Suite JUnit 5 + WireMock | ⏳ Pendente |
| `docs/architecture/casino-proxy/trace-matrices/pragmatic-play-trace-matrix.yaml` | Trace matrix YAML | ⏳ Pendente |
| `docs/architecture/casino-proxy/validation-gates/pragmatic-play-validation-report.md` | Gate report | ⏳ Pendente |

---

## Definition of Done — Epic CASINO-3

- [ ] 8 suites Java criadas (1 por provider) com 50+ testes cada
- [ ] Todos os testes passam 100% contra PHP legado
- [ ] 8 trace matrices YAML criadas com cobertura 100% de regras BR-*
- [ ] 8 validation reports aprovados por @po
- [ ] Suite é executável contra Go sem modificação de código
- [ ] CASINO-4 desbloqueado

---

## Paralelismo com CASINO-2

Este epic é executado **em pipeline** com o CASINO-2:

```
CASINO-2                          CASINO-3
────────────────────────────────  ───────────────────────────────
PP Fase 1-2 ✅ (Done)             PP Fase 3-5 ← pode iniciar AGORA
Evo Fase 1-2 (em andamento) →     Evo Fase 3-5 (aguarda CASINO-2.7 ✅)
PG Soft Fase 1-2 →                PG Soft Fase 3-5 (aguarda CASINO-2.12 ✅)
...                               ...
```

Dev A documenta o próximo provider enquanto Dev B constrói oracle do provider anterior.

---

## Riscos

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| PHP acessível apenas em ambiente específico | Média | Alto | Documentar env setup no README do oracle; usar WireMock para desenvolvimento |
| WireMock stubs divergem do PHP real | Média | Alto | Validar stubs contra PHP antes de fechar CASINO-3.1 |
| Edge cases não cobertos em Fase 2 docs | Baixa | Médio | Fase 3 reporta gaps → volta para CASINO-2 doc para correção |
| Java não é skill do time Go | Baixa | Médio | Oracle é standalone; Go team só executa `mvn test`, não mantém |

---

## Histórico de Execução

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Epic CASINO-3 criado como Draft — aguarda refinamento com @po |
