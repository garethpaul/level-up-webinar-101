# Changes

## 2026-06-09

- Split the local Make verification into `make lint`, `make test`, and
  `make build` targets while keeping `make check` as the full gate.
- Added `MESSAGE_BODY` length validation that rejects oversized bodies by name
  without echoing configured content.
- Rejected ambiguous `DRY_RUN` values so typos fail closed instead of being
  treated as real sends.
- Rejected all-zero Twilio Account SID and Auth Token placeholder-shaped
  credentials without echoing configured values.
- Rejected matching sender and recipient phone numbers without echoing the
  shared value.
- Avoided duplicate phone-number field names when malformed sender and
  recipient values match.

## 2026-06-08

- Added Go module metadata and checksum locking for the Twilio SDK dependency.
- Added environment validation for recipient, sender, account SID, and auth token configuration.
- Added E.164-style phone number validation without echoing configured phone number values in errors.
- Added Twilio Account SID shape validation for real sends without echoing configured values in errors.
- Added Twilio Auth Token shape validation for real sends without echoing configured values in errors.
- Added `DRY_RUN` support for webinar demos that should validate configuration without sending SMS.
- Added an injectable run path so send success and Twilio failures can be tested without contacting Twilio.
- Added Go unit tests and a `make check` verification target that covers formatting, module verification, tests, and builds.
- Updated generated project context for the Go module, tests, and dry-run baseline.
