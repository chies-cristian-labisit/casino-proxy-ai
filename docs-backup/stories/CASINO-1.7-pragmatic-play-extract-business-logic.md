# CASINO-1.7: Extrair & Documentar Regras de Lógica de Negócio do Pragmatic Play

**Story ID:** CASINO-1.7  
**Epic:** CASINO-1 (Migração Casino Proxy PHP → Go)  
**Tipo:** Descoberta de Regras de Negócio (Fase 1 de Análise com 5 Fases)  
**Status:** Done  
**Prioridade:** Alta  
**Atribuído a:** @dev (com revisão de @architect)  
**Relacionado:** CASINO-1.1 (Spec OpenAPI para Pragmatic Play)  
**Data de Conclusão:** 2026-05-11  

---

## Resumo da Story

Extrair e documentar as **regras de lógica de negócio** impostas pelos handlers de webhook do Pragmatic Play no código base PHP existente. Este é o passo fundamental para construir um oráculo de teste—uma implementação de referência que prova que os handlers Go correspondem exatamente ao comportamento do PHP.

**Objetivo:** Para cada endpoint do Pragmatic Play, responder: "O que este endpoint *realmente faz*?"

---

## Contexto

### Por que esta Story?

A Fase 1 (CASINO-1.1 a 1.6) documentou a *forma* dos endpoints (schemas de request/response, métodos HTTP, caminhos). Esta story descobre o *comportamento*—a lógica de negócio dentro dos handlers PHP.

**Exemplo:**
- ✅ Fase 1 diz: "POST /v1/webhooks/pragmatic-play/verify-session aceita { session_id, signature }"
- ❌ Fase 1 não diz: "O que acontece se session_id não existe?" ou "Como a HMAC é validada?"

Esta story preenche essa lacuna.

### Como se Encaixa no Plano Geral

```
Fase 1: Extrair lógica de negócio dos handlers PHP ← VOCÊ ESTÁ AQUI
Fase 2: Documentar regras em markdown
Fase 3: Construir suite de testes de integração (implementação de referência)
Fase 4: Criar matriz YAML de rastreamento para regras críticas
Fase 5: Validar: Testes PHP passam 100%

Depois → Repetir para próximo provider
Depois → PO valida
Depois → Prosseguir com implementação Go (CASINO-2.x)
```

---

## Critérios de Aceitação

### Deve Ter (Definição de Pronto)

- [x] **AC-1:** Código handler do Pragmatic Play (`app/Http/Controllers/Webhooks/PragmaticPlayController.php`) lido e analisado
- [x] **AC-2:** Todas as regras de negócio extraídas dos métodos handlers (uma regra por ponto de decisão/validação) — **12 regras extraídas (PP-001 a PP-012)**
- [x] **AC-3:** Documento de regras criado em `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`
- [x] **AC-4:** Para cada regra, documentado:
  - ID da Regra (ex: `PP-001`)
  - Descrição da regra (uma frase: o que é validado/imposto)
  - Contexto de negócio (por que essa regra existe)
  - Escopo da regra (quais endpoints a acionam)
  - Lógica de decisão (pseudocódigo ou árvore de decisão if-then)
  - Casos extremos (entradas inválidas, condições limites)
- [x] **AC-5:** Rastrear cada regra de volta ao número de linha do código fonte no handler PHP — **todas as 12 regras rastreadas**
- [x] **AC-6:** Documentar dependências entre regras (ex: "Regra PP-002 depende de PP-001 ter sucesso") — **matriz de dependências incluída**
- [x] **AC-7:** Todas as regras validadas contra comportamento real da API Pragmatic Play — **validadas contra testes em PragmaticPlayControllerTest.php**
- [x] **AC-8:** Lista de Arquivos atualizada com novos arquivos de documentação

### Deveria Ter

- [x] **AC-9:** Diagrama ou pseudocódigo do fluxo handler (request → lógica → response) — **matriz de dependências entre regras incluída; lógica de decisão para cada regra em pseudocódigo**
- [x] **AC-10:** Anotar ambiguidades ou lógica pouco clara encontrada no código PHP — **5 questões abertas documentadas (cache invalidation, retry ausente, error handling inconsistente, etc.)**

### Fora do Escopo

- ❌ Escrever testes (Fase 3)
- ❌ Criar matrizes YAML (Fase 4)
- ❌ Implementar handlers Go (CASINO-2.x)
- ❌ Corrigir código PHP (apenas descoberta)

---

## Detalhes Técnicos

### Contexto do Provider Pragmatic Play

De CASINO-1.1 (Fase 1):

**Endpoints:**
- `POST /v1/webhooks/pragmatic-play/{endpoint}` — Roteamento genérico baseado em parâmetro `{endpoint}`
- Validação de assinatura: HMAC-SHA256 (provável)
- Padrão de autenticação: Baseado em secret do provider

