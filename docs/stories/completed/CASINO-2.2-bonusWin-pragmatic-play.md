# CASINO-2.2-bonusWin: Documentar Endpoint /bonusWin do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-bonusWin  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Done  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7, CASINO-2.2-result (Ready)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/bonusWin` do Pragmatic Play seguindo o padrão estabelecido em `pragmatic-play-result.md`. O `/bonusWin` é o **segundo membro da família handleResult()** — estruturalmente idêntico ao `/result`, diferindo apenas no argumento `'bonusWin'` passado para `handleResult()` e na URL de destino (`bonusWin.html`).

**Objetivo:** Produzir `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bonusWin.md` completo — fluxo 8 fases, 9 regras mapeadas, exemplos request/response e security checklist.

---

## Contexto

### Por que esta Story?

O `/bonusWin` registra o pagamento de um bônus ganho pelo jogador. É um evento distinto de `/result` (que registra o resultado padrão de uma rodada) mas compartilha exatamente a mesma implementação via `handleResult()`.

Documentar `/bonusWin` como story separada garante:
1. **Rastreabilidade por endpoint** — cada URL tem seu próprio documento técnico
2. **Contexto de negócio diferenciado** — bônus têm regras de negócio do operador distintas (frequência, limites) mesmo que a implementação PHP seja idêntica
3. **Completude para o handler Go** — o desenvolvedor Go precisa saber que `/bonusWin` é um endpoint real e independente

### Como se Encaixa no Plano

```
Fase 2: Documentar endpoints
  ├─ /balance      ✅ (template)
  ├─ /authenticate ✅ Ready
  ├─ /bet          ✅ Ready
  ├─ /refund       ✅ Ready
  ├─ /result       ✅ Ready (introduz handleResult())
  ├─ /bonusWin     ← ESTA STORY  (handleResult — segundo membro)
  ├─ /jackpotWin   ┐ handleResult() family — mesmo padrão
  ├─ /promoWin     ┘
  └─ /adjustment
```

### Diferencial do /bonusWin

| Característica | /result | **/bonusWin** |
|---------------|---------|--------------|
| Identificador | `userId` | `userId` |
| Regras exclusivas | Nenhuma | Nenhuma |
| Transforma response | Não (passthrough) | Não (passthrough) |
| Implementação | Thin wrapper → handleResult('result') | Thin wrapper → **handleResult('bonusWin')** |
| URL de destino | `.../result.html` | `.../bonusWin.html` |
| Contexto de negócio | Resultado de rodada | **Pagamento de bônus** |
| Fonte PHP wrapper | `result()` linhas ~94-97 | `bonusWin()` linhas ~99-102 |
| Lógica compartilhada | `handleResult()` ~161-175 | `handleResult()` ~161-175 |

> O `/bonusWin` é tecnicamente o endpoint mais simples da série — adaptação direta de result.md com apenas substituição do nome do endpoint e contexto de negócio.

---

## Critérios de Aceitação

### Deve Ter

- [x] **AC-1:** Endpoint `/bonusWin` analisado — 9 regras BR-* confirmadas. Wrapper `bonusWin()` (~99-102) e delegação para `handleResult()` (~161-175) documentados
- [x] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [x] **AC-3:** 9 regras mapeadas às fases corretas — destaque para uso de `userId` (Fase 2) e passthrough de response (Fase 8)
- [x] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [x] **AC-5:** Exemplo completo request → response mostrando sanitização do `userId` e passthrough da resposta
- [x] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation)
- [x] **AC-7:** Arquivo criado em `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bonusWin.md`
- [x] **AC-8:** File List desta story atualizada

### Deveria Ter

- [x] **AC-9:** Seção de contexto de negócio explicando a diferença semântica entre `/result` (resultado de rodada) e `/bonusWin` (pagamento de bônus), e quando cada um é chamado
- [x] **AC-10:** Nota explícita de que `/bonusWin` é membro da família handleResult() com referência cruzada para `pragmatic-play-result.md`

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
URL:    /v1/webhooks/pragmatic-play/bonusWin
Função: Registrar pagamento de bônus ganho pelo jogador
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php
        - Wrapper: método bonusWin()   linhas ~99-102
        - Lógica:  método handleResult() linhas ~161-175
```

### Arquitetura: Wrapper → handleResult()

```php
// Thin wrapper público (linhas ~99-102)
public function bonusWin($data) {
    return $this->handleResult('bonusWin', $data);
}

// handleResult() é idêntico ao chamado por result() — linha ~161
$data['userId'] = $this->removeTenant($data['userId']);
// ... idêntico ao /result, apenas URL muda: bonusWin.html
```

### Regras Aplicáveis (9 total — idênticas ao /result)

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
| 9 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/bonusWin.html` | 7 | Não |

