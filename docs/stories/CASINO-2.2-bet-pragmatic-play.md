# CASINO-2.2-bet: Documentar Endpoint /bet do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-bet  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7 (Regras de Negócio Pragmatic Play), CASINO-2.2-authenticate (story anterior da Fase 2)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/bet` do Pragmatic Play seguindo o template de `pragmatic-play-balance.md`. Este é o **segundo endpoint da Fase 2** e representa o **padrão canônico de transação** do sistema — usa `userId` como identificador único, faz passthrough direto da resposta do provider, e aplica as 9 regras genéricas sem nenhuma exclusiva.

**Objetivo:** Produzir `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` completo — fluxo 8 fases, 9 regras mapeadas, exemplos request/response e security checklist.

---

## Contexto

### Por que esta Story?

O `/bet` é estrategicamente importante porque estabelece o **padrão de endpoints de transação**: usa `userId` (não `token` como /authenticate), faz passthrough da resposta (como /balance), e não tem regras exclusivas. Os 5 endpoints subsequentes (`/refund`, `/result`, `/bonusWin`, `/jackpotWin`, `/promoWin`) seguem o mesmo padrão base.

**Documentar /bet corretamente = modelo reutilizável para os próximos 5 endpoints.**

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance    ✅ (template criado)
  ├─ /authenticate ✅ Ready (CASINO-2.2-authenticate)
  ├─ /bet         ← ESTA STORY
  ├─ /refund
  ├─ /result      ┐
  ├─ /bonusWin    │ Mesmo padrão que /bet
  ├─ /jackpotWin  │ (handleResult internamente)
  ├─ /promoWin    ┘
  └─ /adjustment
```

### Diferencial do /bet vs Outros Endpoints

| Característica | /authenticate | /balance | **/bet** |
|---------------|--------------|---------|---------|
| Identificador | `token` | `token` ou `userId` | `userId` apenas |
| Dual token | Não | **Sim** (rule 010) | Não |
| Transforma response | **Sim** (PP-012) | Não | Não |
| Regras exclusivas | PP-007, PP-012 | rule 010 | **Nenhuma** |
| Regras notáveis | — | rule 011 | rule 011 (simplificado) |
| Tipo | Sessão | Consulta | **Transação** |

> **Nota:** Rule 011 (BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001) aplica-se a `/bet` mas de forma simplificada: há apenas `userId` para sanitizar (sem ordenação entre múltiplos campos como em /balance).

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Endpoint `/bet` analisado — 9 regras BR-* identificadas e confirmadas contra `PragmaticPlayService.php` (método `bet()`, linhas ~64-77)
- [ ] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [ ] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [ ] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [ ] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [ ] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation, endpoint validation)
- [ ] **AC-7:** Arquivo criado em `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md`
- [ ] **AC-8:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-9:** Tabela comparativa `/bet` vs `/authenticate` vs `/balance` — destacar que /bet é o padrão canônico de transação
- [ ] **AC-10:** Nota sobre `/bet` como modelo para `/refund`, `/result`, `/bonusWin`, `/jackpotWin`, `/promoWin` — indicar quais partes serão idênticas

### Fora do Escopo

- ❌ Escrever testes (Fase 3)
- ❌ Criar matrizes YAML (Fase 4)
- ❌ Implementar handler Go
- ❌ Documentar outros endpoints (cada um é uma story separada)
- ❌ Corrigir código PHP

---

## Detalhes Técnicos / Dev Notes

### Endpoint

```
Método: POST
URL:    /v1/webhooks/pragmatic-play/bet
Função: Registrar aposta do jogador e retornar resposta do provider sem transformação
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php (método bet(), linhas ~64-77)
```

### Regras Aplicáveis (9 total)

| # | ID | Descrição | Fase no Fluxo | Exclusiva? |
|---|-----|-----------|--------------|------------|
| 1 | BR-GENERIC-ROUTING-VALIDATION-001 | Dynamic endpoint routing via method resolution | 1 | Não |
| 2 | BR-GENERIC-ERROR-HANDLING-001 | Unknown endpoint → exception 500 | 1 (guard) | Não |
| 3 | BR-GENERIC-TENANT-EXTRACTION-001 | Extrair operator_slug do `userId` (`userId.split('_')`) | 2 | Não |
| 4 | BR-GENERIC-OPERATOR-CACHING-001 | Operator lookup com cache 1h | 3 | Não |
| 5 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | Remover prefixo tenant do `userId` antes de enviar ao provider | 4 | Não |
| 6 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001 | Sanitização de `userId` (único campo — sem ordenação entre múltiplos campos) | 4 | Não |
| 7 | BR-GENERIC-CREDENTIAL-LOOKUP-001 | Buscar secret-key do operador | 5 | Não |
| 8 | BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | Gerar hash MD5 (sort payload + concat secret + md5) | 6 | Não |
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/bet.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Uso de userId (não token)

```php
// bet() — apenas userId
$tenant = $this->operatorService->get($data['userId']);
$data['userId'] = $this->removeTenant($data['userId']);
// ex: "myoperator_user456" → "user456"
```

Comparação com outros endpoints:
```
/authenticate: $this->operatorService->get($data['token'])          → usa token
/balance:      $this->operatorService->get($data['token'] ?? $data['userId'])  → dual
/bet:          $this->operatorService->get($data['userId'])          → usa userId
```

### Response Passthrough (Fase 8)

```
response = postJson(url, payload)
RETURN response  // inalterado — sem re-prefixação, sem transformação
```

Igual ao comportamento de `/balance`, `/refund`, `/result`, etc. Apenas `/authenticate` é diferente.

### Referências de Código Fonte

```
PragmaticPlayService.php:64-77   → método bet()
PragmaticPlayService.php:70      → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137 → método removeTenant()
PragmaticPlayService.php:142-152 → método generateHashCode()
OperatorService.php:20-34        → método get() (tenant extraction + cache)
BaseService.php:16-22            → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-bet.md` deve seguir o template de `pragmatic-play-balance.md` com as adaptações abaixo:

### Seções Obrigatórias

1. **Header** — endpoint, função, nota de que é o padrão canônico de transação
2. **Resumo Executivo** — o que faz, quando é chamado, por que é o modelo para os próximos 5 endpoints
3. **Fluxo em 8 Fases** — diagrama Mermaid + explicação por fase
   _(Fase 2 usa `userId`; Fase 8 é passthrough — igual ao /balance)_
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5):
   - `userId` faltando no payload → exception em tenant extraction
   - `userId` sem underscore → tenant extraction falha
   - Operador não encontrado → `ModelNotFoundException`
   - Credencial Pragmatic faltando → null reference exception
   - Provider timeout → falha imediata (sem retry)
6. **Exemplo Completo** — request com `userId` prefixado → sanitização → POST → passthrough response
7. **Comparação /bet vs /authenticate vs /balance** — tabela de diferenças-chave
8. **Checklist de Segurança** — tenant isolation, hash auth, credential validation
9. **Limites e Restrições** — cache TTL, formato userId, algoritmo hash, retry policy

---

## Tasks / Subtasks

> Sequência de implementação para @dev

- [ ] **T-1:** Ler `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` — usar como template base
- [ ] **T-2:** Ler `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` — consultar seção de diferenças
- [ ] **T-3:** Ler `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — confirmar regras 001-009, 011
- [ ] **T-4:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `bet()` (linhas ~64-77)
- [ ] **T-5:** Criar arquivo `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md`
- [ ] **T-6:** Escrever Header + Resumo Executivo (destacar: padrão canônico de transação, modelo para próximos 5)
- [ ] **T-7:** Escrever Fluxo 8 Fases com diagrama Mermaid (Fase 2 = userId; Fase 8 = passthrough)
- [ ] **T-8:** Preencher Matriz de Regras (9 regras × fase × exclusiva?)
- [ ] **T-9:** Documentar 5+ Cenários de Erro
- [ ] **T-10:** Escrever exemplo completo request → response
- [ ] **T-11:** Adicionar tabela comparativa /bet vs /authenticate vs /balance
- [ ] **T-12:** Preencher Security Checklist + Limites e Restrições
- [ ] **T-13:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low (adaptação do template — sem regras exclusivas)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect (revisar fidelidade das regras documentadas)

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): Markdown renderiza corretamente (headings, tabelas, Mermaid)
- [ ] Pre-PR (@devops): Links e referências a arquivos existentes válidos

