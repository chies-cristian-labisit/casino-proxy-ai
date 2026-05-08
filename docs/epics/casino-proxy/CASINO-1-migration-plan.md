# CASINO-1: Casino Proxy PHP → Go Microservices Migration

**Epic ID:** CASINO-1  
**Type:** Brownfield Migration (Large Scale)  
**Status:** In Progress (Phase 1: 3/6 stories complete - 50%)  
**Created:** 2026-05-08  
**Owner:** Morgan (PM)  
**Progress:** Stories 1.1-1.3 ✅ merged | Stories 1.4-1.6 in queue  

---

## Executive Summary

Migrate **Casino Proxy** from monolithic Laravel PHP to a **microservices architecture in Go**. Each gaming provider gets a dedicated Go service, communicating with a unified Go gateway. Run PHP and Go in parallel during transition, then migrate traffic provider-by-provider to zero-downtime cutover.

**Key Metrics:**
- **Target:** 8 provider integrations
- **Architecture:** Microservices (1 service per provider + 1 gateway service)
- **Database:** New PostgreSQL with GORM
- **Deployment:** Hybrid (PHP + Go) → Full Go
- **Quality:** Engineering best practices, zero technical debt

---

## Project Analysis

### Existing System Context

**Current Technology Stack:**
- Framework: Laravel 11
- PHP Version: 8.2
- Dependencies: Sanctum (auth), Intervention Image (image processing)
- Database: Likely PostgreSQL (inferred from Laravel conventions)

**Current Functionality:**
- Gaming provider webhook gateway (8 providers)
- Session/player management
- Real-time transaction processing (bets, wins, refunds)
- Operator admin dashboard
- Image management for operators

**Integration Points:**
- **Inbound:** 8 gaming provider webhooks + admin dashboard
- **Outbound:** Likely callback APIs to gaming providers, database persistence
- **Auth:** Session-based (dashboard) + provider-specific signature validation (webhooks)

### Gaming Providers Integrated

1. **Pragmatic Play** - Generic endpoint routing
2. **Mancala** - Generic endpoint routing
3. **Digitain RGS** - Generic endpoint routing
4. **PG Soft** - Specific methods (VerifySession, GetCash, TransferInOut, Adjustment)
5. **Evoplay** - Generic webhook
6. **Evolution Gaming** - Specific methods (Auth, Debit, Credit, Rollback, Token)
7. **OpenBox** - PUT-based endpoints (Balance, Bet, Win, Refund)
8. **Alternar** - Redirect handling

### Current API Endpoints Inventory

**Public API (v1):**
```
POST   /v1/entry                                    - Session entry
GET    /v1/status                                   - Health check
POST   /v1/image                                    - Image upload

Admin (internal):
GET    /v1/internal/operator/                       - List operators
POST   /v1/internal/operator/store                  - Create operator
DELETE /v1/internal/operator/delete/{slug}         - Delete operator

Webhooks (primary traffic):
POST   /v1/webhooks/pragmatic-play/{endpoint}      - Pragmatic Play
POST   /v1/webhooks/mancala/{endpoint}             - Mancala
POST   /v1/webhooks/digitain-rgs/{endpoint}        - Digitain
POST   /v1/webhooks/pgsoft/{verb}                  - PG Soft (4 methods)
POST   /v1/webhooks/evoplay                        - Evoplay
POST   /v1/webhooks/evolution/{verb}               - Evolution (5 methods)
PUT    /v1/webhooks/openbox/seamless/{verb}        - OpenBox (4 methods)
POST   /v1/webhooks/alternar/                      - Alternar
```

**Analysis:** ❌ NO OpenAPI documentation exists. All endpoints must be reverse-engineered from controllers.

---

## Epic Goal

Create comprehensive OpenAPI documentation for all 8 provider integrations, then migrate to a microservices Go architecture with dedicated services per provider, running alongside PHP during transition for zero-downtime provider migration.

---

## Migration Phases

### **Phase 1: Endpoint Discovery & OpenAPI Generation** (PRIORITY)

**Goal:** Extract all endpoints from PHP codebase and generate OpenAPI 3.0 spec for each provider.

**Deliverables:**
- [x] Complete endpoint inventory with request/response schemas (Pragmatic Play ✅)
- [ ] OpenAPI 3.0 spec file per provider (1/8 complete - Pragmatic Play)
- [x] Provider authentication patterns documented (Pragmatic Play ✅)
- [x] Request/response examples for each endpoint (Pragmatic Play ✅)
- [x] Error handling and status codes mapped (Pragmatic Play ✅)

