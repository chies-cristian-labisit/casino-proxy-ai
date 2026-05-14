# CASINO-2.2-refund: Documentar Endpoint /refund do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-refund  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** InReview  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7, CASINO-2.2-authenticate (Ready), CASINO-2.2-bet (Ready)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/refund` do Pragmatic Play seguindo o padrão estabelecido em `pragmatic-play-bet.md`. O `/refund` é o **reverso de uma aposta** — cancela ou estorna uma transação de bet — e segue o **mesmo padrão estrutural do /bet**: usa `userId`, passthrough de response, 9 regras genéricas, sem regras exclusivas.

**Objetivo:** Produzir `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-refund.md` completo — fluxo 8 fases, 9 regras mapeadas, exemplos request/response e security checklist.

---

## Contexto

### Por que esta Story?

O `/refund` fecha o ciclo transacional junto com `/bet`: se `/bet` registra uma aposta, `/refund` a cancela. Documentar os dois juntos como par (bet + refund) é essencial para que o handler Go implemente a lógica de rollback corretamente.

Do ponto de vista técnico, `/refund` é **estruturalmente idêntico ao /bet** — mesmo identificador, mesmas regras, mesmo passthrough. A única diferença relevante é o **contexto de negócio** (estorno vs aposta) e a **URL de destino** (`refund.html`).

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance      ✅ (template)
  ├─ /authenticate ✅ Ready
  ├─ /bet          ✅ Ready
  ├─ /refund        ← ESTA STORY
  ├─ /result        ┐
  ├─ /bonusWin      │ handleResult() internamente
  ├─ /jackpotWin    │
  ├─ /promoWin      ┘
  └─ /adjustment
