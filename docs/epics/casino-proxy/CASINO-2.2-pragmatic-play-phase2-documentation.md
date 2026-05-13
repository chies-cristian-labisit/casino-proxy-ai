# CASINO-2.2: Pragmatic Play — Fase 2 Technical Documentation

**Epic ID:** CASINO-2.2  
**Tipo:** Fase 2 de 5 — Technical Documentation (Pragmatic Play)  
**Status:** 🟡 Ready for Implementation — Stories completas, @dev pendente  
**Criado:** 2026-05-12  
**Atualizado:** 2026-05-13  
**Owner:** @sm (River)  
**Executor:** @dev (Dex)  
**Validação:** @po (Pax)  
**Depende De:** CASINO-2.1 (Fase 1 — Extract) ✅ Done  
**Bloqueia:** CASINO-2.3 (Fase 3 — Test Oracle)

---

## Objetivo

Produzir documentação técnica completa para os **9 endpoints do Pragmatic Play**, mostrando como as 12 regras de negócio (BR-*) extraídas na Fase 1 se combinam em fluxos de 8 fases por endpoint. Os documentos de output são a especificação que guia a implementação dos handlers Go na CASINO-3.

```
Input:   docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md (12 regras BR-*)
Output:  docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-{endpoint}.md (9 docs)
```

---

## Contexto

### Por que esta Fase?

A Fase 1 extraiu *o que o PHP faz*. A Fase 2 documenta *como o PHP faz* — transformando regras brutas em fluxos de 8 fases consumíveis por um @dev Go que nunca viu o código PHP. Sem esta documentação, a Fase 3 (Test Oracle) não tem base para escrever os testes.

### Estrutura dos Documentos de Output

Cada `pragmatic-play-{endpoint}.md` contém:
1. Resumo Executivo
2. Fluxo em 8 Fases (Mermaid diagram)
3. Matriz de Regras (9 regras × fase × exclusiva?)
4. Cenários de Erro (mínimo 5)
5. Exemplo Completo request → response
6. Seção de contexto de negócio
7. Checklist de Segurança
8. Limites e Restrições

---

## Arquitetura dos 9 Endpoints

Os 9 endpoints se dividem em **4 grupos de implementação**, fundamentais para o @dev Go:

| Grupo | Endpoints | Padrão PHP | Identificador | Response |
|-------|-----------|-----------|---------------|----------|
| **Sessão** | `authenticate` | Inline + transform | `token` | Re-prefixa `userId` se `error==0` |
| **Consulta** | `balance` | Inline + dual token | `token` OU `userId` | Passthrough |
| **Transação inline** | `bet`, `refund`, `adjustment` | Lógica direta no método | `userId` | Passthrough |
| **handleResult() family** | `result`, `bonusWin`, `jackpotWin`, `promoWin` | Thin wrapper → `handleResult()` | `userId` | Passthrough |

> **Regra de ouro para Go:** 8 dos 9 endpoints são passthrough. Apenas `/authenticate` transforma a resposta.

---

## Fluxo Canônico — 8 Fases

Todos os endpoints (com variações mínimas) seguem este fluxo:

```
Fase 1 → ROTEAMENTO     BR-GENERIC-ROUTING-VALIDATION-001
Fase 2 → TENANT         BR-GENERIC-TENANT-EXTRACTION-001
Fase 3 → OPERADOR       BR-GENERIC-OPERATOR-CACHING-001
Fase 4 → SANITIZAÇÃO    BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 + ORDER-001
Fase 5 → CREDENCIAL     BR-GENERIC-CREDENTIAL-LOOKUP-001
Fase 6 → HASH MD5       BR-GENERIC-AUTHENTICATION-HMAC-MD5-001
Fase 7 → HTTP POST      BR-GENERIC-PROVIDER-INTEGRATION-001
Fase 8 → RESPONSE       Passthrough (7 endpoints) | Transform userId (authenticate)
```

---

## Stories — Status de Execução

### Visão Geral

| Métrica | Valor |
|---------|-------|
| Total de stories | 9 |
| Stories criadas | 9 / 9 ✅ |
| Stories validadas (@po) | 9 / 9 ✅ |
| Documentos de output criados (@dev) | 0 / 9 ⏳ |
| Fase 2 completa | ❌ Pendente implementação |

---

### Detalhe por Story

#### Template Base

| Campo | Valor |
|-------|-------|
| **Story** | CASINO-2.2-balance (template) |
| **Endpoint** | `/balance` |
| **Arquivo Story** | — (não possui story separada; serviu como template) |
| **Documento Output** | `pragmatic-play-balance.md` |
| **Status Story** | ✅ Completo |
| **Status Output** | ✅ Criado (439 linhas — template canônico) |
| **Validação @po** | ✅ Aprovado |
| **Diferencial** | Dual token (`token` OU `userId`); rule 010 exclusiva |

---

