# CASINO-2.2-promoWin: Documentar Endpoint /promoWin do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-promoWin  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7, CASINO-2.2-result (Ready), CASINO-2.2-jackpotWin (Ready)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/promoWin` do Pragmatic Play seguindo o padrão estabelecido em `pragmatic-play-result.md`. O `/promoWin` é o **quarto e último membro da família handleResult()** — estruturalmente idêntico a `/result`, `/bonusWin` e `/jackpotWin`, diferindo apenas no argumento `'promoWin'` passado para `handleResult()` e na URL de destino (`promoWin.html`).

**Objetivo:** Produzir `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-promoWin.md` completo — fluxo 8 fases, 9 regras mapeadas, exemplos request/response e security checklist. Esta é a última story da família handleResult().

---

## Contexto

### Por que esta Story?

O `/promoWin` registra o pagamento de um prêmio promocional — um evento distinto de jackpot e bônus, vinculado a campanhas específicas do operador. Completar esta story encerra a documentação da família handleResult() e desbloqueia `/adjustment`, o único endpoint restante da Fase 2.

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance      ✅ (template)
  ├─ /authenticate ✅ Ready
  ├─ /bet          ✅ Ready
  ├─ /refund       ✅ Ready
  ├─ /result       ✅ Ready (introduz handleResult())
  ├─ /bonusWin     ✅ Ready
  ├─ /jackpotWin   ✅ Ready
  ├─ /promoWin     ← ESTA STORY  (handleResult — último membro)
  └─ /adjustment
```

### Diferencial do /promoWin

| Característica | /result | /bonusWin | /jackpotWin | **/promoWin** |
|---------------|---------|-----------|-------------|--------------|
| Identificador | `userId` | `userId` | `userId` | `userId` |
| Regras exclusivas | Nenhuma | Nenhuma | Nenhuma | Nenhuma |
| Transforma response | Não | Não | Não | Não |
| handleResult() arg | `'result'` | `'bonusWin'` | `'jackpotWin'` | **`'promoWin'`** |
| URL de destino | `.../result.html` | `.../bonusWin.html` | `.../jackpotWin.html` | `.../promoWin.html` |
| Contexto de negócio | Resultado de rodada | Pagamento de bônus | Pagamento de jackpot | **Pagamento promocional** |
| Fonte PHP wrapper | ~94-97 | ~99-102 | ~104-107 | **~109-112** |

> Encerra a família handleResult() — 4 endpoints, 1 implementação compartilhada.

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Endpoint `/promoWin` analisado — 9 regras BR-* confirmadas. Wrapper `promoWin()` (~109-112) e delegação para `handleResult()` (~161-175) documentados
- [ ] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [ ] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [ ] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [ ] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [ ] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation)
- [ ] **AC-7:** Arquivo criado em `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-promoWin.md`
- [ ] **AC-8:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-9:** Seção de contexto de negócio explicando quando `/promoWin` é chamado vs. `/bonusWin` e `/jackpotWin` — distinção de prêmio promocional como evento de campanha do operador
- [ ] **AC-10:** Nota de encerramento da família handleResult() — referência cruzada para `pragmatic-play-result.md` como documento canônico da família

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
URL:    /v1/webhooks/pragmatic-play/promoWin
Função: Registrar pagamento de prêmio promocional do jogador
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php
        - Wrapper: método promoWin()    linhas ~109-112
        - Lógica:  método handleResult() linhas ~161-175
```

### Arquitetura: Wrapper → handleResult()

```php
// Thin wrapper público (linhas ~109-112)
public function promoWin($data) {
    return $this->handleResult('promoWin', $data);
}

// handleResult() idêntico ao chamado por result(), bonusWin(), jackpotWin()
// URL destino: {tenant_url}/pragmatic-play/promoWin.html
```

### Regras Aplicáveis (9 total)

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
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/promoWin.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Referências de Código Fonte