> **Fase 8:** Passthrough direto — resposta do provider retornada **sem nenhuma transformação**.  
> **Fonte das regras:** `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Referências de Código Fonte

```
PragmaticPlayService.php:99-102    → método bonusWin() (wrapper)
PragmaticPlayService.php:161-175   → método handleResult() (lógica compartilhada)
PragmaticPlayService.php:161       → $data['userId'] = removeTenant($data['userId'])
PragmaticPlayService.php:132-137   → método removeTenant()
PragmaticPlayService.php:142-152   → método generateHashCode()
OperatorService.php:20-34          → método get()
BaseService.php:16-22              → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-bonusWin.md` deve seguir **exatamente** o template de `pragmatic-play-result.md`, substituindo:
- "result" → "bonusWin" em textos
- `result.html` → `bonusWin.html` na URL de destino
- Contexto de negócio: "resultado de rodada" → "pagamento de bônus"
- Linhas do wrapper: ~94-97 → ~99-102
- Seção handleResult(): manter referência cruzada para result.md

### Seções Obrigatórias

1. **Header** — endpoint, função, nota sobre família handleResult()
2. **Resumo Executivo** — o que faz, quando é chamado, relação com /result via handleResult()
3. **Fluxo em 8 Fases** — Mermaid + explicação (Fase 2 = userId; Fase 8 = passthrough)
4. **Matriz de Regras** — 9 regras × fase × exclusiva?
5. **Cenários de Erro** (mínimo 5):
   - `userId` faltando → exception em tenant extraction
   - `userId` sem underscore → falha
   - Operador não encontrado → `ModelNotFoundException`
   - Credencial faltando → null reference exception
   - Provider timeout → falha imediata
6. **Exemplo Completo** — request com `userId` prefixado → sanitização → POST → passthrough response
7. **Contexto result vs bonusWin** — distinção semântica de negócio
8. **Checklist de Segurança**
9. **Limites e Restrições**

---

## Tasks / Subtasks

- [x] **T-1:** Ler `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` — usar como template direto (estrutura idêntica)
- [x] **T-2:** Ler `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — confirmar regras e tabela de endpoints
- [x] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `bonusWin()` (~99-102) — confirmar delegação para handleResult()
- [x] **T-4:** Criar arquivo `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bonusWin.md`
- [x] **T-5:** Adaptar conteúdo de result.md — substituir referências result → bonusWin, ajustar URLs e linhas
- [x] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid
- [x] **T-7:** Preencher Matriz de Regras (9 regras)
- [x] **T-8:** Documentar 5+ Cenários de Erro
- [x] **T-9:** Escrever exemplo completo request → response
- [x] **T-10:** Adicionar seção contextual result vs bonusWin
- [x] **T-11:** Preencher Security Checklist + Limites e Restrições
- [x] **T-12:** Atualizar File List desta story

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
- [x] Pre-Commit (@dev): [N/A — CodeRabbit disabled] — Markdown renderiza corretamente; sem referências residuais a "result" onde deveria ser "bonusWin"
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
- Consistência result → bonusWin nas substituições de texto
- URL de destino correta: `bonusWin.html` (não `result.html`)
- Linhas de código fonte corretas: wrapper ~99-102
- Contexto de negócio do bônus presente e claro

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bonusWin.md` | Documentação técnica do endpoint /bonusWin | ✅ Criado |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-bonusWin.md` | Output principal desta story | ✅ Criado |
| `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-result.md` | Template direto a seguir | ✅ Existe (Ready) |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [x] Arquivo `pragmatic-play-bonusWin.md` criado e completo
- [x] Diagrama Mermaid renderiza corretamente
- [x] `userId` como identificador documentado
- [x] Passthrough de response documentado (Fase 8)
- [x] URL de destino `bonusWin.html`
- [x] Linhas de código PHP corretas (wrapper ~99-102, handleResult() ~161-175)
- [x] Contexto result vs bonusWin presente
- [x] Security checklist preenchido
- [x] File List desta story atualizada

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa fidelidade ao código PHP e contexto de negócio do bônus.  
**Próxima Fase:** Fase 3 (CASINO-2.3) criará testes Java para as regras documentadas.

---

## Métricas de Sucesso

- **Correção:** URL `bonusWin.html`, linhas PHP corretas, sem erros de cópia de result.md
- **Contexto:** Distinção semântica result vs bonusWin clara para implementação Go
- **Rastreabilidade:** 9 regras mapeadas por fase

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 1 hora (adaptação direta de result.md — menor esforço possível)
- **Depende De:** CASINO-1.7 ✅, CASINO-2.2-result ✅ (Ready)
- **Bloqueia:** CASINO-2.2-jackpotWin
- **Sequência Fase 2:** authenticate ✅ → bet ✅ → refund ✅ → result ✅ → **bonusWin** → jackpotWin → promoWin → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (9/10) — Status: Draft → Ready. Story completa; contexto de negócio bônus documentado. |
| 2026-05-14 | @dev (Dex) | Implementação completa — `pragmatic-play-bonusWin.md` criado (9 seções, 6 cenários de erro, exemplo completo, tabela comparativa result vs bonusWin). Todos T-1..T-12 concluídos. Status: InReview. |

| 2026-05-15 | @qa (Quinn) | QA Gate PASS — Todos os ACs atendidos, output file completo, CodeRabbit N/A. Status: InReview → Done |

## QA Results

### Review Date: 2026-05-15

### Reviewed By: Quinn (Test Architect)

All ACs verified complete. Output file exists and passes documentation quality checks. CodeRabbit disabled (N/A).

### Gate Status

Gate: PASS → docs/qa/gates/casino-2.2-bonuswin.yml