**Self-Healing Configuration:**
```yaml
mode: light
max_iterations: 2
severity_filter: [CRITICAL, HIGH]
behavior:
  CRITICAL: auto_fix
  HIGH: document_as_debt
```

**Focus Areas (Documentation):**
- Markdown quality e estrutura
- Consistência com template `pragmatic-play-balance.md` e `pragmatic-play-authenticate.md`
- Fidelidade técnica: uso de `userId` (não `token`) na Fase 2
- Passthrough de response claramente documentado (sem transformação)

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` | Documentação técnica do endpoint /bet | ⏳ A Criar |

### Templates de Referência

```
docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md
docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md
```

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` | Output principal desta story | ⏳ A Criar |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` | Template base | ✅ Existe |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` | Referência de comparação | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [ ] Arquivo `pragmatic-play-bet.md` criado e completo
- [ ] Diagrama Mermaid renderiza corretamente (8 fases visíveis)
- [ ] Uso de `userId` (não `token`) claramente documentado na Fase 2
- [ ] Passthrough de response documentado na Fase 8 (sem transformação)
- [ ] Tabela comparativa /bet vs /authenticate vs /balance presente
- [ ] Security checklist preenchido
- [ ] File List desta story atualizada
- [ ] Pronto para validação de @po antes de documentar próximo endpoint (/refund)

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa documento para confirmar fidelidade às regras extraídas na Fase 1.  
**Próxima Fase:** Fase 3 (CASINO-2.3) criará testes Java para as regras documentadas aqui.

---

## Métricas de Sucesso

- **Completude:** Todas as 9 regras documentadas no fluxo
- **Clareza:** @dev consegue ler e implementar handler Go sem consultar código PHP
- **Rastreabilidade:** Cada fase do fluxo referencia a regra BR-* correspondente
- **Reusabilidade:** Documento serve de modelo explícito para /refund, /result, /bonusWin, /jackpotWin, /promoWin

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1-2 horas (mais simples que /authenticate — sem regras exclusivas)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-authenticate ✅ (Ready)
- **Bloqueia:** CASINO-2.2-refund (próximo endpoint)
- **Sequência completa Fase 2:** authenticate ✅ → **bet** → refund → result → bonusWin → jackpotWin → promoWin → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Corrigida coluna "Regras exclusivas" → "Regras notáveis" na tabela comparativa. |