```

### Diferencial do /refund

| Característica | /bet | **/refund** |
|---------------|------|------------|
| Identificador | `userId` | `userId` |
| Regras exclusivas | Nenhuma | Nenhuma |
| Transforma response | Não (passthrough) | Não (passthrough) |
| Regras notáveis | rule 011 | rule 011 |
| URL de destino | `.../bet.html` | `.../refund.html` |
| Contexto de negócio | Registra aposta | **Estorna aposta** |
| Fonte PHP | `bet()` linhas ~64-77 | `refund()` linhas ~79-92 |

> O `/refund` é tecnicamente o endpoint mais simples da série — idêntico ao /bet em estrutura, diferindo apenas no propósito e URL.

---

## Critérios de Aceitação

### Deve Ter

- [x] **AC-1:** Endpoint `/refund` analisado — 9 regras BR-* confirmadas contra `PragmaticPlayService.php` (método `refund()`, linhas ~79-92)
- [x] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [x] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [x] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [x] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [x] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation)
- [x] **AC-7:** Arquivo criado em `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-refund.md`
- [x] **AC-8:** File List desta story atualizada

### Deveria Ter

- [x] **AC-9:** Seção de contexto de negócio explicando a relação bet ↔ refund (par transacional) e quando cada um é chamado
- [x] **AC-10:** Nota explícita de que `/refund` é estruturalmente idêntico ao `/bet`, servindo de confirmação do padrão canônico

### Fora do Escopo

- ❌ Escrever testes (Fase 3)
- ❌ Criar matrizes YAML (Fase 4)
- ❌ Implementar handler Go
- ❌ Documentar outros endpoints
- ❌ Corrigir código PHP

---

## Detalhes Técnicos / Dev Notes

### Endpoint

```
Método: POST
URL:    /v1/webhooks/pragmatic-play/refund
Função: Estornar/cancelar uma aposta previamente registrada
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php (método refund(), linhas ~79-92)
```

### Regras Aplicáveis (9 total — idênticas ao /bet)

| # | ID | Descrição | Fase no Fluxo | Exclusiva? |
|---|-----|-----------|--------------|------------|
| 1 | BR-GENERIC-ROUTING-VALIDATION-001 | Dynamic endpoint routing via method resolution | 1 | Não |
| 2 | BR-GENERIC-ERROR-HANDLING-001 | Unknown endpoint → exception 500 | 1 (guard) | Não |
| 3 | BR-GENERIC-TENANT-EXTRACTION-001 | Extrair operator_slug do `userId` | 2 | Não |
| 4 | BR-GENERIC-OPERATOR-CACHING-001 | Operator lookup com cache 1h | 3 | Não |
| 5 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | Remover prefixo tenant do `userId` | 4 | Não |
| 6 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-ORDER-001 | Sanitização de `userId` (campo único) | 4 | Não |
| 7 | BR-GENERIC-CREDENTIAL-LOOKUP-001 | Buscar secret-key do operador | 5 | Não |
| 8 | BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | Gerar hash MD5 | 6 | Não |
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/refund.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Uso de userId

```php
// refund() — apenas userId (idêntico ao bet())
$tenant = $this->operatorService->get($data['userId']);
$data['userId'] = $this->removeTenant($data['userId']);
// ex: "myoperator_user456" → "user456"
```

**Código Fonte:** `PragmaticPlayService.php:85`

### Referências de Código Fonte

```
PragmaticPlayService.php:79-92   → método refund()
PragmaticPlayService.php:85      → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137 → método removeTenant()
PragmaticPlayService.php:142-152 → método generateHashCode()
OperatorService.php:20-34        → método get()
BaseService.php:16-22            → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-refund.md` deve seguir **exatamente** o template de `pragmatic-play-bet.md`, substituindo:
- "bet" → "refund" em textos
- `bet.html` → `refund.html` na URL de destino
- Contexto de negócio: "aposta" → "estorno de aposta"
- Linhas de código fonte: ~64-77 → ~79-92

### Seções Obrigatórias

1. **Header** — endpoint, função, nota sobre par bet↔refund
2. **Resumo Executivo** — o que faz, quando é chamado, relação com /bet
3. **Fluxo em 8 Fases** — Mermaid + explicação (Fase 2 = userId; Fase 8 = passthrough)
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5):
   - `userId` faltando → exception em tenant extraction
   - `userId` sem underscore → falha
   - Operador não encontrado → `ModelNotFoundException`
   - Credencial faltando → null reference exception
   - Provider timeout → falha imediata
6. **Exemplo Completo** — request com `userId` prefixado → sanitização → POST → passthrough response
7. **Contexto bet↔refund** — quando cada um é chamado, fluxo de ciclo de vida da transação
8. **Checklist de Segurança**
9. **Limites e Restrições**

---

## Tasks / Subtasks

> Sequência de implementação para @dev

- [x] **T-1:** Ler `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` — usar como template direto (estrutura idêntica)
- [x] **T-2:** Ler `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — confirmar regras aplicáveis ao refund
- [x] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `refund()` (linhas ~79-92) — confirmar identidade com bet()
- [x] **T-4:** Criar arquivo `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-refund.md`
- [x] **T-5:** Adaptar conteúdo de bet.md — substituir referências bet → refund, ajustar URLs e linhas de código
- [x] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid
- [x] **T-7:** Preencher Matriz de Regras (9 regras)
- [x] **T-8:** Documentar 5+ Cenários de Erro
- [x] **T-9:** Escrever exemplo completo request → response
- [x] **T-10:** Adicionar seção contextual bet↔refund (par transacional)
- [x] **T-11:** Preencher Security Checklist + Limites e Restrições
- [x] **T-12:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low (adaptação direta de bet.md — estrutura idêntica)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): Markdown renderiza corretamente; sem referências residuais a "bet" onde deveria ser "refund"
- [ ] Pre-PR (@devops): Links e paths válidos

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
- Consistência bet → refund nas substituições de texto
- URL de destino correta: `refund.html` (não `bet.html`)
- Linhas de código fonte corretas: ~79-92 (não ~64-77)
- Contexto de negócio do estorno presente e claro

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-refund.md` | Documentação técnica do endpoint /refund | ✅ Criado |

### Template de Referência

```
docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md  ← principal
docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md
```

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-refund.md` | Output principal desta story | ✅ Criado |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` | Template direto a seguir | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [x] Arquivo `pragmatic-play-refund.md` criado e completo
- [x] Diagrama Mermaid renderiza corretamente
- [x] `userId` como identificador documentado (não `token`)
- [x] Passthrough de response documentado (Fase 8)
- [x] Seção contexto bet↔refund presente
- [x] URL de destino `refund.html` (não `bet.html`)
- [x] Linhas de código PHP corretas (~79-92)
- [x] Security checklist preenchido
- [x] File List desta story atualizada

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa fidelidade ao código PHP e contexto de negócio do estorno.  
**Próxima Fase:** Fase 3 (CASINO-2.3) criará testes Java para as regras documentadas.

---

## Métricas de Sucesso

- **Correção:** URL `refund.html`, linhas PHP ~79-92, sem erros de cópia de bet.md
- **Contexto:** Relação bet↔refund clara para implementação Go
- **Rastreabilidade:** 9 regras mapeadas por fase
- **Clareza:** @dev Go consegue implementar o handler sem consultar PHP

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1 hora (adaptação direta de bet.md — menor esforço da série)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-bet ✅ (Ready)
- **Bloqueia:** CASINO-2.2-result
- **Sequência Fase 2:** authenticate ✅ → bet ✅ → **refund** → result → bonusWin → jackpotWin → promoWin → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Story completa e rastreável; riscos de copy-paste cobertos no CodeRabbit Focus Areas. |
| 2026-05-14 | @dev (Dex) | Implementação completa — `pragmatic-play-refund.md` criado (9 seções, 6 cenários de erro, exemplo completo, seção contexto bet↔refund com side-by-side PHP, tabela ciclo transacional). Todos T-1..T-12 concluídos. Status: InReview. |
