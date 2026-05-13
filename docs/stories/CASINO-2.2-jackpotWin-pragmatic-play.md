# CASINO-2.2-jackpotWin: Documentar Endpoint /jackpotWin do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-jackpotWin  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7, CASINO-2.2-result (Ready), CASINO-2.2-bonusWin  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/jackpotWin` do Pragmatic Play seguindo o padrão estabelecido em `pragmatic-play-result.md`. O `/jackpotWin` é o **terceiro membro da família handleResult()** — estruturalmente idêntico ao `/result` e `/bonusWin`, diferindo apenas no argumento `'jackpotWin'` passado para `handleResult()` e na URL de destino (`jackpotWin.html`).

**Objetivo:** Produzir `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-jackpotWin.md` completo — fluxo 8 fases, 9 regras mapeadas, exemplos request/response e security checklist.

---

## Contexto

### Por que esta Story?

O `/jackpotWin` registra o pagamento de um jackpot ganho pelo jogador — um evento financeiramente significativo que requer rastreabilidade independente. Embora a implementação PHP seja idêntica aos demais membros da família handleResult(), o contexto de negócio é distinto:

- **Impacto financeiro:** Jackpots envolvem valores geralmente maiores que resultados normais
- **Auditoria:** Operadores e reguladores exigem rastreabilidade separada de pagamentos de jackpot
- **Contexto para Go:** O handler Go deve rotear `/jackpotWin` como endpoint independente

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance      ✅ (template)
  ├─ /authenticate ✅ Ready
  ├─ /bet          ✅ Ready
  ├─ /refund       ✅ Ready
  ├─ /result       ✅ Ready (introduz handleResult())
  ├─ /bonusWin     📝 Draft
  ├─ /jackpotWin   ← ESTA STORY  (handleResult — terceiro membro)
  ├─ /promoWin     ┐ handleResult() family — último membro
  └─ /adjustment
```

### Diferencial do /jackpotWin

| Característica | /result | /bonusWin | **/jackpotWin** |
|---------------|---------|-----------|----------------|
| Identificador | `userId` | `userId` | `userId` |
| Regras exclusivas | Nenhuma | Nenhuma | Nenhuma |
| Transforma response | Não | Não | Não |
| handleResult() arg | `'result'` | `'bonusWin'` | **`'jackpotWin'`** |
| URL de destino | `.../result.html` | `.../bonusWin.html` | `.../jackpotWin.html` |
| Contexto de negócio | Resultado de rodada | Pagamento de bônus | **Pagamento de jackpot** |
| Fonte PHP wrapper | ~94-97 | ~99-102 | **~104-107** |

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Endpoint `/jackpotWin` analisado — 9 regras BR-* confirmadas. Wrapper `jackpotWin()` (~104-107) e delegação para `handleResult()` (~161-175) documentados
- [ ] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [ ] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [ ] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [ ] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [ ] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation)
- [ ] **AC-7:** Arquivo criado em `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-jackpotWin.md`
- [ ] **AC-8:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-9:** Seção de contexto de negócio explicando quando `/jackpotWin` é chamado vs. `/result` e `/bonusWin` — distinção de jackpot como evento de alta magnitude financeira
- [ ] **AC-10:** Nota explícita de que `/jackpotWin` é membro da família handleResult() com referência cruzada para `pragmatic-play-result.md`

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
URL:    /v1/webhooks/pragmatic-play/jackpotWin
Função: Registrar pagamento de jackpot ganho pelo jogador
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php
        - Wrapper: método jackpotWin()  linhas ~104-107
        - Lógica:  método handleResult() linhas ~161-175
```

### Arquitetura: Wrapper → handleResult()

```php
// Thin wrapper público (linhas ~104-107)
public function jackpotWin($data) {
    return $this->handleResult('jackpotWin', $data);
}

// handleResult() é idêntico ao chamado por result() e bonusWin()
// Diferença: URL destino = {tenant_url}/pragmatic-play/jackpotWin.html
```

