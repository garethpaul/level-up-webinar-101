# Twilio SMS Baseline Plan

status: completed

## Context

This repository is a compact Go webinar sample that sends an SMS with the
Twilio Go SDK. Before this pass it imported Twilio directly from `main.go`
without module metadata, tests, or validation of the environment variables the
sample depends on.

## Risks

- Running the sample with missing environment variables could send an invalid
  API request and produce a Twilio error that is harder for learners to
  understand.
- Demos needed a no-send path so setup can be checked without contacting Twilio
  or delivering a real SMS.
- Without `go.mod` and `go.sum`, dependency resolution was not reproducible.
- Without tests, credential and dry-run behavior could regress silently.

## Work Completed

- Added `go.mod` and `go.sum` with the Twilio Go SDK dependency.
- Added `loadConfig` validation for `TO_PHONE_NUMBER`, `TWILIO_PHONE_NUMBER`,
  `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`, optional `MESSAGE_BODY`, and
  `DRY_RUN`.
- Added a dry-run branch that validates required non-secret inputs and exits
  without sending SMS or printing phone numbers/secrets.
- Added unit tests for missing values, dry-run behavior, custom message body,
  trimming, truthy dry-run values, successful sends, and sender failures without
  contacting Twilio.
- Added `make check` to run formatting checks, module verification, unit tests,
  and a local build.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
