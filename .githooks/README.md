# Git Hooks — Versionable Custom Hooks

This directory contains git hooks that are **committed to the repository** and automatically installed on clone.

## How It Works

1. **Versionable:** Unlike `.git/hooks/`, this directory is committed to git
2. **Auto-installed:** `scripts/setup.sh` configures `git config core.hooksPath .githooks`
3. **No tools needed:** Pure bash hooks, no external dependencies (unlike Husky)

## Current Hooks

### `pre-push`
**Purpose:** Block direct push to protected branches (`master`, `main`)

**Behavior:**
```bash
git push origin master  # ❌ BLOCKED
git push origin feature/my-feature  # ✅ Allowed
```

**Location:** `.githooks/pre-push`

**Documentation:** `.claude/rules/git-push-restrictions.md`

---

### `commit-msg`
**Purpose:** Block commits with LLM (Claude, ChatGPT, etc) attribution

**Behavior:**
```bash
git commit -m "feat: Add feature

Co-Authored-By: Claude <claude@anthropic.com>"  # ❌ BLOCKED

git commit -m "feat: Add feature

Co-Authored-By: John Doe <john@example.com>"  # ✅ Allowed
```

**Location:** `.githooks/commit-msg`

**Documentation:** `.claude/rules/git-commit-restrictions.md`

**Blocked patterns:**
- `Co-Authored-By: Claude`
- `Co-Authored-By: ChatGPT`
- `Co-Authored-By: Copilot`
- `Co-Authored-By: Gemini`
- `authored by: AI`
- And others (see hook for full list)

## Adding New Hooks

To add a new git hook (e.g., `commit-msg`):

1. Create `.githooks/commit-msg`
2. Make it executable: `chmod +x .githooks/commit-msg`
3. Add implementation
4. Commit and push
5. On next clone, it's automatically installed

### Hook Template

```bash
#!/bin/bash
# Git hook: {hook-name}
# Purpose: {what this does}

# Your logic here

exit 0  # or exit 1 to block
```

## Testing a Hook

```bash
# Test pre-push hook without actually pushing
bash .githooks/pre-push

# Or manually trigger it
git push  # Will run the hook
```

## Troubleshooting

**Hook not running:**
```bash
# Verify git config
git config core.hooksPath

# Should output: .githooks
```

**Make hook executable:**
```bash
chmod +x .githooks/{hook-name}
```

**Manual setup (if needed):**
```bash
git config core.hooksPath .githooks
```

## Related

- **Setup Script:** `scripts/setup.sh` (installs on clone)
- **Git Push Rules:** `.claude/rules/git-push-restrictions.md`
- **CLAUDE.md:** `.claude/CLAUDE.md` (push authority section)

---

**Format:** Bash scripts (no external dependencies)  
**Lifecycle:** Committed to git, auto-installed via `scripts/setup.sh`  
**Last Updated:** 2026-05-12
