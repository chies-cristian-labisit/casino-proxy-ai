# CASINO-2.2-authenticate: Documentar Endpoint /authenticate do Pragmatic Play — Fase 2

**Story ID:** CASINO-2.2-authenticate  
**Epic:** CASINO-2 (Business Rules Discovery & Test Oracle)  
**Tipo:** Documentação Técnica (Fase 2 de 5 — Technical Documentation)  
**Status:** Done  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.7 (Regras de Negócio Pragmatic Play extraídas), CASINO-2.2 (Template /balance)  
**Data de Criação:** 2026-05-12  

---

## Resumo da Story

Documentar o endpoint `/authenticate` do Pragmatic Play seguindo o template estabelecido em `pragmatic-play-balance.md`. Este é o **primeiro dos 8 endpoints restantes** da Fase 2 (após `/balance`), e também o **mais crítico**: é o único endpoint que **transforma a resposta** do provider (re-prefixação do `userId`), diferente de todos os outros que fazem passthrough direto.

**Objetivo:** Produzir `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` completo — fluxo 8 fases, 9 regras mapeadas, exemplos request/response e security checklist.

---

## Contexto

### Por que esta Story?

A Fase 2 documenta o **comportamento técnico de cada endpoint individualmente** — como as regras de negócio (BR-*) extraídas na Fase 1 se combinam para formar o fluxo completo de processamento.

O `/authenticate` é estrategicamente o primeiro a documentar porque:
1. **É único:** possui 2 regras exclusivas (PP-007, PP-012) que nenhum outro endpoint tem
2. **É crítico para Go:** o novo handler Go deve replicar a re-prefixação do `userId` corretamente
3. **Serve de referência:** expõe as diferenças entre endpoints com e sem transformação de response

### Como se Encaixa no Plano

```
Fase 1: Extrair regras ✅ (CASINO-1.7 — 12 regras BR-* documentadas)
Fase 2: Documentar endpoints ← VOCÊ ESTÁ AQUI
  ├─ /balance ✅ (template criado em pragmatic-play-balance.md)
  ├─ /authenticate ← ESTA STORY
  ├─ /bet
  ├─ /refund
  ├─ /result
  ├─ /bonusWin
  ├─ /jackpotWin
  ├─ /promoWin
  └─ /adjustment
Fase 3: Construir Test Oracle (Java + WireMock)
Fase 4: Criar matrizes YAML de rastreamento
Fase 5: Validar — 100% testes PHP passam
```

### Diferencial do /authenticate

> ⚠️ **Este é o único endpoint do Pragmatic Play que NÃO faz passthrough da resposta.**
>
> Todos os outros endpoints (balance, bet, refund, etc.) retornam a resposta do provider inalterada.
> O `/authenticate` **modifica a resposta**: re-prefixa o `userId` com o `operator_slug` (ex: `"12345"` → `"myoperator_12345"`), mas **apenas se `error == 0`** (sucesso).
>
> Regras envolvidas: **PP-007** (mecanismo de re-prefixação) e **PP-012** (autenticação é o único endpoint com essa transformação).

---

## Critérios de Aceitação

### Deve Ter

- [x] **AC-1:** Endpoint `/authenticate` analisado — 9 regras BR-* identificadas e confirmadas contra `PragmaticPlayService.php`
- [x] **AC-2:** Fluxo de 8 fases documentado com diagrama Mermaid renderizável
- [x] **AC-3:** Regras genéricas (7) mapeadas às fases corretas, com destaque para PP-007 e PP-012 na Fase 8 (exclusivas de authenticate)
- [x] **AC-4:** Mínimo 5 cenários de erro documentados com causa raiz e comportamento esperado
- [x] **AC-5:** Exemplo completo request → response mostrando a re-prefixação do `userId` (caso sucesso) e ausência de re-prefixação (caso erro)
- [x] **AC-6:** Security checklist preenchido (tenant isolation, hash auth, operator/credential validation, endpoint validation)
- [x] **AC-7:** Arquivo criado em `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md`
- [x] **AC-8:** File List desta story atualizada

### Deveria Ter

