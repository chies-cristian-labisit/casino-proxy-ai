# ✅ Quality Assurance (QA)

## Purpose

This folder contains **Quality Assurance gates**, test results, and code reviews that validate story completion. QA is the final checkpoint before stories are marked as "Done" and merged to the topic branch.

## Contents

### QA Gates (`gates/` folder)
Quality assurance verdicts for each completed story.

**File format:** `{storyId}-qa-gate.yml`

**Example content:**
```yaml
storyId: TEMPL-001.1
verdict: PASS | CONCERNS | FAIL | WAIVED
issues:
  - severity: low | medium | high
    category: code | tests | requirements | performance | security | docs
    description: "..."
    recommendation: "..."
```

### CodeRabbit Reports (`coderabbit-reports/` folder)
Automated code quality reviews from CodeRabbit CLI.

**File format:** `{storyId}-coderabbit-report.md`

## AIOX Workflow: QA Gate

### Phase 4: Quality Assurance (@qa)

**Trigger:** Story status is **InReview** (code implementation complete)

**Task:** `*qa-gate` command

**7 Quality Checks:**
1. ✅ **Code Review** — Patterns, readability, maintainability
2. ✅ **Unit Tests** — Coverage adequate, all passing
3. ✅ **Acceptance Criteria** — All story AC met
4. ✅ **No Regressions** — Existing functionality preserved
5. ✅ **Performance** — Within acceptable limits
6. ✅ **Security** — OWASP basics verified
7. ✅ **Documentation** — Updated if necessary

### QA Gate Verdicts

| Verdict | Score | Action |
|---------|-------|--------|
| **PASS** | All checks OK | Story moves to Done, ready for merge |
| **CONCERNS** | Minor issues | Approve with observations documented |
| **FAIL** | HIGH/CRITICAL issues | Return to @dev with feedback |
| **WAIVED** | Issues accepted | Approve with waiver documented (rare) |

### Pre-QA: CodeRabbit Self-Healing Loop

Before @qa gate, @dev runs CodeRabbit automated review:

1. **@dev** executes CodeRabbit (light mode)
   ```bash
   wsl bash -c 'cd /path && ~/.local/bin/coderabbit --prompt-only -t uncommitted'
   ```

2. **CRITICAL issues found?**
   - YES → Auto-fix (up to 2 iterations), re-run
   - NO → Document HIGH issues, proceed to @qa gate

3. **CRITICAL persist after 2 iterations?**
   - YES → HALT, manual intervention required
   - NO → Ready for @qa gate

### Story Status After QA

**If verdict is PASS or CONCERNS:**
- Story status updated to: **Done**
- Story file moved to: `docs/stories/completed/`
- Story ready for merge to topic branch

**If verdict is FAIL:**
- Story status returns to: **InReview**
- Story file stays in: `docs/stories/active/`
- @dev receives feedback, fixes issues
- Loop returns to QA gate (max 5 iterations)

## ✅ Checklist: QA Gate

Before @qa gate review:
- [ ] All story tasks completed and tested
- [ ] CodeRabbit pre-commit review passed
- [ ] Unit tests passing (`npm test` or `go test ./...`)
- [ ] Acceptance criteria verified
- [ ] File List in story updated
- [ ] Story status: InReview
- [ ] No console warnings or errors

During @qa gate:
- [ ] All 7 checks performed
- [ ] Code matches requirements
- [ ] No regressions detected
- [ ] Security basics validated
- [ ] Documentation complete
- [ ] Verdict recorded in `qa/gates/{storyId}-qa-gate.yml`

After QA approval:
- [ ] Story status updated to: Done
- [ ] Story moved to: `docs/stories/completed/`
- [ ] Ready for @devops to create PR and merge

## 📌 QA Loop (Iterative Review-Fix)

If QA gate fails, an iterative loop begins:

```
@qa review (verdict: FAIL)
    ↓
@dev fixes issues (up to 5 iterations)
    ↓
@qa re-review (verdict: PASS/CONCERNS/FAIL)
    ↓
If FAIL again → loop repeats
If PASS/CONCERNS → story moves to Done
```

**Command:** `*qa-loop {storyId}` — Orchestrate full loop

## 🔗 Related Folders

- **Previous:** `docs/stories/` — Stories being reviewed
- **Reference:** `docs/framework/` — Coding standards validated
- **Reference:** `docs/architecture/` — Architecture patterns verified
- **Next Phase:** Code merge to topic branch (handled by @devops)

## Key Commands

- `@qa *qa-gate {storyId}` — Execute quality gate review
- `@qa *qa-loop {storyId}` — Start review-fix loop
- `@dev *coderabbit-review` — Pre-commit CodeRabbit check
- `@devops *push` — Push approved story to remote

## Success Metrics

**Story passes QA when:**
- ✅ All 7 checks pass
- ✅ CodeRabbit CRITICAL issues resolved
- ✅ Test coverage adequate
- ✅ Acceptance criteria verified
- ✅ No regressions detected
- ✅ Security validated
- ✅ Documentation complete

---

## 📋 No Next Phase

**This is the final documentation phase.**

After QA approval:
1. Story moves to `docs/stories/completed/`
2. Code is merged to topic branch (handled by @devops)
3. Topic branch opens PR to main
4. Release tag applied to main branch

QA is the final gate before production. ✅

---

*AIOX Framework: Story-Driven Development*  
*Phase 4: Quality Assurance & Final Gate*
