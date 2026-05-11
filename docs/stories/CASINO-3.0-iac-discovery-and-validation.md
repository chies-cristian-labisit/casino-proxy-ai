# CASINO-3.0: Infrastructure as Code Discovery & Validation

**Story ID:** CASINO-3.0 (Fase 0 de CASINO-3)  
**Epic:** CASINO-3 (Casino Proxy Go Microservices Migration)  
**Tipo:** Infrastructure Discovery & Setup  
**Status:** Planejamento  
**Prioridade:** CRÍTICA (bloqueia CASINO-3.1+)  
**Atribuído a:** @architect + @devops  
**Relacionado:** CASINO-1 (OpenAPI completo), CASINO-2 (em paralelo)  

---

## Resumo da Story

Descobrir, avaliar e estabelecer a **Infrastructure as Code (IaC)** que suportará toda a migração Casino Proxy para Go. Antes de implementar serviços, arquitetura ou testes, precisamos saber **como vamos provisionar, configurar e orquestrar a infraestrutura**.

**Objetivo:** Validar escolha de IaC com deploy funcional + health check endpoints respondendo.

---

## Contexto

### Por que Antes de Tudo?

A IaC é a **fundação** para:
- Provisionar ambientes (dev, staging, prod)
- Configurar banco de dados PostgreSQL novo
- Deploy de serviços Go (containers, orchestration)
- Dual-write durante migração híbrida
- CI/CD pipelines
- Monitoring & alerting

**Sem IaC decidida agora:**
- Fases 2-5 ficam presas esperando infraestrutura
- Risco de decisões arquiteturais inconsistentes
- Migração fica manual e frágil

---

## Critérios de Aceitação

### Deve Ter (Definition of Done)

- [ ] **AC-1:** Avaliar 3+ opções de IaC (Terraform, CloudFormation, Pulumi, etc.)
- [ ] **AC-2:** Documentar matriz de comparação (features, curva aprendizado, suporte, custo)
- [ ] **AC-3:** Escolher IaC baseado em análise + decisão arquitetural (@architect)
- [ ] **AC-4:** Setup inicial de IaC no projeto (`/infrastructure/` ou pasta apropriada)
- [ ] **AC-5:** Criar configuração básica:
  - VPC/networking (se cloud)
  - RDS PostgreSQL (ou managed database)
  - ECR/container registry (para images Go)
  - Security groups/network policies
  - IAM roles (se AWS/GCP/Azure)
- [ ] **AC-6:** Deploy de stack IaC para ambiente **staging/dev**
- [ ] **AC-7:** Criar health check endpoints de teste
  - GET `/health` → 200 OK com timestamp
  - GET `/health/db` → 200 OK se DB conectado
  - GET `/health/infra` → 200 OK se infraestrutura ready
- [ ] **AC-8:** Validar deployment:
  - Health check endpoints respondendo
  - Database acessível
  - Container registry configurado
  - Logging centralizado
- [ ] **AC-9:** Documentar arquitetura IaC
  - Diagrama da infraestrutura
  - Decisões (por que Terraform vs CloudFormation, etc.)
  - Guia de provisioning
  - Guia de mudanças (how to update infra)
- [ ] **AC-10:** Documentar runbooks básicos
  - Como provisionar novo ambiente
  - Como fazer deploy de serviço
  - Como escalar recursos
  - Como fazer rollback

### Deveria Ter

- [ ] **AC-11:** Configurar CI/CD pipeline básico para aplicar IaC
  - GitHub Actions (ou similar) para `terraform apply`
  - Validações: `terraform validate`, `terraform plan`
  - Approval gates antes de apply em prod
  
- [ ] **AC-12:** Setup de monitoring básico (Prometheus, CloudWatch, etc.)
  - Métricas de infraestrutura
  - Logs centralizados
  - Alertas iniciais

### Fora do Escopo

- ❌ Implementar serviços Go (Fase 3)
- ❌ Esquema de banco de dados completo (Fase 2)
- ❌ CI/CD completo para aplicação (Fase 3)
- ❌ Production deployment (Fase 5)

---

## Detalhes Técnicos

### Opções de IaC para Avaliar

| Opção | Pros | Cons | Curva Aprendizado | Custo |
|-------|------|------|-------------------|-------|
| **Terraform** | Cloud-agnostic, popular, HCL legível | State management complexo | Média | Grátis |
| **CloudFormation** | Native AWS, integrado console | AWS-only, JSON/YAML verboso | Média-Alta | Grátis |
| **Pulumi** | Linguagens de programação (Python/Go/JS), moderno | Menos adoção, Beta em alguns recursos | Média | Grátis |
| **CDK (AWS)** | Python/TypeScript, abstrations altas | AWS-only | Baixa | Grátis |
| **Docker Compose** | Simples para dev | Não é IaC real, não scale | Baixa | Grátis |

**Recomendação Prévia:** Terraform (cloud-agnostic, matura, comunidade grande)

### Infraestrutura Base Necessária

```yaml
Networking:
  - VPC (ou Virtual Network)
  - Subnets (público/privado)
  - Security Groups (inbound/outbound rules)
  - NAT Gateway (para egress)

Database:
  - RDS PostgreSQL (ou Cloud SQL)
  - Subnet privada (não exposta)
  - Automated backups
  - Monitoring

Container Registry:
  - ECR (AWS) / GCR (GCP) / ACR (Azure)
  - Autenticação para CI/CD
  - Image scanning (segurança)

Orchestration:
  - ECS Fargate (AWS) / GKE (GCP) / AKS (Azure)
  - Auto-scaling configuration
  - Load balancer (ALB/NLB)

Monitoring:
  - Centralized logging (CloudWatch / ELK / Datadog)
  - Metrics (Prometheus / CloudWatch)
  - Alerting rules

Storage (se necessário):
  - S3 (AWS) / GCS (GCP) para backups/logs
  - Encryption at rest

Secrets Management:
  - Secrets Manager (ou Vault)
  - Rotating credentials
  - Audit logging
```

