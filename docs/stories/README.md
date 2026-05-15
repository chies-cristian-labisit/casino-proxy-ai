# 📖 Development Stories

## Purpose

This folder contains **user stories** that define specific development work. Stories are created from epics and follow the AIOX Story-Driven Development (SDD) workflow. Each story is a complete, actionable unit of work assigned to @dev for implementation.

## Contents

### Active Stories (`active/` folder)
Stories currently being worked on or ready for development.

### Completed Stories (`completed/` folder)
Stories that have passed QA gates and are merged to main.

### Epic Overviews (`epics/` folder)
Summary pages for each epic (e.g., `1-go-backend-template.md`).

### Individual Stories (root)
Stories following naming: `{epicNum}.{storyNum}.{title}.story.md`

**Example:**
- `TEMPL-001.1.migrate-git-remote.story.md` — Story 1 of epic TEMPL-001

## Story Format & Sections

Each story file MUST include:

```markdown
# Story {ID}: {Title}

## Status
Draft | Ready | InProgress | InReview | Done

## Story
As a [user], I want [capability], so that [benefit]

## Acceptance Criteria
1. [Testable requirement]
2. [Testable requirement]
...

## Tasks / Subtasks
- [ ] Task 1 — Description (AC: X)
- [ ] Task 2 — Description (AC: X)

## Dev Notes
- Architecture references
- Technical constraints
- Dependencies

## File List
- Modified files tracked here
```

## AIOX Workflow: Story Lifecycle

### Phase 1: Create (@sm)
**Input:** Epic from `docs/prd/epics/`  
**Task:** `*draft` command  
**Output:** Story file in `docs/stories/` with status: **Draft**

1. @sm reads epic from `docs/prd/epics/`
2. @sm creates story file with template
3. @sm loads architecture docs from `docs/framework/` & `docs/architecture/`
4. @sm populates tasks and technical details
5. Story saved with status: **Draft**

### Phase 2: Validate (@po)
**Input:** Story file (Draft status)  
**Task:** `*validate-story-draft` command  
**Output:** Story status: **Ready** (if approved)

1. @po reviews story against 10-point checklist
2. @po verifies acceptance criteria are testable
3. If GO: updates status to **Ready** in story file
4. If NO-GO: returns to @sm with required fixes

### Phase 3: Implement (@dev)
**Input:** Story file (Ready status)  
**Task:** `*develop` command  
**Output:** Completed story + code changes

1. @dev reads story from `docs/stories/`
2. @dev implements tasks sequentially
3. @dev runs all tests and validations
4. @dev updates File List section
5. @dev marks story status: **InReview**

### Phase 4: QA Gate (@qa)
**Input:** Story file (InReview status) + code changes  
**Task:** `*qa-gate` command  
**Output:** QA verdict + story status update

1. @qa reviews code against 7 quality checks
2. @qa verifies acceptance criteria met
3. If PASS: updates status to **Done**, story moves to `completed/`
4. If FAIL: returns to @dev with feedback

## 📌 Story Status Progression

```
Draft → Ready → InProgress → InReview → Done
  ↓       ↓         ↓           ↓
  @sm     @po       @dev        @qa
create  validate  implement   quality-gate
```

## ✅ Checklist: Creating a New Story

Before story creation by @sm:
- [ ] Epic is approved and in `docs/prd/epics/`
- [ ] Epic has clear breakdown into stories
- [ ] Architecture docs available in `docs/framework/` & `docs/architecture/`
- [ ] Team has reviewed epic scope

During story creation by @sm:
- [ ] Story title is clear and specific
- [ ] User story format (As a... I want... so that...)
- [ ] Acceptance criteria are testable (Given/When/Then)
- [ ] Tasks are sequential and implementation-ready
- [ ] Dev Notes include architecture references
- [ ] File List tracks all affected files
- [ ] Status set to: Draft

Before handoff to @dev:
- [ ] @po has validated story (status: Ready)
- [ ] Story branch created: `feat/{storyId}-*`
- [ ] Story assigned to @dev

## 🔗 Related Folders

- **Previous:** `docs/prd/epics/` — Epics that stories are created from
- **Next:** `docs/qa/` — QA gates after story implementation
- **Reference:** `docs/framework/` — Coding standards to follow in stories
- **Reference:** `docs/architecture/` — Technical context for story tasks

## Key Commands

- `@sm *draft` — Create new story from epic
- `@po *validate-story-draft` — Validate story readiness
- `@dev *develop` — Implement story
- `@qa *qa-gate` — Quality assurance review
- `@devops *push` — Push completed story to remote

---

*AIOX Framework: Story-Driven Development*  
*Phase 2: Story Preparation & Implementation*
