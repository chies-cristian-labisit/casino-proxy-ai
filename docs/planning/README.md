# Planning & Historical Documents

This folder contains pre-AIOX planning documents, architectural decisions, and historical context that informed the project's design and strategy.

## When to Consult These Documents

- **Understanding design rationale** — Why we chose certain approaches
- **Historical context** — Decisions made before adopting AIOX workflow
- **Architectural reference** — Pre-implementation architecture decisions
- **Migration strategy** — How we approached the monolith-to-microservices transition

## Documents

### RFQ-CG-MIGRATION-202603001.md
**Request for Quotation** — Migration requirements and scope from Cometa Gaming. Defines the business drivers and initial technical requirements for the backend template project.

### ADR-001-arquitetura-migracao-microservicos.md
**Architecture Decision Record #1** — Microservices migration architecture. Documents the decision to move from monolithic architecture to a microservices-based approach and the high-level architectural patterns chosen.

### ADR-002-go-backend-template-clean-architecture.md
**Architecture Decision Record #2** — Go backend template clean architecture. Defines the clean architecture approach used in this template, including layer definitions, dependency rules, and structural patterns for generated projects.

### application-first-prd.md
**Product Requirements Document** — Application-first approach PRD. Outlines the product vision, requirements, and feature set that guided the template design.

### aiox-implementation-plan.md
**AIOX Implementation Plan** — Strategic plan for integrating the AIOX agent framework into the template. Documents phased implementation approach, tooling decisions, and governance model.

## How These Relate to Current Development

These documents provide **reference context** for development work in the AIOX-driven workflow:

- **ADR references** — When implementing stories, architectural decisions documented here inform acceptance criteria
- **PRD alignment** — Features in development should trace back to requirements in the PRD
- **Migration context** — Understanding the original migration strategy helps when optimizing microservice boundaries

## Active Development

For current work, refer to:
- **`docs/stories/`** — Active development stories (story-driven workflow)
- **`docs/architecture/`** — Current architectural decisions and patterns
- **`docs/guides/`** — Developer guides and how-tos

---

**Note:** These are historical planning documents. Current architectural decisions and story-driven development happen in the AIOX workflow documented in `docs/README.md`.
