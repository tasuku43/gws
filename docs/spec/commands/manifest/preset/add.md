---
title: "gwst manifest preset add"
status: planned
aliases:
  - "gwst manifest pre add"
  - "gwst manifest p add"
---

## Synopsis
`gwst manifest preset add [<name>] [--repo <repo> ...] [--no-prompt]`

## Intent
Create a preset entry in `gwst.yaml` without manual YAML editing.

## Notes
- This is the manifest-first replacement for the legacy `gwst preset add`.
- This command is inventory-only and does not run `gwst apply`.
