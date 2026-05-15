# CASINO-2.3-setup: Setup do Projeto Casino Proxy Test Oracle

**Story ID:** CASINO-2.3-setup  
**Epic:** CASINO-2.3 (Pragmatic Play — Fase 3 Test Oracle)  
**Tipo:** Infraestrutura de Testes — Setup de Projeto Java  
**Status:** Ready  
**Prioridade:** Alta  
**Atribuído a:** @dev (Dex)  
**Relacionado:** CASINO-2.2 (Fase 2 — Documentação Técnica ✅), CASINO-2.3 (Epic Test Oracle)  
**Data de Criação:** 2026-05-15  

---

## Resumo da Story

Criar o projeto Maven `casino-proxy-test-oracle/` com toda a infraestrutura de testes configurada: JUnit 5, WireMock 3.x, RestAssured e AssertJ. O projeto deve buildar limpo com `mvn clean test` e ter um health-check test passando, provando que o ambiente está funcional antes das stories de implementação de testes.

**Objetivo:** Estrutura de projeto Java pronta e buildando — base para todas as 4 stories subsequentes do CASINO-2.3.

---

## Contexto

### Por que esta Story?

As 4 stories seguintes (regras genéricas, regras Pragmatic, endpoints, CI/CD) dependem de uma base Java sólida. Separar o setup em story própria garante que o ambiente esteja funcionando antes de qualquer escrita de teste, evitando que problemas de configuração bloqueiem o progresso.

### Como se Encaixa no Plano

```
CASINO-2.3-setup          ← ESTA STORY (infraestrutura base)
CASINO-2.3-generic-rules  (aguarda setup)
CASINO-2.3-pragmatic-rules (aguarda generic-rules)
CASINO-2.3-endpoint-tests  (aguarda pragmatic-rules)
CASINO-2.3-ci-cd           (aguarda endpoint-tests)
```

### Princípio Core do Test Oracle

O projeto é **agnóstico de implementação** — os testes enviam HTTP requests para qualquer sistema via `base_url` configurável. Mudar de PHP para Go = alterar apenas `application.properties`. Zero mudança no código de teste.

---

## Critérios de Aceitação

### Deve Ter

- [ ] **AC-1:** Diretório `casino-proxy-test-oracle/` criado na raiz do repositório
- [ ] **AC-2:** `pom.xml` configurado com dependências: JUnit 5.x, WireMock 3.x, RestAssured 5.x, AssertJ 3.x, JDK 21+
- [ ] **AC-3:** Estrutura de pacotes criada: `src/main/java/com/casino/oracle/{client,mock,assertions,data,config}` e `src/test/java/com/casino/oracle/{rules,integration}`
- [ ] **AC-4:** `TestConfig.java` criado com `base_url` lido de `application.properties` (não hardcoded)
- [ ] **AC-5:** `application.properties` com `oracle.base_url=http://localhost:8000` como padrão configurável
- [ ] **AC-6:** `ProviderMockServer.java` criado com WireMock standalone configurável na porta `8081`
- [ ] **AC-7:** `HealthCheckTest.java` criado em `src/test/java/com/casino/oracle/` — testa que WireMock sobe e responde
- [ ] **AC-8:** `mvn clean test` executa com BUILD SUCCESS (apenas HealthCheckTest)
- [ ] **AC-9:** File List desta story atualizada

### Deveria Ter

- [ ] **AC-10:** `.gitignore` configurado para excluir `target/`, `.mvn/wrapper/`, `*.class`
- [ ] **AC-11:** `HttpClientFactory.java` com método estático para criar RestAssured `RequestSpecification` com `base_url`

### Fora do Escopo

- ❌ Escrever testes de regras BR-* (story CASINO-2.3-generic-rules)
- ❌ Escrever testes de endpoints (story CASINO-2.3-endpoint-tests)
- ❌ Configurar CI/CD GitHub Actions (story CASINO-2.3-ci-cd)
- ❌ Criar stubs WireMock JSON (stories seguintes)
- ❌ Implementar lógica de teste — apenas infraestrutura

---

## Detalhes Técnicos / Dev Notes

### Stack Definida

