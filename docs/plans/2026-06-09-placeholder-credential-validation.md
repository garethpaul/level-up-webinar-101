# Placeholder Credential Validation

status: completed

## Context

Real sends validate Twilio Account SID and Auth Token shape before constructing
the Twilio client. All-zero values can still satisfy those length and hexadecimal
checks while representing common placeholder-shaped credentials.

## Objectives

- Reject all-zero Twilio Account SID bodies before real sends.
- Reject all-zero Twilio Auth Token values before real sends.
- Keep dry-run credential behavior unchanged.
- Keep validation errors limited to environment variable names rather than
  configured credential values.
- Document the guard beside the existing Twilio credential validation baseline.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
