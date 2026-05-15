# 🏛️ Architecture & Technical Decisions

## Purpose

This folder contains **architecture documentation** including:
- Technical system design
- Architecture Decision Records (ADRs)
- Technology selections and justifications
- System components and interactions
- Technical constraints and considerations

Architecture docs provide the technical context needed for story creation, implementation, and design decisions.

## Contents

### Project Decisions (`project-decisions/` folder)
**Architecture Decision Records (ADRs)** for this specific project.

**File format:** `ADR-XXX-{title}.md`

**Example structure:**
```markdown
# ADR-001: Use Fiber for HTTP Framework

## Status
Accepted | Proposed | Deprecated

## Context
Why we needed to make this decision

## Decision
What we decided to do

## Rationale
Why we chose this option

## Consequences
Positive and negative impacts
```

## AIOX Usage

Architecture docs are consulted during:

### 1️⃣ Story Creation (@sm)
- **When:** Creating story from epic
- **What:** Read architecture docs to inform Dev Notes and tasks
- **Which:** 
  - tech-stack.md (for technology context)
  - system-architecture.md (for system design context)
  - ADRs in project-decisions/ (for design patterns)

### 2️⃣ Story Implementation (@dev)
- **When:** Implementing story tasks
- **What:** Follow architectural patterns documented
- **Which:**
  - Source-tree.md (for file organization)
  - ADRs (for design decisions)
  - System components documentation

### 3️⃣ QA Review (@qa)
- **When:** Validating story code architecture
- **What:** Verify implementation follows documented patterns
- **Which:** All architecture docs + ADRs

### 4️⃣ Design Decisions (@architect)
- **When:** Making new architecture decisions
- **What:** Document decision in ADR format
- **Which:** Create new file in project-decisions/

## Architecture Decision Records (ADRs)

### When to Write an ADR

Write an ADR when making decisions about:
- Technology choices (frameworks, libraries, databases)
- System design patterns (layering, messaging, caching)
- Naming conventions or code organization
- Integration patterns
- Scalability approaches

### ADR Format

```markdown
# ADR-XXX: [Concise Title]

## Status
Accepted | Proposed | Deprecated

## Context
What is the issue/question that forced this decision?
What are the driving forces?

## Decision
What is the decision/solution?
State it clearly and concisely.

## Rationale
Why did we choose this approach?
What are the trade-offs?
What alternatives did we consider?

## Consequences
What are the positive impacts?
What are the potential downsides?
What follow-up actions are needed?

## Related ADRs
Links to related decisions
```

## Framework vs Project Architecture

### Framework (`docs/framework/`)
- ✅ Company-wide tech stack
- ✅ Standard patterns everyone uses
- ✅ Doesn't change often
- ✅ Applies to all projects

### Project Architecture (this folder)
- ✅ Project-specific decisions
- ✅ Deviations from framework
- ✅ Technical solutions for business problems
- ✅ Evolves with project

## ✅ Checklist: Architecture Docs

Architecture documentation should:
- [ ] Justify every major decision with an ADR
- [ ] Include diagrams for complex systems (if helpful)
- [ ] Explain the WHY, not just the WHAT
- [ ] Document constraints and trade-offs
- [ ] Be reviewed by @architect before implementation
- [ ] Be updateable as architecture evolves
- [ ] Reference framework docs when applicable

## Key Principles

### 1. Decisions Must Be Documented
Before implementing a significant technical decision, write an ADR.

### 2. Context Is Critical
Help future readers understand what problem was being solved.

### 3. Trade-offs Matter
Document what you're gaining AND losing.

### 4. Reference, Don't Duplicate
ADRs reference framework docs, don't duplicate them.

## 🔗 Related Folders

- **Consumed by:** `docs/prd/` (context for epic design)
- **Consumed by:** `docs/stories/` (context for story Dev Notes)
- **Consumed by:** `docs/qa/` (validating architectural compliance)
- **References:** `docs/framework/` (for technology choices)

## Key Commands

- `@architect *create-adr` — Create new Architecture Decision Record
- `@architect *review-adr` — Review proposed ADR before implementation

---

*AIOX Framework: Architecture & Technical Decisions*  
*Reference documentation for technical decisions and system design*