#### Story 1 — /authenticate

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-authenticate |
| **Arquivo Story** | `docs/stories/CASINO-2.2-pragmatic-play-authenticate.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` |
| **Status Story** | ✅ Ready (GO 8/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-12 — GO 8/10 |
| **Diferencial** | Único endpoint com response transformation (re-prefixa `userId` se `error==0`); regras exclusivas PP-007 + PP-012 |
| **Estimativa @dev** | 2-3 horas |
| **Dependência** | CASINO-1.7 ✅ |

---

#### Story 2 — /bet

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-bet |
| **Arquivo Story** | `docs/stories/CASINO-2.2-bet-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-12 — GO 9/10. Fix: coluna "Regras exclusivas" → "Regras notáveis" |
| **Diferencial** | Padrão canônico de transação; usa `userId`; sem regras exclusivas; serve de modelo para refund/result/bonusWin/jackpotWin/promoWin |
| **Estimativa @dev** | 1-2 horas |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-authenticate ✅ |

---

#### Story 3 — /refund

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-refund |
| **Arquivo Story** | `docs/stories/CASINO-2.2-refund-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-refund.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-12 — GO 9/10 |
| **Diferencial** | Reverso do /bet; estruturalmente idêntico (única diferença: URL `refund.html`, contexto estorno, linhas ~79-92) |
| **Estimativa @dev** | 1 hora |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-bet ✅ |

---

#### Story 4 — /result

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-result |
| **Arquivo Story** | `docs/stories/CASINO-2.2-result-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-12 — GO 9/10 |
| **Diferencial** | Introduz padrão `handleResult()` — thin wrapper `result()` (~94-97) delega para método privado compartilhado (~161-175); serve de referência canônica para bonusWin/jackpotWin/promoWin |
| **Estimativa @dev** | 1-2 horas |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-refund ✅ |

---

#### Story 5 — /bonusWin

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-bonusWin |
| **Arquivo Story** | `docs/stories/CASINO-2.2-bonusWin-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bonusWin.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-13 — GO 9/10 |
| **Diferencial** | handleResult() family — wrapper `bonusWin()` (~99-102) passa `'bonusWin'` para handleResult(); URL `bonusWin.html`; contexto: pagamento de bônus |
| **Estimativa @dev** | 1 hora |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-result ✅ |

---

#### Story 6 — /jackpotWin

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-jackpotWin |
| **Arquivo Story** | `docs/stories/CASINO-2.2-jackpotWin-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-jackpotWin.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-13 — GO 9/10 |
| **Diferencial** | handleResult() family — wrapper `jackpotWin()` (~104-107); URL `jackpotWin.html`; contexto: evento de alta magnitude financeira com requisito de auditoria |
| **Estimativa @dev** | 1 hora |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-result ✅ |

---

#### Story 7 — /promoWin

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-promoWin |
| **Arquivo Story** | `docs/stories/CASINO-2.2-promoWin-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-promoWin.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-13 — GO 9/10 |
| **Diferencial** | handleResult() family — último membro; wrapper `promoWin()` (~109-112); URL `promoWin.html`; contexto: prêmio de campanha promocional do operador |
| **Estimativa @dev** | 1 hora |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-result ✅ |

---

#### Story 8 — /adjustment

| Campo | Valor |
|-------|-------|
| **Story ID** | CASINO-2.2-adjustment |
| **Arquivo Story** | `docs/stories/CASINO-2.2-adjustment-pragmatic-play.md` |
| **Documento Output** | `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-adjustment.md` |
| **Status Story** | ✅ Ready (GO 9/10) |
| **Status Output** | ⏳ A criar — @dev pendente |
| **Validação @po** | ✅ 2026-05-13 — GO 9/10 |
| **Diferencial** | Único endpoint **não**-handleResult() fora do grupo bet/refund; lógica inline (~114-127), userId em linha ~120; iniciado pelo operador (não pelo jogador); ajuste administrativo de saldo; inclui tabela de grupos dos 9 endpoints (fecha Fase 2) |
| **Estimativa @dev** | 1-2 horas |
| **Dependência** | CASINO-1.7 ✅, CASINO-2.2-bet ✅ |

---

## Kanban de Execução

```
STORIES (Backlog)       READY ✅             IN PROGRESS         DONE
────────────────────    ─────────────────    ────────────────    ────────────────
                        authenticate         —                   balance (template)
                        bet
                        refund
                        result
                        bonusWin
                        jackpotWin
                        promoWin
                        adjustment
```

> **Estado atual (2026-05-13):** Todas as 8 stories em Ready. @dev pode iniciar implementação em qualquer ordem — a dependência real é apenas que `pragmatic-play-bet.md` exista antes de refund/adjustment usá-lo como template.

### Ordem Recomendada de Implementação (@dev)

```
1. authenticate  (mais complexo — response transform; 2-3h)
2. bet           (padrão canônico — base para os demais; 1-2h)
3. refund        (clone de bet; 1h)
4. result        (introduz handleResult(); 1-2h)
5. bonusWin      (clone de result; 1h)
6. jackpotWin    (clone de result; 1h)
7. promoWin      (clone de result; 1h)
8. adjustment    (clone de bet + seção contexto; 1-2h)
```

**Total estimado:** 9-13 horas de implementação @dev

---

## Artefatos

### Input (Fase 1 — já existem)

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | 12 regras BR-* com rastreabilidade PHP | ✅ Completo |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` | Template canônico da Fase 2 | ✅ Completo |

### Output (Fase 2 — a criar pelo @dev)

| Arquivo | Endpoint | Status |
|---------|----------|--------|
| `pragmatic-play-authenticate.md` | `/authenticate` | ⏳ Pendente |
| `pragmatic-play-bet.md` | `/bet` | ⏳ Pendente |
| `pragmatic-play-refund.md` | `/refund` | ⏳ Pendente |
| `pragmatic-play-result.md` | `/result` | ⏳ Pendente |
| `pragmatic-play-bonusWin.md` | `/bonusWin` | ⏳ Pendente |
| `pragmatic-play-jackpotWin.md` | `/jackpotWin` | ⏳ Pendente |
| `pragmatic-play-promoWin.md` | `/promoWin` | ⏳ Pendente |
| `pragmatic-play-adjustment.md` | `/adjustment` | ⏳ Pendente |

### Stories (todas em `docs/stories/`)

| Arquivo | Status |
|---------|--------|
| `CASINO-2.2-pragmatic-play-authenticate.md` | ✅ Ready |
| `CASINO-2.2-bet-pragmatic-play.md` | ✅ Ready |
| `CASINO-2.2-refund-pragmatic-play.md` | ✅ Ready |
| `CASINO-2.2-result-pragmatic-play.md` | ✅ Ready |
| `CASINO-2.2-bonusWin-pragmatic-play.md` | ✅ Ready |
| `CASINO-2.2-jackpotWin-pragmatic-play.md` | ✅ Ready |
| `CASINO-2.2-promoWin-pragmatic-play.md` | ✅ Ready |
| `CASINO-2.2-adjustment-pragmatic-play.md` | ✅ Ready |

---

## Definition of Done — Epic CASINO-2.2

- [ ] 9 documentos `pragmatic-play-{endpoint}.md` criados em `phase-2-technical-documentation/`
- [ ] Todos os diagramas Mermaid renderizam corretamente
- [ ] Cada documento referencia as regras BR-* corretas por fase
- [ ] Passthrough documentado em 8/9 endpoints; response transform em 1/9 (authenticate)
- [ ] Padrão handleResult() documentado em result.md e referenciado em bonusWin/jackpotWin/promoWin
- [ ] Tabela de grupos dos 9 endpoints presente em adjustment.md
- [ ] @po revisa e aprova cada documento antes de CASINO-2.3 iniciar
- [ ] CASINO-2.3 desbloqueado

---

## Riscos

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| PHP source não acessível (submodule) | Alta | Médio | Regras já extraídas em phase-1-business-rules/pragmatic-play-rules.md — suficiente para documentação |
| Copy-paste residual entre stories similares | Média | Baixo | CodeRabbit Focus Areas detecta; @po valida cada output |
| Linhas PHP aproximadas (~) divergem do real | Baixa | Baixo | @dev confirma em T-3 de cada story ao ler o fonte |

---

## Histórico de Execução

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-11 | @dev (Dex) | CASINO-1.7 implementado — 12 regras BR-* extraídas |
| 2026-05-12 | @po (Pax) | CASINO-1.7 validado GO 8/10 — status: Done |
| 2026-05-12 | @sm (River) | CASINO-2.2-authenticate criada (Draft) |
| 2026-05-12 | @po (Pax) | CASINO-2.2-authenticate validada GO 8/10 — Ready |
| 2026-05-12 | @sm (River) | CASINO-2.2-bet criada (Draft) |
| 2026-05-12 | @po (Pax) | CASINO-2.2-bet validada GO 9/10 — Ready |
| 2026-05-12 | @sm (River) | CASINO-2.2-refund criada (Draft) |
| 2026-05-12 | @po (Pax) | CASINO-2.2-refund validada GO 9/10 — Ready |
| 2026-05-12 | @sm (River) | CASINO-2.2-result criada (Draft) |
| 2026-05-12 | @po (Pax) | CASINO-2.2-result validada GO 9/10 — Ready |
| 2026-05-12 | @sm (River) | CASINO-2.2-bonusWin e CASINO-2.2-jackpotWin criadas (Draft) |
| 2026-05-13 | @po (Pax) | CASINO-2.2-bonusWin validada GO 9/10 — Ready |
| 2026-05-13 | @po (Pax) | CASINO-2.2-jackpotWin validada GO 9/10 — Ready |
| 2026-05-12 | @sm (River) | CASINO-2.2-promoWin e CASINO-2.2-adjustment criadas (Draft) |
| 2026-05-13 | @po (Pax) | CASINO-2.2-promoWin validada GO 9/10 — Ready |
| 2026-05-13 | @po (Pax) | CASINO-2.2-adjustment validada GO 9/10 — Ready |
| 2026-05-13 | @devops (Gage) | Branch `docs/CASINO-2.2-phase2-pragmatic-play-stories` publicado |
| 2026-05-13 | @sm (River) | Epic CASINO-2.2 criado com tracking completo |
