# CASINO-3.3: Validation Gate — Pragmatic Play Oracle 100% PHP PASS

**Story ID:** CASINO-3.3  
**Epic:** CASINO-3 (Test Oracle — Business Rules Validation Suite)  
**Tipo:** Fase 5 de 5 — Validation Gate  
**Status:** Draft  
**Prioridade:** Alta (gate: GO libera Evolution Gaming CASINO-3.4)  
**Atribuído a:** @dev + @po (gate)  
**Relacionado:** CASINO-3.1 (oracle suite), CASINO-3.2 (trace matrix)  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Executar a suite oracle completa do Pragmatic Play contra o PHP legado, produzir o relatório de validação e obter aprovação do @po. **Este gate desbloqueia o início do oracle para Evolution Gaming (CASINO-3.4).**

**Objetivo:** 100% dos testes passam contra PHP. @po aprova. `pragmatic-play-validation-report.md` criado e assinado.

---

## Contexto

### Por que este Gate?

A validação formal garante que:
1. O oracle não tem falsos positivos (testes que passam mas medem a coisa errada)
2. O comportamento PHP documentado em CASINO-2.2 está corretamente testado
3. O critério de aceitação para CASINO-4 é confiável

### Sequência no Provider

```
CASINO-3.1: Suite criada (50+ testes, WireMock stubs) ✅
CASINO-3.2: Trace matrix criada (12 regras mapeadas) ✅
CASINO-3.3: ← VOCÊ ESTÁ AQUI
  ├─ Executar: mvn clean test (contra PHP real)
  ├─ Resultado: 100% PASS
  ├─ Relatório: pragmatic-play-validation-report.md
  └─ Gate: @po aprova → CASINO-3.4 (Evolution Gaming) desbloqueado
```

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Suite completa executada contra PHP legado: `mvn clean test -Dphp.base.url={php-url}`
- [ ] **AC-2:** Resultado: 0 failures, 0 errors — todos os 50+ testes PASS
- [ ] **AC-3:** Relatório criado em `docs/architecture/casino-proxy/validation-gates/pragmatic-play-validation-report.md` contendo:
  - Total de testes executados
  - Resultado: BUILD SUCCESS
  - Cobertura de regras BR-* (checklist 12/12)
  - Cobertura de endpoints (checklist 9/9)
  - Output completo do Maven (ou link para CI run)
  - Data e ambiente de execução
- [ ] **AC-4:** @po revisa relatório e aprova formalmente (sign-off no relatório)
- [ ] **AC-5:** Trace matrix (CASINO-3.2) verificada — 100% de regras com `test_ids` preenchidos
- [ ] **AC-6:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-7:** CI run documentado (GitHub Actions ou similar) mostrando verde
- [ ] **AC-8:** Qualquer falha encontrada e corrigida documentada no Change Log desta story

### Fora do Escopo

- ❌ Executar contra Go (isso é critério de aceitação do CASINO-4)
- ❌ Outros providers

---

## Estrutura do Validation Report

```markdown
# Validation Report — Pragmatic Play Oracle
**Data:** YYYY-MM-DD
**Ambiente:** {PHP URL, versão, branch}
**Executado por:** @dev

## Resultado

| Métrica | Valor |
|---------|-------|
| Total de testes | 5X |
| PASS | 5X |
| FAIL | 0 |
| ERRORS | 0 |
| BUILD | ✅ SUCCESS |

## Cobertura de Regras BR-*
- [x] BR-GENERIC-ROUTING-VALIDATION-001
- [x] BR-GENERIC-TENANT-EXTRACTION-001
- [x] BR-GENERIC-OPERATOR-CACHING-001
- [x] BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001
- [x] BR-PRAGMATIC-PARAMETER-ORDER-001
- [x] BR-GENERIC-CREDENTIAL-LOOKUP-001
- [x] BR-GENERIC-AUTHENTICATION-HMAC-MD5-001
- [x] BR-GENERIC-PROVIDER-INTEGRATION-001
- [x] BR-GENERIC-RESPONSE-PASSTHROUGH-001
- [x] BR-PRAGMATIC-BALANCE-DUAL-TOKEN-001
- [x] BR-PRAGMATIC-AUTH-USERID-REPREFIX-001
- [x] BR-PRAGMATIC-AUTH-TRANSFORM-EXCLUSIVE-001

## Cobertura de Endpoints
- [x] /authenticate
- [x] /balance
- [x] /bet
- [x] /refund
- [x] /result
- [x] /bonusWin
- [x] /jackpotWin
- [x] /promoWin
- [x] /adjustment

## @po Sign-Off
**Aprovado por:** Pax (@po)
**Data:** YYYY-MM-DD
**Gate Decision:** ✅ GO — Evolution Gaming (CASINO-3.4) desbloqueado
```

---

## File List

| Arquivo | Ação | Status |
|---------|------|--------|
| `docs/architecture/casino-proxy/validation-gates/pragmatic-play-validation-report.md` | Criar | ⏳ |

---

## Estimativa

**2–4 horas** (@dev + @po):
- Executar suite + corrigir eventuais falhas de ambiente: 1–2h
- Redigir relatório: 30min
- @po review: 30min–1h
- Buffer: 30min

---

## Dependências & Bloqueadores

**Precisa:**
- CASINO-3.1 ✅ (suite executável)
- CASINO-3.2 ✅ (trace matrix completa para verificar cobertura)
- Acesso ao PHP legado (mesmo ambiente do CASINO-3.1)

**Bloqueia:**
- CASINO-3.4 (Evolution Gaming Test Oracle — só inicia após este gate GO)

---

## Change Log

| Data | Agente | Alteração |
|------|--------|-----------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