```xml
<!-- pom.xml — dependências principais -->
<dependencies>
  <dependency>
    <groupId>org.junit.jupiter</groupId>
    <artifactId>junit-jupiter</artifactId>
    <version>5.10.x</version>
    <scope>test</scope>
  </dependency>
  <dependency>
    <groupId>com.github.tomakehurst</groupId>
    <artifactId>wiremock-standalone</artifactId>
    <version>3.x</version>
    <scope>test</scope>
  </dependency>
  <dependency>
    <groupId>io.rest-assured</groupId>
    <artifactId>rest-assured</artifactId>
    <version>5.x</version>
    <scope>test</scope>
  </dependency>
  <dependency>
    <groupId>org.assertj</groupId>
    <artifactId>assertj-core</artifactId>
    <version>3.x</version>
    <scope>test</scope>
  </dependency>
</dependencies>
```

### Estrutura de Diretórios Completa

```
casino-proxy-test-oracle/
├── pom.xml
├── .gitignore
├── README.md                                    (placeholder — detalhado na story ci-cd)
└── src/
    ├── main/
    │   ├── java/com/casino/oracle/
    │   │   ├── client/
    │   │   │   ├── HttpClientFactory.java
    │   │   │   └── PayloadBuilder.java          (placeholder vazio)
    │   │   ├── mock/
    │   │   │   ├── ProviderMockServer.java
    │   │   │   └── PragmaticPlayMocks.java       (placeholder vazio)
    │   │   ├── assertions/
    │   │   │   ├── ResponseAssertions.java       (placeholder vazio)
    │   │   │   ├── RuleAssertions.java           (placeholder vazio)
    │   │   │   └── SecurityAssertions.java       (placeholder vazio)
    │   │   ├── data/
    │   │   │   ├── Fixtures.java                 (placeholder vazio)
    │   │   │   └── PragmaticPlayFixtures.java    (placeholder vazio)
    │   │   └── config/
    │   │       └── TestConfig.java
    │   └── resources/
    │       └── application.properties
    └── test/
        └── java/com/casino/oracle/
            ├── HealthCheckTest.java
            ├── rules/                            (diretório vazio)
            └── integration/                      (diretório vazio)
```

### TestConfig.java — Comportamento Esperado

```java
// Lê base_url de application.properties
// Fallback: http://localhost:8000
// Nunca hardcode URLs no código de teste
public class TestConfig {
    public static String getBaseUrl() {
        // lê oracle.base_url de properties
    }
}
```

### HealthCheckTest.java — Comportamento Esperado

```java
// Verifica que WireMock sobe corretamente
// Registra stub GET /health → 200 OK
// Faz request e confirma 200
// PASS = infraestrutura funcionando
@Test
void wireMockStartsAndResponds() {
    // stub + request + assert 200
}
```

---

## Tasks / Subtasks

- [ ] **T-1:** Criar diretório `casino-proxy-test-oracle/` na raiz do repositório
- [ ] **T-2:** Criar `pom.xml` com groupId `com.casino`, artifactId `casino-proxy-test-oracle`, JDK 21, dependências completas
- [ ] **T-3:** Criar estrutura de pacotes Java completa (todos os diretórios conforme estrutura acima)
- [ ] **T-4:** Criar `src/main/resources/application.properties` com `oracle.base_url=http://localhost:8000`
- [ ] **T-5:** Criar `TestConfig.java` — lê `oracle.base_url` de properties
- [ ] **T-6:** Criar `ProviderMockServer.java` — WireMock server configurável (porta 8081)
- [ ] **T-7:** Criar `HttpClientFactory.java` — factory com `RequestSpecification` baseado em `TestConfig.getBaseUrl()`
- [ ] **T-8:** Criar placeholders vazios para classes restantes (PayloadBuilder, PragmaticPlayMocks, assertions, fixtures)
- [ ] **T-9:** Criar `HealthCheckTest.java` — stub WireMock GET /health → 200, assert passa
- [ ] **T-10:** Criar `.gitignore` para o projeto Java
- [ ] **T-11:** Criar `README.md` placeholder (1 linha: "Casino Proxy Test Oracle — ver CASINO-2.3-ci-cd para documentação completa")
- [ ] **T-12:** Executar `mvn clean test` — confirmar BUILD SUCCESS
- [ ] **T-13:** Atualizar File List desta story

