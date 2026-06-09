# Account SID Validation Plan

status: completed

## Context

`level-up-webinar-101` validates phone numbers before dry runs and real Twilio
sends. Real sends also depend on `TWILIO_ACCOUNT_SID`, but the sample only
checked that the value was present, allowing placeholder or malformed values to
reach the Twilio client.

## Objectives

- Validate `TWILIO_ACCOUNT_SID` shape before real sends.
- Keep `DRY_RUN=1` usable without Twilio credentials.
- Return environment variable names for malformed Account SIDs without echoing
  configured values.
- Cover valid and invalid Account SID shapes in unit tests.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
