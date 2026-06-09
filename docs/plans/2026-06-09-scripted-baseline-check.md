# Scripted Baseline Check

status: completed

## Context

The repository had focused Go tests and Make targets, but it did not have a
scriptable baseline guard for required files, verification documentation, module
metadata, docs-plan metadata, or local secret/editor hygiene.

## Completed Scope

- Added `scripts/check-baseline.sh`.
- Wired the script into `make check` after lint, test, and build.
- Added ignore coverage for local editor metadata.
- Checked that local secret and editor metadata files are not tracked.
- Updated README, VISION, and CHANGES so the scripted gate is visible.

## Verification

- `scripts/check-baseline.sh`
- `go test ./...`
- `make check`
- `git diff --check`

## Follow-Up Candidates

- Add `go vet` if the project decides to broaden static analysis beyond the
  existing format, module, test, build, and baseline checks.