---

## CodeRabbit Integration

**Story Type Analysis:**
- Tipo primário: `Infrastructure`
- Complexidade: Low (boilerplate Maven + config)
- Tipo secundário: N/A

**Specialized Agents:**
- Executor primário: @dev
- Quality Gate: @qa

**Quality Gate Tasks:**
- [ ] Pre-Commit (@dev): `mvn clean test` com BUILD SUCCESS
- [ ] Pre-PR (@devops): Validar que `.gitignore` exclui artefatos de build

**Self-Healing Configuration:**
```yaml
mode: light
max_iterations: 2
severity_filter: [CRITICAL, HIGH]
behavior:
  CRITICAL: auto_fix
  HIGH: document_as_debt
```

**Focus Areas (Infrastructure):**
- `pom.xml` versões de dependências válidas e compatíveis
- Nenhuma URL hardcoded (tudo via `TestConfig`)
- Estrutura de pacotes segue convenção `com.casino.oracle`

---

## Entregáveis

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `casino-proxy-test-oracle/pom.xml` | Build config Maven | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/config/TestConfig.java` | Configuração centralizada | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/mock/ProviderMockServer.java` | WireMock base | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/client/HttpClientFactory.java` | HTTP client factory | ⏳ A Criar |
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/HealthCheckTest.java` | Smoke test | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/resources/application.properties` | Config | ⏳ A Criar |

---

## Lista de Arquivos

| Arquivo | Propósito | Status |
|---------|-----------|--------|
| `casino-proxy-test-oracle/pom.xml` | Maven build | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/resources/application.properties` | Base URL config | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/config/TestConfig.java` | URL config reader | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/mock/ProviderMockServer.java` | WireMock server | ⏳ A Criar |
| `casino-proxy-test-oracle/src/main/java/com/casino/oracle/client/HttpClientFactory.java` | RestAssured factory | ⏳ A Criar |
| `casino-proxy-test-oracle/src/test/java/com/casino/oracle/HealthCheckTest.java` | Health smoke test | ⏳ A Criar |
| `docs/epics/casino-proxy/CASINO-2.3-pragmatic-play-test-oracle.md` | Epic de referência | ✅ Existe |

---

## Definição de Pronto

- [ ] Diretório `casino-proxy-test-oracle/` criado
- [ ] `pom.xml` com dependências corretas (JUnit 5, WireMock 3, RestAssured 5, AssertJ 3, JDK 21)
- [ ] `TestConfig.java` lê `base_url` de properties (sem hardcode)
- [ ] `ProviderMockServer.java` sobe WireMock corretamente
- [ ] `HealthCheckTest.java` passa com `mvn clean test`
- [ ] Estrutura de pacotes completa criada
- [ ] `.gitignore` configurado
- [ ] File List atualizada
- [ ] Pronto para validação @po antes de CASINO-2.3-generic-rules iniciar

---

## Estratégia de Teste

**Validação desta story:** `mvn clean test` — BUILD SUCCESS com `HealthCheckTest` passando.  
**Validação de @po:** Confirmar que estrutura de projeto está correta e agnóstica de implementação (`base_url` configurável).  
**Próxima Story:** CASINO-2.3-generic-rules (implementa 7 testes de regras BR-GENERIC-*).

---

## Métricas de Sucesso

- **Build:** `mvn clean test` retorna BUILD SUCCESS
- **Isolamento:** Nenhuma URL hardcoded no código Java
- **Extensibilidade:** Estrutura de pacotes permite adicionar novos providers sem reorganizar
- **Documentação:** README placeholder presente (será completado em CASINO-2.3-ci-cd)

---

## Notas

- **Criado:** 2026-05-15
- **Estimado:** 3-4 horas
- **Depende De:** CASINO-2.2 (documentação dos endpoints ✅)
- **Bloqueia:** CASINO-2.3-generic-rules (próxima story)

---

## Change Log

| Data | Agente | Ação |
|------|--------|------|
| 2026-05-15 | @sm (River) | Story criada — Draft |
| 2026-05-15 | @po (Pax) | Validação GO (8/10) — Status: Draft → Ready. Should-fix: adicionar seção Riscos em revisão futura. |
