# Active Stories

This folder contains stories that are currently in progress — status is **Ready**, **InProgress**, or **InReview**.

## Lifecycle

```
docs/stories/active/     ← story lives here while being worked
    ↓ (QA gate: PASS)
docs/stories/completed/  ← story moves here when Done
```

## File Naming

`{EPIC-ID}.{STORY-NUM}.{kebab-title}.story.md`

Example: `TEMPL-001.6.k8s-kustomize-foundation.story.md`

## Status Values

| Status | Meaning |
|--------|---------|
| **Ready** | Validated by @po, waiting for @dev |
| **InProgress** | @dev is implementing |
| **InReview** | @qa is reviewing |

Once @qa issues a PASS or CONCERNS verdict, the story moves to `docs/stories/completed/` and its status becomes **Done**.
