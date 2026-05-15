# CASINO-2.3-ci-cd: CI/CD Pipeline e Documentação README do Test Oracle

**Story ID:** CASINO-2.3-ci-cd  
**Epic:** CASINO-2.3 (Pragmatic Play — Fase 3 Test Oracle)  
**Tipo:** DevOps / Documentação — Pipeline e README  
**Status:** Ready  
**Prioridade:** Média  
**Atribuído a:** @dev (Dex)  
**Relacionado:** CASINO-2.3-endpoint-tests (pré-requisito — todos os testes implementados)  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Criar o pipeline CI/CD GitHub Actions para o Test Oracle e escrever o README completo do projeto `casino-proxy-test-oracle/`. Esta story finaliza o epic CASINO-2.3 — após ela, o Test Oracle está pronto para rodar automaticamente e qualquer desenvolvedor sabe como usá-lo.

**Objetivo:** Test Oracle rodando em CI/CD e documentado para extensão por novos providers e endpoints.

---

## Contexto

### Por que esta Story?

Um Test Oracle que só roda localmente é frágil. O CI/CD garante que os testes rodem a cada commit, detectando regressões automaticamente. O README transforma o projeto em algo que **qualquer desenvolvedor pode usar e estender** sem precisar ler o código.

### Como se Encaixa no Plano

```
CASINO-2.3-setup           ✅ (pré-requisito)
CASINO-2.3-generic-rules   ✅ (pré-requisito)
CASINO-2.3-pragmatic-rules ✅ (pré-requisito)
CASINO-2.3-endpoint-tests  ✅ (pré-requisito)
CASINO-2.3-ci-cd           ← ESTA STORY (fecha o epic)
```

### O que esta Story Desbloqueia

Ao completar esta story:
- Epic CASINO-2.3 está **Done** ✅
- CASINO-2.4 (Trace Matrix) pode iniciar
- O Test Oracle está pronto para receber Evolution Gaming (CASINO-2.8) no futuro

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** `.github/workflows/test-oracle.yml` criado — executa `mvn clean test` no diretório `casino-proxy-test-oracle/` a cada push em qualquer branch
- [ ] **AC-2:** Workflow dispara em `push` (todas as branches) e `pull_request` (para master)
- [ ] **AC-3:** Workflow usa Java 21 (`actions/setup-java@v4` com `distribution: 'temurin'`)
- [ ] **AC-4:** Workflow faz cache de dependências Maven (`~/.m2`) para builds mais rápidos
- [ ] **AC-5:** `casino-proxy-test-oracle/README.md` completo com as seções listadas em Dev Notes
- [ ] **AC-6:** README documenta como alterar `base_url` para rodar contra PHP (HML) ou Go (futura URL)
- [ ] **AC-7:** README inclui seção "Como adicionar novo provider" com passos concretos
- [ ] **AC-8:** README inclui seção "Como adicionar novo endpoint" com passos concretos
- [ ] **AC-9:** Badge de CI adicionado no README (`![Test Oracle](https://github.com/.../badge.svg)`)
- [ ] **AC-10:** `mvn clean test` no CI passa com BUILD SUCCESS (≥ 50 testes)
- [ ] **AC-11:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-12:** Workflow publica relatório de testes como artefato (`actions/upload-artifact`) — arquivo Surefire XML
- [ ] **AC-13:** README inclui seção de troubleshooting com os 3 erros mais comuns e como resolver

### Fora do Escopo

- ❌ Deploy automático para qualquer ambiente
- ❌ Configuração de secrets ou tokens no CI
- ❌ Testes de performance ou carga no pipeline
- ❌ Integração com SonarQube ou ferramentas de análise estática

---

## Detalhes Técnicos / Dev Notes

### test-oracle.yml — Estrutura Esperada

```yaml
name: Casino Proxy Test Oracle

on:
  push:
    branches: ["**"]
    paths:
      - 'casino-proxy-test-oracle/**'
  pull_request:
    branches: [master]
    paths:
      - 'casino-proxy-test-oracle/**'

jobs:
  test:
    name: Run Test Oracle
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Java 21
        uses: actions/setup-java@v4
        with:
          java-version: '21'
          distribution: 'temurin'
          cache: 'maven'

      - name: Run Tests
        working-directory: casino-proxy-test-oracle
        run: mvn clean test

      - name: Upload Test Reports
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: surefire-reports
          path: casino-proxy-test-oracle/target/surefire-reports/
```

> **Nota:** O workflow só dispara quando arquivos dentro de `casino-proxy-test-oracle/**` são alterados, evitando runs desnecessários.

### README.md — Estrutura Completa

