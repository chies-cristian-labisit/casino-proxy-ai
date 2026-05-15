# 📋 Product Requirements & Epics (PRD)

## Purpose

This folder contains **Product Requirements Documents (PRDs)** and **Epic definitions** that drive the development roadmap. Epics are high-level initiatives that are broken down into individual user stories.

## Contents

### Epics (`epics/` folder)
- **epic-TEMPL-001-standardize-go-templates.md** — Epic for standardizing Go backend templates across the organization
  - Contains business goals, scope, 8 stories, timeline, and success criteria
  - Status: Planning

## AIOX Workflow: PRD → Stories

### Documents in this folder should:
1. ✅ **Define business goals** — Why this epic matters
2. ✅ **List all stories** — Breakdown into implementable chunks
3. ✅ **Set acceptance criteria** — How to measure success
4. ✅ **Identify timeline** — Deadlines and milestones
5. ✅ **Document scope** — What's included and excluded

## 📌 Next Phase: Story Creation

**When an epic is approved:**

1. **@sm (Scrum Master)** executes `*draft` command
   - Creates individual story files in `docs/stories/`
   - Each story gets unique ID (e.g., TEMPL-001.1, TEMPL-001.2, etc.)
   - Stories reference the epic for context

2. **Stories are created with:**
   - User story format (As a... I want... so that...)
   - Acceptance criteria from epic
   - Technical tasks for implementation
   - Dev Notes from architecture docs

3. **Stories flow to `docs/stories/` folder**
   - Each story file: `{epicNum}.{storyNum}.{title}.story.md`
   - Stories enter "Draft" status

## ✅ Checklist: Before Moving to Story Creation

- [ ] Epic has clear title and description
- [ ] Business goals documented
- [ ] Scope clearly defined (what's IN and OUT)
- [ ] All stories listed with brief descriptions
- [ ] Timeline and milestones set
- [ ] Success criteria defined
- [ ] Epic approved by @po (Product Owner)

## 🔗 Related Folders

- **Next:** `docs/stories/` — Individual user stories (created from epics here)
- **Reference:** `docs/framework/` — Coding standards for story formatting
- **Reference:** `docs/architecture/` — Technical context for epic design

---

*AIOX Framework: Story-Driven Development*  
*Phase 1: Epics & Requirements*
