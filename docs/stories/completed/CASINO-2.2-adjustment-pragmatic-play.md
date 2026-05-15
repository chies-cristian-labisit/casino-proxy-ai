# CASINO-2.2-adjustment: Documentar Endpoint /adjustment do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-adjustment  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Done  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7, CASINO-2.2-promoWin, CASINO-2.2-bet (Ready)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/adjustment` do Pragmatic Play seguindo o padrão estabelecido em `pragmatic-play-bet.md`. O `/adjustment` é o **último endpoint da Fase 2** e o único que **não pertence à família handleResult()**. É estruturalmente idêntico ao `/bet` e `/refund` — usa `userId`, lógica inline, passthrough de response, 9 regras genéricas — mas com contexto de negócio distinto: **correção/ajuste administrativo de saldo**, não um evento de jogo.

**Objetivo:** Produzir `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-adjustment.md` completo e encerrar a documentação de todos os 9 endpoints do Pragmatic Play na Fase 2.

---

## Contexto

### Por que esta Story?

O `/adjustment` fecha a Fase 2 completa. É o endpoint de uso menos frequente — acionado por operadores para corrigir saldos por motivos administrativos (erros de processamento, estornos manuais, reconciliações). Diferente de todos os outros endpoints:

- **Não é um evento de jogo:** Não corresponde a uma ação do jogador
- **Iniciado pelo operador:** É um ajuste administrativo, não reação a uma rodada
- **Não usa handleResult():** Tem lógica inline como bet/refund, apesar de ser o endpoint mais "tardio" na sequência PHP

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance      ✅ (template)
  ├─ /authenticate ✅ Ready
  ├─ /bet          ✅ Ready
  ├─ /refund       ✅ Ready
  ├─ /result       ✅ Ready
  ├─ /bonusWin     ✅ Ready
  ├─ /jackpotWin   ✅ Ready
  ├─ /promoWin     📝 Draft
  └─ /adjustment   ← ESTA STORY  (último — fecha Fase 2)
```

### Diferencial do /adjustment — Standalone (não-handleResult)

| Característica | handleResult() family | **/adjustment** |
|---------------|----------------------|----------------|
| Identificador | `userId` | `userId` |
| Implementação | Thin wrapper → handleResult() | **Lógica inline** (como bet/refund) |
| Regras exclusivas | Nenhuma | Nenhuma |
| Transforma response | Não (passthrough) | Não (passthrough) |
| URL de destino | `*.html` variável | `.../adjustment.html` |
| Contexto de negócio | Evento de jogo | **Ajuste administrativo de saldo** |
| Iniciador | Pragmatic Play (provider) | **Operador** |
| Fonte PHP | Wrapper + handleResult() | `adjustment()` linhas ~114-127 |
| userId removeTenant() | linha ~161 | **linha ~120** |

> O `/adjustment` é o **único endpoint não-handleResult() após o /refund** na sequência PHP. O @dev deve tratar como o terceiro membro do grupo bet/refund/adjustment — todos com lógica inline.

---

## Critérios de Aceitação

### Deve Ter

- [x] **AC-1:** Endpoint `/adjustment` analisado — 9 regras BR-* confirmadas contra `PragmaticPlayService.php` (método `adjustment()`, linhas ~114-127, userId em linha ~120)
- [x] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [x] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [x] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [x] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [x] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation)
- [x] **AC-7:** Arquivo criado em `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-adjustment.md`
- [x] **AC-8:** File List desta story atualizada

### Deveria Ter

- [x] **AC-9:** Seção de contexto de negócio explicando: quando `/adjustment` é chamado, quem o inicia (operador vs. provider), e distinção de todos os outros endpoints (ajuste administrativo vs. evento de jogo)
- [x] **AC-10:** Nota explícita de que `/adjustment` é standalone (lógica inline) — não usa handleResult() — com referência ao grupo bet/refund/adjustment como padrão de lógica inline

### Fora do Escopo

- ❌ Escrever testes (Fase 3)
- ❌ Criar matrizes YAML (Fase 4)
- ❌ Implementar handler Go
- ❌ Documentar outros endpoints (Fase 2 encerrada após esta story)
- ❌ Corrigir código PHP

---

## Detalhes Técnicos / Dev Notes

### Endpoint

```
Método: POST
URL:    /v1/webhooks/pragmatic-play/adjustment
Função: Aplicar ajuste/correção administrativa de saldo do jogador
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php
        - Método: adjustment()  linhas ~114-127
        - userId removeTenant() linha ~120
