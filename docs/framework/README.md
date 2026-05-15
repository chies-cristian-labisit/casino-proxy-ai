# 🏗️ Framework Documentation

## Purpose

This folder contains **framework-level documentation** that establishes standards and conventions for all development work. These documents are consulted throughout the AIOX workflow and should be read by all team members.

## Contents

### Essential Framework Docs (to be created/maintained)

- **coding-standards.md** — Code style, naming conventions, import rules
- **tech-stack.md** — Technology choices, versions, and rationale
- **source-tree.md** — Project directory structure and file organization

## AIOX Usage

Framework docs are **reference material** consulted during:

### 1️⃣ Story Creation (@sm)
- **When:** Creating story from epic
- **What:** Read framework docs to populate Dev Notes
- **Files:** coding-standards.md, tech-stack.md, source-tree.md

### 2️⃣ Story Implementation (@dev)
- **When:** Implementing story tasks
- **What:** Follow coding-standards during development
- **Files:** coding-standards.md (naming, structure, patterns)

### 3️⃣ QA Review (@qa)
- **When:** Validating story code quality
- **What:** Verify code follows coding-standards
- **Files:** coding-standards.md (patterns, structure, best practices)

### 4️⃣ Architecture Decisions (@architect)
- **When:** Evaluating technical choices
- **What:** Reference tech-stack for approved technologies
- **Files:** tech-stack.md, source-tree.md

## Framework vs Project Decisions

### Framework Docs (this folder)
- ✅ Company-wide standards
- ✅ Coding conventions
- ✅ Technology stack decisions
- ✅ Project structure
- ✅ Stable, long-lived

### Project Docs (`docs/architecture/project-decisions/`)
- ✅ Project-specific Architecture Decision Records (ADRs)
- ✅ Justifications for deviations from framework
- ✅ Project-specific patterns and constraints
- ✅ Evolves with project

## ✅ Checklist: Framework Docs Quality

Framework documentation should:
- [ ] Be clear and specific (no vague guidance)
- [ ] Include examples where possible
- [ ] Document rationale (the WHY, not just the WHAT)
- [ ] Be version-controlled and reviewed
- [ ] Be updated when standards change
- [ ] Be enforced in QA gates

## 🔗 Related Folders

- **Consumed by:** `docs/prd/` (story creation)
- **Consumed by:** `docs/stories/` (story tasks and dev notes)
- **Consulted by:** `docs/architecture/` (when creating project decisions)
- **Validated by:** `docs/qa/` (QA gates verify compliance)

## Key Guidelines

When writing framework documentation:

1. **Be Prescriptive** — "DO X" not "you might consider X"
2. **Explain Why** — Help developers understand the reasoning
3. **Show Examples** — Concrete examples beat abstract descriptions
4. **Keep Current** — Update when standards evolve
5. **Make Searchable** — Use clear headings and organization

---

*AIOX Framework: Reference Documentation*  
*Framework docs are consulted throughout development lifecycle*
