# Redacted Twilio Send Errors

status: completed

## Context

The command currently wraps the raw Twilio sender error into its returned and
CLI-visible message. Provider errors can include phone numbers, request data,
or other operational details that should not be printed by a public sample.

## Priorities

1. Return a stable generic send-failure message to users.
2. Preserve the underlying error through Go error unwrapping for callers and
   tests that need programmatic diagnosis.
3. Prove that sensitive provider text is absent from the displayed error.
4. Keep dry-run, validation, timeout, and successful-send behavior unchanged.

## Implementation Units

### Redacted Error Type

File: `main.go`

Add a small error wrapper with a generic `Error()` result and an `Unwrap()`
method. Use it when the injected or Twilio sender returns an error.

### Focused Tests

File: `main_test.go`

Assert the generic message, absence of a sample phone number/provider detail,
and successful `errors.Is` matching of the underlying cause.

### Static Contract And Documentation

Files:

- `scripts/check-baseline.sh`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `CHANGES.md`
- `docs/plans/2026-06-12-redacted-twilio-send-errors.md`

Keep the redaction and unwrap contracts visible to the no-network baseline.

## Verification

- `gofmt -w main.go main_test.go`
- `go test ./...`
- `go vet ./...`
- `go mod verify`
- `make lint`
- `make test`
- `make build`
- `make check`
- hostile mutations exposing the cause or removing unwrapping
- `git diff --check`
- hosted push and pull-request checks

## Boundaries

- Do not make a real Twilio request during tests or CI.
- Do not log the raw provider error elsewhere.
- Do not remove access to the underlying cause through `errors.Is`/`errors.As`.
