# CASINO-2.2-result: Documentar Endpoint /result do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-result  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Done  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7, CASINO-2.2-refund (Ready), CASINO-2.2-bet (Ready)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/result` do Pragmatic Play seguindo o padrão estabelecido em `pragmatic-play-bet.md`. O `/result` é o **primeiro membro da família handleResult()** — usa `userId`, passthrough de response, 9 regras genéricas, sem regras exclusivas. A característica que o distingue de `/bet` e `/refund` é que **o método `result()` é um thin wrapper** que delega toda a lógica para o método compartilhado `handleResult()`, o mesmo usado por `/bonusWin`, `/jackpotWin` e `/promoWin`.

**Objetivo:** Produzir `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` completo — fluxo 8 fases, 9 regras mapeadas, documentação do padrão handleResult(), exemplos request/response e security checklist.

---

## Contexto

### Por que esta Story?

O `/result` fecha o loop de uma rodada de jogo: se `/bet` registra a aposta, `/result` registra o resultado/pagamento. É o endpoint mais frequente no fluxo de jogo (toda rodada termina com um `/result`).

Documentar `/result` corretamente é estratégico porque:
1. **Introduz o padrão handleResult()** — compartilhado com 3 endpoints subsequentes (bonusWin, jackpotWin, promoWin)
2. **Documenta a separação wrapper/lógica** — o wrapper público delega para implementação privada compartilhada
3. **Serve de referência** para as 3 stories seguintes, que são adaptações mínimas deste mesmo padrão

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance      ✅ (template)
  ├─ /authenticate ✅ Ready
  ├─ /bet          ✅ Ready
  ├─ /refund        ✅ Ready
  ├─ /result        ← ESTA STORY  (handleResult — primeiro membro)
  ├─ /bonusWin      ┐
  ├─ /jackpotWin    │ handleResult() family — mesmo padrão
  ├─ /promoWin      ┘
  └─ /adjustment
```

### Diferencial do /result — Arquitetura Wrapper + handleResult()

| Característica | /bet | /refund | **/result** |
|---------------|------|---------|------------|
| Identificador | `userId` | `userId` | `userId` |
| Regras exclusivas | Nenhuma | Nenhuma | Nenhuma |
| Transforma response | Não (passthrough) | Não (passthrough) | Não (passthrough) |
| Implementação | Inline (~14 linhas) | Inline (~14 linhas) | **Thin wrapper → handleResult()** |
| Compartilhado com | — | — | bonusWin, jackpotWin, promoWin |
| URL de destino | `.../bet.html` | `.../refund.html` | `.../result.html` |
| Contexto de negócio | Registra aposta | Estorna aposta | **Registra resultado/pagamento** |
| Fonte PHP wrapper | `bet()` linhas ~64-77 | `refund()` linhas ~79-92 | `result()` linhas ~94-97 |
| Fonte PHP lógica | — | — | `handleResult()` linhas ~161-175 |

> **O diferencial arquitetural do /result:** A lógica de negócio reside em `handleResult()` (método privado compartilhado), não no método público `result()`. Documentar esse padrão é essencial para que o handler Go implemente a mesma separação.

---

## Critérios de Aceitação

### Deve Ter

- [x] **AC-1:** Endpoint `/result` analisado — 9 regras BR-* confirmadas. Tanto o wrapper `result()` (~linhas 94-97) quanto a lógica compartilhada `handleResult()` (~linhas 161-175) documentados
- [x] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [x] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [x] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [x] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [x] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation)
- [x] **AC-7:** Arquivo criado em `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md`
- [x] **AC-8:** File List desta story atualizada

### Deveria Ter

- [x] **AC-9:** Seção sobre o padrão handleResult() — explicar a arquitetura wrapper → shared logic, quando é chamado, e que bonusWin/jackpotWin/promoWin são membros da mesma família
- [x] **AC-10:** Nota explícita de que `/result` é o **primeiro da família handleResult()**, servindo de referência canônica para os próximos 3 endpoints da série

### Fora do Escopo

- ❌ Escrever testes (Fase 3)
- ❌ Criar matrizes YAML (Fase 4)
- ❌ Implementar handler Go
- ❌ Documentar outros endpoints (/bonusWin, /jackpotWin, /promoWin, /adjustment)
- ❌ Corrigir código PHP

---

## Detalhes Técnicos / Dev Notes

### Endpoint

```
Método: POST
URL:    /v1/webhooks/pragmatic-play/result
Função: Registrar resultado/pagamento de uma rodada de jogo
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php
        - Wrapper: método result()   linhas ~94-97
        - Lógica:  método handleResult() linhas ~161-175