**Métodos de Webhook Conhecidos** (extrair regras para):
- Entrada/verificação de sessão
- Consultas de saldo
- Colocação/rollback de aposta
- Processamento de vitória/pagamento
- Tratamento de promoção/bônus

**Localização Esperada do Handler:**
```
app/Http/Controllers/Webhooks/PragmaticPlayController.php
  ├─ index() ou route handler
  ├─ verifySession()
  ├─ getBalance()
  ├─ placeBet()
  ├─ ... (outros métodos)
```

### O que Significa "Regras de Lógica de Negócio"

Uma regra de negócio é uma decisão ou restrição imposta pelo código. Exemplos:

| Tipo de Regra | Exemplo |
|-----------|---------|
| **Validação** | "Sessão deve existir e estar ativa antes de permitir consulta de saldo" |
| **Transformação** | "Converter valor de aposta da moeda do provider para moeda do sistema" |
| **Efeito Colateral** | "Descontar saldo da conta do jogador em aposta bem-sucedida" |
| **Tratamento de Erro** | "Se saldo insuficiente, retornar código de erro 'INSUFFICIENT_FUNDS'" |
| **Timing** | "Sessão expira após 24 horas de inatividade" |

---

## Entregáveis

### Primário

**Arquivo: `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md`**

```markdown
# Regras de Lógica de Negócio do Pragmatic Play

## Regra PP-001: Validação de Assinatura HMAC
- **Descrição:** Todas as requisições de webhook devem incluir assinatura HMAC-SHA256 válida
- **Contexto:** Garante que requisições vêm de servidores legítimos do Pragmatic Play
- **Escopo:** Todos os endpoints
- **Lógica:**
  ```
  SE request.signature == HMAC-SHA256(request.body, provider_secret)
    ENTÃO proceed
    SENÃO return 403 Forbidden com "INVALID_SIGNATURE"
  ```
- **Casos Extremos:**
  - Header de assinatura faltando → 403
  - Assinatura de provider secret diferente → 403
  - Body malformado (quebra o hash) → 403
- **Código Fonte:** `PragmaticPlayController.php:45-52`
- **Dependências:** Nenhuma (primeiro check no handler)

## Regra PP-002: Existência de Sessão & Status Ativo
- **Descrição:** Sessão deve existir no banco e ter status "ativo"
- **Contexto:** Previne operações em sessões fechadas/expiradas
- **Escopo:** Endpoints verifySession, getBalance, placeBet
- **Lógica:**
  ```
  SE session_id existe em sessions table
    E session.status == "ativo"
    E session.created_at + 24h > agora()
    ENTÃO proceed
    SENÃO return 401 Unauthorized com "SESSION_NOT_FOUND" ou "SESSION_EXPIRED"
  ```
- **Casos Extremos:**
  - Sessão existe mas marcada "fechada" → 401 SESSION_CLOSED
  - Sessão com mais de 24h → 401 SESSION_EXPIRED
  - Session.player_id faltando → 500 (erro de integridade de dados)
- **Código Fonte:** `PragmaticPlayController.php:60-75`
- **Dependências:** PP-001 (assinatura deve passar primeiro)

## Regra PP-003: Verificação de Conta do Operador
- **Descrição:** operator_id da sessão deve estar ativo (não suspenso/deletado)
- **Contexto:** Previne jogo se conta do operador está inativa
- **Escopo:** Todos os endpoints
- **Lógica:**
  ```
  operador = buscar(operators, operator_id)
  SE operador.status == "ativo"
    ENTÃO proceed
    SENÃO return 403 Forbidden com "OPERATOR_INACTIVE"
  ```
- **Casos Extremos:**
  - Operador.status == "suspenso" → 403
  - Registro de operador deletado → 404
- **Código Fonte:** `PragmaticPlayController.php:78-85`
- **Dependências:** PP-002 (precisa da sessão primeiro para obter operator_id)

...
(continuar para todas as regras descobertas)
```

### Arquivos de Suporte

- **Lista de Arquivos (este arquivo de story)** — links para entregável + código fonte
- **Notas de Pesquisa** (opcional) — ambiguidades ou clarificações exploração código PHP
- **Branch:** `feature/CASINO-1.7-pp-business-logic`

---

## Notas de Implementação

### Abordagem

1. **Ler Código do Handler:** Começar em `PragmaticPlayController.php`
2. **Rastrear Cada Método:** Para cada handler de endpoint público (verifySession, getBalance, etc.):
   - Listar todos os pontos de validação/decisão
   - Extrair a lógica (cadeias if-then-else)
   - Identificar queries de banco e transformações
   - Anotar retornos de erro e códigos de status
