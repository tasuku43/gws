---
title: "gwst manifest gc"
status: planned
aliases:
  - "gwst man gc"
  - "gwst m gc"
pending:
  - rules-implementation
  - reason-format
  - confirmation-ux
---

## Synopsis
`gwst manifest gc [--no-apply] [--no-prompt]`

## Intent
Conservatively remove workspace entries from `gwst.yaml` that are highly likely safe to delete, then (by default) run `gwst apply` to reconcile the filesystem.

This command is intentionally separated from manual removal flows (`gwst manifest rm`), which remain the explicit/human-judgment path.

## Scope / Non-goals
- **GC**: automatic, bulk, conservative. Exclude when in doubt.
- No implicit `git fetch` / `git remote prune`.
- No per-item interactive selection (single bulk decision + apply confirmation).

## Definitions
- **Clean**: no uncommitted changes in any repo.
- **Unpushed**: local branch is ahead of upstream.
- **Unknown**: status cannot be determined (e.g., git error, no upstream, detached HEAD).

## Target branch selection (per repo)
For each repo, determine a merge target:
1) If `repos[].base_ref` is set in `gwst.yaml`, use it (`origin/<branch>`).
2) Otherwise, use `origin/<default>` resolved from `refs/remotes/origin/HEAD`.

## Base exclusions (per workspace)
Any workspace containing any repo in one of these states is excluded from GC:
- Dirty
- Unpushed
- Diverged
- Unknown

## Safe-to-remove Rule (initial, extensible)
Rules are predicates that return `(matched bool, reason string)` and are evaluated over a shared per-repo snapshot (avoid re-running expensive git commands per rule).

Initial rule (single rule):
1) **Strict merged into target**: repo `HEAD` is an ancestor of `origin/<target>` and `HEAD != origin/<target>`.
   - This prevents deleting "created-only" workspaces where no commits have been made (even if `HEAD` equals the target).
   - Reason: `merged`

A workspace is a candidate only if:
- all repos pass base exclusions, and
- every repo matches at least one rule (initially: strict merged).

## Behavior
- Scans workspaces present in `gwst.yaml`.
- For each workspace:
  - Loads per-repo state (clean/unpushed/etc).
  - Computes per-repo rule results and reasons.
- Prints candidates with reasons (always shown before manifest mutation).
- Updates `gwst.yaml` by removing all candidates.
- By default, runs `gwst apply` once for the entire root.
- If apply is canceled/declined at confirmation, restores the previous `gwst.yaml`.

## Flags
- `--no-apply`: update `gwst.yaml` and exit (do not run `gwst apply`).
- `--no-prompt`: forwarded to `gwst apply` when apply is run (behavior follows `gwst apply` spec).

## Output
- `Inputs`/`Info`: scanned / candidates / skipped counts.
- Candidate list: workspace id + short reasons (e.g., `[merged]`) and target context.
- `Plan`/`Apply`/`Result`: delegated to `gwst apply` when apply is run.

## Failure Modes
- Any git status or rule error => treat as unknown, skip, and report warning.
- Manifest write failure.
- Apply failure (manifest remains updated; users can re-run `gwst apply`).
