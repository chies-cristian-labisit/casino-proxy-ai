---
paths: "**/*"
---

# Git Commit Restrictions — LLM Attribution Policy

## Rule: No LLM Co-Authors in Commits

**Severity:** CRITICAL (enforced by `.githooks/commit-msg` hook)

### What This Means

LLMs (Claude, ChatGPT, Copilot, etc.) are **tools**, not co-authors. They must NOT be credited as contributors in commit messages.

❌ **Blocked patterns:**
```
Co-Authored-By: Claude <claude@anthropic.com>
Co-Authored-By: ChatGPT <gpt@openai.com>
authored by: Claude AI
Assisted by: GPT-4
```

✅ **Allowed approach:**
```
# List ONLY human co-authors
Co-Authored-By: John Doe <john@example.com>

# Mention tool usage in commit body if relevant (optional):
Uses Claude Code for code generation and refactoring.
```

### Why?

- **Attribution Clarity:** Commits must show who is RESPONSIBLE for the code
- **Legal/Audit:** Human ownership is legally and auditably clear
- **Professional Standards:** Industry norm is to credit humans, not tools
- **Accountability:** Humans answer for code quality and correctness

### Enforcement Mechanism

**Git Commit-Msg Hook:**
- File: `.githooks/commit-msg`
- Installed: Automatically via `scripts/setup.sh`
- Triggered: On every `git commit`
- Blocked patterns: Claude, ChatGPT, Copilot, Gemini, LLaMA, Grok, etc.

**Example block message:**
```
❌ COMMIT BLOCKED: LLM attribution detected

   ❗ LLMs are TOOLS, not co-authors
   ℹ️  Remove any 'Co-Authored-By' lines crediting AI tools

   Found pattern: Co-Authored-By.*Claude

   Valid commit message:
   - Only list HUMAN authors/co-authors
   - Mention tools in commit body if relevant (not as co-authors)
```

### Correct Commit Format

```bash
# GOOD: Human co-author only
git commit -m "feat: Add new feature

Co-Authored-By: Alice Smith <alice@example.com>"

# ALSO GOOD: Mention tool usage in body (not as author)
git commit -m "refactor: Improve code quality

Uses Claude Code for refactoring suggestions.

Co-Authored-By: Bob Johnson <bob@example.com>"

# BAD: Will be blocked
git commit -m "fix: Bug fix

Co-Authored-By: Claude <noreply@anthropic.com>"
```

### Rationale

In this project:
- **Human developers** make architectural decisions and own code quality
- **AI tools** (Claude Code, Copilot, etc.) assist with writing, refactoring, documentation
- **Clear attribution** means humans are accountable; tools are acknowledged in commit body if relevant

This mirrors industry practices at major tech companies where LLM-assisted code commits credit the human developer, not the tool.

### Exception Process

If there's a **legitimate reason** to document LLM assistance (research, documentation, non-production context):

1. Put it in the **commit body**, not as co-author:
   ```
   Research paper generated with assistance from Claude
   ```

2. **Never use** `Co-Authored-By` for tools

3. If blocked unfairly, contact @devops for hook adjustment

---

## Related Documents

- `.githooks/commit-msg` — Implementation
- `.claude/rules/git-push-restrictions.md` — Push restrictions
- `.claude/CLAUDE.md` — Git authority and workflows

---

**Effective:** 2026-05-12  
**Enforced By:** `.githooks/commit-msg` hook  
**Status:** ✅ Active