### Regras Aplicáveis (9 total — idênticas aos demais membros da família)

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
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/jackpotWin.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Referências de Código Fonte

```
PragmaticPlayService.php:104-107   → método jackpotWin() (wrapper)
PragmaticPlayService.php:161-175   → método handleResult() (lógica compartilhada)
PragmaticPlayService.php:161       → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137   → método removeTenant()
PragmaticPlayService.php:142-152   → método generateHashCode()
OperatorService.php:20-34          → método get()
BaseService.php:16-22              → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-jackpotWin.md` deve seguir **exatamente** o template de `pragmatic-play-result.md`, substituindo:
- "result" → "jackpotWin" em textos
- `result.html` → `jackpotWin.html` na URL de destino
- Contexto de negócio: "resultado de rodada" → "pagamento de jackpot"
- Linhas do wrapper: ~94-97 → ~104-107

### Seções Obrigatórias

1. **Header** — endpoint, função, nota sobre família handleResult()
2. **Resumo Executivo** — o que faz, quando é chamado, relação com família via handleResult()
3. **Fluxo em 8 Fases** — Mermaid + explicação
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5)
6. **Exemplo Completo** — request → response com passthrough
7. **Contexto de negócio** — distinção jackpotWin vs. result vs. bonusWin
8. **Checklist de Segurança**
9. **Limites e Restrições**

---

## Tasks / Subtasks

- [ ] **T-1:** Ler `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` — usar como template direto
- [ ] **T-2:** Ler `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — confirmar regras
- [ ] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `jackpotWin()` (~104-107) — confirmar delegação
- [ ] **T-4:** Criar arquivo `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-jackpotWin.md`
- [ ] **T-5:** Adaptar conteúdo de result.md — substituir referências result → jackpotWin
- [ ] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid
- [ ] **T-7:** Preencher Matriz de Regras (9 regras)
- [ ] **T-8:** Documentar 5+ Cenários de Erro
- [ ] **T-9:** Escrever exemplo completo request → response
- [ ] **T-10:** Adicionar seção contexto jackpotWin vs. result/bonusWin
- [ ] **T-11:** Preencher Security Checklist + Limites e Restrições
- [ ] **T-12:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low (adaptação direta de result.md)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): Markdown renderiza corretamente; sem referências residuais a "result" ou "bonusWin" onde deveria ser "jackpotWin"
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
- Consistência result/bonusWin → jackpotWin nas substituições de texto
- URL de destino correta: `jackpotWin.html`
- Linhas de código fonte corretas: wrapper ~104-107
- Contexto de negócio do jackpot presente e claro

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-jackpotWin.md` | Documentação técnica do endpoint /jackpotWin | ⏳ A Criar |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-jackpotWin.md` | Output principal desta story | ⏳ A Criar |
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` | Template direto a seguir | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [ ] Arquivo `pragmatic-play-jackpotWin.md` criado e completo
- [ ] Diagrama Mermaid renderiza corretamente
- [ ] `userId` como identificador documentado
- [ ] Passthrough de response documentado (Fase 8)
- [ ] URL de destino `jackpotWin.html`
- [ ] Linhas de código PHP corretas (wrapper ~104-107, handleResult() ~161-175)
- [ ] Contexto de negócio jackpot presente
- [ ] Security checklist preenchido
- [ ] File List desta story atualizada

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa fidelidade ao código PHP e contexto de negócio do jackpot.

---

## Métricas de Sucesso

- **Correção:** URL `jackpotWin.html`, linhas PHP corretas
- **Contexto:** Distinção de jackpot como evento de alta magnitude financeira clara
- **Rastreabilidade:** 9 regras mapeadas por fase

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1 hora (adaptação direta de result.md)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-result ✅ (Ready)
- **Bloqueia:** CASINO-2.2-promoWin
- **Sequência Fase 2:** authenticate ✅ → bet ✅ → refund ✅ → result ✅ → bonusWin → **jackpotWin** → promoWin → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Story completa; contexto jackpot como evento de alta magnitude financeira documentado. |
