# Casino Proxy Migration Documentation

**Purpose:** Centralized documentation for the Casino Proxy PHP → Go migration (CASINO epics)

**Status:** CASINO-1 ✅ Complete | CASINO-2 🚀 In Progress | CASINO-3 🟡 In Pipeline | CASINO-4 ⏸️ Waiting

---

## 📚 Directory Structure

```
docs/casino-proxy/
├── README.md                                    (this file)
├── TEMPLATE-PROVIDER-BUSINESS-RULES.md         ✅ MANDATORY TEMPLATE
├── CASINO-1-migration-plan.md                  (Epic overview: OpenAPI docs)
├── EPIC-REORGANIZATION-PLAN.md                 (3-epic structure + timeline)
├── CASINO-DELIVERY-PLAN.md                     (Complete roadmap with templates)
│
├── phase-1-business-rules/                     (Business Rules - CASINO-2 Phase 1)
│   ├── pragmatic-play-rules.md                 ✅ Template Reference Implementation
│   ├── evolution-gaming-rules.md               (To create)
│   ├── pgsoft-rules.md                         (To create)
│   └── ... (6 more providers)
│
├── phase-2-technical-documentation/            (Endpoint Documentation - CASINO-2 Phase 2)
│   ├── pragmatic-play-balance.md               ✅ Template Reference Implementation
│   ├── pragmatic-play-authenticate.md          (To create)
│   ├── pragmatic-play-bet.md                   (To create)
│   └── ... (all endpoints for all 8 providers)
│
├── trace-matrices/                             (Traceability - CASINO-2 Phase 4)
│   └── (To create: rule → spec → code → test → Go impl)
│
└── validation-gates/                           (Test Validation - CASINO-2 Phase 5)
    └── (To create: test reports per provider)
```

---

## 🎯 How to Use This Documentation

### For CASINO-2 Phase 1-2 (Business Rules & Endpoint Docs)

**MANDATORY:** Follow `TEMPLATE-PROVIDER-BUSINESS-RULES.md` exactly

1. **Create Phase 1:** Extract business rules from PHP service
   - Output: `phase-1-business-rules/{provider}-rules.md`
   - Template section: Part 1: Business Rules Extraction
   - Reference: `pragmatic-play-rules.md`

2. **Create Phase 2:** Document each endpoint flow
   - Output: `phase-2-technical-documentation/{provider}-{endpoint}.md` (one file per endpoint)
   - Template section: Part 2: Technical Endpoint Documentation
   - Reference: `pragmatic-play-balance.md`

3. **Quality Gate:** Before proceeding to Phase 3
   - All rules documented with BR-* nomenclature
   - All endpoints documented with Mermaid flows
   - All error scenarios documented
   - Security checklist complete

### For CASINO-2 Phase 3 (Test Oracle)

**Input:** Completed Phase 1-2 documentation
**Output:** Java test module with 50+ test cases per provider
**Location:** `casino-proxy-test-oracle/` (separate module)

### For CASINO-2 Phase 4 (Trace Matrices)

**Input:** Completed Phase 1-2 documentation + Phase 3 tests
**Output:** YAML trace matrices mapping rule → spec → code → test → Go impl
**Location:** `trace-matrices/{provider}-trace-matrix.yaml`

### For CASINO-2 Phase 5 (Validation)

**Input:** Completed Phase 1-4 + Test Oracle
**Output:** Test validation reports proving PHP → Go parity
**Location:** `validation-gates/{provider}-validation-report.md`

---

## 📖 Key Documents Reference

### Epics & Overview
- **EPIC-REORGANIZATION-PLAN.md** — Why we split into 3 epics, how they sequence
- **CASINO-1-migration-plan.md** — OpenAPI documentation phase (complete)
- **CASINO-DELIVERY-PLAN.md** — Complete roadmap with all 5 phases + full endpoint inventory

### Templates (MANDATORY)
- **TEMPLATE-PROVIDER-BUSINESS-RULES.md** — Master template for all providers & all phases
  - Part 1: Business Rules Extraction format
  - Part 2: Endpoint Documentation format
  - Quality checklist for completion

