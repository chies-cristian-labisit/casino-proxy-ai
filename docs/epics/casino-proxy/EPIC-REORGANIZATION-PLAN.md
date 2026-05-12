# Plano de Reorganização de Epics - Casino Proxy Migration

**Data:** 2026-05-11  
**Propósito:** Dividir o grande epic CASINO-1 em 3 epics focados, cada um com ciclo de vida e validação independentes.

---

## Estrutura Atual (ANTES)

```
CASINO-1: Casino Proxy PHP → Go Migration
├─ Fase 1: OpenAPI Discovery (CASINO-1.1-1.6) ✅ PRONTO
├─ Fase 1.5: Business Rules Discovery (CASINO-1.7-1.14) ← NOVO
├─ Fase 2: Microservices Architecture (CASINO-2.1-2.4)
├─ Fase 3: Core Services Implementation (CASINO-3.1-3.7)
├─ Fase 4: Database & Data Migration (CASINO-4.1-4.5)
└─ Fase 5: Testing & Hybrid Deployment (CASINO-5.1-5.6)

TOTAL: 42 stories em um mega-epic
```

---

## Estrutura Nova (DEPOIS)

```
CASINO-1: OpenAPI Documentation Generation ✅ READY
├─ 1.1: Analyze & Document Pragmatic Play ✅
├─ 1.2: Analyze & Document Evolution Gaming ✅
├─ 1.3: Analyze & Document PG Soft ✅
├─ 1.4: Analyze & Document Remaining Providers ✅
├─ 1.5: Create Master OpenAPI Spec ✅
└─ 1.6: Document Admin API ✅
└─ Status: ✅ COMPLETO (pode fazer merge quando quiser)

CASINO-2: Business Rules Discovery & Test Oracle ← VOCÊ ESTÁ AQUI
├─ 2.1-2.5: Pragmatic Play (Extract → Doc → Test → Matrix → Validate)
├─ 2.6-2.10: Evolution Gaming (5 fases)
├─ 2.11-2.15: PG Soft (5 fases)
├─ 2.16-2.20: Mancala, Digitain, Evoplay (5 fases)
└─ 2.21-2.25: OpenBox, Alternar (5 fases)
└─ Status: 🚀 INICIANDO
└─ Bloqueia: CASINO-3 (não pode implementar Go até ter oráculo de testes)

CASINO-3: Go Microservices Implementation & Migration
├─ Fase 2: Microservices Architecture (CASINO-3.1-3.4)
├─ Fase 3: Core Services Implementation (CASINO-3.5-3.11)
├─ Fase 4: Database & Data Migration (CASINO-3.12-3.16)
└─ Fase 5: Testing & Hybrid Deployment (CASINO-3.17-3.22)
└─ Status: ⏸️ AGUARDANDO (bloqueado por CASINO-2)
```

---

## Timeline & Dependências

```
Semana 1           Semana 2-3         Semana 4+           Semana 8+
│                  │                  │                   │
└─ CASINO-2.1 ──┬─ CASINO-2.2/3/4/5  │                   │
   (Extract)     │  (por provider)    │                   │
                 │  (validar PO)      │                   │
                 │                    │                   │
                 └──────────┬─────────┘                   │
                            │                             │
                     CASINO-2 ✅ 100%                     │
                     (Oráculo de Testes Pronto)          │
                            │                             │
                            ├─────────────┬───────────────┴─ CASINO-3 Fases 2-5
                            │             │   (Implementação Go)
                     CASINO-3.1-3.4      │   (Pode começar paralelo com CASINO-2 Fase 2+)
                     (Arquitetura)       │
                                         └─ CASINO-3.5+ (Implementação)
```

---

## Responsabilidades por Agent

| Epic | Foco | Agente Principal | Agente de Revisão | Status |
|------|------|------------------|------------------|--------|
| **CASINO-1** | Documentação OpenAPI | @dev | @po | ✅ COMPLETO |
| **CASINO-2** | Descoberta de Regras & Oráculo de Testes | @dev | @po (validar entre fases) | 🚀 INICIANDO |
| **CASINO-3** | Implementação & Migração Go | @dev + @architect | @qa + @devops | ⏸️ AGUARDANDO CASINO-2 |

---

## Detalhes de Cada Epic

### CASINO-1: OpenAPI Documentation Generation

**Objetivo:** Documentar 100% dos endpoints em OpenAPI 3.0

**Escopo:**
- ✅ 6 stories completas (Pragmatic Play, Evolution, PG Soft, Remaining, Master Spec, Admin API)
- Toda integração com 8 providers documentada
- Pronto para referência de implementação Go

**Resultado:** `/docs/casino-proxy/openapi/` com specs de cada provider

**Timeline:** ✅ PRONTO AGORA

**Definition of Done:**
- [ ] Mesclar branch feature/openapi-documentation-viewers para master
- [ ] Validar com PO que documentação é 100% completa
- [ ] Tag release (v1.0-openapi)

---

### CASINO-2: Business Rules Discovery & Test Oracle

**Objetivo:** Descobrir & documentar regras de negócio PHP, construir oráculo agnóstico de testes

**Escopo:**
- **Fase 1:** Extrair lógica de negócio de cada handler PHP (CASINO-2.1-2.5 / 2.6-2.10 / etc.)
- **Fase 2:** Documentar em markdown (rules.md + endpoint docs por provider)
- **Fase 3:** Construir Test Oracle Java (agnóstico a linguagem, com WireMock para mocks)
- **Fase 4:** Criar YAML trace matrix (regra → spec → código → teste)
- **Fase 5:** Validar: testes passam 100% contra PHP legado