- [x] **AC-9:** Tabela comparativa `authenticate vs balance` destacando as diferenças (response passthrough vs transformação, token único vs dual token)
- [x] **AC-10:** Documentar os casos extremos de PP-007 — o que acontece se `userId` estiver ausente na resposta do provider, ou se `error` field estiver faltando

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
URL:    /v1/webhooks/pragmatic-play/authenticate
Função: Autenticar jogador e retornar userId prefixado com tenant
Fonte:  legacy/casino-proxy/app/Services/PragmaticPlayService.php (método authenticate(), linhas ~26-44)
```

### Regras Aplicáveis (9 total)

| # | ID | Descrição | Fase no Fluxo | Exclusiva? |
|---|-----|-----------|--------------|------------|
| 1 | BR-GENERIC-ROUTING-VALIDATION-001 | Dynamic endpoint routing via method resolution | 1 | Não |
| 2 | BR-GENERIC-ERROR-HANDLING-001 | Unknown endpoint → exception 500 | 1 (guard) | Não |
| 3 | BR-GENERIC-TENANT-EXTRACTION-001 | Extrair operator_slug do token (`token.split('_')`) | 2 | Não |
| 4 | BR-GENERIC-OPERATOR-CACHING-001 | Operator lookup com cache 1h (`cache_key = 'operator_' + slug`) | 3 | Não |
| 5 | BR-PRAGMATIC-BALANCE-TOKEN-SANITIZATION-001 | Remover prefixo tenant do token antes de enviar ao provider | 4 | Não |
| 6 | BR-GENERIC-CREDENTIAL-LOOKUP-001 | Buscar secret-key do operador (`credentials.where('key','secret-key').first()`) | 5 | Não |
| 7 | BR-GENERIC-AUTHENTICATION-HMAC-MD5-001 | Gerar hash MD5 (sort payload + concat secret + md5) | 6 | Não |
| 8 | BR-GENERIC-PROVIDER-INTEGRATION-001 | HTTP POST para `{tenant_url}/pragmatic-play/authenticate.html` | 7 | Não |
| 9 | **PP-012** (inclui PP-007) | **Re-prefixa `userId` se `error == 0`** — ÚNICO endpoint com transformação | **8** | **SIM** |

> **Fonte das regras:** `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`

### Comportamento de Re-prefixação (PP-007 / PP-012)

```
response = postJson(url, payload)

SE response['error'] == 0:   // sucesso
  response['userId'] = operator_slug + '_' + response['userId']
  // ex: "12345" → "myoperator_12345"

// SE error != 0: response retornada inalterada (sem prefixo)

RETURN response
```

**Código Fonte:** `PragmaticPlayService.php:40-42`

### Nota sobre Token

O `/authenticate` usa **apenas `token`** (sem suporte dual token como `/balance`).

```php
// authenticate() — apenas token
$tenant = $this->operatorService->get($data['token']);
$data['token'] = $this->removeTenant($data['token']);
```

Comparação:
```
/balance:  aceita token OU userId  (BR-PRAGMATIC-BALANCE-DUAL-TOKEN-SUPPORT-001)
/authenticate: apenas token
```

### Referências de Código Fonte

```
PragmaticPlayService.php:26-44  → método authenticate()
PragmaticPlayService.php:132-137 → método removeTenant()
PragmaticPlayService.php:142-152 → método generateHashCode()
OperatorService.php:20-34       → método get() (tenant extraction + cache)
BaseService.php:16-22           → método postJson()
```

---

## Estrutura do Documento de Output

O arquivo `pragmatic-play-authenticate.md` deve seguir exatamente o template de `pragmatic-play-balance.md`, com as adaptações abaixo:

### Seções Obrigatórias

1. **Header** — endpoint, função, nota de destaque sobre response transformation
2. **Resumo Executivo** — o que faz, quando é chamado, por que é único (re-prefixação)
3. **Fluxo em 8 Fases** — diagrama Mermaid + explicação detalhada por fase  
   _(Fase 8 é diferente do /balance: transformação em vez de passthrough)_
4. **Matriz de Regras** — tabela: regra × fase × impacto × exclusiva?
5. **Cenários de Erro** (mínimo 5):
   - Token faltando / sem underscore → exception em tenant extraction
   - Operador não encontrado → `ModelNotFoundException`
   - Credencial Pragmatic faltando → null reference exception
   - Provider timeout → falha imediata (sem retry)
   - Hash inválido → provider retorna error != 0
6. **Exemplo Completo** — 2 examples: sucesso (com re-prefixação) + falha (sem re-prefixação)
7. **Comparação authenticate vs balance** — tabela de diferenças-chave
8. **Checklist de Segurança** — tenant isolation, hash auth, operator/credential validation
9. **Limites e Restrições** — cache TTL, formato token, algoritmo hash, retry policy

---

## Tasks / Subtasks

> Sequência de implementação para @dev

- [x] **T-1:** Ler `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` completo — internalizar estrutura do template
- [x] **T-2:** Ler `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` — focar nas regras PP-007 e PP-012 (linhas `authenticate()` seção)
- [x] **T-3:** Ler `legacy/casino-proxy/app/Services/PragmaticPlayService.php` método `authenticate()` (linhas ~26-44) — confirmar fluxo real
- [x] **T-4:** Criar arquivo `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md`
- [x] **T-5:** Escrever Header + Resumo Executivo (destacar diferencial de response transformation)
- [x] **T-6:** Escrever Fluxo 8 Fases com diagrama Mermaid (Fase 8 = transformação, não passthrough)
- [x] **T-7:** Preencher Matriz de Regras (9 regras × fase × exclusiva?)
- [x] **T-8:** Documentar 5+ Cenários de Erro com causa raiz e comportamento esperado
- [x] **T-9:** Escrever 2 exemplos completos request → response (sucesso com prefixo + erro sem prefixo)
- [x] **T-10:** Adicionar tabela comparativa `authenticate vs balance`
- [x] **T-11:** Preencher Security Checklist + Limites e Restrições
- [x] **T-12:** Atualizar File List desta story

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Low (adaptação de template existente)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect (revisar fidelidade das regras documentadas)

**Quality Gate Tasks:**
- [x] Pre-Commit (@dev): [N/A — CodeRabbit disabled] — Verificar que markdown renderiza corretamente (headings, tabelas, Mermaid)
- [x] Pre-PR (@devops): [N/A — CodeRabbit disabled] — Validar links e referências a arquivos existentes

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
- Referências a arquivos existentes (paths válidos)
- Consistência com template `pragmatic-play-balance.md`
- Fidelidade técnica das regras BR-* documentadas

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` | Documentação técnica do endpoint /authenticate | ✅ Criado |