### Reference Implementations (✅ Approved)
- **pragmatic-play-rules.md** — Approved example of Phase 1 (12 business rules extracted)
- **pragmatic-play-balance.md** — Approved example of Phase 2 (/balance endpoint documented)

---

## 🔄 Workflow: Add New Provider

### Step 1: Create Phase 1 (1-2 days per provider)
```bash
# Copy template structure and fill in
cp TEMPLATE-PROVIDER-BUSINESS-RULES.md phase-1-business-rules/{PROVIDER}-rules.md

# Extract from PHP code:
# app/Services/{PROVIDER}Service.php
# app/Http/Controllers/{PROVIDER}Controller.php

# Document all BR-* rules with dependencies and edge cases
```

### Step 2: Create Phase 2 (2-3 days per provider)
```bash
# For EACH endpoint, create a file:
# phase-2-technical-documentation/{provider}-{endpoint}.md

# For each file:
# - Mermaid flowchart (phases)
# - Map rules to phases
# - Document 5+ error scenarios
# - Provide complete request/response examples
# - Security checklist
```

### Step 3: Validate Phases 1-2
- [ ] All endpoints documented
- [ ] All rules extracted
- [ ] Code locations exact (file + line numbers)
- [ ] Examples accurate
- [ ] Error scenarios comprehensive
- [ ] Security checklist complete
- Ready for Phase 3 (Test Oracle)

### Step 4: Create Phase 3-5
- Phases 3, 4, 5 flow automatically once Phase 1-2 complete

---

## 📊 Provider Status

| Provider | Phase 1 (Rules) | Phase 2 (Endpoints) | Phase 3 (Tests) | Phase 4 (Matrix) | Phase 5 (Validation) |
|----------|---|---|---|---|---|
| **Pragmatic Play** | ✅ Complete | ✅ 1/9 | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **Evolution Gaming** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **PG Soft** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **Mancala** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **Digitain RGS** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **Evoplay** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **OpenBox** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |
| **Alternar** | ❌ To create | ❌ To create | ⏳ Pending | ⏳ Pending | ⏳ Pending |

---

## 🎯 Current Priority

**CASINO-2 Phase 1-2 for Evolution Gaming**
- Extract business rules from PHP code
- Document all 5 endpoints using template
- Reference: `TEMPLATE-PROVIDER-BUSINESS-RULES.md`
- Model: `pragmatic-play-rules.md` + `pragmatic-play-balance.md`

---

## ✅ Quality Gates

### Before Phase 3 (Test Oracle):
- [ ] Phase 1 complete: All rules extracted + documented
- [ ] Phase 2 complete: All endpoints documented with flows
- [ ] Rule IDs standardized: BR-[TYPE]-[ENDPOINT]-[CONCERN]-[SEQ]
- [ ] Code locations exact: File + line numbers
- [ ] Examples complete: Request/response pairs
- [ ] Error scenarios: 5+ per endpoint
- [ ] Security checklist: All items verified

### Before Phase 4 (Trace Matrices):
- [ ] Phase 3 complete: 50+ test cases per provider
- [ ] Tests passing: 100% against PHP legacy

### Before Phase 5 (Validation):
- [ ] Phase 4 complete: YAML trace matrices
- [ ] Go implementation ready
- [ ] Final test run: PHP → Go parity verification

---

## 📞 Support & Questions

**Reference Documents:**
- Template: `TEMPLATE-PROVIDER-BUSINESS-RULES.md`
- Examples: `pragmatic-play-rules.md`, `pragmatic-play-balance.md`
- Roadmap: `CASINO-DELIVERY-PLAN.md`
- Epic structure: `EPIC-REORGANIZATION-PLAN.md`

**Template Compliance:**
- All Phase 1-2 docs MUST follow the template
- No deviations without explicit approval
- Template is the single source of truth for Phases 1-2

---

**Last Updated:** 2026-05-12  
**Template Status:** v1.0 - APPROVED & MANDATORY  
**Next Update:** When first provider (Evolution Gaming) Phase 1-2 complete