```

### Lógica Inline (como bet/refund)

```php
// adjustment() — lógica inline, sem delegação (linhas ~114-127)
public function adjustment($data) {
    $tenant = $this->operatorService->get($data['userId']);        // ~115
    $data['userId'] = $this->removeTenant($data['userId']);        // ~120
    $secret = $tenant->credentials()
        ->where('name', 'pragmatic')
        ->where('key', 'secret-key')
        ->first()->value;                                          // ~121
    $data['hash'] = $this->generateHashCode($data, $secret);      // ~123
    return $this->postJson(
        $tenant['url'] . '/pragmatic-play/adjustment.html', $data // ~125
    );
}
```

**Nota:** Estrutura idêntica a `bet()` (linhas 64-77) e `refund()` (linhas 79-92) — apenas o nome do método e a URL de destino diferem.

### Regras Aplicáveis (9 total — idênticas ao /bet e /refund)

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
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/adjustment.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Uso de userId em adjustment()

```php
// adjustment() — apenas userId (linha ~120 do PragmaticPlayService.php)
$tenant = $this->operatorService->get($data['userId']);
$data['userId'] = $this->removeTenant($data['userId']);
// ex: "myoperator_user456" → "user456"
```

**Código Fonte:** `PragmaticPlayService.php:120`

### Grupos de Implementação — Visão Geral dos 9 Endpoints

| Grupo | Endpoints | Implementação |
|-------|-----------|---------------|
| Sessão | `authenticate` | Inline + response transform (PP-007/012) |
| Consulta | `balance` | Inline + dual token (rule 010) |
| Transação inline | `bet`, `refund`, **`adjustment`** | Inline — lógica direta no método |
| handleResult() family | `result`, `bonusWin`, `jackpotWin`, `promoWin` | Thin wrapper → handleResult() |

### Referências de Código Fonte

```
PragmaticPlayService.php:114-127   → método adjustment()
PragmaticPlayService.php:120       → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137   → método removeTenant()
PragmaticPlayService.php:142-152   → método generateHashCode()
OperatorService.php:20-34          → método get()
BaseService.php:16-22              → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-adjustment.md` deve seguir o template de `pragmatic-play-bet.md`, substituindo:
- "bet" → "adjustment" em textos
- `bet.html` → `adjustment.html` na URL de destino
- Contexto de negócio: "registra aposta" → "aplica ajuste administrativo de saldo"
- Linhas de código fonte: ~64-77 → ~114-127, userId em ~70 → ~120
- **Adicionar seção de contexto de negócio** com distinção ajuste administrativo vs. evento de jogo

### Seções Obrigatórias

1. **Header** — endpoint, função, nota de encerramento da Fase 2
2. **Resumo Executivo** — o que faz, quem o inicia (operador), distinção dos demais
3. **Fluxo em 8 Fases** — Mermaid + explicação (Fase 2 = userId; Fase 8 = passthrough)
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5):
   - `userId` faltando → exception em tenant extraction
   - `userId` sem underscore → falha
   - Operador não encontrado → `ModelNotFoundException`
   - Credencial faltando → null reference exception
   - Provider timeout → falha imediata
6. **Exemplo Completo** — request com `userId` prefixado → sanitização → POST → passthrough response
7. **Contexto de negócio** — ajuste administrativo vs. evento de jogo; quando é chamado; quem inicia
8. **Tabela de grupos de implementação** — visão dos 9 endpoints organizados por padrão
9. **Checklist de Segurança**
10. **Limites e Restrições**

---

## Tasks / Subtasks

- [x] **T-1:** Ler `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` — usar como template direto (lógica inline idêntica)
- [x] **T-2:** Ler `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — confirmar regras e tabela de endpoints (adjustment)
- [x] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `adjustment()` (linhas ~114-127) — confirmar lógica inline e userId em ~120
- [x] **T-4:** Criar arquivo `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-adjustment.md`
- [x] **T-5:** Adaptar conteúdo de bet.md — substituir referências bet → adjustment, ajustar URLs e linhas de código
- [x] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid
- [x] **T-7:** Preencher Matriz de Regras (9 regras)
- [x] **T-8:** Documentar 5+ Cenários de Erro
- [x] **T-9:** Escrever exemplo completo request → response
- [x] **T-10:** Adicionar seção contexto de negócio (ajuste administrativo, iniciado por operador) + tabela de grupos de implementação dos 9 endpoints
- [x] **T-11:** Preencher Security Checklist + Limites e Restrições
- [x] **T-12:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low-Medium (adaptação de bet.md + seção de contexto administrativo + tabela de grupos)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect

