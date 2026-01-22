---
title: "gwst preset validate"
status: legacy
superseded_by: "gwst manifest preset validate"
migrated_from: "docs/spec/commands/preset-validate.md"
---

## Synopsis
`gwst preset validate`

## Intent
Validate `gwst.yaml` to catch malformed presets before use.

## Behavior
- Loads `<root>/gwst.yaml`; missing or unreadable file is reported as an issue.
- Parses YAML and reports errors if invalid.
- Checks for required fields:
  - top-level `presets` mapping exists.
  - each preset entry includes a non-empty `repos` list.
- Detects duplicate preset names in the YAML source.
- Validates preset names using the same rules as `gwst preset add`.
- Validates each repo spec via the existing repo spec normalization rules.
- Output uses the standard “Result” section with one bullet per issue; when no issues are found, prints “no issues found”.

## Success Criteria
- Returns success when `gwst.yaml` is valid.

## Failure Modes
- `gwst.yaml` missing/unreadable.
- YAML parse error.
- Missing required fields, duplicate preset names, invalid preset names, or invalid repo specs.
