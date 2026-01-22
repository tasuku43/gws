---
title: "gwst manifest ls"
status: planned
aliases:
  - "gwst man ls"
  - "gwst m ls"
migrated_from: "docs/spec/commands/manifest-ls.md"
---

## Synopsis
`gwst manifest ls [--root <path>] [--no-prompt]`

## Intent
List the workspace inventory in `gwst.yaml` (desired state) and show a lightweight per-workspace drift indicator by scanning the filesystem (actual state).

This is the primary "what do I have and is it applied?" command.

## Behavior
- Loads `<root>/gwst.yaml`; errors if missing or invalid.
- Scans `<root>/workspaces` to build the current filesystem state.
- For each workspace in the manifest, computes a status summary:
  - `applied`: no diff.
  - `missing`: present in manifest, missing on filesystem (would be `add` in plan/apply).
  - `drift`: present in both but differs (would be `update` in plan/apply).
- Also detects filesystem-only workspaces (present on filesystem, missing in manifest) and reports them as `extra`.
  - `extra` entries are informational only; use `gwst import` to capture them into the manifest, or `gwst apply` (with confirmation) to remove them.
- No changes are made (read-only).
- `--no-prompt` is accepted but has no effect (kept for CLI consistency).

## Output
Uses the common sectioned layout. No interactive UI is required.

- `Info` (optional): counts for `applied`, `missing`, `drift`, `extra`.
- `Result`: workspace list in inventory order (or sorted by ID), each with a short status tag.

Example:
```
Info
  • extra: PROJ-OLD

Result
  • PROJ-123 (applied)
  • PROJ-124 (drift)
  • PROJ-125 (missing)
```

## Success Criteria
- Inventory workspaces are listed and drift is accurately classified.

## Failure Modes
- Manifest missing or invalid.
- Filesystem or git errors while scanning workspaces (reported as warnings where possible).
