# Phone Number Validation Plan

status: completed

## Context

`level-up-webinar-101` validates required Twilio environment variables before sending or dry-running an SMS. The phone number values were only checked for presence, so malformed values could pass local validation and fail later in Twilio.

## Objectives

- Validate recipient and sender phone numbers before dry runs and real sends.
- Keep the validation simple and readable for a webinar sample.
- Return environment variable names for malformed values without echoing configured phone numbers.
- Cover valid and invalid phone number shapes in unit tests.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
