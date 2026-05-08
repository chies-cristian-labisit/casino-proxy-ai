# AIOX Story Development Cycle (SDC) — Mandatory Workflow

**CRITICAL:** This workflow is NOT optional. Every story MUST follow these phases in order.

---

## 4-Phase Story Development Cycle (SDC)

```
Phase 1: Create (@sm)
    ↓ Draft stories
Phase 2: Validate (@po)
    ↓ 10-point checklist, GO/NO-GO decision
Phase 3: Implement (@dev)
    ↓ Code implementation, feature branches
Phase 4: QA Gate (@qa)
    ↓ Testing, verdict (PASS/CONCERNS/FAIL)
    ↓ 
PR & Push (@devops)
    ↓ Create PR, merge to main
Done
```

**Reference:** `.claude/rules/workflow-execution.md` (Workflow Execution — Detailed Rules)

---

## Phase Details

### Phase 1: Create (@sm)
**Who:** Scrum Master (@sm / River)  
**What:** Create story markdown files from epic requirements  
**Output:** `docs/stories/{N}.{M}.story.md` (Draft status)  
**Command:** `@sm *create-story` or `@sm *draft`

### Phase 2: Validate (@po)
**Who:** Product Owner (@po / Pax)  
**What:** Apply 10-point validation checklist  
**Criteria:**
1. Clear title
2. Acceptance criteria (well-defined, measurable)
3. Testable outcomes
4. No ambiguity
5. Business value
6. Dependencies (clear or none)
7. Technical feasibility
8. Proper scope
9. Measurable success
10. Independent (or explicit blockers)

**Decision:** GO (>=7 points) or NO-GO (fixes required)  
**Command:** `@po *validate-story-draft`

### Phase 3: Implement (@dev)
**Who:** Developer (@dev / Dex)  
**What:** Code implementation on feature branch  
**Branch naming:** `feat/epic-{N}/story-{N.M}-{slug}`  
**Modes:**
- **Interactive** — Step-by-step with @dev guidance
- **YOLO** — Rapid implementation, self-paced
- **Pre-Flight** — Planned approach with @architect first

**Status flow:** Draft → Ready → InProgress  
**Commits:** Use `[Story N.M]` tag in commit messages  
**Command:** `@dev` (then work in YOLO mode)

### Phase 4: QA Gate (@qa)
**Who:** QA Lead (@qa / Quinn)  
**What:** Quality assurance testing and verification  
**7 Quality Checks:**
1. Acceptance criteria met (all checkboxes)
2. File list completeness
3. Functionality testing (manual or automated)
4. Code quality (linting, types)
5. Documentation accuracy
6. Edge cases handled
7. No regressions

**Verdict Options:**
- **PASS** → Story ready to merge
- **CONCERNS** → Minor issues, proceed with caution
- **FAIL** → Blocker found, return to @dev
- **WAIVED** → Explicit exemption (documented)

**If FAIL:** Use QA Loop (`@qa *qa-loop {storyId}`) for iterative fixes  
**Status:** InProgress → InReview → Done  
**Command:** `@qa *qa-gate {storyId}`

---

## Recommended: Option B — Batch QA

**Fastest while maintaining quality:**

1. **Implement all stories** (1.1, 1.2, 1.3, 1.4, 1.5)
   - Each on separate feature branch
   - All committed locally
   - NO PRs yet

2. **QA gate all stories** (batch testing)
   - Run `@qa *qa-gate 1.1`, `@qa *qa-gate 1.2`, etc.
   - Parallel testing possible
   - Document any FAIL verdicts

3. **Fix any FAIL verdicts** (if needed)
   - @dev fixes reported issues
   - Re-run QA gate
   - Max 5 iterations per story (then escalate)

4. **Create PRs and push** (all at once or staged)
   - `@devops` creates PRs
   - Reference story IDs in PR descriptions
   - Merge to main when all QA PASS

---

## Before Every Step: Verify Workflow

**MANDATORY CHECKLIST** — Answer these before proceeding:

- [ ] What phase am I in? (Create / Validate / Implement / QA / PR)
- [ ] What agent should lead this phase? (@sm / @po / @dev / @qa / @devops)
- [ ] Have all prior phases completed? (No skipping)
- [ ] What is the expected status change? (Draft → Ready → InProgress → InReview → Done)
- [ ] What command should I run? (*create-story / *validate-story-draft / [implement] / *qa-gate / *push)
- [ ] Is my branch properly named? (feat/epic-{N}/story-{N.M}-{slug})
- [ ] Have I referenced story ID in commits/PRs? ([Story N.M])

---

## Common Mistakes (Avoid These!)

❌ **Committing to main** — Always use feature branches  
❌ **Skipping QA gate** — Every story MUST be QA-tested  
❌ **Implementing without validation** — @po must GO before @dev starts  
❌ **Creating PR before QA approval** — QA verdict must be PASS  
❌ **Forgetting story ID in commits** — Link commits to stories  
❌ **Mixing stories in one branch** — One story per feature branch  
❌ **@dev pushing to main** — Only @devops can push (has exclusive authority)

---

## Workflow State Machine

```
┌─────────┐
│ Draft   │  Created by @sm
└────┬────┘
     │ @po validates
     ↓
┌─────────┐
│ Ready   │  Approved for implementation
└────┬────┘
     │ @dev implements on feat/ branch
     ↓
┌─────────────┐
│ InProgress  │  Feature branch commits
└────┬────────┘
     │ All commits done
     ↓
┌─────────────┐
│ InReview    │  Awaiting QA gate
└────┬────────┘
     │ @qa runs QA gate
     ├─ FAIL → @dev fixes, loop back
     ├─ CONCERNS → Proceed with note
     ↓
┌─────────┐
│ Done    │  QA PASS, ready for PR
└────┬────┘
     │ @devops creates PR & merges
     ↓
┌─────────────┐
│ Merged      │  On main branch
└─────────────┘
```

---

## Story Status Values

| Status | Meaning | Who Sets | Next Action |
|--------|---------|----------|------------|
| Draft | Created, not validated | @sm | @po validates |
| Ready | Validated, approved for dev | @po | @dev implements |
| InProgress | Implementation started | @dev | Complete work, request QA |
| InReview | Awaiting QA testing | @qa | Run QA gate |
| Done | QA approved, ready to merge | @qa | @devops creates PR |
| Merged | On main branch | @devops | Archive |

---

## Remember

**This workflow exists to maintain code quality and traceability.**

- Every phase has a reason
- Every agent has authority in their phase
- No skipping, no shortcuts
- Document decisions in story files
- Reference story IDs everywhere

**When in doubt, check `.claude/rules/workflow-execution.md`**

---

*Last Updated: 2026-04-30*  
*Enforced by AIOX Constitution Article III — Story-Driven Development*