```markdown
# Casino Proxy Test Oracle

Framework de testes Java agnóstico de implementação para validar o Casino Proxy.
Os mesmos testes rodam contra PHP legado e Go (futura implementação).

## Quick Start

\`\`\`bash
cd casino-proxy-test-oracle
mvn clean test
\`\`\`

## Configuração

Edite `src/main/resources/application.properties`:
\`\`\`properties
# PHP legado (HML/DEV)
oracle.base_url=http://localhost:8000

# Go (futura implementação)
# oracle.base_url=http://localhost:9000
\`\`\`

## Resultados

- ≥ 50 test cases cobrindo Pragmatic Play
- 12 regras BR-* validadas
- 9 endpoints com happy path + error path

## Arquitetura

[diagrama ou descrição da arquitetura agnóstica]

## Como Adicionar Novo Provider

1. Criar `src/main/java/com/casino/oracle/mock/{Provider}Mocks.java`
2. Criar `src/main/java/com/casino/oracle/data/{Provider}Fixtures.java`
3. Criar `src/main/resources/wiremock/{provider}/` com stubs JSON
4. Criar `src/test/java/com/casino/oracle/rules/` com classes por regra
5. Criar `src/test/java/com/casino/oracle/integration/` com classes por endpoint

## Como Adicionar Novo Endpoint

1. Adicionar payload em `PayloadBuilder.java`
2. Criar stubs JSON em `wiremock/{provider}/`
3. Adicionar stub em `{Provider}Mocks.java`
4. Criar `{Endpoint}EndpointTest.java` em `integration/`

## Troubleshooting

| Erro | Causa | Solução |
|------|-------|---------|
| Connection refused | Sistema sob teste não está rodando | Iniciar servidor na porta configurada |
| WireMock port conflict | Porta 8081 em uso | Alterar porta em `ProviderMockServer.java` |
| BUILD FAILURE (compile) | JDK < 21 | `java -version` — instalar JDK 21+ |

## CI/CD

Badge: [link]

Pipeline: `.github/workflows/test-oracle.yml`
- Dispara em: push (todos os branches), PR para master
- Artefato: relatório Surefire em `target/surefire-reports/`
```

### Localização do Workflow

```
.github/
└── workflows/
    └── test-oracle.yml    ← arquivo a criar
```

---

## Tasks / Subtasks

- [ ] **T-1:** Verificar se diretório `.github/workflows/` existe na raiz do repositório; criar se necessário
- [ ] **T-2:** Criar `.github/workflows/test-oracle.yml` conforme estrutura em Dev Notes
- [ ] **T-3:** Configurar `paths` filter no workflow para `casino-proxy-test-oracle/**`
- [ ] **T-4:** Escrever `casino-proxy-test-oracle/README.md` completo (Quick Start, Configuração, Arquitetura, Como adicionar provider/endpoint, Troubleshooting)
- [ ] **T-5:** Adicionar badge de CI no topo do README (substituir placeholder com URL real após primeiro run)
- [ ] **T-6:** Fazer push de um commit de teste para verificar que o workflow dispara e passa (coordenar com @devops para push)
- [ ] **T-7:** Confirmar que artefato Surefire é publicado na aba Actions do GitHub
- [ ] **T-8:** Atualizar File List desta story

---

## CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `DevOps` + `Documentation`
- Complexidade: Low (YAML de CI + Markdown)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Push/CI validation: @devops (coordenação para trigger do workflow)
- Quality Gate: @qa

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): Validar YAML do workflow com `yamllint` ou GitHub Actions validator
- [ ] Pre-PR (@devops): Confirmar que workflow aparece na aba Actions e passa

**Self-Healing Configuration:**
```yaml
mode: light
max_iterations: 1
severity_filter: [CRITICAL]
behavior:
  CRITICAL: auto_fix
```

**Focus Areas (DevOps + Documentation):**
- YAML do workflow: indentação correta, `working-directory` correto
- Cache Maven: reduz build time de ~3min para ~45s
- README: seções "Como adicionar provider/endpoint" devem ser passos acionáveis, não teoria

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `.github/workflows/test-oracle.yml` | Pipeline CI/CD GitHub Actions | ⏳ A Criar |
| `casino-proxy-test-oracle/README.md` | Documentação completa do projeto | ⏳ A Criar (substituir placeholder) |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `.github/workflows/test-oracle.yml` | CI/CD pipeline | ⏳ A Criar |
| `casino-proxy-test-oracle/README.md` | README completo | ⏳ A Atualizar (placeholder existe) |
| `casino-proxy-test-oracle/` (projeto completo) | Test Oracle — todos os testes | ✅ CASINO-2.3-endpoint-tests |

---

## Definição de Pronto

- [ ] `.github/workflows/test-oracle.yml` criado e válido
- [ ] Workflow dispara em push e PR para master
- [ ] README completo: Quick Start, Configuração, Como adicionar provider/endpoint, Troubleshooting
- [ ] Badge CI adicionado no README
- [ ] Workflow passou com BUILD SUCCESS em execução real no GitHub Actions
- [ ] Artefato Surefire publicado
- [ ] File List atualizada
- [ ] **Epic CASINO-2.3 Done** → CASINO-2.4 desbloqueado

---

## Estratégia de Teste

**Validação desta story:** Workflow executa e passa no GitHub Actions (coordenar com @devops para o push).  
**Validação de @po:** README é claro o suficiente para que um novo desenvolvedor consiga rodar o Test Oracle sem ajuda externa.  
**Próxima Story:** CASINO-2.4 (Trace Matrix — YAML rastreando cada BR-* por 4 camadas).

---

## Notas

- **Criado:** 2026-05-15
- **Estimado:** 2-3 horas
- **Depende De:** CASINO-2.3-endpoint-tests (todos os ≥50 testes implementados)
- **Bloqueia:** CASINO-2.4 (Trace Matrix)
- **Coordenação:** T-6 (verificar workflow no GitHub Actions) requer que @devops faça o push

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
| 2026-05-15 | @po (Pax) | Validação GO (8/10) — Status: Draft → Ready. Should-fix: adicionar seção Riscos em revisão futura. |