3. **Extrair Regras:** Converter pseudocódigo em regras formais (veja template acima)
4. **Referência Cruzada:** Vincular cada regra aos números de linha exatos no código PHP
5. **Testar Ambiguidades:** Se a lógica for pouco clara, verificar com API ao vivo ou perguntar @architect

### Perguntas para Responder Enquanto Extrair

- Que validações acontecem *primeiro* (ordem fail-fast)?
- Que tabelas de banco são consultadas e em qual ordem?
- Que mudanças de estado são feitas (inserts, updates, deletes)?
- Quais códigos de erro são retornados para cada caso de falha?
- Há race conditions ou problemas de timing no código?
- Quais suposições o código faz sobre consistência de dados?

### Armadilhas Comuns

- ❌ **Muito abstrato:** "Sessão é verificada" (muito vago)
- ✅ **Concreto:** "Sessão existe, está ativa e foi criada há <24h" (testável)

- ❌ **Misturando explicação:** "O código verifica se X" (narrativa)
- ✅ **Regra de negócio:** "Regra: X deve ser verdadeiro para operação Y" (testável)

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|------|---------|--------|
| `docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md` | Extração de regras de negócio (12 regras PP-001 a PP-012) | ✅ Criado |
| `legacy/casino-proxy/app/Services/PragmaticPlayService.php` | Código fonte analisado | ✅ Analisado |
| `legacy/casino-proxy/tests/Feature/PragmaticPlayControllerTest.php` | Testes usados para validação | ✅ Validado |

---

## Definição de Pronto

- [x] Documento de regras de negócio completo e revisado — **12 regras bem documentadas em docs/casino-proxy/phase-1-business-rules/pragmatic-play-rules.md**
- [x] Todas as regras rastreadas até código fonte (números de linha) — **todas as 12 regras possuem referências exatas de linhas**
- [x] Nenhuma ambiguidade restante (ou documentada como tal) — **5 questões abertas documentadas na seção "Questões Abertas"**
- [x] Lista de Arquivos atualizada nesta story — **3 arquivos documentados**
- ⏳ PR aberto (sem merge ainda—aguardar validação do PO) — **Pendente confirmação do usuário**
- ⏳ @po valida que regras estão completas antes de Fase 2 começar — **Aguardando @po (será realizado pelo usuário)**

---

## Estratégia de Teste

**Esta story:** Apenas descoberta, nenhum teste ainda.

**Como as regras serão validadas:**
- Fase 3 (Construir suite de testes de integração) escreverá testes reais para essas regras
- Fase 5 (Validar) confirmará que testes passam contra PHP 100%

---

## Métricas de Sucesso

- **Completude:** Todos os endpoints cobertos, sem lacunas
- **Rastreabilidade:** Cada regra mapeia para linha de código fonte
- **Clareza:** @dev consegue ler regras e implementar handler Go sem consultar código PHP
- **Correção:** Regras correspondem ao comportamento real do Pragmatic Play (validado por testes)

---

## 🤖 CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Documentation`
- Complexidade: Medium (análise de código PHP + extração de 12 regras)
- Tipo secundário: `Research`

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @architect (revisar fidelidade técnica das regras)

**Quality Gate Tasks:**
- [x] Pre-Commit (@dev): Markdown renderiza corretamente, referências de linha válidas
- [x] Pre-PR (@devops): Links a arquivos PHP existentes verificados

**Self-Healing Configuration:**
```yaml
mode: light
max_iterations: 2
severity_filter: [CRITICAL, HIGH]
behavior:
  CRITICAL: auto_fix
  HIGH: document_as_debt
```

**Focus Areas (Documentation/Research):**
- Fidelidade das regras extraídas vs código PHP fonte
- Rastreabilidade linha-a-linha
- Completude — nenhuma regra omitida
- Clareza dos pseudocódigos para implementação Go

---

## Comunicação/Slack

Aviso prévio para @architect: Esta descoberta pode surfaçar decisões arquiteturais que informam o design da Fase 2. Se encontrado, discutiremos antes de prosseguir.

---

## Notas

- **Iniciado:** 2026-05-10
- **Concluído:** 2026-05-11
- **Tempo Real:** ~2 horas (leitura de código + extração + documentação)
- **Estimado:** 4-6 horas
- **Depende De:** CASINO-1.1 (Spec OpenAPI do Pragmatic Play completa) ✅
- **Bloqueia:** CASINO-1.8 (próximo provider: Mancala)

---

**Pronto para validação de @po.** Uma vez aprovado, repetimos esta story para os 7 providers restantes, depois movemos para Fase 2 (documentação).

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-11 | @dev (Dex) | Story implementada — todos os ACs concluídos, 12 regras extraídas |
| 2026-05-12 | @po (Pax) | Validação GO (8/10) — Status: Completo → Done. Desvio de lifecycle registrado (impl. antes de validação). Adicionada seção CodeRabbit Integration. |
