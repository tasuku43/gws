---
title: "gwst manifest preset rm"
status: planned
aliases:
  - "gwst manifest pre rm"
  - "gwst manifest p rm"
---

## Synopsis
`gwst manifest preset rm [<name> ...] [--no-prompt]`

## Intent
Remove preset entries from `gwst.yaml`.

## Notes
- This is the manifest-first replacement for the legacy `gwst preset rm`.
- This command is inventory-only and does not run `gwst apply`.
