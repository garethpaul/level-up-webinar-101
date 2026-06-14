# Integrate Webinar Security Stack

status: pending

## Context

The location-independent verification stack and the open security integration
stack diverged from the same earlier base. The current PR #5 line therefore
does not contain PR #2's Twilio error redaction, patched Go toolchain, or hosted
`govulncheck` enforcement even though tracker text claimed those boundaries.

## Requirements

- Integrate exact security head `093bf99dd2355c7eb84ad6d54318aa601c9ef809`
  onto the current PR #5 head.
- Preserve credential-free checkout and location-independent Make execution.
- Preserve Twilio timeout, error redaction, patched Go, vulnerability scanning,
  tests, workflow limits, and completed evidence contracts.
- Validate the combined tree from repository and external directories.

## Scope Boundaries

- Do not merge or close existing pull requests.
- Do not weaken tests, scanner enforcement, credential handling, or timeout
  coverage to resolve integration conflicts.
- Do not use live Twilio credentials or send a real SMS.

## Planned Verification

- Run formatting, vet, module verification, tests, race tests, coverage, build,
  vulnerability analysis, and all four Make gates.
- Run `make check` through the absolute Makefile path from an external directory.
- Run diff, artifact, conflict-marker, and changed-line credential audits.
- Record exact conflict resolution and completed verification evidence.