**Resultado:** 
- `/docs/casino-proxy/phase-1-business-rules/` (regras markdown com BR-* nomenclature)
- `/docs/casino-proxy/phase-2-technical-documentation/` (endpoint documentation por provider)
- `casino-proxy-test-oracle/` (módulo Java independente com JUnit 5 + WireMock)
  - Reutilizável para testar PHP legado E Go futuro
  - 50+ testes por provider (1 por regra BR-* + por endpoint)
- `/docs/casino-proxy/trace-matrices/` (YAML trace matrices)

**Definition of Done:**
- [ ] 8 providers com todos os handlers mapeados
- [ ] Regras de negócio documentadas + validadas (BR-* nomenclature)
- [ ] Endpoint documentation (8 fases cada) como templates
- [ ] Test Oracle implementado em Java
- [ ] 50+ testes por provider rodando contra PHP legado 100% PASS
- [ ] PO valida completude antes de liberar CASINO-3

**Bloqueia:** CASINO-3 (não pode implementar Go sem oráculo)

---

### CASINO-3: Go Microservices Implementation & Migration

**Objetivo:** Implementar serviços Go que replicam 100% comportamento PHP (com descoberta prévia de IaC)

**Escopo:**
- **Fase 0 (CASINO-3.0):** Descoberta & Validação de Infrastructure as Code
  - Avaliar opções de IaC (Terraform, CloudFormation, Pulumi, etc.)
  - Escolher stack IaC apropriado para projeto
  - Setup inicial de infrastructure code
  - Deploy e validação com health check endpoints
  - Documentar decisões arquiteturais de infraestrutura
  
- **Fase 2 (CASINO-3.1-3.4):** Arquitetura microservices + design de banco de dados
  - Usar IaC definida em Fase 0 para toda infraestrutura
  - Microservices architecture design
  - Database schema design (PostgreSQL + GORM)
  - Gateway service architecture
  
- **Fase 3 (CASINO-3.5-3.11):** Implementar serviços Go (1 por provider + gateway + admin)
  - Usar IaC para provisionar ambientes
  - Implementar cada serviço Go
  - Configurar CI/CD pipeline
  
- **Fase 4 (CASINO-3.12-3.16):** Migração de dados + dual-write
  - Usar IaC para provisionar nova infraestrutura PostgreSQL
  - Implementar migração de dados
  - Setup dual-write entre PHP e Go
  
- **Fase 5 (CASINO-3.17-3.22):** Testes, deployment híbrido, migração tráfego
  - Deploy híbrido (PHP + Go) via IaC
  - Testes end-to-end
  - Migração gradual de tráfego
  - Cutover final

**Sequência:** Pode começar Fase 0 (IaC) em paralelo com CASINO-2. Fases 2-5 iniciam quando CASINO-2 Fase 1 completa.

**Resultado:**
- IaC completa (`/infrastructure/` com Terraform/CloudFormation)
- Serviços Go em `/services/` (1 pasta por provider + gateway)
- Schema PostgreSQL novo (migrations GORM)
- Suite de testes Go (match CASINO-2 tests)
- CI/CD pipelines configured
- Documentação operacional

**Validação:** 
- **Fase 0:** Deploy de infraestrutura com health check endpoints respondendo
- **Fases 2-5:** Go tests passam contra CASINO-2 PHP tests (prova de parity)

**Definition of Done:**
- [ ] Fase 0: IaC escolhida, validada e documentada
- [ ] Todos 8 providers migrados para Go em produção
- [ ] Infraestrutura gerenciada 100% via IaC
- [ ] Zero downtime durante migração por provider
- [ ] Performance meets or exceeds PHP baseline
- [ ] PHP decommissioned
- [ ] Monitoring & alerting configured via IaC

---

## Resumo Executivo

| Métrica | CASINO-1 | CASINO-2 | CASINO-3 |
|---------|----------|----------|----------|
| **Stories** | 6 | 40 | 22 |
| **Status** | COMPLETO | INICIANDO | AGUARDANDO |
| **Depende De** | - | CASINO-1 | CASINO-2 |
| **Bloqueia** | - | CASINO-3 | - |
| **Agente** | @dev | @dev | @dev + @architect |
| **Validação** | @po | @po (por fase) | @qa + @devops |

---

## Vantagens da Reorganização

✅ **Clareza Narrativa:** Cada epic responde uma pergunta diferente
- CASINO-1: "Como documentamos os endpoints?"
- CASINO-2: "O que cada endpoint realmente faz?"
- CASINO-3: "Como reconstruímos em Go?"

✅ **Ciclos Independentes:** Cada epic tem seu próprio DoD e validação

✅ **Execução Paralela:** CASINO-2 Fase 2+ pode rodar enquanto CASINO-3 Fase 2 começa

✅ **Qualidade Garantida:** CASINO-2 tests = acceptance criteria para CASINO-3

✅ **Manageability:** Nenhum epic > 40 stories (fácil de rastrear)

✅ **Escalabilidade:** Se um provider pegar fogo, isolado no CASINO-2, não impacta arquitetura (CASINO-3)

---

## Próximas Ações

- [ ] Renomear/revalidar CASINO-1 (apenas OpenAPI, 6 stories)
- [ ] Criar CASINO-2 epic (Business Rules Discovery, 40 stories)
- [ ] Criar CASINO-3 epic (Go Implementation, 22 stories)
- [ ] Atualizar `CASINO-1-migration-plan.md` com nova estrutura
- [ ] Vincular epics com "Depends On" / "Blocks" relationships
- [ ] Começar CASINO-2.1 (Pragmatic Play extraction)

---

**Proposto por:** River (Scrum Master)  
**Data:** 2026-05-11  
**Status:** Pronto para PO validar & aprovar reorganização
