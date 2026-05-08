# Project Conventions — Branch Naming & Development Workflow

This file documents conventions that all developers should follow in this project.

---

## Branch Naming Convention

**All development work MUST use feature branches.** Never commit directly to `main`.

### Format

```
feat/epic-{EPIC_ID}/story-{STORY_ID}-{slug}
```

### Examples

```
feat/epic-1/story-1.1-update-readme
feat/epic-1/story-1.2-create-setup-script
feat/epic-2/story-2.3-add-kubernetes-deployment
```

### Components

| Component | Format | Example | Notes |
|-----------|--------|---------|-------|
| **Prefix** | `feat/` | `feat/` | Always use `feat/` (feature), not `fix/`, `docs/`, etc. |
| **Epic ID** | `epic-{N}` | `epic-1` | From epic file name (e.g., `docs/epics/1.aiox...`) |
| **Story ID** | `story-{N.M}` | `story-1.1` | From story file name (e.g., `docs/stories/1.1.story.md`) |
| **Slug** | `kebab-case` | `update-readme` | Short description, 2-4 words, lowercase with hyphens |

### Why This Matters

- **Traceability:** Branch name links directly to story and epic
- **Consistency:** All developers follow the same pattern
- **CI/CD Integration:** Automated tools parse branch names to link commits to stories
- **Code Review:** Reviewers immediately understand the scope of the PR

---

## Workflow Summary

1. **Find story** → `docs/stories/1.1.story.md`
2. **Create branch** → `git checkout -b feat/epic-1/story-1.1-update-readme`
3. **Implement** → Work on acceptance criteria
4. **Commit** → `git commit -m "[Story 1.1] feat: add AIOX prerequisites to README"`
5. **Push** → `git push origin feat/epic-1/story-1.1-update-readme`
6. **PR** → Create pull request (only @devops pushes to main)

---

## Commit Messages

Follow **Conventional Commits** format with **Story ID at the start**:

```
[Story <N.M>] <type>: <description>
```

### Format Rules

1. **Story ID first** — `[Story 1.1]` format
2. **Branch name must match** — Commit story ID = branch story ID
3. **Type after ID** — `feat:`, `fix:`, `docs:`, etc.
4. **Clear description** — What changed, not why

### Types

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `test:` — Tests
- `chore:` — Maintenance
- `refactor:` — Code refactoring

### Examples

**Branch:** `feat/epic-1/story-1.1-update-readme`
```
[Story 1.1] feat: add AIOX prerequisites to README
[Story 1.1] docs: add branch naming and development conventions
```

**Branch:** `feat/epic-1/story-1.2-create-setup-script`
```
[Story 1.2] feat: create npm setup script for AIOX initialization
[Story 1.2] docs: mark Story 1.2 as complete
```

**Branch:** `feat/epic-1/story-1.3-update-local-dev`
```
[Story 1.3] docs: add story workflow section to LOCAL_DEVELOPMENT.md
```

### Why Story ID at Start?

- **Consistent:** Commit log is scannable (all story IDs visible immediately)
- **Linked:** Commit ID = Story ID = Branch name (perfect traceability)
- **Sortable:** Git log can be sorted/filtered by story ID
- **Tools-friendly:** Automation can parse commits reliably
- **Human-friendly:** First thing you see when reading history

---

## Pushing & Code Review

**Important:** Only `@devops` agent can push to `main` branch.

1. Create feature branch locally
2. Push feature branch: `git push origin feat/epic-1/story-1.1-...`
3. Create PR via `@devops`: `@devops *push` or `gh pr create`
4. @devops reviews and merges to main

---

## Remember

- ✅ **DO:** Create feature branch with proper naming (feat/epic-{N}/story-{N.M}-{slug})
- ✅ **DO:** Put Story ID at START of commit message: `[Story 1.1] feat: description`
- ✅ **DO:** Match commit story ID to branch story ID
- ✅ **DO:** Push to feature branch first
- ❌ **DON'T:** Commit directly to main
- ❌ **DON'T:** Put Story ID at end of commit (wrong: `feat: description [Story 1.1]`)
- ❌ **DON'T:** Use `fix/`, `docs/`, `refactor/` prefixes for story work
- ❌ **DON'T:** Push to main yourself (only @devops)

---

**This convention is enforced by CI/CD.** Commits not following this pattern may be rejected by pre-push hooks.

For questions, refer to `CLAUDE.md` or ask your Scrum Master (@sm).
