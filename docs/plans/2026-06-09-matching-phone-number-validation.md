# Matching Phone Number Validation

status: completed

## Context

The sample validates sender and recipient phone number shape, but a copied
`TO_PHONE_NUMBER` value could still be reused as `TWILIO_PHONE_NUMBER`. That
configuration is not useful for webinar demos and should fail before any Twilio
request is attempted.

## Completed Scope

- Rejected configurations where `TO_PHONE_NUMBER` and `TWILIO_PHONE_NUMBER`
  are the same valid phone-number value.
- Kept malformed matching values reported once per phone-number variable.
- Limited matching-number rejection to values that already pass phone shape
  validation so malformed matching values report each field once.
- Reported only environment variable names, not the shared phone number.
- Added unit coverage for the validation.
- Updated README, VISION, SECURITY, and CHANGES to document the guardrail.

## Verification

- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