### Template de Referência

```
docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md
```

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-authenticate.md` | Output principal desta story | ✅ Criado |
| `docs/architecture/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Fonte das 9 regras BR-* | ✅ Existe |
| `docs/architecture/casino-proxy/phase-2-technical-documentation/pragmatic-play-balance.md` | Template a seguir | ✅ Existe |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte de referência | ✅ Existe |

---

## Definição de Pronto

- [x] Arquivo `pragmatic-play-authenticate.md` criado e completo
- [x] Diagrama Mermaid renderiza corretamente (8 fases visíveis)
- [x] PP-007 e PP-012 documentadas como exclusivas do authenticate com exemplos
- [x] Dois exemplos request/response: sucesso (userId prefixado) e erro (userId não prefixado)
- [x] Tabela comparativa authenticate vs balance presente
- [x] Security checklist preenchido
- [x] File List desta story atualizada
- [x] Pronto para validação de @po antes de documentar próximo endpoint (/bet)

---

## Estratégia de Teste

**Esta story:** Apenas documentação, sem código ou testes.  
**Validação:** @po revisa documento para confirmar fidelidade às regras extraídas na Fase 1.  
**Próxima Fase:** Fase 3 (CASINO-2.3) criará testes Java para as regras documentadas aqui.

---

## Métricas de Sucesso

- **Completude:** Todas as 9 regras documentadas no fluxo
- **Clareza:** @dev consegue ler o documento e implementar handler Go sem consultar código PHP
- **Rastreabilidade:** Cada fase do fluxo referencia a regra BR-* correspondente
- **Correção:** Exemplos request/response correspondem ao comportamento real (validado contra código PHP)
- **Unicidade:** Diferencial da re-prefixação claramente explicado e distinguido dos demais endpoints

---

## Notas

- **Criado:** 2026-05-12
- **Estimado:** 2-3 horas (leitura do template + adaptação + escrita)
- **Depende De:** CASINO-1.7 (regras extraídas ✅), CASINO-2.2 /balance (template ✅)
- **Bloqueia:** CASINO-2.2-bet (próximo endpoint a documentar)
- **Sequência completa Fase 2:** authenticate → bet → refund → result → bonusWin → jackpotWin → promoWin → adjustment

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-12 | @sm (River) | Story criada — Draft |
| 2026-05-12 | @po (Pax) | Validação GO (8/10) — Status: Draft → Ready. Adicionadas seções Tasks/Subtasks e CodeRabbit Integration. |
| 2026-05-12 | @dev (Dex) | Implementação completa — `pragmatic-play-authenticate.md` criado (10 seções, 6 cenários de erro, 2 exemplos, tabela comparativa, Mermaid 8 fases). Todos T-1..T-12 concluídos. Status: Ready → InProgress → aguarda QA. |
| 2026-05-12 | @qa (Quinn) | QA Gate PASS — 7/7 checks OK. Todos os ACs (10/10) atendidos. Observação minor: atualizações administrativas da story aplicadas (checkboxes, File List, DoD). Status: InReview. |
| 2026-05-15 | @qa (Quinn) | QA Gate PASS (retroativo) — Gate file criado. Status: InReview → Done | @qa |

## QA Results

### Review Date: 2026-05-15

### Reviewed By: Quinn (Test Architect)

All 10 ACs verified as complete. Output file exists and is complete. Prior QA PASS (2026-05-12) confirmed valid.

### Gate Status

Gate: PASS → docs/qa/gates/casino-2.2-authenticate.yml