```

### Arquitetura: Wrapper → handleResult()

```php
// Thin wrapper público (linhas ~94-97)
public function result($data) {
    return $this->handleResult('result', $data);
}

// Lógica compartilhada (linhas ~161-175)
private function handleResult($endpoint, $data) {
    $tenant = $this->operatorService->get($data['userId']);
    $data['userId'] = $this->removeTenant($data['userId']);
    // ... credential lookup, hash, postJson ...
    return postJson("{tenant_url}/pragmatic-play/{$endpoint}.html", $data);
}
```

> **Nota Go:** O handler Go para `/result` deve replicar essa separação ou implementar handleResult() como função helper compartilhada invocada pelos 4 endpoints da família.

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
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/result.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Uso de userId em handleResult()

```php
// handleResult() — apenas userId (linha ~161 do PragmaticPlayService.php)
$tenant = $this->operatorService->get($data['userId']);
$data['userId'] = $this->removeTenant($data['userId']);
// ex: "myoperator_user456" → "user456"
```

**Código Fonte:** `PragmaticPlayService.php:161`

### Família handleResult() — 4 Endpoints

| Endpoint | Wrapper (linha aprox.) | handleResult() endpoint arg |
|----------|----------------------|----------------------------|
| `result` | ~94-97 | `'result'` |
| `bonusWin` | ~99-102 | `'bonusWin'` |
| `jackpotWin` | ~104-107 | `'jackpotWin'` |
| `promoWin` | ~109-112 | `'promoWin'` |

> Todos os 4 diferem apenas no argumento endpoint passado para `handleResult()`. A lógica é 100% idêntica — diferem apenas na URL de destino.

### Referências de Código Fonte

```
PragmaticPlayService.php:94-97     → método result() (wrapper)
PragmaticPlayService.php:161-175   → método handleResult() (lógica compartilhada)
PragmaticPlayService.php:161       → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137   → método removeTenant()
PragmaticPlayService.php:142-152   → método generateHashCode()
OperatorService.php:20-34          → método get()
BaseService.php:16-22              → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-result.md` deve seguir o template de `pragmatic-play-bet.md`, substituindo:
- "bet" → "result" em textos
- `bet.html` → `result.html` na URL de destino
- Contexto de negócio: "aposta" → "resultado/pagamento de rodada"
- Linhas do wrapper: ~64-77 → ~94-97
- **Adicionar seção sobre arquitetura handleResult()** (diferencial vs. bet/refund)

### Seções Obrigatórias

1. **Header** — endpoint, função, nota sobre família handleResult()
2. **Resumo Executivo** — o que faz, quando é chamado, arquitetura wrapper → handleResult()
3. **Fluxo em 8 Fases** — Mermaid + explicação (Fase 2 = userId; Fase 8 = passthrough)
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5):
   - `userId` faltando → exception em tenant extraction
   - `userId` sem underscore → falha
   - Operador não encontrado → `ModelNotFoundException`
   - Credencial faltando → null reference exception
   - Provider timeout → falha imediata
6. **Exemplo Completo** — request com `userId` prefixado → sanitização → POST → passthrough response
7. **Seção: Família handleResult()** — padrão arquitetural, 4 endpoints membros, tabela de wrappers
8. **Checklist de Segurança**
9. **Limites e Restrições**

---

## Tasks / Subtasks

> Sequência de implementação para @dev

- [x] **T-1:** Ler `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` — usar como template direto (estrutura idêntica para as 8 fases e regras)
- [x] **T-2:** Ler `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — focar na seção handleResult() (linha ~161) e na tabela de endpoints
- [x] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` métodos `result()` (~94-97) e `handleResult()` (~161-175) — confirmar arquitetura wrapper + shared logic
- [x] **T-4:** Criar arquivo `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md`
- [x] **T-5:** Adaptar conteúdo de bet.md — substituir referências bet → result, ajustar URLs e linhas de código
- [x] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid
- [x] **T-7:** Preencher Matriz de Regras (9 regras)
- [x] **T-8:** Documentar 5+ Cenários de Erro
- [x] **T-9:** Escrever exemplo completo request → response
- [x] **T-10:** Adicionar seção "Família handleResult()" — tabela dos 4 endpoints, arquitetura wrapper → shared logic, nota para implementação Go
- [x] **T-11:** Preencher Security Checklist + Limites e Restrições
- [x] **T-12:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low-Medium (adaptação de bet.md + seção nova sobre handleResult())
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect

**Quality Gate Tasks:**
- [x] Pre-Commit (@dev): [N/A — CodeRabbit disabled] — Markdown renderiza corretamente; sem referências residuais a "bet" onde deveria ser "result"; handleResult() documentado
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
- Consistência bet → result nas substituições de texto
- URL de destino correta: `result.html` (não `bet.html` ou `refund.html`)
- Linhas de código fonte corretas: wrapper ~94-97, handleResult() ~161-175
- Seção handleResult() presente e clara (diferencial desta story vs. bet/refund)
- Contexto de negócio do resultado/pagamento presente e claro

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` | Documentação técnica do endpoint /result | ✅ Criado |

