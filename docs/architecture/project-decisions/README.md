# Architecture Decision Records (ADRs)

This folder contains Architecture Decision Records — lightweight documents that capture significant technical decisions, the context in which they were made, and their consequences.

## Format

Files follow the naming convention: `ADR-{NNN}-{kebab-title}.md`

Each ADR contains:
- **Status** — Draft / Accepted / Superseded / Deprecated
- **Context** — The situation that forced a decision
- **Decision** — What was decided
- **Consequences** — Trade-offs and implications

## Index

| ADR | Title | Status |
|-----|-------|--------|
| ADR-001 | Kubernetes Manifest Strategy | Accepted |

## When to Create an ADR

Create an ADR when making a decision that:
- Is hard to reverse (infrastructure, framework choice, API design)
- Has significant trade-offs
- Future team members will need to understand the "why"
- Affects multiple services or teams