```
PragmaticPlayService.php:109-112   → método promoWin() (wrapper)
PragmaticPlayService.php:161-175   → método handleResult() (lógica compartilhada)
PragmaticPlayService.php:161       → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137   → método removeTenant()
PragmaticPlayService.php:142-152   → método generateHashCode()
OperatorService.php:20-34          → método get()
BaseService.php:16-22              → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-promoWin.md` deve seguir **exatamente** o template de `pragmatic-play-result.md`, substituindo:
- "result" → "promoWin" em textos
- `result.html` → `promoWin.html` na URL de destino
- Contexto de negócio: "resultado de rodada" → "pagamento de prêmio promocional"
- Linhas do wrapper: ~94-97 → ~109-112
- Nota de encerramento da família handleResult()

### Seções Obrigatórias

1. **Header** — endpoint, função, nota de encerramento da família handleResult()
2. **Resumo Executivo** — o que faz, quando é chamado, relação com família
3. **Fluxo em 8 Fases** — Mermaid + explicação
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5)
6. **Exemplo Completo** — request → response com passthrough
7. **Contexto de negócio** — distinção promoWin vs. bonusWin/jackpotWin
8. **Checklist de Segurança**
9. **Limites e Restrições**

---

## Tasks / Subtasks

- [ ] **T-1:** Ler `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` — usar como template direto
- [ ] **T-2:** Ler `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — confirmar regras
- [ ] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `promoWin()` (~109-112) — confirmar delegação para handleResult()
- [ ] **T-4:** Criar arquivo `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-promoWin.md`
- [ ] **T-5:** Adaptar conteúdo de result.md — substituir referências result → promoWin
- [ ] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid
- [ ] **T-7:** Preencher Matriz de Regras (9 regras)
- [ ] **T-8:** Documentar 5+ Cenários de Erro
- [ ] **T-9:** Escrever exemplo completo request → response
- [ ] **T-10:** Adicionar seção contexto promoWin vs. bonusWin/jackpotWin + nota de encerramento da família
- [ ] **T-11:** Preencher Security Checklist + Limites e Restrições
- [ ] **T-12:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low (adaptação direta de result.md — último membro da família)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): Markdown renderiza corretamente; sem referências residuais a "result", "bonusWin" ou "jackpotWin" onde deveria ser "promoWin"
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
- Consistência → promoWin nas substituições de texto (4 endpoints anteriores como fonte de erro de cópia)
- URL de destino correta: `promoWin.html`
- Linhas de código fonte corretas: wrapper ~109-112
- Nota de encerramento da família handleResult() presente

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-promoWin.md` | Documentação técnica do endpoint /promoWin | ⏳ A Criar |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-promoWin.md` | Output principal desta story | ⏳ A Criar |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` | Template direto a seguir | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [ ] Arquivo `pragmatic-play-promoWin.md` criado e completo
- [ ] Diagrama Mermaid renderiza corretamente
- [ ] `userId` como identificador documentado
- [ ] Passthrough de response documentado (Fase 8)
- [ ] URL de destino `promoWin.html`
- [ ] Linhas de código PHP corretas (wrapper ~109-112, handleResult() ~161-175)
- [ ] Nota de encerramento da família handleResult() presente
- [ ] Security checklist preenchido
- [ ] File List desta story atualizada

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa fidelidade ao código PHP e contexto de negócio do prêmio promocional.

---

## Métricas de Sucesso

- **Correção:** URL `promoWin.html`, linhas PHP corretas, sem resíduos dos endpoints anteriores
- **Contexto:** Distinção promoWin vs. bonusWin/jackpotWin clara
- **Fechamento:** Família handleResult() completa e documentada (result → bonusWin → jackpotWin → promoWin)

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1 hora (adaptação direta de result.md)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-result ✅ (Ready)
- **Bloqueia:** CASINO-2.2-adjustment (último endpoint da Fase 2)
- **Sequência Fase 2:** authenticate ✅ → bet ✅ → refund ✅ → result ✅ → bonusWin ✅ → jackpotWin ✅ → **promoWin** → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Encerra família handleResult(); nota de fechamento exigida em AC-10. |