**Quality Gate Tasks:**
- [x] Pre-Commit (@dev): [N/A — CodeRabbit disabled] — Markdown renderiza corretamente; sem referências residuais a "bet" onde deveria ser "adjustment"; tabela de grupos presente
- [x] Pre-PR (@devops): [N/A — CodeRabbit disabled] — Links e paths válidos

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
- Consistência bet → adjustment nas substituições de texto
- URL de destino correta: `adjustment.html`
- Linhas de código PHP corretas: ~114-127, userId em ~120 (não ~70 do bet)
- Contexto de negócio administrativo (não evento de jogo) presente e claro
- Tabela de grupos dos 9 endpoints presente (encerra Fase 2 com visão completa)

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-adjustment.md` | Documentação técnica do endpoint /adjustment | ✅ Criado |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-adjustment.md` | Output principal desta story | ✅ Criado |
| `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` | Template direto a seguir | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [x] Arquivo `pragmatic-play-adjustment.md` criado e completo
- [x] Diagrama Mermaid renderiza corretamente
- [x] `userId` como identificador documentado (linha PHP ~120)
- [x] Passthrough de response documentado (Fase 8)
- [x] URL de destino `adjustment.html`
- [x] Linhas de código PHP corretas (~114-127)
- [x] Contexto de negócio administrativo presente (distinção de evento de jogo)
- [x] Tabela de grupos dos 9 endpoints presente
- [x] Security checklist preenchido
- [x] File List desta story atualizada

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa fidelidade ao código PHP, contexto administrativo e completude da visão geral dos 9 endpoints.  
**Conclusão da Fase 2:** Após esta story aprovada, todos os 9 endpoints do Pragmatic Play estão documentados — Fase 2 completa.

---

## Métricas de Sucesso

- **Correção:** URL `adjustment.html`, linhas PHP corretas (~114-127 / ~120), sem erros de cópia de bet.md
- **Contexto:** Distinção ajuste administrativo vs. evento de jogo clara para implementação Go
- **Completude:** Tabela de grupos dos 9 endpoints fecha a documentação da Fase 2 com visão arquitetural completa
- **Rastreabilidade:** 9 regras mapeadas por fase

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1-2 horas (adaptação de bet.md + seção de contexto + tabela de grupos)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-bet ✅ (Ready)
- **Bloqueia:** CASINO-2.3 (próxima fase — Test Oracle)
- **Sequência Fase 2:** authenticate ✅ → bet ✅ → refund ✅ → result ✅ → bonusWin ✅ → jackpotWin ✅ → promoWin → **adjustment** ← ÚLTIMO

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Fecha Fase 2; tabela de grupos dos 9 endpoints e contexto administrativo exigidos em ACs. |
| 2026-05-14 | @dev (Dex) | Implementação completa — `pragmatic-play-adjustment.md` criado (10 seções, 6 cenários, tabela de grupos dos 9 endpoints, nota standalone/não-handleResult(), encerramento Fase 2). Todos T-1..T-12 concluídos. Status: InReview. |

| 2026-05-15 | @qa (Quinn) | QA Gate PASS — Todos os ACs atendidos, output file completo, CodeRabbit N/A. Status: InReview → Done |

## QA Results

### Review Date: 2026-05-15

### Reviewed By: Quinn (Test Architect)

All ACs verified complete. Output file exists and passes documentation quality checks. CodeRabbit disabled (N/A).

### Gate Status

Gate: PASS → docs/qa/gates/casino-2.2-adjustment.yml
