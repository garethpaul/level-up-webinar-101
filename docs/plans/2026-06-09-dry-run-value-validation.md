# Dry-Run Value Validation Plan

status: completed

## Context

`DRY_RUN` is the safety switch for webinar setup checks that should validate
configuration without sending SMS. The previous parser treated any unrecognized
value as false, so a typo could become a real-send path when credentials were
present.

## Objectives

- Accept documented true and false dry-run values.
- Reject ambiguous `DRY_RUN` values before building the SMS config.
- Keep errors limited to the `DRY_RUN` environment variable name.
- Extend unit tests and docs for strict dry-run parsing.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