### Estrutura de Pastas IaC

```
infrastructure/
├─ main.tf              # Configuração principal
├─ variables.tf         # Variáveis (ambiente, tamanhos, etc)
├─ outputs.tf           # Outputs (URLs, IPs, etc)
├─ provider.tf          # Provider configuration (AWS/GCP/Azure)
├─ terraform.tfvars     # Values por ambiente (dev, staging, prod)
│
├─ modules/
│  ├─ networking/       # VPC, subnets, security groups
│  ├─ database/         # RDS/Cloud SQL setup
│  ├─ container_registry/ # ECR/GCR configuration
│  ├─ orchestration/    # ECS/EKS/GKE setup
│  └─ monitoring/       # Logging, metrics, alerting
│
├─ envs/
│  ├─ dev.tfvars
│  ├─ staging.tfvars
│  └─ prod.tfvars
│
└─ docs/
   ├─ architecture.md
   ├─ provisioning-guide.md
   └─ runbooks/
      ├─ deploy-service.md
      ├─ scale-resources.md
      └─ rollback.md
```

---

## Entregáveis

### Primários

1. **IaC Code** (`/infrastructure/`)
   - Terraform modules (ou CloudFormation templates)
   - Variables e outputs documentados
   - Configuração para dev/staging

2. **Documentação Arquitetura** (`/infrastructure/docs/architecture.md`)
   - Diagrama da infraestrutura
   - Decisão: por que Terraform (ou alternativa)
   - Trade-offs avaliados
   - Assumptions e constraints

3. **Health Check Service** (Serviço Go mínimo)
   - Endpoints: `/health`, `/health/db`, `/health/infra`
   - Deployável via IaC (container)
   - Responde status da infraestrutura

4. **Runbooks** (`/infrastructure/docs/runbooks/`)
   - Como provisionar ambiente novo
   - Como fazer deploy de serviço
   - Como escalar recursos
   - Como fazer rollback

5. **Validation Report**
   - Deploy bem-sucedido em staging
   - Health checks respondendo 200 OK
   - Database acessível
   - Container registry funcional

### Secundários

- CI/CD pipeline básico (GitHub Actions para `terraform apply`)
- Monitoring setup inicial (métricas, logs)
- Security review da IaC (IAM roles, network policies)

---

## Notas de Implementação

### Abordagem

1. **Pesquisa (2-3h):** Avaliar opções IaC, criar matriz comparativa
2. **Decisão (1-2h):** @architect escolhe + documenta racional
3. **Setup (4-6h):** Terraform modules básicos
4. **Deploy (2-4h):** Provisionar environment staging, testar connectivity
5. **Documentação (2-3h):** Architecture, runbooks, health check setup
6. **Validação (1-2h):** Deploy + health checks passando
7. **Buffer (2h):** Revisão, ajustes, correções

**Total Estimado:** 14-20 horas (1.75-2.5 dias de trabalho)

### Questões para Responder

- Qual provider cloud? (AWS, GCP, Azure, ou on-prem?)
- Qual orquestração? (ECS Fargate, EKS, GKE, ou Docker Compose?)
- Qual database managed? (RDS, Cloud SQL, ou self-hosted?)
- Qual monitoring? (CloudWatch, Prometheus, Datadog?)
- Qual secrets management? (Secrets Manager, Vault, ou simples?)
- Qual budget mensal esperado para infraestrutura?
- Há infraestrutura existente que reutilizar?

### Riscos

- **IaC errada escolhida:** Pode bloquear fases futuras → Decisão rápida + validação com deploy
- **Infraestrutura incompleta:** Missing database, registry, etc → AC checklist
- **Health checks inadequados:** Não detecta falhas reais → Testar com falhas simuladas

---

## Métricas de Sucesso

- ✅ IaC escolhida e documentada
- ✅ Infraestrutura deployável com single command (`terraform apply`)
- ✅ Health check endpoints respondendo 200 OK
- ✅ Database PostgreSQL acessível
- ✅ Container registry funcional
- ✅ Toda equipe pode provisionar ambiente nova independentemente
- ✅ Runbooks claros e testados

---

## Dependências & Bloqueadores

**Precisa:**
- Decisão: Cloud provider (AWS/GCP/Azure/on-prem)?
- Acesso: Conta cloud / credenciais admin
- Conhecimento: Alguém na equipe familiar com IaC

**Bloqueia:**
- CASINO-3.1-3.4 (Architecture design - precisa saber infra options)
- CASINO-3.5-3.11 (Go implementation - precisa provisionar envs)

---

## Comunicação

Aviso prévio para @devops: Vamos escolher IaC que vocês vão manter. Input sobre ferramentas, standards corporativos appreciated.

---

## Notas

- **Estimado:** 1.75-2.5 dias (14-20 horas)
- **Depende De:** CASINO-2 Fase 1 (pode começar em paralelo, mas Fase 2+ de CASINO-3 bloqueada)
- **Bloqueia:** CASINO-3.1+
- **Tag:** infrastructure, iac, devops, critical

---

**Pronto para validação @architect + @devops.** Uma vez aprovado, começar descoberta IaC imediatamente (paralelo com CASINO-2).
