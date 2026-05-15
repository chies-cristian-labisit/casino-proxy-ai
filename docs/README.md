# 📚 Documentation Structure

Welcome to the Cometa Gaming Go Backend Template documentation. This folder contains all project documentation organized according to the **AIOX Framework** for story-driven development.

## Quick Navigation

| Folder | Purpose | Phase | Read First |
|--------|---------|-------|-----------|
| **prd/** | Epics & product requirements | Planning | Epic definitions |
| **stories/** | Development stories & tasks | Preparation → Implementation | Story templates |
| **qa/** | Quality assurance & gates | Validation | QA verdicts |
| **planning/** | Historical planning & architectural decisions | Reference | Project context |
| **framework/** | Company standards & conventions | Reference | Coding standards |
| **architecture/** | Technical decisions & design | Reference | System design |

## 🔄 AIOX Workflow Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    PRODUCT PLANNING                         │
│                    (docs/prd/)                              │
│  • Define epics and product requirements                    │
│  • Set business goals and scope                             │
│  • Create acceptance criteria                               │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    STORY CREATION                           │
│              (docs/stories/) Phase 1                        │
│  • @sm creates stories from epics                           │
│  • Stories: Draft status                                    │
│  • Tasks and technical details defined                      │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    STORY VALIDATION                         │
│              (docs/stories/) Phase 2                        │
│  • @po reviews and validates stories                        │
│  • Stories: Ready status (if approved)                      │
│  • Ready for implementation                                 │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    IMPLEMENTATION                           │
│              (docs/stories/) Phase 3                        │
│  • @dev implements story tasks                              │
│  • Code changes committed to task branch                    │
│  • Stories: InReview status                                 │
│  • CodeRabbit pre-commit review                             │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    QUALITY ASSURANCE                        │
│                    (docs/qa/)                               │
│  • @qa performs 7 quality checks                            │
│  • Stories: Done status (if passed)                         │
│  • Stories moved to docs/stories/completed/                 │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    CODE MERGE & RELEASE                     │
│                  (DevOps / @devops)                         │
│  • @devops creates PR to topic branch                       │
│  • Topic branch merges to main                              │
│  • Release tag applied                                      │
└─────────────────────────────────────────────────────────────┘
```

## 📖 Reading Guide by Role

### Product Owners (@po, @pm)
1. Start: `docs/prd/README.md` — Understand epic structure
2. Read: `docs/prd/epics/` — Review epics and requirements
3. Read: `docs/stories/README.md` — Validate stories before handoff
4. Consult: `docs/architecture/README.md` — For technical context

### Scrum Masters (@sm)
1. Start: `docs/prd/README.md` — Understand epics
2. Read: `docs/stories/README.md` — Story creation workflow
3. Consult: `docs/framework/README.md` — Coding standards for story tasks
4. Consult: `docs/architecture/README.md` — Technical context

### Developers (@dev)
1. Start: `docs/stories/README.md` — Understand story format
2. Read: Active story in `docs/stories/active/`
3. Consult: `docs/framework/README.md` — Coding standards
4. Consult: `docs/architecture/README.md` — Technical design
5. Implement according to story tasks

### QA (@qa)
1. Start: `docs/qa/README.md` — Understand QA gates
2. Read: Story in `docs/stories/active/`
3. Consult: `docs/framework/README.md` — Coding standards
4. Record: Verdict in `docs/qa/gates/`

### Architects (@architect)
1. Start: `docs/architecture/README.md` — Overview
2. Read: `docs/architecture/project-decisions/` — ADRs
3. Consult: `docs/planning/README.md` — Historical architecture decisions that informed current design
4. Consult: `docs/framework/README.md` — Framework standards
5. Create: New ADRs for major decisions

### DevOps (@devops)
1. Start: `docs/stories/README.md` — Understand story branches
2. Read: Story file for branch info
3. Create: PRs from task branches → topic branch
4. Manage: Release tags and merges

## 🗂️ Folder Structure

```
docs/
├── README.md (this file)
│
├── prd/                          # PHASE 1: Product Requirements
│   ├── README.md                # Getting started with epics
│   └── epics/                    # Epic definitions
│       └── epic-TEMPL-001-*.md
│
├── stories/                      # PHASE 2-3: Story Preparation & Implementation
│   ├── README.md                # Story workflow & lifecycle
│   ├── epics/                    # Epic summary pages
│   │   └── 1-go-backend-template.md
│   ├── active/                   # In-progress stories
│   ├── completed/                # Finished stories (Done status)
│   └── {storyId}.{title}.story.md  # Individual story files
│
├── qa/                           # PHASE 4: Quality Assurance
│   ├── README.md                # QA gates & verdicts
│   ├── gates/                    # QA gate results
│   │   └── {storyId}-qa-gate.yml
│   └── coderabbit-reports/       # Automated code reviews
│       └── {storyId}-report.md
│
├── planning/                     # REFERENCE: Pre-AIOX Planning
│   ├── README.md                # Planning docs overview
│   ├── RFQ-CG-MIGRATION-*.md    # Project requirements & scope
│   ├── ADR-001-*.md             # Migration architecture decisions
│   ├── ADR-002-*.md             # Template architecture decisions
│   ├── application-first-prd.md # Product vision & requirements
│   └── aiox-implementation-plan.md # Strategic implementation plan
│
├── framework/                    # REFERENCE: Company Standards
│   ├── README.md                # Framework docs overview
│   ├── coding-standards.md       # Code style & conventions
│   ├── tech-stack.md            # Technology choices
│   └── source-tree.md           # Directory structure
│
└── architecture/                 # REFERENCE: Technical Design
    ├── README.md                # Architecture overview
    └── project-decisions/        # Architecture Decision Records
        └── ADR-XXX-*.md
```

## ✅ Key Principles

### 1. Story-Driven Development
Everything starts with a **story**. Stories are the unit of work in this framework.

### 2. Clear Handoff Points
Each phase has clear entrance/exit criteria. Use the folder READMEs as checklists.

### 3. Reference vs Implementation
- **Reference docs** (`framework/`, `architecture/`, `planning/`) — Read for context and decision rationale
- **Implementation docs** (`stories/`) — Read during specific phase
- **Planning docs** (`prd/`) — Read to understand business context
- **Historical planning** (`planning/`) — Read to understand pre-AIOX decisions and project evolution

### 4. Status Is Sacred
Story status (Draft → Ready → InProgress → InReview → Done) drives the workflow.

### 5. File Organization Mirrors Workflow
Folders are organized by **phase**, not by type. This guides you naturally through the process.

## 🎯 Getting Started

### If you're creating a new initiative:
1. Go to `docs/prd/README.md` and create an epic
2. Reference `docs/framework/README.md` for standards

### If you're implementing a story:
1. Find your story in `docs/stories/`
2. Read story format in `docs/stories/README.md`
3. Follow `docs/framework/README.md` while coding
4. Reference `docs/architecture/README.md` for design patterns

### If you're reviewing a story:
1. Read QA checklist in `docs/qa/README.md`
2. Consult `docs/framework/README.md` for standards
3. Record verdict in `docs/qa/gates/`

## 🔗 Quick Links

| Task | Location |
|------|----------|
| Understand AIOX framework | This file |
| Create an epic | `docs/prd/README.md` |
| Create a story | `docs/stories/README.md` |
| Implement a story | Story file in `docs/stories/` |
| Review code quality | `docs/qa/README.md` |
| Check coding standards | `docs/framework/README.md` |
| Understand architecture | `docs/architecture/README.md` |
| Read project context & decisions | `docs/planning/README.md` |

---

## 📋 Quick Checklist: Before You Start

- [ ] You've read the README in the **folder you're working in**
- [ ] You understand the **next phase** and what documents go there
- [ ] You've consulted **reference docs** (framework & architecture) if needed
- [ ] Your work matches the **standards** documented here
- [ ] You're ready to hand off to the **next phase**

---

*AIOX Framework: Story-Driven Development*  
*Documentation organized by phase and purpose*  
*Last updated: 2026-05-13*
