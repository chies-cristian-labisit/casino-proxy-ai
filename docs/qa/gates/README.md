# QA Gates

This folder stores the quality gate verdict files for each completed story. A QA gate is the final checkpoint before a story is marked Done and merged.

## Format

Files follow the naming convention: `{storyId}-qa-gate.yml`

```yaml
storyId: TEMPL-001.1
verdict: PASS          # PASS | CONCERNS | FAIL | WAIVED
reviewedBy: "@qa"
reviewedAt: "2026-05-14"
issues:
  - severity: low      # low | medium | high | critical
    category: docs     # code | tests | requirements | performance | security | docs
    description: "Missing inline comment on complex algorithm"
    recommendation: "Add a brief comment explaining the approach"
```

## Verdict Meanings

| Verdict | Meaning |
|---------|---------|
| **PASS** | All checks OK — story moves to Done |
| **CONCERNS** | Minor issues noted — approved with observations |
| **FAIL** | HIGH or CRITICAL issues — returned to @dev |
| **WAIVED** | Issues accepted by stakeholder — rarely used |
