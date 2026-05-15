# CASINO-3.2: Criar Trace Matrix YAML — Pragmatic Play

**Story ID:** CASINO-3.2  
**Epic:** CASINO-3 (Test Oracle — Business Rules Validation Suite)  
**Tipo:** Fase 4 de 5 — Trace Matrix  
**Status:** Draft  
**Prioridade:** Alta  
**Atribuído a:** @dev  
**Relacionado:** CASINO-3.1 (oracle tests — IDs necessários), CASINO-2.2 (docs), CASINO-1.7 (regras)  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Criar o arquivo YAML de rastreabilidade que mapeia cada regra de negócio BR-* do Pragmatic Play através de todas as camadas: spec OpenAPI → código PHP (arquivo + linhas) → test IDs do oracle → implementação Go (TBD).

**Objetivo:** Produzir `docs/architecture/casino-proxy/trace-matrices/pragmatic-play-trace-matrix.yaml` com 100% das 12 regras BR-* mapeadas.

---

## Contexto

### Por que esta Story?

A trace matrix é o **artefato de rastreabilidade** que liga todas as fases:

```
OpenAPI Spec  →  BR-* Rule  →  PHP Code  →  Oracle Test  →  Go Implementation
(CASINO-1)       (CASINO-1.7)  (CASINO-2.1)  (CASINO-3.1)   (CASINO-4, TBD)
```

Ela serve como:
1. **Auditoria:** Prova que nenhuma regra ficou sem teste e sem implementação
2. **Handoff para CASINO-4:** O time Go usa a matrix para saber exatamente o que implementar e qual teste validar
3. **Rastreamento de progresso:** Campo `go_implementation.status` vai de `pending` → `implemented` durante CASINO-4

### Por que depende de CASINO-3.1?

O campo `test_coverage.test_ids` referencia os IDs dos testes JUnit gerados em CASINO-3.1. Sem os testes existindo, o campo não pode ser preenchido com valores reais — e uma matrix com IDs fictícios seria inútil para auditoria.

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Arquivo `docs/architecture/casino-proxy/trace-matrices/pragmatic-play-trace-matrix.yaml` criado
- [ ] **AC-2:** Todas as 12 regras BR-* mapeadas — nenhuma sem entrada na matrix
- [ ] **AC-3:** Cada entrada contém todos os campos obrigatórios:
  - `openapi_reference` — path no spec OpenAPI correspondente
  - `php_code.file` — arquivo PHP onde a regra é implementada
  - `php_code.lines` — linhas aproximadas (conforme CASINO-1.7)
  - `test_coverage.file` — arquivo Java do oracle
  - `test_coverage.test_ids` — lista de IDs dos testes JUnit (de CASINO-3.1)
  - `go_implementation.status` — `"pending"` (preenchido durante CASINO-4)
- [ ] **AC-4:** 100% de cobertura — cada BR-* tem pelo menos 1 `test_id` referenciado
- [ ] **AC-5:** Arquivo YAML válido (sem erros de parse)
- [ ] **AC-6:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-7:** Campo `endpoints` em cada regra listando em quais dos 9 endpoints ela se aplica
- [ ] **AC-8:** Campo `notes` para regras com comportamento especial (ex: PP-007 — apenas authenticate)

### Fora do Escopo

- ❌ Preencher `go_implementation` além de `status: pending` — isso é CASINO-4
- ❌ Outros providers (stories separadas em CASINO-3.5, 3.8, etc.)

---

## Estrutura do Arquivo de Output

```yaml
# pragmatic-play-trace-matrix.yaml
# Trace Matrix — Pragmatic Play Business Rules
# Gerado: 2026-05-XX | Autor: @dev | Revisado: @po

provider: pragmatic-play
epic_extract: CASINO-1.7
epic_document: CASINO-2.2
epic_oracle: CASINO-3.1
epic_go: CASINO-4 (TBD)

rules:
  - id: BR-GENERIC-ROUTING-VALIDATION-001
    description: "Valida parâmetros obrigatórios da requisição antes do processamento"
    endpoints: [authenticate, balance, bet, refund, result, bonusWin, jackpotWin, promoWin, adjustment]
    openapi_reference: "application/docs/openapi/pragmatic-play.yaml#/paths/~1authenticate"
    php_code:
      file: "app/Services/PragmaticPlayService.php"
      lines: "~45-52"
    test_coverage:
      file: "casino-proxy-test-oracle/pragmatic-play/src/test/java/.../PragmaticPlayOracleTest.java"
      test_ids:
        - "routing_missingToken_returns400"
        - "routing_invalidHash_returns401"
    go_implementation:
      status: "pending"
      file: "TBD"
      lines: "TBD"

  # ... (11 regras restantes)
```

---

## File List

| Arquivo | Ação | Status |
|---------|------|--------|
| `docs/architecture/casino-proxy/trace-matrices/pragmatic-play-trace-matrix.yaml` | Criar | ⏳ |

---

## Estimativa

**2–4 horas** (@dev):
- Ler CASINO-3.1 para coletar test IDs reais: 30min
- Preencher 12 entradas de regras: 1–2h
- Validar YAML + verificar 100% cobertura: 30min
- Revisão com @po: 30min

---

## Dependências & Bloqueadores

**Precisa:**
- CASINO-3.1 ✅ (test IDs dos JUnit tests — campo `test_coverage.test_ids`)
- CASINO-1.7 ✅ (linhas PHP — campo `php_code.lines`)
- CASINO-2.2 ✅ (referências OpenAPI — campo `openapi_reference`)

**Bloqueia:**
- CASINO-3.3 (validation gate verifica que matrix está 100% preenchida)
- CASINO-4 (Go team usa matrix como handoff document)

---

## Change Log

| Data | Agente | Alteração |
|------|--------|-----------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