**Stories:**
1. **CASINO-1.1** - Analyze & Document Pragmatic Play Provider ✅ **DONE** (Merged PR #1)
2. **CASINO-1.2** - Analyze & Document Evolution Gaming Provider ✅ **DONE** (Merged PR #3)
3. **CASINO-1.3** - Analyze & Document PG Soft Provider ✅ **DONE** (Merged PR #4)
4. **CASINO-1.4** - Analyze & Document Remaining Providers (Mancala, Digitain, Evoplay, OpenBox, Alternar) (Pending)
5. **CASINO-1.5** - Create Master OpenAPI Spec (Gateway + All Services) (Pending - depends on 1.1-1.4)
6. **CASINO-1.6** - Document Admin API (Operator Management) (Pending)

---

### **Phase 2: Microservices Architecture Design**

**Goal:** Design Go microservices architecture for each provider, database schema, and gateway routing logic.

**Deliverables:**
- [ ] Architecture diagram (Provider services + Gateway)
- [ ] Service interface specifications (gRPC/HTTP)
- [ ] Data models and database schema (PostgreSQL + GORM)
- [ ] Authentication/authorization strategy (provider keys, JWT)
- [ ] Deployment topology (containers, orchestration)

**Stories:**
7. **CASINO-2.1** - Design Microservices Architecture & Service Boundaries
8. **CASINO-2.2** - Design Database Schema & Data Models (GORM)
9. **CASINO-2.3** - Design Gateway Service (routing, authentication, middleware)
10. **CASINO-2.4** - Design CI/CD Pipeline & Deployment Strategy

---

### **Phase 3: Core Services Implementation**

**Goal:** Build Go services for each provider with full webhook handling, database persistence, and transaction processing.

**Deliverables per Provider:**
- [ ] Go service with webhook handlers
- [ ] Database models (GORM entities)
- [ ] Request validation & response formatting
- [ ] Error handling & logging
- [ ] Unit & integration tests

**Stories:**
11. **CASINO-3.1** - Implement Gateway Service (core routing, auth middleware)
12. **CASINO-3.2** - Implement Pragmatic Play Service
13. **CASINO-3.3** - Implement Evolution Gaming Service
14. **CASINO-3.4** - Implement PG Soft Service
15. **CASINO-3.5** - Implement Remaining Services (Mancala, Digitain, Evoplay, OpenBox, Alternar)
16. **CASINO-3.6** - Implement Admin Service (Operator CRUD)
17. **CASINO-3.7** - Implement Session/Player Management Service

---

### **Phase 4: Database & Data Migration**

**Goal:** Create new PostgreSQL schema, migrate existing data from PHP, handle ongoing sync during hybrid mode.

**Deliverables:**
- [ ] PostgreSQL schema with GORM migrations
- [ ] Data migration scripts (PHP DB → PostgreSQL)
- [ ] Data validation & reconciliation
- [ ] Dual-write strategy for hybrid mode
- [ ] Rollback procedures

**Stories:**
18. **CASINO-4.1** - Create PostgreSQL Schema & GORM Migrations
19. **CASINO-4.2** - Implement Data Migration Scripts
20. **CASINO-4.3** - Implement Dual-Write Logic (PHP writes to both DBs)
21. **CASINO-4.4** - Data Validation & Reconciliation
22. **CASINO-4.5** - Implement Rollback Procedures

---

### **Phase 5: Testing & Hybrid Deployment**

**Goal:** Full testing, performance validation, production hybrid deployment, and gradual traffic migration.

**Deliverables:**
- [ ] Comprehensive test suite (unit, integration, end-to-end)
- [ ] Performance benchmarks (latency, throughput)
- [ ] Load testing & capacity planning
- [ ] Production deployment (PHP + Go side-by-side)
- [ ] Traffic migration strategy per provider
- [ ] Monitoring & alerting setup

**Stories:**
23. **CASINO-5.1** - Create Comprehensive Test Suite
24. **CASINO-5.2** - Performance Testing & Benchmarking
25. **CASINO-5.3** - Deploy Hybrid System (PHP + Go)
26. **CASINO-5.4** - Migrate Traffic (Provider-by-Provider)
27. **CASINO-5.5** - Monitor & Optimize Production
28. **CASINO-5.6** - Full Cutover & PHP Decommission

---

## Microservices Architecture

```
┌─────────────────────────────────────────────────────────┐
│                 Casino Proxy Gateway (Go)               │
│  ├─ Request Routing (by provider)                      │
│  ├─ Authentication (provider keys, operator tokens)    │
│  ├─ Rate Limiting & Circuit Breaking                  │
│  └─ Logging & Monitoring                               │
└────────────┬────────────────────────────────────────────┘
             │
    ┌────────┴────────────────────────────────────────┐
    │                                                 │
┌───▼────────────┐ ┌──────────────────┐ ┌───────────▼───┐
│ Pragmatic Play │ │ Evolution Gaming │ │ PG Soft Svc  │
│   Service      │ │   Service        │ │              │
│   (Go)         │ │   (Go)           │ │ (Go)         │
└────────────────┘ └──────────────────┘ └──────────────┘
        │                   │                    │
┌───────▼───────────────────▼────────────────────▼────────┐
│          PostgreSQL (New Database)                       │
│  ├─ Operators                                           │
│  ├─ Players/Sessions                                    │
│  ├─ Transactions (bets, wins, refunds)                 │
│  ├─ Provider Integrations                              │
│  └─ Audit Logs                                         │
└────────────────────────────────────────────────────────┘

Hybrid Mode (During Transition):
┌─ PHP Laravel ─────────────┐  ┌─ Go Services ──────────┐
│  (existing traffic)        │  │  (new traffic)         │
│  • Dashboard              │  │  • Webhooks (migrated) │
│  • Some webhooks          │  │  • Admin API           │
│  • Dual-write to Go DB    │  │                        │
└────────────────────────────┘  └────────────────────────┘
                    ▼                     ▼
            PostgreSQL (Synced Dual Write)
```

---

## Technology Stack

| Layer | Tech | Notes |
|-------|------|-------|
| **Language** | Go 1.22+ | High performance, concurrent |
| **Framework** | Echo or Gin | Lightweight HTTP, fast routing |
| **Database** | PostgreSQL + GORM | Type-safe, migration support |
| **Authentication** | JWT + HMAC | Provider signatures + operator tokens |
| **Logging** | Structured (JSON) | Zap or slog |
| **Monitoring** | Prometheus + Grafana | Metrics, dashboards |
| **Container** | Docker + Kubernetes | Orchestration |
| **CI/CD** | GitHub Actions | Automated testing, deployment |

---

## Success Criteria

### Phase 1 (OpenAPI Documentation)
- ✓ All 8 providers fully documented in OpenAPI 3.0
- ✓ Request/response schemas match actual PHP handlers
- ✓ Provider authentication patterns clearly specified
- ✓ Documentation validated against live PHP API

### Phase 2 (Architecture)
- ✓ Microservices architecture approved by tech lead
- ✓ GORM schema supports all provider requirements
- ✓ Gateway design handles routing, auth, logging
- ✓ Deployment strategy supports zero-downtime migration

### Phase 3 (Implementation)
- ✓ All Go services handle 100% of PHP functionality
- ✓ Database models match schema design
- ✓ Error responses match PHP API contracts
- ✓ Code coverage >80% (unit + integration tests)

### Phase 4 (Data Migration)
- ✓ Historical data migrated without loss
- ✓ Dual-write working correctly during hybrid mode
- ✓ Data validation passing 100%
- ✓ Rollback procedures tested

### Phase 5 (Production)
- ✓ All providers migrated to Go successfully
- ✓ Performance metrics meet or exceed PHP baseline
- ✓ Zero downtime during provider migration
- ✓ Monitoring alerts configured and tested
- ✓ PHP completely decommissioned

---

## Quality & Engineering Standards

### Code Quality
- **Language:** Go 1.22+
- **Testing:** Unit (>80%), Integration (key flows), E2E (provider webhooks)
- **Linting:** golangci-lint, strict rules
- **Documentation:** Godoc + OpenAPI specs
- **Reviews:** Code review + architecture validation

### Performance
- **Latency Target:** <100ms p95 webhook handling
- **Throughput:** 10k+ requests/sec per service
- **Database:** Query optimization, connection pooling
- **Memory:** <500MB per service (tuned)

### Reliability
- **Availability:** 99.9% uptime SLA
- **Error Handling:** Graceful degradation, circuit breakers
- **Monitoring:** Real-time alerts, dashboards
- **Rollback:** Instant provider failover to PHP

### Security
- **OWASP Compliance:** Input validation, SQL injection prevention
- **Provider Keys:** Encrypted at rest, rotated regularly
- **JWT Tokens:** Signed, short-lived, refresh supported
- **Audit Logging:** All transactions logged with provider context

---

## Risk Management

### Primary Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Data loss during migration | CRITICAL | Dual-write validation, reconciliation scripts, backups |
| Provider webhook validation fails | HIGH | Extensive testing, live traffic replay, gradual cutover |
| Performance degradation | HIGH | Load testing, benchmarking, capacity planning |
| Hybrid mode synchronization issues | MEDIUM | Dual-write reconciliation, monitoring alerts |
| Production cutover downtime | MEDIUM | Instant rollback, feature flags, provider-by-provider |

### Rollback Plans

**During Phase 1-4:** Roll back to PHP entirely (no traffic to Go yet)

**During Phase 5 (Hybrid):** 
- Single provider issue? → Failover that provider back to PHP
- Database issue? → Restore from backup, replay transactions from dual-write log
- Gateway issue? → Instant DNS failover to PHP load balancer

**Post-Migration:** Keep PHP in standby for 30 days, then decommission

---

## Definition of Done (Epic Level)

- [x] All endpoints documented (OpenAPI Phase 1)
- [ ] Architecture designed & approved (Phase 2)
- [ ] All Go services implemented & tested (Phase 3)
- [ ] Data migrated & validated (Phase 4)
- [ ] All providers migrated to Go in production (Phase 5)
- [ ] PHP decommissioned safely
- [ ] Documentation updated for operations team
- [ ] Monitoring & alerting fully configured
- [ ] Performance meets or exceeds PHP baseline
- [ ] Incident response procedures documented

---

## Stakeholders & Communication

| Stakeholder | Role | Update Cadence |
|-------------|------|---|
| Engineering Lead | Approval, architecture | Weekly sync |
| Ops/DevOps | Deployment, monitoring | Daily during Phase 5 |
| Provider Relations | Coordination, cutover | As needed |
| QA | Testing, validation | Per phase |

---

## Budget & Timeline Estimate

| Phase | Stories | Duration | Notes |
|-------|---------|----------|-------|
| **1. Discovery & OpenAPI** | 6 | 2-3 weeks | Parallel analysis |
| **2. Architecture** | 4 | 2-3 weeks | Design reviews |
| **3. Implementation** | 7 | 8-12 weeks | Depends on complexity per provider |
| **4. Data Migration** | 5 | 2-3 weeks | Validation critical |
| **5. Testing & Cutover** | 6 | 3-4 weeks | Gradual per provider |
| **TOTAL** | 28 stories | **17-25 weeks** | ~4-6 months aggressive |

---

## Next Steps

### Immediate (This Week)
1. ✅ Complete Phase 1 story planning (CASINO-1.1 to CASINO-1.6)
2. ✅ Assign story creators from @sm (Scrum Master)
3. ✅ Start Phase 1 story development

### Week 2-3
- Complete OpenAPI documentation for all 8 providers
- Review & validate against live PHP API
- Create master OpenAPI spec

### Week 4+
- Proceed with Phase 2 (Architecture)
- Assign @architect for design review
- Begin Phase 3 (Implementation)

---

## Appendix: Provider Integration Map

### Quick Reference - Provider Webhook Patterns

| Provider | Methods | Auth Pattern | Key Endpoints |
|----------|---------|--------------|---|
| Pragmatic Play | Generic routing | Provider signature | `/v1/webhooks/pragmatic-play/{endpoint}` |
| Evolution | 5 specific | Provider signature | `/v1/webhooks/evolution/{verb}` |
| PG Soft | 4 specific | Provider signature | `/v1/webhooks/pgsoft/{verb}` |
| Mancala | Generic routing | Provider signature | `/v1/webhooks/mancala/{endpoint}` |
| Digitain | Generic routing | Provider signature | `/v1/webhooks/digitain-rgs/{endpoint}` |
| Evoplay | Single endpoint | Provider signature | `/v1/webhooks/evoplay` |
| OpenBox | 4 PUT endpoints | Provider signature | `/v1/webhooks/openbox/seamless/{verb}` |
| Alternar | Redirect | Provider signature | `/v1/webhooks/alternar/` |

*All provider signatures use HMAC-SHA256 or similar secret-based validation (to be confirmed in Phase 1)*

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-05-08 | Morgan (PM) | Initial epic creation |

---

**Epic Ready for Story Breakdown → Assign to @sm for story creation**