### Template de Referência

```
docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md  ← principal
docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md
```

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` | Output principal desta story | ✅ Criado |
| `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* e handleResult() | ✅ Existe |
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bet.md` | Template direto a seguir | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [x] Arquivo `pragmatic-play-result.md` criado e completo
- [x] Diagrama Mermaid renderiza corretamente
- [x] `userId` como identificador documentado (não `token`)
- [x] Passthrough de response documentado (Fase 8)
- [x] Seção família handleResult() presente com tabela dos 4 endpoints membros
- [x] URL de destino `result.html` (não `bet.html`)
- [x] Linhas de código PHP corretas (wrapper ~94-97, handleResult() ~161-175)
- [x] Security checklist preenchido
- [x] File List desta story atualizada

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa fidelidade ao código PHP e contexto de negócio do resultado.  
**Próxima Fase:** Fase 3 (CASINO-2.3) criará testes Java para as regras documentadas.

---

## Métricas de Sucesso

- **Correção:** URL `result.html`, linhas PHP corretas, sem erros de cópia de bet.md
- **Arquitetura:** handleResult() documentado como padrão — reutilizável para bonusWin/jackpotWin/promoWin
- **Rastreabilidade:** 9 regras mapeadas por fase
- **Clareza:** @dev Go consegue implementar o handler sem consultar PHP

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1-2 horas (adaptação de bet.md + seção handleResult() adicional)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-refund ✅ (Ready)
- **Bloqueia:** CASINO-2.2-bonusWin (próximo endpoint da família)
- **Sequência Fase 2:** authenticate ✅ → bet ✅ → refund ✅ → **result** → bonusWin → jackpotWin → promoWin → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Padrão handleResult() bem articulado; referência canônica para família bonusWin/jackpotWin/promoWin. |
| 2026-05-14 | @dev (Dex) | Implementação completa — `pragmatic-play-result.md` criado (9 seções, 6 cenários de erro, exemplo completo, seção família handleResult() com PHP side-by-side e snippet Go). Todos T-1..T-12 concluídos. Status: InReview. |

| 2026-05-15 | @qa (Quinn) | QA Gate PASS — Todos os ACs atendidos, output file completo, CodeRabbit N/A. Status: InReview → Done |

## QA Results

### Review Date: 2026-05-15

### Reviewed By: Quinn (Test Architect)

All ACs verified complete. Output file exists and passes documentation quality checks. CodeRabbit disabled (N/A).

### Gate Status

Gate: PASS → docs/qa/gates/casino-2.2-result.yml
