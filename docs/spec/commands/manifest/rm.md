---
title: "gwst manifest rm"
status: planned
aliases:
  - "gwst man rm"
  - "gwst m rm"
pending:
  - interactive-selection-from-rm
  - apply-handoff
---

## Synopsis
`gwst manifest rm [<WORKSPACE_ID> ...] [--no-apply] [--no-prompt]`

## Intent
Remove workspace entries from the inventory (`gwst.yaml`) using an interactive UX (same intent as the legacy `gwst rm` selection), then reconcile the filesystem via `gwst apply` by default.

## Behavior (high level)
- Targets one or more workspaces:
  - With args: treat as the selected workspace IDs.
  - Without args: interactive multi-select (same UX as legacy `gwst rm`).
- Updates `<root>/gwst.yaml` by removing the selected workspace entries.
- By default, runs `gwst apply` to reconcile the filesystem with the updated manifest.
  - Destructive behavior is enforced by `gwst apply` (and `--no-prompt` must error if removals exist).
- With `--no-apply`, stops after rewriting `gwst.yaml` and prints a suggestion to run `gwst apply` next.

## Output (IA)
- Always uses the common sectioned layout from `docs/spec/ui/UI.md`.
- `Inputs`: selection inputs (workspace ids, warning indicators).
- `Plan`/`Apply`/`Result`: delegated to `gwst apply` when apply is run.

## Success Criteria
- Selected workspace entries are removed from `gwst.yaml`.
- When apply is run and confirmed, filesystem no longer contains removed workspaces.

## Failure Modes
- Workspace selection empty/canceled (interactive).
- Manifest write failure.
- `gwst apply` failure (git/filesystem).
