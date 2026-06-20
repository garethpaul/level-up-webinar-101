# Checkout Credential Boundary

status: completed

## Context

Canonical PR #3 runs a read-only, secret-free Go gate and performs no
authenticated Git operation after checkout, but the action default retained the
workflow token in the runner's Git configuration.

## Implementation

- Set `persist-credentials: false` on the one commit-pinned checkout step.
- Require exactly one checkout action and only the canonical workflow file.
- Preserve Ubuntu 24.04, Go 1.24, read-only permission, timeout, concurrency,
  and `make check` command.
- Preserve the no-live-Twilio, injected sender, dry-run, credential validation,
  UTF-8, message-size, and request-timeout contracts.

## Verification

- `make lint`, `make test`, `make build`, and `make check` passed.
- The checker passed from an external working directory.
- Workflow YAML parsing, shell syntax, Go formatting/vetting/module checks, and
  `git diff --check` passed.
- Focused hostile mutations rejected a missing or true credential setting,
  duplicate checkout action, extra workflow file, incomplete plan, and stale
  documentation; all hostile mutations rejected.
- Exact-head hosted verification remains pending until this successor is
  pushed.

## Boundaries

- Do not contact Twilio or introduce real credentials/phone numbers in tests.
- Do not weaken input validation, secret-safe diagnostics, or request timeout.
- Do not add post-checkout pushes, tags, or authenticated Git fetches.
