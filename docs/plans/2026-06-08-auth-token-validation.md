# Auth Token Validation Plan

status: completed

## Context

Real sends already validate phone number shape and Twilio Account SID shape before contacting Twilio. `TWILIO_AUTH_TOKEN` was required, but malformed values were accepted until the Twilio client rejected the request.

## Objectives

- Validate real-send Twilio Auth Tokens as 32-character hexadecimal values.
- Keep `DRY_RUN=1` available without requiring valid Twilio credentials.
- Keep validation errors limited to environment variable names rather than configured credential values.
- Cover valid and malformed Auth Token values in Go unit tests.

## Verification

- `go test ./...`
- `make check`
- `git diff --check`
